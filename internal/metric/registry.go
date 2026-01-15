package metric

import (
	"fmt"

	"github.com/neox5/obsbox/internal/config"
	"github.com/neox5/obsbox/internal/generator"
)

// Registry holds protocol-agnostic metric definitions.
type Registry struct {
	metrics []Descriptor
}

// New creates a registry from configuration.
func New(cfg *config.Config, gen *generator.Generator) (*Registry, error) {
	var metrics []Descriptor

	for _, metricCfg := range cfg.Metrics {
		val, exists := gen.GetValue(metricCfg.Value)
		if !exists {
			return nil, fmt.Errorf("value %q not found for metric %q", metricCfg.Value, metricCfg.Name)
		}

		metrics = append(metrics, Descriptor{
			Name:  metricCfg.Name,
			Type:  MetricType(metricCfg.Type),
			Help:  metricCfg.Help,
			Value: val,
		})
	}

	return &Registry{metrics: metrics}, nil
}

// Metrics returns all registered metric descriptors.
func (r *Registry) Metrics() []Descriptor {
	return r.metrics
}
