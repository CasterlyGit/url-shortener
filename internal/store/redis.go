package store

import (
    "context"
    "encoding/json"
    "fmt"
    "time"

    "github.com/CasterlyGit/url-shortener/internal/model"
    "github.com/go-redis/redis/v8"
)

type RedisCache struct {
    client *redis.Client
    ttl    time.Duration
}

func NewRedisCache(redisURL string) (*RedisCache, error) {
    opts, err := redis.ParseURL(redisURL)
    if err != nil {
        return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
    }

    client := redis.NewClient(opts)
    
    // Test connection
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := client.Ping(ctx).Err(); err != nil {
        return nil, fmt.Errorf("failed to connect to Redis: %w", err)
    }

    return &RedisCache{
        client: client,
        ttl:    24 * time.Hour, // Cache for 24 hours
    }, nil
}

func (r *RedisCache) GetURL(ctx context.Context, shortCode string) (*model.URL, error) {
    data, err := r.client.Get(ctx, "url:"+shortCode).Result()
    if err == redis.Nil {
        return nil, nil // Cache miss
    }
    if err != nil {
        return nil, fmt.Errorf("failed to get from Redis: %w", err)
    }

    var url model.URL
    if err := json.Unmarshal([]byte(data), &url); err != nil {
        return nil, fmt.Errorf("failed to unmarshal URL: %w", err)
    }

    return &url, nil
}

func (r *RedisCache) SetURL(ctx context.Context, url *model.URL) error {
    data, err := json.Marshal(url)
    if err != nil {
        return fmt.Errorf("failed to marshal URL: %w", err)
    }

    err = r.client.Set(ctx, "url:"+url.ShortCode, data, r.ttl).Err()
    if err != nil {
        return fmt.Errorf("failed to set Redis key: %w", err)
    }

    return nil
}

func (r *RedisCache) Close() error {
    return r.client.Close()
}