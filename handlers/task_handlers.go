package handlers

import (
	"net/http"
	"time"

	apperrors "blockstracker_backend/internal/errors"
	"blockstracker_backend/internal/repositories"
	"blockstracker_backend/internal/utils"
	"blockstracker_backend/messages"
	"blockstracker_backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type TaskHandler struct {
	taskRepo *repositories.TaskRepository
	logger   *zap.SugaredLogger
}

func NewTaskHandler(
	taskRepo *repositories.TaskRepository,
	logger *zap.SugaredLogger,
) *TaskHandler {
	return &TaskHandler{
		taskRepo: taskRepo,
		logger:   logger,
	}
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		utils.SendErrorResponse(c, h.logger, messages.ErrTaskCreationFailed,
			"User ID not found in context", apperrors.ErrInternalServerError)
		return
	}

	var req models.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		invalidReqErr := apperrors.NewInvalidReqErr(err.Error())
		utils.SendErrorResponse(c, h.logger, messages.ErrTaskCreationFailed,
			err.Error(), invalidReqErr)
		return
	}

	task := models.Task{
		IsActive:                 req.IsActive,
		Title:                    req.Title,
		Description:              req.Description,
		Schedule:                 req.Schedule,
		Priority:                 req.Priority,
		CompletionStatus:         req.CompletionStatus,
		DueDate:                  req.DueDate,
		ShouldBeScored:           req.ShouldBeScored,
		Score:                    req.Score,
		TimeOfDay:                req.TimeOfDay,
		RepetitiveTaskTemplateID: req.RepetitiveTaskTemplateID,
		CreatedAt:                req.CreatedAt,
		ModifiedAt:               req.ModifiedAt,
		Tags:                     req.Tags,
		SpaceID:                  req.SpaceID,
		UserID:                   userID.(uuid.UUID),
	}

	if err := h.taskRepo.CreateTask(&task); err != nil {
		utils.SendErrorResponse(c, h.logger, messages.ErrTaskCreationFailed,
			err.Error(), apperrors.ErrInternalServerError)
		return
	}

	c.JSON(http.StatusOK, utils.CreateJSONResponse(
		messages.MsgTaskCreationSuccess, messages.Success, task))
}

