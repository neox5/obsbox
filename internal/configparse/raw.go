package configparse

import (
	"fmt"
	"time"

	"go.yaml.in/yaml/v4"
)

// RawConfig represents unparsed YAML structure
type RawConfig struct {
	Templates RawTemplates      `yaml:"templates"`
	Instances RawInstances      `yaml:"instances"`
	Metrics   []RawMetricConfig `yaml:"metrics"`
	Export    RawExportConfig   `yaml:"export"`
	Settings  RawSettingsConfig `yaml:"settings"`
}

// RawTemplates holds all template definitions
type RawTemplates struct {
	Clocks  map[string]RawClockReference   `yaml:"clocks,omitempty"`
	Sources map[string]RawSourceReference  `yaml:"sources,omitempty"`
	Values  map[string]RawValueReference   `yaml:"values,omitempty"`
	Metrics map[string]RawMetricDefinition `yaml:"metrics,omitempty"`
}

// RawInstances holds all instance definitions
type RawInstances struct {
	Clocks  map[string]RawClockReference  `yaml:"clocks,omitempty"`
	Sources map[string]RawSourceReference `yaml:"sources,omitempty"`
	Values  map[string]RawValueReference  `yaml:"values,omitempty"`
}

// RawMetricConfig with polymorphic value field
type RawMetricConfig struct {
	Name        RawMetricNameConfig `yaml:"name"`
	Type        string              `yaml:"type"`
	Description string              `yaml:"description"`
	Value       RawValueReference   `yaml:"value"`
	Attributes  map[string]string   `yaml:"attributes,omitempty"`
}

// RawMetricDefinition - always full object, no self-reference
type RawMetricDefinition struct {
	Type       string             `yaml:"type"`
	Value      *RawValueReference `yaml:"value,omitempty"`
	Attributes map[string]string  `yaml:"attributes,omitempty"`
}

// RawMetricNameConfig supports both short and full forms for metric names
type RawMetricNameConfig struct {
	Simple     string
	Prometheus string
	OTEL       string
}

