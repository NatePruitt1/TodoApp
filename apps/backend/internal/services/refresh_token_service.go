package services

import (
	"backend/internal/config"
	"backend/internal/models"
	"backend/internal/repository"
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"
)

type RefreshTokenService interface {
	CreateToken(userId uuid.UUID) (models.RefreshToken, error)
	CheckToken(id string) (models.RefreshToken, error)
	CheckTokenForUser(userId uuid.UUID) (models.RefreshToken, error)
	UpdateToken(id string) (models.RefreshToken, error)
	UpdateTokenByHash(hash string) (models.RefreshToken, error)
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

func generateRefreshToken() (string, string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", "", err
	}
	raw := hex.EncodeToString(bytes)
	sum := sha512.Sum512([]byte(raw))
	hash := hex.EncodeToString(sum[:])
	return raw, hash, nil
}

func (rts *RefreshTokenServiceImpl) CreateToken(userId uuid.UUID) (models.RefreshToken, error) {
	raw, id, err := generateRefreshToken()
	if err != nil {
		return models.RefreshToken{}, err
	}

	token := models.RefreshToken{
		Raw:       raw,
		Hash:      id,
		UserId:    userId,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		IssuedAt:  time.Now(),
	}

	err = rts.repo.Create(token)
	if err != nil {
		return models.RefreshToken{}, err
	}

	return token, nil
}

func (rts *RefreshTokenServiceImpl) CheckTokenForUser(userId uuid.UUID) (models.RefreshToken, error) {
	token, err := rts.repo.GetByUser(userId)
	if err != nil {
		return models.RefreshToken{}, err
	}

	if time.Now().After(token.ExpiresAt) {
		rts.repo.Delete(token.Hash)
		return models.RefreshToken{}, errors.New("Token is expired.")
	}

	return token, nil
}

func (rts *RefreshTokenServiceImpl) CheckToken(id string) (models.RefreshToken, error) {
	sum := sha512.Sum512([]byte(id))
	hash := hex.EncodeToString(sum[:])

	token, err := rts.repo.Get(hash)
	if err != nil {
		return models.RefreshToken{}, err
	}

	if time.Now().After(token.ExpiresAt) {
		rts.repo.Delete(hash)
		return models.RefreshToken{}, errors.New("Token is expired.")
	}

	return token, nil
}

func (rts *RefreshTokenServiceImpl) UpdateTokenByHash(hash string) (models.RefreshToken, error) {
	token, err := rts.repo.Get(hash)
	if err != nil {
		return models.RefreshToken{}, err
	}

	rts.repo.Delete(hash)

	if time.Now().After(token.ExpiresAt) {
		return models.RefreshToken{}, errors.New("Token is expired.")
	}

	return rts.CreateToken(token.UserId)
}

func (rts *RefreshTokenServiceImpl) UpdateToken(id string) (models.RefreshToken, error) {
	token, err := rts.CheckToken(id)
	if err != nil {
		return models.RefreshToken{}, err
	}

	rts.repo.Delete(token.Hash)
	return rts.CreateToken(token.UserId)
}
