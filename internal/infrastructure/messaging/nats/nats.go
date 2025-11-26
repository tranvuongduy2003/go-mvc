package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"

	"github.com/tranvuongduy2003/go-mvc/internal/domain/messaging"
	"github.com/tranvuongduy2003/go-mvc/internal/infrastructure/config"
)

// NATSBroker implements the MessageBroker interface using NATS
type NATSBroker struct {
	conn          *nats.Conn
	config        config.NATSConfig
	logger        *zap.Logger
	mu            sync.RWMutex
	subscriptions map[string]*NATSSubscription
}

// NATSMessage wraps a NATS message
type NATSMessage struct {
	msg *nats.Msg
}

// NATSSubscription wraps a NATS subscription
type NATSSubscription struct {
	sub     *nats.Subscription
	subject string
}

// NewNATSBroker creates a new NATS message broker
func NewNATSBroker(config config.NATSConfig, logger *zap.Logger) *NATSBroker {
	return &NATSBroker{
		config:        config,
		logger:        logger,
		subscriptions: make(map[string]*NATSSubscription),
	}
}

// Connect establishes connection to NATS server
func (n *NATSBroker) Connect() error {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.conn != nil && n.conn.IsConnected() {
		return nil
	}

	opts := []nats.Option{
		nats.Name("go-mvc"),
		nats.MaxReconnects(n.config.MaxReconnects),
		nats.ReconnectWait(n.config.ReconnectWait),
		nats.Timeout(n.config.Timeout),
		nats.DrainTimeout(n.config.DrainTimeout),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			if err != nil {
				n.logger.Warn("NATS disconnected", zap.Error(err))
			}
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			n.logger.Info("NATS reconnected", zap.String("url", nc.ConnectedUrl()))
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			n.logger.Info("NATS connection closed")
		}),
	}

	conn, err := nats.Connect(n.config.URL, opts...)
	if err != nil {
		return fmt.Errorf("failed to connect to NATS: %w", err)
	}

	n.conn = conn
	n.logger.Info("Connected to NATS", zap.String("url", n.config.URL))

	return nil
}

// Close closes the NATS connection
func (n *NATSBroker) Close() error {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.conn == nil {
		return nil
	}

	// Close all subscriptions
	for _, sub := range n.subscriptions {
		if err := sub.Unsubscribe(); err != nil {
			n.logger.Warn("Failed to unsubscribe", zap.String("subject", sub.subject), zap.Error(err))
		}
	}

	// Drain and close connection
	if err := n.conn.Drain(); err != nil {
		n.logger.Warn("Failed to drain NATS connection", zap.Error(err))
	}

	n.conn = nil
	n.subscriptions = make(map[string]*NATSSubscription)
	n.logger.Info("NATS connection closed")

	return nil
}

// IsConnected returns whether the connection is active
func (n *NATSBroker) IsConnected() bool {
	n.mu.RLock()
	defer n.mu.RUnlock()

	return n.conn != nil && n.conn.IsConnected()
}

// Health returns the health status of the connection
func (n *NATSBroker) Health() error {
	if !n.IsConnected() {
		return fmt.Errorf("NATS connection is not active")
	}

	// Try a simple ping
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := n.conn.FlushWithContext(ctx); err != nil {
		return fmt.Errorf("NATS health check failed: %w", err)
	}

	return nil
}

// Publish publishes a message to the specified subject
func (n *NATSBroker) Publish(ctx context.Context, subject string, data []byte) error {
	if !n.IsConnected() {
		return fmt.Errorf("NATS connection is not active")
	}

	if err := n.conn.Publish(subject, data); err != nil {
		n.logger.Error("Failed to publish message",
			zap.String("subject", subject),
			zap.Error(err))
		return fmt.Errorf("failed to publish to subject %s: %w", subject, err)
	}

	n.logger.Debug("Message published",
		zap.String("subject", subject),
		zap.Int("size", len(data)))

	return nil
}

