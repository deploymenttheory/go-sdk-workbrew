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
	runs, _, err := workbrewClient.BrewCommands.ListRunsByLabelV0(ctx, brewCommandLabel)
	if err != nil {
		log.Fatalf("Failed to list brew command runs: %v", err)
	}

	fmt.Printf("Retrieved %d runs for brew command '%s'\n", len(*runs), brewCommandLabel)
	for i, run := range *runs {
		fmt.Printf("\nRun %d:\n", i+1)
		fmt.Printf("  Command: %s\n", run.Command)
		fmt.Printf("  Label: %s\n", run.Label)
		fmt.Printf("  Device: %s\n", run.Device)
		fmt.Printf("  Success: %t\n", run.Success)
		fmt.Printf("  Started At: %s\n", run.StartedAt.String())
		fmt.Printf("  Finished At: %s\n", run.FinishedAt.String())
	}
}
