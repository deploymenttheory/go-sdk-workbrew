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
	brewCommandLabel := os.Getenv("BREW_COMMAND_LABEL") // e.g., "outdated"

	if apiKey == "" || workspace == "" || brewCommandLabel == "" {
		log.Fatal("WORKBREW_API_KEY, WORKBREW_WORKSPACE, and BREW_COMMAND_LABEL environment variables must be set")
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
	csvData, _, err := workbrewClient.BrewCommands.ListRunsByLabelCSVV0(ctx, brewCommandLabel)
	if err != nil {
		log.Fatalf("Failed to list brew command runs CSV: %v", err)
	}

	fmt.Printf("Brew Command Runs CSV for '%s' (%d bytes):\n", brewCommandLabel, len(csvData))
	fmt.Println(string(csvData))
}
