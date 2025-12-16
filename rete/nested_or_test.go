// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete

import (
	"testing"

	"github.com/treivax/tsd/constraint"
)

// TestAnalyzeNestedOR_Simple vérifie l'analyse d'expressions OR simples
func TestAnalyzeNestedOR_Simple(t *testing.T) {
	t.Log("🧪 TEST ANALYSE OR SIMPLE")
	t.Log("=========================")
	// Expression simple: p.age > 18
	simpleExpr := constraint.BinaryOperation{
		Type: "binaryOperation",
		Left: constraint.FieldAccess{
			Type:   "fieldAccess",
			Object: "p",
			Field:  "age",
		},
		Operator: ">",
		Right: constraint.NumberLiteral{
			Type:  "numberLiteral",
			Value: 18.0,
		},
	}
	analysis, err := AnalyzeNestedOR(simpleExpr)
	if err != nil {
		t.Fatalf("❌ Erreur analyse: %v", err)
	}
	if analysis.Complexity != ComplexitySimple {
		t.Errorf("❌ Complexité incorrecte: attendu %v, obtenu %v",
			ComplexitySimple, analysis.Complexity)
	}
	if analysis.NestingDepth != 0 {
		t.Errorf("❌ Profondeur incorrecte: attendu 0, obtenu %d", analysis.NestingDepth)
	}
	if analysis.RequiresDNF || analysis.RequiresFlattening {
		t.Error("❌ Expression simple ne devrait pas nécessiter de transformation")
	}
	t.Log("✅ Analyse OR simple correcte")
}

// TestAnalyzeNestedOR_Flat vérifie l'analyse d'expressions OR plates
func TestAnalyzeNestedOR_Flat(t *testing.T) {
	t.Log("🧪 TEST ANALYSE OR PLAT")
	t.Log("=======================")
	// Expression: p.age > 18 OR p.status == "VIP" OR p.country == "FR"
	flatExpr := constraint.LogicalExpression{
		Type: "logicalExpr",
		Left: constraint.BinaryOperation{
			Type: "binaryOperation",
			Left: constraint.FieldAccess{
				Type:   "fieldAccess",
				Object: "p",
				Field:  "age",
			},
			Operator: ">",
			Right: constraint.NumberLiteral{
				Type:  "numberLiteral",
				Value: 18.0,
			},
		},
		Operations: []constraint.LogicalOperation{
			{
				Op: "OR",
				Right: constraint.BinaryOperation{
					Type: "binaryOperation",
					Left: constraint.FieldAccess{
						Type:   "fieldAccess",
						Object: "p",
						Field:  "status",
					},
					Operator: "==",
					Right: constraint.StringLiteral{
						Type:  "stringLiteral",
						Value: "VIP",
					},
				},
			},
			{
				Op: "OR",
				Right: constraint.BinaryOperation{
					Type: "binaryOperation",
					Left: constraint.FieldAccess{
						Type:   "fieldAccess",
						Object: "p",
						Field:  "country",
					},
					Operator: "==",
					Right: constraint.StringLiteral{
						Type:  "stringLiteral",
						Value: "FR",
					},
				},
			},
		},
	}
	analysis, err := AnalyzeNestedOR(flatExpr)
	if err != nil {
		t.Fatalf("❌ Erreur analyse: %v", err)
	}
	if analysis.Complexity != ComplexityFlat {
		t.Errorf("❌ Complexité incorrecte: attendu %v, obtenu %v",
			ComplexityFlat, analysis.Complexity)
	}
	if analysis.ORTermCount != 2 {
		t.Errorf("❌ Comptage OR incorrect: attendu 2, obtenu %d", analysis.ORTermCount)
	}
	if analysis.RequiresDNF {
		t.Error("❌ Expression OR plate ne devrait pas nécessiter DNF")
	}
	t.Log("✅ Analyse OR plat correcte")
	t.Logf("  Termes OR: %d", analysis.ORTermCount)
}

