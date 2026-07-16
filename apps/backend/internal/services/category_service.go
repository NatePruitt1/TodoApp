package services

import (
	"backend/internal/models"
	"errors"

	"github.com/google/uuid"
)

func (ps *ProjectServiceImpl) CheckCategoryOwner(categoryId, userId uuid.UUID) error {
	return ps.ProjectRepository.CheckCategoryOwner(categoryId, userId)
}

func (ps *ProjectServiceImpl) GetCategory(categoryId uuid.UUID) (*models.Category, error) {
	return ps.ProjectRepository.GetCategoryByID(categoryId)
}

func (ps *ProjectServiceImpl) AddCategory(projectId uuid.UUID, categoryName string) (*models.Category, error) {
	project, err := ps.ProjectRepository.GetAggregate(projectId)
	if err != nil {
		return nil, err
	}

	nextIndex := len(project.Categories)

	newCategory := models.Category{
		ID:        uuid.New(),
		ProjectID: project.ID,
		Name:      categoryName,
		Index:     nextIndex,
	}

	err = ps.ProjectRepository.SaveCategory(&newCategory)
	if err != nil {
		return nil, err
	}

	return &newCategory, nil
}

func (ps *ProjectServiceImpl) DeleteCategory(categoryId uuid.UUID) error {
	return ps.ProjectRepository.Delete(categoryId)
}

func (ps *ProjectServiceImpl) MoveCategory(categoryId uuid.UUID, index int) (*models.Category, error) {
	category, err := ps.ProjectRepository.GetCategoryByID(categoryId)
	if err != nil {
		return nil, err
	}

	project, err := ps.ProjectRepository.GetAggregate(category.ProjectID)
	if err != nil {
		return nil, err
	}

	// Check if the index is < len project.categories
	if index < len(project.Categories) {
		for c := index; c < len(project.Categories); c += 1 {
			project.Categories[c].Index += 1
			ps.ProjectRepository.SaveCategory(project.Categories[c])
		}

		category.Index = index
		ps.ProjectRepository.SaveCategory(category)

		return category, nil
	} else {
		return nil, errors.New("Index is outside of bounds.")
	}
}

func (ps *ProjectServiceImpl) RenameCategory(categoryId uuid.UUID, name string) (*models.Category, error) {
	category, err := ps.ProjectRepository.GetCategoryByID(categoryId)
	if err != nil {
		return nil, err
	}

	category.Name = name
	err = ps.ProjectRepository.SaveCategory(category)
	if err != nil {
		return nil, err
	}
	return category, nil
}
