package messaging

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/messaging"
	messagingPorts "github.com/tranvuongduy2003/go-mvc/internal/domain/ports/messaging"
)

// InboxService handles the inbox pattern for message deduplication on consumer side
type InboxService struct {
	inboxRepo                messagingPorts.InboxRepository
	messageDeduplicationRepo messagingPorts.MessageDeduplicationRepository
}

// NewInboxService creates a new inbox service
func NewInboxService(
	inboxRepo messagingPorts.InboxRepository,
	messageDeduplicationRepo messagingPorts.MessageDeduplicationRepository,
) *InboxService {
	return &InboxService{
		inboxRepo:                inboxRepo,
		messageDeduplicationRepo: messageDeduplicationRepo,
	}
}

// ProcessMessageWithInbox processes a message using the full inbox pattern
// Returns true if message should be processed, false if it's a duplicate
func (s *InboxService) ProcessMessageWithInbox(ctx context.Context, tx interface{}, messageID uuid.UUID, eventType, consumerID string) (bool, error) {
	// Check if message already exists
	exists, err := s.inboxRepo.Exists(ctx, messageID, consumerID)
	if err != nil {
		return false, fmt.Errorf("failed to check if message exists: %w", err)
	}

	if exists {
		return false, nil // Message already processed, skip
	}

	// Create inbox message record
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

// MarkAsProcessed marks a message as successfully processed in the inbox
func (s *InboxService) MarkAsProcessed(ctx context.Context, messageID uuid.UUID, consumerID string) error {
	return s.inboxRepo.MarkAsProcessed(ctx, messageID, consumerID)
}

// ProcessMessageWithDeduplication processes a message using lightweight deduplication
// Returns true if message should be processed, false if it's a duplicate
func (s *InboxService) ProcessMessageWithDeduplication(ctx context.Context, messageID uuid.UUID, eventType, consumerID string, ttl time.Duration) (bool, error) {
	// Try to create deduplication record
	created, err := s.messageDeduplicationRepo.CreateIfNotExists(ctx, messageID, consumerID, eventType, ttl)
	if err != nil {
		return false, fmt.Errorf("failed to create deduplication record: %w", err)
	}

	return created, nil // True if created (new message), false if already exists (duplicate)
}

// IsMessageProcessed checks if a message has been processed using inbox pattern
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

// IsMessageDuplicate checks if a message is a duplicate using deduplication table
func (s *InboxService) IsMessageDuplicate(ctx context.Context, messageID uuid.UUID, consumerID string) (bool, error) {
	return s.messageDeduplicationRepo.Exists(ctx, messageID, consumerID)
}

// CleanupOldInboxMessages removes old inbox messages
func (s *InboxService) CleanupOldInboxMessages(ctx context.Context, olderThanDays int) error {
	olderThan := time.Now().AddDate(0, 0, -olderThanDays).Unix()
	return s.inboxRepo.DeleteOldMessages(ctx, olderThan)
}

// CleanupExpiredDeduplicationRecords removes expired deduplication records
func (s *InboxService) CleanupExpiredDeduplicationRecords(ctx context.Context) error {
	return s.messageDeduplicationRepo.DeleteExpired(ctx)
}

// ProcessWithIdempotency combines deduplication check with business logic execution
// The businessLogic function will only be called if the message is not a duplicate
func (s *InboxService) ProcessWithIdempotency(
	ctx context.Context,
	messageID uuid.UUID,
	eventType, consumerID string,
	ttl time.Duration,
	businessLogic func(ctx context.Context) error,
) error {
	// Check for duplicate and create deduplication record atomically
	shouldProcess, err := s.ProcessMessageWithDeduplication(ctx, messageID, eventType, consumerID, ttl)
	if err != nil {
		return fmt.Errorf("deduplication check failed: %w", err)
	}

	if !shouldProcess {
		// Message is a duplicate, skip processing
		return nil
	}

	// Execute business logic
	return businessLogic(ctx)
}

// ProcessWithInboxPattern combines inbox pattern with business logic execution
func (s *InboxService) ProcessWithInboxPattern(
	ctx context.Context,
	tx interface{},
	messageID uuid.UUID,
	eventType, consumerID string,
	businessLogic func(ctx context.Context, tx interface{}) error,
) error {
	// Check for duplicate and create inbox record
	shouldProcess, err := s.ProcessMessageWithInbox(ctx, tx, messageID, eventType, consumerID)
	if err != nil {
		return fmt.Errorf("inbox processing failed: %w", err)
	}

	if !shouldProcess {
		// Message is a duplicate, skip processing
		return nil
	}

	// Execute business logic
	err = businessLogic(ctx, tx)
	if err != nil {
		return err
	}

	// Mark as processed
	return s.MarkAsProcessed(ctx, messageID, consumerID)
}
