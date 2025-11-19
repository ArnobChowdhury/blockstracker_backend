package handlers

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"

	_ "blockstracker_backend/docs"
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
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"
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

// SignupUser godoc
// @Summary      Sign up a new user
// @Description  Signs up a new user with email and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body models.SignUpRequest true "User sign up request"
// @Success      200  {object}  models.GenericSuccessResponse "User creation successful"
// @Failure      400  {object}  models.GenericErrorResponse "Malformed Request"
// @Failure      500  {object}  models.GenericErrorResponse "Internal Server Error"
// @Router       /auth/signup [post]
func (h *AuthHandler) SignupUser(c *gin.Context) {
	var req models.SignUpRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorw(messages.ErrInvalidRequestBody, messages.Error, err)

		if validationErrors, ok := err.(validator.ValidationErrors); ok && len(validationErrors) > 0 {
			msg := validators.GetCustomMessage(validationErrors[0], req)

			authError := apperrors.NewInvalidReqErr(msg)
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

// EmailSignIn godoc
// @Summary      Sign in with email and password
// @Description  Sign in with email and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body models.EmailSignInRequest true "User sign in request"
// @Success      200  {object}  models.SignInSuccessResponse "User sign in successful"
// @Failure      400  {object}  models.GenericErrorResponse  "Malformed Request"
// @Failure      401  {object}  models.GenericErrorResponse  "Invalid Credentials"
// @Failure      500  {object}  models.GenericErrorResponse "Internal Server Error"
// @Router       /auth/signin [post]
func (h *AuthHandler) EmailSignIn(c *gin.Context) {
	var req models.EmailSignInRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, h.logger, messages.ErrInvalidRequestBody,
			err.Error(), apperrors.ErrMalformedTaskRequest)
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

// @Summary      Sign out user
// @Description  Invalidates the user's access and refresh tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.GenericSuccessResponse "User sign out successful"
// @Failure      500  {object}  models.GenericErrorResponse "Internal Server Error"
// @Router       /auth/signout [post]
func (h *AuthHandler) Signout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	token, _ := utils.ExtractBearerToken(authHeader)

	deletedCount, err := h.tokenRepo.InvalidateAccessAndRefreshTokens(token)
	if err != nil {
		utils.SendErrorResponse(c, h.logger, apperrors.ErrRedisKeyNotFound.Error(),
			err.Error(), apperrors.ErrInternalServerError)
		return
	}
	if deletedCount == 0 {
		utils.SendErrorResponse(c, h.logger, messages.ErrTokenNotFoundDuringLogout, "", apperrors.ErrUnauthorized)
		return
	}

	h.logger.Infow(messages.MsgSignOutSuccessful, "token", token)
	c.JSON(http.StatusOK, utils.CreateJSONResponse(messages.Success, messages.MsgSignOutSuccessful, nil))
}

// RefreshToken godoc
// @Summary      Refresh access token
// @Description  Refreshes the access token using a valid refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body models.RefreshTokenRequest true "Refresh token request"
// @Success      200  {object}  models.SignInSuccessResponse "Token refresh successful"
// @Failure      400  {object}  models.GenericErrorResponse "Malformed Request"
// @Failure      401  {object}  models.GenericErrorResponse "Unauthorized"
// @Failure      500  {object}  models.GenericErrorResponse "Internal Server Error"
// @Router       /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, h.logger, messages.ErrInvalidRequestBody, err.Error(), apperrors.ErrMalformedRequest)
		return
	}

	refreshTokenInRedis, err := h.tokenRepo.GetRefreshToken(req.AccessToken)
	if err != nil {
		h.logger.Errorw(messages.ErrGettingRefreshTokenFromRedis, err)
	} else {
		hashedRefreshToken := sha256.Sum256([]byte(req.RefreshToken))
		hashedRefreshTokenString := hex.EncodeToString(hashedRefreshToken[:])
		if hashedRefreshTokenString != refreshTokenInRedis {
			utils.SendErrorResponse(c, h.logger, messages.ErrRefreshTokenDidNotMatchWithCachedToken,
				"refresh token mismatch", apperrors.ErrUnauthorized)
			return
		}
	}

	parsedClaims, err := utils.ParseToken(req.RefreshToken, h.authConfig.RefreshSecret)
	if err != nil {
		utils.SendErrorResponse(c, h.logger, messages.ErrJWTParsingError,
			err.Error(), apperrors.ErrUnauthorized)
		return
	}

	user := &models.User{
		ID:    parsedClaims.UserID,
		Email: parsedClaims.Email,
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

	deletedCount, err := h.tokenRepo.InvalidateAccessAndRefreshTokens(req.AccessToken)
	if err != nil {
		h.logger.Errorw(messages.ErrInvalidatingOldTokens, messages.Error, err)
	}
	if deletedCount == 0 {
		h.logger.Errorw(messages.ErrTokenNotFoundDuringRefreshing, "token", req.AccessToken)
	}

	c.JSON(http.StatusOK, utils.CreateJSONResponse(messages.Success, messages.MsgSuccessfulTokenRefresh, gin.H{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	}))
}

