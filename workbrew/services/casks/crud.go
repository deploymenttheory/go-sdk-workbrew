package casks

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/constants"
	"resty.dev/v3"
)

// CasksServiceInterface defines the interface for casks operations
//
// Workbrew API docs: https://console.workbrew.com/documentation/api
type CasksServiceInterface interface {
	// ListCasks returns a list of Casks
	//
	// Returns installed Homebrew casks with their names, assigned devices, outdated status, deprecation info, and versions.
	ListCasks(ctx context.Context) (*CasksResponse, *resty.Response, error)

	// ListCasksCSV returns a list of Casks in CSV format
	//
	// Returns cask data as CSV with columns: name, devices, outdated, deprecated, homebrew_cask_version.
	ListCasksCSV(ctx context.Context) ([]byte, *resty.Response, error)
}

// Casks handles communication with the casks
// related methods of the Workbrew API.
//
// Workbrew API docs: https://console.workbrew.com/documentation/api
type Casks struct {
	client client.Client
}

// Ensure Casks implements CasksServiceInterface
var _ CasksServiceInterface = (*Casks)(nil)

// NewCasks creates a new casks service
func NewCasks(client client.Client) *Casks {
	return &Casks{
		client: client,
	}
}

// ListCasks retrieves all casks in JSON format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/casks.json
func (s *Casks) ListCasks(ctx context.Context) (*CasksResponse, *resty.Response, error) {
	var result CasksResponse
	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetResult(&result).
		Get(constants.EndpointCasksJSON)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// ListCasksCSV retrieves all casks in CSV format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/casks.csv
func (s *Casks) ListCasksCSV(ctx context.Context) ([]byte, *resty.Response, error) {
	resp, csvData, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.TextCSV).
		GetBytes(constants.EndpointCasksCSV)
	if err != nil {
		return nil, resp, err
	}

	return csvData, resp, nil
}
