// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"strings"
	"testing"
)

func TestParseXupleSpace_Valid(t *testing.T) {
	t.Log("🧪 TEST PARSING XUPLE-SPACE VALID")
	t.Log("===================================")

	tests := []struct {
		name                      string
		input                     string
		expectedName              string
		expectedSelection         string
		expectedConsumption       string
		expectedConsumptionLimit  int
		expectedRetention         string
		expectedRetentionDuration int
	}{
		{
			name: "complete configuration with all policies",
			input: `xuple-space agents-commands {
				selection: fifo
				consumption: once
				retention: unlimited
			}`,
			expectedName:        "agents-commands",
			expectedSelection:   "fifo",
			expectedConsumption: "once",
			expectedRetention:   "unlimited",
		},
		{
			name: "random selection",
			input: `xuple-space notifications {
				selection: random
				consumption: per-agent
				retention: duration(5m)
			}`,
			expectedName:              "notifications",
			expectedSelection:         "random",
			expectedConsumption:       "per-agent",
			expectedRetention:         "duration",
			expectedRetentionDuration: 300, // 5 minutes in seconds
		},
		{
			name: "lifo selection",
			input: `xuple-space stack {
				selection: lifo
				consumption: limited(3)
				retention: duration(1h)
			}`,
			expectedName:              "stack",
			expectedSelection:         "lifo",
			expectedConsumption:       "limited",
			expectedConsumptionLimit:  3,
			expectedRetention:         "duration",
			expectedRetentionDuration: 3600, // 1 hour in seconds
		},
		{
			name:                "minimal configuration with defaults",
			input:               `xuple-space minimal {}`,
			expectedName:        "minimal",
			expectedSelection:   "fifo",
			expectedConsumption: "once",
			expectedRetention:   "unlimited",
		},
		{
			name: "duration in seconds",
			input: `xuple-space cache {
				selection: fifo
				consumption: limited(10)
				retention: duration(30s)
			}`,
			expectedName:              "cache",
			expectedSelection:         "fifo",
			expectedConsumption:       "limited",
			expectedConsumptionLimit:  10,
			expectedRetention:         "duration",
			expectedRetentionDuration: 30,
		},
		{
			name: "duration in days",
			input: `xuple-space archive {
				selection: fifo
				consumption: once
				retention: duration(7d)
			}`,
			expectedName:              "archive",
			expectedSelection:         "fifo",
			expectedConsumption:       "once",
			expectedRetention:         "duration",
			expectedRetentionDuration: 604800, // 7 days in seconds
		},
		{
			name: "mixed properties order",
			input: `xuple-space mixed {
				retention: duration(2h)
				selection: random
				consumption: limited(5)
			}`,
			expectedName:              "mixed",
			expectedSelection:         "random",
			expectedConsumption:       "limited",
			expectedConsumptionLimit:  5,
			expectedRetention:         "duration",
			expectedRetentionDuration: 7200, // 2 hours
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("📝 Testing: %s", tt.name)

			result, err := ParseConstraint("test.tsd", []byte(tt.input))
			if err != nil {
				t.Fatalf("❌ Parse error: %v", err)
			}

			program, err := ConvertResultToProgram(result)
			if err != nil {
				t.Fatalf("❌ ConvertResultToProgram error: %v", err)
			}

			if len(program.XupleSpaces) != 1 {
				t.Fatalf("❌ Expected 1 xuple-space, got %d", len(program.XupleSpaces))
			}

			xs := program.XupleSpaces[0]

			// Validate name
			if xs.Name != tt.expectedName {
				t.Errorf("❌ Name: expected '%s', got '%s'", tt.expectedName, xs.Name)
			}

			// Validate selection policy
			if xs.SelectionPolicy != tt.expectedSelection {
				t.Errorf("❌ SelectionPolicy: expected '%s', got '%s'", tt.expectedSelection, xs.SelectionPolicy)
			}

			// Validate consumption policy
			if xs.ConsumptionPolicy.Type != tt.expectedConsumption {
				t.Errorf("❌ ConsumptionPolicy.Type: expected '%s', got '%s'", tt.expectedConsumption, xs.ConsumptionPolicy.Type)
			}

			if tt.expectedConsumptionLimit > 0 && xs.ConsumptionPolicy.Limit != tt.expectedConsumptionLimit {
				t.Errorf("❌ ConsumptionPolicy.Limit: expected %d, got %d", tt.expectedConsumptionLimit, xs.ConsumptionPolicy.Limit)
			}

			// Validate retention policy
			if xs.RetentionPolicy.Type != tt.expectedRetention {
				t.Errorf("❌ RetentionPolicy.Type: expected '%s', got '%s'", tt.expectedRetention, xs.RetentionPolicy.Type)
			}

			if tt.expectedRetentionDuration > 0 && xs.RetentionPolicy.Duration != tt.expectedRetentionDuration {
				t.Errorf("❌ RetentionPolicy.Duration: expected %d, got %d", tt.expectedRetentionDuration, xs.RetentionPolicy.Duration)
			}

			t.Log("✅ Test passed")
		})
	}
}

