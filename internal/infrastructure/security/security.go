package security

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// PasswordHasher handles password hashing and verification
type PasswordHasher struct {
	cost int
}

// NewPasswordHasher creates a new password hasher
func NewPasswordHasher(cost int) *PasswordHasher {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = bcrypt.DefaultCost
	}
	return &PasswordHasher{cost: cost}
}

// Hash hashes a password using bcrypt
func (h *PasswordHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(bytes), nil
}

// Verify verifies a password against its hash
func (h *PasswordHasher) Verify(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// TokenGenerator generates secure random tokens
type TokenGenerator struct{}

// NewTokenGenerator creates a new token generator
func NewTokenGenerator() *TokenGenerator {
	return &TokenGenerator{}
}

// Generate generates a random token of specified length
func (g *TokenGenerator) Generate(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random token: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateAPIKey generates a secure API key
func (g *TokenGenerator) GenerateAPIKey() (string, error) {
	return g.Generate(32) // 256-bit key
}

// GenerateSessionToken generates a session token
func (g *TokenGenerator) GenerateSessionToken() (string, error) {
	return g.Generate(16) // 128-bit token
}

// SecureCompare performs constant-time string comparison
func SecureCompare(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

// Sanitizer provides input sanitization
type Sanitizer struct{}

// NewSanitizer creates a new sanitizer
func NewSanitizer() *Sanitizer {
	return &Sanitizer{}
}

// SanitizeString removes potentially dangerous characters from a string
func (s *Sanitizer) SanitizeString(input string) string {
	// Basic sanitization - remove null bytes and control characters
	result := ""
	for _, char := range input {
		if char >= 32 && char <= 126 || char >= 160 {
			result += string(char)
		}
	}
	return result
}

// ValidatePassword validates password strength
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasDigit = true
		case char >= 32 && char <= 126:
			hasSpecial = true
		}
	}

	if !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	if !hasDigit {
		return fmt.Errorf("password must contain at least one digit")
	}
	if !hasSpecial {
		return fmt.Errorf("password must contain at least one special character")
	}

	return nil
}

// RateLimiter provides simple rate limiting
type RateLimiter struct {
	requests map[string][]int64
	limit    int
	window   int64
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(limit int, windowSeconds int64) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]int64),
		limit:    limit,
		window:   windowSeconds,
	}
}

// Allow checks if a request is allowed for the given key
func (rl *RateLimiter) Allow(key string, timestamp int64) bool {
	requests, exists := rl.requests[key]
	if !exists {
		requests = []int64{}
	}

	// Remove old requests outside the window
	validRequests := []int64{}
	for _, req := range requests {
		if timestamp-req < rl.window {
			validRequests = append(validRequests, req)
		}
	}

	// Check if we can add a new request
	if len(validRequests) < rl.limit {
		validRequests = append(validRequests, timestamp)
		rl.requests[key] = validRequests
		return true
	}

	rl.requests[key] = validRequests
	return false
}

// CSRFToken generates CSRF tokens
type CSRFToken struct {
	generator *TokenGenerator
}

// NewCSRFToken creates a new CSRF token generator
func NewCSRFToken() *CSRFToken {
	return &CSRFToken{
		generator: NewTokenGenerator(),
	}
}

// Generate generates a new CSRF token
func (c *CSRFToken) Generate() (string, error) {
	return c.generator.Generate(16)
}

// Validate validates a CSRF token (basic implementation)
func (c *CSRFToken) Validate(token, expected string) bool {
	return SecureCompare(token, expected)
}
