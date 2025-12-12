// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete
import (
	"testing"
)
func TestEvaluateArithmeticOperation(t *testing.T) {
	evaluator := &AlphaConditionEvaluator{}
	tests := []struct {
		name     string
		left     interface{}
		operator string
		right    interface{}
		expected interface{}
		wantErr  bool
	}{
		{
			name:     "addition",
			left:     10.0,
			operator: "+",
			right:    5.0,
			expected: 15.0,
			wantErr:  false,
		},
		{
			name:     "subtraction",
			left:     10.0,
			operator: "-",
			right:    5.0,
			expected: 5.0,
			wantErr:  false,
		},
		{
			name:     "multiplication",
			left:     10.0,
			operator: "*",
			right:    5.0,
			expected: 50.0,
			wantErr:  false,
		},
		{
			name:     "division",
			left:     10.0,
			operator: "/",
			right:    5.0,
			expected: 2.0,
			wantErr:  false,
		},
		{
			name:     "modulo",
			left:     10.0,
			operator: "%",
			right:    3.0,
			expected: 1.0,
			wantErr:  false,
		},
		{
			name:     "division by zero",
			left:     10.0,
			operator: "/",
			right:    0.0,
			wantErr:  true,
		},
		{
			name:     "modulo by zero",
			left:     10.0,
			operator: "%",
			right:    0.0,
			wantErr:  true,
		},
		{
			name:     "unsupported operator",
			left:     10.0,
			operator: "^",
			right:    5.0,
			wantErr:  true,
		},
		{
			name:     "non-numeric left",
			left:     "hello",
			operator: "+",
			right:    5.0,
			wantErr:  true,
		},
		{
			name:     "non-numeric right",
			left:     10.0,
			operator: "+",
			right:    "world",
			wantErr:  true,
		},
		{
			name:     "negative numbers",
			left:     -10.0,
			operator: "+",
			right:    -5.0,
			expected: -15.0,
			wantErr:  false,
		},
		{
			name:     "float division",
			left:     7.0,
			operator: "/",
			right:    2.0,
			expected: 3.5,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evaluator.evaluateArithmeticOperation(tt.left, tt.operator, tt.right)
			if (err != nil) != tt.wantErr {
				t.Errorf("evaluateArithmeticOperation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("evaluateArithmeticOperation() = %v, want %v", result, tt.expected)
			}
		})
	}
}
func TestEvaluateContains(t *testing.T) {
	evaluator := &AlphaConditionEvaluator{}
	tests := []struct {
		name     string
		left     interface{}
		right    interface{}
		expected bool
		wantErr  bool
	}{
		{
			name:     "contains substring",
			left:     "hello world",
			right:    "world",
			expected: true,
			wantErr:  false,
		},
		{
			name:     "does not contain",
			left:     "hello world",
			right:    "xyz",
			expected: false,
			wantErr:  false,
		},
		{
			name:     "contains at beginning",
			left:     "hello world",
			right:    "hello",
			expected: true,
			wantErr:  false,
		},
		{
			name:     "contains at end",
			left:     "hello world",
			right:    "world",
			expected: true,
			wantErr:  false,
		},
		{
			name:     "contains empty string",
			left:     "hello",
			right:    "",
			expected: true,
			wantErr:  false,
		},
		{
			name:     "empty string contains empty",
			left:     "",
			right:    "",
			expected: true,
			wantErr:  false,
		},
		{
			name:     "case sensitive",
			left:     "Hello World",
			right:    "hello",
			expected: false,
			wantErr:  false,
		},
		{
			name:    "non-string left",
			left:    123,
			right:   "hello",
			wantErr: true,
		},
		{
			name:    "non-string right",
			left:    "hello",
			right:   123,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evaluator.evaluateContains(tt.left, tt.right)
			if (err != nil) != tt.wantErr {
				t.Errorf("evaluateContains() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("evaluateContains() = %v, want %v", result, tt.expected)
			}
		})
	}
}
func TestEvaluateIn(t *testing.T) {
	evaluator := &AlphaConditionEvaluator{}
	tests := []struct {
		name     string
		left     interface{}
		right    interface{}
		expected bool
		wantErr  bool
	}{
		{
			name:     "string in slice",
			left:     "apple",
			right:    []interface{}{"apple", "banana", "orange"},
			expected: true,
			wantErr:  false,
		},
		{
			name:     "string not in slice",
			left:     "grape",
			right:    []interface{}{"apple", "banana", "orange"},
			expected: false,
			wantErr:  false,
		},
		{
			name:     "number in slice",
			left:     5.0,
			right:    []interface{}{1.0, 2.0, 5.0, 10.0},
			expected: true,
			wantErr:  false,
		},
		{
			name:     "number not in slice",
			left:     7.0,
			right:    []interface{}{1.0, 2.0, 5.0, 10.0},
			expected: false,
			wantErr:  false,
		},
		{
			name:     "in string slice",
			left:     "hello",
			right:    []string{"hello", "world"},
			expected: true,
			wantErr:  false,
		},
		{
			name:     "in int slice",
			left:     5,
			right:    []int{1, 2, 5, 10},
			expected: true,
			wantErr:  false,
		},
		{
			name:     "in float64 slice",
			left:     5.5,
			right:    []float64{1.1, 2.2, 5.5, 10.0},
			expected: true,
			wantErr:  false,
		},
		{
			name:     "empty slice",
			left:     "test",
			right:    []interface{}{},
			expected: false,
			wantErr:  false,
		},
		{
			name:    "non-slice right",
			left:    "test",
			right:   "not a slice",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evaluator.evaluateIn(tt.left, tt.right)
			if (err != nil) != tt.wantErr {
				t.Errorf("evaluateIn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("evaluateIn() = %v, want %v", result, tt.expected)
			}
		})
	}
}
func TestEvaluateLike(t *testing.T) {
	evaluator := &AlphaConditionEvaluator{}
	tests := []struct {
		name     string
		left     interface{}
		right    interface{}
		expected bool
		wantErr  bool
	}{
		{
			name:     "exact match",
			left:     "hello",
			right:    "hello",
			expected: true,
			wantErr:  false,
		},
		{
			name:     "percent wildcard at end",
			left:     "hello world",
			right:    "hello%",
			expected: true,
			wantErr:  false,
		},
		{
			name:     "percent wildcard at start",
			left:     "hello world",
			right:    "%world",
			expected: true,
			wantErr:  false,
		},
		{
			name:     "percent wildcard both sides",
			left:     "hello world test",
			right:    "%world%",
			expected: true,
			wantErr:  false,
		},
		{
			name:     "underscore single char",
			left:     "hello",
			right:    "hel_o",
			expected: true,
			wantErr:  false,
		},
		{
			name:     "underscore multiple",
			left:     "hello",
			right:    "h___o",
			expected: true,
			wantErr:  false,
		},
		{
			name:     "mixed wildcards",
			left:     "hello world",
			right:    "h%_ld",
			expected: true,
			wantErr:  false,
		},
		{
			name:     "no match",
			left:     "hello",
			right:    "world",
			expected: false,
			wantErr:  false,
		},
		{
			name:     "no match with wildcard",
			left:     "hello",
			right:    "world%",
			expected: false,
			wantErr:  false,
		},
		{
			name:     "case sensitive",
			left:     "Hello",
			right:    "hello",
			expected: false,
			wantErr:  false,
		},
		{
			name:    "non-string left",
			left:    123,
			right:   "hello",
			wantErr: true,
		},
		{
			name:    "non-string right",
			left:    "hello",
			right:   123,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evaluator.evaluateLike(tt.left, tt.right)
			if (err != nil) != tt.wantErr {
				t.Errorf("evaluateLike() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("evaluateLike() = %v, want %v", result, tt.expected)
			}
		})
	}
}
func TestEvaluateMatches(t *testing.T) {
	evaluator := &AlphaConditionEvaluator{}
	tests := []struct {
		name     string
		left     interface{}
		right    interface{}
		expected bool
		wantErr  bool
	}{
		{
			name:     "simple pattern match",
			left:     "hello123",
			right:    "^hello[0-9]+$",
			expected: true,
			wantErr:  false,
		},
		{
			name:     "email pattern",
			left:     "test@example.com",
			right:    "^[a-z]+@[a-z]+\\.[a-z]+$",
			expected: true,
			wantErr:  false,
		},
		{
			name:     "no match",
			left:     "hello",
			right:    "^world$",
			expected: false,
			wantErr:  false,
		},
		{
			name:     "digit pattern",
			left:     "12345",
			right:    "^[0-9]+$",
			expected: true,
			wantErr:  false,
		},
		{
			name:     "partial match with anchors",
			left:     "hello world",
			right:    "world",
			expected: true,
			wantErr:  false,
		},
		{
			name:     "case insensitive flag",
			left:     "HELLO",
			right:    "(?i)^hello$",
			expected: true,
			wantErr:  false,
		},
		{
			name:    "invalid regex pattern",
			left:    "hello",
			right:   "[invalid",
			wantErr: true,
		},
		{
			name:    "non-string left",
			left:    123,
			right:   "^[0-9]+$",
			wantErr: true,
		},
		{
			name:    "non-string right",
			left:    "hello",
			right:   123,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evaluator.evaluateMatches(tt.left, tt.right)
			if (err != nil) != tt.wantErr {
				t.Errorf("evaluateMatches() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("evaluateMatches() = %v, want %v", result, tt.expected)
			}
		})
	}
}