package config

import "go.yaml.in/yaml/v4"

// TransformConfig defines a transform with optional parameters.
type TransformConfig struct {
	Type    string                 `yaml:"type"`
	Options map[string]interface{} `yaml:"options,omitempty"`
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

// ValueConfig defines a simv value with transforms or derivation.
type ValueConfig struct {
	Source     string            `yaml:"source,omitempty"`
	Clone      string            `yaml:"clone,omitempty"`
	Transforms []TransformConfig `yaml:"transforms,omitempty"`
	Reset      ResetConfig       `yaml:"reset,omitempty"`
}
