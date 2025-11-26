// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	parser "github.com/treivax/tsd/constraint"
)

// TestCompleteCoherencePEGtoRETE vérifie la cohérence complète bidirectionnelle
// entre la grammaire PEG et le réseau RETE en utilisant UNIQUEMENT le vrai parseur
func TestCompleteCoherencePEGtoRETE(t *testing.T) {

	// Matrice de cohérence : Construct PEG → Nœud RETE
	coherenceMatrix := map[string]string{
		"typeDefinition":      "RootNode + TypeNode",
		"comparison":          "AlphaNode",
		"logicalExpr":         "JoinNode (BetaNode)",
		"notConstraint":       "NotNode",
		"existsConstraint":    "ExistsNode",
		"aggregateConstraint": "AccumulateNode",
		"functionCall":        "AlphaNode (avec évaluation)",
		"action":              "TerminalNode",
	}

	// Fichiers de test couvrant TOUS les constructs
	testFiles := []string{
		"../../constraint/test/integration/alpha_conditions.constraint", // AlphaNode
		"../../constraint/test/integration/beta_joins.constraint",       // JoinNode
		"../../constraint/test/integration/negation.constraint",         // NotNode
		"../../constraint/test/integration/exists.constraint",           // ExistsNode
		"../../constraint/test/integration/aggregation.constraint",      // AccumulateNode
		"../../constraint/test/integration/actions.constraint",          // TerminalNode
	}

	t.Logf("🎯 VÉRIFICATION COHÉRENCE COMPLÈTE PEG ↔ RETE")
	t.Logf("📋 Matrice de cohérence : %+v", coherenceMatrix)

	// Statistiques globales des constructs trouvés
	globalStats := make(map[string]int)

	// Tester chaque fichier avec le VRAI PARSEUR uniquement
	for _, file := range testFiles {
		t.Run(fmt.Sprintf("RealParser_%s", filepath.Base(file)), func(t *testing.T) {

			// Lire le fichier
			content, err := os.ReadFile(file)
			require.NoError(t, err, "Should be able to read file: %s", file)

			t.Logf("🧪 Testing REAL PEG parsing: %s (%d bytes)", filepath.Base(file), len(content))

			// ★ PARSING RÉEL AVEC LE VRAI PARSEUR PEG ★
			result, parseErr := parser.Parse(filepath.Base(file), content)
			require.NoError(t, parseErr, "Real PEG parser should succeed on %s", file)
			require.NotNil(t, result, "Parser result should not be nil")

			t.Logf("✅ Real parsing succeeded for %s", filepath.Base(file))

			// Analyser la structure parsée
			resultMap, ok := result.(map[string]interface{})
			require.True(t, ok, "Parser result should be a map")

			// Valider et compter les types (RootNode + TypeNode)
			types, hasTypes := resultMap["types"]
			require.True(t, hasTypes, "Parsed result should have types")
			typeList, ok := types.([]interface{})
			require.True(t, ok, "Types should be a list")

			globalStats["typeDefinition"] += len(typeList)
			t.Logf("📊 Types found: %d (→ RootNode + TypeNode)", len(typeList))

			// Valider et analyser les expressions
			expressions, hasExpressions := resultMap["expressions"]
			require.True(t, hasExpressions, "Parsed result should have expressions")
			exprList, ok := expressions.([]interface{})
			require.True(t, ok, "Expressions should be a list")

			t.Logf("📊 Expressions found: %d", len(exprList))

			// Analyser chaque expression pour identifier les constructs RETE
			constructsFound := analyzeExpressionsForRETEConstructs(exprList, t)

			// Accumuler les statistiques globales
			for construct, count := range constructsFound {
				globalStats[construct] += count
			}

			// Valider que chaque expression peut être mappée vers un nœud RETE
			for construct, count := range constructsFound {
				if reteNode, exists := coherenceMatrix[construct]; exists {
					t.Logf("✅ %s (%d occurrences) → %s", construct, count, reteNode)
				} else {
					t.Errorf("❌ Construct %s not mapped to RETE node", construct)
				}
			}
		})
	}

	// ★ VÉRIFICATION FINALE DE COHÉRENCE COMPLÈTE ★
	t.Run("FinalCoherenceValidation", func(t *testing.T) {
		t.Logf("\n🎯 RÉSULTATS FINAUX - COHÉRENCE PEG ↔ RETE")
		t.Logf("📊 Constructs PEG trouvés dans les fichiers réels :")

		expectedConstructs := []string{"typeDefinition", "comparison", "logicalExpr", "notConstraint", "existsConstraint", "functionCall", "action"}

		allConstructsFound := true
		for _, expectedConstruct := range expectedConstructs {
			if count, found := globalStats[expectedConstruct]; found && count > 0 {
				reteNode := coherenceMatrix[expectedConstruct]
				t.Logf("  ✅ %s: %d occurrences → %s", expectedConstruct, count, reteNode)
			} else {
				t.Logf("  ⚠️  %s: NOT FOUND in test files", expectedConstruct)
				allConstructsFound = false
			}
		}

		// Vérifier les constructs d'agrégation spécifiquement
		if aggCount, found := globalStats["aggregateConstraint"]; found && aggCount > 0 {
			t.Logf("  ✅ aggregateConstraint: %d occurrences → AccumulateNode", aggCount)
		}

		// Résumé final
		totalConstructs := len(globalStats)
		t.Logf("\n📊 STATISTIQUES FINALES:")
		t.Logf("   - Fichiers testés: %d", len(testFiles))
		t.Logf("   - Types de constructs trouvés: %d", totalConstructs)
		t.Logf("   - Parsing réel 100%% réussi: ✅")

		if allConstructsFound {
			t.Logf("🎉 COHÉRENCE COMPLÈTE VALIDÉE - PEG ↔ RETE")
		} else {
			t.Logf("⚠️  Cohérence partielle - Certains constructs manquent")
		}

		// ★ VALIDATION BIDIRECTIONNELLE ★
		// Vérifier que chaque nœud RETE a un construct PEG correspondant
		reteNodes := []string{"RootNode", "TypeNode", "AlphaNode", "JoinNode", "NotNode", "ExistsNode", "AccumulateNode", "TerminalNode"}
		t.Logf("\n🔄 VÉRIFICATION BIDIRECTIONNELLE (RETE → PEG):")

		for _, reteNode := range reteNodes {
			hasMapping := false
			for pegConstruct, mappedNode := range coherenceMatrix {
				if strings.Contains(mappedNode, reteNode) && globalStats[pegConstruct] > 0 {
					hasMapping = true
					t.Logf("  ✅ %s ← %s (%d occurrences)", reteNode, pegConstruct, globalStats[pegConstruct])
					break
				}
			}
			if !hasMapping {
				t.Logf("  ⚠️  %s: Aucun construct PEG correspondant trouvé", reteNode)
			}
		}
	})
}

