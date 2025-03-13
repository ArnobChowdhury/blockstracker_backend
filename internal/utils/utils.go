package utils

import (
	"blockstracker_backend/config"
	apperrors "blockstracker_backend/internal/errors"
	"blockstracker_backend/messages"
	"blockstracker_backend/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

const (
	Issuer             = "api.blocks-tracker.com"
	AccessTokenExpiry  = 30 * time.Minute
	RefreshTokenExpiry = 7 * 24 * time.Hour
)

func CreateJSONResponse(status string, message string, data interface{}) gin.H {
	result := gin.H{
		"status":  status,
		"message": message,
	}

	if data != nil {
		result["data"] = data
	}

	return gin.H{"result": result}
}

func GenerateJWT(claims models.Claims, secretKey string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func GetClaims(user *models.User, tokenType string) models.Claims {
	expiresAt := time.Now().Add(AccessTokenExpiry)

	if tokenType == "refresh" {
		expiresAt = time.Now().Add(RefreshTokenExpiry)
	}

	claims := models.Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    Issuer,
		},
	}

	return claims
}

func ExtractBearerToken(header string) (string, *apperrors.AuthError) {
	splitToken := strings.Split(header, " ")
	if len(splitToken) != 2 || strings.ToLower(splitToken[0]) != "bearer" {
		return "", apperrors.ErrNoAuthorizationHeader
	}
	return splitToken[1], nil
}

func ParseToken(tokenString string, authConfig *config.AuthConfig) (*models.Claims, error) {
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

func SendErrorResponse(c *gin.Context, logger *zap.SugaredLogger, logTitle string,
	logErrMsg string, resErr *apperrors.AuthError) {

	logger.Errorw(logTitle, messages.Error, logErrMsg)
	c.JSON(resErr.StatusCode(), CreateJSONResponse(messages.Error, resErr.Error(), nil))
}
