// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete

import (
	"testing"
)

// TestCircularDependency_IntegrationWithDecomposer teste l'intégration avec le décomposeur
func TestCircularDependency_IntegrationWithDecomposer(t *testing.T) {
	decomposer := NewArithmeticExpressionDecomposer()
	detector := NewCircularDependencyDetector()
	// Expression valide: (c.qte * 23) > 100
	condition := map[string]interface{}{
		"type":     "comparison",
		"operator": ">",
		"left": map[string]interface{}{
			"type":     "binaryOp",
			"operator": "*",
			"left": map[string]interface{}{
				"type":     "fieldAccess",
				"field":    "qte",
				"variable": "c",
			},
			"right": map[string]interface{}{
				"type":  "number",
				"value": 23,
			},
		},
		"right": map[string]interface{}{
			"type":  "number",
			"value": 100,
		},
	}
	// Décomposer
	decomposed, err := decomposer.DecomposeToDecomposedConditions(condition)
	if err != nil {
		t.Fatalf("Decomposition failed: %v", err)
	}
	t.Logf("Decomposed into %d steps", len(decomposed))
	for i, step := range decomposed {
		t.Logf("  Step %d: %s (deps: %v)", i+1, step.ResultName, step.Dependencies)
	}
	// Valider avec le détecteur
	err = detector.ValidateDecomposedConditions(decomposed)
	if err != nil {
		t.Errorf("Validation failed: %v", err)
	}
	report := detector.Validate()
	if !report.Valid {
		t.Errorf("Expected valid decomposition, got error: %s", report.ErrorMessage)
	}
	if report.HasCircularDeps {
		t.Error("Expected no circular dependencies")
	}
	t.Logf("Validation report: Valid=%v, MaxDepth=%d, TotalNodes=%d",
		report.Valid, report.MaxDepth, report.TotalNodes)
}

// TestCircularDependency_ComplexExpression teste une expression complexe
func TestCircularDependency_ComplexExpression(t *testing.T) {
	decomposer := NewArithmeticExpressionDecomposer()
	detector := NewCircularDependencyDetector()
	// Expression complexe: (c.qte * 23 - 10 + c.remise * 43) > 0
	condition := map[string]interface{}{
		"type":     "comparison",
		"operator": ">",
		"left": map[string]interface{}{
			"type":     "binaryOp",
			"operator": "+",
			"left": map[string]interface{}{
				"type":     "binaryOp",
				"operator": "-",
				"left": map[string]interface{}{
					"type":     "binaryOp",
					"operator": "*",
					"left": map[string]interface{}{
						"type":     "fieldAccess",
						"field":    "qte",
						"variable": "c",
					},
					"right": map[string]interface{}{
						"type":  "number",
						"value": 23,
					},
				},
				"right": map[string]interface{}{
					"type":  "number",
					"value": 10,
				},
			},
			"right": map[string]interface{}{
				"type":     "binaryOp",
				"operator": "*",
				"left": map[string]interface{}{
					"type":     "fieldAccess",
					"field":    "remise",
					"variable": "c",
				},
				"right": map[string]interface{}{
					"type":  "number",
					"value": 43,
				},
			},
		},
		"right": map[string]interface{}{
			"type":  "number",
			"value": 0,
		},
	}
	// Décomposer
	decomposed, err := decomposer.DecomposeToDecomposedConditions(condition)
	if err != nil {
		t.Fatalf("Decomposition failed: %v", err)
	}
	t.Logf("Complex expression decomposed into %d steps", len(decomposed))
	for i, step := range decomposed {
		t.Logf("  Step %d: %s (deps: %v, atomic: %v)",
			i+1, step.ResultName, step.Dependencies, step.IsAtomic)
	}
	// Valider
	err = detector.ValidateDecomposedConditions(decomposed)
	if err != nil {
		t.Errorf("Validation failed: %v", err)
	}
	report := detector.Validate()
	if !report.Valid {
		t.Errorf("Expected valid decomposition, got error: %s", report.ErrorMessage)
	}
	// Vérifier le tri topologique (ordre d'exécution)
	sorted, err := detector.GetTopologicalSort()
	if err != nil {
		t.Fatalf("Topological sort failed: %v", err)
	}
	t.Logf("Execution order (topological): %v", sorted)
	// Vérifier que le tri est cohérent avec les dépendances
	nodeIndex := make(map[string]int)
	for i, node := range sorted {
		nodeIndex[node] = i
	}
	for _, step := range decomposed {
		for _, dep := range step.Dependencies {
			if nodeIndex[dep] >= nodeIndex[step.ResultName] {
				t.Errorf("Dependency ordering violation: %s depends on %s but appears before it",
					step.ResultName, dep)
			}
		}
	}
}

