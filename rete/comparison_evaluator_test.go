// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"testing"
)

func TestComparisonEvaluator_CompareFactIDs(t *testing.T) {
	t.Log("🧪 TEST COMPARISON EVALUATOR - IDS DE FAITS")
	t.Log("============================================")

	evaluator := NewComparisonEvaluator(nil)

	tests := []struct {
		name     string
		left     interface{}
		right    interface{}
		operator string
		expected bool
		wantErr  bool
	}{
		{
			name:     "IDs égaux avec ==",
			left:     "User~Alice",
			right:    "User~Alice",
			operator: "==",
			expected: true,
			wantErr:  false,
		},
		{
			name:     "IDs différents avec ==",
			left:     "User~Alice",
			right:    "User~Bob",
			operator: "==",
			expected: false,
			wantErr:  false,
		},
		{
			name:     "IDs différents avec !=",
			left:     "User~Alice",
			right:    "User~Bob",
			operator: "!=",
			expected: true,
			wantErr:  false,
		},
		{
			name:     "IDs égaux avec !=",
			left:     "User~Alice",
			right:    "User~Alice",
			operator: "!=",
			expected: false,
			wantErr:  false,
		},
		{
			name:     "opérateur < non supporté pour faits",
			left:     "User~Alice",
			right:    "User~Bob",
			operator: "<",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evaluator.compareFactIDs(tt.left, tt.right, tt.operator)

			if tt.wantErr {
				if err == nil {
					t.Errorf("❌ Attendu une erreur, reçu nil")
				} else {
					t.Logf("✅ Erreur attendue: %v", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("❌ Erreur inattendue: %v", err)
			}

			if result != tt.expected {
				t.Errorf("❌ Résultat attendu %v, reçu %v", tt.expected, result)
			} else {
				t.Logf("✅ Comparaison correcte: %v %s %v = %v", tt.left, tt.operator, tt.right, result)
			}
		})
	}
}

func TestComparisonEvaluator_ComparePrimitives(t *testing.T) {
	t.Log("🧪 TEST COMPARISON EVALUATOR - PRIMITIFS")
	t.Log("=========================================")

	evaluator := NewComparisonEvaluator(nil)

	tests := []struct {
		name     string
		left     interface{}
		right    interface{}
		operator string
		expected bool
		wantErr  bool
	}{
		// Strings
		{
			name:     "strings égaux",
			left:     "alice",
			right:    "alice",
			operator: "==",
			expected: true,
		},
		{
			name:     "strings différents",
			left:     "alice",
			right:    "bob",
			operator: "!=",
			expected: true,
		},
		{
			name:     "string < string",
			left:     "alice",
			right:    "bob",
			operator: "<",
			expected: true,
		},

		// Numbers
		{
			name:     "numbers égaux",
			left:     42.0,
			right:    42.0,
			operator: "==",
			expected: true,
		},
		{
			name:     "numbers différents",
			left:     10.0,
			right:    20.0,
			operator: "<",
			expected: true,
		},
		{
			name:     "int et float64",
			left:     int(42),
			right:    42.0,
			operator: "==",
			expected: true,
		},

		// Booleans
		{
			name:     "booleans égaux",
			left:     true,
			right:    true,
			operator: "==",
			expected: true,
		},
		{
			name:     "booleans différents",
			left:     true,
			right:    false,
			operator: "!=",
			expected: true,
		},
		{
			name:     "boolean < interdit",
			left:     true,
			right:    false,
			operator: "<",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evaluator.comparePrimitives(tt.left, tt.right, tt.operator)

			if tt.wantErr {
				if err == nil {
					t.Errorf("❌ Attendu une erreur, reçu nil")
				} else {
					t.Logf("✅ Erreur attendue: %v", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("❌ Erreur inattendue: %v", err)
			}

			if result != tt.expected {
				t.Errorf("❌ Résultat attendu %v, reçu %v", tt.expected, result)
			} else {
				t.Logf("✅ Comparaison correcte: %v %s %v = %v", tt.left, tt.operator, tt.right, result)
			}
		})
	}
}

func TestComparisonEvaluator_EvaluateComparison(t *testing.T) {
	t.Log("🧪 TEST COMPARISON EVALUATOR - GLOBAL")
	t.Log("======================================")

	evaluator := NewComparisonEvaluator(nil)

	tests := []struct {
		name      string
		left      interface{}
		right     interface{}
		operator  string
		leftType  string
		rightType string
		expected  bool
		wantErr   bool
	}{
		{
			name:      "comparaison de faits",
			left:      "User~Alice",
			right:     "User~Alice",
			operator:  "==",
			leftType:  FieldTypeFact,
			rightType: FieldTypeFact,
			expected:  true,
		},
		{
			name:      "comparaison de primitifs",
			left:      "alice",
			right:     "alice",
			operator:  "==",
			leftType:  FieldTypePrimitive,
			rightType: FieldTypePrimitive,
			expected:  true,
		},
		{
			name:      "types incompatibles",
			left:      "User~Alice",
			right:     "alice",
			operator:  "==",
			leftType:  FieldTypeFact,
			rightType: FieldTypePrimitive,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evaluator.EvaluateComparison(tt.left, tt.right, tt.operator, tt.leftType, tt.rightType)

			if tt.wantErr {
				if err == nil {
					t.Errorf("❌ Attendu une erreur, reçu nil")
				} else {
					t.Logf("✅ Erreur attendue: %v", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("❌ Erreur inattendue: %v", err)
			}

			if result != tt.expected {
				t.Errorf("❌ Résultat attendu %v, reçu %v", tt.expected, result)
			} else {
				t.Logf("✅ Évaluation correcte")
			}
		})
	}
}

func TestToFloat64(t *testing.T) {
	t.Log("🧪 TEST HELPER - convertToFloat64")
	t.Log("==================================")

	tests := []struct {
		name     string
		input    interface{}
		expected float64
		ok       bool
	}{
		{"float64", float64(42.5), 42.5, true},
		{"int", int(42), 42.0, true},
		{"int32", int32(42), 42.0, true},
		{"int64", int64(42), 42.0, true},
		{"float32", float32(42.5), 42.5, true},
		{"string", "42", 42.0, true},
		{"bool", true, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := convertToFloat64(tt.input)

			if ok != tt.ok {
				t.Errorf("❌ ok attendu %v, reçu %v", tt.ok, ok)
				return
			}

			if !ok {
				t.Logf("✅ Conversion échouée comme attendu pour %T", tt.input)
				return
			}

			if result != tt.expected {
				t.Errorf("❌ Valeur attendue %v, reçu %v", tt.expected, result)
			} else {
				t.Logf("✅ Conversion réussie: %v -> %v", tt.input, result)
			}
		})
	}
}
