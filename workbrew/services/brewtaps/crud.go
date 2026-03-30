package brewtaps

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/constants"
	"resty.dev/v3"
)

// BrewTapsServiceInterface defines the interface for brew taps operations
//
// Workbrew API docs: https://console.workbrew.com/documentation/api
type BrewTapsServiceInterface interface {
	// ListBrewTaps returns a list of Taps
	//
	// Returns Homebrew taps with their names, assigned devices, installed formulae/casks counts, and available packages.
	ListBrewTaps(ctx context.Context) (*BrewTapsResponse, *resty.Response, error)

	// ListBrewTapsCSV returns a list of Taps in CSV format
	//
	// Returns tap data as CSV with columns: tap, devices, formulae_installed, casks_installed, available_packages.
	ListBrewTapsCSV(ctx context.Context) ([]byte, *resty.Response, error)
}

// BrewTaps handles communication with the brew taps
// related methods of the Workbrew API.
//
// Workbrew API docs: https://console.workbrew.com/documentation/api
type BrewTaps struct {
	client client.Client
}

// Ensure BrewTaps implements BrewTapsServiceInterface
var _ BrewTapsServiceInterface = (*BrewTaps)(nil)

// NewBrewTaps creates a new brew taps service
func NewBrewTaps(client client.Client) *BrewTaps {
	return &BrewTaps{
		client: client,
	}
}

// ListBrewTaps retrieves all brew taps in JSON format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/brew_taps.json
func (s *BrewTaps) ListBrewTaps(ctx context.Context) (*BrewTapsResponse, *resty.Response, error) {
	var result BrewTapsResponse
	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetResult(&result).
		Get(constants.EndpointBrewTapsJSON)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// ListBrewTapsCSV retrieves all brew taps in CSV format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/brew_taps.csv
func (s *BrewTaps) ListBrewTapsCSV(ctx context.Context) ([]byte, *resty.Response, error) {
	resp, csvData, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.TextCSV).
		GetBytes(constants.EndpointBrewTapsCSV)
	if err != nil {
		return nil, resp, err
	}

	return csvData, resp, nil
}
