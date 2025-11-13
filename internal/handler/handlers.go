package handler

import (
    "encoding/json"
    "html/template"
    "net/http"
	"context"
    
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

func (h *Handler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
    var req model.CreateURLRequest
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    // Generate short code with retry logic for collision
    var shortCode string
    var err error
    for i := 0; i < 3; i++ { // Retry up to 3 times on collision
        shortCode, err = shortcode.GenerateRandom()
        if err != nil {
            http.Error(w, "Failed to generate short code", http.StatusInternalServerError)
            return
        }
        
        // Check if short code already exists
        existing, err := h.store.GetURLByShortCode(r.Context(), shortCode)
        if err != nil {
            http.Error(w, "Database error", http.StatusInternalServerError)
            return
        }
        if existing == nil {
            break // Short code is available
        }
        // If we get here, there was a collision - try again
    }
    
    if shortCode == "" {
        http.Error(w, "Failed to generate unique short code", http.StatusInternalServerError)
        return
    }
    
    // Create URL record
    url := &model.URL{
        ShortCode: shortCode,
        LongURL:   req.LongURL,
    }
    
    if err := h.store.CreateURL(r.Context(), url); err != nil {
        http.Error(w, "Failed to create short URL", http.StatusInternalServerError)
        return
    }
    
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
    
    // Increment click count in background
    go func() {
        // Using background context since original request context might be cancelled
        ctx := context.Background()
        h.store.IncrementClickCount(ctx, shortCode)
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