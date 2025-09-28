package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/tranvuongduy2003/go-mvc/internal/application/services/messaging"
	"github.com/tranvuongduy2003/go-mvc/pkg/response"
)

const (
	IdempotencyKeyHeader = "Idempotency-Key"
	DefaultTTL           = 24 * time.Hour
)

// IdempotencyMiddleware provides idempotency for HTTP requests
type IdempotencyMiddleware struct {
	inboxService *messaging.InboxService
	logger       *zap.Logger
	ttl          time.Duration
}

// NewIdempotencyMiddleware creates a new idempotency middleware
func NewIdempotencyMiddleware(
	inboxService *messaging.InboxService,
	logger *zap.Logger,
	ttl time.Duration,
) *IdempotencyMiddleware {
	if ttl == 0 {
		ttl = DefaultTTL
	}

	return &IdempotencyMiddleware{
		inboxService: inboxService,
		logger:       logger,
		ttl:          ttl,
	}
}

// Handler returns the Gin middleware handler
func (m *IdempotencyMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only apply to non-idempotent methods
		if !m.shouldApplyIdempotency(c.Request.Method) {
			c.Next()
			return
		}

		idempotencyKey := c.GetHeader(IdempotencyKeyHeader)
		if idempotencyKey == "" {
			// No idempotency key provided, continue normally
			c.Next()
			return
		}

		// Create unique message ID from idempotency key + request details
		messageID := m.generateMessageID(c, idempotencyKey)
		consumerID := m.generateConsumerID(c)
		eventType := m.generateEventType(c)

		ctx := c.Request.Context()

		// Check if this request has already been processed
		shouldProcess, err := m.inboxService.ProcessMessageWithDeduplication(
			ctx,
			messageID,
			eventType,
			consumerID,
			m.ttl,
		)

		if err != nil {
			if m.logger != nil {
				m.logger.Error("Idempotency check failed",
					zap.Error(err),
					zap.String("idempotency_key", idempotencyKey),
					zap.String("message_id", messageID.String()))
			}

			c.JSON(http.StatusInternalServerError, response.APIResponse{
				Success: false,
				Error: &response.ErrorInfo{
					Type:    "internal_error",
					Message: "Idempotency check failed",
				},
				Timestamp: time.Now().UTC(),
			})
			c.Abort()
			return
		}

		if !shouldProcess {
			// Request already processed, return cached response or 409 Conflict
			if m.logger != nil {
				m.logger.Info("Duplicate request detected",
					zap.String("idempotency_key", idempotencyKey),
					zap.String("message_id", messageID.String()),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method))
			}

			c.JSON(http.StatusConflict, response.APIResponse{
				Success: false,
				Error: &response.ErrorInfo{
					Type:    "conflict",
					Message: "Request already processed",
				},
				Timestamp: time.Now().UTC(),
			})
			c.Abort()
			return
		}

		// Store idempotency context for handlers to access
		c.Set("idempotency_key", idempotencyKey)
		c.Set("message_id", messageID.String())

		if m.logger != nil {
			m.logger.Info("Processing request with idempotency",
				zap.String("idempotency_key", idempotencyKey),
				zap.String("message_id", messageID.String()),
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method))
		}

		c.Next()
	}
}

// shouldApplyIdempotency determines if idempotency should be applied to this method
func (m *IdempotencyMiddleware) shouldApplyIdempotency(method string) bool {
	// Apply idempotency to non-idempotent HTTP methods
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return true
	default:
		return false
	}
}

// generateMessageID creates a deterministic message ID from the request
func (m *IdempotencyMiddleware) generateMessageID(c *gin.Context, idempotencyKey string) uuid.UUID {
	// Create a deterministic UUID from idempotency key + request signature
	hasher := sha256.New()
	hasher.Write([]byte(idempotencyKey))
	hasher.Write([]byte(c.Request.Method))
	hasher.Write([]byte(c.Request.URL.Path))

	// Include query parameters for uniqueness
	hasher.Write([]byte(c.Request.URL.RawQuery))

	hash := hex.EncodeToString(hasher.Sum(nil))

	// Create a deterministic UUID from the hash
	return uuid.NewSHA1(uuid.Nil, []byte(hash))
}

