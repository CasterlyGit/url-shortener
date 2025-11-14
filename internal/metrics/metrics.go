package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    // HTTP metrics
    HttpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
        Name: "http_requests_total",
        Help: "Total number of HTTP requests",
    }, []string{"service", "method", "endpoint", "status_code"})

    HttpRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
        Name:    "http_request_duration_seconds",
        Help:    "HTTP request duration in seconds",
        Buckets: prometheus.DefBuckets,
    }, []string{"service", "method", "endpoint"})

    // URL shortener specific metrics
    ShortURLCreated = promauto.NewCounter(prometheus.CounterOpts{
        Name: "short_url_created_total",
        Help: "Total number of short URLs created",
    })

    RedirectsTotal = promauto.NewCounter(prometheus.CounterOpts{
        Name: "redirects_total", 
        Help: "Total number of URL redirects",
    })

    // Cache metrics
    CacheHits = promauto.NewCounter(prometheus.CounterOpts{
        Name: "cache_hits_total",
        Help: "Total number of cache hits",
    })

    CacheMisses = promauto.NewCounter(prometheus.CounterOpts{
        Name: "cache_misses_total",
        Help: "Total number of cache misses",
    })

    // Database metrics
    DatabaseQueriesTotal = promauto.NewCounterVec(prometheus.CounterOpts{
        Name: "database_queries_total",
        Help: "Total number of database queries",
    }, []string{"operation"})
)