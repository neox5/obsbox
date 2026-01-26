package config

// SeedConfig defines optional seed configuration for reproducible simulations.
type SeedConfig struct {
	MasterSeed *uint64 `yaml:"seed,omitempty"`
}

// IsConfigured returns true if a master seed was explicitly provided.
func (s *SeedConfig) IsConfigured() bool {
	return s.MasterSeed != nil
}

// Value returns the configured master seed value.
// Only call after checking IsConfigured().
func (s *SeedConfig) Value() uint64 {
	if s.MasterSeed == nil {
		return 0
	}
	return *s.MasterSeed
}
