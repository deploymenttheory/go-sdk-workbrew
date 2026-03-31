package brewfiles

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/constants"
	"resty.dev/v3"
)

type (
	// Brewfiles handles communication with the brewfiles-related methods of the Workbrew API.
	//
	// Workbrew API docs: https://console.workbrew.com/documentation/api
	Brewfiles struct {
		client client.Client
	}
)

func NewBrewfiles(client client.Client) *Brewfiles {
	return &Brewfiles{client: client}
}

// ListV0 retrieves all brewfiles in JSON format.
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/brewfiles.json
func (s *Brewfiles) ListV0(ctx context.Context) (*BrewfilesResponse, *resty.Response, error) {
	var result BrewfilesResponse
	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetResult(&result).
		Get(constants.EndpointBrewfilesJSON)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// ListCSVV0 retrieves all brewfiles in CSV format.
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/brewfiles.csv
func (s *Brewfiles) ListCSVV0(ctx context.Context) ([]byte, *resty.Response, error) {
	resp, csvData, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.TextCSV).
		GetBytes(constants.EndpointBrewfilesCSV)
	if err != nil {
		return nil, resp, err
	}

	return csvData, resp, nil
}

// CreateV0 creates a new brewfile.
// URL: POST https://console.workbrew.com/workspaces/{workspace_name}/brewfiles.json
func (s *Brewfiles) CreateV0(ctx context.Context, request *CreateBrewfileRequest) (*BrewfileMessageResponse, *resty.Response, error) {
	var result BrewfileMessageResponse
	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetBody(request).
		SetResult(&result).
		Post(constants.EndpointBrewfilesJSON)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// UpdateByLabelV0 updates an existing brewfile.
// URL: PUT https://console.workbrew.com/workspaces/{workspace_name}/brewfiles/{label}.json
//
// Response codes:
//   - 200: Brewfile updated successfully
//   - 422: Validation error
func (s *Brewfiles) UpdateByLabelV0(ctx context.Context, label string, request *UpdateBrewfileRequest) (*BrewfileMessageResponse, *resty.Response, error) {
	if label == "" {
		return nil, nil, fmt.Errorf("brewfile label is required")
	}

	endpoint := fmt.Sprintf(constants.EndpointBrewfileLabelFormat, label)

	var result BrewfileMessageResponse
	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetBody(request).
		SetResult(&result).
		Put(endpoint)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// DeleteByLabelV0 deletes a brewfile.
// URL: DELETE https://console.workbrew.com/workspaces/{workspace_name}/brewfiles/{label}.json
//
// Response codes:
//   - 200: Brewfile deleted successfully
func (s *Brewfiles) DeleteByLabelV0(ctx context.Context, label string) (*BrewfileMessageResponse, *resty.Response, error) {
	if label == "" {
		return nil, nil, fmt.Errorf("brewfile label is required")
	}

	endpoint := fmt.Sprintf(constants.EndpointBrewfileLabelFormat, label)

	var result BrewfileMessageResponse
	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetResult(&result).
		Delete(endpoint)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// ListRunsByLabelV0 retrieves all runs for a specific brewfile in JSON format.
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/brewfiles/{label}/runs.json
func (s *Brewfiles) ListRunsByLabelV0(ctx context.Context, label string) (*BrewfileRunsResponse, *resty.Response, error) {
	if label == "" {
		return nil, nil, fmt.Errorf("brewfile label is required")
	}

	endpoint := fmt.Sprintf(constants.EndpointBrewfileRunsJSONFormat, label)

	var result BrewfileRunsResponse
	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetResult(&result).
		Get(endpoint)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// ListRunsByLabelCSVV0 retrieves all runs for a specific brewfile in CSV format.
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/brewfiles/{label}/runs.csv
func (s *Brewfiles) ListRunsByLabelCSVV0(ctx context.Context, label string) ([]byte, *resty.Response, error) {
	if label == "" {
		return nil, nil, fmt.Errorf("brewfile label is required")
	}

	endpoint := fmt.Sprintf(constants.EndpointBrewfileRunsCSVFormat, label)

	resp, csvData, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.TextCSV).
		GetBytes(endpoint)
	if err != nil {
		return nil, resp, err
	}

	return csvData, resp, nil
}
