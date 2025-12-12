// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import "fmt"

// alpha_activation_helpers.go contient des fonctions helper pour l'activation des nœuds alpha.
// Ces fonctions ont été extraites de ActivateWithContext() pour réduire la complexité.

// verifyDependencies vérifie que toutes les dépendances sont satisfaites dans le contexte
func verifyDependencies(dependencies []string, context *EvaluationContext, nodeID string) error {
	for _, dep := range dependencies {
		if !context.HasIntermediateResult(dep) {
			return fmt.Errorf("dependency %s not satisfied for node %s", dep, nodeID)
		}
	}
	return nil
}

// buildDependenciesMap construit une map des dépendances depuis le contexte
func buildDependenciesMap(dependencies []string, context *EvaluationContext) map[string]interface{} {
	dependenciesMap := make(map[string]interface{})
	for _, dep := range dependencies {
		if val, exists := context.GetIntermediateResult(dep); exists {
			dependenciesMap[dep] = val
		}
	}
	return dependenciesMap
}

// EvaluationResult contient le résultat d'une évaluation de condition
type EvaluationResult struct {
	Result    interface{}
	FromCache bool
	Error     error
}

// tryGetFromCache tente de récupérer le résultat depuis le cache
func tryGetFromCache(an *AlphaNode, context *EvaluationContext) (interface{}, bool) {
	if an.ResultName == "" || context.Cache == nil {
		return nil, false
	}

	dependencies := buildDependenciesMap(an.Dependencies, context)
	if cachedResult, found := context.Cache.GetWithDependencies(an.ResultName, dependencies); found {
		return cachedResult, true
	}

	return nil, false
}

// evaluateConditionWithContext évalue une condition avec le contexte
func evaluateConditionWithContext(an *AlphaNode, fact *Fact, context *EvaluationContext) (interface{}, error) {
	evaluator := NewConditionEvaluator(an.Storage)
	result, err := evaluator.EvaluateWithContext(an.Condition, fact, context)
	if err != nil {
		return nil, fmt.Errorf("error evaluating condition with context in node %s: %w", an.ID, err)
	}
	return result, nil
}

// storeInCache stocke le résultat dans le cache si approprié
func storeInCache(an *AlphaNode, context *EvaluationContext, result interface{}) {
	if an.ResultName == "" || context.Cache == nil {
		return
	}

	dependencies := buildDependenciesMap(an.Dependencies, context)
	context.Cache.SetWithDependencies(an.ResultName, dependencies, result)
}

// evaluateAtomicCondition évalue une condition atomique avec cache
func evaluateAtomicCondition(an *AlphaNode, fact *Fact, context *EvaluationContext) (*EvaluationResult, error) {
	// Essayer le cache d'abord
	if cachedResult, found := tryGetFromCache(an, context); found {
		return &EvaluationResult{
			Result:    cachedResult,
			FromCache: true,
		}, nil
	}

	// Évaluer la condition
	result, err := evaluateConditionWithContext(an, fact, context)
	if err != nil {
		return nil, err
	}

	// Stocker dans le cache
	storeInCache(an, context, result)

	return &EvaluationResult{
		Result:    result,
		FromCache: false,
	}, nil
}

// evaluateNonAtomicCondition évalue une condition non-atomique
func evaluateNonAtomicCondition(an *AlphaNode, fact *Fact) (bool, error) {
	if an.Condition == nil {
		return true, nil
	}

	evaluator := NewAlphaConditionEvaluator()
	passed, err := evaluator.EvaluateCondition(an.Condition, fact, an.VariableName)
	if err != nil {
		return false, fmt.Errorf("error evaluating condition in node %s: %w", an.ID, err)
	}

	return passed, nil
}

// storeIntermediateResult stocke un résultat intermédiaire dans le contexte
func storeIntermediateResult(an *AlphaNode, context *EvaluationContext, result interface{}) {
	if an.ResultName != "" {
		context.SetIntermediateResult(an.ResultName, result)
	}
}

// shouldPropagateResult détermine si le résultat doit être propagé
func shouldPropagateResult(condition interface{}, result interface{}) bool {
	if !isComparisonCondition(condition) {
		return true
	}

	// Pour les conditions de comparaison, vérifier si le résultat est false
	if boolResult, ok := result.(bool); ok && !boolResult {
		return false
	}

	return true
}

// addFactToMemory ajoute un fait à la mémoire du nœud (opération idempotente)
func addFactToMemory(an *AlphaNode, fact *Fact) (bool, error) {
	an.mutex.Lock()
	defer an.mutex.Unlock()

	internalID := fact.GetInternalID()
	_, alreadyExists := an.Memory.Facts[internalID]

	if !alreadyExists {
		if err := an.Memory.AddFact(fact); err != nil {
			return false, fmt.Errorf("error adding fact to alpha node: %w", err)
		}
	}

	return alreadyExists, nil
}

// isPassthroughRightNode détermine si le nœud est un passthrough RIGHT
func isPassthroughRightNode(condition interface{}) bool {
	condMap, ok := condition.(map[string]interface{})
	if !ok {
		return false
	}

	condType, exists := condMap["type"].(string)
	if !exists || condType != "passthrough" {
		return false
	}

	side, sideExists := condMap["side"].(string)
	return sideExists && side == "right"
}

// propagateToAlphaChild propage le fait à un enfant AlphaNode
func propagateToAlphaChild(alphaChild *AlphaNode, fact *Fact, context *EvaluationContext) error {
	if err := alphaChild.ActivateWithContext(fact, context); err != nil {
		return fmt.Errorf("error propagating to alpha child %s: %w", alphaChild.GetID(), err)
	}
	return nil
}

// propagateToNonAlphaChild propage le fait à un enfant non-alpha
func propagateToNonAlphaChild(an *AlphaNode, child Node, fact *Fact) error {
	if isPassthroughRightNode(an.Condition) {
		// Passthrough RIGHT: use ActivateRight for JoinNode
		if err := child.ActivateRight(fact); err != nil {
			return fmt.Errorf("error propagating fact to %s: %w", child.GetID(), err)
		}
	} else {
		// Passthrough LEFT ou nœud atomique final: créer token et utiliser ActivateLeft
		token := &Token{
			ID:       fmt.Sprintf("token_%s_%s", an.ID, fact.ID),
			Facts:    []*Fact{fact},
			NodeID:   an.ID,
			Bindings: NewBindingChainWith(an.VariableName, fact),
		}
		if err := child.ActivateLeft(token); err != nil {
			return fmt.Errorf("error propagating token to %s: %w", child.GetID(), err)
		}
	}
	return nil
}

// propagateToChildren propage le fait à tous les enfants du nœud
func propagateToChildren(an *AlphaNode, fact *Fact, context *EvaluationContext) error {
	for _, child := range an.GetChildren() {
		if alphaChild, ok := child.(*AlphaNode); ok {
			// Enfant alpha: propager avec contexte
			if err := propagateToAlphaChild(alphaChild, fact, context); err != nil {
				return err
			}
		} else {
			// Enfant non-alpha (JoinNode, TerminalNode, etc.)
			if err := propagateToNonAlphaChild(an, child, fact); err != nil {
				return err
			}
		}
	}
	return nil
}
