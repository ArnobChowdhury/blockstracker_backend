package utils

import (
	"blockstracker_backend/models"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	Issuer             = "api.blocks-tracker.com"
	AccessTokenExpiry  = 30 * time.Minute
	RefreshTokenExpiry = 7 * 24 * time.Hour
)

func CreateJSONResponse(status string, message string, data interface{}) gin.H {
	response := gin.H{
		"result": gin.H{
			"status":  status,
			"message": message,
		},
	}

	if data != nil {
		if ginData, ok := data.(gin.H); ok {
			response["result"].(gin.H)["data"] = ginData
		} else {
			response["data"] = data
		}
	} else {
		fmt.Printf("Warning: Data is not a gin.H. Type: %T\n", data)
	}

	return response
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
