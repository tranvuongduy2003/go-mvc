package messaging

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/messaging"
)

type OutboxService struct {
	outboxRepo messaging.OutboxRepository
}

func NewOutboxService(outboxRepo messaging.OutboxRepository) *OutboxService {
	return &OutboxService{
		outboxRepo: outboxRepo,
	}
}

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

func (s *OutboxService) StoreMessageWithID(ctx context.Context, tx interface{}, message *messaging.OutboxMessage) error {
	if tx != nil {
		return s.outboxRepo.CreateWithTx(ctx, tx, message)
	}

	return s.outboxRepo.Create(ctx, message)
}

func (s *OutboxService) GetPendingMessages(ctx context.Context, limit int) ([]*messaging.OutboxMessage, error) {
	return s.outboxRepo.GetPendingMessages(ctx, limit)
}

func (s *OutboxService) MarkAsProcessed(ctx context.Context, messageID string) error {
	return s.outboxRepo.UpdateStatus(ctx,
		parseUUID(messageID),
		messaging.OutboxMessageStatusProcessed)
}

func (s *OutboxService) MarkAsFailed(ctx context.Context, messageID, errorMessage string) error {
	message, err := s.outboxRepo.GetByID(ctx, parseUUID(messageID))
	if err != nil {
		return err
	}

	if message == nil {
		return fmt.Errorf("message not found: %s", messageID)
	}

	message.MarkAsFailed(errorMessage)
	message.IncrementRetry()

	return s.outboxRepo.Update(ctx, message)
}

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

func (s *OutboxService) CleanupOldMessages(ctx context.Context, olderThanDays int) error {
	olderThan := int64(olderThanDays * 24 * 3600) // Convert days to seconds
	return s.outboxRepo.DeleteOldProcessedMessages(ctx, olderThan)
}

func parseUUID(s string) uuid.UUID {
	id, _ := uuid.Parse(s)
	return id
}
