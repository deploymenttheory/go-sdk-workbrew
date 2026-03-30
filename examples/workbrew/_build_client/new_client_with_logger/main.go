package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew"
	
	"go.uber.org/zap"
)

// This example demonstrates creating a production-ready client with custom configuration.
//
// Use this approach when:
// - Running in production environments
// - You need structured logging
// - You want to customize timeouts and retries
// - You need to add custom headers
// - You want fine-grained control over client behavior
//
// This example shows:
// - Structured logging with zap
// - Custom timeout configuration
// - Retry policy tuning
// - Custom headers for request tracking
// - Debug mode for development

func main() {
	// Check API key and workspace from environment
	apiKey := os.Getenv("WORKBREW_API_KEY")
	if apiKey == "" {
		log.Fatal("WORKBREW_API_KEY environment variable is required")
	}

	workspace := os.Getenv("WORKBREW_WORKSPACE")
	if workspace == "" {
		log.Fatal("WORKBREW_WORKSPACE environment variable is required")
	}

	// Step 1: Create a structured logger
	// Use NewProduction() for production, NewDevelopment() for local dev
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	// Step 2: Create client with custom configuration
	workbrewClient, err := workbrew.NewClient(
		apiKey,
		workspace,

		// Structured logging for production observability
		workbrew.WithLogger(logger),

		// Custom timeout for slow networks or large operations
		workbrew.WithTimeout(60*time.Second),

		// Retry configuration for better reliability
		workbrew.WithRetryCount(5),                    // Retry up to 5 times
		workbrew.WithRetryWaitTime(3*time.Second),     // Initial wait time
		workbrew.WithRetryMaxWaitTime(30*time.Second), // Maximum wait time

		// Add custom headers for request tracking
		workbrew.WithGlobalHeader("X-Application-Name", "MyDeviceManager"),
		workbrew.WithGlobalHeader("X-Application-Version", "1.0.0"),
		workbrew.WithGlobalHeader("X-Environment", "production"),

		// Uncomment to enable debug mode (only for development!)
		// workbrew.WithDebug(),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	logger.Info("Workbrew client created",
		zap.String("workspace", workspace),
		zap.String("timeout", "60s"),
		zap.Int("retry_count", 5))

	// Step 3: Use the client with structured logging
	ctx := context.Background()

	logger.Info("Fetching device groups")

	deviceGroups, resp, err := workbrewClient.DeviceGroups.ListDeviceGroups(ctx)
	if err != nil {
		logger.Error("Failed to list device groups",
			zap.Error(err),
			zap.Int("status_code", resp.StatusCode))
		log.Fatalf("API call failed: %v", err)
	}

	// Log successful operation
	logger.Info("Device groups retrieved",
		zap.Int("status_code", resp.StatusCode),
		zap.Duration("duration", resp.Duration),
		zap.Int64("response_size", resp.Size),
		zap.Int("group_count", len(*deviceGroups)))

	// Display results
	fmt.Printf("\n✓ Production-ready client created successfully\n\n")
	fmt.Printf("Configuration:\n")
	fmt.Printf("  Timeout: 60s\n")
	fmt.Printf("  Retry Count: 5\n")
	fmt.Printf("  Logger: zap (production)\n")
	fmt.Printf("  Custom Headers: 3\n")

	fmt.Printf("\nDevice Groups:\n")
	fmt.Printf("  Total Groups: %d\n", len(*deviceGroups))
	if len(*deviceGroups) > 0 {
		fmt.Printf("\nFirst 3 Groups:\n")
		for i, group := range (*deviceGroups)[:min(3, len(*deviceGroups))] {
			deviceCount := 0
			if group.DeviceCount != nil {
				deviceCount = *group.DeviceCount
			}
			fmt.Printf("  %d. %s (%d devices)\n", i+1, group.Name, deviceCount)
		}
	}

	fmt.Printf("\nAPI Response:\n")
	fmt.Printf("  Status Code: %d\n", resp.StatusCode)
	fmt.Printf("  Duration: %v\n", resp.Duration)
	fmt.Printf("  Response Size: %d bytes\n", resp.Size)

	fmt.Printf("\n✓ Custom client with logging example completed successfully!\n")
	fmt.Printf("\n💡 Production Tips:\n")
	fmt.Printf("   Use structured logging for better observability\n")
	fmt.Printf("   Adjust timeouts based on your network conditions\n")
	fmt.Printf("   Configure retries for transient failures\n")
	fmt.Printf("   Add custom headers for request tracking\n")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
