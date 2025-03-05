package routes

import (
	"blockstracker_backend/handlers"
	"blockstracker_backend/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(rg *gin.RouterGroup, authHandler *handlers.AuthHandler, authMiddleware *middleware.AuthMiddleware) error {
	authGroup := rg.Group("/auth")

	{
		authGroup.POST("/signup", authHandler.SignupUser)
		authGroup.POST("/signin", authHandler.EmailSignIn)
		authGroup.Use(authMiddleware.Handle).POST("/signout", authHandler.Signout)
	}
	return nil
}
