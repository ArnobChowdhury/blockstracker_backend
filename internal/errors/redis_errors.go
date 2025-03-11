package apperrors

import (
	"fmt"
	"net/http"
)

type RedisError struct {
	code       string
	message    string
	statusCode int
}

func NewRedisError(code, message string, statusCode int) *RedisError {
	return &RedisError{
		code:       code,
		message:    message,
		statusCode: statusCode,
	}
}

func (e *RedisError) StatusCode() int {
	return e.statusCode
}

func (e *RedisError) Error() string {
	return e.message
}

func (e *RedisError) LogError() string {
	return fmt.Sprintf("RedisError - Code: %s, Message: %s, Status Code: %d", e.code, e.message, e.statusCode)
}

var (
	ErrRedisSet         = NewRedisError("REDIS_SET", "Redis set error", http.StatusInternalServerError)
	ErrRedisKeyNotFound = NewRedisError("KEY_NOT_FOUND", "Key not found", http.StatusNotFound)
)
