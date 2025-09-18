package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/tranvuongduy2003/go-mvc/internal/application/dto"
	"github.com/tranvuongduy2003/go-mvc/internal/application/services"
	"github.com/tranvuongduy2003/go-mvc/internal/application/validators"
	"github.com/tranvuongduy2003/go-mvc/internal/shared/logger"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	userService   *services.UserApplicationService
	userValidator *validators.UserValidator
	logger        *logger.Logger
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(
	userService *services.UserApplicationService,
	userValidator *validators.UserValidator,
	logger *logger.Logger,
) *AuthHandler {
	return &AuthHandler{
		userService:   userService,
		userValidator: userValidator,
		logger:        logger,
	}
}

// Login handles POST /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	h.logger.Info("Processing login request")

	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorf("Failed to bind login request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Validate request
	if validationErrors := h.userValidator.ValidateLoginRequest(&req); len(validationErrors) > 0 {
		h.logger.Errorf("Login validation failed: %v", validationErrors)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": validationErrors,
		})
		return
	}

	// Authenticate user
	response, err := h.userService.Login(c.Request.Context(), &req)
	if err != nil {
		h.logger.Errorf("Login failed: %v", err)
		if strings.Contains(err.Error(), "invalid credentials") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid email or password",
			})
			return
		}
		if strings.Contains(err.Error(), "user not found") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid email or password",
			})
			return
		}
		if strings.Contains(err.Error(), "user not active") {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Account is not active",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Login failed",
		})
		return
	}

	h.logger.Infof("User logged in successfully: %s", response.User.Email)
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"data":    response,
	})
}

// Register handles POST /auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	h.logger.Info("Processing registration request")

	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorf("Failed to bind registration request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Validate request
	if validationErrors := h.userValidator.ValidateCreateUserRequest(&req); len(validationErrors) > 0 {
		h.logger.Errorf("Registration validation failed: %v", validationErrors)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": validationErrors,
		})
		return
	}

	// Create user
	user, err := h.userService.CreateUser(c.Request.Context(), &req)
	if err != nil {
		h.logger.Errorf("Failed to register user: %v", err)
		if strings.Contains(err.Error(), "already exists") {
			c.JSON(http.StatusConflict, gin.H{
				"error": "User with this email or username already exists",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Registration failed",
		})
		return
	}

	h.logger.Infof("User registered successfully: %s", user.Email)
	c.JSON(http.StatusCreated, gin.H{
		"message": "Registration successful",
		"data":    user,
	})
}

// RefreshToken handles POST /auth/refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	h.logger.Info("Processing token refresh request")

	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Errorf("Failed to bind refresh token request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	if req.RefreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Refresh token is required",
		})
		return
	}

	// For now, return not implemented
	h.logger.Info("Refresh token feature not implemented yet")
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "Refresh token feature not implemented yet",
	})
}

// Logout handles POST /auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	h.logger.Info("Processing logout request")

	// For now, just return success
	h.logger.Info("User logged out successfully")
	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
	})
}

// GetProfile handles GET /auth/profile
func (h *AuthHandler) GetProfile(c *gin.Context) {
	h.logger.Info("Getting user profile")

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Error("User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		h.logger.Error("Invalid user ID format in context")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}

	// Parse user ID
	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.logger.Errorf("Invalid user ID format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	// Get user by ID
	user, err := h.userService.GetUserByID(c.Request.Context(), userUUID)
	if err != nil {
		h.logger.Errorf("Failed to get profile: %v", err)
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get profile",
		})
		return
	}

	h.logger.Infof("Profile retrieved successfully for user: %s", userIDStr)
	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

// NotImplemented returns not implemented for features not yet available
func (h *AuthHandler) NotImplemented(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "Feature not implemented yet",
	})
}
