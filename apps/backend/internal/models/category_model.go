package models

import "github.com/google/uuid"

type Category struct {
	ID        uuid.UUID `json:"id"`
	ProjectID uuid.UUID `json:"project_id"`
	Name      string    `json:"name"`
	Index     int       `json:"index"`
	Cards     []*Card
}
