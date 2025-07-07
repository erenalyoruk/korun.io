package models

import (
	"errors"
	"regexp"
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID           string    `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	IsVerified   bool      `json:"is_verified" db:"is_verified"`
}

type CreateAccountRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

var (
	ErrAccountNotFound    = errors.New("account not found")
	ErrAccountExists      = errors.New("account already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

func (a *Account) Validate() error {
	if a.Email == "" {
		return errors.New("email is required")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(a.Email) {
		return errors.New("invalid email format")
	}

	if len(a.PasswordHash) < 8 || len(a.PasswordHash) > 64 {
		return errors.New("password must be between 8 and 64 characters")
	}

	return nil
}

func NewAccount(email string) *Account {
	return &Account{
		ID:         uuid.New().String(),
		Email:      email,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		IsVerified: false,
	}
}
