package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	apperrors "blockstracker_backend/internal/errors"
	"blockstracker_backend/internal/repositories"
	messages "blockstracker_backend/messages"

	"blockstracker_backend/config"
	"blockstracker_backend/internal/utils"
	"blockstracker_backend/internal/validators"
	"blockstracker_backend/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthHandler struct {
	userRepo    *repositories.UserRepository
	logger      *zap.SugaredLogger
	authConfig  *config.AuthConfig
	redisClient *redis.Client
}

func NewAuthHandler(
	userRepo *repositories.UserRepository,
	logger *zap.SugaredLogger,
	authConfig *config.AuthConfig,
	redisClient *redis.Client) *AuthHandler {

	return &AuthHandler{
		userRepo:    userRepo,
		logger:      logger,
		authConfig:  authConfig,
		redisClient: redisClient,
	}
}

func (h *AuthHandler) respondWithError(c *gin.Context, status int, logMsg string, err error, clientMsg string) {
	h.logger.Errorw(logMsg, messages.Error, err)
	c.JSON(status, utils.CreateJSONResponse(messages.Error, clientMsg, nil))
}

func (h *AuthHandler) SignupUser(c *gin.Context) {
	var req models.SignUpRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw(messages.ErrInvalidRequestBody, messages.Error, err)

		if validationErrors, ok := err.(validator.ValidationErrors); ok && len(validationErrors) > 0 {
			msg := validators.GetCustomMessage(validationErrors[0], req)
			c.JSON(http.StatusBadRequest,
				utils.CreateJSONResponse(messages.Error, msg, nil))
			return
		}

		c.JSON(http.StatusBadRequest,
			utils.CreateJSONResponse(messages.Error, messages.ErrMalformedRequest, nil))
		return
	}

	hashedPassword, pwHashingError := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if pwHashingError != nil {
		h.respondWithError(c, http.StatusInternalServerError,
			messages.ErrHashingPassword, pwHashingError, messages.ErrInternalServerError)
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
			c.JSON(http.StatusBadRequest,
				utils.CreateJSONResponse(messages.Error, messages.ErrUserCreationFailed, nil))
			return

		} else {
			h.logger.Errorw(messages.ErrUnexpectedErrorDuringUserCreation, "email", user.Email, messages.Error, err)
			c.JSON(http.StatusInternalServerError,
				utils.CreateJSONResponse(messages.Error, messages.ErrInternalServerError, nil))
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
		h.respondWithError(c, http.StatusBadRequest,
			messages.ErrInvalidRequestBody, err, messages.ErrMalformedRequest)
		return
	}

	user, err := h.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.logger.Errorw(messages.ErrEmailNotFoundDuringSignIn, "email", req.Email)
			c.JSON(http.StatusUnauthorized,
				utils.CreateJSONResponse(messages.Error, messages.ErrInvalidCredentials, nil))
			return
		} else {
			h.respondWithError(c, http.StatusInternalServerError,
				messages.ErrUnexpectedErrorDuringUserRetrieval, err, messages.ErrInternalServerError)
			return
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(req.Password))
	if err != nil {
		h.respondWithError(c, http.StatusUnauthorized,
			messages.ErrMismatchingPasswordDuringSignIn, err, messages.ErrInvalidCredentials)
		return
	}

	accessTokenClaims := utils.GetClaims(user, "access")
	refreshTokenClaims := utils.GetClaims(user, "refresh")

	accessToken, err := utils.GenerateJWT(accessTokenClaims, h.authConfig.AccessSecret)
	if err != nil {
		h.respondWithError(c, http.StatusInternalServerError, messages.ErrGeneratingJWT, err, messages.ErrInternalServerError)
		return
	}

	refreshToken, err := utils.GenerateJWT(refreshTokenClaims, h.authConfig.RefreshSecret)
	if err != nil {
		h.respondWithError(c, http.StatusInternalServerError, messages.ErrGeneratingJWT, err, messages.ErrInternalServerError)
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

	invalidateAccessTokenCtx, invalidateAccessTokenCtxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer invalidateAccessTokenCtxCancel()
	// for simplicity, we are not calculating the expiry time, and using the default expirty time we used to generate the token
	err := h.redisClient.Set(invalidateAccessTokenCtx, "accessToken:"+token, "invalid", utils.AccessTokenExpiry).Err()
	if err != nil {
		h.respondWithError(c, http.StatusInternalServerError,
			apperrors.ErrRedisSet.Error(),
			err, messages.ErrInternalServerError)
		return
	}

	refreshToken, err := h.redisClient.Get(context.Background(), "accessToRefresh:"+token).Result()
	if err != nil {
		h.respondWithError(c, http.StatusInternalServerError,
			apperrors.ErrRedisGet.Error(),
			err, messages.ErrInternalServerError)
		return
	}

	invalidateRefreshTokenCtx, invalidateRefreshTokenCtxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer invalidateRefreshTokenCtxCancel()
	err = h.redisClient.Set(invalidateRefreshTokenCtx, "refreshToken:"+refreshToken, "invalid", utils.RefreshTokenExpiry).Err()
	if err != nil {
		h.respondWithError(c, http.StatusInternalServerError,
			apperrors.ErrRedisSet.Error(),
			err, messages.ErrInternalServerError)
		return
	}

	h.logger.Infow(messages.MsgSignOutSuccessful, "token", token)
	c.JSON(http.StatusOK, utils.CreateJSONResponse(messages.Success, messages.MsgSignOutSuccessful, nil))
}
