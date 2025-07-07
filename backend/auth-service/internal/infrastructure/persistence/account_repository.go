package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
	"korun.io/auth-service/internal/domain"
	"korun.io/shared/models"
)

type PostgresAuthRepository struct {
	db *sqlx.DB
}

func NewPostgresAuthRepository(db *sqlx.DB) domain.AuthRepository {
	return &PostgresAuthRepository{db: db}
}

func (r *PostgresAuthRepository) CreateAccount(ctx context.Context, account *models.Account) error {
	query := `
		INSERT INTO accounts (id, email, password_hash, created_at, updated_at, is_verified)
		VALUES (:id, :email, :password_hash, :created_at, :updated_at, :is_verified)`

	_, err := r.db.NamedExecContext(ctx, query, account)
	if err != nil {
		slog.Error("Failed to create account", "error", err, "email", account.Email)
		return fmt.Errorf("failed to create account: %w", err)
	}

	return nil
}

func (r *PostgresAuthRepository) GetAccountByEmail(ctx context.Context, email string) (*models.Account, error) {
	query := `
    SELECT id, email, password_hash, created_at, updated_at
    FROM accounts
    WHERE email = $1
  `

	account := models.Account{}
	err := r.db.GetContext(ctx, &account, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrAccountNotFound
		}
		slog.Error("Failed to get account by email", "error", err, "email", email)
		return nil, fmt.Errorf("failed to get account by email: %w", err)
	}

	return &account, nil
}

func (r *PostgresAuthRepository) GetAccountByID(ctx context.Context, id string) (*models.Account, error) {
	query := `
    SELECT id, email, password_hash, created_at, updated_at
    FROM accounts
    WHERE id = $1
  `

	account := models.Account{}
	err := r.db.GetContext(ctx, &account, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrAccountNotFound
		}
		slog.Error("Failed to get account by ID", "error", err, "id", id)
		return nil, fmt.Errorf("failed to get account by ID: %w", err)
	}

	return &account, nil
}

func (r *PostgresAuthRepository) UpdateAccount(ctx context.Context, account *models.Account) error {
	query := `
    UPDATE accounts
    SET email = $1, password_hash = $2, updated_at = $3
    WHERE id = $4
  `

	_, err := r.db.ExecContext(ctx, query, account.Email, account.PasswordHash, account.UpdatedAt, account.ID)
	if err != nil {
		slog.Error("Failed to update account", "error", err, "account_id", account.ID)
		return fmt.Errorf("failed to update account: %w", err)
	}

	slog.Info("Account updated successfully", "account_id", account.ID, "email", account.Email)
	return nil
}

func (r *PostgresAuthRepository) DeleteAccount(ctx context.Context, id string) error {
	query := `
    DELETE FROM accounts
    WHERE id = $1
  `

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		slog.Error("Failed to delete account", "error", err, "account_id", id)
		return fmt.Errorf("failed to delete account: %w", err)
	}

	slog.Info("Account deleted successfully", "account_id", id)
	return nil
}
