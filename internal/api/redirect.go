package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/adityapadgal/go-url-shortener/internal/shortener"
)

func RedirectHandler(store *shortener.InMemoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := chi.URLParam(r, "code")
		if code == "" {
			http.Error(w, "code required", http.StatusBadRequest)
			return
		}

		data, err := store.GetURLWithAnalytics(code)
		if err != nil {
			http.Error(w, "link expired or not found", http.StatusNotFound)
			return
		}

		http.Redirect(w, r, data.OriginalURL, http.StatusFound) // 302 redirect
	}
}
