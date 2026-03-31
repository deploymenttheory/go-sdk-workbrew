package acceptance

import (
	"testing"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/analytics"
	"github.com/stretchr/testify/assert"
)

// TestAcceptance_Analytics_ListAnalytics tests retrieving analytics in JSON format
func TestAcceptance_Analytics_ListAnalytics(t *testing.T) {
	RequireClient(t)

	RateLimitedTest(t, func(t *testing.T) {
		ctx, cancel := NewContext()
		defer cancel()

		service := analytics.NewAnalytics(Client)

		LogTestStage(t, "📈 List Analytics", "Testing ListAnalytics")

		result, resp, err := service.ListV0(ctx)
		AssertNoError(t, err, "ListAnalytics should not return an error")
		AssertNotNil(t, result, "ListAnalytics result should not be nil")
		AssertNotNil(t, resp, "Response should not be nil")
		assert.Equal(t, 200, resp.StatusCode, "Status code should be 200")

		// Validate response structure
		assert.NotNil(t, result, "Analytics list should not be nil")

		analyticsCount := len(*result)
		LogTestSuccess(t, "Retrieved %d analytics entries", analyticsCount)

		// If analytics exist, validate structure
		if analyticsCount > 0 {
			analytic := (*result)[0]
			assert.NotEmpty(t, analytic.Device, "Device should not be empty")

			LogResponse(t, "  Sample analytic - Device: %s, Command: %s, Count: %d",
				analytic.Device,
				analytic.Command,
				analytic.Count)
		}
	})
}

// TestAcceptance_Analytics_ListAnalyticsCSV tests retrieving analytics in CSV format
func TestAcceptance_Analytics_ListAnalyticsCSV(t *testing.T) {
	RequireClient(t)

	RateLimitedTest(t, func(t *testing.T) {
		ctx, cancel := NewContext()
		defer cancel()

		service := analytics.NewAnalytics(Client)

		LogTestStage(t, "📊 List Analytics CSV", "Testing ListAnalyticsCSV")

		csvData, resp, err := service.ListCSVV0(ctx)
		AssertNoError(t, err, "ListAnalyticsCSV should not return an error")
		AssertNotNil(t, csvData, "CSV data should not be nil")
		AssertNotNil(t, resp, "Response should not be nil")
		assert.Equal(t, 200, resp.StatusCode, "Status code should be 200")

		// Validate CSV data
		assert.Greater(t, len(csvData), 0, "CSV data should not be empty")

		// CSV should start with headers
		csvString := string(csvData)
		assert.Contains(t, csvString, "device_serial", "CSV should contain device_serial header")

		LogTestSuccess(t, "Retrieved CSV data with %d bytes", len(csvData))
		LogResponse(t, "  CSV preview: %s", csvString[:min(100, len(csvString))])
	})
}
