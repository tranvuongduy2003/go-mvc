package security

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher struct {
	cost int
}

func NewPasswordHasher(cost int) *PasswordHasher {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = bcrypt.DefaultCost
	}
	return &PasswordHasher{cost: cost}
}

func (h *PasswordHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(bytes), nil
}

func (h *PasswordHasher) Verify(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type TokenGenerator struct{}

func NewTokenGenerator() *TokenGenerator {
	return &TokenGenerator{}
}

func (g *TokenGenerator) Generate(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random token: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

func (g *TokenGenerator) GenerateAPIKey() (string, error) {
	return g.Generate(32) // 256-bit key
}

func (g *TokenGenerator) GenerateSessionToken() (string, error) {
	return g.Generate(16) // 128-bit token
}

func SecureCompare(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

type Sanitizer struct{}

func NewSanitizer() *Sanitizer {
	return &Sanitizer{}
}

func (s *Sanitizer) SanitizeString(input string) string {
	result := ""
	for _, char := range input {
		if char >= 32 && char <= 126 || char >= 160 {
			result += string(char)
		}
	}
	return result
}

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

type RateLimiter struct {
	requests map[string][]int64
	limit    int
	window   int64
}

func NewRateLimiter(limit int, windowSeconds int64) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]int64),
		limit:    limit,
		window:   windowSeconds,
	}
}

func (rl *RateLimiter) Allow(key string, timestamp int64) bool {
	requests, exists := rl.requests[key]
	if !exists {
		requests = []int64{}
	}

	validRequests := []int64{}
	for _, req := range requests {
		if timestamp-req < rl.window {
			validRequests = append(validRequests, req)
		}
	}

	if len(validRequests) < rl.limit {
		validRequests = append(validRequests, timestamp)
		rl.requests[key] = validRequests
		return true
	}

	rl.requests[key] = validRequests
	return false
}

type CSRFToken struct {
	generator *TokenGenerator
}

func NewCSRFToken() *CSRFToken {
	return &CSRFToken{
		generator: NewTokenGenerator(),
	}
}

func (c *CSRFToken) Generate() (string, error) {
	return c.generator.Generate(16)
}

func (c *CSRFToken) Validate(token, expected string) bool {
	return SecureCompare(token, expected)
}
