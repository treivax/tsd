package domain

import (
	"testing"
)

func TestConstraintError(t *testing.T) {
	t.Run("NewParseError", func(t *testing.T) {
		err := NewParseError("syntax error", "test.txt", 10, 5)
		
		if err.Type != ParseError {
			t.Errorf("Expected error type ParseError, got %v", err.Type)
		}
		
		if err.Message != "syntax error" {
			t.Errorf("Expected message 'syntax error', got '%s'", err.Message)
		}
		
		if err.Context.File != "test.txt" {
			t.Errorf("Expected file 'test.txt', got '%s'", err.Context.File)
		}
		
		if err.Context.Line != 10 {
			t.Errorf("Expected line 10, got %d", err.Context.Line)
		}
		
		if err.Context.Column != 5 {
			t.Errorf("Expected column 5, got %d", err.Context.Column)
		}
		
		expectedMsg := "PARSE_ERROR at test.txt:10:5: syntax error"
		if err.Error() != expectedMsg {
			t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
		}
	})

	t.Run("NewValidationError", func(t *testing.T) {
		ctx := Context{
			Field: "test_field",
			Type:  "test_type",
		}
		err := NewValidationError("validation failed", ctx)
		
		if err.Type != ValidationError {
			t.Errorf("Expected error type ValidationError, got %v", err.Type)
		}
		
		if err.Context.Field != "test_field" {
			t.Errorf("Expected field 'test_field', got '%s'", err.Context.Field)
		}
	})

	t.Run("NewTypeMismatchError", func(t *testing.T) {
		ctx := Context{}
		err := NewTypeMismatchError("string", "number", ctx)
		
		if err.Type != TypeMismatchError {
			t.Errorf("Expected error type TypeMismatchError, got %v", err.Type)
		}
		
		if err.Context.Expected != "string" {
			t.Errorf("Expected 'string', got '%s'", err.Context.Expected)
		}
		
		if err.Context.Actual != "number" {
			t.Errorf("Expected 'number', got '%s'", err.Context.Actual)
		}
		
		expectedMsg := "type mismatch: expected string, got number"
		if err.Message != expectedMsg {
			t.Errorf("Expected message '%s', got '%s'", expectedMsg, err.Message)
		}
	})

	t.Run("NewFieldNotFoundError", func(t *testing.T) {
		ctx := Context{}
		err := NewFieldNotFoundError("age", "Person", ctx)
		
		if err.Type != FieldNotFoundError {
			t.Errorf("Expected error type FieldNotFoundError, got %v", err.Type)
		}
		
		if err.Context.Field != "age" {
			t.Errorf("Expected field 'age', got '%s'", err.Context.Field)
		}
		
		if err.Context.Type != "Person" {
			t.Errorf("Expected type 'Person', got '%s'", err.Context.Type)
		}
		
		expectedMsg := "field 'age' not found in type 'Person'"
		if err.Message != expectedMsg {
			t.Errorf("Expected message '%s', got '%s'", expectedMsg, err.Message)
		}
	})

	t.Run("NewUnknownTypeError", func(t *testing.T) {
		ctx := Context{}
		err := NewUnknownTypeError("UnknownType", ctx)
		
		if err.Type != UnknownTypeError {
			t.Errorf("Expected error type UnknownTypeError, got %v", err.Type)
		}
		
		if err.Context.Type != "UnknownType" {
			t.Errorf("Expected type 'UnknownType', got '%s'", err.Context.Type)
		}
		
		expectedMsg := "unknown type 'UnknownType'"
		if err.Message != expectedMsg {
			t.Errorf("Expected message '%s', got '%s'", expectedMsg, err.Message)
		}
	})

	t.Run("NewConstraintError", func(t *testing.T) {
		ctx := Context{Field: "constraint"}
		err := NewConstraintError("invalid constraint", ctx)
		
		if err.Type != ConstraintValidationError {
			t.Errorf("Expected error type ConstraintValidationError, got %v", err.Type)
		}
		
		if err.Message != "invalid constraint" {
			t.Errorf("Expected message 'invalid constraint', got '%s'", err.Message)
		}
	})

	t.Run("NewActionError", func(t *testing.T) {
		ctx := Context{Field: "action"}
		err := NewActionError("invalid action", ctx)
		
		if err.Type != ActionError {
			t.Errorf("Expected error type ActionError, got %v", err.Type)
		}
		
		if err.Message != "invalid action" {
			t.Errorf("Expected message 'invalid action', got '%s'", err.Message)
		}
	})
}

