package config

import (
	"fmt"
	"log/slog"
	"strings"
)

// Resolver handles template and instance resolution
type Resolver struct {
	raw *RawConfig

	// Namespace tracking (all entity names)
	registeredNames map[string]string // name -> entity type

	// Resolved templates (temporary, discarded after final config built)
	templateClocks  map[string]ClockConfig
	templateSources map[string]SourceConfig
	templateValues  map[string]ValueConfig
	templateMetrics map[string]MetricConfig

	// Resolved instances (kept in final config)
	instanceClocks  map[string]ClockConfig
	instanceSources map[string]SourceConfig
	instanceValues  map[string]ValueConfig
}

// newResolver creates a new resolver
func newResolver(raw *RawConfig) *Resolver {
	return &Resolver{
		raw:             raw,
		registeredNames: make(map[string]string),
		templateClocks:  make(map[string]ClockConfig),
		templateSources: make(map[string]SourceConfig),
		templateValues:  make(map[string]ValueConfig),
		templateMetrics: make(map[string]MetricConfig),
		instanceClocks:  make(map[string]ClockConfig),
		instanceSources: make(map[string]SourceConfig),
		instanceValues:  make(map[string]ValueConfig),
	}
}

// Resolve performs hierarchical template and instance resolution and builds final config
func Resolve(raw *RawConfig) (*Config, error) {
	// Expansion must happen before calling Resolve
	// This is enforced by Load() pipeline

	slog.Debug("--- Template and Instance Resolution ---")
	resolver := newResolver(raw)

	// Clocks (no dependencies)
	slog.Debug("resolving clocks")
	if err := resolver.resolveTemplateClocks(); err != nil {
		return nil, err
	}
	if err := resolver.resolveInstanceClocks(); err != nil {
		return nil, err
	}

	// Sources (depend on clocks)
	slog.Debug("resolving sources")
	if err := resolver.resolveTemplateSources(); err != nil {
		return nil, err
	}
	if err := resolver.resolveInstanceSources(); err != nil {
		return nil, err
	}

	// Values (depend on sources)
	slog.Debug("resolving values")
	if err := resolver.resolveTemplateValues(); err != nil {
		return nil, err
	}
	if err := resolver.resolveInstanceValues(); err != nil {
		return nil, err
	}

	// Metrics (depend on values)
	slog.Debug("resolving metrics")
	if err := resolver.resolveTemplateMetrics(); err != nil {
		return nil, err
	}

	// Phase 3: Metric resolution
	slog.Debug("--- Metric Resolution ---")
	metrics, err := resolver.resolveMetrics()
	if err != nil {
		return nil, err
	}

	// Phase 4: Export resolution
	export, err := resolveExport(&raw.Export)
	if err != nil {
		return nil, err
	}

	// Phase 5: Settings resolution
	settings, err := resolveSettings(&raw.Settings)
	if err != nil {
		return nil, err
	}

	// Phase 6: Assemble final config
	return buildConfig(resolver, metrics, export, settings), nil
}

// buildConfig assembles the final configuration
func buildConfig(
	resolver *Resolver,
	metrics []MetricConfig,
	export ExportConfig,
	settings SettingsConfig,
) *Config {
	return &Config{
		Instances: InstanceRegistry{
			Clocks:  resolver.instanceClocks,
			Sources: resolver.instanceSources,
			Values:  resolver.instanceValues,
		},
		Metrics:  metrics,
		Export:   export,
		Settings: settings,
	}
}

// registerName validates namespace uniqueness and registers the name
func (r *Resolver) registerName(name string, entityType string) error {
	if existingType, exists := r.registeredNames[name]; exists {
		return fmt.Errorf("name %q already used by %s, cannot reuse for %s",
			name, existingType, entityType)
	}
	r.registeredNames[name] = entityType
	return nil
}

// getStringValue safely dereferences a string pointer
func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// resolveContext tracks resolution path for error messages
type resolveContext []string

func (ctx resolveContext) push(component, name string) resolveContext {
	return append(ctx, fmt.Sprintf("%s %q", component, name))
}

func (ctx resolveContext) error(msg string) error {
	if len(ctx) == 0 {
		return fmt.Errorf("%s", msg)
	}

	var b strings.Builder
	b.WriteString(msg)
	// Print stack top-down (metric → templates → error)
	for i := len(ctx) - 1; i >= 0; i-- {
		b.WriteString("\n  in ")
		b.WriteString(ctx[i])
	}
	return fmt.Errorf("%s", b.String())
}
