// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"testing"
)

func TestConditionSplitter_SingleAlphaCondition(t *testing.T) {
	splitter := NewConditionSplitter()

	// Alpha condition: c.qte > 5 (only variable "c")
	condition := map[string]interface{}{
		"type": "comparison",
		"left": map[string]interface{}{
			"type":   "fieldAccess",
			"object": "c",
			"field":  "qte",
		},
		"operator": ">",
		"right": map[string]interface{}{
			"type":  "number",
			"value": 5.0,
		},
	}

	alphas, betas, err := splitter.SplitConditions(condition)
	if err != nil {
		t.Fatalf("Error splitting condition: %v", err)
	}

	if len(alphas) != 1 {
		t.Errorf("Expected 1 alpha condition, got %d", len(alphas))
	}

	if len(betas) != 0 {
		t.Errorf("Expected 0 beta conditions, got %d", len(betas))
	}

	if len(alphas) > 0 {
		if alphas[0].Type != ConditionTypeAlpha {
			t.Errorf("Expected alpha type, got %v", alphas[0].Type)
		}
		if len(alphas[0].Variables) != 1 {
			t.Errorf("Expected 1 variable, got %d", len(alphas[0].Variables))
		}
		if alphas[0].Variable != "c" {
			t.Errorf("Expected variable 'c', got '%s'", alphas[0].Variable)
		}
	}

	t.Log("✅ Single alpha condition correctly classified")
}

func TestConditionSplitter_SingleBetaCondition(t *testing.T) {
	splitter := NewConditionSplitter()

	// Beta condition: c.produit_id == p.id (variables "c" and "p")
	condition := map[string]interface{}{
		"type": "comparison",
		"left": map[string]interface{}{
			"type":   "fieldAccess",
			"object": "c",
			"field":  "produit_id",
		},
		"operator": "==",
		"right": map[string]interface{}{
			"type":   "fieldAccess",
			"object": "p",
			"field":  "id",
		},
	}

	alphas, betas, err := splitter.SplitConditions(condition)
	if err != nil {
		t.Fatalf("Error splitting condition: %v", err)
	}

	if len(alphas) != 0 {
		t.Errorf("Expected 0 alpha conditions, got %d", len(alphas))
	}

	if len(betas) != 1 {
		t.Errorf("Expected 1 beta condition, got %d", len(betas))
	}

	if len(betas) > 0 {
		if betas[0].Type != ConditionTypeBeta {
			t.Errorf("Expected beta type, got %v", betas[0].Type)
		}
		if len(betas[0].Variables) != 2 {
			t.Errorf("Expected 2 variables, got %d", len(betas[0].Variables))
		}
	}

	t.Log("✅ Single beta condition correctly classified")
}

