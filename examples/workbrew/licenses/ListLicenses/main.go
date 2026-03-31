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
	licensesList, _, err := workbrewClient.Licenses.ListV0(ctx)
	if err != nil {
		log.Fatalf("Failed to list licenses: %v", err)
	}

	fmt.Printf("Retrieved %d licenses\n", len(*licensesList))
	for i, license := range *licensesList {
		fmt.Printf("\nLicense %d:\n", i+1)
		fmt.Printf("  Name: %s\n", license.Name)
		fmt.Printf("  Device Count: %d\n", license.DeviceCount)
		fmt.Printf("  Formula Count: %d\n", license.FormulaCount)
	}
}
