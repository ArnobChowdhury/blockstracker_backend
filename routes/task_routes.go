package routes

import (
	"blockstracker_backend/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterTaskRoutes(r *gin.Engine) {
	taskGroup := r.Group("/task")
	{
		taskGroup.POST("/create", handlers.CreateTask)
		taskGroup.PUT("/update", handlers.UpdateTask)
		taskGroup.PUT("/update-repetitive", handlers.UpdateRepetitiveTask)
	}
}
