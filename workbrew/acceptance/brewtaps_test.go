package acceptance

import (
	"testing"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/brewtaps"
	"github.com/stretchr/testify/assert"
)

// TestAcceptance_BrewTaps_ListBrewTaps tests retrieving brew taps in JSON format
func TestAcceptance_BrewTaps_ListBrewTaps(t *testing.T) {
	RequireClient(t)

	RateLimitedTest(t, func(t *testing.T) {
		ctx, cancel := NewContext()
		defer cancel()

		service := brewtaps.NewBrewTaps(Client)

		LogTestStage(t, "🚰 List Taps", "Testing ListV0")

		result, resp, err := service.ListV0(ctx)
		AssertNoError(t, err, "ListV0 should not return an error")
		AssertNotNil(t, result, "ListV0 result should not be nil")
		AssertNotNil(t, resp, "Response should not be nil")
		assert.Equal(t, 200, resp.StatusCode, "Status code should be 200")

		// Validate response structure
		assert.NotNil(t, result, "Brew taps list should not be nil")

		tapCount := len(*result)
		LogTestSuccess(t, "Retrieved %d brew taps", tapCount)

		// If taps exist, validate structure
		if tapCount > 0 {
			tap := (*result)[0]
			assert.NotEmpty(t, tap.Tap, "Tap name should not be empty")

			LogResponse(t, "  Sample tap - Tap: %s, Devices: %d, Formulae: %d, Casks: %d",
				tap.Tap,
				len(tap.Devices),
				tap.FormulaeInstalled,
				tap.CasksInstalled)
		}
	})
}

// TestAcceptance_BrewTaps_ListBrewTapsCSV tests retrieving brew taps in CSV format
func TestAcceptance_BrewTaps_ListBrewTapsCSV(t *testing.T) {
	RequireClient(t)

	RateLimitedTest(t, func(t *testing.T) {
		ctx, cancel := NewContext()
		defer cancel()

		service := brewtaps.NewBrewTaps(Client)

		LogTestStage(t, "📊 List Taps CSV", "Testing ListCSVV0")

		csvData, resp, err := service.ListCSVV0(ctx)
		AssertNoError(t, err, "ListCSVV0 should not return an error")
		AssertNotNil(t, csvData, "CSV data should not be nil")
		AssertNotNil(t, resp, "Response should not be nil")
		assert.Equal(t, 200, resp.StatusCode, "Status code should be 200")

		// Validate CSV data
		assert.Greater(t, len(csvData), 0, "CSV data should not be empty")

		// CSV should start with headers
		csvString := string(csvData)
		assert.Contains(t, csvString, "name", "CSV should contain name header")

		LogTestSuccess(t, "Retrieved CSV data with %d bytes", len(csvData))
		LogResponse(t, "  CSV preview: %s", csvString[:min(100, len(csvString))])
	})
}