// TestConditionSplitter_ArithmeticExpressions verifies that arithmetic expressions
// with single variable are classified as alpha conditions
func TestConditionSplitter_ArithmeticExpressions(t *testing.T) {
	splitter := NewConditionSplitter()

	tests := []struct {
		name      string
		condition map[string]interface{}
		wantAlpha int
		wantBeta  int
		variables []string
	}{
		{
			name: "simple arithmetic - single variable",
			condition: map[string]interface{}{
				"type": "comparison",
				"left": map[string]interface{}{
					"type": "binaryOp",
					"left": map[string]interface{}{
						"type":   "fieldAccess",
						"object": "c",
						"field":  "qte",
					},
					"operator": "*",
					"right": map[string]interface{}{
						"type":  "number",
						"value": 23.0,
					},
				},
				"operator": ">",
				"right": map[string]interface{}{
					"type":  "number",
					"value": 100.0,
				},
			},
			wantAlpha: 1,
			wantBeta:  0,
			variables: []string{"c"},
		},
		{
			name: "complex arithmetic - single variable",
			condition: map[string]interface{}{
				"type": "comparison",
				"left": map[string]interface{}{
					"type": "binaryOp",
					"left": map[string]interface{}{
						"type": "binaryOp",
						"left": map[string]interface{}{
							"type":   "fieldAccess",
							"object": "c",
							"field":  "qte",
						},
						"operator": "*",
						"right": map[string]interface{}{
							"type":  "number",
							"value": 23.0,
						},
					},
					"operator": "-",
					"right": map[string]interface{}{
						"type":  "number",
						"value": 10.0,
					},
				},
				"operator": ">",
				"right": map[string]interface{}{
					"type":  "number",
					"value": 0.0,
				},
			},
			wantAlpha: 1,
			wantBeta:  0,
			variables: []string{"c"},
		},
		{
			name: "arithmetic with two variables - should be beta",
			condition: map[string]interface{}{
				"type": "comparison",
				"left": map[string]interface{}{
					"type": "binaryOp",
					"left": map[string]interface{}{
						"type":   "fieldAccess",
						"object": "c",
						"field":  "qte",
					},
					"operator": "*",
					"right": map[string]interface{}{
						"type":   "fieldAccess",
						"object": "p",
						"field":  "price",
					},
				},
				"operator": ">",
				"right": map[string]interface{}{
					"type":  "number",
					"value": 100.0,
				},
			},
			wantAlpha: 0,
			wantBeta:  1,
			variables: []string{"c", "p"},
		},
		{
			name: "nested arithmetic - single variable",
			condition: map[string]interface{}{
				"type": "comparison",
				"left": map[string]interface{}{
					"type": "binaryOp",
					"left": map[string]interface{}{
						"type": "binaryOp",
						"left": map[string]interface{}{
							"type": "binaryOp",
							"left": map[string]interface{}{
								"type":   "fieldAccess",
								"object": "o",
								"field":  "quantity",
							},
							"operator": "*",
							"right": map[string]interface{}{
								"type":   "fieldAccess",
								"object": "o",
								"field":  "price",
							},
						},
						"operator": "-",
						"right": map[string]interface{}{
							"type":   "fieldAccess",
							"object": "o",
							"field":  "discount",
						},
					},
					"operator": "/",
					"right": map[string]interface{}{
						"type":  "number",
						"value": 2.0,
					},
				},
				"operator": ">",
				"right": map[string]interface{}{
					"type":  "number",
					"value": 50.0,
				},
			},
			wantAlpha: 1,
			wantBeta:  0,
			variables: []string{"o"},
		},
		{
			name: "arithmetic with modulo - single variable",
			condition: map[string]interface{}{
				"type": "comparison",
				"left": map[string]interface{}{
					"type": "binaryOp",
					"left": map[string]interface{}{
						"type":   "fieldAccess",
						"object": "n",
						"field":  "value",
					},
					"operator": "%",
					"right": map[string]interface{}{
						"type":  "number",
						"value": 2.0,
					},
				},
				"operator": "==",
				"right": map[string]interface{}{
					"type":  "number",
					"value": 0.0,
				},
			},
			wantAlpha: 1,
			wantBeta:  0,
			variables: []string{"n"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			alphas, betas, err := splitter.SplitConditions(tt.condition)
			if err != nil {
				t.Fatalf("Error splitting condition: %v", err)
			}

			if len(alphas) != tt.wantAlpha {
				t.Errorf("Expected %d alpha conditions, got %d", tt.wantAlpha, len(alphas))
			}

			if len(betas) != tt.wantBeta {
				t.Errorf("Expected %d beta conditions, got %d", tt.wantBeta, len(betas))
			}

			// Verify variables extracted correctly
			if len(alphas) > 0 {
				vars := alphas[0].Variables
				if len(vars) != len(tt.variables) {
					t.Errorf("Expected %d variables, got %d: %v", len(tt.variables), len(vars), vars)
				}
			}

			if len(betas) > 0 {
				vars := betas[0].Variables
				if len(vars) != len(tt.variables) {
					t.Errorf("Expected %d variables, got %d: %v", len(tt.variables), len(vars), vars)
				}
			}
		})
	}

	t.Log("✅ Arithmetic expressions correctly classified")
}

// TestConditionSplitter_Division verifies division operations
func TestConditionSplitter_Division(t *testing.T) {
	splitter := NewConditionSplitter()

	// c.value / 10 > 5
	condition := map[string]interface{}{
		"type": "comparison",
		"left": map[string]interface{}{
			"type": "binaryOp",
			"left": map[string]interface{}{
				"type":   "fieldAccess",
				"object": "c",
				"field":  "value",
			},
			"operator": "/",
			"right": map[string]interface{}{
				"type":  "number",
				"value": 10.0,
			},
		},
		"operator": ">",
		"right": map[string]interface{}{
			"type":  "number",
			"value": 5.0,
		},
	}

	alphas, betas, err := splitter.SplitConditions(condition)
	if err != nil {
		t.Fatalf("Error splitting condition: %v", err)
	}

	if len(alphas) != 1 {
		t.Errorf("Expected 1 alpha condition for division, got %d", len(alphas))
	}

	if len(betas) != 0 {
		t.Errorf("Expected 0 beta conditions, got %d", len(betas))
	}

	t.Log("✅ Division operations correctly handled")
}

