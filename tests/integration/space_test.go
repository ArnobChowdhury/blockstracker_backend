package integration

import (
	"blockstracker_backend/messages"
	"blockstracker_backend/models"
	"blockstracker_backend/tests/integration/testutils"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateSpaceIntegration(t *testing.T) {
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
			name: "Success - Valid Space Creation",
			requestBody: map[string]interface{}{
				"name":       "Test Space",
				"createdAt":  time.Now().UTC().Format(time.RFC3339Nano),
				"modifiedAt": time.Now().UTC().Format(time.RFC3339Nano),
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Failure - Missing Name",
			requestBody: map[string]interface{}{
				"createdAt":  time.Now().UTC().Format(time.RFC3339Nano),
				"modifiedAt": time.Now().UTC().Format(time.RFC3339Nano),
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Failure - Missing CreatedAt",
			requestBody: map[string]interface{}{
				"name":       "Test Space",
				"modifiedAt": time.Now().UTC().Format(time.RFC3339Nano),
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Failure - Missing ModifiedAt",
			requestBody: map[string]interface{}{
				"name":      "Test Space",
				"createdAt": time.Now().UTC().Format(time.RFC3339Nano),
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Failure - Invalid CreatedAt Format",
			requestBody: map[string]interface{}{
				"name":       "Test Space",
				"createdAt":  "invalid-date",
				"modifiedAt": time.Now().UTC().Format(time.RFC3339Nano),
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Failure - Invalid ModifiedAt Format",
			requestBody: map[string]interface{}{
				"name":       "Test Space",
				"createdAt":  time.Now().UTC().Format(time.RFC3339Nano),
				"modifiedAt": "invalid-date",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Failure - Empty Name",
			requestBody: map[string]interface{}{
				"name":       "",
				"createdAt":  time.Now().UTC().Format(time.RFC3339Nano),
				"modifiedAt": time.Now().UTC().Format(time.RFC3339Nano),
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := testutils.CreateRequest(
				http.MethodPost,
				"/spaces/",
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
				assert.Equal(t, messages.MsgSpaceCreationSuccess, result["message"])
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

func TestUpdateSpaceIntegration(t *testing.T) {
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

	// Create a second user to simulate updating a space that doesn't belong to the first user
	signUpReqBody2 := map[string]string{"email": "test4@example.com", "password": "StrongPassword123!"}
	signUpReq2, err := testutils.CreateRequest(http.MethodPost, "/signup", signUpReqBody2)
	if err != nil {
		t.Fatalf("Error creating sign-up request for second user: %v", err)
	}
	signUpResp2 := httptest.NewRecorder()
	router.ServeHTTP(signUpResp2, signUpReq2)
	assert.Equal(t, http.StatusOK, signUpResp2.Code, "Sign-up failed for second user")
	signInReqBody2 := map[string]string{"email": "test4@example.com", "password": "StrongPassword123!"}
	signInReq2, err := testutils.CreateRequest(http.MethodPost, "/signin", signInReqBody2)
	if err != nil {
		t.Fatalf("Error creating sign-in request for second user: %v", err)
	}
	signInResp2 := httptest.NewRecorder()
	router.ServeHTTP(signInResp2, signInReq2)
	assert.Equal(t, http.StatusOK, signInResp2.Code, "Sign-in failed for second user")

	// Create a space to update
	createSpaceReqBody := map[string]interface{}{
		"name":       "Initial Space",
		"createdAt":  time.Now().UTC().Format(time.RFC3339Nano),
		"modifiedAt": time.Now().UTC().Format(time.RFC3339Nano),
	}
	createSpaceReq, err := testutils.CreateRequest(http.MethodPost, "/spaces/", createSpaceReqBody, testutils.WithAccessToken(accessToken))
	if err != nil {
		t.Fatalf("Error creating create space request: %v", err)
	}
	createSpaceResp := httptest.NewRecorder()
	router.ServeHTTP(createSpaceResp, createSpaceReq)
	assert.Equal(t, http.StatusOK, createSpaceResp.Code, "Create space failed")

	var createSpaceResponseBody map[string]interface{}
	err = json.Unmarshal(createSpaceResp.Body.Bytes(), &createSpaceResponseBody)
	assert.NoError(t, err)
	createSpaceResult, ok := createSpaceResponseBody["result"].(map[string]interface{})
	assert.True(t, ok)
	createSpaceData, ok := createSpaceResult["data"].(map[string]interface{})
	assert.True(t, ok)
	spaceID, ok := createSpaceData["id"].(string)
	assert.True(t, ok)

	var signInResponseBody2 map[string]interface{}
	err = json.Unmarshal(signInResp2.Body.Bytes(), &signInResponseBody2)
	assert.NoError(t, err)
	result2, ok := signInResponseBody2["result"].(map[string]interface{})
	assert.True(t, ok)
	data2, ok := result2["data"].(map[string]interface{})
	assert.True(t, ok)

	testCases := []struct {
		name                string
		requestBody         map[string]interface{}
		expectedStatus      int
		expectedResponseMsg string
		accessToken         string
		checkMessage        bool
	}{
		{
			name: "Success - Update Name",
			requestBody: map[string]interface{}{
				"name":       "Updated Space Name",
				"createdAt":  createSpaceReqBody["createdAt"],
				"modifiedAt": time.Now().UTC().Format(time.RFC3339Nano),
			},
			expectedStatus:      http.StatusOK,
			expectedResponseMsg: messages.MsgSpaceUpdateSuccess,
			accessToken:         accessToken,
			checkMessage:        true,
		},
		{
			name: "Failure - Invalid Space ID",
			requestBody: map[string]interface{}{
				"name":       "Updated Space Name",
				"createdAt":  createSpaceReqBody["createdAt"],
				"modifiedAt": time.Now().UTC().Format(time.RFC3339Nano),
			},
			expectedStatus:      http.StatusBadRequest,
			expectedResponseMsg: messages.ErrSpaceUpdateFailed,
			accessToken:         accessToken,
			checkMessage:        false,
		},
		{
			name: "Failure - Missing ModifiedAt",
			requestBody: map[string]interface{}{
				"name":      "Updated Space Name",
				"createdAt": createSpaceReqBody["createdAt"],
			},
			expectedStatus:      http.StatusBadRequest,
			expectedResponseMsg: messages.ErrSpaceUpdateFailed,
			accessToken:         accessToken,
			checkMessage:        false,
		},
		{
			name: "Failure - Invalid ModifiedAt Format",
			requestBody: map[string]interface{}{
				"name":       "Updated Space Name",
				"createdAt":  createSpaceReqBody["createdAt"],
				"modifiedAt": "invalid-date",
			},
			expectedStatus:      http.StatusBadRequest,
			expectedResponseMsg: messages.ErrSpaceUpdateFailed,
			accessToken:         accessToken,
			checkMessage:        false,
		},
		{
			name: "Failure - Space Not Found",
			requestBody: map[string]interface{}{
				"name":       "Updated Space Name",
				"createdAt":  createSpaceReqBody["createdAt"],
				"modifiedAt": time.Now().UTC().Format(time.RFC3339Nano),
			},
			expectedStatus:      http.StatusUnauthorized,
			expectedResponseMsg: messages.ErrUnauthorized,
			accessToken:         accessToken,
			checkMessage:        true,
		},
		{
			name: "Failure - Update Space That Doesn't Belong To User",
			requestBody: map[string]interface{}{
				"name":       "Updated Space Name",
				"createdAt":  createSpaceReqBody["createdAt"],
				"modifiedAt": time.Now().UTC().Format(time.RFC3339Nano),
			},
			expectedStatus:      http.StatusUnauthorized,
			expectedResponseMsg: messages.ErrUnauthorized,
			accessToken:         data2["accessToken"].(string),
			checkMessage:        true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var req *http.Request
			var err error
			if tc.name == "Failure - Space Not Found" {
				req, err = testutils.CreateRequest(
					http.MethodPut,
					fmt.Sprintf("/spaces/%s", uuid.New().String()),
					tc.requestBody,
					testutils.WithAccessToken(tc.accessToken),
				)
			} else if tc.name == "Failure - Invalid Space ID" {
				req, err = testutils.CreateRequest(
					http.MethodPut,
					fmt.Sprintf("/spaces/%s", "invalid-space-id"),
					tc.requestBody,
					testutils.WithAccessToken(tc.accessToken),
				)
			} else {
				req, err = testutils.CreateRequest(
					http.MethodPut,
					fmt.Sprintf("/spaces/%s", spaceID),
					tc.requestBody,
					testutils.WithAccessToken(tc.accessToken),
				)
			}

			if err != nil {
				t.Fatalf("Error creating request: %v", err)
			}

			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, tc.expectedStatus, resp.Code, "Unexpected status code")

			if tc.expectedStatus == http.StatusOK {
				if tc.checkMessage {
					assert.Contains(t, resp.Body.String(), tc.expectedResponseMsg, "Expected error message not found")
				}
				var responseBody map[string]interface{}
				err = json.Unmarshal(resp.Body.Bytes(), &responseBody)
				assert.NoError(t, err)
				result, ok := responseBody["result"].(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, messages.Success, result["status"])
				assert.Equal(t, messages.MsgSpaceUpdateSuccess, result["message"])
				data, ok := result["data"].(map[string]interface{})
				assert.True(t, ok)
				assert.NotEmpty(t, data["id"])

				updatedSpace := &models.Space{ID: uuid.MustParse(spaceID)}
				err = TestDB.First(updatedSpace).Error
				assert.NoError(t, err)

				if name, ok := tc.requestBody["name"].(string); ok {
					assert.Equal(t, name, updatedSpace.Name)
				}
			} else {
				if tc.checkMessage {
					assert.Contains(t, resp.Body.String(), tc.expectedResponseMsg, "Expected error message not found")
				}
			}
		})
	}
}
