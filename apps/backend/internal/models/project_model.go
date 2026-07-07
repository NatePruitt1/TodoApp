package models

import (
	"github.com/google/uuid"
)

type Project struct {
	ID          uuid.UUID
	OwnerID     uuid.UUID
	Name        string
	Description string
}
