// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"testing"
)

// BenchmarkWorkingMemory_AddFactWithPKSimple benchmarks adding facts with simple PK-based IDs
func BenchmarkWorkingMemory_AddFactWithPKSimple(b *testing.B) {
	wm := &WorkingMemory{
		NodeID: "benchmark_node",
		Facts:  make(map[string]*Fact),
		Tokens: make(map[string]*Token),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fact := &Fact{
			ID:   fmt.Sprintf("Person~User%d", i),
			Type: "Person",
			Fields: map[string]interface{}{
				"nom": fmt.Sprintf("User%d", i),
				"age": 30,
			},
		}
		_ = wm.AddFact(fact)
	}
}

// BenchmarkWorkingMemory_AddFactWithPKComposite benchmarks adding facts with composite PK-based IDs
func BenchmarkWorkingMemory_AddFactWithPKComposite(b *testing.B) {
	wm := &WorkingMemory{
		NodeID: "benchmark_node",
		Facts:  make(map[string]*Fact),
		Tokens: make(map[string]*Token),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fact := &Fact{
			ID:   fmt.Sprintf("Person~User%d_Company%d", i, i%100),
			Type: "Person",
			Fields: map[string]interface{}{
				"user_id":    i,
				"company_id": i % 100,
				"name":       fmt.Sprintf("User%d", i),
			},
		}
		_ = wm.AddFact(fact)
	}
}

// BenchmarkWorkingMemory_AddFactWithHashID benchmarks adding facts with hash-based IDs
func BenchmarkWorkingMemory_AddFactWithHashID(b *testing.B) {
	wm := &WorkingMemory{
		NodeID: "benchmark_node",
		Facts:  make(map[string]*Fact),
		Tokens: make(map[string]*Token),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fact := &Fact{
			ID:   fmt.Sprintf("Event~%016x", i),
			Type: "Event",
			Fields: map[string]interface{}{
				"timestamp": i,
				"message":   fmt.Sprintf("Message %d", i),
			},
		}
		_ = wm.AddFact(fact)
	}
}

// BenchmarkWorkingMemory_GetFactByTypeAndID benchmarks retrieving facts by type and ID
func BenchmarkWorkingMemory_GetFactByTypeAndID(b *testing.B) {
	wm := &WorkingMemory{
		NodeID: "benchmark_node",
		Facts:  make(map[string]*Fact),
		Tokens: make(map[string]*Token),
	}

	// Prepopulate with facts
	numFacts := 1000
	for i := 0; i < numFacts; i++ {
		fact := &Fact{
			ID:   fmt.Sprintf("Person~User%d", i),
			Type: "Person",
			Fields: map[string]interface{}{
				"nom": fmt.Sprintf("User%d", i),
			},
		}
		_ = wm.AddFact(fact)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		factID := fmt.Sprintf("Person~User%d", i%numFacts)
		_, _ = wm.GetFactByTypeAndID("Person", factID)
	}
}

// BenchmarkWorkingMemory_RemoveFact benchmarks fact removal with new ID formats
func BenchmarkWorkingMemory_RemoveFact(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		wm := &WorkingMemory{
			NodeID: "benchmark_node",
			Facts:  make(map[string]*Fact),
			Tokens: make(map[string]*Token),
		}
		fact := &Fact{
			ID:   fmt.Sprintf("Person~User%d", i),
			Type: "Person",
			Fields: map[string]interface{}{
				"nom": fmt.Sprintf("User%d", i),
			},
		}
		_ = wm.AddFact(fact)
		internalID := fact.GetInternalID()
		b.StartTimer()

		wm.RemoveFact(internalID)
	}
}

