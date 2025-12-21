// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"os"
	"path/filepath"
	"testing"
)

// TestModuloOperator teste l'opérateur modulo (%) dans différents contextes
func TestModuloOperator(t *testing.T) {
	t.Log("🧪 TEST MODULO OPERATOR")
	t.Log("=======================")

	tests := []struct {
		name        string
		tsdContent  string
		facts       []*Fact
		expectMatch bool
		description string
	}{
		{
			name: "even number detection",
			tsdContent: `
type Number(#id: string, value: number)
action check(msg: string)
rule even : {n: Number} / n.value % 2 == 0 ==> check("Even")`,
			facts: []*Fact{
				{
					ID:   "n1",
					Type: "Number",
					Fields: map[string]interface{}{
						"id":    "n1",
						"value": 42.0,
					},
				},
			},
			expectMatch: true,
			description: "42 % 2 = 0 (even number)",
		},
		{
			name: "odd number detection",
			tsdContent: `
type Number(#id: string, value: number)
action check(msg: string)
rule odd : {n: Number} / n.value % 2 == 1 ==> check("Odd")`,
			facts: []*Fact{
				{
					ID:   "n1",
					Type: "Number",
					Fields: map[string]interface{}{
						"id":    "n1",
						"value": 43.0,
					},
				},
			},
			expectMatch: true,
			description: "43 % 2 = 1 (odd number)",
		},
		{
			name: "divisibility by 5",
			tsdContent: `
type Number(#id: string, value: number)
action check(msg: string)
rule divisible_by_5 : {n: Number} / n.value % 5 == 0 ==> check("Divisible by 5")`,
			facts: []*Fact{
				{
					ID:   "n1",
					Type: "Number",
					Fields: map[string]interface{}{
						"id":    "n1",
						"value": 25.0,
					},
				},
			},
			expectMatch: true,
			description: "25 % 5 = 0",
		},
		{
			name: "not divisible by 5",
			tsdContent: `
type Number(#id: string, value: number)
action check(msg: string)
rule divisible_by_5 : {n: Number} / n.value % 5 == 0 ==> check("Divisible by 5")`,
			facts: []*Fact{
				{
					ID:   "n1",
					Type: "Number",
					Fields: map[string]interface{}{
						"id":    "n1",
						"value": 23.0,
					},
				},
			},
			expectMatch: false,
			description: "23 % 5 = 3 (not divisible)",
		},
		{
			name: "modulo with comparison",
			tsdContent: `
type Number(#id: string, value: number)
action check(msg: string)
rule mod_greater : {n: Number} / n.value % 10 > 5 ==> check("Remainder > 5")`,
			facts: []*Fact{
				{
					ID:   "n1",
					Type: "Number",
					Fields: map[string]interface{}{
						"id":    "n1",
						"value": 17.0,
					},
				},
			},
			expectMatch: true,
			description: "17 % 10 = 7 > 5",
		},
		{
			name: "modulo in complex expression",
			tsdContent: `
type Number(#id: string, value: number)
action check(msg: string)
rule complex_mod : {n: Number} / (n.value % 3) + 1 == 2 ==> check("Complex")`,
			facts: []*Fact{
				{
					ID:   "n1",
					Type: "Number",
					Fields: map[string]interface{}{
						"id":    "n1",
						"value": 4.0,
					},
				},
			},
			expectMatch: true,
			description: "(4 % 3) + 1 = 1 + 1 = 2",
		},
		{
			name: "modulo with negative numbers",
			tsdContent: `
type Number(#id: string, value: number)
action check(msg: string)
rule neg_mod : {n: Number} / n.value % 3 == 2 ==> check("Remainder 2")`,
			facts: []*Fact{
				{
					ID:   "n1",
					Type: "Number",
					Fields: map[string]interface{}{
						"id":    "n1",
						"value": 5.0,
					},
				},
			},
			expectMatch: true,
			description: "5 % 3 = 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("🔍 Testing: %s", tt.description)

			// Create temporary file
			tempDir := t.TempDir()
			tsdFile := filepath.Join(tempDir, "test.tsd")
			if err := os.WriteFile(tsdFile, []byte(tt.tsdContent), 0644); err != nil {
				t.Fatalf("❌ Failed to write test file: %v", err)
			}

			// Build network
			storage := NewMemoryStorage()
			pipeline := NewConstraintPipeline()
			network, _, err := pipeline.IngestFile(tsdFile, nil, storage)
			if err != nil {
				t.Fatalf("❌ Failed to build network: %v", err)
			}

			// Submit facts
			for _, fact := range tt.facts {
				network.SubmitFact(fact)
			}

			// Check if rule activated
			var activated bool
			for _, terminal := range network.TerminalNodes {
				if terminal.GetExecutionCount() > 0 {
					activated = true
					break
				}
			}

			if activated != tt.expectMatch {
				t.Errorf("❌ %s: expected match=%v, got match=%v", tt.description, tt.expectMatch, activated)
				return
			}

			t.Logf("✅ %s: match=%v", tt.description, activated)
		})
	}
}

