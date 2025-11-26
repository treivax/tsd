// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package domain

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

// ===== Program Tests =====

func TestNewProgram(t *testing.T) {
	program := NewProgram()

	if program == nil {
		t.Fatal("NewProgram should not return nil")
	}
	if program.Types == nil {
		t.Error("Types should be initialized")
	}
	if program.Expressions == nil {
		t.Error("Expressions should be initialized")
	}
	if program.Metadata == nil {
		t.Error("Metadata should be initialized")
	}
	if program.Metadata.Version != "1.0" {
		t.Errorf("Expected version '1.0', got '%s'", program.Metadata.Version)
	}
	if program.Metadata.CreatedAt.IsZero() {
		t.Error("CreatedAt should be set")
	}
	if time.Since(program.Metadata.CreatedAt) > time.Second {
		t.Error("CreatedAt should be recent")
	}
}

func TestProgram_GetTypeByName(t *testing.T) {
	program := NewProgram()
	typeDef := NewTypeDefinition("Person")
	typeDef.AddField("name", "string")
	program.Types = append(program.Types, *typeDef)

	// Existing type
	found := program.GetTypeByName("Person")
	if found == nil {
		t.Error("GetTypeByName should find existing type")
	}
	if found != nil && found.Name != "Person" {
		t.Errorf("Expected type name 'Person', got '%s'", found.Name)
	}

	// Non-existing type
	notFound := program.GetTypeByName("NonExistent")
	if notFound != nil {
		t.Error("GetTypeByName should return nil for non-existing type")
	}
}

func TestProgram_String(t *testing.T) {
	program := NewProgram()
	typeDef := NewTypeDefinition("Person")
	typeDef.AddField("name", "string")
	program.Types = append(program.Types, *typeDef)

	result := program.String()
	if result == "" {
		t.Error("String() should return non-empty string")
	}

	// Should be valid JSON
	var data map[string]interface{}
	err := json.Unmarshal([]byte(result), &data)
	if err != nil {
		t.Errorf("String() should return valid JSON: %v", err)
	}
}

func TestProgram_String_EmptyProgram(t *testing.T) {
	program := &Program{
		Types:       []TypeDefinition{},
		Expressions: []Expression{},
	}

	result := program.String()
	if result == "" {
		t.Error("String() should return non-empty string for empty program")
	}
}

// ===== TypeDefinition Tests =====

func TestNewTypeDefinition(t *testing.T) {
	typeDef := NewTypeDefinition("Person")

	if typeDef == nil {
		t.Fatal("NewTypeDefinition should not return nil")
	}
	if typeDef.Type != "typeDefinition" {
		t.Errorf("Expected Type 'typeDefinition', got '%s'", typeDef.Type)
	}
	if typeDef.Name != "Person" {
		t.Errorf("Expected Name 'Person', got '%s'", typeDef.Name)
	}
	if typeDef.Fields == nil {
		t.Error("Fields should be initialized")
	}
	if len(typeDef.Fields) != 0 {
		t.Errorf("Expected 0 fields, got %d", len(typeDef.Fields))
	}
}

func TestTypeDefinition_AddField(t *testing.T) {
	typeDef := NewTypeDefinition("Person")

	typeDef.AddField("name", "string")
	if len(typeDef.Fields) != 1 {
		t.Errorf("Expected 1 field, got %d", len(typeDef.Fields))
	}

	typeDef.AddField("age", "integer")
	if len(typeDef.Fields) != 2 {
		t.Errorf("Expected 2 fields, got %d", len(typeDef.Fields))
	}

	// Verify fields
	if typeDef.Fields[0].Name != "name" || typeDef.Fields[0].Type != "string" {
		t.Error("First field not correctly added")
	}
	if typeDef.Fields[1].Name != "age" || typeDef.Fields[1].Type != "integer" {
		t.Error("Second field not correctly added")
	}
}

func TestTypeDefinition_AddField_AllTypes(t *testing.T) {
	typeDef := NewTypeDefinition("AllTypes")

	typeDef.AddField("str", "string")
	typeDef.AddField("num", "number")
	typeDef.AddField("int", "integer")
	typeDef.AddField("bool", "bool")

	if len(typeDef.Fields) != 4 {
		t.Errorf("Expected 4 fields, got %d", len(typeDef.Fields))
	}

	expectedTypes := map[string]string{
		"str":  "string",
		"num":  "number",
		"int":  "integer",
		"bool": "bool",
	}

	for _, field := range typeDef.Fields {
		expectedType, exists := expectedTypes[field.Name]
		if !exists {
			t.Errorf("Unexpected field name '%s'", field.Name)
			continue
		}
		if field.Type != expectedType {
			t.Errorf("Field '%s': expected type '%s', got '%s'", field.Name, expectedType, field.Type)
		}
	}
}

