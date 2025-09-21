package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/tranvuongduy2003/go-mvc/internal/application/dto"
	"github.com/tranvuongduy2003/go-mvc/internal/application/services"
	"github.com/tranvuongduy2003/go-mvc/internal/application/validators"
	"github.com/tranvuongduy2003/go-mvc/internal/handlers/http/middleware"
	"github.com/tranvuongduy2003/go-mvc/internal/shared/logger"
	"github.com/tranvuongduy2003/go-mvc/internal/shared/tracing"
	"go.opentelemetry.io/otel/attribute"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userService   *services.UserApplicationService
	userValidator *validators.UserValidator
	logger        *logger.Logger
	tracing       *tracing.TracingService
}

// NewUserHandler creates a new user handler
func NewUserHandler(
	userService *services.UserApplicationService,
	userValidator *validators.UserValidator,
	logger *logger.Logger,
	tracing *tracing.TracingService,
) *UserHandler {
	return &UserHandler{
		userService:   userService,
		userValidator: userValidator,
		logger:        logger,
		tracing:       tracing,
	}
}

// CreateUser handles POST /users
func (h *UserHandler) CreateUser(c *gin.Context) {
	ctx := middleware.TraceContext(c)
	ctx, span := h.tracing.StartHTTPSpan(ctx, c.Request.Method, c.FullPath())
	defer span.End()

	h.logger.Infof("Creating user")

	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorf("Failed to bind request: %v", err)
		h.tracing.RecordError(span, err)
		span.SetAttributes(attribute.String("error.type", "validation_error"))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Add request attributes to span
	span.SetAttributes(
		attribute.String("user.email", req.Email),
		attribute.String("user.username", req.Username),
	)

	// Validate request
	if validationErrors := h.userValidator.ValidateCreateUserRequest(&req); len(validationErrors) > 0 {
		h.logger.Errorf("Validation failed: %v", validationErrors)
		span.SetAttributes(attribute.String("error.type", "validation_failed"))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": validationErrors,
		})
		return
	}

	// Create user
	user, err := h.userService.CreateUser(ctx, &req)
	if err != nil {
		h.logger.Errorf("Failed to create user: %v", err)
		h.tracing.RecordError(span, err)
		if strings.Contains(err.Error(), "already exists") {
			span.SetAttributes(attribute.String("error.type", "user_exists"))
			c.JSON(http.StatusConflict, gin.H{
				"error": "User already exists",
			})
			return
		}
		span.SetAttributes(attribute.String("error.type", "internal_error"))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	h.logger.Infof("User created successfully: %s", user.ID)
	span.SetAttributes(
		attribute.String("user.id", user.ID.String()),
		attribute.String("response.status", "success"),
	)
	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"data":    user,
	})
}

// GetUser handles GET /users/:id
func (h *UserHandler) GetUser(c *gin.Context) {
	userIDStr := c.Param("id")
	h.logger.Infof("Getting user by ID: %s", userIDStr)

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.logger.Errorf("Invalid user ID format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		h.logger.Errorf("Failed to get user: %v", err)
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user",
		})
		return
	}

	h.logger.Infof("User retrieved successfully: %s", user.ID)
	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

// UpdateUser handles PUT /users/:id
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userIDStr := c.Param("id")
	h.logger.Infof("Updating user: %s", userIDStr)

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.logger.Errorf("Invalid user ID format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorf("Failed to bind request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Validate request
	if validationErrors := h.userValidator.ValidateUpdateUserRequest(&req); len(validationErrors) > 0 {
		h.logger.Errorf("Validation failed: %v", validationErrors)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": validationErrors,
		})
		return
	}

	// Update user
	user, err := h.userService.UpdateUser(c.Request.Context(), userID, &req)
	if err != nil {
		h.logger.Errorf("Failed to update user: %v", err)
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update user",
		})
		return
	}

	h.logger.Infof("User updated successfully: %s", user.ID)
	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"data":    user,
	})
}

// DeleteUser handles DELETE /users/:id
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userIDStr := c.Param("id")
	h.logger.Infof("Deleting user: %s", userIDStr)

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.logger.Errorf("Invalid user ID format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	// Delete user
	if err := h.userService.DeleteUser(c.Request.Context(), userID); err != nil {
		h.logger.Errorf("Failed to delete user: %v", err)
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete user",
		})
		return
	}

	h.logger.Infof("User deleted successfully: %s", userID)
	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}

