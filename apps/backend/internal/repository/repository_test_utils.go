package repository

import (
	"backend/internal/models"
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func setTestDB(t *testing.T) (*pgxpool.Pool, *migrate.Migrate) {
	t.Helper()

	ctx := context.Background()
	ctr, err := postgres.Run(
		ctx,
		"postgres:alpine",
		postgres.WithDatabase("kanban"),
		postgres.WithUsername("nate"),
		postgres.WithPassword("TestPassword"),
		postgres.BasicWaitStrategies(),
	)
	require.NoError(t, err)

	connStr, err := ctr.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	pool, err := pgxpool.New(ctx, connStr)
	require.NoError(t, err)

	migrationsDir, err := filepath.Abs(filepath.Join("..", "db", "migrations"))
	require.NoError(t, err)

	m, err := migrate.New("file://"+filepath.ToSlash(migrationsDir), connStr)
	require.NoError(t, err)

	require.NoError(t, m.Up())

	return pool, m
}

func makeTestUser(t *testing.T, pool *pgxpool.Pool, username, passwordHash string) models.User {
	userRepo := NewUserRepository(pool)
	user := models.User{
		ID:           uuid.New(),
		Username:     username,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
		LastLogin:    nil,
	}

	err := userRepo.Create(user)
	require.NoError(t, err)

	return user
}

func makeTestProject(t *testing.T, pool *pgxpool.Pool, name string, owner models.User, description string) models.Project {
	projectRepo := NewProjectRepository(pool)
	proj := models.Project{
		ID:          uuid.New(),
		Name:        name,
		OwnerID:     owner.ID,
		Description: description,
		Categories:  nil,
	}
	require.NoError(t, projectRepo.Save(&proj))
	return proj
}
