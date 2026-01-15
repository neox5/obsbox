package config

import "fmt"

// SimulationConfig defines the simulation domain configuration.
type SimulationConfig struct {
	Clocks  map[string]ClockConfig  `yaml:"clocks"`
	Sources map[string]SourceConfig `yaml:"sources"`
	Values  map[string]ValueConfig  `yaml:"values"`
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

	// Validate value references (source or clone, not both)
	for valName, val := range s.Values {
		if val.Source == "" && val.Clone == "" {
			return fmt.Errorf("value %q must specify either source or clone", valName)
		}
		if val.Source != "" && val.Clone != "" {
			return fmt.Errorf("value %q cannot specify both source and clone", valName)
		}
		if val.Source != "" {
			if _, exists := s.Sources[val.Source]; !exists {
				return fmt.Errorf("value %q references unknown source %q", valName, val.Source)
			}
		}
		if val.Clone != "" {
			if _, exists := s.Values[val.Clone]; !exists {
				return fmt.Errorf("value %q references unknown clone %q", valName, val.Clone)
			}
		}
	}

	return nil
}
