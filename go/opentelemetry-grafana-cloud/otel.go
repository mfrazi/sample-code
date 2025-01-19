package main

import (
	"context"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/semconv/v1.12.0"
)

const (
	grafanaCloudOTLPHTTPEndpoint = "" // Replace with OTEL_EXPORTER_OTLP_ENDPOINT
	grafanaAuthorization         = "" // Replace with OTEL_EXPORTER_OTLP_HEADERS - Authorization (Something like Basic xkjfasdi...)
)

func InitOTLPMetricsExporter(ctx context.Context) (*sdkmetric.MeterProvider, error) {
	httpOpts := []otlpmetrichttp.Option{
		otlpmetrichttp.WithEndpointURL(grafanaCloudOTLPHTTPEndpoint),
		otlpmetrichttp.WithHeaders(map[string]string{
			"Authorization": grafanaAuthorization,
		}),
	}
	exporter, err := otlpmetrichttp.New(ctx, httpOpts...)
	if err != nil {
		return nil, err
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("application"),
			semconv.ServiceVersionKey.String("1.0.0"),
		),
	)
	if err != nil {
		return nil, err
	}

	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithInterval(5*time.Second))),
		sdkmetric.WithResource(res),
	)

	otel.SetMeterProvider(provider)

	return provider, nil
}

func pushHistogram(ctx context.Context, latency int64, tags map[string]string) {
	meter := otel.Meter("http")
	histogram, err := meter.Int64Histogram("http_handler")
	if err != nil {
		log.Printf("Failed to create histogram: %v", err)
		return
	}

	var labels []attribute.KeyValue
	for key, value := range tags {
		labels = append(labels, attribute.String(key, value))
	}

	histogram.Record(ctx, latency, metric.WithAttributes(labels...))
}
