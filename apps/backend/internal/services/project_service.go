package services

import (
	"backend/internal/dto"
	"backend/internal/models"
	"backend/internal/repository"
	"errors"
	"slices"

	"github.com/google/uuid"
)

type ProjectService interface {
	GetAllUserProjects(userId uuid.UUID) (dto.AllUserProjectsDTO, error)
	DeleteProject(userId, projectId uuid.UUID) error
	AddProject(userId uuid.UUID, name, description string) (dto.NewProjectResponseDTO, error)
}

type ProjectServiceImpl struct {
	ProjectRepository  repository.ProjectRepository
	CategoryRepository repository.CategoryRepository
	CardRepository     repository.CardRepository
}

func NewProjectService(repo repository.ProjectRepository) ProjectService {
	return &ProjectServiceImpl{
		ProjectRepository: repo,
	}
}

func (ps *ProjectServiceImpl) GetAllUserProjects(userId uuid.UUID) (dto.AllUserProjectsDTO, error) {
	projects, err := ps.ProjectRepository.GetAllUserProjects(userId)
	if err != nil {
		return dto.AllUserProjectsDTO{}, err
	}

	return dto.AllUserProjectsDTO{
		Projects: projects,
	}, nil
}

func (ps *ProjectServiceImpl) DeleteProject(userId, projectId uuid.UUID) error {
	project, err := ps.ProjectRepository.GetByID(projectId)
	if err != nil {
		return err
	}

	if project.OwnerID != userId {
		return errors.New("User does not own project.")
	}

	return ps.ProjectRepository.Delete(projectId)
}

func (ps *ProjectServiceImpl) AddProject(userId uuid.UUID, name, description string) (dto.NewProjectResponseDTO, error) {
	projects, err := ps.ProjectRepository.GetAllUserProjects(userId)
	if err != nil {
		return dto.NewProjectResponseDTO{}, err
	}

	exists := slices.ContainsFunc(projects, func(p models.Project) bool { return p.Name == name })
	if exists {
		return dto.NewProjectResponseDTO{}, errors.New("Project name taken.")
	}

	newProject := models.Project{
		ID:          uuid.New(),
		OwnerID:     userId,
		Name:        name,
		Description: description,
	}

	err = ps.ProjectRepository.Create(newProject)
	if err != nil {
		return dto.NewProjectResponseDTO{}, err
	}

	return dto.NewProjectResponseDTO{
		Project: newProject,
	}, nil
}
