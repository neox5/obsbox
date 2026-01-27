package config

import "fmt"

// resolveClockTemplates resolves clock templates (no dependencies)
func (r *Resolver) resolveClockTemplates() error {
	for name, raw := range r.raw.Templates.Clocks {
		ctx := resolveContext{}.push("clock template", name)

		resolved := ClockConfig{
			Type:     raw.Type,
			Interval: raw.Interval,
		}

		// Validate
		if resolved.Type == "" {
			return ctx.error("type required")
		}
		if resolved.Interval == 0 {
			return ctx.error("interval required")
		}

		r.clocks[name] = resolved
	}
	return nil
}

// resolveSourceTemplates resolves source templates (may reference clock templates)
func (r *Resolver) resolveSourceTemplates() error {
	for name, raw := range r.raw.Templates.Sources {
		ctx := resolveContext{}.push("source template", name)

		resolved := SourceConfig{
			Type: raw.Type,
		}

		// Resolve clock reference if present
		if raw.Clock != nil {
			clock, err := r.resolveClock(raw.Clock, ctx)
			if err != nil {
				return err
			}
			resolved.Clock = clock
		}

		// Copy optional fields
		if raw.Min != nil {
			resolved.Min = *raw.Min
		}
		if raw.Max != nil {
			resolved.Max = *raw.Max
		}

		// Validate
		if resolved.Type == "" {
			return ctx.error("type required")
		}

		r.sources[name] = resolved
	}
	return nil
}

// resolveValueTemplates resolves value templates (may reference source + clock templates)
func (r *Resolver) resolveValueTemplates() error {
	for name, raw := range r.raw.Templates.Values {
		ctx := resolveContext{}.push("value template", name)

		resolved := ValueConfig{}

		// Resolve source reference if present
		if raw.Source != nil {
			source, err := r.resolveSource(raw.Source, ctx)
			if err != nil {
				return err
			}
			resolved.Source = source
		}

		// Resolve clock reference if present (optional override)
		if raw.Clock != nil {
			clock, err := r.resolveClock(raw.Clock, ctx)
			if err != nil {
				return err
			}
			resolved.Clock = clock
		}

		// Copy transforms and reset
		resolved.Transforms = raw.Transforms
		resolved.Reset = raw.Reset

		// Validate
		if err := r.validateValue(resolved, ctx); err != nil {
			return err
		}

		r.values[name] = resolved
	}
	return nil
}

// resolveMetricTemplates resolves metric templates (may reference value templates)
func (r *Resolver) resolveMetricTemplates() error {
	for name, raw := range r.raw.Templates.Metrics {
		ctx := resolveContext{}.push("metric template", name)

		resolved := MetricConfig{
			Type: MetricType(raw.Type),
		}

		// Resolve value reference if present
		if raw.Value != nil {
			value, err := r.resolveValue(raw.Value, ctx)
			if err != nil {
				return err
			}
			resolved.Value = value
		}

		// Copy attributes (can be nil)
		if raw.Attributes != nil {
			resolved.Attributes = make(map[string]string, len(raw.Attributes))
			for k, v := range raw.Attributes {
				resolved.Attributes[k] = v
			}
		}

		// Validate
		if resolved.Type == "" {
			return ctx.error("type required")
		}

		r.metrics[name] = resolved
	}
	return nil
}

// resolveClock resolves a clock reference or inline definition
func (r *Resolver) resolveClock(raw *RawClockConfig, ctx resolveContext) (ClockConfig, error) {
	// Template reference
	if raw.Template != "" {
		template, exists := r.clocks[raw.Template]
		if !exists {
			return ClockConfig{}, ctx.error(fmt.Sprintf("clock template %q not found", raw.Template))
		}
		return template, nil
	}

	// Inline definition
	if raw.Inline != nil {
		resolved := ClockConfig{
			Type:     raw.Inline.Type,
			Interval: raw.Inline.Interval,
		}

		// Validate
		if resolved.Type == "" {
			return ClockConfig{}, ctx.error("clock type required")
		}
		if resolved.Interval == 0 {
			return ClockConfig{}, ctx.error("clock interval required")
		}

		return resolved, nil
	}

	return ClockConfig{}, ctx.error("clock must be template reference or inline definition")
}

// resolveSource resolves a source reference or inline definition
func (r *Resolver) resolveSource(raw *RawSourceConfig, ctx resolveContext) (SourceConfig, error) {
	// Template reference
	if raw.Template != "" {
		template, exists := r.sources[raw.Template]
		if !exists {
			return SourceConfig{}, ctx.error(fmt.Sprintf("source template %q not found", raw.Template))
		}
		return template, nil
	}

	// Inline definition
	if raw.Inline != nil {
		resolved := SourceConfig{
			Type: raw.Inline.Type,
		}

		// Resolve clock if present
		if raw.Inline.Clock != nil {
			clock, err := r.resolveClock(raw.Inline.Clock, ctx)
			if err != nil {
				return SourceConfig{}, err
			}
			resolved.Clock = clock
		}

		// Copy optional fields
		if raw.Inline.Min != nil {
			resolved.Min = *raw.Inline.Min
		}
		if raw.Inline.Max != nil {
			resolved.Max = *raw.Inline.Max
		}

		// Validate
		if resolved.Type == "" {
			return SourceConfig{}, ctx.error("source type required")
		}

		return resolved, nil
	}

	return SourceConfig{}, ctx.error("source must be template reference or inline definition")
}

// resolveValue resolves a value reference or inline definition with overrides
func (r *Resolver) resolveValue(raw *RawValueReference, ctx resolveContext) (ValueConfig, error) {
	var result ValueConfig

	// Step 1: Apply template if string form
	if raw.Template != "" {
		template, exists := r.values[raw.Template]
		if !exists {
			return ValueConfig{}, ctx.error(fmt.Sprintf("value template %q not found", raw.Template))
		}
		result = template
	}

	// Step 2: Apply inline/overrides if object form
	if raw.Inline != nil {
		if err := r.applyValueOverrides(&result, raw.Inline, ctx); err != nil {
			return ValueConfig{}, err
		}
	}

	// Step 3: Validate final result
	if err := r.validateValue(result, ctx); err != nil {
		return ValueConfig{}, err
	}

	return result, nil
}

// applyValueOverrides applies field overrides to value config
func (r *Resolver) applyValueOverrides(result *ValueConfig, overrides *RawValueConfig, ctx resolveContext) error {
	// Override source (complete replacement)
	if overrides.Source != nil {
		source, err := r.resolveSource(overrides.Source, ctx)
		if err != nil {
			return err
		}
		result.Source = source
	}

	// Override clock (complete replacement, optional)
	if overrides.Clock != nil {
		clock, err := r.resolveClock(overrides.Clock, ctx)
		if err != nil {
			return err
		}
		result.Clock = clock
	}

	// Override transforms (complete replacement)
	if overrides.Transforms != nil {
		result.Transforms = overrides.Transforms
	}

	// Override reset (complete replacement)
	if overrides.Reset.Type != "" {
		result.Reset = overrides.Reset
	}

	return nil
}

// validateValue validates a resolved value config
func (r *Resolver) validateValue(value ValueConfig, ctx resolveContext) error {
	// Source required
	if value.Source.Type == "" {
		return ctx.error("source required")
	}

	// Clock required - either on value (override) or source (inherited)
	hasValueClock := value.Clock.Type != ""
	hasSourceClock := value.Source.Clock.Type != ""

	if !hasValueClock && !hasSourceClock {
		return ctx.error("clock required (either on value or source)")
	}

	return nil
}
