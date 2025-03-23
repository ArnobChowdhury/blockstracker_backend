package handlers

import (
	"net/http"

	apperrors "blockstracker_backend/internal/errors"
	"blockstracker_backend/internal/repositories"
	"blockstracker_backend/internal/utils"
	"blockstracker_backend/messages"
	"blockstracker_backend/models"

	"github.com/gin-gonic/gin"
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

// CreateTag godoc
// @Summary Create a new tag
// @Description Create a new tag with the given details
// @Tags tags
// @Accept json
// @Produce json
// @Param tag body models.CreateTagRequest true "Tag details"
// @Success 200 {object} models.CreateTagResponseForSwagger
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

	var req models.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		invalidReqErr := apperrors.NewInvalidReqErr(err.Error())
		utils.SendErrorResponse(c, h.logger, messages.ErrTagCreationFailed,
			err.Error(), invalidReqErr)
		return
	}

	tag := models.Tag{
		Name:       req.Name,
		CreatedAt:  req.CreatedAt,
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
