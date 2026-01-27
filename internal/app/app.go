package app

import (
	"fmt"

	"github.com/neox5/obsbox/internal/config"
	"github.com/neox5/obsbox/internal/exporter"
	"github.com/neox5/obsbox/internal/generator"
	"github.com/neox5/obsbox/internal/metric"
	"github.com/neox5/obsbox/internal/simulation"
)

// App holds initialized application components.
type App struct {
	Config             *config.Config
	Generator          *generator.Generator
	Metrics            *metric.Registry
	PrometheusExporter *exporter.PrometheusExporter
	OTELExporter       *exporter.OTELExporter
}

// New initializes the application from configuration.
// Seed must be initialized before calling this function.
func New(cfg *config.Config) (*App, error) {
	// Initialize seed before creating any simv objects
	simulation.InitializeSeed(&cfg.Settings)

	// Create generator from metrics
	gen, err := generator.New(cfg.Metrics)
	if err != nil {
		return nil, fmt.Errorf("failed to create generator: %w", err)
	}

	// Create metrics
	metrics, err := metric.New(cfg, gen)
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics: %w", err)
	}

	var promExporter *exporter.PrometheusExporter
	var otelExporter *exporter.OTELExporter

	// Create Prometheus exporter if enabled
	if cfg.Export.Prometheus != nil && cfg.Export.Prometheus.Enabled {
		promExporter = exporter.NewPrometheusExporter(
			cfg.Export.Prometheus.Port,
			cfg.Export.Prometheus.Path,
			metrics,
			cfg.Settings.InternalMetrics.Enabled,
		)
	}

	// Create OTEL exporter if enabled
	if cfg.Export.OTEL != nil && cfg.Export.OTEL.Enabled {
		otelExporter, err = exporter.NewOTELExporter(
			cfg.Export.OTEL,
			metrics,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create OTEL exporter: %w", err)
		}
	}

	return &App{
		Config:             cfg,
		Generator:          gen,
		Metrics:            metrics,
		PrometheusExporter: promExporter,
		OTELExporter:       otelExporter,
	}, nil
}
