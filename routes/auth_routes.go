package routes

import (
	"blockstracker_backend/di"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(rg *gin.RouterGroup) error {
	authGroup := rg.Group("/auth")
	authHandler, err := di.InitializeAuthHandler()
	if err != nil {
		return err
	}
	authMiddleware, err := di.InitializeAuthMiddleware()
	if err != nil {
		return err
	}

	{
		authGroup.POST("/signup", authHandler.SignupUser)
		authGroup.POST("/signin", authHandler.EmailSignIn)
		authGroup.Use(authMiddleware.Handle).POST("/signout", authHandler.Signout)
	}
	return nil
}
