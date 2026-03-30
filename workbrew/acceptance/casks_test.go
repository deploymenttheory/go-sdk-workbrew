package acceptance

import (
	"testing"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/casks"
	"github.com/stretchr/testify/assert"
)

// TestAcceptance_Casks_ListCasks tests retrieving casks in JSON format
func TestAcceptance_Casks_ListCasks(t *testing.T) {
	RequireClient(t)

	RateLimitedTest(t, func(t *testing.T) {
		ctx, cancel := NewContext()
		defer cancel()

		service := casks.NewCasks(Client)

		LogTestStage(t, "📦 List Casks", "Testing ListCasks")

		result, resp, err := service.ListCasks(ctx)
		AssertNoError(t, err, "ListCasks should not return an error")
		AssertNotNil(t, result, "ListCasks result should not be nil")
		AssertNotNil(t, resp, "Response should not be nil")
		assert.Equal(t, 200, resp.StatusCode, "Status code should be 200")

		// Validate response structure
		assert.NotNil(t, result, "Casks list should not be nil")

		caskCount := len(*result)
		LogTestSuccess(t, "Retrieved %d casks", caskCount)

		// If casks exist, validate structure
		if caskCount > 0 {
			cask := (*result)[0]
			assert.NotEmpty(t, cask.Name, "Cask name should not be empty")

			displayName := "N/A"
			if cask.DisplayName != nil {
				displayName = *cask.DisplayName
			}

			LogResponse(t, "  Sample cask - Name: %s, Display: %s, Devices: %d, Outdated: %v",
				cask.Name,
				displayName,
				len(cask.Devices),
				cask.Outdated)
		}
	})
}

// TestAcceptance_Casks_ListCasksCSV tests retrieving casks in CSV format
func TestAcceptance_Casks_ListCasksCSV(t *testing.T) {
	RequireClient(t)

	RateLimitedTest(t, func(t *testing.T) {
		ctx, cancel := NewContext()
		defer cancel()

		service := casks.NewCasks(Client)

		LogTestStage(t, "📊 List Casks CSV", "Testing ListCasksCSV")

		csvData, resp, err := service.ListCasksCSV(ctx)
		AssertNoError(t, err, "ListCasksCSV should not return an error")
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

// TestAcceptance_Casks_ValidateFields tests cask field validation
func TestAcceptance_Casks_ValidateFields(t *testing.T) {
	RequireClient(t)

	RateLimitedTest(t, func(t *testing.T) {
		ctx, cancel := NewContext()
		defer cancel()

		service := casks.NewCasks(Client)

		LogTestStage(t, "✅ Validate Fields", "Testing cask field validation")

		result, resp, err := service.ListCasks(ctx)
		AssertNoError(t, err, "ListCasks should not return an error")
		AssertNotNil(t, result, "ListCasks result should not be nil")
		assert.Equal(t, 200, resp.StatusCode, "Status code should be 200")

		if len(*result) == 0 {
			t.Skip("No casks in workspace, skipping field validation")
		}

		// Validate first cask has required fields
		cask := (*result)[0]

		// Required fields
		assert.NotEmpty(t, cask.Name, "Cask name is required")

		// Array fields
		assert.IsType(t, []string{}, cask.Devices, "Devices should be array of strings")

		// Boolean fields
		assert.IsType(t, false, cask.Outdated, "Outdated should be boolean")

		// Optional pointer fields
		if cask.DisplayName != nil {
			assert.NotEmpty(t, *cask.DisplayName, "DisplayName should not be empty if set")
		}

		if cask.HomebrewCaskVersion != nil {
			assert.NotEmpty(t, *cask.HomebrewCaskVersion, "HomebrewCaskVersion should not be empty if set")
		}

		LogTestSuccess(t, "Cask fields validated successfully")
	})
}
