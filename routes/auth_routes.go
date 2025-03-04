package routes

import (
	"blockstracker_backend/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(rg *gin.RouterGroup, authHandler *handlers.AuthHandler) {
	authGroup := rg.Group("/auth")

	{
		authGroup.POST("/signup", authHandler.SignupUser)
		authGroup.POST("/signin", authHandler.EmailSignIn)

	}
}
