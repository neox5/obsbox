# Instances Reference

[← Configuration Guide](../configuration.md) | [← Reference Index](README.md)

Detailed reference for instance definitions and sharing behavior.

## Overview

Instances are named, concrete objects that are shared identically across all references. Unlike templates, instances cannot be overridden.

## Instance Types

### Clocks

Clock instances define shared timing sources.

**Syntax:**

```yaml
instances:
  clocks:
    - name: <string> # Required - instance name
      type: <string> # Required - clock type ("periodic")
      interval: <duration> # Required - update interval
```

**Usage:**

```yaml
source:
  clock:
    instance: main_tick # No overrides allowed
```

**Behavior:**

- All references share the same clock instance
- Updates synchronized across all references
- Guarantees same timing for all consumers

### Sources

Source instances define shared data generators.

**Syntax:**

```yaml
instances:
  sources:
    - name: <string> # Required - instance name
      type: <string> # Required - source type ("random_int")
      clock: <clock_reference> # Required - clock reference
      min: <int> # Required for random_int
      max: <int> # Required for random_int
```

**Usage:**

```yaml
value:
  source:
    instance: api_requests # No overrides allowed
```

**Behavior:**

- All references share the same source instance
- Guarantees data coherence across metrics
- Same value generated for all consumers at each clock tick

### Values

Value instances define shared transformation pipelines.

**Syntax:**

```yaml
instances:
  values:
    - name: <string> # Required - instance name
      source: <source_reference> # Required - source reference
      transforms: [<transform>] # Optional - transform pipeline
      reset: <reset_config> # Optional - reset behavior
```

**Usage:**

```yaml
metrics:
  - value:
      instance: total_requests # No overrides allowed
```

**Behavior:**

- All references share the same value instance
- Each metric reading sees identical values

## Instance Restrictions

Instances cannot be overridden. Any attempt to specify additional fields when referencing an instance will cause a configuration error.

**Invalid examples:**

```yaml
# ERROR: Cannot override clock instance
source:
  clock:
    instance: main_tick
    interval: 500ms # Not allowed

# ERROR: Cannot override source instance
value:
  source:
    instance: events
    max: 200 # Not allowed

# ERROR: Cannot override value instance
metrics:
  - value:
      instance: total_events
      reset: on_read # Not allowed
```

## Hierarchical References

Instances can reference templates or other instances:

```yaml
templates:
  clocks:
    - name: base_clock
      type: periodic
      interval: 1s

instances:
  clocks:
    - name: main_tick
      type: periodic
      interval: 500ms

  sources:
    - name: using_template
      type: random_int
      clock:
        template: base_clock # Instance uses template
      min: 0
      max: 100

    - name: using_instance
      type: random_int
      clock:
        instance: main_tick # Instance uses instance
      min: 0
      max: 100
```

## Coherence Guarantee

Instances provide mathematical coherence - multiple metrics derived from the same instance are guaranteed to use identical values.

**Example - Counter/Gauge Coherence:**

```yaml
instances:
  sources:
    - name: events
      type: random_int
      clock:
        type: periodic
        interval: 1s
      min: 0
      max: 100

  values:
    - name: total_events
      source:
        instance: events
      transforms: [accumulate]

    - name: recent_events
      source:
        instance: events # Same source instance
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

Both metrics derive from the same source instance, guaranteeing:

- Same raw values at each clock tick
- Mathematical relationship between counter and gauge is exact
- No drift or inconsistency between related metrics

## Examples

See [testdata/instances.yaml](../../testdata/instances.yaml) for:

- Clock instance sharing
- Source instance sharing for coherence
- Value instances with different transforms
- Multiple metrics from same instances

## See Also

- [Templates Reference](templates.md) - Overridable definitions
- [File Structure Reference](file-structure.md) - Reference type syntax
- [Iterators Reference](iterators.md) - Instance iteration
