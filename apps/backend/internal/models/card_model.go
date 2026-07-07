package models

import "github.com/google/uuid"

type Card struct {
	ID         uuid.UUID
	CategoryID uuid.UUID
	Title      string
	Content    string
	Finished   bool
}
