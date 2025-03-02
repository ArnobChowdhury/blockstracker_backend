package routes

import (
	"blockstracker_backend/handlers"
	"blockstracker_backend/internal/database"
	"blockstracker_backend/internal/repositories"

	"blockstracker_backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(rg *gin.RouterGroup) {
	authGroup := rg.Group("/auth")
	userRepo := repositories.NewUserRepository(database.DB)
	authHandler := handlers.NewAuthHandler(userRepo, logger.Log)

	{
		authGroup.POST("/signup", authHandler.SignupUser)
		authGroup.POST("/signin", authHandler.EmailSignIn)

	}
}
