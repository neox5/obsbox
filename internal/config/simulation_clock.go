package config

import "time"

// ClockConfig defines a clock.
type ClockConfig struct {
	Type     string        `yaml:"type"`
	Interval time.Duration `yaml:"interval"`
}
