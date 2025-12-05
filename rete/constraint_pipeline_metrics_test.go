// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewMetricsCollector tests the creation of a new metrics collector
func TestNewMetricsCollector(t *testing.T) {
	collector := NewMetricsCollector()

	require.NotNil(t, collector, "NewMetricsCollector should return a non-nil collector")
}

// TestMetricsCollector_RecordParsingDuration tests recording parsing duration
func TestMetricsCollector_RecordParsingDuration(t *testing.T) {
	collector := NewMetricsCollector()
	duration := 150 * time.Millisecond

	collector.RecordParsingDuration(duration)
	metrics := collector.Finalize()

	assert.Equal(t, duration, metrics.ParsingDuration, "Parsing duration should be recorded correctly")
}

// TestMetricsCollector_RecordValidationDuration tests recording validation duration
func TestMetricsCollector_RecordValidationDuration(t *testing.T) {
	collector := NewMetricsCollector()
	duration := 75 * time.Millisecond

	collector.RecordValidationDuration(duration)
	metrics := collector.Finalize()

	assert.Equal(t, duration, metrics.ValidationDuration, "Validation duration should be recorded correctly")
}

// TestMetricsCollector_RecordTypeCreationDuration tests recording type creation duration
func TestMetricsCollector_RecordTypeCreationDuration(t *testing.T) {
	collector := NewMetricsCollector()
	duration := 50 * time.Millisecond

	collector.RecordTypeCreationDuration(duration)
	metrics := collector.Finalize()

	assert.Equal(t, duration, metrics.TypeCreationDuration, "Type creation duration should be recorded correctly")
}

// TestMetricsCollector_RecordRuleCreationDuration tests recording rule creation duration
func TestMetricsCollector_RecordRuleCreationDuration(t *testing.T) {
	collector := NewMetricsCollector()
	duration := 200 * time.Millisecond

	collector.RecordRuleCreationDuration(duration)
	metrics := collector.Finalize()

	assert.Equal(t, duration, metrics.RuleCreationDuration, "Rule creation duration should be recorded correctly")
}

// TestMetricsCollector_RecordFactCollectionDuration tests recording fact collection duration
func TestMetricsCollector_RecordFactCollectionDuration(t *testing.T) {
	collector := NewMetricsCollector()
	duration := 100 * time.Millisecond

	collector.RecordFactCollectionDuration(duration)
	metrics := collector.Finalize()

	assert.Equal(t, duration, metrics.FactCollectionDuration, "Fact collection duration should be recorded correctly")
}

// TestMetricsCollector_RecordPropagationDuration tests recording propagation duration
func TestMetricsCollector_RecordPropagationDuration(t *testing.T) {
	collector := NewMetricsCollector()
	duration := 300 * time.Millisecond

	collector.RecordPropagationDuration(duration)
	metrics := collector.Finalize()

	assert.Equal(t, duration, metrics.PropagationDuration, "Propagation duration should be recorded correctly")
}

// TestMetricsCollector_RecordFactSubmissionDuration tests recording fact submission duration
func TestMetricsCollector_RecordFactSubmissionDuration(t *testing.T) {
	collector := NewMetricsCollector()
	duration := 250 * time.Millisecond

	collector.RecordFactSubmissionDuration(duration)
	metrics := collector.Finalize()

	assert.Equal(t, duration, metrics.FactSubmissionDuration, "Fact submission duration should be recorded correctly")
}

// TestMetricsCollector_SetTypesAdded tests setting types added count
func TestMetricsCollector_SetTypesAdded(t *testing.T) {
	collector := NewMetricsCollector()

	collector.SetTypesAdded(5)
	metrics := collector.Finalize()

	assert.Equal(t, 5, metrics.TypesAdded, "Types added should be set correctly")
}

// TestMetricsCollector_SetRulesAdded tests setting rules added count
func TestMetricsCollector_SetRulesAdded(t *testing.T) {
	collector := NewMetricsCollector()

	collector.SetRulesAdded(15)
	metrics := collector.Finalize()

	assert.Equal(t, 15, metrics.RulesAdded, "Rules added should be set correctly")
}

