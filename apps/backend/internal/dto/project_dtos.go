package dto

import "backend/internal/models"

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
