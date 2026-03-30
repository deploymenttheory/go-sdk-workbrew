package acceptance

import (
	"testing"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/formulae"
	"github.com/stretchr/testify/assert"
)

// TestAcceptance_Formulae_ListFormulae tests retrieving formulae in JSON format
func TestAcceptance_Formulae_ListFormulae(t *testing.T) {
	RequireClient(t)

	RateLimitedTest(t, func(t *testing.T) {
		ctx, cancel := NewContext()
		defer cancel()

		service := formulae.NewFormulae(Client)

		LogTestStage(t, "🍺 List Formulae", "Testing ListFormulae")

		result, resp, err := service.ListFormulae(ctx)
		AssertNoError(t, err, "ListFormulae should not return an error")
		AssertNotNil(t, result, "ListFormulae result should not be nil")
		AssertNotNil(t, resp, "Response should not be nil")
		assert.Equal(t, 200, resp.StatusCode, "Status code should be 200")

		// Validate response structure
		assert.NotNil(t, result, "Formulae list should not be nil")

		formulaCount := len(*result)
		LogTestSuccess(t, "Retrieved %d formulae", formulaCount)

		// If formulae exist, validate structure
		if formulaCount > 0 {
			formula := (*result)[0]
			assert.NotEmpty(t, formula.Name, "Formula name should not be empty")

			LogResponse(t, "  Sample formula - Name: %s, Devices: %d, Outdated: %v, Vulnerabilities: %d",
				formula.Name,
				len(formula.Devices),
				formula.Outdated,
				len(formula.Vulnerabilities))
		}
	})
}

// TestAcceptance_Formulae_ListFormulaeCSV tests retrieving formulae in CSV format
func TestAcceptance_Formulae_ListFormulaeCSV(t *testing.T) {
	RequireClient(t)

	RateLimitedTest(t, func(t *testing.T) {
		ctx, cancel := NewContext()
		defer cancel()

		service := formulae.NewFormulae(Client)

		LogTestStage(t, "📊 List Formulae CSV", "Testing ListFormulaeCSV")

		csvData, resp, err := service.ListFormulaeCSV(ctx)
		AssertNoError(t, err, "ListFormulaeCSV should not return an error")
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

// TestAcceptance_Formulae_ValidateFields tests formula field validation
func TestAcceptance_Formulae_ValidateFields(t *testing.T) {
	RequireClient(t)

	RateLimitedTest(t, func(t *testing.T) {
		ctx, cancel := NewContext()
		defer cancel()

		service := formulae.NewFormulae(Client)

		LogTestStage(t, "✅ Validate Fields", "Testing formula field validation")

		result, resp, err := service.ListFormulae(ctx)
		AssertNoError(t, err, "ListFormulae should not return an error")
		AssertNotNil(t, result, "ListFormulae result should not be nil")
		assert.Equal(t, 200, resp.StatusCode, "Status code should be 200")

		if len(*result) == 0 {
			t.Skip("No formulae in workspace, skipping field validation")
		}

		// Validate first formula has required fields
		formula := (*result)[0]

		// Required fields
		assert.NotEmpty(t, formula.Name, "Formula name is required")

		// Array fields
		assert.IsType(t, []string{}, formula.Devices, "Devices should be array of strings")
		assert.IsType(t, []string{}, formula.Vulnerabilities, "Vulnerabilities should be array of strings")

		// Boolean fields
		assert.IsType(t, false, formula.Outdated, "Outdated should be boolean")
		assert.IsType(t, false, formula.InstalledOnRequest, "InstalledOnRequest should be boolean")
		assert.IsType(t, false, formula.InstalledAsDependency, "InstalledAsDependency should be boolean")

		// Optional pointer fields
		if formula.HomebrewCoreVersion != nil {
			assert.NotEmpty(t, *formula.HomebrewCoreVersion, "HomebrewCoreVersion should not be empty if set")
		}

		if formula.License != nil {
			assert.Greater(t, len(*formula.License), 0, "License should not be empty if set")
		}

		LogTestSuccess(t, "Formula fields validated successfully")
	})
}
