package middleware

import (
	"blockstracker_backend/config"
	apperrors "blockstracker_backend/internal/errors"
	"blockstracker_backend/internal/utils"
	"blockstracker_backend/messages"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type AuthMiddleware struct {
	logger     *zap.SugaredLogger
	authConfig *config.AuthConfig
}

func NewAuthMiddleware(logger *zap.SugaredLogger, authConfig *config.AuthConfig) *AuthMiddleware {
	return &AuthMiddleware{
		logger:     logger,
		authConfig: authConfig,
	}
}

func (m *AuthMiddleware) mapAuthError(err error) (logTitle, logErrMsg string, resErr apperrors.AppError) {
	switch {
	case errors.Is(err, jwt.ErrTokenExpired):
		return messages.ErrJWTExpired, err.Error(), apperrors.ErrTokenExpired
	case func() bool { _, ok := err.(*apperrors.AuthError); return ok }():
		return messages.ErrJWTParsingError, err.Error(), apperrors.ErrUnauthorized
	default:
		return messages.ErrJWTParsingError, err.Error(), apperrors.ErrUnauthorized
	}
}

func (m *AuthMiddleware) Handle(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		utils.SendErrorResponse(c, m.logger, "No Auth Header",
			apperrors.ErrNoAuthorizationHeader.LogError(),
			apperrors.ErrUnauthorized)
		c.Abort()
		return
	}

	tokenString, err := utils.ExtractBearerToken(authHeader)
	if err != nil {
		utils.SendErrorResponse(c, m.logger, "Invalid Auth Header",
			err.LogError(), apperrors.ErrUnauthorized)
		c.Abort()
		return
	}

	claims, parseErr := utils.ParseToken(tokenString, m.authConfig.AccessSecret)
	if parseErr != nil {
		logTitle, logErrMsg, resErr := m.mapAuthError(parseErr)
		utils.SendErrorResponse(c, m.logger, logTitle, logErrMsg, resErr)
		c.Abort()
		return
	}

	c.Set("userID", claims.UserID)
	c.Set("email", claims.Email)
	c.Set("is_premium", claims.IsPremium)
	c.Next()
}

func (m *AuthMiddleware) RequirePremium(c *gin.Context) {
	isPremium := c.GetBool("is_premium")
	if !isPremium {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "error",
			"message": "This feature requires a premium subscription.",
			"code":    "PREMIUM_REQUIRED",
		})
		c.Abort()
		return
	}
	c.Next()
}
