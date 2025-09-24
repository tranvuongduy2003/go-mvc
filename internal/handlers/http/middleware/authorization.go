package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tranvuongduy2003/go-mvc/internal/core/ports/services"
	apperrors "github.com/tranvuongduy2003/go-mvc/pkg/errors"
)

// AuthzMiddleware provides authorization middleware functionality
type AuthzMiddleware struct {
	authzService services.AuthorizationService
}

// NewAuthzMiddleware creates a new authorization middleware
func NewAuthzMiddleware(authzService services.AuthorizationService) *AuthzMiddleware {
	return &AuthzMiddleware{
		authzService: authzService,
	}
}

// RequirePermission middleware that requires specific permission
func (m *AuthzMiddleware) RequirePermission(resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := RequireUserID(c)
		if err != nil {
			m.sendAuthzErrorResponse(c, "User authentication required for authorization")
			return
		}

		hasPermission, err := m.authzService.UserHasPermission(c.Request.Context(), userID, resource, action)
		if err != nil {
			m.sendInternalErrorResponse(c, "Failed to check permission")
			return
		}

		if !hasPermission {
			m.sendAuthzErrorResponse(c, "Insufficient permissions to access this resource")
			return
		}

		c.Next()
	}
}

// RequirePermissionByName middleware that requires specific permission by name
func (m *AuthzMiddleware) RequirePermissionByName(permissionName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := RequireUserID(c)
		if err != nil {
			m.sendAuthzErrorResponse(c, "User authentication required for authorization")
			return
		}

		hasPermission, err := m.authzService.UserHasPermissionByName(c.Request.Context(), userID, permissionName)
		if err != nil {
			m.sendInternalErrorResponse(c, "Failed to check permission")
			return
		}

		if !hasPermission {
			m.sendAuthzErrorResponse(c, "Insufficient permissions to access this resource")
			return
		}

		c.Next()
	}
}

// RequireRole middleware that requires specific role
func (m *AuthzMiddleware) RequireRole(roleName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := RequireUserID(c)
		if err != nil {
			m.sendAuthzErrorResponse(c, "User authentication required for authorization")
			return
		}

		hasRole, err := m.authzService.UserHasRole(c.Request.Context(), userID, roleName)
		if err != nil {
			m.sendInternalErrorResponse(c, "Failed to check role")
			return
		}

		if !hasRole {
			m.sendAuthzErrorResponse(c, "Insufficient role to access this resource")
			return
		}

		c.Next()
	}
}

// RequireAnyRole middleware that requires any of the specified roles
func (m *AuthzMiddleware) RequireAnyRole(roleNames ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := RequireUserID(c)
		if err != nil {
			m.sendAuthzErrorResponse(c, "User authentication required for authorization")
			return
		}

		hasRole, err := m.authzService.UserHasAnyRole(c.Request.Context(), userID, roleNames)
		if err != nil {
			m.sendInternalErrorResponse(c, "Failed to check roles")
			return
		}

		if !hasRole {
			m.sendAuthzErrorResponse(c, "Insufficient roles to access this resource")
			return
		}

		c.Next()
	}
}

// RequireAllRoles middleware that requires all specified roles
func (m *AuthzMiddleware) RequireAllRoles(roleNames ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := RequireUserID(c)
		if err != nil {
			m.sendAuthzErrorResponse(c, "User authentication required for authorization")
			return
		}

		hasAllRoles, err := m.authzService.UserHasAllRoles(c.Request.Context(), userID, roleNames)
		if err != nil {
			m.sendInternalErrorResponse(c, "Failed to check roles")
			return
		}

		if !hasAllRoles {
			m.sendAuthzErrorResponse(c, "All required roles are needed to access this resource")
			return
		}

		c.Next()
	}
}

// RequireAdmin middleware that requires admin role
func (m *AuthzMiddleware) RequireAdmin() gin.HandlerFunc {
	return m.RequireRole("admin")
}

// RequireModerator middleware that requires admin or moderator role
func (m *AuthzMiddleware) RequireModerator() gin.HandlerFunc {
	return m.RequireAnyRole("admin", "moderator")
}

// RequireOwnership middleware that checks if user owns a resource
// This middleware extracts resource ID from URL parameters and checks ownership
func (m *AuthzMiddleware) RequireOwnership(paramName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := RequireUserID(c)
		if err != nil {
			m.sendAuthzErrorResponse(c, "User authentication required for authorization")
			return
		}

		resourceOwnerID := c.Param(paramName)
		if resourceOwnerID == "" {
			m.sendAuthzErrorResponse(c, "Resource owner ID not found in request")
			return
		}

		// Check if user is the owner of the resource
		if userID != resourceOwnerID {
			// Allow admin to access any resource
			isAdmin, err := m.authzService.IsAdmin(c.Request.Context(), userID)
			if err != nil {
				m.sendInternalErrorResponse(c, "Failed to check admin status")
				return
			}

			if !isAdmin {
				m.sendAuthzErrorResponse(c, "You can only access your own resources")
				return
			}
		}

		c.Next()
	}
}

