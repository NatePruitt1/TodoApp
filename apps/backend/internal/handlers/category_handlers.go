package handlers

import (
	"backend/internal/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (ph *ProjectHandlerImpl) AddCategory(c *gin.Context) {
	var req dto.CategoryRequest

	userUUID, err := getUserId(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.BadRequestError("failed to retrieve user id.", err.Error()))
		return
	}

	projectUUID, err := getProjectId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequestError("failed to retrieve project id", err.Error()))
		return
	}

	err = ph.ProjectService.CheckProjectOwner(projectUUID, userUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequestError("user does not own project", err.Error()))
		return
	}

	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequestError("failed to parse request body.", err.Error()))
		return
	}

	cat, err := ph.ProjectService.AddCategory(projectUUID, req.Name)
	resp := dto.NewCategoryResponse(cat)

	c.JSON(http.StatusAccepted, gin.H{
		"status": "success",
		"data":   resp,
	})
}

func (ph *ProjectHandlerImpl) UpdateCategory(c *gin.Context) {
	var req dto.CategoryRequest

	userUUID, err := getUserId(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.BadRequestError("failed to retrieve user id.", err.Error()))
		return
	}

	catergoryUUID, err := getCategoryId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequestError("failed to retrieve category id", err.Error()))
		return
	}

	err = ph.ProjectService.CheckCategoryOwner(catergoryUUID, userUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequestError("user does not own project", err.Error()))
		return
	}

	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequestError("failed to parse request body.", err.Error()))
		return
	}

	cat, err := ph.ProjectService.RenameCategory(catergoryUUID, req.Name)
	resp := dto.NewCategoryResponse(cat)

	c.JSON(http.StatusAccepted, gin.H{
		"status": "success",
		"data":   resp,
	})
}

func (ph *ProjectHandlerImpl) DeleteCategory(c *gin.Context) {
	userUUID, err := getUserId(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.BadRequestError("failed to retrieve user id.", err.Error()))
		return
	}

	catergoryUUID, err := getCategoryId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequestError("failed to retrieve category id", err.Error()))
		return
	}

	err = ph.ProjectService.CheckCategoryOwner(catergoryUUID, userUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequestError("user does not own project", err.Error()))
		return
	}

	err = ph.ProjectService.DeleteCategory(catergoryUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BadRequestError("failed to delete category.", err.Error()))
		return
	}

	c.Status(http.StatusNoContent)
}

func (ph *ProjectHandlerImpl) MoveCategory(c *gin.Context) {
	var req dto.CategoryMoveRequest

	userUUID, err := getUserId(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.BadRequestError("failed to retrieve user id.", err.Error()))
		return
	}

	catergoryUUID, err := getCategoryId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequestError("failed to retrieve category id", err.Error()))
		return
	}

	err = ph.ProjectService.CheckCategoryOwner(catergoryUUID, userUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequestError("user does not own project", err.Error()))
		return
	}

	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequestError("failed to parse request body.", err.Error()))
		return
	}

	cat, err := ph.ProjectService.MoveCategory(catergoryUUID, req.Index)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequestError("failed to move category.", err.Error()))
		return
	}
	resp := dto.NewCategoryResponse(cat)

	c.JSON(http.StatusAccepted, gin.H{
		"status": "success",
		"data":   resp,
	})
}
