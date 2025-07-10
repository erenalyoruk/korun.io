package domain

import (
	"context"

	"korun.io/shared/models"
)

type RefreshTokenRepository interface {
	CreateToken(ctx context.Context, token *models.RefreshToken) error
	GetTokenByHash(ctx context.Context, hash string) (*models.RefreshToken, error)
	RevokeToken(ctx context.Context, token *models.RefreshToken) error
	RevokeTokensForAccount(ctx context.Context, accountID string) error
}
