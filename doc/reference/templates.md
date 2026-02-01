# Templates Reference

[← Configuration Guide](../configuration.md) | [← Reference Index](README.md)

Detailed reference for template definitions and override behavior.

## Overview

Templates are reusable configuration definitions that support field overrides when referenced.

## Template Types

### Clocks

Clock templates define reusable timing patterns.

**Syntax:**

```yaml
templates:
  clocks:
    - name: <string> # Required - template name
      type: <string> # Required - clock type ("periodic")
      interval: <duration> # Required - update interval
```

**Usage:**

```yaml
source:
  clock:
    template: tick_1s
    interval: 500ms # Optional override
```

### Sources

Source templates define reusable data generators.

**Syntax:**

```yaml
templates:
  sources:
    - name: <string> # Required - template name
      type: <string> # Required - source type ("random_int")
      clock: <clock_reference> # Required - clock reference
      min: <int> # Required for random_int
      max: <int> # Required for random_int
```

**Usage:**

```yaml
value:
  source:
    template: base_events
    max: 200 # Override max only
```

### Values

Value templates define reusable transformation pipelines.

**Syntax:**

```yaml
templates:
  values:
    - name: <string> # Required - template name
      source: <source_reference> # Required - source reference
      transforms: [<transform>] # Optional - transform pipeline
      reset: <reset_config> # Optional - reset behavior
```

**Usage:**

```yaml
metrics:
  - value:
      template: counter_value
      reset: on_read # Override reset
```

## Override Behavior

**Rules:**

1. Only specified fields are overridden
2. Unspecified fields use template values
3. Nested objects can be partially overridden
4. Arrays are replaced entirely (not merged)

**Examples:**

Clock override:

```yaml
templates:
  clocks:
    - name: base_tick
      type: periodic
      interval: 1s

# Usage - override interval only
source:
  clock:
    template: base_tick
    interval: 500ms
    # type still 'periodic' from template
```

Source override:

```yaml
templates:
  sources:
    - name: base_events
      type: random_int
      clock:
        type: periodic
        interval: 1s
      min: 0
      max: 100

# Usage - override max and clock
value:
  source:
    template: base_events
    max: 200
    clock:
      template: tick_fast
    # min still 0 from template
    # type still random_int from template
```

## Hierarchical References

Templates can reference other templates:

```yaml
templates:
  clocks:
    - name: base_clock
      type: periodic
      interval: 1s

  sources:
    - name: base_source
      type: random_int
      clock:
        template: base_clock # Reference clock template
      min: 0
      max: 100

  values:
    - name: base_value
      source:
        template: base_source # Reference source template
      transforms: [accumulate]
```

Overrides cascade through hierarchy:

```yaml
metrics:
  - value:
      template: base_value
      source:
        template: base_source
        clock:
          template: base_clock
          interval: 500ms # Overrides base_clock interval
        max: 200 # Overrides base_source max
```

## Transform Configuration

### Accumulate Transform

Converts stream to monotonically increasing counter (running sum).

**Short form:**

```yaml
transforms: [accumulate]
```

**Full form:**

```yaml
transforms:
  - type: accumulate
```

## Reset Configuration

Defines when and how values reset.

### Reset On Read

**Short form (reset to zero):**

```yaml
reset: on_read
```

**Full form (reset to specific value):**

```yaml
reset:
  type: on_read
  value: 100
```

**Parameters:**

- `type` (string, required) - Reset trigger ("on_read")
- `value` (int, optional) - Reset target value (default: 0)

**Behavior:** Value resets after each read operation. Useful for gauge semantics (window-based metrics).

## Examples

See [testdata/templates.yaml](../../testdata/templates.yaml) for:

- Clock template references
- Source templates with overrides
- Value templates with inline sources
- Hierarchical template references
- Reset behavior patterns

## See Also

- [Instances Reference](instances.md) - Non-overridable shared objects
- [File Structure Reference](file-structure.md) - Reference type syntax
- [Iterators Reference](iterators.md) - Template iteration
