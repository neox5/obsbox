package config

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
	Export    ExportConfig      `yaml:"export"`
	Settings  SettingsConfig    `yaml:"settings"`
}

// RawInstances holds all instance definitions
type RawInstances struct {
	Clocks  map[string]RawClockReference  `yaml:"clocks,omitempty"`
	Sources map[string]RawSourceReference `yaml:"sources,omitempty"`
	Values  map[string]RawValueReference  `yaml:"values,omitempty"`
}

// RawMetricConfig with polymorphic value field
type RawMetricConfig struct {
	Name        MetricNameConfig  `yaml:"name"`
	Type        string            `yaml:"type"`
	Description string            `yaml:"description"`
	Value       RawValueReference `yaml:"value"`
	Attributes  map[string]string `yaml:"attributes,omitempty"`
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
