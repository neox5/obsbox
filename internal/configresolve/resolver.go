package configresolve

import (
	"github.com/neox5/obsbox/internal/config"
	"github.com/neox5/obsbox/internal/configparse"
)

// Resolver handles template and instance resolution
type Resolver struct {
	raw *configparse.RawConfig

	// Namespace tracking (all entity names)
	registeredNames map[string]string // name -> entity type

	// Resolved templates (temporary, discarded after final config built)
	templateClocks  map[string]config.ClockConfig
	templateSources map[string]config.SourceConfig
	templateValues  map[string]config.ValueConfig
	templateMetrics map[string]config.MetricConfig

	// Resolved instances (kept in final config)
	instanceClocks  map[string]config.ClockConfig
	instanceSources map[string]config.SourceConfig
	instanceValues  map[string]config.ValueConfig
}

// newResolver creates a new resolver
func newResolver(raw *configparse.RawConfig) *Resolver {
	return &Resolver{
		raw:             raw,
		registeredNames: make(map[string]string),
		templateClocks:  make(map[string]config.ClockConfig),
		templateSources: make(map[string]config.SourceConfig),
		templateValues:  make(map[string]config.ValueConfig),
		templateMetrics: make(map[string]config.MetricConfig),
		instanceClocks:  make(map[string]config.ClockConfig),
		instanceSources: make(map[string]config.SourceConfig),
		instanceValues:  make(map[string]config.ValueConfig),
	}
}

// Resolve performs hierarchical template and instance resolution and builds final config
func Resolve(raw *configparse.RawConfig) (*config.Config, error) {
	r := newResolver(raw)

	// Phase 1: Resolve templates hierarchically
	if err := r.resolveTemplateClocks(); err != nil {
		return nil, err
	}
	if err := r.resolveTemplateSources(); err != nil {
		return nil, err
	}
	if err := r.resolveTemplateValues(); err != nil {
		return nil, err
	}
	if err := r.resolveTemplateMetrics(); err != nil {
		return nil, err
	}

	// Phase 2: Resolve instances hierarchically
	if err := r.resolveInstanceClocks(); err != nil {
		return nil, err
	}
	if err := r.resolveInstanceSources(); err != nil {
		return nil, err
	}
	if err := r.resolveInstanceValues(); err != nil {
		return nil, err
	}

	// Phase 3: Resolve metrics
	resolvedMetrics, err := r.resolveMetrics()
	if err != nil {
		return nil, err
	}

	// Phase 4: Resolve export config
	resolvedExport, err := resolveExport(&r.raw.Export)
	if err != nil {
		return nil, err
	}

	// Phase 5: Resolve settings config
	resolvedSettings, err := resolveSettings(&r.raw.Settings)
	if err != nil {
		return nil, err
	}

	// Phase 6: Build final config (templates discarded, instances kept)
	return &config.Config{
		Instances: config.InstanceRegistry{
			Clocks:  r.instanceClocks,
			Sources: r.instanceSources,
			Values:  r.instanceValues,
		},
		Metrics:  resolvedMetrics,
		Export:   resolvedExport,
		Settings: resolvedSettings,
	}, nil
}
