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

// func TestSignupUserIntegration(t *testing.T) {
// 	router, err := setupRouter(TestDB)
// 	if err != nil {
// 		t.Fatalf("Error setting up router: %v", err)
// 	}

// 	t.Run("Success - Valid Request", func(t *testing.T) {
// 		requestBody := map[string]string{
// 			"email":    "test@example.com",
// 			"password": "StrongPassword123!",
// 		}
// 		jsonBody, _ := json.Marshal(requestBody)

// 		req, _ := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(jsonBody))
// 		req.Header.Set("Content-Type", "application/json")
// 		resp := httptest.NewRecorder()

// 		router.ServeHTTP(resp, req)

// 		assert.Equal(t, http.StatusOK, resp.Code)
// 		assert.Contains(t, resp.Body.String(), messages.MsgUserCreationSuccess)

// 		var user models.User
// 		err := TestDB.Where("email = ?", "test@example.com").First(&user).Error
// 		assert.Nil(t, err)
// 		assert.Equal(t, "test@example.com", user.Email)
// 	})

// 	t.Run("Failure - Duplicate Email", func(t *testing.T) {
// 		requestBody := map[string]string{
// 			"email":    "test@example.com",
// 			"password": "StrongPassword123!",
// 		}
// 		jsonBody, _ := json.Marshal(requestBody)

// 		req, _ := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(jsonBody))
// 		req.Header.Set("Content-Type", "application/json")
// 		resp := httptest.NewRecorder()

// 		router.ServeHTTP(resp, req)

// 		assert.Equal(t, http.StatusBadRequest, resp.Code)
// 		assert.Contains(t, resp.Body.String(), messages.ErrUserCreationFailed)
// 	})

// 	t.Run("Failure - Weak Password", func(t *testing.T) {
// 		requestBody := map[string]string{
// 			"email":    "test@example.com",
// 			"password": "weakpassword",
// 		}
// 		jsonBody, _ := json.Marshal(requestBody)

// 		req, _ := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(jsonBody))
// 		req.Header.Set("Content-Type", "application/json")
// 		resp := httptest.NewRecorder()

// 		router.ServeHTTP(resp, req)

// 		assert.Equal(t, http.StatusBadRequest, resp.Code)
// 		assert.Contains(t, resp.Body.String(), messages.ErrNotStrongPassword)
// 	})
// }

// func TestSigninUserIntegration(t *testing.T) {
// 	router, err := setupRouter(TestDB)
// 	if err != nil {
// 		t.Fatalf("Error setting up router: %v", err)
// 	}

// 	t.Run("Success - Valid Credentials", func(t *testing.T) {
// 		requestBody := map[string]string{
// 			"email":    "test@example.com",
// 			"password": "StrongPassword123!",
// 		}

// 		jsonBody, _ := json.Marshal(requestBody)

// 		req, _ := http.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(jsonBody))
// 		req.Header.Set("Content-Type", "application/json")
// 		resp := httptest.NewRecorder()

// 		router.ServeHTTP(resp, req)

// 		assert.Equal(t, http.StatusOK, resp.Code)
// 		var responseBody map[string]interface{}
// 		err = json.Unmarshal(resp.Body.Bytes(), &responseBody)
// 		assert.NoError(t, err)
// 		result, ok := responseBody["result"].(map[string]interface{})
// 		assert.True(t, ok)
// 		data, ok := result["data"].(map[string]interface{})
// 		assert.True(t, ok)
// 		_, accessTokenExists := data["accessToken"].(string)
// 		_, refreshTokenExists := data["refreshToken"].(string)
// 		assert.True(t, accessTokenExists && refreshTokenExists)
// 		assert.Contains(t, resp.Body.String(), messages.MsgSignInSuccessful)
// 	})

// 	t.Run("Failure - User Not Found", func(t *testing.T) {
// 		requestBody := map[string]string{
// 			"email":    "nonexistent@example.com",
// 			"password": "anypassword",
// 		}
// 		jsonBody, _ := json.Marshal(requestBody)

// 		req, _ := http.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(jsonBody))
// 		req.Header.Set("Content-Type", "application/json")
// 		resp := httptest.NewRecorder()

// 		router.ServeHTTP(resp, req)

// 		assert.Equal(t, http.StatusUnauthorized, resp.Code)
// 		assert.Contains(t, resp.Body.String(), messages.ErrInvalidCredentials)
// 	})

