package common

type ErrorType string

const (
	NotFound      ErrorType = "NOT_FOUND"
	ValidationErr ErrorType = "VALIDATION_ERROR"
	DatabaseErr   ErrorType = "DATABASE_ERROR"
	InternalErr   ErrorType = "INTERNAL_ERROR"
)

type AppError struct {
	Type    ErrorType
	Message string
	Err     error
}

func (e AppError) Error() string {
	return e.Message
}

// Error constructors
func NewNotFoundError(message string) AppError {
	return AppError{
		Type:    NotFound,
		Message: message,
	}
}

func NewValidationError(message string) AppError {
	return AppError{
		Type:    ValidationErr,
		Message: message,
	}
}

func NewDatabaseError(err error) AppError {
	return AppError{
		Type:    DatabaseErr,
		Message: "Database error occurred",
		Err:     err,
	}
}

func NewInternalError(err error) AppError {
	return AppError{
		Type:    InternalErr,
		Message: "Internal server error",
		Err:     err,
	}
}
