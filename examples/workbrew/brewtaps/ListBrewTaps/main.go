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
	taps, _, err := workbrewClient.BrewTaps.ListBrewTaps(ctx)
	if err != nil {
		log.Fatalf("Failed to list brew taps: %v", err)
	}

	fmt.Printf("Retrieved %d brew taps\n", len(*taps))
	for i, tap := range *taps {
		fmt.Printf("\nTap %d:\n", i+1)
		fmt.Printf("  Tap: %s\n", tap.Tap)
		fmt.Printf("  Devices: %v\n", tap.Devices)
		fmt.Printf("  Formulae Installed: %d\n", tap.FormulaeInstalled)
		fmt.Printf("  Casks Installed: %d\n", tap.CasksInstalled)
		fmt.Printf("  Available Packages: %s\n", tap.AvailablePackages)
	}
}
