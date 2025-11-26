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

type IdempotencyMiddleware struct {
	inboxService *messaging.InboxService
	logger       *zap.Logger
	ttl          time.Duration
}

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

func (m *IdempotencyMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !m.shouldApplyIdempotency(c.Request.Method) {
			c.Next()
			return
		}

		idempotencyKey := c.GetHeader(IdempotencyKeyHeader)
		if idempotencyKey == "" {
			c.Next()
			return
		}

		messageID := m.generateMessageID(c, idempotencyKey)
		consumerID := m.generateConsumerID(c)
		eventType := m.generateEventType(c)

		ctx := c.Request.Context()

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

func (m *IdempotencyMiddleware) shouldApplyIdempotency(method string) bool {
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return true
	default:
		return false
	}
}

func (m *IdempotencyMiddleware) generateMessageID(c *gin.Context, idempotencyKey string) uuid.UUID {
	hasher := sha256.New()
	hasher.Write([]byte(idempotencyKey))
	hasher.Write([]byte(c.Request.Method))
	hasher.Write([]byte(c.Request.URL.Path))

	hasher.Write([]byte(c.Request.URL.RawQuery))

	hash := hex.EncodeToString(hasher.Sum(nil))

	return uuid.NewSHA1(uuid.Nil, []byte(hash))
}

func (m *IdempotencyMiddleware) generateConsumerID(c *gin.Context) string {
	return "http-api"
}

func (m *IdempotencyMiddleware) generateEventType(c *gin.Context) string {
	return fmt.Sprintf("http.%s.%s",
		c.Request.Method,
		c.FullPath()) // Use route pattern, not actual path
}

type IdempotencyOptions struct {
	TTL            time.Duration
	RequireKey     bool
	CacheResponses bool
	IgnoredPaths   []string
	IgnoredMethods []string
}

func (m *IdempotencyMiddleware) WithOptions(options IdempotencyOptions) gin.HandlerFunc {
	if options.TTL != 0 {
		m.ttl = options.TTL
	}

	return func(c *gin.Context) {
		if m.shouldIgnorePath(c.Request.URL.Path, options.IgnoredPaths) {
			c.Next()
			return
		}

		if m.shouldIgnoreMethod(c.Request.Method, options.IgnoredMethods) {
			c.Next()
			return
		}

		m.handleWithOptions(c, options)
	}
}

func (m *IdempotencyMiddleware) shouldIgnorePath(path string, ignoredPaths []string) bool {
	for _, ignored := range ignoredPaths {
		if path == ignored {
			return true
		}
	}
	return false
}

func (m *IdempotencyMiddleware) shouldIgnoreMethod(method string, ignoredMethods []string) bool {
	for _, ignored := range ignoredMethods {
		if method == ignored {
			return true
		}
	}
	return false
}

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

	m.processWithIdempotency(c, idempotencyKey, options.TTL)
}

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
