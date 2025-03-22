package integration

import (
	"blockstracker_backend/messages"
	"blockstracker_backend/tests/integration/testutils"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTagIntegration(t *testing.T) {
	// First, sign in a user to get a valid access token
	signInReqBody := map[string]string{"email": "test@example.com", "password": "StrongPassword123!"}
	signInReq, err := testutils.CreateRequest(http.MethodPost, "/signin", signInReqBody)
	if err != nil {
		t.Fatalf("Error creating sign-in request: %v", err)
	}
	signInResp := httptest.NewRecorder()
	router.ServeHTTP(signInResp, signInReq)
	assert.Equal(t, http.StatusOK, signInResp.Code, "Sign-in failed")

	var signInResponseBody map[string]interface{}
	err = json.Unmarshal(signInResp.Body.Bytes(), &signInResponseBody)
	assert.NoError(t, err)
	result, ok := signInResponseBody["result"].(map[string]interface{})
	assert.True(t, ok)
	data, ok := result["data"].(map[string]interface{})
	assert.True(t, ok)
	accessToken, ok := data["accessToken"].(string)
	assert.True(t, ok)

	testCases := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
	}{
		{
			name: "Success - Valid Tag Creation",
			requestBody: map[string]interface{}{
				"name":       "Test Tag",
				"createdAt":  "2025-03-21T06:33:44.183Z",
				"modifiedAt": "2025-03-21T06:33:44.184Z",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Failure - Missing Name",
			requestBody: map[string]interface{}{
				"createdAt":  "2025-03-21T06:33:44.183Z",
				"modifiedAt": "2025-03-21T06:33:44.184Z",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Failure - Missing createdAt",
			requestBody: map[string]interface{}{
				"name":       "Test Tag",
				"modifiedAt": "2025-03-21T06:33:44.184Z",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Failure - Missing modifiedAt",
			requestBody: map[string]interface{}{
				"name":      "Test Tag",
				"createdAt": "2025-03-21T06:33:44.183Z",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Failure - Invalid createdAt",
			requestBody: map[string]interface{}{
				"name":       "Test Tag",
				"createdAt":  "invalid-date",
				"modifiedAt": "2025-03-21T06:33:44.184Z",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Failure - Invalid modifiedAt",
			requestBody: map[string]interface{}{
				"name":       "Test Tag",
				"createdAt":  "invalid-date",
				"modifiedAt": "invalid-date",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := testutils.CreateRequest(
				http.MethodPost,
				"/tags/",
				tc.requestBody,
				testutils.WithAccessToken(accessToken),
			)
			if err != nil {
				t.Fatalf("Error creating request: %v", err)
			}

			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, tc.expectedStatus, resp.Code, "Unexpected status code")

			if tc.expectedStatus == http.StatusOK {
				var responseBody map[string]interface{}
				err = json.Unmarshal(resp.Body.Bytes(), &responseBody)
				assert.NoError(t, err)
				result, ok := responseBody["result"].(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, messages.Success, result["status"])
				assert.Equal(t, messages.MsgTagCreationSuccess, result["message"])
				data, ok := result["data"].(map[string]interface{})
				assert.True(t, ok)
				assert.NotEmpty(t, data["id"])
				assert.Equal(t, tc.requestBody["name"], data["name"])
				assert.Equal(t, tc.requestBody["createdAt"], data["createdAt"])
				assert.Equal(t, tc.requestBody["modifiedAt"], data["modifiedAt"])
				assert.NotEmpty(t, data["userId"])
			}
		})
	}
}
