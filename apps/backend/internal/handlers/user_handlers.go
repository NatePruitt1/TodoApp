package handlers

import (
	"backend/internal/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoginHandler(c *gin.Context) {
	var req dto.LoginRequestDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequestError("Failed to login, bad body data.", err.Error()))
		return
	}

}
