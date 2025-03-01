package messages

const (
	Error                                = "Error"
	ErrUserCreationFailed                = "User creation failed"
	ErrInvalidEmail                      = "Invalid email address"
	ErrUnexpectedErrorDuringUserCreation = "Unexpected error during user creation"
	ErrNotStrongPassword                 = "Password must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter and one number"
	ErrMalformedRequest                  = "Malformed Request"
	ErrInternalServerError               = "Internal Server Error"
	ErrInvalidTaskError                  = "Invalid task data provided"
	ErrTaskNotFoundError                 = "Task not found"
	ErrDatabaseConnectionError           = "Failed to connect to the database"
	ErrHashingPassword                   = "Error hashing password"
	ErrInvalidRequestBody                = "Invalid request body"
	ErrUniqueConstraintFailed            = "Unique constraint failed"
)
