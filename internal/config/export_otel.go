package config

import (
	"fmt"
	"time"
)

const (
	// OTEL defaults
	DefaultOTELReadInterval = 1 * time.Second
	DefaultOTELPushInterval = 1 * time.Second
	DefaultOTELTransport    = "grpc"
	DefaultOTELHost         = "localhost"
	DefaultOTELPortGRPC     = 4317
	DefaultOTELPortHTTP     = 4318
	DefaultServiceName      = "obsbox"
	DefaultServiceVersion   = "dev"
)

// OTELExportConfig defines OTEL push settings.
type OTELExportConfig struct {
	Enabled   bool
	Transport string
	Host      string
	Port      int
	Interval  IntervalConfig
	Resource  map[string]string
	Headers   map[string]string
}

// IntervalConfig defines read and push intervals for OTEL.
type IntervalConfig struct {
	Read time.Duration
	Push time.Duration
}

// Validate applies defaults and validates OTEL configuration.
func (c *OTELExportConfig) Validate() error {
	if !c.Enabled {
		return nil
	}

	// Apply transport default
	if c.Transport == "" {
		c.Transport = DefaultOTELTransport
	}

	// Validate transport
	if c.Transport != "grpc" && c.Transport != "http" {
		return fmt.Errorf("invalid transport: %s (must be grpc or http)", c.Transport)
	}

	// Apply host default
	if c.Host == "" {
		c.Host = DefaultOTELHost
	}

	// Apply port default based on transport
	if c.Port == 0 {
		if c.Transport == "grpc" {
			c.Port = DefaultOTELPortGRPC
		} else {
			c.Port = DefaultOTELPortHTTP
		}
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

// GetEndpoint returns the full endpoint address.
func (c *OTELExportConfig) GetEndpoint() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
