package repository

import (
	"backend/internal/models"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	Create(user models.User) error
	GetByUsername(username string) (models.User, error)
	UpdateLastLogin(userID uuid.UUID, at time.Time) error
}

type UserRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		db: pool,
	}
}

func (r *UserRepositoryImpl) Create(user models.User) error {
	const q = `
		INSERT INTO users (id, username, password_hash, created_at, last_login)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(context.Background(), q,
		user.ID,
		user.Username,
		user.PasswordHash,
		user.CreatedAt,
		user.LastLogin,
	)

	return err
}

func (r *UserRepositoryImpl) Delete(id uuid.UUID) error {
	const q = `
		DELETE FROM users
		WHERE id = $1
	`

	tag, err := r.db.Exec(context.Background(), q, id)
	if err != nil {
		return err
	} else if tag.RowsAffected() == 0 {
		return ErrUserNotFound
	} else {
		return nil
	}
}

func (r *UserRepositoryImpl) GetbyUsername(username string) (models.User, error) {
	const q = `
		SELECT id, username, password_hash, created_at, last_login 
		FROM users
		WHERE username = $1
	`

	var user models.User

	err := r.db.QueryRow(context.Background(), q, username).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.LastLogin,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return models.User{}, ErrUserNotFound
	} else if err != nil {
		return models.User{}, err
	} else {
		return user, nil
	}
}

func (r *UserRepositoryImpl) GetByID(id uuid.UUID) (models.User, error) {
	const q = `
		SELECT id, username, password_hash, created_at, last_login
		FROM users
		WHERE id = $1
	`

	var user models.User

	err := r.db.QueryRow(context.Background(), q, id).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.LastLogin,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return models.User{}, ErrUserNotFound
	} else if err != nil {
		return models.User{}, err
	} else {
		return user, nil
	}
}

func (r *UserRepositoryImpl) UpdateLastLogin(userID uuid.UUID, at time.Time) error {
	const q = `
		UPDATE users
		SET last_login = $1
		WHERE id = $2
	`

	tag, err := r.db.Exec(context.Background(), q, at, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}