// TestAnalyzeNestedOR_Nested vérifie l'analyse d'expressions OR imbriquées
func TestAnalyzeNestedOR_Nested(t *testing.T) {
	t.Log("🧪 TEST ANALYSE OR IMBRIQUÉ")
	t.Log("===========================")
	// Expression: p.age > 18 OR (p.status == "VIP" OR p.country == "FR")
	nestedExpr := constraint.LogicalExpression{
		Type: "logicalExpr",
		Left: constraint.BinaryOperation{
			Type: "binaryOperation",
			Left: constraint.FieldAccess{
				Type:   "fieldAccess",
				Object: "p",
				Field:  "age",
			},
			Operator: ">",
			Right: constraint.NumberLiteral{
				Type:  "numberLiteral",
				Value: 18.0,
			},
		},
		Operations: []constraint.LogicalOperation{
			{
				Op: "OR",
				Right: constraint.LogicalExpression{
					Type: "logicalExpr",
					Left: constraint.BinaryOperation{
						Type: "binaryOperation",
						Left: constraint.FieldAccess{
							Type:   "fieldAccess",
							Object: "p",
							Field:  "status",
						},
						Operator: "==",
						Right: constraint.StringLiteral{
							Type:  "stringLiteral",
							Value: "VIP",
						},
					},
					Operations: []constraint.LogicalOperation{
						{
							Op: "OR",
							Right: constraint.BinaryOperation{
								Type: "binaryOperation",
								Left: constraint.FieldAccess{
									Type:   "fieldAccess",
									Object: "p",
									Field:  "country",
								},
								Operator: "==",
								Right: constraint.StringLiteral{
									Type:  "stringLiteral",
									Value: "FR",
								},
							},
						},
					},
				},
			},
		},
	}
	analysis, err := AnalyzeNestedOR(nestedExpr)
	if err != nil {
		t.Fatalf("❌ Erreur analyse: %v", err)
	}
	if analysis.Complexity != ComplexityNestedOR {
		t.Errorf("❌ Complexité incorrecte: attendu %v, obtenu %v",
			ComplexityNestedOR, analysis.Complexity)
	}
	if analysis.NestingDepth < 1 {
		t.Errorf("❌ Profondeur incorrecte: attendu >= 1, obtenu %d", analysis.NestingDepth)
	}
	if !analysis.RequiresFlattening {
		t.Error("❌ Expression imbriquée devrait nécessiter aplatissement")
	}
	t.Log("✅ Analyse OR imbriqué correcte")
	t.Logf("  Profondeur: %d", analysis.NestingDepth)
	t.Logf("  Termes OR: %d", analysis.ORTermCount)
}

// TestAnalyzeNestedOR_MixedANDOR vérifie l'analyse d'expressions mixtes
func TestAnalyzeNestedOR_MixedANDOR(t *testing.T) {
	t.Log("🧪 TEST ANALYSE MIXTE AND/OR")
	t.Log("============================")
	// Expression: (p.age > 18 OR p.status == "VIP") AND p.country == "FR"
	mixedExpr := constraint.LogicalExpression{
		Type: "logicalExpr",
		Left: constraint.BinaryOperation{
			Type: "binaryOperation",
			Left: constraint.FieldAccess{
				Type:   "fieldAccess",
				Object: "p",
				Field:  "age",
			},
			Operator: ">",
			Right: constraint.NumberLiteral{
				Type:  "numberLiteral",
				Value: 18.0,
			},
		},
		Operations: []constraint.LogicalOperation{
			{
				Op: "OR",
				Right: constraint.BinaryOperation{
					Type: "binaryOperation",
					Left: constraint.FieldAccess{
						Type:   "fieldAccess",
						Object: "p",
						Field:  "status",
					},
					Operator: "==",
					Right: constraint.StringLiteral{
						Type:  "stringLiteral",
						Value: "VIP",
					},
				},
			},
			{
				Op: "AND",
				Right: constraint.BinaryOperation{
					Type: "binaryOperation",
					Left: constraint.FieldAccess{
						Type:   "fieldAccess",
						Object: "p",
						Field:  "country",
					},
					Operator: "==",
					Right: constraint.StringLiteral{
						Type:  "stringLiteral",
						Value: "FR",
					},
				},
			},
		},
	}
	analysis, err := AnalyzeNestedOR(mixedExpr)
	if err != nil {
		t.Fatalf("❌ Erreur analyse: %v", err)
	}
	if analysis.Complexity < ComplexityMixedANDOR {
		t.Errorf("❌ Complexité incorrecte: attendu >= %v, obtenu %v",
			ComplexityMixedANDOR, analysis.Complexity)
	}
	if analysis.ORTermCount < 1 || analysis.ANDTermCount < 1 {
		t.Errorf("❌ Comptage incorrect: OR=%d, AND=%d",
			analysis.ORTermCount, analysis.ANDTermCount)
	}
	t.Log("✅ Analyse mixte AND/OR correcte")
	t.Logf("  Complexité: %v", analysis.Complexity)
	t.Logf("  Termes OR: %d, AND: %d", analysis.ORTermCount, analysis.ANDTermCount)
}

