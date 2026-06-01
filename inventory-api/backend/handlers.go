// handlers.go contains the HTTP handler functions for each route
// Each handler reads the request and writes a JSON response
package main

import (
	"encoding/json"
	"net/http"
	"time"
)

// mock data 
var mockProducts = []Product{
	{ID: 1, Name: "Widget", Description: "A small widget", Quantity: 100, Price: 9.99, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	{ID: 2, Name: "Gadget", Description: "A useful gadget", Quantity: 50, Price: 24.99, CreatedAt: time.Now(), UpdatedAt: time.Now()},
}

// healthHandler returns a simple JSON status, used by Kubernetes liveness/readiness probes
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// getProducts returns all products in the inventory
func getProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mockProducts)
}

// getProduct returns a single product by ID
func getProduct(w http.ResponseWriter, r *http.Request) {
	// PathValue extracts {id} from the URL 
	w.Header().Set("Content-Type", "application/json")
	// ignoring id for now, always returns first product
	json.NewEncoder(w).Encode(map[string]any{"id": id, "product": mockProducts[0]})
}

// decodes a JSON request body and returns the created product with a 201 status
func createProduct(w http.ResponseWriter, r *http.Request) {
	var p Product
	// decode the JSON request body into a Product struct
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	// mock: assign a fake ID and timestamps
	p.ID = len(mockProducts) + 1
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201 Created
	json.NewEncoder(w).Encode(p)
}
