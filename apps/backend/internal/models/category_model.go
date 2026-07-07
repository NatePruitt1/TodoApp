package models

import "github.com/google/uuid"

type Category struct {
	ID        uuid.UUID
	ProjectID uuid.UUID
	Name      string
	Index     int
}
