package acceptance

import (
	"testing"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/brewcommands"
	"github.com/stretchr/testify/assert"
)

// TestAcceptance_BrewCommands_ListBrewCommands tests retrieving brew commands in JSON format
func TestAcceptance_BrewCommands_ListBrewCommands(t *testing.T) {
	RequireClient(t)

	RateLimitedTest(t, func(t *testing.T) {
		ctx, cancel := NewContext()
		defer cancel()

		service := brewcommands.NewBrewCommands(Client)

		LogTestStage(t, "⚡ List Commands", "Testing ListBrewCommands")

		result, resp, err := service.ListBrewCommands(ctx)
		AssertNoError(t, err, "ListBrewCommands should not return an error")
		AssertNotNil(t, result, "ListBrewCommands result should not be nil")
		AssertNotNil(t, resp, "Response should not be nil")
		assert.Equal(t, 200, resp.StatusCode, "Status code should be 200")

		// Validate response structure
		assert.NotNil(t, result, "Brew commands list should not be nil")

		commandCount := len(*result)
		LogTestSuccess(t, "Retrieved %d brew commands", commandCount)

		// If commands exist, validate structure
		if commandCount > 0 {
			cmd := (*result)[0]
			assert.NotEmpty(t, cmd.Label, "Command label should not be empty")

			LogResponse(t, "  Sample command - Label: %s, Command: %s, Devices: %d, Runs: %d",
				cmd.Label,
				cmd.Command,
				len(cmd.Devices),
				cmd.RunCount)
		}
	})
}

// TestAcceptance_BrewCommands_ListBrewCommandsCSV tests retrieving brew commands in CSV format
func TestAcceptance_BrewCommands_ListBrewCommandsCSV(t *testing.T) {
	RequireClient(t)

	RateLimitedTest(t, func(t *testing.T) {
		ctx, cancel := NewContext()
		defer cancel()

		service := brewcommands.NewBrewCommands(Client)

		LogTestStage(t, "📊 List Commands CSV", "Testing ListBrewCommandsCSV")

		csvData, resp, err := service.ListBrewCommandsCSV(ctx)
		AssertNoError(t, err, "ListBrewCommandsCSV should not return an error")
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

// TestAcceptance_BrewCommands_CreateCommand tests creating a brew command
// NOTE: This test creates actual commands on devices - use with caution
func TestAcceptance_BrewCommands_CreateCommand(t *testing.T) {
	RequireClient(t)

	// Skip if no known device configured (safer to not run arbitrary commands)
	if Config.KnownDeviceSerial == "" {
		t.Skip("WORKBREW_TEST_DEVICE_SERIAL not set, skipping command creation test")
	}

	RateLimitedTest(t, func(t *testing.T) {
		ctx, cancel := NewContext()
		defer cancel()

		service := brewcommands.NewBrewCommands(Client)

		LogTestStage(t, "➕ Create Command", "Creating test command: brew --version")

		// Use a safe, read-only command
		// Arguments should be just the command without "brew" prefix
		createRequest := &brewcommands.CreateBrewCommandRequest{
			Arguments: "--version",
			DeviceIDs: &Config.KnownDeviceSerial,
		}

		result, resp, err := service.CreateBrewCommand(ctx, createRequest)
		AssertNoError(t, err, "CreateBrewCommand should not return an error")
		AssertNotNil(t, result, "CreateBrewCommand result should not be nil")
		assert.Equal(t, 201, resp.StatusCode, "Status code should be 201 for creation")

		// Validate response
		assert.NotEmpty(t, result.Message, "Response message should not be empty")
		LogTestSuccess(t, "Command created: %s", result.Message)

		LogTestWarning(t, "Command execution may take time to complete on device")
	})
}