// TestConditionSplitter_NegativeNumbers verifies negative number operations
func TestConditionSplitter_NegativeNumbers(t *testing.T) {
	splitter := NewConditionSplitter()

	// c.balance * -1 > 100
	condition := map[string]interface{}{
		"type": "comparison",
		"left": map[string]interface{}{
			"type": "binaryOp",
			"left": map[string]interface{}{
				"type":   "fieldAccess",
				"object": "c",
				"field":  "balance",
			},
			"operator": "*",
			"right": map[string]interface{}{
				"type":  "number",
				"value": -1.0,
			},
		},
		"operator": ">",
		"right": map[string]interface{}{
			"type":  "number",
			"value": 100.0,
		},
	}

	alphas, betas, err := splitter.SplitConditions(condition)
	if err != nil {
		t.Fatalf("Error splitting condition: %v", err)
	}

	if len(alphas) != 1 {
		t.Errorf("Expected 1 alpha condition for negative arithmetic, got %d", len(alphas))
	}

	if len(betas) != 0 {
		t.Errorf("Expected 0 beta conditions, got %d", len(betas))
	}

	t.Log("✅ Negative number operations correctly handled")
}

func TestConditionSplitter_MixedAlphaAndBeta(t *testing.T) {
	splitter := NewConditionSplitter()

	// Mixed condition: c.produit_id == p.id AND c.qte > 5
	// This is the structure from the TSD parser with logicalExpr
	condition := map[string]interface{}{
		"type": "logicalExpr",
		"left": map[string]interface{}{
			"type": "comparison",
			"left": map[string]interface{}{
				"type":   "fieldAccess",
				"object": "c",
				"field":  "produit_id",
			},
			"operator": "==",
			"right": map[string]interface{}{
				"type":   "fieldAccess",
				"object": "p",
				"field":  "id",
			},
		},
		"operations": []interface{}{
			map[string]interface{}{
				"op": "AND",
				"right": map[string]interface{}{
					"type": "comparison",
					"left": map[string]interface{}{
						"type":   "fieldAccess",
						"object": "c",
						"field":  "qte",
					},
					"operator": ">",
					"right": map[string]interface{}{
						"type":  "number",
						"value": 5.0,
					},
				},
			},
		},
	}

	alphas, betas, err := splitter.SplitConditions(condition)
	if err != nil {
		t.Fatalf("Error splitting condition: %v", err)
	}

	t.Logf("Alpha conditions: %d", len(alphas))
	t.Logf("Beta conditions: %d", len(betas))

	if len(alphas) != 1 {
		t.Errorf("Expected 1 alpha condition, got %d", len(alphas))
	}

	if len(betas) != 1 {
		t.Errorf("Expected 1 beta condition, got %d", len(betas))
	}

	// Verify alpha condition
	if len(alphas) > 0 {
		if alphas[0].Type != ConditionTypeAlpha {
			t.Errorf("Expected alpha type for first alpha condition")
		}
		if alphas[0].Variable != "c" {
			t.Errorf("Expected variable 'c' for alpha condition, got '%s'", alphas[0].Variable)
		}
	}

	// Verify beta condition
	if len(betas) > 0 {
		if betas[0].Type != ConditionTypeBeta {
			t.Errorf("Expected beta type for beta condition")
		}
		if len(betas[0].Variables) != 2 {
			t.Errorf("Expected 2 variables for beta condition, got %d", len(betas[0].Variables))
		}
	}

	t.Log("✅ Mixed alpha/beta conditions correctly separated")
}

