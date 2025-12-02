// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
)

// ArithmeticNode représente un nœud dans l'arbre de décomposition arithmétique
type ArithmeticNode struct {
	Type     string      // "binaryOp", "fieldAccess", "number", "comparison"
	Operator string      // "+", "-", "*", "/", ">", "<", "==", etc.
	Left     interface{} // Peut être un autre ArithmeticNode ou une valeur
	Right    interface{} // Peut être un autre ArithmeticNode ou une valeur
	Hash     string      // Hash unique pour le partage
}

// DecomposedExpression représente une expression décomposée en étapes atomiques
type DecomposedExpression struct {
	Steps    []SimpleCondition // Liste ordonnée des étapes de calcul
	OpType   string            // Type d'opération globale ("ARITHMETIC", "COMPARISON")
	RootHash string            // Hash de l'expression racine
}

// ArithmeticExpressionDecomposer décompose les expressions arithmétiques complexes
// en séquence d'opérations atomiques
type ArithmeticExpressionDecomposer struct {
	stepCounter int
}

// NewArithmeticExpressionDecomposer crée un nouveau décomposeur
func NewArithmeticExpressionDecomposer() *ArithmeticExpressionDecomposer {
	return &ArithmeticExpressionDecomposer{
		stepCounter: 0,
	}
}

// DecomposeExpression décompose une expression arithmétique en étapes atomiques
// Exemple: (a * 2 + b * 3) > 10
// Devient:
//
//	Step 1: a * 2 → temp1
//	Step 2: b * 3 → temp2
//	Step 3: temp1 + temp2 → temp3
//	Step 4: temp3 > 10 → result
func (aed *ArithmeticExpressionDecomposer) DecomposeExpression(expr interface{}) (*DecomposedExpression, error) {
	aed.stepCounter = 0

	steps := make([]SimpleCondition, 0)
	_, extractedSteps, err := aed.decomposeRecursive(expr, &steps)
	if err != nil {
		return nil, err
	}

	if len(extractedSteps) == 0 {
		return nil, fmt.Errorf("no steps extracted from expression")
	}

	return &DecomposedExpression{
		Steps:    extractedSteps,
		OpType:   "ARITHMETIC",
		RootHash: extractedSteps[len(extractedSteps)-1].Hash,
	}, nil
}

// decomposeRecursive décompose récursivement une expression
// Retourne: (interface représentant le résultat, steps extraites, erreur)
func (aed *ArithmeticExpressionDecomposer) decomposeRecursive(
	expr interface{},
	steps *[]SimpleCondition,
) (interface{}, []SimpleCondition, error) {

	switch e := expr.(type) {
	case map[string]interface{}:
		return aed.decomposeMap(e, steps)

	case string, int, int64, float64, bool:
		// Valeur primitive - retourner telle quelle
		return expr, *steps, nil

	default:
		return nil, *steps, fmt.Errorf("type non supporté pour décomposition: %T", expr)
	}
}

// decomposeMap décompose une expression sous forme de map
func (aed *ArithmeticExpressionDecomposer) decomposeMap(
	expr map[string]interface{},
	steps *[]SimpleCondition,
) (interface{}, []SimpleCondition, error) {

	exprType, ok := expr["type"].(string)
	if !ok {
		return expr, *steps, fmt.Errorf("type manquant dans map")
	}

	switch exprType {
	case "binaryOp", "binaryOperation":
		return aed.decomposeBinaryOp(expr, steps)

	case "comparison":
		return aed.decomposeComparison(expr, steps)

	case "fieldAccess":
		// Accès direct à un champ - pas de décomposition
		return expr, *steps, nil

	case "number", "numberLiteral":
		// Nombre - pas de décomposition
		return expr, *steps, nil

	case "string", "stringLiteral":
		// String - pas de décomposition
		return expr, *steps, nil

	default:
		// Type inconnu - retourner tel quel
		return expr, *steps, nil
	}
}

