// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import "testing"

func TestChangeType_String(t *testing.T) {
	tests := []struct {
		ct   ChangeType
		want string
	}{
		{ChangeTypeModified, "Modified"},
		{ChangeTypeAdded, "Added"},
		{ChangeTypeRemoved, "Removed"},
		{ChangeType(999), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.ct.String(); got != tt.want {
				t.Errorf("ChangeType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValueType_String(t *testing.T) {
	tests := []struct {
		vt   ValueType
		want string
	}{
		{ValueTypeString, "string"},
		{ValueTypeNumber, "number"},
		{ValueTypeBool, "bool"},
		{ValueTypeObject, "object"},
		{ValueTypeArray, "array"},
		{ValueTypeNull, "null"},
		{ValueTypeUnknown, "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.vt.String(); got != tt.want {
				t.Errorf("ValueType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInferValueType(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		want  ValueType
	}{
		{"nil", nil, ValueTypeNull},
		{"string", "hello", ValueTypeString},
		{"int", 42, ValueTypeNumber},
		{"float64", 3.14, ValueTypeNumber},
		{"bool", true, ValueTypeBool},
		{"object", map[string]interface{}{"a": 1}, ValueTypeObject},
		{"array", []interface{}{1, 2, 3}, ValueTypeArray},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := inferValueType(tt.value); got != tt.want {
				t.Errorf("inferValueType(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}
