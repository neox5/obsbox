package simulation

import (
	"fmt"

	"github.com/neox5/obsbox/internal/config"
	"github.com/neox5/simv/clock"
)

// CreateClock creates a clock from configuration.
func CreateClock(cfg config.ClockConfig) (clock.Clock, error) {
	switch cfg.Type {
	case "periodic":
		return clock.NewPeriodicClock(cfg.Interval), nil
	default:
		return nil, fmt.Errorf("unknown clock type: %s", cfg.Type)
	}
}
