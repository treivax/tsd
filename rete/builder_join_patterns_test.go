// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"testing"
)

// TestBuildJoinPatterns_3Variables validates buildJoinPatterns for 3 variables
func TestBuildJoinPatterns_3Variables(t *testing.T) {
	t.Log("🧪 TEST buildJoinPatterns - 3 variables")
	t.Log("========================================")

	variableNames := []string{"u", "o", "p"}
	variableTypes := []string{"User", "Order", "Product"}

	// Create a minimal builder with utils
	storage := NewMemoryStorage()
	utils := NewBuilderUtils(storage)
	builder := &JoinRuleBuilder{utils: utils}

	patterns := builder.buildJoinPatterns(variableNames, variableTypes, nil)

	// Must create 2 patterns (N-1)
	if len(patterns) != 2 {
		t.Fatalf("❌ Attendu 2 patterns, got %d", len(patterns))
	}

	// Pattern 1: [u] + [o] = [u, o]
	p1 := patterns[0]
	if !slicesEqual(p1.LeftVars, []string{"u"}) {
		t.Errorf("❌ Pattern 1 LeftVars incorrect: %v (attendu [u])", p1.LeftVars)
	}
	if !slicesEqual(p1.RightVars, []string{"o"}) {
		t.Errorf("❌ Pattern 1 RightVars incorrect: %v (attendu [o])", p1.RightVars)
	}
	if !slicesEqual(p1.AllVars, []string{"u", "o"}) {
		t.Errorf("❌ Pattern 1 AllVars incorrect: %v (attendu [u, o])", p1.AllVars)
	}

	// Verify VarTypes contains all three variables
	if len(p1.VarTypes) != 3 {
		t.Errorf("❌ Pattern 1 VarTypes devrait contenir 3 variables, got %d: %v", len(p1.VarTypes), p1.VarTypes)
	}
	expectedTypes := map[string]string{"u": "User", "o": "Order", "p": "Product"}
	for varName, expectedType := range expectedTypes {
		if actualType, exists := p1.VarTypes[varName]; !exists {
			t.Errorf("❌ Pattern 1 VarTypes manque la variable '%s'", varName)
		} else if actualType != expectedType {
			t.Errorf("❌ Pattern 1 VarTypes['%s'] = '%s', attendu '%s'", varName, actualType, expectedType)
		}
	}

	// Pattern 2: [u, o] + [p] = [u, o, p]
	p2 := patterns[1]
	if !slicesEqual(p2.LeftVars, []string{"u", "o"}) {
		t.Errorf("❌ Pattern 2 LeftVars incorrect: %v (attendu [u, o])", p2.LeftVars)
	}
	if !slicesEqual(p2.RightVars, []string{"p"}) {
		t.Errorf("❌ Pattern 2 RightVars incorrect: %v (attendu [p])", p2.RightVars)
	}
	if !slicesEqual(p2.AllVars, []string{"u", "o", "p"}) {
		t.Errorf("❌ Pattern 2 AllVars incorrect: %v (attendu [u, o, p])", p2.AllVars)
	}

	// Verify VarTypes contains all three variables
	if len(p2.VarTypes) != 3 {
		t.Errorf("❌ Pattern 2 VarTypes devrait contenir 3 variables, got %d: %v", len(p2.VarTypes), p2.VarTypes)
	}
	for varName, expectedType := range expectedTypes {
		if actualType, exists := p2.VarTypes[varName]; !exists {
			t.Errorf("❌ Pattern 2 VarTypes manque la variable '%s'", varName)
		} else if actualType != expectedType {
			t.Errorf("❌ Pattern 2 VarTypes['%s'] = '%s', attendu '%s'", varName, actualType, expectedType)
		}
	}

	t.Log("✅ Patterns corrects pour 3 variables")
}

