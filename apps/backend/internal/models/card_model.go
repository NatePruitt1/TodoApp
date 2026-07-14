package models

import "github.com/google/uuid"

type Card struct {
	ID         uuid.UUID `json:"id"`
	CategoryID uuid.UUID `json:"category_id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Finished   bool      `json:"finished"`
}
