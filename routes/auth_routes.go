package routes

import (
	"blockstracker_backend/handlers"
	"blockstracker_backend/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(rg *gin.RouterGroup, authHandler *handlers.AuthHandler, authMiddleware *middleware.AuthMiddleware) {
	authGroup := rg.Group("/auth")

	{
		authGroup.POST("/signup", authHandler.SignupUser)
		authGroup.POST("/signin", authHandler.EmailSignIn)
		authGroup.POST("/refresh", authHandler.RefreshToken)
		authGroup.POST("/google/mobile", authHandler.GoogleSignInMobile)
		authGroup.POST("/google/desktop", authHandler.GoogleSignInDesktop)
		// probably has a problem since we are using auth middleware. What if the user is not authenticated?
		// not a problem for now, since both front ends will try to automatically refresh the token and then log out
		// but it could have been straightforward
		authGroup.Use(authMiddleware.Handle).POST("/signout", authHandler.Signout)
	}
}
