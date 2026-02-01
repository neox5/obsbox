# Configuration Guide

[← Back to README](../README.md)

Complete guide to understanding and configuring otelbox for telemetry signal generation.

## Overview

otelbox uses YAML configuration files to define how metrics are generated and exposed. The configuration system is built around a few core concepts that work together to provide flexibility and reusability.

**Core Concepts:**

- **Iterators** - Generate multiple similar configurations from patterns
- **Templates** - Reusable definitions that can be customized when referenced
- **Instances** - Named, shared objects used identically across references
- **Metrics** - Map generated values to exposed metrics
- **Export** - Configure how metrics are exposed (Prometheus/OTEL)
- **Settings** - Application-level configuration

## Configuration File Structure

```yaml
iterators: # Generate configurations from patterns (optional)
templates: # Reusable definitions with override support (optional)
instances: # Named, shared objects (optional)
metrics: # Metric definitions (required)
export: # Metric exposition configuration (required)
settings: # Application settings (optional)
```

Only `metrics` and `export` are required. The others are organizational tools for managing complex configurations.

→ Full structure: [reference/file-structure.md](reference/file-structure.md)

## Quick Start Example

Minimal configuration generating a single counter metric:

```yaml
instances:
  clocks:
    tick:
      type: periodic
      interval: 1s

  sources:
    events:
      type: random_int
      clock:
        instance: tick
      min: 0
      max: 10

  values:
    total_events:
      source:
        instance: events
      transforms: [accumulate]

metrics:
  - name: app_events_total
    type: counter
    description: "Total events processed"
    value:
      instance: total_events
    attributes:
      service: myapp

export:
  prometheus:
    enabled: true
    port: 9090
    path: /metrics
```

This creates a counter that increments by 0-10 every second, exposed at `http://localhost:9090/metrics`.

## Iterators

Generate multiple similar configurations from patterns using `{placeholder}` syntax.

**Two types:**

- `range`: Sequential numbers (0, 1, 2, ...)
- `list`: Explicit values (us, eu, asia, ...)

**Example:** Generate 6 queue metrics (2 regions × 3 queues):

```yaml
iterators:
  - name: region
    type: list
    values: [us, eu]
  - name: queue_id
    type: range
    start: 1
    end: 3

instances:
  sources:
    queue_{region}_{queue_id}_depth:
      type: random_int
      clock:
        type: periodic
        interval: 1s
      min: 0
      max: 1000

metrics:
  - name: queue_depth
    type: gauge
    description: "Current queue depth"
    value:
      source:
        instance: queue_{region}_{queue_id}_depth
    attributes:
      region: "{region}"
      queue: "QUEUE_{queue_id}"
```

Expands to 6 metrics with attributes like `{region=us, queue=QUEUE_1}`, `{region=us, queue=QUEUE_2}`, etc.

→ Full syntax: [reference/iterators.md](reference/iterators.md)  
→ Complete example: [testdata/iterators.yaml](../testdata/iterators.yaml)

## Templates

Templates are reusable configuration definitions that can be customized when referenced.

**Key characteristics:**

- Define once, reference many times
- Support field overrides when referenced
- Can reference other templates

**Example:** Create variations from a base pattern:

```yaml
templates:
  sources:
    base_events:
      type: random_int
      clock:
        type: periodic
        interval: 1s
      min: 0
      max: 100

metrics:
  - name: low_rate_events
    type: counter
    description: "Low rate events"
    value:
      source:
        template: base_events
        max: 50 # Override max
      transforms: [accumulate]

  - name: high_rate_events
    type: counter
    description: "High rate events"
    value:
      source:
        template: base_events
        max: 200 # Different override
      transforms: [accumulate]
```

**Override rules:**

- Only specified fields are overridden
- Unspecified fields use template values
- Nested objects can be partially overridden

→ Full syntax: [reference/templates.md](reference/templates.md)  
→ Complete example: [testdata/templates.yaml](../testdata/templates.yaml)

## Instances

Instances are named, concrete objects that are shared identically across all references. Unlike templates, instances cannot be overridden.

**Key characteristics:**

- Reference by name without modification
- Guarantee identical behavior across all references
- No overrides allowed

**When to use instances:**

Use instances when you need mathematical coherence - multiple metrics derived from the exact same underlying data:

```yaml
instances:
  sources:
    api_requests:
      type: random_int
      clock:
        type: periodic
        interval: 1s
      min: 0
      max: 100

  values:
    total_requests:
      source:
        instance: api_requests
      transforms: [accumulate]

    recent_requests:
      source:
        instance: api_requests # Same source - guaranteed coherence
      transforms: [accumulate]
      reset: on_read

metrics:
  - name: requests_total
    type: counter
    value:
      instance: total_requests

  - name: requests_current
    type: gauge
    value:
      instance: recent_requests
```

