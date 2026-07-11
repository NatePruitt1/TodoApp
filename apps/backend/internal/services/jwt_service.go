package services

import (
	"backend/internal/models"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenClaims struct {
	UserID   string `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateRefreshToken(refresh models.RefreshToken, secret string) (string, error) {
	claims := jwt.RegisteredClaims{
		ID:        refresh.Id.String(),
		Subject:   refresh.UserId.String(),
		ExpiresAt: jwt.NewNumericDate(refresh.ExpiresAt),
		IssuedAt:  jwt.NewNumericDate(refresh.IssuedAt),
		NotBefore: jwt.NewNumericDate(refresh.IssuedAt),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateToken(user models.User, secret string) (string, error) {
	fmt.Println("Starting Gen token")
	claims := TokenClaims{
		UserID:   user.ID.String(),
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseRefreshToken(tokenString, secret string) (models.RefreshToken, error) {
	claims, err := jwt.ParseWithClaims(
		tokenString,
		jwt.RegisteredClaims{},
		func(t *jwt.Token) (any, error) {
			if t.Method.Alg() != jwt.SigningMethodHS512.Alg() {
				return nil, errors.New("unexpected signing method.")
			}
			return []byte(secret), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS512.Alg()}),
		jwt.WithLeeway(5*time.Second),
	)

	if err != nil {
		return models.RefreshToken{}, err
	}

	regClaims, ok := claims.Claims.(*jwt.RegisteredClaims)
	if !ok || !claims.Valid {
		return models.RefreshToken{}, errors.New("invalid token")
	}

	Id, err := uuid.Parse(regClaims.ID)
	if err != nil {
		return models.RefreshToken{}, errors.New("invalid token")
	}

	UserId, err := uuid.Parse(regClaims.Subject)
	if err != nil {
		return models.RefreshToken{}, errors.New("invalid token")
	}

	return models.RefreshToken{
		Id:        Id,
		UserId:    UserId,
		ExpiresAt: regClaims.ExpiresAt.Time,
		IssuedAt:  regClaims.IssuedAt.Time,
	}, nil
}

func ParseToken(tokenString, secret string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&TokenClaims{},
		func(t *jwt.Token) (any, error) {
			if t.Method.Alg() != jwt.SigningMethodHS512.Alg() {
				return nil, errors.New("unexpected signing method.")
			}
			return []byte(secret), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS512.Alg()}),
		jwt.WithLeeway(5*time.Second),
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