func TestTypeDefinition_GetFieldByName(t *testing.T) {
	typeDef := NewTypeDefinition("Person")
	typeDef.AddField("name", "string")
	typeDef.AddField("age", "integer")
	typeDef.AddField("active", "bool")

	// Existing fields
	nameField := typeDef.GetFieldByName("name")
	if nameField == nil {
		t.Error("GetFieldByName should find 'name' field")
	}
	if nameField != nil && nameField.Type != "string" {
		t.Errorf("Expected type 'string', got '%s'", nameField.Type)
	}

	ageField := typeDef.GetFieldByName("age")
	if ageField == nil {
		t.Error("GetFieldByName should find 'age' field")
	}
	if ageField != nil && ageField.Type != "integer" {
		t.Errorf("Expected type 'integer', got '%s'", ageField.Type)
	}

	// Non-existing field
	notFound := typeDef.GetFieldByName("nonexistent")
	if notFound != nil {
		t.Error("GetFieldByName should return nil for non-existing field")
	}
}

func TestTypeDefinition_HasField(t *testing.T) {
	typeDef := NewTypeDefinition("Person")
	typeDef.AddField("name", "string")
	typeDef.AddField("age", "integer")

	if !typeDef.HasField("name") {
		t.Error("HasField should return true for 'name'")
	}
	if !typeDef.HasField("age") {
		t.Error("HasField should return true for 'age'")
	}
	if typeDef.HasField("nonexistent") {
		t.Error("HasField should return false for 'nonexistent'")
	}
}

// ===== Expression Tests =====

func TestNewExpression(t *testing.T) {
	expr := NewExpression()

	if expr == nil {
		t.Fatal("NewExpression should not return nil")
	}
	if expr.Type != "expression" {
		t.Errorf("Expected Type 'expression', got '%s'", expr.Type)
	}
	if expr.Set.Type != "set" {
		t.Errorf("Expected Set.Type 'set', got '%s'", expr.Set.Type)
	}
	if expr.Set.Variables == nil {
		t.Error("Variables should be initialized")
	}
	if len(expr.Set.Variables) != 0 {
		t.Errorf("Expected 0 variables, got %d", len(expr.Set.Variables))
	}
}

func TestExpression_AddVariable(t *testing.T) {
	expr := NewExpression()

	expr.AddVariable("p", "Person")
	if len(expr.Set.Variables) != 1 {
		t.Errorf("Expected 1 variable, got %d", len(expr.Set.Variables))
	}

	expr.AddVariable("pr", "Product")
	if len(expr.Set.Variables) != 2 {
		t.Errorf("Expected 2 variables, got %d", len(expr.Set.Variables))
	}

	// Verify variables
	if expr.Set.Variables[0].Name != "p" || expr.Set.Variables[0].DataType != "Person" {
		t.Error("First variable not correctly added")
	}
	if expr.Set.Variables[1].Name != "pr" || expr.Set.Variables[1].DataType != "Product" {
		t.Error("Second variable not correctly added")
	}
	if expr.Set.Variables[0].Type != "typedVariable" {
		t.Error("Variable Type should be 'typedVariable'")
	}
}

// ===== Constraint Tests =====

func TestNewConstraint(t *testing.T) {
	left := &FieldAccess{Object: "p", Field: "age"}
	right := &IntegerLiteral{Value: 18}

	constraint := NewConstraint(left, ">=", right)

	if constraint == nil {
		t.Fatal("NewConstraint should not return nil")
	}
	if constraint.Type != "constraint" {
		t.Errorf("Expected Type 'constraint', got '%s'", constraint.Type)
	}
	if constraint.Operator != ">=" {
		t.Errorf("Expected Operator '>=', got '%s'", constraint.Operator)
	}
	if constraint.Left == nil {
		t.Error("Left should not be nil")
	}
	if constraint.Right == nil {
		t.Error("Right should not be nil")
	}
}

// ===== FieldAccess Tests =====

func TestNewFieldAccess(t *testing.T) {
	fa := NewFieldAccess("person", "name")

	if fa == nil {
		t.Fatal("NewFieldAccess should not return nil")
	}
	if fa.Type != "fieldAccess" {
		t.Errorf("Expected Type 'fieldAccess', got '%s'", fa.Type)
	}
	if fa.Object != "person" {
		t.Errorf("Expected Object 'person', got '%s'", fa.Object)
	}
	if fa.Field != "name" {
		t.Errorf("Expected Field 'name', got '%s'", fa.Field)
	}
}

