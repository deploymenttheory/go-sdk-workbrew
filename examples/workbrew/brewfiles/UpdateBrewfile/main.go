package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew"
	
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/brewfiles"
	"go.uber.org/zap"
)

func main() {
	apiKey := os.Getenv("WORKBREW_API_KEY")
	workspace := os.Getenv("WORKBREW_WORKSPACE")
	brewfileLabel := os.Getenv("BREWFILE_LABEL") // e.g., "my-brewfile"

	if apiKey == "" || workspace == "" || brewfileLabel == "" {
		log.Fatal("WORKBREW_API_KEY, WORKBREW_WORKSPACE, and BREWFILE_LABEL environment variables must be set")
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

	// Update brewfile request
	request := &brewfiles.UpdateBrewfileRequest{
		Content: "brew \"wget\"\nbrew \"htop\"\nbrew \"curl\"",
	}

	ctx := context.Background()
	response, _, err := workbrewClient.Brewfiles.UpdateByLabelV0(ctx, brewfileLabel, request)
	if err != nil {
		log.Fatalf("Failed to update brewfile: %v", err)
	}

	fmt.Printf("Success: %s\n", response.Message)
}
