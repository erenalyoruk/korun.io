package models

import (
	"crypto/rand"
	"encoding/base64"
	"net"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type RefreshToken struct {
	ID        string      `json:"id" db:"id"`
	AccountID string      `json:"account_id" db:"account_id"`
	Token     string      `json:"-" db:"token"`
	TokenHash string      `json:"token_hash" db:"token_hash"`
	CreatedAt time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt time.Time   `json:"updated_at" db:"updated_at"`
	ExpiresAt time.Time   `json:"expires_at" db:"expires_at"`
	RevokedAt *time.Time  `json:"revoked_at" db:"revoked_at"`
	ClientIP  *net.IPAddr `json:"client_ip_address" db:"client_ip_address"`
	UserAgent *string     `json:"user_agent" db:"user_agent"`
}

func NewRefreshToken(accountID string, ttl time.Duration, clientIP *net.IPAddr, userAgent *string) (*RefreshToken, error) {
	tokenBytes := make([]byte, 32) // 256 bits
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return nil, err
	}

	token := base64.StdEncoding.EncodeToString(tokenBytes)

	tokenHash, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &RefreshToken{
		ID:        uuid.NewString(),
		AccountID: accountID,
		Token:     token,
		TokenHash: string(tokenHash),
		ExpiresAt: time.Now().Add(ttl),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ClientIP:  clientIP,
		UserAgent: userAgent,
	}, nil
}

func (rt *RefreshToken) CompareToken(token string) error {
	return bcrypt.CompareHashAndPassword([]byte(rt.TokenHash), []byte(token))
}

func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}

func (rt *RefreshToken) IsRevoked() bool {
	return rt.RevokedAt != nil && rt.RevokedAt.Before(time.Now())
}

type TokenClaims struct {
	AccountID string `json:"account_id"`
	Email     string `json:"email"`
	jwt.RegisteredClaims
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
