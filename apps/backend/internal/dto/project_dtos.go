package dto

import (
	"github.com/google/uuid"
)

type ProjectResponse struct {
	ID          uuid.UUID          `json:"id"`
	OwnerID     uuid.UUID          `json:"owner_id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Categories  []CategoryResponse `json:"categories"`
}

type CategoryResponse struct {
	ID        uuid.UUID      `json:"id"`
	ProjectID uuid.UUID      `json:"project_id"`
	Name      string         `json:"name"`
	Index     int            `json:"index"`
	Cards     []CardResponse `json:"cards"`
}

type CardResponse struct {
	ID         uuid.UUID `json:"id"`
	CategoryID uuid.UUID `json:"category_id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Finished   bool      `json:"finished"`
}

type ProjectListResponse struct {
	Projects []ProjectResponse `json:"projects"`
}

type ProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CategoryRequest struct {
	Name string `json:"name"`
}

type CategoryMoveRequest struct {
	Index int `json:"index"`
}

type CardRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type CardMoveRequest struct {
	CategoryID uuid.UUID `json:"category_id"`
	Index      int       `json:"index"`
}

type CardFinishRequest struct {
	Finished bool `json:"finished"`
}
