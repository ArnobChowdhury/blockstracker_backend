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

type TagHandler struct {
	tagRepo    *repositories.TagRepository
	changeRepo *repositories.ChangeRepository
	db         *gorm.DB
	logger     *zap.SugaredLogger
}

func NewTagHandler(
	tagRepo *repositories.TagRepository,
	changeRepo *repositories.ChangeRepository,
	db *gorm.DB,
	logger *zap.SugaredLogger,
) *TagHandler {
	return &TagHandler{
		tagRepo:    tagRepo,
		changeRepo: changeRepo,
		db:         db,
		logger:     logger,
	}
}

// CreateTag godoc
// @Summary Create a new tag
// @Description Create a new tag with the given details
// @Tags tags
// @Accept json
// @Produce json
// @Param tag body models.TagRequest true "Tag details"
// @Success 200 {object} models.TagResponseForSwagger
// @Failure 400 {object} models.GenericErrorResponse
// @Failure 500 {object} models.GenericErrorResponse
// @Router /tags [post]
func (h *TagHandler) CreateTag(c *gin.Context) {
	uid, err := utils.ExtractUIDFromGinContext(c)
	if err != nil {
		utils.SendErrorResponse(c, h.logger, messages.ErrTagCreationFailed, err.LogError(),
			apperrors.ErrInternalServerError)
		return
	}

	var req models.TagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		invalidReqErr := apperrors.NewInvalidReqErr(err.Error())
		utils.SendErrorResponse(c, h.logger, messages.ErrTagCreationFailed,
			err.Error(), invalidReqErr)
		return
	}

	tag := models.Tag{
		ID:         req.ID,
		Name:       req.Name,
		CreatedAt:  req.CreatedAt,
		ModifiedAt: req.ModifiedAt,
		UserID:     uid,
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

	if err := h.tagRepo.CreateTag(tx, &tag); err != nil {
		tx.Rollback()
		utils.SendErrorResponse(c, h.logger, messages.ErrTagCreationFailed,
			err.Error(), apperrors.ErrInternalServerError)
		return
	}

	change := models.Change{
		UserID:     uid,
		EntityType: "tag",
		EntityID:   tag.ID,
		Operation:  "create",
	}
	if err := h.changeRepo.CreateChange(tx, &change); err != nil {
		tx.Rollback()
		utils.SendErrorResponse(c, h.logger, "Failed to create change record",
			err.Error(), apperrors.ErrInternalServerError)
		return
	}

	if err := tx.Model(&tag).Update("last_change_id", change.ChangeID).Error; err != nil {
		tx.Rollback()
		utils.SendErrorResponse(c, h.logger, "Failed to update tag with change ID",
			err.Error(), apperrors.ErrInternalServerError)
		return
	}

	if err := tx.Commit().Error; err != nil {
		utils.SendErrorResponse(c, h.logger, "Failed to commit transaction",
			err.Error(), apperrors.ErrInternalServerError)
		return
	}
	tag.LastChangeID = change.ChangeID
	c.JSON(http.StatusOK, utils.CreateJSONResponse(messages.Success, messages.MsgTagCreationSuccess, tag))
}

// UpdateTag godoc
// @Summary Update an existing tag
// @Description Update an existing tag with the given details
// @Tags tags
// @Accept json
// @Produce json
// @Param id path string true "Tag ID"
// @Param tag body models.TagRequest true "Tag details"
// @Success 200 {object} models.TagResponseForSwagger
// @Failure 400 {object} models.GenericErrorResponse
// @Failure 500 {object} models.GenericErrorResponse
// @Router /tags/{id} [put]
func (h *TagHandler) UpdateTag(c *gin.Context) {
	uid, err := utils.ExtractUIDFromGinContext(c)
	if err != nil {
		utils.SendErrorResponse(c, h.logger, messages.ErrTagUpdateFailed,
			err.LogError(), apperrors.ErrInternalServerError)
		return
	}

	tagIDStr := c.Param("id")
	tagID, parseErr := uuid.Parse(tagIDStr)
	if parseErr != nil {
		utils.SendErrorResponse(c, h.logger, messages.ErrTagUpdateFailed,
			fmt.Sprintf("Invalid tag ID format: %s", tagIDStr), apperrors.NewInvalidReqErr("Invalid tag ID"))
		return
	}

	var req models.TagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		invalidReqErr := apperrors.NewInvalidReqErr(err.Error())
		utils.SendErrorResponse(c, h.logger, messages.ErrTagUpdateFailed, err.Error(), invalidReqErr)
		return
	}

	tag := models.Tag{
		ID:         tagID,
		Name:       req.Name,
		CreatedAt:  req.CreatedAt,
		ModifiedAt: req.ModifiedAt,
		UserID:     uid,
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

	existingTag, fetchErr := h.tagRepo.GetTagByID(tx, tagID, uid)
	if fetchErr != nil {
		tx.Rollback()
		if errors.Is(fetchErr, gorm.ErrRecordNotFound) {
			utils.SendErrorResponse(c, h.logger, messages.ErrTagUpdateFailed, "Tag not found or does not belong to user", apperrors.ErrNotFound)
		} else {
			utils.SendErrorResponse(c, h.logger, messages.ErrTagUpdateFailed,
				fetchErr.Error(), apperrors.ErrInternalServerError)
		}
		return
	}

	if time.Time(req.ModifiedAt).Before(time.Time(existingTag.ModifiedAt)) {
		tx.Rollback()
		logMsg := fmt.Sprintf("Stale update rejected for tag_id: %s. Incoming timestamp: %s, Database timestamp: %s",
			tagID, time.Time(req.ModifiedAt).Format(time.RFC3339), time.Time(existingTag.ModifiedAt).Format(time.RFC3339))

		utils.SendErrorResponse(c, h.logger, messages.ErrTagUpdateFailed, logMsg, apperrors.ErrStaleData)
		return
	}

	if err := h.tagRepo.UpdateTag(tx, &tag); err != nil {
		tx.Rollback()
		utils.SendErrorResponse(c, h.logger, messages.ErrTagUpdateFailed, err.Error(), apperrors.ErrInternalServerError)
		return
	}

	change := models.Change{
		UserID:     uid,
		EntityType: "tag",
		EntityID:   tag.ID,
		Operation:  "update",
	}
	if err := h.changeRepo.CreateChange(tx, &change); err != nil {
		tx.Rollback()
		utils.SendErrorResponse(c, h.logger, "Failed to create change record",
			err.Error(), apperrors.ErrInternalServerError)
		return
	}

	if err := tx.Model(&tag).Update("last_change_id", change.ChangeID).Error; err != nil {
		tx.Rollback()
		utils.SendErrorResponse(c, h.logger, "Failed to update tag with change ID",
			err.Error(), apperrors.ErrInternalServerError)
		return
	}

	if err := tx.Commit().Error; err != nil {
		utils.SendErrorResponse(c, h.logger, "Failed to commit transaction",
			err.Error(), apperrors.ErrInternalServerError)
		return
	}
	tag.LastChangeID = change.ChangeID
	c.JSON(http.StatusOK, utils.CreateJSONResponse(messages.Success, messages.MsgTagUpdateSuccess, tag))
}

func (h *TagHandler) GetTagsFromVersion(c *gin.Context) {
}
