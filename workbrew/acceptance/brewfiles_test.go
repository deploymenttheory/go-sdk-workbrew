package acceptance

import (
	"fmt"
	"testing"
	"time"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/brewfiles"
	"github.com/stretchr/testify/assert"
)

// TestAcceptance_Brewfiles_ListBrewfiles tests retrieving brewfiles in JSON format
func TestAcceptance_Brewfiles_ListBrewfiles(t *testing.T) {
	RequireClient(t)

	RateLimitedTest(t, func(t *testing.T) {
		ctx, cancel := NewContext()
		defer cancel()

		service := brewfiles.NewBrewfiles(Client)

		LogTestStage(t, "📝 List Brewfiles", "Testing ListV0")

		result, resp, err := service.ListV0(ctx)
		AssertNoError(t, err, "ListV0 should not return an error")
		AssertNotNil(t, result, "ListV0 result should not be nil")
		AssertNotNil(t, resp, "Response should not be nil")
		assert.Equal(t, 200, resp.StatusCode, "Status code should be 200")

		// Validate response structure
		assert.NotNil(t, result, "Brewfiles list should not be nil")

		brewfileCount := len(*result)
		LogTestSuccess(t, "Retrieved %d brewfiles", brewfileCount)

		// If brewfiles exist, validate structure
		if brewfileCount > 0 {
			brewfile := (*result)[0]
			assert.NotEmpty(t, brewfile.Label, "Brewfile label should not be empty")

			LogResponse(t, "  Sample brewfile - Label: %s, Devices: %d, Runs: %d",
				brewfile.Label,
				len(brewfile.Devices),
				brewfile.RunCount)
		}
	})
}

// TestAcceptance_Brewfiles_ListBrewfilesCSV tests retrieving brewfiles in CSV format
func TestAcceptance_Brewfiles_ListBrewfilesCSV(t *testing.T) {
	RequireClient(t)

	RateLimitedTest(t, func(t *testing.T) {
		ctx, cancel := NewContext()
		defer cancel()

		service := brewfiles.NewBrewfiles(Client)

		LogTestStage(t, "📊 List Brewfiles CSV", "Testing ListCSVV0")

		csvData, resp, err := service.ListCSVV0(ctx)
		AssertNoError(t, err, "ListCSVV0 should not return an error")
		AssertNotNil(t, csvData, "CSV data should not be nil")
		AssertNotNil(t, resp, "Response should not be nil")
		assert.Equal(t, 200, resp.StatusCode, "Status code should be 200")

		// Validate CSV data
		assert.Greater(t, len(csvData), 0, "CSV data should not be empty")

		// CSV should start with headers
		csvString := string(csvData)
		assert.Contains(t, csvString, "label", "CSV should contain label header")

		LogTestSuccess(t, "Retrieved CSV data with %d bytes", len(csvData))
		LogResponse(t, "  CSV preview: %s", csvString[:min(100, len(csvString))])
	})
}

