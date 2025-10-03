package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// Set up REST API health check endpoint
func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Request received\n"))
	log.Println("200 OK\n")
}

// Set up alert receiver endpoint
func handleAlertReceiver(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Timestamp in ISO 8601 format
	timestamp := time.Now().Format(time.RFC3339)

	// Print request method and URL
	fmt.Printf("Received %s request for %s\n", r.Method, r.URL.Path)

	// Print Connection information
	fmt.Println("Connection Information")
	fmt.Printf("  Time: %s\n", timestamp)
	fmt.Printf("  Client IP/Port: %s\n", r.RemoteAddr)
	fmt.Printf("  Protocol: %s\n", r.Proto)
	fmt.Printf("  URI: %s\n", r.RequestURI)

	// Print headers
	fmt.Println("Headers:")
	for name, values := range r.Header {
		for _, value := range values {
			fmt.Printf("  %s: %s\n", name, value)
		}
	}

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	} else {
		fmt.Println("Body:")
		fmt.Println(string(body))
	}

	//dump, err := httputil.DumpRequest(r, true)
	//if err != nil {
	//	http.Error(w, "Failed to dump request: "+err.Error(), http.StatusInternalServerError)
	//	return
	//}

	//fmt.Println("Full HTTP Request:")
	//fmt.Println(string(dump))

	// Respond to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Request received\n"))

	defer r.Body.Close()
}

// Check for ENV and return
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	defer func() {
		log.Println("Exiting application.")
	}()

	http.HandleFunc("/health", handleHealth)
	http.HandleFunc("/receiver", handleAlertReceiver)

	port := getEnv("PORT", "8080")
	host := getEnv("HOST", "127.0.0.1")
	log.Printf("Starting server on %s port %s", host, port)
	log.Printf("Health endpoint: http://%s:%s/health", host, port)
	log.Printf("Alert Receiver endpoint: http://%s:%s/receiver", host, port)

	if err := http.ListenAndServe(host+":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
