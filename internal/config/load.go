package config

import (
	"fmt"

	"github.com/neox5/obsbox/internal/configparse"
	"github.com/neox5/obsbox/internal/configresolve"
)

// Load reads and resolves a YAML configuration file
func Load(path string) (*Config, error) {
	raw, err := configparse.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	cfg, err := configresolve.Resolve(raw)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve config: %w", err)
	}

	return cfg, nil
}
