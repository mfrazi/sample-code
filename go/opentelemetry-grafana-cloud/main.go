package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

func main() {
	ctx := context.Background()

	// Init opentelemetry
	provider, err := InitOTLPMetricsExporter(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize OTLP exporter: %v", err)
	}

	defer func() {
		if err := provider.Shutdown(ctx); err != nil {
			log.Printf("Failed to shut down MeterProvider: %v", err)
		}
	}()

	// HTTP Server
	http.HandleFunc("/divide", handleDivide)

	port := ":8080"
	fmt.Printf("Server is running on http://localhost%s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
