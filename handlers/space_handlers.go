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

type SpaceHandler struct {
	SpaceRepo  *repositories.SpaceRepository
	changeRepo *repositories.ChangeRepository
	db         *gorm.DB
	logger     *zap.SugaredLogger
}

func NewSpaceHandler(
	SpaceRepo *repositories.SpaceRepository,
	changeRepo *repositories.ChangeRepository,
	db *gorm.DB,
	logger *zap.SugaredLogger,
) *SpaceHandler {
	return &SpaceHandler{
		SpaceRepo:  SpaceRepo,
		changeRepo: changeRepo,
		db:         db,
		logger:     logger,
	}
}

// CreateSpace godoc
// @Summary Create a new Space
// @Description Create a new Space with the given details
// @Tags spaces
// @Accept json
// @Produce json
// @Param Space body models.SpaceRequest true "Space details"
// @Success 200 {object} models.SpaceResponseForSwagger
// @Failure 400 {object} models.GenericErrorResponse
// @Failure 500 {object} models.GenericErrorResponse
// @Router /spaces [post]
func (h *SpaceHandler) CreateSpace(c *gin.Context) {
	uid, err := utils.ExtractUIDFromGinContext(c)
	if err != nil {
		utils.SendErrorResponse(c, h.logger, messages.ErrSpaceCreationFailed, err.LogError(),
			apperrors.ErrInternalServerError)
		return
	}

	var req models.SpaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		invalidReqErr := apperrors.NewInvalidReqErr(err.Error())
		utils.SendErrorResponse(c, h.logger, messages.ErrSpaceCreationFailed,
			err.Error(), invalidReqErr)
		return
	}

	Space := models.Space{
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

	if err := h.SpaceRepo.CreateSpace(tx, &Space); err != nil {
		tx.Rollback()
		utils.SendErrorResponse(c, h.logger, messages.ErrSpaceCreationFailed,
			err.Error(), apperrors.ErrInternalServerError)
		return
	}

	change := models.Change{
		UserID:     uid,
		EntityType: "space",
		EntityID:   Space.ID,
		Operation:  "create",
	}
	if err := h.changeRepo.CreateChange(tx, &change); err != nil {
		tx.Rollback()
		utils.SendErrorResponse(c, h.logger, "Failed to create change record",
			err.Error(), apperrors.ErrInternalServerError)
		return
	}

	if err := tx.Model(&Space).Update("last_change_id", change.ChangeID).Error; err != nil {
		tx.Rollback()
		utils.SendErrorResponse(c, h.logger, "Failed to update space with change ID",
			err.Error(), apperrors.ErrInternalServerError)
		return
	}

	if err := tx.Commit().Error; err != nil {
		utils.SendErrorResponse(c, h.logger, "Failed to commit transaction",
			err.Error(), apperrors.ErrInternalServerError)
		return
	}
	Space.LastChangeID = change.ChangeID
	c.JSON(http.StatusOK, utils.CreateJSONResponse(messages.Success, messages.MsgSpaceCreationSuccess, Space))
}

// UpdateSpace godoc
// @Summary Update an existing Space
// @Description Update an existing Space with the given details
// @Tags spaces
// @Accept json
// @Produce json
// @Param id path string true "Space ID"
// @Param space body models.SpaceRequest true "Space details"
// @Success 200 {object} models.SpaceResponseForSwagger
// @Failure 400 {object} models.GenericErrorResponse
// @Failure 404 {object} models.GenericErrorResponse
// @Failure 500 {object} models.GenericErrorResponse
// @Router /spaces/{id} [put]
func (h *SpaceHandler) UpdateSpace(c *gin.Context) {
	uid, err := utils.ExtractUIDFromGinContext(c)
	if err != nil {
		utils.SendErrorResponse(c, h.logger, messages.ErrSpaceUpdateFailed, err.LogError(),
			apperrors.ErrInternalServerError)
		return
	}

	spaceIDStr := c.Param("id")
	spaceID, parseErr := uuid.Parse(spaceIDStr)
	if parseErr != nil {
		utils.SendErrorResponse(c, h.logger, messages.ErrSpaceUpdateFailed,
			fmt.Sprintf("Invalid space ID format: %s", spaceIDStr),
			apperrors.NewInvalidReqErr("Invalid space ID"))
		return
	}

	var req models.SpaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		invalidReqErr := apperrors.NewInvalidReqErr(err.Error())
		utils.SendErrorResponse(c, h.logger, messages.ErrSpaceUpdateFailed,
			err.Error(), invalidReqErr)
		return
	}

	updateData := map[string]any{
		"name":        req.Name,
		"modified_at": req.ModifiedAt,
		"user_id":     uid,
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

	existingSpace, fetchErr := h.SpaceRepo.GetSpaceByID(tx, spaceID, uid)
	if fetchErr != nil {
		tx.Rollback()
		if errors.Is(fetchErr, gorm.ErrRecordNotFound) {
			utils.SendErrorResponse(c, h.logger, messages.ErrSpaceUpdateFailed,
				"Space not found or does not belong to user", apperrors.ErrNotFound)
		} else {
			utils.SendErrorResponse(c, h.logger, messages.ErrSpaceUpdateFailed,
				fetchErr.Error(), apperrors.ErrInternalServerError)
		}
		return
	}

	if time.Time(req.ModifiedAt).Before(time.Time(existingSpace.ModifiedAt)) {
		tx.Rollback()
		utils.SendErrorResponse(c, h.logger, messages.ErrSpaceUpdateFailed, "Stale data", apperrors.ErrStaleData)
		return
	}

	if err := h.SpaceRepo.UpdateSpace(tx, spaceID, uid, updateData); err != nil {
		tx.Rollback()
		utils.SendErrorResponse(c, h.logger, messages.ErrSpaceUpdateFailed,
			err.Error(), apperrors.ErrInternalServerError)
		return
	}

	change := models.Change{
		UserID:     uid,
		EntityType: "space",
		EntityID:   spaceID,
		Operation:  "update",
	}
	if err := h.changeRepo.CreateChange(tx, &change); err != nil {
		tx.Rollback()
		utils.SendErrorResponse(c, h.logger, "Failed to create change record",
			err.Error(), apperrors.ErrInternalServerError)
		return
	}

	if err := tx.Model(&models.Space{}).Where("id = ?", spaceID).Update("last_change_id", change.ChangeID).Error; err != nil {
		tx.Rollback()
		utils.SendErrorResponse(c, h.logger, "Failed to update space with change ID",
			err.Error(), apperrors.ErrInternalServerError)
		return
	}

	if err := tx.Commit().Error; err != nil {
		utils.SendErrorResponse(c, h.logger, "Failed to commit transaction",
			err.Error(), apperrors.ErrInternalServerError)
		return
	}

	updatedSpace, getErr := h.SpaceRepo.GetSpaceByID(h.db, spaceID, uid)
	if getErr != nil {
		utils.SendErrorResponse(c, h.logger, "Update succeeded, but failed to fetch the updated record for response.",
			getErr.Error(), apperrors.ErrInternalServerError)
		return
	}

	updatedSpace.LastChangeID = change.ChangeID
	c.JSON(http.StatusOK, utils.CreateJSONResponse(messages.Success, messages.MsgSpaceUpdateSuccess, updatedSpace))
}

func (h *SpaceHandler) GetSpacesFromVersion(c *gin.Context) {
}
