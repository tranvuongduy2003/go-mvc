package messaging

import (
	"context"
)

// Publisher defines the interface for publishing messages to a message broker
type Publisher interface {
	// Publish publishes a message to the specified subject/topic
	Publish(ctx context.Context, subject string, data []byte) error

	// PublishWithReply publishes a message and waits for a reply
	PublishWithReply(ctx context.Context, subject string, data []byte, timeout context.Context) ([]byte, error)
}

// Subscriber defines the interface for subscribing to messages from a message broker
type Subscriber interface {
	// Subscribe creates a subscription to a subject/topic with a message handler
	Subscribe(subject string, handler MessageHandler) (Subscription, error)

	// QueueSubscribe creates a queue subscription (load balanced across multiple consumers)
	QueueSubscribe(subject, queue string, handler MessageHandler) (Subscription, error)
}

// MessageBroker combines Publisher and Subscriber interfaces
type MessageBroker interface {
	Publisher
	Subscriber

	// Connect establishes connection to the message broker
	Connect() error

	// Close closes the connection to the message broker
	Close() error

	// IsConnected returns whether the connection is active
	IsConnected() bool

	// Health returns the health status of the connection
	Health() error
}

// BrokerMessage represents a message received from the broker
type BrokerMessage interface {
	// Data returns the message payload
	Data() []byte

	// Subject returns the subject/topic the message was sent to
	Subject() string

	// Reply returns the reply subject if this is a request message
	Reply() string

	// Ack acknowledges the message
	Ack() error

	// Nack negatively acknowledges the message
	Nack() error

	// Headers returns message headers (if supported)
	Headers() map[string]string
}

// Subscription represents an active subscription
type Subscription interface {
	// Subject returns the subscription subject
	Subject() string

	// Unsubscribe removes the subscription
	Unsubscribe() error

	// IsValid returns whether the subscription is still active
	IsValid() bool
}

// MessageHandler is the function signature for handling incoming messages
type MessageHandler func(msg BrokerMessage) error

// Event represents a domain event that can be published
type Event interface {
	// EventType returns the type of the event
	EventType() string

	// EventData returns the event payload as bytes
	EventData() ([]byte, error)

	// AggregateID returns the ID of the aggregate that generated this event
	AggregateID() string

	// Version returns the event version
	Version() int

	// Timestamp returns when the event occurred
	Timestamp() int64
}

// EventBus defines high-level event publishing and subscribing
type EventBus interface {
	// PublishEvent publishes a domain event
	PublishEvent(ctx context.Context, event Event) error

	// SubscribeToEvent subscribes to a specific event type
	SubscribeToEvent(eventType string, handler EventHandler) (Subscription, error)

	// SubscribeToEvents subscribes to multiple event types
	SubscribeToEvents(eventTypes []string, handler EventHandler) ([]Subscription, error)
}

// EventHandler handles domain events
type EventHandler func(ctx context.Context, event Event) error
