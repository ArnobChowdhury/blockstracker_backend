package routes

import (
	"blockstracker_backend/di"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(rg *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	authGroup := rg.Group("/auth")
	authHandler := di.InitializeAuthHandler()

	{
		authGroup.POST("/signup", authHandler.SignupUser)
		authGroup.POST("/signin", authHandler.EmailSignIn)
		authGroup.Use(authMiddleware).POST("/signout", authHandler.Signout)
	}
}
