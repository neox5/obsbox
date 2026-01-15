package metric

import "github.com/neox5/simv/value"

// MetricType defines the semantic type of a metric.
type MetricType string

const (
	MetricTypeCounter MetricType = "counter"
	MetricTypeGauge   MetricType = "gauge"
)

// Descriptor holds protocol-agnostic metric metadata and value reference.
type Descriptor struct {
	Name  string
	Type  MetricType
	Help  string
	Value value.Value[int]
}
