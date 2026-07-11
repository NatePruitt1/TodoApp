package repository

import (
	"backend/internal/models"
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshTokenRepository interface {
	Create(token models.RefreshToken) error
	Get(hash string) (models.RefreshToken, error)
	Delete(hash string) error
}

type RefreshTokenRespositoryImpl struct {
	db *pgxpool.Pool
}

func NewRefreshTokenRepository(db *pgxpool.Pool) RefreshTokenRepository {
	return &RefreshTokenRespositoryImpl{
		db: db,
	}
}

func (rtr *RefreshTokenRespositoryImpl) Create(token models.RefreshToken) error {
	const q = `
		INSERT INTO refresh_tokens (id, user_id, expires_at, created_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := rtr.db.Exec(context.Background(), q,
		token.Hash,
		token.UserId,
		token.ExpiresAt,
		token.IssuedAt,
	)

	return err
}

func (rtr *RefreshTokenRespositoryImpl) Get(hash string) (models.RefreshToken, error) {
	const q = `
		SELECT (id, user_id, expires_at, created_at) FROM refresh_tokens WHERE id = $1
	`

	var token models.RefreshToken

	err := rtr.db.QueryRow(context.Background(), q, hash).Scan(
		&token.Hash,
		&token.UserId,
		&token.ExpiresAt,
		&token.IssuedAt,
	)

	if err != nil {
		return models.RefreshToken{}, err
	}

	return token, nil
}

func (rtr *RefreshTokenRespositoryImpl) Delete(hash string) error {
	const q = `
		DELETE FROM refresh_tokens
		WHERE id = $1
	`

	tag, err := rtr.db.Exec(context.Background(), q, hash)
	if err != nil {
		return err
	} else if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	} else {
		return nil
	}

}
