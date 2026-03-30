package acceptance

import (
	"testing"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/events"
	"github.com/stretchr/testify/assert"
)

// TestAcceptance_Events_ListEvents tests retrieving events in JSON format
func TestAcceptance_Events_ListEvents(t *testing.T) {
	RequireClient(t)

	RateLimitedTest(t, func(t *testing.T) {
		ctx, cancel := NewContext()
		defer cancel()

		service := events.NewEvents(Client)

		LogTestStage(t, "📋 List Events", "Testing ListEvents")

		result, resp, err := service.ListEvents(ctx, nil)
		AssertNoError(t, err, "ListEvents should not return an error")
		AssertNotNil(t, result, "ListEvents result should not be nil")
		AssertNotNil(t, resp, "Response should not be nil")
		assert.Equal(t, 200, resp.StatusCode, "Status code should be 200")

		// Validate response structure
		assert.NotNil(t, result, "Events list should not be nil")

		eventCount := len(*result)
		LogTestSuccess(t, "Retrieved %d events", eventCount)

		// If events exist, validate structure
		if eventCount > 0 {
			event := (*result)[0]
			assert.NotEmpty(t, event.ID, "Event ID should not be empty")

			LogResponse(t, "  Sample event - ID: %s, Type: %s, OccurredAt: %v",
				event.ID,
				event.EventType,
				event.OccurredAt)
		}
	})
}

// TestAcceptance_Events_ListEventsCSV tests retrieving events in CSV format
func TestAcceptance_Events_ListEventsCSV(t *testing.T) {
	RequireClient(t)

	RateLimitedTest(t, func(t *testing.T) {
		ctx, cancel := NewContext()
		defer cancel()

		service := events.NewEvents(Client)

		LogTestStage(t, "📊 List Events CSV", "Testing ListEventsCSV")

		csvData, resp, err := service.ListEventsCSV(ctx, nil)
		AssertNoError(t, err, "ListEventsCSV should not return an error")
		AssertNotNil(t, csvData, "CSV data should not be nil")
		AssertNotNil(t, resp, "Response should not be nil")
		assert.Equal(t, 200, resp.StatusCode, "Status code should be 200")

		// Validate CSV data
		assert.Greater(t, len(csvData), 0, "CSV data should not be empty")

		// CSV should start with headers
		csvString := string(csvData)
		assert.Contains(t, csvString, "id", "CSV should contain id header")

		LogTestSuccess(t, "Retrieved CSV data with %d bytes", len(csvData))
		LogResponse(t, "  CSV preview: %s", csvString[:min(100, len(csvString))])
	})
}
