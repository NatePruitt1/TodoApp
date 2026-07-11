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
	UserService         services.UserService
	RefreshTokenService services.RefreshTokenService
}

// Creates and returns a new (*UserHandlerImpl) for use as a UserHandler.
func NewUserHandler(service services.UserService, refreshTokService services.RefreshTokenService) *UserHandlerImpl {
	return &UserHandlerImpl{
		RefreshTokenService: refreshTokService,
		UserService:         service,
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

	refreshToken, err := uh.RefreshTokenService.CreateToken(resp.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BadRequestError("Failed to create refresh token.", err.Error()))
		return
	}

	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("refresh_token",
		refreshToken.Raw,
		7*24*60*60,
		"/api/v0/refresh",
		"",
		true,
		true,
	)

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

	refreshToken, err := uh.RefreshTokenService.CreateToken(resp.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BadRequestError("Failed to create refresh token.", err.Error()))
		return
	}

	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("refresh_token",
		refreshToken.Raw,
		7*24*60*60,
		"/api/v0/refresh",
		"",
		true,
		true,
	)

	c.Header("Authorization", "Bearer "+token)
	c.JSON(http.StatusAccepted, gin.H{
		"status": "success",
		"data":   resp,
	})
}

func (uh *UserHandlerImpl) RefreshHandler(c *gin.Context) {
	raw, ok := c.Get("refresh_token")
	if !ok {
		c.JSON(http.StatusBadRequest, dto.BadRequestError("No refresh token provided.", ""))
		return
	}

	tokenId, ok := raw.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, dto.BadRequestError("Invalid refresh token provided.", "Could not parse refresh token string."))
		return
	}

	_, err := uh.RefreshTokenService.CheckToken(tokenId)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.BadRequestError("Refresh token not accepted.", err.Error()))
		return
	}

	token, err := uh.RefreshTokenService.UpdateToken(tokenId)

	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("refresh_token",
		token.Raw,
		7*24*60*60,
		"/api/v0/refresh",
		"",
		true,
		true,
	)

	resp, authToken, err := uh.UserService.RefreshAccount(context.Background(), token.UserId)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequestError("Failed to log into account using refresh.", err.Error()))
		return
	}

	c.Header("Authorization", "Bearer "+authToken)
	c.JSON(http.StatusAccepted, gin.H{
		"status": "success",
		"data":   resp,
	})

}