// TestMetricsCollector_SetFactsSubmitted tests setting facts submitted count
func TestMetricsCollector_SetFactsSubmitted(t *testing.T) {
	collector := NewMetricsCollector()

	collector.SetFactsSubmitted(100)
	metrics := collector.Finalize()

	assert.Equal(t, 100, metrics.FactsSubmitted, "Facts submitted should be set correctly")
}

// TestMetricsCollector_SetExistingFactsCollected tests setting existing facts collected count
func TestMetricsCollector_SetExistingFactsCollected(t *testing.T) {
	collector := NewMetricsCollector()

	collector.SetExistingFactsCollected(50)
	metrics := collector.Finalize()

	assert.Equal(t, 50, metrics.ExistingFactsCollected, "Existing facts collected should be set correctly")
}

// TestMetricsCollector_SetFactsPropagated tests setting facts propagated count
func TestMetricsCollector_SetFactsPropagated(t *testing.T) {
	collector := NewMetricsCollector()

	collector.SetFactsPropagated(30)
	metrics := collector.Finalize()

	assert.Equal(t, 30, metrics.FactsPropagated, "Facts propagated should be set correctly")
}

// TestMetricsCollector_SetNewTerminalsAdded tests setting new terminals added count
func TestMetricsCollector_SetNewTerminalsAdded(t *testing.T) {
	collector := NewMetricsCollector()

	collector.SetNewTerminalsAdded(8)
	metrics := collector.Finalize()

	assert.Equal(t, 8, metrics.NewTerminalsAdded, "New terminals added should be set correctly")
}

// TestMetricsCollector_SetPropagationTargets tests setting propagation targets count
func TestMetricsCollector_SetPropagationTargets(t *testing.T) {
	collector := NewMetricsCollector()

	collector.SetPropagationTargets(25)
	metrics := collector.Finalize()

	assert.Equal(t, 25, metrics.PropagationTargets, "Propagation targets should be set correctly")
}

// TestMetricsCollector_SetWasReset tests setting was reset flag
func TestMetricsCollector_SetWasReset(t *testing.T) {
	collector := NewMetricsCollector()

	collector.SetWasReset(true)
	metrics := collector.Finalize()

	assert.True(t, metrics.WasReset, "Was reset should be true after setting")
}

// TestMetricsCollector_SetWasIncremental tests setting was incremental flag
func TestMetricsCollector_SetWasIncremental(t *testing.T) {
	collector := NewMetricsCollector()

	collector.SetWasIncremental(true)
	metrics := collector.Finalize()

	assert.True(t, metrics.WasIncremental, "Was incremental should be true after setting")
}

// TestMetricsCollector_SetValidationSkipped tests setting validation skipped flag
func TestMetricsCollector_SetValidationSkipped(t *testing.T) {
	collector := NewMetricsCollector()

	collector.SetValidationSkipped(true)
	metrics := collector.Finalize()

	assert.True(t, metrics.ValidationSkipped, "Validation skipped should be true after setting")
}

// TestMetricsCollector_Finalize tests the finalization of metrics
func TestMetricsCollector_Finalize(t *testing.T) {
	collector := NewMetricsCollector()

	// Set some metrics
	collector.RecordParsingDuration(100 * time.Millisecond)
	collector.RecordValidationDuration(50 * time.Millisecond)
	collector.RecordTypeCreationDuration(75 * time.Millisecond)
	collector.RecordRuleCreationDuration(150 * time.Millisecond)
	collector.RecordFactCollectionDuration(80 * time.Millisecond)
	collector.RecordPropagationDuration(200 * time.Millisecond)
	collector.RecordFactSubmissionDuration(120 * time.Millisecond)

	collector.SetTypesAdded(3)
	collector.SetRulesAdded(7)
	collector.SetFactsSubmitted(50)
	collector.SetExistingFactsCollected(40)
	collector.SetFactsPropagated(30)
	collector.SetNewTerminalsAdded(5)
	collector.SetPropagationTargets(20)
	collector.SetWasReset(false)
	collector.SetWasIncremental(true)
	collector.SetValidationSkipped(false)

	// Wait a bit to ensure elapsed time is non-zero
	time.Sleep(10 * time.Millisecond)

	// Finalize
	metrics := collector.Finalize()

	// Verify all values are preserved
	require.NotNil(t, metrics, "Finalize should return non-nil metrics")
	assert.Equal(t, 100*time.Millisecond, metrics.ParsingDuration)
	assert.Equal(t, 50*time.Millisecond, metrics.ValidationDuration)
	assert.Equal(t, 75*time.Millisecond, metrics.TypeCreationDuration)
	assert.Equal(t, 150*time.Millisecond, metrics.RuleCreationDuration)
	assert.Equal(t, 80*time.Millisecond, metrics.FactCollectionDuration)
	assert.Equal(t, 200*time.Millisecond, metrics.PropagationDuration)
	assert.Equal(t, 120*time.Millisecond, metrics.FactSubmissionDuration)

	assert.Equal(t, 3, metrics.TypesAdded)
	assert.Equal(t, 7, metrics.RulesAdded)
	assert.Equal(t, 50, metrics.FactsSubmitted)
	assert.Equal(t, 40, metrics.ExistingFactsCollected)
	assert.Equal(t, 30, metrics.FactsPropagated)
	assert.Equal(t, 5, metrics.NewTerminalsAdded)
	assert.Equal(t, 20, metrics.PropagationTargets)

	assert.False(t, metrics.WasReset)
	assert.True(t, metrics.WasIncremental)
	assert.False(t, metrics.ValidationSkipped)

	// Verify total duration is calculated
	assert.Greater(t, metrics.TotalDuration, 10*time.Millisecond, "Total duration should include elapsed time")
}