func TestParseXupleSpace_Invalid(t *testing.T) {
	t.Log("🧪 TEST PARSING XUPLE-SPACE INVALID")
	t.Log("====================================")

	tests := []struct {
		name        string
		input       string
		expectedErr string
	}{
		{
			name: "invalid selection policy",
			input: `xuple-space bad {
				selection: priority
			}`,
			expectedErr: "no match found",
		},
		{
			name: "consumption limit zero",
			input: `xuple-space bad {
				consumption: limited(0)
			}`,
			expectedErr: "consumption limit must be greater than zero",
		},
		{
			name: "negative consumption limit",
			input: `xuple-space bad {
				consumption: limited(-5)
			}`,
			expectedErr: "no match found",
		},
		{
			name: "duration with zero value",
			input: `xuple-space bad {
				retention: duration(0s)
			}`,
			expectedErr: "duration must be positive",
		},
		{
			name: "duration with negative value",
			input: `xuple-space bad {
				retention: duration(-5m)
			}`,
			expectedErr: "no match found",
		},
		{
			name: "invalid time unit",
			input: `xuple-space bad {
				retention: duration(5x)
			}`,
			expectedErr: "no match found",
		},
		{
			name: "missing parentheses for limited",
			input: `xuple-space bad {
				consumption: limited 5
			}`,
			expectedErr: "no match found",
		},
		{
			name: "missing parentheses for duration",
			input: `xuple-space bad {
				retention: duration 5m
			}`,
			expectedErr: "no match found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("📝 Testing: %s", tt.name)

			_, err := ParseConstraint("test.tsd", []byte(tt.input))
			if err == nil {
				t.Fatal("❌ Expected error but parsing succeeded")
			}

			if !strings.Contains(err.Error(), tt.expectedErr) {
				t.Errorf("❌ Expected error containing '%s', got: %v", tt.expectedErr, err)
			}

			t.Logf("✅ Error correctly detected: %v", err)
		})
	}
}

