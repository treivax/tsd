// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"fmt"
	"testing"
	"time"
)

func TestNewFieldDelta(t *testing.T) {
	tests := []struct {
		name           string
		fieldName      string
		oldValue       interface{}
		newValue       interface{}
		wantChangeType ChangeType
		wantValueType  ValueType
	}{
		{
			name:           "modification string",
			fieldName:      "status",
			oldValue:       "active",
			newValue:       "inactive",
			wantChangeType: ChangeTypeModified,
			wantValueType:  ValueTypeString,
		},
		{
			name:           "ajout champ",
			fieldName:      "email",
			oldValue:       nil,
			newValue:       "user@example.com",
			wantChangeType: ChangeTypeAdded,
			wantValueType:  ValueTypeString,
		},
		{
			name:           "suppression champ",
			fieldName:      "temp_field",
			oldValue:       42,
			newValue:       nil,
			wantChangeType: ChangeTypeRemoved,
			wantValueType:  ValueTypeNumber,
		},
		{
			name:           "modification nombre",
			fieldName:      "price",
			oldValue:       100.0,
			newValue:       150.0,
			wantChangeType: ChangeTypeModified,
			wantValueType:  ValueTypeNumber,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delta := NewFieldDelta(tt.fieldName, tt.oldValue, tt.newValue)

			if delta.FieldName != tt.fieldName {
				t.Errorf("FieldName = %v, want %v", delta.FieldName, tt.fieldName)
			}
			if delta.ChangeType != tt.wantChangeType {
				t.Errorf("ChangeType = %v, want %v", delta.ChangeType, tt.wantChangeType)
			}
			if delta.ValueType != tt.wantValueType {
				t.Errorf("ValueType = %v, want %v", delta.ValueType, tt.wantValueType)
			}
		})
	}
}

func TestFieldDelta_String(t *testing.T) {
	tests := []struct {
		name      string
		delta     FieldDelta
		wantMatch string
	}{
		{
			name:      "modification",
			delta:     NewFieldDelta("price", 100, 150),
			wantMatch: "price: (100 → 150)",
		},
		{
			name:      "ajout",
			delta:     NewFieldDelta("email", nil, "test@example.com"),
			wantMatch: "email: (nil → test@example.com)",
		},
		{
			name:      "suppression",
			delta:     NewFieldDelta("temp", 42, nil),
			wantMatch: "temp: (42 → nil)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.delta.String()
			if got != tt.wantMatch {
				t.Errorf("String() = %v, want %v", got, tt.wantMatch)
			}
		})
	}
}

func TestNewFactDelta(t *testing.T) {
	before := time.Now()
	delta := NewFactDelta("Product~123", "Product")
	after := time.Now()

	if delta.FactID != "Product~123" {
		t.Errorf("FactID = %v, want Product~123", delta.FactID)
	}
	if delta.FactType != "Product" {
		t.Errorf("FactType = %v, want Product", delta.FactType)
	}
	if delta.Fields == nil {
		t.Error("Fields map should be initialized")
	}
	if delta.Timestamp.Before(before) || delta.Timestamp.After(after) {
		t.Errorf("Timestamp %v not in expected range", delta.Timestamp)
	}
}

func TestFactDelta_AddFieldChange(t *testing.T) {
	delta := NewFactDelta("Product~123", "Product")

	delta.AddFieldChange("price", 100.0, 150.0)
	delta.AddFieldChange("status", "active", "inactive")

	if len(delta.Fields) != 2 {
		t.Errorf("Expected 2 fields, got %d", len(delta.Fields))
	}

	priceChange, ok := delta.Fields["price"]
	if !ok {
		t.Fatal("price field not found")
	}
	if priceChange.OldValue != 100.0 || priceChange.NewValue != 150.0 {
		t.Errorf("price change incorrect: %v → %v", priceChange.OldValue, priceChange.NewValue)
	}
}

func TestFactDelta_IsEmpty(t *testing.T) {
	delta := NewFactDelta("Product~123", "Product")

	if !delta.IsEmpty() {
		t.Error("New delta should be empty")
	}

	delta.AddFieldChange("price", 100.0, 150.0)

	if delta.IsEmpty() {
		t.Error("Delta with changes should not be empty")
	}
}

