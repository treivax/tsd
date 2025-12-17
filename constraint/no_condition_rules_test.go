// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestParser_NoConditionRules tests parsing of rules without conditions
func TestParser_NoConditionRules(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name: "simple rule without condition",
			input: `
type Person(#id: string, name: string, age: number)

action log(message: string)

rule assertion_user : {p: Person} / ==> log(p.name)
`,
			wantErr: false,
		},
		{
			name: "multiple rules without conditions",
			input: `
type Person(#id: string, name: string)
type Product(#id: string, name: string, price: number)

action notify(msg: string)
action track(entity: string)

rule on_person : {p: Person} / ==> notify(p.name)
rule on_product : {pr: Product} / ==> track(pr.name)
`,
			wantErr: false,
		},
		{
			name: "rule without condition with multiple actions",
			input: `
type Order(#id: string, customerId: string, amount: number)

action log(msg: string)
action notify(recipient: string)
action process(orderId: string)

rule process_order : {o: Order} / ==>
    log(o.id),
    notify("admin"),
    process(o.id)
`,
			wantErr: false,
		},
		{
			name: "rule without condition accessing multiple fields",
			input: `
type Employee(#id: string, name: string, department: string, salary: number)

action record(id: string, name: string, dept: string, sal: number)

rule track_employee : {e: Employee} / ==> record(e.id, e.name, e.department, e.salary)
`,
			wantErr: false,
		},
		{
			name: "mix of rules with and without conditions",
			input: `
type Person(#id: string, name: string, age: number)

action log(msg: string)
action notify_adult(name: string)

rule log_all_persons : {p: Person} / ==> log(p.name)
rule notify_adults : {p: Person} / p.age >= 18 ==> notify_adult(p.name)
`,
			wantErr: false,
		},
		{
			name: "rule without condition with arithmetic expression in action",
			input: `
type Product(#id: string, price: number, tax: number)

action calculate(total: number)

rule compute_total : {p: Product} / ==> calculate(p.price + p.tax)
`,
			wantErr: false,
		},
		{
			name: "rule without condition with string concatenation",
			input: `
type User(#id: string, firstName: string, lastName: string)

action greet(message: string)

rule welcome_user : {u: User} / ==> greet(u.firstName)
`,
			wantErr: false,
		},
		{
			name: "multi-pattern rule without condition",
			input: `
type Order(#id: string, customerId: string)
type Customer(#id: string, name: string)

action match(orderId: string, customerName: string)

rule match_order : {o: Order} / {c: Customer} / ==> match(o.id, c.name)
`,
			wantErr: false,
		},
		{
			name: "rule without condition with boolean field",
			input: `
type Account(#id: string, active: bool)

action log(msg: string, status: bool)

rule log_account : {a: Account} / ==> log("Account status", a.active)
`,
			wantErr: false,
		},
		{
			name: "rule without condition with number field",
			input: `
type Counter(#id: string, value: number)

action report(val: number)

rule report_counter : {c: Counter} / ==> report(c.value)
`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse("test", []byte(tt.input))
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err, "Failed to parse: %v", err)
				assert.NotNil(t, result)

				// Verify the parsed program structure
				program, ok := result.(map[string]interface{})
				require.True(t, ok, "Result should be a map")

				// Check that expressions were parsed
				expressions, hasExpr := program["expressions"]
				require.True(t, hasExpr, "Program should have expressions")
				require.NotNil(t, expressions, "Expressions should not be nil")

				exprList, ok := expressions.([]interface{})
				require.True(t, ok, "Expressions should be a list")
				require.Greater(t, len(exprList), 0, "Should have at least one expression")

				// Verify first expression structure
				expr, ok := exprList[0].(map[string]interface{})
				require.True(t, ok, "Expression should be a map")

				// Check expression type
				assert.Equal(t, "expression", expr["type"])

				// Check that ruleId exists
				_, hasRuleId := expr["ruleId"]
				assert.True(t, hasRuleId, "Expression should have ruleId")

				// Check that action exists
				_, hasAction := expr["action"]
				assert.True(t, hasAction, "Expression should have action")
			}
		})
	}
}

// TestParser_NoConditionRulesValidation tests that rules without conditions are properly validated
func TestParser_NoConditionRulesValidation(t *testing.T) {
	input := `
type Person(#personId: string, name: string, age: number)

action log(message: string)

rule log_person : {p: Person} / ==> log(p.name)

Person(personId: "p1", name: "Alice", age: 30)
`

	result, err := Parse("test", []byte(input))
	require.NoError(t, err)
	assert.NotNil(t, result)

	program, ok := result.(map[string]interface{})
	require.True(t, ok)

	// Verify expressions
	expressions, ok := program["expressions"].([]interface{})
	require.True(t, ok)
	require.Len(t, expressions, 1)

	expr := expressions[0].(map[string]interface{})

	// The constraints field should be nil or not present for no-condition rules
	constraints := expr["constraints"]
	assert.Nil(t, constraints, "Constraints should be nil for no-condition rules")
}

