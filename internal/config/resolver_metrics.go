package config

import "fmt"

// resolveMetrics resolves final metrics from raw config
func (r *Resolver) resolveMetrics() ([]MetricConfig, error) {
	var metrics []MetricConfig

	for _, raw := range r.raw.Metrics {
		promName := raw.Name.GetPrometheusName()
		ctx := resolveContext{}.push("metric", promName)

		metric, err := r.resolveMetric(&raw, ctx)
		if err != nil {
			return nil, err
		}

		metrics = append(metrics, metric)
	}

	return metrics, nil
}

// resolveMetric resolves a single metric with template + overrides
func (r *Resolver) resolveMetric(raw *RawMetricConfig, ctx resolveContext) (MetricConfig, error) {
	result := MetricConfig{
		PrometheusName: raw.Name.GetPrometheusName(),
		OTELName:       raw.Name.GetOTELName(),
		Type:           MetricType(raw.Type),
		Description:    raw.Description,
	}

	// Always resolve to full ValueConfig
	value, err := r.resolveValue(&raw.Value, ctx)
	if err != nil {
		return MetricConfig{}, err
	}
	result.Value = value

	// Apply attribute overrides (complete replacement if specified)
	if raw.Attributes != nil {
		result.Attributes = make(map[string]string, len(raw.Attributes))
		for k, v := range raw.Attributes {
			result.Attributes[k] = v
		}
	}

	// Validate final metric
	if err := r.validateMetric(result, ctx); err != nil {
		return MetricConfig{}, err
	}

	return result, nil
}

// resolveValue resolves a value reference into fully populated ValueConfig.
// Handles three cases: instance reference, template with overrides, inline definition.
func (r *Resolver) resolveValue(raw *RawValueReference, ctx resolveContext) (ValueConfig, error) {
	// Case 1: Instance reference - return stored config
	if raw.Instance != "" {
		instance, exists := r.instanceValues[raw.Instance]
		if !exists {
			return ValueConfig{}, ctx.error(fmt.Sprintf("value instance %q not found", raw.Instance))
		}

		// No overrides allowed for instances
		if raw.Template != "" || raw.Source != nil ||
			len(raw.Transforms) > 0 || raw.Reset.Type != "" {
			return ValueConfig{}, ctx.error("cannot override instance value")
		}

		return instance, nil // Returns full config with references preserved
	}

	// Case 2: Template reference with optional overrides
	if raw.Template != "" {
		template, exists := r.templateValues[raw.Template]
		if !exists {
			return ValueConfig{}, ctx.error(fmt.Sprintf("value template %q not found", raw.Template))
		}

		// Start with template, apply overrides
		result := template

		if raw.Source != nil {
			source, sourceRef, err := r.resolveSourceReference(raw.Source, ctx)
			if err != nil {
				return ValueConfig{}, err
			}
			result.Source = source
			result.SourceRef = sourceRef // Preserve reference tracking
		}

		if len(raw.Transforms) > 0 {
			result.Transforms = raw.Transforms
		}

		if raw.Reset.Type != "" {
			result.Reset = raw.Reset
		}

		return result, nil
	}

	// Case 3: Inline definition - must have source
	if raw.Source == nil {
		return ValueConfig{}, ctx.error("value must reference instance, template, or provide inline source")
	}

	result := ValueConfig{}

	source, sourceRef, err := r.resolveSourceReference(raw.Source, ctx)
	if err != nil {
		return ValueConfig{}, err
	}
	result.Source = source
	result.SourceRef = sourceRef // Preserve reference tracking

	result.Transforms = raw.Transforms
	result.Reset = raw.Reset

	return result, nil
}

// resolveSourceReference resolves a source reference
func (r *Resolver) resolveSourceReference(raw *RawSourceReference, ctx resolveContext) (SourceConfig, *string, error) {
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
		return instance, &raw.Instance, nil // Return instance ref
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
		return result, nil, nil // No instance ref for templates
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

// validateMetric validates a resolved metric config
func (r *Resolver) validateMetric(metric MetricConfig, ctx resolveContext) error {
	// Names validated during raw syntax validation

	// Type required
	if metric.Type == "" {
		return ctx.error("type required")
	}

	// Validate type is valid
	if metric.Type != MetricTypeCounter && metric.Type != MetricTypeGauge {
		return ctx.error(fmt.Sprintf("invalid type: %s (must be counter or gauge)", metric.Type))
	}

	// Description required
	if metric.Description == "" {
		return ctx.error("description required")
	}

	// Value must be populated
	if metric.Value.Source.Type == "" {
		return ctx.error("value source required")
	}

	return nil
}
