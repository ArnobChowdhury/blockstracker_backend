package handlers

import (
	"blockstracker_backend/config"
	apperrors "blockstracker_backend/internal/errors"
	"blockstracker_backend/internal/repositories"
	"blockstracker_backend/internal/utils"
	"blockstracker_backend/messages"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/api/androidpublisher/v3"
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

	var expiry time.Time

	// 1. Development/Mock Flow // TODO: Remove before production
	// This matches the "mock-purchase-token" we set in the React Native IAPService
	if req.PurchaseToken == "mock-purchase-token" {
		// Directly update the user status for testing
		// For mock purchases, we grant 1 year of premium
		expiry = time.Now().AddDate(1, 0, 0)
	} else {
		// 2. Production Flow
		// Requires GOOGLE_APPLICATION_CREDENTIALS env var to be set to the path of the service account JSON file
		ctx := context.Background()
		service, err := androidpublisher.NewService(ctx)
		if err != nil {
			utils.SendErrorResponse(c, h.logger, "Failed to initialize Google Play service", err.Error(), apperrors.ErrInternalServerError)
			return
		}

		// TODO: Move package name to config
		packageName := "com.blockstrackerapp"

		sub, err := service.Purchases.Subscriptions.Get(packageName, req.ProductId, req.PurchaseToken).Do()
		if err != nil {
			h.logger.Errorw("Google Play API Error", "error", err)
			utils.SendErrorResponse(c, h.logger, "Failed to verify subscription with Google", err.Error(), apperrors.ErrInvalidCredentials)
			return
		}

		expiry = time.Unix(sub.ExpiryTimeMillis/1000, 0)
		if time.Now().After(expiry) {
			utils.SendErrorResponse(c, h.logger, "Subscription expired", "Subscription expired", apperrors.ErrUnauthorized)
			return
		}
	}

	if err := h.userRepo.UpdatePremiumExpiry(userId, expiry); err != nil {
		utils.SendErrorResponse(c, h.logger, "Failed to update user premium status", err.Error(), apperrors.ErrInternalServerError)
		return
	}

	user, userGetErr := h.userRepo.GetUserByID(userId)
	if userGetErr != nil {
		utils.SendErrorResponse(c, h.logger, "Failed to fetch updated user", userGetErr.Error(), apperrors.ErrInternalServerError)
		return
	}

	accessTokenClaims := utils.GetClaims(user, "access")
	refreshTokenClaims := utils.GetClaims(user, "refresh")

	accessToken, jwtGenErr := utils.GenerateJWT(accessTokenClaims, h.authConfig.AccessSecret)
	if jwtGenErr != nil {
		utils.SendErrorResponse(c, h.logger, messages.ErrGeneratingJWT, jwtGenErr.Error(), apperrors.ErrInternalServerError)
		return
	}

	refreshToken, jwtGenErr := utils.GenerateJWT(refreshTokenClaims, h.authConfig.RefreshSecret)
	if jwtGenErr != nil {
		utils.SendErrorResponse(c, h.logger, messages.ErrGeneratingJWT, jwtGenErr.Error(), apperrors.ErrInternalServerError)
		return
	}

	if err := h.tokenRepo.StoreAccessTokenAndRefreshToken(accessToken, refreshToken); err != nil {
		utils.SendErrorResponse(c, h.logger, apperrors.ErrRedisSet.LogError(), err.Error(), apperrors.ErrInternalServerError)
		return
	}

	c.JSON(http.StatusOK, utils.CreateJSONResponse(messages.Success, "Purchase verified successfully", gin.H{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	}))
}
