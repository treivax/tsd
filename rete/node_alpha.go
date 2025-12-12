// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
)

type AlphaNode struct {
	BaseNode
	Condition    interface{} `json:"condition"`
	VariableName string      `json:"variable_name"`

	// Decomposition support fields
	ResultName   string   `json:"result_name,omitempty"`  // Name of intermediate result produced (e.g., "temp_1")
	IsAtomic     bool     `json:"is_atomic,omitempty"`    // true if atomic operation (single step)
	Dependencies []string `json:"dependencies,omitempty"` // Required intermediate results
}

// NewAlphaNode crée un nouveau nœud alpha
func NewAlphaNode(nodeID string, condition interface{}, variableName string, storage Storage) *AlphaNode {
	return &AlphaNode{
		BaseNode: BaseNode{
			ID:       nodeID,
			Type:     "alpha",
			Memory:   &WorkingMemory{NodeID: nodeID, Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
			Children: make([]Node, 0),
			Storage:  storage,
		},
		Condition:    condition,
		VariableName: variableName,
		Dependencies: make([]string, 0),
	}
}

// ActivateLeft propage le token aux enfants (utilisé dans les cascades de jointures)
func (an *AlphaNode) ActivateLeft(token *Token) error {
	// Dans les cascades de jointures, l'AlphaNode doit propager le token
	// pour préserver tous les bindings accumulés dans les jointures précédentes

	// Si l'AlphaNode a une condition de filtrage, l'appliquer
	if an.Condition != nil {
		// Vérifier si c'est un passthrough
		if condMap, ok := an.Condition.(map[string]interface{}); ok {
			if condType, exists := condMap["type"].(string); exists && condType == "passthrough" {
				// Mode passthrough: propager le token tel quel
				return an.PropagateToChildren(nil, token)
			}
		}

		// Si l'AlphaNode a une vraie condition, évaluer sur chaque fait du token
		// Pour l'instant, on propage tel quel (le filtrage sera fait ailleurs si nécessaire)
	}

	// Propager le token aux enfants
	return an.PropagateToChildren(nil, token)
}

// ActivateRetract retire le fait de la mémoire alpha et propage aux enfants
// factID doit être l'identifiant interne (Type_ID)
func (an *AlphaNode) ActivateRetract(factID string) error {
	an.mutex.Lock()
	_, exists := an.Memory.GetFact(factID)
	if exists {
		an.Memory.RemoveFact(factID)
	}
	an.mutex.Unlock()
	if !exists {
		return nil
	}
	fmt.Printf("🗑️  [ALPHA_%s] Rétractation du fait: %s\n", an.ID, factID)
	return an.PropagateRetractToChildren(factID)
}

// ActivateRight teste la condition sur le fait
func (an *AlphaNode) ActivateRight(fact *Fact) error {
	// Log désactivé pour les performances
	// fmt.Printf("[ALPHA_%s] Test condition sur fait: %s\n", an.ID, fact.String())

	// Cas spécial: passthrough pour les JoinNodes - pas de filtrage
	if an.Condition != nil {
		if condMap, ok := an.Condition.(map[string]interface{}); ok {
			if condType, exists := condMap["type"].(string); exists && condType == "passthrough" {
				// Mode pass-through: convertir le fait en token et propager selon le côté
				an.mutex.Lock()
				if err := an.Memory.AddFact(fact); err != nil {
					an.mutex.Unlock()
					return fmt.Errorf("erreur ajout fait dans alpha node: %w", err)
				}
				an.mutex.Unlock() // Créer un token pour le fait avec la variable correspondante
				token := &Token{
					ID:       fmt.Sprintf("alpha_token_%s_%s", an.ID, fact.ID),
					Facts:    []*Fact{fact},
					NodeID:   an.ID,
					Bindings: map[string]*Fact{an.VariableName: fact},
				}

				// Déterminer le côté et propager selon l'architecture RETE
				side, sideExists := condMap["side"].(string)
				if sideExists && side == "left" {
					return an.PropagateToChildren(nil, token) // ActivateLeft
				} else {
					return an.PropagateToChildren(fact, nil) // ActivateRight
				}
			}
		}
	}

	// Évaluation normale de condition Alpha
	if an.Condition != nil {
		evaluator := NewAlphaConditionEvaluator()
		passed, err := evaluator.EvaluateCondition(an.Condition, fact, an.VariableName)
		if err != nil {
			return fmt.Errorf("erreur évaluation condition Alpha: %w", err)
		}

		// Si la condition n'est pas satisfaite, ignorer le fait
		if !passed {
			// Log désactivé pour les performances
			// fmt.Printf("[ALPHA_%s] Condition non satisfaite pour le fait: %s\n", an.ID, fact.String())
			return nil
		}
	}

	// Log désactivé pour les performances
	// fmt.Printf("[ALPHA_%s] Condition satisfaite pour le fait: %s\n", an.ID, fact.String())

	// Vérifier si le fait existe déjà (idempotence pour les propagations multiples)
	an.mutex.Lock()
	internalID := fact.GetInternalID()
	_, alreadyExists := an.Memory.Facts[internalID]
	if !alreadyExists {
		if err := an.Memory.AddFact(fact); err != nil {
			an.mutex.Unlock()
			return fmt.Errorf("erreur ajout fait dans alpha node: %w", err)
		}
	}
	an.mutex.Unlock()

	// Si le fait existait déjà, ne pas propager à nouveau
	if alreadyExists {
		return nil
	}

	// Persistance désactivée pour les performances

	// Propager aux enfants
	// Dans une chaîne d'AlphaNodes, propager le fait directement via ActivateRight
	// Pour les autres types de nœuds (Terminal, Join), créer un token et propager via ActivateLeft
	for _, child := range an.GetChildren() {
		childType := child.GetType()

		if childType == "alpha" {
			// Propager le fait directement aux AlphaNodes enfants (chaîne)
			if err := child.ActivateRight(fact); err != nil {
				return fmt.Errorf("erreur propagation fait vers %s: %w", child.GetID(), err)
			}
		} else {
			// Pour les autres types de nœuds, créer un token et propager via ActivateLeft
			token := &Token{
				ID:       fmt.Sprintf("token_%s_%s", an.ID, fact.ID),
				Facts:    []*Fact{fact},
				NodeID:   an.ID,
				Bindings: map[string]*Fact{an.VariableName: fact},
			}
			if err := child.ActivateLeft(token); err != nil {
				return fmt.Errorf("erreur propagation token vers %s: %w", child.GetID(), err)
			}
		}
	}

	return nil
}

// ActivateWithContext activates the node with an evaluation context for decomposed expressions.
// This method supports intermediate result propagation through the context.
func (an *AlphaNode) ActivateWithContext(fact *Fact, context *EvaluationContext) error {
	// Verify all dependencies are satisfied
	for _, dep := range an.Dependencies {
		if !context.HasIntermediateResult(dep) {
			return fmt.Errorf("dependency %s not satisfied for node %s", dep, an.ID)
		}
	}

	// For atomic nodes, evaluate condition with context
	var result interface{}
	var err error

	if an.IsAtomic && an.Condition != nil {
		// Try cache first if ResultName is set
		var fromCache bool
		if an.ResultName != "" && context.Cache != nil {
			// Build dependencies map from context
			dependencies := make(map[string]interface{})
			for _, dep := range an.Dependencies {
				if val, exists := context.GetIntermediateResult(dep); exists {
					dependencies[dep] = val
				}
			}

			// Try to get from cache
			if cachedResult, found := context.Cache.GetWithDependencies(an.ResultName, dependencies); found {
				result = cachedResult
				fromCache = true
			}
		}

		// If not from cache, evaluate
		if !fromCache {
			// Use ConditionEvaluator for context-aware evaluation
			evaluator := NewConditionEvaluator(an.Storage)
			result, err = evaluator.EvaluateWithContext(an.Condition, fact, context)
			if err != nil {
				return fmt.Errorf("error evaluating condition with context in node %s: %w", an.ID, err)
			}

			// Store in cache if ResultName is set
			if an.ResultName != "" && context.Cache != nil {
				dependencies := make(map[string]interface{})
				for _, dep := range an.Dependencies {
					if val, exists := context.GetIntermediateResult(dep); exists {
						dependencies[dep] = val
					}
				}
				context.Cache.SetWithDependencies(an.ResultName, dependencies, result)
			}
		}

		// Store intermediate result if this node produces one
		if an.ResultName != "" {
			context.SetIntermediateResult(an.ResultName, result)
		}

		// For comparison conditions, check if result is false
		if isComparisonCondition(an.Condition) {
			if boolResult, ok := result.(bool); ok && !boolResult {
				// Condition not satisfied, don't propagate
				return nil
			}
		}
	} else {
		// Non-atomic node: use standard evaluation
		if an.Condition != nil {
			evaluator := NewAlphaConditionEvaluator()
			passed, err := evaluator.EvaluateCondition(an.Condition, fact, an.VariableName)
			if err != nil {
				return fmt.Errorf("error evaluating condition in node %s: %w", an.ID, err)
			}
			if !passed {
				return nil
			}
		}
	}

	// Add fact to memory (idempotent)
	an.mutex.Lock()
	internalID := fact.GetInternalID()
	_, alreadyExists := an.Memory.Facts[internalID]
	if !alreadyExists {
		if err := an.Memory.AddFact(fact); err != nil {
			an.mutex.Unlock()
			return fmt.Errorf("error adding fact to alpha node: %w", err)
		}
	}
	an.mutex.Unlock()

	if alreadyExists {
		return nil
	}

	// Propagate to children with context
	for _, child := range an.GetChildren() {
		if alphaChild, ok := child.(*AlphaNode); ok {
			// Propagate with context to alpha children
			if err := alphaChild.ActivateWithContext(fact, context); err != nil {
				return fmt.Errorf("error propagating to alpha child %s: %w", child.GetID(), err)
			}
		} else {
			// For non-alpha nodes (JoinNode, TerminalNode, etc.)
			// Check if this is a passthrough node to determine activation method
			isPassthroughRight := false
			if an.Condition != nil {
				if condMap, ok := an.Condition.(map[string]interface{}); ok {
					if condType, exists := condMap["type"].(string); exists && condType == "passthrough" {
						if side, sideExists := condMap["side"].(string); sideExists && side == "right" {
							isPassthroughRight = true
						}
					}
				}
			}

			if isPassthroughRight {
				// Passthrough RIGHT: use ActivateRight for JoinNode
				if err := child.ActivateRight(fact); err != nil {
					return fmt.Errorf("error propagating fact to %s: %w", child.GetID(), err)
				}
			} else {
				// Passthrough LEFT or final atomic node: create token and use ActivateLeft
				token := &Token{
					ID:       fmt.Sprintf("token_%s_%s", an.ID, fact.ID),
					Facts:    []*Fact{fact},
					NodeID:   an.ID,
					Bindings: map[string]*Fact{an.VariableName: fact},
				}
				if err := child.ActivateLeft(token); err != nil {
					return fmt.Errorf("error propagating token to %s: %w", child.GetID(), err)
				}
			}
		}
	}

	return nil
}

// isComparisonCondition checks if a condition is a comparison operation
func isComparisonCondition(condition interface{}) bool {
	if condMap, ok := condition.(map[string]interface{}); ok {
		if condType, exists := condMap["type"].(string); exists {
			return condType == "comparison"
		}
	}
	return false
}

// Clone crée une copie profonde de l'AlphaNode
func (an *AlphaNode) Clone() *AlphaNode {
	clone := &AlphaNode{
		BaseNode: BaseNode{
			ID:       an.ID,
			Type:     an.Type,
			Memory:   an.Memory.Clone(),
			Children: make([]Node, len(an.Children)),
			Storage:  an.Storage,
		},
		Condition:    an.Condition,
		VariableName: an.VariableName,
		ResultName:   an.ResultName,
		IsAtomic:     an.IsAtomic,
		Dependencies: make([]string, len(an.Dependencies)),
	}

	// Copier les enfants
	copy(clone.Children, an.Children)

	// Copier les dépendances
	copy(clone.Dependencies, an.Dependencies)

	return clone
}
