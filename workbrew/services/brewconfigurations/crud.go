package brewconfigurations

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/constants"
	"resty.dev/v3"
)

type (
	// BrewConfigurations handles communication with the brew configurations-related methods of the Workbrew API.
	//
	// Workbrew API docs: https://console.workbrew.com/documentation/api
	BrewConfigurations struct {
		client client.Client
	}
)

func NewBrewConfigurations(client client.Client) *BrewConfigurations {
	return &BrewConfigurations{client: client}
}

// ListV0 retrieves all brew configurations in JSON format.
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/brew_configurations.json
func (s *BrewConfigurations) ListV0(ctx context.Context) (*BrewConfigurationsResponse, *resty.Response, error) {
	var result BrewConfigurationsResponse
	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetResult(&result).
		Get(constants.EndpointBrewConfigurationsJSON)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// ListCSVV0 retrieves all brew configurations in CSV format.
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/brew_configurations.csv
func (s *BrewConfigurations) ListCSVV0(ctx context.Context) ([]byte, *resty.Response, error) {
	resp, csvData, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.TextCSV).
		GetBytes(constants.EndpointBrewConfigurationsCSV)
	if err != nil {
		return nil, resp, err
	}

	return csvData, resp, nil
}
