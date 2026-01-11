package roadmap

import (
	"errors"
	"fmt"
)

// Sentinel errors for roadmap operations.
var (
	// ErrInvalidIRVersion indicates an unsupported IR version.
	ErrInvalidIRVersion = errors.New("invalid or unsupported IR version")

	// ErrMissingRequiredField indicates a required field is missing.
	ErrMissingRequiredField = errors.New("required field is missing")

	// ErrDuplicateID indicates a duplicate ID was found.
	ErrDuplicateID = errors.New("duplicate ID")

	// ErrInvalidStatus indicates an invalid status value.
	ErrInvalidStatus = errors.New("invalid status")

	// ErrInvalidReference indicates a reference to a non-existent item.
	ErrInvalidReference = errors.New("invalid reference")

	// ErrInvalidFormat indicates an invalid format for a field value.
	ErrInvalidFormat = errors.New("invalid format")

	// ErrInvalidType indicates an invalid change type.
	ErrInvalidType = errors.New("invalid change type")

	// ErrParseJSON indicates a JSON parsing error.
	ErrParseJSON = errors.New("failed to parse JSON")

	// ErrReadFile indicates a file read error.
	ErrReadFile = errors.New("failed to read file")

	// ErrWriteFile indicates a file write error.
	ErrWriteFile = errors.New("failed to write file")
)

// ParseError wraps a parsing error with context.
type ParseError struct {
	Op  string // Operation (e.g., "read", "parse", "unmarshal")
	Err error  // Underlying error
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("%s: %v", e.Op, e.Err)
}

func (e *ParseError) Unwrap() error {
	return e.Err
}

// FieldError represents a validation error for a specific field.
type FieldError struct {
	Field   string // Field path (e.g., "items[0].id")
	Message string // Error message
	Err     error  // Underlying sentinel error
}

func (e *FieldError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

func (e *FieldError) Unwrap() error {
	return e.Err
}

// NewFieldError creates a new FieldError with the given sentinel error.
func NewFieldError(field, message string, sentinel error) *FieldError {
	return &FieldError{
		Field:   field,
		Message: message,
		Err:     sentinel,
	}
}
