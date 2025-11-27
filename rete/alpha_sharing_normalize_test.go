// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"reflect"
	"testing"
)

// TestNormalizeConditionForSharing_Unwrap teste le déballe des wrappers "constraint"
func TestNormalizeConditionForSharing_Unwrap(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name: "Unwrap single constraint wrapper",
			input: map[string]interface{}{
				"type": "constraint",
				"constraint": map[string]interface{}{
					"type":     "binaryOperation",
					"operator": ">",
					"left":     map[string]interface{}{"type": "field", "name": "age"},
					"right":    map[string]interface{}{"type": "literal", "value": 18.0},
				},
			},
			expected: map[string]interface{}{
				"type":     "binaryOperation",
				"operator": ">",
				"left":     map[string]interface{}{"type": "field", "name": "age"},
				"right":    map[string]interface{}{"type": "literal", "value": 18.0},
			},
		},
		{
			name: "Unwrap nested constraint wrappers",
			input: map[string]interface{}{
				"type": "constraint",
				"constraint": map[string]interface{}{
					"type": "constraint",
					"constraint": map[string]interface{}{
						"type":     "binaryOperation",
						"operator": "==",
						"left":     map[string]interface{}{"type": "field", "name": "name"},
						"right":    map[string]interface{}{"type": "literal", "value": "toto"},
					},
				},
			},
			expected: map[string]interface{}{
				"type":     "binaryOperation",
				"operator": "==",
				"left":     map[string]interface{}{"type": "field", "name": "name"},
				"right":    map[string]interface{}{"type": "literal", "value": "toto"},
			},
		},
		{
			name: "No wrapper - return as is",
			input: map[string]interface{}{
				"type":     "binaryOperation",
				"operator": ">",
				"left":     map[string]interface{}{"type": "field", "name": "salary"},
				"right":    map[string]interface{}{"type": "literal", "value": 1000.0},
			},
			expected: map[string]interface{}{
				"type":     "binaryOperation",
				"operator": ">",
				"left":     map[string]interface{}{"type": "field", "name": "salary"},
				"right":    map[string]interface{}{"type": "literal", "value": 1000.0},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeConditionForSharing(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("normalizeConditionForSharing() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestNormalizeConditionForSharing_TypeNormalization teste la normalisation des types
func TestNormalizeConditionForSharing_TypeNormalization(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name: "Normalize comparison to binaryOperation",
			input: map[string]interface{}{
				"type":     "comparison",
				"operator": ">",
				"left":     map[string]interface{}{"type": "field", "name": "age"},
				"right":    map[string]interface{}{"type": "literal", "value": 18.0},
			},
			expected: map[string]interface{}{
				"type":     "binaryOperation",
				"operator": ">",
				"left":     map[string]interface{}{"type": "field", "name": "age"},
				"right":    map[string]interface{}{"type": "literal", "value": 18.0},
			},
		},
		{
			name: "binaryOperation stays binaryOperation",
			input: map[string]interface{}{
				"type":     "binaryOperation",
				"operator": "==",
				"left":     map[string]interface{}{"type": "field", "name": "status"},
				"right":    map[string]interface{}{"type": "literal", "value": "VIP"},
			},
			expected: map[string]interface{}{
				"type":     "binaryOperation",
				"operator": "==",
				"left":     map[string]interface{}{"type": "field", "name": "status"},
				"right":    map[string]interface{}{"type": "literal", "value": "VIP"},
			},
		},
		{
			name: "Nested comparison normalization",
			input: map[string]interface{}{
				"type":     "comparison",
				"operator": ">",
				"left": map[string]interface{}{
					"type":     "comparison",
					"operator": "+",
					"left":     map[string]interface{}{"type": "field", "name": "x"},
					"right":    map[string]interface{}{"type": "field", "name": "y"},
				},
				"right": map[string]interface{}{"type": "literal", "value": 100.0},
			},
			expected: map[string]interface{}{
				"type":     "binaryOperation",
				"operator": ">",
				"left": map[string]interface{}{
					"type":     "binaryOperation",
					"operator": "+",
					"left":     map[string]interface{}{"type": "field", "name": "x"},
					"right":    map[string]interface{}{"type": "field", "name": "y"},
				},
				"right": map[string]interface{}{"type": "literal", "value": 100.0},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeConditionForSharing(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("normalizeConditionForSharing() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestNormalizeConditionForSharing_Combined teste les cas combinant unwrap et normalisation
func TestNormalizeConditionForSharing_Combined(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name: "Unwrap constraint wrapper AND normalize comparison to binaryOperation",
			input: map[string]interface{}{
				"type": "constraint",
				"constraint": map[string]interface{}{
					"type":     "comparison",
					"operator": ">",
					"left":     map[string]interface{}{"type": "field", "name": "age"},
					"right":    map[string]interface{}{"type": "literal", "value": 18.0},
				},
			},
			expected: map[string]interface{}{
				"type":     "binaryOperation",
				"operator": ">",
				"left":     map[string]interface{}{"type": "field", "name": "age"},
				"right":    map[string]interface{}{"type": "literal", "value": 18.0},
			},
		},
		{
			name: "Multiple wrappers with comparison normalization",
			input: map[string]interface{}{
				"type": "constraint",
				"constraint": map[string]interface{}{
					"type": "constraint",
					"constraint": map[string]interface{}{
						"type":     "comparison",
						"operator": "==",
						"left":     map[string]interface{}{"type": "field", "name": "name"},
						"right":    map[string]interface{}{"type": "literal", "value": "toto"},
					},
				},
			},
			expected: map[string]interface{}{
				"type":     "binaryOperation",
				"operator": "==",
				"left":     map[string]interface{}{"type": "field", "name": "name"},
				"right":    map[string]interface{}{"type": "literal", "value": "toto"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeConditionForSharing(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("normalizeConditionForSharing() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestNormalizeConditionForSharing_Slices teste la normalisation des slices
func TestNormalizeConditionForSharing_Slices(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name: "Normalize slice of conditions",
			input: []interface{}{
				map[string]interface{}{
					"type":     "comparison",
					"operator": ">",
					"left":     map[string]interface{}{"type": "field", "name": "age"},
					"right":    map[string]interface{}{"type": "literal", "value": 18.0},
				},
				map[string]interface{}{
					"type":     "comparison",
					"operator": "==",
					"left":     map[string]interface{}{"type": "field", "name": "status"},
					"right":    map[string]interface{}{"type": "literal", "value": "active"},
				},
			},
			expected: []interface{}{
				map[string]interface{}{
					"type":     "binaryOperation",
					"operator": ">",
					"left":     map[string]interface{}{"type": "field", "name": "age"},
					"right":    map[string]interface{}{"type": "literal", "value": 18.0},
				},
				map[string]interface{}{
					"type":     "binaryOperation",
					"operator": "==",
					"left":     map[string]interface{}{"type": "field", "name": "status"},
					"right":    map[string]interface{}{"type": "literal", "value": "active"},
				},
			},
		},
		{
			name: "Normalize slice with wrapped conditions",
			input: []interface{}{
				map[string]interface{}{
					"type": "constraint",
					"constraint": map[string]interface{}{
						"type":     "comparison",
						"operator": ">",
						"left":     map[string]interface{}{"type": "field", "name": "x"},
						"right":    map[string]interface{}{"type": "literal", "value": 0.0},
					},
				},
			},
			expected: []interface{}{
				map[string]interface{}{
					"type":     "binaryOperation",
					"operator": ">",
					"left":     map[string]interface{}{"type": "field", "name": "x"},
					"right":    map[string]interface{}{"type": "literal", "value": 0.0},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeConditionForSharing(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("normalizeConditionForSharing() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestNormalizeConditionForSharing_Primitives teste les types primitifs
func TestNormalizeConditionForSharing_Primitives(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name:     "String primitive",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "Number primitive",
			input:    42.0,
			expected: 42.0,
		},
		{
			name:     "Boolean primitive",
			input:    true,
			expected: true,
		},
		{
			name:     "Nil value",
			input:    nil,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeConditionForSharing(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("normalizeConditionForSharing() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestNormalizeConditionForSharing_ComplexNested teste les structures complexes imbriquées
func TestNormalizeConditionForSharing_ComplexNested(t *testing.T) {
	input := map[string]interface{}{
		"type": "constraint",
		"constraint": map[string]interface{}{
			"type":     "comparison",
			"operator": "AND",
			"left": map[string]interface{}{
				"type":     "comparison",
				"operator": ">",
				"left":     map[string]interface{}{"type": "field", "name": "age"},
				"right":    map[string]interface{}{"type": "literal", "value": 18.0},
			},
			"right": map[string]interface{}{
				"type": "constraint",
				"constraint": map[string]interface{}{
					"type":     "comparison",
					"operator": "==",
					"left":     map[string]interface{}{"type": "field", "name": "status"},
					"right":    map[string]interface{}{"type": "literal", "value": "VIP"},
				},
			},
		},
	}

	expected := map[string]interface{}{
		"type":     "binaryOperation",
		"operator": "AND",
		"left": map[string]interface{}{
			"type":     "binaryOperation",
			"operator": ">",
			"left":     map[string]interface{}{"type": "field", "name": "age"},
			"right":    map[string]interface{}{"type": "literal", "value": 18.0},
		},
		"right": map[string]interface{}{
			"type":     "binaryOperation",
			"operator": "==",
			"left":     map[string]interface{}{"type": "field", "name": "status"},
			"right":    map[string]interface{}{"type": "literal", "value": "VIP"},
		},
	}

	result := normalizeConditionForSharing(input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("normalizeConditionForSharing() with complex nested structure failed")
		t.Errorf("Got: %+v", result)
		t.Errorf("Want: %+v", expected)
	}
}

// TestNormalizeConditionForSharing_RealWorldScenarios teste des scénarios réels
func TestNormalizeConditionForSharing_RealWorldScenarios(t *testing.T) {
	t.Run("Simple rule condition (wrapped)", func(t *testing.T) {
		// Condition d'une règle simple (comme r1 ou r4 dans les tests d'intégration)
		input := map[string]interface{}{
			"type": "constraint",
			"constraint": map[string]interface{}{
				"type":     "comparison",
				"operator": ">",
				"left":     map[string]interface{}{"type": "field", "name": "age"},
				"right":    map[string]interface{}{"type": "literal", "value": 18.0},
			},
		}

		expected := map[string]interface{}{
			"type":     "binaryOperation",
			"operator": ">",
			"left":     map[string]interface{}{"type": "field", "name": "age"},
			"right":    map[string]interface{}{"type": "literal", "value": 18.0},
		}

		result := normalizeConditionForSharing(input)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Simple rule normalization failed")
		}
	})

	t.Run("Chain condition (unwrapped)", func(t *testing.T) {
		// Condition d'une chaîne (déjà décomposée, pas de wrapper)
		input := map[string]interface{}{
			"type":     "binaryOperation",
			"operator": ">",
			"left":     map[string]interface{}{"type": "field", "name": "age"},
			"right":    map[string]interface{}{"type": "literal", "value": 18.0},
		}

		expected := map[string]interface{}{
			"type":     "binaryOperation",
			"operator": ">",
			"left":     map[string]interface{}{"type": "field", "name": "age"},
			"right":    map[string]interface{}{"type": "literal", "value": 18.0},
		}

		result := normalizeConditionForSharing(input)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Chain condition normalization failed")
		}
	})

	t.Run("Both should produce same normalized form", func(t *testing.T) {
		// Une condition de règle simple (wrapped + comparison)
		simpleRuleCondition := map[string]interface{}{
			"type": "constraint",
			"constraint": map[string]interface{}{
				"type":     "comparison",
				"operator": ">",
				"left":     map[string]interface{}{"type": "field", "name": "age"},
				"right":    map[string]interface{}{"type": "literal", "value": 18.0},
			},
		}

		// La même condition dans une chaîne (unwrapped + binaryOperation)
		chainCondition := map[string]interface{}{
			"type":     "binaryOperation",
			"operator": ">",
			"left":     map[string]interface{}{"type": "field", "name": "age"},
			"right":    map[string]interface{}{"type": "literal", "value": 18.0},
		}

		result1 := normalizeConditionForSharing(simpleRuleCondition)
		result2 := normalizeConditionForSharing(chainCondition)

		if !reflect.DeepEqual(result1, result2) {
			t.Errorf("Simple rule and chain conditions should normalize to same form")
			t.Errorf("Simple rule result: %+v", result1)
			t.Errorf("Chain result: %+v", result2)
		}
	})
}

// TestNormalizeConditionForSharing_Idempotence teste l'idempotence de la normalisation
func TestNormalizeConditionForSharing_Idempotence(t *testing.T) {
	input := map[string]interface{}{
		"type": "constraint",
		"constraint": map[string]interface{}{
			"type":     "comparison",
			"operator": ">",
			"left":     map[string]interface{}{"type": "field", "name": "age"},
			"right":    map[string]interface{}{"type": "literal", "value": 18.0},
		},
	}

	// Normaliser une première fois
	result1 := normalizeConditionForSharing(input)

	// Normaliser le résultat (devrait être idempotent)
	result2 := normalizeConditionForSharing(result1)

	if !reflect.DeepEqual(result1, result2) {
		t.Errorf("normalizeConditionForSharing() should be idempotent")
		t.Errorf("First pass: %+v", result1)
		t.Errorf("Second pass: %+v", result2)
	}
}

// TestNormalizeConditionForSharing_EdgeCases teste les cas limites
func TestNormalizeConditionForSharing_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name:     "Empty map",
			input:    map[string]interface{}{},
			expected: map[string]interface{}{},
		},
		{
			name:     "Empty slice",
			input:    []interface{}{},
			expected: []interface{}{},
		},
		{
			name: "Map with only type field",
			input: map[string]interface{}{
				"type": "comparison",
			},
			expected: map[string]interface{}{
				"type": "binaryOperation",
			},
		},
		{
			name: "Constraint wrapper with no constraint field",
			input: map[string]interface{}{
				"type": "constraint",
			},
			expected: map[string]interface{}{
				"type": "constraint",
			},
		},
		{
			name: "Constraint wrapper with nil constraint",
			input: map[string]interface{}{
				"type":       "constraint",
				"constraint": nil,
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeConditionForSharing(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("normalizeConditionForSharing() = %v, want %v", result, tt.expected)
			}
		})
	}
}
