package licenses

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/constants"
	"resty.dev/v3"
)

type (
	// Licenses handles communication with the licenses-related methods of the Workbrew API.
	//
	// Workbrew API docs: https://console.workbrew.com/documentation/api
	Licenses struct {
		client client.Client
	}
)

func NewLicenses(client client.Client) *Licenses {
	return &Licenses{client: client}
}

// ListV0 retrieves all licenses in JSON format.
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/licenses.json
func (s *Licenses) ListV0(ctx context.Context) (*LicensesResponse, *resty.Response, error) {
	var result LicensesResponse
	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetResult(&result).
		Get(constants.EndpointLicensesJSON)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// ListCSVV0 retrieves all licenses in CSV format.
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/licenses.csv
func (s *Licenses) ListCSVV0(ctx context.Context) ([]byte, *resty.Response, error) {
	resp, csvData, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.TextCSV).
		GetBytes(constants.EndpointLicensesCSV)
	if err != nil {
		return nil, resp, err
	}

	return csvData, resp, nil
}
