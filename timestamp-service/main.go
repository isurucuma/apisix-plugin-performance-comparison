package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type response struct {
	Time string `json:"time"`
}

func timestampHandler(w http.ResponseWriter, r *http.Request) {
	resp := response{
		Time: time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081" // Default port if PORT environment variable is not set
	}

	http.HandleFunc("/", timestampHandler)

	log.Printf("Timestamp service listening on port %s", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Error starting timestamp service: %v", err)
	}
}
