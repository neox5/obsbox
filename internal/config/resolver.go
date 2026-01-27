package config

import (
	"fmt"
	"strings"
)

// resolveContext tracks resolution path for error messages
type resolveContext []string

func (ctx resolveContext) push(component, name string) resolveContext {
	return append(ctx, fmt.Sprintf("%s %q", component, name))
}

func (ctx resolveContext) error(msg string) error {
	if len(ctx) == 0 {
		return fmt.Errorf(msg)
	}

	var b strings.Builder
	b.WriteString(msg)
	// Print stack top-down (metric → templates → error)
	for i := len(ctx) - 1; i >= 0; i-- {
		b.WriteString("\n  in ")
		b.WriteString(ctx[i])
	}
	return fmt.Errorf(b.String())
}

// Resolver handles template resolution
type Resolver struct {
	raw *RawConfig

	// Resolved templates (temporary, discarded after final config built)
	clocks  map[string]ClockConfig
	sources map[string]SourceConfig
	values  map[string]ValueConfig
	metrics map[string]MetricConfig
}

// NewResolver creates a new resolver
func NewResolver(raw *RawConfig) *Resolver {
	return &Resolver{
		raw:     raw,
		clocks:  make(map[string]ClockConfig),
		sources: make(map[string]SourceConfig),
		values:  make(map[string]ValueConfig),
		metrics: make(map[string]MetricConfig),
	}
}

// Resolve performs hierarchical template resolution and builds final config
func (r *Resolver) Resolve() (*Config, error) {
	// Phase 1: Resolve templates hierarchically
	if err := r.resolveClockTemplates(); err != nil {
		return nil, err
	}
	if err := r.resolveSourceTemplates(); err != nil {
		return nil, err
	}
	if err := r.resolveValueTemplates(); err != nil {
		return nil, err
	}
	if err := r.resolveMetricTemplates(); err != nil {
		return nil, err
	}

	// Phase 2: Resolve metrics
	resolvedMetrics, err := r.resolveMetrics()
	if err != nil {
		return nil, err
	}

	// Phase 3: Build final config (templates discarded)
	return &Config{
		Metrics:  resolvedMetrics,
		Export:   r.raw.Export,
		Settings: r.raw.Settings,
	}, nil
}
