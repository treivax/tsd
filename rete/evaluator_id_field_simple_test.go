// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"testing"
)

// TestEvaluator_IDFieldAccess_BasicComparisons teste l'accès basique au champ id
func TestEvaluator_IDFieldAccess_BasicComparisons(t *testing.T) {
	t.Log("🧪 TEST: Evaluator - Comparaisons basiques sur champ id")
	t.Log("=======================================================")

	tests := []struct {
		name       string
		fact       *Fact
		expression map[string]interface{}
		varName    string
		wantResult bool
		wantErr    bool
	}{
		{
			name: "Égalité id PK simple",
			fact: &Fact{
				ID:   "Person~Alice",
				Type: "Person",
				Fields: map[string]interface{}{
					"nom": "Alice",
				},
			},
			expression: map[string]interface{}{
				"type": "binary_op",
				"op":   "==",
				"left": map[string]interface{}{
					"type":   "field_access",
					"object": "p",
					"field":  "id",
				},
				"right": map[string]interface{}{
					"type":  "string",
					"value": "Person~Alice",
				},
			},
			varName:    "p",
			wantResult: true,
			wantErr:    false,
		},
		{
			name: "Inégalité id PK simple",
			fact: &Fact{
				ID:   "Person~Alice",
				Type: "Person",
				Fields: map[string]interface{}{
					"nom": "Alice",
				},
			},
			expression: map[string]interface{}{
				"type": "binary_op",
				"op":   "!=",
				"left": map[string]interface{}{
					"type":   "field_access",
					"object": "p",
					"field":  "id",
				},
				"right": map[string]interface{}{
					"type":  "string",
					"value": "Person~Bob",
				},
			},
			varName:    "p",
			wantResult: true,
			wantErr:    false,
		},
		{
			name: "Égalité id PK composite",
			fact: &Fact{
				ID:   "Person~Alice_Dupont",
				Type: "Person",
				Fields: map[string]interface{}{
					"prenom": "Alice",
					"nom":    "Dupont",
				},
			},
			expression: map[string]interface{}{
				"type": "binary_op",
				"op":   "==",
				"left": map[string]interface{}{
					"type":   "field_access",
					"object": "p",
					"field":  "id",
				},
				"right": map[string]interface{}{
					"type":  "string",
					"value": "Person~Alice_Dupont",
				},
			},
			varName:    "p",
			wantResult: true,
			wantErr:    false,
		},
		{
			name: "Égalité id hash",
			fact: &Fact{
				ID:   "Event~a1b2c3d4e5f6g7h8",
				Type: "Event",
				Fields: map[string]interface{}{
					"timestamp": 1234567890,
				},
			},
			expression: map[string]interface{}{
				"type": "binary_op",
				"op":   "==",
				"left": map[string]interface{}{
					"type":   "field_access",
					"object": "e",
					"field":  "id",
				},
				"right": map[string]interface{}{
					"type":  "string",
					"value": "Event~a1b2c3d4e5f6g7h8",
				},
			},
			varName:    "e",
			wantResult: true,
			wantErr:    false,
		},
		{
			name: "CONTAINS sur id",
			fact: &Fact{
				ID:   "Person~Alice_Dupont",
				Type: "Person",
				Fields: map[string]interface{}{
					"prenom": "Alice",
					"nom":    "Dupont",
				},
			},
			expression: map[string]interface{}{
				"type": "binary_op",
				"op":   "CONTAINS",
				"left": map[string]interface{}{
					"type":   "field_access",
					"object": "p",
					"field":  "id",
				},
				"right": map[string]interface{}{
					"type":  "string",
					"value": "Alice",
				},
			},
			varName:    "p",
			wantResult: true,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eval := NewAlphaConditionEvaluator()

			result, err := eval.EvaluateCondition(tt.expression, tt.fact, tt.varName)

			if (err != nil) != tt.wantErr {
				t.Fatalf("❌ EvaluateCondition() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil && result != tt.wantResult {
				t.Errorf("❌ EvaluateCondition() = %v, attendu %v", result, tt.wantResult)
			} else if err == nil {
				t.Log("✅ Test réussi")
			}
		})
	}

	t.Log("✅ Test complet: Comparaisons basiques sur champ id")
}