func (h *TaskHandler) CreateRepetitiveTaskTemplate(c *gin.Context) {
	var req models.RepetitiveTaskTemplate
	if err := c.ShouldBindJSON(&req); err != nil {

		utils.SendErrorResponse(c, h.logger, "task creation failed", err.Error(), apperrors.ErrMalformedRequest)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": messages.MsgTaskCreationSuccess})
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	var req models.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": messages.MsgTaskUpdatedSuccess})
}

func (h *TaskHandler) UpdateRepetitiveTask(c *gin.Context) {
	var req models.UpdateRepetitiveTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": messages.MsgTaskUpdatedSuccess})
}

func (h *TaskHandler) GetTasksToday(c *gin.Context) {
	// Mock data, assuming you have a Task model and a service that fetches real tasks
	tasks := []map[string]interface{}{
		{
			"id":         1,
			"title":      "Sample Task 1",
			"due_date":   time.Now().Format("2006-01-02"),
			"completion": false,
		},
		{
			"id":         2,
			"title":      "Sample Task 2",
			"due_date":   time.Now().Format("2006-01-02"),
			"completion": true,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Tasks for today fetched successfully",
		"tasks":   tasks,
	})
}

func (h *TaskHandler) ToggleTaskCompletionStatus(c *gin.Context) {
	var req struct {
		TaskID    int  `json:"task_id" binding:"required"`
		Completed bool `json:"completed" binding:"required"`
	}

	// Parse JSON body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Here, you would toggle the completion status in the database
	// This is just a mock response
	c.JSON(http.StatusOK, gin.H{
		"message":   "Task completion status updated successfully",
		"task_id":   req.TaskID,
		"completed": req.Completed,
	})
}

func (h *TaskHandler) GetDailyTasksMonthlyReport(c *gin.Context) {
	// Mock data, assuming you are pulling actual report data from a service
	report := []map[string]interface{}{
		{
			"date":       time.Now().Format("2006-01-02"),
			"completed":  5, // tasks completed on this day
			"incomplete": 2, // tasks incomplete on this day
		},
		{
			"date":       time.Now().AddDate(0, 0, -1).Format("2006-01-02"),
			"completed":  3,
			"incomplete": 4,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Daily tasks report for the month fetched successfully",
		"report":  report,
	})
}

func (h *TaskHandler) GetSpecificDaysInAWeekTasksMonthlyReport(c *gin.Context) {
	// Mock data, assuming you are pulling actual report data from a service
	report := []map[string]interface{}{
		{
			"date":       time.Now().Format("2006-01-02"),
			"completed":  5, // tasks completed on this day
			"incomplete": 2, // tasks incomplete on this day
		},
		{
			"date":       time.Now().AddDate(0, 0, -1).Format("2006-01-02"),
			"completed":  3,
			"incomplete": 4,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Specific days in a week tasks report for the month fetched successfully",
		"report":  report,
	})
}

func (h *TaskHandler) GetOverdueTasks(c *gin.Context) {
	// Get today's date to compare with due dates
	// today := time.Now().Format("2006-01-02")

	// Sample overdue tasks (mock data)
	// overdueTasks := []models.Task{
	// 	{
	// 		// ID:        1,
	// 		Title: "Overdue Task 1",
	// 		// DueDate:   "2025-02-10", // Past due date
	// 		// Completed: false,
	// 	},
	// 	{
	// 		// ID:        2,
	// 		Title: "Overdue Task 2",
	// 		// DueDate:   "2025-02-08", // Past due date
	// 		// Completed: false,
	// 	},
	// }

	// Filter overdue tasks (tasks that are not completed and have a due date in the past)
	var filteredTasks []models.Task
	// for _, task := range overdueTasks {
	// 	if task.DueDate < today && !task.Completed {
	// 		filteredTasks = append(filteredTasks, task)
	// 	}
	// }

	// Return the overdue tasks
	c.JSON(http.StatusOK, gin.H{
		"message": "Overdue tasks fetched successfully",
		"tasks":   filteredTasks,
	})
}

func (h *TaskHandler) MarkTaskAsFailure(c *gin.Context) {
	// Fetch the task ID from the URL parameter
	// taskID := c.Param("id")

	// [mock] make a call to the db with the taskID to find the original task
	// task := models.Task{
	// 	ID:        1,
	// 	Title:     "Overdue Task 1",
	// 	DueDate:   "2025-02-10", // Past due date
	// 	Completed: false,
	// }

	// Respond with the updated task
	// c.JSON(http.StatusOK, gin.H{
	// 	"message": "Task marked as failure successfully",
	// 	"task":    task,
	// })
}

func (h *TaskHandler) RescheduleTask(c *gin.Context) {
	// Extract task ID and new due date from the request
	var req struct {
		TaskID     int    `json:"task_id" binding:"required"`
		NewDueDate string `json:"new_due_date" binding:"required"`
	}

	// Bind the JSON request body to the struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Respond with success message
	c.JSON(http.StatusOK, gin.H{"message": "Task rescheduled successfully"})
}

func (h *TaskHandler) GetAllUnscheduledActiveTasks(c *gin.Context) {
	// Helper: Implement logic to fetch unscheduled active tasks from the database or service
}

func (h *TaskHandler) GetAllOneOffActiveTasks(c *gin.Context) {
	// Helper: Implement logic to fetch one-off active tasks from the database or service
}

func (h *TaskHandler) GetAllDailyActiveTasks(c *gin.Context) {
	// Helper: Implement logic to fetch daily active tasks from the database or service
}

func (h *TaskHandler) GetAllSpecificDaysInAWeekActiveTasks(c *gin.Context) {
	// Helper: Implement logic to fetch active tasks for specific days in a week from the database or service
}

func (h *TaskHandler) BulkTaskFailure(c *gin.Context) {
	var taskIDs []int

	// Bind the request body to the taskIDs slice
	if err := c.ShouldBindJSON(&taskIDs); err != nil {
		// If there is an error with the JSON, return a bad request response
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Here, you would process the task IDs, like marking them as failed in the database
	// For now, we just return the received IDs for demonstration
	c.JSON(http.StatusOK, gin.H{
		"message":         "Tasks marked as failed successfully",
		"failed_task_ids": taskIDs,
	})
}

// GetTaskDetails retrieves details of a task by ID
func (h *TaskHandler) GetTaskDetails(c *gin.Context) {
	// c.Param("id") - Get task ID from the URL
	// Implement logic here
	c.JSON(http.StatusOK, gin.H{"message": "Task details fetched successfully"})
}

// GetRepetitiveTaskDetails retrieves details of a repetitive task by ID
func (h *TaskHandler) GetRepetitiveTaskDetails(c *gin.Context) {
	// c.Param("id") - Get repetitive task ID from the URL
	// Implement logic here
	c.JSON(http.StatusOK, gin.H{"message": "Repetitive task details fetched successfully"})
}

// StopRepetitiveTask marks a repetitive task as inactive
func (h *TaskHandler) StopRepetitiveTask(c *gin.Context) {
	// c.Param("id") - Get task ID from the URL
	// Implement logic here
	c.JSON(http.StatusOK, gin.H{"message": "Repetitive task stopped successfully"})
}

// GetAllTags retrieves all tags
func (h *TaskHandler) GetAllTags(c *gin.Context) {
	// Implement logic here
	c.JSON(http.StatusOK, gin.H{"message": "All tags fetched successfully"})
}

// GetAllSpaces retrieves all spaces
func (h *TaskHandler) GetAllSpaces(c *gin.Context) {
	// Implement logic here
	c.JSON(http.StatusOK, gin.H{"message": "All spaces fetched successfully"})
}

// CreateTag creates a new tag
func (h *TaskHandler) CreateTag(c *gin.Context) {
	// Implement logic here
	c.JSON(http.StatusOK, gin.H{"message": "Tag created successfully"})
}

// CreateSpace creates a new space
func (h *TaskHandler) CreateSpace(c *gin.Context) {
	// Implement logic here
	c.JSON(http.StatusOK, gin.H{"message": "Space created successfully"})
}

// GetUnscheduledActiveTasksWithSpaceID retrieves unscheduled active tasks by space ID
func (h *TaskHandler) GetUnscheduledActiveTasksWithSpaceID(c *gin.Context) {
	// c.Param("spaceId") - Get space ID from the URL
	// Implement logic here
	c.JSON(http.StatusOK, gin.H{"message": "Unscheduled active tasks for space fetched successfully"})
}

// GetOneOffActiveTasksWithSpaceID retrieves one-off active tasks by space ID
func (h *TaskHandler) GetOneOffActiveTasksWithSpaceID(c *gin.Context) {
	// c.Param("spaceId") - Get space ID from the URL
	// Implement logic here
	c.JSON(http.StatusOK, gin.H{"message": "One-off active tasks for space fetched successfully"})
}

// GetDailyActiveTasksWithSpaceID retrieves daily active tasks by space ID
func (h *TaskHandler) GetDailyActiveTasksWithSpaceID(c *gin.Context) {
	// c.Param("spaceId") - Get space ID from the URL
	// Implement logic here
	c.JSON(http.StatusOK, gin.H{"message": "Daily active tasks for space fetched successfully"})
}

// GetSpecificDaysInAWeekActiveTasksWithSpaceID retrieves tasks based on specific days in a week and space ID
func (h *TaskHandler) GetSpecificDaysInAWeekActiveTasksWithSpaceID(c *gin.Context) {
	// c.Param("spaceId") - Get space ID from the URL
	// Implement logic here
	c.JSON(http.StatusOK, gin.H{"message": "Tasks for specific days in a week for space fetched successfully"})
}

// GetUnscheduledActiveTasksWithoutSpace retrieves unscheduled active tasks without a space
func (h *TaskHandler) GetUnscheduledActiveTasksWithoutSpace(c *gin.Context) {
	// Implement logic here
	c.JSON(http.StatusOK, gin.H{"message": "Unscheduled active tasks without space fetched successfully"})
}

// GetOneOffActiveTasksWithoutSpace retrieves one-off active tasks without a space
func (h *TaskHandler) GetOneOffActiveTasksWithoutSpace(c *gin.Context) {
	// Implement logic here
	c.JSON(http.StatusOK, gin.H{"message": "One-off active tasks without space fetched successfully"})
}

// GetDailyActiveTasksWithoutSpace retrieves daily active tasks without a space
func (h *TaskHandler) GetDailyActiveTasksWithoutSpace(c *gin.Context) {
	// Implement logic here
	c.JSON(http.StatusOK, gin.H{"message": "Daily active tasks without space fetched successfully"})
}

// GetSpecificDaysInAWeekActiveTasksWithSpaceID retrieves tasks based on specific days in a week and space ID
func (h *TaskHandler) GetSpecificDaysInAWeekActiveTasksWithoutSpace(c *gin.Context) {
	// c.Param("spaceId") - Get space ID from the URL
	// Implement logic here
	c.JSON(http.StatusOK, gin.H{"message": "Tasks for specific days in a week for space fetched successfully"})
}
