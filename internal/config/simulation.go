package config

import "fmt"

// SimulationConfig defines the simulation domain configuration.
type SimulationConfig struct {
	Seed    *uint64                 `yaml:"seed,omitempty"`
	Clocks  map[string]ClockConfig  `yaml:"clocks"`
	Sources map[string]SourceConfig `yaml:"sources"`
	Values  map[string]ValueConfig  `yaml:"values"`
}

// HasSeed returns true if a master seed was explicitly configured.
func (s *SimulationConfig) HasSeed() bool {
	return s.Seed != nil
}

// GetSeed returns the configured master seed value.
// Only call after checking HasSeed().
func (s *SimulationConfig) GetSeed() uint64 {
	if s.Seed == nil {
		return 0
	}
	return *s.Seed
}

// Validate checks simulation configuration consistency.
func (s *SimulationConfig) Validate() error {
	// Validate component existence
	if len(s.Clocks) == 0 {
		return fmt.Errorf("at least one clock must be defined")
	}
	if len(s.Sources) == 0 {
		return fmt.Errorf("at least one source must be defined")
	}
	if len(s.Values) == 0 {
		return fmt.Errorf("at least one value must be defined")
	}

	// Validate source clock references
	for srcName, src := range s.Sources {
		if _, exists := s.Clocks[src.Clock]; !exists {
			return fmt.Errorf("source %q references unknown clock %q", srcName, src.Clock)
		}
	}

	// Validate value source references
	for valName, val := range s.Values {
		if val.Source == "" {
			return fmt.Errorf("value %q must specify source", valName)
		}
		if _, exists := s.Sources[val.Source]; !exists {
			return fmt.Errorf("value %q references unknown source %q", valName, val.Source)
		}
	}

	return nil
}
