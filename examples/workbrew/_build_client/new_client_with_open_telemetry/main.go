package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew"
	
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/zap"
)

// This example demonstrates creating a client with OpenTelemetry tracing for observability.
//
// Use this approach when:
// - You need distributed tracing across microservices
// - You want to monitor API performance and latency
// - You need to track errors and failures in production
// - You're using observability platforms (Jaeger, Zipkin, DataDog, etc.)
// - You want complete visibility into API call chains
//
// This example shows:
// - OpenTelemetry tracer provider setup
// - Client instrumentation with tracing
// - Automatic span creation for HTTP requests
// - Trace export to stdout (replace with your backend)
// - Combined logging and tracing for full observability

func main() {
	// Check API key and workspace from environment
	apiKey := os.Getenv("WORKBREW_API_KEY")
	if apiKey == "" {
		log.Fatal("WORKBREW_API_KEY environment variable is required")
	}

	workspace := os.Getenv("WORKBREW_WORKSPACE")
	if workspace == "" {
		log.Fatal("WORKBREW_WORKSPACE environment variable is required")
	}

	// Step 1: Create structured logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	// Step 2: Initialize OpenTelemetry exporter
	// In production, replace stdout exporter with:
	// - OTLP exporter for OpenTelemetry Collector
	// - Jaeger exporter for Jaeger
	// - Zipkin exporter for Zipkin
	// - DataDog, Honeycomb, New Relic, etc.
	exporter, err := stdouttrace.New(
		stdouttrace.WithPrettyPrint(),
	)
	if err != nil {
		log.Fatalf("Failed to create trace exporter: %v", err)
	}

	// Step 3: Create tracer provider
	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		// Add resource attributes for better trace context
		// trace.WithResource(resource.NewWithAttributes(
		//     semconv.SchemaURL,
		//     semconv.ServiceName("workbrew-client"),
		//     semconv.ServiceVersion("1.0.0"),
		// )),
	)
	defer func() {
		if err := tracerProvider.Shutdown(context.Background()); err != nil {
			logger.Error("Failed to shutdown tracer provider", zap.Error(err))
		}
	}()

	// Set as global tracer provider
	otel.SetTracerProvider(tracerProvider)

	// Step 4: Create client with tracing enabled
	workbrewClient, err := workbrew.NewClient(
		apiKey,
		workspace,

		// Enable structured logging
		workbrew.WithLogger(logger),

		// Enable OpenTelemetry tracing - this automatically instruments all HTTP calls
		workbrew.WithTracing(nil), // nil uses default config with global tracer provider

		// Or use custom tracing configuration:
		// workbrew.WithTracing(&workbrew.OTelConfig{
		//     TracerProvider: tracerProvider,
		//     ServiceName:    "my-workbrew-app",
		//     SpanNameFormatter: func(operation string, req *http.Request) string {
		//         return fmt.Sprintf("Workbrew: %s %s", req.Method, req.URL.Path)
		//     },
		// }),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	logger.Info("Workbrew client created with OpenTelemetry tracing",
		zap.String("workspace", workspace))

	// Step 5: Use the client - tracing happens automatically!
	ctx := context.Background()

	// This API call will automatically create spans with:
	// - HTTP method, URL, and status code
	// - Request/response timing
	// - Error details (if any)
	// - All OpenTelemetry semantic conventions
	logger.Info("Fetching brewfiles")

	brewfiles, resp, err := workbrewClient.Brewfiles.ListV0(ctx)
	if err != nil {
		logger.Error("Brewfile list failed",
			zap.Error(err),
			zap.Int("status_code", resp.StatusCode))
		log.Fatalf("API call failed: %v", err)
	}

	logger.Info("Brewfiles retrieved",
		zap.Int("brewfile_count", len(*brewfiles)),
		zap.Int("status_code", resp.StatusCode),
		zap.Duration("duration", resp.Duration))

	// Display results
	fmt.Printf("\n✓ Client created with OpenTelemetry tracing\n\n")
	fmt.Printf("Observability Setup:\n")
	fmt.Printf("  Tracing: Enabled (OpenTelemetry)\n")
	fmt.Printf("  Exporter: stdout (replace with your backend)\n")
	fmt.Printf("  Logging: zap (structured)\n")
	fmt.Printf("  Service: workbrew-client\n")

	fmt.Printf("\nBrewfiles:\n")
	fmt.Printf("  Total Brewfiles: %d\n", len(*brewfiles))
	fmt.Printf("  Status Code: %d\n", resp.StatusCode)
	fmt.Printf("  Duration: %v\n", resp.Duration)

	if len(*brewfiles) > 0 {
		fmt.Printf("\nFirst Brewfile:\n")
		brewfile := (*brewfiles)[0]
		fmt.Printf("  Label: %s\n", brewfile.Label)
		if brewfile.DeviceCount != nil {
			fmt.Printf("  Device Count: %d\n", *brewfile.DeviceCount)
		}
		if brewfile.CreatedAt != nil {
			fmt.Printf("  Created: %s\n", brewfile.CreatedAt.String())
		}
	}

	fmt.Printf("\n📊 Trace Information:\n")
	fmt.Printf("Check the output above for detailed trace spans.\n")
	fmt.Printf("Each HTTP request is automatically instrumented with:\n")
	fmt.Printf("  - HTTP method and URL\n")
	fmt.Printf("  - Request/response timing\n")
	fmt.Printf("  - Status codes and errors\n")
	fmt.Printf("  - OpenTelemetry semantic conventions\n")

	fmt.Printf("\n✓ OpenTelemetry client example completed successfully!\n")
	fmt.Printf("\nNext Steps:\n")
	fmt.Printf("  1. Replace stdout exporter with your observability backend\n")
	fmt.Printf("  2. Add resource attributes for better trace context\n")
	fmt.Printf("  3. Configure sampling for production workloads\n")
	fmt.Printf("  4. View traces in your observability platform\n")
}
