package application

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"korun.io/auth-service/internal/config"
	"korun.io/auth-service/internal/domain"
	"korun.io/shared/models"
)

type TokenService struct {
	jwtConfig        *config.JWTConfig
	refreshTokenRepo domain.RefreshTokenRepository
}

type TokenClaims struct {
	AccountID string `json:"account_id"`
	Email     string `json:"email"`
	jwt.RegisteredClaims
}

func NewTokenService(refreshTokenRepo domain.RefreshTokenRepository, jwtConfig *config.JWTConfig) *TokenService {
	return &TokenService{
		refreshTokenRepo: refreshTokenRepo,
		jwtConfig:        jwtConfig,
	}
}

func (s *TokenService) GenerateAuthTokens(ctx context.Context, account *models.Account, clientIP, userAgent string) (string, string, error) {
	accessToken, err := s.generateAccessToken(account)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.generateAndStoreRefreshToken(ctx, account, clientIP, userAgent)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate and store refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

func (s *TokenService) generateAccessToken(account *models.Account) (string, error) {
	claims := &TokenClaims{
		AccountID: account.ID,
		Email:     account.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.jwtConfig.AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   account.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtConfig.Secret))
}

func (s *TokenService) generateAndStoreRefreshToken(ctx context.Context, account *models.Account, clientIP, userAgent string) (string, error) {
	refreshToken, err := models.NewRefreshToken(account.ID, s.jwtConfig.RefreshTokenTTL, clientIP, userAgent)
	if err != nil {
		return "", fmt.Errorf("failed to create new refresh token: %w", err)
	}

	tokenHash := s.HashToken(refreshToken.Token)
	refreshToken.TokenHash = tokenHash

	if err := s.refreshTokenRepo.CreateToken(ctx, refreshToken); err != nil {
		return "", fmt.Errorf("failed to store refresh token: %w", err)
	}

	return refreshToken.Token, nil
}

func (s *TokenService) RevokeTokensForAccount(ctx context.Context, accountID string) error {
	return s.refreshTokenRepo.RevokeTokensForAccount(ctx, accountID)
}

func (s *TokenService) ValidateAccessToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtConfig.Secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (s *TokenService) HashToken(token string) string {
	hasher := sha256.New()
	hasher.Write([]byte(token))
	return hex.EncodeToString(hasher.Sum(nil))
}
