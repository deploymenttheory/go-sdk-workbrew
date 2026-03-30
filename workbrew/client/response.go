package client

import (
	"fmt"
	"net/http"
	"strings"

	"go.uber.org/zap"
	"resty.dev/v3"
)

// GetRateLimitHeaders extracts Workbrew API rate limiting headers from the response.
// These headers indicate the API quota limits and current usage.
//
// Returns limit, remaining, reset, and retryAfter header values.
func GetRateLimitHeaders(resp *resty.Response) (limit, remaining, reset, retryAfter string) {
	if resp == nil {
		return
	}
	return resp.Header().Get("X-Api-Quota-Limit"),
		resp.Header().Get("X-Api-Quota-Remaining"),
		resp.Header().Get("X-Api-Quota-Reset"),
		resp.Header().Get("Retry-After")
}

// GetResponseHeader retrieves a single header value from the response by key.
// Header lookup is case-insensitive following HTTP standards.
func GetResponseHeader(resp *resty.Response, key string) string {
	if resp == nil {
		return ""
	}
	return resp.Header().Get(key)
}

// GetResponseHeaders returns all HTTP headers from the response.
func GetResponseHeaders(resp *resty.Response) http.Header {
	if resp == nil {
		return make(http.Header)
	}
	return resp.Header()
}

// validateResponse validates the HTTP response before processing.
// Checks for unexpected Content-Type on successful JSON responses.
func (t *Transport) validateResponse(resp *resty.Response, method, path string) error {
	// Handle empty responses (204 No Content, etc.)
	bodyLen := len(resp.Bytes())
	if resp.Header().Get("Content-Length") == "0" || bodyLen == 0 {
		t.logger.Debug("Empty response received",
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status_code", resp.StatusCode()))
		return nil
	}

	// For non-error responses with content, validate Content-Type is JSON.
	// Error responses skip validation (handled by error parser).
	if !resp.IsError() && bodyLen > 0 {
		contentType := resp.Header().Get("Content-Type")

		// Allow responses without Content-Type header (some endpoints don't set it)
		if contentType != "" && !strings.HasPrefix(contentType, "application/json") {
			t.logger.Warn("Unexpected Content-Type in response",
				zap.String("method", method),
				zap.String("path", path),
				zap.String("content_type", contentType),
				zap.String("expected", "application/json"))

			return fmt.Errorf("unexpected response Content-Type from %s %s: got %q, expected application/json",
				method, path, contentType)
		}
	}

	return nil
}
