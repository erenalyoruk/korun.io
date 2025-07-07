package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"
	"korun.io/auth-service/internal/config"
	"korun.io/auth-service/internal/repository"
	"korun.io/auth-service/internal/validators"
	"korun.io/shared/events"
	"korun.io/shared/messaging"
	"korun.io/shared/models"
)

var (
	ErrInvalidEmail    = errors.New("invalid email format")
	ErrPasswordTooWeak = errors.New("password is too weak")
)

type AuthService struct {
	accountRepo  repository.AuthRepository
	tokenService *TokenService
	producer     *messaging.KafkaProducer
	topicConfig  *config.KafkaTopicsConfig
}

func NewAuthService(
	accountRepo repository.AuthRepository,
	tokenService *TokenService,
	producer *messaging.KafkaProducer,
	topicConfig *config.KafkaTopicsConfig,
) *AuthService {
	return &AuthService{
		accountRepo:  accountRepo,
		tokenService: tokenService,
		producer:     producer,
		topicConfig:  topicConfig,
	}
}

func (s *AuthService) Register(ctx context.Context, req *models.CreateAccountRequest) (*models.AuthResponse, error) {
	if !validators.IsValidEmail(req.Email) {
		return nil, ErrInvalidEmail
	}

	if !validators.IsStrongPassword(req.Password) {
		return nil, ErrPasswordTooWeak
	}

	_, err := s.accountRepo.GetAccountByEmail(ctx, req.Email)
	if err == nil {
		return nil, models.ErrAccountExists
	}
	if !errors.Is(err, models.ErrAccountNotFound) {
		return nil, fmt.Errorf("failed to check account existence: %w", err)
	}

	account := models.NewAccount(req.Email)

	hashedPassword, err := s.hashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	account.PasswordHash = hashedPassword

	if err := account.Validate(); err != nil {
		return nil, fmt.Errorf("account validation failed: %w", err)
	}

	if err := s.accountRepo.CreateAccount(ctx, account); err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	accessToken, refreshToken, err := s.tokenService.GenerateAuthTokens(ctx, account, "", "")
	if err != nil {
		return nil, fmt.Errorf("failed to generate auth tokens: %w", err)
	}

	s.publishAccountRegisteredEvent(ctx, account)

	return &models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Account:      account,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req *models.LoginRequest, clientIP, userAgent string) (*models.AuthResponse, error) {
	if !validators.IsValidEmail(req.Email) {
		return nil, ErrInvalidEmail
	}

	account, err := s.accountRepo.GetAccountByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, models.ErrAccountNotFound) {
			return nil, models.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	if !s.checkPassword(account.PasswordHash, req.Password) {
		return nil, models.ErrInvalidCredentials
	}

	accessToken, refreshToken, err := s.tokenService.GenerateAuthTokens(ctx, account, clientIP, userAgent)
	if err != nil {
		return nil, fmt.Errorf("failed to generate auth tokens: %w", err)
	}

	s.publishAccountLoggedInEvent(ctx, account)

	return &models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Account:      account,
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, req *models.RefreshRequest, clientIP, userAgent string) (*models.AuthResponse, error) {
	oldTokenHash := s.tokenService.HashToken(req.RefreshToken)

	oldToken, err := s.tokenService.refreshTokenRepo.GetTokenByHash(ctx, oldTokenHash)
	if err != nil {
		if errors.Is(err, models.ErrRefreshTokenNotFound) {
			return nil, models.ErrInvalidRefreshToken
		}

		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}

	if oldToken.IsExpired() {
		return nil, models.ErrInvalidRefreshToken
	}

	if oldToken.IsRevoked() {
		// Token has been used before, which could indicate a token theft attempt.
		// As a security measure, we revoke all tokens for the user.
		slog.Warn("Attempted to use a revoked refresh token", "account_id", oldToken.AccountID)
		if err := s.tokenService.RevokeTokensForAccount(ctx, oldToken.AccountID); err != nil {
			slog.Error("Failed to revoke tokens for account after detecting token reuse", "error", err, "account_id", oldToken.AccountID)
		}
		return nil, models.ErrInvalidRefreshToken
	}

	if err := s.tokenService.refreshTokenRepo.RevokeToken(ctx, oldToken.ID); err != nil {
		slog.Error("Failed to revoke old refresh token during refresh flow", "error", err, "token_id", oldToken.ID)
		return nil, fmt.Errorf("failed to revoke old refresh token: %w", err)
	}

	account, err := s.accountRepo.GetAccountByID(ctx, oldToken.AccountID)
	if err != nil {
		slog.Error("Account not found for valid refresh token", "error", err, "account_id", oldToken.AccountID)
		return nil, errors.New("associated account not found")
	}

	accessToken, refreshToken, err := s.tokenService.GenerateAuthTokens(ctx, account, clientIP, userAgent)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new auth tokens: %w", err)
	}

	slog.Info("Refresh token successfully used", "account_id", account.ID, "client_ip_address", clientIP, "user_agent", userAgent)

	return &models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Account:      account,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, req *models.LogoutRequest) error {
	claims, err := s.tokenService.ValidateAccessToken(req.AccessToken)
	if err != nil {
		return fmt.Errorf("failed to validate access token: %w", err)
	}

	if err := s.tokenService.RevokeTokensForAccount(ctx, claims.AccountID); err != nil {
		return fmt.Errorf("failed to revoke tokens for account: %w", err)
	}

	slog.Info("User logged out successfully", "account_id", claims.AccountID)
	return nil
}

func (s *AuthService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *AuthService) checkPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func (s *AuthService) publishAccountRegisteredEvent(ctx context.Context, account *models.Account) {
	eventData := map[string]any{
		"account_id": account.ID,
		"email":      account.Email,
	}

	event := events.NewEvent(events.AccountRegisteredEvent, "auth-service", eventData)

	if err := s.producer.PublishEvent(ctx, s.topicConfig.AuthEvents, event); err != nil {
		slog.Error("Failed to publish account registered event", "error", err)
	}
}

func (s *AuthService) publishAccountLoggedInEvent(ctx context.Context, account *models.Account) {
	eventData := map[string]any{
		"account_id": account.ID,
		"email":      account.Email,
		"login_time": time.Now().UTC(),
	}

	event := events.NewEvent(events.AccountLoggedInEvent, "auth-service", eventData)

	if err := s.producer.PublishEvent(ctx, s.topicConfig.AuthEvents, event); err != nil {
		slog.Error("Failed to publish account logged in event", "error", err)
	}
}
