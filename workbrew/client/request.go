package client

import (
	"fmt"

	"go.uber.org/zap"
	"resty.dev/v3"
)

// executeRequest is the centralized request executor used by the RequestBuilder.
// It handles error processing and returns both the resty response and any error.
func (t *Transport) executeRequest(req *resty.Request, method, path string) (*resty.Response, error) {
	t.logger.Debug("Executing API request",
		zap.String("method", method),
		zap.String("path", path))

	resp, err := req.Execute(method, path)

	if err != nil {
		t.logger.Error("Request failed",
			zap.String("method", method),
			zap.String("path", path),
			zap.Error(err))
		return resp, fmt.Errorf("request failed: %w", err)
	}

	if err := t.validateResponse(resp, method, path); err != nil {
		return resp, err
	}

	if resp.IsError() {
		return resp, ParseErrorResponse(
			resp.Bytes(),
			resp.StatusCode(),
			resp.Status(),
			method,
			path,
			t.logger,
		)
	}

	t.logger.Debug("Request completed successfully",
		zap.String("method", method),
		zap.String("path", path),
		zap.Int("status_code", resp.StatusCode()))

	return resp, nil
}
