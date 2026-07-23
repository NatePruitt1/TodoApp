package config

import (
	"errors"
	"net/url"
	"os"
	"strings"
)

type Config struct {
	JwtSecret   string
	DatabaseURL string
}

var ErrNoJWTSecret = errors.New("Could not load JWT token.")
var ErrIncompleteDBConfig = errors.New("Could not load database configuration: POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_HOST and POSTGRES_DB must all be set.")

// LoadConfigFromEnv builds the application config from environment variables.
//
// The database connection is assembled from the same POSTGRES_USER,
// POSTGRES_PASSWORD, POSTGRES_HOST, POSTGRES_PORT and POSTGRES_DB variables
// that provision the postgres-db container (see internal/db/.env.db and
// docker-compose.yml), so the backend always authenticates with the same
// credentials the database was actually created with. DATABASE_URL can
// still be set directly to override this (e.g. for pointing at an external
// database), but it is no longer the source of truth for local/docker use.
func LoadConfigFromEnv() (*Config, error) {
	jwtSecret, ok := os.LookupEnv("JWT_SECRET")
	if !ok || strings.TrimSpace(jwtSecret) == "" {
		return nil, ErrNoJWTSecret
	}

	databaseURL, err := buildDatabaseURL()
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		JwtSecret:   jwtSecret,
		DatabaseURL: databaseURL,
	}

	return cfg, nil
}

func buildDatabaseURL() (string, error) {
	if explicit, ok := os.LookupEnv("DATABASE_URL"); ok && strings.TrimSpace(explicit) != "" {
		return explicit, nil
	}

	user := strings.TrimSpace(os.Getenv("POSTGRES_USER"))
	password := os.Getenv("POSTGRES_PASSWORD")
	host := strings.TrimSpace(os.Getenv("POSTGRES_HOST"))
	dbName := strings.TrimSpace(os.Getenv("POSTGRES_DB"))
	port := strings.TrimSpace(os.Getenv("POSTGRES_PORT"))
	if port == "" {
		port = "5432"
	}

	if user == "" || strings.TrimSpace(password) == "" || host == "" || dbName == "" {
		return "", ErrIncompleteDBConfig
	}

	dsn := url.URL{
		Scheme:   "postgresql",
		User:     url.UserPassword(user, password),
		Host:     host + ":" + port,
		Path:     "/" + dbName,
		RawQuery: "sslmode=disable",
	}

	return dsn.String(), nil
}
