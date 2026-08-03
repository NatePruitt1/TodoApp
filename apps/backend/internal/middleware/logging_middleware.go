package middleware

import (
	"backend/internal/dto"
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware() gin.HandlerFunc {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	return func(ctx *gin.Context) {
		start := time.Now()

		ctx.Set(dto.LoggerKey, logger)
		ctx.Next()

		logger.Info("request",
			"method", ctx.Request.Method,
			"path", ctx.Request.URL.Path,
			"status", ctx.Writer.Status(),
			"duration", time.Since(start).String(),
		)
	}
}