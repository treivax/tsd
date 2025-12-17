// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License

package constraint

import (
	"strings"
	"testing"
)

// TestAggregationWithJoinSyntax tests the new parser support for aggregation with join patterns
func TestAggregationWithJoinSyntax(t *testing.T) {
	input := `type Department(#id: string, name:string)
type Employee(#id: string, deptId: string, salary:number)

rule dept_avg_salary : {d: Department, avg_sal: AVG(e.salary)} / {e: Employee} / e.deptId == d.id ==> print("Avg salary")
`

	result, err := Parse("test", []byte(input))
	if err != nil {
		t.Fatalf("Failed to parse aggregation with join syntax: %v", err)
	}

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected result to be map[string]interface{}, got %T", result)
	}

	// Verify types were parsed
	types, ok := resultMap["types"].([]interface{})
	if !ok || len(types) != 2 {
		t.Errorf("Expected 2 types, got %v", types)
	}

	// Verify expressions were parsed
	expressions, ok := resultMap["expressions"].([]interface{})
	if !ok || len(expressions) != 1 {
		t.Fatalf("Expected 1 expression, got %v", expressions)
	}

	// Verify the expression structure
	expr, ok := expressions[0].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected expression to be map, got %T", expressions[0])
	}

	// Check rule ID
	if ruleID, ok := expr["ruleId"].(string); !ok || ruleID != "dept_avg_salary" {
		t.Errorf("Expected ruleId 'dept_avg_salary', got %v", expr["ruleId"])
	}

	// Check that patterns field exists (new multi-pattern syntax)
	patterns, hasPatterns := expr["patterns"].([]interface{})
	if !hasPatterns {
		t.Fatalf("Expected 'patterns' field in expression, got: %+v", expr)
	}

	// Should have 2 pattern blocks
	if len(patterns) != 2 {
		t.Errorf("Expected 2 pattern blocks, got %d", len(patterns))
	}

	// Verify first pattern block (aggregation result)
	pattern1, ok := patterns[0].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected first pattern to be map, got %T", patterns[0])
	}

	vars1, ok := pattern1["variables"].([]interface{})
	if !ok || len(vars1) != 2 {
		t.Errorf("Expected 2 variables in first pattern, got %v", vars1)
	}

	// Check first variable (d: Department)
	var1, ok := vars1[0].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected variable to be map, got %T", vars1[0])
	}
	if var1["name"] != "d" || var1["dataType"] != "Department" {
		t.Errorf("Expected variable 'd: Department', got name=%v, dataType=%v", var1["name"], var1["dataType"])
	}

	// Check second variable (avg_sal: AVG(e.salary)) - this is an aggregation variable
	var2, ok := vars1[1].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected aggregation variable to be map, got %T", vars1[1])
	}
	if var2["type"] != "aggregationVariable" {
		t.Errorf("Expected type 'aggregationVariable', got %v", var2["type"])
	}
	if var2["name"] != "avg_sal" {
		t.Errorf("Expected name 'avg_sal', got %v", var2["name"])
	}
	if var2["function"] != "AVG" {
		t.Errorf("Expected function 'AVG', got %v", var2["function"])
	}

	// Check the field access in aggregation
	field, ok := var2["field"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected field to be map, got %T", var2["field"])
	}
	if field["object"] != "e" || field["field"] != "salary" {
		t.Errorf("Expected field access 'e.salary', got object=%v, field=%v", field["object"], field["field"])
	}

	// Verify second pattern block (source pattern)
	pattern2, ok := patterns[1].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected second pattern to be map, got %T", patterns[1])
	}

	vars2, ok := pattern2["variables"].([]interface{})
	if !ok || len(vars2) != 1 {
		t.Errorf("Expected 1 variable in second pattern, got %v", vars2)
	}

	// Check the Employee variable
	empVar, ok := vars2[0].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected employee variable to be map, got %T", vars2[0])
	}
	if empVar["name"] != "e" || empVar["dataType"] != "Employee" {
		t.Errorf("Expected variable 'e: Employee', got name=%v, dataType=%v", empVar["name"], empVar["dataType"])
	}

	// Verify constraints (join condition)
	constraints, ok := expr["constraints"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected constraints to be map, got %T", expr["constraints"])
	}

	// The constraint should be a comparison: e.deptId == d.id
	if constraints["type"] != "comparison" {
		t.Errorf("Expected constraint type 'comparison', got %v", constraints["type"])
	}

	t.Logf("✅ Successfully parsed aggregation with join syntax")
}

