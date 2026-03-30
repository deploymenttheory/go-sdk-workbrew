package client

import (
	"fmt"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/constants"
	"go.uber.org/zap"
	"resty.dev/v3"
)

// AuthConfig holds authentication configuration for the Workbrew API.
type AuthConfig struct {
	// APIKey is the bearer token for authentication
	APIKey string

	// APIVersion is the API version (defaults to v0)
	APIVersion string
}

// Validate checks if the authentication configuration is valid.
func (a *AuthConfig) Validate() error {
	if a.APIKey == "" {
		return fmt.Errorf("API key is required")
	}
	return nil
}

// SetupAuthentication configures the HTTP client with bearer token authentication.
// The API key is set once during client initialization and cannot be changed.
//
// Parameters:
//   - client: The resty HTTP client to configure
//   - authConfig: Authentication configuration containing API key and version
//   - logger: Logger instance for logging authentication setup
//
// Returns:
//   - error: Any error encountered during authentication setup
func SetupAuthentication(client *resty.Client, authConfig *AuthConfig, logger *zap.Logger) error {
	if err := authConfig.Validate(); err != nil {
		logger.Error("Authentication validation failed", zap.Error(err))
		return fmt.Errorf("authentication validation failed: %w", err)
	}

	// Set bearer token authentication
	client.SetAuthScheme(constants.BearerScheme)
	client.SetAuthToken(authConfig.APIKey)

	// Set API version header
	apiVersion := authConfig.APIVersion
	if apiVersion == "" {
		apiVersion = constants.DefaultAPIVersion
	}
	client.SetHeader(constants.APIVersionHeader, apiVersion)

	logger.Info("Authentication configured",
		zap.String("api_version", apiVersion))

	return nil
}