// TestMetricsCollector_FinalizeWithZeroValues tests finalization with default values
func TestMetricsCollector_FinalizeWithZeroValues(t *testing.T) {
	collector := NewMetricsCollector()

	// Finalize immediately without setting any values
	metrics := collector.Finalize()

	require.NotNil(t, metrics, "Finalize should return non-nil metrics even with zero values")
	assert.Equal(t, time.Duration(0), metrics.ParsingDuration)
	assert.Equal(t, time.Duration(0), metrics.ValidationDuration)
	assert.Equal(t, time.Duration(0), metrics.TypeCreationDuration)
	assert.Equal(t, time.Duration(0), metrics.RuleCreationDuration)
	assert.Equal(t, time.Duration(0), metrics.FactCollectionDuration)
	assert.Equal(t, time.Duration(0), metrics.PropagationDuration)
	assert.Equal(t, time.Duration(0), metrics.FactSubmissionDuration)

	assert.Equal(t, 0, metrics.TypesAdded)
	assert.Equal(t, 0, metrics.RulesAdded)
	assert.Equal(t, 0, metrics.FactsSubmitted)
	assert.Equal(t, 0, metrics.ExistingFactsCollected)
	assert.Equal(t, 0, metrics.FactsPropagated)
	assert.Equal(t, 0, metrics.NewTerminalsAdded)
	assert.Equal(t, 0, metrics.PropagationTargets)

	assert.False(t, metrics.WasReset)
	assert.False(t, metrics.WasIncremental)
	assert.False(t, metrics.ValidationSkipped)

	// Total duration should still be non-zero (time elapsed since creation)
	assert.GreaterOrEqual(t, metrics.TotalDuration, time.Duration(0), "Total duration should be non-negative")
}

// TestMetricsCollector_MultipleRecordings tests updating metrics multiple times
func TestMetricsCollector_MultipleRecordings(t *testing.T) {
	collector := NewMetricsCollector()

	// Record multiple times (last value should be kept)
	collector.RecordParsingDuration(100 * time.Millisecond)
	collector.RecordParsingDuration(200 * time.Millisecond)
	collector.RecordParsingDuration(150 * time.Millisecond)

	collector.SetTypesAdded(5)
	collector.SetTypesAdded(10)
	collector.SetTypesAdded(7)

	metrics := collector.Finalize()

	assert.Equal(t, 150*time.Millisecond, metrics.ParsingDuration, "Last recorded parsing duration should be used")
	assert.Equal(t, 7, metrics.TypesAdded, "Last set types added should be used")
}

// TestMetricsCollector_LargeValues tests handling large metric values
func TestMetricsCollector_LargeValues(t *testing.T) {
	collector := NewMetricsCollector()

	// Set large values
	collector.RecordParsingDuration(10 * time.Minute)
	collector.SetFactsSubmitted(1000000)
	collector.SetPropagationTargets(50000)

	metrics := collector.Finalize()

	assert.Equal(t, 10*time.Minute, metrics.ParsingDuration)
	assert.Equal(t, 1000000, metrics.FactsSubmitted)
	assert.Equal(t, 50000, metrics.PropagationTargets)
}

