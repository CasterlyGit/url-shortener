package main

import (
    "log"
    "net/http"
    "os"
    
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    
    "github.com/CasterlyGit/url-shortener/internal/handler"
	"github.com/CasterlyGit/url-shortener/internal/store"
)

func main() {
    // Get configuration from environment variables
    dbConnStr := getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/url_shortener?sslmode=disable")
    port := getEnv("PORT", "8080")
    baseURL := getEnv("BASE_URL", "http://localhost:8080")
    
    // Initialize database store
    store, err := store.NewPostgresStore(dbConnStr)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer store.Close()
    
    // Initialize handler
    handler, err := handler.NewHandler(store, baseURL)
    if err != nil {
        log.Fatalf("Failed to create handler: %v", err)
    }
    
    // Setup router
    r := chi.NewRouter()
    
    // Middleware
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    r.Use(middleware.RealIP)
    
    // Routes
    r.Get("/", handler.HomePage)
    r.Post("/api/shorten", handler.CreateShortURL)
    r.Get("/{shortCode}", handler.RedirectToURL)
    
    // Serve static files
    r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
    
    log.Printf("Server starting on :%s", port)
    log.Fatal(http.ListenAndServe(":"+port, r))
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}