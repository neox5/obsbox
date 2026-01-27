package simulation

import (
	"fmt"

	"github.com/neox5/obsbox/internal/config"
	"github.com/neox5/simv/source"
	"github.com/neox5/simv/transform"
	"github.com/neox5/simv/value"
)

// ValueWrapper wraps simv Value for easier management
type ValueWrapper struct {
	*value.Value[int]
}

// CreateValue creates a value from configuration.
// The value is started and ready to receive updates.
func CreateValue(
	cfg config.ValueConfig,
	src source.Publisher[int],
) (*ValueWrapper, error) {
	if src == nil {
		return nil, fmt.Errorf("source required for value")
	}

	// Create value
	val := value.New(src)

	// Add transforms
	if len(cfg.Transforms) > 0 {
		transforms, err := buildTransforms(cfg.Transforms)
		if err != nil {
			return nil, err
		}
		for _, t := range transforms {
			val.AddTransform(t)
		}
	}

	// Apply reset behavior
	if cfg.Reset.Type == "on_read" {
		val.EnableResetOnRead(cfg.Reset.Value)
	}

	// Start the value (begins receiving updates)
	val.Start()

	return &ValueWrapper{Value: val}, nil
}

// buildTransforms creates transform instances from configuration.
func buildTransforms(transformCfgs []config.TransformConfig) ([]transform.Transformation[int], error) {
	var transforms []transform.Transformation[int]

	for _, tfCfg := range transformCfgs {
		switch tfCfg.Type {
		case "accumulate":
			transforms = append(transforms, transform.NewAccumulate[int]())
		case "":
			return nil, fmt.Errorf("transform type cannot be empty")
		default:
			return nil, fmt.Errorf("unknown transform type: %q", tfCfg.Type)
		}
	}

	return transforms, nil
}

// GetEffectiveClock returns the clock to use for this value.
// If value has explicit clock (override), use it.
// Otherwise, inherit from source.
func GetEffectiveClock(cfg config.ValueConfig) config.ClockConfig {
	if cfg.Clock.Type != "" {
		return cfg.Clock
	}
	return cfg.Source.Clock
}
