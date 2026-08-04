package dto

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

// LoggerKey is the gin context key the logger is stored under by the logging middleware.
const LoggerKey = "logger"

// Error codes are namespaced "<CATEGORY>-<number>" so the frontend can key off
// a stable identifier instead of parsing the human-readable issue text.
const (
	// AUTH: authentication/authorization failures.
	ErrAuthMissingHeader  = "AUTH-100" // Authorization header not provided.
	ErrAuthMalformedToken = "AUTH-101" // Authorization header isn't "Bearer {token}".
	ErrAuthInvalidToken   = "AUTH-102" // Access token missing, invalid, or expired.
	ErrAuthInvalidUser    = "AUTH-103" // User identity couldn't be resolved from the request.
	ErrAuthForbidden      = "AUTH-104" // User does not own/have access to the requested resource.
	ErrAuthBadCredentials = "AUTH-105" // Login/account creation rejected by the auth service.
	ErrAuthRefreshToken   = "AUTH-106" // Refresh token missing or rejected.

	ErrAuthUsernameSpace    = "AUTH-000" // Username has a space.
	ErrAuthUsernameSC       = "AUTH-001" // Username contains a banned special character (or a special character first/last)
	ErrAuthUsernameTooShort = "AUTH-002" // Username too short
	ErrAuthUsernameTooLong  = "AUTH-003" // Username too long

	ErrAuthPasswordSpace       = "AUTH-200" // Password contains a space.
	ErrAuthPasswordNSC         = "AUTH-201" // Password does not contain a special character
	ErrAuthPasswordTooShort    = "AUTH-202" // Password is too short
	ErrAuthPasswordTooLong     = "AUTH-203" // Password is too long
	ErrAuthPasswordNoUppercase = "AUTH-204" // Password does not contain an uppercase character

	// VAL: request validation failures.
	ErrValInvalidID   = "VAL-100" // Route ID parameter missing or not a valid UUID.
	ErrValBadBody     = "VAL-101" // Request body missing/malformed JSON.
	ErrValCredentials = "VAL-102" // Username/password fail complexity requirements.

	// DB: failures surfaced from a service/repository call.
	ErrDBCreateFailed = "DB-100"
	ErrDBUpdateFailed = "DB-101"
	ErrDBDeleteFailed = "DB-102"
	ErrDBFetchFailed  = "DB-103"

	ErrCodeNotIMPL = "ERR-000" // Error code has not been implemented.
)

func BadRequestError(ctx *gin.Context, issue, details string) gin.H {
	return BadRequestErrorWithCode(ctx, ErrCodeNotIMPL, issue, details)
}

func BadRequestErrorWithCode(ctx *gin.Context, ec, issue, details string) gin.H {
	if logger, ok := ctx.Get(LoggerKey); ok {
		logger.(*slog.Logger).Error(issue, "details", details, "error_code", ec)
	}

	return gin.H{
		"code":  ec,
		"issue": issue,
	}
}
