package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"korun.io/secret-service/internal/interfaces"
	"korun.io/shared/models"
)

type PostgresSecretRepository struct {
	db *pgxpool.Pool
}

func NewSecretRepository(db *pgxpool.Pool) interfaces.SecretRepository {
	return &PostgresSecretRepository{db: db}
}

func (r *PostgresSecretRepository) Create(ctx context.Context, secret *models.Secret) error {
	err := r.db.QueryRow(ctx, `
    INSERT INTO secrets (name, value, description)
    FROM VALUES ($1, $2, $3)
    RETURNING id
  `, secret.Name, secret.Value, secret.Description).Scan(&secret.ID)

	return err
}

func (r *PostgresSecretRepository) GetSecrets(ctx context.Context) ([]models.Secret, error) {
	rows, err := r.db.Query(ctx, `
    SELECT id, name, value, description
    FROM secrets
  `)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var secrets []models.Secret
	for rows.Next() {
		var s models.Secret
		if err := rows.Scan(&s.ID, &s.Name, &s.Value, &s.Description); err != nil {
			return nil, err
		}

		secrets = append(secrets, s)
	}

	return secrets, nil
}

func (r *PostgresSecretRepository) GetSecretByID(ctx context.Context, id uint) (*models.Secret, error) {
	var s models.Secret
	err := r.db.QueryRow(ctx, `
    SELECT id, name, value, description
    FROM secrets
    WHERE id = $1
  `, id).Scan(&s.ID, &s.Name, &s.Value, &s.Description)

	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (r *PostgresSecretRepository) GetSecretByName(ctx context.Context, name string) (*models.Secret, error) {
	var s models.Secret
	err := r.db.QueryRow(ctx, `
    SELECT id, name, value, description
    FROM secrets
    WHERE name = $1
  `, name).Scan(&s.ID, &s.Name, &s.Value, &s.Description)

	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (r *PostgresSecretRepository) Update(ctx context.Context, secret *models.Secret) error {
	cmd, err := r.db.Exec(ctx, `
    UPDATE secrets
    SET name = $1, value = $2, description = $3
    WHERE id = $4
  `, secret.Name, secret.Value, secret.Description, secret.ID)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *PostgresSecretRepository) Delete(ctx context.Context, id uint) error {
	cmd, err := r.db.Exec(ctx, `
    DELETE FROM secrets
    WHERE id = $1
  `, id)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
