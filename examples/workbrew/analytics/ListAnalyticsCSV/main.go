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

	// Create HTTP client
	workbrewClient, err := workbrew.NewClient(apiKey, workspace,
		workbrew.WithLogger(logger),
		workbrew.WithBaseURL("https://console.workbrew.com"),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Create analytics service

	// List analytics as CSV
	ctx := context.Background()
	csvData, _, err := workbrewClient.Analytics.ListCSVV0(ctx)
	if err != nil {
		log.Fatalf("Failed to list analytics CSV: %v", err)
	}

	// Print CSV data
	fmt.Printf("Analytics CSV (%d bytes):\n", len(csvData))
	fmt.Println(string(csvData))
}
