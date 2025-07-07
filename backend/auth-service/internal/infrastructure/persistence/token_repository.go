package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"korun.io/auth-service/internal/infrastructure/redis"
	"korun.io/auth-service/internal/domain"
	"korun.io/shared/models"
)

type RedisRefreshTokenRepository struct {
	redisClient *redis.Client
}

func NewRedisRefreshTokenRepository(redisClient *redis.Client) domain.RefreshTokenRepository {
	return &RedisRefreshTokenRepository{
		redisClient: redisClient,
	}
}

func (r *RedisRefreshTokenRepository) CreateToken(ctx context.Context, token *models.RefreshToken) error {
	tokenKey := fmt.Sprintf("refresh_token:%s", token.TokenHash)
	tokenBytes, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("failed to marshal refresh token: %w", err)
	}

	if err := r.redisClient.Set(ctx, tokenKey, tokenBytes, time.Until(token.ExpiresAt)); err != nil {
		slog.Error("Failed to store refresh token in Redis", "error", err, "token_id", token.ID)
		return fmt.Errorf("failed to store refresh token in Redis: %w", err)
	}

	accountTokensKey := fmt.Sprintf("account_tokens:%s", token.AccountID)
	if err := r.redisClient.SAdd(ctx, accountTokensKey, token.TokenHash); err != nil {
		slog.Error("Failed to add token hash to account tokens set in Redis", "error", err, "account_id", token.AccountID)
		return fmt.Errorf("failed to add token hash to account tokens set in Redis: %w", err)
	}

	slog.Debug("Refresh token created successfully", "token_id", token.ID, "account_id", token.AccountID)
	return nil
}

func (r *RedisRefreshTokenRepository) GetTokenByHash(ctx context.Context, tokenHash string) (*models.RefreshToken, error) {
	tokenKey := fmt.Sprintf("refresh_token:%s", tokenHash)
	val, err := r.redisClient.Get(ctx, tokenKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get refresh token from Redis: %w", err)
	}

	if val == "" {
		return nil, models.ErrRefreshTokenNotFound
	}

	var token models.RefreshToken
	if err := json.Unmarshal([]byte(val), &token); err != nil {
		return nil, fmt.Errorf("failed to unmarshal refresh token from Redis: %w", err)
	}

	return &token, nil
}

func (r *RedisRefreshTokenRepository) RevokeToken(ctx context.Context, tokenHash string) error {
	tokenKey := fmt.Sprintf("refresh_token:%s", tokenHash)

	var token models.RefreshToken
	val, err := r.redisClient.Get(ctx, tokenKey)
	if err != nil {
		slog.Error("Failed to get refresh token from Redis for revocation (might be already gone)", "error", err, "token_hash", tokenHash)
	} else if val != "" {
		if err := json.Unmarshal([]byte(val), &token); err != nil {
			slog.Error("Failed to unmarshal refresh token from Redis for revocation", "error", err, "token_hash", tokenHash)
		}
	}

	if err := r.redisClient.Del(ctx, tokenKey); err != nil {
		slog.Error("Failed to delete refresh token from Redis", "error", err, "token_hash", tokenHash)
		return fmt.Errorf("failed to delete refresh token from Redis: %w", err)
	}

	if token.AccountID != "" {
		accountTokensKey := fmt.Sprintf("account_tokens:%s", token.AccountID)
		if err := r.redisClient.SRem(ctx, accountTokensKey, tokenHash); err != nil {
			slog.Error("Failed to remove token hash from account tokens set in Redis", "error", err, "token_hash", tokenHash, "account_id", token.AccountID)
		}
	}

	slog.Info("Refresh token revoked successfully (or was already revoked)", "token_hash", tokenHash)
	return nil
}

func (r *RedisRefreshTokenRepository) RevokeTokensForAccount(ctx context.Context, accountID string) error {
	accountTokensKey := fmt.Sprintf("account_tokens:%s", accountID)
	tokenHashes, err := r.redisClient.SMembers(ctx, accountTokensKey)
	if err != nil {
		slog.Error("Failed to get token hashes from account tokens set in Redis", "error", err, "account_id", accountID)
		return fmt.Errorf("failed to get token hashes from account tokens set in Redis: %w", err)
	}

	if len(tokenHashes) > 0 {
		keysToDelete := make([]string, len(tokenHashes))
		for i, hash := range tokenHashes {
			keysToDelete[i] = fmt.Sprintf("refresh_token:%s", hash)
		}

		if err := r.redisClient.Del(ctx, keysToDelete...); err != nil {
			slog.Error("Failed to delete individual refresh tokens from Redis", "error", err, "account_id", accountID)
			return fmt.Errorf("failed to delete individual refresh tokens from Redis: %w", err)
		}
	}

	if err := r.redisClient.Del(ctx, accountTokensKey); err != nil {
		slog.Error("Failed to delete account tokens set from Redis", "error", err, "account_id", accountID)
		return fmt.Errorf("failed to delete account tokens set from Redis: %w", err)
	}

	slog.Info("Refresh tokens revoked successfully for account", "account_id", accountID)
	return nil
}
