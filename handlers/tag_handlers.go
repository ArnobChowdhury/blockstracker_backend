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
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
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
			"Invalid tag ID", apperrors.NewInvalidReqErr("Invalid tag ID"))
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

	if err := h.tagRepo.UpdateTag(&tag); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SendErrorResponse(c, h.logger, messages.ErrTagUpdateFailed,
				"Tag not found or does not belong to user", apperrors.ErrUnauthorized)
			return
		}

		utils.SendErrorResponse(c, h.logger, messages.ErrTagUpdateFailed,
			err.Error(), apperrors.ErrInternalServerError)
		return
	}

	c.JSON(http.StatusOK, utils.CreateJSONResponse(
		messages.Success, messages.MsgTagUpdateSuccess, tag))
}

func (h *TagHandler) GetTagsFromVersion(c *gin.Context) {
}
