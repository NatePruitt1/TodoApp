package handlers

import (
	"backend/internal/dto"
	"backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProjectHandler interface {
	GetProjectsHandler(c *gin.Context)
	DeleteProjectHandler(c *gin.Context)
	AddProjectHandler(c *gin.Context)
	GetProjectHandler(c *gin.Context)

	// Category Handlers
	AddCategory(c *gin.Context)
	DeleteCategory(c *gin.Context)
	MoveCategory(c *gin.Context)
	RemameCategory(c *gin.Context)

	// Card Handlers
	AddCard(c *gin.Context)
	DeleteCard(c *gin.Context)
	MoveCard(c *gin.Context)
	RenameCard(c *gin.Context)
	EditCard(c *gin.Context)
	FinishCard(c *gin.Context)
}

type ProjectHandlerImpl struct {
	ProjectService services.ProjectService
}

func NewProjectHandler(service services.ProjectService) ProjectHandler {
	return &ProjectHandlerImpl{
		ProjectService: service,
	}
}

func (ph *ProjectHandlerImpl) GetProjectHandler(c *gin.Context) {
	var projectRequest dto.ProjectRequest

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

	if err = c.ShouldBindJSON(&projectRequest); err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequestError("Bad request body.", err.Error()))
		return
	}

	project, err := ph.ProjectService.GetProject(userUUID, projectRequest.ProjectId)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequestError("Failed to get project.", err.Error()))
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"status": "success",
		"data":   project,
	})
}

func (ph *ProjectHandlerImpl) GetProjectsHandler(c *gin.Context) {
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

	projects, err := ph.ProjectService.GetProjects(userUUID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.BadRequestError("Error getting projects.", err.Error()))
		return
	}

	projectList := dto.ProjectListResponse{
		Projects: make([]dto.ProjectResponse, len(projects)),
	}
	for p := range projects {
		project := projects[p]
		projectList.Projects[p] = dto.NewProjectResponse(project)
	}

	c.JSON(http.StatusAccepted, gin.H{
		"status": "success",
		"data":   projectList,
	})
}

func (ph *ProjectHandlerImpl) DeleteProjectHandler(c *gin.Context) {
	var req dto.ProjectRequest

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

	err = ph.ProjectService.DeleteProject(userUUID, req.ProjectId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BadRequestError("Failed to delete project.", err.Error()))
		return
	}

	c.Status(http.StatusAccepted)
}

func (ph *ProjectHandlerImpl) AddProjectHandler(c *gin.Context) {
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

	project, err := ph.ProjectService.AddProject(userUUID, req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BadRequestError("Failed to create project.", err.Error()))
		return
	}

	resp := dto.NewProjectResponse(project)

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   resp,
	})
}
