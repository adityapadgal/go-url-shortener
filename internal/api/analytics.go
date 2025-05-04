package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/adityapadgal/go-url-shortener/internal/shortener"
)

type AnalyticsResponse struct {
	OriginalURL string `json:"original_url"`
	CreatedAt   string `json:"created_at"`
	ExpiresAt   string `json:"expires_at"`
	AccessCount int    `json:"access_count"`
	LastAccess  string `json:"last_access"`
}

func AnalyticsHandler(store *shortener.InMemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := chi.URLParam(r, "code")
		data, err := store.GetURL(code)
		if err != nil {
			http.Error(w, "not found or expired", http.StatusNotFound)
			return
		}

		resp := AnalyticsResponse{
			OriginalURL: data.OriginalURL,
			CreatedAt:   data.CreatedAt.Format(time.RFC3339),
			ExpiresAt:   data.ExpiresAt.Format(time.RFC3339),
			AccessCount: data.AccessCount,
			LastAccess:  data.LastAccess.Format(time.RFC3339),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
