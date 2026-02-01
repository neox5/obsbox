# Settings Reference

[← Configuration Guide](../configuration.md) | [← Reference Index](README.md)

Detailed reference for application settings.

## Settings Configuration

**Syntax:**

```yaml
settings:
  seed: <uint64> # Optional
  internal_metrics:
    enabled: <bool> # Optional
    format: <naming_format> # Optional
```

## Seed

Optional seed for reproducible simulations.

**Parameters:**

- `seed` (uint64, optional) - Master seed for random number generation (range: 0 to 18,446,744,073,709,551,615)

**Example:**

```yaml
settings:
  seed: 12345
```

**Behavior:**

- Same seed produces identical value sequences across runs
- When omitted, uses time-based seed (logged at startup)

**Use cases:**

- Reproducible testing - Same seed generates same metrics
- Debugging - Recreate exact conditions from previous run
- Deterministic scenarios - Known metric patterns for testing

**Omitting seed:**

When seed is not specified, otelbox uses current time (nanoseconds since epoch) as seed and logs it:

```
INFO seed initialized master=1738425850123456789 stream=0 explicit=false
```

This logged seed can be used to reproduce the run:

```yaml
settings:
  seed: 1738425850123456789
```

## Internal Metrics

otelbox self-monitoring metrics for observing operational health.

**Parameters:**

- `enabled` (bool, optional) - Enable internal metrics (default: false)
- `format` (string, optional) - Naming convention ("native", "underscore", "dot", default: "native")

**Example:**

```yaml
settings:
  internal_metrics:
    enabled: true
    format: native
```

### Naming Format

**Native (default):**

```yaml
format: native
```

Uses each exporter's native convention:

- Prometheus: `otelbox_metric_name`
- OTEL: `otelbox.metric.name`

**Underscore:**

```yaml
format: underscore
```

Forces underscore-separated names for all exporters: `otelbox_metric_name`

**Dot:**

```yaml
format: dot
```

Forces dot-separated names for all exporters: `otelbox.metric.name`

**When to use:**

- `native` - Let each protocol use its convention (recommended)
- `underscore` - Need consistent naming across protocols
- `dot` - Prefer hierarchical naming across protocols

## Complete Examples

### Reproducible Simulation

```yaml
settings:
  seed: 12345
  internal_metrics:
    enabled: false
```

### Development with Monitoring

```yaml
settings:
  internal_metrics:
    enabled: true
    format: native
```

### Production with Fixed Seed

```yaml
settings:
  seed: 9876543210
  internal_metrics:
    enabled: true
    format: dot
```

### Minimal (all defaults)

```yaml
settings: {}
```

Equivalent to:

```yaml
settings:
  # seed: <time-based>
  internal_metrics:
    enabled: false
    format: native
```

## See Also

- [Export Reference](export.md) - Export configuration
