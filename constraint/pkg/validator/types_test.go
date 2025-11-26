// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package validator

import (
	"sync"
	"testing"

	"github.com/treivax/tsd/constraint/pkg/domain"
)

// TestNewTypeRegistry teste la création d'un nouveau registre de types
func TestNewTypeRegistry(t *testing.T) {
	registry := NewTypeRegistry()
	if registry == nil {
		t.Fatal("NewTypeRegistry should not return nil")
	}
	if registry.types == nil {
		t.Error("TypeRegistry.types map should be initialized")
	}
	if len(registry.types) != 0 {
		t.Errorf("Expected empty registry, got %d types", len(registry.types))
	}
}

// TestTypeRegistry_RegisterType teste l'enregistrement de types
func TestTypeRegistry_RegisterType(t *testing.T) {
	tests := []struct {
		name        string
		typeDef     domain.TypeDefinition
		wantErr     bool
		errContains string
	}{
		{
			name: "register valid type",
			typeDef: domain.TypeDefinition{
				Type: "typeDefinition",
				Name: "Person",
				Fields: []domain.Field{
					{Name: "name", Type: "string"},
					{Name: "age", Type: "integer"},
				},
			},
			wantErr: false,
		},
		{
			name: "register type with single field",
			typeDef: domain.TypeDefinition{
				Type: "typeDefinition",
				Name: "Product",
				Fields: []domain.Field{
					{Name: "price", Type: "number"},
				},
			},
			wantErr: false,
		},
		{
			name: "register type with all field types",
			typeDef: domain.TypeDefinition{
				Type: "typeDefinition",
				Name: "AllTypes",
				Fields: []domain.Field{
					{Name: "str", Type: "string"},
					{Name: "num", Type: "number"},
					{Name: "int", Type: "integer"},
					{Name: "bool", Type: "bool"},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registry := NewTypeRegistry()
			err := registry.RegisterType(tt.typeDef)
			if (err != nil) != tt.wantErr {
				t.Errorf("RegisterType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errContains != "" {
				if !containsString(err.Error(), tt.errContains) {
					t.Errorf("Expected error to contain %q, got %q", tt.errContains, err.Error())
				}
			}
			if !tt.wantErr {
				// Vérifier que le type a été enregistré
				if !registry.HasType(tt.typeDef.Name) {
					t.Errorf("Type %s should be registered", tt.typeDef.Name)
				}
			}
		})
	}
}

// TestTypeRegistry_RegisterType_Duplicate teste l'enregistrement de types dupliqués
func TestTypeRegistry_RegisterType_Duplicate(t *testing.T) {
	registry := NewTypeRegistry()
	typeDef := domain.TypeDefinition{
		Type: "typeDefinition",
		Name: "Person",
		Fields: []domain.Field{
			{Name: "name", Type: "string"},
		},
	}

	// Premier enregistrement - devrait réussir
	err := registry.RegisterType(typeDef)
	if err != nil {
		t.Fatalf("First registration should succeed: %v", err)
	}

	// Deuxième enregistrement - devrait échouer
	err = registry.RegisterType(typeDef)
	if err == nil {
		t.Error("Registering duplicate type should fail")
	}
	if !domain.IsValidationError(err) {
		t.Errorf("Expected ValidationError, got %T", err)
	}
	if !containsString(err.Error(), "already exists") {
		t.Errorf("Expected error about duplicate, got: %v", err)
	}
}

// TestTypeRegistry_GetType teste la récupération de types
func TestTypeRegistry_GetType(t *testing.T) {
	registry := NewTypeRegistry()
	typeDef := domain.TypeDefinition{
		Type: "typeDefinition",
		Name: "Person",
		Fields: []domain.Field{
			{Name: "name", Type: "string"},
			{Name: "age", Type: "integer"},
		},
	}

	// Enregistrer le type
	err := registry.RegisterType(typeDef)
	if err != nil {
		t.Fatalf("RegisterType failed: %v", err)
	}

	// Récupérer le type
	retrieved, err := registry.GetType("Person")
	if err != nil {
		t.Fatalf("GetType failed: %v", err)
	}
	if retrieved == nil {
		t.Fatal("GetType returned nil")
	}
	if retrieved.Name != typeDef.Name {
		t.Errorf("Expected name %q, got %q", typeDef.Name, retrieved.Name)
	}
	if len(retrieved.Fields) != len(typeDef.Fields) {
		t.Errorf("Expected %d fields, got %d", len(typeDef.Fields), len(retrieved.Fields))
	}
}

// TestTypeRegistry_GetType_NotFound teste la récupération d'un type inexistant
func TestTypeRegistry_GetType_NotFound(t *testing.T) {
	registry := NewTypeRegistry()

	retrieved, err := registry.GetType("NonExistent")
	if err == nil {
		t.Error("GetType should return error for non-existent type")
	}
	if retrieved != nil {
		t.Error("GetType should return nil for non-existent type")
	}
	if !domain.IsUnknownTypeError(err) {
		t.Errorf("Expected UnknownTypeError, got %T", err)
	}
}

// TestTypeRegistry_HasType teste la vérification d'existence de types
func TestTypeRegistry_HasType(t *testing.T) {
	registry := NewTypeRegistry()
	typeDef := domain.TypeDefinition{
		Type: "typeDefinition",
		Name: "Person",
		Fields: []domain.Field{
			{Name: "name", Type: "string"},
		},
	}

	// Type non enregistré
	if registry.HasType("Person") {
		t.Error("HasType should return false for non-existent type")
	}

	// Enregistrer le type
	err := registry.RegisterType(typeDef)
	if err != nil {
		t.Fatalf("RegisterType failed: %v", err)
	}

	// Type enregistré
	if !registry.HasType("Person") {
		t.Error("HasType should return true for existing type")
	}

	// Autre type non enregistré
	if registry.HasType("Product") {
		t.Error("HasType should return false for non-existent type")
	}
}

// TestTypeRegistry_ListTypes teste la liste de tous les types
func TestTypeRegistry_ListTypes(t *testing.T) {
	registry := NewTypeRegistry()

	// Liste vide au début
	types := registry.ListTypes()
	if len(types) != 0 {
		t.Errorf("Expected empty list, got %d types", len(types))
	}

	// Enregistrer plusieurs types
	typeDefPerson := domain.TypeDefinition{
		Type: "typeDefinition",
		Name: "Person",
		Fields: []domain.Field{
			{Name: "name", Type: "string"},
		},
	}
	typeDefProduct := domain.TypeDefinition{
		Type: "typeDefinition",
		Name: "Product",
		Fields: []domain.Field{
			{Name: "price", Type: "number"},
		},
	}

	registry.RegisterType(typeDefPerson)
	registry.RegisterType(typeDefProduct)

	// Vérifier la liste
	types = registry.ListTypes()
	if len(types) != 2 {
		t.Errorf("Expected 2 types, got %d", len(types))
	}

	// Vérifier que les types sont présents
	foundPerson := false
	foundProduct := false
	for _, td := range types {
		if td.Name == "Person" {
			foundPerson = true
		}
		if td.Name == "Product" {
			foundProduct = true
		}
	}
	if !foundPerson {
		t.Error("Person type not found in list")
	}
	if !foundProduct {
		t.Error("Product type not found in list")
	}
}

// TestTypeRegistry_GetTypeFields teste la récupération des champs d'un type
func TestTypeRegistry_GetTypeFields(t *testing.T) {
	registry := NewTypeRegistry()
	typeDef := domain.TypeDefinition{
		Type: "typeDefinition",
		Name: "Person",
		Fields: []domain.Field{
			{Name: "name", Type: "string"},
			{Name: "age", Type: "integer"},
			{Name: "active", Type: "bool"},
		},
	}

	// Enregistrer le type
	err := registry.RegisterType(typeDef)
	if err != nil {
		t.Fatalf("RegisterType failed: %v", err)
	}

	// Récupérer les champs
	fields, err := registry.GetTypeFields("Person")
	if err != nil {
		t.Fatalf("GetTypeFields failed: %v", err)
	}
	if len(fields) != 3 {
		t.Errorf("Expected 3 fields, got %d", len(fields))
	}

	// Vérifier les champs
	expectedFields := map[string]string{
		"name":   "string",
		"age":    "integer",
		"active": "bool",
	}
	for fieldName, expectedType := range expectedFields {
		actualType, exists := fields[fieldName]
		if !exists {
			t.Errorf("Field %q not found", fieldName)
			continue
		}
		if actualType != expectedType {
			t.Errorf("Field %q: expected type %q, got %q", fieldName, expectedType, actualType)
		}
	}
}

// TestTypeRegistry_GetTypeFields_NotFound teste la récupération des champs d'un type inexistant
func TestTypeRegistry_GetTypeFields_NotFound(t *testing.T) {
	registry := NewTypeRegistry()

	fields, err := registry.GetTypeFields("NonExistent")
	if err == nil {
		t.Error("GetTypeFields should return error for non-existent type")
	}
	if fields != nil {
		t.Error("GetTypeFields should return nil for non-existent type")
	}
	if !domain.IsUnknownTypeError(err) {
		t.Errorf("Expected UnknownTypeError, got %T", err)
	}
}

// TestTypeRegistry_Clear teste le nettoyage du registre
func TestTypeRegistry_Clear(t *testing.T) {
	registry := NewTypeRegistry()

	// Enregistrer des types
	registry.RegisterType(domain.TypeDefinition{
		Type:   "typeDefinition",
		Name:   "Person",
		Fields: []domain.Field{{Name: "name", Type: "string"}},
	})
	registry.RegisterType(domain.TypeDefinition{
		Type:   "typeDefinition",
		Name:   "Product",
		Fields: []domain.Field{{Name: "price", Type: "number"}},
	})

	// Vérifier qu'ils sont présents
	if len(registry.ListTypes()) != 2 {
		t.Error("Expected 2 types before Clear")
	}

	// Nettoyer
	registry.Clear()

	// Vérifier qu'ils sont supprimés
	if len(registry.ListTypes()) != 0 {
		t.Error("Expected 0 types after Clear")
	}
	if registry.HasType("Person") {
		t.Error("Person should not exist after Clear")
	}
	if registry.HasType("Product") {
		t.Error("Product should not exist after Clear")
	}
}

// TestTypeRegistry_Concurrent teste l'accès concurrent au registre
func TestTypeRegistry_Concurrent(t *testing.T) {
	registry := NewTypeRegistry()
	var wg sync.WaitGroup

	// Enregistrer des types en concurrence
	numGoroutines := 10
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			typeDef := domain.TypeDefinition{
				Type:   "typeDefinition",
				Name:   "Type" + string(rune('A'+id)),
				Fields: []domain.Field{{Name: "field", Type: "string"}},
			}
			registry.RegisterType(typeDef)
		}(i)
	}
	wg.Wait()

	// Vérifier les lectures concurrentes
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			typeName := "Type" + string(rune('A'+id))
			registry.HasType(typeName)
			registry.GetType(typeName)
		}(i)
	}
	wg.Wait()

	// Vérifier ListTypes en concurrence
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			registry.ListTypes()
		}()
	}
	wg.Wait()
}