func TestConditionSplitter_MultipleAlphaConditions(t *testing.T) {
	splitter := NewConditionSplitter()

	// Multiple alpha conditions: c.qte > 5 AND c.statut == "actif"
	condition := map[string]interface{}{
		"type": "logicalExpr",
		"left": map[string]interface{}{
			"type": "comparison",
			"left": map[string]interface{}{
				"type":   "fieldAccess",
				"object": "c",
				"field":  "qte",
			},
			"operator": ">",
			"right": map[string]interface{}{
				"type":  "number",
				"value": 5.0,
			},
		},
		"operations": []interface{}{
			map[string]interface{}{
				"op": "AND",
				"right": map[string]interface{}{
					"type": "comparison",
					"left": map[string]interface{}{
						"type":   "fieldAccess",
						"object": "c",
						"field":  "statut",
					},
					"operator": "==",
					"right": map[string]interface{}{
						"type":  "string",
						"value": "actif",
					},
				},
			},
		},
	}

	alphas, betas, err := splitter.SplitConditions(condition)
	if err != nil {
		t.Fatalf("Error splitting condition: %v", err)
	}

	if len(alphas) != 2 {
		t.Errorf("Expected 2 alpha conditions, got %d", len(alphas))
	}

	if len(betas) != 0 {
		t.Errorf("Expected 0 beta conditions, got %d", len(betas))
	}

	// Both should be alpha conditions on variable "c"
	for i, alpha := range alphas {
		if alpha.Type != ConditionTypeAlpha {
			t.Errorf("Alpha condition %d: expected alpha type", i)
		}
		if alpha.Variable != "c" {
			t.Errorf("Alpha condition %d: expected variable 'c', got '%s'", i, alpha.Variable)
		}
	}

	t.Log("✅ Multiple alpha conditions correctly identified")
}

func TestConditionSplitter_ClassifyCondition(t *testing.T) {
	splitter := NewConditionSplitter()

	// Test alpha condition
	alphaCondition := map[string]interface{}{
		"type": "comparison",
		"left": map[string]interface{}{
			"type":   "fieldAccess",
			"object": "c",
			"field":  "qte",
		},
		"operator": ">",
		"right": map[string]interface{}{
			"type":  "number",
			"value": 5.0,
		},
	}

	if !splitter.IsAlphaCondition(alphaCondition) {
		t.Errorf("Expected alpha condition to be classified as alpha")
	}

	if splitter.IsBetaCondition(alphaCondition) {
		t.Errorf("Expected alpha condition NOT to be classified as beta")
	}

	// Test beta condition
	betaCondition := map[string]interface{}{
		"type": "comparison",
		"left": map[string]interface{}{
			"type":   "fieldAccess",
			"object": "c",
			"field":  "produit_id",
		},
		"operator": "==",
		"right": map[string]interface{}{
			"type":   "fieldAccess",
			"object": "p",
			"field":  "id",
		},
	}

	if splitter.IsAlphaCondition(betaCondition) {
		t.Errorf("Expected beta condition NOT to be classified as alpha")
	}

	if !splitter.IsBetaCondition(betaCondition) {
		t.Errorf("Expected beta condition to be classified as beta")
	}

	t.Log("✅ Condition classification works correctly")
}

func TestConditionSplitter_ExtractVariables(t *testing.T) {
	splitter := NewConditionSplitter()

	tests := []struct {
		name          string
		condition     map[string]interface{}
		expectedCount int
	}{
		{
			name: "Single variable",
			condition: map[string]interface{}{
				"type": "comparison",
				"left": map[string]interface{}{
					"type":   "fieldAccess",
					"object": "c",
					"field":  "qte",
				},
				"operator": ">",
				"right": map[string]interface{}{
					"type":  "number",
					"value": 5.0,
				},
			},
			expectedCount: 1,
		},
		{
			name: "Two variables",
			condition: map[string]interface{}{
				"type": "comparison",
				"left": map[string]interface{}{
					"type":   "fieldAccess",
					"object": "c",
					"field":  "produit_id",
				},
				"operator": "==",
				"right": map[string]interface{}{
					"type":   "fieldAccess",
					"object": "p",
					"field":  "id",
				},
			},
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vars := splitter.ExtractVariables(tt.condition)
			if len(vars) != tt.expectedCount {
				t.Errorf("Expected %d variables, got %d: %v", tt.expectedCount, len(vars), vars)
			}
		})
	}

	t.Log("✅ Variable extraction works correctly")
}

func TestConditionSplitter_GetPrimaryVariable(t *testing.T) {
	splitter := NewConditionSplitter()

	// Alpha condition with single variable
	alphaCondition := map[string]interface{}{
		"type": "comparison",
		"left": map[string]interface{}{
			"type":   "fieldAccess",
			"object": "c",
			"field":  "qte",
		},
		"operator": ">",
		"right": map[string]interface{}{
			"type":  "number",
			"value": 5.0,
		},
	}

	primaryVar := splitter.GetPrimaryVariable(alphaCondition)
	if primaryVar != "c" {
		t.Errorf("Expected primary variable 'c', got '%s'", primaryVar)
	}

	// Beta condition (should return empty string)
	betaCondition := map[string]interface{}{
		"type": "comparison",
		"left": map[string]interface{}{
			"type":   "fieldAccess",
			"object": "c",
			"field":  "produit_id",
		},
		"operator": "==",
		"right": map[string]interface{}{
			"type":   "fieldAccess",
			"object": "p",
			"field":  "id",
		},
	}

	primaryVar = splitter.GetPrimaryVariable(betaCondition)
	if primaryVar != "" {
		t.Errorf("Expected empty string for beta condition, got '%s'", primaryVar)
	}

	t.Log("✅ Primary variable extraction works correctly")
}

