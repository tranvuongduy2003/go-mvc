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

type NATSBroker struct {
	conn          *nats.Conn
	config        config.NATSConfig
	logger        *zap.Logger
	mu            sync.RWMutex
	subscriptions map[string]*NATSSubscription
}

type NATSMessage struct {
	msg *nats.Msg
}

type NATSSubscription struct {
	sub     *nats.Subscription
	subject string
}

func NewNATSBroker(config config.NATSConfig, logger *zap.Logger) *NATSBroker {
	return &NATSBroker{
		config:        config,
		logger:        logger,
		subscriptions: make(map[string]*NATSSubscription),
	}
}

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

func (n *NATSBroker) Close() error {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.conn == nil {
		return nil
	}

	for _, sub := range n.subscriptions {
		if err := sub.Unsubscribe(); err != nil {
			n.logger.Warn("Failed to unsubscribe", zap.String("subject", sub.subject), zap.Error(err))
		}
	}

	if err := n.conn.Drain(); err != nil {
		n.logger.Warn("Failed to drain NATS connection", zap.Error(err))
	}

	n.conn = nil
	n.subscriptions = make(map[string]*NATSSubscription)
	n.logger.Info("NATS connection closed")

	return nil
}

func (n *NATSBroker) IsConnected() bool {
	n.mu.RLock()
	defer n.mu.RUnlock()

	return n.conn != nil && n.conn.IsConnected()
}

func (n *NATSBroker) Health() error {
	if !n.IsConnected() {
		return fmt.Errorf("NATS connection is not active")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := n.conn.FlushWithContext(ctx); err != nil {
		return fmt.Errorf("NATS health check failed: %w", err)
	}

	return nil
}

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

func (m *NATSMessage) Data() []byte {
	return m.msg.Data
}

func (m *NATSMessage) Subject() string {
	return m.msg.Subject
}

func (m *NATSMessage) Reply() string {
	return m.msg.Reply
}

func (m *NATSMessage) Ack() error {
	return nil
}

func (m *NATSMessage) Nack() error {
	return nil
}

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

func (s *NATSSubscription) Subject() string {
	return s.subject
}

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

func (s *NATSSubscription) IsValid() bool {
	return s.sub != nil && s.sub.IsValid()
}

type NATSEventBus struct {
	broker *NATSBroker
	logger *zap.Logger
}

func NewNATSEventBus(broker *NATSBroker, logger *zap.Logger) *NATSEventBus {
	return &NATSEventBus{
		broker: broker,
		logger: logger,
	}
}

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

func (e *NATSEventBus) SubscribeToEvent(eventType string, handler messaging.EventHandler) (messaging.Subscription, error) {
	subject := fmt.Sprintf("events.%s", eventType)

	messageHandler := func(msg messaging.BrokerMessage) error {
		var eventData map[string]interface{}
		if err := json.Unmarshal(msg.Data(), &eventData); err != nil {
			return fmt.Errorf("failed to parse event data: %w", err)
		}

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

func (e *NATSEventBus) SubscribeToEvents(eventTypes []string, handler messaging.EventHandler) ([]messaging.Subscription, error) {
	subscriptions := make([]messaging.Subscription, 0, len(eventTypes))

	for _, eventType := range eventTypes {
		sub, err := e.SubscribeToEvent(eventType, handler)
		if err != nil {
			for _, existingSub := range subscriptions {
				existingSub.Unsubscribe()
			}
			return nil, fmt.Errorf("failed to subscribe to event %s: %w", eventType, err)
		}
		subscriptions = append(subscriptions, sub)
	}

	return subscriptions, nil
}

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