// TestTypeRegistry_EmptyFieldList teste l'enregistrement d'un type sans champs
func TestTypeRegistry_EmptyFieldList(t *testing.T) {
	registry := NewTypeRegistry()
	typeDef := domain.TypeDefinition{
		Type:   "typeDefinition",
		Name:   "EmptyType",
		Fields: []domain.Field{},
	}

	// L'enregistrement devrait réussir (la validation est faite par le validateur)
	err := registry.RegisterType(typeDef)
	if err != nil {
		t.Errorf("RegisterType with empty fields should succeed: %v", err)
	}

	// Vérifier GetTypeFields
	fields, err := registry.GetTypeFields("EmptyType")
	if err != nil {
		t.Fatalf("GetTypeFields failed: %v", err)
	}
	if len(fields) != 0 {
		t.Errorf("Expected 0 fields, got %d", len(fields))
	}
}

// TestTypeRegistry_ComplexType teste l'enregistrement de types complexes
func TestTypeRegistry_ComplexType(t *testing.T) {
	registry := NewTypeRegistry()
	typeDef := domain.TypeDefinition{
		Type: "typeDefinition",
		Name: "ComplexType",
		Fields: []domain.Field{
			{Name: "id", Type: "integer"},
			{Name: "name", Type: "string"},
			{Name: "price", Type: "number"},
			{Name: "available", Type: "bool"},
			{Name: "count", Type: "integer"},
			{Name: "description", Type: "string"},
			{Name: "rating", Type: "number"},
			{Name: "active", Type: "bool"},
		},
	}

	err := registry.RegisterType(typeDef)
	if err != nil {
		t.Fatalf("RegisterType failed: %v", err)
	}

	// Vérifier tous les champs
	fields, err := registry.GetTypeFields("ComplexType")
	if err != nil {
		t.Fatalf("GetTypeFields failed: %v", err)
	}
	if len(fields) != 8 {
		t.Errorf("Expected 8 fields, got %d", len(fields))
	}
}

