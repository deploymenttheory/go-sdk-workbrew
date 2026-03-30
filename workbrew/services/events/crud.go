package events

import (
	"context"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/constants"
	"resty.dev/v3"
)

// EventsServiceInterface defines the interface for events operations
//
// Workbrew API docs: https://console.workbrew.com/documentation/api
type EventsServiceInterface interface {
	// ListEvents returns a list of audit log events
	//
	// Returns audit log events with IDs, event types, timestamps, actor information, and target details.
	// Supports filtering by actor type (user, system, or all) via query options.
	ListEvents(ctx context.Context, opts *RequestQueryOptions) (*EventsResponse, *resty.Response, error)

	// ListEventsCSV returns audit log events as CSV
	//
	// Returns audit log event data as CSV with columns: id, event_type, occurred_at, actor_id, actor_type, target_id, target_type, target_identifier.
	// Supports filtering by actor type and optional download parameter via query options.
	ListEventsCSV(ctx context.Context, opts *RequestQueryOptions) ([]byte, *resty.Response, error)
}

// Events handles communication with the events
// related methods of the Workbrew API.
//
// Workbrew API docs: https://console.workbrew.com/documentation/api
type Events struct {
	client client.Client
}

// Ensure Events implements EventsServiceInterface
var _ EventsServiceInterface = (*Events)(nil)

// NewEvents creates a new events service
func NewEvents(client client.Client) *Events {
	return &Events{
		client: client,
	}
}

// ListEvents retrieves all events in JSON format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/events.json
//
// Parameters:
//   - opts: Optional query parameters (filter by actor type: user, system, all)
func (s *Events) ListEvents(ctx context.Context, opts *RequestQueryOptions) (*EventsResponse, *resty.Response, error) {
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

// ListEventsCSV retrieves all events in CSV format
// URL: GET https://console.workbrew.com/workspaces/{workspace_name}/events.csv
//
// Parameters:
//   - opts: Optional query parameters (filter by actor type, download flag)
func (s *Events) ListEventsCSV(ctx context.Context, opts *RequestQueryOptions) ([]byte, *resty.Response, error) {
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
