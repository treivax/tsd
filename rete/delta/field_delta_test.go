// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
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
