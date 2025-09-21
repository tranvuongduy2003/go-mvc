package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/tranvuongduy2003/go-mvc/internal/core/domain/rbac"
	"github.com/tranvuongduy2003/go-mvc/pkg/jwt"
)

// AuthMiddleware handles JWT authentication
type AuthMiddleware struct {
	jwtService  jwt.JWTService
	rbacService rbac.RBACService
	logger      *zap.Logger
}

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware(jwtService jwt.JWTService, rbacService rbac.RBACService, logger *zap.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService:  jwtService,
		rbacService: rbacService,
		logger:      logger,
	}
}

// AuthRequired middleware that requires valid JWT token
func (a *AuthMiddleware) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := a.extractToken(c)
		if token == "" {
			a.logger.Warn("Missing authorization token",
				zap.String("request_id", requestid.Get(c)),
				zap.String("path", c.Request.URL.Path),
			)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Missing or invalid authorization token",
				"code":    "MISSING_TOKEN",
			})
			c.Abort()
			return
		}

		claims, err := a.jwtService.ValidateToken(token)
		if err != nil {
			a.logger.Warn("Invalid token",
				zap.String("request_id", requestid.Get(c)),
				zap.String("path", c.Request.URL.Path),
				zap.Error(err),
			)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Invalid or expired token",
				"code":    "INVALID_TOKEN",
			})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("token_type", claims.Type)

		// Log successful authentication
		a.logger.Info("User authenticated",
			zap.String("request_id", requestid.Get(c)),
			zap.String("user_id", claims.UserID.String()),
			zap.String("email", claims.Email),
			zap.String("path", c.Request.URL.Path),
		)

		c.Next()
	}
}

// OptionalAuth middleware that allows both authenticated and unauthenticated requests
func (a *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := a.extractToken(c)
		if token != "" {
			claims, err := a.jwtService.ValidateToken(token)
			if err == nil {
				// Set user information if token is valid
				c.Set("user_id", claims.UserID)
				c.Set("user_email", claims.Email)
				c.Set("token_type", claims.Type)
				c.Set("authenticated", true)
			}
		}

		c.Next()
	}
}

// RequireRole middleware that requires specific role
func (a *AuthMiddleware) RequireRole(roleName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "User not authenticated",
				"code":    "NOT_AUTHENTICATED",
			})
			c.Abort()
			return
		}

		userID, ok := userIDInterface.(uuid.UUID)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Internal server error",
				"message": "Invalid user ID format",
				"code":    "INVALID_USER_ID",
			})
			c.Abort()
			return
		}

		hasRole, err := a.rbacService.HasAnyRole(context.Background(), userID, []string{roleName})
		if err != nil {
			a.logger.Error("Failed to check user role",
				zap.String("request_id", requestid.Get(c)),
				zap.String("user_id", userID.String()),
				zap.String("role", roleName),
				zap.Error(err),
			)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Internal server error",
				"message": "Failed to verify user permissions",
				"code":    "PERMISSION_CHECK_FAILED",
			})
			c.Abort()
			return
		}

		if !hasRole {
			a.logger.Warn("User lacks required role",
				zap.String("request_id", requestid.Get(c)),
				zap.String("user_id", userID.String()),
				zap.String("required_role", roleName),
				zap.String("path", c.Request.URL.Path),
			)
			c.JSON(http.StatusForbidden, gin.H{
				"error":         "Forbidden",
				"message":       "Insufficient permissions",
				"code":          "INSUFFICIENT_PERMISSIONS",
				"required_role": roleName,
			})
			c.Abort()
			return
		}

		c.Set("user_role", roleName)
		c.Next()
	}
}

