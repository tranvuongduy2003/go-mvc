package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/contracts"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/user"
	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/cache"
	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/external"
	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/logger"
	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/security"
	apperrors "github.com/tranvuongduy2003/go-mvc/pkg/errors"
	"github.com/tranvuongduy2003/go-mvc/pkg/jwt"
)

type AuthService struct {
	userRepo        user.UserRepository
	jwtService      jwt.JWTService
	passwordHasher  *security.PasswordHasher
	tokenGenerator  *security.TokenGenerator
	cacheService    *cache.Service
	smtpService     *external.SMTPService
	logger          *logger.Logger
	tokenBlacklist  string // Redis key prefix for blacklisted tokens
	verificationTTL time.Duration
	resetTokenTTL   time.Duration
}

var _ contracts.AuthService = (*AuthService)(nil)
var _ contracts.TokenManagementService = (*AuthService)(nil)
var _ contracts.PasswordManagementService = (*AuthService)(nil)
var _ contracts.EmailVerificationService = (*AuthService)(nil)

func NewAuthService(
	userRepo user.UserRepository,
	jwtService jwt.JWTService,
	passwordHasher *security.PasswordHasher,
	cacheService *cache.Service,
	smtpService *external.SMTPService,
	logger *logger.Logger,
) *AuthService {
	return &AuthService{
		userRepo:        userRepo,
		jwtService:      jwtService,
		passwordHasher:  passwordHasher,
		tokenGenerator:  security.NewTokenGenerator(),
		cacheService:    cacheService,
		smtpService:     smtpService,
		logger:          logger,
		tokenBlacklist:  "blacklist:token:",
		verificationTTL: 24 * time.Hour, // Email verification valid for 24 hours
		resetTokenTTL:   1 * time.Hour,  // Password reset valid for 1 hour
	}
}