// ===== Action Tests =====

func TestNewAction(t *testing.T) {
	tests := []struct {
		name     string
		jobName  string
		args     []interface{}
		wantArgs int
	}{
		{
			name:     "action without args",
			jobName:  "notify",
			args:     []interface{}{},
			wantArgs: 0,
		},
		{
			name:     "action with string args",
			jobName:  "sendEmail",
			args:     []interface{}{"user@test.com", "Subject"},
			wantArgs: 2,
		},
		{
			name:     "action with mixed args",
			jobName:  "process",
			args:     []interface{}{"id", 42, true, 3.14},
			wantArgs: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			action := NewAction(tt.jobName, tt.args...)

			if action == nil {
				t.Fatal("NewAction should not return nil")
			}
			if action.Type != "action" {
				t.Errorf("Expected Type 'action', got '%s'", action.Type)
			}
			if action.Job.Type != "jobCall" {
				t.Errorf("Expected Job.Type 'jobCall', got '%s'", action.Job.Type)
			}
			if action.Job.Name != tt.jobName {
				t.Errorf("Expected Job.Name '%s', got '%s'", tt.jobName, action.Job.Name)
			}
			if len(action.Job.Args) != tt.wantArgs {
				t.Errorf("Expected %d args, got %d", tt.wantArgs, len(action.Job.Args))
			}
		})
	}
}

// ===== Validation Helper Tests =====

func TestIsValidOperator(t *testing.T) {
	validOperators := []string{"==", "!=", "<", ">", "<=", ">=", "AND", "OR", "NOT", "+", "-", "*", "/", "%"}
	invalidOperators := []string{"INVALID", "&&", "||", "===", "!==", "++", "--", "xor", ""}

	for _, op := range validOperators {
		if !IsValidOperator(op) {
			t.Errorf("IsValidOperator should return true for '%s'", op)
		}
	}

	for _, op := range invalidOperators {
		if IsValidOperator(op) {
			t.Errorf("IsValidOperator should return false for '%s'", op)
		}
	}
}

func TestIsValidType(t *testing.T) {
	validTypes := []string{"string", "number", "bool", "integer"}
	invalidTypes := []string{"invalid", "float", "double", "char", "byte", "array", ""}

	for _, tp := range validTypes {
		if !IsValidType(tp) {
			t.Errorf("IsValidType should return true for '%s'", tp)
		}
	}

	for _, tp := range invalidTypes {
		if IsValidType(tp) {
			t.Errorf("IsValidType should return false for '%s'", tp)
		}
	}
}

// ===== Error Tests =====

func TestNewValidationError(t *testing.T) {
	ctx := Context{
		Field:    "age",
		Expected: "positive",
		Actual:   "negative",
	}
	err := NewValidationError("invalid age", ctx)

	if err == nil {
		t.Fatal("NewValidationError should not return nil")
	}
	if err.Type != ValidationError {
		t.Errorf("Expected Type ValidationError, got %v", err.Type)
	}
	if err.Message != "invalid age" {
		t.Errorf("Expected Message 'invalid age', got '%s'", err.Message)
	}
	if err.Context.Field != "age" {
		t.Errorf("Expected Context.Field 'age', got '%s'", err.Context.Field)
	}

	errMsg := err.Error()
	if errMsg == "" {
		t.Error("Error() should return non-empty string")
	}
	if !strings.Contains(errMsg, "invalid age") {
		t.Error("Error message should contain 'invalid age'")
	}
}

func TestNewTypeMismatchError(t *testing.T) {
	ctx := Context{Variable: "x"}
	err := NewTypeMismatchError("string", "integer", ctx)

	if err == nil {
		t.Fatal("NewTypeMismatchError should not return nil")
	}
	if err.Type != TypeMismatchError {
		t.Errorf("Expected Type TypeMismatchError, got %v", err.Type)
	}
	if err.Context.Expected != "string" {
		t.Errorf("Expected Context.Expected 'string', got '%s'", err.Context.Expected)
	}
	if err.Context.Actual != "integer" {
		t.Errorf("Expected Context.Actual 'integer', got '%s'", err.Context.Actual)
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, "type mismatch") {
		t.Error("Error message should contain 'type mismatch'")
	}
	if !strings.Contains(errMsg, "string") || !strings.Contains(errMsg, "integer") {
		t.Error("Error message should contain both types")
	}
}

