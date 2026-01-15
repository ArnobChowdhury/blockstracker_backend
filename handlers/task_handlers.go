package handlers

import (
	"errors"
	"fmt"
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
		IsActive:                 *req.IsActive,
		Title:                    req.Title,
		Description:              req.Description,
		Schedule:                 req.Schedule,
		Priority:                 *req.Priority,
		CompletionStatus:         req.CompletionStatus,
		DueDate:                  req.DueDate,
		ShouldBeScored:           req.ShouldBeScored,
		Score:                    req.Score,
		TimeOfDay:                req.TimeOfDay,
		RepetitiveTaskTemplateID: req.RepetitiveTaskTemplateID,
		CreatedAt:                req.CreatedAt,
		ModifiedAt:               req.ModifiedAt,
		// Tags:                     req.Tags,
		SpaceID: req.SpaceID,
		UserID:  uid,
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

	tx.SavePoint("before_create")

	if err := h.taskRepo.CreateTask(tx, &task); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			tx.RollbackTo("before_create")

			// 1. Check for ID collision (Hydration/Restore case)
			existingTask, fetchErr := h.taskRepo.GetTaskByID(tx, task.ID, uid)
			if fetchErr == nil {
				// ID exists. Check timestamps.
				if time.Time(task.ModifiedAt).After(time.Time(existingTask.ModifiedAt)) {
					// Incoming is newer. Update.
					updateData := map[string]interface{}{
						"is_active":                   task.IsActive,
						"title":                       task.Title,
						"description":                 task.Description,
						"schedule":                    task.Schedule,
						"priority":                    task.Priority,
						"completion_status":           task.CompletionStatus,
						"due_date":                    task.DueDate,
						"should_be_scored":            task.ShouldBeScored,
						"score":                       task.Score,
						"time_of_day":                 task.TimeOfDay,
						"repetitive_task_template_id": task.RepetitiveTaskTemplateID,
						"modified_at":                 task.ModifiedAt,
						"space_id":                    task.SpaceID,
						"user_id":                     uid,
					}

					if err := h.taskRepo.UpdateTask(tx, task.ID, uid, updateData); err != nil {
						tx.Rollback()
						utils.SendErrorResponse(c, h.logger, messages.ErrTaskUpdateFailed, err.Error(), apperrors.ErrInternalServerError)
						return
					}

					change := models.Change{UserID: uid, EntityType: "task", EntityID: task.ID, Operation: "update"}
					if err := h.changeRepo.CreateChange(tx, &change); err != nil {
						tx.Rollback()
						utils.SendErrorResponse(c, h.logger, "Failed to create change record", err.Error(), apperrors.ErrInternalServerError)
						return
					}
					if err := tx.Model(&models.Task{}).Where("id = ?", task.ID).Update("last_change_id", change.ChangeID).Error; err != nil {
						tx.Rollback()
						utils.SendErrorResponse(c, h.logger, "Failed to update task with change ID", err.Error(), apperrors.ErrInternalServerError)
						return
					}
					task.LastChangeID = change.ChangeID
				} else {
					// Incoming is older or equal. Server wins. Return existing state.
					task = *existingTask
				}

				if err := tx.Commit().Error; err != nil {
					utils.SendErrorResponse(c, h.logger, "Failed to commit transaction", err.Error(), apperrors.ErrInternalServerError)
					return
				}
				c.JSON(http.StatusOK, utils.CreateJSONResponse(messages.Success, "Task synced successfully (upsert)", task))
				return
			}

			// 2. Check for Logic Collision (Unique Constraint on TemplateID + DueDate)
			if task.RepetitiveTaskTemplateID != nil && task.DueDate != nil {
				existingTask, fetchErr := h.taskRepo.GetTaskByRepetitiveTemplateIDAndDueDate(tx, *task.RepetitiveTaskTemplateID, time.Time(*task.DueDate), uid)
				if fetchErr == nil {
					tx.Rollback()
					utils.SendErrorResponse(c, h.logger, messages.ErrTaskCreationFailed,
						"Duplicate task creation attempt",
						apperrors.ErrDuplicateEntity,
						gin.H{"canonical_id": existingTask.ID.String()})
					return
				}
				h.logger.Errorw("Failed to find existing duplicate task after unique constraint violation",
					"error", fetchErr.Error(), "repetitiveTaskTemplateId", task.RepetitiveTaskTemplateID, "dueDate", task.DueDate)
			}
			tx.Rollback()
			utils.SendErrorResponse(c, h.logger, messages.ErrTaskCreationFailed,
				"Duplicate task creation attempt (could not determine canonical ID)",
				apperrors.ErrDuplicateEntity)
			return
		}
		tx.Rollback()
		utils.SendErrorResponse(c, h.logger, messages.ErrTaskCreationFailed, err.Error(), apperrors.ErrInternalServerError)
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

	if err := tx.Model(&task).Update("last_change_id", change.ChangeID).Error; err != nil {
		tx.Rollback()
		utils.SendErrorResponse(c, h.logger, "Failed to update task with change ID",
			err.Error(), apperrors.ErrInternalServerError)
		return
	}

	if err := tx.Commit().Error; err != nil {
		utils.SendErrorResponse(c, h.logger, "Failed to commit transaction",
			err.Error(), apperrors.ErrInternalServerError)
		return
	}

	task.LastChangeID = change.ChangeID
	c.JSON(http.StatusOK, utils.CreateJSONResponse(messages.Success, messages.MsgTaskCreationSuccess, task))
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

	updateData := map[string]interface{}{
		"is_active":                   *req.IsActive,
		"title":                       req.Title,
		"description":                 req.Description,
		"schedule":                    req.Schedule,
		"priority":                    *req.Priority,
		"completion_status":           req.CompletionStatus,
		"due_date":                    req.DueDate,
		"should_be_scored":            req.ShouldBeScored,
		"score":                       req.Score,
		"time_of_day":                 req.TimeOfDay,
		"repetitive_task_template_id": req.RepetitiveTaskTemplateID,
		"modified_at":                 req.ModifiedAt,
		"space_id":                    req.SpaceID,
		"user_id":                     uid,
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
				"Task not found or does not belong to user", apperrors.ErrNotFound)
		} else {
			utils.SendErrorResponse(c, h.logger, messages.ErrTaskUpdateFailed,
				fetchErr.Error(), apperrors.ErrInternalServerError)
		}
		return
	}

	if time.Time(req.ModifiedAt).Before(time.Time(existingTask.ModifiedAt)) {
		tx.Rollback()
		logMsg := fmt.Sprintf("Stale update rejected for task_id: %s. Incoming timestamp: %s, Database timestamp: %s", taskID, time.Time(req.ModifiedAt).Format(time.RFC3339), time.Time(existingTask.ModifiedAt).Format(time.RFC3339))

		utils.SendErrorResponse(c, h.logger, messages.ErrTaskUpdateFailed, logMsg, apperrors.ErrStaleData)
		return
	}

	if err := h.taskRepo.UpdateTask(tx, taskID, uid, updateData); err != nil {
		tx.Rollback()
		utils.SendErrorResponse(c, h.logger, messages.ErrTaskUpdateFailed,
			err.Error(), apperrors.ErrInternalServerError)
		return
	}

	change := models.Change{
		UserID:     uid,
		EntityType: "task",
		EntityID:   taskID,
		Operation:  "update",
	}
	if err := h.changeRepo.CreateChange(tx, &change); err != nil {
		tx.Rollback()
		utils.SendErrorResponse(c, h.logger, "Failed to create change record",
			err.Error(), apperrors.ErrInternalServerError)
		return
	}

	if err := tx.Model(&models.Task{}).Where("id = ?", taskID).Update("last_change_id", change.ChangeID).Error; err != nil {
		tx.Rollback()
		utils.SendErrorResponse(c, h.logger, "Failed to update task with change ID",
			err.Error(), apperrors.ErrInternalServerError)
		return
	}

	if err := tx.Commit().Error; err != nil {
		utils.SendErrorResponse(c, h.logger, "Failed to commit transaction",
			err.Error(), apperrors.ErrInternalServerError)
		return
	}

	updatedTask, getErr := h.taskRepo.GetTaskByID(h.db, taskID, uid)
	if getErr != nil {
		utils.SendErrorResponse(c, h.logger, "Update succeeded, but failed to fetch the updated record for response.",
			getErr.Error(), apperrors.ErrInternalServerError)
		return
	}

	updatedTask.LastChangeID = change.ChangeID
	c.JSON(http.StatusOK, utils.CreateJSONResponse(messages.Success, messages.MsgTaskUpdateSuccess, updatedTask))
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
		IsActive:                 *req.IsActive,
		Title:                    req.Title,
		Description:              req.Description,
		Schedule:                 req.Schedule,
		Priority:                 *req.Priority,
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
		// Tags:                     req.Tags,
		SpaceID: req.SpaceID,
		UserID:  uid,
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

	tx.SavePoint("before_create")

	if err := h.taskRepo.CreateRepetitiveTaskTemplate(tx, &repetitiveTaskTemplate); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			tx.RollbackTo("before_create")

			// Check for ID collision
			existingTemplate, fetchErr := h.taskRepo.GetRepetitiveTaskTemplateByID(tx, repetitiveTaskTemplate.ID, uid)
			if fetchErr == nil {
				if time.Time(repetitiveTaskTemplate.ModifiedAt).After(time.Time(existingTemplate.ModifiedAt)) {
					updateData := map[string]any{
						"is_active":                    repetitiveTaskTemplate.IsActive,
						"title":                        repetitiveTaskTemplate.Title,
						"description":                  repetitiveTaskTemplate.Description,
						"schedule":                     repetitiveTaskTemplate.Schedule,
						"priority":                     repetitiveTaskTemplate.Priority,
						"should_be_scored":             repetitiveTaskTemplate.ShouldBeScored,
						"monday":                       repetitiveTaskTemplate.Monday,
						"tuesday":                      repetitiveTaskTemplate.Tuesday,
						"wednesday":                    repetitiveTaskTemplate.Wednesday,
						"thursday":                     repetitiveTaskTemplate.Thursday,
						"friday":                       repetitiveTaskTemplate.Friday,
						"saturday":                     repetitiveTaskTemplate.Saturday,
						"sunday":                       repetitiveTaskTemplate.Sunday,
						"time_of_day":                  repetitiveTaskTemplate.TimeOfDay,
						"last_date_of_task_generation": repetitiveTaskTemplate.LastDateOfTaskGeneration,
						"modified_at":                  repetitiveTaskTemplate.ModifiedAt,
						"space_id":                     repetitiveTaskTemplate.SpaceID,
						"user_id":                      uid,
					}

					if err := h.taskRepo.UpdateRepetitiveTaskTemplate(tx, repetitiveTaskTemplate.ID, uid, updateData); err != nil {
						tx.Rollback()
						utils.SendErrorResponse(c, h.logger, messages.ErrRepetitiveTaskTemplateUpdateFailed, err.Error(), apperrors.ErrInternalServerError)
						return
					}

					change := models.Change{UserID: uid, EntityType: "repetitive_task_template", EntityID: repetitiveTaskTemplate.ID, Operation: "update"}
					if err := h.changeRepo.CreateChange(tx, &change); err != nil {
						tx.Rollback()
						utils.SendErrorResponse(c, h.logger, "Failed to create change record", err.Error(), apperrors.ErrInternalServerError)
						return
					}
					if err := tx.Model(&models.RepetitiveTaskTemplate{}).Where("id = ?", repetitiveTaskTemplate.ID).Update("last_change_id", change.ChangeID).Error; err != nil {
						tx.Rollback()
						utils.SendErrorResponse(c, h.logger, "Failed to update template with change ID", err.Error(), apperrors.ErrInternalServerError)
						return
					}
					repetitiveTaskTemplate.LastChangeID = change.ChangeID
				} else {
					repetitiveTaskTemplate = *existingTemplate
				}

				if err := tx.Commit().Error; err != nil {
					utils.SendErrorResponse(c, h.logger, "Failed to commit transaction", err.Error(), apperrors.ErrInternalServerError)
					return
				}
				c.JSON(http.StatusOK, utils.CreateJSONResponse(messages.Success, "Repetitive Task Template synced successfully (upsert)", repetitiveTaskTemplate))
				return
			}
		}

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

	if err := tx.Model(&repetitiveTaskTemplate).Update("last_change_id", change.ChangeID).Error; err != nil {
		tx.Rollback()
		utils.SendErrorResponse(c, h.logger, "Failed to update repetitive task template with change ID",
			err.Error(), apperrors.ErrInternalServerError)
		return
	}

	if err := tx.Commit().Error; err != nil {
		utils.SendErrorResponse(c, h.logger, "Failed to commit transaction",
			err.Error(), apperrors.ErrInternalServerError)
		return
	}

	repetitiveTaskTemplate.LastChangeID = change.ChangeID
	c.JSON(http.StatusOK, utils.CreateJSONResponse(messages.Success, messages.MsgRepetitiveTaskTemplateCreationSuccess, repetitiveTaskTemplate))
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

	updateData := map[string]any{
		"is_active":                    *req.IsActive,
		"title":                        req.Title,
		"description":                  req.Description,
		"schedule":                     req.Schedule,
		"priority":                     *req.Priority,
		"should_be_scored":             req.ShouldBeScored,
		"monday":                       req.Monday,
		"tuesday":                      req.Tuesday,
		"wednesday":                    req.Wednesday,
		"thursday":                     req.Thursday,
		"friday":                       req.Friday,
		"saturday":                     req.Saturday,
		"sunday":                       req.Sunday,
		"time_of_day":                  req.TimeOfDay,
		"last_date_of_task_generation": req.LastDateOfTaskGeneration,
		"modified_at":                  req.ModifiedAt,
		// Tags need to be handled separately if they are being updated
		"space_id": req.SpaceID,
		"user_id":  uid,
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
				"Repetitive task template not found or does not belong to user", apperrors.ErrNotFound)
		} else {
			utils.SendErrorResponse(c, h.logger, messages.ErrRepetitiveTaskTemplateUpdateFailed,
				fetchErr.Error(), apperrors.ErrInternalServerError)
		}
		return
	}

	if time.Time(req.ModifiedAt).Before(time.Time(existingTemplate.ModifiedAt)) {
		tx.Rollback()
		logMsg := fmt.Sprintf("Stale update rejected for repetitive_task_template_id: %s. Incoming timestamp: %s, Database timestamp: %s",
			repetitiveTaskTemplateID, req.ModifiedAt, existingTemplate.ModifiedAt)
		utils.SendErrorResponse(c, h.logger, messages.ErrRepetitiveTaskTemplateUpdateFailed, logMsg, apperrors.ErrStaleData)
		return
	}

	if err := h.taskRepo.UpdateRepetitiveTaskTemplate(tx, repetitiveTaskTemplateID, uid, updateData); err != nil {
		tx.Rollback()
		utils.SendErrorResponse(c, h.logger, messages.ErrRepetitiveTaskTemplateUpdateFailed,
			err.Error(), apperrors.ErrInternalServerError)
		return
	}

	change := models.Change{
		UserID:     uid,
		EntityType: "repetitive_task_template",
		EntityID:   repetitiveTaskTemplateID,
		Operation:  "update",
	}
	if err := h.changeRepo.CreateChange(tx, &change); err != nil {
		tx.Rollback()
		utils.SendErrorResponse(c, h.logger, "Failed to create change record",
			err.Error(), apperrors.ErrInternalServerError)
		return
	}

	if err := tx.Model(&models.RepetitiveTaskTemplate{}).Where("id = ?", repetitiveTaskTemplateID).Update("last_change_id", change.ChangeID).Error; err != nil {
		tx.Rollback()
		utils.SendErrorResponse(c, h.logger, "Failed to update repetitive task template with change ID",
			err.Error(), apperrors.ErrInternalServerError)
		return
	}

	if err := tx.Commit().Error; err != nil {
		utils.SendErrorResponse(c, h.logger, "Failed to commit transaction",
			err.Error(), apperrors.ErrInternalServerError)
		return
	}

	updatedTemplate, getErr := h.taskRepo.GetRepetitiveTaskTemplateByID(h.db, repetitiveTaskTemplateID, uid)
	if getErr != nil {
		utils.SendErrorResponse(c, h.logger, "Update succeeded, but failed to fetch the updated record for response.",
			getErr.Error(), apperrors.ErrInternalServerError)
		return
	}

	updatedTemplate.LastChangeID = change.ChangeID
	c.JSON(http.StatusOK, utils.CreateJSONResponse(messages.Success, messages.MsgRepetitiveTaskTemplateUpdateSuccess, updatedTemplate))
}