func (s *AuthService) Register(ctx context.Context, req *contracts.RegisterRequest) (*contracts.AuthenticatedUser, error) {
	if err := s.validateRegisterRequest(req); err != nil {
		return nil, fmt.Errorf("invalid registration data: %w", err)
	}

	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	userEntity, err := user.NewUser(
		req.Email,
		req.Name,
		req.Phone,
		req.Password, // Use plain password, domain will hash it
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user entity: %w", err)
	}

	if err := s.userRepo.Create(ctx, userEntity); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	userID, err := uuid.Parse(userEntity.ID())
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	tokens, err := s.generateTokens(userID, userEntity.Email())
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return &contracts.AuthenticatedUser{
		User:   userEntity,
		Tokens: tokens,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, credentials *contracts.LoginCredentials) (*contracts.AuthenticatedUser, error) {
	if err := s.validateLoginCredentials(credentials); err != nil {
		return nil, apperrors.NewValidationError("invalid credentials", err)
	}

	userEntity, err := s.userRepo.GetByEmail(ctx, credentials.Email)
	if err != nil {
		return nil, apperrors.NewInternalError("failed to get user", err)
	}
	if userEntity == nil {
		return nil, apperrors.NewUnauthorizedError("invalid email or password")
	}

	if !userEntity.IsActive() {
		return nil, apperrors.NewUnauthorizedError("user account is inactive")
	}

	if !userEntity.VerifyPassword(credentials.Password) {
		return nil, apperrors.NewUnauthorizedError("invalid email or password")
	}

	userID, err := uuid.Parse(userEntity.ID())
	if err != nil {
		return nil, apperrors.NewInternalError("invalid user ID format", err)
	}

	tokens, err := s.generateTokens(userID, userEntity.Email())
	if err != nil {
		return nil, apperrors.NewInternalError("failed to generate tokens", err)
	}

	return &contracts.AuthenticatedUser{
		User:   userEntity,
		Tokens: tokens,
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*contracts.AuthTokens, error) {
	if blacklisted, err := s.IsTokenBlacklisted(ctx, refreshToken); err != nil {
		return nil, apperrors.NewInternalError("failed to check token blacklist", err)
	} else if blacklisted {
		return nil, apperrors.NewUnauthorizedError("refresh token is invalid")
	}

	newAccessToken, err := s.jwtService.RefreshAccessToken(refreshToken)
	if err != nil {
		return nil, err // Already an AppError from jwt service
	}

	claims, err := s.jwtService.ValidateToken(refreshToken)
	if err != nil {
		return nil, err // Already an AppError from jwt service
	}

	return &contracts.AuthTokens{
		AccessToken:           newAccessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  time.Unix(s.jwtService.GetAccessTokenExpirationTime(), 0),
		RefreshTokenExpiresAt: claims.ExpiresAt.Time,
		TokenType:             "Bearer",
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, userID string) error {
	return nil
}

func (s *AuthService) LogoutAll(ctx context.Context, userID string) error {
	return nil
}

func (s *AuthService) ValidateToken(ctx context.Context, accessToken string) (*user.User, error) {
	if blacklisted, err := s.IsTokenBlacklisted(ctx, accessToken); err != nil {
		return nil, apperrors.NewInternalError("failed to check token blacklist", err)
	} else if blacklisted {
		return nil, apperrors.NewUnauthorizedError("token is invalid")
	}

	claims, err := s.jwtService.ValidateToken(accessToken)
	if err != nil {
		return nil, apperrors.NewUnauthorizedError("invalid token")
	}

	if claims.Type != "access" {
		return nil, apperrors.NewUnauthorizedError("token is not an access token")
	}

	userEntity, err := s.userRepo.GetByID(ctx, claims.UserID.String())
	if err != nil {
		return nil, apperrors.NewInternalError("failed to get user", err)
	}
	if userEntity == nil {
		return nil, apperrors.NewUnauthorizedError("user not found")
	}

	if !userEntity.IsActive() {
		return nil, apperrors.NewUnauthorizedError("user account is inactive")
	}

	return userEntity, nil
}

func (s *AuthService) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	userEntity, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if userEntity == nil {
		return fmt.Errorf("user not found")
	}

	if !userEntity.VerifyPassword(oldPassword) {
		return fmt.Errorf("invalid old password")
	}

	if err := s.validatePassword(newPassword); err != nil {
		return fmt.Errorf("invalid new password: %w", err)
	}

	if err := userEntity.ChangePassword(newPassword); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	if err := s.userRepo.Update(ctx, userEntity); err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}

	return nil
}

func (s *AuthService) ResetPassword(ctx context.Context, email string) error {
	userEntity, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return apperrors.NewInternalError("failed to get user", err)
	}
	if userEntity == nil {
		return nil
	}

	resetToken, err := s.tokenGenerator.Generate(32)
	if err != nil {
		return apperrors.NewInternalError("failed to generate reset token", err)
	}

	cacheKey := fmt.Sprintf("password_reset:%s", resetToken)
	cacheOptions := &cache.CacheOptions{TTL: s.resetTokenTTL}

	if err := s.cacheService.Set(ctx, cacheKey, userEntity.ID(), cacheOptions); err != nil {
		return apperrors.NewInternalError("failed to store reset token", err)
	}

	if err := s.smtpService.SendPasswordResetEmail(ctx, userEntity.Email(), userEntity.Name(), resetToken); err != nil {
		s.logger.Errorf("Failed to send password reset email: %v", err)
	}

	return nil
}

func (s *AuthService) ConfirmPasswordReset(ctx context.Context, token, newPassword string) error {
	if err := s.validatePassword(newPassword); err != nil {
		return apperrors.NewValidationError("invalid password", err)
	}

	cacheKey := fmt.Sprintf("password_reset:%s", token)
	var userID string
	if err := s.cacheService.Get(ctx, cacheKey, &userID); err != nil {
		return apperrors.NewValidationError("invalid or expired reset token", err)
	}

	userEntity, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return apperrors.NewInternalError("failed to get user", err)
	}
	if userEntity == nil {
		return apperrors.NewNotFoundError("user not found")
	}

	if err := userEntity.ChangePassword(newPassword); err != nil {
		return apperrors.NewInternalError("failed to update password", err)
	}

	if err := s.userRepo.Update(ctx, userEntity); err != nil {
		return apperrors.NewInternalError("failed to save user", err)
	}

	_ = s.cacheService.Delete(ctx, cacheKey) // Ignore error, not critical

	return nil
}

