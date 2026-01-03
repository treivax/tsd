// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import "testing"

// TestCompareSignedIntegers teste la fonction compareSignedIntegers avec tous les types signés
func TestCompareSignedIntegers(t *testing.T) {
	tests := []struct {
		name   string
		a      interface{}
		b      interface{}
		wantEq bool
		wantOk bool
	}{
		// int
		{"int égaux", int(42), int(42), true, true},
		{"int différents", int(42), int(43), false, true},
		{"int vs mauvais type", int(42), "42", false, true},

		// int64
		{"int64 égaux", int64(1000), int64(1000), true, true},
		{"int64 différents", int64(1000), int64(2000), false, true},
		{"int64 vs int", int64(42), int(42), false, true},

		// int32
		{"int32 égaux", int32(100), int32(100), true, true},
		{"int32 différents", int32(100), int32(200), false, true},
		{"int32 vs int64", int32(100), int64(100), false, true},

		// int16
		{"int16 égaux", int16(10), int16(10), true, true},
		{"int16 différents", int16(10), int16(20), false, true},
		{"int16 vs int32", int16(10), int32(10), false, true},

		// int8
		{"int8 égaux", int8(5), int8(5), true, true},
		{"int8 différents", int8(5), int8(6), false, true},
		{"int8 vs int16", int8(5), int16(5), false, true},

		// Type non-signé (non géré)
		{"uint non géré", uint(42), uint(42), false, false},
		{"string non géré", "hello", "hello", false, false},
		{"float64 non géré", 1.5, 1.5, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEq, gotOk := compareSignedIntegers(tt.a, tt.b)
			if gotEq != tt.wantEq {
				t.Errorf("compareSignedIntegers() gotEq = %v, want %v", gotEq, tt.wantEq)
			}
			if gotOk != tt.wantOk {
				t.Errorf("compareSignedIntegers() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

// TestCompareUnsignedIntegers teste la fonction compareUnsignedIntegers avec tous les types non-signés
func TestCompareUnsignedIntegers(t *testing.T) {
	tests := []struct {
		name   string
		a      interface{}
		b      interface{}
		wantEq bool
		wantOk bool
	}{
		// uint
		{"uint égaux", uint(42), uint(42), true, true},
		{"uint différents", uint(42), uint(43), false, true},
		{"uint vs mauvais type", uint(42), "42", false, true},

		// uint64
		{"uint64 égaux", uint64(1000), uint64(1000), true, true},
		{"uint64 différents", uint64(1000), uint64(2000), false, true},
		{"uint64 vs uint", uint64(42), uint(42), false, true},

		// uint32
		{"uint32 égaux", uint32(100), uint32(100), true, true},
		{"uint32 différents", uint32(100), uint32(200), false, true},
		{"uint32 vs uint64", uint32(100), uint64(100), false, true},

		// uint16
		{"uint16 égaux", uint16(10), uint16(10), true, true},
		{"uint16 différents", uint16(10), uint16(20), false, true},
		{"uint16 vs uint32", uint16(10), uint32(10), false, true},

		// uint8
		{"uint8 égaux", uint8(5), uint8(5), true, true},
		{"uint8 différents", uint8(5), uint8(6), false, true},
		{"uint8 vs uint16", uint8(5), uint16(5), false, true},

		// Type signé (non géré)
		{"int non géré", int(42), int(42), false, false},
		{"string non géré", "hello", "hello", false, false},
		{"float64 non géré", 1.5, 1.5, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEq, gotOk := compareUnsignedIntegers(tt.a, tt.b)
			if gotEq != tt.wantEq {
				t.Errorf("compareUnsignedIntegers() gotEq = %v, want %v", gotEq, tt.wantEq)
			}
			if gotOk != tt.wantOk {
				t.Errorf("compareUnsignedIntegers() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

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
