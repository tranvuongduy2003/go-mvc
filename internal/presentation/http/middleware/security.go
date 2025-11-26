package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SecurityHeaders adds security headers to responses
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// X-Content-Type-Options: Prevent MIME sniffing
		c.Header("X-Content-Type-Options", "nosniff")

		// X-Frame-Options: Prevent clickjacking
		c.Header("X-Frame-Options", "DENY")

		// X-XSS-Protection: Enable XSS protection
		c.Header("X-XSS-Protection", "1; mode=block")

		// Strict-Transport-Security: Force HTTPS
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		// Content-Security-Policy: Prevent XSS and injection attacks
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self'; connect-src 'self'; frame-ancestors 'none';")

		// Referrer-Policy: Control referrer information
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// Permissions-Policy: Control browser features
		c.Header("Permissions-Policy", "camera=(), microphone=(), geolocation=(), interest-cohort=()")

		// X-Permitted-Cross-Domain-Policies: Control Adobe Flash/PDF cross-domain access
		c.Header("X-Permitted-Cross-Domain-Policies", "none")

		// Cache-Control: Prevent caching of sensitive data
		if c.Request.URL.Path != "/health" && c.Request.URL.Path != "/metrics" {
			c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
			c.Header("Pragma", "no-cache")
			c.Header("Expires", "0")
		}

		c.Next()
	}
}

// SecureHeaders adds comprehensive security headers for production
func SecureHeaders(config SecurityConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Basic security headers
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", config.FrameOptions)
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", config.ReferrerPolicy)
		c.Header("X-Permitted-Cross-Domain-Policies", "none")

		// HSTS (only if HTTPS)
		if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
			c.Header("Strict-Transport-Security", config.HSTSMaxAge)
		}

		// CSP
		if config.CSP != "" {
			c.Header("Content-Security-Policy", config.CSP)
		}

		// Permissions Policy
		if config.PermissionsPolicy != "" {
			c.Header("Permissions-Policy", config.PermissionsPolicy)
		}

		// Custom headers
		for key, value := range config.CustomHeaders {
			c.Header(key, value)
		}

		c.Next()
	}
}

// SecurityConfig represents security middleware configuration
type SecurityConfig struct {
	FrameOptions      string
	ReferrerPolicy    string
	HSTSMaxAge        string
	CSP               string
	PermissionsPolicy string
	CustomHeaders     map[string]string
}

// DefaultSecurityConfig returns default security configuration
func DefaultSecurityConfig() SecurityConfig {
	return SecurityConfig{
		FrameOptions:      "DENY",
		ReferrerPolicy:    "strict-origin-when-cross-origin",
		HSTSMaxAge:        "max-age=31536000; includeSubDomains",
		CSP:               "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self'; connect-src 'self'; frame-ancestors 'none';",
		PermissionsPolicy: "camera=(), microphone=(), geolocation=(), interest-cohort=()",
		CustomHeaders: map[string]string{
			"X-API-Version": "v1",
		},
	}
}

// APIKeyAuthMiddleware validates API keys
func APIKeyAuthMiddleware(validAPIKeys map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Missing API key",
				"message": "X-API-Key header is required",
				"code":    "MISSING_API_KEY",
			})
			c.Abort()
			return
		}

		// Validate API key
		if clientName, exists := validAPIKeys[apiKey]; exists {
			// Set client information in context
			c.Set("api_key", apiKey)
			c.Set("client_name", clientName)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid API key",
				"message": "The provided API key is not valid",
				"code":    "INVALID_API_KEY",
			})
			c.Abort()
		}
	}
}

// BasicAuthMiddleware provides basic HTTP authentication
func BasicAuthMiddleware(users map[string]string) gin.HandlerFunc {
	return gin.BasicAuth(users)
}

// IPWhitelistMiddleware restricts access to whitelisted IPs
func IPWhitelistMiddleware(allowedIPs []string) gin.HandlerFunc {
	ipMap := make(map[string]bool)
	for _, ip := range allowedIPs {
		ipMap[ip] = true
	}

	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		if !ipMap[clientIP] {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Access denied",
				"message": "Your IP address is not allowed to access this resource",
				"code":    "IP_NOT_ALLOWED",
				"ip":      clientIP,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// SizeLimit middleware limits request body size
func SizeLimit(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > maxSize {
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"error":    "Request too large",
				"message":  "Request body exceeds maximum allowed size",
				"code":     "REQUEST_TOO_LARGE",
				"max_size": maxSize,
			})
			c.Abort()
			return
		}

		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)
		c.Next()
	}
}

// ContentTypeMiddleware enforces content type validation
func ContentTypeMiddleware(allowedTypes []string) gin.HandlerFunc {
	typeMap := make(map[string]bool)
	for _, contentType := range allowedTypes {
		typeMap[contentType] = true
	}

	return func(c *gin.Context) {
		if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut || c.Request.Method == http.MethodPatch {
			contentType := c.GetHeader("Content-Type")
			if contentType == "" || !typeMap[contentType] {
				c.JSON(http.StatusUnsupportedMediaType, gin.H{
					"error":         "Unsupported media type",
					"message":       "Content-Type header must be one of the allowed types",
					"code":          "UNSUPPORTED_MEDIA_TYPE",
					"allowed_types": allowedTypes,
					"received_type": contentType,
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
