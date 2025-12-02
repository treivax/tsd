// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"sync"
	"time"
)

// EvaluationContext stores intermediate results during alpha chain evaluation.
// It enables decomposition of complex arithmetic expressions into atomic steps
// where each step can reference results from previous steps.
//
// Thread-safe for concurrent access.
type EvaluationContext struct {
	// OriginalFact is the fact being evaluated through the alpha chain
	OriginalFact *Fact

	// IntermediateResults stores computed values indexed by result name (e.g., "temp_1")
	// Values can be numeric results from arithmetic operations or boolean results from comparisons
	IntermediateResults map[string]interface{}

	// EvaluationPath tracks the sequence of steps executed, useful for debugging
	EvaluationPath []string

	// Timestamp records when this context was created
	Timestamp time.Time

	// Metadata stores arbitrary debugging or profiling information
	Metadata map[string]interface{}

	// mutex protects concurrent access to maps and slices
	mutex sync.RWMutex
}

// NewEvaluationContext creates a new evaluation context for the given fact.
// The context is initialized with empty maps and the current timestamp.
func NewEvaluationContext(fact *Fact) *EvaluationContext {
	return &EvaluationContext{
		OriginalFact:        fact,
		IntermediateResults: make(map[string]interface{}),
		EvaluationPath:      make([]string, 0),
		Timestamp:           time.Now(),
		Metadata:            make(map[string]interface{}),
	}
}

// SetIntermediateResult stores a computed value with the given key.
// The key is typically a generated name like "temp_1", "temp_2", etc.
// Also appends the key to the evaluation path for tracing.
//
// Thread-safe.
func (ec *EvaluationContext) SetIntermediateResult(key string, value interface{}) {
	ec.mutex.Lock()
	defer ec.mutex.Unlock()

	ec.IntermediateResults[key] = value
	ec.EvaluationPath = append(ec.EvaluationPath, key)
}

// GetIntermediateResult retrieves a previously stored intermediate result.
// Returns the value and true if found, nil and false if not found.
//
// Thread-safe.
func (ec *EvaluationContext) GetIntermediateResult(key string) (interface{}, bool) {
	ec.mutex.RLock()
	defer ec.mutex.RUnlock()

	value, exists := ec.IntermediateResults[key]
	return value, exists
}

// HasIntermediateResult checks if a result with the given key exists.
//
// Thread-safe.
func (ec *EvaluationContext) HasIntermediateResult(key string) bool {
	ec.mutex.RLock()
	defer ec.mutex.RUnlock()

	_, exists := ec.IntermediateResults[key]
	return exists
}

// Clone creates a deep copy of the evaluation context.
// Useful when branching evaluation paths or for parallel processing.
//
// Thread-safe.
func (ec *EvaluationContext) Clone() *EvaluationContext {
	ec.mutex.RLock()
	defer ec.mutex.RUnlock()

	clone := &EvaluationContext{
		OriginalFact:        ec.OriginalFact,
		IntermediateResults: make(map[string]interface{}, len(ec.IntermediateResults)),
		EvaluationPath:      make([]string, len(ec.EvaluationPath)),
		Timestamp:           ec.Timestamp,
		Metadata:            make(map[string]interface{}, len(ec.Metadata)),
	}

	// Deep copy intermediate results
	for k, v := range ec.IntermediateResults {
		clone.IntermediateResults[k] = v
	}

	// Copy evaluation path
	copy(clone.EvaluationPath, ec.EvaluationPath)

	// Deep copy metadata
	for k, v := range ec.Metadata {
		clone.Metadata[k] = v
	}

	return clone
}

// SetMetadata stores arbitrary metadata for debugging or profiling.
//
// Thread-safe.
func (ec *EvaluationContext) SetMetadata(key string, value interface{}) {
	ec.mutex.Lock()
	defer ec.mutex.Unlock()

	ec.Metadata[key] = value
}

// GetMetadata retrieves metadata by key.
// Returns the value and true if found, nil and false if not found.
//
// Thread-safe.
func (ec *EvaluationContext) GetMetadata(key string) (interface{}, bool) {
	ec.mutex.RLock()
	defer ec.mutex.RUnlock()

	value, exists := ec.Metadata[key]
	return value, exists
}

// GetEvaluationPathString returns a formatted string showing the evaluation path.
// Useful for debugging and logging.
//
// Thread-safe.
func (ec *EvaluationContext) GetEvaluationPathString() string {
	ec.mutex.RLock()
	defer ec.mutex.RUnlock()

	if len(ec.EvaluationPath) == 0 {
		return "(empty path)"
	}

	result := ""
	for i, step := range ec.EvaluationPath {
		if i > 0 {
			result += " → "
		}
		result += step
	}
	return result
}

// String returns a human-readable representation of the context.
// Includes fact ID, number of intermediate results, and evaluation path.
func (ec *EvaluationContext) String() string {
	ec.mutex.RLock()
	defer ec.mutex.RUnlock()

	factID := "(nil)"
	if ec.OriginalFact != nil {
		factID = ec.OriginalFact.ID
	}

	return fmt.Sprintf("EvaluationContext{Fact: %s, Results: %d, Path: %s}",
		factID,
		len(ec.IntermediateResults),
		ec.GetEvaluationPathString())
}

// Reset clears all intermediate results and evaluation path.
// Useful for reusing a context object.
// Does NOT reset the original fact, timestamp, or metadata.
//
// Thread-safe.
func (ec *EvaluationContext) Reset() {
	ec.mutex.Lock()
	defer ec.mutex.Unlock()

	ec.IntermediateResults = make(map[string]interface{})
	ec.EvaluationPath = make([]string, 0)
}

// Size returns the number of intermediate results stored.
//
// Thread-safe.
func (ec *EvaluationContext) Size() int {
	ec.mutex.RLock()
	defer ec.mutex.RUnlock()

	return len(ec.IntermediateResults)
}

// GetAllResults returns a copy of all intermediate results.
// Returns a new map to prevent external modification.
//
// Thread-safe.
func (ec *EvaluationContext) GetAllResults() map[string]interface{} {
	ec.mutex.RLock()
	defer ec.mutex.RUnlock()

	results := make(map[string]interface{}, len(ec.IntermediateResults))
	for k, v := range ec.IntermediateResults {
		results[k] = v
	}
	return results
}
