package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/tranvuongduy2003/go-mvc/internal/application/services/messaging"
	domainMessaging "github.com/tranvuongduy2003/go-mvc/internal/domain/messaging"
)

type OutboxProcessorJob struct {
	outboxService    *messaging.OutboxService
	messagePublisher domainMessaging.Publisher // Interface for publishing messages (NATS, etc.)
	batchSize        int
	retryDelay       time.Duration
}

func NewOutboxProcessorJob(
	outboxService *messaging.OutboxService,
	messagePublisher domainMessaging.Publisher,
	batchSize int,
	retryDelay time.Duration,
) *OutboxProcessorJob {
	return &OutboxProcessorJob{
		outboxService:    outboxService,
		messagePublisher: messagePublisher,
		batchSize:        batchSize,
		retryDelay:       retryDelay,
	}
}

func (j *OutboxProcessorJob) Execute(ctx context.Context) error {
	log.Printf("Starting outbox message processing...")

	messages, err := j.outboxService.GetPendingMessages(ctx, j.batchSize)
	if err != nil {
		return fmt.Errorf("failed to get pending messages: %w", err)
	}

	if len(messages) == 0 {
		log.Printf("No pending messages to process")
		return nil
	}

	log.Printf("Processing %d pending messages", len(messages))

	for _, message := range messages {
		err := j.processMessage(ctx, message)
		if err != nil {
			log.Printf("Failed to process message %s: %v", message.ID.String(), err)
		}
	}

	return nil
}

func (j *OutboxProcessorJob) processMessage(ctx context.Context, message *domainMessaging.OutboxMessage) error {
	if !message.CanRetry() && message.Status == domainMessaging.OutboxMessageStatusFailed {
		log.Printf("Message %s has exceeded max retries, skipping", message.ID.String())
		return nil
	}

	publishMessage := &domainMessaging.Message{
		ID:          message.MessageID,
		EventType:   message.EventType,
		AggregateID: message.AggregateID,
		Payload:     message.Payload,
		Timestamp:   message.CreatedAt,
		Metadata: map[string]interface{}{
			"outbox_message_id": message.ID.String(),
			"retry_count":       message.Retries,
		},
	}

	err := j.publishMessage(ctx, publishMessage)
	if err != nil {
		return j.handlePublishError(ctx, message, err)
	}

	message.MarkAsProcessed()
	err = j.outboxService.MarkAsProcessed(ctx, message.ID.String())
	if err != nil {
		log.Printf("Failed to mark message %s as processed: %v", message.ID.String(), err)
		return err
	}

	log.Printf("Successfully processed outbox message %s", message.ID.String())
	return nil
}

func (j *OutboxProcessorJob) publishMessage(ctx context.Context, message *domainMessaging.Message) error {
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to serialize message: %w", err)
	}

	topic := j.getTopicForEventType(message.EventType)

	return j.messagePublisher.Publish(ctx, topic, data)
}

func (j *OutboxProcessorJob) handlePublishError(ctx context.Context, message *domainMessaging.OutboxMessage, publishErr error) error {
	errorMessage := publishErr.Error()

	err := j.outboxService.MarkAsFailed(ctx, message.ID.String(), errorMessage)
	if err != nil {
		log.Printf("Failed to mark message %s as failed: %v", message.ID.String(), err)
		return err
	}

	log.Printf("Failed to publish message %s (retry %d/%d): %v",
		message.ID.String(), message.Retries+1, message.MaxRetries, publishErr)

	return publishErr
}

func (j *OutboxProcessorJob) getTopicForEventType(eventType string) string {
	topicMap := map[string]string{
		"user.created":    "users.events",
		"user.updated":    "users.events",
		"user.deleted":    "users.events",
		"auth.login":      "auth.events",
		"auth.logout":     "auth.events",
		"auth.registered": "auth.events",
	}

	if topic, exists := topicMap[eventType]; exists {
		return topic
	}

	return "default.events"
}

func (j *OutboxProcessorJob) ExecuteWithRetry(ctx context.Context) error {
	err := j.Execute(ctx)
	if err != nil {
		return err
	}

	return j.retryFailedMessages(ctx)
}

func (j *OutboxProcessorJob) retryFailedMessages(ctx context.Context) error {
	retryableMessages, err := j.outboxService.RetryFailedMessages(ctx, j.batchSize)
	if err != nil {
		return fmt.Errorf("failed to get retryable messages: %w", err)
	}

	if len(retryableMessages) == 0 {
		return nil
	}

	log.Printf("Retrying %d failed messages", len(retryableMessages))

	time.Sleep(j.retryDelay)

	for _, message := range retryableMessages {
		err := j.processMessage(ctx, message)
		if err != nil {
			log.Printf("Retry failed for message %s: %v", message.ID.String(), err)
		}
	}

	return nil
}

func (j *OutboxProcessorJob) CleanupOldMessages(ctx context.Context, olderThanDays int) error {
	return j.outboxService.CleanupOldMessages(ctx, olderThanDays)
}
