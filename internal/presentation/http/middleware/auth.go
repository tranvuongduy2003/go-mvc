package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/contracts"
	apperrors "github.com/tranvuongduy2003/go-mvc/pkg/errors"
)

const (
	UserContextKey = "user"

	UserIDContextKey = "user_id"

	AuthorizationHeader = "Authorization"

	BearerPrefix = "Bearer "
)

type ErrorResponse struct {
	Success   bool       `json:"success"`
	Error     *ErrorInfo `json:"error,omitempty"`
	Timestamp time.Time  `json:"timestamp"`
}

type ErrorInfo struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

type AuthMiddleware struct {
	authService contracts.AuthService
}

func NewAuthMiddleware(authService contracts.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

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

		c.Set(UserContextKey, user)
		c.Set(UserIDContextKey, user.ID())

		c.Next()
	}
}

func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := m.extractTokenFromHeader(c)
		if err != nil {
			c.Next()
			return
		}

		user, err := m.authService.ValidateToken(c.Request.Context(), token)
		if err != nil {
			c.Next()
			return
		}

		c.Set(UserContextKey, user)
		c.Set(UserIDContextKey, user.ID())

		c.Next()
	}
}

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

		if !user.IsActive() {
			m.sendErrorResponse(c, http.StatusForbidden, "Account is inactive", "Your account has been deactivated. Please contact support.")
			return
		}

		c.Set(UserContextKey, user)
		c.Set(UserIDContextKey, user.ID())

		c.Next()
	}
}

func (m *AuthMiddleware) TokenRefresh() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := m.extractTokenFromHeader(c)
		if err != nil {
			m.sendErrorResponse(c, http.StatusUnauthorized, "Refresh token required", err.Error())
			return
		}

		c.Set("refresh_token", token)

		c.Next()
	}
}

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

type AuthError struct {
	Message string
}

func NewAuthError(message string) *AuthError {
	return &AuthError{Message: message}
}

func (e *AuthError) Error() string {
	return e.Message
}

func GetUserFromContext(c *gin.Context) (interface{}, bool) {
	user, exists := c.Get(UserContextKey)
	return user, exists
}

func GetUserIDFromContext(c *gin.Context) (string, bool) {
	userID, exists := c.Get(UserIDContextKey)
	if !exists {
		return "", false
	}

	id, ok := userID.(string)
	return id, ok
}

func RequireUserID(c *gin.Context) (string, error) {
	userID, exists := GetUserIDFromContext(c)
	if !exists {
		return "", NewAuthError("user ID not found in context")
	}
	return userID, nil
}

type APIKeyMiddleware struct {
	validAPIKeys map[string]string // API key -> service name
}

func NewAPIKeyMiddleware(apiKeys map[string]string) *APIKeyMiddleware {
	return &APIKeyMiddleware{
		validAPIKeys: apiKeys,
	}
}

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

		c.Set("service_name", serviceName)
		c.Set("auth_type", "api_key")

		c.Next()
	}
}

func FlexibleAuth(authMiddleware *AuthMiddleware, apiKeyMiddleware *APIKeyMiddleware) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := authMiddleware.extractTokenFromHeader(c)
		if err == nil {
			user, err := authMiddleware.authService.ValidateToken(c.Request.Context(), token)
			if err == nil {
				c.Set(UserContextKey, user)
				c.Set(UserIDContextKey, user.ID())
				c.Set("auth_type", "jwt")
				c.Next()
				return
			}
		}

		apiKey := c.GetHeader("X-API-Key")
		if apiKey != "" {
			serviceName, valid := apiKeyMiddleware.validAPIKeys[apiKey]
			if valid {
				c.Set("service_name", serviceName)
				c.Set("auth_type", "api_key")
				c.Next()
				return
			}
		}

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