func TestConditionSplitter_ReconstructBetaCondition(t *testing.T) {
	splitter := NewConditionSplitter()

	// Create a beta condition
	betaCond := SplitCondition{
		Type: ConditionTypeBeta,
		Condition: map[string]interface{}{
			"type": "comparison",
			"left": map[string]interface{}{
				"type":   "fieldAccess",
				"object": "c",
				"field":  "produit_id",
			},
			"operator": "==",
			"right": map[string]interface{}{
				"type":   "fieldAccess",
				"object": "p",
				"field":  "id",
			},
		},
		Variables: []string{"c", "p"},
	}

	// Test with single beta condition
	reconstructed := splitter.ReconstructBetaCondition([]SplitCondition{betaCond})
	if reconstructed == nil {
		t.Errorf("Expected reconstructed condition, got nil")
	}

	// Test with empty list
	reconstructed = splitter.ReconstructBetaCondition([]SplitCondition{})
	if reconstructed != nil {
		t.Errorf("Expected nil for empty list, got %v", reconstructed)
	}

	t.Logf("✅ Beta condition reconstruction works correctly")
}

// TestConditionSplitter_DebugJoinWithAlpha tests the exact structure from the failing diagnostic test
func TestConditionSplitter_DebugJoinWithAlpha(t *testing.T) {
	// This is the exact condition structure from: p.id == o.personId AND o.amount > 100
	condition := map[string]interface{}{
		"type": "constraint",
		"constraint": map[string]interface{}{
			"type": "logicalExpr",
			"left": map[string]interface{}{
				"type": "comparison",
				"left": map[string]interface{}{
					"type":   "fieldAccess",
					"object": "p",
					"field":  "id",
				},
				"operator": "==",
				"right": map[string]interface{}{
					"type":   "fieldAccess",
					"object": "o",
					"field":  "personId",
				},
			},
			"operations": []interface{}{
				map[string]interface{}{
					"op": "AND",
					"right": map[string]interface{}{
						"type": "comparison",
						"left": map[string]interface{}{
							"type":   "fieldAccess",
							"object": "o",
							"field":  "amount",
						},
						"operator": ">",
						"right":    100.0,
					},
				},
			},
		},
	}

	splitter := NewConditionSplitter()
	alphaConditions, betaConditions, err := splitter.SplitConditions(condition)
	if err != nil {
		t.Fatalf("Failed to split conditions: %v", err)
	}

	t.Logf("Alpha conditions: %d", len(alphaConditions))
	for i, alpha := range alphaConditions {
		t.Logf("  Alpha[%d]: variable=%s, vars=%v", i, alpha.Variable, alpha.Variables)
	}

	t.Logf("Beta conditions: %d", len(betaConditions))
	for i, beta := range betaConditions {
		t.Logf("  Beta[%d]: vars=%v", i, beta.Variables)
	}

	// Expected: 1 alpha (o.amount > 100) and 1 beta (p.id == o.personId)
	if len(alphaConditions) != 1 {
		t.Errorf("Expected 1 alpha condition (o.amount > 100), got %d", len(alphaConditions))
	}

	if len(betaConditions) != 1 {
		t.Errorf("Expected 1 beta condition (p.id == o.personId), got %d", len(betaConditions))
	}

	if len(alphaConditions) > 0 {
		if alphaConditions[0].Variable != "o" {
			t.Errorf("Expected alpha condition on variable 'o', got '%s'", alphaConditions[0].Variable)
		}
	}

	if len(betaConditions) > 0 {
		if len(betaConditions[0].Variables) != 2 {
			t.Errorf("Expected beta condition with 2 variables, got %d", len(betaConditions[0].Variables))
		}
	}

	t.Logf("✅ Join with alpha condition correctly separated")
}

