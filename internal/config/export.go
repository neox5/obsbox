package config

import "fmt"

// ExportConfig defines how metrics are exposed.
type ExportConfig struct {
	Prometheus *PrometheusExportConfig `yaml:"prometheus,omitempty"`
	OTEL       *OTELExportConfig       `yaml:"otel,omitempty"`
}

// Validate applies defaults and validates export configuration.
func (e *ExportConfig) Validate() error {
	// Default to Prometheus enabled if no exporters configured
	if e.Prometheus == nil && e.OTEL == nil {
		e.Prometheus = &PrometheusExportConfig{
			Enabled: true,
			Port:    DefaultPrometheusPort,
			Path:    DefaultPrometheusPath,
		}
		return nil
	}

	// Validate individual exporters
	if e.Prometheus != nil && e.Prometheus.Enabled {
		if err := e.Prometheus.Validate(); err != nil {
			return err
		}
	}

	if e.OTEL != nil && e.OTEL.Enabled {
		if err := e.OTEL.Validate(); err != nil {
			return err
		}
	}

	// Verify at least one exporter enabled
	promEnabled := e.Prometheus != nil && e.Prometheus.Enabled
	otelEnabled := e.OTEL != nil && e.OTEL.Enabled

	if !promEnabled && !otelEnabled {
		return fmt.Errorf("at least one exporter must be enabled")
	}

	// Verify only one exporter enabled (prevent read conflicts)
	if promEnabled && otelEnabled {
		return fmt.Errorf("only one exporter can be enabled at a time (prometheus or otel)")
	}

	return nil
}
