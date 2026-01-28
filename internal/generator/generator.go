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
	// Lifecycle management - unique objects only
	clocks  []clock.Clock
	sources []source.Publisher[int]
	values  []*simulation.ValueWrapper

	// Instance sharing - named references
	clockInstances  map[string]clock.Clock
	sourceInstances map[string]source.Publisher[int]
	valueInstances  map[string]*simulation.ValueWrapper

	// Metric indexing - fast lookup by metric index
	metricValues []*simulation.ValueWrapper
}

// New creates a generator from metric configurations.
// Creates separate clock/source/value instances for each metric.
// Reuses instances when referenced by name via *Ref fields.
func New(metrics []config.MetricConfig) (*Generator, error) {
	g := &Generator{
		clockInstances:  make(map[string]clock.Clock),
		sourceInstances: make(map[string]source.Publisher[int]),
		valueInstances:  make(map[string]*simulation.ValueWrapper),
		metricValues:    make([]*simulation.ValueWrapper, len(metrics)),
	}

	for i, metric := range metrics {
		// Get or create clock
		clk, err := g.getOrCreateClock(metric.Value.Source)
		if err != nil {
			return nil, fmt.Errorf("metric %d (%s): failed to create clock: %w",
				i, metric.PrometheusName, err)
		}

		// Get or create source
		src, err := g.getOrCreateSource(metric.Value, clk)
		if err != nil {
			return nil, fmt.Errorf("metric %d (%s): failed to create source: %w",
				i, metric.PrometheusName, err)
		}

		// Get or create value
		val, err := g.getOrCreateValue(metric.Value, src)
		if err != nil {
			return nil, fmt.Errorf("metric %d (%s): failed to create value: %w",
				i, metric.PrometheusName, err)
		}

		// Store for metric lookup (allows duplicates)
		g.metricValues[i] = val
	}

	return g, nil
}

// getOrCreateClock returns cached clock if ClockRef is set, otherwise creates new.
// Adds unique clocks to lifecycle management.
func (g *Generator) getOrCreateClock(sourceCfg config.SourceConfig) (clock.Clock, error) {
	// Check if clock is shared instance
	if sourceCfg.ClockRef != nil {
		instanceName := *sourceCfg.ClockRef

		// Return cached clock if already created
		if clk, exists := g.clockInstances[instanceName]; exists {
			return clk, nil
		}

		// Create new clock
		clk, err := simulation.CreateClock(sourceCfg.Clock)
		if err != nil {
			return nil, fmt.Errorf("clock instance %q: %w", instanceName, err)
		}

		// Cache for sharing
		g.clockInstances[instanceName] = clk

		// Add to lifecycle management
		g.clocks = append(g.clocks, clk)

		return clk, nil
	}

	// Unique clock - create new without caching
	clk, err := simulation.CreateClock(sourceCfg.Clock)
	if err != nil {
		return nil, err
	}

	// Add to lifecycle management
	g.clocks = append(g.clocks, clk)

	return clk, nil
}

// getOrCreateSource returns cached source if SourceRef is set, otherwise creates new.
// Adds unique sources to lifecycle management.
func (g *Generator) getOrCreateSource(valueCfg config.ValueConfig, clk clock.Clock) (source.Publisher[int], error) {
	// Check if source is shared instance
	if valueCfg.SourceRef != nil {
		instanceName := *valueCfg.SourceRef

		// Return cached source if already created
		if src, exists := g.sourceInstances[instanceName]; exists {
			return src, nil
		}

		// Create new source
		src, err := simulation.CreateSource(valueCfg.Source, clk)
		if err != nil {
			return nil, fmt.Errorf("source instance %q: %w", instanceName, err)
		}

		// Cache for sharing
		g.sourceInstances[instanceName] = src

		// Add to lifecycle management
		g.sources = append(g.sources, src)

		return src, nil
	}

	// Unique source - create new without caching
	src, err := simulation.CreateSource(valueCfg.Source, clk)
	if err != nil {
		return nil, err
	}

	// Add to lifecycle management
	g.sources = append(g.sources, src)

	return src, nil
}

// getOrCreateValue creates or returns cached value.
// Values are always added to lifecycle management.
func (g *Generator) getOrCreateValue(valueCfg config.ValueConfig, src source.Publisher[int]) (*simulation.ValueWrapper, error) {
	// Note: Value instance sharing not yet implemented in config
	// This structure supports future value instance sharing

	// Create value
	val, err := simulation.CreateValue(valueCfg, src)
	if err != nil {
		return nil, err
	}

	// Add to lifecycle management
	g.values = append(g.values, val)

	return val, nil
}

// Start begins value generation by starting all unique clocks.
func (g *Generator) Start() {
	// Start each unique clock exactly once
	for _, clk := range g.clocks {
		clk.Start()
	}
}

// Stop halts value generation and releases resources.
func (g *Generator) Stop() {
	// Stop unique clocks
	for _, clk := range g.clocks {
		clk.Stop()
	}

	// Stop unique values
	for _, val := range g.values {
		val.Stop()
	}
}

// GetValue returns the value at the specified metric index.
func (g *Generator) GetValue(index int) *simulation.ValueWrapper {
	if index < 0 || index >= len(g.metricValues) {
		return nil
	}
	return g.metricValues[index]
}
