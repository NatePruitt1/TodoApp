package main

import (
	"backend/internal/config"
	"backend/internal/handlers"
	"backend/internal/repository"
	"backend/internal/services"
	"context"
	"fmt"

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

	userRepo := repository.NewUserRepository(pool)
	userService := services.NewUserService(userRepo, cfg)
	userHandler := handlers.NewUserHandler(userService)

	api := router.Group("/api/v0")
	api.POST("/create", userHandler.CreateAccountHandler)
	api.POST("/login", userHandler.LoginHandler)

	fmt.Println("Running server on :8080")
	router.Run(":8080")
}
