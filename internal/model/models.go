package model

import (
    "time"
)

type URL struct {
    ID          int64     `json:"id" db:"id"`
    ShortCode   string    `json:"short_code" db:"short_code"`
    LongURL     string    `json:"long_url" db:"long_url"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
    ClickCount  int64     `json:"click_count" db:"click_count"`
}

type CreateURLRequest struct {
    LongURL string `json:"long_url" validate:"required,url"`
}

type CreateURLResponse struct {
    ShortURL string `json:"short_url"`
}