func (s *AuthService) VerifyEmail(ctx context.Context, token string) error {
	cacheKey := fmt.Sprintf("email_verification:%s", token)
	var userID string
	if err := s.cacheService.Get(ctx, cacheKey, &userID); err != nil {
		return apperrors.NewValidationError("invalid or expired verification token", err)
	}

	userEntity, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return apperrors.NewInternalError("failed to get user", err)
	}
	if userEntity == nil {
		return apperrors.NewNotFoundError("user not found")
	}

	userEntity.Activate()

	if err := s.userRepo.Update(ctx, userEntity); err != nil {
		return apperrors.NewInternalError("failed to save user", err)
	}

	_ = s.cacheService.Delete(ctx, cacheKey) // Ignore error, not critical

	return nil
}

func (s *AuthService) ResendVerificationEmail(ctx context.Context, email string) error {
	userEntity, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return apperrors.NewInternalError("failed to get user", err)
	}
	if userEntity == nil {
		return nil
	}

	verificationToken, err := s.tokenGenerator.Generate(32)
	if err != nil {
		return apperrors.NewInternalError("failed to generate verification token", err)
	}

	cacheKey := fmt.Sprintf("email_verification:%s", verificationToken)
	cacheOptions := &cache.CacheOptions{TTL: s.verificationTTL}

	if err := s.cacheService.Set(ctx, cacheKey, userEntity.ID(), cacheOptions); err != nil {
		return apperrors.NewInternalError("failed to store verification token", err)
	}

	if err := s.smtpService.SendVerificationEmail(ctx, userEntity.Email(), userEntity.Name(), verificationToken); err != nil {
		s.logger.Errorf("Failed to send verification email: %v", err)
	}

	return nil
}

func (s *AuthService) GetUserFromToken(ctx context.Context, token string) (*user.User, error) {
	return s.ValidateToken(ctx, token)
}

func (s *AuthService) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	cacheKey := s.tokenBlacklist + token
	var exists bool
	if err := s.cacheService.Get(ctx, cacheKey, &exists); err != nil {
		return false, nil
	}
	return exists, nil
}

func (s *AuthService) BlacklistToken(ctx context.Context, token string) error {
	claims, err := s.jwtService.ValidateToken(token)
	if err != nil {
		return fmt.Errorf("invalid token: %w", err)
	}

	ttl := time.Until(claims.ExpiresAt.Time)
	if ttl <= 0 {
		return nil
	}

	cacheKey := s.tokenBlacklist + token
	cacheOptions := &cache.CacheOptions{TTL: ttl}

	if err := s.cacheService.Set(ctx, cacheKey, true, cacheOptions); err != nil {
		return fmt.Errorf("failed to blacklist token: %w", err)
	}

	return nil
}

func (s *AuthService) generateTokens(userID uuid.UUID, email string) (*contracts.AuthTokens, error) {
	accessToken, err := s.jwtService.GenerateAccessToken(userID, email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(userID, email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &contracts.AuthTokens{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  time.Unix(s.jwtService.GetAccessTokenExpirationTime(), 0),
		RefreshTokenExpiresAt: time.Unix(s.jwtService.GetRefreshTokenExpirationTime(), 0),
		TokenType:             "Bearer",
	}, nil
}

func (s *AuthService) validateRegisterRequest(req *contracts.RegisterRequest) error {
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

func (s *AuthService) validateLoginCredentials(credentials *contracts.LoginCredentials) error {
	if credentials.Email == "" {
		return fmt.Errorf("email is required")
	}
	if credentials.Password == "" {
		return fmt.Errorf("password is required")
	}
	return nil
}

func (s *AuthService) validatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}
	return nil
}
