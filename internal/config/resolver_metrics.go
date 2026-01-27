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

	// Check if using metric template
	if raw.Value.Template != "" && raw.Value.Inline == nil {
		// Pure template reference (check if it's a metric template)
		if template, exists := r.metrics[raw.Value.Template]; exists {
			// Apply metric template
			ctx = ctx.push("metric template", raw.Value.Template)
			result = r.applyMetricTemplate(result, template)
		}
	}

	// Resolve value (handles template/inline/overrides)
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

// applyMetricTemplate applies metric template to result
func (r *Resolver) applyMetricTemplate(result MetricConfig, template MetricConfig) MetricConfig {
	// Type from template (can be overridden by metric definition)
	if result.Type == "" {
		result.Type = template.Type
	}

	// Value from template (will be overridden if metric specifies value)
	result.Value = template.Value

	// Attributes from template (will be overridden if metric specifies attributes)
	if result.Attributes == nil && template.Attributes != nil {
		result.Attributes = make(map[string]string, len(template.Attributes))
		for k, v := range template.Attributes {
			result.Attributes[k] = v
		}
	}

	return result
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

	// Value required and already validated by resolveValue

	return nil
}
