package store

import (
    "context"
    "github.com/CasterlyGit/url-shortener/internal/model"
)

type URLStore interface {
    CreateURL(ctx context.Context, url *model.URL) error
    GetURLByShortCode(ctx context.Context, shortCode string) (*model.URL, error)
    IncrementClickCount(ctx context.Context, shortCode string) error
    Close() error
}