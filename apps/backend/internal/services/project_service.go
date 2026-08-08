package services

import (
	"backend/internal/models"
	"backend/internal/repository"
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
	GetCategory(userId, categoryId uuid.UUID) (*models.Category, error)
	AddCategory(userId, projectId uuid.UUID, categoryName string) (*models.Category, error)
	DeleteCategory(userId, categoryId uuid.UUID) error
	MoveCategory(userId, categoryId uuid.UUID, index int) (*models.Category, error)
	RenameCategory(userId, categoryId uuid.UUID, name string) (*models.Category, error)
	CheckCategoryOwner(categoryId, userId uuid.UUID) error

	// Card Services
	GetCard(userId, cardId uuid.UUID) (*models.Card, error)
	AddCard(userId, categoryId uuid.UUID, name, content string) (*models.Card, error)
	DeleteCard(userId, cardId uuid.UUID) error
	MoveCard(userId, cardId, categoryId uuid.UUID) (*models.Card, error)
	RenameCard(userId, cardId uuid.UUID, name string) (*models.Card, error)
	EditCard(userId, cardId uuid.UUID, content string) (*models.Card, error)
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
		return repository.ErrNotResourceOwner
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
		return nil, repository.ErrProjectNameConflict
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
		return nil, repository.ErrNotResourceOwner
	}

	return project, nil
}

func (ps *ProjectServiceImpl) CheckProjectOwner(projectId, userId uuid.UUID) error {
	return ps.ProjectRepository.CheckProjectOwner(projectId, userId)
}
