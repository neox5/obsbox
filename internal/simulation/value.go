package simulation

import (
	"fmt"

	"github.com/neox5/obsbox/internal/config"
	"github.com/neox5/simv/source"
	"github.com/neox5/simv/transform"
	"github.com/neox5/simv/value"
)

// CreateValue creates a value from configuration.
func CreateValue(
	cfg config.ValueConfig,
	src source.Publisher[int],
	baseValue value.Value[int],
) (value.Value[int], error) {
	var val value.Value[int]

	if cfg.Clone != "" {
		// Clone from base value
		if baseValue == nil {
			return nil, fmt.Errorf("base value required for clone")
		}
		val = baseValue.Clone()

		// Extend transforms if specified
		if len(cfg.Transforms) > 0 {
			transforms := buildTransforms(cfg.Transforms)
			val = val.WithTransforms(transforms...)
		}
	} else {
		// Create from source
		val = createValueFromSource(src, cfg.Transforms)
	}

	// Apply reset behavior
	if cfg.Reset.Type == "on_read" {
		val = value.NewResetOnRead(val, cfg.Reset.Value)
	}

	return val, nil
}

// createValueFromSource creates a value from a source with transforms.
func createValueFromSource(src source.Publisher[int], transformCfgs []config.TransformConfig) value.Value[int] {
	transforms := buildTransforms(transformCfgs)
	return value.New(src, transforms...)
}

// buildTransforms creates transform instances from configuration.
func buildTransforms(transformCfgs []config.TransformConfig) []transform.Transformation[int] {
	var transforms []transform.Transformation[int]

	for _, tfCfg := range transformCfgs {
		switch tfCfg.Type {
		case "accumulate":
			transforms = append(transforms, transform.NewAccumulate[int]())
			// Future transforms with options can be added here
			// case "moving_average":
			//     window := tfCfg.Options["window"].(int)
			//     transforms = append(transforms, transform.NewMovingAverage[int](window))
		}
	}

	return transforms
}
