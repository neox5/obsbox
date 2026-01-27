package config

import (
	"time"
)

// RawTemplates holds all template definitions
type RawTemplates struct {
	Clocks  map[string]RawClockTemplate  `yaml:"clocks,omitempty"`
	Sources map[string]RawSourceTemplate `yaml:"sources,omitempty"`
	Values  map[string]RawValueTemplate  `yaml:"values,omitempty"`
	Metrics map[string]RawMetricTemplate `yaml:"metrics,omitempty"`
}

// RawClockTemplate - always full object, no self-reference
type RawClockTemplate struct {
	Type     string        `yaml:"type"`
	Interval time.Duration `yaml:"interval"`
}

// RawSourceTemplate - always full object, no self-reference
type RawSourceTemplate struct {
	Type  string          `yaml:"type"`
	Clock *RawClockConfig `yaml:"clock,omitempty"`
	Min   *int            `yaml:"min,omitempty"`
	Max   *int            `yaml:"max,omitempty"`
}

// RawValueTemplate - always full object, no self-reference
type RawValueTemplate struct {
	Source     *RawSourceConfig  `yaml:"source,omitempty"`
	Clock      *RawClockConfig   `yaml:"clock,omitempty"`
	Transforms []TransformConfig `yaml:"transforms,omitempty"`
	Reset      ResetConfig       `yaml:"reset,omitempty"`
}

// RawMetricTemplate - always full object, no self-reference
type RawMetricTemplate struct {
	Type       string             `yaml:"type"`
	Value      *RawValueReference `yaml:"value,omitempty"`
	Attributes map[string]string  `yaml:"attributes,omitempty"`
}
