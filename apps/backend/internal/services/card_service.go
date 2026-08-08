package services

import (
	"backend/internal/models"
	"errors"

	"github.com/google/uuid"
)

func (ps *ProjectServiceImpl) CheckCardOwner(cardId, userId uuid.UUID) error {
	return ps.ProjectRepository.CheckCardOwner(cardId, userId)
}

func (ps *ProjectServiceImpl) GetCard(userId, cardId uuid.UUID) (*models.Card, error) {
	err := ps.CheckCardOwner(cardId, userId)
	if err != nil {
		return nil, err
	}
	return ps.ProjectRepository.GetCardByID(cardId)
}

func (ps *ProjectServiceImpl) AddCard(userId, categoryId uuid.UUID, name, content string) (*models.Card, error) {
	err := ps.CheckCategoryOwner(categoryId, userId)
	if err != nil {
		return nil, err
	}
	newCard := models.Card{
		ID:         uuid.New(),
		CategoryID: categoryId,
		Title:      name,
		Content:    content,
		Finished:   false,
	}

	err = ps.ProjectRepository.SaveCard(&newCard)
	if err != nil {
		return nil, err
	}

	return &newCard, nil
}

func (ps *ProjectServiceImpl) DeleteCard(userId, cardId uuid.UUID) error {
	err := ps.CheckCardOwner(cardId, userId)
	if err != nil {
		return err
	}
	return ps.ProjectRepository.DeleteCard(cardId)
}

func (ps *ProjectServiceImpl) MoveCard(userId, cardId, categoryId uuid.UUID) (*models.Card, error) {
	err := ps.CheckCategoryOwner(categoryId, userId)
	if err != nil {
		return nil, err
	}

	err = ps.CheckCardOwner(cardId, userId)
	if err != nil {
		return nil, err
	}

	card, err := ps.ProjectRepository.GetCardByID(cardId)
	if err != nil {
		return nil, err
	}

	currentCategory, err := ps.ProjectRepository.GetCategoryByID(card.CategoryID)
	if err != nil {
		return nil, err
	}

	desiredCategory, err := ps.ProjectRepository.GetCategoryByID(categoryId)
	if err != nil {
		return nil, err
	}

	if currentCategory.ProjectID != desiredCategory.ProjectID {
		return nil, errors.New("Categories must be in the same project.")
	}

	card.CategoryID = categoryId
	err = ps.ProjectRepository.SaveCard(card)
	if err != nil {
		return nil, err
	}

	return card, nil
}

func (ps *ProjectServiceImpl) RenameCard(userId, cardId uuid.UUID, name string) (*models.Card, error) {
	err := ps.CheckCardOwner(cardId, userId)
	if err != nil {
		return nil, err
	}

	card, err := ps.ProjectRepository.GetCardByID(cardId)
	if err != nil {
		return nil, err
	}

	card.Title = name
	err = ps.ProjectRepository.SaveCard(card)
	if err != nil {
		return nil, err
	}

	return card, nil
}

func (ps *ProjectServiceImpl) EditCard(userId, cardId uuid.UUID, content string) (*models.Card, error) {
	err := ps.CheckCardOwner(cardId, userId)
	if err != nil {
		return nil, err
	}

	card, err := ps.ProjectRepository.GetCardByID(cardId)
	if err != nil {
		return nil, err
	}

	card.Content = content
	err = ps.ProjectRepository.SaveCard(card)
	if err != nil {
		return nil, err
	}

	return card, nil
}
