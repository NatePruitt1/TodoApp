package models

import (
	"github.com/google/uuid"
)

type Project struct {
	ID          uuid.UUID `json:"id"`
	OwnerID     uuid.UUID `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Categories  []*Category
}
