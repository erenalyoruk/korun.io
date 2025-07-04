package service

import (
	"context"

	"korun.io/secret-service/internal/interfaces"
	"korun.io/shared/models"
)

type SecretService struct {
	repo interfaces.SecretRepository
}

func NewSecretService(repo interfaces.SecretRepository) *SecretService {
	return &SecretService{repo: repo}
}

func (s *SecretService) CreateSecret(ctx context.Context, secret *models.Secret) error {
	return s.repo.Create(ctx, secret)
}

func (s *SecretService) GetSecrets(ctx context.Context) ([]models.Secret, error) {
	return s.repo.GetSecrets(ctx)
}

func (s *SecretService) GetSecretByID(ctx context.Context, id uint) (*models.Secret, error) {
	return s.repo.GetSecretByID(ctx, id)
}

func (s *SecretService) GetSecretByName(ctx context.Context, name string) (*models.Secret, error) {
	return s.repo.GetSecretByName(ctx, name)
}

func (s *SecretService) UpdateSecret(ctx context.Context, secret *models.Secret) error {
	return s.repo.Update(ctx, secret)
}

func (s *SecretService) DeleteSecret(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
