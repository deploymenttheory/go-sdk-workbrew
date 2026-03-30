package devices

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/constants"
	"resty.dev/v3"
)

// DevicesServiceInterface defines the interface for devices operations
//
// Workbrew API docs: https://console.workbrew.com/documentation/api
type DevicesServiceInterface interface {
	// ListDevices returns a list of devices
	//
	// Returns devices with serial numbers, group assignments, MDM names, last seen timestamps, device types,
	// OS versions, Homebrew/Workbrew versions, and installed package counts.
	ListDevices(ctx context.Context) (*DevicesResponse, *resty.Response, error)

	// ListDevicesCSV returns a list of devices in CSV format
	//
	// Returns device data as CSV with columns: serial_number, groups, mdm_user_or_device_name, last_seen_at,
	// command_last_run_at, device_type, os_version, homebrew_prefix, homebrew_version, workbrew_version, formulae_count, casks_count.
	ListDevicesCSV(ctx context.Context) ([]byte, *resty.Response, error)
}

// Devices handles communication with the devices
// related methods of the Workbrew API.
//
// Workbrew API docs: https://console.workbrew.com/documentation/api
type Devices struct {
	client client.Client
}

// Ensure Devices implements DevicesServiceInterface
var _ DevicesServiceInterface = (*Devices)(nil)

// NewDevices creates a new devices service
func NewDevices(client client.Client) *Devices {
	return &Devices{
		client: client,
	}
}

// ListDevices retrieves all devices in JSON format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/devices.json
func (s *Devices) ListDevices(ctx context.Context) (*DevicesResponse, *resty.Response, error) {
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

// ListDevicesCSV retrieves all devices in CSV format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/devices.csv
func (s *Devices) ListDevicesCSV(ctx context.Context) ([]byte, *resty.Response, error) {
	resp, csvData, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.TextCSV).
		GetBytes(constants.EndpointDevicesCSV)
	if err != nil {
		return nil, resp, err
	}

	return csvData, resp, nil
}