// 	t.Run("Failure - Invalid Credentials", func(t *testing.T) {
// 		requestBody := map[string]string{
// 			"email":    "test@example.com",
// 			"password": "wrongpassword",
// 		}
// 		jsonBody, _ := json.Marshal(requestBody)

// 		req, _ := http.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(jsonBody))
// 		req.Header.Set("Content-Type", "application/json")
// 		resp := httptest.NewRecorder()

// 		router.ServeHTTP(resp, req)

// 		assert.Equal(t, http.StatusUnauthorized, resp.Code)
// 		assert.Contains(t, resp.Body.String(), messages.ErrInvalidCredentials)
// 	})

// 	t.Run("Failure - Empty Email", func(t *testing.T) {
// 		requestBody := map[string]string{
// 			"email":    "",
// 			"password": "StrongPassword123!",
// 		}
// 		jsonBody, _ := json.Marshal(requestBody)

// 		req, _ := http.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(jsonBody))
// 		req.Header.Set("Content-Type", "application/json")
// 		resp := httptest.NewRecorder()

// 		router.ServeHTTP(resp, req)

// 		assert.Equal(t, http.StatusBadRequest, resp.Code)
// 		assert.Contains(t, resp.Body.String(), messages.ErrMalformedRequest)
// 	})

// 	t.Run("Failure - Invalid Email Format", func(t *testing.T) {
// 		requestBody := map[string]string{
// 			"email":    "test",
// 			"password": "StrongPassword123!",
// 		}
// 		jsonBody, _ := json.Marshal(requestBody)

// 		req, _ := http.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(jsonBody))
// 		req.Header.Set("Content-Type", "application/json")
// 		resp := httptest.NewRecorder()

// 		router.ServeHTTP(resp, req)

// 		assert.Equal(t, http.StatusBadRequest, resp.Code)
// 		assert.Contains(t, resp.Body.String(), messages.ErrMalformedRequest)
// 	})
// }

func createRequest(method, path string, body interface{}) (*http.Request, error) {
	var jsonBody []byte
	var err error

	if body != nil {
		jsonBody, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("error marshaling request body: %w", err)
		}
	}

	req, err := http.NewRequest(method, path, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func TestSigninUserIntegration(t *testing.T) {
	router, err := setupRouter(TestDB)
	if err != nil {
		t.Fatalf("Error setting up router: %v", err)
	}

	testCases := []struct {
		name           string
		requestBody    map[string]string
		expectedStatus int
		expectedErrMsg string
	}{
		{
			name:           "Success - Valid Credentials",
			requestBody:    map[string]string{"email": "test@example.com", "password": "StrongPassword123!"},
			expectedStatus: http.StatusOK,
			expectedErrMsg: "",
		},
		{
			name:           "Failure - User Not Found",
			requestBody:    map[string]string{"email": "nonexistent@example.com", "password": "anypassword"},
			expectedStatus: http.StatusUnauthorized,
			expectedErrMsg: messages.ErrInvalidCredentials,
		},
		{
			name:           "Failure - Invalid Credentials",
			requestBody:    map[string]string{"email": "test@example.com", "password": "wrongpassword"},
			expectedStatus: http.StatusUnauthorized,
			expectedErrMsg: messages.ErrInvalidCredentials,
		},
		{
			name:           "Failure - Empty Email",
			requestBody:    map[string]string{"email": "", "password": "StrongPassword123!"},
			expectedStatus: http.StatusBadRequest,
			expectedErrMsg: messages.ErrMalformedRequest,
		},
		{
			name:           "Failure - Invalid Email Format",
			requestBody:    map[string]string{"email": "test", "password": "StrongPassword123!"},
			expectedStatus: http.StatusBadRequest,
			expectedErrMsg: messages.ErrMalformedRequest,
		},
		{
			name:           "Failure - Empty Password",
			requestBody:    map[string]string{"email": "test@example.com", "password": ""},
			expectedStatus: http.StatusBadRequest,
			expectedErrMsg: messages.ErrMalformedRequest, // Adjust if a specific error message is returned
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := createRequest(http.MethodPost, "/signin", tc.requestBody)
			if err != nil {
				t.Fatalf("Error creating request: %v", err)
			}
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, tc.expectedStatus, resp.Code, "Unexpected status code")
			if tc.expectedErrMsg != "" {
				assert.Contains(t, resp.Body.String(), tc.expectedErrMsg, "Expected error message not found")
			} else {
				//Check for access and refresh tokens in success case.
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

			}
		})
	}
}
