# Iterators Reference

[← Configuration Guide](../configuration.md) | [← Reference Index](README.md)

Detailed reference for iterator definitions and expansion behavior.

## Syntax

```yaml
iterators:
  - name: <string> # Required - placeholder name
    type: <type> # Required - "range" or "list"

    # For type: range
    start: <int> # Required - first value (inclusive)
    end: <int> # Required - last value (inclusive)

    # For type: list
    values: [<string>] # Required - explicit values
```

## Iterator Types

### Range Iterator

Generates sequential integer values from start to end (inclusive).

**Parameters:**

- `name` (string, required) - Iterator name used in placeholders
- `type` (string, required) - Must be `range`
- `start` (int, required) - First value (inclusive)
- `end` (int, required) - Last value (inclusive)

**Example:**

```yaml
iterators:
  - name: shard
    type: range
    start: 0
    end: 2
```

Generates: `0`, `1`, `2`

### List Iterator

Generates values from an explicit list.

**Parameters:**

- `name` (string, required) - Iterator name used in placeholders
- `type` (string, required) - Must be `list`
- `values` (array[string], required) - List of values to generate

**Example:**

```yaml
iterators:
  - name: region
    type: list
    values: [us-east, us-west, eu-central]
```

Generates: `us-east`, `us-west`, `eu-central`

## Expansion Rules

**Placeholder syntax:** `{iterator_name}` in any string field

**Cartesian product:** Multiple iterators generate all combinations

**Expansion targets:**

- Template/instance names
- Metric names
- Attribute values
- Any configuration string field

## Expansion Behavior

### Single Iterator

Each placeholder occurrence generates that many objects:

```yaml
iterators:
  - name: shard
    type: range
    start: 0
    end: 1

instances:
  clocks:
    tick_shard_{shard}:
      type: periodic
      interval: 1s
```

Expands to 2 clocks: `tick_shard_0`, `tick_shard_1`

### Multiple Iterators (Cartesian Product)

All combinations are generated:

```yaml
iterators:
  - name: region
    type: list
    values: [us, eu]
  - name: shard
    type: range
    start: 0
    end: 1

instances:
  sources:
    events_{region}_{shard}:
      type: random_int
      clock:
        type: periodic
        interval: 1s
      min: 0
      max: 100
```

Expands to 4 sources (2 regions × 2 shards):

- `events_us_0`
- `events_us_1`
- `events_eu_0`
- `events_eu_1`

## Examples

See [testdata/iterators.yaml](../../testdata/iterators.yaml) for:

- Single iterator expansion
- Multiple iterator combinations (Cartesian product)
- Usage in templates, instances, and metrics
- Attribute placeholder patterns

## See Also

- [Templates Reference](templates.md) - Template iteration
- [Instances Reference](instances.md) - Instance iteration
- [Metrics Reference](metrics.md) - Metric iteration
