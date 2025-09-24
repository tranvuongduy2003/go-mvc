package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tranvuongduy2003/go-mvc/internal/adapters/cache"
	"github.com/tranvuongduy2003/go-mvc/internal/core/domain/user"
	"github.com/tranvuongduy2003/go-mvc/internal/core/ports/repositories"
	"github.com/tranvuongduy2003/go-mvc/internal/core/ports/services"
	"github.com/tranvuongduy2003/go-mvc/internal/shared/security"
	"github.com/tranvuongduy2003/go-mvc/pkg/jwt"
)

// authService implements the AuthService interface
type authService struct {
	userRepo        repositories.UserRepository
	jwtService      jwt.JWTService
	passwordHasher  *security.PasswordHasher
	tokenGenerator  *security.TokenGenerator
	cacheService    *cache.Service
	tokenBlacklist  string // Redis key prefix for blacklisted tokens
	verificationTTL time.Duration
	resetTokenTTL   time.Duration
}

// NewAuthService creates a new authentication service
func NewAuthService(
	userRepo repositories.UserRepository,
	jwtService jwt.JWTService,
	passwordHasher *security.PasswordHasher,
	cacheService *cache.Service,
) services.AuthService {
	return &authService{
		userRepo:        userRepo,
		jwtService:      jwtService,
		passwordHasher:  passwordHasher,
		tokenGenerator:  security.NewTokenGenerator(),
		cacheService:    cacheService,
		tokenBlacklist:  "blacklist:token:",
		verificationTTL: 24 * time.Hour, // Email verification valid for 24 hours
		resetTokenTTL:   1 * time.Hour,  // Password reset valid for 1 hour
	}
}

// Register creates a new user account
func (s *authService) Register(ctx context.Context, req *services.RegisterRequest) (*services.AuthenticatedUser, error) {
	// Validate request
	if err := s.validateRegisterRequest(req); err != nil {
		return nil, fmt.Errorf("invalid registration data: %w", err)
	}

	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	// Create user domain entity
	userEntity, err := user.NewUser(
		req.Email,
		req.Name,
		req.Phone,
		req.Password, // Use plain password, domain will hash it
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user entity: %w", err)
	}

	// Save user to repository
	if err := s.userRepo.Create(ctx, userEntity); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Parse user ID as UUID for token generation
	userID, err := uuid.Parse(userEntity.ID())
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	// Generate tokens
	tokens, err := s.generateTokens(userID, userEntity.Email())
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &services.AuthenticatedUser{
		User:   userEntity,
		Tokens: tokens,
	}, nil
}

// Login authenticates a user with email and password
func (s *authService) Login(ctx context.Context, credentials *services.LoginCredentials) (*services.AuthenticatedUser, error) {
	// Validate credentials
	if err := s.validateLoginCredentials(credentials); err != nil {
		return nil, fmt.Errorf("invalid credentials: %w", err)
	}

	// Get user by email
	userEntity, err := s.userRepo.GetByEmail(ctx, credentials.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if userEntity == nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Check if user is active
	if !userEntity.IsActive() {
		return nil, fmt.Errorf("user account is inactive")
	}

	// Verify password
	if !userEntity.VerifyPassword(credentials.Password) {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Parse user ID as UUID
	userID, err := uuid.Parse(userEntity.ID())
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	// Generate tokens
	tokens, err := s.generateTokens(userID, userEntity.Email())
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &services.AuthenticatedUser{
		User:   userEntity,
		Tokens: tokens,
	}, nil
}

// RefreshToken generates new access token using refresh token
func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*services.AuthTokens, error) {
	// Check if token is blacklisted
	if blacklisted, err := s.IsTokenBlacklisted(ctx, refreshToken); err != nil {
		return nil, fmt.Errorf("failed to check token blacklist: %w", err)
	} else if blacklisted {
		return nil, fmt.Errorf("refresh token is invalid")
	}

	// Validate refresh token and generate new access token
	newAccessToken, err := s.jwtService.RefreshAccessToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	// Parse refresh token to get expiry
	claims, err := s.jwtService.ValidateToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to parse refresh token: %w", err)
	}

	return &services.AuthTokens{
		AccessToken:           newAccessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  time.Unix(s.jwtService.GetAccessTokenExpirationTime(), 0),
		RefreshTokenExpiresAt: claims.ExpiresAt.Time,
		TokenType:             "Bearer",
	}, nil
}

