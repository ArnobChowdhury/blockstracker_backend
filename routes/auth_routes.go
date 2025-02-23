package routes

import (
	"blockstracker_backend/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(rg *gin.RouterGroup) {
	taskGroup := rg.Group("/auth")
	{
		taskGroup.POST("/signup", handlers.SignupUser)

	}
}
