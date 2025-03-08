package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"blockstracker_backend/config"
	"blockstracker_backend/handlers"
	"blockstracker_backend/internal/repositories"
	"blockstracker_backend/messages"
	"blockstracker_backend/models"
	"blockstracker_backend/pkg/logger"

	// "blockstracker_backend/tests"
	// "blockstracker_backend/tests/integration"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupRouter(db *gorm.DB) (*gin.Engine, error) {
	gin.SetMode(gin.TestMode)

	userRepo := repositories.NewUserRepository(db)
	authConfig, err := config.LoadAuthConfig()
	if err != nil {
		return nil, fmt.Errorf("Error loading auth config: %v", err)
	}
	authHandler := handlers.NewAuthHandler(userRepo, logger.Log, authConfig)

	router := gin.Default()
	router.POST("/signup", authHandler.SignupUser)
	router.POST("/signin", authHandler.EmailSignIn)
	return router, nil
}

func TestSignupUserIntegration(t *testing.T) {
	router, err := setupRouter(TestDB)
	if err != nil {
		t.Fatalf("Error setting up router: %v", err)
	}

	t.Run("Success - Valid Request", func(t *testing.T) {
		requestBody := map[string]string{
			"email":    "test@example.com",
			"password": "StrongPassword123!",
		}
		jsonBody, _ := json.Marshal(requestBody)

		req, _ := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), messages.MsgUserCreationSuccess)

		var user models.User
		err := TestDB.Where("email = ?", "test@example.com").First(&user).Error
		assert.Nil(t, err)
		assert.Equal(t, "test@example.com", user.Email)
	})

	t.Run("Failure - Duplicate Email", func(t *testing.T) {
		requestBody := map[string]string{
			"email":    "test@example.com",
			"password": "StrongPassword123!",
		}
		jsonBody, _ := json.Marshal(requestBody)

		req, _ := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Contains(t, resp.Body.String(), messages.ErrUserCreationFailed)
	})

	t.Run("Failure - Weak Password", func(t *testing.T) {
		requestBody := map[string]string{
			"email":    "test@example.com",
			"password": "weakpassword",
		}
		jsonBody, _ := json.Marshal(requestBody)

		req, _ := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Contains(t, resp.Body.String(), messages.ErrNotStrongPassword)
	})
}

func TestSigninUserIntegration(t *testing.T) {
	router, err := setupRouter(TestDB)
	if err != nil {
		t.Fatalf("Error setting up router: %v", err)
	}

	t.Run("Success - Valid Credentials", func(t *testing.T) {
		requestBody := map[string]string{
			"email":    "test@example.com",
			"password": "StrongPassword123!",
		}

		jsonBody, _ := json.Marshal(requestBody)

		req, _ := http.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		var responseBody map[string]interface{}
		err = json.Unmarshal(resp.Body.Bytes(), &responseBody)
		assert.NoError(t, err)
		result, ok := responseBody["result"].(map[string]interface{})
		assert.True(t, ok)
		data, ok := result["data"].(map[string]interface{})
		assert.True(t, ok)
		_, accessTokenExists := data["accessToken"].(string)
		_, refreshTokenExists := data["refreshToken"].(string)
		assert.True(t, accessTokenExists && refreshTokenExists)
		assert.Contains(t, resp.Body.String(), messages.MsgSignInSuccessful)
	})

	t.Run("Failure - User Not Found", func(t *testing.T) {
		requestBody := map[string]string{
			"email":    "nonexistent@example.com",
			"password": "anypassword",
		}
		jsonBody, _ := json.Marshal(requestBody)

		req, _ := http.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
		assert.Contains(t, resp.Body.String(), messages.ErrInvalidCredentials)
	})

	t.Run("Failure - Invalid Credentials", func(t *testing.T) {
		requestBody := map[string]string{
			"email":    "test@example.com",
			"password": "wrongpassword",
		}
		jsonBody, _ := json.Marshal(requestBody)

		req, _ := http.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
		assert.Contains(t, resp.Body.String(), messages.ErrInvalidCredentials)
	})

	t.Run("Failure - Empty Email", func(t *testing.T) {
		requestBody := map[string]string{
			"email":    "",
			"password": "StrongPassword123!",
		}
		jsonBody, _ := json.Marshal(requestBody)

		req, _ := http.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Contains(t, resp.Body.String(), messages.ErrMalformedRequest)
	})

	t.Run("Failure - Invalid Email Format", func(t *testing.T) {
		requestBody := map[string]string{
			"email":    "test",
			"password": "StrongPassword123!",
		}
		jsonBody, _ := json.Marshal(requestBody)

		req, _ := http.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Contains(t, resp.Body.String(), messages.ErrMalformedRequest)
	})
}
