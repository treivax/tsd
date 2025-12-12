// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete

import (
	"sync"
	"testing"
)

// TestEvaluationContext_NewContext tests context creation
func TestEvaluationContext_NewContext(t *testing.T) {
	fact := &Fact{
		ID:     "test1",
		Type:   "TestType",
		Fields: map[string]interface{}{"value": 42},
	}
	ctx := NewEvaluationContext(fact)
	if ctx == nil {
		t.Fatal("NewEvaluationContext returned nil")
	}
	if ctx.OriginalFact != fact {
		t.Error("OriginalFact not set correctly")
	}
	if ctx.IntermediateResults == nil {
		t.Error("IntermediateResults not initialized")
	}
	if ctx.EvaluationPath == nil {
		t.Error("EvaluationPath not initialized")
	}
	if ctx.Metadata == nil {
		t.Error("Metadata not initialized")
	}
	if ctx.Timestamp.IsZero() {
		t.Error("Timestamp not set")
	}
	if len(ctx.IntermediateResults) != 0 {
		t.Error("IntermediateResults should be empty initially")
	}
	if len(ctx.EvaluationPath) != 0 {
		t.Error("EvaluationPath should be empty initially")
	}
}

// TestEvaluationContext_SetGet tests basic set/get operations
func TestEvaluationContext_SetGet(t *testing.T) {
	fact := &Fact{
		ID:     "test1",
		Type:   "TestType",
		Fields: map[string]interface{}{},
	}
	ctx := NewEvaluationContext(fact)
	// Test setting and getting a numeric value
	ctx.SetIntermediateResult("temp_1", 42.0)
	value, exists := ctx.GetIntermediateResult("temp_1")
	if !exists {
		t.Fatal("temp_1 should exist")
	}
	if value != 42.0 {
		t.Errorf("Expected 42.0, got %v", value)
	}
	// Test setting and getting a boolean value
	ctx.SetIntermediateResult("temp_2", true)
	value, exists = ctx.GetIntermediateResult("temp_2")
	if !exists {
		t.Fatal("temp_2 should exist")
	}
	if value != true {
		t.Errorf("Expected true, got %v", value)
	}
	// Test getting non-existent value
	value, exists = ctx.GetIntermediateResult("nonexistent")
	if exists {
		t.Error("nonexistent should not exist")
	}
	if value != nil {
		t.Error("nonexistent value should be nil")
	}
}

// TestEvaluationContext_HasIntermediateResult tests existence checks
func TestEvaluationContext_HasIntermediateResult(t *testing.T) {
	fact := &Fact{
		ID:     "test1",
		Type:   "TestType",
		Fields: map[string]interface{}{},
	}
	ctx := NewEvaluationContext(fact)
	// Initially should not have any results
	if ctx.HasIntermediateResult("temp_1") {
		t.Error("temp_1 should not exist initially")
	}
	// After setting, should exist
	ctx.SetIntermediateResult("temp_1", 42.0)
	if !ctx.HasIntermediateResult("temp_1") {
		t.Error("temp_1 should exist after setting")
	}
	// Other keys should still not exist
	if ctx.HasIntermediateResult("temp_2") {
		t.Error("temp_2 should not exist")
	}
}

// TestEvaluationContext_EvaluationPath tests path tracking
func TestEvaluationContext_EvaluationPath(t *testing.T) {
	fact := &Fact{
		ID:     "test1",
		Type:   "TestType",
		Fields: map[string]interface{}{},
	}
	ctx := NewEvaluationContext(fact)
	// Set multiple intermediate results
	ctx.SetIntermediateResult("temp_1", 10.0)
	ctx.SetIntermediateResult("temp_2", 20.0)
	ctx.SetIntermediateResult("temp_3", 30.0)
	// Check evaluation path
	if len(ctx.EvaluationPath) != 3 {
		t.Errorf("Expected 3 steps in path, got %d", len(ctx.EvaluationPath))
	}
	expectedPath := []string{"temp_1", "temp_2", "temp_3"}
	for i, step := range ctx.EvaluationPath {
		if step != expectedPath[i] {
			t.Errorf("Path step %d: expected %s, got %s", i, expectedPath[i], step)
		}
	}
	// Test path string formatting
	pathString := ctx.GetEvaluationPathString()
	expected := "temp_1 → temp_2 → temp_3"
	if pathString != expected {
		t.Errorf("Expected path string %q, got %q", expected, pathString)
	}
}