// TestCircularDependency_MultipleExpressions teste plusieurs expressions
func TestCircularDependency_MultipleExpressions(t *testing.T) {
	decomposer := NewArithmeticExpressionDecomposer()
	detector := NewCircularDependencyDetector()
	expressions := []map[string]interface{}{
		// Expression 1: c.qte * 10
		{
			"type":     "binaryOp",
			"operator": "*",
			"left": map[string]interface{}{
				"type":     "fieldAccess",
				"field":    "qte",
				"variable": "c",
			},
			"right": map[string]interface{}{
				"type":  "number",
				"value": 10,
			},
		},
		// Expression 2: c.price + 5
		{
			"type":     "binaryOp",
			"operator": "+",
			"left": map[string]interface{}{
				"type":     "fieldAccess",
				"field":    "price",
				"variable": "c",
			},
			"right": map[string]interface{}{
				"type":  "number",
				"value": 5,
			},
		},
	}
	allSteps := make([]DecomposedCondition, 0)
	// Décomposer toutes les expressions
	for i, expr := range expressions {
		decomposed, err := decomposer.DecomposeToDecomposedConditions(expr)
		if err != nil {
			t.Fatalf("Decomposition of expression %d failed: %v", i+1, err)
		}
		allSteps = append(allSteps, decomposed...)
	}
	t.Logf("Multiple expressions decomposed into %d total steps", len(allSteps))
	// Valider globalement
	err := detector.ValidateDecomposedConditions(allSteps)
	if err != nil {
		t.Errorf("Global validation failed: %v", err)
	}
	report := detector.Validate()
	if !report.Valid {
		t.Errorf("Expected valid decomposition, got error: %s", report.ErrorMessage)
	}
	stats := detector.GetStatistics()
	t.Logf("Statistics: %+v", stats)
}

// TestCircularDependency_WithAlphaChainBuilder teste l'intégration complète
func TestCircularDependency_WithAlphaChainBuilder(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	// Définir le type
	typeDef := TypeDefinition{
		Type: "typeDefinition",
		Name: "Order",
		Fields: []Field{
			{Name: "qte", Type: "number"},
			{Name: "price", Type: "number"},
		},
	}
	network.Types = append(network.Types, typeDef)
	typeNode := NewTypeNode("Order", typeDef, storage)
	typeNode.SetNetwork(network)
	network.TypeNodes["Order"] = typeNode
	// Expression: (o.qte * o.price) > 100
	condition := map[string]interface{}{
		"type":     "comparison",
		"operator": ">",
		"left": map[string]interface{}{
			"type":     "binaryOp",
			"operator": "*",
			"left": map[string]interface{}{
				"type":     "fieldAccess",
				"field":    "qte",
				"variable": "o",
			},
			"right": map[string]interface{}{
				"type":     "fieldAccess",
				"field":    "price",
				"variable": "o",
			},
		},
		"right": map[string]interface{}{
			"type":  "number",
			"value": 100,
		},
	}
	// Décomposer
	decomposer := NewArithmeticExpressionDecomposer()
	decomposed, err := decomposer.DecomposeToDecomposedConditions(condition)
	if err != nil {
		t.Fatalf("Decomposition failed: %v", err)
	}
	// Valider avec le détecteur AVANT de construire la chaîne
	detector := NewCircularDependencyDetector()
	err = detector.ValidateDecomposedConditions(decomposed)
	if err != nil {
		t.Fatalf("Validation failed before building chain: %v", err)
	}
	report := detector.Validate()
	t.Logf("Pre-build validation: Valid=%v, MaxDepth=%d", report.Valid, report.MaxDepth)
	// Construire la chaîne alpha
	chainBuilder := NewAlphaChainBuilder(network, storage)
	alphaChain, err := chainBuilder.BuildDecomposedChain(decomposed, "o", typeNode, "test_rule")
	if err != nil {
		t.Fatalf("Building chain failed: %v", err)
	}
	t.Logf("Built alpha chain with %d nodes", len(alphaChain.Nodes))
	// Valider la chaîne construite
	detector2 := NewCircularDependencyDetector()
	report2 := detector2.ValidateAlphaChain(alphaChain.Nodes)
	if !report2.Valid {
		t.Errorf("Chain validation failed: %s", report2.ErrorMessage)
	}
	if report2.HasCircularDeps {
		t.Errorf("Circular dependencies detected in built chain: %v", report2.CyclePath)
	}
	t.Logf("Post-build validation: Valid=%v, MaxDepth=%d, TotalNodes=%d",
		report2.Valid, report2.MaxDepth, report2.TotalNodes)
}

