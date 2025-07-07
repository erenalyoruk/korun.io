package application_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"korun.io/auth-service/internal/application"
	"korun.io/auth-service/internal/config"
	sharedConfig "korun.io/shared/config"
	"korun.io/shared/events"
	"korun.io/shared/models"
)

type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) CreateAccount(ctx context.Context, account *models.Account) error {
	args := m.Called(ctx, account)
	return args.Error(0)
}

func (m *MockAuthRepository) GetAccountByEmail(ctx context.Context, email string) (*models.Account, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*models.Account), args.Error(1)
}

func (m *MockAuthRepository) GetAccountByID(ctx context.Context, id string) (*models.Account, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Account), args.Error(1)
}

func (m *MockAuthRepository) UpdateAccount(ctx context.Context, account *models.Account) error {
	args := m.Called(ctx, account)
	return args.Error(0)
}

func (m *MockAuthRepository) DeleteAccount(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockRefreshTokenRepository struct {
	mock.Mock
}

func (m *MockRefreshTokenRepository) CreateToken(ctx context.Context, token *models.RefreshToken) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) GetTokenByHash(ctx context.Context, tokenHash string) (*models.RefreshToken, error) {
	args := m.Called(ctx, tokenHash)
	return args.Get(0).(*models.RefreshToken), args.Error(1)
}

func (m *MockRefreshTokenRepository) RevokeToken(ctx context.Context, tokenHash string) error {
	args := m.Called(ctx, tokenHash)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) RevokeTokensForAccount(ctx context.Context, accountID string) error {
	args := m.Called(ctx, accountID)
	return args.Error(0)
}

type MockKafkaProducer struct {
	mock.Mock
}

func (m *MockKafkaProducer) PublishEvent(ctx context.Context, topic string, event *events.Event) error {
	args := m.Called(ctx, topic, event)
	return args.Error(0)
}

