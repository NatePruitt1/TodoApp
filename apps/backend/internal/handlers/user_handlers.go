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

const PASSWORD_MIN_LEN = 6
const PASSWORD_MAX_LEN = 256
const USERNAME_MIN_LEN = 3
const USERNAME_MAX_LEN = 16

const USERNAME_ALLOWED_SPECIALS = ".-_"

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

func hasSpace(s string) bool {
	for _, i := range s {
		if unicode.IsSpace(i) {
			return true
		}
	}

	return false
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

// hasBannedUsernameChar reports whether s contains anything other than
// letters, digits, underscores, or hyphens.
func hasBannedUsernameChar(s string) bool {
	for _, i := range s {
		if !unicode.IsLetter(i) && !unicode.IsDigit(i) && !strings.ContainsRune(USERNAME_ALLOWED_SPECIALS, i) {
			return true
		}
	}

	return false
}

func validatePassword(password string, ctx *gin.Context) (bool, gin.H) {
	//Check at least one uppercase
	if hasSpace(password) {
		return false, dto.BadRequestErrorWithCode(ctx, dto.ErrAuthPasswordNSC, "Password may not contain spaces.", "")
	}

	if strippedLen(password) < PASSWORD_MIN_LEN {
		return false, dto.BadRequestErrorWithCode(ctx, dto.ErrAuthPasswordTooShort, fmt.Sprintf("Password must be atleast %d characters.", PASSWORD_MIN_LEN), "")
	}

	if strippedLen(password) >= PASSWORD_MAX_LEN {
		return false, dto.BadRequestErrorWithCode(ctx, dto.ErrAuthPasswordTooLong, fmt.Sprintf("Password must be shorter that %d characters.", PASSWORD_MAX_LEN), "")
	}

	if !hasUppercase(password) {
		return false, dto.BadRequestErrorWithCode(ctx, dto.ErrAuthPasswordNoUppercase, "Password must have atleast 1 uppercase character.", "")
	}

	if !hasSpecial(password) {
		return false, dto.BadRequestErrorWithCode(ctx, dto.ErrAuthPasswordNSC, "Password must have atleast 1 special character.", "")
	}

	return true, gin.H{}
}

func validateUsername(username string, ctx *gin.Context) (bool, gin.H) {
	if hasSpace(username) {
		return false, dto.BadRequestErrorWithCode(ctx, dto.ErrAuthUsernameSpace, "Username may not contain spaces.", "")
	}

	if strippedLen(username) < USERNAME_MIN_LEN {
		return false, dto.BadRequestErrorWithCode(ctx, dto.ErrAuthUsernameTooShort, fmt.Sprintf("Username must be atleast %d characters.", USERNAME_MIN_LEN), "")
	}

	if strippedLen(username) >= USERNAME_MAX_LEN {
		return false, dto.BadRequestErrorWithCode(ctx, dto.ErrAuthUsernameTooLong, fmt.Sprintf("Username must be shorter that %d characters.", USERNAME_MAX_LEN), "")
	}

	trimmed := []rune(strings.TrimSpace(username))
	firstOrLast := strings.ContainsRune(USERNAME_ALLOWED_SPECIALS, trimmed[0]) || strings.ContainsRune(USERNAME_ALLOWED_SPECIALS, trimmed[len(trimmed)-1])
	if hasBannedUsernameChar(string(trimmed)) || firstOrLast {
		return false, dto.BadRequestErrorWithCode(ctx, dto.ErrAuthUsernameSC, "Username contains a banned character, or a special character at the start/end.", "")
	}

	return true, gin.H{}
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

	validUsername, errResp := validateUsername(req.Username, c)
	if !validUsername {
		c.JSON(http.StatusBadRequest, errResp)
		uh.deleteRefreshToken(c)
		return
	}

	validPassword, errResp := validatePassword(req.Password, c)
	if !validPassword {
		c.JSON(http.StatusBadRequest, errResp)
		uh.deleteRefreshToken(c)
		return
	}

	resp, token, err := uh.UserService.CreateAccount(context.Background(), req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequestErrorWithCode(c, dto.ErrAuthBadCredentials, "Failed to create account.", err.Error()))
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

	validUsername, errResp := validateUsername(req.Username, c)
	if !validUsername {
		c.JSON(http.StatusBadRequest, errResp)
		uh.deleteRefreshToken(c)
		return
	}

	validPassword, errResp := validatePassword(req.Password, c)
	if !validPassword {
		c.JSON(http.StatusBadRequest, errResp)
		uh.deleteRefreshToken(c)
		return
	}

	resp, token, err := uh.UserService.AuthenticateAccount(context.Background(), req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.BadRequestErrorWithCode(c, dto.ErrAuthBadCredentials, "Failed to log in to account.", err.Error()))
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
