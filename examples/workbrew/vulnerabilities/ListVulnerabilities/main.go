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
	vulnList, _, err := workbrewClient.Vulnerabilities.ListVulnerabilities(ctx)
	if err != nil {
		log.Fatalf("Failed to list vulnerabilities: %v", err)
	}

	fmt.Printf("Retrieved %d vulnerabilities\n", len(*vulnList))
	for i, vuln := range *vulnList {
		fmt.Printf("\nVulnerability Group %d:\n", i+1)
		fmt.Printf("  Formula: %s\n", vuln.Formula)
		fmt.Printf("  Outdated Devices: %v\n", vuln.OutdatedDevices)
		fmt.Printf("  Supported: %t\n", vuln.Supported)
		fmt.Printf("  Version: %s\n", vuln.HomebrewCoreVersion)
		fmt.Printf("  Vulnerabilities:\n")
		for j, v := range vuln.Vulnerabilities {
			if v.CVSSScore != nil {
				fmt.Printf("    %d. %s (CVSS: %.1f)\n", j+1, v.CleanID, *v.CVSSScore)
			} else {
				fmt.Printf("    %d. %s (CVSS: N/A)\n", j+1, v.CleanID)
			}
		}
	}
}
