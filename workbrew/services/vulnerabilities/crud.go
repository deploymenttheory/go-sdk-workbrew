package vulnerabilities

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/constants"
	"resty.dev/v3"
)

type (
	// Vulnerabilities handles communication with the vulnerabilities-related methods of the Workbrew API.
	//
	// Workbrew API docs: https://console.workbrew.com/documentation/api
	Vulnerabilities struct {
		client client.Client
	}
)

func NewVulnerabilities(client client.Client) *Vulnerabilities {
	return &Vulnerabilities{client: client}
}

// ListV0 retrieves all vulnerabilities in JSON format.
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/vulnerabilities.json
//
// Note: This endpoint may return 403 on Free tier plans.
func (s *Vulnerabilities) ListV0(ctx context.Context) (*VulnerabilitiesResponse, *resty.Response, error) {
	var result VulnerabilitiesResponse
	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetResult(&result).
		Get(constants.EndpointVulnerabilitiesJSON)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// ListCSVV0 retrieves all vulnerabilities in CSV format.
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/vulnerabilities.csv
//
// Note: This endpoint may return 403 on Free tier plans.
func (s *Vulnerabilities) ListCSVV0(ctx context.Context) ([]byte, *resty.Response, error) {
	resp, csvData, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.TextCSV).
		GetBytes(constants.EndpointVulnerabilitiesCSV)
	if err != nil {
		return nil, resp, err
	}

	return csvData, resp, nil
}
