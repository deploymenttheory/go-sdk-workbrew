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
	formulaeList, _, err := workbrewClient.Formulae.ListFormulae(ctx)
	if err != nil {
		log.Fatalf("Failed to list formulae: %v", err)
	}

	fmt.Printf("Retrieved %d formulae\n", len(*formulaeList))
	for i, formula := range *formulaeList {
		fmt.Printf("\nFormula %d:\n", i+1)
		fmt.Printf("  Name: %s\n", formula.Name)
		fmt.Printf("  Devices: %v\n", formula.Devices)
		fmt.Printf("  Outdated: %t\n", formula.Outdated)
		fmt.Printf("  Installed On Request: %t\n", formula.InstalledOnRequest)
		fmt.Printf("  Installed As Dependency: %t\n", formula.InstalledAsDependency)
		fmt.Printf("  Vulnerabilities: %v\n", formula.Vulnerabilities)
		if formula.Deprecated != nil {
			fmt.Printf("  Deprecated: %s\n", *formula.Deprecated)
		}
		if formula.License != nil {
			fmt.Printf("  License: %v\n", *formula.License)
		}
		if formula.HomebrewCoreVersion != nil {
			fmt.Printf("  Version: %s\n", *formula.HomebrewCoreVersion)
		}
	}
}