// TestTypeRegistry_MultipleRegistrationAttempts teste plusieurs tentatives d'enregistrement
func TestTypeRegistry_MultipleRegistrationAttempts(t *testing.T) {
	registry := NewTypeRegistry()
	typeDef := domain.TypeDefinition{
		Type:   "typeDefinition",
		Name:   "Test",
		Fields: []domain.Field{{Name: "field", Type: "string"}},
	}

	// Premier enregistrement
	err := registry.RegisterType(typeDef)
	if err != nil {
		t.Fatalf("First registration failed: %v", err)
	}

	// Tentatives suivantes devraient toutes échouer
	for i := 0; i < 5; i++ {
		err = registry.RegisterType(typeDef)
		if err == nil {
			t.Errorf("Registration attempt %d should have failed", i+2)
		}
	}

	// Vérifier qu'il n'y a toujours qu'un seul type
	types := registry.ListTypes()
	if len(types) != 1 {
		t.Errorf("Expected 1 type, got %d", len(types))
	}
}

// Helper function
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && contains(s, substr)))
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// ===== TypeChecker Tests =====

// TestNewTypeChecker teste la création d'un nouveau vérificateur de types
func TestNewTypeChecker(t *testing.T) {
	registry := NewTypeRegistry()
	checker := NewTypeChecker(registry)
	if checker == nil {
		t.Fatal("NewTypeChecker should not return nil")
	}
	if checker.registry == nil {
		t.Error("TypeChecker.registry should be set")
	}
}