func (m *MockKafkaProducer) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestAuthService_Register(t *testing.T) {
	// setup mocks
	mockAccountRepo := new(MockAuthRepository)
	mockRefreshTokenRepo := new(MockRefreshTokenRepository)
	mockProducer := new(MockKafkaProducer)

	jwtConfig := &config.JWTConfig{
		Secret:          "supersecretjwtkey",
		AccessTokenTTL:  time.Minute * 15,
		RefreshTokenTTL: time.Hour * 24 * 7,
	}

	// setup kafka topics config
	infraConfig := &sharedConfig.InfrastructureConfig{
		Kafka: sharedConfig.KafkaConfig{
			AuthEvents: "auth_events_topic",
		},
	}

	// setup services
	tokenService := application.NewTokenService(mockRefreshTokenRepo, jwtConfig)
	authService := application.NewAuthService(mockAccountRepo, tokenService, mockProducer, infraConfig)

	ctx := context.Background()

	// successful registration
	t.Run("Successful Registration", func(t *testing.T) {
		req := &models.CreateAccountRequest{
			Email:    "test@example.com",
			Password: "StrongPassword1!",
		}

		mockAccountRepo.On("GetAccountByEmail", ctx, req.Email).Return(&models.Account{}, models.ErrAccountNotFound).Once()
		mockAccountRepo.On("CreateAccount", ctx, mock.AnythingOfType("*models.Account")).Return(nil).Once()
		mockRefreshTokenRepo.On("CreateToken", ctx, mock.AnythingOfType("*models.RefreshToken")).Return(nil).Once()
		mockProducer.On("PublishEvent", ctx, infraConfig.Kafka.AuthEvents, mock.Anything).Return(nil).Once()

		res, err := authService.Register(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.NotEmpty(t, res.AccessToken)
		assert.NotEmpty(t, res.RefreshToken)
		assert.Equal(t, req.Email, res.Account.Email)

		mockAccountRepo.AssertExpectations(t)
		mockRefreshTokenRepo.AssertExpectations(t)
		mockProducer.AssertExpectations(t)
	})

	// registration with existing email
	t.Run("Registration with existing email", func(t *testing.T) {
		req := &models.CreateAccountRequest{
			Email:    "existing@example.com",
			Password: "StrongPassword1!",
		}

		existingAccount := &models.Account{ID: "123", Email: req.Email}
		mockAccountRepo.On("GetAccountByEmail", ctx, req.Email).Return(existingAccount, nil).Once()

		res, err := authService.Register(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, models.ErrAccountExists, err)

		mockAccountRepo.AssertExpectations(t)
		mockRefreshTokenRepo.AssertNotCalled(t, "CreateToken")
		mockProducer.AssertNotCalled(t, "PublishEvent")
	})

	// registration with invalid email format
	t.Run("Registration with invalid email format", func(t *testing.T) {
		req := &models.CreateAccountRequest{
			Email:    "invalid-email",
			Password: "StrongPassword1!",
		}

		res, err := authService.Register(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, application.ErrInvalidEmail, err)

		mockAccountRepo.AssertNotCalled(t, "GetAccountByEmail")
		mockRefreshTokenRepo.AssertNotCalled(t, "CreateToken")
		mockProducer.AssertNotCalled(t, "PublishEvent")
	})

	// registration with weak password
	t.Run("Registration with weak password", func(t *testing.T) {
		req := &models.CreateAccountRequest{
			Email:    "test@example.com",
			Password: "weak",
		}

		res, err := authService.Register(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, application.ErrPasswordTooWeak, err)

		mockAccountRepo.AssertNotCalled(t, "GetAccountByEmail")
		mockRefreshTokenRepo.AssertNotCalled(t, "CreateToken")
		mockProducer.AssertNotCalled(t, "PublishEvent")
	})
}

func TestAuthService_Login(t *testing.T) {
	// setup mocks
	mockAccountRepo := new(MockAuthRepository)
	mockRefreshTokenRepo := new(MockRefreshTokenRepository)
	mockProducer := new(MockKafkaProducer)

	jwtConfig := &config.JWTConfig{
		Secret:          "supersecretjwtkey",
		AccessTokenTTL:  time.Minute * 15,
		RefreshTokenTTL: time.Hour * 24 * 7,
	}

	// setup kafka topics config
	infraConfig := &sharedConfig.InfrastructureConfig{
		Kafka: sharedConfig.KafkaConfig{
			AuthEvents: "auth_events_topic",
		},
	}

	tokenService := application.NewTokenService(mockRefreshTokenRepo, jwtConfig)
	authService := application.NewAuthService(mockAccountRepo, tokenService, mockProducer, infraConfig)

	ctx := context.Background()

	hashPassword := func(password string) string {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		return string(hashed)
	}

	// successful login
	t.Run("Successful Login", func(t *testing.T) {
		req := &models.LoginRequest{
			Email:    "test@example.com",
			Password: "StrongPassword1!",
		}
		clientIP := "127.0.0.1"
		userAgent := "test-agent"

		existingAccount := &models.Account{
			ID:           "account-123",
			Email:        req.Email,
			PasswordHash: hashPassword(req.Password),
		}

		mockAccountRepo.On("GetAccountByEmail", ctx, req.Email).Return(existingAccount, nil).Once()
		mockRefreshTokenRepo.On("CreateToken", ctx, mock.AnythingOfType("*models.RefreshToken")).Return(nil).Once()
		mockProducer.On("PublishEvent", ctx, infraConfig.Kafka.AuthEvents, mock.Anything).Return(nil).Once()

		res, err := authService.Login(ctx, req, clientIP, userAgent)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.NotEmpty(t, res.AccessToken)
		assert.NotEmpty(t, res.RefreshToken)
		assert.Equal(t, req.Email, res.Account.Email)

		mockAccountRepo.AssertExpectations(t)
		mockRefreshTokenRepo.AssertExpectations(t)
		mockProducer.AssertExpectations(t)
	})

	// login with invalid credentials (wrong password)
	t.Run("Login with invalid credentials (wrong password)", func(t *testing.T) {
		req := &models.LoginRequest{
			Email:    "test@example.com",
			Password: "WrongPassword!",
		}
		clientIP := "127.0.0.1"
		userAgent := "test-agent"

		existingAccount := &models.Account{
			ID:           "account-123",
			Email:        req.Email,
			PasswordHash: hashPassword("CorrectPassword!"),
		}

		mockAccountRepo.On("GetAccountByEmail", ctx, req.Email).Return(existingAccount, nil).Once()

		res, err := authService.Login(ctx, req, clientIP, userAgent)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, models.ErrInvalidCredentials, err)

		mockAccountRepo.AssertExpectations(t)
		mockRefreshTokenRepo.AssertNotCalled(t, "CreateToken")
		mockProducer.AssertNotCalled(t, "PublishEvent")
	})

	// login with non-existent email
	t.Run("Login with non-existent email", func(t *testing.T) {
		req := &models.LoginRequest{
			Email:    "nonexistent@example.com",
			Password: "AnyPassword!",
		}
		clientIP := "127.0.0.1"
		userAgent := "test-agent"

		mockAccountRepo.On("GetAccountByEmail", ctx, req.Email).Return(&models.Account{}, models.ErrAccountNotFound).Once()

		res, err := authService.Login(ctx, req, clientIP, userAgent)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, models.ErrInvalidCredentials, err)

		mockAccountRepo.AssertExpectations(t)
		mockRefreshTokenRepo.AssertNotCalled(t, "CreateToken")
		mockProducer.AssertNotCalled(t, "PublishEvent")
	})

	// login with invalid email format
	t.Run("Login with invalid email format", func(t *testing.T) {
		req := &models.LoginRequest{
			Email:    "invalid-email",
			Password: "AnyPassword!",
		}
		clientIP := "127.0.0.1"
		userAgent := "test-agent"

		res, err := authService.Login(ctx, req, clientIP, userAgent)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, application.ErrInvalidEmail, err)

		mockAccountRepo.AssertNotCalled(t, "GetAccountByEmail")
		mockRefreshTokenRepo.AssertNotCalled(t, "CreateToken")
		mockProducer.AssertNotCalled(t, "PublishEvent")
	})
}