Both metrics derive from the same source, guaranteeing mathematical consistency.

→ Full syntax: [reference/instances.md](reference/instances.md)  
→ Complete example: [testdata/instances.yaml](../testdata/instances.yaml)

## Reference Types

Configuration fields that reference other objects support three forms:

**1. Instance reference** - Use shared object by name:

```yaml
value:
  instance: total_events
```

**2. Template reference** - Use template with optional overrides:

```yaml
value:
  template: counter_value
  reset: on_read
```

**3. Inline definition** - Define object directly:

```yaml
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

→ Full details: [reference/file-structure.md](reference/file-structure.md#reference-types)

## Metrics

Metrics map generated values to exposed telemetry.

**Basic metric:**

```yaml
metrics:
  - name: events_total
    type: counter
    description: "Total events processed"
    value:
      instance: total_events
    attributes:
      service: myapp
      environment: production
```

**Protocol-specific names:**

```yaml
metrics:
  - name:
      prometheus: app_events_total
      otel: app.events.total
    type: counter
    description: "Total events"
    value:
      instance: total_events
```

**Metric types:**

- `counter` - Monotonically increasing value
- `gauge` - Value that can increase or decrease

→ Full syntax: [reference/metrics.md](reference/metrics.md)

## Export Configuration

Export configuration determines how metrics are exposed to collectors.

**Prometheus (Pull-based):**

```yaml
export:
  prometheus:
    enabled: true
    port: 9090
    path: /metrics
```

**OTEL (Push-based):**

```yaml
export:
  otel:
    enabled: true
    transport: grpc
    host: localhost
    port: 4317
    interval: 10s
    resource:
      service.name: myapp
      deployment.environment: production
```

**Constraints:**

- At least one exporter must be enabled
- Only one exporter can be enabled at a time

→ Full syntax: [reference/export.md](reference/export.md)

## Settings

Application-level configuration for otelbox behavior.

**Seed for reproducible simulations:**

```yaml
settings:
  seed: 12345
```

When omitted, a time-based seed is used (logged at startup for reproduction).

**Internal metrics for self-monitoring:**

```yaml
settings:
  internal_metrics:
    enabled: true
    format: native
```

→ Full syntax: [reference/settings.md](reference/settings.md)

## Common Patterns

### Template with Variations

Create multiple metrics from a base template with different parameters:

```yaml
templates:
  sources:
    base_load:
      type: random_int
      clock:
        type: periodic
        interval: 1s
      min: 0
      max: 100

metrics:
  - name: low_load_total
    type: counter
    value:
      source:
        template: base_load
        max: 50
      transforms: [accumulate]

  - name: high_load_total
    type: counter
    value:
      source:
        template: base_load
        max: 200
      transforms: [accumulate]
```

### Shared Clock Pattern

Multiple sources driven by single clock for synchronized updates:

```yaml
instances:
  clocks:
    main_tick:
      type: periodic
      interval: 1s

  sources:
    source_a:
      type: random_int
      clock:
        instance: main_tick
      min: 0
      max: 100

    source_b:
      type: random_int
      clock:
        instance: main_tick
      min: 0
      max: 50
```

### Coherent Counter/Gauge Pairs

Create counter and gauge from same source for mathematical coherence:

```yaml
instances:
  sources:
    events:
      type: random_int
      clock:
        type: periodic
        interval: 1s
      min: 0
      max: 100

  values:
    total_events:
      source:
        instance: events
      transforms: [accumulate]

    recent_events:
      source:
        instance: events
      transforms: [accumulate]
      reset: on_read

metrics:
  - name: events_total
    type: counter
    value:
      instance: total_events

  - name: events_current
    type: gauge
    value:
      instance: recent_events
```

### Iterator-Based Metric Families

Generate metrics for multiple entities using iterators - see [testdata/iterators.yaml](../testdata/iterators.yaml) for complete example.

### Multiple Independent Update Frequencies

Different metric groups with different update rates - see [testdata/instances.yaml](../testdata/instances.yaml) for complete example.

## Reference Documentation

Complete parameter reference for all configuration sections:

- [File Structure Reference](reference/file-structure.md) - YAML structure and reference types
- [Iterators Reference](reference/iterators.md) - Iterator types and expansion
- [Templates Reference](reference/templates.md) - Template definitions and overrides
- [Instances Reference](reference/instances.md) - Instance definitions and sharing
- [Metrics Reference](reference/metrics.md) - Metric parameters and types
- [Export Reference](reference/export.md) - Prometheus and OTEL configuration
- [Settings Reference](reference/settings.md) - Application settings
