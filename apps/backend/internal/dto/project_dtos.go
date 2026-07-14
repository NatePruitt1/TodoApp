package dto

import (
	"backend/internal/models"
)

type AllUserProjectsDTO struct {
	Projects []models.Project `json:"projects"`
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
