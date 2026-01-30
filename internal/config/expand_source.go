package config

import "fmt"

// deepCopySource creates an independent copy of a source reference
func deepCopySource(src RawSourceReference) RawSourceReference {
	clone := src

	// Deep copy pointer fields
	if src.Type != nil {
		typeCopy := *src.Type
		clone.Type = &typeCopy
	}

	if src.Min != nil {
		minCopy := *src.Min
		clone.Min = &minCopy
	}

	if src.Max != nil {
		maxCopy := *src.Max
		clone.Max = &maxCopy
	}

	// Deep copy nested clock reference
	if src.Clock != nil {
		clockCopy := deepCopyClock(*src.Clock)
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

// expandSources expands source references containing iterator placeholders.
// Returns expanded array with iterator placeholders substituted.
func expandSources(
	sources []RawSourceReference,
	registry *IteratorRegistry,
) ([]RawSourceReference, error) {
	expanded := make([]RawSourceReference, 0)

	for i, source := range sources {
		// Find placeholders using type-specific method (includes nested clock)
		usedIterators := source.FindPlaceholders()

		if len(usedIterators) == 0 {
			// No placeholders - keep source as-is
			expanded = append(expanded, source)
			continue
		}

		// Get iterator instances
		iterators, err := registry.GetIterators(usedIterators)
		if err != nil {
			return nil, fmt.Errorf("source at index %d: %w", i, err)
		}

		// Create combination generator
		gen := NewCombinationGenerator(iterators)

		if gen.Total() == 0 {
			return nil, fmt.Errorf("source at index %d: iterator combination produces zero results", i)
		}

		// Generate one source per combination
		err = gen.ForEach(func(iteratorValues map[string]string) error {
			clone := deepCopySource(source)
			clone.SubstitutePlaceholders(iteratorValues)
			expanded = append(expanded, clone)
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("source at index %d: %w", i, err)
		}
	}

	return expanded, nil
}