// TestAggregationWithMultipleAggregates tests multiple aggregation variables
func TestAggregationWithMultipleAggregates(t *testing.T) {
	input := `type Department(#id:string)
type Employee(#id: string, deptId: string, salary:number)

rule dept_stats : {d: Department, avg_sal: AVG(e.salary), max_sal: MAX(e.salary)} / {e: Employee} / e.deptId == d.id ==> print("Stats")
`

	result, err := Parse("test", []byte(input))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	resultMap := result.(map[string]interface{})
	expressions := resultMap["expressions"].([]interface{})
	expr := expressions[0].(map[string]interface{})
	patterns := expr["patterns"].([]interface{})
	pattern1 := patterns[0].(map[string]interface{})
	vars := pattern1["variables"].([]interface{})

	// Should have 3 variables: d, avg_sal, max_sal
	if len(vars) != 3 {
		t.Errorf("Expected 3 variables, got %d", len(vars))
	}

	// Check aggregation variables
	avgVar := vars[1].(map[string]interface{})
	if avgVar["type"] != "aggregationVariable" || avgVar["function"] != "AVG" {
		t.Errorf("Expected AVG aggregation variable, got %+v", avgVar)
	}

	maxVar := vars[2].(map[string]interface{})
	if maxVar["type"] != "aggregationVariable" || maxVar["function"] != "MAX" {
		t.Errorf("Expected MAX aggregation variable, got %+v", maxVar)
	}

	t.Logf("✅ Successfully parsed multiple aggregation variables")
}

// TestBackwardCompatibilitySinglePattern ensures old syntax still works
func TestBackwardCompatibilitySinglePattern(t *testing.T) {
	input := `type Person(#id: string, age:number)

rule adults : {p: Person} / p.age >= 18 ==> print("Adult")
`

	result, err := Parse("test", []byte(input))
	if err != nil {
		t.Fatalf("Failed to parse backward compatible syntax: %v", err)
	}

	resultMap := result.(map[string]interface{})
	expressions := resultMap["expressions"].([]interface{})
	expr := expressions[0].(map[string]interface{})

	// For single pattern, should use 'set' field (backward compatibility)
	if _, hasSet := expr["set"]; !hasSet {
		t.Errorf("Expected 'set' field for backward compatibility, got: %+v", expr)
	}

	// Should NOT have 'patterns' field
	if _, hasPatterns := expr["patterns"]; hasPatterns {
		t.Errorf("Should not have 'patterns' field for single pattern rules")
	}

	t.Logf("✅ Backward compatibility maintained for single pattern rules")
}

// TestAggregationVariableSyntaxVariants tests different aggregation functions
func TestAggregationVariableSyntaxVariants(t *testing.T) {
	tests := []struct {
		name     string
		funcName string
		input    string
	}{
		{
			name:     "AVG function",
			funcName: "AVG",
			input:    `rule test : {d: Dept, avg: AVG(e.salary)} / {e: Emp} / e.deptId == d.id ==> print("x")`,
		},
		{
			name:     "SUM function",
			funcName: "SUM",
			input:    `rule test : {d: Dept, total: SUM(e.salary)} / {e: Emp} / e.deptId == d.id ==> print("x")`,
		},
		{
			name:     "COUNT function",
			funcName: "COUNT",
			input:    `rule test : {d: Dept, cnt: COUNT(e.id)} / {e: Emp} / e.deptId == d.id ==> print("x")`,
		},
		{
			name:     "MIN function",
			funcName: "MIN",
			input:    `rule test : {d: Dept, min: MIN(e.salary)} / {e: Emp} / e.deptId == d.id ==> print("x")`,
		},
		{
			name:     "MAX function",
			funcName: "MAX",
			input:    `rule test : {d: Dept, max: MAX(e.salary)} / {e: Emp} / e.deptId == d.id ==> print("x")`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fullInput := `type Dept(#id:string)
type Emp(#id: string, deptId: string, salary:number)

` + tt.input

			result, err := Parse("test", []byte(fullInput))
			if err != nil {
				t.Fatalf("Failed to parse %s: %v", tt.funcName, err)
			}

			resultMap := result.(map[string]interface{})
			expressions := resultMap["expressions"].([]interface{})
			expr := expressions[0].(map[string]interface{})
			patterns := expr["patterns"].([]interface{})
			pattern1 := patterns[0].(map[string]interface{})
			vars := pattern1["variables"].([]interface{})

			// Second variable should be the aggregation
			aggVar := vars[1].(map[string]interface{})
			if aggVar["function"] != tt.funcName {
				t.Errorf("Expected function '%s', got %v", tt.funcName, aggVar["function"])
			}

			t.Logf("✅ %s function parsed correctly", tt.funcName)
		})
	}
}

// TestParseError ensures invalid syntax is rejected
func TestAggregationJoinParseError(t *testing.T) {
	// Missing closing paren in aggregation
	input := `type Dept(#id:string)
type Emp(#id: string, salary:number)

rule bad : {d: Dept, avg: AVG(e.salary} / {e: Emp} / e.deptId == d.id ==> print("x")
`

	_, err := Parse("test", []byte(input))
	if err == nil {
		t.Error("Expected parse error for invalid aggregation syntax, got nil")
	}

	if !strings.Contains(err.Error(), "no match found") {
		t.Logf("Parse error as expected: %v", err)
	}
}
