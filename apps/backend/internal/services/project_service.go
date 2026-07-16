package services

import (
	"backend/internal/models"
	"backend/internal/repository"
	"errors"
	"slices"

	"github.com/google/uuid"
)

type ProjectService interface {
	GetProjects(userId uuid.UUID) ([]*models.Project, error)
	GetProject(userId, projectId uuid.UUID) (*models.Project, error)
	DeleteProject(userId, projectId uuid.UUID) error
	AddProject(userId uuid.UUID, name, description string) (*models.Project, error)
	CheckProjectOwner(projectId, userId uuid.UUID) error

	// Category Services
	GetCategory(categoryId uuid.UUID) (*models.Category, error)
	AddCategory(projectId uuid.UUID, categoryName string) (*models.Category, error)
	DeleteCategory(categoryId uuid.UUID) error
	MoveCategory(categoryId uuid.UUID, index int) (*models.Category, error)
	RenameCategory(categoryId uuid.UUID, name string) (*models.Category, error)
	CheckCategoryOwner(categoryId, userId uuid.UUID) error

	// Card Services
	GetCard(cardId uuid.UUID) (*models.Card, error)
	AddCard(categoryId uuid.UUID, name, content string) (*models.Card, error)
	DeleteCard(cardId uuid.UUID) error
	MoveCard(cardId, categoryId uuid.UUID) (*models.Card, error)
	RenameCard(cardId uuid.UUID, name string) (*models.Card, error)
	EditCard(cardId uuid.UUID, content string) (*models.Card, error)
	CheckCardOwner(cardId, userId uuid.UUID) error
}

type ProjectServiceImpl struct {
	ProjectRepository repository.ProjectRepository
}

func NewProjectService(repo repository.ProjectRepository) ProjectService {
	return &ProjectServiceImpl{
		ProjectRepository: repo,
	}
}

func (ps *ProjectServiceImpl) GetProjects(userId uuid.UUID) ([]*models.Project, error) {
	projects, err := ps.ProjectRepository.GetAllUserProjects(userId)
	if err != nil {
		return nil, err
	}

	return projects, nil
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

func (ps *ProjectServiceImpl) AddProject(userId uuid.UUID, name, description string) (*models.Project, error) {
	projects, err := ps.ProjectRepository.GetAllUserProjects(userId)
	if err != nil {
		return nil, err
	}

	exists := slices.ContainsFunc(projects, func(p *models.Project) bool { return p.Name == name })
	if exists {
		return nil, errors.New("Project name taken.")
	}

	newProject := models.Project{
		ID:          uuid.New(),
		OwnerID:     userId,
		Name:        name,
		Description: description,
	}

	err = ps.ProjectRepository.Save(&newProject)
	if err != nil {
		return nil, err
	}

	return &newProject, nil
}

func (ps *ProjectServiceImpl) GetProject(userId, projectId uuid.UUID) (*models.Project, error) {
	project, err := ps.ProjectRepository.GetAggregate(projectId)
	if err != nil {
		return nil, err
	}

	if project.OwnerID != userId {
		return nil, errors.New("User does not own project.")
	}

	return project, nil
}

func (ps *ProjectServiceImpl) CheckProjectOwner(projectId, userId uuid.UUID) error {
	return ps.ProjectRepository.CheckProjectOwner(projectId, userId)
}
