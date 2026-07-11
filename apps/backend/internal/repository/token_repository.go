package repository

import (
	"backend/internal/models"
)

type RefreshTokenRepository interface {
	Create(token models.RefreshToken) error
	Get(id string) (models.RefreshToken, error)
	Delete(id string) error
}
