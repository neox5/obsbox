package config

import (
	"regexp"
	"strings"
	"time"
)

var attributeNameRegex = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

// MetricConfig defines a fully resolved metric
type MetricConfig struct {
	PrometheusName string
	OTELName       string
	Type           MetricType
	Description    string
	Value          ValueConfig
	Attributes     map[string]string
}

// MetricType defines the semantic type of a metric
type MetricType string

const (
	MetricTypeCounter MetricType = "counter"
	MetricTypeGauge   MetricType = "gauge"
)

// isValidAttributeName checks if an attribute name follows conventions
func isValidAttributeName(name string) bool {
	if len(name) == 0 {
		return false
	}
	if strings.HasPrefix(name, "__") {
		return false
	}
	return attributeNameRegex.MatchString(name)
}

// ValueConfig defines a fully resolved value with embedded components.
type ValueConfig struct {
	Source     SourceConfig
	SourceRef  *string // Instance name if source is shared
	Transforms []TransformConfig
	Reset      ResetConfig
}

// SourceConfig defines a fully resolved source with embedded clock
type SourceConfig struct {
	Type     string
	Clock    ClockConfig
	ClockRef *string // Instance name if clock is shared
	Min      int
	Max      int
}

// ClockConfig defines a fully resolved clock
type ClockConfig struct {
	Type     string
	Interval time.Duration
}

// TransformConfig defines a transform operation
type TransformConfig struct {
	Type string
}

// ResetConfig defines reset behavior
type ResetConfig struct {
	Type  string
	Value int
}
