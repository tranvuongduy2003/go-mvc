package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/config"
	apperrors "github.com/tranvuongduy2003/go-mvc/pkg/errors"
)

type JWTService interface {
	GenerateAccessToken(userID uuid.UUID, email string) (string, error)
	GenerateRefreshToken(userID uuid.UUID, email string) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
	RefreshAccessToken(refreshToken string) (string, error)
	GetAccessTokenExpirationTime() int64
	GetRefreshTokenExpirationTime() int64
}

type Service struct {
	secret        []byte
	accessExpiry  time.Duration
	refreshExpiry time.Duration
	issuer        string
	audience      string
}

var _ JWTService = (*Service)(nil)

func NewService(cfg config.JWT) *Service {
	return &Service{
		secret:        []byte(cfg.Secret),
		accessExpiry:  cfg.AccessExpiry,
		refreshExpiry: cfg.RefreshExpiry,
		issuer:        cfg.Issuer,
		audience:      cfg.Audience,
	}
}

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Type   string    `json:"type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

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

func (s *Service) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apperrors.NewUnauthorizedError(fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]))
		}
		return s.secret, nil
	})

	if err != nil {
		return nil, apperrors.NewUnauthorizedError("failed to parse token")
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, apperrors.NewUnauthorizedError("invalid token")
}

func (s *Service) RefreshAccessToken(refreshToken string) (string, error) {
	claims, err := s.ValidateToken(refreshToken)
	if err != nil {
		return "", apperrors.NewUnauthorizedError("invalid refresh token")
	}

	if claims.Type != "refresh" {
		return "", apperrors.NewUnauthorizedError("token is not a refresh token")
	}

	return s.GenerateAccessToken(claims.UserID, claims.Email)
}

func (s *Service) GetAccessTokenExpirationTime() int64 {
	return time.Now().Add(s.accessExpiry).Unix()
}

func (s *Service) GetRefreshTokenExpirationTime() int64 {
	return time.Now().Add(s.refreshExpiry).Unix()
}

func ExtractTokenFromHeader(authHeader string) string {
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}
	return ""
}
