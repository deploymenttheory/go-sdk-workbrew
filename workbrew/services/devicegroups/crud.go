package devicegroups

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/constants"
	"resty.dev/v3"
)

// DeviceGroupsServiceInterface defines the interface for device groups operations
//
// Workbrew API docs: https://console.workbrew.com/documentation/api
type DeviceGroupsServiceInterface interface {
	// ListDeviceGroups returns a list of Device Groups
	//
	// Returns device groups with their IDs, names, and assigned device serial numbers.
	ListDeviceGroups(ctx context.Context) (*DeviceGroupsResponse, *resty.Response, error)

	// ListDeviceGroupsCSV returns a list of Device Groups in CSV format
	//
	// Returns device group data as CSV with columns: id, name, devices.
	ListDeviceGroupsCSV(ctx context.Context) ([]byte, *resty.Response, error)
}

// DeviceGroups handles communication with the device groups
// related methods of the Workbrew API.
//
// Workbrew API docs: https://console.workbrew.com/documentation/api
type DeviceGroups struct {
	client client.Client
}

// Ensure DeviceGroups implements DeviceGroupsServiceInterface
var _ DeviceGroupsServiceInterface = (*DeviceGroups)(nil)

// NewDeviceGroups creates a new device groups service
func NewDeviceGroups(client client.Client) *DeviceGroups {
	return &DeviceGroups{
		client: client,
	}
}

// ListDeviceGroups retrieves all device groups in JSON format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/device_groups.json
func (s *DeviceGroups) ListDeviceGroups(ctx context.Context) (*DeviceGroupsResponse, *resty.Response, error) {
	var result DeviceGroupsResponse
	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetResult(&result).
		Get(constants.EndpointDeviceGroupsJSON)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// ListDeviceGroupsCSV retrieves all device groups in CSV format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/device_groups.csv
func (s *DeviceGroups) ListDeviceGroupsCSV(ctx context.Context) ([]byte, *resty.Response, error) {
	resp, csvData, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.TextCSV).
		GetBytes(constants.EndpointDeviceGroupsCSV)
	if err != nil {
		return nil, resp, err
	}

	return csvData, resp, nil
}
