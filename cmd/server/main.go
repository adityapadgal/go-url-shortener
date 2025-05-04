package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func main() {
	router := chi.NewRouter()
	
	router.Get("/health", func(w http.ResponseWriter, router *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

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