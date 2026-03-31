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
	devicesList, _, err := workbrewClient.Devices.ListV0(ctx)
	if err != nil {
		log.Fatalf("Failed to list devices: %v", err)
	}

	fmt.Printf("Retrieved %d devices\n", len(*devicesList))
	for i, device := range *devicesList {
		fmt.Printf("\nDevice %d:\n", i+1)
		fmt.Printf("  Serial Number: %s\n", device.SerialNumber)
		if device.MDMUserOrDeviceName != nil {
			fmt.Printf("  MDM Name: %s\n", *device.MDMUserOrDeviceName)
		}
		fmt.Printf("  Device Type: %s\n", device.DeviceType)
		fmt.Printf("  OS Version: %s\n", device.OSVersion)
		fmt.Printf("  Homebrew Version: %s\n", device.HomebrewVersion)
		fmt.Printf("  Workbrew Version: %s\n", device.WorkbrewVersion)
		fmt.Printf("  Last Seen: %s\n", device.LastSeenAt.String())
		fmt.Printf("  Command Last Run: %s\n", device.CommandLastRunAt.String())
		fmt.Printf("  Formulae Count: %d\n", device.FormulaeCount)
		fmt.Printf("  Casks Count: %d\n", device.CasksCount)
	}
}
