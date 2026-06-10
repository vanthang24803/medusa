package queue

import (
	"context"
	"crypto/rand"
	"fmt"
	"sync"
	"time"
)

// ProviderType defines the message queue provider type.
type ProviderType string

const (
	ProviderInMemory ProviderType = "inmemory"
	ProviderRabbitMQ ProviderType = "rabbitmq"
	ProviderKafka    ProviderType = "kafka"
	ProviderRedis    ProviderType = "redis"
)

// Config holds the configuration parameters for the Queue.
type Config struct {
	Provider ProviderType
	URL      string // Connection string (e.g., "amqp://...", "localhost:9092", "redis://...")
}

// Message represents the message structure sent/received in the queue.
type Message struct {
	ID        string
	Topic     string
	Payload   []byte
	Timestamp int64
}

// Handler defines the function to process a new message in the queue.
type Handler func(ctx context.Context, msg Message) error

// Queue defines the interface for communicating with message queues (RabbitMQ, Kafka, Redis, etc.).
type Queue interface {
	// Publish sends a message to the queue.
	Publish(ctx context.Context, topic string, payload []byte) error

	// Subscribe subscribes to listen and process messages from a topic/queue.
	Subscribe(ctx context.Context, topic string, handler Handler) error

	// Close closes the queue connection.
	Close() error
}

// NewQueue is a Factory function to create the corresponding Queue provider.
func NewQueue(cfg Config) (Queue, error) {
	switch cfg.Provider {
	case ProviderInMemory:
		return newInmemoryQueue(cfg)
	case ProviderRabbitMQ:
		return newRabbitMQQueue(cfg)
	case ProviderKafka:
		return newKafkaQueue(cfg)
	case ProviderRedis:
		return newRedisQueue(cfg)
	default:
		return nil, fmt.Errorf("queue: unsupported provider: %s", cfg.Provider)
	}
}

// ── InMemory Implementation (Thread-Safe) ───────────────────────────────────

type inmemoryQueue struct {
	mu       sync.RWMutex
	handlers map[string][]Handler
}

func newInmemoryQueue(_ Config) (*inmemoryQueue, error) {
	return &inmemoryQueue{
		handlers: make(map[string][]Handler),
	}, nil
}

func (q *inmemoryQueue) Publish(ctx context.Context, topic string, payload []byte) error {
	q.mu.RLock()
	handlers, ok := q.handlers[topic]
	q.mu.RUnlock()

	if !ok {
		return nil
	}

	msg := Message{
		ID:        genMessageID(),
		Topic:     topic,
		Payload:   payload,
		Timestamp: time.Now().UnixMilli(),
	}

	for _, h := range handlers {
		// Run handler concurrently in a goroutine
		go func(handler Handler) {
			_ = handler(context.Background(), msg)
		}(h)
	}

	return nil
}

func (q *inmemoryQueue) Subscribe(ctx context.Context, topic string, handler Handler) error {
	q.mu.Lock()
	q.handlers[topic] = append(q.handlers[topic], handler)
	q.mu.Unlock()
	return nil
}

func (q *inmemoryQueue) Close() error {
	return nil
}

func genMessageID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// ── RabbitMQ Implementation (Stub) ───────────────────────────────────────────

type rabbitmqQueue struct {
	cfg Config
}

func newRabbitMQQueue(cfg Config) (*rabbitmqQueue, error) {
	// TODO: Initialize actual RabbitMQ connection
	return &rabbitmqQueue{cfg: cfg}, nil
}

func (q *rabbitmqQueue) Publish(ctx context.Context, topic string, payload []byte) error {
	// TODO: Implement actual RabbitMQ publishing logic
	return nil
}

func (q *rabbitmqQueue) Subscribe(ctx context.Context, topic string, handler Handler) error {
	// TODO: Implement actual RabbitMQ subscribing logic
	return nil
}

func (q *rabbitmqQueue) Close() error {
	return nil
}

// ── Kafka Implementation (Stub) ──────────────────────────────────────────────

type kafkaQueue struct {
	cfg Config
}

func newKafkaQueue(cfg Config) (*kafkaQueue, error) {
	// TODO: Initialize actual Kafka connection
	return &kafkaQueue{cfg: cfg}, nil
}

func (q *kafkaQueue) Publish(ctx context.Context, topic string, payload []byte) error {
	// TODO: Implement actual Kafka publishing logic
	return nil
}

func (q *kafkaQueue) Subscribe(ctx context.Context, topic string, handler Handler) error {
	// TODO: Implement actual Kafka subscribing logic
	return nil
}

func (q *kafkaQueue) Close() error {
	return nil
}

// ── Redis Implementation (Stub) ──────────────────────────────────────────────

type redisQueue struct {
	cfg Config
}

func newRedisQueue(cfg Config) (*redisQueue, error) {
	// TODO: Initialize actual Redis queue connection (Pub/Sub or Streams)
	return &redisQueue{cfg: cfg}, nil
}

func (q *redisQueue) Publish(ctx context.Context, topic string, payload []byte) error {
	// TODO: Implement actual Redis queue publishing logic
	return nil
}

func (q *redisQueue) Subscribe(ctx context.Context, topic string, handler Handler) error {
	// TODO: Implement actual Redis queue subscribing logic
	return nil
}

func (q *redisQueue) Close() error {
	return nil
}
