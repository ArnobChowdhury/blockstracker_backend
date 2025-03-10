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
	ErrUnauthorized            = NewAuthError("UNAUTHORIZED", "Unauthorized", http.StatusUnauthorized)
	ErrNoAuthorizationHeader   = NewAuthError("NO_AUTH_HEADER", "No Authorization header", http.StatusUnauthorized)
	ErrInvalidAuthHeader       = NewAuthError("INVALID_AUTH_HEADER", "Invalid Authorization header", http.StatusUnauthorized)
	ErrTokenExpired            = NewAuthError("TOKEN_EXPIRED", "Token expired", http.StatusUnauthorized)
	ErrInvalidToken            = NewAuthError("INVALID_TOKEN", "Invalid token", http.StatusUnauthorized)
	ErrUnexpectedSigningMethod = NewAuthError("UNEXPECTED_SIGNING_METHOD", "Unexpected signing method", http.StatusUnauthorized)
	ErrRedisSet                = NewAuthError("REDIS_SET", "Redis set error", http.StatusInternalServerError)
	ErrRedisGet                = NewAuthError("REDIS_GET", "Redis get error", http.StatusInternalServerError)
)
