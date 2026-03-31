package formulae

import (
	"context"
	"testing"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/formulae/mocks"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func setupMockClient(t *testing.T) (*Formulae, string) {
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

	return NewFormulae(httpClient), baseURL
}

func TestListV0_Success(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.FormulaeMock{}
	mockHandler.RegisterMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	result, _, err := service.ListV0(ctx)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, *result)

	// Verify first formula
	formula := (*result)[0]
	assert.Equal(t, "curl", formula.Name)
	assert.Contains(t, formula.Devices, "TC6R2DHVHG")
	assert.Contains(t, formula.Devices, "1234567890")
	assert.True(t, formula.Outdated)
	assert.False(t, formula.InstalledOnRequest)
	assert.True(t, formula.InstalledAsDependency)
	assert.Contains(t, formula.Vulnerabilities, "CVE-2024-2466")
	assert.NotNil(t, formula.License)
	assert.Contains(t, *formula.License, "curl")
	assert.NotNil(t, formula.HomebrewCoreVersion)
	assert.Equal(t, "8.11.0_1", *formula.HomebrewCoreVersion)

	// Verify we have multiple formulae
	assert.Len(t, *result, 3)

	// Verify second formula (actionlint)
	formula2 := (*result)[1]
	assert.Equal(t, "actionlint", formula2.Name)
	assert.Contains(t, formula2.Devices, "TC6R2DHVHG")
	assert.Empty(t, formula2.Vulnerabilities)

	// Verify third formula (wget) - installed on request
	formula3 := (*result)[2]
	assert.Equal(t, "wget", formula3.Name)
	assert.True(t, formula3.InstalledOnRequest)
	assert.False(t, formula3.InstalledAsDependency)
	assert.Contains(t, formula3.Vulnerabilities, "CVE-2024-10524")

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestListV0_Unauthorized(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.FormulaeMock{}
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
	mockHandler := &mocks.FormulaeMock{}
	mockHandler.RegisterMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	csvData, _, err := service.ListCSVV0(ctx)

	require.NoError(t, err)
	require.NotNil(t, csvData)
	assert.NotEmpty(t, csvData)

	// Verify CSV headers
	csvString := string(csvData)
	assert.Contains(t, csvString, "name,devices,outdated,installed_on_request,installed_as_dependency")
	assert.Contains(t, csvString, "vulnerabilities,deprecated,license,homebrew_core_version")
	assert.Contains(t, csvString, "curl")
	assert.Contains(t, csvString, "actionlint")
	assert.Contains(t, csvString, "wget")
	assert.Contains(t, csvString, "CVE-2024-11053")

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestListCSVV0_Unauthorized(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.FormulaeMock{}
	mockHandler.RegisterErrorMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	result, _, err := service.ListCSVV0(ctx)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "401")

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}
