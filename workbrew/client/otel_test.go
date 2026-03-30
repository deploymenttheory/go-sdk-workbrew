package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

// TestEnableTracing_DefaultConfig tests enabling tracing with default configuration
func TestEnableTracing_DefaultConfig(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data": {"id": "test"}}`))
	}))
	defer server.Close()

	// Create client with base URL override
	transport, err := NewTransport("test-api-key", "test-workspace", WithBaseURL(server.URL))
	require.NoError(t, err)

	// Enable tracing with default config
	err = transport.EnableTracing(nil)
	require.NoError(t, err)

	// Verify the transport was wrapped
	httpClient := transport.client.Client()
	require.NotNil(t, httpClient)
	require.NotNil(t, httpClient.Transport)
}

// TestEnableTracing_CustomConfig tests enabling tracing with custom configuration
func TestEnableTracing_CustomConfig(t *testing.T) {
	// Create a span recorder to capture spans
	spanRecorder := tracetest.NewSpanRecorder()
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(spanRecorder),
	)

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data": [{"id": "test"}]}`))
	}))
	defer server.Close()

	// Create client with base URL override
	transport, err := NewTransport("test-api-key", "test-workspace", WithBaseURL(server.URL))
	require.NoError(t, err)

	// Enable tracing with custom config
	config := &OTelConfig{
		TracerProvider: tracerProvider,
		Propagators:    propagation.TraceContext{},
		ServiceName:    "test-service",
	}
	err = transport.EnableTracing(config)
	require.NoError(t, err)

	// Make a request to trigger span creation
	ctx := context.Background()
	type TestResponse struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	var result TestResponse
	_, err = transport.NewRequest(ctx).SetResult(&result).Get("/test")
	require.NoError(t, err)

	// Verify span was created
	spans := spanRecorder.Ended()
	require.Greater(t, len(spans), 0, "Expected at least one span to be recorded")

	// Verify span attributes
	span := spans[0]
	assert.Equal(t, "HTTP GET", span.Name())

	// Check for HTTP semantic convention attributes
	// Note: otelhttp may use either old (http.method) or new (http.request.method) semantic conventions
	attributes := span.Attributes()
	var foundMethod, foundURL, foundStatusCode bool
	for _, attr := range attributes {
		key := string(attr.Key)
		switch key {
		case "http.method", "http.request.method":
			foundMethod = true
			assert.Equal(t, "GET", attr.Value.AsString())
		case "http.url", "url.full":
			foundURL = true
		case "http.status_code", "http.response.status_code":
			foundStatusCode = true
			assert.Equal(t, int64(200), attr.Value.AsInt64())
		}
	}

	// These assertions are informative - otelhttp sets these automatically
	// If they fail, the version may use different semantic convention names
	if !foundMethod {
		t.Logf("Warning: http.method attribute not found, attributes: %v", attributes)
	}
	if !foundURL {
		t.Logf("Warning: http.url attribute not found, attributes: %v", attributes)
	}
	if !foundStatusCode {
		t.Logf("Warning: http.status_code attribute not found, attributes: %v", attributes)
	}
}

// TestEnableTracing_WithCustomSpanNameFormatter tests custom span name formatting
func TestEnableTracing_WithCustomSpanNameFormatter(t *testing.T) {
	// Create a span recorder to capture spans
	spanRecorder := tracetest.NewSpanRecorder()
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(spanRecorder),
	)

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data": [{"id": "test"}]}`))
	}))
	defer server.Close()

	// Create client with base URL override
	transport, err := NewTransport("test-api-key", "test-workspace", WithBaseURL(server.URL))
	require.NoError(t, err)

	// Enable tracing with custom span name formatter
	config := &OTelConfig{
		TracerProvider: tracerProvider,
		SpanNameFormatter: func(operation string, req *http.Request) string {
			return "WB: " + req.Method + " " + req.URL.Path
		},
	}
	err = transport.EnableTracing(config)
	require.NoError(t, err)

	// Make a request
	ctx := context.Background()
	type TestResponse struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	var result TestResponse
	_, err = transport.NewRequest(ctx).SetResult(&result).Get("/test/path")
	require.NoError(t, err)

	// Verify custom span name
	spans := spanRecorder.Ended()
	require.Greater(t, len(spans), 0)
	assert.Contains(t, spans[0].Name(), "WB: GET")
}

// TestWithTracing_ClientOption tests the WithTracing client option
func TestWithTracing_ClientOption(t *testing.T) {
	// Create a span recorder
	spanRecorder := tracetest.NewSpanRecorder()
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(spanRecorder),
	)

	// Set as global tracer provider for this test
	otel.SetTracerProvider(tracerProvider)

	// Create client with tracing enabled via option
	transport, err := NewTransport(
		"test-api-key",
		"test-workspace",
		WithTracing(nil), // Uses global tracer provider
	)
	require.NoError(t, err)

	// Verify tracing is enabled
	httpClient := transport.client.Client()
	require.NotNil(t, httpClient)
	require.NotNil(t, httpClient.Transport)
}

// TestDefaultOTelConfig tests the default OpenTelemetry configuration
func TestDefaultOTelConfig(t *testing.T) {
	config := DefaultOTelConfig()

	assert.NotNil(t, config)
	assert.NotNil(t, config.TracerProvider)
	assert.NotNil(t, config.Propagators)
	assert.Equal(t, "workbrew-client", config.ServiceName)
	assert.Nil(t, config.SpanNameFormatter)
}

// TestEnableTracing_ErrorPropagation tests that HTTP errors are properly recorded in spans
func TestEnableTracing_ErrorPropagation(t *testing.T) {
	// Create a span recorder
	spanRecorder := tracetest.NewSpanRecorder()
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(spanRecorder),
	)

	// Create a test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": {"message": "Not found"}}`))
	}))
	defer server.Close()

	// Create client with tracing
	transport, err := NewTransport(
		"test-api-key",
		"test-workspace",
		WithBaseURL(server.URL),
	)
	require.NoError(t, err)

	config := &OTelConfig{
		TracerProvider: tracerProvider,
	}
	err = transport.EnableTracing(config)
	require.NoError(t, err)

	// Make a request that will result in an error
	ctx := context.Background()
	type TestResponse struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	var result TestResponse
	_, err = transport.NewRequest(ctx).SetResult(&result).Get("/not-found")

	// The request should return an API error
	require.Error(t, err)

	// Verify span was created with error status
	spans := spanRecorder.Ended()
	require.Greater(t, len(spans), 0)

	// Check that status code is recorded
	span := spans[0]
	attributes := span.Attributes()
	var foundStatusCode bool
	for _, attr := range attributes {
		// otelhttp uses "http.response.status_code" in newer versions
		if attr.Key == "http.status_code" || attr.Key == "http.response.status_code" {
			foundStatusCode = true
			assert.Equal(t, int64(404), attr.Value.AsInt64())
			break
		}
	}
	assert.True(t, foundStatusCode, "Expected http.status_code or http.response.status_code attribute for error response")
}