// TestMetricsCollector_CombinedScenario tests a realistic combined scenario
func TestMetricsCollector_CombinedScenario(t *testing.T) {
	collector := NewMetricsCollector()

	// Simulate an incremental ingestion with no reset
	collector.SetWasIncremental(true)
	collector.SetWasReset(false)
	collector.SetValidationSkipped(false)

	// Record durations for each phase
	collector.RecordParsingDuration(120 * time.Millisecond)
	collector.RecordValidationDuration(80 * time.Millisecond)
	collector.RecordTypeCreationDuration(50 * time.Millisecond)
	collector.RecordRuleCreationDuration(180 * time.Millisecond)
	collector.RecordFactCollectionDuration(90 * time.Millisecond)
	collector.RecordPropagationDuration(250 * time.Millisecond)
	collector.RecordFactSubmissionDuration(140 * time.Millisecond)

	// Set counts
	collector.SetTypesAdded(2)
	collector.SetRulesAdded(5)
	collector.SetExistingFactsCollected(100)
	collector.SetFactsSubmitted(150)
	collector.SetNewTerminalsAdded(3)
	collector.SetPropagationTargets(30)
	collector.SetFactsPropagated(75)

	// Wait to ensure total duration captures elapsed time
	time.Sleep(5 * time.Millisecond)

	metrics := collector.Finalize()

	// Verify flags
	assert.True(t, metrics.WasIncremental, "Should be incremental")
	assert.False(t, metrics.WasReset, "Should not be reset")
	assert.False(t, metrics.ValidationSkipped, "Validation should not be skipped")

	// Verify all durations are present
	assert.Greater(t, metrics.ParsingDuration, time.Duration(0))
	assert.Greater(t, metrics.ValidationDuration, time.Duration(0))
	assert.Greater(t, metrics.TypeCreationDuration, time.Duration(0))
	assert.Greater(t, metrics.RuleCreationDuration, time.Duration(0))
	assert.Greater(t, metrics.FactCollectionDuration, time.Duration(0))
	assert.Greater(t, metrics.PropagationDuration, time.Duration(0))
	assert.Greater(t, metrics.FactSubmissionDuration, time.Duration(0))

	// Verify counts make sense
	assert.Equal(t, 2, metrics.TypesAdded)
	assert.Equal(t, 5, metrics.RulesAdded)
	assert.Equal(t, 100, metrics.ExistingFactsCollected)
	assert.Equal(t, 150, metrics.FactsSubmitted)
	assert.Equal(t, 3, metrics.NewTerminalsAdded)
	assert.Equal(t, 30, metrics.PropagationTargets)
	assert.Equal(t, 75, metrics.FactsPropagated)

	// Total duration should be at least the sleep time plus processing
	assert.Greater(t, metrics.TotalDuration, 5*time.Millisecond)
}

// TestMetricsCollector_AllSetters tests all setter methods together
func TestMetricsCollector_AllSetters(t *testing.T) {
	collector := NewMetricsCollector()

	// Set all possible metrics
	collector.RecordParsingDuration(10 * time.Millisecond)
	collector.RecordValidationDuration(20 * time.Millisecond)
	collector.RecordTypeCreationDuration(30 * time.Millisecond)
	collector.RecordRuleCreationDuration(40 * time.Millisecond)
	collector.RecordFactCollectionDuration(50 * time.Millisecond)
	collector.RecordPropagationDuration(60 * time.Millisecond)
	collector.RecordFactSubmissionDuration(70 * time.Millisecond)

	collector.SetTypesAdded(1)
	collector.SetRulesAdded(2)
	collector.SetFactsSubmitted(3)
	collector.SetExistingFactsCollected(4)
	collector.SetFactsPropagated(5)
	collector.SetNewTerminalsAdded(6)
	collector.SetPropagationTargets(7)

	collector.SetWasReset(true)
	collector.SetWasIncremental(true)
	collector.SetValidationSkipped(true)

	metrics := collector.Finalize()

	// Verify all values
	assert.Equal(t, 10*time.Millisecond, metrics.ParsingDuration)
	assert.Equal(t, 20*time.Millisecond, metrics.ValidationDuration)
	assert.Equal(t, 30*time.Millisecond, metrics.TypeCreationDuration)
	assert.Equal(t, 40*time.Millisecond, metrics.RuleCreationDuration)
	assert.Equal(t, 50*time.Millisecond, metrics.FactCollectionDuration)
	assert.Equal(t, 60*time.Millisecond, metrics.PropagationDuration)
	assert.Equal(t, 70*time.Millisecond, metrics.FactSubmissionDuration)

	assert.Equal(t, 1, metrics.TypesAdded)
	assert.Equal(t, 2, metrics.RulesAdded)
	assert.Equal(t, 3, metrics.FactsSubmitted)
	assert.Equal(t, 4, metrics.ExistingFactsCollected)
	assert.Equal(t, 5, metrics.FactsPropagated)
	assert.Equal(t, 6, metrics.NewTerminalsAdded)
	assert.Equal(t, 7, metrics.PropagationTargets)

	assert.True(t, metrics.WasReset)
	assert.True(t, metrics.WasIncremental)
	assert.True(t, metrics.ValidationSkipped)
}

