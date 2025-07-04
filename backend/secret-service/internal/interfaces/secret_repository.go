package interfaces

import (
	"context"

	"korun.io/shared/models"
)

type SecretRepository interface {
	Create(ctx context.Context, secret *models.Secret) error

	GetSecrets(ctx context.Context) ([]models.Secret, error)
	GetSecretByName(ctx context.Context, name string) (*models.Secret, error)
	GetSecretByID(ctx context.Context, id uint) (*models.Secret, error)

	Update(ctx context.Context, secret *models.Secret) error

	Delete(ctx context.Context, id uint) error
}
