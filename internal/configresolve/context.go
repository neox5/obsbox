package configresolve

import (
	"fmt"
	"strings"
)

// resolveContext tracks resolution path for error messages
type resolveContext []string

func (ctx resolveContext) push(component, name string) resolveContext {
	return append(ctx, fmt.Sprintf("%s %q", component, name))
}

func (ctx resolveContext) error(msg string) error {
	if len(ctx) == 0 {
		return fmt.Errorf(msg)
	}

	var b strings.Builder
	b.WriteString(msg)
	// Print stack top-down (metric → templates → error)
	for i := len(ctx) - 1; i >= 0; i-- {
		b.WriteString("\n  in ")
		b.WriteString(ctx[i])
	}
	return fmt.Errorf(b.String())
}
