package middleware

import (
	"blockstracker_backend/config"
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

func abortUnauthorized(c *gin.Context, logger *zap.SugaredLogger, errMsg string, err error) {
	logger.Errorw(errMsg, "error", err)
	c.AbortWithStatusJSON(http.StatusUnauthorized, utils.CreateJSONResponse(messages.Error, errMsg, nil))
}

func (m *AuthMiddleware) Handle(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		abortUnauthorized(c, m.logger, messages.ErrNoAuthorizationHeader, errors.New(messages.ErrUnauthorized))
		return
	}

	tokenString, err := extractBearerToken(authHeader)
	if err != nil {
		abortUnauthorized(c, m.logger, messages.ErrInvalidAuthorizationHeader, errors.New(messages.ErrUnauthorized))
		return
	}

	claims, err := parseToken(tokenString, m.authConfig)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			abortUnauthorized(c, m.logger, messages.ErrTokenExpired, err)
			return
		}
		// Optionally check for token expiration here if needed
		abortUnauthorized(c, m.logger, err.Error(), err)
		return
	}

	c.Set("userID", claims.UserID)
	c.Set("email", claims.Email)
	c.Next()

}

func extractBearerToken(header string) (string, error) {
	splitToken := strings.Split(header, " ")
	if len(splitToken) != 2 || strings.ToLower(splitToken[0]) != "bearer" {
		return "", errors.New(messages.ErrInvalidAuthorizationHeader)
	}
	return splitToken[1], nil
}

func parseToken(tokenString string, authConfig *config.AuthConfig) (*models.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(messages.ErrUnexpectedSigningMethod)
		}
		return []byte(authConfig.AccessSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New(messages.ErrInvalidToken)
}
