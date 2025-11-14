package handler

import (
	"context"
    "encoding/json"
    "html/template"
    "net/http"
    
    "github.com/CasterlyGit/url-shortener/internal/metrics"  
    "github.com/CasterlyGit/url-shortener/internal/model"
    "github.com/CasterlyGit/url-shortener/internal/shortcode"
    "github.com/CasterlyGit/url-shortener/internal/store"
)

type Handler struct {
    store    store.URLStore
    baseURL  string
    template *template.Template
}

func NewHandler(store store.URLStore, baseURL string) (*Handler, error) {
    h := &Handler{
        store:   store,
        baseURL: baseURL,
    }
    
    // Parse templates
    tmpl, err := template.ParseGlob("web/template/*.html")
    if err != nil {
        return nil, err
    }
    h.template = tmpl
    
    return h, nil
}

// In CreateShortURL method, use Snowflake properly:
func (h *Handler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
    var req model.CreateURLRequest
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    // Generate Snowflake ID (returns int64)
    snowflakeID, err := shortcode.GenerateFromSnowflake()
    if err != nil {
        http.Error(w, "Failed to generate ID", http.StatusInternalServerError)
        return
    }
    
    // Generate short code from Snowflake ID
    shortCode := shortcode.EncodeBase62(snowflakeID)
    
    // Create URL record with Snowflake ID
    url := &model.URL{
        ID:        snowflakeID,  // int64 Snowflake ID
        ShortCode: shortCode,    // string short code
        LongURL:   req.LongURL,
    }
    
    if err := h.store.CreateURL(r.Context(), url); err != nil {
        http.Error(w, "Failed to create short URL", http.StatusInternalServerError)
        return
    }
    
    // Track metrics
    metrics.ShortURLCreated.Inc()
    metrics.DatabaseQueriesTotal.WithLabelValues("create_url").Inc()
    
    resp := model.CreateURLResponse{
        ShortURL: h.baseURL + "/" + shortCode,
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(resp)
}

func (h *Handler) RedirectToURL(w http.ResponseWriter, r *http.Request) {
    shortCode := r.URL.Path[1:] // Remove leading slash
    
    if shortCode == "" {
        http.Redirect(w, r, "/", http.StatusFound)
        return
    }
    
    url, err := h.store.GetURLByShortCode(r.Context(), shortCode)
    if err != nil {
        http.Error(w, "Database error", http.StatusInternalServerError)
        return
    }
    
    if url == nil {
        http.NotFound(w, r)
        return
    }
    
    // TRACK METRICS - MAKE SURE THESE LINES ARE PRESENT
    metrics.RedirectsTotal.Inc()
    metrics.DatabaseQueriesTotal.WithLabelValues("get_url").Inc()
    
    // Increment click count in background
    go func() {
        ctx := context.Background()
        h.store.IncrementClickCount(ctx, shortCode)
        metrics.DatabaseQueriesTotal.WithLabelValues("increment_click").Inc()
    }()
    
    http.Redirect(w, r, url.LongURL, http.StatusFound)
}

func (h *Handler) HomePage(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }
    
    h.template.ExecuteTemplate(w, "index.html", nil)
}