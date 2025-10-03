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
	taskRepo   *repositories.TaskRepository
	changeRepo *repositories.ChangeRepository
	db         *gorm.DB
	logger     *zap.SugaredLogger
}

func NewTaskHandler(
	taskRepo *repositories.TaskRepository,
	changeRepo *repositories.ChangeRepository,
	db *gorm.DB,
	logger *zap.SugaredLogger,
) *TaskHandler {
	return &TaskHandler{
		taskRepo:   taskRepo,
		changeRepo: changeRepo,
		db:         db,
		logger:     logger,
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
		ID:                       req.ID,
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

	tx := h.db.Begin()
	if tx.Error != nil {
		utils.SendErrorResponse(c, h.logger, "Failed to begin transaction",
			tx.Error.Error(), apperrors.ErrInternalServerError)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := h.taskRepo.CreateTask(tx, &task); err != nil {
		tx.Rollback()
		utils.SendErrorResponse(c, h.logger, messages.ErrTaskCreationFailed,
			err.Error(), apperrors.ErrInternalServerError)
		return
	}

	change := models.Change{
		UserID:     uid,
		EntityType: "task",
		EntityID:   task.ID,
		Operation:  "create",
	}
	if err := h.changeRepo.CreateChange(tx, &change); err != nil {
		tx.Rollback()
		utils.SendErrorResponse(c, h.logger, "Failed to create change record",
			err.Error(), apperrors.ErrInternalServerError)
		return
	}

	task.LastChangeID = change.ChangeID
	if err := tx.Save(&task).Commit().Error; err != nil {
		tx.Rollback()
		utils.SendErrorResponse(c, h.logger, "Failed to commit transaction",
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

	tx := h.db.Begin()
	if tx.Error != nil {
		utils.SendErrorResponse(c, h.logger, "Failed to begin transaction",
			tx.Error.Error(), apperrors.ErrInternalServerError)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	existingTask, fetchErr := h.taskRepo.GetTaskByID(tx, taskID, uid)
	if fetchErr != nil {
		tx.Rollback()
		if errors.Is(fetchErr, gorm.ErrRecordNotFound) {
			utils.SendErrorResponse(c, h.logger, messages.ErrTaskUpdateFailed,
				"Task not found or does not belong to user", apperrors.ErrUnauthorized)
		} else {
			utils.SendErrorResponse(c, h.logger, messages.ErrTaskUpdateFailed,
				fetchErr.Error(), apperrors.ErrInternalServerError)
		}
		return
	}

	if req.ModifiedAt.Before(existingTask.ModifiedAt) {
		tx.Rollback()
		logMsg := fmt.Sprintf("Stale update rejected for task_id: %s. Incoming timestamp: %s, Database timestamp: %s",
			taskID, req.ModifiedAt, existingTask.ModifiedAt)
		utils.SendErrorResponse(c, h.logger, messages.ErrTaskUpdateFailed, logMsg, apperrors.ErrStaleData)
		return
	}

	if err := h.taskRepo.UpdateTask(tx, &task); err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SendErrorResponse(c, h.logger, messages.ErrTaskUpdateFailed,
				"Task not found or does not belong to user", apperrors.ErrUnauthorized)
		} else {
			utils.SendErrorResponse(c, h.logger, messages.ErrTaskUpdateFailed,
				err.Error(), apperrors.ErrInternalServerError)
		}
		return
	}

	change := models.Change{
		UserID:     uid,
		EntityType: "task",
		EntityID:   task.ID,
		Operation:  "update",
	}
	if err := h.changeRepo.CreateChange(tx, &change); err != nil {
		tx.Rollback()
		utils.SendErrorResponse(c, h.logger, "Failed to create change record",
			err.Error(), apperrors.ErrInternalServerError)
		return
	}

	task.LastChangeID = change.ChangeID
	if err := tx.Model(&task).Update("last_change_id", change.ChangeID).Commit().Error; err != nil {
		tx.Rollback()
		utils.SendErrorResponse(c, h.logger, "Failed to commit transaction",
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
// @Param task body models.RepetitiveTaskTemplateRequest true "Repetitive task template details"
// @Success 200 {object} models.RepetitiveTaskTemplateResponseForSwagger
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

	var req models.RepetitiveTaskTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		invalidReqErr := apperrors.NewInvalidReqErr(err.Error())
		utils.SendErrorResponse(c, h.logger, messages.ErrRepetitiveTaskTemplateCreationFailed,
			err.Error(), invalidReqErr)
		return
	}

	repetitiveTaskTemplate := models.RepetitiveTaskTemplate{
		ID:                       req.ID,
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

	tx := h.db.Begin()
	if tx.Error != nil {
		utils.SendErrorResponse(c, h.logger, "Failed to begin transaction",
			tx.Error.Error(), apperrors.ErrInternalServerError)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := h.taskRepo.CreateRepetitiveTaskTemplate(tx, &repetitiveTaskTemplate); err != nil {
		tx.Rollback()
		utils.SendErrorResponse(c, h.logger, messages.ErrRepetitiveTaskTemplateCreationFailed,
			err.Error(), apperrors.ErrInternalServerError)
		return
	}

	change := models.Change{
		UserID:     uid,
		EntityType: "repetitive_task_template",
		EntityID:   repetitiveTaskTemplate.ID,
		Operation:  "create",
	}
	if err := h.changeRepo.CreateChange(tx, &change); err != nil {
		tx.Rollback()
		utils.SendErrorResponse(c, h.logger, "Failed to create change record",
			err.Error(), apperrors.ErrInternalServerError)
		return
	}

	repetitiveTaskTemplate.LastChangeID = change.ChangeID
	if err := tx.Save(&repetitiveTaskTemplate).Commit().Error; err != nil {
		tx.Rollback()
		utils.SendErrorResponse(c, h.logger, "Failed to commit transaction",
			err.Error(), apperrors.ErrInternalServerError)
		return
	}
	c.JSON(http.StatusOK, utils.CreateJSONResponse(messages.Success,
		messages.MsgRepetitiveTaskTemplateCreationSuccess, repetitiveTaskTemplate))
}

// UpdateRepetitiveTaskTemplate godoc
// @Summary Update an existing repetitive task template
// @Description Update an existing repetitive task template with the given details
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Repetitive Task Template ID"
// @Param task body models.RepetitiveTaskTemplateRequest true "Repetitive task template details"
// @Success 200 {object} models.RepetitiveTaskTemplateResponseForSwagger
// @Failure 400 {object} models.GenericErrorResponse
// @Failure 404 {object} models.GenericErrorResponse
// @Failure 500 {object} models.GenericErrorResponse
// @Router /tasks/repetitive/{id} [put]
func (h *TaskHandler) UpdateRepetitiveTaskTemplate(c *gin.Context) {
	uid, err := utils.ExtractUIDFromGinContext(c)
	if err != nil {
		utils.SendErrorResponse(c, h.logger, messages.ErrRepetitiveTaskTemplateUpdateFailed,
			err.LogError(), apperrors.ErrInternalServerError)
		return
	}

	repetitiveTaskTemplateIDStr := c.Param("id")
	repetitiveTaskTemplateID, taskIdParseErr := uuid.Parse(repetitiveTaskTemplateIDStr)
	if taskIdParseErr != nil {
		utils.SendErrorResponse(c, h.logger, messages.ErrRepetitiveTaskTemplateUpdateFailed,
			fmt.Sprintf("Invalid repetitive task template ID format: %s", repetitiveTaskTemplateIDStr),
			apperrors.ErrMalformedRepetitiveTaskTemplateRequest)
		return
	}

	var req models.RepetitiveTaskTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		invalidReqErr := apperrors.NewInvalidReqErr(err.Error())
		utils.SendErrorResponse(c, h.logger, messages.ErrRepetitiveTaskTemplateUpdateFailed,
			err.Error(), invalidReqErr)
		return
	}

	repetitiveTaskTemplate := models.RepetitiveTaskTemplate{
		ID:                       repetitiveTaskTemplateID,
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
		ModifiedAt:               req.ModifiedAt,
		Tags:                     req.Tags,
		SpaceID:                  req.SpaceID,
		UserID:                   uid,
	}

	tx := h.db.Begin()
	if tx.Error != nil {
		utils.SendErrorResponse(c, h.logger, "Failed to begin transaction",
			tx.Error.Error(), apperrors.ErrInternalServerError)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	existingTemplate, fetchErr := h.taskRepo.GetRepetitiveTaskTemplateByID(tx, repetitiveTaskTemplateID, uid)
	if fetchErr != nil {
		tx.Rollback()
		if errors.Is(fetchErr, gorm.ErrRecordNotFound) {
			utils.SendErrorResponse(c, h.logger, messages.ErrRepetitiveTaskTemplateUpdateFailed,
				"Repetitive task template not found or does not belong to user", apperrors.ErrUnauthorized)
		} else {
			utils.SendErrorResponse(c, h.logger, messages.ErrRepetitiveTaskTemplateUpdateFailed,
				fetchErr.Error(), apperrors.ErrInternalServerError)
		}
		return
	}

	if req.ModifiedAt.Before(existingTemplate.ModifiedAt) {
		tx.Rollback()
		logMsg := fmt.Sprintf("Stale update rejected for repetitive_task_template_id: %s. Incoming timestamp: %s, Database timestamp: %s",
			repetitiveTaskTemplateID, req.ModifiedAt, existingTemplate.ModifiedAt)
		utils.SendErrorResponse(c, h.logger, messages.ErrRepetitiveTaskTemplateUpdateFailed, logMsg, apperrors.ErrStaleData)
		return
	}

	if err := h.taskRepo.UpdateRepetitiveTaskTemplate(tx, &repetitiveTaskTemplate); err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SendErrorResponse(c, h.logger, messages.ErrRepetitiveTaskTemplateUpdateFailed,
				"Repetitive task template not found or does not belong to user", apperrors.ErrUnauthorized)
		} else {
			utils.SendErrorResponse(c, h.logger, messages.ErrRepetitiveTaskTemplateUpdateFailed,
				err.Error(), apperrors.ErrInternalServerError)
		}
		return
	}

	change := models.Change{
		UserID:     uid,
		EntityType: "repetitive_task_template",
		EntityID:   repetitiveTaskTemplate.ID,
		Operation:  "update",
	}
	if err := h.changeRepo.CreateChange(tx, &change); err != nil {
		tx.Rollback()
		utils.SendErrorResponse(c, h.logger, "Failed to create change record",
			err.Error(), apperrors.ErrInternalServerError)
		return
	}

	repetitiveTaskTemplate.LastChangeID = change.ChangeID
	if err := tx.Model(&repetitiveTaskTemplate).Update("last_change_id", change.ChangeID).Commit().Error; err != nil {
		tx.Rollback()
		utils.SendErrorResponse(c, h.logger, "Failed to commit transaction",
			err.Error(), apperrors.ErrInternalServerError)
		return
	}
	c.JSON(http.StatusOK, utils.CreateJSONResponse(
		messages.Success, messages.MsgRepetitiveTaskTemplateUpdateSuccess, repetitiveTaskTemplate))
}