func TestAuthService_RefreshToken(t *testing.T) {
	mockAccountRepo := new(MockAuthRepository)
	mockRefreshTokenRepo := new(MockRefreshTokenRepository)
	mockProducer := new(MockKafkaProducer)

	jwtConfig := &config.JWTConfig{
		Secret:          "supersecretjwtkey",
		AccessTokenTTL:  time.Minute * 15,
		RefreshTokenTTL: time.Hour * 24 * 7,
	}

	// setup kafka topics config
	infraConfig := &sharedConfig.InfrastructureConfig{
		Kafka: sharedConfig.KafkaConfig{
			AuthEvents: "auth_events_topic",
		},
	}

	// setup services
	tokenService := application.NewTokenService(mockRefreshTokenRepo, jwtConfig)
	authService := application.NewAuthService(mockAccountRepo, tokenService, mockProducer, infraConfig)

	ctx := context.Background()

	// successful token refresh
	t.Run("Successful Token Refresh", func(t *testing.T) {
		accountID := "test-account-id"
		refreshTokenString := "valid-refresh-token"
		clientIP := "127.0.0.1"
		userAgent := "test-agent"

		oldRefreshToken := &models.RefreshToken{
			ID:        "old-token-id",
			AccountID: accountID,
			TokenHash: tokenService.HashToken(refreshTokenString),
			ExpiresAt: time.Now().Add(time.Hour),
			RevokedAt: nil,
		}
		existingAccount := &models.Account{ID: accountID, Email: "test@example.com"}

		mockRefreshTokenRepo.On("GetTokenByHash", ctx, tokenService.HashToken(refreshTokenString)).Return(oldRefreshToken, nil).Once()
		mockRefreshTokenRepo.On("RevokeToken", ctx, oldRefreshToken.TokenHash).Return(nil).Once()
		mockAccountRepo.On("GetAccountByID", ctx, accountID).Return(existingAccount, nil).Once()
		mockRefreshTokenRepo.On("CreateToken", ctx, mock.AnythingOfType("*models.RefreshToken")).Return(nil).Once()

		req := &models.RefreshRequest{RefreshToken: refreshTokenString}
		res, err := authService.RefreshToken(ctx, req, clientIP, userAgent)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.NotEmpty(t, res.AccessToken)
		assert.NotEmpty(t, res.RefreshToken)
		assert.Equal(t, existingAccount.Email, res.Account.Email)

		mockRefreshTokenRepo.AssertExpectations(t)
		mockAccountRepo.AssertExpectations(t)
	})

	// refresh with invalid refresh token
	t.Run("Refresh with invalid refresh token", func(t *testing.T) {
		refreshTokenString := "invalid-refresh-token"
		clientIP := "127.0.0.1"
		userAgent := "test-agent"

		mockRefreshTokenRepo.On("GetTokenByHash", ctx, tokenService.HashToken(refreshTokenString)).Return(&models.RefreshToken{}, models.ErrRefreshTokenNotFound).Once()

		req := &models.RefreshRequest{RefreshToken: refreshTokenString}
		res, err := authService.RefreshToken(ctx, req, clientIP, userAgent)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, models.ErrInvalidRefreshToken, err)

		mockRefreshTokenRepo.AssertExpectations(t)
		mockAccountRepo.AssertNotCalled(t, "GetAccountByID")
	})

	// refresh with expired refresh token
	t.Run("Refresh with expired refresh token", func(t *testing.T) {
		accountID := "test-account-id"
		refreshTokenString := "expired-refresh-token"
		clientIP := "127.0.0.1"
		userAgent := "test-agent"

		oldRefreshToken := &models.RefreshToken{
			ID:        "expired-token-id",
			AccountID: accountID,
			TokenHash: tokenService.HashToken(refreshTokenString),
			ExpiresAt: time.Now().Add(-time.Hour), // Expired
			RevokedAt: nil,
		}

		mockRefreshTokenRepo.On("GetTokenByHash", ctx, tokenService.HashToken(refreshTokenString)).Return(oldRefreshToken, nil).Once()

		req := &models.RefreshRequest{RefreshToken: refreshTokenString}
		res, err := authService.RefreshToken(ctx, req, clientIP, userAgent)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, models.ErrInvalidRefreshToken, err)

		mockRefreshTokenRepo.AssertExpectations(t)
		mockAccountRepo.AssertNotCalled(t, "GetAccountByID")
	})

	// refresh with revoked refresh token
	t.Run("Refresh with revoked refresh token", func(t *testing.T) {
		accountID := "test-account-id"
		refreshTokenString := "revoked-refresh-token"
		clientIP := "127.0.0.1"
		userAgent := "test-agent"

		oldRefreshToken := &models.RefreshToken{
			ID:        "revoked-token-id",
			AccountID: accountID,
			TokenHash: tokenService.HashToken(refreshTokenString),
			ExpiresAt: time.Now().Add(time.Hour),
			RevokedAt: &[]time.Time{time.Now()}[0], // Revoked
		}

		mockRefreshTokenRepo.On("GetTokenByHash", ctx, tokenService.HashToken(refreshTokenString)).Return(oldRefreshToken, nil).Once()
		mockRefreshTokenRepo.On("RevokeTokensForAccount", ctx, accountID).Return(nil).Once() // Expect this to be called

		req := &models.RefreshRequest{RefreshToken: refreshTokenString}
		res, err := authService.RefreshToken(ctx, req, clientIP, userAgent)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, models.ErrInvalidRefreshToken, err)

		mockRefreshTokenRepo.AssertExpectations(t)
		mockAccountRepo.AssertNotCalled(t, "GetAccountByID")
	})
}

