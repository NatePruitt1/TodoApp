package services

import (
	"backend/internal/models"
	"errors"
	"fmt"

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
	return ps.ProjectRepository.DeleteCategory(categoryId)
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

	fmt.Printf("Category Index: %d\n", category.Index)

	//Check if new index is valid.
	if index >= 0 && index < len(project.Categories) {
		element := project.Categories[category.Index]

		if element.Index < index {
			copy(project.Categories[element.Index:index], project.Categories[element.Index+1:index+1])
		} else {
			copy(project.Categories[index+1:element.Index+1], project.Categories[index:element.Index])
		}

		project.Categories[index] = element

		for c := range project.Categories {
			if project.Categories[c].Index != c {
				project.Categories[c].Index = c
				ps.ProjectRepository.SaveCategory(project.Categories[c])
			}
		}

		return element, nil
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