// TestCircularDependency_DeepNesting teste une expression profondément imbriquée
func TestCircularDependency_DeepNesting(t *testing.T) {
	decomposer := NewArithmeticExpressionDecomposer()
	detector := NewCircularDependencyDetector()
	// Expression profonde: ((((a + b) * c) - d) / e) > 0
	condition := map[string]interface{}{
		"type":     "comparison",
		"operator": ">",
		"left": map[string]interface{}{
			"type":     "binaryOp",
			"operator": "/",
			"left": map[string]interface{}{
				"type":     "binaryOp",
				"operator": "-",
				"left": map[string]interface{}{
					"type":     "binaryOp",
					"operator": "*",
					"left": map[string]interface{}{
						"type":     "binaryOp",
						"operator": "+",
						"left": map[string]interface{}{
							"type":     "fieldAccess",
							"field":    "a",
							"variable": "x",
						},
						"right": map[string]interface{}{
							"type":     "fieldAccess",
							"field":    "b",
							"variable": "x",
						},
					},
					"right": map[string]interface{}{
						"type":     "fieldAccess",
						"field":    "c",
						"variable": "x",
					},
				},
				"right": map[string]interface{}{
					"type":     "fieldAccess",
					"field":    "d",
					"variable": "x",
				},
			},
			"right": map[string]interface{}{
				"type":     "fieldAccess",
				"field":    "e",
				"variable": "x",
			},
		},
		"right": map[string]interface{}{
			"type":  "number",
			"value": 0,
		},
	}
	// Décomposer
	decomposed, err := decomposer.DecomposeToDecomposedConditions(condition)
	if err != nil {
		t.Fatalf("Decomposition failed: %v", err)
	}
	t.Logf("Deep expression decomposed into %d steps", len(decomposed))
	// Valider
	err = detector.ValidateDecomposedConditions(decomposed)
	if err != nil {
		t.Errorf("Validation failed: %v", err)
	}
	report := detector.Validate()
	if !report.Valid {
		t.Errorf("Expected valid decomposition, got error: %s", report.ErrorMessage)
	}
	t.Logf("Deep nesting validation: Valid=%v, MaxDepth=%d, TotalNodes=%d",
		report.Valid, report.MaxDepth, report.TotalNodes)
	// Vérifier l'avertissement de profondeur
	if report.MaxDepth > 10 {
		if len(report.Warnings) == 0 {
			t.Error("Expected warning for excessive depth")
		}
		t.Logf("Warnings: %v", report.Warnings)
	}
}

// TestCircularDependency_Statistics teste les statistiques du graphe
func TestCircularDependency_Statistics(t *testing.T) {
	decomposer := NewArithmeticExpressionDecomposer()
	detector := NewCircularDependencyDetector()
	// Expression avec branches multiples: (a + b) * (c - d)
	condition := map[string]interface{}{
		"type":     "binaryOp",
		"operator": "*",
		"left": map[string]interface{}{
			"type":     "binaryOp",
			"operator": "+",
			"left": map[string]interface{}{
				"type":     "fieldAccess",
				"field":    "a",
				"variable": "x",
			},
			"right": map[string]interface{}{
				"type":     "fieldAccess",
				"field":    "b",
				"variable": "x",
			},
		},
		"right": map[string]interface{}{
			"type":     "binaryOp",
			"operator": "-",
			"left": map[string]interface{}{
				"type":     "fieldAccess",
				"field":    "c",
				"variable": "x",
			},
			"right": map[string]interface{}{
				"type":     "fieldAccess",
				"field":    "d",
				"variable": "x",
			},
		},
	}
	// Décomposer
	decomposed, err := decomposer.DecomposeToDecomposedConditions(condition)
	if err != nil {
		t.Fatalf("Decomposition failed: %v", err)
	}
	// Valider et collecter statistiques
	err = detector.ValidateDecomposedConditions(decomposed)
	if err != nil {
		t.Errorf("Validation failed: %v", err)
	}
	stats := detector.GetStatistics()
	t.Logf("Graph statistics:")
	t.Logf("  Total nodes: %d", stats["total_nodes"])
	t.Logf("  Total edges: %d", stats["total_edges"])
	t.Logf("  Average out-degree: %.2f", stats["average_outdegree"])
	t.Logf("  Max depth: %d", stats["max_depth"])
	t.Logf("  Has cycles: %v", stats["has_cycles"])
	// Vérifier que les statistiques sont cohérentes
	if stats["has_cycles"].(bool) {
		t.Error("Expected no cycles")
	}
	if stats["total_nodes"].(int) != len(decomposed) {
		t.Errorf("Expected %d nodes, got %d", len(decomposed), stats["total_nodes"])
	}
}

// BenchmarkCircularDependency_Integration benchmark l'intégration complète
func BenchmarkCircularDependency_Integration(b *testing.B) {
	decomposer := NewArithmeticExpressionDecomposer()
	// Expression standard
	condition := map[string]interface{}{
		"type":     "comparison",
		"operator": ">",
		"left": map[string]interface{}{
			"type":     "binaryOp",
			"operator": "*",
			"left": map[string]interface{}{
				"type":     "fieldAccess",
				"field":    "qte",
				"variable": "c",
			},
			"right": map[string]interface{}{
				"type":  "number",
				"value": 23,
			},
		},
		"right": map[string]interface{}{
			"type":  "number",
			"value": 100,
		},
	}
	// Décomposer une fois
	decomposed, _ := decomposer.DecomposeToDecomposedConditions(condition)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		detector := NewCircularDependencyDetector()
		_ = detector.ValidateDecomposedConditions(decomposed)
		_ = detector.Validate()
	}
}
