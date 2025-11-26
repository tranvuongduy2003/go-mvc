package messaging

import (
	"context"
)

type Publisher interface {
	Publish(ctx context.Context, subject string, data []byte) error

	PublishWithReply(ctx context.Context, subject string, data []byte, timeout context.Context) ([]byte, error)
}

type Subscriber interface {
	Subscribe(subject string, handler MessageHandler) (Subscription, error)

	QueueSubscribe(subject, queue string, handler MessageHandler) (Subscription, error)
}

type MessageBroker interface {
	Publisher
	Subscriber

	Connect() error

	Close() error

	IsConnected() bool

	Health() error
}

type BrokerMessage interface {
	Data() []byte

	Subject() string

	Reply() string

	Ack() error

	Nack() error

	Headers() map[string]string
}

type Subscription interface {
	Subject() string

	Unsubscribe() error

	IsValid() bool
}

type MessageHandler func(msg BrokerMessage) error

type Event interface {
	EventType() string

	EventData() ([]byte, error)

	AggregateID() string

	Version() int

	Timestamp() int64
}

type EventBus interface {
	PublishEvent(ctx context.Context, event Event) error

	SubscribeToEvent(eventType string, handler EventHandler) (Subscription, error)

	SubscribeToEvents(eventTypes []string, handler EventHandler) ([]Subscription, error)
}

type EventHandler func(ctx context.Context, event Event) error
