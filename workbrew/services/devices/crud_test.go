package devices

import (
	"context"
	"testing"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/devices/mocks"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// setupMockClient creates a client with httpmock enabled
func setupMockClient(t *testing.T) (*Devices, string) {
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

	// Create devices service
	return NewDevices(httpClient), baseURL
}

func TestListV0_Success(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.DevicesMock{}
	mockHandler.RegisterMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	result, _, err := service.ListV0(ctx)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, *result)

	// Verify the first device
	device := (*result)[0]
	assert.Equal(t, "TC6R2DHVHG", device.SerialNumber)
	assert.NotNil(t, device.MDMUserOrDeviceName)
	assert.Equal(t, "Mike's MacBook Pro", *device.MDMUserOrDeviceName)
	assert.Equal(t, "MacBook Pro", device.DeviceType)
	assert.Equal(t, "macOS 14.0 (23A344)", device.OSVersion)
	assert.Equal(t, "/opt/homebrew", device.HomebrewPrefix)
	assert.Equal(t, "4.1.15-24-g5e78ba3", device.HomebrewVersion)
	assert.Equal(t, "0.2.1", device.WorkbrewVersion)
	assert.Equal(t, 9, device.FormulaeCount)
	assert.Equal(t, 3, device.CasksCount)

	// Verify groups
	assert.NotEmpty(t, device.Groups)
	assert.Contains(t, device.Groups, "OSX 14")

	// Verify timestamps
	assert.NotNil(t, device.LastSeenAt)
	assert.NotNil(t, device.CommandLastRunAt)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestListV0_Unauthorized(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.DevicesMock{}
	mockHandler.RegisterErrorMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	result, _, err := service.ListV0(ctx)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "401")

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestListCSVV0_Success(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.DevicesMock{}
	mockHandler.RegisterMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	result, _, err := service.ListCSVV0(ctx)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, result)

	// Verify CSV content
	csvContent := string(result)
	assert.Contains(t, csvContent, "serial_number")
	assert.Contains(t, csvContent, "TC6R2DHVHG")
	assert.Contains(t, csvContent, "MacBook Pro")

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestListCSVV0_Forbidden(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.DevicesMock{}
	mockHandler.RegisterErrorMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	result, _, err := service.ListCSVV0(ctx)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "403")

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}
