package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// Set up REST API health check endpoint
func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "200 OK\n")
}

// Set up alert receiver endpoint
func handleAlertReceiver(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	fmt.Println(body)
}

// Check for ENV and return
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	http.HandleFunc("/health", handleHealth)

	port := getEnv("PORT", "8080")
	host := getEnv("HOST", "127.0.0.1")
	if host == "127.0.0.1" || host == "localhost" {
		log.Printf("Starting server on port %s", port)
	} else {
		log.Printf("Starting server on %s port %s", host, port)
	}
	log.Printf("Health endpoint: http://%s:%s/health", host, port)

	if err := http.ListenAndServe(host+":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
