package utils_test

import (
	"blockstracker_backend/config"
	apperrors "blockstracker_backend/internal/errors"
	"blockstracker_backend/internal/utils"
	"blockstracker_backend/messages"
	"blockstracker_backend/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateJSONResponse(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		message  string
		data     any
		expected gin.H
	}{
		{
			name:    "Success with data",
			status:  messages.Success,
			message: messages.MsgSignInSuccessful,
			data:    gin.H{"key": "value"},
			expected: gin.H{
				"result": gin.H{
					"status":  messages.Success,
					"message": messages.MsgSignInSuccessful,
					"data":    gin.H{"key": "value"},
				},
			},
		},
		{
			name:    "Success without data",
			status:  messages.Success,
			message: messages.MsgSignInSuccessful,
			data:    nil,
			expected: gin.H{
				"result": gin.H{
					"status":  messages.Success,
					"message": messages.MsgSignInSuccessful,
				},
			},
		},
		{
			name:    "Error",
			status:  messages.Error,
			message: messages.ErrUnauthorized,
			data:    nil,
			expected: gin.H{
				"result": gin.H{
					"status":  messages.Error,
					"message": messages.ErrUnauthorized,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.CreateJSONResponse(tt.status, tt.message, tt.data)
			assert.Equal(t, tt.expected, result)
		})
	}
}
func TestGenerateJWT(t *testing.T) {
	tests := []struct {
		name          string
		email         string
		secretKey     string
		expectedError bool
	}{
		{
			name:          "Valid JWT",
			email:         "test@example.com",
			secretKey:     "testsecret",
			expectedError: false,
		},
		{
			name:          "Empty Secret Key",
			email:         "test@example.com",
			secretKey:     "",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims := utils.GetClaims(&models.User{ID: uuid.New(), Email: tt.email}, "access")
			token, err := utils.GenerateJWT(claims, tt.secretKey)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			}
		})
	}
}

func TestGetClaims(t *testing.T) {
	user := &models.User{
		ID:    uuid.New(),
		Email: "test@example.com",
	}

	tests := []struct {
		name      string
		user      *models.User
		tokenType string
		wantErr   bool
	}{
		{
			name:      "Access Token Claims",
			user:      user,
			tokenType: "access",
			wantErr:   false,
		},
		{
			name:      "Refresh Token Claims",
			user:      user,
			tokenType: "refresh",
			wantErr:   false,
		},
		{
			name:      "Invalid Token Type",
			user:      user,
			tokenType: "invalid",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims := utils.GetClaims(tt.user, tt.tokenType)
			assert.Equal(t, tt.user.ID, claims.UserID)
			assert.Equal(t, tt.user.Email, claims.Email)
			assert.Equal(t, utils.Issuer, claims.Issuer)
			assert.NotNil(t, claims.ExpiresAt)
			assert.NotNil(t, claims.IssuedAt)
		})
	}
}

func TestExtractBearerToken(t *testing.T) {
	tests := []struct {
		name          string
		header        string
		expectedToken string
		expectedError bool
	}{
		{
			name:          "Valid Bearer Token",
			header:        "Bearer validtoken123",
			expectedToken: "validtoken123",
			expectedError: false,
		},
		{
			name:          "No Bearer Prefix",
			header:        "invalidtoken123",
			expectedToken: "",
			expectedError: true,
		},
		{
			name:          "Empty Header",
			header:        "",
			expectedToken: "",
			expectedError: true,
		},
		{
			name:          "Multiple Spaces",
			header:        "Bearer  validtoken123",
			expectedToken: "",
			expectedError: true,
		},
		{
			name:          "Lowercase bearer",
			header:        "bearer validtoken123",
			expectedToken: "validtoken123",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, tokenErr := utils.ExtractBearerToken(tt.header)
			if tt.expectedError {
				assert.Error(t, tokenErr)
			} else {
				assert.Equal(t, tokenErr, (*apperrors.AuthError)(nil))
				assert.Equal(t, tt.expectedToken, token)
			}
		})
	}
}
func TestParseToken(t *testing.T) {
	user := &models.User{
		ID:    uuid.New(),
		Email: "test@example.com",
	}
	accessClaims := utils.GetClaims(user, "access")
	validToken, _ := utils.GenerateJWT(accessClaims, "testsecret")

	tests := []struct {
		name          string
		tokenString   string
		secret        string
		expectedError bool
	}{
		{
			name:          "Valid Token",
			tokenString:   validToken,
			secret:        "testsecret",
			expectedError: false,
		},
		{
			name:          "Invalid Token",
			tokenString:   "invalidtoken",
			secret:        "testsecret",
			expectedError: true,
		},
		{
			name:          "Wrong Secret",
			tokenString:   validToken,
			secret:        "wrongsecret",
			expectedError: true,
		},
		{
			name:          "Empty Token",
			tokenString:   "",
			secret:        "testsecret",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := utils.ParseToken(tt.tokenString, &config.AuthConfig{AccessSecret: tt.secret})
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, user.ID, claims.UserID)
				assert.Equal(t, user.Email, claims.Email)
			}
		})
	}
}