type entityUpdater func(tx *gorm.DB, id, userID uuid.UUID, data map[string]any) error

type entityGetter[P models.TimeStampedEntity] func(tx *gorm.DB, id, userID uuid.UUID) (P, error)

func updateEntity[P interface {
	models.TimeStampedEntity
	*E
}, E any](
	c *gin.Context,
	h *TaskHandler,
	entityID uuid.UUID,
	uid uuid.UUID,
	updateData map[string]any,
	modifiedAt models.JSONTime,
	entityType string,
	updater entityUpdater,
	getter entityGetter[P],
) {
	tx := h.db.Begin()
	if tx.Error != nil {
		utils.SendErrorResponse(c, h.logger, "Failed to begin transaction", tx.Error.Error(), apperrors.ErrInternalServerError)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	existingEntity, fetchErr := getter(tx, entityID, uid)
	if fetchErr != nil {
		tx.Rollback()
		if errors.Is(fetchErr, gorm.ErrRecordNotFound) {
			utils.SendErrorResponse(c, h.logger, "Update failed", "Entity not found or does not belong to user", apperrors.ErrNotFound)
		} else {
			utils.SendErrorResponse(c, h.logger, "Update failed", fetchErr.Error(), apperrors.ErrInternalServerError)
		}
		return
	}

	if time.Time(modifiedAt).Before(time.Time(existingEntity.GetModifiedAt())) {
		tx.Rollback()
		logMsg := fmt.Sprintf("Stale update rejected for %s_id: %s. Incoming: %s, DB: %s",
			entityType, entityID, modifiedAt, existingEntity.GetModifiedAt())
		utils.SendErrorResponse(c, h.logger, "Update failed", logMsg, apperrors.ErrStaleData)
		return
	}

	if err := updater(tx, entityID, uid, updateData); err != nil {
		tx.Rollback()
		utils.SendErrorResponse(c, h.logger, "Update failed", err.Error(), apperrors.ErrInternalServerError)
		return
	}

	change := models.Change{
		UserID:     uid,
		EntityType: entityType,
		EntityID:   entityID,
		Operation:  "update",
	}
	if err := h.changeRepo.CreateChange(tx, &change); err != nil {
		tx.Rollback()
		utils.SendErrorResponse(c, h.logger, "Failed to create change record", err.Error(), apperrors.ErrInternalServerError)
		return
	}

	if err := tx.Model(new(E)).Where("id = ?", entityID).Update("last_change_id", change.ChangeID).Error; err != nil {
		tx.Rollback()
		utils.SendErrorResponse(c, h.logger, "Failed to update entity with change ID", err.Error(), apperrors.ErrInternalServerError)
		return
	}

	if err := tx.Commit().Error; err != nil {
		utils.SendErrorResponse(c, h.logger, "Failed to commit transaction", err.Error(), apperrors.ErrInternalServerError)
		return
	}

	updatedEntity, getErr := getter(h.db, entityID, uid)
	if getErr != nil {
		utils.SendErrorResponse(c, h.logger, "Update succeeded, but failed to fetch updated record.", getErr.Error(), apperrors.ErrInternalServerError)
		return
	}

	updatedEntity.SetLastChangeID(change.ChangeID)
	c.JSON(http.StatusOK, utils.CreateJSONResponse(messages.Success, "Update successful", updatedEntity))
}

// UpdateRepetitiveTaskTemplateLastGenDate godoc
// @Summary      Update a repetitive task template's last generation date
// @Description  Partially updates a repetitive task template, specifically its lastDateOfTaskGeneration field. This is used by the system after generating due tasks.
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id path string true "Repetitive Task Template ID"
// @Param        lastGenDate body models.UpdateRepetitiveTaskTemplateLastGenDateRequest true "Last generation date details"
// @Success      200 {object} models.RepetitiveTaskTemplateResponseForSwagger
// @Failure      400 {object} models.GenericErrorResponse
// @Failure      404 {object} models.GenericErrorResponse
// @Failure      500 {object} models.GenericErrorResponse
// @Router       /tasks/repetitive/{id}/last-gen-date [put]
func (h *TaskHandler) UpdateRepetitiveTaskTemplateLastGenDate(c *gin.Context) {
	uid, err := utils.ExtractUIDFromGinContext(c)
	if err != nil {
		utils.SendErrorResponse(c, h.logger, messages.ErrRepetitiveTaskTemplateUpdateFailed, err.LogError(), apperrors.ErrInternalServerError)
		return
	}

	repetitiveTaskTemplateIDStr := c.Param("id")
	repetitiveTaskTemplateID, parseErr := uuid.Parse(repetitiveTaskTemplateIDStr)
	if parseErr != nil {
		utils.SendErrorResponse(c, h.logger, messages.ErrRepetitiveTaskTemplateUpdateFailed, fmt.Sprintf("Invalid ID format: %s", repetitiveTaskTemplateIDStr), apperrors.NewInvalidReqErr("Invalid ID"))
		return
	}

	var req models.UpdateRepetitiveTaskTemplateLastGenDateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		invalidReqErr := apperrors.NewInvalidReqErr(err.Error())
		utils.SendErrorResponse(c, h.logger, messages.ErrRepetitiveTaskTemplateUpdateFailed, err.Error(), invalidReqErr)
		return
	}

	updateData := map[string]any{
		"last_date_of_task_generation": req.LastDateOfTaskGeneration,
		"modified_at":                  req.ModifiedAt,
	}

	updateEntity(
		c, h,
		repetitiveTaskTemplateID, uid,
		updateData, req.ModifiedAt,
		"repetitive_task_template",
		h.taskRepo.UpdateRepetitiveTaskTemplate,
		h.taskRepo.GetRepetitiveTaskTemplateByID,
	)
}
