package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	JwtSecret   string
	DatabaseURL string
}

var ErrNoJWTSecret = errors.New("Could not load JWT token.")
var ErrNoDatabaseURL = errors.New("Could not load Database url.")

func LoadConfigFromEnv() (*Config, error) {
	jwtSecret, ok := os.LookupEnv("JWT_SECRET")
	if !ok || strings.TrimSpace(jwtSecret) == "" {
		return nil, ErrNoJWTSecret
	}

	databaseURL, ok := os.LookupEnv("DATABASE_URL")
	if !ok || strings.TrimSpace(databaseURL) == "" {
		return nil, ErrNoDatabaseURL
	}
	cfg := &Config{
		JwtSecret:   jwtSecret,
		DatabaseURL: databaseURL,
	}

	fmt.Printf("config %v\n", cfg)

	return cfg, nil
}
