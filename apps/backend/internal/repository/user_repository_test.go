package repository

import (
	"backend/internal/models"

	"errors"
	"testing"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_Create(t *testing.T) {
	pool, m := setTestDB(t)
	defer func() {
		require.NoError(t, m.Down())
	}()
	repo := NewUserRepository(pool)

	now := time.Now()
	id := uuid.New()

	e := repo.Create(models.User{
		ID:           id,
		Username:     "username",
		PasswordHash: "abcd",
		CreatedAt:    now,
		LastLogin:    nil,
	})

	require.NoError(t, e)

	u, e := repo.GetByUsername("username")
	require.NoError(t, e)

	assert.Equal(t, u.Username, "username")
	assert.Equal(t, u.PasswordHash, "abcd")
	assert.Equal(t, u.ID, id)
	assert.WithinDuration(t, u.CreatedAt, now, 1*time.Second)
}

func TestUserRepository_Delete(t *testing.T) {
	pool, m := setTestDB(t)
	defer func() {
		require.NoError(t, m.Down())
	}()
	repo := NewUserRepository(pool)

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
	pool, m := setTestDB(t)
	defer func() {
		require.NoError(t, m.Down())
	}()
	repo := NewUserRepository(pool)

	_, e := repo.GetByUsername("UserDoesNotExist")
	if e == nil {
		t.Fatal("Bad retrieval should have returned an error.")
	}

	if !errors.Is(e, ErrUserNotFound) {
		t.Fatalf("Error should have been pgx.ErrNoRows. Instead got %v", e)
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	pool, m := setTestDB(t)
	defer func() {
		require.NoError(t, m.Down())
	}()
	repo := NewUserRepository(pool)

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
	pool, m := setTestDB(t)
	defer func() {
		require.NoError(t, m.Down())
	}()
	repo := NewUserRepository(pool)

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

	u, e := repo.GetByUsername("username")
	if e != nil {
		t.Fatalf("Failed to fetch user. %v", e)
	}

	got := u.LastLogin.UTC().Truncate(time.Microsecond)

	if !got.Equal(now) {
		t.Fatalf("Inequal times. u.LastLogin: %v\t now: %v", u.LastLogin, now)
	}
}
