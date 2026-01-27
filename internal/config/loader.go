package config

import (
	"fmt"
	"os"

	"go.yaml.in/yaml/v4"
)

// Load reads and parses a YAML configuration file using two-step process
func Load(path string) (*Config, error) {
	// Step 1: Load raw config (syntactic validation)
	raw, err := loadRaw(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Step 2: Resolve templates (semantic validation)
	resolver := NewResolver(raw)
	config, err := resolver.Resolve()
	if err != nil {
		return nil, fmt.Errorf("failed to resolve config: %w", err)
	}

	return config, nil
}

// loadRaw reads YAML file and performs syntactic validation
func loadRaw(path string) (*RawConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var raw RawConfig
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Syntactic validation
	if err := validateRawSyntax(&raw); err != nil {
		return nil, err
	}

	// Validate export configuration
	if err := raw.Export.Validate(); err != nil {
		return nil, err
	}

	// Validate settings configuration
	if err := raw.Settings.Validate(); err != nil {
		return nil, err
	}

	return &raw, nil
}
