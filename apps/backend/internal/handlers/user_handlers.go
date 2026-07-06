package handlers

import (
	"backend/internal/dto"
	"backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	LoginHandler(c *gin.Context)
}

type UserHandlerImpl struct {
	UserService *services.UserService
}

func (uh *UserHandlerImpl) LoginHandler(c *gin.Context) {
	var req dto.LoginRequestDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequestError("Failed to login, bad body data.", err.Error()))
		return
	}

}
