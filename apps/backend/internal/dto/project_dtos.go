package dto

import (
	"backend/internal/models"

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

type NewProjectResponseDTO struct {
	Project models.Project `json:"project"`
}

type NewProjectRequestDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ProjectDTO struct {
	Project    models.Project `json:"project"`
	Categories []CategoryDTO  `json:"categories"`
}

type CategoryDTO struct {
	Category models.Category `json:"category"`
	Cards    []models.Card   `json:"cards"`
}

type ProjectRequest struct {
	ProjectId uuid.UUID `json:"id"`
}