// TestEvaluationContext_Clone tests deep copying
func TestEvaluationContext_Clone(t *testing.T) {
	fact := &Fact{
		ID:     "test1",
		Type:   "TestType",
		Fields: map[string]interface{}{"value": 100},
	}
	ctx := NewEvaluationContext(fact)
	// Populate original context
	ctx.SetIntermediateResult("temp_1", 42.0)
	ctx.SetIntermediateResult("temp_2", true)
	ctx.SetMetadata("key1", "value1")
	// Clone the context
	clone := ctx.Clone()
	// Verify clone has same data
	if clone.OriginalFact != ctx.OriginalFact {
		t.Error("Clone should reference same OriginalFact")
	}
	if clone.Timestamp != ctx.Timestamp {
		t.Error("Clone should have same Timestamp")
	}
	value, exists := clone.GetIntermediateResult("temp_1")
	if !exists || value != 42.0 {
		t.Error("Clone should have temp_1")
	}
	value, exists = clone.GetIntermediateResult("temp_2")
	if !exists || value != true {
		t.Error("Clone should have temp_2")
	}
	metaValue, exists := clone.GetMetadata("key1")
	if !exists || metaValue != "value1" {
		t.Error("Clone should have metadata key1")
	}
	// Verify modifications to clone don't affect original
	clone.SetIntermediateResult("temp_3", 99.0)
	if ctx.HasIntermediateResult("temp_3") {
		t.Error("Original should not have temp_3 after clone modification")
	}
	if !clone.HasIntermediateResult("temp_3") {
		t.Error("Clone should have temp_3")
	}
	// Verify modifications to original don't affect clone
	ctx.SetIntermediateResult("temp_4", 88.0)
	if clone.HasIntermediateResult("temp_4") {
		t.Error("Clone should not have temp_4 after original modification")
	}
	if !ctx.HasIntermediateResult("temp_4") {
		t.Error("Original should have temp_4")
	}
}

// TestEvaluationContext_Metadata tests metadata operations
func TestEvaluationContext_Metadata(t *testing.T) {
	fact := &Fact{
		ID:     "test1",
		Type:   "TestType",
		Fields: map[string]interface{}{},
	}
	ctx := NewEvaluationContext(fact)
	// Set metadata
	ctx.SetMetadata("debug", true)
	ctx.SetMetadata("rule_id", "rule_123")
	ctx.SetMetadata("complexity", 5)
	// Get metadata
	value, exists := ctx.GetMetadata("debug")
	if !exists || value != true {
		t.Error("Metadata 'debug' not set correctly")
	}
	value, exists = ctx.GetMetadata("rule_id")
	if !exists || value != "rule_123" {
		t.Error("Metadata 'rule_id' not set correctly")
	}
	value, exists = ctx.GetMetadata("complexity")
	if !exists || value != 5 {
		t.Error("Metadata 'complexity' not set correctly")
	}
	// Get non-existent metadata
	value, exists = ctx.GetMetadata("nonexistent")
	if exists {
		t.Error("Non-existent metadata should not exist")
	}
	if value != nil {
		t.Error("Non-existent metadata value should be nil")
	}
}

// TestEvaluationContext_Reset tests context reset
func TestEvaluationContext_Reset(t *testing.T) {
	fact := &Fact{
		ID:     "test1",
		Type:   "TestType",
		Fields: map[string]interface{}{},
	}
	ctx := NewEvaluationContext(fact)
	// Populate context
	ctx.SetIntermediateResult("temp_1", 42.0)
	ctx.SetIntermediateResult("temp_2", 84.0)
	ctx.SetMetadata("key", "value")
	originalFact := ctx.OriginalFact
	originalTimestamp := ctx.Timestamp
	// Reset
	ctx.Reset()
	// Verify results and path cleared
	if len(ctx.IntermediateResults) != 0 {
		t.Error("IntermediateResults should be empty after reset")
	}
	if len(ctx.EvaluationPath) != 0 {
		t.Error("EvaluationPath should be empty after reset")
	}
	// Verify fact, timestamp, and metadata NOT cleared
	if ctx.OriginalFact != originalFact {
		t.Error("OriginalFact should not change after reset")
	}
	if ctx.Timestamp != originalTimestamp {
		t.Error("Timestamp should not change after reset")
	}
	metaValue, exists := ctx.GetMetadata("key")
	if !exists || metaValue != "value" {
		t.Error("Metadata should not be cleared by reset")
	}
}

