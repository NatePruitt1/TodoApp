package models

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	Raw       string
	Hash      string
	UserId    uuid.UUID
	ExpiresAt time.Time
	IssuedAt  time.Time
}
