// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import "testing"

func TestValuesEqual(t *testing.T) {
	tests := []struct {
		name    string
		a       interface{}
		b       interface{}
		epsilon float64
		want    bool
	}{
		// Cas simples
		{"int égaux", 42, 42, 0, true},
		{"int différents", 42, 43, 0, false},
		{"string égaux", "hello", "hello", 0, true},
		{"string différents", "hello", "world", 0, false},
		{"bool égaux", true, true, 0, true},
		{"bool différents", true, false, 0, false},

		// nil
		{"nil == nil", nil, nil, 0, true},
		{"nil != value", nil, 42, 0, false},
		{"value != nil", 42, nil, 0, false},

		// Floats avec epsilon
		{"float64 égaux", 1.0, 1.0, DefaultFloatEpsilon, true},
		{"float64 proches", 1.0, 1.0000000001, DefaultFloatEpsilon, true},
		{"float64 différents", 1.0, 1.1, DefaultFloatEpsilon, false},
		{"float32 égaux", float32(1.0), float32(1.0), DefaultFloatEpsilon, true},

		// Types différents
		{"int vs float", 1, 1.0, 0, false},
		{"string vs int", "1", 1, 0, false},

		// Maps
		{
			"maps égaux",
			map[string]interface{}{"a": 1, "b": 2},
			map[string]interface{}{"a": 1, "b": 2},
			0,
			true,
		},
		{
			"maps différents",
			map[string]interface{}{"a": 1, "b": 2},
			map[string]interface{}{"a": 1, "b": 3},
			0,
			false,
		},

		// Slices
		{
			"slices égaux",
			[]interface{}{1, 2, 3},
			[]interface{}{1, 2, 3},
			0,
			true,
		},
		{
			"slices différents",
			[]interface{}{1, 2, 3},
			[]interface{}{1, 2, 4},
			0,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValuesEqual(tt.a, tt.b, tt.epsilon)
			if got != tt.want {
				t.Errorf("ValuesEqual(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestFactsEqual(t *testing.T) {
	tests := []struct {
		name  string
		fact1 map[string]interface{}
		fact2 map[string]interface{}
		want  bool
	}{
		{
			name:  "faits identiques",
			fact1: map[string]interface{}{"id": "123", "price": 100.0, "status": "active"},
			fact2: map[string]interface{}{"id": "123", "price": 100.0, "status": "active"},
			want:  true,
		},
		{
			name:  "faits différents (1 champ)",
			fact1: map[string]interface{}{"id": "123", "price": 100.0, "status": "active"},
			fact2: map[string]interface{}{"id": "123", "price": 150.0, "status": "active"},
			want:  false,
		},
		{
			name:  "faits de tailles différentes",
			fact1: map[string]interface{}{"id": "123", "price": 100.0},
			fact2: map[string]interface{}{"id": "123", "price": 100.0, "status": "active"},
			want:  false,
		},
		{
			name:  "faits avec champs manquants",
			fact1: map[string]interface{}{"id": "123", "price": 100.0},
			fact2: map[string]interface{}{"id": "123", "quantity": 10},
			want:  false,
		},
		{
			name:  "faits vides",
			fact1: map[string]interface{}{},
			fact2: map[string]interface{}{},
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FactsEqual(tt.fact1, tt.fact2)
			if got != tt.want {
				t.Errorf("FactsEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkValuesEqual_Simple(b *testing.B) {
	v1 := 42
	v2 := 42
	for i := 0; i < b.N; i++ {
		_ = ValuesEqual(v1, v2, 0)
	}
}

func BenchmarkValuesEqual_Float(b *testing.B) {
	v1 := 3.14159
	v2 := 3.14159
	for i := 0; i < b.N; i++ {
		_ = ValuesEqual(v1, v2, DefaultFloatEpsilon)
	}
}

func BenchmarkFactsEqual(b *testing.B) {
	fact1 := map[string]interface{}{
		"id":     "123",
		"price":  100.0,
		"status": "active",
		"qty":    10,
	}
	fact2 := map[string]interface{}{
		"id":     "123",
		"price":  100.0,
		"status": "active",
		"qty":    10,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FactsEqual(fact1, fact2)
	}
}
