package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/tranvuongduy2003/go-mvc/internal/shared/config"
)

// JWTService defines the interface for JWT operations
type JWTService interface {
	GenerateAccessToken(userID uuid.UUID, email string) (string, error)
	GenerateRefreshToken(userID uuid.UUID, email string) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
	RefreshAccessToken(refreshToken string) (string, error)
	GetAccessTokenExpirationTime() int64
	GetRefreshTokenExpirationTime() int64
}

// Service handles JWT operations
type Service struct {
	secret        []byte
	accessExpiry  time.Duration
	refreshExpiry time.Duration
	issuer        string
	audience      string
}

// Ensure Service implements JWTService interface
var _ JWTService = (*Service)(nil)

// NewService creates a new JWT service
func NewService(cfg config.JWT) *Service {
	return &Service{
		secret:        []byte(cfg.Secret),
		accessExpiry:  cfg.AccessExpiry,
		refreshExpiry: cfg.RefreshExpiry,
		issuer:        cfg.Issuer,
		audience:      cfg.Audience,
	}
}

// Claims represents JWT claims
type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Type   string    `json:"type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

// GenerateAccessToken generates an access token
func (s *Service) GenerateAccessToken(userID uuid.UUID, email string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		Email:  email,
		Type:   "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Audience:  jwt.ClaimStrings{s.audience},
			Subject:   userID.String(),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessExpiry)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

// GenerateRefreshToken generates a refresh token
func (s *Service) GenerateRefreshToken(userID uuid.UUID, email string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		Email:  email,
		Type:   "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Audience:  jwt.ClaimStrings{s.audience},
			Subject:   userID.String(),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshExpiry)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

// ValidateToken validates a token and returns claims
func (s *Service) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// RefreshAccessToken generates a new access token from a valid refresh token
func (s *Service) RefreshAccessToken(refreshToken string) (string, error) {
	claims, err := s.ValidateToken(refreshToken)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %w", err)
	}

	if claims.Type != "refresh" {
		return "", fmt.Errorf("token is not a refresh token")
	}

	// Generate new access token
	return s.GenerateAccessToken(claims.UserID, claims.Email)
}

// GetAccessTokenExpirationTime returns the expiration time for access tokens
func (s *Service) GetAccessTokenExpirationTime() int64 {
	return time.Now().Add(s.accessExpiry).Unix()
}

// GetRefreshTokenExpirationTime returns the expiration time for refresh tokens
func (s *Service) GetRefreshTokenExpirationTime() int64 {
	return time.Now().Add(s.refreshExpiry).Unix()
}

// ExtractTokenFromHeader extracts token from Authorization header
func ExtractTokenFromHeader(authHeader string) string {
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}
	return ""
}
