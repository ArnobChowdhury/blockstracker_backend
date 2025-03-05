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

func NewAuthMiddleware(logger *zap.SugaredLogger, authConfig *config.AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			logger.Errorw(messages.ErrNoAuthorizationHeader)
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				utils.CreateJSONResponse(messages.Error, messages.ErrUnauthorized, nil))
			return
		}

		splitToken := strings.Split(authHeader, " ")
		if len(splitToken) != 2 || strings.ToLower(splitToken[0]) != "bearer" {
			logger.Errorw(messages.ErrInvalidAuthorizationHeader)
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				utils.CreateJSONResponse(messages.Error, messages.ErrUnauthorized, nil))
			return
		}

		tokenString := splitToken[1]

		token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New(messages.ErrUnexpectedSigningMethod)
			}
			return []byte(authConfig.AccessSecret), nil
		})

		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				logger.Errorw(messages.ErrTokenExpired, messages.Error, err)
				c.AbortWithStatusJSON(http.StatusUnauthorized,
					utils.CreateJSONResponse(messages.Error, messages.ErrTokenExpired, nil))
				return
			}

			logger.Errorw(messages.ErrInvalidToken, messages.Error, err)
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				utils.CreateJSONResponse(messages.Error, messages.ErrUnauthorized, nil))
			return
		}

		if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
			c.Set("userID", claims.UserID)
			c.Set("email", claims.Email)

			c.Next()
		} else {
			logger.Errorw(messages.ErrInvalidToken)
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				utils.CreateJSONResponse(messages.Error, messages.ErrUnauthorized, nil))
		}

	}

}
