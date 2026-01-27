package generator

import (
	"fmt"

	"github.com/neox5/obsbox/internal/config"
	"github.com/neox5/obsbox/internal/simulation"
	"github.com/neox5/simv/clock"
)

// Generator manages simv components and value generation.
type Generator struct {
	clocks []clock.Clock              // Parallel to values (one per metric)
	values []*simulation.ValueWrapper // Index-aligned with metrics
}

// New creates a generator from metric configurations.
// Creates separate clock/source/value instances for each metric.
func New(metrics []config.MetricConfig) (*Generator, error) {
	clocks := make([]clock.Clock, len(metrics))
	values := make([]*simulation.ValueWrapper, len(metrics))

	for i, metric := range metrics {
		// Determine effective clock (value override or source clock)
		effectiveClock := simulation.GetEffectiveClock(metric.Value)

		// Create clock for this metric
		clk, err := simulation.CreateClock(effectiveClock)
		if err != nil {
			return nil, fmt.Errorf("metric %d (%s): failed to create clock: %w",
				i, metric.PrometheusName, err)
		}
		clocks[i] = clk

		// Create source for this metric (using effective clock)
		src, err := simulation.CreateSource(metric.Value.Source, clk)
		if err != nil {
			return nil, fmt.Errorf("metric %d (%s): failed to create source: %w",
				i, metric.PrometheusName, err)
		}

		// Create value for this metric
		val, err := simulation.CreateValue(metric.Value, src)
		if err != nil {
			return nil, fmt.Errorf("metric %d (%s): failed to create value: %w",
				i, metric.PrometheusName, err)
		}

		values[i] = val
	}

	return &Generator{
		clocks: clocks,
		values: values,
	}, nil
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
