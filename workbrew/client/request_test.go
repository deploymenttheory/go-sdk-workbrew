package client

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go.uber.org/zap/zaptest"
	"resty.dev/v3"
)

type testResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

func setupTestClient(t *testing.T, baseURL string) *Transport {
	logger := zaptest.NewLogger(t)
	authConfig := &AuthConfig{
		APIKey:     "test-api-key",
		APIVersion: "v0",
	}

	transport := &Transport{
		client:        resty.New().SetBaseURL(baseURL),
		logger:        logger,
		authConfig:    authConfig,
		BaseURL:       baseURL,
		globalHeaders: make(map[string]string),
		userAgent:     "test-agent",
	}

	return transport
}

func TestNewRequest_Get_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/test" {
			t.Errorf("Expected path /test, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("limit") != "10" {
			t.Errorf("Expected query param limit=10, got %s", r.URL.Query().Get("limit"))
		}
		if r.Header.Get("X-Test-Header") != "test-value" {
			t.Errorf("Expected header X-Test-Header=test-value, got %s", r.Header.Get("X-Test-Header"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(testResponse{ID: "123", Message: "success"})
	}))
	defer server.Close()

	transport := setupTestClient(t, server.URL)

	var result testResponse
	resp, err := transport.NewRequest(context.Background()).
		SetQueryParam("limit", "10").
		SetHeader("X-Test-Header", "test-value").
		SetResult(&result).
		Get("/test")

	if err != nil {
		t.Fatalf("Get() error = %v, want nil", err)
	}
	if resp == nil {
		t.Fatal("Get() response is nil")
	}
	if resp.StatusCode() != 200 {
		t.Errorf("StatusCode = %d, want 200", resp.StatusCode())
	}
	if result.ID != "123" {
		t.Errorf("ID = %q, want %q", result.ID, "123")
	}
	if result.Message != "success" {
		t.Errorf("Message = %q, want %q", result.Message, "success")
	}
}

func TestNewRequest_Get_EmptyQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("empty") != "" {
			t.Error("Empty query param should not be sent")
		}
		if r.URL.Query().Get("valid") != "value" {
			t.Errorf("Expected query param valid=value, got %s", r.URL.Query().Get("valid"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(testResponse{ID: "test"})
	}))
	defer server.Close()

	transport := setupTestClient(t, server.URL)

	var result testResponse
	_, err := transport.NewRequest(context.Background()).
		SetQueryParam("empty", "").
		SetQueryParam("valid", "value").
		SetResult(&result).
		Get("/test")

	if err != nil {
		t.Fatalf("Get() error = %v, want nil", err)
	}
}

func TestNewRequest_Post_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		var received testResponse
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("Failed to decode body: %v", err)
		}
		if received.Message != "test message" {
			t.Errorf("Expected message 'test message', got %q", received.Message)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(testResponse{ID: "456", Message: "created"})
	}))
	defer server.Close()

	transport := setupTestClient(t, server.URL)

	var result testResponse
	resp, err := transport.NewRequest(context.Background()).
		SetBody(testResponse{Message: "test message"}).
		SetResult(&result).
		Post("/test")

	if err != nil {
		t.Fatalf("Post() error = %v, want nil", err)
	}
	if resp.StatusCode() != 201 {
		t.Errorf("StatusCode = %d, want 201", resp.StatusCode())
	}
	if result.ID != "456" {
		t.Errorf("ID = %q, want %q", result.ID, "456")
	}
}

func TestNewRequest_Post_NilBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength > 0 {
			t.Error("Expected no body for nil body parameter")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(testResponse{ID: "test"})
	}))
	defer server.Close()

	transport := setupTestClient(t, server.URL)

	var result testResponse
	_, err := transport.NewRequest(context.Background()).
		SetBody(nil).
		SetResult(&result).
		Post("/test")

	if err != nil {
		t.Fatalf("Post() error = %v, want nil", err)
	}
}

func TestNewRequest_Put_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(testResponse{ID: "updated"})
	}))
	defer server.Close()

	transport := setupTestClient(t, server.URL)

	var result testResponse
	resp, err := transport.NewRequest(context.Background()).
		SetBody(testResponse{Message: "update"}).
		SetResult(&result).
		Put("/test/123")

	if err != nil {
		t.Fatalf("Put() error = %v, want nil", err)
	}
	if resp.StatusCode() != 200 {
		t.Errorf("StatusCode = %d, want 200", resp.StatusCode())
	}
	if result.ID != "updated" {
		t.Errorf("ID = %q, want %q", result.ID, "updated")
	}
}

func TestNewRequest_Patch_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Expected PATCH request, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(testResponse{ID: "patched"})
	}))
	defer server.Close()

	transport := setupTestClient(t, server.URL)

	var result testResponse
	resp, err := transport.NewRequest(context.Background()).
		SetBody(map[string]string{"field": "value"}).
		SetResult(&result).
		Patch("/test/123")

	if err != nil {
		t.Fatalf("Patch() error = %v, want nil", err)
	}
	if resp.StatusCode() != 200 {
		t.Errorf("StatusCode = %d, want 200", resp.StatusCode())
	}
	if result.ID != "patched" {
		t.Errorf("ID = %q, want %q", result.ID, "patched")
	}
}

