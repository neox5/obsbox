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

// DeepCopy creates an independent copy of the source reference
func (s *RawSourceReference) DeepCopy() RawSourceReference {
	clone := *s

	// Deep copy pointer fields
	if s.Type != nil {
		typeCopy := *s.Type
		clone.Type = &typeCopy
	}

	if s.Min != nil {
		minCopy := *s.Min
		clone.Min = &minCopy
	}

	if s.Max != nil {
		maxCopy := *s.Max
		clone.Max = &maxCopy
	}

	// Deep copy nested clock reference
	if s.Clock != nil {
		clockCopy := s.Clock.DeepCopy()
		clone.Clock = &clockCopy
	}

	return clone
}

// FindPlaceholders implements IteratorExpandable for RawSourceReference
func (s *RawSourceReference) FindPlaceholders() []string {
	found := make(map[string]bool)

	// Scan own string fields
	for _, name := range extractPlaceholderNames(s.Name) {
		found[name] = true
	}
	for _, name := range extractPlaceholderNames(s.Instance) {
		found[name] = true
	}
	for _, name := range extractPlaceholderNames(s.Template) {
		found[name] = true
	}

	// Recursively scan nested clock
	if s.Clock != nil {
		for _, name := range s.Clock.FindPlaceholders() {
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

// SubstitutePlaceholders implements IteratorExpandable for RawSourceReference
func (s *RawSourceReference) SubstitutePlaceholders(iteratorValues map[string]string) {
	s.Name = substitutePlaceholders(s.Name, iteratorValues)
	s.Instance = substitutePlaceholders(s.Instance, iteratorValues)
	s.Template = substitutePlaceholders(s.Template, iteratorValues)

	// Recursively substitute in nested clock
	if s.Clock != nil {
		s.Clock.SubstitutePlaceholders(iteratorValues)
	}
}
