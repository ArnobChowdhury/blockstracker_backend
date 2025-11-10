package utils

import (
	apperrors "blockstracker_backend/internal/errors"
	"blockstracker_backend/messages"
	"blockstracker_backend/models"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	Issuer = "api.blocks-tracker.com"
	// this probably should be in the config
	AccessTokenExpiry  = 30 * time.Minute
	RefreshTokenExpiry = 7 * 24 * time.Hour
)

func CreateJSONResponse(status string, message string, data any, code ...string) gin.H {
	result := gin.H{
		"status":  status,
		"message": message,
	}

	if data != nil {
		result["data"] = data
	}

	if len(code) > 0 && code[0] != "" {
		result["code"] = code[0]
	}

	return gin.H{"result": result}
}

func GenerateJWT(claims *models.Claims, secretKey string) (string, error) {
	if secretKey == "" {
		return "", fmt.Errorf("secret key is empty")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err // Consider wrapping the error for more context
	}

	return signedToken, nil
}

func GetClaims(user *models.User, tokenType string) *models.Claims {
	expiresAt := time.Now().Add(AccessTokenExpiry)

	if tokenType == "refresh" {
		expiresAt = time.Now().Add(RefreshTokenExpiry)
	}

	claims := &models.Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    Issuer,
			ID:        uuid.NewString(),
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

func ParseToken(tokenString string, secretKey string) (*models.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apperrors.ErrUnexpectedSigningMethod
		}
		return []byte(secretKey), nil // Use the provided secretKey
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
	logErrMsg string, resErr apperrors.AppError, data ...any) {

	logger.Errorw(logTitle, messages.Error, logErrMsg)
	var responseData any
	if len(data) > 0 {
		responseData = data[0]
	}

	c.JSON(resErr.StatusCode(), CreateJSONResponse(messages.Error, resErr.Error(), responseData, resErr.Code()))
}

func ExtractUIDFromGinContext(c *gin.Context) (uuid.UUID, *apperrors.CommonError) {
	userID, ok := c.Get("userID")
	if !ok {
		return uuid.Nil, apperrors.ErrUserIDNotFoundInContext
	}

	uid, ok := userID.(uuid.UUID)
	if !ok {
		return uuid.Nil, apperrors.ErrUserIDNotValidType
	}

	return uid, nil
}
