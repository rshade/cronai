// Package errors provides error handling utilities.
package errors

import (
	"errors"
	"fmt"

	"github.com/rshade/cronai/internal/logger"
)

// Standard error types
var (
	ErrNotFound       = errors.New("resource not found")
	ErrInvalidInput   = errors.New("invalid input")
	ErrInternal       = errors.New("internal error")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden")
	ErrUnavailable    = errors.New("service unavailable")
	ErrTimeout        = errors.New("operation timed out")
	ErrNotImplemented = errors.New("not implemented")
)

// Category represents the category of an error
type Category int

const (
	// CategoryUnknown is used when the error category is not known
	CategoryUnknown Category = iota
	// CategoryConfiguration is used for configuration-related errors
	CategoryConfiguration
	// CategoryValidation is used for input validation errors
	CategoryValidation
	// CategoryExternal is used for errors from external services
	CategoryExternal
	// CategorySystem is used for system-level errors
	CategorySystem
	// CategoryApplication is used for application-level errors
	CategoryApplication
	// CategorySecurity is used for security-related errors
	CategorySecurity
)

// String returns the string representation of the error category.
func (c Category) String() string {
	switch c {
	case CategoryUnknown:
		return "UNKNOWN"
	case CategoryConfiguration:
		return "CONFIGURATION"
	case CategoryValidation:
		return "VALIDATION"
	case CategoryExternal:
		return "EXTERNAL"
	case CategorySystem:
		return "SYSTEM"
	case CategoryApplication:
		return "APPLICATION"
	case CategorySecurity:
		return "SECURITY"
	default:
		return fmt.Sprintf("CATEGORY(%d)", c)
	}
}

// Error is a custom error with additional context information.
type Error struct {
	category Category
	err      error
	context  logger.Fields
}

// New creates a new Error with the given category and message.
func New(category Category, msg string) *Error {
	return &Error{
		category: category,
		err:      errors.New(msg),
		context:  make(logger.Fields),
	}
}

// Wrap wraps an existing error with a category and additional context.
func Wrap(category Category, err error, msg string) *Error {
	if err == nil {
		return nil
	}
	return &Error{
		category: category,
		err:      fmt.Errorf("%s: %w", msg, err),
		context:  make(logger.Fields),
	}
}

// Error returns the error message.
func (e *Error) Error() string {
	return e.err.Error()
}

// Unwrap returns the wrapped error.
func (e *Error) Unwrap() error {
	return e.err
}

// Category returns the error category.
func (e *Error) Category() Category {
	return e.category
}

// WithContext adds context information to the error.
func (e *Error) WithContext(key string, value interface{}) *Error {
	e.context[key] = value
	return e
}

// Context returns the error context information.
func (e *Error) Context() logger.Fields {
	return e.context
}

// Is reports whether any error in err's chain matches target.
func Is(err error, target error) bool {
	return errors.Is(err, target)
}

// As finds the first error in err's chain that matches target and sets target to that error value.
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// GetCategory returns the category of an error.
// If the error is not a *Error, CategoryUnknown is returned.
func GetCategory(err error) Category {
	var e *Error
	if errors.As(err, &e) {
		return e.Category()
	}
	return CategoryUnknown
}

// GetContext returns the context of an error.
// If the error is not a *Error, nil is returned.
func GetContext(err error) logger.Fields {
	var e *Error
	if errors.As(err, &e) {
		return e.Context()
	}
	return nil
}