func TestErrorTypeCheckers(t *testing.T) {
	t.Run("IsParseError", func(t *testing.T) {
		parseErr := NewParseError("test", "file", 1, 1)
		validationErr := NewValidationError("test", Context{})
		
		if !IsParseError(parseErr) {
			t.Error("IsParseError should return true for parse error")
		}
		
		if IsParseError(validationErr) {
			t.Error("IsParseError should return false for validation error")
		}
		
		// Test avec une erreur standard
		stdErr := &Error{Type: "OTHER"}
		if IsParseError(stdErr) {
			t.Error("IsParseError should return false for other error types")
		}
	})

	t.Run("IsValidationError", func(t *testing.T) {
		validationErr := NewValidationError("test", Context{})
		typeMismatchErr := NewTypeMismatchError("string", "number", Context{})
		fieldNotFoundErr := NewFieldNotFoundError("field", "Type", Context{})
		unknownTypeErr := NewUnknownTypeError("Type", Context{})
		parseErr := NewParseError("test", "file", 1, 1)
		
		// Tous ces types sont considérés comme des erreurs de validation
		if !IsValidationError(validationErr) {
			t.Error("IsValidationError should return true for validation error")
		}
		
		if !IsValidationError(typeMismatchErr) {
			t.Error("IsValidationError should return true for type mismatch error")
		}
		
		if !IsValidationError(fieldNotFoundErr) {
			t.Error("IsValidationError should return true for field not found error")
		}
		
		if !IsValidationError(unknownTypeErr) {
			t.Error("IsValidationError should return true for unknown type error")
		}
		
		if IsValidationError(parseErr) {
			t.Error("IsValidationError should return false for parse error")
		}
	})

	t.Run("IsTypeMismatchError", func(t *testing.T) {
		typeMismatchErr := NewTypeMismatchError("string", "number", Context{})
		validationErr := NewValidationError("test", Context{})
		
		if !IsTypeMismatchError(typeMismatchErr) {
			t.Error("IsTypeMismatchError should return true for type mismatch error")
		}
		
		if IsTypeMismatchError(validationErr) {
			t.Error("IsTypeMismatchError should return false for general validation error")
		}
	})

	t.Run("IsFieldNotFoundError", func(t *testing.T) {
		fieldNotFoundErr := NewFieldNotFoundError("field", "Type", Context{})
		validationErr := NewValidationError("test", Context{})
		
		if !IsFieldNotFoundError(fieldNotFoundErr) {
			t.Error("IsFieldNotFoundError should return true for field not found error")
		}
		
		if IsFieldNotFoundError(validationErr) {
			t.Error("IsFieldNotFoundError should return false for general validation error")
		}
	})

	t.Run("IsUnknownTypeError", func(t *testing.T) {
		unknownTypeErr := NewUnknownTypeError("Type", Context{})
		validationErr := NewValidationError("test", Context{})
		
		if !IsUnknownTypeError(unknownTypeErr) {
			t.Error("IsUnknownTypeError should return true for unknown type error")
		}
		
		if IsUnknownTypeError(validationErr) {
			t.Error("IsUnknownTypeError should return false for general validation error")
		}
	})
}

func TestErrorCollection(t *testing.T) {
	t.Run("NewErrorCollection", func(t *testing.T) {
		ec := NewErrorCollection()
		
		if ec == nil {
			t.Fatal("NewErrorCollection should not return nil")
		}
		
		if len(ec.Errors) != 0 {
			t.Errorf("New error collection should be empty, got %d errors", len(ec.Errors))
		}
		
		if ec.HasErrors() {
			t.Error("New error collection should not have errors")
		}
	})

	t.Run("Add", func(t *testing.T) {
		ec := NewErrorCollection()
		err1 := NewValidationError("error 1", Context{})
		err2 := NewParseError("error 2", "file", 1, 1)
		
		ec.Add(err1)
		if len(ec.Errors) != 1 {
			t.Errorf("Expected 1 error after adding, got %d", len(ec.Errors))
		}
		
		if !ec.HasErrors() {
			t.Error("Collection should have errors after adding")
		}
		
		ec.Add(err2)
		if len(ec.Errors) != 2 {
			t.Errorf("Expected 2 errors after adding second, got %d", len(ec.Errors))
		}
	})

	t.Run("First", func(t *testing.T) {
		ec := NewErrorCollection()
		
		// Collection vide
		first := ec.First()
		if first != nil {
			t.Error("First() should return nil for empty collection")
		}
		
		// Ajouter des erreurs
		err1 := NewValidationError("first error", Context{})
		err2 := NewValidationError("second error", Context{})
		
		ec.Add(err1)
		ec.Add(err2)
		
		first = ec.First()
		if first == nil {
			t.Fatal("First() should not return nil for non-empty collection")
		}
		
		if first.Message != "first error" {
			t.Errorf("Expected first error message 'first error', got '%s'", first.Message)
		}
	})

	t.Run("Error", func(t *testing.T) {
		ec := NewErrorCollection()
		
		// Collection vide
		if ec.Error() != "no errors" {
			t.Errorf("Expected 'no errors' for empty collection, got '%s'", ec.Error())
		}
		
		// Une erreur
		err1 := NewValidationError("single error", Context{})
		ec.Add(err1)
		
		errorMsg := ec.Error()
		if errorMsg != err1.Error() {
			t.Errorf("Expected single error message, got '%s'", errorMsg)
		}
		
		// Plusieurs erreurs
		err2 := NewValidationError("second error", Context{})
		ec.Add(err2)
		
		errorMsg = ec.Error()
		expectedStart := "2 errors:"
		if len(errorMsg) < len(expectedStart) || errorMsg[:len(expectedStart)] != expectedStart {
			t.Errorf("Expected error message to start with '%s', got '%s'", expectedStart, errorMsg)
		}
	})
}

func TestErrorUnwrapping(t *testing.T) {
	t.Run("Unwrap", func(t *testing.T) {
		originalErr := NewValidationError("original", Context{})
		wrappedErr := &Error{
			Type:    ValidationError,
			Message: "wrapped",
			Context: Context{},
			Cause:   originalErr,
		}
		
		unwrapped := wrappedErr.Unwrap()
		if unwrapped != originalErr {
			t.Error("Unwrap should return the original error")
		}
		
		// Test sans cause
		noWrapErr := NewValidationError("no wrap", Context{})
		if noWrapErr.Unwrap() != nil {
			t.Error("Unwrap should return nil when no cause is set")
		}
	})

	t.Run("Is", func(t *testing.T) {
		err1 := NewValidationError("test", Context{})
		err2 := NewValidationError("other", Context{})
		err3 := NewParseError("parse", "file", 1, 1)
		
		if !err1.Is(err2) {
			t.Error("Errors of same type should match with Is()")
		}
		
		if err1.Is(err3) {
			t.Error("Errors of different types should not match with Is()")
		}
	})
}