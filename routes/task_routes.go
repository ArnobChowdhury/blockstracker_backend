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
		taskGroup.GET("/tasks-today", handlers.GetTasksToday)
		taskGroup.PUT("/toggle-completion-status", handlers.ToggleTaskCompletionStatus)
		taskGroup.GET("/daily-tasks-monthly-report", handlers.GetDailyTasksMonthlyReport)
		taskGroup.GET("/specific-days-in-a-weeek-tasks-monthly-report", handlers.GetSpecificDaysInAWeekTasksMonthlyReport)
		taskGroup.GET("/tasks-overdue", handlers.GetOverdueTasks)
		taskGroup.PUT("/task-failure/:id", handlers.MarkTaskAsFailure)
		taskGroup.PUT("/reschedule/:id", handlers.RescheduleTask)
		taskGroup.GET("/active/unscheduled", handlers.GetAllUnscheduledActiveTasks)                     // Unscheduling tasks
		taskGroup.GET("/active/once", handlers.GetAllOneOffActiveTasks)                                 // One-time tasks
		taskGroup.GET("/active/daily", handlers.GetAllDailyActiveTasks)                                 // Daily tasks
		taskGroup.GET("/active/specific-days-in-a-week", handlers.GetAllSpecificDaysInAWeekActiveTasks) // Tasks for specific days of the week
		taskGroup.PUT("/bulk-failure", handlers.BulkTaskFailure)
		taskGroup.GET("/details/:id", handlers.GetTaskDetails)
		taskGroup.GET("/repetitive-details/:id", handlers.GetRepetitiveTaskDetails)
		taskGroup.PUT("/stop-repetitive/:id", handlers.StopRepetitiveTask)
		taskGroup.GET("/tags", handlers.GetAllTags)
		taskGroup.GET("/spaces", handlers.GetAllSpaces)
		taskGroup.POST("/tags", handlers.CreateTag)
		taskGroup.POST("/spaces", handlers.CreateSpace)
		taskGroup.GET("/active/unscheduled/space/:spaceId", handlers.GetUnscheduledActiveTasksWithSpaceID)
		taskGroup.GET("/active/one-off/space/:spaceId", handlers.GetOneOffActiveTasksWithSpaceID)
		taskGroup.GET("/active/daily/space/:spaceId", handlers.GetDailyActiveTasksWithSpaceID)
		taskGroup.GET("/active/specific-days-in-a-week/space/:spaceId", handlers.GetSpecificDaysInAWeekActiveTasksWithSpaceID)
		taskGroup.GET("/active/unscheduled/without-space", handlers.GetUnscheduledActiveTasksWithoutSpace)
		taskGroup.GET("/active/one-off/without-space", handlers.GetOneOffActiveTasksWithoutSpace)
		taskGroup.GET("/active/daily/without-space", handlers.GetDailyActiveTasksWithoutSpace)
		taskGroup.GET("/active/specific-days-in-a-week/without-space", handlers.GetSpecificDaysInAWeekActiveTasksWithoutSpace)

	}
}
