package acceptance

import (
	"testing"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/brewconfigurations"
	"github.com/stretchr/testify/assert"
)

// TestAcceptance_BrewConfigurations_ListBrewConfigurations tests retrieving brew configurations in JSON format
func TestAcceptance_BrewConfigurations_ListBrewConfigurations(t *testing.T) {
	RequireClient(t)

	RateLimitedTest(t, func(t *testing.T) {
		ctx, cancel := NewContext()
		defer cancel()

		service := brewconfigurations.NewBrewConfigurations(Client)

		LogTestStage(t, "⚙️  List Configurations", "Testing ListBrewConfigurations")

		result, resp, err := service.ListBrewConfigurations(ctx)
		AssertNoError(t, err, "ListBrewConfigurations should not return an error")
		AssertNotNil(t, result, "ListBrewConfigurations result should not be nil")
		AssertNotNil(t, resp, "Response should not be nil")
		assert.Equal(t, 200, resp.StatusCode, "Status code should be 200")

		// Validate response structure
		assert.NotNil(t, result, "Brew configurations list should not be nil")

		configCount := len(*result)
		LogTestSuccess(t, "Retrieved %d brew configurations", configCount)

		// If configurations exist, validate structure
		if configCount > 0 {
			config := (*result)[0]
			assert.NotEmpty(t, config.Key, "Configuration key should not be empty")

			LogResponse(t, "  Sample config - Key: %s, Value: %s, Group: %s",
				config.Key,
				config.Value,
				config.DeviceGroup)
		}
	})
}

// TestAcceptance_BrewConfigurations_ListBrewConfigurationsCSV tests retrieving brew configurations in CSV format
func TestAcceptance_BrewConfigurations_ListBrewConfigurationsCSV(t *testing.T) {
	RequireClient(t)

	RateLimitedTest(t, func(t *testing.T) {
		ctx, cancel := NewContext()
		defer cancel()

		service := brewconfigurations.NewBrewConfigurations(Client)

		LogTestStage(t, "📊 List Configurations CSV", "Testing ListBrewConfigurationsCSV")

		csvData, resp, err := service.ListBrewConfigurationsCSV(ctx)
		AssertNoError(t, err, "ListBrewConfigurationsCSV should not return an error")
		AssertNotNil(t, csvData, "CSV data should not be nil")
		AssertNotNil(t, resp, "Response should not be nil")
		assert.Equal(t, 200, resp.StatusCode, "Status code should be 200")

		// Validate CSV data
		assert.Greater(t, len(csvData), 0, "CSV data should not be empty")

		// CSV should start with headers
		csvString := string(csvData)
		assert.Contains(t, csvString, "key", "CSV should contain key header")

		LogTestSuccess(t, "Retrieved CSV data with %d bytes", len(csvData))
		LogResponse(t, "  CSV preview: %s", csvString[:min(100, len(csvString))])
	})
}