// TestAnalyzeNestedOR_DNFCandidate vérifie la détection de candidats DNF
func TestAnalyzeNestedOR_DNFCandidate(t *testing.T) {
	t.Log("🧪 TEST DÉTECTION CANDIDAT DNF")
	t.Log("==============================")
	// Expression: (p.age > 18 OR p.age < 10) AND (p.status == "VIP" OR p.status == "PREMIUM")
	// Candidat idéal pour DNF
	dnfCandidateExpr := constraint.LogicalExpression{
		Type: "logicalExpr",
		Left: constraint.BinaryOperation{
			Type: "binaryOperation",
			Left: constraint.FieldAccess{
				Type:   "fieldAccess",
				Object: "p",
				Field:  "age",
			},
			Operator: ">",
			Right: constraint.NumberLiteral{
				Type:  "numberLiteral",
				Value: 18.0,
			},
		},
		Operations: []constraint.LogicalOperation{
			{
				Op: "OR",
				Right: constraint.BinaryOperation{
					Type: "binaryOperation",
					Left: constraint.FieldAccess{
						Type:   "fieldAccess",
						Object: "p",
						Field:  "age",
					},
					Operator: "<",
					Right: constraint.NumberLiteral{
						Type:  "numberLiteral",
						Value: 10.0,
					},
				},
			},
			{
				Op: "AND",
				Right: constraint.LogicalExpression{
					Type: "logicalExpr",
					Left: constraint.BinaryOperation{
						Type: "binaryOperation",
						Left: constraint.FieldAccess{
							Type:   "fieldAccess",
							Object: "p",
							Field:  "status",
						},
						Operator: "==",
						Right: constraint.StringLiteral{
							Type:  "stringLiteral",
							Value: "VIP",
						},
					},
					Operations: []constraint.LogicalOperation{
						{
							Op: "OR",
							Right: constraint.BinaryOperation{
								Type: "binaryOperation",
								Left: constraint.FieldAccess{
									Type:   "fieldAccess",
									Object: "p",
									Field:  "status",
								},
								Operator: "==",
								Right: constraint.StringLiteral{
									Type:  "stringLiteral",
									Value: "PREMIUM",
								},
							},
						},
					},
				},
			},
		},
	}
	analysis, err := AnalyzeNestedOR(dnfCandidateExpr)
	if err != nil {
		t.Fatalf("❌ Erreur analyse: %v", err)
	}
	if analysis.Complexity != ComplexityDNFCandidate {
		t.Errorf("❌ Complexité incorrecte: attendu %v, obtenu %v",
			ComplexityDNFCandidate, analysis.Complexity)
	}
	if !analysis.RequiresDNF {
		t.Error("❌ Expression devrait être détectée comme candidat DNF")
	}
	if analysis.OptimizationHint == "" {
		t.Error("❌ Hint d'optimisation devrait être fourni")
	}
	t.Log("✅ Détection candidat DNF correcte")
	t.Logf("  Hint: %s", analysis.OptimizationHint)
}

