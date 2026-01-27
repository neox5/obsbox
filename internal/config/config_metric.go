package config

import (
	"regexp"
	"strings"
	"time"

	"go.yaml.in/yaml/v4"
)

var attributeNameRegex = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

// MetricConfig defines a fully resolved metric
type MetricConfig struct {
	PrometheusName string
	OTELName       string
	Type           MetricType
	Description    string
	Value          ValueConfig
	Attributes     map[string]string
}

// MetricType defines the semantic type of a metric
type MetricType string

const (
	MetricTypeCounter MetricType = "counter"
	MetricTypeGauge   MetricType = "gauge"
)

// MetricNameConfig supports both short and full forms for metric names
type MetricNameConfig struct {
	Simple     string
	Prometheus string
	OTEL       string
}

// UnmarshalYAML handles both string and object forms for metric names
func (m *MetricNameConfig) UnmarshalYAML(value *yaml.Node) error {
	// Try string form first (short form)
	var simple string
	if err := value.Decode(&simple); err == nil {
		m.Simple = simple
		return nil
	}

	// Try full form (object)
	type nameConfig struct {
		Prometheus string `yaml:"prometheus"`
		OTEL       string `yaml:"otel"`
	}
	var full nameConfig
	if err := value.Decode(&full); err != nil {
		return err
	}
	m.Prometheus = full.Prometheus
	m.OTEL = full.OTEL
	return nil
}

// GetPrometheusName returns the Prometheus metric name
func (m *MetricNameConfig) GetPrometheusName() string {
	if m.Simple != "" {
		return m.Simple
	}
	return m.Prometheus
}

// GetOTELName returns the OTEL metric name
func (m *MetricNameConfig) GetOTELName() string {
	if m.Simple != "" {
		return m.Simple
	}
	return m.OTEL
}

// isValidAttributeName checks if an attribute name follows conventions
func isValidAttributeName(name string) bool {
	if len(name) == 0 {
		return false
	}
	if strings.HasPrefix(name, "__") {
		return false
	}
	return attributeNameRegex.MatchString(name)
}

// ValueConfig defines a fully resolved value with embedded components.
// Clock is optional - if not specified, inherited from Source.Clock.
// If specified, overrides Source.Clock.
type ValueConfig struct {
	Source     SourceConfig
	Clock      ClockConfig // Optional - overrides Source.Clock if specified
	Transforms []TransformConfig
	Reset      ResetConfig
}

// SourceConfig defines a fully resolved source with embedded clock
type SourceConfig struct {
	Type  string
	Clock ClockConfig
	Min   int
	Max   int
}

// ClockConfig defines a fully resolved clock
type ClockConfig struct {
	Type     string
	Interval time.Duration
}