func (h *AuthHandler) GoogleSignInMobile(c *gin.Context) {
	var req models.GoogleSignInMobileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, h.logger, messages.ErrInvalidRequestBody, err.Error(), apperrors.ErrMalformedRequest)
		return
	}

	payload, err := idtoken.Validate(context.Background(), req.Token, h.authConfig.GoogleWebClientID)
	if err != nil {
		utils.SendErrorResponse(c, h.logger, "Failed to verify Google ID token",
			err.Error(), apperrors.ErrUnauthorized)
		return
	}

	email := payload.Claims["email"].(string)
	user, err := h.userRepo.GetUserByEmail(email)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			provider := "google"
			newUser := models.User{
				Email:    email,
				Provider: &provider,
			}
			if creationErr := h.userRepo.CreateUser(&newUser); creationErr != nil {
				utils.SendErrorResponse(c, h.logger, messages.ErrUnexpectedErrorDuringUserCreation,
					creationErr.Error(), apperrors.ErrInternalServerError)
				return
			}
			user = &newUser
		} else {
			utils.SendErrorResponse(c, h.logger, messages.ErrUnexpectedErrorDuringUserRetrieval,
				err.Error(), apperrors.ErrInternalServerError)
			return
		}
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

func (h *AuthHandler) GoogleSignInDesktop(c *gin.Context) {
	var req models.GoogleSignInDesktopRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, h.logger, messages.ErrInvalidRequestBody, err.Error(), apperrors.ErrMalformedRequest)
		return
	}

	oauth2Config := &oauth2.Config{
		ClientID:     h.authConfig.GoogleWebClientID,
		ClientSecret: h.authConfig.GoogleWebClientSecret,
		RedirectURL:  req.RedirectURI,
		Endpoint:     google.Endpoint,
		Scopes: []string{
			"openid",
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
	}

	token, err := oauth2Config.Exchange(
		context.Background(),
		req.Code,
		oauth2.SetAuthURLParam("code_verifier", req.CodeVerifier))
	if err != nil {
		utils.SendErrorResponse(c, h.logger, "Failed to exchange authorization code for tokens",
			err.Error(), apperrors.ErrUnauthorized)
		return
	}

	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		utils.SendErrorResponse(c, h.logger, "ID token not found in Google response",
			"id_token missing", apperrors.ErrUnauthorized)
		return
	}

	payload, err := idtoken.Validate(context.Background(), idToken, h.authConfig.GoogleWebClientID)
	if err != nil {
		utils.SendErrorResponse(c, h.logger, "Failed to verify Google ID token",
			err.Error(), apperrors.ErrUnauthorized)
		return
	}

	email := payload.Claims["email"].(string)
	user, err := h.userRepo.GetUserByEmail(email)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			provider := "google"
			newUser := models.User{
				Email:    email,
				Provider: &provider,
			}
			if creationErr := h.userRepo.CreateUser(&newUser); creationErr != nil {
				utils.SendErrorResponse(c, h.logger, messages.ErrUnexpectedErrorDuringUserCreation,
					creationErr.Error(), apperrors.ErrInternalServerError)
				return
			}
			user = &newUser
		} else {
			utils.SendErrorResponse(c, h.logger, messages.ErrUnexpectedErrorDuringUserRetrieval,
				err.Error(), apperrors.ErrInternalServerError)
			return
		}
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