// TestMetricsCollector_RecordNetworkState tests recording network state
func TestMetricsCollector_RecordNetworkState(t *testing.T) {
	collector := NewMetricsCollector()
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)

	// Add some nodes to the network
	network.TypeNodes["Person"] = NewTypeNode("Person", TypeDefinition{Name: "Person"}, storage)
	network.TypeNodes["Order"] = NewTypeNode("Order", TypeDefinition{Name: "Order"}, storage)

	network.TerminalNodes["term1"] = NewTerminalNode("term1", nil, storage)
	network.TerminalNodes["term2"] = NewTerminalNode("term2", nil, storage)
	network.TerminalNodes["term3"] = NewTerminalNode("term3", nil, storage)

	alphaNode := NewAlphaNode("alpha1", nil, "var", storage)
	network.AlphaNodes["alpha1"] = alphaNode

	joinNode := NewJoinNode("join1", nil, []string{}, []string{}, map[string]string{}, storage)
	network.BetaNodes["beta1"] = joinNode
	network.BetaNodes["beta2"] = joinNode // Add duplicate to test count

	// Record network state
	collector.RecordNetworkState(network)
	metrics := collector.Finalize()

	// Verify counts
	assert.Equal(t, 2, metrics.TotalTypeNodes, "Should have 2 type nodes")
	assert.Equal(t, 3, metrics.TotalTerminalNodes, "Should have 3 terminal nodes")
	assert.Equal(t, 1, metrics.TotalAlphaNodes, "Should have 1 alpha node")
	assert.Equal(t, 2, metrics.TotalBetaNodes, "Should have 2 beta nodes")
}

// TestMetricsCollector_GetMetrics tests getting a copy of current metrics
func TestMetricsCollector_GetMetrics(t *testing.T) {
	collector := NewMetricsCollector()

	collector.SetTypesAdded(5)
	collector.SetRulesAdded(10)
	collector.RecordParsingDuration(100 * time.Millisecond)

	// Get metrics before finalization
	metrics := collector.GetMetrics()

	require.NotNil(t, metrics, "GetMetrics should return non-nil")
	assert.Equal(t, 5, metrics.TypesAdded)
	assert.Equal(t, 10, metrics.RulesAdded)
	assert.Equal(t, 100*time.Millisecond, metrics.ParsingDuration)

	// Modify returned metrics - should not affect collector
	metrics.TypesAdded = 999

	// Get again and verify original value preserved
	metrics2 := collector.GetMetrics()
	assert.Equal(t, 5, metrics2.TypesAdded, "Original value should be preserved")
}

