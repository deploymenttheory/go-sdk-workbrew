package workbrew

import (
	"fmt"
	"os"

	"go.uber.org/zap"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/analytics"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/brewcommands"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/brewconfigurations"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/brewfiles"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/brewtaps"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/casks"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/devicegroups"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/devices"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/events"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/formulae"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/licenses"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/vulnerabilities"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/vulnerabilitychanges"
)

// Client is the main entry point for the Workbrew API SDK.
// It aggregates all service clients and provides a unified interface.
// Users should interact with the API exclusively through the provided service methods.
type Client struct {
	// transport is the internal HTTP transport layer (not exposed to users)
	transport *client.Transport

	// Services - users should only call methods on these services
	Analytics            *analytics.Analytics
	BrewCommands         *brewcommands.BrewCommands
	BrewConfigurations   *brewconfigurations.BrewConfigurations
	Brewfiles            *brewfiles.Brewfiles
	BrewTaps             *brewtaps.BrewTaps
	Casks                *casks.Casks
	DeviceGroups         *devicegroups.DeviceGroups
	Devices              *devices.Devices
	Events               *events.Events
	Formulae             *formulae.Formulae
	Licenses             *licenses.Licenses
	Vulnerabilities      *vulnerabilities.Vulnerabilities
	VulnerabilityChanges *vulnerabilitychanges.VulnerabilityChanges
}

// NewClient creates a new Workbrew API client
//
// Parameters:
//   - apiKey: The bearer token for authentication
//   - workspaceName: The workspace slug to operate on
//   - options: Optional client configuration options
//
// Example:
//
//	client, err := workbrew.NewClient(
//	    "your-api-key",
//	    "your-workspace",
//	    workbrew.WithDebug(),
//	)
func NewClient(apiKey string, workspaceName string, options ...ClientOption) (*Client, error) {
	// Create base HTTP transport
	transport, err := client.NewTransport(apiKey, workspaceName, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP transport: %w", err)
	}

	// Initialize service clients
	c := &Client{
		transport:            transport,
		Analytics:            analytics.NewAnalytics(transport),
		BrewCommands:         brewcommands.NewBrewCommands(transport),
		BrewConfigurations:   brewconfigurations.NewBrewConfigurations(transport),
		Brewfiles:            brewfiles.NewBrewfiles(transport),
		BrewTaps:             brewtaps.NewBrewTaps(transport),
		Casks:                casks.NewCasks(transport),
		DeviceGroups:         devicegroups.NewDeviceGroups(transport),
		Devices:              devices.NewDevices(transport),
		Events:               events.NewEvents(transport),
		Formulae:             formulae.NewFormulae(transport),
		Licenses:             licenses.NewLicenses(transport),
		Vulnerabilities:      vulnerabilities.NewVulnerabilities(transport),
		VulnerabilityChanges: vulnerabilitychanges.NewVulnerabilityChanges(transport),
	}

	return c, nil
}

// NewClientFromEnv creates a new client using environment variables
//
// Required environment variables:
//   - WORKBREW_API_KEY: The bearer token for authentication
//   - WORKBREW_WORKSPACE: The workspace slug
//
// Optional environment variables:
//   - WORKBREW_BASE_URL: Custom base URL (defaults to https://console.workbrew.com)
//   - WORKBREW_API_VERSION: API version (defaults to v0)
//
// Example:
//
//	client, err := workbrew.NewClientFromEnv()
func NewClientFromEnv(options ...ClientOption) (*Client, error) {
	apiKey := os.Getenv("WORKBREW_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("WORKBREW_API_KEY environment variable is required")
	}

	workspaceName := os.Getenv("WORKBREW_WORKSPACE")
	if workspaceName == "" {
		return nil, fmt.Errorf("WORKBREW_WORKSPACE environment variable is required")
	}

	// Check for optional environment variables and append to options
	if baseURL := os.Getenv("WORKBREW_BASE_URL"); baseURL != "" {
		options = append(options, client.WithBaseURL(baseURL))
	}

	if apiVersion := os.Getenv("WORKBREW_API_VERSION"); apiVersion != "" {
		options = append(options, client.WithAPIVersion(apiVersion))
	}

	return NewClient(apiKey, workspaceName, options...)
}

// SetWorkspace changes the active workspace for all subsequent API calls.
// This updates the base URL to target the specified workspace.
//
// Parameters:
//   - workspaceName: The name of the workspace to switch to
//
// Example:
//
//	client.SetWorkspace("production-workspace")
func (c *Client) SetWorkspace(workspaceName string) {
	c.transport.SetWorkspace(workspaceName)
}

// GetLogger returns the configured zap logger instance.
// Use this to add custom logging within your application using the same logger.
//
// Returns:
//   - *zap.Logger: The configured logger instance
func (c *Client) GetLogger() *zap.Logger {
	return c.transport.GetLogger()
}
