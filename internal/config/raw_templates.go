package config

// RawTemplates holds all template definitions
type RawTemplates struct {
	Clocks  map[string]RawClockReference   `yaml:"clocks,omitempty"`
	Sources map[string]RawSourceReference  `yaml:"sources,omitempty"`
	Values  map[string]RawValueReference   `yaml:"values,omitempty"`
	Metrics map[string]RawMetricDefinition `yaml:"metrics,omitempty"`
}

// RawMetricDefinition - always full object, no self-reference
type RawMetricDefinition struct {
	Type       string             `yaml:"type"`
	Value      *RawValueReference `yaml:"value,omitempty"`
	Attributes map[string]string  `yaml:"attributes,omitempty"`
}
