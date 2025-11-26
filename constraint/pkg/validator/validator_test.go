// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package validator

import (
	"testing"

	"github.com/treivax/tsd/constraint/pkg/domain"
)

// ===== ConstraintValidator Tests =====

// TestNewConstraintValidator teste la création d'un nouveau validateur
func TestNewConstraintValidator(t *testing.T) {
	registry := NewTypeRegistry()
	checker := NewTypeChecker(registry)
	validator := NewConstraintValidator(registry, checker)

	if validator == nil {
		t.Fatal("NewConstraintValidator should not return nil")
	}
	if validator.typeRegistry == nil {
		t.Error("ConstraintValidator.typeRegistry should be set")
	}
	if validator.typeChecker == nil {
		t.Error("ConstraintValidator.typeChecker should be set")
	}
	if validator.config.StrictMode != true {
		t.Error("Default config should have StrictMode enabled")
	}
	if len(validator.config.AllowedOperators) == 0 {
		t.Error("Default config should have allowed operators")
	}
}

// TestConstraintValidator_ValidateTypes teste la validation des types
func TestConstraintValidator_ValidateTypes(t *testing.T) {
	registry := NewTypeRegistry()
	checker := NewTypeChecker(registry)
	validator := NewConstraintValidator(registry, checker)

	tests := []struct {
		name        string
		types       []domain.TypeDefinition
		wantErr     bool
		errContains string
	}{
		{
			name: "valid single type",
			types: []domain.TypeDefinition{
				{
					Type: "typeDefinition",
					Name: "Person",
					Fields: []domain.Field{
						{Name: "name", Type: "string"},
						{Name: "age", Type: "integer"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid multiple types",
			types: []domain.TypeDefinition{
				{
					Type: "typeDefinition",
					Name: "Person",
					Fields: []domain.Field{
						{Name: "name", Type: "string"},
					},
				},
				{
					Type: "typeDefinition",
					Name: "Product",
					Fields: []domain.Field{
						{Name: "price", Type: "number"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "duplicate type name",
			types: []domain.TypeDefinition{
				{
					Type: "typeDefinition",
					Name: "Person",
					Fields: []domain.Field{
						{Name: "name", Type: "string"},
					},
				},
				{
					Type: "typeDefinition",
					Name: "Person",
					Fields: []domain.Field{
						{Name: "age", Type: "integer"},
					},
				},
			},
			wantErr:     true,
			errContains: "duplicate type name",
		},
		{
			name: "empty type name",
			types: []domain.TypeDefinition{
				{
					Type: "typeDefinition",
					Name: "",
					Fields: []domain.Field{
						{Name: "field", Type: "string"},
					},
				},
			},
			wantErr:     true,
			errContains: "name cannot be empty",
		},
		{
			name: "type with no fields",
			types: []domain.TypeDefinition{
				{
					Type:   "typeDefinition",
					Name:   "Empty",
					Fields: []domain.Field{},
				},
			},
			wantErr:     true,
			errContains: "must have at least one field",
		},
		{
			name: "duplicate field name",
			types: []domain.TypeDefinition{
				{
					Type: "typeDefinition",
					Name: "Person",
					Fields: []domain.Field{
						{Name: "name", Type: "string"},
						{Name: "name", Type: "string"},
					},
				},
			},
			wantErr:     true,
			errContains: "duplicate field name",
		},
		{
			name: "empty field name",
			types: []domain.TypeDefinition{
				{
					Type: "typeDefinition",
					Name: "Person",
					Fields: []domain.Field{
						{Name: "", Type: "string"},
					},
				},
			},
			wantErr:     true,
			errContains: "field name cannot be empty",
		},
		{
			name: "invalid field type",
			types: []domain.TypeDefinition{
				{
					Type: "typeDefinition",
					Name: "Person",
					Fields: []domain.Field{
						{Name: "field", Type: "invalid_type"},
					},
				},
			},
			wantErr:     true,
			errContains: "invalid field type",
		},
		{
			name: "all valid field types",
			types: []domain.TypeDefinition{
				{
					Type: "typeDefinition",
					Name: "AllTypes",
					Fields: []domain.Field{
						{Name: "str", Type: "string"},
						{Name: "num", Type: "number"},
						{Name: "int", Type: "integer"},
						{Name: "bool", Type: "bool"},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateTypes(tt.types)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTypes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errContains != "" {
				if !containsString(err.Error(), tt.errContains) {
					t.Errorf("Expected error to contain %q, got %q", tt.errContains, err.Error())
				}
			}
		})
	}
}

// TestConstraintValidator_ValidateExpression teste la validation des expressions
func TestConstraintValidator_ValidateExpression(t *testing.T) {
	registry := NewTypeRegistry()
	checker := NewTypeChecker(registry)
	validator := NewConstraintValidator(registry, checker)

	// Préparer les types
	types := []domain.TypeDefinition{
		{
			Type: "typeDefinition",
			Name: "Person",
			Fields: []domain.Field{
				{Name: "name", Type: "string"},
				{Name: "age", Type: "integer"},
			},
		},
		{
			Type: "typeDefinition",
			Name: "Product",
			Fields: []domain.Field{
				{Name: "price", Type: "number"},
			},
		},
	}

	tests := []struct {
		name        string
		expr        domain.Expression
		wantErr     bool
		errContains string
	}{
		{
			name: "valid expression with action",
			expr: domain.Expression{
				Type: "expression",
				Set: domain.Set{
					Type: "set",
					Variables: []domain.TypedVariable{
						{Type: "typedVariable", Name: "p", DataType: "Person"},
					},
				},
				Action: &domain.Action{
					Type: "action",
					Job: domain.JobCall{
						Type: "jobCall",
						Name: "notify",
						Args: []interface{}{"message"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "expression with no variables",
			expr: domain.Expression{
				Type: "expression",
				Set: domain.Set{
					Type:      "set",
					Variables: []domain.TypedVariable{},
				},
				Action: &domain.Action{
					Type: "action",
					Job:  domain.JobCall{Type: "jobCall", Name: "test"},
				},
			},
			wantErr:     true,
			errContains: "must have at least one variable",
		},
		{
			name: "expression with unknown type",
			expr: domain.Expression{
				Type: "expression",
				Set: domain.Set{
					Type: "set",
					Variables: []domain.TypedVariable{
						{Type: "typedVariable", Name: "x", DataType: "UnknownType"},
					},
				},
				Action: &domain.Action{
					Type: "action",
					Job:  domain.JobCall{Type: "jobCall", Name: "test"},
				},
			},
			wantErr:     true,
			errContains: "UnknownType",
		},
		{
			name: "expression with multiple variables",
			expr: domain.Expression{
				Type: "expression",
				Set: domain.Set{
					Type: "set",
					Variables: []domain.TypedVariable{
						{Type: "typedVariable", Name: "p", DataType: "Person"},
						{Type: "typedVariable", Name: "pr", DataType: "Product"},
					},
				},
				Action: &domain.Action{
					Type: "action",
					Job:  domain.JobCall{Type: "jobCall", Name: "test"},
				},
			},
			wantErr: false,
		},
		{
			name: "expression without action",
			expr: domain.Expression{
				Type: "expression",
				Set: domain.Set{
					Type: "set",
					Variables: []domain.TypedVariable{
						{Type: "typedVariable", Name: "p", DataType: "Person"},
					},
				},
				Action: nil,
			},
			wantErr:     true,
			errContains: "action manquante",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateExpression(tt.expr, types)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateExpression() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errContains != "" {
				if !containsString(err.Error(), tt.errContains) {
					t.Errorf("Expected error to contain %q, got %q", tt.errContains, err.Error())
				}
			}
		})
	}
}

// TestConstraintValidator_ValidateProgram teste la validation complète d'un programme
func TestConstraintValidator_ValidateProgram(t *testing.T) {
	tests := []struct {
		name        string
		program     interface{}
		wantErr     bool
		errContains string
	}{
		{
			name: "valid complete program",
			program: &domain.Program{
				Types: []domain.TypeDefinition{
					{
						Type: "typeDefinition",
						Name: "Person",
						Fields: []domain.Field{
							{Name: "name", Type: "string"},
							{Name: "age", Type: "integer"},
						},
					},
				},
				Expressions: []domain.Expression{
					{
						Type: "expression",
						Set: domain.Set{
							Type: "set",
							Variables: []domain.TypedVariable{
								{Type: "typedVariable", Name: "p", DataType: "Person"},
							},
						},
						Action: &domain.Action{
							Type: "action",
							Job: domain.JobCall{
								Type: "jobCall",
								Name: "notify",
								Args: []interface{}{"test"},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:        "invalid program type",
			program:     "invalid",
			wantErr:     true,
			errContains: "invalid program type",
		},
		{
			name: "program with invalid type definition",
			program: &domain.Program{
				Types: []domain.TypeDefinition{
					{
						Type:   "typeDefinition",
						Name:   "Empty",
						Fields: []domain.Field{},
					},
				},
				Expressions: []domain.Expression{},
			},
			wantErr:     true,
			errContains: "must have at least one field",
		},
		{
			name: "program with duplicate types",
			program: &domain.Program{
				Types: []domain.TypeDefinition{
					{
						Type: "typeDefinition",
						Name: "Person",
						Fields: []domain.Field{
							{Name: "name", Type: "string"},
						},
					},
					{
						Type: "typeDefinition",
						Name: "Person",
						Fields: []domain.Field{
							{Name: "age", Type: "integer"},
						},
					},
				},
				Expressions: []domain.Expression{},
			},
			wantErr:     true,
			errContains: "duplicate type name",
		},
		{
			name: "program with expression referencing unknown type",
			program: &domain.Program{
				Types: []domain.TypeDefinition{
					{
						Type: "typeDefinition",
						Name: "Person",
						Fields: []domain.Field{
							{Name: "name", Type: "string"},
						},
					},
				},
				Expressions: []domain.Expression{
					{
						Type: "expression",
						Set: domain.Set{
							Type: "set",
							Variables: []domain.TypedVariable{
								{Type: "typedVariable", Name: "x", DataType: "UnknownType"},
							},
						},
						Action: &domain.Action{
							Type: "action",
							Job:  domain.JobCall{Type: "jobCall", Name: "test"},
						},
					},
				},
			},
			wantErr:     true,
			errContains: "UnknownType",
		},
		{
			name: "program with multiple valid expressions",
			program: &domain.Program{
				Types: []domain.TypeDefinition{
					{
						Type: "typeDefinition",
						Name: "Person",
						Fields: []domain.Field{
							{Name: "name", Type: "string"},
							{Name: "age", Type: "integer"},
						},
					},
					{
						Type: "typeDefinition",
						Name: "Product",
						Fields: []domain.Field{
							{Name: "price", Type: "number"},
						},
					},
				},
				Expressions: []domain.Expression{
					{
						Type: "expression",
						Set: domain.Set{
							Type: "set",
							Variables: []domain.TypedVariable{
								{Type: "typedVariable", Name: "p", DataType: "Person"},
							},
						},
						Action: &domain.Action{
							Type: "action",
							Job:  domain.JobCall{Type: "jobCall", Name: "action1"},
						},
					},
					{
						Type: "expression",
						Set: domain.Set{
							Type: "set",
							Variables: []domain.TypedVariable{
								{Type: "typedVariable", Name: "pr", DataType: "Product"},
							},
						},
						Action: &domain.Action{
							Type: "action",
							Job:  domain.JobCall{Type: "jobCall", Name: "action2"},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "program with empty types",
			program: &domain.Program{
				Types: []domain.TypeDefinition{},
				Expressions: []domain.Expression{
					{
						Type: "expression",
						Set: domain.Set{
							Type: "set",
							Variables: []domain.TypedVariable{
								{Type: "typedVariable", Name: "x", DataType: "AnyType"},
							},
						},
						Action: &domain.Action{
							Type: "action",
							Job:  domain.JobCall{Type: "jobCall", Name: "test"},
						},
					},
				},
			},
			wantErr:     true,
			errContains: "AnyType",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registry := NewTypeRegistry()
			checker := NewTypeChecker(registry)
			validator := NewConstraintValidator(registry, checker)

			err := validator.ValidateProgram(tt.program)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateProgram() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errContains != "" {
				if !containsString(err.Error(), tt.errContains) {
					t.Errorf("Expected error to contain %q, got %q", tt.errContains, err.Error())
				}
			}
		})
	}
}

// TestConstraintValidator_SetConfig teste la configuration du validateur
func TestConstraintValidator_SetConfig(t *testing.T) {
	registry := NewTypeRegistry()
	checker := NewTypeChecker(registry)
	validator := NewConstraintValidator(registry, checker)

	newConfig := domain.ValidatorConfig{
		StrictMode:       false,
		AllowedOperators: []string{"==", "!="},
		MaxDepth:         5,
	}

	validator.SetConfig(newConfig)
	gotConfig := validator.GetConfig()

	if gotConfig.StrictMode != newConfig.StrictMode {
		t.Errorf("StrictMode = %v, want %v", gotConfig.StrictMode, newConfig.StrictMode)
	}
	if len(gotConfig.AllowedOperators) != len(newConfig.AllowedOperators) {
		t.Errorf("AllowedOperators length = %d, want %d", len(gotConfig.AllowedOperators), len(newConfig.AllowedOperators))
	}
	if gotConfig.MaxDepth != newConfig.MaxDepth {
		t.Errorf("MaxDepth = %d, want %d", gotConfig.MaxDepth, newConfig.MaxDepth)
	}
}

// TestConstraintValidator_GetConfig teste la récupération de la configuration
func TestConstraintValidator_GetConfig(t *testing.T) {
	registry := NewTypeRegistry()
	checker := NewTypeChecker(registry)
	validator := NewConstraintValidator(registry, checker)

	config := validator.GetConfig()

	if !config.StrictMode {
		t.Error("Default StrictMode should be true")
	}
	if len(config.AllowedOperators) == 0 {
		t.Error("Default AllowedOperators should not be empty")
	}
	if config.MaxDepth != 10 {
		t.Errorf("Default MaxDepth = %d, want 10", config.MaxDepth)
	}
}

// ===== ActionValidator Tests =====

// TestNewActionValidator teste la création d'un nouveau validateur d'actions
func TestNewActionValidator(t *testing.T) {
	validator := NewActionValidator()
	if validator == nil {
		t.Fatal("NewActionValidator should not return nil")
	}
}

// TestActionValidator_ValidateAction teste la validation des actions
func TestActionValidator_ValidateAction(t *testing.T) {
	validator := NewActionValidator()

	tests := []struct {
		name        string
		action      *domain.Action
		wantErr     bool
		errContains string
	}{
		{
			name: "valid action with job name",
			action: &domain.Action{
				Type: "action",
				Job: domain.JobCall{
					Type: "jobCall",
					Name: "notify",
				},
			},
			wantErr: false,
		},
		{
			name: "valid action with job and args",
			action: &domain.Action{
				Type: "action",
				Job: domain.JobCall{
					Type: "jobCall",
					Name: "sendEmail",
					Args: []interface{}{"user@example.com", "Subject"},
				},
			},
			wantErr: false,
		},
		{
			name:        "nil action",
			action:      nil,
			wantErr:     true,
			errContains: "action cannot be nil",
		},
		{
			name: "action with empty job name",
			action: &domain.Action{
				Type: "action",
				Job: domain.JobCall{
					Type: "jobCall",
					Name: "",
				},
			},
			wantErr:     true,
			errContains: "job name cannot be empty",
		},
		{
			name: "action with whitespace job name",
			action: &domain.Action{
				Type: "action",
				Job: domain.JobCall{
					Type: "jobCall",
					Name: "   ",
				},
			},
			wantErr:     true,
			errContains: "job name cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateAction(tt.action)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errContains != "" {
				if !containsString(err.Error(), tt.errContains) {
					t.Errorf("Expected error to contain %q, got %q", tt.errContains, err.Error())
				}
			}
		})
	}
}

// TestActionValidator_ValidateJobCall teste la validation des appels de job
func TestActionValidator_ValidateJobCall(t *testing.T) {
	validator := NewActionValidator()

	tests := []struct {
		name        string
		jobCall     domain.JobCall
		wantErr     bool
		errContains string
	}{
		{
			name: "valid job without args",
			jobCall: domain.JobCall{
				Type: "jobCall",
				Name: "processData",
				Args: []interface{}{},
			},
			wantErr: false,
		},
		{
			name: "valid job with string args",
			jobCall: domain.JobCall{
				Type: "jobCall",
				Name: "sendNotification",
				Args: []interface{}{"message", "recipient"},
			},
			wantErr: false,
		},
		{
			name: "valid job with mixed args",
			jobCall: domain.JobCall{
				Type: "jobCall",
				Name: "updateRecord",
				Args: []interface{}{"id123", 42, true},
			},
			wantErr: false,
		},
		{
			name: "valid job with complex object args",
			jobCall: domain.JobCall{
				Type: "jobCall",
				Name: "processObject",
				Args: []interface{}{
					map[string]interface{}{"key": "value"},
				},
			},
			wantErr: false,
		},
		{
			name: "empty job name",
			jobCall: domain.JobCall{
				Type: "jobCall",
				Name: "",
			},
			wantErr:     true,
			errContains: "job name cannot be empty",
		},
		{
			name: "whitespace job name",
			jobCall: domain.JobCall{
				Type: "jobCall",
				Name: "  \t\n  ",
			},
			wantErr:     true,
			errContains: "job name cannot be empty",
		},
		{
			name: "job with empty string arg",
			jobCall: domain.JobCall{
				Type: "jobCall",
				Name: "test",
				Args: []interface{}{""},
			},
			wantErr:     true,
			errContains: "job argument 0 cannot be empty",
		},
		{
			name: "job with nil arg",
			jobCall: domain.JobCall{
				Type: "jobCall",
				Name: "test",
				Args: []interface{}{nil},
			},
			wantErr:     true,
			errContains: "job argument 0 cannot be nil",
		},
		{
			name: "job with whitespace string arg",
			jobCall: domain.JobCall{
				Type: "jobCall",
				Name: "test",
				Args: []interface{}{"valid", "   "},
			},
			wantErr:     true,
			errContains: "job argument 1 cannot be empty",
		},
		{
			name: "job with multiple valid args",
			jobCall: domain.JobCall{
				Type: "jobCall",
				Name: "complexJob",
				Args: []interface{}{
					"string",
					123,
					true,
					3.14,
					map[string]interface{}{"nested": "object"},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateJobCall(tt.jobCall)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJobCall() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errContains != "" {
				if !containsString(err.Error(), tt.errContains) {
					t.Errorf("Expected error to contain %q, got %q", tt.errContains, err.Error())
				}
			}
		})
	}
}

// TestActionValidator_ErrorTypes teste que les erreurs sont du bon type
func TestActionValidator_ErrorTypes(t *testing.T) {
	validator := NewActionValidator()

	// Test avec action nil
	err := validator.ValidateAction(nil)
	if err == nil {
		t.Fatal("Expected error for nil action")
	}
	if _, ok := err.(*domain.Error); !ok {
		t.Errorf("Expected *domain.Error, got %T", err)
	}
	domainErr := err.(*domain.Error)
	if domainErr.Type != domain.ActionError {
		t.Errorf("Expected ActionError type, got %v", domainErr.Type)
	}

	// Test avec job name vide
	err = validator.ValidateJobCall(domain.JobCall{Name: ""})
	if err == nil {
		t.Fatal("Expected error for empty job name")
	}
	if _, ok := err.(*domain.Error); !ok {
		t.Errorf("Expected *domain.Error, got %T", err)
	}
	domainErr = err.(*domain.Error)
	if domainErr.Type != domain.ActionError {
		t.Errorf("Expected ActionError type, got %v", domainErr.Type)
	}
}

// TestConstraintValidator_ValidateConstraint teste la validation des contraintes
func TestConstraintValidator_ValidateConstraint(t *testing.T) {
	registry := NewTypeRegistry()
	checker := NewTypeChecker(registry)
	validator := NewConstraintValidator(registry, checker)

	// Préparer les types
	types := []domain.TypeDefinition{
		{
			Type: "typeDefinition",
			Name: "Person",
			Fields: []domain.Field{
				{Name: "age", Type: "integer"},
			},
		},
	}

	variables := []domain.TypedVariable{
		{Type: "typedVariable", Name: "p", DataType: "Person"},
	}

	// La validation des contraintes délègue au type checker
	// Ici on teste juste que ça ne crash pas
	err := validator.ValidateConstraint(nil, variables, types)
	// On s'attend à ce que ça retourne nil ou une erreur valide (pas de panic)
	if err != nil && !domain.IsValidationError(err) {
		t.Errorf("ValidateConstraint should return nil or ValidationError, got %T", err)
	}
}
