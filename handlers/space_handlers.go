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

type SpaceHandler struct {
	SpaceRepo *repositories.SpaceRepository
	logger    *zap.SugaredLogger
}

func NewSpaceHandler(
	SpaceRepo *repositories.SpaceRepository,
	logger *zap.SugaredLogger,
) *SpaceHandler {
	return &SpaceHandler{
		SpaceRepo: SpaceRepo,
		logger:    logger,
	}
}

// CreateSpace godoc
// @Summary Create a new Space
// @Description Create a new Space with the given details
// @Tags Spaces
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
		Name:       req.Name,
		CreatedAt:  req.CreatedAt,
		ModifiedAt: req.ModifiedAt,
		UserID:     uid,
	}

	if err := h.SpaceRepo.CreateSpace(&Space); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			utils.SendErrorResponse(c, h.logger, messages.ErrUniqueConstraintFailed,
				err.Error(), apperrors.ErrSpaceDuplicateKey)
			return
		}
		utils.SendErrorResponse(c, h.logger, messages.ErrSpaceCreationFailed,
			err.Error(), apperrors.ErrInternalServerError)
		return
	}

	c.JSON(http.StatusOK, utils.CreateJSONResponse(
		messages.Success, messages.MsgSpaceCreationSuccess, Space))
}

// UpdateSpace godoc
// @Summary Update an existing Space
// @Description Update an existing Space with the given details
// @Tags Spaces
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

	space := models.Space{
		ID:         spaceID,
		Name:       req.Name,
		CreatedAt:  req.CreatedAt,
		ModifiedAt: req.ModifiedAt,
		UserID:     uid,
	}

	if err := h.SpaceRepo.UpdateSpace(&space); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SendErrorResponse(c, h.logger, messages.ErrSpaceUpdateFailed,
				"Space not found or does not belong to user", apperrors.ErrUnauthorized)
			return
		}

		utils.SendErrorResponse(c, h.logger, messages.ErrSpaceUpdateFailed,
			err.Error(), apperrors.ErrInternalServerError)
		return
	}

	c.JSON(http.StatusOK, utils.CreateJSONResponse(
		messages.Success, messages.MsgSpaceUpdateSuccess, space))
}

func (h *SpaceHandler) GetSpacesFromVersion(c *gin.Context) {
}