// BenchmarkEvaluator_IDFieldAccess benchmarks evaluating expressions with id field access
func BenchmarkEvaluator_IDFieldAccess(b *testing.B) {
	eval := NewAlphaConditionEvaluator()
	fact := &Fact{
		ID:   "Person~Alice",
		Type: "Person",
		Fields: map[string]interface{}{
			"nom": "Alice",
			"age": 30,
		},
	}

	expression := map[string]interface{}{
		"type": "binaryOp",
		"op":   "==",
		"left": map[string]interface{}{
			"type":   "fieldAccess",
			"object": "p",
			"field":  "id",
		},
		"right": map[string]interface{}{
			"type":  "stringLiteral",
			"value": "Person~Alice",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = eval.EvaluateCondition(expression, fact, "p")
	}
}

// BenchmarkEvaluator_IDFieldAccess_Contains benchmarks string functions on id field
func BenchmarkEvaluator_IDFieldAccess_Contains(b *testing.B) {
	eval := NewAlphaConditionEvaluator()
	fact := &Fact{
		ID:   "Person~Alice_Dupont",
		Type: "Person",
		Fields: map[string]interface{}{
			"prenom": "Alice",
			"nom":    "Dupont",
		},
	}

	expression := map[string]interface{}{
		"type": "functionCall",
		"name": "contains",
		"args": []interface{}{
			map[string]interface{}{
				"type":   "fieldAccess",
				"object": "p",
				"field":  "id",
			},
			map[string]interface{}{
				"type":  "stringLiteral",
				"value": "Alice",
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = eval.EvaluateCondition(expression, fact, "p")
	}
}

// BenchmarkEvaluator_IDComparison_BetweenFacts benchmarks comparing IDs between facts
func BenchmarkEvaluator_IDComparison_BetweenFacts(b *testing.B) {
	eval := NewAlphaConditionEvaluator()

	fact1 := &Fact{
		ID:     "Person~Alice",
		Type:   "Person",
		Fields: map[string]interface{}{"nom": "Alice"},
	}

	fact2 := &Fact{
		ID:     "Person~Alice",
		Type:   "Person",
		Fields: map[string]interface{}{"nom": "Alice"},
	}

	eval.variableBindings["p1"] = fact1
	eval.variableBindings["p2"] = fact2

	expression := map[string]interface{}{
		"type": "binaryOp",
		"op":   "==",
		"left": map[string]interface{}{
			"type":   "fieldAccess",
			"object": "p1",
			"field":  "id",
		},
		"right": map[string]interface{}{
			"type":   "fieldAccess",
			"object": "p2",
			"field":  "id",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = eval.evaluateExpression(expression)
	}
}

// BenchmarkFact_GetInternalID benchmarks internal ID generation
func BenchmarkFact_GetInternalID(b *testing.B) {
	fact := &Fact{
		ID:   "Person~Alice_Dupont",
		Type: "Person",
		Fields: map[string]interface{}{
			"prenom": "Alice",
			"nom":    "Dupont",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = fact.GetInternalID()
	}
}

// BenchmarkMakeInternalID benchmarks manual internal ID construction
func BenchmarkMakeInternalID(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MakeInternalID("Person", "Person~Alice_Dupont")
	}
}

// BenchmarkParseInternalID benchmarks internal ID parsing
func BenchmarkParseInternalID(b *testing.B) {
	internalID := "Person_Person~Alice_Dupont"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = ParseInternalID(internalID)
	}
}

// BenchmarkWorkingMemory_LargeScale benchmarks working memory with many facts
func BenchmarkWorkingMemory_LargeScale(b *testing.B) {
	sizes := []int{100, 1000, 10000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("Facts_%d", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				wm := &WorkingMemory{
					NodeID: "benchmark_node",
					Facts:  make(map[string]*Fact),
					Tokens: make(map[string]*Token),
				}
				b.StartTimer()

				for j := 0; j < size; j++ {
					fact := &Fact{
						ID:   fmt.Sprintf("Person~User%d", j),
						Type: "Person",
						Fields: map[string]interface{}{
							"nom": fmt.Sprintf("User%d", j),
							"age": j % 100,
						},
					}
					_ = wm.AddFact(fact)
				}
			}
		})
	}
}

// BenchmarkWorkingMemory_MixedOperations benchmarks mixed add/get/remove operations
func BenchmarkWorkingMemory_MixedOperations(b *testing.B) {
	wm := &WorkingMemory{
		NodeID: "benchmark_node",
		Facts:  make(map[string]*Fact),
		Tokens: make(map[string]*Token),
	}

	// Prepopulate with some facts
	for i := 0; i < 500; i++ {
		fact := &Fact{
			ID:   fmt.Sprintf("Person~User%d", i),
			Type: "Person",
			Fields: map[string]interface{}{
				"nom": fmt.Sprintf("User%d", i),
			},
		}
		_ = wm.AddFact(fact)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Add
		addFact := &Fact{
			ID:   fmt.Sprintf("Person~User%d", 500+i),
			Type: "Person",
			Fields: map[string]interface{}{
				"nom": fmt.Sprintf("User%d", 500+i),
			},
		}
		_ = wm.AddFact(addFact)

		// Get
		_, _ = wm.GetFactByTypeAndID("Person", fmt.Sprintf("Person~User%d", i%500))

		// Remove
		removeID := fmt.Sprintf("Person_Person~User%d", i%500)
		wm.RemoveFact(removeID)
	}
}
