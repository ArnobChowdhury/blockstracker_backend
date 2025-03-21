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
		taskGroup.PUT("/update", taskHandler.UpdateTask)
		taskGroup.PUT("/update-repetitive", taskHandler.UpdateRepetitiveTask)
		taskGroup.GET("/tasks-today", taskHandler.GetTasksToday)
		taskGroup.PUT("/toggle-completion-status", taskHandler.ToggleTaskCompletionStatus)
		taskGroup.GET("/daily-tasks-monthly-report", taskHandler.GetDailyTasksMonthlyReport)
		taskGroup.GET("/specific-days-in-a-week-tasks-monthly-report", taskHandler.GetSpecificDaysInAWeekTasksMonthlyReport)
		taskGroup.GET("/tasks-overdue", taskHandler.GetOverdueTasks)
		taskGroup.PUT("/task-failure/:id", taskHandler.MarkTaskAsFailure)
		taskGroup.PUT("/reschedule/:id", taskHandler.RescheduleTask)
		taskGroup.GET("/active/unscheduled", taskHandler.GetAllUnscheduledActiveTasks)                     // Unscheduling tasks
		taskGroup.GET("/active/once", taskHandler.GetAllOneOffActiveTasks)                                 // One-time tasks
		taskGroup.GET("/active/daily", taskHandler.GetAllDailyActiveTasks)                                 // Daily tasks
		taskGroup.GET("/active/specific-days-in-a-week", taskHandler.GetAllSpecificDaysInAWeekActiveTasks) // Tasks for specific days of the week
		taskGroup.PUT("/bulk-failure", taskHandler.BulkTaskFailure)
		taskGroup.GET("/details/:id", taskHandler.GetTaskDetails)
		taskGroup.GET("/repetitive-details/:id", taskHandler.GetRepetitiveTaskDetails)
		taskGroup.PUT("/stop-repetitive/:id", taskHandler.StopRepetitiveTask)
		taskGroup.GET("/tags", taskHandler.GetAllTags)
		taskGroup.GET("/spaces", taskHandler.GetAllSpaces)
		taskGroup.POST("/tags", taskHandler.CreateTag)
		taskGroup.POST("/spaces", taskHandler.CreateSpace)
		taskGroup.GET("/active/unscheduled/space/:spaceId", taskHandler.GetUnscheduledActiveTasksWithSpaceID)
		taskGroup.GET("/active/one-off/space/:spaceId", taskHandler.GetOneOffActiveTasksWithSpaceID)
		taskGroup.GET("/active/daily/space/:spaceId", taskHandler.GetDailyActiveTasksWithSpaceID)
		taskGroup.GET("/active/specific-days-in-a-week/space/:spaceId", taskHandler.GetSpecificDaysInAWeekActiveTasksWithSpaceID)
		taskGroup.GET("/active/unscheduled/without-space", taskHandler.GetUnscheduledActiveTasksWithoutSpace)
		taskGroup.GET("/active/one-off/without-space", taskHandler.GetOneOffActiveTasksWithoutSpace)
		taskGroup.GET("/active/daily/without-space", taskHandler.GetDailyActiveTasksWithoutSpace)
		taskGroup.GET("/active/specific-days-in-a-week/without-space", taskHandler.GetSpecificDaysInAWeekActiveTasksWithoutSpace)
	}
}
