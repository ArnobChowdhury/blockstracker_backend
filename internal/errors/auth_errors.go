package apperrors

import (
	"fmt"
	"net/http"
)

type AuthError struct {
	code       string
	message    string
	statusCode int
}

func NewAuthError(code, message string, statusCode int) *AuthError {
	return &AuthError{
		code:       code,
		message:    message,
		statusCode: statusCode,
	}
}

func (e *AuthError) StatusCode() int {
	return e.statusCode
}

func (e *AuthError) Error() string {
	return e.message
}

func (e *AuthError) LogError() string {
	return fmt.Sprintf("AuthError - Code: %s, Message: %s, Status Code: %d", e.code, e.message, e.statusCode)
}

var (
	ErrUnauthorized             = NewAuthError("UNAUTHORIZED", "Unauthorized", http.StatusUnauthorized)
	ErrNoAuthorizationHeader    = NewAuthError("NO_AUTH_HEADER", "No Authorization header", http.StatusUnauthorized)
	ErrInvalidAuthHeader        = NewAuthError("INVALID_AUTH_HEADER", "Invalid Authorization header", http.StatusUnauthorized)
	ErrTokenExpired             = NewAuthError("TOKEN_EXPIRED", "Token expired", http.StatusUnauthorized)
	ErrInvalidToken             = NewAuthError("INVALID_TOKEN", "Invalid token", http.StatusUnauthorized)
	ErrUnexpectedSigningMethod  = NewAuthError("UNEXPECTED_SIGNING_METHOD", "Unexpected signing method", http.StatusUnauthorized)
	ErrMalformedRequest         = NewAuthError("BAD_REQUEST", "Malformed request", http.StatusBadRequest)
	ErrDBUniqueConstraintFailed = NewAuthError("DB_UNIQUE_CONSTRAINT_FAILED", "Internal server error", http.StatusInternalServerError)
	ErrUniqueConstraintFailed   = NewAuthError("UNIQUE_CONSTRAINT_FAILED", "Unique constraint failed", http.StatusBadRequest)
	ErrInvalidRequestBody       = NewAuthError("INVALID_REQUEST_BODY", "Invalid request body", http.StatusBadRequest)
	ErrUserCreationFailed       = NewAuthError("USER_CREATION_FAILED", "User creation failed", http.StatusBadRequest)
	ErrInvalidCredentials       = NewAuthError("INVALID_CREDENTIALS", "Invalid credentials", http.StatusUnauthorized)
)
