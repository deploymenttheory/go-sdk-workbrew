package events

import (
	"context"
	"testing"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/events/mocks"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// setupMockClient creates a client with httpmock enabled
func setupMockClient(t *testing.T) (*Events, string) {
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

	// Create events service
	return NewEvents(httpClient), baseURL
}

func TestListV0_Success_NoFilter(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.EventsMock{}
	mockHandler.RegisterMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	result, _, err := service.ListV0(ctx, nil) // Test nil options

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, *result)

	// Verify the first event
	event := (*result)[0]
	assert.NotEmpty(t, event.ID)
	assert.Equal(t, "device.created", event.EventType)
	assert.NotNil(t, event.ActorType)
	assert.Equal(t, "User", *event.ActorType)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestListV0_Success_WithFilter(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.EventsMock{}
	mockHandler.RegisterMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	// Test with user filter
	result, _, err := service.ListV0(ctx, &RequestQueryOptions{
		Filter: "user",
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, *result)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestListV0_Success_EmptyOptions(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.EventsMock{}
	mockHandler.RegisterMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	// Test with empty options struct
	result, _, err := service.ListV0(ctx, &RequestQueryOptions{})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, *result)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestListV0_Unauthorized(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.EventsMock{}
	mockHandler.RegisterErrorMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	result, _, err := service.ListV0(ctx, nil)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "401")

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestListCSVV0_Success_NoOptions(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.EventsMock{}
	mockHandler.RegisterMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	csvData, _, err := service.ListCSVV0(ctx, nil)

	require.NoError(t, err)
	require.NotNil(t, csvData)
	assert.NotEmpty(t, csvData)

	// Verify CSV headers
	csvString := string(csvData)
	assert.Contains(t, csvString, "id,event_type,occurred_at")

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestListCSVV0_Success_WithFilter(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.EventsMock{}
	mockHandler.RegisterMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	csvData, _, err := service.ListCSVV0(ctx, &RequestQueryOptions{
		Filter: "system",
	})

	require.NoError(t, err)
	require.NotNil(t, csvData)
	assert.NotEmpty(t, csvData)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestListCSVV0_Success_WithDownload(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.EventsMock{}
	mockHandler.RegisterMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	csvData, _, err := service.ListCSVV0(ctx, &RequestQueryOptions{
		Filter:   "all",
		Download: true,
	})

	require.NoError(t, err)
	require.NotNil(t, csvData)
	assert.NotEmpty(t, csvData)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestListCSVV0_Success_AllFilters(t *testing.T) {
	testCases := []struct {
		name   string
		filter string
	}{
		{"User Filter", "user"},
		{"System Filter", "system"},
		{"All Filter", "all"},
		{"Empty Filter", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service, baseURL := setupMockClient(t)
			mockHandler := &mocks.EventsMock{}
			mockHandler.RegisterMocks(baseURL)
			defer mockHandler.CleanupMockState()

			ctx := context.Background()
			csvData, _, err := service.ListCSVV0(ctx, &RequestQueryOptions{
				Filter: tc.filter,
			})

			require.NoError(t, err)
			require.NotNil(t, csvData)
			assert.NotEmpty(t, csvData)
		})
	}
}

func TestListCSVV0_Unauthorized(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.EventsMock{}
	mockHandler.RegisterErrorMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	result, _, err := service.ListCSVV0(ctx, nil)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "401")

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestRequestQueryOptions_NilSafety(t *testing.T) {
	// This test verifies that nil options are handled gracefully
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.EventsMock{}
	mockHandler.RegisterMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	// Test ListV0 with nil
	events, _, err := service.ListV0(ctx, nil)
	require.NoError(t, err)
	require.NotNil(t, events)

	// Test ListCSVV0 with nil
	csv, _, err := service.ListCSVV0(ctx, nil)
	require.NoError(t, err)
	require.NotNil(t, csv)
}
