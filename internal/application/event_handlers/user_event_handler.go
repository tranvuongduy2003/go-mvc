package event_handlers

import (
	"context"
	"encoding/json"

	"go.uber.org/zap"

	"github.com/tranvuongduy2003/go-mvc/internal/domain/messaging"
	"github.com/tranvuongduy2003/go-mvc/internal/domain/shared/events"
)

// UserEventHandler handles user domain events
type UserEventHandler struct {
	logger *zap.Logger
}

// NewUserEventHandler creates a new user event handler
func NewUserEventHandler(logger *zap.Logger) *UserEventHandler {
	return &UserEventHandler{
		logger: logger,
	}
}

// HandleUserCreated handles user created events
func (h *UserEventHandler) HandleUserCreated(ctx context.Context, event messaging.Event) error {
	h.logger.Info("Handling user created event",
		zap.String("event_type", event.EventType()),
		zap.String("aggregate_id", event.AggregateID()))

	// Parse event data
	var userCreated events.UserCreatedEvent
	data, err := event.EventData()
	if err != nil {
		h.logger.Error("Failed to get event data", zap.Error(err))
		return err
	}

	if err := json.Unmarshal(data, &userCreated); err != nil {
		h.logger.Error("Failed to unmarshal user created event", zap.Error(err))
		return err
	}

	// Handle the event - could send welcome email, create user profile, etc.
	h.logger.Info("User created successfully",
		zap.String("user_id", userCreated.UserID),
		zap.String("email", userCreated.Email),
		zap.String("full_name", userCreated.FullName))

	// TODO: Add business logic here:
	// - Send welcome email
	// - Create user profile
	// - Initialize user settings
	// - Update analytics

	return nil
}

// HandleUserUpdated handles user updated events
func (h *UserEventHandler) HandleUserUpdated(ctx context.Context, event messaging.Event) error {
	h.logger.Info("Handling user updated event",
		zap.String("event_type", event.EventType()),
		zap.String("aggregate_id", event.AggregateID()))

	// Parse event data
	var userUpdated events.UserUpdatedEvent
	data, err := event.EventData()
	if err != nil {
		h.logger.Error("Failed to get event data", zap.Error(err))
		return err
	}

	if err := json.Unmarshal(data, &userUpdated); err != nil {
		h.logger.Error("Failed to unmarshal user updated event", zap.Error(err))
		return err
	}

	// Handle the event
	h.logger.Info("User updated successfully",
		zap.String("user_id", userUpdated.UserID),
		zap.String("email", userUpdated.Email),
		zap.String("full_name", userUpdated.FullName))

	// TODO: Add business logic here:
	// - Update search index
	// - Invalidate cache
	// - Notify related services
	// - Update analytics

	return nil
}

// HandleUserAvatarUploaded handles user avatar uploaded events
func (h *UserEventHandler) HandleUserAvatarUploaded(ctx context.Context, event messaging.Event) error {
	h.logger.Info("Handling user avatar uploaded event",
		zap.String("event_type", event.EventType()),
		zap.String("aggregate_id", event.AggregateID()))

	// Parse event data
	var avatarUploaded events.UserAvatarUploadedEvent
	data, err := event.EventData()
	if err != nil {
		h.logger.Error("Failed to get event data", zap.Error(err))
		return err
	}

	if err := json.Unmarshal(data, &avatarUploaded); err != nil {
		h.logger.Error("Failed to unmarshal user avatar uploaded event", zap.Error(err))
		return err
	}

	// Handle the event
	h.logger.Info("User avatar uploaded successfully",
		zap.String("user_id", avatarUploaded.UserID),
		zap.String("avatar_url", avatarUploaded.AvatarURL),
		zap.String("file_key", avatarUploaded.FileKey))

	// TODO: Add business logic here:
	// - Generate thumbnails
	// - Update user profile
	// - Invalidate cache
	// - Update CDN cache
	// - Notify related services

	return nil
}

// SetupEventSubscriptions sets up event subscriptions for user events
func (h *UserEventHandler) SetupEventSubscriptions(eventBus messaging.EventBus) error {
	// Subscribe to user created events
	if _, err := eventBus.SubscribeToEvent("user.created", h.HandleUserCreated); err != nil {
		h.logger.Error("Failed to subscribe to user.created events", zap.Error(err))
		return err
	}

	// Subscribe to user updated events
	if _, err := eventBus.SubscribeToEvent("user.updated", h.HandleUserUpdated); err != nil {
		h.logger.Error("Failed to subscribe to user.updated events", zap.Error(err))
		return err
	}

	// Subscribe to user avatar uploaded events
	if _, err := eventBus.SubscribeToEvent("user.avatar.uploaded", h.HandleUserAvatarUploaded); err != nil {
		h.logger.Error("Failed to subscribe to user.avatar.uploaded events", zap.Error(err))
		return err
	}

	h.logger.Info("User event subscriptions setup successfully")
	return nil
}
