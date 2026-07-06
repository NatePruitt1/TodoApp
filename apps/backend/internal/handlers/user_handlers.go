package handlers

import (
	"backend/internal/dto"
	"backend/internal/services"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserHandler interface allows DI for the user handler functionality.
type UserHandler interface {
	LoginHandler(c *gin.Context)
}

// UserHandlerImpl is the main UserHandler implementation.
// The UserHandler interface is implemented on (*UserHandlerImpl)
type UserHandlerImpl struct {
	UserService services.UserService
}

// Creates and returns a new (*UserHandlerImpl) for use as a UserHandler.
func NewUserHandler(service services.UserService) *UserHandlerImpl {
	return &UserHandlerImpl{
		UserService: service,
	}
}

// Handles the creation or rejection of a new account on the create endpoint.
func (uh *UserHandlerImpl) CreateAccountHandler(c *gin.Context) {
	var req dto.LoginRequestDTO
	var resp dto.LoginResponseDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequestError("Failed to create account, bad body data.", err.Error()))
		return
	}

	resp, token, err := uh.UserService.CreateAccount(context.Background(), req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequestError("Failed to create account.", err.Error()))
		return
	}

	c.Header("Authorization", "Bearer "+token)
	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   resp,
	})
}

// Handles login requests to allow users to authenticate.
func (uh *UserHandlerImpl) LoginHandler(c *gin.Context) {
	var req dto.LoginRequestDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequestError("Failed to login, bad body data.", err.Error()))
		return
	}

	resp, token, err := uh.UserService.AuthenticateAccount(context.Background(), req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequestError("Failed to create account.", err.Error()))
		return
	}

	c.Header("Authorization", "Bearer "+token)
	c.JSON(http.StatusAccepted, gin.H{
		"status": "success",
		"data":   resp,
	})
}
