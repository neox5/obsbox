package config

import "fmt"

// deepCopyClock creates an independent copy of a clock reference
func deepCopyClock(src RawClockReference) RawClockReference {
	clone := src

	// Deep copy pointer fields
	if src.Type != nil {
		typeCopy := *src.Type
		clone.Type = &typeCopy
	}

	return clone
}

// FindPlaceholders implements IteratorExpandable for RawClockReference
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

// SubstitutePlaceholders implements IteratorExpandable for RawClockReference
func (c *RawClockReference) SubstitutePlaceholders(iteratorValues map[string]string) {
	c.Name = substitutePlaceholders(c.Name, iteratorValues)
	c.Instance = substitutePlaceholders(c.Instance, iteratorValues)
	c.Template = substitutePlaceholders(c.Template, iteratorValues)
}

// expandClocks expands clock references containing iterator placeholders.
// Returns expanded array with iterator placeholders substituted.
func expandClocks(
	clocks []RawClockReference,
	registry *IteratorRegistry,
) ([]RawClockReference, error) {
	expanded := make([]RawClockReference, 0)

	for i, clock := range clocks {
		// Find placeholders using type-specific method
		usedIterators := clock.FindPlaceholders()

		if len(usedIterators) == 0 {
			// No placeholders - keep clock as-is
			expanded = append(expanded, clock)
			continue
		}

		// Get iterator instances
		iterators, err := registry.GetIterators(usedIterators)
		if err != nil {
			return nil, fmt.Errorf("clock at index %d: %w", i, err)
		}

		// Create combination generator
		gen := NewCombinationGenerator(iterators)

		if gen.Total() == 0 {
			return nil, fmt.Errorf("clock at index %d: iterator combination produces zero results", i)
		}

		// Generate one clock per combination
		err = gen.ForEach(func(iteratorValues map[string]string) error {
			clone := deepCopyClock(clock)
			clone.SubstitutePlaceholders(iteratorValues)
			expanded = append(expanded, clone)
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("clock at index %d: %w", i, err)
		}
	}

	return expanded, nil
}
