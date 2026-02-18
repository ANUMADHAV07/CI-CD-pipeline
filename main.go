package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type Response struct {
	Message   string `json:"message"`
	Version   string `json:"version"`
	Hostname  string `json:"hostname"`
	Timestamp string `json:"timestamp"`
}

type HealthResponse struct {
	Status string `json:"status"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	version := os.Getenv("APP_VERSION")
	if version == "" {
		version = "1.0.0"
	}

	hostname, _ := os.Hostname()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := Response{
			Message:   "Hello from Kubernetes (Updated via Rolling Update)!",
			Version:   version,
			Hostname:  hostname,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
		json.NewEncoder(w).Encode(response)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(HealthResponse{Status: "healthy"})
	})

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
