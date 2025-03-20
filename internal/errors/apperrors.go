package apperrors

type AppError interface {
	Error() string
	StatusCode() int
	LogError() string
}
