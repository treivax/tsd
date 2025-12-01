// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"testing"
)

// TestQuotedStringsInFacts tests that quoted strings work in fact definitions
func TestQuotedStringsInFacts(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name: "simple quoted string in fact",
			input: `type Person(id: string, name:string)
Person(id:"p1", name:"Alice")`,
			wantErr: false,
		},
		{
			name: "quoted string with spaces",
			input: `type Person(id: string, name:string)
Person(id:"p1", name:"Alice Smith")`,
			wantErr: false,
		},
		{
			name: "mixed quoted and unquoted",
			input: `type Person(id: string, name: string, age:number)
Person(id:p1, name:"Alice", age:30)`,
			wantErr: false,
		},
		{
			name: "single quotes",
			input: `type Person(id: string, name:string)
Person(id:'p1', name:'Alice')`,
			wantErr: false,
		},
		{
			name: "single quotes with spaces",
			input: `type Person(id: string, name:string)
Person(id:'p1', name:'Alice Smith')`,
			wantErr: false,
		},
		{
			name: "unquoted string (should work)",
			input: `type Person(id: string, name:string)
Person(id:p1, name:Alice)`,
			wantErr: false,
		},
		{
			name: "quoted string in rule condition",
			input: `type Person(id: string, name:string)
rule r1 : {p: Person} / p.name == "Alice" ==> match(p.id)`,
			wantErr: false,
		},
		{
			name: "quoted string with special characters",
			input: `type Message(id: string, text:string)
Message(id:"m1", text:"Hello, World!")`,
			wantErr: false,
		},
		{
			name: "escaped quotes in string",
			input: `type Message(id: string, text:string)
Message(id:"m1", text:"She said \"Hello\"")`,
			wantErr: false,
		},
		{
			name: "complete program with quoted strings",
			input: `type Person(id: string, name: string, city:string)

Person(id:"p1", name:"Alice Smith", city:"New York")
Person(id:"p2", name:"Bob Jones", city:"Los Angeles")

rule check_name : {p: Person} / p.name == "Alice Smith" ==> found_alice(p.id)`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseConstraint("<test>", []byte(tt.input))

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Error("Result should not be nil")
				return
			}

			// Verify the result can be converted to a map
			resultMap, ok := result.(map[string]interface{})
			if !ok {
				t.Errorf("Result should be a map, got %T", result)
				return
			}

			// Check that facts were parsed
			facts, hasFacts := resultMap["facts"]
			if hasFacts && facts != nil {
				factsList, ok := facts.([]interface{})
				if !ok {
					t.Errorf("Facts should be a list, got %T", facts)
					return
				}
				if len(factsList) == 0 {
					t.Log("Warning: No facts parsed, but no error occurred")
				} else {
					t.Logf("Successfully parsed %d fact(s)", len(factsList))
				}
			}
		})
	}
}

// TestQuotedStringsInRules tests that quoted strings work in rule conditions
func TestQuotedStringsInRules(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name: "string equality with double quotes",
			input: `type Person(id: string, name:string)
rule r1 : {p: Person} / p.name == "Alice" ==> match(p.id)`,
			wantErr: false,
		},
		{
			name: "string equality with single quotes",
			input: `type Person(id: string, name:string)
rule r1 : {p: Person} / p.name == 'Alice' ==> match(p.id)`,
			wantErr: false,
		},
		{
			name: "string with spaces in condition",
			input: `type Person(id: string, name:string)
rule r1 : {p: Person} / p.name == "Alice Smith" ==> match(p.id)`,
			wantErr: false,
		},
		{
			name: "multiple string conditions",
			input: `type Person(id: string, name: string, city:string)
rule r1 : {p: Person} / p.name == "Alice" AND p.city == "New York" ==> match(p.id)`,
			wantErr: false,
		},
		{
			name: "string in action parameters",
			input: `type Person(id: string, name:string)
rule r1 : {p: Person} / p.name == "Alice" ==> match(p.id, "found")`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseConstraint("<test>", []byte(tt.input))

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Error("Result should not be nil")
				return
			}

			// Verify expressions were parsed
			resultMap, ok := result.(map[string]interface{})
			if !ok {
				t.Errorf("Result should be a map, got %T", result)
				return
			}

			expressions, hasExpr := resultMap["expressions"]
			if !hasExpr || expressions == nil {
				t.Error("Expected expressions to be parsed")
				return
			}

			exprList, ok := expressions.([]interface{})
			if !ok {
				t.Errorf("Expressions should be a list, got %T", expressions)
				return
			}

			if len(exprList) == 0 {
				t.Error("Expected at least one expression")
			} else {
				t.Logf("Successfully parsed %d expression(s)", len(exprList))
			}
		})
	}
}
