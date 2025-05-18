package errors

import (
	"fmt"
	"testing"
)

func TestErrorCategories(t *testing.T) {
	tests := []struct {
		category Category
		expected string
	}{
		{CategoryUnknown, "UNKNOWN"},
		{CategoryConfiguration, "CONFIGURATION"},
		{CategoryValidation, "VALIDATION"},
		{CategoryExternal, "EXTERNAL"},
		{CategorySystem, "SYSTEM"},
		{CategoryApplication, "APPLICATION"},
		{Category(999), "CATEGORY(999)"},
	}

	for _, test := range tests {
		if test.category.String() != test.expected {
			t.Errorf("expected category %d to be %s, got %s", test.category, test.expected, test.category.String())
		}
	}
}

func TestNew(t *testing.T) {
	err := New(CategoryValidation, "invalid input")
	if err.Error() != "invalid input" {
		t.Errorf("expected error message to be 'invalid input', got %s", err.Error())
	}
	if err.Category() != CategoryValidation {
		t.Errorf("expected category to be VALIDATION, got %s", err.Category())
	}
	if len(err.Context()) != 0 {
		t.Errorf("expected context to be empty, got %v", err.Context())
	}
}

func TestWrap(t *testing.T) {
	original := fmt.Errorf("original error")
	err := Wrap(CategoryExternal, original, "wrapped error")
	if err.Error() != "wrapped error: original error" {
		t.Errorf("expected error message to be 'wrapped error: original error', got %s", err.Error())
	}
	if err.Category() != CategoryExternal {
		t.Errorf("expected category to be EXTERNAL, got %s", err.Category())
	}

	// Test unwrapping
	unwrapped := err.Unwrap()
	if unwrapped.Error() != "wrapped error: original error" {
		t.Errorf("expected unwrapped error to be 'wrapped error: original error', got %s", unwrapped.Error())
	}

	// Test wrapping nil
	nilErr := Wrap(CategoryExternal, nil, "wrapped nil")
	if nilErr != nil {
		t.Errorf("expected wrapping nil to return nil, got %v", nilErr)
	}
}

func TestWithContext(t *testing.T) {
	err := New(CategoryValidation, "invalid input")
	err = err.WithContext("field", "username").WithContext("max_length", 50)

	if err.Context()["field"] != "username" {
		t.Errorf("expected context[field] to be 'username', got %v", err.Context()["field"])
	}
	if err.Context()["max_length"] != 50 {
		t.Errorf("expected context[max_length] to be 50, got %v", err.Context()["max_length"])
	}
}

func TestIs(t *testing.T) {
	// Test with standard errors
	err := Wrap(CategoryValidation, ErrInvalidInput, "wrapped invalid input")
	if !Is(err, ErrInvalidInput) {
		t.Errorf("expected err to be ErrInvalidInput")
	}
	if Is(err, ErrNotFound) {
		t.Errorf("expected err not to be ErrNotFound")
	}

	// Test with custom errors
	original := New(CategoryValidation, "original")
	wrapped := Wrap(CategoryApplication, original, "wrapped")
	if !Is(wrapped, original) {
		t.Errorf("expected wrapped to be original")
	}
}

func TestAs(t *testing.T) {
	original := New(CategoryValidation, "original")
	err := Wrap(CategoryApplication, original, "wrapped")

	var customErr *Error
	if !As(err, &customErr) {
		t.Errorf("expected As to succeed")
	}
	if customErr.Category() != CategoryApplication {
		t.Errorf("expected category to be APPLICATION, got %s", customErr.Category())
	}

	// Test with standard error
	stdErr := fmt.Errorf("standard error")
	if As(stdErr, &customErr) {
		t.Errorf("expected As to fail with standard error")
	}
}

func TestGetCategory(t *testing.T) {
	err := New(CategoryValidation, "invalid input")
	if GetCategory(err) != CategoryValidation {
		t.Errorf("expected category to be VALIDATION, got %s", GetCategory(err))
	}

	stdErr := fmt.Errorf("standard error")
	if GetCategory(stdErr) != CategoryUnknown {
		t.Errorf("expected category to be UNKNOWN for standard error, got %s", GetCategory(stdErr))
	}
}

func TestGetContext(t *testing.T) {
	err := New(CategoryValidation, "invalid input").WithContext("field", "username")
	ctx := GetContext(err)
	if ctx["field"] != "username" {
		t.Errorf("expected context[field] to be 'username', got %v", ctx["field"])
	}

	stdErr := fmt.Errorf("standard error")
	if GetContext(stdErr) != nil {
		t.Errorf("expected context to be nil for standard error, got %v", GetContext(stdErr))
	}
}