// ListUsers handles GET /users
func (h *UserHandler) ListUsers(c *gin.Context) {
	h.logger.Info("Listing users")

	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	sort := c.Query("sort")
	order := c.Query("order")
	search := c.Query("search")
	role := c.Query("role")

	var isActive *bool
	if isActiveStr := c.Query("is_active"); isActiveStr != "" {
		if active, err := strconv.ParseBool(isActiveStr); err == nil {
			isActive = &active
		}
	}

	var isDeleted *bool
	if isDeletedStr := c.Query("is_deleted"); isDeletedStr != "" {
		if deleted, err := strconv.ParseBool(isDeletedStr); err == nil {
			isDeleted = &deleted
		}
	}

	req := &dto.UserListRequest{
		Page:      page,
		Limit:     limit,
		Sort:      sort,
		Order:     order,
		Search:    search,
		IsActive:  isActive,
		IsDeleted: isDeleted,
	}

	if role != "" {
		req.Role = &role
	}

	// Validate request
	if validationErrors := h.userValidator.ValidateUserListRequest(req); len(validationErrors) > 0 {
		h.logger.Errorf("Validation failed: %v", validationErrors)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": validationErrors,
		})
		return
	}

	// List users
	response, err := h.userService.ListUsers(c.Request.Context(), req)
	if err != nil {
		h.logger.Errorf("Failed to list users: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list users",
		})
		return
	}

	h.logger.Infof("Listed %d users", len(response.Users))
	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

// ChangePassword handles PUT /users/:id/password
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userIDStr := c.Param("id")
	h.logger.Infof("Changing password for user: %s", userIDStr)

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.logger.Errorf("Invalid user ID format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorf("Failed to bind request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Validate request
	if validationErrors := h.userValidator.ValidateChangePasswordRequest(&req); len(validationErrors) > 0 {
		h.logger.Errorf("Validation failed: %v", validationErrors)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": validationErrors,
		})
		return
	}

	// Change password
	if err := h.userService.ChangePassword(c.Request.Context(), userID, &req); err != nil {
		h.logger.Errorf("Failed to change password: %v", err)
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}
		if strings.Contains(err.Error(), "incorrect password") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Current password is incorrect",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to change password",
		})
		return
	}

	h.logger.Infof("Password changed successfully for user: %s", userID)
	c.JSON(http.StatusOK, gin.H{
		"message": "Password changed successfully",
	})
}

// ActivateUser handles PUT /users/:id/activate
func (h *UserHandler) ActivateUser(c *gin.Context) {
	userIDStr := c.Param("id")
	h.logger.Infof("Activating user: %s", userIDStr)

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.logger.Errorf("Invalid user ID format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	// Activate user
	if err := h.userService.ActivateUser(c.Request.Context(), userID); err != nil {
		h.logger.Errorf("Failed to activate user: %v", err)
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to activate user",
		})
		return
	}

	h.logger.Infof("User activated successfully: %s", userID)
	c.JSON(http.StatusOK, gin.H{
		"message": "User activated successfully",
	})
}

// DeactivateUser handles PUT /users/:id/deactivate
func (h *UserHandler) DeactivateUser(c *gin.Context) {
	userIDStr := c.Param("id")
	h.logger.Infof("Deactivating user: %s", userIDStr)

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.logger.Errorf("Invalid user ID format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	// Deactivate user
	if err := h.userService.DeactivateUser(c.Request.Context(), userID); err != nil {
		h.logger.Errorf("Failed to deactivate user: %v", err)
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to deactivate user",
		})
		return
	}

	h.logger.Infof("User deactivated successfully: %s", userID)
	c.JSON(http.StatusOK, gin.H{
		"message": "User deactivated successfully",
	})
}

// ChangeRole handles PUT /users/:id/role
func (h *UserHandler) ChangeRole(c *gin.Context) {
	userIDStr := c.Param("id")
	h.logger.Infof("Changing role for user: %s", userIDStr)

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.logger.Errorf("Invalid user ID format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	var req dto.ChangeRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorf("Failed to bind request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Validate request
	if validationErrors := h.userValidator.ValidateChangeRoleRequest(&req); len(validationErrors) > 0 {
		h.logger.Errorf("Validation failed: %v", validationErrors)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": validationErrors,
		})
		return
	}

	// Change role
	if err := h.userService.ChangeRole(c.Request.Context(), userID, &req); err != nil {
		h.logger.Errorf("Failed to change role: %v", err)
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to change role",
		})
		return
	}

	h.logger.Infof("Role changed successfully for user: %s", userID)
	c.JSON(http.StatusOK, gin.H{
		"message": "Role changed successfully",
	})
}

// CheckAvailability handles POST /users/check-availability
func (h *UserHandler) CheckAvailability(c *gin.Context) {
	h.logger.Info("Checking availability")

	var req dto.CheckAvailabilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorf("Failed to bind request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Check availability
	response, err := h.userService.CheckAvailability(c.Request.Context(), &req)
	if err != nil {
		h.logger.Errorf("Failed to check availability: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to check availability",
		})
		return
	}

	h.logger.Info("Availability checked successfully")
	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

// SearchUsers handles GET /users/search
func (h *UserHandler) SearchUsers(c *gin.Context) {
	h.logger.Info("Searching users")

	query := c.Query("q")
	searchType := c.DefaultQuery("type", "all")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Search query is required",
		})
		return
	}

	// Search users
	response, err := h.userService.SearchUsers(c.Request.Context(), query, searchType, page, limit)
	if err != nil {
		h.logger.Errorf("Failed to search users: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to search users",
		})
		return
	}

	h.logger.Infof("Found %d users for search query", len(response.Users))
	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}
