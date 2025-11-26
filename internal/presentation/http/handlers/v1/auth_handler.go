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

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

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

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

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

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	result, err := h.refreshTokenHandler.Handle(c.Request.Context(), authCommands.RefreshTokenCommand{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "Token refreshed successfully", result)
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, errors.New("user not authenticated"))
		return
	}

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

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	result, err := h.resetPasswordHandler.Handle(c.Request.Context(), authCommands.ResetPasswordCommand{
		Email: req.Email,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "Reset instructions sent", result)
}

func (h *AuthHandler) ConfirmPasswordReset(c *gin.Context) {
	var req dto.ConfirmPasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

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

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var req dto.VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	result, err := h.verifyEmailHandler.Handle(c.Request.Context(), authCommands.VerifyEmailCommand{
		Token: req.Token,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "Email verified successfully", result)
}

func (h *AuthHandler) ResendVerificationEmail(c *gin.Context) {
	var req dto.ResendVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}

	result, err := h.resendVerificationHandler.Handle(c.Request.Context(), authCommands.ResendVerificationEmailCommand{
		Email: req.Email,
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "Verification email sent", result)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, errors.New("user not authenticated"))
		return
	}

	result, err := h.logoutHandler.Handle(c.Request.Context(), authCommands.LogoutCommand{
		UserID: userID.(string),
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "Logged out successfully", result)
}

func (h *AuthHandler) LogoutAllDevices(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, errors.New("user not authenticated"))
		return
	}

	result, err := h.logoutAllDevicesHandler.Handle(c.Request.Context(), authCommands.LogoutAllDevicesCommand{
		UserID: userID.(string),
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "Logged out from all devices", result)
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, errors.New("user not authenticated"))
		return
	}

	result, err := h.getUserProfileHandler.Handle(c.Request.Context(), authQueries.GetUserProfileQuery{
		UserID: userID.(string),
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "Profile retrieved successfully", result)
}

func (h *AuthHandler) GetPermissions(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, errors.New("user not authenticated"))
		return
	}

	result, err := h.getUserPermissionsHandler.Handle(c.Request.Context(), authQueries.GetUserPermissionsQuery{
		UserID: userID.(string),
	})
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "Permissions retrieved successfully", result)
}
