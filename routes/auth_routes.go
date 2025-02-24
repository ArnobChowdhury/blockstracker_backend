package routes

import (
	"blockstracker_backend/handlers"
	"blockstracker_backend/internal/database"
	"blockstracker_backend/internal/repositories"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(rg *gin.RouterGroup) {
	taskGroup := rg.Group("/auth")
	userRepo := repositories.NewUserRepository(database.DB)
	authHandler := handlers.NewAuthHandler(userRepo)

	{
		taskGroup.POST("/signup", authHandler.SignupUser)

	}
}
