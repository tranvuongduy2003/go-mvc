package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/ports/services"
	apperrors "github.com/tranvuongduy2003/go-mvc/pkg/errors"
)

const (
	// UserContextKey is the key used to store user info in Gin context
	UserContextKey = "user"

	// UserIDContextKey is the key used to store user ID in Gin context
	UserIDContextKey = "user_id"

	// AuthorizationHeader is the header name for authorization
	AuthorizationHeader = "Authorization"

	// BearerPrefix is the expected prefix for bearer tokens
	BearerPrefix = "Bearer "
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success   bool       `json:"success"`
	Error     *ErrorInfo `json:"error,omitempty"`
	Timestamp time.Time  `json:"timestamp"`
}

// ErrorInfo contains error details
type ErrorInfo struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

// AuthMiddleware provides authentication middleware functionality
type AuthMiddleware struct {
	authService services.AuthService
}

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware(authService services.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// RequireAuth middleware that requires valid authentication
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := m.extractTokenFromHeader(c)
		if err != nil {
			m.sendErrorResponse(c, http.StatusUnauthorized, "Authentication required", err.Error())
			return
		}

		user, err := m.authService.ValidateToken(c.Request.Context(), token)
		if err != nil {
			m.sendErrorResponse(c, http.StatusUnauthorized, "Invalid or expired token", err.Error())
			return
		}

		// Store user information in context for downstream handlers
		c.Set(UserContextKey, user)
		c.Set(UserIDContextKey, user.ID())

		c.Next()
	}
}

// OptionalAuth middleware that extracts user info if token is present
func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := m.extractTokenFromHeader(c)
		if err != nil {
			// No token or invalid format, continue without user context
			c.Next()
			return
		}

		user, err := m.authService.ValidateToken(c.Request.Context(), token)
		if err != nil {
			// Invalid token, continue without user context
			c.Next()
			return
		}

		// Store user information in context for downstream handlers
		c.Set(UserContextKey, user)
		c.Set(UserIDContextKey, user.ID())

		c.Next()
	}
}

// RequireActiveUser middleware that requires an active authenticated user
func (m *AuthMiddleware) RequireActiveUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := m.extractTokenFromHeader(c)
		if err != nil {
			m.sendErrorResponse(c, http.StatusUnauthorized, "Authentication required", err.Error())
			return
		}

		user, err := m.authService.ValidateToken(c.Request.Context(), token)
		if err != nil {
			m.sendErrorResponse(c, http.StatusUnauthorized, "Invalid or expired token", err.Error())
			return
		}

		// Check if user is active
		if !user.IsActive() {
			m.sendErrorResponse(c, http.StatusForbidden, "Account is inactive", "Your account has been deactivated. Please contact support.")
			return
		}

		// Store user information in context for downstream handlers
		c.Set(UserContextKey, user)
		c.Set(UserIDContextKey, user.ID())

		c.Next()
	}
}

// TokenRefresh middleware for refresh token endpoints
func (m *AuthMiddleware) TokenRefresh() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := m.extractTokenFromHeader(c)
		if err != nil {
			m.sendErrorResponse(c, http.StatusUnauthorized, "Refresh token required", err.Error())
			return
		}

		// Store the refresh token in context for the handler
		c.Set("refresh_token", token)

		c.Next()
	}
}

// sendErrorResponse sends a standardized error response
func (m *AuthMiddleware) sendErrorResponse(c *gin.Context, statusCode int, message string, details string) {
	var errorType string
	switch statusCode {
	case http.StatusUnauthorized:
		errorType = string(apperrors.ErrorTypeUnauthorized)
	case http.StatusForbidden:
		errorType = string(apperrors.ErrorTypeForbidden)
	default:
		errorType = string(apperrors.ErrorTypeInternal)
	}

	c.JSON(statusCode, ErrorResponse{
		Success: false,
		Error: &ErrorInfo{
			Type:    errorType,
			Message: message,
		},
		Timestamp: time.Now().UTC(),
	})
	c.Abort()
}