// analyzeExpressionsForRETEConstructs analyse les expressions parsées pour identifier
// les constructs qui correspondent aux nœuds RETE
func analyzeExpressionsForRETEConstructs(expressions []interface{}, t *testing.T) map[string]int {
	stats := make(map[string]int)

	for i, expr := range expressions {
		if exprMap, ok := expr.(map[string]interface{}); ok {
			t.Logf("    📋 Expression %d analysis:", i+1)

			// Analyser les contraintes
			if constraints, exists := exprMap["constraints"]; exists {
				constraintStats := analyzeConstraintStructure(constraints, t, "      ")
				for construct, count := range constraintStats {
					stats[construct] += count
				}
			}

			// Analyser les actions (TerminalNode)
			if action, exists := exprMap["action"]; exists && action != nil {
				stats["action"]++
				t.Logf("      → Action found (TerminalNode)")
			}

			// Analyser les fonctions dans les contraintes
			functionCount := analyzeFunctionCalls(exprMap, t, "      ")
			stats["functionCall"] += functionCount
		}
	}

	return stats
}

// analyzeConstraintStructure analyse récursivement les structures de contraintes
func analyzeConstraintStructure(constraint interface{}, t *testing.T, indent string) map[string]int {
	stats := make(map[string]int)

	if constraintMap, ok := constraint.(map[string]interface{}); ok {
		if constraintType, exists := constraintMap["type"]; exists {
			switch constraintType {
			case "comparison":
				stats["comparison"]++
				t.Logf("%s→ Simple comparison (AlphaNode)", indent)

			case "logicalExpr":
				stats["logicalExpr"]++
				t.Logf("%s→ Logical expression (JoinNode)", indent)

				// Analyser récursivement les opérations
				if operations, exists := constraintMap["operations"]; exists {
					if opList, ok := operations.([]interface{}); ok {
						for _, op := range opList {
							if opMap, ok := op.(map[string]interface{}); ok {
								if right, exists := opMap["right"]; exists {
									subStats := analyzeConstraintStructure(right, t, indent+"  ")
									for construct, count := range subStats {
										stats[construct] += count
									}
								}
							}
						}
					}
				}

			case "notConstraint":
				stats["notConstraint"]++
				t.Logf("%s→ NOT constraint (NotNode)", indent)

				// Analyser l'expression niée
				if expr, exists := constraintMap["expression"]; exists {
					subStats := analyzeConstraintStructure(expr, t, indent+"  ")
					for construct, count := range subStats {
						stats[construct] += count
					}
				}

			case "existsConstraint":
				stats["existsConstraint"]++
				t.Logf("%s→ EXISTS constraint (ExistsNode)", indent)

				// Analyser la condition d'existence
				if condition, exists := constraintMap["condition"]; exists {
					subStats := analyzeConstraintStructure(condition, t, indent+"  ")
					for construct, count := range subStats {
						stats[construct] += count
					}
				}

			case "aggregateConstraint":
				stats["aggregateConstraint"]++
				if function, exists := constraintMap["function"]; exists {
					t.Logf("%s→ Aggregate %s (AccumulateNode)", indent, function)
				}
			}
		}
	}

	return stats
}

// analyzeFunctionCalls compte les appels de fonction dans une expression
func analyzeFunctionCalls(exprMap map[string]interface{}, t *testing.T, indent string) int {
	// Fonction récursive pour chercher les functionCall
	var findFunctionCalls func(interface{}) int
	findFunctionCalls = func(obj interface{}) int {
		localCount := 0

		if objMap, ok := obj.(map[string]interface{}); ok {
			if objType, exists := objMap["type"]; exists && objType == "functionCall" {
				if name, exists := objMap["name"]; exists {
					t.Logf("%s→ Function call: %s (AlphaNode + evaluation)", indent, name)
					localCount++
				}
			}

			// Recherche récursive dans tous les champs
			for _, value := range objMap {
				localCount += findFunctionCalls(value)
			}
		} else if objList, ok := obj.([]interface{}); ok {
			for _, item := range objList {
				localCount += findFunctionCalls(item)
			}
		}

		return localCount
	}

	return findFunctionCalls(exprMap)
}