// RequireOwnershipOrRole middleware that checks ownership or specific role
func (m *AuthzMiddleware) RequireOwnershipOrRole(paramName, roleName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := RequireUserID(c)
		if err != nil {
			m.sendAuthzErrorResponse(c, "User authentication required for authorization")
			return
		}

		resourceOwnerID := c.Param(paramName)
		if resourceOwnerID == "" {
			m.sendAuthzErrorResponse(c, "Resource owner ID not found in request")
			return
		}

		// Check if user is the owner of the resource
		if userID == resourceOwnerID {
			c.Next()
			return
		}

		// Check if user has the required role
		hasRole, err := m.authzService.UserHasRole(c.Request.Context(), userID, roleName)
		if err != nil {
			m.sendInternalErrorResponse(c, "Failed to check role")
			return
		}

		if !hasRole {
			m.sendAuthzErrorResponse(c, "You can only access your own resources or need appropriate role")
			return
		}

		c.Next()
	}
}

// DynamicPermissionCheck middleware that extracts resource and action from request
func (m *AuthzMiddleware) DynamicPermissionCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := RequireUserID(c)
		if err != nil {
			m.sendAuthzErrorResponse(c, "User authentication required for authorization")
			return
		}

		// Extract resource from URL path
		resource := m.extractResourceFromPath(c.Request.URL.Path)
		if resource == "" {
			m.sendAuthzErrorResponse(c, "Unable to determine resource from request")
			return
		}

		// Map HTTP method to action
		action := m.mapMethodToAction(c.Request.Method)
		if action == "" {
			m.sendAuthzErrorResponse(c, "Unable to determine action from request method")
			return
		}

		// Check permission
		hasPermission, err := m.authzService.UserHasPermission(c.Request.Context(), userID, resource, action)
		if err != nil {
			m.sendInternalErrorResponse(c, "Failed to check permission")
			return
		}

		if !hasPermission {
			m.sendAuthzErrorResponse(c, "Insufficient permissions to perform this action on this resource")
			return
		}

		c.Next()
	}
}

// ConditionalAccess middleware that applies different rules based on conditions
func (m *AuthzMiddleware) ConditionalAccess(conditions map[string]gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check conditions and apply appropriate middleware
		for condition, handler := range conditions {
			if m.evaluateCondition(c, condition) {
				handler(c)
				return
			}
		}

		// Default: deny access
		m.sendAuthzErrorResponse(c, "Access denied: no matching access conditions")
	}
}

// Helper methods

// extractResourceFromPath extracts resource name from URL path
func (m *AuthzMiddleware) extractResourceFromPath(path string) string {
	// Remove leading/trailing slashes and split by "/"
	path = strings.Trim(path, "/")
	parts := strings.Split(path, "/")

	// Look for resource patterns like "/api/v1/users" -> "users"
	if len(parts) >= 3 && parts[0] == "api" {
		return parts[2] // Return the resource part
	}

	// Fallback: return the first part
	if len(parts) > 0 {
		return parts[0]
	}

	return ""
}

// mapMethodToAction maps HTTP methods to RBAC actions
func (m *AuthzMiddleware) mapMethodToAction(method string) string {
	switch strings.ToUpper(method) {
	case "GET":
		return "read"
	case "POST":
		return "create"
	case "PUT", "PATCH":
		return "update"
	case "DELETE":
		return "delete"
	default:
		return ""
	}
}

// evaluateCondition evaluates access conditions
func (m *AuthzMiddleware) evaluateCondition(c *gin.Context, condition string) bool {
	switch condition {
	case "authenticated":
		_, exists := GetUserIDFromContext(c)
		return exists
	case "admin":
		userID, exists := GetUserIDFromContext(c)
		if !exists {
			return false
		}
		isAdmin, err := m.authzService.IsAdmin(c.Request.Context(), userID)
		return err == nil && isAdmin
	case "moderator":
		userID, exists := GetUserIDFromContext(c)
		if !exists {
			return false
		}
		isModerator, err := m.authzService.IsModerator(c.Request.Context(), userID)
		return err == nil && isModerator
	default:
		return false
	}
}

// Error response methods

func (m *AuthzMiddleware) sendAuthzErrorResponse(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, ErrorResponse{
		Success: false,
		Error: &ErrorInfo{
			Type:    string(apperrors.ErrorTypeForbidden),
			Message: message,
		},
		Timestamp: time.Now().UTC(),
	})
	c.Abort()
}

func (m *AuthzMiddleware) sendInternalErrorResponse(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Success: false,
		Error: &ErrorInfo{
			Type:    string(apperrors.ErrorTypeInternal),
			Message: message,
		},
		Timestamp: time.Now().UTC(),
	})
	c.Abort()
}

// Utility functions for creating combined middleware

// RequireAuthAndPermission combines authentication and permission checking
func RequireAuthAndPermission(authMiddleware *AuthMiddleware, authzMiddleware *AuthzMiddleware, resource, action string) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// First authenticate
		authMiddleware.RequireAuth()(c)
		if c.IsAborted() {
			return
		}

		// Then authorize
		authzMiddleware.RequirePermission(resource, action)(c)
	})
}

// RequireAuthAndRole combines authentication and role checking
func RequireAuthAndRole(authMiddleware *AuthMiddleware, authzMiddleware *AuthzMiddleware, roleName string) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// First authenticate
		authMiddleware.RequireAuth()(c)
		if c.IsAborted() {
			return
		}

		// Then check role
		authzMiddleware.RequireRole(roleName)(c)
	})
}

// RequireAuthAndOwnership combines authentication and ownership checking
func RequireAuthAndOwnership(authMiddleware *AuthMiddleware, authzMiddleware *AuthzMiddleware, paramName string) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// First authenticate
		authMiddleware.RequireAuth()(c)
		if c.IsAborted() {
			return
		}

		// Then check ownership
		authzMiddleware.RequireOwnership(paramName)(c)
	})
}
