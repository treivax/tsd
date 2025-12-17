// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"testing"
)

func TestFactValue_Unwrap(t *testing.T) {
	t.Log("🧪 TEST FACTVALUE UNWRAP")
	t.Log("=========================")

	tests := []struct {
		name     string
		factVal  FactValue
		expected interface{}
	}{
		{
			name: "valeur simple string",
			factVal: FactValue{
				Type:  "string",
				Value: "test",
			},
			expected: "test",
		},
		{
			name: "valeur simple number",
			factVal: FactValue{
				Type:  "number",
				Value: 42.0,
			},
			expected: 42.0,
		},
		{
			name: "valeur simple bool",
			factVal: FactValue{
				Type:  "bool",
				Value: true,
			},
			expected: true,
		},
		{
			name: "valeur wrappée dans map",
			factVal: FactValue{
				Type: "string",
				Value: map[string]interface{}{
					"value": "wrapped",
				},
			},
			expected: "wrapped",
		},
		{
			name: "valeur wrappée avec number",
			factVal: FactValue{
				Type: "number",
				Value: map[string]interface{}{
					"value": 123.45,
				},
			},
			expected: 123.45,
		},
		{
			name: "map sans clé value",
			factVal: FactValue{
				Type: "string",
				Value: map[string]interface{}{
					"other": "data",
				},
			},
			expected: map[string]interface{}{
				"other": "data",
			},
		},
		{
			name: "valeur nil",
			factVal: FactValue{
				Type:  "string",
				Value: nil,
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.factVal.Unwrap()

			// Compare values
			switch exp := tt.expected.(type) {
			case string:
				if str, ok := result.(string); !ok || str != exp {
					t.Errorf("❌ Attendu string '%s', reçu %T %v", exp, result, result)
				} else {
					t.Logf("✅ Unwrap: %v → '%s'", tt.factVal.Value, str)
				}
			case float64:
				if num, ok := result.(float64); !ok || num != exp {
					t.Errorf("❌ Attendu number %f, reçu %T %v", exp, result, result)
				} else {
					t.Logf("✅ Unwrap: %v → %f", tt.factVal.Value, num)
				}
			case bool:
				if b, ok := result.(bool); !ok || b != exp {
					t.Errorf("❌ Attendu bool %v, reçu %T %v", exp, result, result)
				} else {
					t.Logf("✅ Unwrap: %v → %v", tt.factVal.Value, b)
				}
			case nil:
				if result != nil {
					t.Errorf("❌ Attendu nil, reçu %T %v", result, result)
				} else {
					t.Log("✅ Unwrap: nil → nil")
				}
			case map[string]interface{}:
				if resMap, ok := result.(map[string]interface{}); !ok {
					t.Errorf("❌ Attendu map, reçu %T %v", result, result)
				} else {
					// Compare map contents
					if len(resMap) != len(exp) {
						t.Errorf("❌ Taille map différente: attendu %d, reçu %d", len(exp), len(resMap))
					}
					t.Logf("✅ Unwrap: %v → map[%d]", tt.factVal.Value, len(resMap))
				}
			}
		})
	}
}

func TestFact_BuildFieldMap(t *testing.T) {
	t.Log("🧪 TEST FACT BUILD FIELD MAP")
	t.Log("=============================")

	tests := []struct {
		name         string
		fact         Fact
		expectedSize int
		checkField   string
		checkValue   interface{}
	}{
		{
			name: "fact avec plusieurs champs",
			fact: Fact{
				TypeName: "Person",
				Fields: []FactField{
					{Name: "nom", Value: FactValue{Type: "string", Value: "Alice"}},
					{Name: "age", Value: FactValue{Type: "number", Value: 30.0}},
					{Name: "actif", Value: FactValue{Type: "bool", Value: true}},
				},
			},
			expectedSize: 3,
			checkField:   "nom",
			checkValue:   "Alice",
		},
		{
			name: "fact vide",
			fact: Fact{
				TypeName: "Empty",
				Fields:   []FactField{},
			},
			expectedSize: 0,
		},
		{
			name: "fact avec un seul champ",
			fact: Fact{
				TypeName: "Single",
				Fields: []FactField{
					{Name: "id", Value: FactValue{Type: "number", Value: 42.0}},
				},
			},
			expectedSize: 1,
			checkField:   "id",
			checkValue:   42.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fieldMap := tt.fact.BuildFieldMap()

			if len(fieldMap) != tt.expectedSize {
				t.Errorf("❌ Taille attendue %d, reçu %d", tt.expectedSize, len(fieldMap))
			}

			if tt.checkField != "" {
				fieldVal, exists := fieldMap[tt.checkField]
				if !exists {
					t.Errorf("❌ Champ '%s' non trouvé dans la map", tt.checkField)
				} else if fieldVal.Value != tt.checkValue {
					t.Errorf("❌ Valeur du champ '%s': attendu %v, reçu %v",
						tt.checkField, tt.checkValue, fieldVal.Value)
				} else {
					t.Logf("✅ Map construite: %d champs, '%s' = %v",
						len(fieldMap), tt.checkField, tt.checkValue)
				}
			} else {
				t.Logf("✅ Map construite: %d champs", len(fieldMap))
			}
		})
	}
}

func TestConvertFactFieldValue_Deprecated(t *testing.T) {
	t.Log("🧪 TEST CONVERT FACT FIELD VALUE (DEPRECATED)")
	t.Log("==============================================")

	// Test que la fonction deprecated fonctionne toujours
	val := FactValue{
		Type:  "string",
		Value: "test",
	}

	result := convertFactFieldValue(val)
	if result != "test" {
		t.Errorf("❌ Fonction deprecated cassée: attendu 'test', reçu %v", result)
	} else {
		t.Log("✅ Fonction deprecated fonctionne (backward compatibility)")
	}

	// Test avec valeur wrappée
	valWrapped := FactValue{
		Type: "number",
		Value: map[string]interface{}{
			"value": 42.0,
		},
	}

	resultWrapped := convertFactFieldValue(valWrapped)
	if resultWrapped != 42.0 {
		t.Errorf("❌ Fonction deprecated avec wrap cassée: attendu 42.0, reçu %v", resultWrapped)
	} else {
		t.Log("✅ Fonction deprecated avec wrap fonctionne")
	}
}
