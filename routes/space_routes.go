package routes

import (
	"blockstracker_backend/handlers"
	"blockstracker_backend/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterSpaceRoutes(rg *gin.RouterGroup, spaceHandler *handlers.SpaceHandler, authMiddleware *middleware.AuthMiddleware) {

	spaceGroup := rg.Group("/spaces")
	spaceGroup.Use(authMiddleware.Handle)
	spaceGroup.Use(authMiddleware.RequirePremium)

	{
		spaceGroup.POST("/", spaceHandler.CreateSpace)
		spaceGroup.PUT("/:id", spaceHandler.UpdateSpace)
		spaceGroup.GET("/", spaceHandler.GetSpacesFromVersion)
	}
}
