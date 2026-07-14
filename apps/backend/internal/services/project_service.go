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
	GetProjects(userId uuid.UUID) (dto.ProjectListResponse, error)
	GetProject(userId, projectId uuid.UUID) (dto.ProjectResponse, error)
	DeleteProject(userId, projectId uuid.UUID) error
	AddProject(userId uuid.UUID, name, description string) (dto.NewProjectResponseDTO, error)
}

type ProjectServiceImpl struct {
	ProjectRepository repository.ProjectRepository
}

func NewProjectService(repo repository.ProjectRepository) ProjectService {
	return &ProjectServiceImpl{
		ProjectRepository: repo,
	}
}

func (ps *ProjectServiceImpl) GetProjects(userId uuid.UUID) (dto.ProjectListResponse, error) {
	projects, err := ps.ProjectRepository.GetAllUserProjects(userId)
	if err != nil {
		return dto.ProjectListResponse{}, err
	}

	projectList := make([]dto.ProjectListItem, len(projects))
	for p := range projectList {
		project := projects[p]
		projectList[p] = dto.ProjectListItem{
			ID:      project.ID,
			OwnerID: project.OwnerID,
			Name:    project.Name,
		}
	}

	return dto.ProjectListResponse{
		Projects: projectList,
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

	exists := slices.ContainsFunc(projects, func(p *models.Project) bool { return p.Name == name })
	if exists {
		return dto.NewProjectResponseDTO{}, errors.New("Project name taken.")
	}

	newProject := models.Project{
		ID:          uuid.New(),
		OwnerID:     userId,
		Name:        name,
		Description: description,
	}

	err = ps.ProjectRepository.Save(&newProject)
	if err != nil {
		return dto.NewProjectResponseDTO{}, err
	}

	return dto.NewProjectResponseDTO{
		Project: newProject,
	}, nil
}

func (ps *ProjectServiceImpl) GetProject(userId, projectId uuid.UUID) (dto.ProjectResponse, error) {
	project, err := ps.ProjectRepository.GetAggregate(projectId)
	if err != nil {
		return dto.ProjectResponse{}, err
	}

	if project.OwnerID != userId {
		return dto.ProjectResponse{}, errors.New("User does not own project.")
	}

	return dto.NewProjectResponse(project), nil
}
