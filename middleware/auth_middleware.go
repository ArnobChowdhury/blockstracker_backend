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

	splitToken := strings.Split(authHeader, " ")
	if len(splitToken) != 2 || strings.ToLower(splitToken[0]) != "bearer" {
		abortUnauthorized(c, m.logger, messages.ErrInvalidAuthorizationHeader, errors.New(messages.ErrUnauthorized))
		return
	}

	tokenString := splitToken[1]

	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(messages.ErrUnexpectedSigningMethod)
		}
		return []byte(m.authConfig.AccessSecret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			abortUnauthorized(c, m.logger, messages.ErrTokenExpired, err)
			return
		}

		abortUnauthorized(c, m.logger, messages.ErrInvalidToken, err)
		return
	}

	if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)

		c.Next()
	} else {
		abortUnauthorized(c, m.logger, messages.ErrInvalidToken, errors.New(messages.ErrUnauthorized))
	}

}
