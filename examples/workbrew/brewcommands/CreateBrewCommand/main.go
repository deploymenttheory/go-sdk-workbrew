package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew"
	
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/brewcommands"
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


	// Create brew command request
	recurrence := "once"
	request := &brewcommands.CreateBrewCommandRequest{
		Arguments:  "install wget",
		Recurrence: &recurrence,
		// DeviceIDs: nil to run on all devices
		// RunAfterDatetime: nil to run immediately
		// To run on specific devices, use: DeviceIDs: ptr("device-uuid-1,device-uuid-2")
	}

	ctx := context.Background()
	response, _, err := workbrewClient.BrewCommands.CreateBrewCommand(ctx, request)
	if err != nil {
		log.Fatalf("Failed to create brew command: %v", err)
	}

	fmt.Printf("Success: %s\n", response.Message)
}
