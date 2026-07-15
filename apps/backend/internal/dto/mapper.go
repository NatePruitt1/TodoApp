package dto

import "backend/internal/models"

func NewProjectResponse(project *models.Project) ProjectResponse {

	var categories []CategoryResponse = nil
	if project.Categories != nil {
		categories = make([]CategoryResponse, len(project.Categories))
		for c := range categories {
			category := project.Categories[c]
			categories[c] = NewCategoryResponse(category)
		}
	}

	return ProjectResponse{
		ID:         project.ID,
		OwnerID:    project.OwnerID,
		Name:       project.Name,
		Categories: categories,
	}
}

func NewCategoryResponse(category *models.Category) CategoryResponse {
	var cards []CardResponse = nil
	if category.Cards != nil {
		cards = make([]CardResponse, len(category.Cards))
		for c := range cards {
			card := category.Cards[c]
			cards[c] = NewCardResponse(card)
		}
	}

	return CategoryResponse{
		ID:        category.ID,
		ProjectID: category.ProjectID,
		Name:      category.Name,
		Index:     category.Index,
		Cards:     cards,
	}
}

func NewCardResponse(card *models.Card) CardResponse {
	return CardResponse{
		ID:         card.ID,
		CategoryID: card.CategoryID,
		Title:      card.Title,
		Content:    card.Content,
		Finished:   card.Finished,
	}
}

func NewProjectListItem(project *models.Project) ProjectListItem {
	return ProjectListItem{
		ID:      project.ID,
		OwnerID: project.OwnerID,
		Name:    project.Name,
	}
}
