package config

import "go.yaml.in/yaml/v4"

// TransformConfig defines a transform with optional parameters.
type TransformConfig struct {
	Type    string                 `yaml:"type"`
	Options map[string]interface{} `yaml:"options,omitempty"`
}

// UnmarshalYAML handles both string and object forms for transform config.
func (t *TransformConfig) UnmarshalYAML(value *yaml.Node) error {
	// Try simple string form first
	var simple string
	if err := value.Decode(&simple); err == nil {
		t.Type = simple
		t.Options = nil
		return nil
	}

	// Fall back to full form
	type transformConfig TransformConfig
	var full transformConfig
	if err := value.Decode(&full); err != nil {
		return err
	}
	*t = TransformConfig(full)
	return nil
}

// ResetConfig defines reset behavior for values.
type ResetConfig struct {
	Type  string `yaml:"type,omitempty"`
	Value int    `yaml:"value,omitempty"`
}

// UnmarshalYAML handles both string and object forms for reset config.
func (r *ResetConfig) UnmarshalYAML(value *yaml.Node) error {
	// Try string form first (short form)
	var shortForm string
	if err := value.Decode(&shortForm); err == nil {
		r.Type = shortForm
		r.Value = 0 // default
		return nil
	}

	// Fall back to full form (object)
	type resetConfig ResetConfig // Avoid recursion
	var fullForm resetConfig
	if err := value.Decode(&fullForm); err != nil {
		return err
	}
	*r = ResetConfig(fullForm)
	return nil
}

// ValueConfig defines a simv value with transforms.
type ValueConfig struct {
	Source     string            `yaml:"source"`
	Transforms []TransformConfig `yaml:"transforms,omitempty"`
	Reset      ResetConfig       `yaml:"reset,omitempty"`
}
