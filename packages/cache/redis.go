package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrCacheMiss = errors.New("cache: key not found")

// ProviderType defines the cache provider type.
type ProviderType string

const (
	ProviderRedis    ProviderType = "redis"
	ProviderInMemory ProviderType = "inmemory"
)

// Config holds configuration parameters for the cache.
type Config struct {
	Provider ProviderType
	URL      string // Connection URL (e.g., "redis://localhost:6379")
}

// Cache defines the interface for the cache system.
type Cache interface {
	Get(ctx context.Context, key string, dest any) error
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	Close() error
}

// NewCache is a Factory function to create the corresponding Cache provider.
func NewCache(cfg Config) (Cache, error) {
	switch cfg.Provider {
	case ProviderRedis:
		return NewRedisCache(cfg.URL)
	case ProviderInMemory:
		return NewInMemoryCache(), nil
	default:
		return nil, fmt.Errorf("cache: unsupported provider: %s", cfg.Provider)
	}
}

// ── Redis Cache Implementation ──────────────────────────────────────────────

// RedisCache implements Cache interface using Redis.
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache initializes the connection to Redis.
func NewRedisCache(url string) (*RedisCache, error) {
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opt)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisCache{client: client}, nil
}

// Get retrieves the value from cache and automatically unmarshals it into dest.
func (c *RedisCache) Get(ctx context.Context, key string, dest any) error {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrCacheMiss
		}
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

// Set stores the value in cache after marshaling it to JSON.
func (c *RedisCache) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, bytes, expiration).Err()
}

// Delete deletes a key from cache.
func (c *RedisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (c *RedisCache) Close() error {
	return c.client.Close()
}

// ── InMemory Cache Implementation (Thread-Safe with TTL) ─────────────────────

type inmemoryItem struct {
	value     []byte
	expiresAt time.Time
}

type inmemoryCache struct {
	mu   sync.RWMutex
	data map[string]inmemoryItem
}

// NewInMemoryCache creates a new in-memory cache instance.
func NewInMemoryCache() *inmemoryCache {
	return &inmemoryCache{
		data: make(map[string]inmemoryItem),
	}
}

func (c *inmemoryCache) Get(ctx context.Context, key string, dest any) error {
	c.mu.RLock()
	item, ok := c.data[key]
	c.mu.RUnlock()

	if !ok {
		return ErrCacheMiss
	}

	if !item.expiresAt.IsZero() && time.Now().After(item.expiresAt) {
		// Clean up expired item lazily
		c.mu.Lock()
		delete(c.data, key)
		c.mu.Unlock()
		return ErrCacheMiss
	}

	return json.Unmarshal(item.value, dest)
}

func (c *inmemoryCache) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	var expiresAt time.Time
	if expiration > 0 {
		expiresAt = time.Now().Add(expiration)
	}

	c.mu.Lock()
	c.data[key] = inmemoryItem{
		value:     bytes,
		expiresAt: expiresAt,
	}
	c.mu.Unlock()
	return nil
}

func (c *inmemoryCache) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	delete(c.data, key)
	c.mu.Unlock()
	return nil
}

func (c *inmemoryCache) Close() error {
	return nil
}
