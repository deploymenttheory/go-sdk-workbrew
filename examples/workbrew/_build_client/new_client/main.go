package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew"
)

// This example demonstrates the most basic way to create a Workbrew client.
//
// IMPORTANT SECURITY NOTE:
// This example shows both environment variable (recommended) and hardcoded
// API key approaches. Always use environment variables in real code!
//
// Use this approach when:
// - You want the simplest possible client setup
// - You don't need custom configuration
// - You're getting started with the SDK
//
// The client uses sensible defaults:
// - 120 second timeout
// - 3 retries with exponential backoff
// - Production-level logging

func main() {
	// OPTION 1: From environment variables (RECOMMENDED)
	// This is the recommended approach - never hardcode API keys!
	apiKey := os.Getenv("WORKBREW_API_KEY")
	if apiKey == "" {
		log.Fatal("WORKBREW_API_KEY environment variable is required")
	}

	workspace := os.Getenv("WORKBREW_WORKSPACE")
	if workspace == "" {
		log.Fatal("WORKBREW_WORKSPACE environment variable is required")
	}

	// OPTION 2: Hardcoded (NOT RECOMMENDED - for demonstration only)
	// Never do this in production! Only for local testing/learning.
	// Hardcoded keys can be accidentally committed to version control.
	// apiKey := "your-api-key-here"     // ⚠️ DON'T DO THIS IN REAL CODE!
	// workspace := "your-workspace-id"  // ⚠️ DON'T DO THIS IN REAL CODE!

	// Create the simplest possible client - just pass the API key and workspace
	client, err := workbrew.NewClient(apiKey, workspace)
	if err != nil {
		log.Fatalf("Failed to create Workbrew client: %v", err)
	}

	// Use the client to make a simple API call
	ctx := context.Background()

	devices, resp, err := client.Devices.ListV0(ctx)
	if err != nil {
		log.Fatalf("Failed to list devices: %v", err)
	}

	// Display results
	fmt.Printf("✓ Client created successfully\n\n")
	fmt.Printf("Device List:\n")
	fmt.Printf("  Total Devices: %d\n", len(*devices))
	fmt.Printf("  Status Code: %d\n", resp.StatusCode)
	fmt.Printf("  Request Duration: %v\n", resp.Duration)
	
	if len(*devices) > 0 {
		device := (*devices)[0]
		fmt.Printf("\nFirst Device:\n")
		fmt.Printf("  Serial Number: %s\n", device.SerialNumber)
		if device.MDMUserOrDeviceName != nil {
			fmt.Printf("  Name: %s\n", *device.MDMUserOrDeviceName)
		}
		fmt.Printf("  Device Type: %s\n", device.DeviceType)
		fmt.Printf("  OS Version: %s\n", device.OSVersion)
		fmt.Printf("  Homebrew Version: %s\n", device.HomebrewVersion)
	}

	fmt.Printf("\n✓ Basic client example completed successfully!\n")
	fmt.Printf("\n💡 Security Reminder:\n")
	fmt.Printf("   Always use environment variables for API keys\n")
	fmt.Printf("   Never hardcode credentials in your source code\n")
	fmt.Printf("   Add .env files to .gitignore\n")
}
