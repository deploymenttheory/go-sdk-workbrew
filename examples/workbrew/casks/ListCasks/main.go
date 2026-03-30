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
	casksList, _, err := workbrewClient.Casks.ListCasks(ctx)
	if err != nil {
		log.Fatalf("Failed to list casks: %v", err)
	}

	fmt.Printf("Retrieved %d casks\n", len(*casksList))
	for i, cask := range *casksList {
		fmt.Printf("\nCask %d:\n", i+1)
		fmt.Printf("  Name: %s\n", cask.Name)
		if cask.DisplayName != nil {
			fmt.Printf("  Display Name: %s\n", *cask.DisplayName)
		}
		fmt.Printf("  Devices: %v\n", cask.Devices)
		fmt.Printf("  Outdated: %t\n", cask.Outdated)
		if cask.Deprecated != nil {
			fmt.Printf("  Deprecated: %s\n", *cask.Deprecated)
		}
		if cask.HomebrewCaskVersion != nil {
			fmt.Printf("  Version: %s\n", *cask.HomebrewCaskVersion)
		}
	}
}
