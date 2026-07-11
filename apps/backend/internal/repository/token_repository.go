package repository

import (
	"backend/internal/models"

	"github.com/google/uuid"
)

type RefreshTokenRepository interface {
	Create(id uuid.UUID, token string) error
	Get(id uuid.UUID) (models.RefreshToken, error)
	Delete(id uuid.UUID) error
}
