package config

import (
	"fmt"
)

// resolveTemplateClocks resolves clock templates (no dependencies)
func (r *Resolver) resolveTemplateClocks() error {
	for name, raw := range r.raw.Templates.Clocks {
		if err := r.registerName(name, "template clock"); err != nil {
			return err
		}

		ctx := resolveContext{}.push("clock template", name)

		resolved := ClockConfig{
			Type:     getStringValue(raw.Type),
			Interval: raw.Interval,
		}

		// Validate
		if resolved.Type == "" {
			return ctx.error("type required")
		}
		if resolved.Interval == 0 {
			return ctx.error("interval required")
		}

		r.templateClocks[name] = resolved
	}
	return nil
}

// resolveTemplateSources resolves source templates (may reference clock templates)
func (r *Resolver) resolveTemplateSources() error {
	for name, raw := range r.raw.Templates.Sources {
		if err := r.registerName(name, "template source"); err != nil {
			return err
		}

		ctx := resolveContext{}.push("source template", name)

		resolved := SourceConfig{
			Type: getStringValue(raw.Type),
		}

		// Resolve clock (inline only for templates)
		if raw.Clock != nil {
			clock, clockRef, err := r.resolveClockReference(raw.Clock, ctx)
			if err != nil {
				return err
			}
			resolved.Clock = clock
			resolved.ClockRef = clockRef
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

		r.templateSources[name] = resolved
	}
	return nil
}

// resolveTemplateValues resolves value templates (may reference source templates)
func (r *Resolver) resolveTemplateValues() error {
	for name, raw := range r.raw.Templates.Values {
		if err := r.registerName(name, "template value"); err != nil {
			return err
		}

		ctx := resolveContext{}.push("value template", name)

		resolved := ValueConfig{}

		// Resolve source (inline only for templates)
		if raw.Source != nil {
			source, sourceRef, err := r.resolveSourceFromReference(raw.Source, ctx)
			if err != nil {
				return err
			}
			resolved.Source = source
			resolved.SourceRef = sourceRef
		}

		// Copy transforms and reset
		resolved.Transforms = raw.Transforms
		resolved.Reset = raw.Reset

		// Validate
		if err := r.validateValue(resolved, ctx); err != nil {
			return err
		}

		r.templateValues[name] = resolved
	}
	return nil
}

// resolveTemplateMetrics resolves metric templates (may reference value templates)
func (r *Resolver) resolveTemplateMetrics() error {
	for name, raw := range r.raw.Templates.Metrics {
		if err := r.registerName(name, "template metric"); err != nil {
			return err
		}

		ctx := resolveContext{}.push("metric template", name)

		resolved := MetricConfig{
			Type: MetricType(raw.Type),
		}

		// Resolve value reference if present
		if raw.Value != nil {
			value, err := r.resolveValueFromReference(raw.Value, ctx)
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

		r.templateMetrics[name] = resolved
	}
	return nil
}

// resolveInstanceClocks resolves clock instances
func (r *Resolver) resolveInstanceClocks() error {
	for name, raw := range r.raw.Instances.Clocks {
		if err := r.registerName(name, "instance clock"); err != nil {
			return err
		}

		ctx := resolveContext{}.push("clock instance", name)

		resolved := ClockConfig{
			Type:     getStringValue(raw.Type),
			Interval: raw.Interval,
		}

		// Validate
		if resolved.Type == "" {
			return ctx.error("type required")
		}
		if resolved.Interval == 0 {
			return ctx.error("interval required")
		}

		r.instanceClocks[name] = resolved
	}
	return nil
}

// resolveInstanceSources resolves source instances (may reference template/instance clocks)
func (r *Resolver) resolveInstanceSources() error {
	for name, raw := range r.raw.Instances.Sources {
		if err := r.registerName(name, "instance source"); err != nil {
			return err
		}

		ctx := resolveContext{}.push("source instance", name)

		resolved := SourceConfig{
			Type: getStringValue(raw.Type),
		}

		// Resolve clock reference if present
		if raw.Clock != nil {
			clock, clockRef, err := r.resolveClockReference(raw.Clock, ctx)
			if err != nil {
				return err
			}
			resolved.Clock = clock
			resolved.ClockRef = clockRef
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

		r.instanceSources[name] = resolved
	}
	return nil
}

// resolveInstanceValues resolves value instances (may reference template/instance sources)
func (r *Resolver) resolveInstanceValues() error {
	for name, raw := range r.raw.Instances.Values {
		if err := r.registerName(name, "instance value"); err != nil {
			return err
		}

		ctx := resolveContext{}.push("value instance", name)

		resolved := ValueConfig{}

		// Resolve source reference if present
		if raw.Source != nil {
			source, sourceRef, err := r.resolveSourceFromReference(raw.Source, ctx)
			if err != nil {
				return err
			}
			resolved.Source = source
			resolved.SourceRef = sourceRef
		}

		// Copy transforms and reset
		resolved.Transforms = raw.Transforms
		resolved.Reset = raw.Reset

		// Validate
		if err := r.validateValue(resolved, ctx); err != nil {
			return err
		}

		r.instanceValues[name] = resolved
	}
	return nil
}

// resolveSourceFromReference resolves a source from RawSourceReference (supports instance/template/inline)
func (r *Resolver) resolveSourceFromReference(raw *RawSourceReference, ctx resolveContext) (SourceConfig, *string, error) {
	// Instance reference
	if raw.Instance != "" {
		instance, exists := r.instanceSources[raw.Instance]
		if !exists {
			return SourceConfig{}, nil, ctx.error(fmt.Sprintf("source instance %q not found", raw.Instance))
		}
		// No overrides allowed for instances
		if raw.Template != "" || raw.Type != nil || raw.Clock != nil || raw.Min != nil || raw.Max != nil {
			return SourceConfig{}, nil, ctx.error("cannot override instance source")
		}
		return instance, &raw.Instance, nil
	}

	// Template reference (with optional overrides)
	if raw.Template != "" {
		template, exists := r.templateSources[raw.Template]
		if !exists {
			return SourceConfig{}, nil, ctx.error(fmt.Sprintf("source template %q not found", raw.Template))
		}

		// Apply overrides
		result := template
		if raw.Type != nil {
			result.Type = *raw.Type
		}
		if raw.Clock != nil {
			clock, clockRef, err := r.resolveClockReference(raw.Clock, ctx)
			if err != nil {
				return SourceConfig{}, nil, err
			}
			result.Clock = clock
			result.ClockRef = clockRef
		}
		if raw.Min != nil {
			result.Min = *raw.Min
		}
		if raw.Max != nil {
			result.Max = *raw.Max
		}
		return result, nil, nil
	}

	// Inline definition
	if raw.Type != nil {
		result := SourceConfig{}
		result.Type = *raw.Type

		// Resolve clock if present
		if raw.Clock != nil {
			clock, clockRef, err := r.resolveClockReference(raw.Clock, ctx)
			if err != nil {
				return SourceConfig{}, nil, err
			}
			result.Clock = clock
			result.ClockRef = clockRef
		}

		// Copy optional fields
		if raw.Min != nil {
			result.Min = *raw.Min
		}
		if raw.Max != nil {
			result.Max = *raw.Max
		}

		// Validate
		if result.Type == "" {
			return SourceConfig{}, nil, ctx.error("source type required")
		}

		return result, nil, nil
	}

	return SourceConfig{}, nil, ctx.error("source must reference instance, template, or provide inline definition")
}

// resolveValueFromReference resolves a value from RawValueReference (supports instance/template/inline)
func (r *Resolver) resolveValueFromReference(raw *RawValueReference, ctx resolveContext) (ValueConfig, error) {
	// Instance reference
	if raw.Instance != "" {
		instance, exists := r.instanceValues[raw.Instance]
		if !exists {
			return ValueConfig{}, ctx.error(fmt.Sprintf("value instance %q not found", raw.Instance))
		}
		// No overrides allowed for instances
		if raw.Template != "" || raw.Source != nil || len(raw.Transforms) > 0 || raw.Reset.Type != "" {
			return ValueConfig{}, ctx.error("cannot override instance value")
		}
		return instance, nil
	}

	// Template reference (with optional overrides)
	if raw.Template != "" {
		template, exists := r.templateValues[raw.Template]
		if !exists {
			return ValueConfig{}, ctx.error(fmt.Sprintf("value template %q not found", raw.Template))
		}

		// Apply overrides
		result := template
		if raw.Source != nil {
			source, sourceRef, err := r.resolveSourceFromReference(raw.Source, ctx)
			if err != nil {
				return ValueConfig{}, err
			}
			result.Source = source
			result.SourceRef = sourceRef
		}
		if len(raw.Transforms) > 0 {
			result.Transforms = raw.Transforms
		}
		if raw.Reset.Type != "" {
			result.Reset = raw.Reset
		}
		return result, nil
	}

	// Inline definition
	if raw.Source != nil {
		result := ValueConfig{}

		source, sourceRef, err := r.resolveSourceFromReference(raw.Source, ctx)
		if err != nil {
			return ValueConfig{}, err
		}
		result.Source = source
		result.SourceRef = sourceRef

		result.Transforms = raw.Transforms
		result.Reset = raw.Reset

		return result, nil
	}

	return ValueConfig{}, ctx.error("value must reference instance, template, or provide inline definition")
}

// resolveClockReference resolves a clock reference (supports instance/template/inline)
func (r *Resolver) resolveClockReference(raw *RawClockReference, ctx resolveContext) (ClockConfig, *string, error) {
	// Instance reference
	if raw.Instance != "" {
		instance, exists := r.instanceClocks[raw.Instance]
		if !exists {
			return ClockConfig{}, nil, ctx.error(fmt.Sprintf("clock instance %q not found", raw.Instance))
		}
		// No overrides allowed for instances
		if raw.Template != "" || raw.Type != nil || raw.Interval != 0 {
			return ClockConfig{}, nil, ctx.error("cannot override instance clock")
		}
		return instance, &raw.Instance, nil
	}

	// Template reference (with optional overrides)
	if raw.Template != "" {
		template, exists := r.templateClocks[raw.Template]
		if !exists {
			return ClockConfig{}, nil, ctx.error(fmt.Sprintf("clock template %q not found", raw.Template))
		}

		// Apply overrides
		result := template
		if raw.Type != nil {
			result.Type = *raw.Type
		}
		if raw.Interval != 0 {
			result.Interval = raw.Interval
		}
		return result, nil, nil
	}

	// Inline definition
	if raw.Type != nil {
		resolved := ClockConfig{
			Type:     *raw.Type,
			Interval: raw.Interval,
		}

		// Validate
		if resolved.Type == "" {
			return ClockConfig{}, nil, ctx.error("clock type required")
		}
		if resolved.Interval == 0 {
			return ClockConfig{}, nil, ctx.error("clock interval required")
		}

		return resolved, nil, nil
	}

	return ClockConfig{}, nil, ctx.error("clock must reference instance, template, or provide inline definition")
}

// validateValue validates a resolved value config
func (r *Resolver) validateValue(value ValueConfig, ctx resolveContext) error {
	// Source required
	if value.Source.Type == "" {
		return ctx.error("source required")
	}

	// Clock required in source
	if value.Source.Clock.Type == "" {
		return ctx.error("clock required in source")
	}

	return nil
}

// getStringValue safely dereferences a string pointer
func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
