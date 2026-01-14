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
	// Validate server
	if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", cfg.Server.Port)
	}
	if cfg.Server.Path == "" {
		return fmt.Errorf("server path cannot be empty")
	}

	// Validate simulation config exists
	if len(cfg.Simulation.Clocks) == 0 {
		return fmt.Errorf("at least one clock must be defined")
	}
	if len(cfg.Simulation.Sources) == 0 {
		return fmt.Errorf("at least one source must be defined")
	}
	if len(cfg.Simulation.Values) == 0 {
		return fmt.Errorf("at least one value must be defined")
	}

	// Validate source clock references
	for srcName, src := range cfg.Simulation.Sources {
		if _, exists := cfg.Simulation.Clocks[src.Clock]; !exists {
			return fmt.Errorf("source %q references unknown clock %q", srcName, src.Clock)
		}
	}

	// Validate value references (source or clone, not both)
	for valName, val := range cfg.Simulation.Values {
		if val.Source == "" && val.Clone == "" {
			return fmt.Errorf("value %q must specify either source or clone", valName)
		}
		if val.Source != "" && val.Clone != "" {
			return fmt.Errorf("value %q cannot specify both source and clone", valName)
		}
		if val.Source != "" {
			if _, exists := cfg.Simulation.Sources[val.Source]; !exists {
				return fmt.Errorf("value %q references unknown source %q", valName, val.Source)
			}
		}
		if val.Clone != "" {
			if _, exists := cfg.Simulation.Values[val.Clone]; !exists {
				return fmt.Errorf("value %q references unknown clone %q", valName, val.Clone)
			}
		}
	}

	// Validate metrics reference valid values
	for _, metric := range cfg.Metrics {
		if _, exists := cfg.Simulation.Values[metric.Value]; !exists {
			return fmt.Errorf("metric %q references unknown value %q", metric.Name, metric.Value)
		}
	}

	return nil
}
