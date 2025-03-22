package handlers

import (
	"net/http"

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

// CreateTask godoc
// @Summary Create a new task
// @Description Create a new task with the given details
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body models.CreateTaskRequest true "Task details"
// @Success 200 {object} models.CreateTaskResponseForSwagger
// @Failure 400 {object} models.GenericErrorResponse
// @Failure 500 {object} models.GenericErrorResponse
// @Router /task/create [post]
func (h *TaskHandler) CreateTask(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		utils.SendErrorResponse(c, h.logger, messages.ErrTaskCreationFailed,
			"User ID not found in context", apperrors.ErrInternalServerError)
		return
	}

	uid, ok := userID.(uuid.UUID)
	if !ok {
		utils.SendErrorResponse(c, h.logger, messages.ErrTaskCreationFailed,
			"User ID is not of valid type", apperrors.ErrInternalServerError)
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
		UserID:                   uid,
	}

	if err := h.taskRepo.CreateTask(&task); err != nil {
		utils.SendErrorResponse(c, h.logger, messages.ErrTaskCreationFailed,
			err.Error(), apperrors.ErrInternalServerError)
		return
	}

	c.JSON(http.StatusOK, utils.CreateJSONResponse(
		messages.Success, messages.MsgTaskCreationSuccess, task))
}
