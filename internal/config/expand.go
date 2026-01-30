package config

import (
	"regexp"
	"strings"
)

// iteratorPattern matches {iterator_name} placeholders in strings
var iteratorPattern = regexp.MustCompile(`\{([a-zA-Z_][a-zA-Z0-9_]*)\}`)

// IteratorExpandable types can find and substitute {placeholder} patterns
type IteratorExpandable interface {
	FindPlaceholders() []string
	SubstitutePlaceholders(iteratorValues map[string]string)
}

// substitutePlaceholders replaces {name} patterns in a string with values
func substitutePlaceholders(s string, iteratorValues map[string]string) string {
	result := s
	for name, value := range iteratorValues {
		placeholder := "{" + name + "}"
		result = strings.ReplaceAll(result, placeholder, value)
	}
	return result
}

// extractPlaceholderNames extracts placeholder names from {name} patterns in a string
func extractPlaceholderNames(s string) []string {
	matches := iteratorPattern.FindAllStringSubmatch(s, -1)
	if len(matches) == 0 {
		return nil
	}

	names := make([]string, len(matches))
	for i, match := range matches {
		names[i] = match[1] // Capture group 1 contains the iterator name
	}
	return names
}
