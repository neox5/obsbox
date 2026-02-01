package exporter

import (
	"github.com/neox5/otelbox/internal/metric"
	"github.com/prometheus/client_golang/prometheus"
)

// createPrometheusRegistry creates and populates a Prometheus registry.
func createPrometheusRegistry(metrics *metric.Registry) *prometheus.Registry {
	promRegistry := prometheus.NewRegistry()

	// Create and register collector
	c := newCollector(metrics)
	promRegistry.MustRegister(c)

	return promRegistry
}