// TestEvaluationContext_Size tests size reporting
func TestEvaluationContext_Size(t *testing.T) {
	fact := &Fact{
		ID:     "test1",
		Type:   "TestType",
		Fields: map[string]interface{}{},
	}
	ctx := NewEvaluationContext(fact)
	if ctx.Size() != 0 {
		t.Error("Initial size should be 0")
	}
	ctx.SetIntermediateResult("temp_1", 42.0)
	if ctx.Size() != 1 {
		t.Errorf("Size should be 1, got %d", ctx.Size())
	}
	ctx.SetIntermediateResult("temp_2", 84.0)
	if ctx.Size() != 2 {
		t.Errorf("Size should be 2, got %d", ctx.Size())
	}
	// Setting same key again should not increase size
	ctx.SetIntermediateResult("temp_1", 99.0)
	if ctx.Size() != 2 {
		t.Errorf("Size should still be 2, got %d", ctx.Size())
	}
}

// TestEvaluationContext_GetAllResults tests getting all results
func TestEvaluationContext_GetAllResults(t *testing.T) {
	fact := &Fact{
		ID:     "test1",
		Type:   "TestType",
		Fields: map[string]interface{}{},
	}
	ctx := NewEvaluationContext(fact)
	ctx.SetIntermediateResult("temp_1", 10.0)
	ctx.SetIntermediateResult("temp_2", 20.0)
	ctx.SetIntermediateResult("temp_3", 30.0)
	allResults := ctx.GetAllResults()
	if len(allResults) != 3 {
		t.Errorf("Expected 3 results, got %d", len(allResults))
	}
	if allResults["temp_1"] != 10.0 {
		t.Error("temp_1 not in results")
	}
	if allResults["temp_2"] != 20.0 {
		t.Error("temp_2 not in results")
	}
	if allResults["temp_3"] != 30.0 {
		t.Error("temp_3 not in results")
	}
	// Verify modifications to returned map don't affect context
	allResults["temp_4"] = 40.0
	if ctx.HasIntermediateResult("temp_4") {
		t.Error("Context should not be affected by modifications to returned map")
	}
}

// TestEvaluationContext_String tests string representation
func TestEvaluationContext_String(t *testing.T) {
	fact := &Fact{
		ID:     "test1",
		Type:   "TestType",
		Fields: map[string]interface{}{},
	}
	ctx := NewEvaluationContext(fact)
	ctx.SetIntermediateResult("temp_1", 42.0)
	ctx.SetIntermediateResult("temp_2", 84.0)
	str := ctx.String()
	// Should contain fact ID
	if len(str) == 0 {
		t.Error("String should not be empty")
	}
	// Basic validation - just check it doesn't panic
	t.Logf("Context string: %s", str)
}

// TestEvaluationContext_EmptyPath tests empty evaluation path
func TestEvaluationContext_EmptyPath(t *testing.T) {
	fact := &Fact{
		ID:     "test1",
		Type:   "TestType",
		Fields: map[string]interface{}{},
	}
	ctx := NewEvaluationContext(fact)
	pathString := ctx.GetEvaluationPathString()
	expected := "(empty path)"
	if pathString != expected {
		t.Errorf("Expected %q for empty path, got %q", expected, pathString)
	}
}

// TestEvaluationContext_Concurrent tests thread-safety
func TestEvaluationContext_Concurrent(t *testing.T) {
	fact := &Fact{
		ID:     "test1",
		Type:   "TestType",
		Fields: map[string]interface{}{},
	}
	ctx := NewEvaluationContext(fact)
	var wg sync.WaitGroup
	numGoroutines := 10
	numOperations := 100
	// Concurrent writers
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := "temp_" + string(rune(id))
				ctx.SetIntermediateResult(key, float64(j))
			}
		}(i)
	}
	// Concurrent readers
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := "temp_" + string(rune(id))
				ctx.GetIntermediateResult(key)
				ctx.HasIntermediateResult(key)
			}
		}(i)
	}
	// Concurrent metadata operations
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := "meta_" + string(rune(id))
				ctx.SetMetadata(key, id)
				ctx.GetMetadata(key)
			}
		}(i)
	}
	// Concurrent readers of size and all results
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				ctx.Size()
				ctx.GetAllResults()
				ctx.GetEvaluationPathString()
			}
		}()
	}
	// Wait for all goroutines to complete
	wg.Wait()
	// If we get here without panicking, thread-safety is working
	t.Log("Concurrent operations completed successfully")
}

