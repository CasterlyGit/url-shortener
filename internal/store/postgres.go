package store

import (
    "context"
    "database/sql"
    "fmt"
    "github.com/CasterlyGit/url-shortener/internal/model"
    _ "github.com/lib/pq"
)

type PostgresStore struct {
    db *sql.DB
}

func NewPostgresStore(connStr string) (*PostgresStore, error) {
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }
    
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }
    
    return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) CreateURL(ctx context.Context, url *model.URL) error {
    query := `
        INSERT INTO urls (short_code, long_url) 
        VALUES ($1, $2) 
        RETURNING id, created_at`
    
    err := s.db.QueryRowContext(
        ctx, 
        query, 
        url.ShortCode, 
        url.LongURL,
    ).Scan(&url.ID, &url.CreatedAt)
    
    if err != nil {
        return fmt.Errorf("failed to create URL: %w", err)
    }
    
    return nil
}

func (s *PostgresStore) GetURLByShortCode(ctx context.Context, shortCode string) (*model.URL, error) {
    query := `
        SELECT id, short_code, long_url, created_at, click_count 
        FROM urls 
        WHERE short_code = $1`
    
    var url model.URL
    err := s.db.QueryRowContext(ctx, query, shortCode).Scan(
        &url.ID,
        &url.ShortCode,
        &url.LongURL,
        &url.CreatedAt,
        &url.ClickCount,
    )
    
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, fmt.Errorf("failed to get URL: %w", err)
    }
    
    return &url, nil
}

func (s *PostgresStore) IncrementClickCount(ctx context.Context, shortCode string) error {
    query := `UPDATE urls SET click_count = click_count + 1 WHERE short_code = $1`
    
    _, err := s.db.ExecContext(ctx, query, shortCode)
    if err != nil {
        return fmt.Errorf("failed to increment click count: %w", err)
    }
    
    return nil
}

func (s *PostgresStore) Close() error {
    return s.db.Close()
}