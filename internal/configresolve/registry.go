package configresolve

import "fmt"

// registerName validates namespace uniqueness and registers the name
func (r *Resolver) registerName(name string, entityType string) error {
	if existingType, exists := r.registeredNames[name]; exists {
		return fmt.Errorf("name %q already used by %s, cannot reuse for %s",
			name, existingType, entityType)
	}
	r.registeredNames[name] = entityType
	return nil
}

// getStringValue safely dereferences a string pointer
func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