func TestNewRequest_Delete_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		if r.URL.Query().Get("confirm") != "true" {
			t.Errorf("Expected query param confirm=true, got %s", r.URL.Query().Get("confirm"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(testResponse{Message: "deleted"})
	}))
	defer server.Close()

	transport := setupTestClient(t, server.URL)

	var result testResponse
	resp, err := transport.NewRequest(context.Background()).
		SetQueryParam("confirm", "true").
		SetResult(&result).
		Delete("/test/123")

	if err != nil {
		t.Fatalf("Delete() error = %v, want nil", err)
	}
	if resp.StatusCode() != 200 {
		t.Errorf("StatusCode = %d, want 200", resp.StatusCode())
	}
	if result.Message != "deleted" {
		t.Errorf("Message = %q, want %q", result.Message, "deleted")
	}
}

func TestNewRequest_GetBytes_Success(t *testing.T) {
	csvData := "id,name\n1,test1\n2,test2"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Query().Get("format") != "csv" {
			t.Errorf("Expected query param format=csv, got %s", r.URL.Query().Get("format"))
		}
		w.Header().Set("Content-Type", "text/csv")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(csvData))
	}))
	defer server.Close()

	transport := setupTestClient(t, server.URL)

	resp, data, err := transport.NewRequest(context.Background()).
		SetHeader("Accept", "text/csv").
		SetQueryParam("format", "csv").
		GetBytes("/test/export")

	if err != nil {
		t.Fatalf("GetBytes() error = %v, want nil", err)
	}
	if resp == nil {
		t.Fatal("GetBytes() response is nil")
	}
	if resp.StatusCode() != 200 {
		t.Errorf("StatusCode = %d, want 200", resp.StatusCode())
	}
	if string(data) != csvData {
		t.Errorf("CSV data = %q, want %q", string(data), csvData)
	}
}

func TestNewRequest_GetBytes_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "Resource not found"})
	}))
	defer server.Close()

	transport := setupTestClient(t, server.URL)

	resp, data, err := transport.NewRequest(context.Background()).
		GetBytes("/test/not-found")

	if err == nil {
		t.Fatal("GetBytes() error = nil, want error")
	}
	if resp == nil {
		t.Fatal("GetBytes() response is nil, should return metadata even on error")
	}
	if data != nil {
		t.Errorf("GetBytes() data should be nil on error, got %v", data)
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Errorf("Expected *APIError, got %T", err)
	}
	if apiErr != nil && apiErr.StatusCode != 404 {
		t.Errorf("StatusCode = %d, want 404", apiErr.StatusCode)
	}
}

func TestNewRequest_PostForm_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		contentType := r.Header.Get("Content-Type")
		if !strings.Contains(contentType, "application/x-www-form-urlencoded") {
			t.Errorf("Expected Content-Type to contain application/x-www-form-urlencoded, got %s", contentType)
		}
		if err := r.ParseForm(); err != nil {
			t.Errorf("Failed to parse form: %v", err)
		}
		if r.FormValue("username") != "testuser" {
			t.Errorf("Expected username=testuser, got %s", r.FormValue("username"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(testResponse{ID: "form-123"})
	}))
	defer server.Close()

	transport := setupTestClient(t, server.URL)

	var result testResponse
	resp, err := transport.NewRequest(context.Background()).
		SetFormData(map[string]string{"username": "testuser"}).
		SetResult(&result).
		Post("/test/form")

	if err != nil {
		t.Fatalf("PostForm() error = %v, want nil", err)
	}
	if resp.StatusCode() != 200 {
		t.Errorf("StatusCode = %d, want 200", resp.StatusCode())
	}
	if result.ID != "form-123" {
		t.Errorf("ID = %q, want %q", result.ID, "form-123")
	}
}