// TestAcceptance_Brewfiles_CRUD tests full CRUD operations for brewfiles
func TestAcceptance_Brewfiles_CRUD(t *testing.T) {
	RequireClient(t)

	RateLimitedTest(t, func(t *testing.T) {
		ctx, cancel := NewContext()
		defer cancel()

		service := brewfiles.NewBrewfiles(Client)

		// Generate unique label for test
		testLabel := fmt.Sprintf("test-brewfile-%d", time.Now().Unix())
		testContent := `# Test Brewfile
brew "git"
brew "wget"
`

		// Step 1: Create brewfile
		LogTestStage(t, "➕ Create", "Creating test brewfile: %s", testLabel)

		createRequest := &brewfiles.CreateBrewfileRequest{
			Label:   testLabel,
			Content: testContent,
		}

		createResult, createResp, createErr := service.CreateV0(ctx, createRequest)
		AssertNoError(t, createErr, "CreateV0 should not return an error")
		AssertNotNil(t, createResult, "CreateV0 result should not be nil")
		assert.Equal(t, 201, createResp.StatusCode, "Status code should be 201 for creation")

		LogTestSuccess(t, "Brewfile created: %s", testLabel)

		// Register cleanup
		Cleanup(t, func() {
			LogTestStage(t, "🧹 Cleanup", "Deleting test brewfile: %s", testLabel)
			cleanupCtx, cleanupCancel := NewContext()
			defer cleanupCancel()

			_, _, cleanupErr := service.DeleteByLabelV0(cleanupCtx, testLabel)
			if cleanupErr != nil {
				LogTestWarning(t, "Failed to cleanup brewfile %s: %v", testLabel, cleanupErr)
			} else {
				LogTestSuccess(t, "Cleaned up test brewfile: %s", testLabel)
			}
		})

		// Step 2: List brewfiles and verify it exists
		LogTestStage(t, "📋 Verify", "Verifying brewfile exists in list")

		listResult, listResp, listErr := service.ListV0(ctx)
		AssertNoError(t, listErr, "ListV0 should not return an error")
		assert.Equal(t, 200, listResp.StatusCode, "Status code should be 200")

		found := false
		for _, bf := range *listResult {
			if bf.Label == testLabel {
				found = true
				assert.Equal(t, testContent, bf.Content, "Content should match")
				break
			}
		}
		assert.True(t, found, "Created brewfile should be in the list")
		LogTestSuccess(t, "Brewfile found in list")

		// Step 3: Update brewfile
		LogTestStage(t, "✏️  Update", "Updating brewfile content")

		updatedContent := `# Updated Test Brewfile
brew "git"
brew "wget"
brew "curl"
`
		updateRequest := &brewfiles.UpdateBrewfileRequest{
			Content: updatedContent,
		}

		updateResult, updateResp, updateErr := service.UpdateByLabelV0(ctx, testLabel, updateRequest)
		AssertNoError(t, updateErr, "UpdateByLabelV0 should not return an error")
		AssertNotNil(t, updateResult, "UpdateByLabelV0 result should not be nil")
		assert.Equal(t, 200, updateResp.StatusCode, "Status code should be 200 for update")

		LogTestSuccess(t, "Brewfile updated")

		// Step 4: Verify update
		LogTestStage(t, "🔍 Verify Update", "Verifying brewfile content updated")

		listResultAfterUpdate, _, listErrAfterUpdate := service.ListV0(ctx)
		AssertNoError(t, listErrAfterUpdate, "ListV0 should not return an error")

		foundUpdated := false
		for _, bf := range *listResultAfterUpdate {
			if bf.Label == testLabel {
				foundUpdated = true
				assert.Equal(t, updatedContent, bf.Content, "Content should be updated")
				break
			}
		}
		assert.True(t, foundUpdated, "Updated brewfile should be in the list")
		LogTestSuccess(t, "Brewfile content verified")

		// Step 5: Delete brewfile (will also be cleaned up, but test explicit delete)
		LogTestStage(t, "🗑️  Delete", "Deleting brewfile")

		deleteResult, deleteResp, deleteErr := service.DeleteByLabelV0(ctx, testLabel)
		AssertNoError(t, deleteErr, "DeleteByLabelV0 should not return an error")
		AssertNotNil(t, deleteResult, "DeleteByLabelV0 result should not be nil")
		assert.Equal(t, 200, deleteResp.StatusCode, "Status code should be 200 for deletion")

		LogTestSuccess(t, "Brewfile deleted")

		// Step 6: Verify deletion
		LogTestStage(t, "✅ Verify Deletion", "Verifying brewfile removed")

		listResultAfterDelete, _, listErrAfterDelete := service.ListV0(ctx)
		AssertNoError(t, listErrAfterDelete, "ListV0 should not return an error")

		foundAfterDelete := false
		for _, bf := range *listResultAfterDelete {
			if bf.Label == testLabel {
				foundAfterDelete = true
				break
			}
		}
		assert.False(t, foundAfterDelete, "Deleted brewfile should not be in the list")
		LogTestSuccess(t, "Brewfile deletion verified")
	})
}

// TestAcceptance_Brewfiles_ListBrewfileRuns tests retrieving brewfile runs
func TestAcceptance_Brewfiles_ListBrewfileRuns(t *testing.T) {
	RequireClient(t)

	// Skip if no known brewfile name configured
	if Config.KnownBrewfileName == "" {
		t.Skip("WORKBREW_TEST_BREWFILE_NAME not set, skipping brewfile runs test")
	}

	RateLimitedTest(t, func(t *testing.T) {
		ctx, cancel := NewContext()
		defer cancel()

		service := brewfiles.NewBrewfiles(Client)

		LogTestStage(t, "🏃 List Runs", "Testing ListRunsByLabelV0 for: %s", Config.KnownBrewfileName)

		result, resp, err := service.ListRunsByLabelV0(ctx, Config.KnownBrewfileName)
		AssertNoError(t, err, "ListRunsByLabelV0 should not return an error")
		AssertNotNil(t, result, "ListRunsByLabelV0 result should not be nil")
		AssertNotNil(t, resp, "Response should not be nil")
		assert.Equal(t, 200, resp.StatusCode, "Status code should be 200")

		runCount := len(*result)
		LogTestSuccess(t, "Retrieved %d brewfile runs", runCount)

		// If runs exist, validate structure
		if runCount > 0 {
			run := (*result)[0]
			assert.NotEmpty(t, run.Label, "Run label should not be empty")

			LogResponse(t, "  Sample run - Label: %s, Device: %s, Success: %v",
				run.Label,
				run.Device,
				run.Success)
		}
	})
}
