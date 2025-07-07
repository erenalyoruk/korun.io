package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
	"korun.io/shared/models"
)

type PostgresRefreshTokenRepository struct {
	db *sqlx.DB
}

func NewPostgresRefreshTokenRepository(db *sqlx.DB) RefreshTokenRepository {
	return &PostgresRefreshTokenRepository{db: db}
}

func (r *PostgresRefreshTokenRepository) CreateToken(ctx context.Context, token *models.RefreshToken) error {
	query := `
    INSERT INTO refresh_tokens (id, account_id, token_hash, expires_at, revoked_at, created_at, updated_at, client_ip_address, user_agent)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
  `

	if token.ClientIP != nil && *token.ClientIP == "" {
		token.ClientIP = nil
	}

	_, err := r.db.ExecContext(ctx, query,
		token.ID,
		token.AccountID,
		token.TokenHash,
		token.ExpiresAt,
		token.RevokedAt,
		token.CreatedAt,
		token.UpdatedAt,
		token.ClientIP,
		token.UserAgent,
	)
	if err != nil {
		slog.Error("Failed to create refresh token", "error", err, "account_id", token.AccountID)
		return fmt.Errorf("failed to create refresh token: %w", err)
	}

	slog.Debug("Refresh token created successfully", "token_id", token.ID, "account_id", token.AccountID)
	return nil
}

func (r *PostgresRefreshTokenRepository) GetTokenByHash(ctx context.Context, tokenHash string) (*models.RefreshToken, error) {
	query := `
    SELECT id, account_id, token_hash, expires_at, revoked_at, created_at, updated_at, client_ip_address, user_agent
    FROM refresh_tokens
    WHERE token_hash = $1
  `

	token := models.RefreshToken{}
	err := r.db.GetContext(ctx, &token, query, tokenHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrRefreshTokenNotFound
		}

		slog.Error("Failed to get refresh token by hash", "error", err)
		return nil, fmt.Errorf("failed to get refresh token by hash: %w", err)
	}

	return &token, nil
}

func (r *PostgresRefreshTokenRepository) RevokeToken(ctx context.Context, tokenID string) error {
	now := time.Now().UTC()
	query := `
    UPDATE refresh_tokens SET revoked_at = $1, updated_at = $2
    WHERE id = $3 AND revoked_at IS NULL
  `

	res, err := r.db.ExecContext(ctx, query, now, now, tokenID)
	if err != nil {
		slog.Error("Failed to revoke refresh token", "error", err, "token_id", tokenID)
		return fmt.Errorf("failed to revoke refresh token: %w", err)
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return models.ErrTokenAlreadyRevoked
	}

	slog.Info("Refresh token revoked successfully", "token_id", tokenID, "revoked_at", now)
	return nil
}

func (r *PostgresRefreshTokenRepository) RevokeTokensForAccount(ctx context.Context, accountID string) error {
	now := time.Now().UTC()
	query := `
    UPDATE refresh_tokens SET revoked_at = $1, updated_at = $2
    WHERE account_id = $3 AND revoked_at IS NULL
  `

	_, err := r.db.ExecContext(ctx, query, now, now, accountID)
	if err != nil {
		slog.Error("Failed to revoke refresh tokens for account", "error", err, "account_id", accountID)
		return fmt.Errorf("failed to revoke refresh tokens for account: %w", err)
	}

	slog.Info("Refresh tokens revoked successfully for account", "account_id", accountID, "revoked_at", now)
	return nil
}
