package brewtaps

import (
	"context"
	"testing"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/brewtaps/mocks"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func setupMockClient(t *testing.T) (*BrewTaps, string) {
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

	return NewBrewTaps(httpClient), baseURL
}

func TestListV0_Success(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.BrewTapsMock{}
	mockHandler.RegisterMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	result, _, err := service.ListV0(ctx)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, *result)

	// Verify first tap (Homebrew/homebrew-core)
	tap := (*result)[0]
	assert.Equal(t, "Homebrew/homebrew-core", tap.Tap)
	assert.Contains(t, tap.Devices, "TC6R2DHVHG")
	assert.Contains(t, tap.Devices, "1234567890")
	assert.Equal(t, 10, tap.FormulaeInstalled)
	assert.Equal(t, 0, tap.CasksInstalled)
	assert.Equal(t, "7388 Formulae", tap.AvailablePackages)

	// Verify multiple taps
	assert.Len(t, *result, 4)

	// Verify second tap (Homebrew/homebrew-cask)
	tap2 := (*result)[1]
	assert.Equal(t, "Homebrew/homebrew-cask", tap2.Tap)
	assert.Equal(t, 0, tap2.FormulaeInstalled)
	assert.Equal(t, 2, tap2.CasksInstalled)
	assert.Contains(t, tap2.AvailablePackages, "Casks")

	// Verify third tap (apple/apple)
	tap3 := (*result)[2]
	assert.Equal(t, "apple/apple", tap3.Tap)
	assert.Equal(t, ">=1 Packages", tap3.AvailablePackages)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestListV0_Unauthorized(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.BrewTapsMock{}
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
	mockHandler := &mocks.BrewTapsMock{}
	mockHandler.RegisterMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	csvData, _, err := service.ListCSVV0(ctx)

	require.NoError(t, err)
	require.NotNil(t, csvData)
	assert.NotEmpty(t, csvData)

	// Verify CSV headers and content
	csvString := string(csvData)
	assert.Contains(t, csvString, "tap,devices,formulae_installed,casks_installed,available_packages")
	assert.Contains(t, csvString, "Homebrew/homebrew-core")
	assert.Contains(t, csvString, "Homebrew/homebrew-cask")
	assert.Contains(t, csvString, "apple/apple")
	assert.Contains(t, csvString, "workbrew/private")

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestListCSVV0_Unauthorized(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.BrewTapsMock{}
	mockHandler.RegisterErrorMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	result, _, err := service.ListCSVV0(ctx)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "401")

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}
