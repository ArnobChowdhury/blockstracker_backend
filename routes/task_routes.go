package routes

import (
	"blockstracker_backend/handlers"
	"blockstracker_backend/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterTaskRoutes(rg *gin.RouterGroup, taskHandler *handlers.TaskHandler, authMiddleware *middleware.AuthMiddleware) {
	taskGroup := rg.Group("/task")

	{
		taskGroup.Use(authMiddleware.Handle)
		taskGroup.POST("/create", taskHandler.CreateTask)
	}
}