// PublishWithReply publishes a message and waits for a reply
func (n *NATSBroker) PublishWithReply(ctx context.Context, subject string, data []byte, timeout context.Context) ([]byte, error) {
	if !n.IsConnected() {
		return nil, fmt.Errorf("NATS connection is not active")
	}

	msg, err := n.conn.RequestWithContext(ctx, subject, data)
	if err != nil {
		n.logger.Error("Failed to publish request message",
			zap.String("subject", subject),
			zap.Error(err))
		return nil, fmt.Errorf("failed to publish request to subject %s: %w", subject, err)
	}

	n.logger.Debug("Request message published and replied",
		zap.String("subject", subject),
		zap.Int("request_size", len(data)),
		zap.Int("reply_size", len(msg.Data)))

	return msg.Data, nil
}

// Subscribe creates a subscription to a subject
func (n *NATSBroker) Subscribe(subject string, handler messaging.MessageHandler) (messaging.Subscription, error) {
	if !n.IsConnected() {
		return nil, fmt.Errorf("NATS connection is not active")
	}

	sub, err := n.conn.Subscribe(subject, func(msg *nats.Msg) {
		natsMsg := &NATSMessage{msg: msg}
		if err := handler(natsMsg); err != nil {
			n.logger.Error("Message handler failed",
				zap.String("subject", subject),
				zap.Error(err))
		}
	})

	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to subject %s: %w", subject, err)
	}

	natsSub := &NATSSubscription{
		sub:     sub,
		subject: subject,
	}

	n.mu.Lock()
	n.subscriptions[subject] = natsSub
	n.mu.Unlock()

	n.logger.Info("Subscribed to subject", zap.String("subject", subject))

	return natsSub, nil
}

// QueueSubscribe creates a queue subscription (load balanced)
func (n *NATSBroker) QueueSubscribe(subject, queue string, handler messaging.MessageHandler) (messaging.Subscription, error) {
	if !n.IsConnected() {
		return nil, fmt.Errorf("NATS connection is not active")
	}

	sub, err := n.conn.QueueSubscribe(subject, queue, func(msg *nats.Msg) {
		natsMsg := &NATSMessage{msg: msg}
		if err := handler(natsMsg); err != nil {
			n.logger.Error("Queue message handler failed",
				zap.String("subject", subject),
				zap.String("queue", queue),
				zap.Error(err))
		}
	})

	if err != nil {
		return nil, fmt.Errorf("failed to queue subscribe to subject %s: %w", subject, err)
	}

	natsSub := &NATSSubscription{
		sub:     sub,
		subject: subject,
	}

	subscriptionKey := fmt.Sprintf("%s:%s", subject, queue)
	n.mu.Lock()
	n.subscriptions[subscriptionKey] = natsSub
	n.mu.Unlock()

	n.logger.Info("Queue subscribed to subject",
		zap.String("subject", subject),
		zap.String("queue", queue))

	return natsSub, nil
}

// NATSMessage implementation

// Data returns the message payload
func (m *NATSMessage) Data() []byte {
	return m.msg.Data
}

// Subject returns the subject the message was sent to
func (m *NATSMessage) Subject() string {
	return m.msg.Subject
}

// Reply returns the reply subject if this is a request message
func (m *NATSMessage) Reply() string {
	return m.msg.Reply
}

// Ack acknowledges the message (no-op for basic NATS, used for JetStream)
func (m *NATSMessage) Ack() error {
	// Basic NATS doesn't require explicit ack
	return nil
}

// Nack negatively acknowledges the message (no-op for basic NATS)
func (m *NATSMessage) Nack() error {
	// Basic NATS doesn't support explicit nack
	return nil
}

// Headers returns message headers (NATS 2.2+ feature)
func (m *NATSMessage) Headers() map[string]string {
	headers := make(map[string]string)
	if m.msg.Header != nil {
		for key, values := range m.msg.Header {
			if len(values) > 0 {
				headers[key] = values[0]
			}
		}
	}
	return headers
}

// NATSSubscription implementation

// Subject returns the subscription subject
func (s *NATSSubscription) Subject() string {
	return s.subject
}

