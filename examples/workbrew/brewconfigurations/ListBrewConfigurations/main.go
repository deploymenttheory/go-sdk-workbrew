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

	ctx := context.Background()
	configs, _, err := workbrewClient.BrewConfigurations.ListBrewConfigurations(ctx)
	if err != nil {
		log.Fatalf("Failed to list brew configurations: %v", err)
	}

	fmt.Printf("Retrieved %d brew configurations\n", len(*configs))
	for i, config := range *configs {
		fmt.Printf("\nConfiguration %d:\n", i+1)
		fmt.Printf("  Key: %s\n", config.Key)
		fmt.Printf("  Value: %s\n", config.Value)
		fmt.Printf("  Last Updated By: %s\n", config.LastUpdatedByUser)
		fmt.Printf("  Device Group: %s\n", config.DeviceGroup)
	}
}
