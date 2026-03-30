package client

const (
	// DefaultBaseURL is the default base URL for the Workbrew API
	DefaultBaseURL = "https://console.workbrew.com"

	// UserAgentBase is the base user agent string prefix
	UserAgentBase = "go-api-sdk-workbrew"

	// DefaultTimeout is the default HTTP client timeout in seconds
	DefaultTimeout = 120

	// MaxRetries is the maximum number of retries for failed requests
	MaxRetries = 3

	// RetryWaitTime is the wait time between retries in seconds
	RetryWaitTime = 2

	// RetryMaxWaitTime is the maximum wait time between retries in seconds
	RetryMaxWaitTime = 10
)

// HTTP Status Codes
const (
	// Success status codes
	StatusOK      = 200 // Successful GET requests
	StatusCreated = 201 // Successful POST requests

	// Client error status codes
	StatusBadRequest          = 400 // Bad request, invalid arguments
	StatusUnauthorized        = 401 // Authentication required, invalid API key
	StatusForbidden           = 403 // Forbidden operation
	StatusNotFound            = 404 // Resource not found
	StatusConflict            = 409 // Resource already exists
	StatusUnprocessableEntity = 422 // Validation errors
	StatusFailedDependency    = 424 // Request depended on another request that failed
	StatusTooManyRequests     = 429 // Rate limit exceeded

	// Server error status codes
	StatusInternalServerError = 500 // Server-side error
	StatusBadGateway          = 502 // Gateway error
	StatusServiceUnavailable  = 503 // Transient error, service temporarily unavailable
	StatusGatewayTimeout      = 504 // Deadline exceeded
)

// Response format constants
const (
	FormatJSON = "json"
	FormatCSV  = "csv"
)