// decomposeBinaryOp décompose une opération binaire (*, +, -, /)
func (aed *ArithmeticExpressionDecomposer) decomposeBinaryOp(
	expr map[string]interface{},
	steps *[]SimpleCondition,
) (interface{}, []SimpleCondition, error) {

	operator, ok := expr["operator"].(string)
	if !ok {
		return expr, *steps, fmt.Errorf("operator manquant dans binaryOp")
	}

	left := expr["left"]
	right := expr["right"]

	// Décomposer récursivement left et right
	leftResult, _, err := aed.decomposeRecursive(left, steps)
	if err != nil {
		return nil, *steps, err
	}

	rightResult, _, err := aed.decomposeRecursive(right, steps)
	if err != nil {
		return nil, *steps, err
	}

	// Créer une nouvelle étape pour cette opération
	aed.stepCounter++
	stepName := fmt.Sprintf("temp_%d", aed.stepCounter)

	// Normaliser l'opérateur si nécessaire (décoder les opérateurs encodés)
	normalizedOp := normalizeOperator(operator)

	condition := NewSimpleCondition(
		"binaryOp",
		leftResult,
		normalizedOp,
		rightResult,
	)

	*steps = append(*steps, condition)

	// Retourner une référence symbolique au résultat de cette étape
	return map[string]interface{}{
		"type":      "tempResult",
		"step_name": stepName,
		"step_idx":  len(*steps) - 1,
		"hash":      condition.Hash,
	}, *steps, nil
}

// DecomposeToDecomposedConditions decomposes an expression into DecomposedCondition steps
// with full metadata (ResultName, Dependencies, IsAtomic)
func (aed *ArithmeticExpressionDecomposer) DecomposeToDecomposedConditions(expr interface{}) ([]DecomposedCondition, error) {
	aed.stepCounter = 0
	steps := make([]SimpleCondition, 0)

	_, simpleSteps, err := aed.decomposeRecursive(expr, &steps)
	if err != nil {
		return nil, err
	}

	// Convert SimpleCondition to DecomposedCondition with metadata
	decomposedSteps := make([]DecomposedCondition, 0, len(simpleSteps))

	for i, step := range simpleSteps {
		resultName := fmt.Sprintf("temp_%d", i+1)
		dependencies := extractDependenciesFromCondition(step)

		decomposed := DecomposedCondition{
			SimpleCondition: step,
			ResultName:      resultName,
			Dependencies:    dependencies,
			IsAtomic:        true,
		}

		decomposedSteps = append(decomposedSteps, decomposed)
	}

	return decomposedSteps, nil
}

// extractDependenciesFromCondition extracts tempResult dependencies from a condition
func extractDependenciesFromCondition(cond SimpleCondition) []string {
	deps := make([]string, 0)

	// Check left operand
	if leftMap, ok := cond.Left.(map[string]interface{}); ok {
		if leftMap["type"] == "tempResult" {
			if stepName, ok := leftMap["step_name"].(string); ok {
				deps = append(deps, stepName)
			}
		}
	}

	// Check right operand
	if rightMap, ok := cond.Right.(map[string]interface{}); ok {
		if rightMap["type"] == "tempResult" {
			if stepName, ok := rightMap["step_name"].(string); ok {
				deps = append(deps, stepName)
			}
		}
	}

	return deps
}

// decomposeComparison décompose une comparaison (>, <, ==, etc.)
func (aed *ArithmeticExpressionDecomposer) decomposeComparison(
	expr map[string]interface{},
	steps *[]SimpleCondition,
) (interface{}, []SimpleCondition, error) {

	operator, ok := expr["operator"].(string)
	if !ok {
		return expr, *steps, fmt.Errorf("operator manquant dans comparison")
	}

	left := expr["left"]
	right := expr["right"]

	// Décomposer récursivement left et right
	leftResult, _, err := aed.decomposeRecursive(left, steps)
	if err != nil {
		return nil, *steps, err
	}

	rightResult, _, err := aed.decomposeRecursive(right, steps)
	if err != nil {
		return nil, *steps, err
	}

	// Créer la condition de comparaison finale
	condition := NewSimpleCondition(
		"comparison",
		leftResult,
		operator,
		rightResult,
	)

	*steps = append(*steps, condition)

	return map[string]interface{}{
		"type":        "comparisonResult",
		"step_idx":    len(*steps) - 1,
		"hash":        condition.Hash,
		"is_terminal": true,
	}, *steps, nil
}