func TestFactDelta_FieldsChanged(t *testing.T) {
	delta := NewFactDelta("Product~123", "Product")
	delta.AddFieldChange("price", 100.0, 150.0)
	delta.AddFieldChange("status", "active", "inactive")

	fields := delta.FieldsChanged()

	if len(fields) != 2 {
		t.Errorf("Expected 2 changed fields, got %d", len(fields))
	}

	// Vérifier que les champs sont présents (ordre non garanti)
	found := make(map[string]bool)
	for _, field := range fields {
		found[field] = true
	}
	if !found["price"] || !found["status"] {
		t.Errorf("Expected fields [price, status], got %v", fields)
	}
}

func TestFactDelta_ChangeRatio(t *testing.T) {
	tests := []struct {
		name          string
		fieldsChanged int
		totalFields   int
		wantRatio     float64
	}{
		{"2 sur 10", 2, 10, 0.2},
		{"1 sur 5", 1, 5, 0.2},
		{"5 sur 10", 5, 10, 0.5},
		{"0 sur 10", 0, 10, 0.0},
		{"1 sur 0", 1, 0, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delta := NewFactDelta("Test~1", "Test")
			delta.FieldCount = tt.totalFields

			// Ajouter tt.fieldsChanged champs
			for i := 0; i < tt.fieldsChanged; i++ {
				delta.AddFieldChange(string(rune('a'+i)), i, i+1)
			}

			ratio := delta.ChangeRatio()
			if ratio != tt.wantRatio {
				t.Errorf("ChangeRatio() = %v, want %v", ratio, tt.wantRatio)
			}
		})
	}
}

func TestFactDelta_String(t *testing.T) {
	t.Run("empty delta", func(t *testing.T) {
		delta := NewFactDelta("Product~123", "Product")
		str := delta.String()
		if str != "FactDelta[Product:Product~123] (no changes)" {
			t.Errorf("String() = %v", str)
		}
	})

	t.Run("delta with changes", func(t *testing.T) {
		delta := NewFactDelta("Product~123", "Product")
		delta.AddFieldChange("price", 100, 150)

		str := delta.String()
		if str == "" {
			t.Error("String() should not be empty")
		}
		if len(str) < 10 {
			t.Errorf("String() too short: %v", str)
		}
	})
}

// Test avec valeurs nil
func TestFieldDelta_NilValues(t *testing.T) {
	tests := []struct {
		name     string
		oldValue interface{}
		newValue interface{}
		wantType ChangeType
	}{
		{"nil → value", nil, "test", ChangeTypeAdded},
		{"value → nil", "test", nil, ChangeTypeRemoved},
		{"nil → nil", nil, nil, ChangeTypeModified},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delta := NewFieldDelta("field", tt.oldValue, tt.newValue)
			if delta.ChangeType != tt.wantType {
				t.Errorf("❌ Expected %v, got %v", tt.wantType, delta.ChangeType)
			}
		})
	}
}

// Test avec unicode et caractères spéciaux
func TestFieldDelta_UnicodeFields(t *testing.T) {
	tests := []struct {
		name      string
		fieldName string
		oldValue  interface{}
		newValue  interface{}
	}{
		{"japonais", "名前", "古い値", "新しい値"},
		{"emoji", "emoji", "😀", "😎"},
		{"français", "texte", "Voilà un café", "Zürich est belle"},
		{"cyrillique", "текст", "Привет", "Пока"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delta := NewFieldDelta(tt.fieldName, tt.oldValue, tt.newValue)

			if delta.FieldName != tt.fieldName {
				t.Errorf("❌ FieldName = %v, want %v", delta.FieldName, tt.fieldName)
			}

			if delta.OldValue != tt.oldValue {
				t.Errorf("❌ OldValue = %v, want %v", delta.OldValue, tt.oldValue)
			}

			if delta.NewValue != tt.newValue {
				t.Errorf("❌ NewValue = %v, want %v", delta.NewValue, tt.newValue)
			}
		})
	}
}

