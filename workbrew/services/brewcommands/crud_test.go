package brewcommands

import (
	"context"
	"testing"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/brewcommands/mocks"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// setupMockClient creates a client with httpmock enabled
func setupMockClient(t *testing.T) (*BrewCommands, string) {
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

	// Create brew commands service
	return NewBrewCommands(httpClient), baseURL
}

func TestListV0_Success(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.BrewCommandsMock{}
	mockHandler.RegisterMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	result, _, err := service.ListV0(ctx)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, *result)

	// Verify the first brew command
	command := (*result)[0]
	assert.Equal(t, "brew outdated", command.Command)
	assert.Equal(t, "outdated", command.Label)
	assert.Equal(t, "mikemcquaid", command.LastUpdatedByUser)
	assert.Equal(t, 2, command.RunCount)

	// Verify devices
	assert.NotEmpty(t, command.Devices)
	assert.Contains(t, command.Devices, "TC6R2DHVHG")

	// Verify timestamps
	assert.NotNil(t, command.StartedAt)
	assert.NotNil(t, command.FinishedAt)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestListV0_Unauthorized(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.BrewCommandsMock{}
	mockHandler.RegisterErrorMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	result, _, err := service.ListV0(ctx)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "401")

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestCreateV0_Success(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.BrewCommandsMock{}
	mockHandler.RegisterMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	request := &CreateBrewCommandRequest{
		Arguments: "install wget",
	}

	result, _, err := service.CreateV0(ctx, request)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Brew Command was successfully created.", result.Message)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestCreateV0_FreeTier(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.BrewCommandsMock{}
	mockHandler.RegisterErrorMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	request := &CreateBrewCommandRequest{
		Arguments: "install wget",
	}

	result, _, err := service.CreateV0(ctx, request)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "403")
	// Verify it's a free tier error
	assert.True(t, client.IsFreeTierError(err))

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}