func TestNewFieldNotFoundError(t *testing.T) {
	ctx := Context{Variable: "p"}
	err := NewFieldNotFoundError("name", "Person", ctx)

	if err == nil {
		t.Fatal("NewFieldNotFoundError should not return nil")
	}
	if err.Type != FieldNotFoundError {
		t.Errorf("Expected Type FieldNotFoundError, got %v", err.Type)
	}
	if err.Context.Field != "name" {
		t.Errorf("Expected Context.Field 'name', got '%s'", err.Context.Field)
	}
	if err.Context.Type != "Person" {
		t.Errorf("Expected Context.Type 'Person', got '%s'", err.Context.Type)
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, "field") && !strings.Contains(errMsg, "not found") {
		t.Error("Error message should indicate field not found")
	}
}

func TestNewUnknownTypeError(t *testing.T) {
	ctx := Context{Field: "dataType"}
	err := NewUnknownTypeError("UnknownType", ctx)

	if err == nil {
		t.Fatal("NewUnknownTypeError should not return nil")
	}
	if err.Type != UnknownTypeError {
		t.Errorf("Expected Type UnknownTypeError, got %v", err.Type)
	}
	if err.Context.Type != "UnknownType" {
		t.Errorf("Expected Context.Type 'UnknownType', got '%s'", err.Context.Type)
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, "unknown type") {
		t.Error("Error message should contain 'unknown type'")
	}
}

func TestNewConstraintError(t *testing.T) {
	ctx := Context{Field: "constraint"}
	err := NewConstraintError("invalid constraint", ctx)

	if err == nil {
		t.Fatal("NewConstraintError should not return nil")
	}
	if err.Type != ConstraintValidationError {
		t.Errorf("Expected Type ConstraintValidationError, got %v", err.Type)
	}
}

func TestNewActionError(t *testing.T) {
	ctx := Context{Field: "job"}
	err := NewActionError("invalid action", ctx)

	if err == nil {
		t.Fatal("NewActionError should not return nil")
	}
	if err.Type != ActionError {
		t.Errorf("Expected Type ActionError, got %v", err.Type)
	}
}

func TestNewParseError(t *testing.T) {
	err := NewParseError("syntax error", "test.tsd", 10, 5)

	if err == nil {
		t.Fatal("NewParseError should not return nil")
	}
	if err.Type != ParseError {
		t.Errorf("Expected Type ParseError, got %v", err.Type)
	}
	if err.Context.File != "test.tsd" {
		t.Errorf("Expected Context.File 'test.tsd', got '%s'", err.Context.File)
	}
	if err.Context.Line != 10 {
		t.Errorf("Expected Context.Line 10, got %d", err.Context.Line)
	}
	if err.Context.Column != 5 {
		t.Errorf("Expected Context.Column 5, got %d", err.Context.Column)
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, "test.tsd") {
		t.Error("Error message should contain filename")
	}
	if !strings.Contains(errMsg, "10") || !strings.Contains(errMsg, "5") {
		t.Error("Error message should contain line and column")
	}
}

func TestError_Unwrap(t *testing.T) {
	cause := NewValidationError("cause", Context{})
	err := &Error{
		Type:    ValidationError,
		Message: "wrapper",
		Cause:   cause,
	}

	unwrapped := err.Unwrap()
	if unwrapped != cause {
		t.Error("Unwrap should return the cause")
	}
}

func TestError_Is(t *testing.T) {
	err1 := &Error{Type: ValidationError, Message: "error1"}
	err2 := &Error{Type: ValidationError, Message: "error2"}
	err3 := &Error{Type: TypeMismatchError, Message: "error3"}

	if !err1.Is(err2) {
		t.Error("Is should return true for same error type")
	}
	if err1.Is(err3) {
		t.Error("Is should return false for different error type")
	}
}

// ===== Error Type Checkers =====

func TestIsParseError(t *testing.T) {
	parseErr := NewParseError("syntax error", "test.tsd", 1, 1)
	validationErr := NewValidationError("invalid", Context{})

	if !IsParseError(parseErr) {
		t.Error("IsParseError should return true for ParseError")
	}
	if IsParseError(validationErr) {
		t.Error("IsParseError should return false for non-ParseError")
	}
}

