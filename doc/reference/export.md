# Export Reference

[← Configuration Guide](../configuration.md) | [← Reference Index](README.md)

Detailed reference for export configuration.

## Export Configuration

**Syntax:**

```yaml
export:
  prometheus: # Optional
    enabled: <bool>
    port: <int>
    path: <string>

  otel: # Optional
    enabled: <bool>
    transport: <string>
    host: <string>
    port: <int>
    interval: <interval_config>
    resource: <map>
    headers: <map>
```

**Constraints:**

- At least one exporter must be enabled
- Only one exporter can be enabled at a time (prevents read conflicts)

## Prometheus Export

Pull-based HTTP endpoint for Prometheus scraping.

**Parameters:**

- `enabled` (bool, required) - Enable Prometheus exporter
- `port` (int, optional) - HTTP port (default: 9090, range: 1-65535)
- `path` (string, optional) - Metrics endpoint path (default: `/metrics`)

**Example:**

```yaml
export:
  prometheus:
    enabled: true
    port: 9090
    path: /metrics
```

Metrics available at: `http://localhost:9090/metrics`

**Prometheus Configuration:**

```yaml
scrape_configs:
  - job_name: otelbox
    static_configs:
      - targets: ["localhost:9090"]
```

## OTEL Export

Push-based OTLP export to collectors.

**Parameters:**

- `enabled` (bool, required) - Enable OTEL exporter
- `transport` (string, optional) - OTLP transport ("grpc" or "http", default: "grpc")
- `host` (string, optional) - OTLP endpoint host (default: "localhost")
- `port` (int, optional) - OTLP endpoint port (default: 4317 for grpc, 4318 for http)
- `interval` (interval_config, required) - Export intervals
- `resource` (map[string]string, optional) - Resource attributes
- `headers` (map[string]string, optional) - Custom HTTP headers

### Transport Types

**gRPC Transport:**

```yaml
export:
  otel:
    enabled: true
    transport: grpc
    host: localhost
    port: 4317
    interval: 10s
```

- Default port: 4317
- Binary protocol, typically more efficient

**HTTP Transport:**

```yaml
export:
  otel:
    enabled: true
    transport: http
    host: localhost
    port: 4318
    interval: 10s
```

- Default port: 4318
- JSON-based protocol

### Interval Configuration

**Simple form (same for collection and push):**

```yaml
export:
  otel:
    enabled: true
    interval: 10s
```

**Detailed form (different intervals):**

```yaml
export:
  otel:
    enabled: true
    interval:
      read: 1s # Collect metrics every second
      push: 10s # Push batch every 10 seconds
```

**Parameters:**

- `read` (duration) - How often to collect metric values internally
- `push` (duration) - How often to push batches to collector

**Use cases:**

- Same interval: Simple configuration, immediate push
- Different intervals: Batch multiple collections before pushing (reduces network overhead)

### Resource Attributes

Resource attributes identify the source of metrics.

**Default:**

```yaml
service.name: otelbox
service.version: dev
```

**Custom:**

```yaml
export:
  otel:
    enabled: true
    interval: 10s
    resource:
      service.name: myapp
      service.version: 1.2.3
      deployment.environment: production
      cloud.provider: aws
      cloud.region: us-east-1
```

Follow OpenTelemetry semantic conventions for standard attributes.

### Custom Headers

Add custom HTTP headers to OTLP requests:

```yaml
export:
  otel:
    enabled: true
    interval: 10s
    headers:
      Authorization: Bearer secret-token
      X-Custom-Header: custom-value
```

**Use cases:**

- Authentication tokens
- API keys
- Custom routing headers

## Complete Examples

### Prometheus Only

```yaml
export:
  prometheus:
    enabled: true
    port: 9090
    path: /metrics
```

### OTEL gRPC with Batching

```yaml
export:
  otel:
    enabled: true
    transport: grpc
    host: localhost
    port: 4317
    interval:
      read: 1s # Collect every second
      push: 10s # Push batch every 10 seconds
    resource:
      service.name: otelbox
```

### OTEL HTTP with Authentication

```yaml
export:
  otel:
    enabled: true
    transport: http
    host: collector.example.com
    port: 4318
    interval:
      read: 1s
      push: 30s
    resource:
      service.name: otelbox
      service.version: 1.0.0
      deployment.environment: production
    headers:
      Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

## See Also

- [Settings Reference](settings.md) - Application settings
- [Metrics Reference](metrics.md) - Metric definitions
