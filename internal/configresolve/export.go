package configresolve

import (
	"github.com/neox5/obsbox/internal/config"
	"github.com/neox5/obsbox/internal/configparse"
)

// resolveExport converts raw export config to resolved export config
func resolveExport(raw *configparse.RawExportConfig) (config.ExportConfig, error) {
	result := config.ExportConfig{}

	// Convert Prometheus config if present
	if raw.Prometheus != nil {
		result.Prometheus = &config.PrometheusExportConfig{
			Enabled: raw.Prometheus.Enabled,
			Port:    raw.Prometheus.Port,
			Path:    raw.Prometheus.Path,
		}
	}

	// Convert OTEL config if present
	if raw.OTEL != nil {
		result.OTEL = &config.OTELExportConfig{
			Enabled:   raw.OTEL.Enabled,
			Transport: raw.OTEL.Transport,
			Host:      raw.OTEL.Host,
			Port:      raw.OTEL.Port,
			Interval: config.IntervalConfig{
				Read: raw.OTEL.Interval.Read,
				Push: raw.OTEL.Interval.Push,
			},
			Resource: copyStringMap(raw.OTEL.Resource),
			Headers:  copyStringMap(raw.OTEL.Headers),
		}
	}

	// Validate converted config
	if err := result.Validate(); err != nil {
		return config.ExportConfig{}, err
	}

	return result, nil
}

// copyStringMap creates a copy of a string map (handles nil)
func copyStringMap(src map[string]string) map[string]string {
	if src == nil {
		return nil
	}
	dst := make(map[string]string, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
