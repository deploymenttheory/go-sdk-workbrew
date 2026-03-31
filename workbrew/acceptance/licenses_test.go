package acceptance

import (
	"testing"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/licenses"
	"github.com/stretchr/testify/assert"
)

// TestAcceptance_Licenses_ListLicenses tests retrieving licenses in JSON format
func TestAcceptance_Licenses_ListLicenses(t *testing.T) {
	RequireClient(t)

	RateLimitedTest(t, func(t *testing.T) {
		ctx, cancel := NewContext()
		defer cancel()

		service := licenses.NewLicenses(Client)

		LogTestStage(t, "📜 List Licenses", "Testing ListV0")

		result, resp, err := service.ListV0(ctx)
		AssertNoError(t, err, "ListV0 should not return an error")
		AssertNotNil(t, result, "ListV0 result should not be nil")
		AssertNotNil(t, resp, "Response should not be nil")
		assert.Equal(t, 200, resp.StatusCode, "Status code should be 200")

		// Validate response structure
		assert.NotNil(t, result, "Licenses list should not be nil")

		licenseCount := len(*result)
		LogTestSuccess(t, "Retrieved %d licenses", licenseCount)

		// If licenses exist, validate structure
		if licenseCount > 0 {
			license := (*result)[0]
			assert.NotEmpty(t, license.Name, "License name should not be empty")

			LogResponse(t, "  Sample license - Name: %s, Devices: %d, Formulae: %d",
				license.Name,
				license.DeviceCount,
				license.FormulaCount)
		}
	})
}

// TestAcceptance_Licenses_ListLicensesCSV tests retrieving licenses in CSV format
func TestAcceptance_Licenses_ListLicensesCSV(t *testing.T) {
	RequireClient(t)

	RateLimitedTest(t, func(t *testing.T) {
		ctx, cancel := NewContext()
		defer cancel()

		service := licenses.NewLicenses(Client)

		LogTestStage(t, "📊 List Licenses CSV", "Testing ListCSVV0")

		csvData, resp, err := service.ListCSVV0(ctx)
		AssertNoError(t, err, "ListCSVV0 should not return an error")
		AssertNotNil(t, csvData, "CSV data should not be nil")
		AssertNotNil(t, resp, "Response should not be nil")
		assert.Equal(t, 200, resp.StatusCode, "Status code should be 200")

		// Validate CSV data
		assert.Greater(t, len(csvData), 0, "CSV data should not be empty")

		// CSV should start with headers
		csvString := string(csvData)
		assert.Contains(t, csvString, "license", "CSV should contain license header")

		LogTestSuccess(t, "Retrieved CSV data with %d bytes", len(csvData))
		LogResponse(t, "  CSV preview: %s", csvString[:min(100, len(csvString))])
	})
}