// TestBuildJoinPatterns_NVariables validates buildJoinPatterns for N variables
func TestBuildJoinPatterns_NVariables(t *testing.T) {
	t.Log("🧪 TEST buildJoinPatterns - N variables")
	t.Log("=======================================")

	storage := NewMemoryStorage()
	utils := NewBuilderUtils(storage)
	builder := &JoinRuleBuilder{utils: utils}

	for n := 2; n <= 5; n++ {
		t.Run(fmt.Sprintf("n=%d", n), func(t *testing.T) {
			// Generate N variables
			vars := make([]string, n)
			types := make([]string, n)
			expectedVarTypes := make(map[string]string)
			for i := 0; i < n; i++ {
				vars[i] = fmt.Sprintf("v%d", i)
				types[i] = fmt.Sprintf("Type%d", i)
				expectedVarTypes[vars[i]] = types[i]
			}

			patterns := builder.buildJoinPatterns(vars, types, nil)

			// Verify number of patterns
			if len(patterns) != n-1 {
				t.Errorf("❌ Pour %d variables, attendu %d patterns, got %d", n, n-1, len(patterns))
			}

			// Verify each pattern
			for i, pattern := range patterns {
				expectedAllVars := i + 2 // Pattern i joins (i+2) variables
				if len(pattern.AllVars) != expectedAllVars {
					t.Errorf("❌ Pattern %d: attendu %d AllVars, got %d (AllVars=%v)",
						i, expectedAllVars, len(pattern.AllVars), pattern.AllVars)
				}

				// Verify LeftVars count
				expectedLeftVars := i + 1 // Pattern i has (i+1) left variables
				if len(pattern.LeftVars) != expectedLeftVars {
					t.Errorf("❌ Pattern %d: attendu %d LeftVars, got %d (LeftVars=%v)",
						i, expectedLeftVars, len(pattern.LeftVars), pattern.LeftVars)
				}

				// Verify RightVars count (always 1 in cascade)
				if len(pattern.RightVars) != 1 {
					t.Errorf("❌ Pattern %d: attendu 1 RightVar, got %d (RightVars=%v)",
						i, len(pattern.RightVars), pattern.RightVars)
				}

				// Verify VarTypes contains ALL variables
				if len(pattern.VarTypes) != n {
					t.Errorf("❌ Pattern %d: VarTypes devrait contenir %d variables, got %d: %v",
						i, n, len(pattern.VarTypes), pattern.VarTypes)
				}

				// Verify all expected types are present
				for varName, expectedType := range expectedVarTypes {
					if actualType, exists := pattern.VarTypes[varName]; !exists {
						t.Errorf("❌ Pattern %d: VarTypes manque la variable '%s'", i, varName)
					} else if actualType != expectedType {
						t.Errorf("❌ Pattern %d: VarTypes['%s'] = '%s', attendu '%s'",
							i, varName, actualType, expectedType)
					}
				}
			}
		})
	}

	t.Log("✅ Patterns corrects pour N variables")
}

// TestBuildJoinPatterns_2Variables validates buildJoinPatterns for 2 variables (simple join)
func TestBuildJoinPatterns_2Variables(t *testing.T) {
	t.Log("🧪 TEST buildJoinPatterns - 2 variables (non-régression)")
	t.Log("========================================================")

	variableNames := []string{"p", "o"}
	variableTypes := []string{"Person", "Order"}

	storage := NewMemoryStorage()
	utils := NewBuilderUtils(storage)
	builder := &JoinRuleBuilder{utils: utils}

	patterns := builder.buildJoinPatterns(variableNames, variableTypes, nil)

	// Must create 1 pattern
	if len(patterns) != 1 {
		t.Fatalf("❌ Attendu 1 pattern, got %d", len(patterns))
	}

	p := patterns[0]
	if !slicesEqual(p.LeftVars, []string{"p"}) {
		t.Errorf("❌ LeftVars incorrect: %v (attendu [p])", p.LeftVars)
	}
	if !slicesEqual(p.RightVars, []string{"o"}) {
		t.Errorf("❌ RightVars incorrect: %v (attendu [o])", p.RightVars)
	}
	if !slicesEqual(p.AllVars, []string{"p", "o"}) {
		t.Errorf("❌ AllVars incorrect: %v (attendu [p, o])", p.AllVars)
	}

	// Verify VarTypes
	if len(p.VarTypes) != 2 {
		t.Errorf("❌ VarTypes devrait contenir 2 variables, got %d: %v", len(p.VarTypes), p.VarTypes)
	}
	if p.VarTypes["p"] != "Person" {
		t.Errorf("❌ VarTypes['p'] = '%s', attendu 'Person'", p.VarTypes["p"])
	}
	if p.VarTypes["o"] != "Order" {
		t.Errorf("❌ VarTypes['o'] = '%s', attendu 'Order'", p.VarTypes["o"])
	}

	t.Log("✅ Pattern correct pour 2 variables")
}

// slicesEqual checks if two string slices are equal
func slicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
