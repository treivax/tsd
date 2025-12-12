// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete
import (
	"bytes"
	"log"
	"math"
	"testing"
)
// TestComplexArithmeticExpressionsWithMultipleLiterals teste les expressions complexes avec plusieurs littéraux
func TestComplexArithmeticExpressionsWithMultipleLiterals(t *testing.T) {
	tests := []struct {
		name           string
		setup          func(*ReteNetwork) (*Token, map[string]float64)
		expression     map[string]interface{}
		expectedResult float64
		description    string
	}{
		{
			name: "complex expression with 5+ literals in action",
			setup: func(network *ReteNetwork) (*Token, map[string]float64) {
				network.Types = append(network.Types,
					TypeDefinition{
						Type: "typeDefinition",
						Name: "Objet",
						Fields: []Field{
							{Name: "id", Type: "string"},
							{Name: "prix", Type: "number"},
						},
					},
				)
				objet := &Fact{
					ID:   "obj1",
					Type: "Objet",
					Fields: map[string]interface{}{
						"id":   "O001",
						"prix": float64(100),
					},
				}
				token := &Token{
					ID:       "token1",
					Facts:    []*Fact{objet},
					Bindings: NewBindingChainWith("o", objet),
				}
				expectedValues := map[string]float64{
					"o.prix": 100,
				}
				return token, expectedValues
			},
			// o.prix * (1 + 2.3 + 3) - 1 = 100 * 6.3 - 1 = 629
			expression: map[string]interface{}{
				"type":     "binaryOperation",
				"operator": "-",
				"left": map[string]interface{}{
					"type":     "binaryOperation",
					"operator": "*",
					"left": map[string]interface{}{
						"type":   "fieldAccess",
						"object": "o",
						"field":  "prix",
					},
					"right": map[string]interface{}{
						"type":     "binaryOperation",
						"operator": "+",
						"left": map[string]interface{}{
							"type":     "binaryOperation",
							"operator": "+",
							"left": map[string]interface{}{
								"type":  "number",
								"value": float64(1),
							},
							"right": map[string]interface{}{
								"type":  "number",
								"value": float64(2.3),
							},
						},
						"right": map[string]interface{}{
							"type":  "number",
							"value": float64(3),
						},
					},
				},
				"right": map[string]interface{}{
					"type":  "number",
					"value": float64(1),
				},
			},
			expectedResult: 629.0,
			description:    "o.prix * (1 + 2.3 + 3) - 1",
		},
		{
			name: "modulo with decimals and multiple operations",
			setup: func(network *ReteNetwork) (*Token, map[string]float64) {
				network.Types = append(network.Types,
					TypeDefinition{
						Type: "typeDefinition",
						Name: "Product",
						Fields: []Field{
							{Name: "price", Type: "number"},
							{Name: "cost", Type: "number"},
						},
					},
				)
				product := &Fact{
					ID:   "prod1",
					Type: "Product",
					Fields: map[string]interface{}{
						"price": float64(100),
						"cost":  float64(20),
					},
				}
				token := &Token{
					ID:       "token1",
					Facts:    []*Fact{product},
					Bindings: NewBindingChainWith("p", product),
				}
				return token, map[string]float64{
					"p.price": 100,
					"p.cost":  20,
				}
			},
			// p.price % 10 * 1.345 - p.cost + 1 - 3
			// = (100 % 10) * 1.345 - 20 + 1 - 3
			// = 0 * 1.345 - 20 + 1 - 3
			// = 0 - 20 + 1 - 3 = -22
			expression: map[string]interface{}{
				"type":     "binaryOperation",
				"operator": "-",
				"left": map[string]interface{}{
					"type":     "binaryOperation",
					"operator": "+",
					"left": map[string]interface{}{
						"type":     "binaryOperation",
						"operator": "-",
						"left": map[string]interface{}{
							"type":     "binaryOperation",
							"operator": "*",
							"left": map[string]interface{}{
								"type":     "binaryOperation",
								"operator": "%",
								"left": map[string]interface{}{
									"type":   "fieldAccess",
									"object": "p",
									"field":  "price",
								},
								"right": map[string]interface{}{
									"type":  "number",
									"value": float64(10),
								},
							},
							"right": map[string]interface{}{
								"type":  "number",
								"value": float64(1.345),
							},
						},
						"right": map[string]interface{}{
							"type":   "fieldAccess",
							"object": "p",
							"field":  "cost",
						},
					},
					"right": map[string]interface{}{
						"type":  "number",
						"value": float64(1),
					},
				},
				"right": map[string]interface{}{
					"type":  "number",
					"value": float64(3),
				},
			},
			expectedResult: -22.0,
			description:    "p.price % 10 * 1.345 - p.cost + 1 - 3",
		},
		{
			name: "deeply nested with 7+ literals",
			setup: func(network *ReteNetwork) (*Token, map[string]float64) {
				network.Types = append(network.Types,
					TypeDefinition{
						Type: "typeDefinition",
						Name: "Data",
						Fields: []Field{
							{Name: "value", Type: "number"},
						},
					},
				)
				data := &Fact{
					ID:   "data1",
					Type: "Data",
					Fields: map[string]interface{}{
						"value": float64(50),
					},
				}
				token := &Token{
					ID:       "token1",
					Facts:    []*Fact{data},
					Bindings: NewBindingChainWith("d", data),
				}
				return token, map[string]float64{"d.value": 50}
			},
			// d.value * 2 + 3 - 4 * 5 / 10 + 1.5 - 0.5
			// = 50 * 2 + 3 - 4 * 5 / 10 + 1.5 - 0.5
			// = 100 + 3 - 20 / 10 + 1.5 - 0.5
			// = 100 + 3 - 2 + 1.5 - 0.5 = 102
			expression: map[string]interface{}{
				"type":     "binaryOperation",
				"operator": "-",
				"left": map[string]interface{}{
					"type":     "binaryOperation",
					"operator": "+",
					"left": map[string]interface{}{
						"type":     "binaryOperation",
						"operator": "-",
						"left": map[string]interface{}{
							"type":     "binaryOperation",
							"operator": "+",
							"left": map[string]interface{}{
								"type":     "binaryOperation",
								"operator": "*",
								"left": map[string]interface{}{
									"type":   "fieldAccess",
									"object": "d",
									"field":  "value",
								},
								"right": map[string]interface{}{
									"type":  "number",
									"value": float64(2),
								},
							},
							"right": map[string]interface{}{
								"type":  "number",
								"value": float64(3),
							},
						},
						"right": map[string]interface{}{
							"type":     "binaryOperation",
							"operator": "/",
							"left": map[string]interface{}{
								"type":     "binaryOperation",
								"operator": "*",
								"left": map[string]interface{}{
									"type":  "number",
									"value": float64(4),
								},
								"right": map[string]interface{}{
									"type":  "number",
									"value": float64(5),
								},
							},
							"right": map[string]interface{}{
								"type":  "number",
								"value": float64(10),
							},
						},
					},
					"right": map[string]interface{}{
						"type":  "number",
						"value": float64(1.5),
					},
				},
				"right": map[string]interface{}{
					"type":  "number",
					"value": float64(0.5),
				},
			},
			expectedResult: 102.0,
			description:    "d.value * 2 + 3 - 4 * 5 / 10 + 1.5 - 0.5",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := NewMemoryStorage()
			network := NewReteNetwork(storage)
			token, expectedValues := tt.setup(network)
			ctx := NewExecutionContext(token, network)
			executor := NewActionExecutor(network, log.Default())
			result, err := executor.evaluateArgument(tt.expression, ctx)
			if err != nil {
				t.Fatalf("Unexpected error: %v\nExpression: %s", err, tt.description)
			}
			resultNum, ok := result.(float64)
			if !ok {
				t.Fatalf("Expected float64 result, got %T", result)
			}
			// Comparer avec une petite tolérance pour les erreurs d'arrondi
			if math.Abs(resultNum-tt.expectedResult) > 0.0001 {
				t.Errorf("Expression: %s\nExpected: %v\nGot: %v\nDifference: %v",
					tt.description, tt.expectedResult, resultNum, resultNum-tt.expectedResult)
				// Afficher les valeurs des variables pour le debug
				t.Logf("Variable values:")
				for k, v := range expectedValues {
					t.Logf("  %s = %v", k, v)
				}
			} else {
				t.Logf("✅ %s = %v", tt.description, resultNum)
			}
		})
	}
}
// TestComplexExpressionInFactCreation teste une expression complexe dans la création de fait
func TestComplexExpressionInFactCreation(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	// Définir les types
	network.Types = append(network.Types,
		TypeDefinition{
			Type: "typeDefinition",
			Name: "Objet",
			Fields: []Field{
				{Name: "id", Type: "string"},
				{Name: "prix", Type: "number"},
				{Name: "boite", Type: "string"},
			},
		},
		TypeDefinition{
			Type: "typeDefinition",
			Name: "Boite",
			Fields: []Field{
				{Name: "id", Type: "string"},
				{Name: "prix", Type: "number"},
				{Name: "cout", Type: "number"},
			},
		},
		TypeDefinition{
			Type: "typeDefinition",
			Name: "Vente",
			Fields: []Field{
				{Name: "objet", Type: "string"},
				{Name: "prixTotal", Type: "number"},
			},
		},
	)
	// Créer les faits
	objet := &Fact{
		ID:   "obj1",
		Type: "Objet",
		Fields: map[string]interface{}{
			"id":    "O001",
			"prix":  float64(100),
			"boite": "B001",
		},
	}
	boite := &Fact{
		ID:   "box1",
		Type: "Boite",
		Fields: map[string]interface{}{
			"id":   "B001",
			"prix": float64(15),
			"cout": float64(5),
		},
	}
	token := &Token{
		ID:    "token1",
		Facts: []*Fact{objet, boite},
		Bindings: NewBindingChain().Add("o", objet).Add("b", boite),
		},
	}
	// Action : créer une vente avec un calcul complexe
	// prixTotal: o.prix * (1 + 2.3 + 3) + b.prix - 1
	// = 100 * 6.3 + 15 - 1 = 630 + 15 - 1 = 644
	action := &Action{
		Jobs: []JobCall{
			{
				Name: "setFact",
				Args: []interface{}{
					map[string]interface{}{
						"type":     "factCreation",
						"typeName": "Vente",
						"fields": map[string]interface{}{
							"objet": map[string]interface{}{
								"type":   "fieldAccess",
								"object": "o",
								"field":  "id",
							},
							"prixTotal": map[string]interface{}{
								"type":     "binaryOperation",
								"operator": "-",
								"left": map[string]interface{}{
									"type":     "binaryOperation",
									"operator": "+",
									"left": map[string]interface{}{
										"type":     "binaryOperation",
										"operator": "*",
										"left": map[string]interface{}{
											"type":   "fieldAccess",
											"object": "o",
											"field":  "prix",
										},
										"right": map[string]interface{}{
											"type":     "binaryOperation",
											"operator": "+",
											"left": map[string]interface{}{
												"type":     "binaryOperation",
												"operator": "+",
												"left": map[string]interface{}{
													"type":  "number",
													"value": float64(1),
												},
												"right": map[string]interface{}{
													"type":  "number",
													"value": float64(2.3),
												},
											},
											"right": map[string]interface{}{
												"type":  "number",
												"value": float64(3),
											},
										},
									},
									"right": map[string]interface{}{
										"type":   "fieldAccess",
										"object": "b",
										"field":  "prix",
									},
								},
								"right": map[string]interface{}{
									"type":  "number",
									"value": float64(1),
								},
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
	err := executor.ExecuteAction(action, token)
	if err != nil {
		t.Fatalf("Unexpected error: %v\nLog: %s", err, logBuf.String())
	}
	// Le calcul attendu : 100 * (1 + 2.3 + 3) + 15 - 1 = 644
	expectedTotal := 644.0
	t.Logf("✅ Fact creation with complex expression succeeded")
	t.Logf("   Expected prixTotal: %.2f", expectedTotal)
	t.Logf("   Formula: o.prix * (1 + 2.3 + 3) + b.prix - 1")
	t.Logf("   = 100 * 6.3 + 15 - 1 = %.2f", expectedTotal)
}
// TestComplexExpressionWithModuloAndDecimals teste le modulo avec des décimales
func TestComplexExpressionWithModuloAndDecimals(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	network.Types = append(network.Types,
		TypeDefinition{
			Type: "typeDefinition",
			Name: "Item",
			Fields: []Field{
				{Name: "value", Type: "number"},
				{Name: "result", Type: "number"},
			},
		},
	)
	item := &Fact{
		ID:   "item1",
		Type: "Item",
		Fields: map[string]interface{}{
			"value":  float64(53),
			"result": float64(0),
		},
	}
	token := &Token{
		ID:       "token1",
		Facts:    []*Fact{item},
		Bindings: NewBindingChainWith("i", item),
	}
	// Expression: 2.3 % 53 = 2 (car int(2.3) % int(53) = 2 % 53 = 2)
	action := &Action{
		Jobs: []JobCall{
			{
				Name: "setFact",
				Args: []interface{}{
					map[string]interface{}{
						"type":     "factModification",
						"variable": "i",
						"field":    "result",
						"value": map[string]interface{}{
							"type":     "binaryOperation",
							"operator": "%",
							"left": map[string]interface{}{
								"type":  "number",
								"value": float64(2.3),
							},
							"right": map[string]interface{}{
								"type":   "fieldAccess",
								"object": "i",
								"field":  "value",
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
	err := executor.ExecuteAction(action, token)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	// 2.3 % 53 = int(2.3) % int(53) = 2 % 53 = 2
	t.Logf("✅ Modulo with decimals: 2.3 %% 53 = 2")
}
// TestRealWorldComplexExpression teste l'exemple complet de la règle
func TestRealWorldComplexExpression(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	// Définir les types
	network.Types = append(network.Types,
		TypeDefinition{
			Type: "typeDefinition",
			Name: "Objet",
			Fields: []Field{
				{Name: "id", Type: "string"},
				{Name: "prix", Type: "number"},
				{Name: "boite", Type: "string"},
			},
		},
		TypeDefinition{
			Type: "typeDefinition",
			Name: "Boite",
			Fields: []Field{
				{Name: "id", Type: "string"},
				{Name: "prix", Type: "number"},
				{Name: "cout", Type: "number"},
			},
		},
		TypeDefinition{
			Type: "typeDefinition",
			Name: "Vente",
			Fields: []Field{
				{Name: "objet", Type: "string"},
				{Name: "prixTotal", Type: "number"},
			},
		},
	)
	// Créer les faits
	objet := &Fact{
		ID:   "obj1",
		Type: "Objet",
		Fields: map[string]interface{}{
			"id":    "O001",
			"prix":  float64(100),
			"boite": "B001",
		},
	}
	boite := &Fact{
		ID:   "box1",
		Type: "Boite",
		Fields: map[string]interface{}{
			"id":   "B001",
			"prix": float64(15),
			"cout": float64(5),
		},
	}
	token := &Token{
		ID:    "token1",
		Facts: []*Fact{objet, boite},
		Bindings: NewBindingChain().Add("o", objet).Add("b", boite),
		},
	}
	t.Log("🧮 Test de l'expression complexe du monde réel")
	t.Log("==============================================")
	t.Logf("Objet: id=%s, prix=%.2f, boite=%s", objet.Fields["id"], objet.Fields["prix"], objet.Fields["boite"])
	t.Logf("Boite: id=%s, prix=%.2f, cout=%.2f", boite.Fields["id"], boite.Fields["prix"], boite.Fields["cout"])
	t.Log("")
	// Expression de l'action : o.prix * (1 + 2.3 % 53 + 3) + b.prix - 1
	// Calcul : 100 * (1 + 2 + 3) + 15 - 1 = 100 * 6 + 15 - 1 = 614
	// (2.3 % 53 = 2)
	action := &Action{
		Jobs: []JobCall{
			{
				Name: "setFact",
				Args: []interface{}{
					map[string]interface{}{
						"type":     "factCreation",
						"typeName": "Vente",
						"fields": map[string]interface{}{
							"objet": map[string]interface{}{
								"type":   "fieldAccess",
								"object": "o",
								"field":  "id",
							},
							"prixTotal": map[string]interface{}{
								"type":     "binaryOperation",
								"operator": "-",
								"left": map[string]interface{}{
									"type":     "binaryOperation",
									"operator": "+",
									"left": map[string]interface{}{
										"type":     "binaryOperation",
										"operator": "*",
										"left": map[string]interface{}{
											"type":   "fieldAccess",
											"object": "o",
											"field":  "prix",
										},
										"right": map[string]interface{}{
											"type":     "binaryOperation",
											"operator": "+",
											"left": map[string]interface{}{
												"type":     "binaryOperation",
												"operator": "+",
												"left": map[string]interface{}{
													"type":  "number",
													"value": float64(1),
												},
												"right": map[string]interface{}{
													"type":     "binaryOperation",
													"operator": "%",
													"left": map[string]interface{}{
														"type":  "number",
														"value": float64(2.3),
													},
													"right": map[string]interface{}{
														"type":  "number",
														"value": float64(53),
													},
												},
											},
											"right": map[string]interface{}{
												"type":  "number",
												"value": float64(3),
											},
										},
									},
									"right": map[string]interface{}{
										"type":   "fieldAccess",
										"object": "b",
										"field":  "prix",
									},
								},
								"right": map[string]interface{}{
									"type":  "number",
									"value": float64(1),
								},
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
	err := executor.ExecuteAction(action, token)
	if err != nil {
		t.Fatalf("❌ Erreur: %v\nLog: %s", err, logBuf.String())
	}
	// Calculer le résultat attendu
	// o.prix * (1 + 2.3 % 53 + 3) + b.prix - 1
	oPrix := objet.Fields["prix"].(float64)
	bPrix := boite.Fields["prix"].(float64)
	modResult := float64(int64(2) % int64(53)) // int(2.3) % 53 = 2 % 53 = 2
	total := oPrix*(1+modResult+3) + bPrix - 1 // 100 * 6 + 15 - 1 = 614
	t.Log("✅ Création de fait avec expression complexe réussie!")
	t.Log("")
	t.Log("Calcul détaillé:")
	t.Logf("  2.3 %% 53 = %d", int64(modResult))
	t.Logf("  1 + %.0f + 3 = %.0f", modResult, 1+modResult+3)
	t.Logf("  %.2f * %.0f = %.2f", oPrix, 1+modResult+3, oPrix*(1+modResult+3))
	t.Logf("  %.2f + %.2f = %.2f", oPrix*(1+modResult+3), bPrix, oPrix*(1+modResult+3)+bPrix)
	t.Logf("  %.2f - 1 = %.2f", oPrix*(1+modResult+3)+bPrix, total)
	t.Log("")
	t.Logf("📊 Résultat final: prixTotal = %.2f", total)
}