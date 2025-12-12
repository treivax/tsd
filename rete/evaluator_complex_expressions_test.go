// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete
import (
	"testing"
)
// TestComplexArithmeticInConstraints teste les expressions arithmétiques complexes dans les contraintes (prémisses)
func TestComplexArithmeticInConstraints(t *testing.T) {
	tests := []struct {
		name           string
		setupFacts     func() map[string]*Fact
		constraint     interface{}
		expectedResult bool
		description    string
	}{
		{
			name: "complex constraint with multiple literals - modulo and operations",
			setupFacts: func() map[string]*Fact {
				return map[string]*Fact{
					"o": {
						ID:   "obj1",
						Type: "Objet",
						Fields: map[string]interface{}{
							"prix": float64(100),
						},
					},
					"b": {
						ID:   "box1",
						Type: "Boite",
						Fields: map[string]interface{}{
							"prix": float64(15),
							"cout": float64(5),
						},
					},
				}
			},
			// b.prix < o.prix % 10 * 1.345 - b.cout + 1 - 3
			// b.prix < (100 % 10) * 1.345 - 5 + 1 - 3
			// 15 < 0 * 1.345 - 5 + 1 - 3
			// 15 < 0 - 5 + 1 - 3
			// 15 < -7 => false
			constraint: map[string]interface{}{
				"type":     "binaryOperation",
				"operator": "<",
				"left": map[string]interface{}{
					"type":   "fieldAccess",
					"object": "b",
					"field":  "prix",
				},
				"right": map[string]interface{}{
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
										"object": "o",
										"field":  "prix",
									},
									"right": map[string]interface{}{
										"type":  "numberLiteral",
										"value": float64(10),
									},
								},
								"right": map[string]interface{}{
									"type":  "numberLiteral",
									"value": float64(1.345),
								},
							},
							"right": map[string]interface{}{
								"type":   "fieldAccess",
								"object": "b",
								"field":  "cout",
							},
						},
						"right": map[string]interface{}{
							"type":  "numberLiteral",
							"value": float64(1),
						},
					},
					"right": map[string]interface{}{
						"type":  "numberLiteral",
						"value": float64(3),
					},
				},
			},
			expectedResult: false,
			description:    "b.prix < o.prix % 10 * 1.345 - b.cout + 1 - 3 => 15 < -7 => false",
		},
		{
			name: "complex constraint with 7+ literals",
			setupFacts: func() map[string]*Fact {
				return map[string]*Fact{
					"d": {
						ID:   "data1",
						Type: "Data",
						Fields: map[string]interface{}{
							"value": float64(50),
						},
					},
				}
			},
			// d.value > 2 + 3 - 4 * 5 / 10 + 1.5 - 0.5
			// 50 > 2 + 3 - 2 + 1.5 - 0.5
			// 50 > 4 => true
			constraint: map[string]interface{}{
				"type":     "binaryOperation",
				"operator": ">",
				"left": map[string]interface{}{
					"type":   "fieldAccess",
					"object": "d",
					"field":  "value",
				},
				"right": map[string]interface{}{
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
									"type":  "numberLiteral",
									"value": float64(2),
								},
								"right": map[string]interface{}{
									"type":  "numberLiteral",
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
										"type":  "numberLiteral",
										"value": float64(4),
									},
									"right": map[string]interface{}{
										"type":  "numberLiteral",
										"value": float64(5),
									},
								},
								"right": map[string]interface{}{
									"type":  "numberLiteral",
									"value": float64(10),
								},
							},
						},
						"right": map[string]interface{}{
							"type":  "numberLiteral",
							"value": float64(1.5),
						},
					},
					"right": map[string]interface{}{
						"type":  "numberLiteral",
						"value": float64(0.5),
					},
				},
			},
			expectedResult: true,
			description:    "d.value > 2 + 3 - 4 * 5 / 10 + 1.5 - 0.5 => 50 > 4 => true",
		},
		{
			name: "modulo in constraint with decimals",
			setupFacts: func() map[string]*Fact {
				return map[string]*Fact{
					"p": {
						ID:   "prod1",
						Type: "Product",
						Fields: map[string]interface{}{
							"quantity": float64(17),
						},
					},
				}
			},
			// p.quantity % 5 == 2
			// 17 % 5 == 2 => true
			constraint: map[string]interface{}{
				"type":     "binaryOperation",
				"operator": "==",
				"left": map[string]interface{}{
					"type":     "binaryOperation",
					"operator": "%",
					"left": map[string]interface{}{
						"type":   "fieldAccess",
						"object": "p",
						"field":  "quantity",
					},
					"right": map[string]interface{}{
						"type":  "numberLiteral",
						"value": float64(5),
					},
				},
				"right": map[string]interface{}{
					"type":  "numberLiteral",
					"value": float64(2),
				},
			},
			expectedResult: true,
			description:    "p.quantity % 5 == 2 => 17 % 5 == 2 => true",
		},
		{
			name: "precedence test - multiplication before addition",
			setupFacts: func() map[string]*Fact {
				return map[string]*Fact{
					"x": {
						ID:   "test1",
						Type: "Test",
						Fields: map[string]interface{}{
							"a": float64(10),
							"b": float64(5),
						},
					},
				}
			},
			// x.a + x.b * 2 == 20
			// 10 + 5 * 2 == 20
			// 10 + 10 == 20 => true
			constraint: map[string]interface{}{
				"type":     "binaryOperation",
				"operator": "==",
				"left": map[string]interface{}{
					"type":     "binaryOperation",
					"operator": "+",
					"left": map[string]interface{}{
						"type":   "fieldAccess",
						"object": "x",
						"field":  "a",
					},
					"right": map[string]interface{}{
						"type":     "binaryOperation",
						"operator": "*",
						"left": map[string]interface{}{
							"type":   "fieldAccess",
							"object": "x",
							"field":  "b",
						},
						"right": map[string]interface{}{
							"type":  "numberLiteral",
							"value": float64(2),
						},
					},
				},
				"right": map[string]interface{}{
					"type":  "numberLiteral",
					"value": float64(20),
				},
			},
			expectedResult: true,
			description:    "x.a + x.b * 2 == 20 => 10 + 10 == 20 => true",
		},
		{
			name: "complex nested expression in constraint",
			setupFacts: func() map[string]*Fact {
				return map[string]*Fact{
					"o": {
						ID:   "obj1",
						Type: "Objet",
						Fields: map[string]interface{}{
							"prix": float64(100),
						},
					},
				}
			},
			// o.prix > 10 * (2 + 3) - 5
			// 100 > 10 * 5 - 5
			// 100 > 50 - 5
			// 100 > 45 => true
			constraint: map[string]interface{}{
				"type":     "binaryOperation",
				"operator": ">",
				"left": map[string]interface{}{
					"type":   "fieldAccess",
					"object": "o",
					"field":  "prix",
				},
				"right": map[string]interface{}{
					"type":     "binaryOperation",
					"operator": "-",
					"left": map[string]interface{}{
						"type":     "binaryOperation",
						"operator": "*",
						"left": map[string]interface{}{
							"type":  "numberLiteral",
							"value": float64(10),
						},
						"right": map[string]interface{}{
							"type":     "binaryOperation",
							"operator": "+",
							"left": map[string]interface{}{
								"type":  "numberLiteral",
								"value": float64(2),
							},
							"right": map[string]interface{}{
								"type":  "numberLiteral",
								"value": float64(3),
							},
						},
					},
					"right": map[string]interface{}{
						"type":  "numberLiteral",
						"value": float64(5),
					},
				},
			},
			expectedResult: true,
			description:    "o.prix > 10 * (2 + 3) - 5 => 100 > 45 => true",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			facts := tt.setupFacts()
			// Créer l'évaluateur et définir les bindings
			evaluator := NewAlphaConditionEvaluator()
			evaluator.variableBindings = facts
			// Évaluer l'expression
			result, err := evaluator.evaluateExpression(tt.constraint)
			if err != nil {
				t.Fatalf("Unexpected error: %v\nConstraint: %s", err, tt.description)
			}
			// Vérifier le résultat
			if result != tt.expectedResult {
				t.Errorf("Constraint: %s\nExpected: %v\nGot: %v",
					tt.description, tt.expectedResult, result)
				// Debug: afficher les valeurs des faits
				t.Log("Fact values:")
				for varName, fact := range facts {
					t.Logf("  %s: %+v", varName, fact.Fields)
				}
			} else {
				t.Logf("✅ %s", tt.description)
			}
		})
	}
}
// TestComplexExpressionConstraintTypes teste différents types d'opérations dans les contraintes
func TestComplexExpressionConstraintTypes(t *testing.T) {
	tests := []struct {
		name        string
		fact        *Fact
		varName     string
		constraint  interface{}
		shouldMatch bool
	}{
		{
			name: "arithmetic in equality check",
			fact: &Fact{
				ID:   "f1",
				Type: "Item",
				Fields: map[string]interface{}{
					"total": float64(150),
					"base":  float64(100),
					"tax":   float64(50),
				},
			},
			varName: "i",
			// total == base + tax (150 == 150)
			constraint: map[string]interface{}{
				"type":     "binaryOperation",
				"operator": "==",
				"left": map[string]interface{}{
					"type":   "fieldAccess",
					"object": "i",
					"field":  "total",
				},
				"right": map[string]interface{}{
					"type":     "binaryOperation",
					"operator": "+",
					"left": map[string]interface{}{
						"type":   "fieldAccess",
						"object": "i",
						"field":  "base",
					},
					"right": map[string]interface{}{
						"type":   "fieldAccess",
						"object": "i",
						"field":  "tax",
					},
				},
			},
			shouldMatch: true,
		},
		{
			name: "division in less than check",
			fact: &Fact{
				ID:   "f2",
				Type: "Ratio",
				Fields: map[string]interface{}{
					"value":   float64(25),
					"divided": float64(100),
					"divisor": float64(4),
				},
			},
			varName: "r",
			// value == divided / divisor (25 == 25)
			constraint: map[string]interface{}{
				"type":     "binaryOperation",
				"operator": "==",
				"left": map[string]interface{}{
					"type":   "fieldAccess",
					"object": "r",
					"field":  "value",
				},
				"right": map[string]interface{}{
					"type":     "binaryOperation",
					"operator": "/",
					"left": map[string]interface{}{
						"type":   "fieldAccess",
						"object": "r",
						"field":  "divided",
					},
					"right": map[string]interface{}{
						"type":   "fieldAccess",
						"object": "r",
						"field":  "divisor",
					},
				},
			},
			shouldMatch: true,
		},
		{
			name: "modulo in constraint with literal",
			fact: &Fact{
				ID:   "f3",
				Type: "Number",
				Fields: map[string]interface{}{
					"value": float64(23),
				},
			},
			varName: "n",
			// value % 10 > 2 (23 % 10 = 3, 3 > 2 = true)
			constraint: map[string]interface{}{
				"type":     "binaryOperation",
				"operator": ">",
				"left": map[string]interface{}{
					"type":     "binaryOperation",
					"operator": "%",
					"left": map[string]interface{}{
						"type":   "fieldAccess",
						"object": "n",
						"field":  "value",
					},
					"right": map[string]interface{}{
						"type":  "numberLiteral",
						"value": float64(10),
					},
				},
				"right": map[string]interface{}{
					"type":  "numberLiteral",
					"value": float64(2),
				},
			},
			shouldMatch: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			facts := map[string]*Fact{tt.varName: tt.fact}
			evaluator := NewAlphaConditionEvaluator()
			evaluator.variableBindings = facts
			result, err := evaluator.evaluateExpression(tt.constraint)
			if err != nil {
				t.Fatalf("Error evaluating constraint: %v", err)
			}
			if result != tt.shouldMatch {
				t.Errorf("Expected %v, got %v", tt.shouldMatch, result)
			} else {
				t.Logf("✅ Test passed: %v", tt.name)
			}
		})
	}
}
// TestRealWorldConstraintExpression teste l'exemple complet de contrainte
func TestRealWorldConstraintExpression(t *testing.T) {
	t.Log("🧮 Test de contrainte complexe du monde réel")
	t.Log("=============================================")
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
	t.Logf("Objet: id=%s, prix=%.2f, boite=%s", objet.Fields["id"], objet.Fields["prix"], objet.Fields["boite"])
	t.Logf("Boite: id=%s, prix=%.2f, cout=%.2f", boite.Fields["id"], boite.Fields["prix"], boite.Fields["cout"])
	t.Log("")
	facts := map[string]*Fact{
		"o": objet,
		"b": boite,
	}
	evaluator := NewAlphaConditionEvaluator()
	evaluator.variableBindings = facts
	// Contrainte : b.prix < o.prix % 10 * 1.345 - b.cout + 1 - 3
	// Calcul : 15 < (100 % 10) * 1.345 - 5 + 1 - 3
	//         15 < 0 * 1.345 - 5 + 1 - 3
	//         15 < -7 => false
	constraint := map[string]interface{}{
		"type":     "binaryOperation",
		"operator": "<",
		"left": map[string]interface{}{
			"type":   "fieldAccess",
			"object": "b",
			"field":  "prix",
		},
		"right": map[string]interface{}{
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
								"object": "o",
								"field":  "prix",
							},
							"right": map[string]interface{}{
								"type":  "numberLiteral",
								"value": float64(10),
							},
						},
						"right": map[string]interface{}{
							"type":  "numberLiteral",
							"value": float64(1.345),
						},
					},
					"right": map[string]interface{}{
						"type":   "fieldAccess",
						"object": "b",
						"field":  "cout",
					},
				},
				"right": map[string]interface{}{
					"type":  "numberLiteral",
					"value": float64(1),
				},
			},
			"right": map[string]interface{}{
				"type":  "numberLiteral",
				"value": float64(3),
			},
		},
	}
	result, err := evaluator.evaluateExpression(constraint)
	if err != nil {
		t.Fatalf("❌ Erreur: %v", err)
	}
	// Calcul manuel du résultat attendu
	oPrix := objet.Fields["prix"].(float64)
	bPrix := boite.Fields["prix"].(float64)
	bCout := boite.Fields["cout"].(float64)
	modResult := float64(int64(oPrix) % int64(10)) // 100 % 10 = 0
	rightSide := modResult*1.345 - bCout + 1 - 3   // 0 - 5 + 1 - 3 = -7
	expected := bPrix < rightSide                  // 15 < -7 = false
	t.Log("Calcul détaillé de la contrainte:")
	t.Logf("  o.prix %% 10 = %.0f", modResult)
	t.Logf("  %.0f * 1.345 = %.3f", modResult, modResult*1.345)
	t.Logf("  %.3f - %.0f = %.3f", modResult*1.345, bCout, modResult*1.345-bCout)
	t.Logf("  %.3f + 1 = %.3f", modResult*1.345-bCout, modResult*1.345-bCout+1)
	t.Logf("  %.3f - 3 = %.3f", modResult*1.345-bCout+1, rightSide)
	t.Log("")
	t.Logf("  b.prix < %.3f", rightSide)
	t.Logf("  %.0f < %.3f = %v", bPrix, rightSide, expected)
	t.Log("")
	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	} else {
		t.Logf("✅ Contrainte évaluée correctement: %v", result)
	}
}