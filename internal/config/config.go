package config

import (
	"time"

	"go.yaml.in/yaml/v4"
)

// Config holds the complete application configuration.
type Config struct {
	Server     ServerConfig     `yaml:"server"`
	Simulation SimulationConfig `yaml:"simulation"`
	Metrics    []MetricConfig   `yaml:"metrics"`
}

// ServerConfig defines HTTP server settings.
type ServerConfig struct {
	Port int    `yaml:"port"`
	Path string `yaml:"path"`
}

// SimulationConfig defines the simulation domain configuration.
type SimulationConfig struct {
	Clocks  map[string]ClockConfig  `yaml:"clocks"`
	Sources map[string]SourceConfig `yaml:"sources"`
	Values  map[string]ValueConfig  `yaml:"values"`
}

// ClockConfig defines a clock.
type ClockConfig struct {
	Type     string        `yaml:"type"`
	Interval time.Duration `yaml:"interval"`
}

// SourceConfig defines a simv source.
type SourceConfig struct {
	Type  string `yaml:"type"`
	Clock string `yaml:"clock"`
	Min   int    `yaml:"min,omitempty"`
	Max   int    `yaml:"max,omitempty"`
}

// TransformConfig defines a transform with optional parameters.
type TransformConfig struct {
	Type    string                 `yaml:"type"`
	Options map[string]interface{} `yaml:"options,omitempty"`
}

// ResetConfig defines reset behavior for values.
type ResetConfig struct {
	Type  string `yaml:"type,omitempty"`
	Value int    `yaml:"value,omitempty"`
}

// UnmarshalYAML handles both string and object forms for reset config.
func (r *ResetConfig) UnmarshalYAML(value *yaml.Node) error {
	// Try string form first (short form)
	var shortForm string
	if err := value.Decode(&shortForm); err == nil {
		r.Type = shortForm
		r.Value = 0 // default
		return nil
	}

	// Fall back to full form (object)
	type resetConfig ResetConfig // Avoid recursion
	var fullForm resetConfig
	if err := value.Decode(&fullForm); err != nil {
		return err
	}
	*r = ResetConfig(fullForm)
	return nil
}

// ValueConfig defines a simv value with transforms or derivation.
type ValueConfig struct {
	Source     string            `yaml:"source,omitempty"`
	Clone      string            `yaml:"clone,omitempty"`
	Transforms []TransformConfig `yaml:"transforms,omitempty"`
	Reset      ResetConfig       `yaml:"reset,omitempty"`
}

// MetricConfig defines a Prometheus metric.
type MetricConfig struct {
	Name  string `yaml:"name"`
	Type  string `yaml:"type"`
	Help  string `yaml:"help"`
	Value string `yaml:"value"`
}