// TestFlattenNestedOR_Simple vérifie l'aplatissement d'OR imbriqués simples
func TestFlattenNestedOR_Simple(t *testing.T) {
	t.Log("🧪 TEST APLATISSEMENT OR SIMPLE")
	t.Log("===============================")
	// Expression: A OR (B OR C) -> doit devenir A OR B OR C
	nestedExpr := constraint.LogicalExpression{
		Type: "logicalExpr",
		Left: constraint.BinaryOperation{
			Type:     "binaryOperation",
			Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "a"},
			Operator: "==",
			Right:    constraint.StringLiteral{Type: "stringLiteral", Value: "A"},
		},
		Operations: []constraint.LogicalOperation{
			{
				Op: "OR",
				Right: constraint.LogicalExpression{
					Type: "logicalExpr",
					Left: constraint.BinaryOperation{
						Type:     "binaryOperation",
						Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "b"},
						Operator: "==",
						Right:    constraint.StringLiteral{Type: "stringLiteral", Value: "B"},
					},
					Operations: []constraint.LogicalOperation{
						{
							Op: "OR",
							Right: constraint.BinaryOperation{
								Type:     "binaryOperation",
								Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "c"},
								Operator: "==",
								Right:    constraint.StringLiteral{Type: "stringLiteral", Value: "C"},
							},
						},
					},
				},
			},
		},
	}
	flattened, err := FlattenNestedOR(nestedExpr)
	if err != nil {
		t.Fatalf("❌ Erreur aplatissement: %v", err)
	}
	flattenedExpr, ok := flattened.(constraint.LogicalExpression)
	if !ok {
		t.Fatalf("❌ Type incorrect après aplatissement: %T", flattened)
	}
	// Vérifier que nous avons 3 termes au total (A, B, C)
	totalTerms := 1 + len(flattenedExpr.Operations)
	if totalTerms != 3 {
		t.Errorf("❌ Nombre de termes incorrect: attendu 3, obtenu %d", totalTerms)
	}
	// Vérifier que tous les opérateurs sont OR
	for _, op := range flattenedExpr.Operations {
		if op.Op != "OR" {
			t.Errorf("❌ Opérateur incorrect après aplatissement: %s", op.Op)
		}
	}
	t.Log("✅ Aplatissement OR simple réussi")
	t.Logf("  Termes après aplatissement: %d", totalTerms)
}

