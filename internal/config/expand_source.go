package config

import "fmt"

// expandSources expands source references containing iterator placeholders.
// Called by Expander.ExpandSources().
func expandSources(
	sources []RawSourceReference,
	registry *IteratorRegistry,
) ([]RawSourceReference, error) {
	expanded := make([]RawSourceReference, 0)

	for i, source := range sources {
		placeholders := findSourcePlaceholders(source)

		if len(placeholders) == 0 {
			expanded = append(expanded, source)
			continue
		}

		iterators, err := registry.GetIterators(placeholders)
		if err != nil {
			return nil, fmt.Errorf("source at index %d: %w", i, err)
		}

		gen := NewCombinationGenerator(iterators)

		if gen.Total() == 0 {
			return nil, fmt.Errorf("source at index %d: iterator combination produces zero results", i)
		}

		err = gen.ForEach(func(iteratorValues map[string]string) error {
			clone := deepCopySource(source)
			substituteSourcePlaceholders(&clone, iteratorValues)
			expanded = append(expanded, clone)
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("source at index %d: %w", i, err)
		}
	}

	return expanded, nil
}

// deepCopySource creates an independent copy of a source reference
func deepCopySource(src RawSourceReference) RawSourceReference {
	clone := src

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

	// Nested clock - call its deep copy
	if src.Clock != nil {
		clockCopy := deepCopyClock(*src.Clock)
		clone.Clock = &clockCopy
	}

	return clone
}

// findSourcePlaceholders scans source for {placeholder} patterns (including nested clock)
func findSourcePlaceholders(s RawSourceReference) []string {
	found := make(map[string]bool)

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
		for _, name := range findClockPlaceholders(*s.Clock) {
			found[name] = true
		}
	}

	result := make([]string, 0, len(found))
	for name := range found {
		result = append(result, name)
	}
	return result
}

// substituteSourcePlaceholders replaces {placeholder} patterns in source (including nested clock)
func substituteSourcePlaceholders(s *RawSourceReference, values map[string]string) {
	s.Name = substitutePlaceholders(s.Name, values)
	s.Instance = substitutePlaceholders(s.Instance, values)
	s.Template = substitutePlaceholders(s.Template, values)

	// Recursively substitute in nested clock
	if s.Clock != nil {
		substituteClockPlaceholders(s.Clock, values)
	}
}
