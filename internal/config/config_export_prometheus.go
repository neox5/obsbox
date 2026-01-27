package config

import "fmt"

const (
	// Prometheus defaults
	DefaultPrometheusPort = 9090
	DefaultPrometheusPath = "/metrics"
)

// PrometheusExportConfig defines Prometheus pull endpoint settings.
type PrometheusExportConfig struct {
	Enabled bool   `yaml:"enabled"`
	Port    int    `yaml:"port"`
	Path    string `yaml:"path"`
}

// Validate applies defaults and validates Prometheus configuration.
func (c *PrometheusExportConfig) Validate() error {
	if !c.Enabled {
		return nil
	}

	// Apply defaults
	if c.Port == 0 {
		c.Port = DefaultPrometheusPort
	}
	if c.Path == "" {
		c.Path = DefaultPrometheusPath
	}

	// Validate port range
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("invalid prometheus port: %d", c.Port)
	}

	return nil
}