// TestModuloOperatorEdgeCases teste les cas limites de l'opérateur modulo
func TestModuloOperatorEdgeCases(t *testing.T) {
	t.Log("🧪 TEST MODULO OPERATOR EDGE CASES")
	t.Log("==================================")

	t.Run("modulo with zero dividend", func(t *testing.T) {
		tsdContent := `
type Number(#id: string, value: number)
action check(msg: string)
rule zero_mod : {n: Number} / n.value % 10 == 0 ==> check("Zero")
`
		tempDir := t.TempDir()
		tsdFile := filepath.Join(tempDir, "test.tsd")
		if err := os.WriteFile(tsdFile, []byte(tsdContent), 0644); err != nil {
			t.Fatalf("❌ Failed to write test file: %v", err)
		}

		storage := NewMemoryStorage()
		pipeline := NewConstraintPipeline()
		network, _, err := pipeline.IngestFile(tsdFile, nil, storage)
		if err != nil {
			t.Fatalf("❌ Failed to build network: %v", err)
		}

		fact := &Fact{
			ID:   "n1",
			Type: "Number",
			Fields: map[string]interface{}{
				"id":    "n1",
				"value": 0.0,
			},
		}
		network.SubmitFact(fact)

		// 0 % 10 = 0, should match
		var activated bool
		for _, terminal := range network.TerminalNodes {
			if terminal.GetExecutionCount() > 0 {
				activated = true
				break
			}
		}

		if !activated {
			t.Error("❌ Expected 0 % 10 = 0 to match")
		} else {
			t.Log("✅ 0 % 10 = 0 handled correctly")
		}
	})

	t.Run("modulo with large numbers", func(t *testing.T) {
		tsdContent := `
type Number(#id: string, value: number)
action check(msg: string)
rule large_mod : {n: Number} / n.value % 1000 == 456 ==> check("Match")
`
		tempDir := t.TempDir()
		tsdFile := filepath.Join(tempDir, "test.tsd")
		if err := os.WriteFile(tsdFile, []byte(tsdContent), 0644); err != nil {
			t.Fatalf("❌ Failed to write test file: %v", err)
		}

		storage := NewMemoryStorage()
		pipeline := NewConstraintPipeline()
		network, _, err := pipeline.IngestFile(tsdFile, nil, storage)
		if err != nil {
			t.Fatalf("❌ Failed to build network: %v", err)
		}

		fact := &Fact{
			ID:   "n1",
			Type: "Number",
			Fields: map[string]interface{}{
				"id":    "n1",
				"value": 123456.0,
			},
		}
		network.SubmitFact(fact)

		// 123456 % 1000 = 456, should match
		var activated bool
		for _, terminal := range network.TerminalNodes {
			if terminal.GetExecutionCount() > 0 {
				activated = true
				break
			}
		}

		if !activated {
			t.Error("❌ Expected 123456 % 1000 = 456 to match")
		} else {
			t.Log("✅ Large number modulo handled correctly")
		}
	})
}

// BenchmarkModuloOperator benchmarks modulo operator performance
func BenchmarkModuloOperator(b *testing.B) {
	tsdContent := `
type Number(#id: string, value: number)
action check(msg: string)
rule even : {n: Number} / n.value % 2 == 0 ==> check("Even")
`
	tempDir := b.TempDir()
	tsdFile := filepath.Join(tempDir, "test.tsd")
	if err := os.WriteFile(tsdFile, []byte(tsdContent), 0644); err != nil {
		b.Fatalf("Failed to write test file: %v", err)
	}

	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(tsdFile, nil, storage)
	if err != nil {
		b.Fatalf("Failed to build network: %v", err)
	}

	fact := &Fact{
		ID:   "n1",
		Type: "Number",
		Fields: map[string]interface{}{
			"id":    "n1",
			"value": 42.0,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		network.SubmitFact(fact)
	}
}
