package brewcommands

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/constants"
	"resty.dev/v3"
)

// BrewCommandsServiceInterface defines the interface for brew commands operations
//
// Workbrew API docs: https://console.workbrew.com/documentation/api
type BrewCommandsServiceInterface interface {
	// ListBrewCommands returns a list of brew commands with their configuration and status
	//
	// Returns brew commands with command, label, last updated user, start/finish timestamps, devices, and run count
	ListBrewCommands(ctx context.Context) (*BrewCommandsResponse, *resty.Response, error)

	// ListBrewCommandsCSV returns a list of brew commands in CSV format
	//
	// Returns the same brew commands data as ListBrewCommands but formatted as CSV
	ListBrewCommandsCSV(ctx context.Context) ([]byte, *resty.Response, error)

	// CreateBrewCommand creates a new brew command with specified arguments and optional device/timing configuration
	//
	// Requires arguments field. Optional fields include device_ids, run_after_datetime, and recurrence (once, daily, weekly, monthly)
	CreateBrewCommand(ctx context.Context, request *CreateBrewCommandRequest) (*CreateBrewCommandResponse, *resty.Response, error)

	// ListBrewCommandRuns returns a list of brew command runs for a specific brew command
	//
	// Returns run history including command, label, device, timestamps, success status, and output for the specified brew command label
	ListBrewCommandRuns(ctx context.Context, brewCommandLabel string) (*BrewCommandRunsResponse, *resty.Response, error)

	// ListBrewCommandRunsCSV returns a list of brew command runs in CSV format
	//
	// Returns the same run data as ListBrewCommandRuns but formatted as CSV
	ListBrewCommandRunsCSV(ctx context.Context, brewCommandLabel string) ([]byte, *resty.Response, error)
}

// BrewCommands handles communication with the brew commands
// related methods of the Workbrew API.
//
// Workbrew API docs: https://console.workbrew.com/documentation/api
type BrewCommands struct {
	client client.Client
}

// Ensure BrewCommands implements BrewCommandsServiceInterface
var _ BrewCommandsServiceInterface = (*BrewCommands)(nil)

// NewBrewCommands creates a new brew commands service
func NewBrewCommands(client client.Client) *BrewCommands {
	return &BrewCommands{
		client: client,
	}
}

// ListBrewCommands retrieves all brew commands in JSON format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/brew_commands.json
func (s *BrewCommands) ListBrewCommands(ctx context.Context) (*BrewCommandsResponse, *resty.Response, error) {
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

// ListBrewCommandsCSV retrieves all brew commands in CSV format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/brew_commands.csv
func (s *BrewCommands) ListBrewCommandsCSV(ctx context.Context) ([]byte, *resty.Response, error) {
	resp, csvData, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.TextCSV).
		GetBytes(constants.EndpointBrewCommandsCSV)
	if err != nil {
		return nil, resp, err
	}

	return csvData, resp, nil
}

// CreateBrewCommand creates a new brew command
// URL: POST https://console.workbrew.com/workspaces/{workspace_name}/brew_commands.json
//
// Response codes:
//   - 201: Brew Command created successfully
//   - 403: On a Free tier plan (requires upgrade)
//   - 422: Validation error (e.g., "Arguments cannot include `&&`")
func (s *BrewCommands) CreateBrewCommand(ctx context.Context, request *CreateBrewCommandRequest) (*CreateBrewCommandResponse, *resty.Response, error) {
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

// ListBrewCommandRuns retrieves all runs for a specific brew command in JSON format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/brew_commands/{brew_command_label}/runs.json
func (s *BrewCommands) ListBrewCommandRuns(ctx context.Context, brewCommandLabel string) (*BrewCommandRunsResponse, *resty.Response, error) {
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

// ListBrewCommandRunsCSV retrieves all runs for a specific brew command in CSV format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/brew_commands/{brew_command_label}/runs.csv
func (s *BrewCommands) ListBrewCommandRunsCSV(ctx context.Context, brewCommandLabel string) ([]byte, *resty.Response, error) {
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
