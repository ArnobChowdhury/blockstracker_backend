package handlers

import (
	"errors"
	"net/http"

	apperrors "blockstracker_backend/internal/errors"
	"blockstracker_backend/internal/repositories"
	messages "blockstracker_backend/messages"

	"blockstracker_backend/config"
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
	userRepo   *repositories.UserRepository
	logger     *zap.SugaredLogger
	authConfig *config.AuthConfig
	tokenRepo  repositories.TokenRepository
}

func NewAuthHandler(
	userRepo *repositories.UserRepository,
	logger *zap.SugaredLogger,
	authConfig *config.AuthConfig,
	tokenRepo repositories.TokenRepository,
) *AuthHandler {

	return &AuthHandler{
		userRepo:   userRepo,
		logger:     logger,
		authConfig: authConfig,
		tokenRepo:  tokenRepo,
	}
}

func (h *AuthHandler) SignupUser(c *gin.Context) {
	var req models.SignUpRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw(messages.ErrInvalidRequestBody, messages.Error, err)

		if validationErrors, ok := err.(validator.ValidationErrors); ok && len(validationErrors) > 0 {
			msg := validators.GetCustomMessage(validationErrors[0], req)

			authError := apperrors.ErrInvalidRequestBody
			authError.SetErrMessage(msg)
			utils.SendErrorResponse(c, h.logger, messages.ErrInvalidRequestBody,
				authError.LogError(), authError)
			return
		}

		utils.SendErrorResponse(c, h.logger, messages.ErrMalformedRequest,
			apperrors.ErrMalformedRequest.LogError(), apperrors.ErrMalformedRequest)
		return
	}

	hashedPassword, pwHashingError := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if pwHashingError != nil {
		utils.SendErrorResponse(c, h.logger, messages.ErrHashingPassword,
			pwHashingError.Error(), apperrors.ErrInternalServerError)
		return
	}

	hashedPasswordString := string(hashedPassword)

	user := models.User{
		Email:    req.Email,
		Password: &hashedPasswordString,
	}

	if err := h.userRepo.CreateUser(&user); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			utils.SendErrorResponse(c, h.logger, messages.ErrUniqueConstraintFailed,
				user.Email, apperrors.ErrUserCreationFailed)
			return

		} else {
			utils.SendErrorResponse(c, h.logger, messages.ErrUnexpectedErrorDuringUserCreation,
				err.Error(), apperrors.ErrInternalServerError)
			return
		}
	}

	h.logger.Infow(messages.MsgUserCreationSuccess, "email", user.Email)
	c.JSON(http.StatusOK,
		utils.CreateJSONResponse(messages.Success, messages.MsgUserCreationSuccess, nil))
}

func (h *AuthHandler) EmailSignIn(c *gin.Context) {
	var req models.EmailSignInRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, h.logger, messages.ErrInvalidRequestBody,
			err.Error(), apperrors.ErrMalformedRequest)
		return
	}

	user, err := h.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SendErrorResponse(c, h.logger, messages.ErrEmailNotFoundDuringSignIn,
				req.Email, apperrors.ErrInvalidCredentials)
			return
		} else {
			utils.SendErrorResponse(c, h.logger, messages.ErrUnexpectedErrorDuringUserRetrieval,
				err.Error(), apperrors.ErrInternalServerError)
			return
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(req.Password))
	if err != nil {
		utils.SendErrorResponse(c, h.logger, messages.ErrMismatchingPasswordDuringSignIn,
			err.Error(), apperrors.ErrInvalidCredentials)
		return
	}

	accessTokenClaims := utils.GetClaims(user, "access")
	refreshTokenClaims := utils.GetClaims(user, "refresh")

	accessToken, err := utils.GenerateJWT(accessTokenClaims, h.authConfig.AccessSecret)
	if err != nil {
		utils.SendErrorResponse(c, h.logger, messages.ErrGeneratingJWT,
			err.Error(), apperrors.ErrInternalServerError)
		return
	}

	refreshToken, err := utils.GenerateJWT(refreshTokenClaims, h.authConfig.RefreshSecret)
	if err != nil {
		utils.SendErrorResponse(c, h.logger, messages.ErrGeneratingJWT,
			err.Error(), apperrors.ErrInternalServerError)
		return
	}

	if err := h.tokenRepo.StoreAccessTokenAndRefreshToken(accessToken, refreshToken); err != nil {
		utils.SendErrorResponse(c, h.logger, apperrors.ErrRedisSet.LogError(),
			err.Error(), apperrors.ErrInternalServerError)
		return
	}

	c.JSON(http.StatusOK,
		utils.CreateJSONResponse(messages.Success, messages.MsgSignInSuccessful, gin.H{
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
		}))
}

func (h *AuthHandler) Signout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	token, _ := utils.ExtractBearerToken(authHeader)

	if err := h.tokenRepo.InvalidateAccessAndRefreshTokens(token); err != nil {
		if redisErr, ok := err.(*apperrors.RedisError); ok {
			h.logger.Infow("Redis key not found", "token", token, messages.Error, redisErr.LogError())
			c.JSON(http.StatusOK, utils.CreateJSONResponse(messages.Success, messages.MsgSignOutSuccessful, nil))
			return
		}

		utils.SendErrorResponse(c, h.logger, apperrors.ErrRedisKeyNotFound.Error(),
			err.Error(), apperrors.ErrInternalServerError)
		return
	}

	h.logger.Infow(messages.MsgSignOutSuccessful, "token", token)
	c.JSON(http.StatusOK, utils.CreateJSONResponse(messages.Success, messages.MsgSignOutSuccessful, nil))
}
