package handlers

import (
	"backend/internal/dto"
	"backend/internal/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func getCardId(c *gin.Context) (uuid.UUID, error) {
	cardId := c.Param("cardid")
	if cardId == "" {
		return uuid.Nil, errors.New("no card ID parameter provided.")
	}

	cardUUID, err := uuid.Parse(cardId)
	if err != nil {
		return uuid.Nil, errors.New("card ID provided invalid.")
	}

	return cardUUID, nil
}

func (ph *ProjectHandlerImpl) AddCard(c *gin.Context) {
	var req dto.CardRequest

	userUUID, err := getUserId(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.BadRequestError("failed to retrieve user id.", err.Error()))
		return
	}

	categoryUUID, err := getCategoryId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequestError("failed to retrieve category id", err.Error()))
		return
	}

	err = ph.ProjectService.CheckCategoryOwner(categoryUUID, userUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequestError("user does not own category", err.Error()))
		return
	}

	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequestError("failed to parse request body.", err.Error()))
		return
	}

	card, err := ph.ProjectService.AddCard(categoryUUID, req.Title, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BadRequestError("failed to create card.", err.Error()))
		return
	}

	resp := dto.NewCardResponse(card)

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   resp,
	})
}

func (ph *ProjectHandlerImpl) RenameCard(c *gin.Context) {
	var req dto.CardRequest

	userUUID, err := getUserId(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.BadRequestError("failed to retrieve user id.", err.Error()))
		return
	}

	cardUUID, err := getCardId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequestError("failed to retrieve card id", err.Error()))
		return
	}

	err = ph.ProjectService.CheckCardOwner(cardUUID, userUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequestError("user does not own card", err.Error()))
		return
	}

	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequestError("failed to parse request body.", err.Error()))
		return
	}

	var card *models.Card

	if req.Title != "" {
		card, err = ph.ProjectService.RenameCard(cardUUID, req.Title)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.BadRequestError("failed to rename card.", err.Error()))
			return
		}
	}

	if req.Content != "" {
		card, err = ph.ProjectService.EditCard(cardUUID, req.Content)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.BadRequestError("failed to edit card content.", err.Error()))
			return
		}
	}

	resp := dto.NewCardResponse(card)

	c.JSON(http.StatusAccepted, gin.H{
		"status": "success",
		"data":   resp,
	})
}

func (ph *ProjectHandlerImpl) EditCard(c *gin.Context) {
	var req dto.CardRequest

	userUUID, err := getUserId(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.BadRequestError("failed to retrieve user id.", err.Error()))
		return
	}

	cardUUID, err := getCardId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequestError("failed to retrieve card id", err.Error()))
		return
	}

	err = ph.ProjectService.CheckCardOwner(cardUUID, userUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequestError("user does not own card", err.Error()))
		return
	}

	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequestError("failed to parse request body.", err.Error()))
		return
	}

	card, err := ph.ProjectService.EditCard(cardUUID, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BadRequestError("failed to edit card.", err.Error()))
		return
	}

	resp := dto.NewCardResponse(card)

	c.JSON(http.StatusAccepted, gin.H{
		"status": "success",
		"data":   resp,
	})
}

func (ph *ProjectHandlerImpl) DeleteCard(c *gin.Context) {
	userUUID, err := getUserId(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.BadRequestError("failed to retrieve user id.", err.Error()))
		return
	}

	cardUUID, err := getCardId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequestError("failed to retrieve card id", err.Error()))
		return
	}

	err = ph.ProjectService.CheckCardOwner(cardUUID, userUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequestError("user does not own card", err.Error()))
		return
	}

	err = ph.ProjectService.DeleteCard(cardUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BadRequestError("failed to delete card.", err.Error()))
		return
	}

	c.Status(http.StatusAccepted)
}

func (ph *ProjectHandlerImpl) MoveCard(c *gin.Context) {
	var req dto.CardMoveRequest

	userUUID, err := getUserId(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.BadRequestError("failed to retrieve user id.", err.Error()))
		return
	}

	cardUUID, err := getCardId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequestError("failed to retrieve card id", err.Error()))
		return
	}

	err = ph.ProjectService.CheckCardOwner(cardUUID, userUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequestError("user does not own card", err.Error()))
		return
	}

	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequestError("failed to parse request body.", err.Error()))
		return
	}

	card, err := ph.ProjectService.MoveCard(cardUUID, req.CategoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BadRequestError("failed to move card.", err.Error()))
		return
	}

	resp := dto.NewCardResponse(card)

	c.JSON(http.StatusAccepted, gin.H{
		"status": "success",
		"data":   resp,
	})
}

func (ph *ProjectHandlerImpl) FinishCard(c *gin.Context) {
	var req dto.CardFinishRequest

	userUUID, err := getUserId(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.BadRequestError("failed to retrieve user id.", err.Error()))
		return
	}

	cardUUID, err := getCardId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequestError("failed to retrieve card id", err.Error()))
		return
	}

	err = ph.ProjectService.CheckCardOwner(cardUUID, userUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequestError("user does not own card", err.Error()))
		return
	}

	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequestError("failed to parse request body.", err.Error()))
		return
	}

	// TODO: Implement a service method to set/toggle card finished status
	// For now, return a success response structure
	c.JSON(http.StatusAccepted, gin.H{
		"status":  "success",
		"message": "finished status update not yet implemented",
	})
}
