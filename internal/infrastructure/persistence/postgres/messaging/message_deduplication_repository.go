package messaging

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/messaging"
	"gorm.io/gorm"
)

type gormMessageDeduplicationRepository struct {
	db *gorm.DB
}

func NewMessageDeduplicationRepository(db *gorm.DB) messaging.MessageDeduplicationRepository {
	return &gormMessageDeduplicationRepository{
		db: db,
	}
}

func (r *gormMessageDeduplicationRepository) Create(ctx context.Context, dedup *messaging.MessageDeduplication) error {
	return r.db.WithContext(ctx).Create(dedup).Error
}

func (r *gormMessageDeduplicationRepository) CreateWithTx(ctx context.Context, tx interface{}, dedup *messaging.MessageDeduplication) error {
	gormTx, ok := tx.(*gorm.DB)
	if !ok {
		return r.Create(ctx, dedup)
	}
	return gormTx.WithContext(ctx).Create(dedup).Error
}

func (r *gormMessageDeduplicationRepository) Exists(ctx context.Context, messageID uuid.UUID, consumerID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&messaging.MessageDeduplication{}).
		Where("message_id = ? AND consumer_id = ? AND expires_at > ?", messageID, consumerID, time.Now()).
		Count(&count).Error

	return count > 0, err
}

func (r *gormMessageDeduplicationRepository) GetByMessageID(ctx context.Context, messageID uuid.UUID, consumerID string) (*messaging.MessageDeduplication, error) {
	var dedup messaging.MessageDeduplication
	err := r.db.WithContext(ctx).
		Where("message_id = ? AND consumer_id = ? AND expires_at > ?", messageID, consumerID, time.Now()).
		First(&dedup).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &dedup, nil
}

func (r *gormMessageDeduplicationRepository) Delete(ctx context.Context, messageID uuid.UUID, consumerID string) error {
	return r.db.WithContext(ctx).
		Delete(&messaging.MessageDeduplication{}, "message_id = ? AND consumer_id = ?", messageID, consumerID).Error
}

func (r *gormMessageDeduplicationRepository) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Delete(&messaging.MessageDeduplication{}, "expires_at <= ?", time.Now()).Error
}

func (r *gormMessageDeduplicationRepository) CreateIfNotExists(ctx context.Context, messageID uuid.UUID, consumerID, eventType string, ttl time.Duration) (bool, error) {
	exists, err := r.Exists(ctx, messageID, consumerID)
	if err != nil {
		return false, err
	}

	if exists {
		return false, nil // Already exists
	}

	dedup := messaging.NewMessageDeduplication(messageID, consumerID, eventType, ttl)
	err = r.db.WithContext(ctx).Create(dedup).Error

	if err != nil {
		return false, nil
	}

	return true, nil // Successfully created
}