// generateConsumerID creates a consumer ID for this service instance
func (m *IdempotencyMiddleware) generateConsumerID(c *gin.Context) string {
	// Use service name + instance ID or just service name
	return "http-api"
}

// generateEventType creates an event type from the HTTP request
func (m *IdempotencyMiddleware) generateEventType(c *gin.Context) string {
	return fmt.Sprintf("http.%s.%s",
		c.Request.Method,
		c.FullPath()) // Use route pattern, not actual path
}

// IdempotencyOptions defines configuration for idempotency middleware
type IdempotencyOptions struct {
	TTL            time.Duration
	RequireKey     bool
	CacheResponses bool
	IgnoredPaths   []string
	IgnoredMethods []string
}

// WithOptions creates middleware with custom options
func (m *IdempotencyMiddleware) WithOptions(options IdempotencyOptions) gin.HandlerFunc {
	if options.TTL != 0 {
		m.ttl = options.TTL
	}

	return func(c *gin.Context) {
		// Check if path should be ignored
		if m.shouldIgnorePath(c.Request.URL.Path, options.IgnoredPaths) {
			c.Next()
			return
		}

		// Check if method should be ignored
		if m.shouldIgnoreMethod(c.Request.Method, options.IgnoredMethods) {
			c.Next()
			return
		}

		// Apply normal idempotency logic with custom options
		m.handleWithOptions(c, options)
	}
}

// shouldIgnorePath checks if the path should be ignored
func (m *IdempotencyMiddleware) shouldIgnorePath(path string, ignoredPaths []string) bool {
	for _, ignored := range ignoredPaths {
		if path == ignored {
			return true
		}
	}
	return false
}

// shouldIgnoreMethod checks if the method should be ignored
func (m *IdempotencyMiddleware) shouldIgnoreMethod(method string, ignoredMethods []string) bool {
	for _, ignored := range ignoredMethods {
		if method == ignored {
			return true
		}
	}
	return false
}

// handleWithOptions applies idempotency with custom options
func (m *IdempotencyMiddleware) handleWithOptions(c *gin.Context, options IdempotencyOptions) {
	idempotencyKey := c.GetHeader(IdempotencyKeyHeader)

	if options.RequireKey && idempotencyKey == "" {
		c.JSON(http.StatusBadRequest, response.APIResponse{
			Success: false,
			Error: &response.ErrorInfo{
				Type:    "validation_error",
				Message: "Idempotency-Key header is required",
			},
			Timestamp: time.Now().UTC(),
		})
		c.Abort()
		return
	}

	if idempotencyKey == "" {
		c.Next()
		return
	}

	// Continue with normal idempotency logic...
	// (Similar to Handler() but with custom TTL from options)
	m.processWithIdempotency(c, idempotencyKey, options.TTL)
}

// processWithIdempotency handles the core idempotency logic
func (m *IdempotencyMiddleware) processWithIdempotency(c *gin.Context, idempotencyKey string, customTTL time.Duration) {
	ttl := m.ttl
	if customTTL != 0 {
		ttl = customTTL
	}

	messageID := m.generateMessageID(c, idempotencyKey)
	consumerID := m.generateConsumerID(c)
	eventType := m.generateEventType(c)

	ctx := c.Request.Context()

	shouldProcess, err := m.inboxService.ProcessMessageWithDeduplication(
		ctx, messageID, eventType, consumerID, ttl)

	if err != nil {
		m.logger.Error("Idempotency processing failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, response.APIResponse{
			Success: false,
			Error: &response.ErrorInfo{
				Type:    "internal_error",
				Message: "Idempotency processing failed",
			},
			Timestamp: time.Now().UTC(),
		})
		c.Abort()
		return
	}

	if !shouldProcess {
		c.JSON(http.StatusConflict, response.APIResponse{
			Success: false,
			Error: &response.ErrorInfo{
				Type:    "conflict",
				Message: "Request already processed",
			},
			Timestamp: time.Now().UTC(),
		})
		c.Abort()
		return
	}

	c.Set("idempotency_key", idempotencyKey)
	c.Set("message_id", messageID.String())
	c.Next()
}
