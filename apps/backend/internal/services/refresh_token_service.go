package services

import (
	"backend/internal/config"
	"backend/internal/models"
	"backend/internal/repository"
	"errors"
	"time"

	"github.com/google/uuid"
)

type RefreshTokenService interface {
	CreateToken(id uuid.UUID, userId uuid.UUID) (models.RefreshToken, string, error)
	CheckToken(id uuid.UUID) (models.RefreshToken, string, error)
	UpdateToken(id uuid.UUID) (models.RefreshToken, string, error)
	GetJWT(token models.RefreshToken) (string, error)
	ParseJWT(token string) (models.RefreshToken, error)
}

type RefreshTokenServiceImpl struct {
	repo   repository.RefreshTokenRepository
	config *config.Config
}

func NewRefreshTokenService(repository repository.RefreshTokenRepository, cfg *config.Config) RefreshTokenService {
	return &RefreshTokenServiceImpl{
		repo:   repository,
		config: cfg,
	}
}

func (rts *RefreshTokenServiceImpl) GetJWT(token models.RefreshToken) (string, error) {
	return GenerateRefreshToken(token, rts.config.JwtSecret)
}

func (rts *RefreshTokenServiceImpl) ParseJWT(token string) (models.RefreshToken, error) {
	return ParseRefreshToken(token, rts.config.JwtSecret)
}

func (rts *RefreshTokenServiceImpl) CreateToken(id uuid.UUID, userId uuid.UUID) (models.RefreshToken, string, error) {
	token := models.RefreshToken{
		Id:        id,
		UserId:    userId,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		IssuedAt:  time.Now(),
	}

	tokenStr, err := rts.GetJWT(token)
	if err != nil {
		return models.RefreshToken{}, "", err
	}

	err = rts.repo.Create(id, tokenStr)
	if err != nil {
		return models.RefreshToken{}, "", err
	}

	return token, tokenStr, nil
}

func (rts *RefreshTokenServiceImpl) CheckToken(id uuid.UUID) (models.RefreshToken, string, error) {
	token, err := rts.repo.Get(id)
	if err != nil {
		return models.RefreshToken{}, "", err
	}

	if token.ExpiresAt.After(time.Now()) {
		rts.repo.Delete(id)
		return models.RefreshToken{}, "", errors.New("Token is expired.")
	}

	tokenStr, err := rts.GetJWT(token)
	if err != nil {
		return models.RefreshToken{}, "", err
	}

	return token, tokenStr, nil
}

func (rts *RefreshTokenServiceImpl) UpdateToken(id uuid.UUID) (models.RefreshToken, string, error) {
	token, _, err := rts.CheckToken(id)
	if err != nil {
		return models.RefreshToken{}, "", err
	}

	rts.repo.Delete(id)
	token.Id = uuid.New()

	return rts.CreateToken(token.Id, token.UserId)
}
