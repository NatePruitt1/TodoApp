package main

import (
	"backend/internal/config"
	"backend/internal/handlers"
	"backend/internal/repository"
	"backend/internal/services"
	"context"
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using system vars only.")
	}

	gin.SetMode(gin.TestMode)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	cfg, err := config.LoadConfigFromEnv()
	if err != nil {
		fmt.Println(err.Error())
		panic(-1)
	}

	fmt.Println(cfg.DatabaseURL)

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		fmt.Println(err.Error())
		panic(-1)
	}

	tokenRepo := repository.NewRefreshTokenRepository(pool)
	tokenService := services.NewRefreshTokenService(tokenRepo, cfg)

	userRepo := repository.NewUserRepository(pool)
	userService := services.NewUserService(userRepo, cfg)
	userHandler := handlers.NewUserHandler(userService, tokenService)

	api := router.Group("/api/v0")
	api.POST("/create", userHandler.CreateAccountHandler)
	api.POST("/login", userHandler.LoginHandler)
	api.POST("/refresh", userHandler.RefreshHandler)

	fmt.Println("Running server on :8080")
	router.Run(":8080")
}
