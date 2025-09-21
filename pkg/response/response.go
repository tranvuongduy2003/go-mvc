package response

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	apperrors "github.com/tranvuongduy2003/go-mvc/pkg/errors"
	"github.com/tranvuongduy2003/go-mvc/pkg/pagination"
)

// APIResponse represents a standard API response
type APIResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Error     *ErrorInfo  `json:"error,omitempty"`
	Meta      *Meta       `json:"meta,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// ErrorInfo contains error details
type ErrorInfo struct {
	Type    string            `json:"type"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
	Code    string            `json:"code,omitempty"`
}

// Meta contains metadata for the response
type Meta struct {
	Pagination *pagination.Pagination `json:"pagination,omitempty"`
	Total      int64                  `json:"total,omitempty"`
	RequestID  string                 `json:"request_id,omitempty"`
}

// Success sends a successful response
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Data:      data,
		Timestamp: time.Now().UTC(),
	})
}

// SuccessWithMessage sends a successful response with message
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UTC(),
	})
}

// SuccessWithPagination sends a successful response with pagination
func SuccessWithPagination(c *gin.Context, data interface{}, pagination *pagination.Pagination) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    data,
		Meta: &Meta{
			Pagination: pagination,
		},
		Timestamp: time.Now().UTC(),
	})
}

// Created sends a created response
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, APIResponse{
		Success:   true,
		Data:      data,
		Timestamp: time.Now().UTC(),
	})
}

// Error sends an error response
func Error(c *gin.Context, err error) {
	var appErr *apperrors.AppError
	var statusCode int
	var errorInfo *ErrorInfo

	if errors.As(err, &appErr) {
		statusCode = appErr.HTTPStatusCode()
		errorInfo = &ErrorInfo{
			Type:    string(appErr.Type),
			Message: appErr.Message,
		}
	} else {
		statusCode = http.StatusInternalServerError
		errorInfo = &ErrorInfo{
			Type:    string(apperrors.ErrorTypeInternal),
			Message: "Internal server error",
		}
	}

	c.JSON(statusCode, APIResponse{
		Success:   false,
		Error:     errorInfo,
		Timestamp: time.Now().UTC(),
	})
}

// ValidationError sends a validation error response
func ValidationError(c *gin.Context, validationErrors map[string]string) {
	c.JSON(http.StatusBadRequest, APIResponse{
		Success: false,
		Error: &ErrorInfo{
			Type:    string(apperrors.ErrorTypeValidation),
			Message: "Validation failed",
			Details: validationErrors,
		},
		Timestamp: time.Now().UTC(),
	})
}
