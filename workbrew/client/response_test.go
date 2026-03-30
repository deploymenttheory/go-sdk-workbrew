package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetResponseHeader(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Custom-Header", "test-value")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"id": "1"})
	}))
	defer server.Close()

	transport := setupTestClient(t, server.URL)
	var result map[string]string
	resp, err := transport.NewRequest(t.Context()).
		SetResult(&result).
		Get("/")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got := GetResponseHeader(resp, "Content-Type"); got != "application/json" {
		t.Errorf("GetResponseHeader(Content-Type) = %q, want %q", got, "application/json")
	}
	if got := GetResponseHeader(resp, "X-Custom-Header"); got != "test-value" {
		t.Errorf("GetResponseHeader(X-Custom-Header) = %q, want %q", got, "test-value")
	}
	if got := GetResponseHeader(resp, "Missing-Header"); got != "" {
		t.Errorf("GetResponseHeader(Missing-Header) = %q, want empty", got)
	}
	if got := GetResponseHeader(nil, "Content-Type"); got != "" {
		t.Errorf("GetResponseHeader(nil, ...) = %q, want empty", got)
	}
}

func TestGetResponseHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Custom", "value")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"id": "1"})
	}))
	defer server.Close()

	transport := setupTestClient(t, server.URL)
	var result map[string]string
	resp, err := transport.NewRequest(t.Context()).SetResult(&result).Get("/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	headers := GetResponseHeaders(resp)
	if headers == nil {
		t.Fatal("GetResponseHeaders() returned nil")
	}
	if len(headers) == 0 {
		t.Error("GetResponseHeaders() returned empty map")
	}

	// Nil response returns empty map
	empty := GetResponseHeaders(nil)
	if len(empty) != 0 {
		t.Errorf("GetResponseHeaders(nil) = %d headers, want 0", len(empty))
	}
}

func TestGetRateLimitHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Api-Quota-Limit", "500")
		w.Header().Set("X-Api-Quota-Remaining", "450")
		w.Header().Set("X-Api-Quota-Reset", "1640000000")
		w.Header().Set("Retry-After", "60")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"id": "1"})
	}))
	defer server.Close()

	transport := setupTestClient(t, server.URL)
	var result map[string]string
	resp, err := transport.NewRequest(t.Context()).SetResult(&result).Get("/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	limit, remaining, reset, retry := GetRateLimitHeaders(resp)
	if limit != "500" {
		t.Errorf("limit = %q, want %q", limit, "500")
	}
	if remaining != "450" {
		t.Errorf("remaining = %q, want %q", remaining, "450")
	}
	if reset != "1640000000" {
		t.Errorf("reset = %q, want %q", reset, "1640000000")
	}
	if retry != "60" {
		t.Errorf("retryAfter = %q, want %q", retry, "60")
	}

	// Nil response returns empty strings
	l, r, rs, ra := GetRateLimitHeaders(nil)
	if l != "" || r != "" || rs != "" || ra != "" {
		t.Error("GetRateLimitHeaders(nil) should return empty strings")
	}
}
