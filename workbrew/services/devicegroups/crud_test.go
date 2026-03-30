package devicegroups

import (
	"context"
	"testing"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/devicegroups/mocks"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func setupMockClient(t *testing.T) (*DeviceGroups, string) {
	logger := zap.NewNop()
	baseURL := "https://console.workbrew.com/workspaces/test-workspace"

	httpClient, err := client.NewTransport("test-api-key", "test-workspace",
		client.WithLogger(logger),
		client.WithBaseURL("https://console.workbrew.com"),
	)
	require.NoError(t, err)

	httpmock.ActivateNonDefault(httpClient.GetHTTPClient().Client())
	t.Cleanup(func() {
		httpmock.DeactivateAndReset()
	})

	return NewDeviceGroups(httpClient), baseURL
}

func TestListDeviceGroups_Success(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.DeviceGroupsMock{}
	mockHandler.RegisterMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	result, _, err := service.ListDeviceGroups(ctx)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, *result)

	// Verify first device group (Admin)
	group := (*result)[0]
	assert.Equal(t, "ddba0af6-bd3c-5abf-8311-e62dc6bd9fbc", group.ID)
	assert.Equal(t, "Admin", group.Name)
	assert.Contains(t, group.Devices, "TC6R2DHVHG")
	assert.Len(t, group.Devices, 1)

	// Verify second device group (OSX 14)
	group2 := (*result)[1]
	assert.Equal(t, "377d8aa2-64cd-56a6-8351-6163bcf7dca1", group2.ID)
	assert.Equal(t, "OSX 14", group2.Name)
	assert.Contains(t, group2.Devices, "TC6R2DHVHG")
	assert.Contains(t, group2.Devices, "1234567890")
	assert.Len(t, group2.Devices, 2)

	assert.Len(t, *result, 2)
	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestListDeviceGroups_Unauthorized(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.DeviceGroupsMock{}
	mockHandler.RegisterErrorMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	result, _, err := service.ListDeviceGroups(ctx)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "401")

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestListDeviceGroupsCSV_Success(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.DeviceGroupsMock{}
	mockHandler.RegisterMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	csvData, _, err := service.ListDeviceGroupsCSV(ctx)

	require.NoError(t, err)
	require.NotNil(t, csvData)
	assert.NotEmpty(t, csvData)

	// Verify CSV headers and content
	csvString := string(csvData)
	assert.Contains(t, csvString, "id,name,devices")
	assert.Contains(t, csvString, "ddba0af6-bd3c-5abf-8311-e62dc6bd9fbc")
	assert.Contains(t, csvString, "Admin")
	assert.Contains(t, csvString, "OSX 14")

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestListDeviceGroupsCSV_Unauthorized(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.DeviceGroupsMock{}
	mockHandler.RegisterErrorMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	result, _, err := service.ListDeviceGroupsCSV(ctx)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "401")

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}
