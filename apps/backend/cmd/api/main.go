package main

import (
	"backend/internal/config"
	"backend/internal/handlers"
	"backend/internal/middleware"
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
	if err := godotenv.Load(".env", "internal/db/.env.db"); err != nil {
		fmt.Println("No .env file found, using system vars only.")
	}

	gin.SetMode(gin.TestMode)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://nathanielpruitt.com", "http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
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

	projectRepo := repository.NewProjectRepository(pool)
	projectService := services.NewProjectService(projectRepo)
	projectHandler := handlers.NewProjectHandler(projectService)

	api := router.Group("/api/v0", middleware.LoggingMiddleware())

	// Auth Endpoints - Public
	api.POST("/auth/register", userHandler.CreateAccountHandler)
	api.POST("/auth/login", userHandler.LoginHandler)
	api.POST("/auth/refresh", userHandler.RefreshHandler)

	// Auth Endpoints - Authenticated
	authGroup := api.Group("", middleware.AuthMiddleware(cfg))

	// Project Endpoints
	authGroup.GET("/projects", projectHandler.GetProjects)
	authGroup.POST("/projects", projectHandler.AddProject)
	authGroup.GET("/projects/:projectid", projectHandler.GetProject)
	authGroup.PATCH("/projects/:projectid", projectHandler.UpdateProject)
	authGroup.DELETE("/projects/:projectid", projectHandler.DeleteProject)

	// Category Endpoints
	authGroup.POST("/projects/:projectid/categories", projectHandler.AddCategory)
	authGroup.PATCH("/categories/:categoryid", projectHandler.UpdateCategory)
	authGroup.DELETE("/categories/:categoryid", projectHandler.DeleteCategory)
	authGroup.PATCH("/categories/:categoryid/position", projectHandler.MoveCategory)

	// Card Endpoints
	authGroup.POST("/categories/:categoryid/cards", projectHandler.AddCard)
	authGroup.PATCH("/cards/:cardid", projectHandler.RenameCard)
	authGroup.DELETE("/cards/:cardid", projectHandler.DeleteCard)
	authGroup.PATCH("/cards/:cardid/move", projectHandler.MoveCard)
	authGroup.PATCH("/cards/:cardid/finish", projectHandler.FinishCard)

	fmt.Println("Running server on :8080")
	router.Run(":8080")
}
