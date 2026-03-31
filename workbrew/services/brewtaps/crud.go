package brewtaps

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/constants"
	"resty.dev/v3"
)

type (
	// BrewTaps handles communication with the brew taps-related methods of the Workbrew API.
	//
	// Workbrew API docs: https://console.workbrew.com/documentation/api
	BrewTaps struct {
		client client.Client
	}
)

func NewBrewTaps(client client.Client) *BrewTaps {
	return &BrewTaps{client: client}
}

// ListV0 retrieves all brew taps in JSON format.
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/brew_taps.json
func (s *BrewTaps) ListV0(ctx context.Context) (*BrewTapsResponse, *resty.Response, error) {
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

// ListCSVV0 retrieves all brew taps in CSV format.
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/brew_taps.csv
func (s *BrewTaps) ListCSVV0(ctx context.Context) ([]byte, *resty.Response, error) {
	resp, csvData, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.TextCSV).
		GetBytes(constants.EndpointBrewTapsCSV)
	if err != nil {
		return nil, resp, err
	}

	return csvData, resp, nil
}
