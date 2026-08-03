package dto

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

// LoggerKey is the gin context key the logger is stored under by the logging middleware.
const LoggerKey = "logger"

func BadRequestError(ctx *gin.Context, issue, details string) gin.H {
	if logger, ok := ctx.Get(LoggerKey); ok {
		logger.(*slog.Logger).Error(issue, "details", details)
	}

	return gin.H{
		"issue": issue,
	}
}
