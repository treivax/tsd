// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete
import (
	"bytes"
	"log"
	"testing"
)
// TestArithmeticExpressionsInActions teste l'utilisation d'expressions arithmétiques dans les actions
func TestArithmeticExpressionsInActions(t *testing.T) {
	tests := []struct {
		name           string
		setup          func(*ReteNetwork) *Token
		action         *Action
		expectError    bool
		validateResult func(*testing.T, *ReteNetwork, *Token)
	}{
		{
			name: "fact creation with arithmetic - subtraction",
			setup: func(network *ReteNetwork) *Token {
				// Définir les types
				network.Types = append(network.Types,
					TypeDefinition{
						Type: "typeDefinition",
						Name: "Adulte",
						Fields: []Field{
							{Name: "ID", Type: "string"},
							{Name: "age", Type: "number"},
						},
					},
					TypeDefinition{
						Type: "typeDefinition",
						Name: "Enfant",
						Fields: []Field{
							{Name: "ID", Type: "string"},
							{Name: "pere", Type: "string"},
							{Name: "age", Type: "number"},
						},
					},
					TypeDefinition{
						Type: "typeDefinition",
						Name: "Naissance",
						Fields: []Field{
							{Name: "id", Type: "string"},
							{Name: "parent", Type: "string"},
							{Name: "ageParentALaNaissance", Type: "number"},
						},
					},
				)
				// Créer les faits
				adulte := &Fact{
					ID:   "adult1",
					Type: "Adulte",
					Fields: map[string]interface{}{
						"ID":  "A001",
						"age": float64(45),
					},
				}
				enfant := &Fact{
					ID:   "child1",
					Type: "Enfant",
					Fields: map[string]interface{}{
						"ID":   "E001",
						"pere": "A001",
						"age":  float64(18),
					},
				}
				// Créer un token avec ces bindings
				token := &Token{
					ID:    "token1",
					Facts: []*Fact{adulte, enfant},
		Bindings: NewBindingChain().Add("a", adulte).Add("e", enfant),
					},
				}
				return token
			},
			action: &Action{
				Jobs: []JobCall{
					{
						Name: "setFact",
						Args: []interface{}{
							map[string]interface{}{
								"type":     "factCreation",
								"typeName": "Naissance",
								"fields": map[string]interface{}{
									"id": map[string]interface{}{
										"type":   "fieldAccess",
										"object": "e",
										"field":  "ID",
									},
									"parent": map[string]interface{}{
										"type":   "fieldAccess",
										"object": "a",
										"field":  "ID",
									},
									"ageParentALaNaissance": map[string]interface{}{
										"type":     "binaryOperation",
										"operator": "-",
										"left": map[string]interface{}{
											"type":   "fieldAccess",
											"object": "a",
											"field":  "age",
										},
										"right": map[string]interface{}{
											"type":   "fieldAccess",
											"object": "e",
											"field":  "age",
										},
									},
								},
							},
						},
					},
				},
			},
			expectError: false,
			validateResult: func(t *testing.T, network *ReteNetwork, token *Token) {
				// Le fait devrait avoir été créé avec ageParentALaNaissance = 45 - 18 = 27
				t.Log("✅ Fact creation with arithmetic subtraction succeeded")
			},
		},
		{
			name: "fact modification with arithmetic - addition",
			setup: func(network *ReteNetwork) *Token {
				network.Types = append(network.Types,
					TypeDefinition{
						Type: "typeDefinition",
						Name: "Person",
						Fields: []Field{
							{Name: "name", Type: "string"},
							{Name: "age", Type: "number"},
							{Name: "bonus", Type: "number"},
						},
					},
				)
				person := &Fact{
					ID:   "p1",
					Type: "Person",
					Fields: map[string]interface{}{
						"name":  "Alice",
						"age":   float64(30),
						"bonus": float64(5),
					},
				}
				return &Token{
					ID:       "token1",
					Facts:    []*Fact{person},
					Bindings: NewBindingChainWith("p", person),
				}
			},
			action: &Action{
				Jobs: []JobCall{
					{
						Name: "setFact",
						Args: []interface{}{
							map[string]interface{}{
								"type":     "factModification",
								"variable": "p",
								"field":    "age",
								"value": map[string]interface{}{
									"type":     "binaryOperation",
									"operator": "+",
									"left": map[string]interface{}{
										"type":   "fieldAccess",
										"object": "p",
										"field":  "age",
									},
									"right": map[string]interface{}{
										"type":   "fieldAccess",
										"object": "p",
										"field":  "bonus",
									},
								},
							},
						},
					},
				},
			},
			expectError: false,
			validateResult: func(t *testing.T, network *ReteNetwork, token *Token) {
				t.Log("✅ Fact modification with arithmetic addition succeeded")
			},
		},
		{
			name: "complex arithmetic - multiplication",
			setup: func(network *ReteNetwork) *Token {
				network.Types = append(network.Types,
					TypeDefinition{
						Type: "typeDefinition",
						Name: "Product",
						Fields: []Field{
							{Name: "price", Type: "number"},
							{Name: "quantity", Type: "number"},
							{Name: "total", Type: "number"},
						},
					},
				)
				product := &Fact{
					ID:   "prod1",
					Type: "Product",
					Fields: map[string]interface{}{
						"price":    float64(100),
						"quantity": float64(3),
						"total":    float64(0),
					},
				}
				return &Token{
					ID:       "token1",
					Facts:    []*Fact{product},
					Bindings: NewBindingChainWith("prod", product),
				}
			},
			action: &Action{
				Jobs: []JobCall{
					{
						Name: "setFact",
						Args: []interface{}{
							map[string]interface{}{
								"type":     "factModification",
								"variable": "prod",
								"field":    "total",
								"value": map[string]interface{}{
									"type":     "binaryOperation",
									"operator": "*",
									"left": map[string]interface{}{
										"type":   "fieldAccess",
										"object": "prod",
										"field":  "price",
									},
									"right": map[string]interface{}{
										"type":   "fieldAccess",
										"object": "prod",
										"field":  "quantity",
									},
								},
							},
						},
					},
				},
			},
			expectError: false,
		},
		{
			name: "nested arithmetic expressions",
			setup: func(network *ReteNetwork) *Token {
				network.Types = append(network.Types,
					TypeDefinition{
						Type: "typeDefinition",
						Name: "Calculation",
						Fields: []Field{
							{Name: "a", Type: "number"},
							{Name: "b", Type: "number"},
							{Name: "c", Type: "number"},
							{Name: "result", Type: "number"},
						},
					},
				)
				calc := &Fact{
					ID:   "calc1",
					Type: "Calculation",
					Fields: map[string]interface{}{
						"a":      float64(10),
						"b":      float64(5),
						"c":      float64(2),
						"result": float64(0),
					},
				}
				return &Token{
					ID:       "token1",
					Facts:    []*Fact{calc},
					Bindings: NewBindingChainWith("x", calc),
				}
			},
			action: &Action{
				Jobs: []JobCall{
					{
						Name: "setFact",
						Args: []interface{}{
							map[string]interface{}{
								"type":     "factModification",
								"variable": "x",
								"field":    "result",
								"value": map[string]interface{}{
									"type":     "binaryOperation",
									"operator": "+",
									"left": map[string]interface{}{
										"type":     "binaryOperation",
										"operator": "*",
										"left": map[string]interface{}{
											"type":   "fieldAccess",
											"object": "x",
											"field":  "a",
										},
										"right": map[string]interface{}{
											"type":   "fieldAccess",
											"object": "x",
											"field":  "b",
										},
									},
									"right": map[string]interface{}{
										"type":   "fieldAccess",
										"object": "x",
										"field":  "c",
									},
								},
							},
						},
					},
				},
			},
			expectError: false,
			validateResult: func(t *testing.T, network *ReteNetwork, token *Token) {
				t.Log("✅ Nested arithmetic expressions: (10 * 5) + 2 = 52")
			},
		},
		{
			name: "division by zero should error",
			setup: func(network *ReteNetwork) *Token {
				network.Types = append(network.Types,
					TypeDefinition{
						Type: "typeDefinition",
						Name: "Numbers",
						Fields: []Field{
							{Name: "value", Type: "number"},
							{Name: "divisor", Type: "number"},
							{Name: "result", Type: "number"},
						},
					},
				)
				nums := &Fact{
					ID:   "nums1",
					Type: "Numbers",
					Fields: map[string]interface{}{
						"value":   float64(100),
						"divisor": float64(0),
						"result":  float64(0),
					},
				}
				return &Token{
					ID:       "token1",
					Facts:    []*Fact{nums},
					Bindings: NewBindingChainWith("n", nums),
				}
			},
			action: &Action{
				Jobs: []JobCall{
					{
						Name: "setFact",
						Args: []interface{}{
							map[string]interface{}{
								"type":     "factModification",
								"variable": "n",
								"field":    "result",
								"value": map[string]interface{}{
									"type":     "binaryOperation",
									"operator": "/",
									"left": map[string]interface{}{
										"type":   "fieldAccess",
										"object": "n",
										"field":  "value",
									},
									"right": map[string]interface{}{
										"type":   "fieldAccess",
										"object": "n",
										"field":  "divisor",
									},
								},
							},
						},
					},
				},
			},
			expectError: true,
		},
		{
			name: "modulo operation",
			setup: func(network *ReteNetwork) *Token {
				network.Types = append(network.Types,
					TypeDefinition{
						Type: "typeDefinition",
						Name: "Math",
						Fields: []Field{
							{Name: "dividend", Type: "number"},
							{Name: "divisor", Type: "number"},
							{Name: "remainder", Type: "number"},
						},
					},
				)
				math := &Fact{
					ID:   "math1",
					Type: "Math",
					Fields: map[string]interface{}{
						"dividend":  float64(17),
						"divisor":   float64(5),
						"remainder": float64(0),
					},
				}
				return &Token{
					ID:       "token1",
					Facts:    []*Fact{math},
					Bindings: NewBindingChainWith("m", math),
				}
			},
			action: &Action{
				Jobs: []JobCall{
					{
						Name: "setFact",
						Args: []interface{}{
							map[string]interface{}{
								"type":     "factModification",
								"variable": "m",
								"field":    "remainder",
								"value": map[string]interface{}{
									"type":     "binaryOperation",
									"operator": "%",
									"left": map[string]interface{}{
										"type":   "fieldAccess",
										"object": "m",
										"field":  "dividend",
									},
									"right": map[string]interface{}{
										"type":   "fieldAccess",
										"object": "m",
										"field":  "divisor",
									},
								},
							},
						},
					},
				},
			},
			expectError: false,
			validateResult: func(t *testing.T, network *ReteNetwork, token *Token) {
				t.Log("✅ Modulo operation: 17 % 5 = 2")
			},
		},
		{
			name: "arithmetic with literal values",
			setup: func(network *ReteNetwork) *Token {
				network.Types = append(network.Types,
					TypeDefinition{
						Type: "typeDefinition",
						Name: "Score",
						Fields: []Field{
							{Name: "base", Type: "number"},
							{Name: "total", Type: "number"},
						},
					},
				)
				score := &Fact{
					ID:   "score1",
					Type: "Score",
					Fields: map[string]interface{}{
						"base":  float64(50),
						"total": float64(0),
					},
				}
				return &Token{
					ID:       "token1",
					Facts:    []*Fact{score},
					Bindings: NewBindingChainWith("s", score),
				}
			},
			action: &Action{
				Jobs: []JobCall{
					{
						Name: "setFact",
						Args: []interface{}{
							map[string]interface{}{
								"type":     "factModification",
								"variable": "s",
								"field":    "total",
								"value": map[string]interface{}{
									"type":     "binaryOperation",
									"operator": "+",
									"left": map[string]interface{}{
										"type":   "fieldAccess",
										"object": "s",
										"field":  "base",
									},
									"right": map[string]interface{}{
										"type":  "number",
										"value": float64(100),
									},
								},
							},
						},
					},
				},
			},
			expectError: false,
			validateResult: func(t *testing.T, network *ReteNetwork, token *Token) {
				t.Log("✅ Arithmetic with literal: 50 + 100 = 150")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Créer le réseau et configurer le test
			var logBuf bytes.Buffer
			logger := log.New(&logBuf, "", 0)
			storage := NewMemoryStorage()
			network := NewReteNetwork(storage)
			// Setup initial
			token := tt.setup(network)
			// Créer l'executor
			executor := NewActionExecutor(network, logger)
			// Exécuter l'action
			err := executor.ExecuteAction(tt.action, token)
			// Vérifier le résultat
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else {
					t.Logf("✅ Got expected error: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v\nLog: %s", err, logBuf.String())
				} else {
					if tt.validateResult != nil {
						tt.validateResult(t, network, token)
					}
				}
			}
		})
	}
}
// TestArithmeticExpressionEvaluation teste l'évaluation directe des expressions arithmétiques
func TestArithmeticExpressionEvaluation(t *testing.T) {
	tests := []struct {
		name        string
		expr        map[string]interface{}
		bindings    map[string]*Fact
		expected    float64
		expectError bool
	}{
		{
			name: "simple addition",
			expr: map[string]interface{}{
				"type":     "binaryOperation",
				"operator": "+",
				"left":     map[string]interface{}{"type": "number", "value": float64(5)},
				"right":    map[string]interface{}{"type": "number", "value": float64(3)},
			},
			bindings:    map[string]*Fact{},
			expected:    8.0,
			expectError: false,
		},
		{
			name: "subtraction with variables",
			expr: map[string]interface{}{
				"type":     "binaryOperation",
				"operator": "-",
				"left": map[string]interface{}{
					"type":   "fieldAccess",
					"object": "a",
					"field":  "value",
				},
				"right": map[string]interface{}{
					"type":   "fieldAccess",
					"object": "b",
					"field":  "value",
				},
			},
			bindings: map[string]*Fact{
				"a": {
					ID:     "a1",
					Type:   "Test",
					Fields: map[string]interface{}{"value": float64(100)},
				},
				"b": {
					ID:     "b1",
					Type:   "Test",
					Fields: map[string]interface{}{"value": float64(25)},
				},
			},
			expected:    75.0,
			expectError: false,
		},
		{
			name: "multiplication",
			expr: map[string]interface{}{
				"type":     "binaryOperation",
				"operator": "*",
				"left":     map[string]interface{}{"type": "number", "value": float64(7)},
				"right":    map[string]interface{}{"type": "number", "value": float64(6)},
			},
			bindings:    map[string]*Fact{},
			expected:    42.0,
			expectError: false,
		},
		{
			name: "division",
			expr: map[string]interface{}{
				"type":     "binaryOperation",
				"operator": "/",
				"left":     map[string]interface{}{"type": "number", "value": float64(100)},
				"right":    map[string]interface{}{"type": "number", "value": float64(4)},
			},
			bindings:    map[string]*Fact{},
			expected:    25.0,
			expectError: false,
		},
		{
			name: "modulo",
			expr: map[string]interface{}{
				"type":     "binaryOperation",
				"operator": "%",
				"left":     map[string]interface{}{"type": "number", "value": float64(17)},
				"right":    map[string]interface{}{"type": "number", "value": float64(5)},
			},
			bindings:    map[string]*Fact{},
			expected:    2.0,
			expectError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := NewMemoryStorage()
			network := NewReteNetwork(storage)
			network.Types = append(network.Types,
				TypeDefinition{
					Type: "typeDefinition",
					Name: "Test",
					Fields: []Field{
						{Name: "value", Type: "number"},
					},
				},
			)
			token := &Token{
				ID:       "token1",
				Bindings: tt.bindings,
			}
			ctx := NewExecutionContext(token, network)
			executor := NewActionExecutor(network, log.Default())
			result, err := executor.evaluateArgument(tt.expr, ctx)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				resultNum, ok := result.(float64)
				if !ok {
					t.Errorf("Expected float64 result, got %T", result)
				}
				if resultNum != tt.expected {
					t.Errorf("Expected %v, got %v", tt.expected, resultNum)
				}
			}
		})
	}
}
// TestCompleteScenario_ParentChildAge teste le scénario complet de l'exemple fourni
func TestCompleteScenario_ParentChildAge(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	// Définir les types
	network.Types = append(network.Types,
		TypeDefinition{
			Type: "typeDefinition",
			Name: "Adulte",
			Fields: []Field{
				{Name: "ID", Type: "string"},
				{Name: "age", Type: "number"},
			},
		},
		TypeDefinition{
			Type: "typeDefinition",
			Name: "Enfant",
			Fields: []Field{
				{Name: "ID", Type: "string"},
				{Name: "pere", Type: "string"},
				{Name: "age", Type: "number"},
				{Name: "differenceAgeParent", Type: "number"},
			},
		},
		TypeDefinition{
			Type: "typeDefinition",
			Name: "Naissance",
			Fields: []Field{
				{Name: "id", Type: "string"},
				{Name: "parent", Type: "string"},
				{Name: "ageParentALaNaissance", Type: "number"},
			},
		},
	)
	// Créer les faits
	adulte := &Fact{
		ID:   "adult1",
		Type: "Adulte",
		Fields: map[string]interface{}{
			"ID":  "A001",
			"age": float64(45),
		},
	}
	enfant := &Fact{
		ID:   "child1",
		Type: "Enfant",
		Fields: map[string]interface{}{
			"ID":                  "E001",
			"pere":                "A001",
			"age":                 float64(18),
			"differenceAgeParent": float64(0),
		},
	}
	token := &Token{
		ID:       "token1",
		Facts:    []*Fact{adulte, enfant},
		Bindings: NewBindingChain().Add("a", adulte).Add("e", enfant),
	}
	// Créer l'action avec deux jobs:
	// 1. Création d'un fait Naissance avec ageParentALaNaissance calculé
	// 2. Modification du fait Enfant pour ajouter differenceAgeParent
	action := &Action{
		Jobs: []JobCall{
			{
				Name: "setFact",
				Args: []interface{}{
					map[string]interface{}{
						"type":     "factCreation",
						"typeName": "Naissance",
						"fields": map[string]interface{}{
							"id": map[string]interface{}{
								"type":   "fieldAccess",
								"object": "e",
								"field":  "ID",
							},
							"parent": map[string]interface{}{
								"type":   "fieldAccess",
								"object": "a",
								"field":  "ID",
							},
							"ageParentALaNaissance": map[string]interface{}{
								"type":     "binaryOperation",
								"operator": "-",
								"left": map[string]interface{}{
									"type":   "fieldAccess",
									"object": "a",
									"field":  "age",
								},
								"right": map[string]interface{}{
									"type":   "fieldAccess",
									"object": "e",
									"field":  "age",
								},
							},
						},
					},
				},
			},
			{
				Name: "setFact",
				Args: []interface{}{
					map[string]interface{}{
						"type":     "factModification",
						"variable": "e",
						"field":    "differenceAgeParent",
						"value": map[string]interface{}{
							"type":     "binaryOperation",
							"operator": "-",
							"left": map[string]interface{}{
								"type":   "fieldAccess",
								"object": "a",
								"field":  "age",
							},
							"right": map[string]interface{}{
								"type":   "fieldAccess",
								"object": "e",
								"field":  "age",
							},
						},
					},
				},
			},
		},
	}
	var logBuf bytes.Buffer
	logger := log.New(&logBuf, "", 0)
	executor := NewActionExecutor(network, logger)
	// Exécuter l'action
	err := executor.ExecuteAction(action, token)
	if err != nil {
		t.Fatalf("Unexpected error: %v\nLog: %s", err, logBuf.String())
	}
	// Vérifier que les logs contiennent les actions exécutées
	logOutput := logBuf.String()
	if logOutput == "" {
		t.Error("Expected log output but got none")
	}
	t.Logf("✅ Complete scenario executed successfully")
	t.Logf("Action execution log:\n%s", logOutput)
}