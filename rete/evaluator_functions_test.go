// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"testing"
)

func TestEvaluateLength(t *testing.T) {
	evaluator := &AlphaConditionEvaluator{}

	tests := []struct {
		name     string
		args     []interface{}
		expected interface{}
		wantErr  bool
	}{
		{
			name:     "empty string",
			args:     []interface{}{""},
			expected: float64(0),
			wantErr:  false,
		},
		{
			name:     "simple string",
			args:     []interface{}{"hello"},
			expected: float64(5),
			wantErr:  false,
		},
		{
			name:     "string with spaces",
			args:     []interface{}{"hello world"},
			expected: float64(11),
			wantErr:  false,
		},
		{
			name:     "unicode string",
			args:     []interface{}{"hello 世界"},
			expected: float64(12), // UTF-8 byte count (len() in Go counts bytes)
			wantErr:  false,
		},
		{
			name:    "too many arguments",
			args:    []interface{}{"hello", "world"},
			wantErr: true,
		},
		{
			name:    "no arguments",
			args:    []interface{}{},
			wantErr: true,
		},
		{
			name:    "non-string argument",
			args:    []interface{}{123},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evaluator.evaluateLength(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("evaluateLength() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("evaluateLength() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestEvaluateUpper(t *testing.T) {
	evaluator := &AlphaConditionEvaluator{}

	tests := []struct {
		name     string
		args     []interface{}
		expected interface{}
		wantErr  bool
	}{
		{
			name:     "lowercase string",
			args:     []interface{}{"hello"},
			expected: "HELLO",
			wantErr:  false,
		},
		{
			name:     "mixed case string",
			args:     []interface{}{"Hello World"},
			expected: "HELLO WORLD",
			wantErr:  false,
		},
		{
			name:     "already uppercase",
			args:     []interface{}{"HELLO"},
			expected: "HELLO",
			wantErr:  false,
		},
		{
			name:     "empty string",
			args:     []interface{}{""},
			expected: "",
			wantErr:  false,
		},
		{
			name:    "too many arguments",
			args:    []interface{}{"hello", "world"},
			wantErr: true,
		},
		{
			name:    "non-string argument",
			args:    []interface{}{123},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evaluator.evaluateUpper(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("evaluateUpper() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("evaluateUpper() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestEvaluateLower(t *testing.T) {
	evaluator := &AlphaConditionEvaluator{}

	tests := []struct {
		name     string
		args     []interface{}
		expected interface{}
		wantErr  bool
	}{
		{
			name:     "uppercase string",
			args:     []interface{}{"HELLO"},
			expected: "hello",
			wantErr:  false,
		},
		{
			name:     "mixed case string",
			args:     []interface{}{"Hello World"},
			expected: "hello world",
			wantErr:  false,
		},
		{
			name:     "already lowercase",
			args:     []interface{}{"hello"},
			expected: "hello",
			wantErr:  false,
		},
		{
			name:     "empty string",
			args:     []interface{}{""},
			expected: "",
			wantErr:  false,
		},
		{
			name:    "too many arguments",
			args:    []interface{}{"HELLO", "WORLD"},
			wantErr: true,
		},
		{
			name:    "non-string argument",
			args:    []interface{}{123},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evaluator.evaluateLower(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("evaluateLower() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("evaluateLower() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestEvaluateAbs(t *testing.T) {
	evaluator := &AlphaConditionEvaluator{}

	tests := []struct {
		name     string
		args     []interface{}
		expected interface{}
		wantErr  bool
	}{
		{
			name:     "positive number",
			args:     []interface{}{5.5},
			expected: 5.5,
			wantErr:  false,
		},
		{
			name:     "negative number",
			args:     []interface{}{-5.5},
			expected: 5.5,
			wantErr:  false,
		},
		{
			name:     "zero",
			args:     []interface{}{0.0},
			expected: 0.0,
			wantErr:  false,
		},
		{
			name:     "large negative",
			args:     []interface{}{-1000.99},
			expected: 1000.99,
			wantErr:  false,
		},
		{
			name:    "too many arguments",
			args:    []interface{}{5.5, 10.0},
			wantErr: true,
		},
		{
			name:    "non-number argument",
			args:    []interface{}{"hello"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evaluator.evaluateAbs(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("evaluateAbs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("evaluateAbs() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestEvaluateRound(t *testing.T) {
	evaluator := &AlphaConditionEvaluator{}

	tests := []struct {
		name     string
		args     []interface{}
		expected interface{}
		wantErr  bool
	}{
		{
			name:     "round up",
			args:     []interface{}{5.6},
			expected: 6.0,
			wantErr:  false,
		},
		{
			name:     "round down",
			args:     []interface{}{5.4},
			expected: 5.0,
			wantErr:  false,
		},
		{
			name:     "round exact half",
			args:     []interface{}{5.5},
			expected: 6.0,
			wantErr:  false,
		},
		{
			name:     "negative round",
			args:     []interface{}{-5.5},
			expected: -6.0,
			wantErr:  false,
		},
		{
			name:     "already integer",
			args:     []interface{}{5.0},
			expected: 5.0,
			wantErr:  false,
		},
		{
			name:    "too many arguments",
			args:    []interface{}{5.5, 10.0},
			wantErr: true,
		},
		{
			name:    "non-number argument",
			args:    []interface{}{"hello"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evaluator.evaluateRound(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("evaluateRound() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("evaluateRound() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestEvaluateFloor(t *testing.T) {
	evaluator := &AlphaConditionEvaluator{}

	tests := []struct {
		name     string
		args     []interface{}
		expected interface{}
		wantErr  bool
	}{
		{
			name:     "floor positive",
			args:     []interface{}{5.9},
			expected: 5.0,
			wantErr:  false,
		},
		{
			name:     "floor negative",
			args:     []interface{}{-5.1},
			expected: -6.0,
			wantErr:  false,
		},
		{
			name:     "floor exact",
			args:     []interface{}{5.0},
			expected: 5.0,
			wantErr:  false,
		},
		{
			name:    "too many arguments",
			args:    []interface{}{5.5, 10.0},
			wantErr: true,
		},
		{
			name:    "non-number argument",
			args:    []interface{}{"hello"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evaluator.evaluateFloor(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("evaluateFloor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("evaluateFloor() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestEvaluateCeil(t *testing.T) {
	evaluator := &AlphaConditionEvaluator{}

	tests := []struct {
		name     string
		args     []interface{}
		expected interface{}
		wantErr  bool
	}{
		{
			name:     "ceil positive",
			args:     []interface{}{5.1},
			expected: 6.0,
			wantErr:  false,
		},
		{
			name:     "ceil negative",
			args:     []interface{}{-5.9},
			expected: -5.0,
			wantErr:  false,
		},
		{
			name:     "ceil exact",
			args:     []interface{}{5.0},
			expected: 5.0,
			wantErr:  false,
		},
		{
			name:    "too many arguments",
			args:    []interface{}{5.5, 10.0},
			wantErr: true,
		},
		{
			name:    "non-number argument",
			args:    []interface{}{"hello"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evaluator.evaluateCeil(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("evaluateCeil() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("evaluateCeil() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestEvaluateSubstring(t *testing.T) {
	evaluator := &AlphaConditionEvaluator{}

	tests := []struct {
		name     string
		args     []interface{}
		expected interface{}
		wantErr  bool
	}{
		{
			name:     "substring with start only",
			args:     []interface{}{"hello world", 6.0},
			expected: "world",
			wantErr:  false,
		},
		{
			name:     "substring with start and length",
			args:     []interface{}{"hello world", 0.0, 5.0},
			expected: "hello",
			wantErr:  false,
		},
		{
			name:     "substring from middle",
			args:     []interface{}{"hello world", 3.0, 5.0},
			expected: "lo wo",
			wantErr:  false,
		},
		{
			name:     "substring start at beginning",
			args:     []interface{}{"hello", 0.0},
			expected: "hello",
			wantErr:  false,
		},
		{
			name:     "substring length exceeds string",
			args:     []interface{}{"hello", 2.0, 100.0},
			expected: "llo",
			wantErr:  false,
		},
		{
			name:     "substring start out of bounds",
			args:     []interface{}{"hello", 10.0},
			expected: "",
			wantErr:  false,
		},
		{
			name:     "substring negative start",
			args:     []interface{}{"hello", -1.0},
			expected: "",
			wantErr:  false,
		},
		{
			name:    "substring too few arguments",
			args:    []interface{}{"hello"},
			wantErr: true,
		},
		{
			name:    "substring too many arguments",
			args:    []interface{}{"hello", 0.0, 5.0, 10.0},
			wantErr: true,
		},
		{
			name:    "substring non-string first arg",
			args:    []interface{}{123, 0.0},
			wantErr: true,
		},
		{
			name:    "substring non-number second arg",
			args:    []interface{}{"hello", "world"},
			wantErr: true,
		},
		{
			name:    "substring non-number third arg",
			args:    []interface{}{"hello", 0.0, "world"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evaluator.evaluateSubstring(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("evaluateSubstring() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("evaluateSubstring() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestEvaluateTrim(t *testing.T) {
	evaluator := &AlphaConditionEvaluator{}

	tests := []struct {
		name     string
		args     []interface{}
		expected interface{}
		wantErr  bool
	}{
		{
			name:     "trim leading spaces",
			args:     []interface{}{"  hello"},
			expected: "hello",
			wantErr:  false,
		},
		{
			name:     "trim trailing spaces",
			args:     []interface{}{"hello  "},
			expected: "hello",
			wantErr:  false,
		},
		{
			name:     "trim both sides",
			args:     []interface{}{"  hello  "},
			expected: "hello",
			wantErr:  false,
		},
		{
			name:     "trim tabs and newlines",
			args:     []interface{}{"\t\nhello\n\t"},
			expected: "hello",
			wantErr:  false,
		},
		{
			name:     "no spaces to trim",
			args:     []interface{}{"hello"},
			expected: "hello",
			wantErr:  false,
		},
		{
			name:     "empty string",
			args:     []interface{}{""},
			expected: "",
			wantErr:  false,
		},
		{
			name:     "only spaces",
			args:     []interface{}{"   "},
			expected: "",
			wantErr:  false,
		},
		{
			name:    "too many arguments",
			args:    []interface{}{"hello", "world"},
			wantErr: true,
		},
		{
			name:    "non-string argument",
			args:    []interface{}{123},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evaluator.evaluateTrim(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("evaluateTrim() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("evaluateTrim() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestEvaluateFunctionCall(t *testing.T) {
	fact := &Fact{
		ID:   "f1",
		Type: "Test",
		Fields: map[string]interface{}{
			"name":   "Alice",
			"age":    30,
			"salary": 50000.75,
		},
	}

	evaluator := NewAlphaConditionEvaluator()
	evaluator.variableBindings = map[string]*Fact{"f": fact}

	tests := []struct {
		name     string
		input    map[string]interface{}
		expected interface{}
		wantErr  bool
	}{
		{
			name: "LENGTH function",
			input: map[string]interface{}{
				"name": "LENGTH",
				"args": []interface{}{"hello"},
			},
			expected: float64(5),
			wantErr:  false,
		},
		{
			name: "UPPER function",
			input: map[string]interface{}{
				"name": "UPPER",
				"args": []interface{}{"hello"},
			},
			expected: "HELLO",
			wantErr:  false,
		},
		{
			name: "LOWER function",
			input: map[string]interface{}{
				"name": "LOWER",
				"args": []interface{}{"HELLO"},
			},
			expected: "hello",
			wantErr:  false,
		},
		{
			name: "ABS function",
			input: map[string]interface{}{
				"name": "ABS",
				"args": []interface{}{-5.5},
			},
			expected: 5.5,
			wantErr:  false,
		},
		{
			name: "ROUND function",
			input: map[string]interface{}{
				"name": "ROUND",
				"args": []interface{}{5.5},
			},
			expected: 6.0,
			wantErr:  false,
		},
		{
			name: "FLOOR function",
			input: map[string]interface{}{
				"name": "FLOOR",
				"args": []interface{}{5.9},
			},
			expected: 5.0,
			wantErr:  false,
		},
		{
			name: "CEIL function",
			input: map[string]interface{}{
				"name": "CEIL",
				"args": []interface{}{5.1},
			},
			expected: 6.0,
			wantErr:  false,
		},
		{
			name: "SUBSTRING function",
			input: map[string]interface{}{
				"name": "SUBSTRING",
				"args": []interface{}{"hello world", 6.0},
			},
			expected: "world",
			wantErr:  false,
		},
		{
			name: "TRIM function",
			input: map[string]interface{}{
				"name": "TRIM",
				"args": []interface{}{"  hello  "},
			},
			expected: "hello",
			wantErr:  false,
		},
		{
			name: "unsupported function",
			input: map[string]interface{}{
				"name": "UNKNOWN",
				"args": []interface{}{},
			},
			wantErr: true,
		},
		{
			name: "invalid function name type",
			input: map[string]interface{}{
				"name": 123,
				"args": []interface{}{},
			},
			wantErr: true,
		},
		{
			name: "function with no args field",
			input: map[string]interface{}{
				"name": "LENGTH",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := evaluator.evaluateFunctionCall(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("evaluateFunctionCall() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("evaluateFunctionCall() = %v, want %v", result, tt.expected)
			}
		})
	}
}
