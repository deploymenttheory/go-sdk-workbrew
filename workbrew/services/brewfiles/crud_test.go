package brewfiles

import (
	"context"
	"testing"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/brewfiles/mocks"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func setupMockClient(t *testing.T) (*Brewfiles, string) {
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

	return NewBrewfiles(httpClient), baseURL
}

func TestListV0_Success(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.BrewfilesMock{}
	mockHandler.RegisterMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	result, _, err := service.ListV0(ctx)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, *result)

	// Verify first brewfile
	brewfile := (*result)[0]
	assert.Equal(t, "my-brewfile", brewfile.Label)
	assert.Equal(t, "my-brewfile", brewfile.Slug)
	assert.Equal(t, "brew \"wget\"", brewfile.Content)
	assert.Equal(t, "onboarded", brewfile.LastUpdatedByUser)
	assert.Equal(t, "Not Started", brewfile.StartedAt)
	assert.Equal(t, "Not Finished", brewfile.FinishedAt)
	assert.Equal(t, 1, brewfile.RunCount)

	// Verify second brewfile
	brewfile2 := (*result)[1]
	assert.Equal(t, "production", brewfile2.Label)
	assert.Contains(t, brewfile2.Content, "brew \"git\"")
	assert.Equal(t, "admin", brewfile2.LastUpdatedByUser)
	assert.Equal(t, 5, brewfile2.RunCount)

	assert.Len(t, *result, 2)
	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestListV0_Unauthorized(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.BrewfilesMock{}
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
	mockHandler := &mocks.BrewfilesMock{}
	mockHandler.RegisterMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	csvData, _, err := service.ListCSVV0(ctx)

	require.NoError(t, err)
	require.NotNil(t, csvData)
	assert.NotEmpty(t, csvData)

	// Verify CSV headers and content
	csvString := string(csvData)
	assert.Contains(t, csvString, "label,last_updated_by_user,started_at,finished_at,devices,run_count")
	assert.Contains(t, csvString, "my-brewfile")
	assert.Contains(t, csvString, "production")
	assert.Contains(t, csvString, "Not Started")

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestListCSVV0_Unauthorized(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.BrewfilesMock{}
	mockHandler.RegisterErrorMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	result, _, err := service.ListCSVV0(ctx)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "401")

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestCreateV0_Success(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.BrewfilesMock{}
	mockHandler.RegisterMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	request := &CreateBrewfileRequest{
		Label:   "test-brewfile",
		Content: "brew \"wget\"",
	}

	result, _, err := service.CreateV0(ctx, request)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Brewfile was successfully created.", result.Message)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestCreateV0_WithDevices(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.BrewfilesMock{}
	mockHandler.RegisterMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	deviceSerials := "TC6R2DHVHG,1234567890"
	request := &CreateBrewfileRequest{
		Label:               "test-brewfile",
		Content:             "brew \"wget\"",
		DeviceSerialNumbers: &deviceSerials,
	}

	result, _, err := service.CreateV0(ctx, request)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Brewfile was successfully created.", result.Message)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestCreateV0_WithDeviceGroup(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.BrewfilesMock{}
	mockHandler.RegisterMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	deviceGroupID := "ddba0af6-bd3c-5abf-8311-e62dc6bd9fbc"
	request := &CreateBrewfileRequest{
		Label:         "test-brewfile",
		Content:       "brew \"wget\"",
		DeviceGroupID: &deviceGroupID,
	}

	result, _, err := service.CreateV0(ctx, request)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Brewfile was successfully created.", result.Message)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestCreateV0_Unauthorized(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.BrewfilesMock{}
	mockHandler.RegisterErrorMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	request := &CreateBrewfileRequest{
		Label:   "test-brewfile",
		Content: "brew \"wget\"",
	}

	result, _, err := service.CreateV0(ctx, request)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "401")

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestCreateV0_Forbidden(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.BrewfilesMock{}
	mockHandler.RegisterForbiddenMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	request := &CreateBrewfileRequest{
		Label:   "test-brewfile",
		Content: "brew \"wget\"",
	}

	result, _, err := service.CreateV0(ctx, request)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "403")

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestUpdateByLabelV0_Success(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.BrewfilesMock{}
	mockHandler.RegisterMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	request := &UpdateBrewfileRequest{
		Content: "brew \"wget\"\nbrew \"htop\"",
	}

	result, _, err := service.UpdateByLabelV0(ctx, "my-brewfile", request)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Brewfile was successfully updated.", result.Message)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestUpdateByLabelV0_EmptyLabel(t *testing.T) {
	service, _ := setupMockClient(t)

	ctx := context.Background()
	request := &UpdateBrewfileRequest{
		Content: "brew \"wget\"",
	}

	result, _, err := service.UpdateByLabelV0(ctx, "", request)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "brewfile label is required")

	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestUpdateByLabelV0_Unauthorized(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.BrewfilesMock{}
	mockHandler.RegisterErrorMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	request := &UpdateBrewfileRequest{
		Content: "brew \"wget\"",
	}

	result, _, err := service.UpdateByLabelV0(ctx, "my-brewfile", request)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "401")

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestUpdateByLabelV0_ValidationError(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.BrewfilesMock{}
	mockHandler.RegisterValidationMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	request := &UpdateBrewfileRequest{
		Content: "tap \"foo/bar/baz\"", // Invalid according to API
	}

	result, _, err := service.UpdateByLabelV0(ctx, "my-brewfile", request)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "422")

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestDeleteByLabelV0_Success(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.BrewfilesMock{}
	mockHandler.RegisterMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, _, err := service.DeleteByLabelV0(ctx, "my-brewfile")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "Brewfile was successfully destroyed.", result.Message)

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestDeleteByLabelV0_EmptyLabel(t *testing.T) {
	service, _ := setupMockClient(t)

	ctx := context.Background()

	result, _, err := service.DeleteByLabelV0(ctx, "")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "brewfile label is required")

	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestDeleteByLabelV0_Unauthorized(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.BrewfilesMock{}
	mockHandler.RegisterErrorMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, _, err := service.DeleteByLabelV0(ctx, "my-brewfile")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "401")

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestListRunsByLabelV0_Success(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.BrewfilesMock{}
	mockHandler.RegisterMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, _, err := service.ListRunsByLabelV0(ctx, "my-brewfile")

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, *result)

	// Verify first run
	run := (*result)[0]
	assert.Equal(t, "my-brewfile", run.Label)
	assert.Equal(t, "TC6R2DHVHG", run.Device)
	assert.True(t, run.Success)
	assert.Contains(t, run.Output, "brew bundle")
	assert.Equal(t, "2023-11-01T12:34:56.000Z", run.StartedAt)
	assert.Equal(t, "2023-11-01T21:43:12.000Z", run.FinishedAt)

	// Verify second run - not started
	run2 := (*result)[1]
	assert.Equal(t, "1234567890", run2.Device)
	assert.False(t, run2.Success)
	assert.Equal(t, "Not Started", run2.StartedAt)
	assert.Equal(t, "Not Finished", run2.FinishedAt)

	assert.Len(t, *result, 2)
	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestListRunsByLabelV0_EmptyLabel(t *testing.T) {
	service, _ := setupMockClient(t)

	ctx := context.Background()

	result, _, err := service.ListRunsByLabelV0(ctx, "")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "brewfile label is required")

	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestListRunsByLabelV0_Unauthorized(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.BrewfilesMock{}
	mockHandler.RegisterErrorMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, _, err := service.ListRunsByLabelV0(ctx, "my-brewfile")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "401")

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestListRunsByLabelCSVV0_Success(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.BrewfilesMock{}
	mockHandler.RegisterMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	csvData, _, err := service.ListRunsByLabelCSVV0(ctx, "my-brewfile")

	require.NoError(t, err)
	require.NotNil(t, csvData)
	assert.NotEmpty(t, csvData)

	// Verify CSV headers and content
	csvString := string(csvData)
	assert.Contains(t, csvString, "label,device,created_at,updated_at,success,output,started_at,finished_at")
	assert.Contains(t, csvString, "my-brewfile")
	assert.Contains(t, csvString, "Not Started")

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestListRunsByLabelCSVV0_EmptyLabel(t *testing.T) {
	service, _ := setupMockClient(t)

	ctx := context.Background()

	result, _, err := service.ListRunsByLabelCSVV0(ctx, "")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "brewfile label is required")

	assert.Equal(t, 0, httpmock.GetTotalCallCount())
}

func TestListRunsByLabelCSVV0_Unauthorized(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.BrewfilesMock{}
	mockHandler.RegisterErrorMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()

	result, _, err := service.ListRunsByLabelCSVV0(ctx, "my-brewfile")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "401")

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}
