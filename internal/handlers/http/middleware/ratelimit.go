package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter represents a rate limiter for different endpoints
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rps int, burst int) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     rate.Limit(rps),
		burst:    burst,
	}
}

// getLimiter returns the rate limiter for a given key (IP address)
func (rl *RateLimiter) getLimiter(key string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.limiters[key]
	if !exists {
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.limiters[key] = limiter
	}

	return limiter
}

// RateLimitMiddleware returns a Gin middleware for rate limiting
func (rl *RateLimiter) RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Use client IP as the key
		key := c.ClientIP()
		limiter := rl.getLimiter(key)

		if !limiter.Allow() {
			c.Header("X-RateLimit-Limit", fmt.Sprintf("%.0f", float64(rl.rate)))
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(time.Second).Unix()))

			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Rate limit exceeded",
				"message": "Too many requests, please try again later",
				"code":    "RATE_LIMIT_EXCEEDED",
			})
			c.Abort()
			return
		}

		// Add rate limit headers
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%.0f", float64(rl.rate)))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%.0f", limiter.Tokens()))

		c.Next()
	}
}

// CleanupOldLimiters removes old limiters to prevent memory leaks
func (rl *RateLimiter) CleanupOldLimiters() {
	ticker := time.NewTicker(time.Hour)
	go func() {
		for range ticker.C {
			rl.mu.Lock()
			// Clear all limiters - they will be recreated as needed
			rl.limiters = make(map[string]*rate.Limiter)
			rl.mu.Unlock()
		}
	}()
}

// GlobalRateLimitMiddleware creates a simple global rate limit middleware
func GlobalRateLimitMiddleware(rps int, burst int) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Limit(rps), burst)

	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Global rate limit exceeded",
				"message": "Server is receiving too many requests, please try again later",
				"code":    "GLOBAL_RATE_LIMIT_EXCEEDED",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// APIKeyRateLimitMiddleware creates rate limit based on API key
func APIKeyRateLimitMiddleware(rps int, burst int) gin.HandlerFunc {
	limiters := make(map[string]*rate.Limiter)
	mu := sync.RWMutex{}

	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			apiKey = c.ClientIP() // Fallback to IP if no API key
		}

		mu.Lock()
		limiter, exists := limiters[apiKey]
		if !exists {
			limiter = rate.NewLimiter(rate.Limit(rps), burst)
			limiters[apiKey] = limiter
		}
		mu.Unlock()

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "API rate limit exceeded",
				"message": "API key rate limit exceeded, please try again later",
				"code":    "API_RATE_LIMIT_EXCEEDED",
				"api_key": apiKey,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
