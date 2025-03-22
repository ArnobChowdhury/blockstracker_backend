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

type TagHandler struct {
	tagRepo *repositories.TagRepository
	logger  *zap.SugaredLogger
}

func NewTagHandler(
	tagRepo *repositories.TagRepository,
	logger *zap.SugaredLogger,
) *TagHandler {
	return &TagHandler{
		tagRepo: tagRepo,
		logger:  logger,
	}
}

func (h *TagHandler) CreateTag(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		utils.SendErrorResponse(c, h.logger, messages.ErrTagCreationFailed,
			"User ID not found in context", apperrors.ErrInternalServerError)
		return
	}

	uid, ok := userID.(uuid.UUID)
	if !ok {
		utils.SendErrorResponse(c, h.logger, messages.ErrTagCreationFailed,
			"User ID is not of valid type", apperrors.ErrInternalServerError)
		return
	}

	var req models.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		invalidReqErr := apperrors.NewInvalidReqErr(err.Error())
		utils.SendErrorResponse(c, h.logger, messages.ErrTagCreationFailed,
			err.Error(), invalidReqErr)
		return
	}

	tag := models.Tag{
		Name:       req.Name,
		ModifiedAt: req.ModifiedAt,
		UserID:     uid,
	}

	if err := h.tagRepo.CreateTag(&tag); err != nil {
		utils.SendErrorResponse(c, h.logger, messages.ErrTagCreationFailed,
			err.Error(), apperrors.ErrInternalServerError)
		return
	}

	c.JSON(http.StatusOK, utils.CreateJSONResponse(
		messages.Success, messages.MsgTagCreationSuccess, tag))
}

func (h *TagHandler) UpdateTag(c *gin.Context) {
}

func (h *TagHandler) GetTagsFromVersion(c *gin.Context) {
}
