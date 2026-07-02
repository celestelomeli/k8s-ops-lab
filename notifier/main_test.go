// Unit tests for the notifier handlers. The service is stateless (no database),
// so these need nothing running — no DB, no container, no AWS.
package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// /health should return 200 with status "ok".
func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	healthHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}
	if resp["status"] != "ok" {
		t.Errorf("expected status %q, got %q", "ok", resp["status"])
	}
}

// /notify should echo the product name in the message and report the service name.
func TestNotifyHandler(t *testing.T) {
	os.Setenv("SERVICE_NAME", "notifier-test")
	defer os.Unsetenv("SERVICE_NAME")

	req := httptest.NewRequest(http.MethodPost, "/notify", strings.NewReader(`{"product_name":"widget"}`))
	rec := httptest.NewRecorder()

	notifyHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}
	if resp["message"] != "product created: widget" {
		t.Errorf("expected message %q, got %q", "product created: widget", resp["message"])
	}
	if resp["from"] != "notifier-test" {
		t.Errorf("expected from %q, got %q", "notifier-test", resp["from"])
	}
}
