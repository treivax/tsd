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
	// Étape 1: Vérifier que toutes les dépendances sont satisfaites
	if err := verifyDependencies(an.Dependencies, context, an.ID); err != nil {
		return err
	}

	// Étape 2: Évaluer la condition selon le type de nœud (atomique ou non)
	if an.IsAtomic && an.Condition != nil {
		// Nœud atomique: évaluation avec cache et contexte
		evalResult, err := evaluateAtomicCondition(an, fact, context)
		if err != nil {
			return err
		}

		// Étape 3: Stocker le résultat intermédiaire si nécessaire
		storeIntermediateResult(an, context, evalResult.Result)

		// Étape 4: Vérifier si le résultat permet la propagation
		if !shouldPropagateResult(an.Condition, evalResult.Result) {
			return nil
		}
	} else {
		// Nœud non-atomique: évaluation standard
		passed, err := evaluateNonAtomicCondition(an, fact)
		if err != nil {
			return err
		}
		if !passed {
			return nil
		}
	}

	// Étape 5: Ajouter le fait à la mémoire (opération idempotente)
	alreadyExists, err := addFactToMemory(an, fact)
	if err != nil {
		return err
	}

	// Si le fait existe déjà, ne pas propager
	if alreadyExists {
		return nil
	}

	// Étape 6: Propager aux enfants avec contexte
	return propagateToChildren(an, fact, context)
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
