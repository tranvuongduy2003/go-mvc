package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	authCommands "github.com/tranvuongduy2003/go-mvc/internal/application/commands/auth"
	dto "github.com/tranvuongduy2003/go-mvc/internal/application/dto/auth"
	authQueries "github.com/tranvuongduy2003/go-mvc/internal/application/queries/auth"
	"github.com/tranvuongduy2003/go-mvc/pkg/response"
)

// AuthHandler handles HTTP requests for authentication operations
type AuthHandler struct {
	loginHandler                *authCommands.LoginCommandHandler
	registerHandler             *authCommands.RegisterCommandHandler
	refreshTokenHandler         *authCommands.RefreshTokenCommandHandler
	changePasswordHandler       *authCommands.ChangePasswordCommandHandler
	resetPasswordHandler        *authCommands.ResetPasswordCommandHandler
	confirmPasswordResetHandler *authCommands.ConfirmPasswordResetCommandHandler
	verifyEmailHandler          *authCommands.VerifyEmailCommandHandler
	resendVerificationHandler   *authCommands.ResendVerificationEmailCommandHandler
	logoutHandler               *authCommands.LogoutCommandHandler
	logoutAllDevicesHandler     *authCommands.LogoutAllDevicesCommandHandler
	getUserProfileHandler       *authQueries.GetUserProfileQueryHandler
	getUserPermissionsHandler   *authQueries.GetUserPermissionsQueryHandler
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(
	loginHandler *authCommands.LoginCommandHandler,
	registerHandler *authCommands.RegisterCommandHandler,
	refreshTokenHandler *authCommands.RefreshTokenCommandHandler,
	changePasswordHandler *authCommands.ChangePasswordCommandHandler,
	resetPasswordHandler *authCommands.ResetPasswordCommandHandler,
	confirmPasswordResetHandler *authCommands.ConfirmPasswordResetCommandHandler,
	verifyEmailHandler *authCommands.VerifyEmailCommandHandler,
	resendVerificationHandler *authCommands.ResendVerificationEmailCommandHandler,
	logoutHandler *authCommands.LogoutCommandHandler,
	logoutAllDevicesHandler *authCommands.LogoutAllDevicesCommandHandler,
	getUserProfileHandler *authQueries.GetUserProfileQueryHandler,
	getUserPermissionsHandler *authQueries.GetUserPermissionsQueryHandler,
) *AuthHandler {
	return &AuthHandler{
		loginHandler:                loginHandler,
		registerHandler:             registerHandler,
		refreshTokenHandler:         refreshTokenHandler,
		changePasswordHandler:       changePasswordHandler,
		resetPasswordHandler:        resetPasswordHandler,
		confirmPasswordResetHandler: confirmPasswordResetHandler,
		verifyEmailHandler:          verifyEmailHandler,
		resendVerificationHandler:   resendVerificationHandler,
		logoutHandler:               logoutHandler,
		logoutAllDevicesHandler:     logoutAllDevicesHandler,
		getUserProfileHandler:       getUserProfileHandler,
		getUserPermissionsHandler:   getUserPermissionsHandler,
	}
}

// Login authenticates a user
// @Summary Authenticate user
// @Description Login with email and password to get JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body dto.LoginRequest true "Login credentials"
// @Success 200 {object} response.APIResponse{data=dto.LoginResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	// Execute command
	result, err := h.loginHandler.Handle(c.Request.Context(), authCommands.LoginCommand{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "Login successful", result)
}

// Register creates a new user account
// @Summary Register new user
// @Description Create a new user account with email verification
// @Tags auth
// @Accept json
// @Produce json
// @Param user body dto.RegisterRequest true "Registration data"
// @Success 201 {object} response.APIResponse{data=dto.RegisterResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 409 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	// Execute command
	result, err := h.registerHandler.Handle(c.Request.Context(), authCommands.RegisterCommand{
		Email:    req.Email,
		Name:     req.Name,
		Phone:    req.Phone,
		Password: req.Password,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Registration successful",
		"data":    result,
	})
}

// RefreshToken refreshes JWT access token
// @Summary Refresh JWT token
// @Description Get new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param token body dto.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} response.APIResponse{data=dto.TokensDTO}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	// Execute command
	result, err := h.refreshTokenHandler.Handle(c.Request.Context(), authCommands.RefreshTokenCommand{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "Token refreshed successfully", result)
}