func TestConditionSplitter_HasMixedConditions(t *testing.T) {
	splitter := NewConditionSplitter()

	// Pure alpha condition
	pureAlpha := map[string]interface{}{
		"type": "comparison",
		"left": map[string]interface{}{
			"type":   "fieldAccess",
			"object": "c",
			"field":  "qte",
		},
		"operator": ">",
		"right": map[string]interface{}{
			"type":  "number",
			"value": 5.0,
		},
	}

	if splitter.HasMixedConditions(pureAlpha) {
		t.Errorf("Pure alpha condition should not be mixed")
	}

	// Mixed condition
	mixed := map[string]interface{}{
		"type": "logicalExpr",
		"left": map[string]interface{}{
			"type": "comparison",
			"left": map[string]interface{}{
				"type":   "fieldAccess",
				"object": "c",
				"field":  "produit_id",
			},
			"operator": "==",
			"right": map[string]interface{}{
				"type":   "fieldAccess",
				"object": "p",
				"field":  "id",
			},
		},
		"operations": []interface{}{
			map[string]interface{}{
				"op": "AND",
				"right": map[string]interface{}{
					"type": "comparison",
					"left": map[string]interface{}{
						"type":   "fieldAccess",
						"object": "c",
						"field":  "qte",
					},
					"operator": ">",
					"right": map[string]interface{}{
						"type":  "number",
						"value": 5.0,
					},
				},
			},
		},
	}

	if !splitter.HasMixedConditions(mixed) {
		t.Errorf("Mixed condition should be detected as mixed")
	}

	t.Log("✅ Mixed condition detection works correctly")
}

func TestConditionSplitter_CountConditions(t *testing.T) {
	splitter := NewConditionSplitter()

	// Single condition
	single := map[string]interface{}{
		"type": "comparison",
		"left": map[string]interface{}{
			"type":   "fieldAccess",
			"object": "c",
			"field":  "qte",
		},
		"operator": ">",
		"right": map[string]interface{}{
			"type":  "number",
			"value": 5.0,
		},
	}

	count := splitter.CountConditions(single)
	if count != 1 {
		t.Errorf("Expected 1 condition, got %d", count)
	}

	// Multiple conditions
	multiple := map[string]interface{}{
		"type": "logicalExpr",
		"left": map[string]interface{}{
			"type": "comparison",
			"left": map[string]interface{}{
				"type":   "fieldAccess",
				"object": "c",
				"field":  "produit_id",
			},
			"operator": "==",
			"right": map[string]interface{}{
				"type":   "fieldAccess",
				"object": "p",
				"field":  "id",
			},
		},
		"operations": []interface{}{
			map[string]interface{}{
				"op": "AND",
				"right": map[string]interface{}{
					"type": "comparison",
					"left": map[string]interface{}{
						"type":   "fieldAccess",
						"object": "c",
						"field":  "qte",
					},
					"operator": ">",
					"right": map[string]interface{}{
						"type":  "number",
						"value": 5.0,
					},
				},
			},
		},
	}

	count = splitter.CountConditions(multiple)
	if count != 2 {
		t.Errorf("Expected 2 conditions, got %d", count)
	}

	t.Log("✅ Condition counting works correctly")
}

func TestConditionSplitter_WithConstraintWrapper(t *testing.T) {
	splitter := NewConditionSplitter()

	// Condition with constraint wrapper (as from TSD parser)
	condition := map[string]interface{}{
		"type": "constraint",
		"constraint": map[string]interface{}{
			"type": "logicalExpr",
			"left": map[string]interface{}{
				"type": "comparison",
				"left": map[string]interface{}{
					"type":   "fieldAccess",
					"object": "c",
					"field":  "produit_id",
				},
				"operator": "==",
				"right": map[string]interface{}{
					"type":   "fieldAccess",
					"object": "p",
					"field":  "id",
				},
			},
			"operations": []interface{}{
				map[string]interface{}{
					"op": "AND",
					"right": map[string]interface{}{
						"type": "comparison",
						"left": map[string]interface{}{
							"type":   "fieldAccess",
							"object": "c",
							"field":  "qte",
						},
						"operator": ">",
						"right": map[string]interface{}{
							"type":  "number",
							"value": 5.0,
						},
					},
				},
			},
		},
	}

	alphas, betas, err := splitter.SplitConditions(condition)
	if err != nil {
		t.Fatalf("Error splitting condition with wrapper: %v", err)
	}

	if len(alphas) != 1 {
		t.Errorf("Expected 1 alpha condition, got %d", len(alphas))
	}

	if len(betas) != 1 {
		t.Errorf("Expected 1 beta condition, got %d", len(betas))
	}

	t.Log("✅ Condition with constraint wrapper correctly handled")
}
