package acceptance

import (
	"testing"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/devices"
	"github.com/stretchr/testify/assert"
)

// TestAcceptance_Devices_ListDevices tests retrieving devices in JSON format
func TestAcceptance_Devices_ListDevices(t *testing.T) {
	RequireClient(t)

	RateLimitedTest(t, func(t *testing.T) {
		ctx, cancel := NewContext()
		defer cancel()

		service := devices.NewDevices(Client)

		LogTestStage(t, "🖥️  List Devices", "Testing ListV0")

		result, resp, err := service.ListV0(ctx)
		AssertNoError(t, err, "ListV0 should not return an error")
		AssertNotNil(t, result, "ListV0 result should not be nil")
		AssertNotNil(t, resp, "Response should not be nil")
		assert.Equal(t, 200, resp.StatusCode, "Status code should be 200")

		// Validate response structure
		assert.NotNil(t, result, "Devices list should not be nil")

		deviceCount := len(*result)
		LogTestSuccess(t, "Retrieved %d devices", deviceCount)

		// If devices exist, validate structure
		if deviceCount > 0 {
			device := (*result)[0]
			assert.NotEmpty(t, device.SerialNumber, "Device serial number should not be empty")

			LogResponse(t, "  Sample device - Serial: %s, Type: %s, OS: %s",
				device.SerialNumber,
				device.DeviceType,
				device.OSVersion)
		}
	})
}

// TestAcceptance_Devices_ListDevicesCSV tests retrieving devices in CSV format
func TestAcceptance_Devices_ListDevicesCSV(t *testing.T) {
	RequireClient(t)

	RateLimitedTest(t, func(t *testing.T) {
		ctx, cancel := NewContext()
		defer cancel()

		service := devices.NewDevices(Client)

		LogTestStage(t, "📊 List Devices CSV", "Testing ListCSVV0")

		csvData, resp, err := service.ListCSVV0(ctx)
		AssertNoError(t, err, "ListCSVV0 should not return an error")
		AssertNotNil(t, csvData, "CSV data should not be nil")
		AssertNotNil(t, resp, "Response should not be nil")
		assert.Equal(t, 200, resp.StatusCode, "Status code should be 200")

		// Validate CSV data
		assert.Greater(t, len(csvData), 0, "CSV data should not be empty")

		// CSV should start with headers
		csvString := string(csvData)
		assert.Contains(t, csvString, "serial_number", "CSV should contain serial_number header")

		LogTestSuccess(t, "Retrieved CSV data with %d bytes", len(csvData))
		LogResponse(t, "  CSV preview: %s", csvString[:min(100, len(csvString))])
	})
}

// TestAcceptance_Devices_ListDevices_ValidateFields tests device field validation
func TestAcceptance_Devices_ListDevices_ValidateFields(t *testing.T) {
	RequireClient(t)

	RateLimitedTest(t, func(t *testing.T) {
		ctx, cancel := NewContext()
		defer cancel()

		service := devices.NewDevices(Client)

		LogTestStage(t, "✅ Validate Fields", "Testing device field validation")

		result, resp, err := service.ListV0(ctx)
		AssertNoError(t, err, "ListV0 should not return an error")
		AssertNotNil(t, result, "ListV0 result should not be nil")
		assert.Equal(t, 200, resp.StatusCode, "Status code should be 200")

		if len(*result) == 0 {
			t.Skip("No devices in workspace, skipping field validation")
		}

		// Validate first device has required fields
		device := (*result)[0]

		// Required fields
		assert.NotEmpty(t, device.SerialNumber, "Serial number is required")
		assert.NotEmpty(t, device.DeviceType, "Device type is required")

		// Array fields
		assert.IsType(t, []string{}, device.Groups, "Groups should be array of strings")

		// TimeOrNever fields
		assert.NotEmpty(t, device.LastSeenAt.String(), "LastSeenAt should have a value")

		// String fields
		if device.HomebrewVersion != "" {
			assert.NotEmpty(t, device.HomebrewVersion, "HomebrewVersion should not be empty if set")
		}

		// Integer fields
		assert.GreaterOrEqual(t, device.FormulaeCount, 0, "FormulaeCount should be non-negative")
		assert.GreaterOrEqual(t, device.CasksCount, 0, "CasksCount should be non-negative")

		LogTestSuccess(t, "Device fields validated successfully")
	})
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
