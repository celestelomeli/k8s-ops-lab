//go:build integration

// The build tag keeps this file out of unit tests 

package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
)

// TestMain runs before the tests to connect to Postgres 
func TestMain(m *testing.M) {
	initDB() // sets the global DB 

	_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS products (
		id          SERIAL PRIMARY KEY,
		name        VARCHAR(255) NOT NULL,
		description TEXT,
		quantity    INT NOT NULL DEFAULT 0,
		price       NUMERIC(10, 2) NOT NULL,
		created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	)`)
	if err != nil {
		panic("failed to ensure products table exists: " + err.Error())
	}

	os.Exit(m.Run())
}

// Creating a valid product should insert a row and return 201 with new record 
func TestCreateProduct_Integration(t *testing.T) {
	body := `{"name":"integration-widget","description":"created by a test","quantity":5,"price":9.99}`
	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(body))
	rec := httptest.NewRecorder()


	createProduct(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d (body: %s)", http.StatusCreated, rec.Code, rec.Body.String())
	}

	var p Product
	if err := json.NewDecoder(rec.Body).Decode(&p); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}
	if p.ID == 0 {
		t.Errorf("expected a database-generated id, got 0")
	}
	if p.Name != "integration-widget" {
		t.Errorf("expected name %q, got %q", "integration-widget", p.Name)
	}
}

// Listing products should query the real table and return 200 with a JSON array
func TestGetProducts_Integration(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	rec := httptest.NewRecorder()

	getProducts(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d (body: %s)", http.StatusOK, rec.Code, rec.Body.String())
	}

	var products []Product
	if err := json.NewDecoder(rec.Body).Decode(&products); err != nil {
		t.Fatalf("response was not a JSON array of products: %v", err)
	}
}

// /health should return 200 when the database is reachable.
func TestHealthHandler_OK_Integration(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	healthHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d (body: %s)", http.StatusOK, rec.Code, rec.Body.String())
	}
}

// Fetching an existing product should return 200 with that product.
func TestGetProduct_Found_Integration(t *testing.T) {
	// seed a product directly, then fetch it through the handler
	var id int
	err := DB.QueryRow(
		`INSERT INTO products (name, description, quantity, price, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING id`,
		"lookup-widget", "for the get test", 3, 4.50,
	).Scan(&id)
	if err != nil {
		t.Fatalf("seed insert failed: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/products/"+strconv.Itoa(id), nil)
	req.SetPathValue("id", strconv.Itoa(id))
	rec := httptest.NewRecorder()

	getProduct(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d (body: %s)", http.StatusOK, rec.Code, rec.Body.String())
	}

	var p Product
	if err := json.NewDecoder(rec.Body).Decode(&p); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}
	if p.ID != id || p.Name != "lookup-widget" {
		t.Errorf("got id=%d name=%q, want id=%d name=%q", p.ID, p.Name, id, "lookup-widget")
	}
}

// Fetching an id that does not exist should return 404.
func TestGetProduct_NotFound_Integration(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/products/999999999", nil)
	req.SetPathValue("id", "999999999")
	rec := httptest.NewRecorder()

	getProduct(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d (body: %s)", http.StatusNotFound, rec.Code, rec.Body.String())
	}
}
