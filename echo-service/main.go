package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type response struct {
	ReceivedBody map[string]string `json:"received_headers"`
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	resp := response{
		ReceivedBody: make(map[string]string),
	}

	for key, values := range r.Header {
		for _, value := range values {
			resp.ReceivedBody[key] = value
		}
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
		port = "8080"
	}

	http.HandleFunc("/", echoHandler)

	log.Printf("Echo service listening on port %s", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Error starting echo service: %v", err)
	}
}