// TestTypeChecker_GetFieldType teste la récupération du type d'un champ
func TestTypeChecker_GetFieldType(t *testing.T) {
	registry := NewTypeRegistry()
	checker := NewTypeChecker(registry)

	// Enregistrer un type
	typeDef := domain.TypeDefinition{
		Type: "typeDefinition",
		Name: "Person",
		Fields: []domain.Field{
			{Name: "name", Type: "string"},
			{Name: "age", Type: "integer"},
			{Name: "salary", Type: "number"},
			{Name: "active", Type: "bool"},
		},
	}
	registry.RegisterType(typeDef)

	// Variables
	variables := []domain.TypedVariable{
		{Type: "typedVariable", Name: "p", DataType: "Person"},
	}

	tests := []struct {
		name        string
		fieldAccess interface{}
		wantType    string
		wantErr     bool
	}{
		{
			name: "get string field type",
			fieldAccess: &domain.FieldAccess{
				Type:   "fieldAccess",
				Object: "p",
				Field:  "name",
			},
			wantType: "string",
			wantErr:  false,
		},
		{
			name: "get integer field type",
			fieldAccess: &domain.FieldAccess{
				Type:   "fieldAccess",
				Object: "p",
				Field:  "age",
			},
			wantType: "integer",
			wantErr:  false,
		},
		{
			name: "get number field type",
			fieldAccess: &domain.FieldAccess{
				Type:   "fieldAccess",
				Object: "p",
				Field:  "salary",
			},
			wantType: "number",
			wantErr:  false,
		},
		{
			name: "get bool field type",
			fieldAccess: &domain.FieldAccess{
				Type:   "fieldAccess",
				Object: "p",
				Field:  "active",
			},
			wantType: "bool",
			wantErr:  false,
		},
		{
			name: "field access as map",
			fieldAccess: map[string]interface{}{
				"type":   "fieldAccess",
				"object": "p",
				"field":  "name",
			},
			wantType: "string",
			wantErr:  false,
		},
		{
			name: "unknown field",
			fieldAccess: &domain.FieldAccess{
				Type:   "fieldAccess",
				Object: "p",
				Field:  "unknown",
			},
			wantErr: true,
		},
		{
			name: "unknown variable",
			fieldAccess: &domain.FieldAccess{
				Type:   "fieldAccess",
				Object: "unknown",
				Field:  "name",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fieldType, err := checker.GetFieldType(tt.fieldAccess, variables, []domain.TypeDefinition{typeDef})
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFieldType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && fieldType != tt.wantType {
				t.Errorf("GetFieldType() = %v, want %v", fieldType, tt.wantType)
			}
		})
	}
}

