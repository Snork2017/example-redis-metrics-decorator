package main

import (
    "context"
    "fmt"
    "time"

    "github.com/go-redis/redis/v8"
)

// StorageInterface defines an interface for interacting with the repository
type StorageInterface interface {
    Get(ctx context.Context, key string) (string, error)
    Set(ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error)
}

// RedisClient wraps the Redis client and implements StorageInterface
type RedisClient struct {
    client *redis.Client
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
    return r.client.Get(ctx, key).Result()
}

func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error) {
    return r.client.Set(ctx, key, value, expiration).Result()
}

// MetricsDecorator adds monitoring to Redis operations
type StorageMetricsDecorator struct {
    storage StorageInterface
    metrics *Metrics
}

// NewMetricsDecorator creates a new StorageMetricsDecorator instance
func NewStorageMetricsDecorator(storage StorageInterface, metrics *Metrics) *StorageMetricsDecorator {
    return &StorageMetricsDecorator{
        storage: storage,
        metrics: metrics,
    }
}

func (m *StorageMetricsDecorator) Get(ctx context.Context, key string) (string, error) {
    start := time.Now()
    result, err := m.storage.Get(ctx, key)
    m.metrics.Record("Get", time.Since(start), err)
    return result, err
}

func (m *StorageMetricsDecorator) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error) {
    start := time.Now()
    result, err := m.storage.Set(ctx, key, value, expiration)
    m.metrics.Record("Set", time.Since(start), err)
    return result, err
}

// Metrics to collect and analyze call data
type Metrics struct{}

func (m *Metrics) Record(method string, duration time.Duration, err error) {
    // Recording metrics: operation duration, number of errors, etc.
    fmt.Printf("%s took %v, error: %v\n", method, duration, err)
}

func main() {
    rdb := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })

    redisStorage := &RedisClient{client: rdb}
    metrics := &Metrics{}
    decoratedStorage := NewStorageMetricsDecorator(redisStorage, metrics)

    // Example usage
    ctx := context.Background()
    _, err := decoratedStorage.Set(ctx, "key", "value", 0)
    if err != nil {
        fmt.Println("Error setting value:", err)
    }

    val, err := decoratedStorage.Get(ctx, "key")
    if err != nil {
        fmt.Println("Error getting value:", err)
    }
    fmt.Println("Got value:", val)
}