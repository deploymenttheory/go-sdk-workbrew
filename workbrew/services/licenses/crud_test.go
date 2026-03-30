package licenses

import (
	"context"
	"testing"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/licenses/mocks"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func setupMockClient(t *testing.T) (*Licenses, string) {
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

	return NewLicenses(httpClient), baseURL
}

func TestListLicenses_Success(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.LicensesMock{}
	mockHandler.RegisterMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	result, _, err := service.ListLicenses(ctx)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, *result)

	// Verify first license (GPL-3.0-or-later)
	license := (*result)[0]
	assert.Equal(t, "GPL-3.0-or-later", license.Name)
	assert.Equal(t, 2, license.DeviceCount)
	assert.Equal(t, 2, license.FormulaCount)

	// Verify second license (MIT)
	license2 := (*result)[1]
	assert.Equal(t, "MIT", license2.Name)
	assert.Equal(t, 2, license2.DeviceCount)
	assert.Equal(t, 2, license2.FormulaCount)

	// Verify third license (curl)
	license3 := (*result)[2]
	assert.Equal(t, "curl", license3.Name)
	assert.Equal(t, 1, license3.DeviceCount)
	assert.Equal(t, 1, license3.FormulaCount)

	assert.Len(t, *result, 4)
	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestListLicenses_Unauthorized(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.LicensesMock{}
	mockHandler.RegisterErrorMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	result, _, err := service.ListLicenses(ctx)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "401")

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestListLicensesCSV_Success(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.LicensesMock{}
	mockHandler.RegisterMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	csvData, _, err := service.ListLicensesCSV(ctx)

	require.NoError(t, err)
	require.NotNil(t, csvData)
	assert.NotEmpty(t, csvData)

	// Verify CSV headers and content
	csvString := string(csvData)
	assert.Contains(t, csvString, "name,device_count,formula_count")
	assert.Contains(t, csvString, "GPL-3.0-or-later")
	assert.Contains(t, csvString, "MIT")
	assert.Contains(t, csvString, "curl")
	assert.Contains(t, csvString, "Artistic-2.0")
	assert.Contains(t, csvString, "AGPL-3.0-only")
	assert.Contains(t, csvString, "Unknown")

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestListLicensesCSV_Unauthorized(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.LicensesMock{}
	mockHandler.RegisterErrorMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	result, _, err := service.ListLicensesCSV(ctx)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "401")

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}
