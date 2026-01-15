package config

import (
	"fmt"
	"os"

	"go.yaml.in/yaml/v4"
)

// Load reads and parses a YAML configuration file.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	if err := validate(&cfg); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

// validate checks configuration consistency.
func validate(cfg *Config) error {
	// Validate export configuration
	if err := cfg.Export.Validate(); err != nil {
		return err
	}

	// Validate simulation configuration
	if err := cfg.Simulation.Validate(); err != nil {
		return err
	}

	// Validate metrics reference valid values
	for _, metric := range cfg.Metrics {
		if _, exists := cfg.Simulation.Values[metric.Value]; !exists {
			return fmt.Errorf("metric %q references unknown value %q", metric.Name, metric.Value)
		}
	}

	return nil
}
