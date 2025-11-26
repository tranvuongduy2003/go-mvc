package messaging

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/messaging"
)

type InboxService struct {
	inboxRepo                messaging.InboxRepository
	messageDeduplicationRepo messaging.MessageDeduplicationRepository
}

func NewInboxService(
	inboxRepo messaging.InboxRepository,
	messageDeduplicationRepo messaging.MessageDeduplicationRepository,
) *InboxService {
	return &InboxService{
		inboxRepo:                inboxRepo,
		messageDeduplicationRepo: messageDeduplicationRepo,
	}
}

func (s *InboxService) ProcessMessageWithInbox(ctx context.Context, tx interface{}, messageID uuid.UUID, eventType, consumerID string) (bool, error) {
	exists, err := s.inboxRepo.Exists(ctx, messageID, consumerID)
	if err != nil {
		return false, fmt.Errorf("failed to check if message exists: %w", err)
	}

	if exists {
		return false, nil // Message already processed, skip
	}

	inboxMessage := messaging.NewInboxMessage(messageID, eventType, consumerID)

	if tx != nil {
		err = s.inboxRepo.CreateWithTx(ctx, tx, inboxMessage)
	} else {
		err = s.inboxRepo.Create(ctx, inboxMessage)
	}

	if err != nil {
		return false, fmt.Errorf("failed to create inbox message: %w", err)
	}

	return true, nil // New message, should be processed
}

func (s *InboxService) MarkAsProcessed(ctx context.Context, messageID uuid.UUID, consumerID string) error {
	return s.inboxRepo.MarkAsProcessed(ctx, messageID, consumerID)
}

func (s *InboxService) ProcessMessageWithDeduplication(ctx context.Context, messageID uuid.UUID, eventType, consumerID string, ttl time.Duration) (bool, error) {
	created, err := s.messageDeduplicationRepo.CreateIfNotExists(ctx, messageID, consumerID, eventType, ttl)
	if err != nil {
		return false, fmt.Errorf("failed to create deduplication record: %w", err)
	}

	return created, nil // True if created (new message), false if already exists (duplicate)
}

func (s *InboxService) IsMessageProcessed(ctx context.Context, messageID uuid.UUID, consumerID string) (bool, error) {
	inboxMessage, err := s.inboxRepo.GetByMessageID(ctx, messageID, consumerID)
	if err != nil {
		return false, err
	}

	if inboxMessage == nil {
		return false, nil // Message not found, not processed
	}

	return inboxMessage.IsProcessed(), nil
}

func (s *InboxService) IsMessageDuplicate(ctx context.Context, messageID uuid.UUID, consumerID string) (bool, error) {
	return s.messageDeduplicationRepo.Exists(ctx, messageID, consumerID)
}

func (s *InboxService) CleanupOldInboxMessages(ctx context.Context, olderThanDays int) error {
	olderThan := time.Now().AddDate(0, 0, -olderThanDays).Unix()
	return s.inboxRepo.DeleteOldMessages(ctx, olderThan)
}

func (s *InboxService) CleanupExpiredDeduplicationRecords(ctx context.Context) error {
	return s.messageDeduplicationRepo.DeleteExpired(ctx)
}

func (s *InboxService) ProcessWithIdempotency(
	ctx context.Context,
	messageID uuid.UUID,
	eventType, consumerID string,
	ttl time.Duration,
	businessLogic func(ctx context.Context) error,
) error {
	shouldProcess, err := s.ProcessMessageWithDeduplication(ctx, messageID, eventType, consumerID, ttl)
	if err != nil {
		return fmt.Errorf("deduplication check failed: %w", err)
	}

	if !shouldProcess {
		return nil
	}

	return businessLogic(ctx)
}

func (s *InboxService) ProcessWithInboxPattern(
	ctx context.Context,
	tx interface{},
	messageID uuid.UUID,
	eventType, consumerID string,
	businessLogic func(ctx context.Context, tx interface{}) error,
) error {
	shouldProcess, err := s.ProcessMessageWithInbox(ctx, tx, messageID, eventType, consumerID)
	if err != nil {
		return fmt.Errorf("inbox processing failed: %w", err)
	}

	if !shouldProcess {
		return nil
	}

	err = businessLogic(ctx, tx)
	if err != nil {
		return err
	}

	return s.MarkAsProcessed(ctx, messageID, consumerID)
}
