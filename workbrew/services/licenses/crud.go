package licenses

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/constants"
	"resty.dev/v3"
)

// LicensesServiceInterface defines the interface for licenses operations
//
// Workbrew API docs: https://console.workbrew.com/documentation/api
type LicensesServiceInterface interface {
	// ListLicenses returns a list of Licenses
	//
	// Returns software licenses found across installed formulae, with license names and counts of affected devices and formulae.
	ListLicenses(ctx context.Context) (*LicensesResponse, *resty.Response, error)

	// ListLicensesCSV returns a list of Licenses in CSV format
	//
	// Returns license data as CSV with columns: name, device_count, formula_count.
	ListLicensesCSV(ctx context.Context) ([]byte, *resty.Response, error)
}

// Licenses handles communication with the licenses
// related methods of the Workbrew API.
//
// Workbrew API docs: https://console.workbrew.com/documentation/api
type Licenses struct {
	client client.Client
}

// Ensure Licenses implements LicensesServiceInterface
var _ LicensesServiceInterface = (*Licenses)(nil)

// NewLicenses creates a new licenses service
func NewLicenses(client client.Client) *Licenses {
	return &Licenses{
		client: client,
	}
}

// ListLicenses retrieves all licenses in JSON format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/licenses.json
//
// Example cURL:
//
//	curl -X GET \
//	  -H "Authorization: Bearer YOUR_API_KEY" \
//	  -H "X-Workbrew-API-Version: v0" \
//	  -H "Accept: application/json" \
//	  "https://console.workbrew.com/workspaces/{workspace}/licenses.json"
func (s *Licenses) ListLicenses(ctx context.Context) (*LicensesResponse, *resty.Response, error) {
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

// ListLicensesCSV retrieves all licenses in CSV format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/licenses.csv
//
// Example cURL:
//
//	curl -X GET \
//	  -H "Authorization: Bearer YOUR_API_KEY" \
//	  -H "X-Workbrew-API-Version: v0" \
//	  -H "Accept: text/csv" \
//	  "https://console.workbrew.com/workspaces/{workspace}/licenses.csv"
func (s *Licenses) ListLicensesCSV(ctx context.Context) ([]byte, *resty.Response, error) {
	resp, csvData, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.TextCSV).
		GetBytes(constants.EndpointLicensesCSV)
	if err != nil {
		return nil, resp, err
	}

	return csvData, resp, nil
}
