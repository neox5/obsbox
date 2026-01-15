package config

// MetricConfig defines a Prometheus metric.
type MetricConfig struct {
	Name  string `yaml:"name"`
	Type  string `yaml:"type"`
	Help  string `yaml:"help"`
	Value string `yaml:"value"`
}