// TestIngestionMetrics_String tests the String method
func TestIngestionMetrics_String(t *testing.T) {
	metrics := &IngestionMetrics{
		ParsingDuration:        100 * time.Millisecond,
		ValidationDuration:     50 * time.Millisecond,
		TypeCreationDuration:   75 * time.Millisecond,
		RuleCreationDuration:   150 * time.Millisecond,
		FactCollectionDuration: 80 * time.Millisecond,
		PropagationDuration:    200 * time.Millisecond,
		FactSubmissionDuration: 120 * time.Millisecond,
		TotalDuration:          775 * time.Millisecond,
		TypesAdded:             3,
		RulesAdded:             7,
		FactsSubmitted:         50,
		ExistingFactsCollected: 40,
		FactsPropagated:        30,
		NewTerminalsAdded:      5,
		PropagationTargets:     20,
		TotalTypeNodes:         10,
		TotalTerminalNodes:     15,
		TotalAlphaNodes:        25,
		TotalBetaNodes:         12,
		WasReset:               false,
		WasIncremental:         true,
		ValidationSkipped:      false,
	}

	str := metrics.String()

	// Verify string contains key information
	assert.Contains(t, str, "100ms", "Should contain parsing duration")
	assert.Contains(t, str, "775ms", "Should contain total duration")
	assert.Contains(t, str, "Types ajoutés", "Should contain types label")
	assert.Contains(t, str, "3", "Should contain types count")
	assert.Contains(t, str, "7", "Should contain rules count")
	assert.Contains(t, str, "50", "Should contain facts submitted")
	assert.Contains(t, str, "true", "Should contain incremental flag")
	assert.NotEmpty(t, str, "String should not be empty")
}

// TestIngestionMetrics_Summary tests the Summary method
func TestIngestionMetrics_Summary(t *testing.T) {
	metrics := &IngestionMetrics{
		TotalDuration:     500 * time.Millisecond,
		TypesAdded:        2,
		RulesAdded:        5,
		FactsSubmitted:    100,
		FactsPropagated:   75,
		NewTerminalsAdded: 3,
	}

	summary := metrics.Summary()

	// Verify summary contains key information
	assert.Contains(t, summary, "500ms", "Should contain total duration")
	assert.Contains(t, summary, "2", "Should contain types count")
	assert.Contains(t, summary, "5", "Should contain rules count")
	assert.Contains(t, summary, "100", "Should contain facts submitted")
	assert.Contains(t, summary, "75", "Should contain facts propagated")
	assert.Contains(t, summary, "3", "Should contain new terminals")
	assert.NotEmpty(t, summary, "Summary should not be empty")
}

// TestIngestionMetrics_IsEfficient tests the IsEfficient method
func TestIngestionMetrics_IsEfficient(t *testing.T) {
	tests := []struct {
		name              string
		metrics           *IngestionMetrics
		expectedEfficient bool
		description       string
	}{
		{
			name: "efficient with targeted propagation",
			metrics: &IngestionMetrics{
				NewTerminalsAdded:      2,
				ExistingFactsCollected: 10,
				FactsPropagated:        15, // Less than 10 * 2 = 20
				PropagationDuration:    50 * time.Millisecond,
				TotalDuration:          500 * time.Millisecond, // 10% ratio
			},
			expectedEfficient: true,
			description:       "Should be efficient with targeted propagation and good ratio",
		},
		{
			name: "inefficient with too many propagations",
			metrics: &IngestionMetrics{
				NewTerminalsAdded:      2,
				ExistingFactsCollected: 10,
				FactsPropagated:        25, // More than 10 * 2 = 20
				PropagationDuration:    50 * time.Millisecond,
				TotalDuration:          500 * time.Millisecond,
			},
			expectedEfficient: false,
			description:       "Should be inefficient when too many propagations occur",
		},
		{
			name: "inefficient with high propagation ratio",
			metrics: &IngestionMetrics{
				NewTerminalsAdded:      2,
				ExistingFactsCollected: 10,
				FactsPropagated:        15,
				PropagationDuration:    200 * time.Millisecond, // 40% of total
				TotalDuration:          500 * time.Millisecond,
			},
			expectedEfficient: false,
			description:       "Should be inefficient when propagation takes too much time",
		},
		{
			name: "efficient with no new terminals",
			metrics: &IngestionMetrics{
				NewTerminalsAdded:      0,
				ExistingFactsCollected: 10,
				FactsPropagated:        0,
				PropagationDuration:    0,
				TotalDuration:          100 * time.Millisecond,
			},
			expectedEfficient: true,
			description:       "Should be efficient with no propagation needed",
		},
		{
			name: "efficient with small duration",
			metrics: &IngestionMetrics{
				NewTerminalsAdded:      1,
				ExistingFactsCollected: 5,
				FactsPropagated:        4,
				PropagationDuration:    500 * time.Microsecond, // < 1ms
				TotalDuration:          1 * time.Millisecond,
			},
			expectedEfficient: true,
			description:       "Should ignore ratio for very small durations",
		},
		{
			name: "efficient with no existing facts",
			metrics: &IngestionMetrics{
				NewTerminalsAdded:      5,
				ExistingFactsCollected: 0,
				FactsPropagated:        0,
				PropagationDuration:    0,
				TotalDuration:          100 * time.Millisecond,
			},
			expectedEfficient: true,
			description:       "Should be efficient when no existing facts to propagate",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.metrics.IsEfficient()
			assert.Equal(t, tt.expectedEfficient, result, tt.description)
		})
	}
}

