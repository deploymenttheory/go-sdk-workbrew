package brewconfigurations

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/constants"
	"resty.dev/v3"
)

// BrewConfigurationsServiceInterface defines the interface for brew configurations operations
//
// Workbrew API docs: https://console.workbrew.com/documentation/api
type BrewConfigurationsServiceInterface interface {
	// ListBrewConfigurations returns a list of Brew Configurations
	//
	// Returns Homebrew environment variable configurations with their keys, values, last updated user, and assigned device groups.
	ListBrewConfigurations(ctx context.Context) (*BrewConfigurationsResponse, *resty.Response, error)

	// ListBrewConfigurationsCSV returns a list of Brew Configurations in CSV format
	//
	// Returns brew configuration data as CSV with columns: key, value, last_updated_by_user, device_group.
	ListBrewConfigurationsCSV(ctx context.Context) ([]byte, *resty.Response, error)
}

// BrewConfigurations handles communication with the brew configurations
// related methods of the Workbrew API.
//
// Workbrew API docs: https://console.workbrew.com/documentation/api
type BrewConfigurations struct {
	client client.Client
}

// Ensure BrewConfigurations implements BrewConfigurationsServiceInterface
var _ BrewConfigurationsServiceInterface = (*BrewConfigurations)(nil)

// NewBrewConfigurations creates a new brew configurations service
func NewBrewConfigurations(client client.Client) *BrewConfigurations {
	return &BrewConfigurations{
		client: client,
	}
}

// ListBrewConfigurations retrieves all brew configurations in JSON format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/brew_configurations.json
func (s *BrewConfigurations) ListBrewConfigurations(ctx context.Context) (*BrewConfigurationsResponse, *resty.Response, error) {
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

// ListBrewConfigurationsCSV retrieves all brew configurations in CSV format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/brew_configurations.csv
func (s *BrewConfigurations) ListBrewConfigurationsCSV(ctx context.Context) ([]byte, *resty.Response, error) {
	resp, csvData, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.TextCSV).
		GetBytes(constants.EndpointBrewConfigurationsCSV)
	if err != nil {
		return nil, resp, err
	}

	return csvData, resp, nil
}
