package config

import "fmt"

// expandClocks expands clock references containing iterator placeholders.
// Called by Expander.ExpandClocks().
func expandClocks(
	clocks []RawClockReference,
	registry *IteratorRegistry,
) ([]RawClockReference, error) {
	expanded := make([]RawClockReference, 0)

	for i, clock := range clocks {
		placeholders := findClockPlaceholders(clock)

		if len(placeholders) == 0 {
			expanded = append(expanded, clock)
			continue
		}

		iterators, err := registry.GetIterators(placeholders)
		if err != nil {
			return nil, fmt.Errorf("clock at index %d: %w", i, err)
		}

		gen := NewCombinationGenerator(iterators)

		if gen.Total() == 0 {
			return nil, fmt.Errorf("clock at index %d: iterator combination produces zero results", i)
		}

		err = gen.ForEach(func(iteratorValues map[string]string) error {
			clone := deepCopyClock(clock)
			substituteClockPlaceholders(&clone, iteratorValues)
			expanded = append(expanded, clone)
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("clock at index %d: %w", i, err)
		}
	}

	return expanded, nil
}

// deepCopyClock creates an independent copy of a clock reference
func deepCopyClock(src RawClockReference) RawClockReference {
	clone := src

	if src.Type != nil {
		typeCopy := *src.Type
		clone.Type = &typeCopy
	}

	return clone
}

// findClockPlaceholders scans clock for {placeholder} patterns
func findClockPlaceholders(c RawClockReference) []string {
	found := make(map[string]bool)

	for _, name := range extractPlaceholderNames(c.Name) {
		found[name] = true
	}
	for _, name := range extractPlaceholderNames(c.Instance) {
		found[name] = true
	}
	for _, name := range extractPlaceholderNames(c.Template) {
		found[name] = true
	}

	result := make([]string, 0, len(found))
	for name := range found {
		result = append(result, name)
	}
	return result
}

// substituteClockPlaceholders replaces {placeholder} patterns in clock
func substituteClockPlaceholders(c *RawClockReference, values map[string]string) {
	c.Name = substitutePlaceholders(c.Name, values)
	c.Instance = substitutePlaceholders(c.Instance, values)
	c.Template = substitutePlaceholders(c.Template, values)
}
