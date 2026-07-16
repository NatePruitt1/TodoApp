package middleware

import (
	"backend/internal/config"
	"backend/internal/dto"
	"backend/internal/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, dto.BadRequestError("Empty Auth header", ""))
			ctx.Abort()
			return
		}

		authHeaderParts := strings.SplitN(authHeader, " ", 2)
		if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, dto.BadRequestError("Auth header must be of the form \"Bearer {token}\"", ""))
			ctx.Abort()
			return
		}

		token, err := services.ParseToken(authHeaderParts[1], cfg.JwtSecret)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, dto.BadRequestError("Invalid token.", err.Error()))
			ctx.Abort()
			return
		}

		ctx.Set("user_id", token.UserID)
		ctx.Set("username", token.Username)
		ctx.Next()
	}
}
