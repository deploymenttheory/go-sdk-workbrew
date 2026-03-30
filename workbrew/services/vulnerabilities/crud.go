package vulnerabilities

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/constants"
	"resty.dev/v3"
)

// VulnerabilitiesServiceInterface defines the interface for vulnerabilities operations
//
// Workbrew API docs: https://console.workbrew.com/documentation/api
type VulnerabilitiesServiceInterface interface {
	// ListVulnerabilities returns a list of Vulnerabilities
	//
	// Returns security vulnerabilities affecting installed formulae, including CVE IDs with CVSS scores,
	// affected formula names, outdated devices, support status, and Homebrew core versions.
	// May return 403 Forbidden on Free tier plans.
	ListVulnerabilities(ctx context.Context) (*VulnerabilitiesResponse, *resty.Response, error)

	// ListVulnerabilitiesCSV returns a list of Vulnerabilities in CSV format
	//
	// Returns vulnerability data as CSV with columns: vulnerabilities, formula, outdated_devices, supported, homebrew_core_version.
	ListVulnerabilitiesCSV(ctx context.Context) ([]byte, *resty.Response, error)
}

// Vulnerabilities handles communication with the vulnerabilities
// related methods of the Workbrew API.
//
// Workbrew API docs: https://console.workbrew.com/documentation/api
type Vulnerabilities struct {
	client client.Client
}

// Ensure Vulnerabilities implements VulnerabilitiesServiceInterface
var _ VulnerabilitiesServiceInterface = (*Vulnerabilities)(nil)

// NewVulnerabilities creates a new vulnerabilities service
func NewVulnerabilities(client client.Client) *Vulnerabilities {
	return &Vulnerabilities{
		client: client,
	}
}

// ListVulnerabilities retrieves all vulnerabilities in JSON format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/vulnerabilities.json
//
// Note: This endpoint may return 403 on Free tier plans
func (s *Vulnerabilities) ListVulnerabilities(ctx context.Context) (*VulnerabilitiesResponse, *resty.Response, error) {
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

// ListVulnerabilitiesCSV retrieves all vulnerabilities in CSV format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/vulnerabilities.csv
//
// Note: This endpoint may return 403 on Free tier plans
func (s *Vulnerabilities) ListVulnerabilitiesCSV(ctx context.Context) ([]byte, *resty.Response, error) {
	resp, csvData, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.TextCSV).
		GetBytes(constants.EndpointVulnerabilitiesCSV)
	if err != nil {
		return nil, resp, err
	}

	return csvData, resp, nil
}
