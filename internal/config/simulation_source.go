package config

// SourceConfig defines a simv source.
type SourceConfig struct {
	Type  string `yaml:"type"`
	Clock string `yaml:"clock"`
	Min   int    `yaml:"min,omitempty"`
	Max   int    `yaml:"max,omitempty"`
}
