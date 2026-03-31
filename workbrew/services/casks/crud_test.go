package casks

import (
	"context"
	"testing"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/casks/mocks"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func setupMockClient(t *testing.T) (*Casks, string) {
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

	return NewCasks(httpClient), baseURL
}

func TestListV0_Success(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.CasksMock{}
	mockHandler.RegisterMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	result, _, err := service.ListV0(ctx)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, *result)

	// Verify first cask (logi-options+)
	cask := (*result)[0]
	assert.Equal(t, "logi-options+", cask.Name)
	assert.NotNil(t, cask.DisplayName)
	assert.Equal(t, "Logitech Options+", *cask.DisplayName)
	assert.Contains(t, cask.Devices, "TC6R2DHVHG")
	assert.Contains(t, cask.Devices, "1234567890")
	assert.True(t, cask.Outdated)
	assert.NotNil(t, cask.Deprecated)
	assert.Equal(t, "", *cask.Deprecated)
	assert.NotNil(t, cask.HomebrewCaskVersion)
	assert.Equal(t, "8.11.0_1", *cask.HomebrewCaskVersion)

	// Verify second cask (1password)
	cask2 := (*result)[1]
	assert.Equal(t, "1password", cask2.Name)
	assert.NotNil(t, cask2.DisplayName)
	assert.Equal(t, "1Password", *cask2.DisplayName)
	assert.True(t, cask2.Outdated)

	// Verify third cask (workbrew) - no display_name
	cask3 := (*result)[2]
	assert.Equal(t, "workbrew/private/workbrew", cask3.Name)
	assert.Nil(t, cask3.DisplayName)
	assert.False(t, cask3.Outdated)

	assert.Len(t, *result, 3)
	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestListV0_Unauthorized(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.CasksMock{}
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
	mockHandler := &mocks.CasksMock{}
	mockHandler.RegisterMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	csvData, _, err := service.ListCSVV0(ctx)

	require.NoError(t, err)
	require.NotNil(t, csvData)
	assert.NotEmpty(t, csvData)

	// Verify CSV headers and content
	csvString := string(csvData)
	assert.Contains(t, csvString, "name,devices,outdated,deprecated,homebrew_cask_version")
	assert.Contains(t, csvString, "1password")
	assert.Contains(t, csvString, "logi-options+")
	assert.Contains(t, csvString, "workbrew/private/workbrew")

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestListCSVV0_Unauthorized(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.CasksMock{}
	mockHandler.RegisterErrorMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	result, _, err := service.ListCSVV0(ctx)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "401")

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}
