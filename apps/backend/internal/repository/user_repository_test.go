package repository

import (
	"backend/internal/models"

	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func setTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	pool, err := pgxpool.New(t.Context(), "postgres://postgres:Wallah54639!@localhost:5432/kanban_test")
	if err != nil {
		t.Fatalf("Failed to connect to test db: %v", err)
	}

	t.Cleanup(func() {
		pool.Close()
	})

	return pool
}

func setTestTx(t *testing.T, pool *pgxpool.Pool) pgx.Tx {
	t.Helper()

	tx, err := pool.Begin(context.Background())
	if err != nil {
		t.Fatalf("failed to being tx: %v", err)
	}

	t.Cleanup(func() {
		_ = tx.Rollback(context.Background())
	})

	return tx
}

func TestUserRepository_Create(t *testing.T) {
	pool := setTestDB(t)
	repo := NewUserRepository(pool)

	e := repo.Create(models.User{
		ID:           uuid.New(),
		Username:     "username",
		PasswordHash: "abcd",
		CreatedAt:    time.Now(),
		LastLogin:    nil,
	})

	if e != nil {
		t.Fatalf("Error creating user. %v", e)
	}

	u, e := repo.GetbyUsername("username")
	if e != nil {
		t.Fatalf("Error getting user. %v", e)
	}

	if u.Username != "username" {
		t.Fatalf("Wrong username")
	}

	pool.Exec(context.Background(), `DELETE FROM users`)
}

func TestUserRepository_Delete(t *testing.T) {
	pool := setTestDB(t)
	repo := NewUserRepository(pool)

	t.Cleanup(func() {
		pool.Exec(context.Background(), `DELETE FROM users`)
	})

	ID := uuid.New()

	e := repo.Create(models.User{
		ID:           ID,
		Username:     "username",
		PasswordHash: "abcd",
		CreatedAt:    time.Now(),
		LastLogin:    nil,
	})
	if e != nil {
		t.Fatalf("Error creating user. %v", e)
	}

	e = repo.Delete(ID)
	if e != nil {
		t.Fatalf("Error deleting user. %v", e)
	}

	u, e := repo.GetByID(ID)
	if e == nil {
		t.Fatalf("No error getting user that should have been deleted. Got user %v", u)
	}
}

func TestUserRepository_GetByUsername(t *testing.T) {
	pool := setTestDB(t)
	repo := NewUserRepository(pool)

	_, e := repo.GetbyUsername("UserDoesNotExist")
	if e == nil {
		t.Fatal("Bad retrieval should have returned an error.")
	}

	if !errors.Is(e, ErrUserNotFound) {
		t.Fatalf("Error should have been pgx.ErrNoRows. Instead got %v", e)
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	pool := setTestDB(t)
	repo := NewUserRepository(pool)

	t.Cleanup(func() {
		pool.Exec(context.Background(), `DELETE FROM users`)
	})

	ID := uuid.New()

	e := repo.Create(models.User{
		ID:           ID,
		Username:     "username",
		PasswordHash: "abcd",
		CreatedAt:    time.Now(),
		LastLogin:    nil,
	})
	if e != nil {
		t.Fatalf("Error creating user. %v", e)
	}

	u, e := repo.GetByID(ID)
	if e != nil {
		t.Fatalf("Failed to get the user by ID. %v", e)
	}

	if u.ID != ID {
		t.Fatalf("Failed to get the correct ID. Wanted: %v. Got %v.", ID, u.ID)
	}
}

func TestUserRepository_UpdateLastLogin(t *testing.T) {
	pool := setTestDB(t)
	repo := NewUserRepository(pool)

	t.Cleanup(func() {
		pool.Exec(context.Background(), `DELETE FROM users`)
	})

	ID := uuid.New()

	e := repo.Create(models.User{
		ID:           ID,
		Username:     "username",
		PasswordHash: "abcd",
		CreatedAt:    time.Now(),
		LastLogin:    nil,
	})

	if e != nil {
		t.Fatalf("Failed to create user. %v", e)
	}

	now := time.Now().UTC().Truncate(time.Microsecond)

	e = repo.UpdateLastLogin(ID, now)
	if e != nil {
		t.Fatalf("Failed to update last login. %v", e)
	}

	u, e := repo.GetbyUsername("username")
	if e != nil {
		t.Fatalf("Failed to fetch user. %v", e)
	}

	got := u.LastLogin.UTC().Truncate(time.Microsecond)

	if !got.Equal(now) {
		t.Fatalf("Inequal times. u.LastLogin: %v\t now: %v", u.LastLogin, now)
	}
}
