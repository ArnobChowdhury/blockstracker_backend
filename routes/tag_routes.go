package routes

import (
	"blockstracker_backend/handlers"
	"blockstracker_backend/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterTagRoutes(rg *gin.RouterGroup, tagHandler *handlers.TagHandler, authMiddleware *middleware.AuthMiddleware) {

	tagGroup := rg.Group("/tags")
	tagGroup.Use(authMiddleware.Handle)

	{
		tagGroup.POST("/", tagHandler.CreateTag)
		tagGroup.PUT("/:id", tagHandler.UpdateTag)
		tagGroup.GET("/", tagHandler.GetTagsFromVersion)
	}
}