// ChangePassword changes user password
// @Summary Change user password
// @Description Change current user's password
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param password body dto.ChangePasswordRequest true "Password change data"
// @Success 200 {object} response.APIResponse{data=dto.StatusResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/auth/change-password [post]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, errors.New("user not authenticated"))
		return
	}

	// Execute command
	result, err := h.changePasswordHandler.Handle(c.Request.Context(), authCommands.ChangePasswordCommand{
		UserID:      userID.(string),
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "Password changed successfully", result)
}

// ResetPassword initiates password reset process
// @Summary Reset password
// @Description Send password reset instructions to email
// @Tags auth
// @Accept json
// @Produce json
// @Param email body dto.ResetPasswordRequest true "Reset password email"
// @Success 200 {object} response.APIResponse{data=dto.StatusResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	// Execute command
	result, err := h.resetPasswordHandler.Handle(c.Request.Context(), authCommands.ResetPasswordCommand{
		Email: req.Email,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "Reset instructions sent", result)
}

// ConfirmPasswordReset completes password reset with token
// @Summary Confirm password reset
// @Description Complete password reset with verification token
// @Tags auth
// @Accept json
// @Produce json
// @Param reset body dto.ConfirmPasswordResetRequest true "Password reset confirmation"
// @Success 200 {object} response.APIResponse{data=dto.StatusResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/auth/reset-password/confirm [post]
func (h *AuthHandler) ConfirmPasswordReset(c *gin.Context) {
	var req dto.ConfirmPasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	// Execute command
	result, err := h.confirmPasswordResetHandler.Handle(c.Request.Context(), authCommands.ConfirmPasswordResetCommand{
		Token:       req.Token,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "Password reset successfully", result)
}

// VerifyEmail verifies user email with token
// @Summary Verify email
// @Description Verify user email address with verification token
// @Tags auth
// @Accept json
// @Produce json
// @Param verification body dto.VerifyEmailRequest true "Email verification token"
// @Success 200 {object} response.APIResponse{data=dto.StatusResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/auth/verify-email [post]
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var req dto.VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	// Execute command
	result, err := h.verifyEmailHandler.Handle(c.Request.Context(), authCommands.VerifyEmailCommand{
		Token: req.Token,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "Email verified successfully", result)
}

// ResendVerificationEmail resends email verification
// @Summary Resend verification email
// @Description Send new email verification link
// @Tags auth
// @Accept json
// @Produce json
// @Param email body dto.ResendVerificationRequest true "Email for verification"
// @Success 200 {object} response.APIResponse{data=dto.StatusResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 429 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/auth/resend-verification [post]
func (h *AuthHandler) ResendVerificationEmail(c *gin.Context) {
	var req dto.ResendVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	// Execute command
	result, err := h.resendVerificationHandler.Handle(c.Request.Context(), authCommands.ResendVerificationEmailCommand{
		Email: req.Email,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "Verification email sent", result)
}

// Logout logs out current user
// @Summary Logout user
// @Description Logout current user by invalidating tokens
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse{data=dto.StatusResponse}
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, errors.New("user not authenticated"))
		return
	}

	// Execute command
	result, err := h.logoutHandler.Handle(c.Request.Context(), authCommands.LogoutCommand{
		UserID: userID.(string),
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "Logged out successfully", result)
}

// LogoutAllDevices logs out user from all devices
// @Summary Logout from all devices
// @Description Logout current user from all devices
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse{data=dto.StatusResponse}
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/auth/logout-all [post]
func (h *AuthHandler) LogoutAllDevices(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, errors.New("user not authenticated"))
		return
	}

	// Execute command
	result, err := h.logoutAllDevicesHandler.Handle(c.Request.Context(), authCommands.LogoutAllDevicesCommand{
		UserID: userID.(string),
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "Logged out from all devices", result)
}

// GetProfile gets current user profile
// @Summary Get user profile
// @Description Get current authenticated user's profile with roles and permissions
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse{data=dto.UserProfileResponse}
// @Failure 401 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, errors.New("user not authenticated"))
		return
	}

	// Execute query
	result, err := h.getUserProfileHandler.Handle(c.Request.Context(), authQueries.GetUserProfileQuery{
		UserID: userID.(string),
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "Profile retrieved successfully", result)
}

// GetPermissions gets current user permissions
// @Summary Get user permissions
// @Description Get current authenticated user's effective permissions
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.APIResponse{data=[]dto.PermissionInfoDTO}
// @Failure 401 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /api/v1/auth/permissions [get]
func (h *AuthHandler) GetPermissions(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, errors.New("user not authenticated"))
		return
	}

	// Execute query
	result, err := h.getUserPermissionsHandler.Handle(c.Request.Context(), authQueries.GetUserPermissionsQuery{
		UserID: userID.(string),
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "Permissions retrieved successfully", result)
}
