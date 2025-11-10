package apperrors

import (
	"fmt"
	"net/http"
)

type SpaceError struct {
	code       string
	message    string
	statusCode int
}

func NewSpaceError(code, message string, statusCode int) *SpaceError {
	return &SpaceError{
		code:       code,
		message:    message,
		statusCode: statusCode,
	}
}

func (e *SpaceError) StatusCode() int {
	return e.statusCode
}

func (e *SpaceError) Error() string {
	return e.message
}

func (e *SpaceError) LogError() string {
	return fmt.Sprintf("SpaceError - Code: %s, Message: %s, Status Code: %d", e.code, e.message, e.statusCode)
}

func (e *SpaceError) Code() string {
	return e.code
}

var (
	ErrSpaceDuplicateKey = NewSpaceError("DUPLICATE_NAME_FOR_SPACE", "Duplicate name for space", http.StatusBadRequest)
)
