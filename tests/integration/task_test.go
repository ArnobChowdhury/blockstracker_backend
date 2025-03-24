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

func TestUpdateTaskIntegration(t *testing.T) {
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

	// Create a task to update
	createTaskReqBody := map[string]interface{}{
		"isActive":                 true,
		"title":                    "Initial Task",
		"description":              "Initial description",
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
	}
	createTaskReq, err := testutils.CreateRequest(http.MethodPost, "/tasks/", createTaskReqBody, testutils.WithAccessToken(accessToken))
	if err != nil {
		t.Fatalf("Error creating create task request: %v", err)
	}
	createTaskResp := httptest.NewRecorder()
	router.ServeHTTP(createTaskResp, createTaskReq)
	assert.Equal(t, http.StatusOK, createTaskResp.Code, "Create task failed")

	var createTaskResponseBody map[string]interface{}
	err = json.Unmarshal(createTaskResp.Body.Bytes(), &createTaskResponseBody)
	assert.NoError(t, err)
	createTaskResult, ok := createTaskResponseBody["result"].(map[string]interface{})
	assert.True(t, ok)
	createTaskData, ok := createTaskResult["data"].(map[string]interface{})
	assert.True(t, ok)
	taskID, ok := createTaskData["id"].(string)
	assert.True(t, ok)

	testCases := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		expectedMsg    string
	}{
		{
			name: "Success - Update Title",
			requestBody: map[string]interface{}{
				"title":                    "Updated Title", // Update the title
				"isActive":                 createTaskReqBody["isActive"],
				"description":              createTaskReqBody["description"],
				"schedule":                 createTaskReqBody["schedule"],
				"priority":                 createTaskReqBody["priority"],
				"completionStatus":         createTaskReqBody["completionStatus"],
				"dueDate":                  createTaskReqBody["dueDate"],
				"shouldBeScored":           createTaskReqBody["shouldBeScored"],
				"score":                    createTaskReqBody["score"],
				"timeOfDay":                createTaskReqBody["timeOfDay"],
				"repetitiveTaskTemplateId": createTaskReqBody["repetitiveTaskTemplateID"],
				"createdAt":                createTaskReqBody["createdAt"],
				"modifiedAt":               time.Now().UTC().Format(time.RFC3339Nano), // Update modifiedAt
				"tags":                     createTaskReqBody["tags"],
				"spaceId":                  createTaskReqBody["spaceID"],
			},
			expectedStatus: http.StatusOK,
			expectedMsg:    messages.MsgTaskUpdateSuccess,
		},
		{
			name: "Success - Update Description",
			requestBody: map[string]interface{}{
				"title":                    createTaskReqBody["title"],
				"isActive":                 createTaskReqBody["isActive"],
				"description":              "Updated description",
				"schedule":                 createTaskReqBody["schedule"],
				"priority":                 createTaskReqBody["priority"],
				"completionStatus":         createTaskReqBody["completionStatus"],
				"dueDate":                  createTaskReqBody["dueDate"],
				"shouldBeScored":           createTaskReqBody["shouldBeScored"],
				"score":                    createTaskReqBody["score"],
				"timeOfDay":                createTaskReqBody["timeOfDay"],
				"repetitiveTaskTemplateId": createTaskReqBody["repetitiveTaskTemplateID"],
				"createdAt":                createTaskReqBody["createdAt"],
				"modifiedAt":               time.Now().UTC().Format(time.RFC3339Nano), // Update modifiedAt
				"tags":                     createTaskReqBody["tags"],
				"spaceId":                  createTaskReqBody["spaceID"],
			},
			expectedStatus: http.StatusOK,
			expectedMsg:    messages.MsgTaskUpdateSuccess,
		},
		{
			name: "Failure - Missing ModifiedAt",
			requestBody: map[string]interface{}{
				"title":                    createTaskReqBody["title"],
				"isActive":                 createTaskReqBody["isActive"],
				"description":              createTaskReqBody["description"],
				"schedule":                 createTaskReqBody["schedule"],
				"priority":                 createTaskReqBody["priority"],
				"completionStatus":         createTaskReqBody["completionStatus"],
				"dueDate":                  createTaskReqBody["dueDate"],
				"shouldBeScored":           createTaskReqBody["shouldBeScored"],
				"score":                    createTaskReqBody["score"],
				"timeOfDay":                createTaskReqBody["timeOfDay"],
				"repetitiveTaskTemplateId": createTaskReqBody["repetitiveTaskTemplateID"],
				"createdAt":                createTaskReqBody["createdAt"],
				"tags":                     createTaskReqBody["tags"],
				"spaceId":                  createTaskReqBody["spaceID"],
			},
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    messages.ErrTaskUpdateFailed,
		},
		{
			name: "Failure - Invalid ModifiedAt Format",
			requestBody: map[string]interface{}{
				"title":                    createTaskReqBody["title"],
				"isActive":                 createTaskReqBody["isActive"],
				"description":              createTaskReqBody["description"],
				"schedule":                 createTaskReqBody["schedule"],
				"priority":                 createTaskReqBody["priority"],
				"completionStatus":         createTaskReqBody["completionStatus"],
				"dueDate":                  createTaskReqBody["dueDate"],
				"shouldBeScored":           createTaskReqBody["shouldBeScored"],
				"score":                    createTaskReqBody["score"],
				"timeOfDay":                createTaskReqBody["timeOfDay"],
				"repetitiveTaskTemplateId": createTaskReqBody["repetitiveTaskTemplateID"],
				"createdAt":                createTaskReqBody["createdAt"],
				"modifiedAt":               "invalid-date",
				"tags":                     createTaskReqBody["tags"],
				"spaceId":                  createTaskReqBody["spaceID"],
			},
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    messages.ErrTaskUpdateFailed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := testutils.CreateRequest(
				http.MethodPut,
				fmt.Sprintf("/tasks/%s", taskID),
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
				assert.Contains(t, resp.Body.String(), tc.expectedMsg, "Expected error message not found")
				var responseBody map[string]interface{}
				err = json.Unmarshal(resp.Body.Bytes(), &responseBody)
				assert.NoError(t, err)
				result, ok := responseBody["result"].(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, messages.Success, result["status"])
				assert.Equal(t, messages.MsgTaskUpdateSuccess, result["message"])
				data, ok := result["data"].(map[string]interface{})
				assert.True(t, ok)
				assert.NotEmpty(t, data["id"])

				// Retrieve the updated task from the database
				updatedTask := &models.Task{ID: uuid.MustParse(taskID)}
				err = TestDB.First(updatedTask).Error
				assert.NoError(t, err)

				// Check if the fields were updated correctly
				if title, ok := tc.requestBody["title"].(string); ok {
					assert.Equal(t, title, updatedTask.Title)
				}
				if description, ok := tc.requestBody["description"].(string); ok {
					assert.Equal(t, description, updatedTask.Description)
				}
			}
		})
	}
}

