package devices

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/constants"
	"resty.dev/v3"
)

type (
	// Devices handles communication with the devices-related methods of the Workbrew API.
	//
	// Workbrew API docs: https://console.workbrew.com/documentation/api
	Devices struct {
		client client.Client
	}
)

func NewDevices(client client.Client) *Devices {
	return &Devices{client: client}
}

// ListV0 retrieves all devices in JSON format.
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/devices.json
func (s *Devices) ListV0(ctx context.Context) (*DevicesResponse, *resty.Response, error) {
	var result DevicesResponse
	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetResult(&result).
		Get(constants.EndpointDevicesJSON)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// ListCSVV0 retrieves all devices in CSV format.
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/devices.csv
func (s *Devices) ListCSVV0(ctx context.Context) ([]byte, *resty.Response, error) {
	resp, csvData, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.TextCSV).
		GetBytes(constants.EndpointDevicesCSV)
	if err != nil {
		return nil, resp, err
	}

	return csvData, resp, nil
}
