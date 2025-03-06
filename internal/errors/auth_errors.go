package apperrors

import "fmt"

type AuthError struct {
	code    string
	message string
}

func (e *AuthError) Error() string {
	return e.message
}

func (e *AuthError) LogError() string {
	return fmt.Sprintf("AuthError - Code: %s, Message: %s", e.code, e.message)
}

var (
	ErrUnauthorized            = &AuthError{"UNAUTHORIZED", "Unauthorized"}
	ErrNoAuthorizationHeader   = &AuthError{"NO_AUTH_HEADER", "No Authorization header"}
	ErrInvalidAuthHeader       = &AuthError{"INVALID_AUTH_HEADER", "Invalid Authorization header"}
	ErrTokenExpired            = &AuthError{"TOKEN_EXPIRED", "Token expired"}
	ErrInvalidToken            = &AuthError{"INVALID_TOKEN", "Invalid token"}
	ErrUnexpectedSigningMethod = &AuthError{"UNEXPECTED_SIGNING_METHOD", "Unexpected signing method"}
)
