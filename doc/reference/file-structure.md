# File Structure Reference

[← Configuration Guide](../configuration.md) | [← Reference Index](README.md)

Detailed reference for otelbox YAML configuration file structure.

## Top-Level Sections

**Syntax:**

```yaml
iterators: # Optional - Iterator definitions
templates: # Optional - Reusable template definitions
instances: # Optional - Named instance definitions
metrics: # Required - Metric definitions
export: # Required - Export configuration
settings: # Optional - Application settings
```

**Required sections:**

- `metrics` - At least one metric must be defined
- `export` - At least one exporter must be enabled

**Optional sections:**

- `iterators` - Used when generating multiple similar configurations
- `templates` - Used for reusable definitions with override support
- `instances` - Used for shared, named objects
- `settings` - Application-level configuration

## Array Syntax

Templates and instances use array syntax with a `name` field:

```yaml
templates:
  clocks:
    - name: tick_1s
      type: periodic
      interval: 1s
    - name: tick_fast
      type: periodic
      interval: 100ms

instances:
  sources:
    - name: events
      type: random_int
      clock:
        type: periodic
        interval: 1s
      min: 0
      max: 100
```

Metrics also use array syntax but without a name field (name is defined inline):

```yaml
metrics:
  - name: events_total
    type: counter
    description: "Total events"
    value:
      instance: total_events
```

## Reference Types

Configuration fields that reference other objects support three forms:

### Instance Reference

Reference a named instance by name. No overrides allowed.

```yaml
clock:
  instance: main_tick

source:
  instance: event_source

value:
  instance: total_events
```

### Template Reference

Reference a template with optional field overrides.

```yaml
clock:
  template: tick_1s
  interval: 500ms # Override interval

source:
  template: base_events
  max: 200 # Override max only

value:
  template: counter_value
  reset: on_read # Override reset behavior
```

### Inline Definition

Define object directly without referencing template or instance.

```yaml
clock:
  type: periodic
  interval: 1s

source:
  type: random_int
  clock:
    type: periodic
    interval: 1s
  min: 0
  max: 100

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

## Data Types

### Duration

Go duration string format:

- Valid units: `ns`, `us` (or `µs`), `ms`, `s`, `m`, `h`
- Examples: `100ms`, `1s`, `5s`, `30s`, `1m`, `2h`
- Can combine units: `1m30s`

### Integer

Signed 64-bit integer:

- Range: -9,223,372,036,854,775,808 to 9,223,372,036,854,775,807
- No quotes in YAML

### String

UTF-8 text:

- Quote if contains special characters
- No length limit

### Boolean

Boolean value:

- Values: `true`, `false`
- Case-insensitive in YAML

### Map

Key-value pairs:

```yaml
key1: value1
key2: value2
```

### Array

Ordered list:

```yaml
- item1
- item2
```

Or inline:

```yaml
[item1, item2]
```

## See Also

- [Iterators Reference](iterators.md) - Iterator expansion
- [Templates Reference](templates.md) - Template definitions
- [Instances Reference](instances.md) - Instance definitions
