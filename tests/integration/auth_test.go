package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"blockstracker_backend/messages"

	"blockstracker_backend/tests/integration/testutils"

	"github.com/stretchr/testify/assert"
)

func TestSignupUserIntegration(t *testing.T) {

	testCases := []struct {
		name           string
		requestBody    map[string]string
		expectedStatus int
		expectedErrMsg string
	}{
		{
			name:           "Success - Valid Request",
			requestBody:    map[string]string{"email": "test@example.com", "password": "StrongPassword123!"},
			expectedStatus: http.StatusOK,
			expectedErrMsg: messages.MsgUserCreationSuccess,
		},
		{
			name:           "Failure - Duplicate Email",
			requestBody:    map[string]string{"email": "test@example.com", "password": "StrongPassword123!"},
			expectedStatus: http.StatusBadRequest,
			expectedErrMsg: messages.ErrUserCreationFailed,
		},
		{
			name:           "Failure - Weak Password",
			requestBody:    map[string]string{"email": "test2@example.com", "password": "weakpassword"},
			expectedStatus: http.StatusBadRequest,
			expectedErrMsg: messages.ErrNotStrongPassword,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := testutils.CreateRequest(http.MethodPost, "/signup", tc.requestBody)
			if err != nil {
				t.Fatalf("Error creating request: %v", err)
			}
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, tc.expectedStatus, resp.Code, "Unexpected status code")
			assert.Contains(t, resp.Body.String(), tc.expectedErrMsg, "Expected error message not found")

		},
		)
	}
}

func TestSigninUserIntegration(t *testing.T) {

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
			req, err := testutils.CreateRequest(http.MethodPost, "/signin", tc.requestBody)
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
