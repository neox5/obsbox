# Reference Documentation

[← Configuration Guide](../configuration.md)

Complete parameter reference for all configuration sections.

## Sections

### [File Structure](file-structure.md)

YAML file structure, top-level sections, reference types (instance/template/inline), and data types.

### [Iterators](iterators.md)

Iterator types (range, list), expansion rules, placeholder syntax, and Cartesian product generation.

### [Templates](templates.md)

Template definitions for clocks, sources, and values. Override behavior and hierarchical references.

### [Instances](instances.md)

Instance definitions for clocks, sources, and values. Sharing behavior and coherence guarantees.

### [Metrics](metrics.md)

Metric naming (simple/protocol-specific), types (counter/gauge), value references, and attributes.

### [Export](export.md)

Prometheus pull configuration and OTEL push configuration (gRPC/HTTP transports, intervals, resources).

### [Settings](settings.md)

Application settings: seed configuration and internal metrics control.

## Examples

Complete working examples in [testdata/](../../testdata/):

- [iterators.yaml](../../testdata/iterators.yaml) - Iterator expansion patterns
- [templates.yaml](../../testdata/templates.yaml) - Template usage and overrides
- [instances.yaml](../../testdata/instances.yaml) - Instance sharing and coherence
- [mq.yaml](../../testdata/mq.yaml) - IBM MQ monitoring scenario
