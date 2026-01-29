package config

// RawSourceReference handles polymorphic source field (instance/template/inline)
type RawSourceReference struct {
	Name     string             `yaml:"name,omitempty"` // Only used in templates/instances arrays
	Instance string             `yaml:"instance,omitempty"`
	Template string             `yaml:"template,omitempty"`
	Type     *string            `yaml:"type,omitempty"`
	Clock    *RawClockReference `yaml:"clock,omitempty"`
	Min      *int               `yaml:"min,omitempty"`
	Max      *int               `yaml:"max,omitempty"`
}