// TestEvaluationContext_MultipleTypes tests storing different value types
func TestEvaluationContext_MultipleTypes(t *testing.T) {
	fact := &Fact{
		ID:     "test1",
		Type:   "TestType",
		Fields: map[string]interface{}{},
	}
	ctx := NewEvaluationContext(fact)
	// Store different types
	ctx.SetIntermediateResult("int_value", 42)
	ctx.SetIntermediateResult("float_value", 3.14)
	ctx.SetIntermediateResult("bool_value", true)
	ctx.SetIntermediateResult("string_value", "hello")
	ctx.SetIntermediateResult("nil_value", nil)
	// Retrieve and verify types
	intVal, _ := ctx.GetIntermediateResult("int_value")
	if intVal != 42 {
		t.Errorf("Expected int 42, got %v", intVal)
	}
	floatVal, _ := ctx.GetIntermediateResult("float_value")
	if floatVal != 3.14 {
		t.Errorf("Expected float 3.14, got %v", floatVal)
	}
	boolVal, _ := ctx.GetIntermediateResult("bool_value")
	if boolVal != true {
		t.Errorf("Expected bool true, got %v", boolVal)
	}
	stringVal, _ := ctx.GetIntermediateResult("string_value")
	if stringVal != "hello" {
		t.Errorf("Expected string 'hello', got %v", stringVal)
	}
	nilVal, exists := ctx.GetIntermediateResult("nil_value")
	if !exists {
		t.Error("nil_value should exist")
	}
	if nilVal != nil {
		t.Errorf("Expected nil, got %v", nilVal)
	}
}

// TestEvaluationContext_OverwriteValue tests overwriting existing values
func TestEvaluationContext_OverwriteValue(t *testing.T) {
	fact := &Fact{
		ID:     "test1",
		Type:   "TestType",
		Fields: map[string]interface{}{},
	}
	ctx := NewEvaluationContext(fact)
	// Set initial value
	ctx.SetIntermediateResult("temp_1", 42.0)
	value, _ := ctx.GetIntermediateResult("temp_1")
	if value != 42.0 {
		t.Errorf("Initial value should be 42.0, got %v", value)
	}
	// Overwrite with new value
	ctx.SetIntermediateResult("temp_1", 99.0)
	value, _ = ctx.GetIntermediateResult("temp_1")
	if value != 99.0 {
		t.Errorf("Overwritten value should be 99.0, got %v", value)
	}
	// Check that evaluation path includes both entries
	if len(ctx.EvaluationPath) != 2 {
		t.Errorf("Expected 2 entries in path, got %d", len(ctx.EvaluationPath))
	}
}

// TestEvaluationContext_NilFact tests context with nil fact
func TestEvaluationContext_NilFact(t *testing.T) {
	ctx := NewEvaluationContext(nil)
	if ctx == nil {
		t.Fatal("NewEvaluationContext should not return nil even with nil fact")
	}
	if ctx.OriginalFact != nil {
		t.Error("OriginalFact should be nil")
	}
	// Should still work normally
	ctx.SetIntermediateResult("temp_1", 42.0)
	value, exists := ctx.GetIntermediateResult("temp_1")
	if !exists || value != 42.0 {
		t.Error("Context should work normally even with nil fact")
	}
	// String should handle nil fact gracefully
	str := ctx.String()
	if len(str) == 0 {
		t.Error("String should not be empty even with nil fact")
	}
	t.Logf("String with nil fact: %s", str)
}

// BenchmarkEvaluationContext_SetGet benchmarks set/get operations
func BenchmarkEvaluationContext_SetGet(b *testing.B) {
	fact := &Fact{
		ID:     "test1",
		Type:   "TestType",
		Fields: map[string]interface{}{},
	}
	ctx := NewEvaluationContext(fact)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := "temp_" + string(rune(i%100))
		ctx.SetIntermediateResult(key, float64(i))
		ctx.GetIntermediateResult(key)
	}
}

// BenchmarkEvaluationContext_Clone benchmarks cloning
func BenchmarkEvaluationContext_Clone(b *testing.B) {
	fact := &Fact{
		ID:     "test1",
		Type:   "TestType",
		Fields: map[string]interface{}{},
	}
	ctx := NewEvaluationContext(fact)
	// Populate with some data
	for i := 0; i < 10; i++ {
		ctx.SetIntermediateResult("temp_"+string(rune(i)), float64(i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ctx.Clone()
	}
}

// BenchmarkEvaluationContext_ConcurrentAccess benchmarks concurrent access
func BenchmarkEvaluationContext_ConcurrentAccess(b *testing.B) {
	fact := &Fact{
		ID:     "test1",
		Type:   "TestType",
		Fields: map[string]interface{}{},
	}
	ctx := NewEvaluationContext(fact)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := "temp_" + string(rune(i%10))
			ctx.SetIntermediateResult(key, float64(i))
			ctx.GetIntermediateResult(key)
			i++
		}
	})
}