func TestCreateRepetitiveTaskTemplateIntegration(t *testing.T) {
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
			name: "Success - Valid Repetitive Task Template Creation",
			requestBody: map[string]interface{}{
				"isActive":                 true,
				"title":                    "Test Repetitive Task Template",
				"description":              "This is a test repetitive task template",
				"schedule":                 "Specific Days in a Week",
				"priority":                 3,
				"shouldBeScored":           true,
				"monday":                   true,
				"tuesday":                  false,
				"wednesday":                true,
				"thursday":                 false,
				"friday":                   true,
				"saturday":                 false,
				"sunday":                   true,
				"timeOfDay":                "morning",
				"lastDateOfTaskGeneration": time.Now().Add(7 * 24 * time.Hour).UTC().Format(time.RFC3339Nano),
				"createdAt":                time.Now().UTC().Format(time.RFC3339Nano),
				"modifiedAt":               time.Now().UTC().Format(time.RFC3339Nano),
				"tags":                     nil,
				"spaceID":                  nil,
			},
			expectedStatus: http.StatusOK,
			expectedErrMsg: messages.MsgRepetitiveTaskTemplateCreationSuccess,
		},
		{
			name: "Failure - Missing Title",
			requestBody: map[string]interface{}{
				"isActive":                 true,
				"description":              "This is a test repetitive task template",
				"schedule":                 "Daily",
				"priority":                 3,
				"shouldBeScored":           true,
				"monday":                   true,
				"tuesday":                  true,
				"wednesday":                true,
				"thursday":                 true,
				"friday":                   true,
				"saturday":                 true,
				"sunday":                   true,
				"timeOfDay":                "morning",
				"lastDateOfTaskGeneration": time.Now().Add(7 * 24 * time.Hour).UTC().Format(time.RFC3339Nano),
				"createdAt":                time.Now().UTC().Format(time.RFC3339Nano),
				"modifiedAt":               time.Now().UTC().Format(time.RFC3339Nano),
				"tags":                     nil,
				"spaceID":                  nil,
			},
			expectedStatus: http.StatusBadRequest,
			expectedErrMsg: messages.ErrRepetitiveTaskTemplateCreationFailed,
		},
		{
			name: "Failure - Invalid Last Date Of Task Generation Format",
			requestBody: map[string]interface{}{
				"isActive":                 true,
				"title":                    "Test Repetitive Task Template",
				"description":              "This is a test repetitive task template",
				"schedule":                 "Daily",
				"priority":                 3,
				"shouldBeScored":           true,
				"monday":                   true,
				"tuesday":                  false,
				"wednesday":                true,
				"thursday":                 false,
				"friday":                   true,
				"saturday":                 false,
				"sunday":                   true,
				"timeOfDay":                "Morning",
				"lastDateOfTaskGeneration": "invalid-date",
				"createdAt":                time.Now().UTC().Format(time.RFC3339Nano),
				"modifiedAt":               time.Now().UTC().Format(time.RFC3339Nano),
				"tags":                     nil,
				"spaceID":                  nil,
			},
			expectedStatus: http.StatusBadRequest,
			expectedErrMsg: messages.ErrRepetitiveTaskTemplateCreationFailed,
		},
		{
			name: "Failure - Invalid Space ID",
			requestBody: map[string]interface{}{
				"isActive":                 true,
				"title":                    "Test Repetitive Task Template",
				"description":              "This is a test repetitive task template",
				"schedule":                 "Daily",
				"priority":                 3,
				"shouldBeScored":           true,
				"monday":                   true,
				"tuesday":                  false,
				"wednesday":                true,
				"thursday":                 false,
				"friday":                   true,
				"saturday":                 false,
				"sunday":                   true,
				"timeOfDay":                "Morning",
				"lastDateOfTaskGeneration": time.Now().Add(7 * 24 * time.Hour).UTC().Format(time.RFC3339Nano),
				"createdAt":                time.Now().UTC().Format(time.RFC3339Nano),
				"modifiedAt":               time.Now().UTC().Format(time.RFC3339Nano),
				"tags":                     nil,
				"spaceID":                  "invalid-space-id",
			},
			expectedStatus: http.StatusBadRequest,
			expectedErrMsg: messages.ErrRepetitiveTaskTemplateCreationFailed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := testutils.CreateRequest(
				http.MethodPost,
				"/tasks/repetitive",
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
				assert.Contains(t, resp.Body.String(), tc.expectedErrMsg, "Expected error message not found")
				var responseBody map[string]interface{}
				err = json.Unmarshal(resp.Body.Bytes(), &responseBody)
				assert.NoError(t, err)
				result, ok := responseBody["result"].(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, messages.Success, result["status"])
				assert.Equal(t, messages.MsgRepetitiveTaskTemplateCreationSuccess, result["message"])
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
				assert.Equal(t, tc.requestBody["shouldBeScored"], data["shouldBeScored"])
				assert.Equal(t, tc.requestBody["monday"], data["monday"])
				assert.Equal(t, tc.requestBody["tuesday"], data["tuesday"])
				assert.Equal(t, tc.requestBody["wednesday"], data["wednesday"])
				assert.Equal(t, tc.requestBody["thursday"], data["thursday"])
				assert.Equal(t, tc.requestBody["friday"], data["friday"])
				assert.Equal(t, tc.requestBody["saturday"], data["saturday"])
				assert.Equal(t, tc.requestBody["sunday"], data["sunday"])
				assert.Equal(t, tc.requestBody["timeOfDay"], data["timeOfDay"])
				assert.Equal(t, tc.requestBody["lastDateOfTaskGeneration"], data["lastDateOfTaskGeneration"])
				assert.Equal(t, tc.requestBody["createdAt"], data["createdAt"])
				assert.Equal(t, tc.requestBody["modifiedAt"], data["modifiedAt"])
				assert.Equal(t, tc.requestBody["tags"], data["tags"])
				assert.Equal(t, tc.requestBody["spaceId"], data["spaceId"])
				assert.NotEmpty(t, data["userId"])
			}
		})
	}
}
