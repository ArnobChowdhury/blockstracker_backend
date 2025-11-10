package apperrors

import (
	"fmt"
	"net/http"
)

type TaskError struct {
	code       string
	message    string
	statusCode int
}

func NewTaskError(code, message string, statusCode int) *TaskError {
	return &TaskError{
		code:       code,
		message:    message,
		statusCode: statusCode,
	}
}

func (e *TaskError) StatusCode() int {
	return e.statusCode
}

func (e *TaskError) Error() string {
	return e.message
}

func (e *TaskError) LogError() string {
	return fmt.Sprintf("TaskError - Code: %s, Message: %s, Status Code: %d", e.code, e.message, e.statusCode)
}

func (e *TaskError) Code() string {
	return e.code
}

var (
	ErrMalformedTaskRequest                   = NewTaskError("BAD_REQUEST", "Malformed request", http.StatusBadRequest)
	ErrMalformedRepetitiveTaskTemplateRequest = NewTaskError("BAD_REQUEST", "Malformed request", http.StatusBadRequest)
)
