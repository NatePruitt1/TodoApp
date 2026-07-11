package models

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	Id        uuid.UUID
	UserId    uuid.UUID
	ExpiresAt time.Time
	IssuedAt  time.Time
}
