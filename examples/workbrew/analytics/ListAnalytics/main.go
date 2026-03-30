package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew"
	
	"go.uber.org/zap"
)

func main() {
	// Retrieve API key and workspace from environment variables
	apiKey := os.Getenv("WORKBREW_API_KEY")
	workspace := os.Getenv("WORKBREW_WORKSPACE")

	if apiKey == "" || workspace == "" {
		log.Fatal("WORKBREW_API_KEY and WORKBREW_WORKSPACE environment variables must be set")
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	// Create Workbrew client with all services ready
	workbrewClient, err := workbrew.NewClient(apiKey, workspace,
		workbrew.WithLogger(logger),
		workbrew.WithBaseURL("https://console.workbrew.com"),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// List analytics using the built-in Analytics service
	ctx := context.Background()
	analyticsData, _, err := workbrewClient.Analytics.ListAnalytics(ctx)
	if err != nil {
		log.Fatalf("Failed to list analytics: %v", err)
	}

	fmt.Printf("Retrieved %d analytics records\n", len(*analyticsData))
	for i, analytic := range *analyticsData {
		fmt.Printf("\nAnalytic %d:\n", i+1)
		fmt.Printf("  Device: %s\n", analytic.Device)
		fmt.Printf("  Command: %s\n", analytic.Command)
		fmt.Printf("  Last Run: %s\n", analytic.LastRun)
		fmt.Printf("  Count: %d\n", analytic.Count)
	}
}
