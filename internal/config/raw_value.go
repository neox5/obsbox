package config

import "go.yaml.in/yaml/v4"

// RawValueReference handles polymorphic value field (instance/template/inline)
type RawValueReference struct {
	Name       string              `yaml:"name,omitempty"` // Only used in templates/instances arrays
	Instance   string              `yaml:"instance,omitempty"`
	Template   string              `yaml:"template,omitempty"`
	Source     *RawSourceReference `yaml:"source,omitempty"`
	Transforms []TransformConfig   `yaml:"transforms,omitempty"`
	Reset      ResetConfig         `yaml:"reset,omitempty"`
}

// DeepCopy creates an independent copy of the value reference
func (v RawValueReference) DeepCopy() RawValueReference {
	clone := v

	// Deep copy nested source reference
	if v.Source != nil {
		sourceCopy := v.Source.DeepCopy()
		clone.Source = &sourceCopy
	}

	// Deep copy transforms slice
	if len(v.Transforms) > 0 {
		clone.Transforms = make([]TransformConfig, len(v.Transforms))
		copy(clone.Transforms, v.Transforms)
	}

	// Reset config is plain struct, no pointers to copy

	return clone
}

// FindPlaceholders implements expandable for RawValueReference
func (v *RawValueReference) FindPlaceholders() []string {
	found := make(map[string]bool)

	// Scan own string fields
	for _, name := range extractPlaceholderNames(v.Name) {
		found[name] = true
	}
	for _, name := range extractPlaceholderNames(v.Instance) {
		found[name] = true
	}
	for _, name := range extractPlaceholderNames(v.Template) {
		found[name] = true
	}

	// Recursively scan nested source
	if v.Source != nil {
		for _, name := range v.Source.FindPlaceholders() {
			found[name] = true
		}
	}

	// Convert to slice
	result := make([]string, 0, len(found))
	for name := range found {
		result = append(result, name)
	}
	return result
}

// SubstitutePlaceholders implements expandable for RawValueReference
func (v *RawValueReference) SubstitutePlaceholders(iteratorValues map[string]string) {
	v.Name = substitutePlaceholders(v.Name, iteratorValues)
	v.Instance = substitutePlaceholders(v.Instance, iteratorValues)
	v.Template = substitutePlaceholders(v.Template, iteratorValues)

	// Recursively substitute in nested source
	if v.Source != nil {
		v.Source.SubstitutePlaceholders(iteratorValues)
	}
}

// TransformConfig defines a transform operation
type TransformConfig struct {
	Type string
}

// UnmarshalYAML handles both string and object forms for transforms
func (t *TransformConfig) UnmarshalYAML(value *yaml.Node) error {
	// Try string form first (shorthand)
	var simple string
	if err := value.Decode(&simple); err == nil {
		t.Type = simple
		return nil
	}

	// Fall back to object form
	type transformConfig struct {
		Type string `yaml:"type"`
	}
	var full transformConfig
	if err := value.Decode(&full); err != nil {
		return err
	}
	t.Type = full.Type
	return nil
}

// ResetConfig defines reset behavior
type ResetConfig struct {
	Type  string
	Value int
}

// UnmarshalYAML handles both string and object forms for reset
func (r *ResetConfig) UnmarshalYAML(value *yaml.Node) error {
	// Try string form first (shorthand)
	var simple string
	if err := value.Decode(&simple); err == nil {
		r.Type = simple
		r.Value = 0 // Default value for shorthand
		return nil
	}

	// Fall back to object form
	type resetConfig struct {
		Type  string `yaml:"type"`
		Value int    `yaml:"value"`
	}
	var full resetConfig
	if err := value.Decode(&full); err != nil {
		return err
	}
	r.Type = full.Type
	r.Value = full.Value
	return nil
}
