package casks

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/constants"
	"resty.dev/v3"
)

type (
	// Casks handles communication with the casks-related methods of the Workbrew API.
	//
	// Workbrew API docs: https://console.workbrew.com/documentation/api
	Casks struct {
		client client.Client
	}
)

func NewCasks(client client.Client) *Casks {
	return &Casks{client: client}
}

// ListV0 retrieves all casks in JSON format.
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/casks.json
func (s *Casks) ListV0(ctx context.Context) (*CasksResponse, *resty.Response, error) {
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

// ListCSVV0 retrieves all casks in CSV format.
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/casks.csv
func (s *Casks) ListCSVV0(ctx context.Context) ([]byte, *resty.Response, error) {
	resp, csvData, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.TextCSV).
		GetBytes(constants.EndpointCasksCSV)
	if err != nil {
		return nil, resp, err
	}

	return csvData, resp, nil
}
