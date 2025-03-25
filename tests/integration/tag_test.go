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
			name: "Failure - Missing createdAt",
			requestBody: map[string]interface{}{
				"name":       "Test Tag",
				"modifiedAt": time.Now().UTC().Format(time.RFC3339Nano),
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Failure - Missing modifiedAt",
			requestBody: map[string]interface{}{
				"name":      "Test Tag",
				"createdAt": time.Now().UTC().Format(time.RFC3339Nano),
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Failure - Invalid createdAt",
			requestBody: map[string]interface{}{
				"name":       "Test Tag",
				"createdAt":  "invalid-date",
				"modifiedAt": time.Now().UTC().Format(time.RFC3339Nano),
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

func TestUpdateTagIntegration(t *testing.T) {
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

	// Create a second user to simulate updating a tag that doesn't belong to the first user
	signUpReqBody2 := map[string]string{"email": "test2@example.com", "password": "StrongPassword123!"}
	signUpReq2, err := testutils.CreateRequest(http.MethodPost, "/signup", signUpReqBody2)
	if err != nil {
		t.Fatalf("Error creating sign-up request for second user: %v", err)
	}
	signUpResp2 := httptest.NewRecorder()
	router.ServeHTTP(signUpResp2, signUpReq2)
	assert.Equal(t, http.StatusOK, signUpResp2.Code, "Sign-up failed for second user")
	signInReqBody2 := map[string]string{"email": "test2@example.com", "password": "StrongPassword123!"}
	signInReq2, err := testutils.CreateRequest(http.MethodPost, "/signin", signInReqBody2)
	if err != nil {
		t.Fatalf("Error creating sign-in request for second user: %v", err)
	}
	signInResp2 := httptest.NewRecorder()
	router.ServeHTTP(signInResp2, signInReq2)
	assert.Equal(t, http.StatusOK, signInResp2.Code, "Sign-in failed for second user")

	// Create a tag to update
	createTagReqBody := map[string]interface{}{
		"name":       "Initial Tag",
		"createdAt":  time.Now().UTC().Format(time.RFC3339Nano),
		"modifiedAt": time.Now().UTC().Format(time.RFC3339Nano),
	}
	createTagReq, err := testutils.CreateRequest(http.MethodPost, "/tags/", createTagReqBody, testutils.WithAccessToken(accessToken))
	if err != nil {
		t.Fatalf("Error creating create tag request: %v", err)
	}
	createTagResp := httptest.NewRecorder()
	router.ServeHTTP(createTagResp, createTagReq)
	assert.Equal(t, http.StatusOK, createTagResp.Code, "Create tag failed")

	var createTagResponseBody map[string]interface{}
	err = json.Unmarshal(createTagResp.Body.Bytes(), &createTagResponseBody)
	assert.NoError(t, err)
	createTagResult, ok := createTagResponseBody["result"].(map[string]interface{})
	assert.True(t, ok)
	createTagData, ok := createTagResult["data"].(map[string]interface{})
	assert.True(t, ok)
	tagID, ok := createTagData["id"].(string)
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
	}{
		{
			name: "Success - Update Name",
			requestBody: map[string]interface{}{
				"name":       "Updated Tag Name",
				"createdAt":  createTagReqBody["createdAt"],
				"modifiedAt": time.Now().UTC().Format(time.RFC3339Nano),
			},
			expectedStatus:      http.StatusOK,
			expectedResponseMsg: messages.MsgTagUpdateSuccess,
		},
		{
			name: "Failure - Invalid Tag ID",
			requestBody: map[string]interface{}{
				"name":       "Updated Tag Name",
				"createdAt":  createTagReqBody["createdAt"],
				"modifiedAt": time.Now().UTC().Format(time.RFC3339Nano),
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Failure - Missing ModifiedAt",
			requestBody: map[string]interface{}{
				"name":      "Updated Tag Name",
				"createdAt": createTagReqBody["createdAt"],
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Failure - Invalid ModifiedAt Format",
			requestBody: map[string]interface{}{
				"name":       "Updated Tag Name",
				"createdAt":  createTagReqBody["createdAt"],
				"modifiedAt": "invalid-date",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Failure - Tag Not Found",
			requestBody: map[string]interface{}{
				"name":       "Updated Tag Name",
				"createdAt":  createTagReqBody["createdAt"],
				"modifiedAt": time.Now().UTC().Format(time.RFC3339Nano),
			},
			expectedStatus: http.StatusUnauthorized,
		},

		{
			name: "Failure - Update Tag That Doesn't Belong To User",
			requestBody: map[string]interface{}{
				"name":       "Updated Tag Name",
				"createdAt":  createTagReqBody["createdAt"],
				"modifiedAt": time.Now().UTC().Format(time.RFC3339Nano),
			},
			expectedStatus:      http.StatusUnauthorized,
			expectedResponseMsg: messages.ErrTagUpdateFailed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var req *http.Request
			var err error
			if tc.name == "Failure - Tag Not Found" || tc.name == "Failure - Invalid Tag ID" {
				req, err = testutils.CreateRequest(
					http.MethodPut,
					fmt.Sprintf("/tags/%s", uuid.New().String()),
					tc.requestBody,
					testutils.WithAccessToken(accessToken),
				)
			} else if tc.name == "Failure - Update Tag That Doesn't Belong To User" {
				req, err = testutils.CreateRequest(
					http.MethodPut,
					fmt.Sprintf("/tags/%s", tagID),
					tc.requestBody,
					testutils.WithAccessToken(data2["accessToken"].(string)),
				)
			} else {
				req, err = testutils.CreateRequest(
					http.MethodPut,
					fmt.Sprintf("/tags/%s", tagID),
					tc.requestBody,
					testutils.WithAccessToken(accessToken),
				)
			}

			if err != nil {
				t.Fatalf("Error creating request: %v", err)
			}

			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, tc.expectedStatus, resp.Code, "Unexpected status code")

			if tc.expectedStatus == http.StatusOK {
				assert.Contains(t, resp.Body.String(), tc.expectedResponseMsg, "Expected error message not found")
				var responseBody map[string]interface{}
				err = json.Unmarshal(resp.Body.Bytes(), &responseBody)
				assert.NoError(t, err)
				result, ok := responseBody["result"].(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, messages.Success, result["status"])
				assert.Equal(t, messages.MsgTagUpdateSuccess, result["message"])
				data, ok := result["data"].(map[string]interface{})
				assert.True(t, ok)
				assert.NotEmpty(t, data["id"])

				updatedTag := &models.Tag{ID: uuid.MustParse(tagID)}
				err = TestDB.First(updatedTag).Error
				assert.NoError(t, err)

				if name, ok := tc.requestBody["name"].(string); ok {
					assert.Equal(t, name, updatedTag.Name)
				}
			}
		})
	}
}
