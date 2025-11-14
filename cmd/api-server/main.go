package main

import (
    "log"
    "net/http"
    "os"
    "time"
    
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/go-chi/httprate"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    
    "github.com/CasterlyGit/url-shortener/internal/handler"
    "github.com/CasterlyGit/url-shortener/internal/metrics"
    "github.com/CasterlyGit/url-shortener/internal/shortcode"
    "github.com/CasterlyGit/url-shortener/internal/store"
)

// Metrics middleware
func metricsMiddleware(serviceName string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            
            // Create a custom ResponseWriter to capture status code
            rw := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
            
            next.ServeHTTP(rw, r)
            
            duration := time.Since(start).Seconds()
            statusCode := http.StatusText(rw.Status())
            
            // Record metrics
            metrics.HttpRequestsTotal.WithLabelValues(
                serviceName,
                r.Method,
                r.URL.Path,
                statusCode,
            ).Inc()
            
            metrics.HttpRequestDuration.WithLabelValues(
                serviceName,
                r.Method, 
                r.URL.Path,
            ).Observe(duration)
        })
    }
}

func main() {
    // Get configuration from environment variables
    dbConnStr := getEnv("DATABASE_URL", "postgres://user:password@db:5432/url_shortener?sslmode=disable")
    redisURL := getEnv("REDIS_URL", "redis://redis:6379")
    port := getEnv("PORT", "8080")  // API service port
    redirectBaseURL := getEnv("REDIRECT_BASE_URL", "http://localhost:8081")  // Points to redirect service
    
    // Initialize stores
    dbStore, err := store.NewPostgresStore(dbConnStr)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer dbStore.Close()

    redisCache, err := store.NewRedisCache(redisURL)
    if err != nil {
        log.Fatalf("Failed to connect to Redis: %v", err)
    }
    defer redisCache.Close()

    cachedStore := store.NewCachedStore(dbStore, redisCache)

    // INITIALIZE SNOWFLAKE
    if err := shortcode.InitSnowflake(1); err != nil { // Node ID 1
        log.Fatalf("Failed to initialize Snowflake: %v", err)
    }

    // Initialize API handler
    handler, err := handler.NewHandler(cachedStore, redirectBaseURL)
    if err != nil {
        log.Fatalf("Failed to create handler: %v", err)
    }
    
    // Setup router
    r := chi.NewRouter()
    
      // Middleware - ALL middleware must be defined BEFORE routes
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    r.Use(middleware.RealIP)
    r.Use(metricsMiddleware("api-service")) 
    r.Use(httprate.Limit(
        1000, // requests
        1 * time.Minute, // per minute
        httprate.WithKeyFuncs(httprate.KeyByEndpoint),
    ))
    
    // Prometheus metrics endpoint
    r.Get("/metrics", promhttp.Handler().ServeHTTP)
    
    // API routes only - AFTER all middleware
    r.Get("/", handler.HomePage)
    r.Post("/api/shorten", handler.CreateShortURL)
    
    // Serve static files
    r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
    
    log.Printf("API Service starting on :%s", port)
    log.Fatal(http.ListenAndServe(":"+port, r))
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}