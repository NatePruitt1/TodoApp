package handlers

import (
	"backend/internal/dto"
	"backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProjectHandler interface {
	GetAllProjectsHandler(c *gin.Context)
	DeleteProjectHandler(c *gin.Context)
	AddProjectHandler(c *gin.Context)
}

type ProjectHandlerImpl struct {
	ProjectService services.ProjectService
}

func NewProjectHandler(service services.ProjectService) ProjectHandler {
	return &ProjectHandlerImpl{
		ProjectService: service,
	}
}

func (ph *ProjectHandlerImpl) GetAllProjectsHandler(c *gin.Context) {
	userId := c.GetString("user_id")
	if userId == "" {
		c.JSON(http.StatusUnauthorized, dto.BadRequestError("No userid provided.", ""))
		return
	}

	userUUID, err := uuid.Parse(userId)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.BadRequestError("Invalid userid provided.", err.Error()))
		return
	}

	projects, err := ph.ProjectService.GetAllUserProjects(userUUID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.BadRequestError("Error getting projects.", err.Error()))
		return
	}

	c.JSON(http.StatusFound, gin.H{
		"status": "success",
		"data":   projects,
	})
}

func (ph *ProjectHandlerImpl) DeleteProjectHandler(c *gin.Context) {

}

func (ph *ProjectHandlerImpl) AddProjectHandler(c *gin.Context) {
	var resp dto.NewProjectResponseDTO
	var req dto.NewProjectRequestDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequestError("Bad request body.", err.Error()))
		return
	}

	userId := c.GetString("user_id")
	if userId == "" {
		c.JSON(http.StatusUnauthorized, dto.BadRequestError("No userid provided.", ""))
		return
	}

	userUUID, err := uuid.Parse(userId)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.BadRequestError("Invalid userid provided.", err.Error()))
		return
	}

	resp, err = ph.ProjectService.AddProject(userUUID, req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BadRequestError("Failed to create project.", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   resp,
	})
}
