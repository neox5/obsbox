package generator

import (
	"fmt"

	"github.com/neox5/obsbox/internal/config"
	"github.com/neox5/obsbox/internal/simulation"
	"github.com/neox5/simv/clock"
	"github.com/neox5/simv/source"
	"github.com/neox5/simv/value"
)

// Generator manages simv components and value generation.
type Generator struct {
	clocks  map[string]clock.Clock
	sources map[string]source.Publisher[int]
	values  map[string]*value.Value[int]
}

// New creates a generator from configuration.
// Values are created and started during initialization.
func New(cfg *config.Config) (*Generator, error) {
	gen := &Generator{
		clocks:  make(map[string]clock.Clock),
		sources: make(map[string]source.Publisher[int]),
		values:  make(map[string]*value.Value[int]),
	}

	// Create clocks
	for name, clockCfg := range cfg.Simulation.Clocks {
		clk, err := simulation.CreateClock(clockCfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create clock %q: %w", name, err)
		}
		gen.clocks[name] = clk
	}

	// Create sources
	for name, srcCfg := range cfg.Simulation.Sources {
		clk, exists := gen.clocks[srcCfg.Clock]
		if !exists {
			return nil, fmt.Errorf("clock %q not found for source %q", srcCfg.Clock, name)
		}

		src, err := simulation.CreateSource(srcCfg, clk)
		if err != nil {
			return nil, fmt.Errorf("failed to create source %q: %w", name, err)
		}
		gen.sources[name] = src
	}

	// Create values
	for name, valCfg := range cfg.Simulation.Values {
		src, exists := gen.sources[valCfg.Source]
		if !exists {
			return nil, fmt.Errorf("source %q not found for value %q", valCfg.Source, name)
		}

		val, err := simulation.CreateValue(valCfg, src)
		if err != nil {
			return nil, fmt.Errorf("failed to create value %q: %w", name, err)
		}

		gen.values[name] = val
	}

	return gen, nil
}

// Start begins value generation by starting clocks.
// Values are already started during creation.
func (g *Generator) Start() {
	for _, clk := range g.clocks {
		clk.Start()
	}
}

// Stop halts value generation and releases resources.
func (g *Generator) Stop() {
	// Stop clocks first (stops generating new values)
	for _, clk := range g.clocks {
		clk.Stop()
	}

	// Stop values (cleanup goroutines and channels)
	for _, val := range g.values {
		val.Stop()
	}
}

// GetValue returns a named value.
func (g *Generator) GetValue(name string) (*value.Value[int], bool) {
	val, exists := g.values[name]
	return val, exists
}