// extractTokenFromHeader extracts the JWT token from the Authorization header
func (m *AuthMiddleware) extractTokenFromHeader(c *gin.Context) (string, error) {
	authHeader := c.GetHeader(AuthorizationHeader)
	if authHeader == "" {
		return "", NewAuthError("missing authorization header")
	}

	if !strings.HasPrefix(authHeader, BearerPrefix) {
		return "", NewAuthError("invalid authorization header format")
	}

	token := strings.TrimPrefix(authHeader, BearerPrefix)
	if token == "" {
		return "", NewAuthError("empty token")
	}

	return token, nil
}

// AuthError represents authentication-related errors
type AuthError struct {
	Message string
}

func NewAuthError(message string) *AuthError {
	return &AuthError{Message: message}
}

func (e *AuthError) Error() string {
	return e.Message
}

// Helper functions for extracting user info from context

// GetUserFromContext extracts the authenticated user from Gin context
func GetUserFromContext(c *gin.Context) (interface{}, bool) {
	user, exists := c.Get(UserContextKey)
	return user, exists
}

// GetUserIDFromContext extracts the user ID from Gin context
func GetUserIDFromContext(c *gin.Context) (string, bool) {
	userID, exists := c.Get(UserIDContextKey)
	if !exists {
		return "", false
	}

	id, ok := userID.(string)
	return id, ok
}

// RequireUserID is a helper that extracts user ID and returns error if not found
func RequireUserID(c *gin.Context) (string, error) {
	userID, exists := GetUserIDFromContext(c)
	if !exists {
		return "", NewAuthError("user ID not found in context")
	}
	return userID, nil
}

// API Key Authentication (for service-to-service communication)

// APIKeyMiddleware provides API key authentication
type APIKeyMiddleware struct {
	validAPIKeys map[string]string // API key -> service name
}

// NewAPIKeyMiddleware creates a new API key middleware
func NewAPIKeyMiddleware(apiKeys map[string]string) *APIKeyMiddleware {
	return &APIKeyMiddleware{
		validAPIKeys: apiKeys,
	}
}

// RequireAPIKey middleware that requires valid API key
func (m *APIKeyMiddleware) RequireAPIKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Success: false,
				Error: &ErrorInfo{
					Type:    string(apperrors.ErrorTypeUnauthorized),
					Message: "API key required",
				},
				Timestamp: time.Now().UTC(),
			})
			c.Abort()
			return
		}

		serviceName, valid := m.validAPIKeys[apiKey]
		if !valid {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Success: false,
				Error: &ErrorInfo{
					Type:    string(apperrors.ErrorTypeUnauthorized),
					Message: "Invalid API key",
				},
				Timestamp: time.Now().UTC(),
			})
			c.Abort()
			return
		}

		// Store service name in context
		c.Set("service_name", serviceName)
		c.Set("auth_type", "api_key")

		c.Next()
	}
}

// Combined Authentication (JWT or API Key)

// FlexibleAuth allows either JWT or API key authentication
func FlexibleAuth(authMiddleware *AuthMiddleware, apiKeyMiddleware *APIKeyMiddleware) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try JWT first
		token, err := authMiddleware.extractTokenFromHeader(c)
		if err == nil {
			user, err := authMiddleware.authService.ValidateToken(c.Request.Context(), token)
			if err == nil {
				// JWT authentication successful
				c.Set(UserContextKey, user)
				c.Set(UserIDContextKey, user.ID())
				c.Set("auth_type", "jwt")
				c.Next()
				return
			}
		}

		// Try API Key
		apiKey := c.GetHeader("X-API-Key")
		if apiKey != "" {
			serviceName, valid := apiKeyMiddleware.validAPIKeys[apiKey]
			if valid {
				// API Key authentication successful
				c.Set("service_name", serviceName)
				c.Set("auth_type", "api_key")
				c.Next()
				return
			}
		}

		// No valid authentication found
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Success: false,
			Error: &ErrorInfo{
				Type:    string(apperrors.ErrorTypeUnauthorized),
				Message: "Authentication required",
			},
			Timestamp: time.Now().UTC(),
		})
		c.Abort()
	}
}