// TestTypeChecker_GetValueType teste la détection du type d'une valeur
func TestTypeChecker_GetValueType(t *testing.T) {
	registry := NewTypeRegistry()
	checker := NewTypeChecker(registry)

	tests := []struct {
		name     string
		value    interface{}
		wantType string
	}{
		{name: "bool true", value: true, wantType: "bool"},
		{name: "bool false", value: false, wantType: "bool"},
		{name: "int", value: 42, wantType: "integer"},
		{name: "int8", value: int8(42), wantType: "integer"},
		{name: "int16", value: int16(42), wantType: "integer"},
		{name: "int32", value: int32(42), wantType: "integer"},
		{name: "int64", value: int64(42), wantType: "integer"},
		{name: "float32", value: float32(3.14), wantType: "number"},
		{name: "float64", value: float64(3.14), wantType: "number"},
		{name: "string", value: "hello", wantType: "string"},
		{name: "empty string", value: "", wantType: "string"},
		{
			name: "map with type booleanLiteral",
			value: map[string]interface{}{
				"type":  "booleanLiteral",
				"value": true,
			},
			wantType: "bool",
		},
		{
			name: "map with type integerLiteral",
			value: map[string]interface{}{
				"type":  "integerLiteral",
				"value": 42,
			},
			wantType: "integer",
		},
		{
			name: "map with type numberLiteral",
			value: map[string]interface{}{
				"type":  "numberLiteral",
				"value": 3.14,
			},
			wantType: "number",
		},
		{
			name: "map with type stringLiteral",
			value: map[string]interface{}{
				"type":  "stringLiteral",
				"value": "hello",
			},
			wantType: "string",
		},
		{
			name: "map with nested value",
			value: map[string]interface{}{
				"value": "test",
			},
			wantType: "string",
		},
		{
			name:     "unknown type",
			value:    struct{}{},
			wantType: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotType := checker.GetValueType(tt.value)
			if gotType != tt.wantType {
				t.Errorf("GetValueType() = %v, want %v", gotType, tt.wantType)
			}
		})
	}
}

