package brewcommands

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/constants"
	"resty.dev/v3"
)

type (
	// BrewCommands handles communication with the brew commands-related methods of the Workbrew API.
	//
	// Workbrew API docs: https://console.workbrew.com/documentation/api
	BrewCommands struct {
		client client.Client
	}
)

func NewBrewCommands(client client.Client) *BrewCommands {
	return &BrewCommands{client: client}
}

// ListV0 retrieves all brew commands in JSON format.
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/brew_commands.json
func (s *BrewCommands) ListV0(ctx context.Context) (*BrewCommandsResponse, *resty.Response, error) {
	var result BrewCommandsResponse
	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetResult(&result).
		Get(constants.EndpointBrewCommandsJSON)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// ListCSVV0 retrieves all brew commands in CSV format.
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/brew_commands.csv
func (s *BrewCommands) ListCSVV0(ctx context.Context) ([]byte, *resty.Response, error) {
	resp, csvData, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.TextCSV).
		GetBytes(constants.EndpointBrewCommandsCSV)
	if err != nil {
		return nil, resp, err
	}

	return csvData, resp, nil
}

// CreateV0 creates a new brew command.
// URL: POST https://console.workbrew.com/workspaces/{workspace_name}/brew_commands.json
//
// Response codes:
//   - 201: Brew Command created successfully
//   - 403: On a Free tier plan (requires upgrade)
//   - 422: Validation error (e.g., "Arguments cannot include `&&`")
func (s *BrewCommands) CreateV0(ctx context.Context, request *CreateBrewCommandRequest) (*CreateBrewCommandResponse, *resty.Response, error) {
	var result CreateBrewCommandResponse
	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetBody(request).
		SetResult(&result).
		Post(constants.EndpointBrewCommandsJSON)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// ListRunsByLabelV0 retrieves all runs for a specific brew command in JSON format.
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/brew_commands/{brew_command_label}/runs.json
func (s *BrewCommands) ListRunsByLabelV0(ctx context.Context, brewCommandLabel string) (*BrewCommandRunsResponse, *resty.Response, error) {
	if brewCommandLabel == "" {
		return nil, nil, fmt.Errorf("brew command label is required")
	}

	endpoint := fmt.Sprintf(constants.EndpointBrewCommandRunsJSONFormat, brewCommandLabel)

	var result BrewCommandRunsResponse
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

// ListRunsByLabelCSVV0 retrieves all runs for a specific brew command in CSV format.
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/brew_commands/{brew_command_label}/runs.csv
func (s *BrewCommands) ListRunsByLabelCSVV0(ctx context.Context, brewCommandLabel string) ([]byte, *resty.Response, error) {
	if brewCommandLabel == "" {
		return nil, nil, fmt.Errorf("brew command label is required")
	}

	endpoint := fmt.Sprintf(constants.EndpointBrewCommandRunsCSVFormat, brewCommandLabel)

	resp, csvData, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.TextCSV).
		GetBytes(endpoint)
	if err != nil {
		return nil, resp, err
	}

	return csvData, resp, nil
}