// RequireAnyRole middleware that requires any of the specified roles
func (a *AuthMiddleware) RequireAnyRole(roleNames ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "User not authenticated",
				"code":    "NOT_AUTHENTICATED",
			})
			c.Abort()
			return
		}

		userID, ok := userIDInterface.(uuid.UUID)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Internal server error",
				"message": "Invalid user ID format",
				"code":    "INVALID_USER_ID",
			})
			c.Abort()
			return
		}

		hasRole, err := a.rbacService.HasAnyRole(context.Background(), userID, roleNames)
		if err != nil {
			a.logger.Error("Failed to check user roles",
				zap.String("request_id", requestid.Get(c)),
				zap.String("user_id", userID.String()),
				zap.Strings("roles", roleNames),
				zap.Error(err),
			)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Internal server error",
				"message": "Failed to verify user permissions",
				"code":    "PERMISSION_CHECK_FAILED",
			})
			c.Abort()
			return
		}

		if !hasRole {
			a.logger.Warn("User lacks required roles",
				zap.String("request_id", requestid.Get(c)),
				zap.String("user_id", userID.String()),
				zap.Strings("required_roles", roleNames),
				zap.String("path", c.Request.URL.Path),
			)
			c.JSON(http.StatusForbidden, gin.H{
				"error":          "Forbidden",
				"message":        "Insufficient permissions",
				"code":           "INSUFFICIENT_PERMISSIONS",
				"required_roles": roleNames,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequirePermission middleware that requires specific permission
func (a *AuthMiddleware) RequirePermission(resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "User not authenticated",
				"code":    "NOT_AUTHENTICATED",
			})
			c.Abort()
			return
		}

		userID, ok := userIDInterface.(uuid.UUID)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Internal server error",
				"message": "Invalid user ID format",
				"code":    "INVALID_USER_ID",
			})
			c.Abort()
			return
		}

		hasPermission, err := a.rbacService.HasPermission(context.Background(), userID, resource, action)
		if err != nil {
			a.logger.Error("Failed to check user permission",
				zap.String("request_id", requestid.Get(c)),
				zap.String("user_id", userID.String()),
				zap.String("resource", resource),
				zap.String("action", action),
				zap.Error(err),
			)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Internal server error",
				"message": "Failed to verify user permissions",
				"code":    "PERMISSION_CHECK_FAILED",
			})
			c.Abort()
			return
		}

		if !hasPermission {
			a.logger.Warn("User lacks required permission",
				zap.String("request_id", requestid.Get(c)),
				zap.String("user_id", userID.String()),
				zap.String("resource", resource),
				zap.String("action", action),
				zap.String("path", c.Request.URL.Path),
			)
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Forbidden",
				"message": "Insufficient permissions",
				"code":    "INSUFFICIENT_PERMISSIONS",
				"required_permission": map[string]string{
					"resource": resource,
					"action":   action,
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// AdminOnly middleware that requires admin role
func (a *AuthMiddleware) AdminOnly() gin.HandlerFunc {
	return a.RequireRole(rbac.RoleAdmin)
}

// ModeratorOrAdmin middleware that requires moderator or admin role
func (a *AuthMiddleware) ModeratorOrAdmin() gin.HandlerFunc {
	return a.RequireAnyRole(rbac.RoleModerator, rbac.RoleAdmin)
}

// OwnerOrAdmin middleware that allows resource owner or admin
func (a *AuthMiddleware) OwnerOrAdmin(resourceIDParam string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "User not authenticated",
				"code":    "NOT_AUTHENTICATED",
			})
			c.Abort()
			return
		}

		userID, ok := userIDInterface.(uuid.UUID)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Internal server error",
				"message": "Invalid user ID format",
				"code":    "INVALID_USER_ID",
			})
			c.Abort()
			return
		}

		// Check if user is admin
		isAdmin, err := a.rbacService.HasAnyRole(context.Background(), userID, []string{rbac.RoleAdmin})
		if err == nil && isAdmin {
			c.Next()
			return
		}

		// Check if user is the resource owner
		resourceID := c.Param(resourceIDParam)
		if resourceID != "" {
			resourceUUID, err := uuid.Parse(resourceID)
			if err == nil && resourceUUID == userID {
				c.Next()
				return
			}
		}

		a.logger.Warn("User is not owner or admin",
			zap.String("request_id", requestid.Get(c)),
			zap.String("user_id", userID.String()),
			zap.String("resource_id", resourceID),
			zap.String("path", c.Request.URL.Path),
		)

		c.JSON(http.StatusForbidden, gin.H{
			"error":   "Forbidden",
			"message": "You can only access your own resources",
			"code":    "OWNER_OR_ADMIN_REQUIRED",
		})
		c.Abort()
	}
}

// extractToken extracts JWT token from Authorization header
func (a *AuthMiddleware) extractToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	// Support both "Bearer token" and "token" formats
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}

	return authHeader
}

// GetCurrentUser gets current authenticated user information
func GetCurrentUser(c *gin.Context) (uuid.UUID, string, bool) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, "", false
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		return uuid.Nil, "", false
	}

	emailInterface, exists := c.Get("user_email")
	if !exists {
		return userID, "", true
	}

	email, ok := emailInterface.(string)
	if !ok {
		return userID, "", true
	}

	return userID, email, true
}

// IsAuthenticated checks if the current request is authenticated
func IsAuthenticated(c *gin.Context) bool {
	_, exists := c.Get("user_id")
	return exists
}

// HasRole checks if current user has specific role
func HasRole(c *gin.Context, rbacService rbac.RBACService, roleName string) bool {
	userID, _, authenticated := GetCurrentUser(c)
	if !authenticated {
		return false
	}

	hasRole, err := rbacService.HasAnyRole(context.Background(), userID, []string{roleName})
	return err == nil && hasRole
}

// HasPermission checks if current user has specific permission
func HasPermission(c *gin.Context, rbacService rbac.RBACService, resource, action string) bool {
	userID, _, authenticated := GetCurrentUser(c)
	if !authenticated {
		return false
	}

	hasPermission, err := rbacService.HasPermission(context.Background(), userID, resource, action)
	return err == nil && hasPermission
}
