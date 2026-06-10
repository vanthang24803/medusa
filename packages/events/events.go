package events

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"
)

// ProviderType defines the event bus provider type.
type ProviderType string

const (
	ProviderLog      ProviderType = "log"
	ProviderInMemory ProviderType = "inmemory"
	ProviderRedis    ProviderType = "redis"
	ProviderKafka    ProviderType = "kafka"
	ProviderRabbitMQ ProviderType = "rabbitmq"
)

// Config holds the configuration parameters for the EventBus.
type Config struct {
	Provider ProviderType
	URL      string // Connection string for the brokers
}

// Event topics — published when a domain event occurs.
const (
	TopicOrderCreated     = "order.created"
	TopicOrderCancelled   = "order.cancelled"
	TopicOrderFulfilled   = "order.fulfilled"
	TopicPaymentCaptured  = "payment.captured"
	TopicPaymentRefunded  = "payment.refunded"
	TopicInventoryUpdated = "inventory.updated"
	TopicCartCompleted    = "cart.completed"
	TopicCustomerCreated  = "customer.created"
	TopicReturnRequested  = "return.requested"
)

// EventBus — publish/subscribe abstraction.
type EventBus interface {
	Publish(ctx context.Context, topic string, payload any) error
	Subscribe(topic string, handler HandlerFunc) error
	Close() error
}

type HandlerFunc func(ctx context.Context, payload []byte) error

// NewEventBus is a Factory function to create the corresponding EventBus.
func NewEventBus(cfg Config, log *zap.Logger) (EventBus, error) {
	switch cfg.Provider {
	case ProviderLog:
		return NewLogBus(log), nil
	case ProviderInMemory:
		return newInmemoryEventBus(), nil
	case ProviderRedis:
		return newRedisEventBus(cfg)
	case ProviderKafka:
		return newKafkaEventBus(cfg)
	case ProviderRabbitMQ:
		return newRabbitMQEventBus(cfg)
	default:
		return nil, fmt.Errorf("events: unsupported provider: %s", cfg.Provider)
	}
}

// ── LogBus Implementation (Console Log Only) ────────────────────────────────

// LogBus — logs events to zap. Sufficient for the skeleton; replace with RedisBus in production.
type LogBus struct {
	log *zap.Logger
}

func NewLogBus(log *zap.Logger) *LogBus {
	return &LogBus{log: log}
}

func (b *LogBus) Publish(ctx context.Context, topic string, payload any) error {
	b.log.Info("event published", zap.String("topic", topic), zap.Any("payload", payload))
	return nil
}

func (b *LogBus) Subscribe(topic string, handler HandlerFunc) error {
	b.log.Info("subscribed", zap.String("topic", topic))
	return nil
}

func (b *LogBus) Close() error {
	return nil
}

// ── InMemory EventBus Implementation (Thread-Safe Pub/Sub) ───────────────────

type inmemoryEventBus struct {
	mu       sync.RWMutex
	handlers map[string][]HandlerFunc
}

func newInmemoryEventBus() *inmemoryEventBus {
	return &inmemoryEventBus{
		handlers: make(map[string][]HandlerFunc),
	}
}

func (b *inmemoryEventBus) Publish(ctx context.Context, topic string, payload any) error {
	b.mu.RLock()
	handlers, ok := b.handlers[topic]
	b.mu.RUnlock()

	if !ok {
		return nil
	}

	// In a real pub/sub, the payload is serialized (typically to JSON) before sending
	// We simulate this by encoding to JSON so handlers receive raw bytes
	var payloadBytes []byte
	if bytes, ok := payload.([]byte); ok {
		payloadBytes = bytes
	} else {
		// Mock serialization if not bytes
		payloadBytes = []byte(fmt.Sprintf("%v", payload))
	}

	for _, h := range handlers {
		go func(handler HandlerFunc) {
			_ = handler(context.Background(), payloadBytes)
		}(h)
	}

	return nil
}

func (b *inmemoryEventBus) Subscribe(topic string, handler HandlerFunc) error {
	b.mu.Lock()
	b.handlers[topic] = append(b.handlers[topic], handler)
	b.mu.Unlock()
	return nil
}

func (b *inmemoryEventBus) Close() error {
	return nil
}

// ── Redis EventBus Implementation (Stub) ─────────────────────────────────────

type redisEventBus struct {
	cfg Config
}

func newRedisEventBus(cfg Config) (*redisEventBus, error) {
	// TODO: Implement actual Redis event bus (Pub/Sub or Streams)
	return &redisEventBus{cfg: cfg}, nil
}

func (b *redisEventBus) Publish(ctx context.Context, topic string, payload any) error {
	return nil
}

func (b *redisEventBus) Subscribe(topic string, handler HandlerFunc) error {
	return nil
}

func (b *redisEventBus) Close() error {
	return nil
}

// ── Kafka EventBus Implementation (Stub) ─────────────────────────────────────

type kafkaEventBus struct {
	cfg Config
}

func newKafkaEventBus(cfg Config) (*kafkaEventBus, error) {
	// TODO: Implement actual Kafka event bus
	return &kafkaEventBus{cfg: cfg}, nil
}

func (b *kafkaEventBus) Publish(ctx context.Context, topic string, payload any) error {
	return nil
}

func (b *kafkaEventBus) Subscribe(topic string, handler HandlerFunc) error {
	return nil
}

func (b *kafkaEventBus) Close() error {
	return nil
}

// ── RabbitMQ EventBus Implementation (Stub) ──────────────────────────────────

type rabbitmqEventBus struct {
	cfg Config
}

func newRabbitMQEventBus(cfg Config) (*rabbitmqEventBus, error) {
	// TODO: Implement actual RabbitMQ event bus
	return &rabbitmqEventBus{cfg: cfg}, nil
}

func (b *rabbitmqEventBus) Publish(ctx context.Context, topic string, payload any) error {
	return nil
}

func (b *rabbitmqEventBus) Subscribe(topic string, handler HandlerFunc) error {
	return nil
}

func (b *rabbitmqEventBus) Close() error {
	return nil
}
