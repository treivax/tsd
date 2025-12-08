// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"strings"
	"testing"
)

func TestValidationError_Error(t *testing.T) {
	tests := []struct {
		name     string
		ve       ValidationError
		expected string
	}{
		{
			name: "error with line number",
			ve: ValidationError{
				File:    "test.constraint",
				Type:    "rule",
				Message: "invalid syntax",
				Line:    42,
			},
			expected: "test.constraint:42: invalid syntax in rule",
		},
		{
			name: "error without line number",
			ve: ValidationError{
				File:    "test.constraint",
				Type:    "fact",
				Message: "missing field",
				Line:    0,
			},
			expected: "test.constraint: missing field in fact",
		},
		{
			name: "error with type definition",
			ve: ValidationError{
				File:    "types.constraint",
				Type:    ErrorTypeType,
				Message: "duplicate type definition",
				Line:    10,
			},
			expected: "types.constraint:10: duplicate type definition in type",
		},
		{
			name: "error with negative line (treated as 0)",
			ve: ValidationError{
				File:    "rules.constraint",
				Type:    ErrorTypeRule,
				Message: "undefined variable",
				Line:    -1,
			},
			expected: "rules.constraint: undefined variable in rule",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.ve.Error()
			if result != tt.expected {
				t.Errorf("ValidationError.Error() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestValidationErrors_Error(t *testing.T) {
	tests := []struct {
		name     string
		ve       ValidationErrors
		contains []string
	}{
		{
			name:     "empty errors",
			ve:       ValidationErrors{},
			contains: []string{"no validation errors"},
		},
		{
			name: "single error",
			ve: ValidationErrors{
				{
					File:    "test.constraint",
					Type:    "rule",
					Message: "invalid syntax",
					Line:    42,
				},
			},
			contains: []string{"test.constraint:42", "invalid syntax", "rule"},
		},
		{
			name: "multiple errors",
			ve: ValidationErrors{
				{
					File:    "test1.constraint",
					Type:    "rule",
					Message: "invalid syntax",
					Line:    10,
				},
				{
					File:    "test2.constraint",
					Type:    "fact",
					Message: "missing field",
					Line:    20,
				},
				{
					File:    "test3.constraint",
					Type:    "type",
					Message: "duplicate definition",
					Line:    0,
				},
			},
			contains: []string{
				"3 validation errors",
				"test1.constraint:10",
				"test2.constraint:20",
				"test3.constraint",
				"invalid syntax",
				"missing field",
				"duplicate definition",
				"1.",
				"2.",
				"3.",
			},
		},
		{
			name: "two errors",
			ve: ValidationErrors{
				{
					File:    "file1.constraint",
					Type:    ErrorTypeFact,
					Message: "error one",
					Line:    5,
				},
				{
					File:    "file2.constraint",
					Type:    ErrorTypeRule,
					Message: "error two",
					Line:    15,
				},
			},
			contains: []string{
				"2 validation errors",
				"file1.constraint:5",
				"file2.constraint:15",
				"error one",
				"error two",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.ve.Error()
			for _, substr := range tt.contains {
				if !strings.Contains(result, substr) {
					t.Errorf("ValidationErrors.Error() = %q, want to contain %q", result, substr)
				}
			}
		})
	}
}

func TestValidationErrors_HasErrors(t *testing.T) {
	tests := []struct {
		name     string
		ve       ValidationErrors
		expected bool
	}{
		{
			name:     "empty errors",
			ve:       ValidationErrors{},
			expected: false,
		},
		{
			name:     "nil errors",
			ve:       nil,
			expected: false,
		},
		{
			name: "single error",
			ve: ValidationErrors{
				{File: "test.constraint", Type: "rule", Message: "error", Line: 1},
			},
			expected: true,
		},
		{
			name: "multiple errors",
			ve: ValidationErrors{
				{File: "test1.constraint", Type: "rule", Message: "error1", Line: 1},
				{File: "test2.constraint", Type: "fact", Message: "error2", Line: 2},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.ve.HasErrors()
			if result != tt.expected {
				t.Errorf("ValidationErrors.HasErrors() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestValidationErrors_Count(t *testing.T) {
	tests := []struct {
		name     string
		ve       ValidationErrors
		expected int
	}{
		{
			name:     "empty errors",
			ve:       ValidationErrors{},
			expected: 0,
		},
		{
			name:     "nil errors",
			ve:       nil,
			expected: 0,
		},
		{
			name: "single error",
			ve: ValidationErrors{
				{File: "test.constraint", Type: "rule", Message: "error", Line: 1},
			},
			expected: 1,
		},
		{
			name: "multiple errors",
			ve: ValidationErrors{
				{File: "test1.constraint", Type: "rule", Message: "error1", Line: 1},
				{File: "test2.constraint", Type: "fact", Message: "error2", Line: 2},
				{File: "test3.constraint", Type: "type", Message: "error3", Line: 3},
			},
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.ve.Count()
			if result != tt.expected {
				t.Errorf("ValidationErrors.Count() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestErrorTypeConstants(t *testing.T) {
	// Test that error type constants have expected values
	if ErrorTypeFact != "fact" {
		t.Errorf("ErrorTypeFact = %q, want %q", ErrorTypeFact, "fact")
	}
	if ErrorTypeRule != "rule" {
		t.Errorf("ErrorTypeRule = %q, want %q", ErrorTypeRule, "rule")
	}
	if ErrorTypeType != "type" {
		t.Errorf("ErrorTypeType = %q, want %q", ErrorTypeType, "type")
	}
}

func TestValidationErrors_AsError(t *testing.T) {
	// Test that ValidationErrors can be used as an error interface
	ve := ValidationErrors{
		{File: "test.constraint", Type: "rule", Message: "test error", Line: 1},
	}
	var err error = ve

	// Verify it implements error interface by checking the error message
	if err.Error() != ve.Error() {
		t.Errorf("error.Error() = %q, want %q", err.Error(), ve.Error())
	}

	// Verify type assertion works
	if _, ok := err.(ValidationErrors); !ok {
		t.Error("ValidationErrors should be assertable from error interface")
	}
}

func TestValidationError_Fields(t *testing.T) {
	// Test that all fields are properly set and accessible
	ve := ValidationError{
		File:    "myfile.constraint",
		Type:    "myrule",
		Message: "mymessage",
		Line:    123,
	}

	if ve.File != "myfile.constraint" {
		t.Errorf("File = %q, want %q", ve.File, "myfile.constraint")
	}
	if ve.Type != "myrule" {
		t.Errorf("Type = %q, want %q", ve.Type, "myrule")
	}
	if ve.Message != "mymessage" {
		t.Errorf("Message = %q, want %q", ve.Message, "mymessage")
	}
	if ve.Line != 123 {
		t.Errorf("Line = %d, want %d", ve.Line, 123)
	}
}
