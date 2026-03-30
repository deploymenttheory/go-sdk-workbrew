package brewfiles

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/constants"
	"resty.dev/v3"
)

// BrewfilesServiceInterface defines the interface for brewfiles operations
//
// Workbrew API docs: https://console.workbrew.com/documentation/api
type BrewfilesServiceInterface interface {
	// ListBrewfiles returns a list of brewfiles with their status and device assignments
	//
	// Returns brewfiles with last updated user, start/finish timestamps, assigned devices, and run count
	ListBrewfiles(ctx context.Context) (*BrewfilesResponse, *resty.Response, error)

	// ListBrewfilesCSV returns a list of brewfiles in CSV format
	//
	// Returns the same brewfiles data as ListBrewfiles but formatted as CSV
	ListBrewfilesCSV(ctx context.Context) ([]byte, *resty.Response, error)

	// CreateBrewfile creates a new brewfile with specified label, content, and device/group assignment
	//
	// Requires label and content fields. Can assign to specific devices via device_serial_numbers or to a device group via device_group_id
	CreateBrewfile(ctx context.Context, request *CreateBrewfileRequest) (*BrewfileMessageResponse, *resty.Response, error)

	// UpdateBrewfile updates an existing brewfile's content and device assignments
	//
	// Updates the brewfile identified by label. Can update content, device_serial_numbers, or device_group_id
	UpdateBrewfile(ctx context.Context, label string, request *UpdateBrewfileRequest) (*BrewfileMessageResponse, *resty.Response, error)

	// DeleteBrewfile deletes a brewfile by its label
	//
	// Permanently removes the brewfile identified by the specified label
	DeleteBrewfile(ctx context.Context, label string) (*BrewfileMessageResponse, *resty.Response, error)

	// ListBrewfileRuns returns a list of brewfile runs for a specific brewfile
	//
	// Returns run history including label, device, timestamps, success status, and output for the specified brewfile label
	ListBrewfileRuns(ctx context.Context, label string) (*BrewfileRunsResponse, *resty.Response, error)

	// ListBrewfileRunsCSV returns a list of brewfile runs in CSV format
	//
	// Returns the same run data as ListBrewfileRuns but formatted as CSV
	ListBrewfileRunsCSV(ctx context.Context, label string) ([]byte, *resty.Response, error)
}

// Brewfiles handles communication with the brewfiles
// related methods of the Workbrew API.
//
// Workbrew API docs: https://console.workbrew.com/documentation/api
type Brewfiles struct {
	client client.Client
}

// Ensure Brewfiles implements BrewfilesServiceInterface
var _ BrewfilesServiceInterface = (*Brewfiles)(nil)

// NewBrewfiles creates a new brewfiles service
func NewBrewfiles(client client.Client) *Brewfiles {
	return &Brewfiles{
		client: client,
	}
}

// ListBrewfiles retrieves all brewfiles in JSON format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/brewfiles.json
func (s *Brewfiles) ListBrewfiles(ctx context.Context) (*BrewfilesResponse, *resty.Response, error) {
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

// ListBrewfilesCSV retrieves all brewfiles in CSV format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/brewfiles.csv
func (s *Brewfiles) ListBrewfilesCSV(ctx context.Context) ([]byte, *resty.Response, error) {
	resp, csvData, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.TextCSV).
		GetBytes(constants.EndpointBrewfilesCSV)
	if err != nil {
		return nil, resp, err
	}

	return csvData, resp, nil
}

// CreateBrewfile creates a new brewfile
// URL: POST https://console.workbrew.com/workspaces/{workspace_name}/brewfiles.json
func (s *Brewfiles) CreateBrewfile(ctx context.Context, request *CreateBrewfileRequest) (*BrewfileMessageResponse, *resty.Response, error) {
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

// UpdateBrewfile updates an existing brewfile
// URL: PUT https://console.workbrew.com/workspaces/{workspace_name}/brewfiles/{label}.json
//
// Response codes:
//   - 200: Brewfile updated successfully
//   - 422: Validation error
func (s *Brewfiles) UpdateBrewfile(ctx context.Context, label string, request *UpdateBrewfileRequest) (*BrewfileMessageResponse, *resty.Response, error) {
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

// DeleteBrewfile deletes a brewfile
// URL: DELETE https://console.workbrew.com/workspaces/{workspace_name}/brewfiles/{label}.json
//
// Response codes:
//   - 200: Brewfile deleted successfully
func (s *Brewfiles) DeleteBrewfile(ctx context.Context, label string) (*BrewfileMessageResponse, *resty.Response, error) {
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

// ListBrewfileRuns retrieves all runs for a specific brewfile in JSON format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/brewfiles/{label}/runs.json
func (s *Brewfiles) ListBrewfileRuns(ctx context.Context, label string) (*BrewfileRunsResponse, *resty.Response, error) {
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

// ListBrewfileRunsCSV retrieves all runs for a specific brewfile in CSV format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/brewfiles/{label}/runs.csv
func (s *Brewfiles) ListBrewfileRunsCSV(ctx context.Context, label string) ([]byte, *resty.Response, error) {
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