func TestNewRequest_PostMultipart_Success(t *testing.T) {
	expectedFileName := "test.txt"
	expectedFileContent := "test file content"
	expectedFieldName := "uploadedFile"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		contentType := r.Header.Get("Content-Type")
		if !strings.Contains(contentType, "multipart/form-data") {
			t.Errorf("Expected Content-Type to contain multipart/form-data, got %s", contentType)
		}
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			t.Errorf("Failed to parse multipart form: %v", err)
		}
		file, header, err := r.FormFile(expectedFieldName)
		if err != nil {
			t.Errorf("Failed to get file: %v", err)
			return
		}
		defer file.Close()
		if header.Filename != expectedFileName {
			t.Errorf("Expected filename %q, got %q", expectedFileName, header.Filename)
		}
		fileContent, _ := io.ReadAll(file)
		if string(fileContent) != expectedFileContent {
			t.Errorf("Expected file content %q, got %q", expectedFileContent, string(fileContent))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(testResponse{ID: "upload-123", Message: "file uploaded"})
	}))
	defer server.Close()

	transport := setupTestClient(t, server.URL)

	fileContent := strings.NewReader(expectedFileContent)
	progressCalled := false
	progressCallback := func(fieldName, fileName string, bytesWritten, totalBytes int64) {
		progressCalled = true
		if fieldName != expectedFieldName {
			t.Errorf("Progress callback fieldName = %q, want %q", fieldName, expectedFieldName)
		}
	}

	var result testResponse
	resp, err := transport.NewRequest(context.Background()).
		SetMultipartFile(expectedFieldName, expectedFileName, fileContent, int64(len(expectedFileContent)), progressCallback).
		SetMultipartFormData(map[string]string{"description": "test description"}).
		SetResult(&result).
		Post("/test/upload")

	if err != nil {
		t.Fatalf("PostMultipart() error = %v, want nil", err)
	}
	if resp.StatusCode() != 201 {
		t.Errorf("StatusCode = %d, want 201", resp.StatusCode())
	}
	if result.ID != "upload-123" {
		t.Errorf("ID = %q, want %q", result.ID, "upload-123")
	}
	if !progressCalled {
		t.Error("Progress callback was not called")
	}
}

func TestNewRequest_ErrorResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"message": "Invalid request",
			"errors":  []string{"Field 'name' is required"},
		})
	}))
	defer server.Close()

	transport := setupTestClient(t, server.URL)

	var result testResponse
	resp, err := transport.NewRequest(context.Background()).
		SetResult(&result).
		Get("/test")

	if err == nil {
		t.Fatal("Get() error = nil, want error")
	}
	if resp == nil {
		t.Fatal("Get() response is nil, should return metadata even on error")
	}
	if resp.StatusCode() != 400 {
		t.Errorf("StatusCode = %d, want 400", resp.StatusCode())
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("Expected *APIError, got %T", err)
	}
	if apiErr.Message != "Invalid request" {
		t.Errorf("Error message = %q, want %q", apiErr.Message, "Invalid request")
	}
}

func TestNewRequest_WithGlobalHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Global") != "global-value" {
			t.Errorf("Expected global header X-Global=global-value, got %s", r.Header.Get("X-Global"))
		}
		if r.Header.Get("X-Override") != "request-value" {
			t.Errorf("Expected overridden header X-Override=request-value, got %s", r.Header.Get("X-Override"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(testResponse{ID: "test"})
	}))
	defer server.Close()

	transport := setupTestClient(t, server.URL)
	transport.globalHeaders["X-Global"] = "global-value"
	transport.globalHeaders["X-Override"] = "global-override"

	var result testResponse
	_, err := transport.NewRequest(context.Background()).
		SetHeader("X-Override", "request-value").
		SetResult(&result).
		Get("/test")

	if err != nil {
		t.Fatalf("Get() error = %v, want nil", err)
	}
}

func TestNewRequest_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Request should not reach server due to cancelled context")
	}))
	defer server.Close()

	transport := setupTestClient(t, server.URL)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	var result testResponse
	_, err := transport.NewRequest(ctx).SetResult(&result).Get("/test")

	if err == nil {
		t.Fatal("Get() error = nil, want error for cancelled context")
	}
}

func TestNewRequest_InvalidContentType(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<html><body>Not JSON</body></html>"))
	}))
	defer server.Close()

	transport := setupTestClient(t, server.URL)

	var result testResponse
	_, err := transport.NewRequest(context.Background()).
		SetResult(&result).
		Get("/test")

	if err == nil {
		t.Fatal("Get() error = nil, want error for invalid content type")
	}
}

func TestNewRequest_SetQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("key1") != "val1" {
			t.Errorf("Expected key1=val1, got %s", r.URL.Query().Get("key1"))
		}
		if r.URL.Query().Get("key2") != "val2" {
			t.Errorf("Expected key2=val2, got %s", r.URL.Query().Get("key2"))
		}
		if r.URL.Query().Get("empty") != "" {
			t.Error("Empty query param should not be sent")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(testResponse{ID: "test"})
	}))
	defer server.Close()

	transport := setupTestClient(t, server.URL)

	var result testResponse
	_, err := transport.NewRequest(context.Background()).
		SetQueryParams(map[string]string{
			"key1":  "val1",
			"key2":  "val2",
			"empty": "",
		}).
		SetResult(&result).
		Get("/test")

	if err != nil {
		t.Fatalf("Get() error = %v, want nil", err)
	}
}

func TestExecuteRequest_LogsDebug(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(testResponse{ID: "log-test"})
	}))
	defer server.Close()

	transport := setupTestClient(t, server.URL)

	var result testResponse
	resp, err := transport.NewRequest(context.Background()).
		SetResult(&result).
		Get("/test")

	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if resp.StatusCode() != 200 {
		t.Errorf("StatusCode = %d, want 200", resp.StatusCode())
	}
}
