package handlers

import (
	"backend/internal/dto"
	"backend/internal/services"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProjectHandler interface {
	GetProjects(c *gin.Context)
	GetProject(c *gin.Context)
	DeleteProject(c *gin.Context)
	AddProject(c *gin.Context)
	UpdateProject(c *gin.Context)

	// Category Handlers
	AddCategory(c *gin.Context)
	UpdateCategory(c *gin.Context)
	DeleteCategory(c *gin.Context)
	MoveCategory(c *gin.Context)

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

func getProjectId(c *gin.Context) (uuid.UUID, error) {
	projectId := c.Param("projectid")
	if projectId == "" {
		return uuid.Nil, errors.New("no project ID parameter provided.")
	}

	projectUUID, err := uuid.Parse(projectId)
	if err != nil {
		return uuid.Nil, errors.New("project ID provided invalid.")
	}

	return projectUUID, nil
}

func getCategoryId(c *gin.Context) (uuid.UUID, error) {
	categoryId := c.Param("categoryid")
	if categoryId == "" {
		return uuid.Nil, errors.New("no category ID parameter provided.")
	}

	categoryUUID, err := uuid.Parse(categoryId)
	if err != nil {
		return uuid.Nil, errors.New("category ID provided invalid.")
	}

	return categoryUUID, nil
}

func getUserId(c *gin.Context) (uuid.UUID, error) {
	userId := c.GetString("user_id")
	if userId == "" {
		return uuid.Nil, errors.New("no user ID provided.")
	}

	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return uuid.Nil, errors.New("invalid user ID provided.")
	}

	return userUUID, nil
}

func (ph *ProjectHandlerImpl) GetProject(c *gin.Context) {
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

	project, err := ph.ProjectService.GetProject(userUUID, projectUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequestError("Failed to get project.", err.Error()))
		return
	}

	projectResp := dto.NewProjectResponse(project)

	c.JSON(http.StatusAccepted, gin.H{
		"status": "success",
		"data":   projectResp,
	})
}

func (ph *ProjectHandlerImpl) GetProjects(c *gin.Context) {
	userUUID, err := getUserId(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.BadRequestError("failed to retrieve user id.", err.Error()))
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
		fmt.Printf("Description of project: %s\n", project.Description)

		projectList.Projects[p] = dto.NewProjectResponse(project)
	}

	c.JSON(http.StatusAccepted, gin.H{
		"status": "success",
		"data":   projectList,
	})
}

func (ph *ProjectHandlerImpl) DeleteProject(c *gin.Context) {
	userUUID, err := getUserId(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.BadRequestError("failed to retrieve user id.", err.Error()))
		return
	}

	projectUUID, err := getProjectId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequestError("Error retrieving project id.", err.Error()))
	}

	err = ph.ProjectService.DeleteProject(userUUID, projectUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BadRequestError("Failed to delete project.", err.Error()))
		return
	}

	c.Status(http.StatusNoContent)
}

func (ph *ProjectHandlerImpl) AddProject(c *gin.Context) {
	var req dto.ProjectRequest

	userUUID, err := getUserId(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.BadRequestError("failed to retrieve user id.", err.Error()))
		return
	}

	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, dto.BadRequestError("failed to parse request.", err.Error()))
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

func (ph *ProjectHandlerImpl) UpdateProject(c *gin.Context) {

}
