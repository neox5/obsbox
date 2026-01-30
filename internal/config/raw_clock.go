package config

import "time"

// RawClockReference handles polymorphic clock field (instance/template/inline)
type RawClockReference struct {
	Name     string        `yaml:"name,omitempty"` // Only used in templates/instances arrays
	Instance string        `yaml:"instance,omitempty"`
	Template string        `yaml:"template,omitempty"`
	Type     *string       `yaml:"type,omitempty"`
	Interval time.Duration `yaml:"interval,omitempty"`
}

// DeepCopy creates an independent copy of the clock reference
func (c RawClockReference) DeepCopy() RawClockReference {
	clone := c

	// Deep copy pointer fields
	if c.Type != nil {
		typeCopy := *c.Type
		clone.Type = &typeCopy
	}

	return clone
}

// FindPlaceholders implements expandable for RawClockReference
func (c *RawClockReference) FindPlaceholders() []string {
	found := make(map[string]bool)

	// Scan string fields for {placeholder} patterns
	for _, name := range extractPlaceholderNames(c.Name) {
		found[name] = true
	}
	for _, name := range extractPlaceholderNames(c.Instance) {
		found[name] = true
	}
	for _, name := range extractPlaceholderNames(c.Template) {
		found[name] = true
	}

	// Convert to slice
	result := make([]string, 0, len(found))
	for name := range found {
		result = append(result, name)
	}
	return result
}

// SubstitutePlaceholders implements expandable for RawClockReference
func (c *RawClockReference) SubstitutePlaceholders(iteratorValues map[string]string) {
	c.Name = substitutePlaceholders(c.Name, iteratorValues)
	c.Instance = substitutePlaceholders(c.Instance, iteratorValues)
	c.Template = substitutePlaceholders(c.Template, iteratorValues)
}
