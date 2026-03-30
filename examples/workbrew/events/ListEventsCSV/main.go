package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew"
	
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/events"
	"go.uber.org/zap"
)

func main() {
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

	workbrewClient, err := workbrew.NewClient(apiKey, workspace,
		workbrew.WithLogger(logger),
		workbrew.WithBaseURL("https://console.workbrew.com"),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Optional: Filter by actor type or add download parameter
	opts := &events.RequestQueryOptions{
		// Filter: "user",
		// Download: "1",
	}

	ctx := context.Background()
	csvData, _, err := workbrewClient.Events.ListEventsCSV(ctx, opts)
	if err != nil {
		log.Fatalf("Failed to list events CSV: %v", err)
	}

	fmt.Printf("Events CSV (%d bytes):\n", len(csvData))
	fmt.Println(string(csvData))
}
