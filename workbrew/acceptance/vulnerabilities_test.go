package acceptance

import (
	"testing"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/vulnerabilities"
	"github.com/stretchr/testify/assert"
)

// TestAcceptance_Vulnerabilities_ListVulnerabilities tests retrieving vulnerabilities in JSON format
func TestAcceptance_Vulnerabilities_ListVulnerabilities(t *testing.T) {
	RequireClient(t)

	RateLimitedTest(t, func(t *testing.T) {
		ctx, cancel := NewContext()
		defer cancel()

		service := vulnerabilities.NewVulnerabilities(Client)

		LogTestStage(t, "🔒 List Vulnerabilities", "Testing ListVulnerabilities")

		result, resp, err := service.ListVulnerabilities(ctx)
		AssertNoError(t, err, "ListVulnerabilities should not return an error")
		AssertNotNil(t, result, "ListVulnerabilities result should not be nil")
		AssertNotNil(t, resp, "Response should not be nil")
		assert.Equal(t, 200, resp.StatusCode, "Status code should be 200")

		// Validate response structure
		assert.NotNil(t, result, "Vulnerabilities list should not be nil")

		vulnCount := len(*result)
		LogTestSuccess(t, "Retrieved %d vulnerabilities", vulnCount)

		// If vulnerabilities exist, validate structure
		if vulnCount > 0 {
			vuln := (*result)[0]
			assert.NotEmpty(t, vuln.Formula, "Formula should not be empty")

			vulnCount := len(vuln.Vulnerabilities)
			LogResponse(t, "  Sample vulnerability - Formula: %s, Vulnerabilities: %d, Devices: %d, Supported: %v",
				vuln.Formula,
				vulnCount,
				len(vuln.OutdatedDevices),
				vuln.Supported)
		}
	})
}

// TestAcceptance_Vulnerabilities_ListVulnerabilitiesCSV tests retrieving vulnerabilities in CSV format
func TestAcceptance_Vulnerabilities_ListVulnerabilitiesCSV(t *testing.T) {
	RequireClient(t)

	RateLimitedTest(t, func(t *testing.T) {
		ctx, cancel := NewContext()
		defer cancel()

		service := vulnerabilities.NewVulnerabilities(Client)

		LogTestStage(t, "📊 List Vulnerabilities CSV", "Testing ListVulnerabilitiesCSV")

		csvData, resp, err := service.ListVulnerabilitiesCSV(ctx)
		AssertNoError(t, err, "ListVulnerabilitiesCSV should not return an error")
		AssertNotNil(t, csvData, "CSV data should not be nil")
		AssertNotNil(t, resp, "Response should not be nil")
		assert.Equal(t, 200, resp.StatusCode, "Status code should be 200")

		// Validate CSV data
		assert.Greater(t, len(csvData), 0, "CSV data should not be empty")

		// CSV should start with headers
		csvString := string(csvData)
		assert.Contains(t, csvString, "cve", "CSV should contain cve header")

		LogTestSuccess(t, "Retrieved CSV data with %d bytes", len(csvData))
		LogResponse(t, "  CSV preview: %s", csvString[:min(100, len(csvString))])
	})
}
