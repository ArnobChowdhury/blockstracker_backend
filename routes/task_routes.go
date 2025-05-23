package routes

import (
	"blockstracker_backend/handlers"
	"blockstracker_backend/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterTaskRoutes(rg *gin.RouterGroup, taskHandler *handlers.TaskHandler, authMiddleware *middleware.AuthMiddleware) {
	taskGroup := rg.Group("/tasks")
	taskGroup.Use(authMiddleware.Handle)

	{
		taskGroup.POST("/", taskHandler.CreateTask)
		taskGroup.PUT("/:id", taskHandler.UpdateTask)

		taskGroup.POST("/repetitive", taskHandler.CreateRepetitiveTaskTemplate)
		taskGroup.PUT("/repetitive/:id", taskHandler.UpdateRepetitiveTaskTemplate)
	}
}
