package handlers

import (
	"errors"
	"fmt"
	"net/http"

	apperrors "blockstracker_backend/internal/errors"
	"blockstracker_backend/internal/repositories"
	"blockstracker_backend/internal/utils"
	"blockstracker_backend/messages"
	"blockstracker_backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
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
// @Param task body models.TaskRequest true "Task details"
// @Success 200 {object} models.TaskResponseForSwagger
// @Failure 400 {object} models.GenericErrorResponse
// @Failure 500 {object} models.GenericErrorResponse
// @Router /tasks [post]
func (h *TaskHandler) CreateTask(c *gin.Context) {
	uid, err := utils.ExtractUIDFromGinContext(c)
	if err != nil {
		utils.SendErrorResponse(c, h.logger, messages.ErrTaskCreationFailed, err.LogError(),
			apperrors.ErrInternalServerError)
		return
	}

	var req models.TaskRequest
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

// UpdateTask godoc
// @Summary Update an existing task
// @Description Update an existing task with the given details
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Param task body models.TaskRequest true "Task details"
// @Success 200 {object} models.TaskResponseForSwagger
// @Failure 400 {object} models.GenericErrorResponse
// @Failure 404 {object} models.GenericErrorResponse
// @Failure 500 {object} models.GenericErrorResponse
// @Router /tasks/{id} [put]
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	uid, err := utils.ExtractUIDFromGinContext(c)
	if err != nil {
		utils.SendErrorResponse(c, h.logger, messages.ErrTaskUpdateFailed, err.LogError(),
			apperrors.ErrInternalServerError)
		return
	}

	taskIDStr := c.Param("id")
	taskID, taskIdParseErr := uuid.Parse(taskIDStr)
	if taskIdParseErr != nil {
		utils.SendErrorResponse(c, h.logger, messages.ErrTaskUpdateFailed,
			fmt.Sprintf("Invalid task ID format: %s", taskIDStr), apperrors.ErrMalformedTaskRequest)
		return
	}

	var req models.TaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		invalidReqErr := apperrors.NewInvalidReqErr(err.Error())
		utils.SendErrorResponse(c, h.logger, messages.ErrTaskUpdateFailed,
			err.Error(), invalidReqErr)
		return
	}

	task := models.Task{
		ID:                       taskID,
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
		ModifiedAt:               req.ModifiedAt,
		Tags:                     req.Tags,
		SpaceID:                  req.SpaceID,
		UserID:                   uid,
	}

	if err := h.taskRepo.UpdateTask(&task); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SendErrorResponse(c, h.logger, messages.ErrTaskUpdateFailed,
				"Task not found or does not belong to user", apperrors.ErrUnauthorized)
			return
		}

		utils.SendErrorResponse(c, h.logger, messages.ErrTaskUpdateFailed,
			err.Error(), apperrors.ErrInternalServerError)
		return
	}

	c.JSON(http.StatusOK, utils.CreateJSONResponse(
		messages.Success, messages.MsgTaskUpdateSuccess, task))
}

// CreateRepetitiveTaskTemplate godoc
// @Summary Create a new repetitive task template
// @Description Create a new repetitive task template with the given details
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body models.CreateRepetitiveTaskTemplateRequest true "Repetitive task template details"
// @Success 200 {object} models.CreateRepetitiveTaskTemplateResponseForSwagger
// @Failure 400 {object} models.GenericErrorResponse
// @Failure 500 {object} models.GenericErrorResponse
// @Router /tasks/repetitive [post]
func (h *TaskHandler) CreateRepetitiveTaskTemplate(c *gin.Context) {
	uid, err := utils.ExtractUIDFromGinContext(c)
	if err != nil {
		utils.SendErrorResponse(c, h.logger, messages.ErrRepetitiveTaskTemplateCreationFailed,
			err.LogError(), apperrors.ErrInternalServerError)
		return
	}

	var req models.CreateRepetitiveTaskTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		invalidReqErr := apperrors.NewInvalidReqErr(err.Error())
		utils.SendErrorResponse(c, h.logger, messages.ErrRepetitiveTaskTemplateCreationFailed,
			err.Error(), invalidReqErr)
		return
	}

	repetitiveTaskTemplate := models.RepetitiveTaskTemplate{
		IsActive:                 req.IsActive,
		Title:                    req.Title,
		Description:              req.Description,
		Schedule:                 req.Schedule,
		Priority:                 req.Priority,
		ShouldBeScored:           req.ShouldBeScored,
		Monday:                   req.Monday,
		Tuesday:                  req.Tuesday,
		Wednesday:                req.Wednesday,
		Thursday:                 req.Thursday,
		Friday:                   req.Friday,
		Saturday:                 req.Saturday,
		Sunday:                   req.Sunday,
		TimeOfDay:                req.TimeOfDay,
		LastDateOfTaskGeneration: req.LastDateOfTaskGeneration,
		CreatedAt:                req.CreatedAt,
		ModifiedAt:               req.ModifiedAt,
		Tags:                     req.Tags,
		SpaceID:                  req.SpaceID,
		UserID:                   uid,
	}

	if err := h.taskRepo.CreateRepetitiveTaskTemplate(&repetitiveTaskTemplate); err != nil {
		utils.SendErrorResponse(c, h.logger, messages.ErrTaskCreationFailed,
			err.Error(), apperrors.ErrInternalServerError)
		return
	}

	c.JSON(http.StatusOK, utils.CreateJSONResponse(messages.Success,
		messages.MsgRepetitiveTaskTemplateCreationSuccess, repetitiveTaskTemplate))
}
