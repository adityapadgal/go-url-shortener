package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/adityapadgal/go-url-shortener/internal/shortener"
)

type ShortenRequest struct {
	URL string `json:"url"`
	TTL int    `json:"ttl_seconds"` // e.g., 3600 for 1 hour
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}

func ShortenHandler(store *shortener.InMemoryStore, baseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ShortenRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		ttl := time.Duration(req.TTL) * time.Second
		if ttl <= 0 {
			ttl = 24 * time.Hour
		}

		code, err := shortener.GenerateShortCode(6)
		if err != nil {
			http.Error(w, "error generating code", http.StatusInternalServerError)
			return
		}

		store.SaveURL(code, req.URL, ttl)

		resp := ShortenResponse{
			ShortURL: baseURL + "/" + code,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
