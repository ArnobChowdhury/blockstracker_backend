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

func (e *CommonError) Code() string {
	return e.code
}

func NewInvalidReqErr(customMessage ...string) *CommonError {
	message := "Invalid request"
	if len(customMessage) > 0 {
		message = customMessage[0]
	}
	return NewCommonError("BAD_REQUEST", message, http.StatusBadRequest)
}

var (
	ErrUnauthorized            = NewCommonError("UNAUTHORIZED", "Unauthorized", http.StatusUnauthorized)
	ErrInternalServerError     = NewCommonError("INTERNAL_SERVER_ERROR", "Internal server error", http.StatusInternalServerError)
	ErrUserIDNotFoundInContext = NewCommonError("USER_ID_NOT_FOUND_IN_CONTEXT", "User ID not found in context", http.StatusInternalServerError)
	ErrUserIDNotValidType      = NewCommonError("USER_ID_NOT_VALID_TYPE", "User ID is not of valid type", http.StatusInternalServerError)
	ErrStaleData               = NewCommonError("STALE_DATA", "Stale data", http.StatusConflict)
	ErrDuplicateEntity         = NewCommonError("DUPLICATE_ENTITY", "Duplicate entity found", http.StatusConflict)
)
