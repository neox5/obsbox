package metric

import (
	"github.com/neox5/simv/value"
	"github.com/prometheus/client_golang/prometheus"
)

// metricDescriptor holds metadata for a Prometheus metric.
type metricDescriptor struct {
	desc      *prometheus.Desc
	valueType prometheus.ValueType
	value     value.Value[int]
}

// Collector implements prometheus.Collector to read simv values on scrape.
type Collector struct {
	metrics []metricDescriptor
}

// NewCollector creates a collector with metric descriptors and simv value references.
func NewCollector(metrics []metricDescriptor) *Collector {
	return &Collector{
		metrics: metrics,
	}
}

// Describe sends metric descriptors to the channel.
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range c.metrics {
		ch <- m.desc
	}
}

// Collect reads simv values and sends metrics to the channel.
// This is called on each Prometheus scrape.
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	for _, m := range c.metrics {
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
