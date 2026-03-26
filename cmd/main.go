package main

import (
	"log"
	"net/http"
	"time"

	"github.com/hferr/pack-api/internal/httpjson"
)

func main() {
	h := httpjson.NewHandler()

	s := http.Server{
		Addr:         ":8080",
		Handler:      h.NewRouter(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	if err := s.ListenAndServe(); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
