package handlers

import (
	"backend/internal/dto"
	"backend/internal/services"
	"context"
	"fmt"
	"net/http"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
)

const RTOKEN_PATH = "/api/v0/auth/refresh"

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

func (uh *UserHandlerImpl) deleteRefreshToken(c *gin.Context) {
	c.SetCookie("refresh_token", "", -1, RTOKEN_PATH, "", false, true)
}

func (uh *UserHandlerImpl) createRefreshToken(c *gin.Context, token string) {

	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("refresh_token",
		token,
		7*24*60*60,
		RTOKEN_PATH,
		"",
		false,
		true,
	)
}

func strippedLen(value string) int {
	return len(strings.TrimSpace(value))
}

func validateCredentialField(fieldName, value string, minLength int) string {
	if strippedLen(value) < minLength {
		return fmt.Sprintf("%s too short", fieldName)
	}

	if strings.IndexFunc(value, unicode.IsSpace) >= 0 {
		return fmt.Sprintf("%s cannot contain spaces", fieldName)
	}

	return ""
}

func hasUppercase(s string) bool {
	for _, i := range s {
		if unicode.IsUpper(i) {
			return true
		}
	}

	return false
}

func hasSpecial(s string) bool {
	for _, i := range s {
		if unicode.IsPunct(i) || unicode.IsSymbol(i) {
			return true
		}
	}

	return false
}

func validatePassword(password string) string {
	//Check at least one uppercase
	if !hasUppercase(password) {
		return "password requires atleast 1 uppercase."
	}

	if !hasSpecial(password) {
		return "password requires atleast 1 special character."
	}

	return ""
}

func validateUsernameAndPassword(username, password string) string {
	if message := validateCredentialField("Username", username, 3); message != "" {
		return message
	}

	if message := validateCredentialField("Password", password, 6); message != "" {
		return message
	}

	return ""
}

// Handles the creation or rejection of a new account on the create endpoint.
func (uh *UserHandlerImpl) CreateAccountHandler(c *gin.Context) {
	var req dto.LoginRequestDTO
	var resp dto.LoginResponseDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequestError(c, "Failed to create account, bad body data.", err.Error()))
		uh.deleteRefreshToken(c)
		return
	}

	if message := validateUsernameAndPassword(req.Username, req.Password); message != "" {
		c.JSON(http.StatusBadRequest, dto.BadRequestError(c, message, ""))
		uh.deleteRefreshToken(c)
		return
	}

	resp, token, err := uh.UserService.CreateAccount(context.Background(), req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequestError(c, "Failed to create account.", err.Error()))
		uh.deleteRefreshToken(c)
		return
	}

	refreshToken, err := uh.RefreshTokenService.CheckTokenForUser(resp.ID)
	if err != nil {
		refreshToken, err = uh.RefreshTokenService.CreateToken(resp.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.BadRequestError(c, "Failed to create refresh token.", err.Error()))
			uh.deleteRefreshToken(c)
			return
		}
	}

	uh.createRefreshToken(c, refreshToken.Raw)

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
		c.JSON(http.StatusBadRequest, dto.BadRequestError(c, "Failed to login, bad body data.", err.Error()))
		uh.deleteRefreshToken(c)
		return
	}

	if message := validateUsernameAndPassword(req.Username, req.Password); message != "" {
		c.JSON(http.StatusBadRequest, dto.BadRequestError(c, message, ""))
		uh.deleteRefreshToken(c)
		return
	}

	resp, token, err := uh.UserService.AuthenticateAccount(context.Background(), req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequestError(c, "Failed to create account.", err.Error()))
		uh.deleteRefreshToken(c)
		return
	}

	refreshToken, err := uh.RefreshTokenService.CheckTokenForUser(resp.ID)
	if err != nil {
		refreshToken, err = uh.RefreshTokenService.CreateToken(resp.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.BadRequestError(c, "Failed to create refresh token.", err.Error()))
			uh.deleteRefreshToken(c)
			return
		}
	} else {
		refreshToken, err = uh.RefreshTokenService.UpdateTokenByHash(refreshToken.Hash)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.BadRequestError(c, "Failed to create refresh token.", err.Error()))
			uh.deleteRefreshToken(c)
			return
		}
	}

	uh.createRefreshToken(c, refreshToken.Raw)

	c.Header("Authorization", "Bearer "+token)
	c.JSON(http.StatusAccepted, gin.H{
		"status": "success",
		"data":   resp,
	})
}

func (uh *UserHandlerImpl) RefreshHandler(c *gin.Context) {
	raw, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequestError(c, "No refresh token provided.", ""))
		return
	}

	fmt.Printf("Raw refresh token: %s\n", raw)

	_, err = uh.RefreshTokenService.CheckToken(raw)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.BadRequestError(c, "Refresh token not accepted.", err.Error()))
		uh.deleteRefreshToken(c)
		return
	}

	token, err := uh.RefreshTokenService.UpdateToken(raw)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.BadRequestError(c, "Failed to update refresh token.", err.Error()))
		return
	}

	resp, authToken, err := uh.UserService.RefreshAccount(context.Background(), token.UserId)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequestError(c, "Failed to log into account using refresh.", err.Error()))
		uh.deleteRefreshToken(c)
		return
	}

	uh.createRefreshToken(c, token.Raw)

	c.Header("Authorization", "Bearer "+authToken)
	c.JSON(http.StatusAccepted, gin.H{
		"status": "success",
		"data":   resp,
	})

}