func TestParseXupleSpace_MultipleDeclarations(t *testing.T) {
	t.Log("🧪 TEST MULTIPLE XUPLE-SPACE DECLARATIONS")
	t.Log("=========================================")

	input := `
		xuple-space first {
			selection: fifo
			consumption: once
			retention: unlimited
		}

		xuple-space second {
			selection: random
			consumption: per-agent
			retention: duration(5m)
		}

		xuple-space third {
			selection: lifo
			consumption: limited(3)
			retention: duration(1h)
		}
	`

	result, err := ParseConstraint("test.tsd", []byte(input))
	if err != nil {
		t.Fatalf("❌ Parse error: %v", err)
	}

	program, err := ConvertResultToProgram(result)
	if err != nil {
		t.Fatalf("❌ ConvertResultToProgram error: %v", err)
	}

	if len(program.XupleSpaces) != 3 {
		t.Fatalf("❌ Expected 3 xuple-spaces, got %d", len(program.XupleSpaces))
	}

	// Verify first xuple-space
	if program.XupleSpaces[0].Name != "first" {
		t.Errorf("❌ XupleSpaces[0].Name: expected 'first', got '%s'", program.XupleSpaces[0].Name)
	}

	// Verify second xuple-space
	if program.XupleSpaces[1].Name != "second" {
		t.Errorf("❌ XupleSpaces[1].Name: expected 'second', got '%s'", program.XupleSpaces[1].Name)
	}

	// Verify third xuple-space
	if program.XupleSpaces[2].Name != "third" {
		t.Errorf("❌ XupleSpaces[2].Name: expected 'third', got '%s'", program.XupleSpaces[2].Name)
	}

	t.Log("✅ All xuple-spaces parsed correctly")
}

func TestParseXupleSpace_MixedWithOtherDeclarations(t *testing.T) {
	t.Log("🧪 TEST XUPLE-SPACE MIXED WITH OTHER DECLARATIONS")
	t.Log("=================================================")

	input := `
		type Person(#id: string, name: string, age: number)

		action notify(message: string)

		xuple-space commands {
			selection: fifo
			consumption: once
			retention: duration(1h)
		}

		rule test: {p: Person} / p.age > 18 ==> notify("Adult")
	`

	result, err := ParseConstraint("test.tsd", []byte(input))
	if err != nil {
		t.Fatalf("❌ Parse error: %v", err)
	}

	program, err := ConvertResultToProgram(result)
	if err != nil {
		t.Fatalf("❌ ConvertResultToProgram error: %v", err)
	}

	// Verify all sections are present
	if len(program.Types) != 1 {
		t.Errorf("❌ Expected 1 type, got %d", len(program.Types))
	}

	if len(program.Actions) != 1 {
		t.Errorf("❌ Expected 1 action, got %d", len(program.Actions))
	}

	if len(program.XupleSpaces) != 1 {
		t.Errorf("❌ Expected 1 xuple-space, got %d", len(program.XupleSpaces))
	}

	if len(program.Expressions) != 1 {
		t.Errorf("❌ Expected 1 expression, got %d", len(program.Expressions))
	}

	// Verify xuple-space details
	if program.XupleSpaces[0].Name != "commands" {
		t.Errorf("❌ XupleSpace name: expected 'commands', got '%s'", program.XupleSpaces[0].Name)
	}

	t.Log("✅ Mixed declarations parsed correctly")
}

func TestParseXupleSpace_DefaultValues(t *testing.T) {
	t.Log("🧪 TEST XUPLE-SPACE DEFAULT VALUES")
	t.Log("===================================")

	input := `xuple-space defaults {}`

	result, err := ParseConstraint("test.tsd", []byte(input))
	if err != nil {
		t.Fatalf("❌ Parse error: %v", err)
	}

	program, err := ConvertResultToProgram(result)
	if err != nil {
		t.Fatalf("❌ ConvertResultToProgram error: %v", err)
	}

	if len(program.XupleSpaces) != 1 {
		t.Fatalf("❌ Expected 1 xuple-space, got %d", len(program.XupleSpaces))
	}

	xs := program.XupleSpaces[0]

	// Verify default values
	if xs.SelectionPolicy != "fifo" {
		t.Errorf("❌ Default SelectionPolicy: expected 'fifo', got '%s'", xs.SelectionPolicy)
	}

	if xs.ConsumptionPolicy.Type != "once" {
		t.Errorf("❌ Default ConsumptionPolicy.Type: expected 'once', got '%s'", xs.ConsumptionPolicy.Type)
	}

	if xs.RetentionPolicy.Type != "unlimited" {
		t.Errorf("❌ Default RetentionPolicy.Type: expected 'unlimited', got '%s'", xs.RetentionPolicy.Type)
	}

	t.Log("✅ Default values correctly applied")
}
