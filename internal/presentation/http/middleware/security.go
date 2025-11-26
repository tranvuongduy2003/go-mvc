package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")

		c.Header("X-Frame-Options", "DENY")

		c.Header("X-XSS-Protection", "1; mode=block")

		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self'; connect-src 'self'; frame-ancestors 'none';")

		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		c.Header("Permissions-Policy", "camera=(), microphone=(), geolocation=(), interest-cohort=()")

		c.Header("X-Permitted-Cross-Domain-Policies", "none")

		if c.Request.URL.Path != "/health" && c.Request.URL.Path != "/metrics" {
			c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
			c.Header("Pragma", "no-cache")
			c.Header("Expires", "0")
		}

		c.Next()
	}
}

func SecureHeaders(config SecurityConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", config.FrameOptions)
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", config.ReferrerPolicy)
		c.Header("X-Permitted-Cross-Domain-Policies", "none")

		if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
			c.Header("Strict-Transport-Security", config.HSTSMaxAge)
		}

		if config.CSP != "" {
			c.Header("Content-Security-Policy", config.CSP)
		}

		if config.PermissionsPolicy != "" {
			c.Header("Permissions-Policy", config.PermissionsPolicy)
		}

		for key, value := range config.CustomHeaders {
			c.Header(key, value)
		}

		c.Next()
	}
}

type SecurityConfig struct {
	FrameOptions      string
	ReferrerPolicy    string
	HSTSMaxAge        string
	CSP               string
	PermissionsPolicy string
	CustomHeaders     map[string]string
}

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

		if clientName, exists := validAPIKeys[apiKey]; exists {
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

func BasicAuthMiddleware(users map[string]string) gin.HandlerFunc {
	return gin.BasicAuth(users)
}

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
