package integration

import (
	"blockstracker_backend/messages"
	"blockstracker_backend/tests/integration/testutils"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateTaskIntegration(t *testing.T) {
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
		expectedErrMsg string
	}{
		{
			name: "Success - Valid Task Creation",
			requestBody: map[string]interface{}{
				"isActive":                 true,
				"title":                    "Test Task",
				"description":              "This is a test task",
				"schedule":                 "Once",
				"priority":                 3,
				"completionStatus":         "INCOMPLETE",
				"dueDate":                  time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339Nano),
				"shouldBeScored":           true,
				"score":                    5,
				"timeOfDay":                "morning",
				"repetitiveTaskTemplateID": nil,
				"createdAt":                time.Now().UTC().Format(time.RFC3339Nano),
				"modifiedAt":               time.Now().UTC().Format(time.RFC3339Nano),
				"tags":                     nil,
				"spaceID":                  nil,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Failure - Missing Title",
			requestBody: map[string]interface{}{
				"isActive":                 true,
				"description":              "This is a test task",
				"schedule":                 "Once",
				"priority":                 3,
				"completionStatus":         "INCOMPLETE",
				"dueDate":                  time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339Nano),
				"shouldBeScored":           true,
				"score":                    5,
				"timeOfDay":                "morning",
				"repetitiveTaskTemplateID": nil,
				"createdAt":                time.Now().UTC().Format(time.RFC3339Nano),
				"modifiedAt":               time.Now().UTC().Format(time.RFC3339Nano),
				"tags":                     nil,
				"spaceID":                  nil,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Failure - Invalid Due Date Format",
			requestBody: map[string]interface{}{
				"isActive":                 true,
				"description":              "This is a test task",
				"schedule":                 "Once",
				"priority":                 3,
				"completionStatus":         "INCOMPLETE",
				"dueDate":                  "invalid-date",
				"shouldBeScored":           true,
				"score":                    5,
				"timeOfDay":                "morning",
				"repetitiveTaskTemplateID": nil,
				"createdAt":                time.Now().UTC().Format(time.RFC3339Nano),
				"modifiedAt":               time.Now().UTC().Format(time.RFC3339Nano),
				"tags":                     nil,
				"spaceID":                  nil,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Failure - Invalid Space ID",
			requestBody: map[string]interface{}{
				"isActive":                 true,
				"description":              "This is a test task",
				"schedule":                 "Once",
				"priority":                 3,
				"completionStatus":         "INCOMPLETE",
				"dueDate":                  "invalid-date",
				"shouldBeScored":           true,
				"score":                    5,
				"timeOfDay":                "morning",
				"repetitiveTaskTemplateID": nil,
				"createdAt":                time.Now().UTC().Format(time.RFC3339Nano),
				"modifiedAt":               time.Now().UTC().Format(time.RFC3339Nano),
				"tags":                     nil,
				"spaceID":                  "invalid-space-id",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := testutils.CreateRequest(
				http.MethodPost,
				"/tasks/",
				tc.requestBody,
				testutils.WithAccessToken(accessToken),
			)
			if err != nil {
				t.Fatalf("Error creating request: %v", err)
			}

			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, tc.expectedStatus, resp.Code, "Unexpected status code")
			assert.Contains(t, resp.Body.String(), tc.expectedErrMsg, "Expected error message not found")

			if tc.expectedStatus == http.StatusOK {
				var responseBody map[string]interface{}
				err = json.Unmarshal(resp.Body.Bytes(), &responseBody)
				assert.NoError(t, err)
				result, ok := responseBody["result"].(map[string]interface{})
				assert.Equal(t, messages.Success, result["status"])
				assert.Equal(t, messages.MsgTaskCreationSuccess, result["message"])
				assert.True(t, ok)
				data, ok := result["data"].(map[string]interface{})

				assert.True(t, ok)
				assert.NotEmpty(t, data["id"])
				assert.Equal(t, tc.requestBody["title"], data["title"])
				assert.Equal(t, tc.requestBody["description"], data["description"])
				assert.Equal(t, tc.requestBody["schedule"], data["schedule"])
				if priorityFloat, ok := data["priority"].(float64); ok {
					assert.Equal(t, tc.requestBody["priority"], int(priorityFloat))
				} else {
					t.Errorf("priority is not a float64: %v", data["priority"])
				}
				assert.Equal(t, tc.requestBody["completionStatus"], data["completionStatus"])
				assert.Equal(t, tc.requestBody["shouldBeScored"], data["shouldBeScored"])
				if scoreFloat, ok := data["score"].(float64); ok {
					assert.Equal(t, tc.requestBody["score"], int(scoreFloat))
				} else {
					t.Errorf("score is not a float64: %v", data["score"])
				}
				assert.Equal(t, tc.requestBody["timeOfDay"], data["timeOfDay"])
				assert.Equal(t, tc.requestBody["tags"], data["tags"])
				assert.Equal(t, tc.requestBody["isActive"], data["isActive"])
				assert.Equal(t, tc.requestBody["createdAt"], data["createdAt"])
				assert.Equal(t, tc.requestBody["modifiedAt"], data["modifiedAt"])
				assert.Equal(t, tc.requestBody["spaceId"], data["spaceId"])
				assert.Equal(t, tc.requestBody["repetitiveTaskTemplateId"], data["repetitiveTaskTemplateId"])
				assert.Equal(t, tc.requestBody["dueDate"], data["dueDate"])
				assert.NotEmpty(t, data["userId"])
			}
		})
	}
}
