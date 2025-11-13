package domain

import (
	"encoding/json"
	"testing"
	"time"
)

func TestProgram(t *testing.T) {
	t.Run("NewProgram", func(t *testing.T) {
		prog := NewProgram()
		if prog == nil {
			t.Fatal("NewProgram() returned nil")
		}

		if prog.Metadata == nil {
			t.Fatal("Program metadata is nil")
		}

		if prog.Metadata.Version != "1.0" {
			t.Errorf("Expected version 1.0, got %s", prog.Metadata.Version)
		}

		if len(prog.Types) != 0 {
			t.Errorf("Expected empty types slice, got %d", len(prog.Types))
		}

		if len(prog.Expressions) != 0 {
			t.Errorf("Expected empty expressions slice, got %d", len(prog.Expressions))
		}
	})

	t.Run("GetTypeByName", func(t *testing.T) {
		prog := NewProgram()

		// Type qui n'existe pas
		typeDef := prog.GetTypeByName("NonExistent")
		if typeDef != nil {
			t.Error("GetTypeByName should return nil for non-existent type")
		}

		// Ajouter un type et le retrouver
		personType := TypeDefinition{
			Name: "Person",
			Type: "typeDefinition",
			Fields: []Field{
				{Name: "name", Type: "string"},
				{Name: "age", Type: "number"},
			},
		}
		prog.Types = append(prog.Types, personType)

		found := prog.GetTypeByName("Person")
		if found == nil {
			t.Fatal("GetTypeByName should return the Person type")
		}

		if found.Name != "Person" {
			t.Errorf("Expected type name 'Person', got '%s'", found.Name)
		}

		if len(found.Fields) != 2 {
			t.Errorf("Expected 2 fields, got %d", len(found.Fields))
		}
	})

	t.Run("String", func(t *testing.T) {
		prog := NewProgram()
		str := prog.String()

		// Vérifier que c'est du JSON valide
		var parsed map[string]interface{}
		if err := json.Unmarshal([]byte(str), &parsed); err != nil {
			t.Errorf("Program.String() should return valid JSON: %v", err)
		}

		// Vérifier que les champs essentiels sont présents
		if _, ok := parsed["types"]; !ok {
			t.Error("JSON should contain 'types' field")
		}

		if _, ok := parsed["expressions"]; !ok {
			t.Error("JSON should contain 'expressions' field")
		}
	})
}

func TestTypeDefinition(t *testing.T) {
	t.Run("NewTypeDefinition", func(t *testing.T) {
		typeDef := NewTypeDefinition("TestType")

		if typeDef.Name != "TestType" {
			t.Errorf("Expected name 'TestType', got '%s'", typeDef.Name)
		}

		if typeDef.Type != "typeDefinition" {
			t.Errorf("Expected type 'typeDefinition', got '%s'", typeDef.Type)
		}

		if len(typeDef.Fields) != 0 {
			t.Errorf("Expected empty fields slice, got %d", len(typeDef.Fields))
		}
	})

	t.Run("AddField", func(t *testing.T) {
		typeDef := NewTypeDefinition("Person")

		typeDef.AddField("name", "string")
		typeDef.AddField("age", "number")

		if len(typeDef.Fields) != 2 {
			t.Errorf("Expected 2 fields after adding, got %d", len(typeDef.Fields))
		}

		// Vérifier le premier champ
		if typeDef.Fields[0].Name != "name" {
			t.Errorf("Expected first field name 'name', got '%s'", typeDef.Fields[0].Name)
		}

		if typeDef.Fields[0].Type != "string" {
			t.Errorf("Expected first field type 'string', got '%s'", typeDef.Fields[0].Type)
		}

		// Vérifier le deuxième champ
		if typeDef.Fields[1].Name != "age" {
			t.Errorf("Expected second field name 'age', got '%s'", typeDef.Fields[1].Name)
		}

		if typeDef.Fields[1].Type != "number" {
			t.Errorf("Expected second field type 'number', got '%s'", typeDef.Fields[1].Type)
		}
	})

	t.Run("GetFieldByName", func(t *testing.T) {
		typeDef := NewTypeDefinition("Person")
		typeDef.AddField("name", "string")
		typeDef.AddField("age", "number")

		// Champ qui existe
		field := typeDef.GetFieldByName("name")
		if field == nil {
			t.Fatal("GetFieldByName should return the field")
		}

		if field.Name != "name" || field.Type != "string" {
			t.Errorf("Expected field {name: string}, got {%s: %s}", field.Name, field.Type)
		}

		// Champ qui n'existe pas
		field = typeDef.GetFieldByName("nonexistent")
		if field != nil {
			t.Error("GetFieldByName should return nil for non-existent field")
		}
	})

	t.Run("HasField", func(t *testing.T) {
		typeDef := NewTypeDefinition("Person")
		typeDef.AddField("name", "string")

		if !typeDef.HasField("name") {
			t.Error("HasField should return true for existing field")
		}

		if typeDef.HasField("nonexistent") {
			t.Error("HasField should return false for non-existent field")
		}
	})
}

