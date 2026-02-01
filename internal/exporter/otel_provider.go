package exporter

import (
	"context"
	"fmt"

	"github.com/neox5/otelbox/internal/config"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
)

// createMeterProvider creates an OTEL meter provider with OTLP exporter.
func createMeterProvider(
	cfg *config.OTELExportConfig,
	res *resource.Resource,
) (*sdkmetric.MeterProvider, error) {
	// Create exporter based on transport type
	var exporter sdkmetric.Exporter
	var err error

	switch cfg.Transport {
	case "grpc":
		exporter, err = createGRPCExporter(cfg)
	case "http":
		exporter, err = createHTTPExporter(cfg)
	default:
		return nil, fmt.Errorf("unsupported transport: %s", cfg.Transport)
	}

	if err != nil {
		return nil, err
	}

	// Create periodic reader with push interval
	reader := sdkmetric.NewPeriodicReader(
		exporter,
		sdkmetric.WithInterval(cfg.Interval.Push),
	)

	// Create meter provider
	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(reader),
	)

	return meterProvider, nil
}

// createGRPCExporter creates an OTLP gRPC exporter.
func createGRPCExporter(cfg *config.OTELExportConfig) (sdkmetric.Exporter, error) {
	opts := []otlpmetricgrpc.Option{
		otlpmetricgrpc.WithEndpoint(cfg.GetEndpoint()),
		otlpmetricgrpc.WithInsecure(), // TODO: Add TLS support later
	}

	// Add custom headers
	if len(cfg.Headers) > 0 {
		opts = append(opts, otlpmetricgrpc.WithHeaders(cfg.Headers))
	}

	exporter, err := otlpmetricgrpc.New(context.Background(), opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP gRPC exporter: %w", err)
	}

	return exporter, nil
}

// createHTTPExporter creates an OTLP HTTP exporter.
func createHTTPExporter(cfg *config.OTELExportConfig) (sdkmetric.Exporter, error) {
	opts := []otlpmetrichttp.Option{
		otlpmetrichttp.WithEndpoint(cfg.GetEndpoint()),
		otlpmetrichttp.WithInsecure(), // TODO: Add TLS support later
	}

	// Add custom headers
	if len(cfg.Headers) > 0 {
		opts = append(opts, otlpmetrichttp.WithHeaders(cfg.Headers))
	}

	exporter, err := otlpmetrichttp.New(context.Background(), opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP HTTP exporter: %w", err)
	}

	return exporter, nil
}
