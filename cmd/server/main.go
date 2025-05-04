package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/adityapadgal/go-url-shortener/internal/api"
	"github.com/adityapadgal/go-url-shortener/internal/shortener"
	"github.com/adityapadgal/go-url-shortener/internal/rateLimiter"

)

func main() {
	router := chi.NewRouter()
	store := shortener.NewStore()

		
	limiter := rateLimiter.NewLimiter(3, time.Minute) // 3 requests/minute
	router.Use(limiter.Limit)
	
	router.Get("/health", func(w http.ResponseWriter, router *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	router.Post("/shorten", api.ShortenHandler(store, "http://localhost:8000"))
	fmt.Println(store)
	router.Get("/{code}", api.RedirectHandler(store))
	router.Get("/analytics/{code}", api.AnalyticsHandler(store))


	port := "8000"
	srv := &http.Server{
		Addr:			":" + port,
		Handler:		router,
		ReadTimeout:	10 * time.Second,
		WriteTimeout:	10 * time.Second,
	}

	fmt.Printf("Server running on port:%s\n", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}