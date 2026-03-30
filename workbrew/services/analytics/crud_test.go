package analytics

import (
	"context"
	"testing"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/analytics/mocks"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// setupMockClient creates a client with httpmock enabled
func setupMockClient(t *testing.T) (*Analytics, string) {
	// Create test logger
	logger := zap.NewNop()

	// Create base URL for testing
	baseURL := "https://console.workbrew.com/workspaces/test-workspace"

	// Create HTTP client
	httpClient, err := client.NewTransport("test-api-key", "test-workspace",
		client.WithLogger(logger),
		client.WithBaseURL("https://console.workbrew.com"),
	)
	require.NoError(t, err)

	// Activate httpmock
	httpmock.ActivateNonDefault(httpClient.GetHTTPClient().Client())

	// Setup cleanup
	t.Cleanup(func() {
		httpmock.DeactivateAndReset()
	})

	// Create analytics service
	return NewAnalytics(httpClient), baseURL
}

func TestListAnalytics_Success(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.AnalyticsMock{}
	mockHandler.RegisterMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	result, _, err := service.ListAnalytics(ctx)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, *result)

	// Verify the first analytic entry
	analytic := (*result)[0]
	assert.Equal(t, "TC6R2DHVHG", analytic.Device)
	assert.Equal(t, "brew install curl", analytic.Command)
	assert.Equal(t, 2, analytic.Count)
	assert.NotNil(t, analytic.LastRun)

	// Verify we have multiple analytics
	assert.Len(t, *result, 4, "Should have 4 analytic entries")

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestListAnalytics_VerifyAllFields(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.AnalyticsMock{}
	mockHandler.RegisterMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	result, _, err := service.ListAnalytics(ctx)

	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify each entry has all required fields
	for i, analytic := range *result {
		assert.NotEmpty(t, analytic.Device, "Entry %d should have device", i)
		assert.NotEmpty(t, analytic.Command, "Entry %d should have command", i)
		assert.NotNil(t, analytic.LastRun, "Entry %d should have last_run", i)
		assert.GreaterOrEqual(t, analytic.Count, 1, "Entry %d should have count >= 1", i)
	}
}

func TestListAnalytics_VerifyDevices(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.AnalyticsMock{}
	mockHandler.RegisterMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	result, _, err := service.ListAnalytics(ctx)

	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify we have analytics for multiple devices
	deviceMap := make(map[string]bool)
	for _, analytic := range *result {
		deviceMap[analytic.Device] = true
	}

	assert.Contains(t, deviceMap, "TC6R2DHVHG", "Should have analytics for device TC6R2DHVHG")
	assert.Contains(t, deviceMap, "1234567890", "Should have analytics for device 1234567890")
}

func TestListAnalytics_Unauthorized(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.AnalyticsMock{}
	mockHandler.RegisterErrorMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	result, _, err := service.ListAnalytics(ctx)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "401")

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestListAnalyticsCSV_Success(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.AnalyticsMock{}
	mockHandler.RegisterMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	csvData, _, err := service.ListAnalyticsCSV(ctx)

	require.NoError(t, err)
	require.NotNil(t, csvData)
	assert.NotEmpty(t, csvData)

	// Verify CSV structure
	csvString := string(csvData)
	assert.Contains(t, csvString, "device,command,last_run,count", "CSV should have headers")
	assert.Contains(t, csvString, "TC6R2DHVHG", "CSV should contain device data")
	assert.Contains(t, csvString, "brew install curl", "CSV should contain command data")

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestListAnalyticsCSV_VerifyFormat(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.AnalyticsMock{}
	mockHandler.RegisterMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	csvData, _, err := service.ListAnalyticsCSV(ctx)

	require.NoError(t, err)
	require.NotNil(t, csvData)

	csvString := string(csvData)

	// Verify all expected data is present
	assert.Contains(t, csvString, "TC6R2DHVHG,brew install curl,2024-01-01T12:34:56Z,2")
	assert.Contains(t, csvString, "TC6R2DHVHG,brew install wget,2024-02-03T08:22:33Z,1")
	assert.Contains(t, csvString, "TC6R2DHVHG,brew info curl,2024-04-15T14:45:22Z,1")
	assert.Contains(t, csvString, "1234567890,brew upgrade git,2024-05-10T09:30:00Z,3")
}

func TestListAnalyticsCSV_Unauthorized(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.AnalyticsMock{}
	mockHandler.RegisterErrorMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	result, _, err := service.ListAnalyticsCSV(ctx)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "401")

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestAnalytics_HTTPMockCallCounts(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.AnalyticsMock{}
	mockHandler.RegisterMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	// Make multiple calls
	_, _, err1 := service.ListAnalytics(ctx)
	_, _, err2 := service.ListAnalyticsCSV(ctx)

	require.NoError(t, err1)
	require.NoError(t, err2)

	// Verify httpmock tracked both calls
	assert.Equal(t, 2, httpmock.GetTotalCallCount(), "Should have made 2 HTTP calls")
}

func TestAnalytics_MultipleSequentialCalls(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.AnalyticsMock{}
	mockHandler.RegisterMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	// First call
	result1, _, err1 := service.ListAnalytics(ctx)
	require.NoError(t, err1)
	require.NotNil(t, result1)

	// Second call should return same data
	result2, _, err2 := service.ListAnalytics(ctx)
	require.NoError(t, err2)
	require.NotNil(t, result2)

	// Verify both results have same length
	assert.Equal(t, len(*result1), len(*result2), "Sequential calls should return consistent data")
}
