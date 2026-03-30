package analytics

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/constants"
	"resty.dev/v3"
)

// AnalyticsServiceInterface defines the interface for analytics operations
//
// Workbrew API docs: https://console.workbrew.com/documentation/api
type AnalyticsServiceInterface interface {
	// ListAnalytics returns a list of analytics data showing command usage statistics per device
	//
	// Returns analytics records with device, command, last run timestamp, and count information
	ListAnalytics(ctx context.Context) (*AnalyticsResponse, *resty.Response, error)

	// ListAnalyticsCSV returns a list of analytics data in CSV format
	//
	// Returns the same analytics data as ListAnalytics but formatted as CSV
	ListAnalyticsCSV(ctx context.Context) ([]byte, *resty.Response, error)
}

// Analytics handles communication with the analytics
// related methods of the Workbrew API.
//
// Workbrew API docs: https://console.workbrew.com/documentation/api
type Analytics struct {
	client client.Client
}

// Ensure Analytics implements AnalyticsServiceInterface
var _ AnalyticsServiceInterface = (*Analytics)(nil)

// NewAnalytics creates a new analytics service
func NewAnalytics(client client.Client) *Analytics {
	return &Analytics{
		client: client,
	}
}

// ListAnalytics retrieves all analytics in JSON format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/analytics.json
func (s *Analytics) ListAnalytics(ctx context.Context) (*AnalyticsResponse, *resty.Response, error) {
	var result AnalyticsResponse
	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetResult(&result).
		Get(constants.EndpointAnalyticsJSON)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// ListAnalyticsCSV retrieves all analytics in CSV format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/analytics.csv
func (s *Analytics) ListAnalyticsCSV(ctx context.Context) ([]byte, *resty.Response, error) {
	resp, csvData, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.TextCSV).
		GetBytes(constants.EndpointAnalyticsCSV)
	if err != nil {
		return nil, resp, err
	}

	return csvData, resp, nil
}