// TestParser_NoConditionRulesWithSpaces tests parsing with various whitespace patterns
func TestParser_NoConditionRulesWithSpaces(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name: "minimal spaces",
			input: `
type Person(#personId: string, name: string)
action log(msg: string)
rule r:{p: Person}/==>log(p.name)
`,
		},
		{
			name: "extra spaces before arrow",
			input: `
type Person(#personId: string, name: string)
action log(msg: string)
rule r : {p: Person} /     ==> log(p.name)
`,
		},
		{
			name: "newlines in rule",
			input: `
type Person(#personId: string, name: string)
action log(msg: string)
rule r : {p: Person} /
  ==>
    log(p.name)
`,
		},
		{
			name: "tabs and spaces mixed",
			input: `
type Person(#personId: string, name: string)
action log(msg: string)
rule r	:	{p: Person}	/	==>	log(p.name)
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse("test", []byte(tt.input))
			require.NoError(t, err, "Failed to parse with %s", tt.name)
			assert.NotNil(t, result)
		})
	}
}

// TestParser_NoConditionRulesErrorCases tests error cases
func TestParser_NoConditionRulesErrorCases(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name: "missing action after arrow",
			input: `
type Person(#personId: string, name: string)
rule r : {p: Person} / ==>
`,
		},
		{
			name: "missing arrow",
			input: `
type Person(#personId: string, name: string)
action log(msg: string)
rule r : {p: Person} / log(p.name)
`,
		},
		{
			name: "missing pattern",
			input: `
type Person(#personId: string, name: string)
action log(msg: string)
rule r : / ==> log("test")
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse("test", []byte(tt.input))
			assert.Error(t, err, "Should fail to parse for: %s", tt.name)
		})
	}
}

// TestNoConditionRulesIntegration tests end-to-end integration with RETE
func TestNoConditionRulesIntegration(t *testing.T) {
	input := `
type Person(#personId: string, name: string, age: number)

action log(message: string)

rule log_all : {p: Person} / ==> log(p.name)

Person(personId: "1", name: "Alice", age: 30)
Person(personId: "2", name: "Bob", age: 25)
`

	// Parse the program
	result, err := Parse("test", []byte(input))
	require.NoError(t, err)
	require.NotNil(t, result)

	// Convert to Program structure
	program, ok := result.(map[string]interface{})
	require.True(t, ok)

	// Verify we have the expected components
	assert.NotNil(t, program["types"])
	assert.NotNil(t, program["actions"])
	assert.NotNil(t, program["expressions"])
	assert.NotNil(t, program["facts"])

	// Verify expressions
	expressions := program["expressions"].([]interface{})
	require.Len(t, expressions, 1)

	expr := expressions[0].(map[string]interface{})
	assert.Equal(t, "expression", expr["type"])
	assert.Equal(t, "log_all", expr["ruleId"])
	assert.Nil(t, expr["constraints"], "No-condition rule should have nil constraints")
	assert.NotNil(t, expr["action"])

	// Verify facts
	facts := program["facts"].([]interface{})
	require.Len(t, facts, 2)
}

// TestNoConditionRulesWithComplexActions tests no-condition rules with complex action expressions
func TestNoConditionRulesWithComplexActions(t *testing.T) {
	input := `
type Product(#productId: string, name: string, price: number, quantity: number)

action log(msg: string)
action notify(recipient: string, message: string)
action calculate(total: number, tax: number)

rule process_product : {p: Product} / ==>
    log(p.name),
    notify("admin", p.name),
    calculate(p.price * p.quantity, p.price * p.quantity * 0.2)

Product(productId: "1", name: "Laptop", price: 1000, quantity: 5)
`

	result, err := Parse("test", []byte(input))
	require.NoError(t, err)
	require.NotNil(t, result)

	program, ok := result.(map[string]interface{})
	require.True(t, ok)

	expressions := program["expressions"].([]interface{})
	require.Len(t, expressions, 1)

	expr := expressions[0].(map[string]interface{})
	assert.Nil(t, expr["constraints"], "Should have nil constraints")
	assert.NotNil(t, expr["action"], "Should have action")
}

// TestNoConditionRulesWithAggregation tests no-condition rules cannot be used with aggregation patterns
func TestNoConditionRulesWithAggregation(t *testing.T) {
	// This tests that we can parse multi-pattern rules without conditions
	input := `
type Employee(#id: string, name: string, salary: number)
type Department(#id: string, name: string)

action report(deptName: string, avgSalary: number)

rule dept_avg : {d: Department, avg_sal: AVG(e.salary)} / {e: Employee} / ==> report(d.name, avg_sal)
`

	result, err := Parse("test", []byte(input))
	require.NoError(t, err, "Should parse multi-pattern rule without explicit join condition")
	require.NotNil(t, result)

	program, ok := result.(map[string]interface{})
	require.True(t, ok)

	expressions := program["expressions"].([]interface{})
	require.Len(t, expressions, 1)

	expr := expressions[0].(map[string]interface{})
	assert.Nil(t, expr["constraints"], "Should have nil constraints")
	assert.NotNil(t, expr["patterns"], "Should have multiple patterns")
}