// TestFlattenNestedOR_Deep vérifie l'aplatissement d'OR profondément imbriqués
func TestFlattenNestedOR_Deep(t *testing.T) {
	t.Log("🧪 TEST APLATISSEMENT OR PROFOND")
	t.Log("================================")
	// Expression: A OR (B OR (C OR D)) -> A OR B OR C OR D
	deepNestedExpr := constraint.LogicalExpression{
		Type: "logicalExpr",
		Left: constraint.BinaryOperation{
			Type:     "binaryOperation",
			Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "a"},
			Operator: "==",
			Right:    constraint.StringLiteral{Type: "stringLiteral", Value: "A"},
		},
		Operations: []constraint.LogicalOperation{
			{
				Op: "OR",
				Right: constraint.LogicalExpression{
					Type: "logicalExpr",
					Left: constraint.BinaryOperation{
						Type:     "binaryOperation",
						Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "b"},
						Operator: "==",
						Right:    constraint.StringLiteral{Type: "stringLiteral", Value: "B"},
					},
					Operations: []constraint.LogicalOperation{
						{
							Op: "OR",
							Right: constraint.LogicalExpression{
								Type: "logicalExpr",
								Left: constraint.BinaryOperation{
									Type:     "binaryOperation",
									Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "c"},
									Operator: "==",
									Right:    constraint.StringLiteral{Type: "stringLiteral", Value: "C"},
								},
								Operations: []constraint.LogicalOperation{
									{
										Op: "OR",
										Right: constraint.BinaryOperation{
											Type:     "binaryOperation",
											Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "d"},
											Operator: "==",
											Right:    constraint.StringLiteral{Type: "stringLiteral", Value: "D"},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	flattened, err := FlattenNestedOR(deepNestedExpr)
	if err != nil {
		t.Fatalf("❌ Erreur aplatissement: %v", err)
	}
	flattenedExpr, ok := flattened.(constraint.LogicalExpression)
	if !ok {
		t.Fatalf("❌ Type incorrect: %T", flattened)
	}
	totalTerms := 1 + len(flattenedExpr.Operations)
	if totalTerms != 4 {
		t.Errorf("❌ Nombre de termes incorrect: attendu 4, obtenu %d", totalTerms)
	}
	t.Log("✅ Aplatissement OR profond réussi")
	t.Logf("  Termes après aplatissement: %d", totalTerms)
}

// TestNormalizeNestedOR_Complete vérifie la normalisation complète
func TestNormalizeNestedOR_Complete(t *testing.T) {
	t.Log("🧪 TEST NORMALISATION COMPLÈTE")
	t.Log("==============================")
	// Expression: C OR (A OR B) -> doit devenir A OR B OR C (après aplatissement et tri)
	expr := constraint.LogicalExpression{
		Type: "logicalExpr",
		Left: constraint.BinaryOperation{
			Type:     "binaryOperation",
			Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "name"},
			Operator: "==",
			Right:    constraint.StringLiteral{Type: "stringLiteral", Value: "C"},
		},
		Operations: []constraint.LogicalOperation{
			{
				Op: "OR",
				Right: constraint.LogicalExpression{
					Type: "logicalExpr",
					Left: constraint.BinaryOperation{
						Type:     "binaryOperation",
						Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "name"},
						Operator: "==",
						Right:    constraint.StringLiteral{Type: "stringLiteral", Value: "A"},
					},
					Operations: []constraint.LogicalOperation{
						{
							Op: "OR",
							Right: constraint.BinaryOperation{
								Type:     "binaryOperation",
								Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "name"},
								Operator: "==",
								Right:    constraint.StringLiteral{Type: "stringLiteral", Value: "B"},
							},
						},
					},
				},
			},
		},
	}
	normalized, err := NormalizeNestedOR(expr)
	if err != nil {
		t.Fatalf("❌ Erreur normalisation: %v", err)
	}
	normalizedExpr, ok := normalized.(constraint.LogicalExpression)
	if !ok {
		t.Fatalf("❌ Type incorrect: %T", normalized)
	}
	// Vérifier que nous avons 3 termes
	totalTerms := 1 + len(normalizedExpr.Operations)
	if totalTerms != 3 {
		t.Errorf("❌ Nombre de termes incorrect: attendu 3, obtenu %d", totalTerms)
	}
	t.Log("✅ Normalisation complète réussie")
	t.Logf("  Termes normalisés: %d", totalTerms)
}

// TestNormalizeNestedOR_OrderIndependent vérifie l'indépendance d'ordre
func TestNormalizeNestedOR_OrderIndependent(t *testing.T) {
	t.Log("🧪 TEST INDÉPENDANCE D'ORDRE")
	t.Log("============================")
	// Deux expressions équivalentes avec ordre différent
	expr1 := constraint.LogicalExpression{
		Type: "logicalExpr",
		Left: constraint.BinaryOperation{
			Type:     "binaryOperation",
			Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "x"},
			Operator: "==",
			Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 1.0},
		},
		Operations: []constraint.LogicalOperation{
			{
				Op: "OR",
				Right: constraint.LogicalExpression{
					Type: "logicalExpr",
					Left: constraint.BinaryOperation{
						Type:     "binaryOperation",
						Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "x"},
						Operator: "==",
						Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 2.0},
					},
					Operations: []constraint.LogicalOperation{
						{
							Op: "OR",
							Right: constraint.BinaryOperation{
								Type:     "binaryOperation",
								Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "x"},
								Operator: "==",
								Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 3.0},
							},
						},
					},
				},
			},
		},
	}
	expr2 := constraint.LogicalExpression{
		Type: "logicalExpr",
		Left: constraint.BinaryOperation{
			Type:     "binaryOperation",
			Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "x"},
			Operator: "==",
			Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 3.0},
		},
		Operations: []constraint.LogicalOperation{
			{
				Op: "OR",
				Right: constraint.LogicalExpression{
					Type: "logicalExpr",
					Left: constraint.BinaryOperation{
						Type:     "binaryOperation",
						Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "x"},
						Operator: "==",
						Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 1.0},
					},
					Operations: []constraint.LogicalOperation{
						{
							Op: "OR",
							Right: constraint.BinaryOperation{
								Type:     "binaryOperation",
								Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "x"},
								Operator: "==",
								Right:    constraint.NumberLiteral{Type: "numberLiteral", Value: 2.0},
							},
						},
					},
				},
			},
		},
	}
	normalized1, err := NormalizeNestedOR(expr1)
	if err != nil {
		t.Fatalf("❌ Erreur normalisation expr1: %v", err)
	}
	normalized2, err := NormalizeNestedOR(expr2)
	if err != nil {
		t.Fatalf("❌ Erreur normalisation expr2: %v", err)
	}
	// Calculer les hashes pour comparaison
	condition1 := map[string]interface{}{
		"type":       "constraint",
		"constraint": normalized1,
	}
	condition2 := map[string]interface{}{
		"type":       "constraint",
		"constraint": normalized2,
	}
	hash1, err := ConditionHash(condition1, "p")
	if err != nil {
		t.Fatalf("❌ Erreur hash1: %v", err)
	}
	hash2, err := ConditionHash(condition2, "p")
	if err != nil {
		t.Fatalf("❌ Erreur hash2: %v", err)
	}
	if hash1 != hash2 {
		t.Errorf("❌ Hashes différents après normalisation:\n  Hash1: %s\n  Hash2: %s",
			hash1, hash2)
	}
	t.Log("✅ Normalisation indépendante de l'ordre")
	t.Logf("  Hash commun: %s", hash1)
}

