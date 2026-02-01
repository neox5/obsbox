# Metrics Reference

[← Configuration Guide](../configuration.md) | [← Reference Index](README.md)

Detailed reference for metric definitions.

## Metric Definition

**Syntax:**

```yaml
metrics:
  - name: <metric_name>              # Simple form
    name:                            # Or full form
      prometheus: <prom_name>
      otel: <otel_name>
    type: <metric_type>              # Required - "counter" or "gauge"
    description: <help_text>         # Required
    value: <value_reference>         # Required
    attributes:                      # Optional
      <key>: <value>
```

## Naming

### Simple Form

Single string used for both Prometheus and OTEL:

```yaml
metrics:
  - name: events_total
    type: counter
    description: "Total events"
```

Both exporters use `events_total`.

### Full Form

Protocol-specific names:

```yaml
metrics:
  - name:
      prometheus: app_events_total
      otel: app.events.total
    type: counter
    description: "Total events"
```

Prometheus uses `app_events_total`, OTEL uses `app.events.total`.

**When to use:**

- Simple form: When naming conventions align
- Full form: When protocols have different conventions (underscores vs dots)

## Metric Types

### Counter

Monotonically increasing value (never decreases).

**Characteristics:**

- Value never decreases
- Used for cumulative metrics
- Examples: total requests, bytes sent, errors encountered

**Example:**

```yaml
metrics:
  - name: requests_total
    type: counter
    description: "Total requests"
    value:
      source:
        type: random_int
        clock:
          type: periodic
          interval: 1s
        min: 0
        max: 100
      transforms: [accumulate]
```

### Gauge

Value that can increase or decrease.

**Characteristics:**

- Value can go up or down
- Used for current state metrics
- Examples: active connections, queue depth, memory usage

**Example:**

```yaml
metrics:
  - name: queue_depth
    type: gauge
    description: "Current queue depth"
    value:
      source:
        type: random_int
        clock:
          type: periodic
          interval: 1s
        min: 0
        max: 1000
```

## Value References

Metrics reference values in three ways:

### Instance Reference

```yaml
metrics:
  - name: events_total
    type: counter
    value:
      instance: total_events
```

### Template Reference

```yaml
metrics:
  - name: events_total
    type: counter
    value:
      template: counter_value
      reset: on_read # Optional override
```

### Inline Definition

```yaml
metrics:
  - name: events_total
    type: counter
    value:
      source:
        type: random_int
        clock:
          type: periodic
          interval: 1s
        min: 0
        max: 100
      transforms: [accumulate]
```

## Attributes

Key-value pairs attached to metrics. Called "labels" in Prometheus, "attributes" in OTEL.

**Syntax:**

```yaml
metrics:
  - name: requests_total
    type: counter
    value:
      instance: total_requests
    attributes:
      service: api
      environment: production
      region: us-east
```

**Constraints:**

- Keys must match pattern: `[a-zA-Z_][a-zA-Z0-9_]*`
- Keys cannot start with `__` (reserved prefix)
- Values are strings

**With Iterators:**

Attributes can use iterator placeholders:

```yaml
iterators:
  - name: region
    type: list
    values: [us, eu]

metrics:
  - name: requests_total
    type: counter
    value:
      source:
        type: random_int
        clock:
          type: periodic
          interval: 1s
        min: 0
        max: 100
      transforms: [accumulate]
    attributes:
      region: "{region}"
```

Generates 2 metrics with different attribute values:

- `requests_total{region="us"}`
- `requests_total{region="eu"}`

## Examples

See [testdata/](../../testdata/) for:

- [instances.yaml](../../testdata/instances.yaml) - Counter/gauge pairs from shared sources
- [iterators.yaml](../../testdata/iterators.yaml) - Metrics with iterator placeholders
- [templates.yaml](../../testdata/templates.yaml) - Metrics using template values
- [mq.yaml](../../testdata/mq.yaml) - Protocol-specific naming

## See Also

- [Instances Reference](instances.md) - Value instance sharing
- [Templates Reference](templates.md) - Value template usage
- [Iterators Reference](iterators.md) - Metric iteration