func TestExpression(t *testing.T) {
	t.Run("NewExpression", func(t *testing.T) {
		expr := NewExpression()

		if expr.Type != "expression" {
			t.Errorf("Expected type 'expression', got '%s'", expr.Type)
		}

		if expr.Set.Type != "set" {
			t.Errorf("Expected set type 'set', got '%s'", expr.Set.Type)
		}

		if len(expr.Set.Variables) != 0 {
			t.Errorf("Expected empty variables slice, got %d", len(expr.Set.Variables))
		}
	})

	t.Run("AddVariable", func(t *testing.T) {
		expr := NewExpression()

		expr.AddVariable("person", "Person")
		expr.AddVariable("company", "Company")

		if len(expr.Set.Variables) != 2 {
			t.Errorf("Expected 2 variables, got %d", len(expr.Set.Variables))
		}

		// Vérifier la première variable
		if expr.Set.Variables[0].Name != "person" {
			t.Errorf("Expected first variable name 'person', got '%s'", expr.Set.Variables[0].Name)
		}

		if expr.Set.Variables[0].DataType != "Person" {
			t.Errorf("Expected first variable dataType 'Person', got '%s'", expr.Set.Variables[0].DataType)
		}
	})
}

func TestConstructors(t *testing.T) {
	t.Run("NewConstraint", func(t *testing.T) {
		left := "variable"
		right := "value"
		constraint := NewConstraint(left, "==", right)

		if constraint.Type != "constraint" {
			t.Errorf("Expected type 'constraint', got '%s'", constraint.Type)
		}

		if constraint.Operator != "==" {
			t.Errorf("Expected operator '==', got '%s'", constraint.Operator)
		}

		if constraint.Left != left {
			t.Errorf("Expected left operand '%v', got '%v'", left, constraint.Left)
		}

		if constraint.Right != right {
			t.Errorf("Expected right operand '%v', got '%v'", right, constraint.Right)
		}
	})

	t.Run("NewFieldAccess", func(t *testing.T) {
		fa := NewFieldAccess("person", "name")

		if fa.Type != "fieldAccess" {
			t.Errorf("Expected type 'fieldAccess', got '%s'", fa.Type)
		}

		if fa.Object != "person" {
			t.Errorf("Expected object 'person', got '%s'", fa.Object)
		}

		if fa.Field != "name" {
			t.Errorf("Expected field 'name', got '%s'", fa.Field)
		}
	})

	t.Run("NewAction", func(t *testing.T) {
		action := NewAction("notify", "email", "sms")

		if action.Type != "action" {
			t.Errorf("Expected type 'action', got '%s'", action.Type)
		}

		if action.Job.Type != "jobCall" {
			t.Errorf("Expected job type 'jobCall', got '%s'", action.Job.Type)
		}

		if action.Job.Name != "notify" {
			t.Errorf("Expected job name 'notify', got '%s'", action.Job.Name)
		}

		if len(action.Job.Args) != 2 {
			t.Errorf("Expected 2 arguments, got %d", len(action.Job.Args))
		}

		if action.Job.Args[0].(string) != "email" {
			t.Errorf("Expected first arg 'email', got '%s'", action.Job.Args[0])
		}
	})
}

func TestValidationHelpers(t *testing.T) {
	t.Run("IsValidOperator", func(t *testing.T) {
		validOps := []string{"==", "!=", "<", ">", "<=", ">=", "AND", "OR", "NOT", "+", "-", "*", "/", "%"}

		for _, op := range validOps {
			if !IsValidOperator(op) {
				t.Errorf("IsValidOperator should return true for '%s'", op)
			}
		}

		invalidOps := []string{"===", "<<", ">>", "XOR", "MOD", "?", "!=="}

		for _, op := range invalidOps {
			if IsValidOperator(op) {
				t.Errorf("IsValidOperator should return false for '%s'", op)
			}
		}
	})

	t.Run("IsValidType", func(t *testing.T) {
		validTypes := []string{"string", "number", "bool", "integer"}

		for _, typ := range validTypes {
			if !IsValidType(typ) {
				t.Errorf("IsValidType should return true for '%s'", typ)
			}
		}

		invalidTypes := []string{"float", "double", "char", "void", "object", "array"}

		for _, typ := range invalidTypes {
			if IsValidType(typ) {
				t.Errorf("IsValidType should return false for '%s'", typ)
			}
		}
	})
}

func TestMetadata(t *testing.T) {
	t.Run("MetadataCreation", func(t *testing.T) {
		prog := NewProgram()

		if prog.Metadata.CreatedAt.IsZero() {
			t.Error("CreatedAt should not be zero time")
		}

		// Vérifier que le timestamp est récent (dans la dernière seconde)
		now := time.Now()
		diff := now.Sub(prog.Metadata.CreatedAt)
		if diff > time.Second {
			t.Errorf("CreatedAt should be recent, got diff of %v", diff)
		}

		// Test de modification des métadonnées
		prog.Metadata.Author = "Test Author"
		prog.Metadata.Description = "Test Program"

		if prog.Metadata.Author != "Test Author" {
			t.Errorf("Expected author 'Test Author', got '%s'", prog.Metadata.Author)
		}

		if prog.Metadata.Description != "Test Program" {
			t.Errorf("Expected description 'Test Program', got '%s'", prog.Metadata.Description)
		}
	})
}
