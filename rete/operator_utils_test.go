// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete
import (
	"testing"
)
func TestDecodeOperator(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// Base64 encoded operators
		{
			name:     "Plus operator (Base64)",
			input:    "Kw==",
			expected: "+",
		},
		{
			name:     "Minus operator (Base64)",
			input:    "LQ==",
			expected: "-",
		},
		{
			name:     "Multiply operator (Base64)",
			input:    "Kg==",
			expected: "*",
		},
		{
			name:     "Divide operator (Base64)",
			input:    "Lw==",
			expected: "/",
		},
		{
			name:     "Modulo operator (Base64)",
			input:    "JQ==",
			expected: "%",
		},
		{
			name:     "Equals operator (Base64)",
			input:    "PT0=",
			expected: "==",
		},
		{
			name:     "Not equals operator (Base64)",
			input:    "IT0=",
			expected: "!=",
		},
		{
			name:     "Less than operator (Base64)",
			input:    "PA==",
			expected: "<",
		},
		{
			name:     "Less than or equal operator (Base64)",
			input:    "PD0=",
			expected: "<=",
		},
		{
			name:     "Greater than operator (Base64)",
			input:    "Pg==",
			expected: ">",
		},
		{
			name:     "Greater than or equal operator (Base64)",
			input:    "Pj0=",
			expected: ">=",
		},
		// Already decoded operators (idempotent)
		{
			name:     "Already decoded plus",
			input:    "+",
			expected: "+",
		},
		{
			name:     "Already decoded minus",
			input:    "-",
			expected: "-",
		},
		{
			name:     "Already decoded multiply",
			input:    "*",
			expected: "*",
		},
		{
			name:     "Already decoded divide",
			input:    "/",
			expected: "/",
		},
		{
			name:     "Already decoded modulo",
			input:    "%",
			expected: "%",
		},
		{
			name:     "Already decoded equals",
			input:    "==",
			expected: "==",
		},
		{
			name:     "Already decoded not equals",
			input:    "!=",
			expected: "!=",
		},
		// String operators
		{
			name:     "CONTAINS operator",
			input:    "CONTAINS",
			expected: "CONTAINS",
		},
		{
			name:     "IN operator",
			input:    "IN",
			expected: "IN",
		},
		{
			name:     "LIKE operator",
			input:    "LIKE",
			expected: "LIKE",
		},
		{
			name:     "MATCHES operator",
			input:    "MATCHES",
			expected: "MATCHES",
		},
		// Logical operators
		{
			name:     "AND operator",
			input:    "AND",
			expected: "AND",
		},
		{
			name:     "OR operator",
			input:    "OR",
			expected: "OR",
		},
		{
			name:     "NOT operator",
			input:    "NOT",
			expected: "NOT",
		},
		// Edge cases
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Invalid Base64 that's not an operator",
			input:    "abc123",
			expected: "abc123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DecodeOperator(tt.input)
			if result != tt.expected {
				t.Errorf("DecodeOperator(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
func TestIsValidOperator(t *testing.T) {
	tests := []struct {
		name     string
		operator string
		expected bool
	}{
		// Valid arithmetic operators
		{"Plus", "+", true},
		{"Minus", "-", true},
		{"Multiply", "*", true},
		{"Divide", "/", true},
		{"Modulo", "%", true},
		// Valid comparison operators
		{"Equals", "==", true},
		{"Not equals", "!=", true},
		{"Less than", "<", true},
		{"Less than or equal", "<=", true},
		{"Greater than", ">", true},
		{"Greater than or equal", ">=", true},
		// Valid string operators
		{"CONTAINS", "CONTAINS", true},
		{"IN", "IN", true},
		{"LIKE", "LIKE", true},
		{"MATCHES", "MATCHES", true},
		// Valid logical operators
		{"AND", "AND", true},
		{"OR", "OR", true},
		{"NOT", "NOT", true},
		// Invalid operators
		{"Invalid symbol", "^", false},
		{"Invalid word", "INVALID", false},
		{"Empty string", "", false},
		{"Random text", "xyz", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidOperator(tt.operator)
			if result != tt.expected {
				t.Errorf("IsValidOperator(%q) = %v, want %v", tt.operator, result, tt.expected)
			}
		})
	}
}
func TestIsArithmeticOperator(t *testing.T) {
	tests := []struct {
		name     string
		operator string
		expected bool
	}{
		{"Plus", "+", true},
		{"Minus", "-", true},
		{"Multiply", "*", true},
		{"Divide", "/", true},
		{"Modulo", "%", true},
		{"Equals (not arithmetic)", "==", false},
		{"CONTAINS (not arithmetic)", "CONTAINS", false},
		{"AND (not arithmetic)", "AND", false},
		{"Invalid", "xyz", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsArithmeticOperator(tt.operator)
			if result != tt.expected {
				t.Errorf("IsArithmeticOperator(%q) = %v, want %v", tt.operator, result, tt.expected)
			}
		})
	}
}
func TestIsComparisonOperator(t *testing.T) {
	tests := []struct {
		name     string
		operator string
		expected bool
	}{
		{"Equals", "==", true},
		{"Not equals", "!=", true},
		{"Less than", "<", true},
		{"Less than or equal", "<=", true},
		{"Greater than", ">", true},
		{"Greater than or equal", ">=", true},
		{"Plus (not comparison)", "+", false},
		{"CONTAINS (not comparison)", "CONTAINS", false},
		{"AND (not comparison)", "AND", false},
		{"Invalid", "xyz", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsComparisonOperator(tt.operator)
			if result != tt.expected {
				t.Errorf("IsComparisonOperator(%q) = %v, want %v", tt.operator, result, tt.expected)
			}
		})
	}
}
func TestIsStringOperator(t *testing.T) {
	tests := []struct {
		name     string
		operator string
		expected bool
	}{
		{"CONTAINS", "CONTAINS", true},
		{"IN", "IN", true},
		{"LIKE", "LIKE", true},
		{"MATCHES", "MATCHES", true},
		{"Plus (not string)", "+", false},
		{"Equals (not string)", "==", false},
		{"AND (not string)", "AND", false},
		{"Invalid", "xyz", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsStringOperator(tt.operator)
			if result != tt.expected {
				t.Errorf("IsStringOperator(%q) = %v, want %v", tt.operator, result, tt.expected)
			}
		})
	}
}
func TestIsLogicalOperator(t *testing.T) {
	tests := []struct {
		name     string
		operator string
		expected bool
	}{
		{"AND", "AND", true},
		{"OR", "OR", true},
		{"NOT", "NOT", true},
		{"Plus (not logical)", "+", false},
		{"Equals (not logical)", "==", false},
		{"CONTAINS (not logical)", "CONTAINS", false},
		{"Invalid", "xyz", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsLogicalOperator(tt.operator)
			if result != tt.expected {
				t.Errorf("IsLogicalOperator(%q) = %v, want %v", tt.operator, result, tt.expected)
			}
		})
	}
}
func TestNormalizeOperator(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  string
		shouldErr bool
	}{
		// Valid operators
		{
			name:      "Base64 plus",
			input:     "Kw==",
			expected:  "+",
			shouldErr: false,
		},
		{
			name:      "Already decoded plus",
			input:     "+",
			expected:  "+",
			shouldErr: false,
		},
		{
			name:      "CONTAINS operator",
			input:     "CONTAINS",
			expected:  "CONTAINS",
			shouldErr: false,
		},
		// Invalid operators
		{
			name:      "Invalid operator",
			input:     "INVALID",
			expected:  "",
			shouldErr: true,
		},
		{
			name:      "Empty string",
			input:     "",
			expected:  "",
			shouldErr: true,
		},
		{
			name:      "Random text",
			input:     "xyz123",
			expected:  "",
			shouldErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NormalizeOperator(tt.input)
			if tt.shouldErr {
				if err == nil {
					t.Errorf("NormalizeOperator(%q) expected error, got nil", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("NormalizeOperator(%q) unexpected error: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("NormalizeOperator(%q) = %q, want %q", tt.input, result, tt.expected)
				}
			}
		})
	}
}
func TestExtractOperatorFromMap(t *testing.T) {
	tests := []struct {
		name      string
		input     map[string]interface{}
		expected  string
		shouldErr bool
	}{
		{
			name: "String operator",
			input: map[string]interface{}{
				"operator": "+",
			},
			expected:  "+",
			shouldErr: false,
		},
		{
			name: "Base64 string operator",
			input: map[string]interface{}{
				"operator": "Kw==",
			},
			expected:  "+",
			shouldErr: false,
		},
		{
			name: "Byte array operator",
			input: map[string]interface{}{
				"operator": []uint8{'+'},
			},
			expected:  "+",
			shouldErr: false,
		},
		{
			name: "Byte slice operator",
			input: map[string]interface{}{
				"operator": []byte("*"),
			},
			expected:  "*",
			shouldErr: false,
		},
		{
			name: "Missing operator key",
			input: map[string]interface{}{
				"other": "+",
			},
			expected:  "",
			shouldErr: true,
		},
		{
			name: "Invalid operator type (int)",
			input: map[string]interface{}{
				"operator": 42,
			},
			expected:  "",
			shouldErr: true,
		},
		{
			name: "Invalid operator type (bool)",
			input: map[string]interface{}{
				"operator": true,
			},
			expected:  "",
			shouldErr: true,
		},
		{
			name:      "Nil map",
			input:     nil,
			expected:  "",
			shouldErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ExtractOperatorFromMap(tt.input)
			if tt.shouldErr {
				if err == nil {
					t.Errorf("ExtractOperatorFromMap() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("ExtractOperatorFromMap() unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("ExtractOperatorFromMap() = %q, want %q", result, tt.expected)
				}
			}
		})
	}
}
// Benchmark tests
func BenchmarkDecodeOperator(b *testing.B) {
	operators := []string{"Kw==", "+", "Kg==", "*", "PT0=", "=="}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		op := operators[i%len(operators)]
		DecodeOperator(op)
	}
}
func BenchmarkIsValidOperator(b *testing.B) {
	operators := []string{"+", "-", "*", "/", "%", "==", "!=", "<", ">", "CONTAINS"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		op := operators[i%len(operators)]
		IsValidOperator(op)
	}
}
func BenchmarkNormalizeOperator(b *testing.B) {
	operators := []string{"Kw==", "+", "Kg==", "*", "PT0=", "=="}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		op := operators[i%len(operators)]
		NormalizeOperator(op)
	}
}