// TestIntegration_NestedOR_SingleAlphaNode vérifie l'intégration complète
func TestIntegration_NestedOR_SingleAlphaNode(t *testing.T) {
	t.Log("🧪 TEST INTÉGRATION OR IMBRIQUÉ")
	t.Log("===============================")
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	// Créer TypeNode
	typeDef := TypeDefinition{
		Type: "typeDefinition",
		Name: "Person",
		Fields: []Field{
			{Name: "name", Type: "string"},
			{Name: "age", Type: "number"},
		},
	}
	typeNode := NewTypeNode("Person", typeDef, storage)
	network.TypeNodes["Person"] = typeNode
	// Expression imbriquée: p.name == "A" OR (p.name == "B" OR p.name == "C")
	nestedExpr := constraint.LogicalExpression{
		Type: "logicalExpr",
		Left: constraint.BinaryOperation{
			Type:     "binaryOperation",
			Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "name"},
			Operator: "==",
			Right:    constraint.StringLiteral{Type: "stringLiteral", Value: "A"},
		},
		Operations: []constraint.LogicalOperation{
			{
				Op: "OR",
				Right: constraint.LogicalExpression{
					Type: "logicalExpr",
					Left: constraint.BinaryOperation{
						Type:     "binaryOperation",
						Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "name"},
						Operator: "==",
						Right:    constraint.StringLiteral{Type: "stringLiteral", Value: "B"},
					},
					Operations: []constraint.LogicalOperation{
						{
							Op: "OR",
							Right: constraint.BinaryOperation{
								Type:     "binaryOperation",
								Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "name"},
								Operator: "==",
								Right:    constraint.StringLiteral{Type: "stringLiteral", Value: "C"},
							},
						},
					},
				},
			},
		},
	}
	action := &Action{
		Type: "log",
		Job: &JobCall{
			Name: "log",
			Args: []interface{}{"Match detected"},
		},
	}
	cp := &ConstraintPipeline{}
	err := cp.createAlphaNodeWithTerminal(network, "rule_nested", nestedExpr, "p", "Person", action, storage)
	if err != nil {
		t.Fatalf("❌ Erreur création AlphaNode: %v", err)
	}
	// Vérifier qu'un seul AlphaNode est créé
	if len(network.AlphaNodes) != 1 {
		t.Errorf("❌ Attendu 1 AlphaNode, obtenu %d", len(network.AlphaNodes))
	}
	// Vérifier qu'un TerminalNode est créé
	if len(network.TerminalNodes) != 1 {
		t.Errorf("❌ Attendu 1 TerminalNode, obtenu %d", len(network.TerminalNodes))
	}
	t.Log("✅ Intégration OR imbriqué réussie")
	t.Logf("  AlphaNodes: %d", len(network.AlphaNodes))
	t.Logf("  TerminalNodes: %d", len(network.TerminalNodes))
}

