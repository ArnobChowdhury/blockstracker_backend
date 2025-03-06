package middleware

import (
	"blockstracker_backend/config"
	apperrors "blockstracker_backend/internal/errors"
	"blockstracker_backend/internal/utils"
	"blockstracker_backend/messages"
	"blockstracker_backend/models"
	"errors"
	"net/http"
	"strings"

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

func (m *AuthMiddleware) abortUnauthorized(c *gin.Context, logTitle string, logErrMsg string, resErrMsg string) {
	m.logger.Errorw(logTitle, messages.Error, logErrMsg)
	c.AbortWithStatusJSON(http.StatusUnauthorized,
		utils.CreateJSONResponse(messages.Error, resErrMsg, nil))
}

func (m *AuthMiddleware) Handle(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		m.abortUnauthorized(c, "No Auth Header",
			apperrors.ErrNoAuthorizationHeader.LogError(),
			apperrors.ErrUnauthorized.Error())
		return
	}

	tokenString, err := extractBearerToken(authHeader)
	if err != nil {
		m.abortUnauthorized(c, "Invalid Auth Header",
			err.LogError(), apperrors.ErrUnauthorized.Error())
		return
	}

	claims, parseErr := parseToken(tokenString, m.authConfig)
	if parseErr != nil {
		if errors.Is(parseErr, jwt.ErrTokenExpired) {
			m.abortUnauthorized(c, "JWT expired",
				parseErr.Error(),
				apperrors.ErrTokenExpired.Error())
			return
		}

		if authError, ok := parseErr.(*apperrors.AuthError); ok {
			m.abortUnauthorized(c, "JWT parsing error",
				authError.LogError(),
				apperrors.ErrTokenExpired.Error())
			return
		}

		m.abortUnauthorized(c, "JWT parsing error", err.Error(), apperrors.ErrUnauthorized.Error())
		return
	}

	c.Set("userID", claims.UserID)
	c.Set("email", claims.Email)
	c.Next()
}

func extractBearerToken(header string) (string, *apperrors.AuthError) {
	splitToken := strings.Split(header, " ")
	if len(splitToken) != 2 || strings.ToLower(splitToken[0]) != "bearer" {
		return "", apperrors.ErrNoAuthorizationHeader
	}
	return splitToken[1], nil
}

func parseToken(tokenString string, authConfig *config.AuthConfig) (*models.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apperrors.ErrUnexpectedSigningMethod
		}
		return []byte(authConfig.AccessSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, apperrors.ErrInvalidToken
}