// normalizeOperator normalise les opérateurs encodés
func normalizeOperator(operator string) string {
	switch operator {
	case "Kg==":
		return "*"
	case "Kw==":
		return "+"
	case "LQ==":
		return "-"
	case "Lw==":
		return "/"
	default:
		return operator
	}
}

// ShouldDecompose détermine si une expression devrait être décomposée
// Retourne true si l'expression contient plusieurs opérations arithmétiques
func (aed *ArithmeticExpressionDecomposer) ShouldDecompose(expr interface{}) bool {
	count := aed.countOperations(expr)
	return count > 1
}

// countOperations compte le nombre d'opérations dans une expression
func (aed *ArithmeticExpressionDecomposer) countOperations(expr interface{}) int {
	exprMap, ok := expr.(map[string]interface{})
	if !ok {
		return 0
	}

	exprType, ok := exprMap["type"].(string)
	if !ok {
		return 0
	}

	switch exprType {
	case "binaryOp", "binaryOperation":
		// Compter cette opération + les opérations dans left et right
		leftCount := 0
		rightCount := 0

		if left := exprMap["left"]; left != nil {
			leftCount = aed.countOperations(left)
		}

		if right := exprMap["right"]; right != nil {
			rightCount = aed.countOperations(right)
		}

		return 1 + leftCount + rightCount

	case "comparison":
		// Compter les opérations dans left et right (mais pas la comparaison elle-même
		// car elle sera toujours la dernière étape)
		leftCount := 0
		rightCount := 0

		if left := exprMap["left"]; left != nil {
			leftCount = aed.countOperations(left)
		}

		if right := exprMap["right"]; right != nil {
			rightCount = aed.countOperations(right)
		}

		return leftCount + rightCount

	default:
		return 0
	}
}

// GetComplexity retourne la complexité d'une expression (nombre total d'opérations)
func (aed *ArithmeticExpressionDecomposer) GetComplexity(expr interface{}) int {
	return aed.countOperations(expr)
}

// SimplifySteps optimise la liste d'étapes en éliminant les redondances
func (aed *ArithmeticExpressionDecomposer) SimplifySteps(steps []SimpleCondition) []SimpleCondition {
	// Pour l'instant, retourner tel quel
	// Future optimisation: éliminer les étapes redondantes ou fusionner les étapes similaires
	return steps
}

// ValidateDecomposition vérifie que la décomposition est valide
func (aed *ArithmeticExpressionDecomposer) ValidateDecomposition(original interface{}, decomposed *DecomposedExpression) error {
	if decomposed == nil {
		return fmt.Errorf("decomposed expression is nil")
	}

	if len(decomposed.Steps) == 0 {
		return fmt.Errorf("no steps in decomposed expression")
	}

	// Vérifier que chaque étape a un hash unique
	seen := make(map[string]bool)
	for i, step := range decomposed.Steps {
		if step.Hash == "" {
			return fmt.Errorf("step %d has empty hash", i)
		}
		if seen[step.Hash] {
			// C'est OK d'avoir des hash dupliqués si les conditions sont identiques
			// (ça signifie qu'on peut partager le nœud)
		}
		seen[step.Hash] = true
	}

	return nil
}

// FormatSteps retourne une représentation lisible des étapes
func (aed *ArithmeticExpressionDecomposer) FormatSteps(steps []SimpleCondition) string {
	result := ""
	for i, step := range steps {
		result += fmt.Sprintf("Step %d: %s\n", i+1, CanonicalString(step))
	}
	return result
}

// GetDecompositionStats retourne des statistiques sur une décomposition
func (aed *ArithmeticExpressionDecomposer) GetDecompositionStats(decomposed *DecomposedExpression) map[string]interface{} {
	if decomposed == nil {
		return map[string]interface{}{
			"valid": false,
		}
	}

	operatorCounts := make(map[string]int)
	for _, step := range decomposed.Steps {
		operatorCounts[step.Operator]++
	}

	return map[string]interface{}{
		"valid":           true,
		"total_steps":     len(decomposed.Steps),
		"operator_counts": operatorCounts,
		"root_hash":       decomposed.RootHash,
	}
}