// TestIngestionMetrics_GetBottleneck tests the GetBottleneck method
func TestIngestionMetrics_GetBottleneck(t *testing.T) {
	tests := []struct {
		name                  string
		metrics               *IngestionMetrics
		expectedBottleneckKey string
		description           string
	}{
		{
			name: "parsing is bottleneck",
			metrics: &IngestionMetrics{
				ParsingDuration:        500 * time.Millisecond,
				ValidationDuration:     50 * time.Millisecond,
				TypeCreationDuration:   30 * time.Millisecond,
				RuleCreationDuration:   40 * time.Millisecond,
				FactCollectionDuration: 20 * time.Millisecond,
				PropagationDuration:    60 * time.Millisecond,
				FactSubmissionDuration: 30 * time.Millisecond,
				TotalDuration:          730 * time.Millisecond,
			},
			expectedBottleneckKey: "Parsing",
			description:           "Parsing should be identified as bottleneck",
		},
		{
			name: "propagation is bottleneck",
			metrics: &IngestionMetrics{
				ParsingDuration:        50 * time.Millisecond,
				ValidationDuration:     30 * time.Millisecond,
				TypeCreationDuration:   20 * time.Millisecond,
				RuleCreationDuration:   40 * time.Millisecond,
				FactCollectionDuration: 25 * time.Millisecond,
				PropagationDuration:    800 * time.Millisecond,
				FactSubmissionDuration: 35 * time.Millisecond,
				TotalDuration:          1000 * time.Millisecond,
			},
			expectedBottleneckKey: "Propagation",
			description:           "Propagation should be identified as bottleneck",
		},
		{
			name: "rule creation is bottleneck",
			metrics: &IngestionMetrics{
				ParsingDuration:        30 * time.Millisecond,
				ValidationDuration:     20 * time.Millisecond,
				TypeCreationDuration:   15 * time.Millisecond,
				RuleCreationDuration:   600 * time.Millisecond,
				FactCollectionDuration: 25 * time.Millisecond,
				PropagationDuration:    50 * time.Millisecond,
				FactSubmissionDuration: 40 * time.Millisecond,
				TotalDuration:          780 * time.Millisecond,
			},
			expectedBottleneckKey: "Création règles",
			description:           "Rule creation should be identified as bottleneck",
		},
		{
			name: "fact submission is bottleneck",
			metrics: &IngestionMetrics{
				ParsingDuration:        20 * time.Millisecond,
				ValidationDuration:     15 * time.Millisecond,
				TypeCreationDuration:   10 * time.Millisecond,
				RuleCreationDuration:   30 * time.Millisecond,
				FactCollectionDuration: 25 * time.Millisecond,
				PropagationDuration:    40 * time.Millisecond,
				FactSubmissionDuration: 700 * time.Millisecond,
				TotalDuration:          840 * time.Millisecond,
			},
			expectedBottleneckKey: "Soumission faits",
			description:           "Fact submission should be identified as bottleneck",
		},
		{
			name: "all zero durations",
			metrics: &IngestionMetrics{
				ParsingDuration:        0,
				ValidationDuration:     0,
				TypeCreationDuration:   0,
				RuleCreationDuration:   0,
				FactCollectionDuration: 0,
				PropagationDuration:    0,
				FactSubmissionDuration: 0,
				TotalDuration:          0,
			},
			expectedBottleneckKey: "Aucun goulot",
			description:           "Should return no bottleneck message when all durations are zero",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bottleneck := tt.metrics.GetBottleneck()
			assert.Contains(t, bottleneck, tt.expectedBottleneckKey, tt.description)
		})
	}
}
