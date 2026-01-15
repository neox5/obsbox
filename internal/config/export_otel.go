package config

import (
	"fmt"
	"time"

	"go.yaml.in/yaml/v4"
)

const (
	// OTEL defaults
	DefaultOTELReadInterval = 1 * time.Second
	DefaultOTELPushInterval = 1 * time.Second
	DefaultServiceName      = "obsbox"
	DefaultServiceVersion   = "dev"
)

// OTELExportConfig defines OTEL push settings.
type OTELExportConfig struct {
	Enabled  bool              `yaml:"enabled"`
	Endpoint string            `yaml:"endpoint"`
	Interval IntervalConfig    `yaml:"interval"`
	Resource map[string]string `yaml:"resource,omitempty"`
	Headers  map[string]string `yaml:"headers,omitempty"`
}

// IntervalConfig defines read and push intervals for OTEL.
type IntervalConfig struct {
	Read time.Duration
	Push time.Duration
}

// UnmarshalYAML handles both simple (10s) and detailed (read/push) forms.
func (i *IntervalConfig) UnmarshalYAML(value *yaml.Node) error {
	// Try simple duration form first
	var simple time.Duration
	if err := value.Decode(&simple); err == nil {
		i.Read = simple
		i.Push = simple
		return nil
	}

	// Fall back to detailed form
	type intervalConfig struct {
		Read time.Duration `yaml:"read"`
		Push time.Duration `yaml:"push"`
	}
	var detailed intervalConfig
	if err := value.Decode(&detailed); err != nil {
		return err
	}
	i.Read = detailed.Read
	i.Push = detailed.Push
	return nil
}

// Validate applies defaults and validates OTEL configuration.
func (c *OTELExportConfig) Validate() error {
	if !c.Enabled {
		return nil
	}

	// Validate required fields
	if c.Endpoint == "" {
		return fmt.Errorf("otel endpoint cannot be empty when enabled")
	}

	// Apply interval defaults
	if c.Interval.Read == 0 {
		c.Interval.Read = DefaultOTELReadInterval
	}
	if c.Interval.Push == 0 {
		c.Interval.Push = DefaultOTELPushInterval
	}

	// Apply resource defaults
	if c.Resource == nil {
		c.Resource = make(map[string]string)
	}
	if _, exists := c.Resource["service.name"]; !exists {
		c.Resource["service.name"] = DefaultServiceName
	}
	if _, exists := c.Resource["service.version"]; !exists {
		c.Resource["service.version"] = DefaultServiceVersion
	}

	return nil
}
