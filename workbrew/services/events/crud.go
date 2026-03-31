package events

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/constants"
	"resty.dev/v3"
)

type (
	// Events handles communication with the events-related methods of the Workbrew API.
	//
	// Workbrew API docs: https://console.workbrew.com/documentation/api
	Events struct {
		client client.Client
	}
)

func NewEvents(client client.Client) *Events {
	return &Events{client: client}
}

// ListV0 retrieves all events in JSON format.
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/events.json
//
// Parameters:
//   - opts: Optional query parameters (filter by actor type: user, system, all)
func (s *Events) ListV0(ctx context.Context, opts *RequestQueryOptions) (*EventsResponse, *resty.Response, error) {
	if opts == nil {
		opts = &RequestQueryOptions{}
	}

	var result EventsResponse
	resp, err := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.ApplicationJSON).
		SetHeader("Content-Type", constants.ApplicationJSON).
		SetQueryParam("filter", opts.Filter).
		SetResult(&result).
		Get(constants.EndpointEventsJSON)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// ListCSVV0 retrieves all events in CSV format.
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/events.csv
//
// Parameters:
//   - opts: Optional query parameters (filter by actor type, download flag)
func (s *Events) ListCSVV0(ctx context.Context, opts *RequestQueryOptions) ([]byte, *resty.Response, error) {
	if opts == nil {
		opts = &RequestQueryOptions{}
	}

	rb := s.client.NewRequest(ctx).
		SetHeader("Accept", constants.TextCSV).
		SetQueryParam("filter", opts.Filter)
	if opts.Download {
		rb.SetQueryParam("download", "1")
	}

	resp, csvData, err := rb.GetBytes(constants.EndpointEventsCSV)
	if err != nil {
		return nil, resp, err
	}

	return csvData, resp, nil
}