// Unsubscribe removes the subscription
func (s *NATSSubscription) Unsubscribe() error {
	if s.sub == nil {
		return nil
	}

	err := s.sub.Unsubscribe()
	if err != nil {
		return fmt.Errorf("failed to unsubscribe from %s: %w", s.subject, err)
	}

	return nil
}

// IsValid returns whether the subscription is still active
func (s *NATSSubscription) IsValid() bool {
	return s.sub != nil && s.sub.IsValid()
}

// NATSEventBus implements EventBus using NATS
type NATSEventBus struct {
	broker *NATSBroker
	logger *zap.Logger
}

// NewNATSEventBus creates a new NATS event bus
func NewNATSEventBus(broker *NATSBroker, logger *zap.Logger) *NATSEventBus {
	return &NATSEventBus{
		broker: broker,
		logger: logger,
	}
}

// PublishEvent publishes a domain event
func (e *NATSEventBus) PublishEvent(ctx context.Context, event messaging.Event) error {
	subject := fmt.Sprintf("events.%s", event.EventType())

	data, err := event.EventData()
	if err != nil {
		return fmt.Errorf("failed to serialize event: %w", err)
	}

	if err := e.broker.Publish(ctx, subject, data); err != nil {
		return fmt.Errorf("failed to publish event %s: %w", event.EventType(), err)
	}

	e.logger.Info("Event published",
		zap.String("event_type", event.EventType()),
		zap.String("aggregate_id", event.AggregateID()),
		zap.Int("version", event.Version()))

	return nil
}

// SubscribeToEvent subscribes to a specific event type
func (e *NATSEventBus) SubscribeToEvent(eventType string, handler messaging.EventHandler) (messaging.Subscription, error) {
	subject := fmt.Sprintf("events.%s", eventType)

	messageHandler := func(msg messaging.BrokerMessage) error {
		// Parse the event from message data
		var eventData map[string]interface{}
		if err := json.Unmarshal(msg.Data(), &eventData); err != nil {
			return fmt.Errorf("failed to parse event data: %w", err)
		}

		// Create a basic event wrapper
		event := &BasicEvent{
			eventType:   eventType,
			data:        msg.Data(),
			aggregateID: getStringValue(eventData, "aggregate_id"),
			version:     getIntValue(eventData, "version"),
			timestamp:   getInt64Value(eventData, "timestamp"),
		}

		return handler(context.Background(), event)
	}

	return e.broker.Subscribe(subject, messageHandler)
}

// SubscribeToEvents subscribes to multiple event types
func (e *NATSEventBus) SubscribeToEvents(eventTypes []string, handler messaging.EventHandler) ([]messaging.Subscription, error) {
	subscriptions := make([]messaging.Subscription, 0, len(eventTypes))

	for _, eventType := range eventTypes {
		sub, err := e.SubscribeToEvent(eventType, handler)
		if err != nil {
			// Clean up already created subscriptions
			for _, existingSub := range subscriptions {
				existingSub.Unsubscribe()
			}
			return nil, fmt.Errorf("failed to subscribe to event %s: %w", eventType, err)
		}
		subscriptions = append(subscriptions, sub)
	}

	return subscriptions, nil
}

// BasicEvent is a simple implementation of the Event interface
type BasicEvent struct {
	eventType   string
	data        []byte
	aggregateID string
	version     int
	timestamp   int64
}

func (e *BasicEvent) EventType() string          { return e.eventType }
func (e *BasicEvent) EventData() ([]byte, error) { return e.data, nil }
func (e *BasicEvent) AggregateID() string        { return e.aggregateID }
func (e *BasicEvent) Version() int               { return e.version }
func (e *BasicEvent) Timestamp() int64           { return e.timestamp }

// Helper functions for parsing event data
func getStringValue(data map[string]interface{}, key string) string {
	if val, ok := data[key].(string); ok {
		return val
	}
	return ""
}

func getIntValue(data map[string]interface{}, key string) int {
	if val, ok := data[key].(float64); ok {
		return int(val)
	}
	return 0
}

func getInt64Value(data map[string]interface{}, key string) int64 {
	if val, ok := data[key].(float64); ok {
		return int64(val)
	}
	return time.Now().Unix()
}
