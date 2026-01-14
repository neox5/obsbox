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
	values  map[string]value.Value[int]
}

// New creates a generator from configuration.
func New(cfg *config.Config) (*Generator, error) {
	gen := &Generator{
		clocks:  make(map[string]clock.Clock),
		sources: make(map[string]source.Publisher[int]),
		values:  make(map[string]value.Value[int]),
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

	// Create values (handles dependencies via multiple passes)
	if err := gen.createValues(cfg.Simulation.Values); err != nil {
		return nil, err
	}

	return gen, nil
}

// createValues creates all values, resolving clone dependencies.
func (g *Generator) createValues(valueCfgs map[string]config.ValueConfig) error {
	created := make(map[string]bool)
	pending := make(map[string]config.ValueConfig)

	// Copy all configs to pending
	for name, cfg := range valueCfgs {
		pending[name] = cfg
	}

	// Process until all values created or deadlock detected
	for len(pending) > 0 {
		progress := false

		for name, valCfg := range pending {
			// Check if dependencies are satisfied
			if valCfg.Clone != "" {
				if !created[valCfg.Clone] {
					continue // Wait for dependency
				}
			}

			// Get source or base value
			var src source.Publisher[int]
			var baseValue value.Value[int]

			if valCfg.Source != "" {
				var exists bool
				src, exists = g.sources[valCfg.Source]
				if !exists {
					return fmt.Errorf("source %q not found for value %q", valCfg.Source, name)
				}
			}

			if valCfg.Clone != "" {
				var exists bool
				baseValue, exists = g.values[valCfg.Clone]
				if !exists {
					return fmt.Errorf("clone base %q not found for value %q", valCfg.Clone, name)
				}
			}

			// Create value
			val, err := simulation.CreateValue(valCfg, src, baseValue)
			if err != nil {
				return fmt.Errorf("failed to create value %q: %w", name, err)
			}

			g.values[name] = val
			created[name] = true
			delete(pending, name)
			progress = true
		}

		if !progress {
			return fmt.Errorf("circular dependency detected in values")
		}
	}

	return nil
}

// Start begins value generation.
func (g *Generator) Start() {
	for _, clk := range g.clocks {
		clk.Start()
	}
}

// Stop halts value generation.
func (g *Generator) Stop() {
	for _, clk := range g.clocks {
		clk.Stop()
	}
}

// GetValue returns a named value.
func (g *Generator) GetValue(name string) (value.Value[int], bool) {
	val, exists := g.values[name]
	return val, exists
}
