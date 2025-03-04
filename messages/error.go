package messages

const (
	Error                                    = "Error"
	ErrUserCreationFailed                    = "User creation failed"
	ErrInvalidEmail                          = "Invalid email address"
	ErrUnexpectedErrorDuringUserCreation     = "Unexpected error during user creation"
	ErrNotStrongPassword                     = "Password must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter and one number"
	ErrMalformedRequest                      = "Malformed Request"
	ErrInternalServerError                   = "Internal Server Error"
	ErrInvalidTaskError                      = "Invalid task data provided"
	ErrTaskNotFoundError                     = "Task not found"
	ErrDatabaseConnectionError               = "Failed to connect to the database"
	ErrHashingPassword                       = "Error hashing password"
	ErrInvalidRequestBody                    = "Invalid request body"
	ErrUniqueConstraintFailed                = "Unique constraint failed"
	ErrEmailNotFoundDuringSignIn             = "Email not found during sign in"
	ErrUnexpectedErrorDuringUserRetrieval    = "Unexpected error during user retrieval"
	ErrInvalidCredentials                    = "Invalid credentials"
	ErrJWTAccessSecretNotFoundInEnvironment  = "JWT_ACCESS_SECRET not found in environment variables"
	ErrJWTRefreshSecretNotFoundInEnvironment = "JWT_REFRESH_SECRET not found in environment variables"
	ErrGeneratingJWT                         = "Failed to generate JWT"
	ErrMismatchingPasswordDuringSignIn       = "Mismatching password during sign in"
)
