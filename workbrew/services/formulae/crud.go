package formulae

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/constants"
	"resty.dev/v3"
)

// FormulaeServiceInterface defines the interface for formulae operations
//
// Workbrew API docs: https://console.workbrew.com/documentation/api
type FormulaeServiceInterface interface {
	// ListFormulae returns a list of Formulae
	//
	// Returns installed Homebrew formulae with names, assigned devices, outdated status, installation type (on request/dependency),
	// known vulnerabilities, deprecation status, licenses, and Homebrew core versions.
	ListFormulae(ctx context.Context) (*FormulaeResponse, *resty.Response, error)

	// ListFormulaeCSV returns a list of Formulae in CSV format
	//
	// Returns formulae data as CSV with columns: name, devices, outdated, installed_on_request, installed_as_dependency,
	// vulnerabilities, deprecated, license, homebrew_core_version.
	ListFormulaeCSV(ctx context.Context) ([]byte, *resty.Response, error)
}

// Formulae handles communication with the formulae
// related methods of the Workbrew API.
//
// Workbrew API docs: https://console.workbrew.com/documentation/api
type Formulae struct {
	client client.Client
}

// Ensure Formulae implements FormulaeServiceInterface
var _ FormulaeServiceInterface = (*Formulae)(nil)

// NewFormulae creates a new formulae service
func NewFormulae(client client.Client) *Formulae {
	return &Formulae{
		client: client,
	}
}

// ListFormulae retrieves all formulae in JSON format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/formulae.json
func (s *Formulae) ListFormulae(ctx context.Context) (*FormulaeResponse, *resty.Response, error) {
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

// ListFormulaeCSV retrieves all formulae in CSV format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/formulae.csv
func (s *Formulae) ListFormulaeCSV(ctx context.Context) ([]byte, *resty.Response, error) {
	resp, csvData, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.TextCSV).
		GetBytes(constants.EndpointFormulaeCSV)
	if err != nil {
		return nil, resp, err
	}

	return csvData, resp, nil
}
