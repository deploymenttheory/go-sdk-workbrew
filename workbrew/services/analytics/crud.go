package analytics

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/constants"
	"resty.dev/v3"
)

type (
	// Analytics handles communication with the analytics-related methods of the Workbrew API.
	//
	// Workbrew API docs: https://console.workbrew.com/documentation/api
	Analytics struct {
		client client.Client
	}
)

func NewAnalytics(client client.Client) *Analytics {
	return &Analytics{client: client}
}

// ListV0 retrieves all analytics in JSON format.
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/analytics.json
func (s *Analytics) ListV0(ctx context.Context) (*AnalyticsResponse, *resty.Response, error) {
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

// ListCSVV0 retrieves all analytics in CSV format.
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/analytics.csv
func (s *Analytics) ListCSVV0(ctx context.Context) ([]byte, *resty.Response, error) {
	resp, csvData, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.TextCSV).
		GetBytes(constants.EndpointAnalyticsCSV)
	if err != nil {
		return nil, resp, err
	}

	return csvData, resp, nil
}
