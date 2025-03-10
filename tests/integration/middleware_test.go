package integration

import (
	apperrors "blockstracker_backend/internal/errors"
	"blockstracker_backend/internal/utils"
	"blockstracker_backend/models"
	"blockstracker_backend/tests/integration/testutils"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	user := &models.User{
		ID:    uuid.New(),
		Email: "test@test.com",
	}

	accessTokenClaims := utils.GetClaims(user, "access")
	accessToken, err := utils.GenerateJWT(accessTokenClaims, testAuthConfig.AccessSecret)
	if err != nil {
		t.Fatalf("Error generating access token: %v", err)
	}
	accessTokenClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(-time.Hour))
	expiredAccessToken, err := utils.GenerateJWT(accessTokenClaims, testAuthConfig.AccessSecret)

	if err != nil {
		t.Fatalf("Error generating access token: %v", err)
	}

	testCases := []struct {
		name             string
		authHeader       string
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:             "Valid Token",
			authHeader:       fmt.Sprintf("Bearer %s", accessToken),
			expectedStatus:   http.StatusOK,
			expectedResponse: `{"message":"success"}`,
		},
		{
			name:             "No Auth Header",
			authHeader:       "",
			expectedStatus:   http.StatusUnauthorized,
			expectedResponse: apperrors.ErrUnauthorized.Error(),
		},
		{
			name:             "Invalid Auth Header",
			authHeader:       "Bearer invalidtoken",
			expectedStatus:   http.StatusUnauthorized,
			expectedResponse: apperrors.ErrUnauthorized.Error(),
		},
		{
			name:             "Expired Token",
			authHeader:       fmt.Sprintf("Bearer %s", expiredAccessToken),
			expectedStatus:   http.StatusUnauthorized,
			expectedResponse: "Token expired",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := testutils.CreateRequest(http.MethodPost, "/protected", nil)
			if err != nil {
				t.Fatalf("Error creating request: %v", err)
			}
			req.Header.Set("Authorization", tc.authHeader)

			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, tc.expectedStatus, resp.Code)
			assert.Contains(t, resp.Body.String(), tc.expectedResponse, "Expected error message not found")
		})
	}
}
