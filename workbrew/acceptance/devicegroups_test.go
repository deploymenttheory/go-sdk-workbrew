package acceptance

import (
	"testing"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/devicegroups"
	"github.com/stretchr/testify/assert"
)

// TestAcceptance_DeviceGroups_ListDeviceGroups tests retrieving device groups in JSON format
func TestAcceptance_DeviceGroups_ListDeviceGroups(t *testing.T) {
	RequireClient(t)

	RateLimitedTest(t, func(t *testing.T) {
		ctx, cancel := NewContext()
		defer cancel()

		service := devicegroups.NewDeviceGroups(Client)

		LogTestStage(t, "👥 List Groups", "Testing ListV0")

		result, resp, err := service.ListV0(ctx)
		AssertNoError(t, err, "ListV0 should not return an error")
		AssertNotNil(t, result, "ListV0 result should not be nil")
		AssertNotNil(t, resp, "Response should not be nil")
		assert.Equal(t, 200, resp.StatusCode, "Status code should be 200")

		// Validate response structure
		assert.NotNil(t, result, "Device groups list should not be nil")

		groupCount := len(*result)
		LogTestSuccess(t, "Retrieved %d device groups", groupCount)

		// If groups exist, validate structure
		if groupCount > 0 {
			group := (*result)[0]
			assert.NotEmpty(t, group.Name, "Group name should not be empty")

			LogResponse(t, "  Sample group - Name: %s, Devices: %d",
				group.Name,
				len(group.Devices))
		}
	})
}

// TestAcceptance_DeviceGroups_ListDeviceGroupsCSV tests retrieving device groups in CSV format
func TestAcceptance_DeviceGroups_ListDeviceGroupsCSV(t *testing.T) {
	RequireClient(t)

	RateLimitedTest(t, func(t *testing.T) {
		ctx, cancel := NewContext()
		defer cancel()

		service := devicegroups.NewDeviceGroups(Client)

		LogTestStage(t, "📊 List Groups CSV", "Testing ListCSVV0")

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

// TestAcceptance_DeviceGroups_ValidateFields tests device group field validation
func TestAcceptance_DeviceGroups_ValidateFields(t *testing.T) {
	RequireClient(t)

	RateLimitedTest(t, func(t *testing.T) {
		ctx, cancel := NewContext()
		defer cancel()

		service := devicegroups.NewDeviceGroups(Client)

		LogTestStage(t, "✅ Validate Fields", "Testing device group field validation")

		result, resp, err := service.ListV0(ctx)
		AssertNoError(t, err, "ListV0 should not return an error")
		AssertNotNil(t, result, "ListV0 result should not be nil")
		assert.Equal(t, 200, resp.StatusCode, "Status code should be 200")

		if len(*result) == 0 {
			t.Skip("No device groups in workspace, skipping field validation")
		}

		// Validate first group has required fields
		group := (*result)[0]

		// Required fields
		assert.NotEmpty(t, group.Name, "Group name is required")
		assert.NotEmpty(t, group.ID, "Group ID is required")

		// Array fields
		assert.IsType(t, []string{}, group.Devices, "Devices should be array of strings")
		assert.GreaterOrEqual(t, len(group.Devices), 0, "Devices count should be non-negative")

		LogTestSuccess(t, "Device group fields validated successfully")
	})
}
