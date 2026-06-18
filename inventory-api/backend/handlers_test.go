// Unit tests for the HTTP handlers

package main

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// A request body that is not valid JSON should return 400 Bad Request
func TestCreateProduct_InvalidBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader("not json"))
	rec := httptest.NewRecorder()

	createProduct(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

// A non-numeric id should return 400 Bad Request.
func TestGetProduct_InvalidID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/products/abc", nil)
	req.SetPathValue("id", "abc") // mimic the router extracting id from the path
	rec := httptest.NewRecorder()

	getProduct(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

// /health should return 503 when db is unreachable (used by Kubernetes) 
func TestHealthHandler_DBDown(t *testing.T) {
	orig := DB
	defer func() { DB = orig }()

	DB, _ = sql.Open("pgx", "host=127.0.0.1 port=1 user=x password=x dbname=x sslmode=disable")

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	healthHandler(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("expected status %d, got %d", http.StatusServiceUnavailable, rec.Code)
	}
}