// Logout invalidates user tokens
func (s *authService) Logout(ctx context.Context, userID string) error {
	// This would typically blacklist the specific token
	// For now, we'll implement a simple approach
	// In a production system, you might want to track active sessions
	return nil
}

// LogoutAll invalidates all tokens for a user across all devices
func (s *authService) LogoutAll(ctx context.Context, userID string) error {
	// This would typically blacklist all tokens for the user
	// Implementation would depend on your session management strategy
	return nil
}

// ValidateToken validates an access token and returns user info
func (s *authService) ValidateToken(ctx context.Context, accessToken string) (*user.User, error) {
	// Check if token is blacklisted
	if blacklisted, err := s.IsTokenBlacklisted(ctx, accessToken); err != nil {
		return nil, fmt.Errorf("failed to check token blacklist: %w", err)
	} else if blacklisted {
		return nil, fmt.Errorf("token is invalid")
	}

	// Validate token
	claims, err := s.jwtService.ValidateToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Check token type
	if claims.Type != "access" {
		return nil, fmt.Errorf("token is not an access token")
	}

	// Get user from repository
	userEntity, err := s.userRepo.GetByID(ctx, claims.UserID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if userEntity == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Check if user is active
	if !userEntity.IsActive() {
		return nil, fmt.Errorf("user account is inactive")
	}

	return userEntity, nil
}

// ChangePassword changes user password
func (s *authService) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	// Get user
	userEntity, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if userEntity == nil {
		return fmt.Errorf("user not found")
	}

	// Verify old password
	if !userEntity.VerifyPassword(oldPassword) {
		return fmt.Errorf("invalid old password")
	}

	// Validate new password
	if err := s.validatePassword(newPassword); err != nil {
		return fmt.Errorf("invalid new password: %w", err)
	}

	// Update password using domain method
	if err := userEntity.ChangePassword(newPassword); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Save updated user
	if err := s.userRepo.Update(ctx, userEntity); err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}

	return nil
}

// ResetPassword initiates password reset process
func (s *authService) ResetPassword(ctx context.Context, email string) error {
	// Get user by email
	userEntity, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if userEntity == nil {
		// Don't reveal that user doesn't exist
		return nil
	}

	// Generate reset token
	resetToken, err := s.tokenGenerator.Generate(32)
	if err != nil {
		return fmt.Errorf("failed to generate reset token: %w", err)
	}

	// Store reset token in cache with TTL
	cacheKey := fmt.Sprintf("password_reset:%s", resetToken)
	cacheOptions := &cache.CacheOptions{TTL: s.resetTokenTTL}

	if err := s.cacheService.Set(ctx, cacheKey, userEntity.ID(), cacheOptions); err != nil {
		return fmt.Errorf("failed to store reset token: %w", err)
	}

	// TODO: Send password reset email
	// This would typically involve sending an email with the reset token
	// For now, we'll just log it (in a real system, integrate with email service)

	return nil
}

// ConfirmPasswordReset completes password reset with token
func (s *authService) ConfirmPasswordReset(ctx context.Context, token, newPassword string) error {
	// Validate new password
	if err := s.validatePassword(newPassword); err != nil {
		return fmt.Errorf("invalid password: %w", err)
	}

	// Get user ID from reset token
	cacheKey := fmt.Sprintf("password_reset:%s", token)
	var userID string
	if err := s.cacheService.Get(ctx, cacheKey, &userID); err != nil {
		return fmt.Errorf("invalid or expired reset token")
	}

	// Get user
	userEntity, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if userEntity == nil {
		return fmt.Errorf("user not found")
	}

	// Update password using domain method
	if err := userEntity.ChangePassword(newPassword); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Save updated user
	if err := s.userRepo.Update(ctx, userEntity); err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}

	// Delete reset token from cache
	_ = s.cacheService.Delete(ctx, cacheKey) // Ignore error, not critical

	return nil
}

