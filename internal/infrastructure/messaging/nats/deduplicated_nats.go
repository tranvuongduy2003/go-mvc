package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/tranvuongduy2003/go-mvc/internal/application/services/messaging"
	domainMessaging "github.com/tranvuongduy2003/go-mvc/internal/domain/messaging"
	messagingPorts "github.com/tranvuongduy2003/go-mvc/internal/domain/ports/messaging"
)

// DeduplicatedNATSBroker extends NATS broker with message deduplication capabilities
type DeduplicatedNATSBroker struct {
	*NATSBroker
	inboxService *messaging.InboxService
	consumerID   string
}

// NewDeduplicatedNATSBroker creates a new NATS broker with deduplication
func NewDeduplicatedNATSBroker(
	natsBroker *NATSBroker,
	inboxService *messaging.InboxService,
	consumerID string,
) *DeduplicatedNATSBroker {
	return &DeduplicatedNATSBroker{
		NATSBroker:   natsBroker,
		inboxService: inboxService,
		consumerID:   consumerID,
	}
}

// MessageWithDeduplication wraps a message with deduplication metadata
type MessageWithDeduplication struct {
	*domainMessaging.Message
	ConsumerID string `json:"consumer_id"`
	TTL        int64  `json:"ttl_seconds"` // Time-to-live for deduplication record
}

// SubscribeWithDeduplication creates a subscription with automatic deduplication
func (d *DeduplicatedNATSBroker) SubscribeWithDeduplication(
	subject string,
	handler messagingPorts.MessageHandler,
	ttl time.Duration,
) (messagingPorts.Subscription, error) {
	// Wrap the original handler with deduplication logic
	deduplicatedHandler := func(msg messagingPorts.Message) error {
		return d.handleWithDeduplication(msg, handler, ttl)
	}

	return d.Subscribe(subject, deduplicatedHandler)
}

// SubscribeWithInbox creates a subscription with full inbox pattern
func (d *DeduplicatedNATSBroker) SubscribeWithInbox(
	subject string,
	handler messagingPorts.MessageHandler,
) (messagingPorts.Subscription, error) {
	// Wrap the original handler with inbox pattern
	inboxHandler := func(msg messagingPorts.Message) error {
		return d.handleWithInbox(msg, handler)
	}

	return d.Subscribe(subject, inboxHandler)
}

// PublishWithMessageID publishes a message with a specific message ID for deduplication
func (d *DeduplicatedNATSBroker) PublishWithMessageID(
	ctx context.Context,
	subject string,
	messageID uuid.UUID,
	eventType string,
	payload interface{},
) error {
	message, err := domainMessaging.NewMessage(eventType, "", payload)
	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}

	// Override the ID with the specific message ID
	message.ID = messageID

	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to serialize message: %w", err)
	}

	return d.Publish(ctx, subject, data)
}

// handleWithDeduplication processes a message with lightweight deduplication
func (d *DeduplicatedNATSBroker) handleWithDeduplication(
	msg messagingPorts.Message,
	handler messagingPorts.MessageHandler,
	ttl time.Duration,
) error {
	// Parse message to extract ID and event type
	var domainMsg domainMessaging.Message
	err := json.Unmarshal(msg.Data(), &domainMsg)
	if err != nil {
		d.logger.Error("Failed to unmarshal message for deduplication",
			zap.Error(err),
			zap.String("subject", msg.Subject()))
		return err
	}

	// Check for duplicates and process if new
	ctx := context.Background()
	err = d.inboxService.ProcessWithIdempotency(
		ctx,
		domainMsg.ID,
		domainMsg.EventType,
		d.consumerID,
		ttl,
		func(ctx context.Context) error {
			return handler(msg)
		},
	)

	if err != nil {
		d.logger.Error("Failed to process message with deduplication",
			zap.Error(err),
			zap.String("message_id", domainMsg.ID.String()),
			zap.String("event_type", domainMsg.EventType))
		return err
	}

	d.logger.Info("Successfully processed message with deduplication",
		zap.String("message_id", domainMsg.ID.String()),
		zap.String("event_type", domainMsg.EventType),
		zap.String("subject", msg.Subject()))

	return nil
}

// handleWithInbox processes a message using the full inbox pattern
func (d *DeduplicatedNATSBroker) handleWithInbox(
	msg messagingPorts.Message,
	handler messagingPorts.MessageHandler,
) error {
	// Parse message to extract ID and event type
	var domainMsg domainMessaging.Message
	err := json.Unmarshal(msg.Data(), &domainMsg)
	if err != nil {
		d.logger.Error("Failed to unmarshal message for inbox processing",
			zap.Error(err),
			zap.String("subject", msg.Subject()))
		return err
	}

	// Process with inbox pattern (no transaction needed for message consumption)
	ctx := context.Background()
	err = d.inboxService.ProcessWithInboxPattern(
		ctx,
		nil, // No transaction for message consumption
		domainMsg.ID,
		domainMsg.EventType,
		d.consumerID,
		func(ctx context.Context, tx interface{}) error {
			return handler(msg)
		},
	)

	if err != nil {
		d.logger.Error("Failed to process message with inbox pattern",
			zap.Error(err),
			zap.String("message_id", domainMsg.ID.String()),
			zap.String("event_type", domainMsg.EventType))
		return err
	}

	d.logger.Info("Successfully processed message with inbox pattern",
		zap.String("message_id", domainMsg.ID.String()),
		zap.String("event_type", domainMsg.EventType),
		zap.String("subject", msg.Subject()))

	return nil
}

// GetMessageMetadata extracts metadata from a NATS message
func (d *DeduplicatedNATSBroker) GetMessageMetadata(msg messagingPorts.Message) (*domainMessaging.Message, error) {
	var domainMsg domainMessaging.Message
	err := json.Unmarshal(msg.Data(), &domainMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal message metadata: %w", err)
	}
	return &domainMsg, nil
}

// CreateIdempotentSubscription is a convenience method for creating idempotent subscriptions
func (d *DeduplicatedNATSBroker) CreateIdempotentSubscription(
	subject string,
	handler messagingPorts.MessageHandler,
	options IdempotencyOptions,
) (messagingPorts.Subscription, error) {
	if options.UseInboxPattern {
		return d.SubscribeWithInbox(subject, handler)
	}

	ttl := options.DeduplicationTTL
	if ttl == 0 {
		ttl = 24 * time.Hour // Default 24 hours
	}

	return d.SubscribeWithDeduplication(subject, handler, ttl)
}

// IdempotencyOptions defines options for idempotent message processing
type IdempotencyOptions struct {
	UseInboxPattern  bool          `json:"use_inbox_pattern"`
	DeduplicationTTL time.Duration `json:"deduplication_ttl"`
	ConsumerID       string        `json:"consumer_id"`
}
