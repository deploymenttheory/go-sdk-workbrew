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


	// Optional: Filter by actor type (user, system, or all)
	opts := &events.RequestQueryOptions{
		// Filter: "user",
	}

	ctx := context.Background()
	eventsList, _, err := workbrewClient.Events.ListEvents(ctx, opts)
	if err != nil {
		log.Fatalf("Failed to list events: %v", err)
	}

	fmt.Printf("Retrieved %d audit log events\n", len(*eventsList))
	for i, event := range *eventsList {
		fmt.Printf("\nEvent %d:\n", i+1)
		fmt.Printf("  ID: %s\n", event.ID)
		fmt.Printf("  Event Type: %s\n", event.EventType)
		fmt.Printf("  Occurred At: %s\n", event.OccurredAt)
		if event.ActorType != nil {
			fmt.Printf("  Actor Type: %s\n", *event.ActorType)
		}
		if event.TargetType != nil {
			fmt.Printf("  Target Type: %s\n", *event.TargetType)
		}
		if event.TargetIdentifier != nil {
			fmt.Printf("  Target Identifier: %s\n", *event.TargetIdentifier)
		}
	}
}
