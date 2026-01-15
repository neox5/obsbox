package exporter

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/neox5/obsbox/internal/metric"
	"github.com/neox5/simv/value"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// PrometheusExporter provides HTTP server for Prometheus metrics.
type PrometheusExporter struct {
	addr         string
	path         string
	server       *http.Server
	promRegistry *prometheus.Registry
}

// metricDescriptor holds metadata for a Prometheus metric.
type metricDescriptor struct {
	desc      *prometheus.Desc
	valueType prometheus.ValueType
	value     value.Value[int]
}

// collector implements prometheus.Collector to read simv values on scrape.
type collector struct {
	descriptors []metricDescriptor
}

// NewPrometheusExporter creates a new Prometheus HTTP exporter.
func NewPrometheusExporter(port int, path string, metrics *metric.Registry) *PrometheusExporter {
	promRegistry := prometheus.NewRegistry()

	// Build Prometheus-specific descriptors
	var descriptors []metricDescriptor
	for _, m := range metrics.Metrics() {
		var valueType prometheus.ValueType
		switch m.Type {
		case metric.MetricTypeCounter:
			valueType = prometheus.CounterValue
		case metric.MetricTypeGauge:
			valueType = prometheus.GaugeValue
		}

		descriptors = append(descriptors, metricDescriptor{
			desc:      prometheus.NewDesc(m.Name, m.Help, nil, nil),
			valueType: valueType,
			value:     m.Value,
		})

		slog.Info("registered prometheus metric", "name", m.Name, "type", m.Type)
	}

	// Register collector
	c := &collector{descriptors: descriptors}
	promRegistry.MustRegister(c)

	// Setup HTTP server
	mux := http.NewServeMux()
	addr := fmt.Sprintf(":%d", port)

	mux.Handle(path, promhttp.HandlerFor(
		promRegistry,
		promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		},
	))

	return &PrometheusExporter{
		addr:         addr,
		path:         path,
		promRegistry: promRegistry,
		server: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
	}
}

// Start begins serving HTTP requests.
func (e *PrometheusExporter) Start(ctx context.Context) error {
	errChan := make(chan error, 1)

	go func() {
		slog.Info("starting prometheus exporter", "addr", e.addr, "path", e.path)
		if err := e.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return e.Stop()
	}
}

// Stop gracefully stops the exporter.
func (e *PrometheusExporter) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	slog.Info("shutting down prometheus exporter")
	return e.server.Shutdown(ctx)
}

// Describe sends metric descriptors to the channel.
func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range c.descriptors {
		ch <- m.desc
	}
}

// Collect reads simv values and sends metrics to the channel.
// This is called on each Prometheus scrape.
func (c *collector) Collect(ch chan<- prometheus.Metric) {
	for _, m := range c.descriptors {
		// Read value from simv (may trigger reset for reset_on_read)
		val := float64(m.value.Value())

		// Create and send metric with current value
		metric, err := prometheus.NewConstMetric(
			m.desc,
			m.valueType,
			val,
		)
		if err != nil {
			continue
		}

		ch <- metric
	}
}
