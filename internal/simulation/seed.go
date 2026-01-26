package simulation

import (
	"log/slog"
	"time"

	"github.com/neox5/obsbox/internal/config"
	"github.com/neox5/simv/seed"
)

// InitializeSeed initializes the simv seed registry (required by simv v0.5.0).
// Must be called before creating any simv objects (clocks, sources, values).
func InitializeSeed(cfg *config.SimulationConfig) {
	var masterSeed uint64
	var explicit bool

	if cfg.HasSeed() {
		masterSeed = cfg.GetSeed()
		explicit = true
	} else {
		masterSeed = uint64(time.Now().UnixNano())
		explicit = false
	}

	seed.Init(masterSeed)

	// Log initialization (stream counter will be 0 at startup)
	master, stream := seed.Current()
	slog.Info("seed initialized", "master", master, "stream", stream, "explicit", explicit)
}
