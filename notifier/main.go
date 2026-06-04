// main.go — tiny stateless notifier service
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /notify", notifyHandler)
	mux.HandleFunc("GET /health", healthHandler)

	log.Println("notifier listening on :9090")
	log.Fatal(http.ListenAndServe(":9090", mux))
}

func notifyHandler(w http.ResponseWriter, r *http.Request) {
	// read the product name from the request body
	var body struct {
		ProductName string `json:"product_name"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	// in a real system this would send an email, Slack message, etc.
	message := fmt.Sprintf("product created: %s", body.ProductName)
	log.Printf("notification sent: %s", message)

	serviceName := os.Getenv("SERVICE_NAME")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": message,
		"from":    serviceName,
	})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
