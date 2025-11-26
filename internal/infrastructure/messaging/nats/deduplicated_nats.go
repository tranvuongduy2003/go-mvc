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
)

type DeduplicatedNATSBroker struct {
	*NATSBroker
	inboxService *messaging.InboxService
	consumerID   string
}

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

type MessageWithDeduplication struct {
	*domainMessaging.Message
	ConsumerID string `json:"consumer_id"`
	TTL        int64  `json:"ttl_seconds"` // Time-to-live for deduplication record
}

func (d *DeduplicatedNATSBroker) SubscribeWithDeduplication(
	subject string,
	handler domainMessaging.MessageHandler,
	ttl time.Duration,
) (domainMessaging.Subscription, error) {
	deduplicatedHandler := func(msg domainMessaging.BrokerMessage) error {
		return d.handleWithDeduplication(msg, handler, ttl)
	}

	return d.Subscribe(subject, deduplicatedHandler)
}

func (d *DeduplicatedNATSBroker) SubscribeWithInbox(
	subject string,
	handler domainMessaging.MessageHandler,
) (domainMessaging.Subscription, error) {
	inboxHandler := func(msg domainMessaging.BrokerMessage) error {
		return d.handleWithInbox(msg, handler)
	}

	return d.Subscribe(subject, inboxHandler)
}

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

	message.ID = messageID

	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to serialize message: %w", err)
	}

	return d.Publish(ctx, subject, data)
}

func (d *DeduplicatedNATSBroker) handleWithDeduplication(
	msg domainMessaging.BrokerMessage,
	handler domainMessaging.MessageHandler,
	ttl time.Duration,
) error {
	var domainMsg domainMessaging.Message
	err := json.Unmarshal(msg.Data(), &domainMsg)
	if err != nil {
		d.logger.Error("Failed to unmarshal message for deduplication",
			zap.Error(err),
			zap.String("subject", msg.Subject()))
		return err
	}

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

func (d *DeduplicatedNATSBroker) handleWithInbox(
	msg domainMessaging.BrokerMessage,
	handler domainMessaging.MessageHandler,
) error {
	var domainMsg domainMessaging.Message
	err := json.Unmarshal(msg.Data(), &domainMsg)
	if err != nil {
		d.logger.Error("Failed to unmarshal message for inbox processing",
			zap.Error(err),
			zap.String("subject", msg.Subject()))
		return err
	}

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

func (d *DeduplicatedNATSBroker) GetMessageMetadata(msg domainMessaging.BrokerMessage) (*domainMessaging.Message, error) {
	var domainMsg domainMessaging.Message
	err := json.Unmarshal(msg.Data(), &domainMsg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal message metadata: %w", err)
	}
	return &domainMsg, nil
}

func (d *DeduplicatedNATSBroker) CreateIdempotentSubscription(
	subject string,
	handler domainMessaging.MessageHandler,
	options IdempotencyOptions,
) (domainMessaging.Subscription, error) {
	if options.UseInboxPattern {
		return d.SubscribeWithInbox(subject, handler)
	}

	ttl := options.DeduplicationTTL
	if ttl == 0 {
		ttl = 24 * time.Hour // Default 24 hours
	}

	return d.SubscribeWithDeduplication(subject, handler, ttl)
}

type IdempotencyOptions struct {
	UseInboxPattern  bool          `json:"use_inbox_pattern"`
	DeduplicationTTL time.Duration `json:"deduplication_ttl"`
	ConsumerID       string        `json:"consumer_id"`
}
