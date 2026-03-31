package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew"
)

// This example demonstrates creating a client using environment variables.
//
// Use this approach when:
// - You want to follow 12-factor app principles
// - Your API key and workspace are stored in environment variables
// - You want the convenience of automatic environment variable handling
// - You need to support optional environment-based configuration
//
// Supported environment variables:
// - WORKBREW_API_KEY (required): Your Workbrew API key
// - WORKBREW_WORKSPACE (required): Your Workbrew workspace ID
// - WORKBREW_BASE_URL (optional): Custom base URL for the API
// - WORKBREW_API_VERSION (optional): Custom API version

func main() {
	// Check required environment variables
	apiKey := os.Getenv("WORKBREW_API_KEY")
	if apiKey == "" {
		log.Fatal("WORKBREW_API_KEY environment variable is required")
	}

	workspace := os.Getenv("WORKBREW_WORKSPACE")
	if workspace == "" {
		log.Fatal("WORKBREW_WORKSPACE environment variable is required")
	}

	// Create client from environment variables
	// This automatically reads WORKBREW_API_KEY, WORKBREW_WORKSPACE,
	// and optional WORKBREW_BASE_URL and WORKBREW_API_VERSION
	client, err := workbrew.NewClientFromEnv()
	if err != nil {
		log.Fatalf("Failed to create Workbrew client: %v", err)
	}

	// Use the client to make an API call
	ctx := context.Background()

	formulae, resp, err := client.Formulae.ListV0(ctx)
	if err != nil {
		log.Fatalf("Failed to list formulae: %v", err)
	}

	// Display results
	fmt.Printf("✓ Client created from environment variables\n\n")
	fmt.Printf("Configuration:\n")
	fmt.Printf("  API Key: %s...%s (redacted)\n", apiKey[:4], apiKey[len(apiKey)-4:])
	fmt.Printf("  Workspace: %s\n", workspace)
	if baseURL := os.Getenv("WORKBREW_BASE_URL"); baseURL != "" {
		fmt.Printf("  Custom Base URL: %s\n", baseURL)
	} else {
		fmt.Printf("  Base URL: https://console.workbrew.com (default)\n")
	}
	if apiVersion := os.Getenv("WORKBREW_API_VERSION"); apiVersion != "" {
		fmt.Printf("  API Version: %s\n", apiVersion)
	} else {
		fmt.Printf("  API Version: v0 (default)\n")
	}

	fmt.Printf("\nFormulae List:\n")
	fmt.Printf("  Total Formulae: %d\n", len(*formulae))
	fmt.Printf("  Status Code: %d\n", resp.StatusCode)
	fmt.Printf("  Request Duration: %v\n", resp.Duration)

	if len(*formulae) > 0 {
		fmt.Printf("\nFirst 3 Formulae:\n")
		for i, formula := range (*formulae)[:min(3, len(*formulae))] {
			fmt.Printf("  %d. %s", i+1, formula.Name)
			if formula.HomebrewCoreVersion != nil {
				fmt.Printf(" (version: %s)", *formula.HomebrewCoreVersion)
			}
			fmt.Printf("\n")
		}
	}

	fmt.Printf("\n✓ Environment-based client example completed successfully!\n")
	fmt.Printf("\n💡 Best Practices:\n")
	fmt.Printf("   Set WORKBREW_API_KEY and WORKBREW_WORKSPACE in your environment\n")
	fmt.Printf("   Use .env files for local development (add to .gitignore)\n")
	fmt.Printf("   Use secrets management in production (AWS, Vault, etc.)\n")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
