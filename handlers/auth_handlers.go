package handlers

import (
	"errors"
	"net/http"

	"blockstracker_backend/internal/repositories"
	messages "blockstracker_backend/messages"

	"blockstracker_backend/internal/utils"
	"blockstracker_backend/internal/validators"
	"blockstracker_backend/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	userRepo *repositories.UserRepository
	logger   *zap.SugaredLogger
}

func NewAuthHandler(userRepo *repositories.UserRepository, logger *zap.SugaredLogger) *AuthHandler {
	return &AuthHandler{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (h *AuthHandler) SignupUser(c *gin.Context) {
	var req models.SignUpRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw(messages.ErrInvalidRequestBody, messages.Error, err)

		if validationErrors, ok := err.(validator.ValidationErrors); ok && len(validationErrors) > 0 {
			msg := validators.GetCustomMessage(validationErrors[0], req)
			c.JSON(http.StatusBadRequest, utils.CreateJSONResponse(messages.Error, msg))
			return
		}

		c.JSON(http.StatusBadRequest, utils.CreateJSONResponse(messages.Error, messages.ErrMalformedRequest))
		return
	}

	hashedPassword, pwHashingError := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if pwHashingError != nil {
		h.logger.Errorw(messages.ErrHashingPassword, messages.Error, pwHashingError)

		c.JSON(http.StatusInternalServerError, utils.CreateJSONResponse(messages.Error, messages.ErrHashingPassword))
		return
	}

	hashedPasswordString := string(hashedPassword)

	user := models.User{
		Email:    req.Email,
		Password: &hashedPasswordString,
	}

	if err := h.userRepo.CreateUser(&user); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			h.logger.Errorw(messages.ErrUniqueConstraintFailed, "email", user.Email)
			c.JSON(http.StatusBadRequest, utils.CreateJSONResponse(messages.Error, messages.ErrUserCreationFailed))
			return

		} else {
			h.logger.Errorw(messages.ErrUnexpectedErrorDuringUserCreation, "email", user.Email, messages.Error, err)
			c.JSON(http.StatusInternalServerError, utils.CreateJSONResponse(messages.Error, messages.ErrInternalServerError))
			return
		}
	}

	h.logger.Infow(messages.MsgUserCreationSuccess, "email", user.Email)
	c.JSON(http.StatusOK, utils.CreateJSONResponse(messages.Success, messages.MsgUserCreationSuccess))
}