// UnmarshalYAML handles both string and object forms for metric names
func (m *RawMetricNameConfig) UnmarshalYAML(value *yaml.Node) error {
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
func (m *RawMetricNameConfig) GetPrometheusName() string {
	if m.Simple != "" {
		return m.Simple
	}
	return m.Prometheus
}

// GetOTELName returns the OTEL metric name
func (m *RawMetricNameConfig) GetOTELName() string {
	if m.Simple != "" {
		return m.Simple
	}
	return m.OTEL
}

// RawClockReference handles polymorphic clock field (instance/template/inline)
type RawClockReference struct {
	Instance string        `yaml:"instance,omitempty"`
	Template string        `yaml:"template,omitempty"`
	Type     *string       `yaml:"type,omitempty"`
	Interval time.Duration `yaml:"interval,omitempty"`
}

// RawSourceReference handles polymorphic source field (instance/template/inline)
type RawSourceReference struct {
	Instance string             `yaml:"instance,omitempty"`
	Template string             `yaml:"template,omitempty"`
	Type     *string            `yaml:"type,omitempty"`
	Clock    *RawClockReference `yaml:"clock,omitempty"`
	Min      *int               `yaml:"min,omitempty"`
	Max      *int               `yaml:"max,omitempty"`
}

// RawValueReference handles polymorphic value field (instance/template/inline)
type RawValueReference struct {
	Instance   string              `yaml:"instance,omitempty"`
	Template   string              `yaml:"template,omitempty"`
	Source     *RawSourceReference `yaml:"source,omitempty"`
	Transforms []TransformConfig   `yaml:"transforms,omitempty"`
	Reset      ResetConfig         `yaml:"reset,omitempty"`
}

// TransformConfig defines a transform operation
type TransformConfig struct {
	Type string
}

// UnmarshalYAML handles both string and object forms for transforms
func (t *TransformConfig) UnmarshalYAML(value *yaml.Node) error {
	// Try string form first (shorthand)
	var simple string
	if err := value.Decode(&simple); err == nil {
		t.Type = simple
		return nil
	}

	// Fall back to object form
	type transformConfig struct {
		Type string `yaml:"type"`
	}
	var full transformConfig
	if err := value.Decode(&full); err != nil {
		return err
	}
	t.Type = full.Type
	return nil
}

// ResetConfig defines reset behavior
type ResetConfig struct {
	Type  string
	Value int
}

// UnmarshalYAML handles both string and object forms for reset
func (r *ResetConfig) UnmarshalYAML(value *yaml.Node) error {
	// Try string form first (shorthand)
	var simple string
	if err := value.Decode(&simple); err == nil {
		r.Type = simple
		r.Value = 0 // Default value for shorthand
		return nil
	}

	// Fall back to object form
	type resetConfig struct {
		Type  string `yaml:"type"`
		Value int    `yaml:"value"`
	}
	var full resetConfig
	if err := value.Decode(&full); err != nil {
		return err
	}
	r.Type = full.Type
	r.Value = full.Value
	return nil
}

// RawExportConfig defines how metrics are exposed
type RawExportConfig struct {
	Prometheus *RawPrometheusExportConfig `yaml:"prometheus,omitempty"`
	OTEL       *RawOTELExportConfig       `yaml:"otel,omitempty"`
}

// RawPrometheusExportConfig defines Prometheus pull endpoint settings
type RawPrometheusExportConfig struct {
	Enabled bool   `yaml:"enabled"`
	Port    int    `yaml:"port"`
	Path    string `yaml:"path"`
}

// RawOTELExportConfig defines OTEL push settings
type RawOTELExportConfig struct {
	Enabled   bool              `yaml:"enabled"`
	Transport string            `yaml:"transport"`
	Host      string            `yaml:"host"`
	Port      int               `yaml:"port"`
	Interval  RawIntervalConfig `yaml:"interval"`
	Resource  map[string]string `yaml:"resource,omitempty"`
	Headers   map[string]string `yaml:"headers,omitempty"`
}

// RawIntervalConfig defines read and push intervals for OTEL
type RawIntervalConfig struct {
	Read time.Duration
	Push time.Duration
}

// UnmarshalYAML handles both simple (10s) and detailed (read/push) forms
func (i *RawIntervalConfig) UnmarshalYAML(value *yaml.Node) error {
	// Try simple duration form first
	var simple time.Duration
	if err := value.Decode(&simple); err == nil {
		i.Read = simple
		i.Push = simple
		return nil
	}

	// Fall back to detailed form
	type intervalConfig struct {
		Read time.Duration `yaml:"read"`
		Push time.Duration `yaml:"push"`
	}
	var detailed intervalConfig
	if err := value.Decode(&detailed); err != nil {
		return err
	}
	i.Read = detailed.Read
	i.Push = detailed.Push
	return nil
}

// RawSettingsConfig holds general application settings
type RawSettingsConfig struct {
	Seed            *uint64                  `yaml:"seed,omitempty"`
	InternalMetrics RawInternalMetricsConfig `yaml:"internal_metrics"`
}

// RawInternalMetricsConfig controls obsbox's self-monitoring metrics
type RawInternalMetricsConfig struct {
	Enabled bool   `yaml:"enabled"`
	Format  string `yaml:"format"`
}

// validateRawSyntax performs basic syntactic validation on raw config
func validateRawSyntax(raw *RawConfig) error {
	// Validate at least one metric defined
	if len(raw.Metrics) == 0 {
		return fmt.Errorf("at least one metric must be defined")
	}

	// Validate metric names
	for i, metric := range raw.Metrics {
		promName := metric.Name.GetPrometheusName()
		otelName := metric.Name.GetOTELName()

		if promName == "" && otelName == "" {
			return fmt.Errorf("metric at index %d: name cannot be empty", i)
		}

		if metric.Type == "" {
			return fmt.Errorf("metric %q: type cannot be empty", promName)
		}

		if metric.Description == "" {
			return fmt.Errorf("metric %q: description cannot be empty", promName)
		}
	}

	return nil
}
