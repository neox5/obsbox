package metric

import (
	"fmt"

	"github.com/neox5/otelbox/internal/config"
	"github.com/neox5/otelbox/internal/generator"
)

// Registry holds protocol-agnostic metric definitions.
type Registry struct {
	metrics []Descriptor
}

// New creates a registry from configuration.
func New(cfg *config.Config, gen *generator.Generator) (*Registry, error) {
	var metrics []Descriptor

	for i, metricCfg := range cfg.Metrics {
		val := gen.GetValue(i)
		if val == nil {
			return nil, fmt.Errorf("metric %d (%s): value not found",
				i, metricCfg.PrometheusName)
		}

		metrics = append(metrics, Descriptor{
			PrometheusName: metricCfg.PrometheusName,
			OTELName:       metricCfg.OTELName,
			Type:           MetricType(metricCfg.Type),
			Description:    metricCfg.Description,
			Attributes:     metricCfg.Attributes,
			Value:          val.Value,
		})
	}

	return &Registry{metrics: metrics}, nil
}

// Metrics returns all registered metric descriptors.
func (r *Registry) Metrics() []Descriptor {
	return r.metrics
}
