package handlers

import (
	"blockstracker_backend/config"
	apperrors "blockstracker_backend/internal/errors"
	"blockstracker_backend/internal/repositories"
	"blockstracker_backend/internal/utils"
	"blockstracker_backend/messages"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// const (
// 	EntityTypeTask                   = "task"
// 	EntityTypeTag                    = "tag"
// 	EntityTypeSpace                  = "space"
// 	EntityTypeRepetitiveTaskTemplate = "repetitive_task_template"

// 	OperationCreate = "create"
// 	OperationUpdate = "update"
// 	OperationDelete = "delete"
// )

type BillingHandler struct {
	db         *gorm.DB
	userRepo   *repositories.UserRepository
	tokenRepo  repositories.TokenRepository
	authConfig *config.AuthConfig
	logger     *zap.SugaredLogger
}

func NewBillingHandler(
	db *gorm.DB,
	userRepo *repositories.UserRepository,
	tokenRepo repositories.TokenRepository,
	authConfig *config.AuthConfig,
	logger *zap.SugaredLogger,
) *BillingHandler {
	return &BillingHandler{
		db:         db,
		userRepo:   userRepo,
		tokenRepo:  tokenRepo,
		authConfig: authConfig,
		logger:     logger,
	}
}

type VerifyGooglePurchaseRequest struct {
	PurchaseToken string `json:"purchaseToken" binding:"required"`
	ProductId     string `json:"productId" binding:"required"`
}

func (h *BillingHandler) Verify(c *gin.Context) {
	var req VerifyGooglePurchaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, h.logger, messages.ErrInvalidRequestBody, err.Error(), apperrors.ErrMalformedRequest)
		return
	}

	uid, err := utils.ExtractUIDFromGinContext(c)
	if err != nil {
		utils.SendErrorResponse(c, h.logger, "User not authenticated", err.LogError(), apperrors.ErrUnauthorized)
		return
	}
	userId := uid.String()

	// 1. Development/Mock Flow
	// This matches the "mock-purchase-token" we set in the React Native IAPService
	if req.PurchaseToken == "mock-purchase-token" {
		// Directly update the user status for testing
		// For mock purchases, we grant 1 year of premium
		expiry := time.Now().AddDate(1, 0, 0)
		if err := h.userRepo.UpdatePremiumExpiry(userId, expiry); err != nil {
			utils.SendErrorResponse(c, h.logger, "Failed to update user premium status", err.Error(), apperrors.ErrInternalServerError)
			return
		}

		// 2. Fetch the updated user to generate new claims
		user, err := h.userRepo.GetUserByID(userId)
		if err != nil {
			utils.SendErrorResponse(c, h.logger, "Failed to fetch updated user", err.Error(), apperrors.ErrInternalServerError)
			return
		}

		// 3. Generate new tokens with is_premium = true
		accessTokenClaims := utils.GetClaims(user, "access")
		refreshTokenClaims := utils.GetClaims(user, "refresh")

		accessToken, err := utils.GenerateJWT(accessTokenClaims, h.authConfig.AccessSecret)
		if err != nil {
			utils.SendErrorResponse(c, h.logger, messages.ErrGeneratingJWT, err.Error(), apperrors.ErrInternalServerError)
			return
		}

		refreshToken, err := utils.GenerateJWT(refreshTokenClaims, h.authConfig.RefreshSecret)
		if err != nil {
			utils.SendErrorResponse(c, h.logger, messages.ErrGeneratingJWT, err.Error(), apperrors.ErrInternalServerError)
			return
		}

		// 4. Store new tokens in Redis
		if err := h.tokenRepo.StoreAccessTokenAndRefreshToken(accessToken, refreshToken); err != nil {
			utils.SendErrorResponse(c, h.logger, apperrors.ErrRedisSet.LogError(), err.Error(), apperrors.ErrInternalServerError)
			return
		}

		c.JSON(http.StatusOK, utils.CreateJSONResponse(messages.Success, "Mock purchase verified successfully", gin.H{
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
		}))
		return
	}

	// 2. Production Flow (TODO)
	// Here you would use the Google Play Developer API to verify the token with Google.
	c.JSON(http.StatusNotImplemented, utils.CreateJSONResponse(messages.Error, "Real Google Play verification is not yet configured. Please use the mock flow.", nil))
}
