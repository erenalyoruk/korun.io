package models

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        string     `json:"id" db:"id"`
	AccountID string     `json:"account_id" db:"account_id"`
	Token     string     `json:"-"`
	TokenHash string     `json:"-" db:"token_hash"`
	ExpiresAt time.Time  `json:"expires_at" db:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at" db:"revoked_at"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	ClientIP  *string    `json:"client_ip_address" db:"client_ip_address"`
	UserAgent *string    `json:"user_agent" db:"user_agent"`
}

func NewRefreshToken(accountID string, ttl time.Duration, clientIP, userAgent string) (*RefreshToken, error) {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return nil, err
	}

	return &RefreshToken{
		ID:        uuid.NewString(),
		AccountID: accountID,
		Token:     base64.StdEncoding.EncodeToString(tokenBytes),
		ExpiresAt: time.Now().UTC().Add(ttl),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		ClientIP:  &clientIP,
		UserAgent: &userAgent,
	}, nil
}

func (rt *RefreshToken) IsExpired() bool {
	return time.Now().UTC().After(rt.ExpiresAt)
}

func (rt *RefreshToken) IsRevoked() bool {
	return rt.RevokedAt != nil
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LogoutRequest struct {
	AccessToken string `json:"access_token" binding:"required"`
}

var (
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
	ErrTokenAlreadyRevoked  = errors.New("token already revoked")
	ErrInvalidRefreshToken  = errors.New("invalid refresh token")
)
