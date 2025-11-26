package rete

// AlphaConditionEvaluator évalue les conditions Alpha sur les faits.
// Cette structure a été refactorisée en plusieurs fichiers pour améliorer la lisibilité:
//   - evaluator_expressions.go: Évaluation des expressions binaires et logiques
//   - evaluator_constraints.go: Évaluation des contraintes
//   - evaluator_values.go: Évaluation des valeurs, champs et variables
//   - evaluator_comparisons.go: Opérations de comparaison
//   - evaluator_operators.go: Opérateurs arithmétiques, chaînes et listes
//   - evaluator_functions.go: Fonctions intégrées (LENGTH, UPPER, ABS, etc.)
type AlphaConditionEvaluator struct {
	variableBindings map[string]*Fact
	partialEvalMode  bool // Mode d'évaluation partielle pour les jointures en cascade
}

// NewAlphaConditionEvaluator crée un nouvel évaluateur de conditions
func NewAlphaConditionEvaluator() *AlphaConditionEvaluator {
	return &AlphaConditionEvaluator{
		variableBindings: make(map[string]*Fact),
		partialEvalMode:  false,
	}
}

// EvaluateCondition évalue une condition sur un fait.
// Il s'agit du point d'entrée principal pour l'évaluation des conditions Alpha.
//
// Parameters:
//   - condition: La condition à évaluer (peut être une map, BinaryOperation, LogicalExpression, etc.)
//   - fact: Le fait sur lequel évaluer la condition
//   - variableName: Le nom de la variable à lier au fait (optionnel)
//
// Returns:
//   - bool: true si la condition est satisfaite, false sinon
//   - error: Une erreur si l'évaluation échoue
func (e *AlphaConditionEvaluator) EvaluateCondition(condition interface{}, fact *Fact, variableName string) (bool, error) {
	// Si c'est un passthrough (agrégation), laisser passer tous les faits
	if condMap, ok := condition.(map[string]interface{}); ok {
		if condType, exists := condMap["type"].(string); exists && condType == "passthrough" {
			return true, nil
		}
	}

	// Lier la variable au fait pour l'évaluation
	if variableName != "" {
		e.variableBindings[variableName] = fact
	}

	return e.evaluateExpression(condition)
}

// ClearBindings efface les liaisons de variables.
// Utilisé pour réinitialiser l'évaluateur entre différentes évaluations.
func (e *AlphaConditionEvaluator) ClearBindings() {
	e.variableBindings = make(map[string]*Fact)
}

// SetPartialEvalMode active ou désactive le mode d'évaluation partielle.
// En mode partiel, les variables non liées renvoient true au lieu d'une erreur.
// Utilisé pour les jointures en cascade où toutes les variables ne sont pas encore disponibles.
func (e *AlphaConditionEvaluator) SetPartialEvalMode(enabled bool) {
	e.partialEvalMode = enabled
}

// GetBindings retourne les liaisons actuelles de variables.
// Utile pour le débogage et l'inspection de l'état de l'évaluateur.
func (e *AlphaConditionEvaluator) GetBindings() map[string]*Fact {
	return e.variableBindings
}

// Note: Les méthodes d'évaluation internes (evaluateExpression, evaluateValue, etc.)
// sont maintenant réparties dans les fichiers suivants pour améliorer la maintenabilité:
//
// evaluator_expressions.go:
//   - evaluateExpression
//   - evaluateMapExpression
//   - evaluateBinaryOperation
//   - evaluateBinaryOperationMap
//   - evaluateLogicalExpression
//   - evaluateLogicalExpressionMap
//
// evaluator_constraints.go:
//   - evaluateConstraint
//   - evaluateConstraintMap
//   - evaluateNegationConstraint
//   - evaluateNotConstraint
//   - evaluateExistsConstraint
//
// evaluator_values.go:
//   - evaluateValue
//   - evaluateValueFromMap
//   - evaluateFieldAccess
//   - evaluateFieldAccessByName
//   - evaluateVariable
//   - evaluateVariableByName
//
// evaluator_comparisons.go:
//   - compareValues
//   - normalizeValue
//   - areEqual
//   - isLess
//   - isGreater
//
// evaluator_operators.go:
//   - evaluateArithmeticOperation
//   - evaluateContains
//   - evaluateIn
//   - evaluateLike
//   - evaluateMatches
//
// evaluator_functions.go:
//   - evaluateFunctionCall
//   - evaluateLength, evaluateUpper, evaluateLower
//   - evaluateAbs, evaluateRound, evaluateFloor, evaluateCeil
//   - evaluateSubstring, evaluateTrim
