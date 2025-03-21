package apperrors

import (
	"fmt"
	"net/http"
)

type CommonError struct {
	code       string
	message    string
	statusCode int
}

func NewCommonError(code, message string, statusCode int) *CommonError {
	return &CommonError{
		code:       code,
		message:    message,
		statusCode: statusCode,
	}
}

func (e *CommonError) StatusCode() int {
	return e.statusCode
}

func (e *CommonError) Error() string {
	return e.message
}

func (e *CommonError) LogError() string {
	return fmt.Sprintf("CommonError - Code: %s, Message: %s, Status Code: %d", e.code, e.message, e.statusCode)
}

func NewInvalidReqErr(customMessage ...string) *CommonError {
	message := "Invalid request"
	if len(customMessage) > 0 {
		message = customMessage[0]
	}
	return NewCommonError("BAD_REQUEST", message, http.StatusBadRequest)
}

var (
	ErrInternalServerError = NewCommonError("INTERNAL_SERVER_ERROR", "Internal server error", http.StatusInternalServerError)
)
