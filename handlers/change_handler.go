package handlers

import (
	apperrors "blockstracker_backend/internal/errors"
	"blockstracker_backend/internal/repositories"
	"blockstracker_backend/internal/utils"
	messages "blockstracker_backend/messages"
	"blockstracker_backend/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	EntityTypeTask                   = "task"
	EntityTypeTag                    = "tag"
	EntityTypeSpace                  = "space"
	EntityTypeRepetitiveTaskTemplate = "repetitive_task_template"

	OperationCreate = "create"
	OperationUpdate = "update"
	OperationDelete = "delete"
)

type ChangeHandler struct {
	db         *gorm.DB
	changeRepo *repositories.ChangeRepository
	taskRepo   *repositories.TaskRepository
	tagRepo    *repositories.TagRepository
	spaceRepo  *repositories.SpaceRepository
	logger     *zap.SugaredLogger
}

func NewChangeHandler(
	db *gorm.DB,
	changeRepo *repositories.ChangeRepository,
	taskRepo *repositories.TaskRepository,
	tagRepo *repositories.TagRepository,
	spaceRepo *repositories.SpaceRepository,
	logger *zap.SugaredLogger,
) *ChangeHandler {
	return &ChangeHandler{
		db:         db,
		changeRepo: changeRepo,
		taskRepo:   taskRepo,
		tagRepo:    tagRepo,
		spaceRepo:  spaceRepo,
		logger:     logger,
	}
}

// SyncChanges godoc
// @Summary      Sync changes
// @Description  Get all entity changes since the last sync.
// @Tags         Sync
// @Accept       json
// @Produce      json
// @Param        last_change_id query int false "The last change ID received by the client. If 0 or omitted, all entities are returned."
// @Success      200  {object}  models.SyncResponse
// @Failure      400  {object}  models.GenericErrorResponse "Invalid last_change_id"
// @Failure      401  {object}  models.GenericErrorResponse "Unauthorized"
// @Failure      500  {object}  models.GenericErrorResponse "Internal Server Error"
// @Router       /changes/sync [get]
func (h *ChangeHandler) SyncChanges(c *gin.Context) {
	uid, uidExtractionErr := utils.ExtractUIDFromGinContext(c)
	if uidExtractionErr != nil {
		utils.SendErrorResponse(c, h.logger, messages.ErrSyncFailed, uidExtractionErr.LogError(),
			apperrors.ErrInternalServerError)
		return
	}

	lastChangeIDStr := c.DefaultQuery("last_change_id", "0")
	lastChangeID, err := strconv.ParseInt(lastChangeIDStr, 10, 64)
	if err != nil {
		utils.SendErrorResponse(c, h.logger, messages.ErrSyncFailed, err.Error(),
			apperrors.NewInvalidReqErr())
		return
	}

	changes, err := h.changeRepo.GetChangesSince(h.db, uid, lastChangeID)
	if err != nil {
		utils.SendErrorResponse(c, h.logger, messages.ErrSyncFailed, err.Error(),
			apperrors.ErrInternalServerError)
		return
	}

	finalChanges := make(map[uuid.UUID]models.Change)
	latestChangeID := lastChangeID
	for _, change := range changes {
		finalChanges[change.EntityID] = change
		if change.ChangeID > latestChangeID {
			latestChangeID = change.ChangeID
		}
	}

	taskIDs := []uuid.UUID{}
	tagIDs := []uuid.UUID{}
	spaceIDs := []uuid.UUID{}
	templateIDs := []uuid.UUID{}

	for _, finalChange := range finalChanges {
		switch finalChange.EntityType {
		case EntityTypeTask:
			taskIDs = append(taskIDs, finalChange.EntityID)
		case EntityTypeTag:
			tagIDs = append(tagIDs, finalChange.EntityID)
		case EntityTypeSpace:
			spaceIDs = append(spaceIDs, finalChange.EntityID)
		case EntityTypeRepetitiveTaskTemplate:
			templateIDs = append(templateIDs, finalChange.EntityID)
		}
	}

	syncResponse := models.SyncResponse{
		LatestChangeID: latestChangeID,
	}

	if len(templateIDs) > 0 {
		templates, err := h.taskRepo.GetRepetitiveTaskTemplatesByIDs(h.db, templateIDs, uid)
		if err != nil {
			utils.SendErrorResponse(c, h.logger, messages.ErrSyncFailed, err.Error(),
				apperrors.ErrInternalServerError)
			return
		}
		syncResponse.RepetitiveTaskTemplates = templates
	}
	if len(taskIDs) > 0 {
		tasks, err := h.taskRepo.GetTasksByIDs(h.db, taskIDs, uid)
		if err != nil {
			utils.SendErrorResponse(c, h.logger, messages.ErrSyncFailed, err.Error(),
				apperrors.ErrInternalServerError)
			return
		}
		syncResponse.Tasks = tasks
	}
	if len(tagIDs) > 0 {
		tags, err := h.tagRepo.GetTagsByIDs(h.db, tagIDs, uid)
		if err != nil {
			utils.SendErrorResponse(c, h.logger, messages.ErrSyncFailed, err.Error(),
				apperrors.ErrInternalServerError)
			return
		}
		syncResponse.Tags = tags
	}
	if len(spaceIDs) > 0 {
		spaces, err := h.spaceRepo.GetSpacesByIDs(h.db, spaceIDs, uid)
		if err != nil {
			utils.SendErrorResponse(c, h.logger, messages.ErrSyncFailed, err.Error(),
				apperrors.ErrInternalServerError)
			return
		}
		syncResponse.Spaces = spaces
	}

	c.JSON(http.StatusOK, utils.CreateJSONResponse(
		messages.Success, messages.MsgSyncSuccessful, syncResponse))
}
