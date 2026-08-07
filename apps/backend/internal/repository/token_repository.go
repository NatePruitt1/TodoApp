package repository

import (
	"backend/internal/models"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
)

type RefreshTokenRepository interface {
	Create(token models.RefreshToken) error
	Get(hash string) (models.RefreshToken, error)
	GetByUser(userId uuid.UUID) (models.RefreshToken, error)
	Delete(hash string) error
}

type RefreshTokenRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewRefreshTokenRepository(db *pgxpool.Pool) RefreshTokenRepository {
	return &RefreshTokenRepositoryImpl{
		db: db,
	}
}

func (rtr *RefreshTokenRepositoryImpl) GetByUser(userId uuid.UUID) (models.RefreshToken, error) {
	const q = `
		SELECT hash, user_id, expires_at, created_at FROM refresh_tokens WHERE user_id = $1
	`

	var token models.RefreshToken

	err := rtr.db.QueryRow(context.Background(), q, userId).Scan(
		&token.Hash,
		&token.UserId,
		&token.ExpiresAt,
		&token.IssuedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return models.RefreshToken{}, ErrRefreshTokenNotFound
	} else if err != nil {
		return models.RefreshToken{}, fmt.Errorf("GetRefreshTokenByUser: %w", err)
	} else {
		return token, nil
	}
}

func (rtr *RefreshTokenRepositoryImpl) Create(token models.RefreshToken) error {
	const q = `
		INSERT INTO refresh_tokens (hash, user_id, expires_at, created_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := rtr.db.Exec(context.Background(), q,
		token.Hash,
		token.UserId,
		token.ExpiresAt,
		token.IssuedAt,
	)

	if err != nil {
		return fmt.Errorf("CreateRefreshToken: %w", err)
	}

	return nil
}

func (rtr *RefreshTokenRepositoryImpl) Get(hash string) (models.RefreshToken, error) {
	const q = `
		SELECT hash, user_id, expires_at, created_at FROM refresh_tokens WHERE hash = $1
	`

	var token models.RefreshToken

	err := rtr.db.QueryRow(context.Background(), q, hash).Scan(
		&token.Hash,
		&token.UserId,
		&token.ExpiresAt,
		&token.IssuedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return models.RefreshToken{}, ErrRefreshTokenNotFound
	} else if err != nil {
		return models.RefreshToken{}, fmt.Errorf("GetRefreshToken: %w", err)
	} else {
		return token, nil
	}
}

func (rtr *RefreshTokenRepositoryImpl) Delete(hash string) error {
	const q = `
		DELETE FROM refresh_tokens
		WHERE hash = $1
	`

	tag, err := rtr.db.Exec(context.Background(), q, hash)
	if err != nil {
		return fmt.Errorf("DeleteRefreshToken: %w", err)
	} else if tag.RowsAffected() == 0 {
		return ErrRefreshTokenNotFound
	} else {
		return nil
	}
}
