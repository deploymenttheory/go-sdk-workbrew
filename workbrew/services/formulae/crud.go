package formulae

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/constants"
	"resty.dev/v3"
)

type (
	// Formulae handles communication with the formulae-related methods of the Workbrew API.
	//
	// Workbrew API docs: https://console.workbrew.com/documentation/api
	Formulae struct {
		client client.Client
	}
)

func NewFormulae(client client.Client) *Formulae {
	return &Formulae{client: client}
}

// ListV0 retrieves all formulae in JSON format.
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/formulae.json
func (s *Formulae) ListV0(ctx context.Context) (*FormulaeResponse, *resty.Response, error) {
	var result FormulaeResponse
	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetResult(&result).
		Get(constants.EndpointFormulaeJSON)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// ListCSVV0 retrieves all formulae in CSV format.
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/formulae.csv
func (s *Formulae) ListCSVV0(ctx context.Context) ([]byte, *resty.Response, error) {
	resp, csvData, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.TextCSV).
		GetBytes(constants.EndpointFormulaeCSV)
	if err != nil {
		return nil, resp, err
	}

	return csvData, resp, nil
}
