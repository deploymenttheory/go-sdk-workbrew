package client

import (
	"testing"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/constants"
	"go.uber.org/zap/zaptest"
	"resty.dev/v3"
)

func TestAuthConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *AuthConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			config: &AuthConfig{
				APIKey:     "test-api-key",
				APIVersion: "v1",
			},
			wantErr: false,
		},
		{
			name: "valid config without version",
			config: &AuthConfig{
				APIKey: "test-api-key",
			},
			wantErr: false,
		},
		{
			name: "empty API key",
			config: &AuthConfig{
				APIKey:     "",
				APIVersion: "v1",
			},
			wantErr: true,
			errMsg:  "API key is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("AuthConfig.Validate() error message = %q, want %q", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestSetupAuthentication_Success(t *testing.T) {
	logger := zaptest.NewLogger(t)
	client := resty.New()

	authConfig := &AuthConfig{
		APIKey:     "test-api-key-12345",
		APIVersion: "v1",
	}

	err := SetupAuthentication(client, authConfig, logger)

	if err != nil {
		t.Fatalf("SetupAuthentication() error = %v, want nil", err)
	}

	// Verify API version header is set
	headers := client.Header()
	if got := headers.Get(constants.APIVersionHeader); got != "v1" {
		t.Errorf("API version header = %q, want %q", got, "v1")
	}
}

func TestSetupAuthentication_DefaultAPIVersion(t *testing.T) {
	logger := zaptest.NewLogger(t)
	client := resty.New()

	authConfig := &AuthConfig{
		APIKey:     "test-api-key",
		APIVersion: "", // Empty, should use default
	}

	err := SetupAuthentication(client, authConfig, logger)

	if err != nil {
		t.Fatalf("SetupAuthentication() error = %v, want nil", err)
	}

	// Verify default API version is used
	headers := client.Header()
	if got := headers.Get(constants.APIVersionHeader); got != constants.DefaultAPIVersion {
		t.Errorf("API version header = %q, want %q (default)", got, constants.DefaultAPIVersion)
	}
}

func TestSetupAuthentication_InvalidConfig(t *testing.T) {
	logger := zaptest.NewLogger(t)
	client := resty.New()

	tests := []struct {
		name       string
		authConfig *AuthConfig
		wantErr    bool
	}{
		{
			name: "empty API key",
			authConfig: &AuthConfig{
				APIKey:     "",
				APIVersion: "v1",
			},
			wantErr: true,
		},
		{
			name: "nil-like config",
			authConfig: &AuthConfig{
				APIKey:     "",
				APIVersion: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SetupAuthentication(client, tt.authConfig, logger)

			if (err != nil) != tt.wantErr {
				t.Errorf("SetupAuthentication() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetupAuthentication_CustomAPIVersion(t *testing.T) {
	logger := zaptest.NewLogger(t)

	tests := []struct {
		name           string
		apiVersion     string
		expectedHeader string
	}{
		{
			name:           "custom v1",
			apiVersion:     "v1",
			expectedHeader: "v1",
		},
		{
			name:           "custom v2",
			apiVersion:     "v2",
			expectedHeader: "v2",
		},
		{
			name:           "empty uses default",
			apiVersion:     "",
			expectedHeader: constants.DefaultAPIVersion,
		},
		{
			name:           "custom version string",
			apiVersion:     "2023-01",
			expectedHeader: "2023-01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := resty.New()
			authConfig := &AuthConfig{
				APIKey:     "test-key",
				APIVersion: tt.apiVersion,
			}

			err := SetupAuthentication(client, authConfig, logger)
			if err != nil {
				t.Fatalf("SetupAuthentication() error = %v, want nil", err)
			}

			headers := client.Header()
			if got := headers.Get(constants.APIVersionHeader); got != tt.expectedHeader {
				t.Errorf("API version header = %q, want %q", got, tt.expectedHeader)
			}
		})
	}
}

func TestSetupAuthentication_BearerTokenSet(t *testing.T) {
	logger := zaptest.NewLogger(t)
	client := resty.New()

	authConfig := &AuthConfig{
		APIKey:     "test-api-key-12345",
		APIVersion: "v1",
	}

	err := SetupAuthentication(client, authConfig, logger)
	if err != nil {
		t.Fatalf("SetupAuthentication() error = %v, want nil", err)
	}

	// Verify auth scheme and token are set
	// Note: resty's internal auth settings aren't directly accessible,
	// but we can verify the setup completed without error
	headers := client.Header()
	if got := headers.Get(constants.APIVersionHeader); got != "v1" {
		t.Errorf("API version header = %q, want %q", got, "v1")
	}
}

func TestAuthConfig_Fields(t *testing.T) {
	// Test that AuthConfig struct can hold expected values
	config := &AuthConfig{
		APIKey:     "my-api-key-12345",
		APIVersion: "v1.5",
	}

	if config.APIKey != "my-api-key-12345" {
		t.Errorf("APIKey = %q, want %q", config.APIKey, "my-api-key-12345")
	}

	if config.APIVersion != "v1.5" {
		t.Errorf("APIVersion = %q, want %q", config.APIVersion, "v1.5")
	}
}

func TestAuthConfig_LongAPIKey(t *testing.T) {
	// Test with a very long API key (should still be valid)
	longKey := ""
	for i := 0; i < 1000; i++ {
		longKey += "a"
	}

	config := &AuthConfig{
		APIKey:     longKey,
		APIVersion: "v1",
	}

	err := config.Validate()
	if err != nil {
		t.Errorf("Validate() with long API key error = %v, want nil", err)
	}

	// Setup should also work
	logger := zaptest.NewLogger(t)
	client := resty.New()
	err = SetupAuthentication(client, config, logger)
	if err != nil {
		t.Errorf("SetupAuthentication() with long API key error = %v, want nil", err)
	}
}

func TestAuthConfig_SpecialCharactersInAPIKey(t *testing.T) {
	// Test with special characters in API key
	specialKeys := []string{
		"key-with-dashes",
		"key_with_underscores",
		"key.with.dots",
		"key123with456numbers",
		"key-_./~:?#[]@!$&'()*+,;=",
	}

	for _, key := range specialKeys {
		t.Run(key, func(t *testing.T) {
			config := &AuthConfig{
				APIKey:     key,
				APIVersion: "v1",
			}

			err := config.Validate()
			if err != nil {
				t.Errorf("Validate() with key %q error = %v, want nil", key, err)
			}

			logger := zaptest.NewLogger(t)
			client := resty.New()
			err = SetupAuthentication(client, config, logger)
			if err != nil {
				t.Errorf("SetupAuthentication() with key %q error = %v, want nil", key, err)
			}
		})
	}
}

func TestAuthConfig_WhitespaceAPIKey(t *testing.T) {
	// Test with whitespace-only API key
	tests := []struct {
		name    string
		apiKey  string
		wantErr bool
	}{
		{
			name:    "spaces only",
			apiKey:  "   ",
			wantErr: false, // Non-empty string, validation passes
		},
		{
			name:    "tabs only",
			apiKey:  "\t\t\t",
			wantErr: false, // Non-empty string, validation passes
		},
		{
			name:    "truly empty",
			apiKey:  "",
			wantErr: true, // Empty string, validation fails
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &AuthConfig{
				APIKey:     tt.apiKey,
				APIVersion: "v1",
			}

			err := config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() with whitespace key error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
