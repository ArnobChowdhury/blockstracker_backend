package handlers

import (
	"errors"
	"net/http"

	apperrors "blockstracker_backend/internal/errors"
	"blockstracker_backend/internal/repositories"
	"blockstracker_backend/internal/utils"
	"blockstracker_backend/messages"
	"blockstracker_backend/models"

	"github.com/gin-gonic/gin"
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
// @Param Space body models.CreateSpaceRequest true "Space details"
// @Success 200 {object} models.CreateSpaceResponseForSwagger
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

	var req models.CreateSpaceRequest
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

func (h *SpaceHandler) UpdateSpace(c *gin.Context) {
}

func (h *SpaceHandler) GetSpacesFromVersion(c *gin.Context) {
}