// TestTypeChecker_ValidateTypeCompatibility_Comparison teste les opérateurs de comparaison
func TestTypeChecker_ValidateTypeCompatibility_Comparison(t *testing.T) {
	registry := NewTypeRegistry()
	checker := NewTypeChecker(registry)

	tests := []struct {
		name      string
		leftType  string
		rightType string
		operator  string
		wantErr   bool
	}{
		// Égalité/Inégalité - tous types si identiques
		{name: "string == string", leftType: "string", rightType: "string", operator: "==", wantErr: false},
		{name: "integer == integer", leftType: "integer", rightType: "integer", operator: "==", wantErr: false},
		{name: "number == number", leftType: "number", rightType: "number", operator: "==", wantErr: false},
		{name: "bool == bool", leftType: "bool", rightType: "bool", operator: "==", wantErr: false},
		{name: "string != string", leftType: "string", rightType: "string", operator: "!=", wantErr: false},
		{name: "string == integer (mismatch)", leftType: "string", rightType: "integer", operator: "==", wantErr: true},
		{name: "bool != integer (mismatch)", leftType: "bool", rightType: "integer", operator: "!=", wantErr: true},

		// Comparaisons ordinales - seulement number, integer, string
		{name: "integer < integer", leftType: "integer", rightType: "integer", operator: "<", wantErr: false},
		{name: "number < number", leftType: "number", rightType: "number", operator: "<", wantErr: false},
		{name: "string < string", leftType: "string", rightType: "string", operator: "<", wantErr: false},
		{name: "integer > integer", leftType: "integer", rightType: "integer", operator: ">", wantErr: false},
		{name: "number <= number", leftType: "number", rightType: "number", operator: "<=", wantErr: false},
		{name: "string >= string", leftType: "string", rightType: "string", operator: ">=", wantErr: false},

		// Comparaisons ordinales - types incompatibles
		{name: "bool < bool (not orderable)", leftType: "bool", rightType: "bool", operator: "<", wantErr: true},
		{name: "integer < number (mismatch)", leftType: "integer", rightType: "number", operator: "<", wantErr: true},
		{name: "string > integer (mismatch)", leftType: "string", rightType: "integer", operator: ">", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checker.ValidateTypeCompatibility(tt.leftType, tt.rightType, tt.operator)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTypeCompatibility() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestTypeChecker_ValidateTypeCompatibility_Logical teste les opérateurs logiques
func TestTypeChecker_ValidateTypeCompatibility_Logical(t *testing.T) {
	registry := NewTypeRegistry()
	checker := NewTypeChecker(registry)

	tests := []struct {
		name      string
		leftType  string
		rightType string
		operator  string
		wantErr   bool
	}{
		{name: "bool AND bool", leftType: "bool", rightType: "bool", operator: "AND", wantErr: false},
		{name: "bool OR bool", leftType: "bool", rightType: "bool", operator: "OR", wantErr: false},
		{name: "NOT bool", leftType: "bool", rightType: "", operator: "NOT", wantErr: false},

		// Types incompatibles
		{name: "integer AND integer", leftType: "integer", rightType: "integer", operator: "AND", wantErr: true},
		{name: "string OR string", leftType: "string", rightType: "string", operator: "OR", wantErr: true},
		{name: "bool AND integer", leftType: "bool", rightType: "integer", operator: "AND", wantErr: true},
		{name: "NOT integer", leftType: "integer", rightType: "", operator: "NOT", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checker.ValidateTypeCompatibility(tt.leftType, tt.rightType, tt.operator)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTypeCompatibility() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestTypeChecker_ValidateTypeCompatibility_Arithmetic teste les opérateurs arithmétiques
func TestTypeChecker_ValidateTypeCompatibility_Arithmetic(t *testing.T) {
	registry := NewTypeRegistry()
	checker := NewTypeChecker(registry)

	tests := []struct {
		name      string
		leftType  string
		rightType string
		operator  string
		wantErr   bool
	}{
		{name: "integer + integer", leftType: "integer", rightType: "integer", operator: "+", wantErr: false},
		{name: "number + number", leftType: "number", rightType: "number", operator: "+", wantErr: false},
		{name: "integer - integer", leftType: "integer", rightType: "integer", operator: "-", wantErr: false},
		{name: "number * number", leftType: "number", rightType: "number", operator: "*", wantErr: false},
		{name: "integer / integer", leftType: "integer", rightType: "integer", operator: "/", wantErr: false},
		{name: "integer % integer", leftType: "integer", rightType: "integer", operator: "%", wantErr: false},

		// Types non numériques
		{name: "string + string", leftType: "string", rightType: "string", operator: "+", wantErr: true},
		{name: "bool + bool", leftType: "bool", rightType: "bool", operator: "+", wantErr: true},
		{name: "string - integer", leftType: "string", rightType: "integer", operator: "-", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checker.ValidateTypeCompatibility(tt.leftType, tt.rightType, tt.operator)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTypeCompatibility() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestTypeChecker_ValidateTypeCompatibility_InvalidOperator teste les opérateurs invalides
func TestTypeChecker_ValidateTypeCompatibility_InvalidOperator(t *testing.T) {
	registry := NewTypeRegistry()
	checker := NewTypeChecker(registry)

	invalidOperators := []string{"INVALID", "&&", "||", "===", "!==", "++", "--"}

	for _, op := range invalidOperators {
		t.Run("invalid operator "+op, func(t *testing.T) {
			err := checker.ValidateTypeCompatibility("integer", "integer", op)
			if err == nil {
				t.Errorf("ValidateTypeCompatibility() should fail for invalid operator %q", op)
			}
			if !domain.IsValidationError(err) {
				t.Errorf("Expected ValidationError for invalid operator, got %T", err)
			}
		})
	}
}

// TestTypeChecker_GetFieldType_UnknownType teste l'accès à un champ d'un type inconnu
func TestTypeChecker_GetFieldType_UnknownType(t *testing.T) {
	registry := NewTypeRegistry()
	checker := NewTypeChecker(registry)

	variables := []domain.TypedVariable{
		{Type: "typedVariable", Name: "x", DataType: "UnknownType"},
	}

	fieldAccess := &domain.FieldAccess{
		Type:   "fieldAccess",
		Object: "x",
		Field:  "field",
	}

	_, err := checker.GetFieldType(fieldAccess, variables, []domain.TypeDefinition{})
	if err == nil {
		t.Error("GetFieldType() should fail for unknown type")
	}
	if !domain.IsUnknownTypeError(err) {
		t.Errorf("Expected UnknownTypeError, got %T", err)
	}
}

// TestTypeChecker_GetFieldType_InvalidFormat teste un format invalide de field access
func TestTypeChecker_GetFieldType_InvalidFormat(t *testing.T) {
	registry := NewTypeRegistry()
	checker := NewTypeChecker(registry)

	variables := []domain.TypedVariable{
		{Type: "typedVariable", Name: "p", DataType: "Person"},
	}

	// Format invalide (ni *FieldAccess ni map)
	invalidAccess := "invalid"

	_, err := checker.GetFieldType(invalidAccess, variables, []domain.TypeDefinition{})
	if err == nil {
		t.Error("GetFieldType() should fail for invalid format")
	}
	if !domain.IsValidationError(err) {
		t.Errorf("Expected ValidationError, got %T", err)
	}
}
