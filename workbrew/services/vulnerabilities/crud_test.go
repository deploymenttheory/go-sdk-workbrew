package vulnerabilities

import (
	"context"
	"testing"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/vulnerabilities/mocks"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func setupMockClient(t *testing.T) (*Vulnerabilities, string) {
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

	return NewVulnerabilities(httpClient), baseURL
}

func TestListVulnerabilities_Success(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.VulnerabilitiesMock{}
	mockHandler.RegisterMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	result, _, err := service.ListVulnerabilities(ctx)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, *result)

	// Verify first vulnerability (curl)
	vuln := (*result)[0]
	assert.Equal(t, "curl", vuln.Formula)
	assert.Contains(t, vuln.OutdatedDevices, "TC6R2DHVHG")
	assert.Contains(t, vuln.OutdatedDevices, "1234567890")
	assert.False(t, vuln.Supported)
	assert.Equal(t, "8.11.0_1", vuln.HomebrewCoreVersion)
	
	// Verify vulnerability details
	assert.Len(t, vuln.Vulnerabilities, 2)
	assert.Equal(t, "CVE-2024-2466", vuln.Vulnerabilities[0].CleanID)
	assert.NotNil(t, vuln.Vulnerabilities[0].CVSSScore)
	assert.Equal(t, 6.5, *vuln.Vulnerabilities[0].CVSSScore)
	
	assert.Equal(t, "THIS-IS-AN-INVALID-CVE-001", vuln.Vulnerabilities[1].CleanID)
	assert.NotNil(t, vuln.Vulnerabilities[1].CVSSScore)
	assert.Equal(t, 8.0, *vuln.Vulnerabilities[1].CVSSScore)

	// Verify second vulnerability (wget)
	vuln2 := (*result)[1]
	assert.Equal(t, "wget", vuln2.Formula)
	assert.Contains(t, vuln2.OutdatedDevices, "TC6R2DHVHG")
	assert.Len(t, vuln2.Vulnerabilities, 1)
	assert.Equal(t, "CVE-2024-10524", vuln2.Vulnerabilities[0].CleanID)
	assert.Nil(t, vuln2.Vulnerabilities[0].CVSSScore)

	assert.Len(t, *result, 2)
	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestListVulnerabilities_Unauthorized(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.VulnerabilitiesMock{}
	mockHandler.RegisterErrorMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	result, _, err := service.ListVulnerabilities(ctx)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "401")

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestListVulnerabilities_Forbidden(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.VulnerabilitiesMock{}
	mockHandler.RegisterForbiddenMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	result, _, err := service.ListVulnerabilities(ctx)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "403")

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestListVulnerabilitiesCSV_Success(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.VulnerabilitiesMock{}
	mockHandler.RegisterMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	csvData, _, err := service.ListVulnerabilitiesCSV(ctx)

	require.NoError(t, err)
	require.NotNil(t, csvData)
	assert.NotEmpty(t, csvData)

	// Verify CSV headers and content
	csvString := string(csvData)
	assert.Contains(t, csvString, "vulnerabilities,formula,outdated_devices,supported,homebrew_core_version")
	assert.Contains(t, csvString, "CVE-2024-11053")
	assert.Contains(t, csvString, "curl")
	assert.Contains(t, csvString, "wget")
	assert.Contains(t, csvString, "ack")
	assert.Contains(t, csvString, "renovate")
	assert.Contains(t, csvString, "GHSA-rqgv-292v-5qgr")

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}

func TestListVulnerabilitiesCSV_Unauthorized(t *testing.T) {
	service, baseURL := setupMockClient(t)
	mockHandler := &mocks.VulnerabilitiesMock{}
	mockHandler.RegisterErrorMocks(baseURL)
	defer mockHandler.CleanupMockState()

	ctx := context.Background()
	result, _, err := service.ListVulnerabilitiesCSV(ctx)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "401")

	assert.Equal(t, 1, httpmock.GetTotalCallCount())
}
