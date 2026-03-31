package devicegroups

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/constants"
	"resty.dev/v3"
)

type (
	// DeviceGroups handles communication with the device groups-related methods of the Workbrew API.
	//
	// Workbrew API docs: https://console.workbrew.com/documentation/api
	DeviceGroups struct {
		client client.Client
	}
)

func NewDeviceGroups(client client.Client) *DeviceGroups {
	return &DeviceGroups{client: client}
}

// ListV0 retrieves all device groups in JSON format.
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/device_groups.json
func (s *DeviceGroups) ListV0(ctx context.Context) (*DeviceGroupsResponse, *resty.Response, error) {
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

// ListCSVV0 retrieves all device groups in CSV format.
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/device_groups.csv
func (s *DeviceGroups) ListCSVV0(ctx context.Context) ([]byte, *resty.Response, error) {
	resp, csvData, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.TextCSV).
		GetBytes(constants.EndpointDeviceGroupsCSV)
	if err != nil {
		return nil, resp, err
	}

	return csvData, resp, nil
}
