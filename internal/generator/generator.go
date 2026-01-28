package generator

import (
	"fmt"

	"github.com/neox5/obsbox/internal/config"
	"github.com/neox5/obsbox/internal/simulation"
	"github.com/neox5/simv/clock"
	"github.com/neox5/simv/source"
)

// Generator manages simv components and value generation.
type Generator struct {
	clocks      []clock.Clock                    // Parallel to values (one per metric)
	values      []*simulation.ValueWrapper       // Index-aligned with metrics
	clockCache  map[string]clock.Clock           // Instance name → shared clock
	sourceCache map[string]source.Publisher[int] // Instance name → shared source
}

// New creates a generator from metric configurations.
// Creates separate clock/source/value instances for each metric.
// Reuses instances when referenced by name via *Ref fields.
func New(metrics []config.MetricConfig) (*Generator, error) {
	g := &Generator{
		clocks:      make([]clock.Clock, len(metrics)),
		values:      make([]*simulation.ValueWrapper, len(metrics)),
		clockCache:  make(map[string]clock.Clock),
		sourceCache: make(map[string]source.Publisher[int]),
	}

	for i, metric := range metrics {
		// Create or reuse clock
		clk, err := g.getOrCreateClock(metric.Value.Source)
		if err != nil {
			return nil, fmt.Errorf("metric %d (%s): failed to create clock: %w",
				i, metric.PrometheusName, err)
		}
		g.clocks[i] = clk

		// Create or reuse source
		src, err := g.getOrCreateSource(metric.Value, clk)
		if err != nil {
			return nil, fmt.Errorf("metric %d (%s): failed to create source: %w",
				i, metric.PrometheusName, err)
		}

		// Create value (always unique - derives from source)
		val, err := simulation.CreateValue(metric.Value, src)
		if err != nil {
			return nil, fmt.Errorf("metric %d (%s): failed to create value: %w",
				i, metric.PrometheusName, err)
		}
		g.values[i] = val
	}

	return g, nil
}

// getOrCreateClock returns cached clock if ClockRef is set, otherwise creates new.
func (g *Generator) getOrCreateClock(sourceCfg config.SourceConfig) (clock.Clock, error) {
	// Check if clock is shared instance
	if sourceCfg.ClockRef != nil {
		instanceName := *sourceCfg.ClockRef

		// Return cached clock if already created
		if clk, exists := g.clockCache[instanceName]; exists {
			return clk, nil
		}

		// Create and cache new clock
		clk, err := simulation.CreateClock(sourceCfg.Clock)
		if err != nil {
			return nil, fmt.Errorf("clock instance %q: %w", instanceName, err)
		}

		g.clockCache[instanceName] = clk
		return clk, nil
	}

	// Unique clock - create new without caching
	return simulation.CreateClock(sourceCfg.Clock)
}

// getOrCreateSource returns cached source if SourceRef is set, otherwise creates new.
func (g *Generator) getOrCreateSource(valueCfg config.ValueConfig, clk clock.Clock) (source.Publisher[int], error) {
	// Check if source is shared instance
	if valueCfg.SourceRef != nil {
		instanceName := *valueCfg.SourceRef

		// Return cached source if already created
		if src, exists := g.sourceCache[instanceName]; exists {
			return src, nil
		}

		// Create and cache new source
		src, err := simulation.CreateSource(valueCfg.Source, clk)
		if err != nil {
			return nil, fmt.Errorf("source instance %q: %w", instanceName, err)
		}

		g.sourceCache[instanceName] = src
		return src, nil
	}

	// Unique source - create new without caching
	return simulation.CreateSource(valueCfg.Source, clk)
}

// Start begins value generation by starting all clocks.
func (g *Generator) Start() {
	// Start each clock (accepts duplicates - simv handles idempotency)
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

// GetValue returns the value at the specified metric index.
func (g *Generator) GetValue(index int) *simulation.ValueWrapper {
	if index < 0 || index >= len(g.values) {
		return nil
	}
	return g.values[index]
}
