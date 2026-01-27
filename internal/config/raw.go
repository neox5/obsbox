package config

import (
	"fmt"
	"time"

	"go.yaml.in/yaml/v4"
)

// RawConfig represents unparsed YAML structure
type RawConfig struct {
	Templates RawTemplates      `yaml:"templates"`
	Metrics   []RawMetricConfig `yaml:"metrics"`
	Export    ExportConfig      `yaml:"export"`
	Settings  SettingsConfig    `yaml:"settings"`
}

// RawMetricConfig with polymorphic value field
type RawMetricConfig struct {
	Name        MetricNameConfig  `yaml:"name"`
	Type        string            `yaml:"type"`
	Description string            `yaml:"description"`
	Value       RawValueReference `yaml:"value"`
	Attributes  map[string]string `yaml:"attributes,omitempty"`
}

// RawValueReference handles polymorphic value field (string or object)
type RawValueReference struct {
	Template string          // Set if string form used
	Inline   *RawValueConfig // Set if object form used
}

func (r *RawValueReference) UnmarshalYAML(value *yaml.Node) error {
	// Try string form first
	var template string
	if err := value.Decode(&template); err == nil {
		if template == "" {
			return fmt.Errorf("template reference cannot be empty string")
		}
		r.Template = template
		return nil
	}

	// Fall back to object form
	var inline RawValueConfig
	if err := value.Decode(&inline); err != nil {
		return fmt.Errorf("invalid value config: %w", err)
	}
	r.Inline = &inline
	return nil
}

// RawValueConfig for metric value with template + overrides
type RawValueConfig struct {
	Template   string            `yaml:"template,omitempty"`
	Source     *RawSourceConfig  `yaml:"source,omitempty"`
	Clock      *RawClockConfig   `yaml:"clock,omitempty"`
	Transforms []TransformConfig `yaml:"transforms,omitempty"`
	Reset      ResetConfig       `yaml:"reset,omitempty"`
}

// RawSourceConfig handles polymorphic source field (string or object)
type RawSourceConfig struct {
	Template string           // Set if string form used
	Inline   *RawSourceInline // Set if object form used
}

func (r *RawSourceConfig) UnmarshalYAML(value *yaml.Node) error {
	// Try string form first
	var template string
	if err := value.Decode(&template); err == nil {
		if template == "" {
			return fmt.Errorf("template reference cannot be empty string")
		}
		r.Template = template
		return nil
	}

	// Fall back to object form
	var inline RawSourceInline
	if err := value.Decode(&inline); err != nil {
		return fmt.Errorf("invalid source config: %w", err)
	}
	r.Inline = &inline
	return nil
}

// RawSourceInline for inline source definition
type RawSourceInline struct {
	Type  string          `yaml:"type"`
	Clock *RawClockConfig `yaml:"clock,omitempty"`
	Min   *int            `yaml:"min,omitempty"`
	Max   *int            `yaml:"max,omitempty"`
}

// RawClockConfig handles polymorphic clock field (string or object)
type RawClockConfig struct {
	Template string          // Set if string form used
	Inline   *RawClockInline // Set if object form used
}

func (r *RawClockConfig) UnmarshalYAML(value *yaml.Node) error {
	// Try string form first
	var template string
	if err := value.Decode(&template); err == nil {
		if template == "" {
			return fmt.Errorf("template reference cannot be empty string")
		}
		r.Template = template
		return nil
	}

	// Fall back to object form
	var inline RawClockInline
	if err := value.Decode(&inline); err != nil {
		return fmt.Errorf("invalid clock config: %w", err)
	}
	r.Inline = &inline
	return nil
}

// RawClockInline for inline clock definition
type RawClockInline struct {
	Type     string        `yaml:"type"`
	Interval time.Duration `yaml:"interval"`
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