// TestIntegration_NestedOR_Sharing vérifie le partage entre règles imbriquées
func TestIntegration_NestedOR_Sharing(t *testing.T) {
	t.Log("🧪 TEST PARTAGE OR IMBRIQUÉ")
	t.Log("===========================")
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	typeDef := TypeDefinition{
		Type: "typeDefinition",
		Name: "Person",
		Fields: []Field{
			{Name: "status", Type: "string"},
		},
	}
	typeNode := NewTypeNode("Person", typeDef, storage)
	network.TypeNodes["Person"] = typeNode
	// Règle 1: A OR (B OR C)
	rule1Expr := constraint.LogicalExpression{
		Type: "logicalExpr",
		Left: constraint.BinaryOperation{
			Type:     "binaryOperation",
			Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "status"},
			Operator: "==",
			Right:    constraint.StringLiteral{Type: "stringLiteral", Value: "A"},
		},
		Operations: []constraint.LogicalOperation{
			{
				Op: "OR",
				Right: constraint.LogicalExpression{
					Type: "logicalExpr",
					Left: constraint.BinaryOperation{
						Type:     "binaryOperation",
						Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "status"},
						Operator: "==",
						Right:    constraint.StringLiteral{Type: "stringLiteral", Value: "B"},
					},
					Operations: []constraint.LogicalOperation{
						{
							Op: "OR",
							Right: constraint.BinaryOperation{
								Type:     "binaryOperation",
								Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "status"},
								Operator: "==",
								Right:    constraint.StringLiteral{Type: "stringLiteral", Value: "C"},
							},
						},
					},
				},
			},
		},
	}
	// Règle 2: (C OR B) OR A (ordre différent mais équivalent après normalisation)
	rule2Expr := constraint.LogicalExpression{
		Type: "logicalExpr",
		Left: constraint.LogicalExpression{
			Type: "logicalExpr",
			Left: constraint.BinaryOperation{
				Type:     "binaryOperation",
				Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "status"},
				Operator: "==",
				Right:    constraint.StringLiteral{Type: "stringLiteral", Value: "C"},
			},
			Operations: []constraint.LogicalOperation{
				{
					Op: "OR",
					Right: constraint.BinaryOperation{
						Type:     "binaryOperation",
						Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "status"},
						Operator: "==",
						Right:    constraint.StringLiteral{Type: "stringLiteral", Value: "B"},
					},
				},
			},
		},
		Operations: []constraint.LogicalOperation{
			{
				Op: "OR",
				Right: constraint.BinaryOperation{
					Type:     "binaryOperation",
					Left:     constraint.FieldAccess{Type: "fieldAccess", Object: "p", Field: "status"},
					Operator: "==",
					Right:    constraint.StringLiteral{Type: "stringLiteral", Value: "A"},
				},
			},
		},
	}
	action1 := &Action{Type: "log", Job: &JobCall{Name: "log", Args: []interface{}{"Rule1"}}}
	action2 := &Action{Type: "log", Job: &JobCall{Name: "log", Args: []interface{}{"Rule2"}}}
	cp := &ConstraintPipeline{}
	err := cp.createAlphaNodeWithTerminal(network, "rule1", rule1Expr, "p", "Person", action1, storage)
	if err != nil {
		t.Fatalf("❌ Erreur règle 1: %v", err)
	}
	err = cp.createAlphaNodeWithTerminal(network, "rule2", rule2Expr, "p", "Person", action2, storage)
	if err != nil {
		t.Fatalf("❌ Erreur règle 2: %v", err)
	}
	// Vérifier le partage: 1 AlphaNode partagé
	if len(network.AlphaNodes) != 1 {
		t.Errorf("❌ Partage échoué: attendu 1 AlphaNode, obtenu %d", len(network.AlphaNodes))
	}
	// Vérifier: 2 TerminalNodes
	if len(network.TerminalNodes) != 2 {
		t.Errorf("❌ Attendu 2 TerminalNodes, obtenu %d", len(network.TerminalNodes))
	}
	t.Log("✅ Partage OR imbriqué réussi")
	t.Logf("  1 AlphaNode -> 2 TerminalNodes")
}
