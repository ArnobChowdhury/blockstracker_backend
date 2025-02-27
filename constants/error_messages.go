package responsemsg

const (
	UserCreationFailed                = "User creation failed"
	InvalidEmail                      = "Invalid email address"
	UnexpectedErrorDuringUserCreation = "Unexpected error during user creation"
	NotStrongPassword                 = "Password must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter and one number"
	MalformedRequest                  = "Malformed Request"
	InternalServerError               = "Internal Server Error"
	InvalidTaskError                  = "Invalid task data provided"
	TaskNotFoundError                 = "Task not found"
	DatabaseConnectionError           = "Failed to connect to the database"
)
