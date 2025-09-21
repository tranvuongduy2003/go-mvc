package errors

import (
	"fmt"
	"net/http"
)

// ErrorType định nghĩa các loại lỗi
type ErrorType string

const (
	ErrorTypeValidation   ErrorType = "VALIDATION_ERROR"
	ErrorTypeNotFound     ErrorType = "NOT_FOUND"
	ErrorTypeConflict     ErrorType = "CONFLICT"
	ErrorTypeUnauthorized ErrorType = "UNAUTHORIZED"
	ErrorTypeForbidden    ErrorType = "FORBIDDEN"
	ErrorTypeInternal     ErrorType = "INTERNAL_ERROR"
)

// AppError là custom error type cho ứng dụng
type AppError struct {
	Type    ErrorType `json:"type"`
	Message string    `json:"message"`
	Code    int       `json:"code"`
	Cause   error     `json:"-"` // Không serialize cause
}

// Error implements error interface
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Type, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Unwrap returns the underlying cause
func (e *AppError) Unwrap() error {
	return e.Cause
}

// HTTPStatusCode returns HTTP status code for the error
func (e *AppError) HTTPStatusCode() int {
	return e.Code
}

// Constructor functions for different error types

func NewValidationError(message string, cause error) *AppError {
	return &AppError{
		Type:    ErrorTypeValidation,
		Message: message,
		Code:    http.StatusBadRequest,
		Cause:   cause,
	}
}

func NewNotFoundError(message string) *AppError {
	return &AppError{
		Type:    ErrorTypeNotFound,
		Message: message,
		Code:    http.StatusNotFound,
	}
}

func NewConflictError(message string, cause error) *AppError {
	return &AppError{
		Type:    ErrorTypeConflict,
		Message: message,
		Code:    http.StatusConflict,
		Cause:   cause,
	}
}

func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Type:    ErrorTypeUnauthorized,
		Message: message,
		Code:    http.StatusUnauthorized,
	}
}

func NewForbiddenError(message string) *AppError {
	return &AppError{
		Type:    ErrorTypeForbidden,
		Message: message,
		Code:    http.StatusForbidden,
	}
}

func NewInternalError(message string, cause error) *AppError {
	return &AppError{
		Type:    ErrorTypeInternal,
		Message: message,
		Code:    http.StatusInternalServerError,
		Cause:   cause,
	}
}
