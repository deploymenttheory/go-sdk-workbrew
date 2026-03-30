package brewconfigurations

import (
	"context"
	"testing"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/brewconfigurations/mocks"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func setupMockClient(t *testing.T) (*BrewConfigurations, string) {
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

	return NewBrewConfigurations(httpClient), baseURL
}

func TestListBrewConfigurations_Success(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.BrewConfigurationsMock{}
	mockHandler.RegisterMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	result, _, err := service.ListBrewConfigurations(ctx)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, *result)

	// Verify first configuration
	config := (*result)[0]
	assert.Equal(t, "HOMEBREW_DEVELOPER", config.Key)
	assert.Equal(t, "1", config.Value)
	assert.Equal(t, "mikemcquaid", config.LastUpdatedByUser)
	assert.Equal(t, "All Devices", config.DeviceGroup)

	assert.Len(t, *result, 4)
	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestListBrewConfigurations_Unauthorized(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.BrewConfigurationsMock{}
	mockHandler.RegisterErrorMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	result, _, err := service.ListBrewConfigurations(ctx)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "401")
}

func TestListBrewConfigurationsCSV_Success(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.BrewConfigurationsMock{}
	mockHandler.RegisterMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	csvData, _, err := service.ListBrewConfigurationsCSV(ctx)

	require.NoError(t, err)
	require.NotNil(t, csvData)

	csvString := string(csvData)
	assert.Contains(t, csvString, "key,value,last_updated_by_user,device_group")
	assert.Contains(t, csvString, "HOMEBREW_DEVELOPER")
}
