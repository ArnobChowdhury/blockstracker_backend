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
