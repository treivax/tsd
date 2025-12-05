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