// VerifyEmail verifies user email with verification token
func (s *authService) VerifyEmail(ctx context.Context, token string) error {
	// Get user ID from verification token
	cacheKey := fmt.Sprintf("email_verification:%s", token)
	var userID string
	if err := s.cacheService.Get(ctx, cacheKey, &userID); err != nil {
		return fmt.Errorf("invalid or expired verification token")
	}

	// Get user
	userEntity, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if userEntity == nil {
		return fmt.Errorf("user not found")
	}

	// For now, just activate the user since we don't have email verification in domain
	// In a complete implementation, you would extend the User domain to handle email verification
	userEntity.Activate()

	// Save updated user
	if err := s.userRepo.Update(ctx, userEntity); err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}

	// Delete verification token from cache
	_ = s.cacheService.Delete(ctx, cacheKey) // Ignore error, not critical

	return nil
}

// ResendVerificationEmail resends email verification
func (s *authService) ResendVerificationEmail(ctx context.Context, email string) error {
	// Get user by email
	userEntity, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if userEntity == nil {
		// Don't reveal that user doesn't exist
		return nil
	}

	// For now, skip the email verification check since it's not implemented in domain
	// In a complete implementation, you would check if email is already verified

	// Generate verification token
	verificationToken, err := s.tokenGenerator.Generate(32)
	if err != nil {
		return fmt.Errorf("failed to generate verification token: %w", err)
	}

	// Store verification token in cache with TTL
	cacheKey := fmt.Sprintf("email_verification:%s", verificationToken)
	cacheOptions := &cache.CacheOptions{TTL: s.verificationTTL}

	if err := s.cacheService.Set(ctx, cacheKey, userEntity.ID(), cacheOptions); err != nil {
		return fmt.Errorf("failed to store verification token: %w", err)
	}

	// TODO: Send verification email
	// This would typically involve sending an email with the verification token

	return nil
}

// GetUserFromToken extracts user information from a valid token
func (s *authService) GetUserFromToken(ctx context.Context, token string) (*user.User, error) {
	return s.ValidateToken(ctx, token)
}

// IsTokenBlacklisted checks if a token is blacklisted
func (s *authService) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	cacheKey := s.tokenBlacklist + token
	var exists bool
	if err := s.cacheService.Get(ctx, cacheKey, &exists); err != nil {
		// If key doesn't exist, token is not blacklisted
		return false, nil
	}
	return exists, nil
}

// BlacklistToken adds a token to blacklist
func (s *authService) BlacklistToken(ctx context.Context, token string) error {
	// Parse token to get expiry
	claims, err := s.jwtService.ValidateToken(token)
	if err != nil {
		return fmt.Errorf("invalid token: %w", err)
	}

	// Calculate TTL based on token expiry
	ttl := time.Until(claims.ExpiresAt.Time)
	if ttl <= 0 {
		// Token already expired, no need to blacklist
		return nil
	}

	// Add to blacklist with TTL
	cacheKey := s.tokenBlacklist + token
	cacheOptions := &cache.CacheOptions{TTL: ttl}

	if err := s.cacheService.Set(ctx, cacheKey, true, cacheOptions); err != nil {
		return fmt.Errorf("failed to blacklist token: %w", err)
	}

	return nil
}

// generateTokens generates both access and refresh tokens
func (s *authService) generateTokens(userID uuid.UUID, email string) (*services.AuthTokens, error) {
	// Generate access token
	accessToken, err := s.jwtService.GenerateAccessToken(userID, email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshToken, err := s.jwtService.GenerateRefreshToken(userID, email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &services.AuthTokens{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  time.Unix(s.jwtService.GetAccessTokenExpirationTime(), 0),
		RefreshTokenExpiresAt: time.Unix(s.jwtService.GetRefreshTokenExpirationTime(), 0),
		TokenType:             "Bearer",
	}, nil
}

// Validation methods

func (s *authService) validateRegisterRequest(req *services.RegisterRequest) error {
	if req.Email == "" {
		return fmt.Errorf("email is required")
	}
	if req.Name == "" {
		return fmt.Errorf("name is required")
	}
	if req.Password == "" {
		return fmt.Errorf("password is required")
	}

	return s.validatePassword(req.Password)
}

func (s *authService) validateLoginCredentials(credentials *services.LoginCredentials) error {
	if credentials.Email == "" {
		return fmt.Errorf("email is required")
	}
	if credentials.Password == "" {
		return fmt.Errorf("password is required")
	}
	return nil
}

func (s *authService) validatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}
	// Add more password validation rules as needed
	return nil
}
