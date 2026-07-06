package dto

import (
	"github.com/google/uuid"
)

type LoginRequestDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponseDTO struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`

	Token string `json:"accessToken"`
}