func TestAuthService_Logout(t *testing.T) {
	mockAccountRepo := new(MockAuthRepository)
	mockRefreshTokenRepo := new(MockRefreshTokenRepository)
	mockProducer := new(MockKafkaProducer)

	jwtConfig := &config.JWTConfig{
		Secret:          "supersecretjwtkey",
		AccessTokenTTL:  time.Minute * 15,
		RefreshTokenTTL: time.Hour * 24 * 7,
	}

	infraConfig := &sharedConfig.InfrastructureConfig{
		Kafka: sharedConfig.KafkaConfig{
			AuthEvents: "auth_events_topic",
		},
	}

	tokenService := application.NewTokenService(mockRefreshTokenRepo, jwtConfig)
	authService := application.NewAuthService(mockAccountRepo, tokenService, mockProducer, infraConfig)

	ctx := context.Background()

	// successful logout
	t.Run("Successful Logout", func(t *testing.T) {
		accountID := "test-account-id"

		mockRefreshTokenRepo.On("RevokeTokensForAccount", ctx, accountID).Return(nil).Once()

		err := authService.Logout(ctx, accountID)

		assert.NoError(t, err)

		mockRefreshTokenRepo.AssertExpectations(t)
	})

	// logout with repository error
	t.Run("Logout with repository error", func(t *testing.T) {
		accountID := "test-account-id"
		repoError := errors.New("database error")

		mockRefreshTokenRepo.On("RevokeTokensForAccount", ctx, accountID).Return(repoError).Once()

		err := authService.Logout(ctx, accountID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to revoke tokens for account")

		mockRefreshTokenRepo.AssertExpectations(t)
	})
}
