package simulation

import (
	"fmt"

	"github.com/neox5/obsbox/internal/config"
	"github.com/neox5/simv/clock"
	"github.com/neox5/simv/source"
)

// CreateSource creates a source from configuration.
func CreateSource(cfg config.SourceConfig, clk clock.Clock) (source.Publisher[int], error) {
	switch cfg.Type {
	case "random_int":
		return source.NewRandomIntSource(clk, cfg.Min, cfg.Max), nil
	default:
		return nil, fmt.Errorf("unknown source type: %s", cfg.Type)
	}
}