func TestIsValidationError(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		wantVal bool
	}{
		{"ValidationError", NewValidationError("msg", Context{}), true},
		{"TypeMismatchError", NewTypeMismatchError("a", "b", Context{}), true},
		{"FieldNotFoundError", NewFieldNotFoundError("f", "t", Context{}), true},
		{"UnknownTypeError", NewUnknownTypeError("t", Context{}), true},
		{"ParseError", NewParseError("msg", "f", 1, 1), false},
		{"ActionError", NewActionError("msg", Context{}), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidationError(tt.err)
			if result != tt.wantVal {
				t.Errorf("IsValidationError() = %v, want %v", result, tt.wantVal)
			}
		})
	}
}

func TestIsTypeMismatchError(t *testing.T) {
	typeMismatchErr := NewTypeMismatchError("a", "b", Context{})
	validationErr := NewValidationError("invalid", Context{})

	if !IsTypeMismatchError(typeMismatchErr) {
		t.Error("IsTypeMismatchError should return true for TypeMismatchError")
	}
	if IsTypeMismatchError(validationErr) {
		t.Error("IsTypeMismatchError should return false for non-TypeMismatchError")
	}
}

func TestIsFieldNotFoundError(t *testing.T) {
	fieldErr := NewFieldNotFoundError("f", "t", Context{})
	validationErr := NewValidationError("invalid", Context{})

	if !IsFieldNotFoundError(fieldErr) {
		t.Error("IsFieldNotFoundError should return true for FieldNotFoundError")
	}
	if IsFieldNotFoundError(validationErr) {
		t.Error("IsFieldNotFoundError should return false for non-FieldNotFoundError")
	}
}

func TestIsUnknownTypeError(t *testing.T) {
	unknownTypeErr := NewUnknownTypeError("t", Context{})
	validationErr := NewValidationError("invalid", Context{})

	if !IsUnknownTypeError(unknownTypeErr) {
		t.Error("IsUnknownTypeError should return true for UnknownTypeError")
	}
	if IsUnknownTypeError(validationErr) {
		t.Error("IsUnknownTypeError should return false for non-UnknownTypeError")
	}
}

// ===== ErrorCollection Tests =====

func TestNewErrorCollection(t *testing.T) {
	ec := NewErrorCollection()

	if ec == nil {
		t.Fatal("NewErrorCollection should not return nil")
	}
	if ec.Errors == nil {
		t.Error("Errors slice should be initialized")
	}
	if len(ec.Errors) != 0 {
		t.Errorf("Expected 0 errors, got %d", len(ec.Errors))
	}
}

func TestErrorCollection_Add(t *testing.T) {
	ec := NewErrorCollection()

	err1 := NewValidationError("error1", Context{})
	err2 := NewValidationError("error2", Context{})

	ec.Add(err1)
	if len(ec.Errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(ec.Errors))
	}

	ec.Add(err2)
	if len(ec.Errors) != 2 {
		t.Errorf("Expected 2 errors, got %d", len(ec.Errors))
	}
}

func TestErrorCollection_HasErrors(t *testing.T) {
	ec := NewErrorCollection()

	if ec.HasErrors() {
		t.Error("HasErrors should return false for empty collection")
	}

	ec.Add(NewValidationError("error", Context{}))

	if !ec.HasErrors() {
		t.Error("HasErrors should return true after adding error")
	}
}

func TestErrorCollection_Error(t *testing.T) {
	ec := NewErrorCollection()

	// No errors
	msg := ec.Error()
	if !strings.Contains(msg, "no errors") {
		t.Error("Error message should indicate no errors")
	}

	// One error
	ec.Add(NewValidationError("error1", Context{}))
	msg = ec.Error()
	if !strings.Contains(msg, "error1") {
		t.Error("Error message should contain the error message")
	}

	// Multiple errors
	ec.Add(NewValidationError("error2", Context{}))
	ec.Add(NewValidationError("error3", Context{}))
	msg = ec.Error()
	if !strings.Contains(msg, "3 errors") {
		t.Error("Error message should indicate count of errors")
	}
	if !strings.Contains(msg, "and 2 more") {
		t.Error("Error message should indicate additional errors")
	}
}

func TestErrorCollection_First(t *testing.T) {
	ec := NewErrorCollection()

	// No errors
	first := ec.First()
	if first != nil {
		t.Error("First should return nil for empty collection")
	}

	// With errors
	err1 := NewValidationError("error1", Context{})
	err2 := NewValidationError("error2", Context{})

	ec.Add(err1)
	ec.Add(err2)

	first = ec.First()
	if first == nil {
		t.Fatal("First should return first error")
	}
	if first.Message != "error1" {
		t.Errorf("Expected 'error1', got '%s'", first.Message)
	}
}
