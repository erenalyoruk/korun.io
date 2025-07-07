package repository

import (
	"context"

	"korun.io/shared/models"
)

type AuthRepository interface {
	CreateAccount(ctx context.Context, account *models.Account) error
	GetAccountByEmail(ctx context.Context, email string) (*models.Account, error)
	GetAccountByID(ctx context.Context, id string) (*models.Account, error)
	UpdateAccount(ctx context.Context, account *models.Account) error
	DeleteAccount(ctx context.Context, id string) error
}

type RefreshTokenRepository interface {
	CreateToken(ctx context.Context, token *models.RefreshToken) error
	GetTokenByHash(ctx context.Context, tokenHash string) (*models.RefreshToken, error)
	RevokeToken(ctx context.Context, tokenID string) error
	RevokeTokensForAccount(ctx context.Context, accountID string) error
}
