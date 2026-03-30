// Package constants defines shared constants used across the Workbrew SDK.
// It has no dependencies and may be imported by any layer — transport, services, or tests.
package constants

// ============================================================================
// Authentication Configuration
// ============================================================================

const (
	// BearerScheme is the HTTP authentication scheme for bearer tokens
	BearerScheme = "Bearer"

	// AuthorizationHeader is the HTTP Authorization header name
	AuthorizationHeader = "Authorization"

	// APIVersionHeader is the header name for the Workbrew API version
	APIVersionHeader = "X-Workbrew-API-Version"

	// DefaultAPIVersion is the default Workbrew API version
	DefaultAPIVersion = "v0"
)
