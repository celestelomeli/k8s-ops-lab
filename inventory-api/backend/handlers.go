// handlers.go contains the HTTP handler functions for each route
// Each handler reads the request and writes a JSON response
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Returns app health status — used by Kubernetes liveness/readiness probes.
// Pings the database so probes reflect real connectivity, not just process status.
func healthHandler(w http.ResponseWriter, r *http.Request) {
	if err := DB.Ping(); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable) // 503 — signals to Kubernetes something is wrong
		json.NewEncoder(w).Encode(map[string]string{"status": "unhealthy", "error": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// Returns all products in the inventory
func getProducts(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Query("SELECT id, name, description, quantity, price, created_at, updated_at FROM products")
	if err != nil {
		http.Error(w, "failed to query products", http.StatusInternalServerError)
		return
	}
	defer rows.Close() // close rows when done

	// iterate over rows and build a slice of products
	products := []Product{}
	for rows.Next() {
		var p Product
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Quantity, &p.Price, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			http.Error(w, "failed to scan product", http.StatusInternalServerError)
			return
		}
		products = append(products, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// Returns a single product by ID
func getProduct(w http.ResponseWriter, r *http.Request) {
	// convert id from string to int
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	// query the product by id, expecting one row
	var p Product
	err = DB.QueryRow(
		"SELECT id, name, description, quantity, price, created_at, updated_at FROM products WHERE id = $1", id,
	).Scan(&p.ID, &p.Name, &p.Description, &p.Quantity, &p.Price, &p.CreatedAt, &p.UpdatedAt)

	if err == sql.ErrNoRows {
		http.Error(w, "product not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "failed to query product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

// Inserts a new product and returns it with a 201 status
func createProduct(w http.ResponseWriter, r *http.Request) {
	var p Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// INSERT and use RETURNING to get the created row back in one query
	err := DB.QueryRow(
		`INSERT INTO products (name, description, quantity, price, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, name, description, quantity, price, created_at, updated_at`,
		p.Name, p.Description, p.Quantity, p.Price, time.Now(), time.Now(),
	).Scan(&p.ID, &p.Name, &p.Description, &p.Quantity, &p.Price, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		http.Error(w, "failed to create product", http.StatusInternalServerError)
		return
	}

	// call the notifier service by DNS name 
	go func() {
		body := fmt.Sprintf(`{"product_name":"%s"}`, p.Name)
		http.Post("http://notifier:9090/notify", "application/json", strings.NewReader(body))
	}()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}