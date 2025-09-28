package messaging

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/tranvuongduy2003/go-mvc/internal/core/domain/messaging"
	messagingPorts "github.com/tranvuongduy2003/go-mvc/internal/core/ports/messaging"
)

// OutboxService handles the outbox pattern for reliable message publishing
type OutboxService struct {
	outboxRepo messagingPorts.OutboxRepository
}

// NewOutboxService creates a new outbox service
func NewOutboxService(outboxRepo messagingPorts.OutboxRepository) *OutboxService {
	return &OutboxService{
		outboxRepo: outboxRepo,
	}
}

// StoreMessage stores a message in the outbox table within the same transaction as business data
func (s *OutboxService) StoreMessage(ctx context.Context, tx interface{}, eventType, aggregateID string, payload interface{}) error {
	outboxMessage, err := messaging.NewOutboxMessage(eventType, aggregateID, payload)
	if err != nil {
		return fmt.Errorf("failed to create outbox message: %w", err)
	}

	if tx != nil {
		return s.outboxRepo.CreateWithTx(ctx, tx, outboxMessage)
	}

	return s.outboxRepo.Create(ctx, outboxMessage)
}

// StoreMessageWithID stores a message with a specific message ID in the outbox table
func (s *OutboxService) StoreMessageWithID(ctx context.Context, tx interface{}, message *messaging.OutboxMessage) error {
	if tx != nil {
		return s.outboxRepo.CreateWithTx(ctx, tx, message)
	}

	return s.outboxRepo.Create(ctx, message)
}

// GetPendingMessages retrieves pending messages from the outbox
func (s *OutboxService) GetPendingMessages(ctx context.Context, limit int) ([]*messaging.OutboxMessage, error) {
	return s.outboxRepo.GetPendingMessages(ctx, limit)
}

// MarkAsProcessed marks a message as successfully processed
func (s *OutboxService) MarkAsProcessed(ctx context.Context, messageID string) error {
	// Implementation would depend on how you want to identify the message
	// This is a simplified version
	return s.outboxRepo.UpdateStatus(ctx,
		parseUUID(messageID),
		messaging.OutboxMessageStatusProcessed)
}

// MarkAsFailed marks a message as failed
func (s *OutboxService) MarkAsFailed(ctx context.Context, messageID, errorMessage string) error {
	// Get the message first
	message, err := s.outboxRepo.GetByID(ctx, parseUUID(messageID))
	if err != nil {
		return err
	}

	if message == nil {
		return fmt.Errorf("message not found: %s", messageID)
	}

	// Mark as failed and increment retry count
	message.MarkAsFailed(errorMessage)
	message.IncrementRetry()

	return s.outboxRepo.Update(ctx, message)
}

// RetryFailedMessages gets failed messages that can be retried
func (s *OutboxService) RetryFailedMessages(ctx context.Context, limit int) ([]*messaging.OutboxMessage, error) {
	allPending, err := s.outboxRepo.GetPendingMessages(ctx, limit*2) // Get more than needed
	if err != nil {
		return nil, err
	}

	var retryableMessages []*messaging.OutboxMessage
	for _, msg := range allPending {
		if msg.ShouldRetry() && len(retryableMessages) < limit {
			retryableMessages = append(retryableMessages, msg)
		}
	}

	return retryableMessages, nil
}

// CleanupOldMessages removes old processed messages
func (s *OutboxService) CleanupOldMessages(ctx context.Context, olderThanDays int) error {
	olderThan := int64(olderThanDays * 24 * 3600) // Convert days to seconds
	return s.outboxRepo.DeleteOldProcessedMessages(ctx, olderThan)
}

// Helper function to parse UUID - in real implementation you'd handle errors properly
func parseUUID(s string) uuid.UUID {
	// This is a simplified implementation - in real code you'd handle parsing errors
	id, _ := uuid.Parse(s)
	return id
}