// Test changement de type
func TestFieldDelta_TypeChange(t *testing.T) {
	tests := []struct {
		name     string
		oldValue interface{}
		newValue interface{}
		oldType  ValueType
		newType  ValueType
	}{
		{"number → string", 42, "42", ValueTypeNumber, ValueTypeString},
		{"string → number", "100", 100.0, ValueTypeString, ValueTypeNumber},
		{"bool → string", true, "true", ValueTypeBool, ValueTypeString},
		{"array → object", []interface{}{1, 2}, map[string]interface{}{"a": 1}, ValueTypeArray, ValueTypeObject},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delta := NewFieldDelta("field", tt.oldValue, tt.newValue)

			// Vérifier que le changement est détecté comme modification
			if delta.ChangeType != ChangeTypeModified {
				t.Errorf("❌ Expected ChangeTypeModified, got %v", delta.ChangeType)
			}

			// La valeur peut être soit l'ancien soit le nouveau type selon l'implémentation
			// On vérifie juste que c'est cohérent
			if delta.ValueType != tt.oldType && delta.ValueType != tt.newType {
				t.Errorf("❌ ValueType = %v, expected %v or %v", delta.ValueType, tt.oldType, tt.newType)
			}
		})
	}
}

// Test avec structures complexes
func TestFieldDelta_ComplexStructures(t *testing.T) {
	t.Run("nested objects", func(t *testing.T) {
		oldValue := map[string]interface{}{
			"level1": map[string]interface{}{
				"level2": "value",
			},
		}
		newValue := map[string]interface{}{
			"level1": map[string]interface{}{
				"level2": "changed",
			},
		}

		delta := NewFieldDelta("nested", oldValue, newValue)
		if delta.ChangeType != ChangeTypeModified {
			t.Errorf("❌ Expected modification detected")
		}
	})

	t.Run("arrays with different content", func(t *testing.T) {
		oldValue := []interface{}{1, 2, 3}
		newValue := []interface{}{1, 2, 4}

		delta := NewFieldDelta("array", oldValue, newValue)
		if delta.ChangeType != ChangeTypeModified {
			t.Errorf("❌ Expected modification detected")
		}
	})

	t.Run("empty arrays", func(t *testing.T) {
		oldValue := []interface{}{}
		newValue := []interface{}{1}

		delta := NewFieldDelta("array", oldValue, newValue)
		if delta.ChangeType != ChangeTypeModified {
			t.Errorf("❌ Expected modification detected")
		}
	})
}

// Test FactDelta avec beaucoup de champs
func TestFactDelta_ManyFields(t *testing.T) {
	delta := NewFactDelta("Test~1", "Test")
	delta.FieldCount = TestFieldCount

	// Ajouter 50 changements
	const changedFields = 50
	for i := 0; i < changedFields; i++ {
		fieldName := fmt.Sprintf("field_%d", i)
		delta.AddFieldChange(fieldName, i, i+1)
	}

	if len(delta.Fields) != changedFields {
		t.Errorf("❌ Expected %d fields, got %d", changedFields, len(delta.Fields))
	}

	ratio := delta.ChangeRatio()
	expectedRatio := float64(changedFields) / float64(TestFieldCount)
	if ratio != expectedRatio {
		t.Errorf("❌ Expected ratio %v, got %v", expectedRatio, ratio)
	}
}

// Test FieldsChanged avec ordre
func TestFactDelta_FieldsChangedConsistency(t *testing.T) {
	delta := NewFactDelta("Test~1", "Test")
	delta.AddFieldChange("alpha", 1, 2)
	delta.AddFieldChange("beta", 3, 4)
	delta.AddFieldChange("gamma", 5, 6)

	// Appeler plusieurs fois pour vérifier cohérence
	fields1 := delta.FieldsChanged()
	fields2 := delta.FieldsChanged()

	if len(fields1) != len(fields2) {
		t.Errorf("❌ FieldsChanged() should return consistent results")
	}

	// Vérifier que les mêmes champs sont présents
	map1 := make(map[string]bool)
	for _, f := range fields1 {
		map1[f] = true
	}

	for _, f := range fields2 {
		if !map1[f] {
			t.Errorf("❌ Field %s not consistent", f)
		}
	}
}
