package store

import (
    "context"
    "fmt"
    "log"

    "github.com/CasterlyGit/url-shortener/internal/model"
)

type CachedStore struct {
    primaryStore URLStore
    cache        *RedisCache
}

func NewCachedStore(primaryStore URLStore, cache *RedisCache) *CachedStore {
    return &CachedStore{
        primaryStore: primaryStore,
        cache:        cache,
    }
}

func (c *CachedStore) CreateURL(ctx context.Context, url *model.URL) error {
    return c.primaryStore.CreateURL(ctx, url)
}

func (c *CachedStore) GetURLByShortCode(ctx context.Context, shortCode string) (*model.URL, error) {
    // Try cache first
    if cachedURL, err := c.cache.GetURL(ctx, shortCode); err != nil {
        log.Printf("Cache error: %v", err)
        // Fall through to primary store
    } else if cachedURL != nil {
        return cachedURL, nil
    }

    // Cache miss - get from primary store
    url, err := c.primaryStore.GetURLByShortCode(ctx, shortCode)
    if err != nil {
        return nil, err
    }
    if url == nil {
        return nil, nil
    }

    // Update cache in background
    go func() {
        bgCtx := context.Background()
        if err := c.cache.SetURL(bgCtx, url); err != nil {
            log.Printf("Failed to cache URL: %v", err)
        }
    }()

    return url, nil
}

func (c *CachedStore) IncrementClickCount(ctx context.Context, shortCode string) error {
    return c.primaryStore.IncrementClickCount(ctx, shortCode)
}

func (c *CachedStore) Close() error {
    if err := c.cache.Close(); err != nil {
        return fmt.Errorf("failed to close cache: %w", err)
    }
    return c.primaryStore.Close()
}