package workbrew

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
	"go.uber.org/zap"
)

// ClientOption configures the Workbrew API transport at construction time.
// Pass one or more ClientOption values to NewClient or NewClientFromEnv.
type ClientOption = client.ClientOption

// OTelConfig is an alias for client.OTelConfig.
// Use this to configure OpenTelemetry tracing for the client.
type OTelConfig = client.OTelConfig

// WithBaseURL sets a custom base URL for the API client.
func WithBaseURL(baseURL string) ClientOption {
	return client.WithBaseURL(baseURL)
}

// WithAPIVersion sets a custom API version.
func WithAPIVersion(version string) ClientOption {
	return client.WithAPIVersion(version)
}

// WithAPIKey allows setting the API key during client initialization.
func WithAPIKey(apiKey string) ClientOption {
	return client.WithAPIKey(apiKey)
}

// WithTimeout sets a custom timeout for HTTP requests.
func WithTimeout(timeout time.Duration) ClientOption {
	return client.WithTimeout(timeout)
}

// WithRetryCount sets the number of retries for failed requests.
func WithRetryCount(count int) ClientOption {
	return client.WithRetryCount(count)
}

// WithRetryWaitTime sets the default wait time between retry attempts.
func WithRetryWaitTime(waitTime time.Duration) ClientOption {
	return client.WithRetryWaitTime(waitTime)
}

// WithRetryMaxWaitTime sets the maximum wait time between retry attempts.
func WithRetryMaxWaitTime(maxWaitTime time.Duration) ClientOption {
	return client.WithRetryMaxWaitTime(maxWaitTime)
}

// WithLogger sets a custom zap logger for the client.
func WithLogger(logger *zap.Logger) ClientOption {
	return client.WithLogger(logger)
}

// WithDebug enables debug mode which logs request and response details.
func WithDebug() ClientOption {
	return client.WithDebug()
}

// WithUserAgent sets a custom user agent string.
func WithUserAgent(userAgent string) ClientOption {
	return client.WithUserAgent(userAgent)
}

// WithCustomAgent appends a custom identifier to the default user agent.
func WithCustomAgent(customAgent string) ClientOption {
	return client.WithCustomAgent(customAgent)
}

// WithGlobalHeader sets a global header included in all requests.
func WithGlobalHeader(key, value string) ClientOption {
	return client.WithGlobalHeader(key, value)
}

// WithGlobalHeaders sets multiple global headers at once.
func WithGlobalHeaders(headers map[string]string) ClientOption {
	return client.WithGlobalHeaders(headers)
}

// WithProxy sets an HTTP proxy for all requests.
func WithProxy(proxyURL string) ClientOption {
	return client.WithProxy(proxyURL)
}

// WithTLSClientConfig sets custom TLS configuration.
func WithTLSClientConfig(tlsConfig *tls.Config) ClientOption {
	return client.WithTLSClientConfig(tlsConfig)
}

// WithClientCertificate sets a client certificate for mutual TLS authentication.
func WithClientCertificate(certFile, keyFile string) ClientOption {
	return client.WithClientCertificate(certFile, keyFile)
}

// WithClientCertificateFromString sets a client certificate from PEM-encoded strings.
func WithClientCertificateFromString(certPEM, keyPEM string) ClientOption {
	return client.WithClientCertificateFromString(certPEM, keyPEM)
}

// WithRootCertificates adds custom root CA certificates for server validation.
func WithRootCertificates(pemFilePaths ...string) ClientOption {
	return client.WithRootCertificates(pemFilePaths...)
}

// WithRootCertificateFromString adds a custom root CA certificate from a PEM string.
func WithRootCertificateFromString(pemContent string) ClientOption {
	return client.WithRootCertificateFromString(pemContent)
}

// WithTransport sets a custom HTTP transport (http.RoundTripper).
func WithTransport(transport http.RoundTripper) ClientOption {
	return client.WithTransport(transport)
}

// WithInsecureSkipVerify disables TLS certificate verification (use only for testing).
func WithInsecureSkipVerify() ClientOption {
	return client.WithInsecureSkipVerify()
}

// WithMinTLSVersion sets the minimum TLS version for connections.
func WithMinTLSVersion(minVersion uint16) ClientOption {
	return client.WithMinTLSVersion(minVersion)
}

// WithTracing enables OpenTelemetry tracing for all HTTP requests.
func WithTracing(config *OTelConfig) ClientOption {
	return client.WithTracing(config)
}
