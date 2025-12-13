// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"time"
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
			ID:          nodeID,
			Type:        "alpha",
			Memory:      &WorkingMemory{NodeID: nodeID, Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
			Children:    make([]Node, 0),
			Storage:     storage,
			createdAt:   time.Now(),
		},
		Condition:    condition,
		VariableName: variableName,
		Dependencies: make([]string, 0),
	}
}

// ActivateLeft propage le token aux enfants (utilisé dans les cascades de jointures)
func (an *AlphaNode) ActivateLeft(token *Token) error {
	// Enregistrer l'activation
	an.recordActivation()
	
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

// ActivateRight teste la condition sur le fait.
// Refactorisé pour réduire la complexité de 18 à <10.
func (an *AlphaNode) ActivateRight(fact *Fact) error {
	// Enregistrer l'activation
	an.recordActivation()
	
	// Mode passthrough: pas de filtrage, conversion directe en token
	if isPassthroughCondition(an.Condition) {
		return an.handlePassthrough(fact)
	}

	// Évaluation normale de condition Alpha
	if !an.evaluateAlphaCondition(fact) {
		return nil // Condition non satisfaite, ignorer le fait
	}

	// Ajouter le fait à la mémoire (avec gestion d'idempotence)
	alreadyExists, err := an.addFactToMemory(fact)
	if err != nil {
		return err
	}

	// Si le fait existait déjà, ne pas propager à nouveau
	if alreadyExists {
		return nil
	}

	// Propager aux enfants selon leur type
	return an.propagateFactToChildren(fact)
}

// isPassthroughCondition vérifie si c'est une condition de passthrough.
func isPassthroughCondition(condition interface{}) bool {
	if condition == nil {
		return false
	}

	condMap, ok := condition.(map[string]interface{})
	if !ok {
		return false
	}

	condType, exists := condMap["type"].(string)
	return exists && condType == "passthrough"
}

// handlePassthrough gère le mode passthrough (convertir fait en token et propager).
func (an *AlphaNode) handlePassthrough(fact *Fact) error {
	condMap := an.Condition.(map[string]interface{})

	// Ajouter le fait à la mémoire
	an.mutex.Lock()
	if err := an.Memory.AddFact(fact); err != nil {
		an.mutex.Unlock()
		return fmt.Errorf("erreur ajout fait dans alpha node: %w", err)
	}
	an.mutex.Unlock()

	// Créer un token pour le fait
	token := &Token{
		ID:       fmt.Sprintf("alpha_token_%s_%s", an.ID, fact.ID),
		Facts:    []*Fact{fact},
		NodeID:   an.ID,
		Bindings: NewBindingChainWith(an.VariableName, fact),
	}

	// Déterminer le côté et propager selon l'architecture RETE
	side, sideExists := condMap["side"].(string)
	if sideExists && side == "left" {
		return an.PropagateToChildren(nil, token) // ActivateLeft
	}

	return an.PropagateToChildren(fact, nil) // ActivateRight
}

// evaluateAlphaCondition évalue la condition alpha sur le fait.
func (an *AlphaNode) evaluateAlphaCondition(fact *Fact) bool {
	if an.Condition == nil {
		return true
	}

	evaluator := NewAlphaConditionEvaluator()
	passed, err := evaluator.EvaluateCondition(an.Condition, fact, an.VariableName)
	if err != nil {
		// Log l'erreur sans retourner d'erreur (comportement existant)
		return false
	}

	return passed
}

// addFactToMemory ajoute un fait à la mémoire alpha avec gestion d'idempotence.
// Retourne (alreadyExists, error).
func (an *AlphaNode) addFactToMemory(fact *Fact) (bool, error) {
	an.mutex.Lock()
	defer an.mutex.Unlock()

	internalID := fact.GetInternalID()
	_, alreadyExists := an.Memory.Facts[internalID]

	if !alreadyExists {
		if err := an.Memory.AddFact(fact); err != nil {
			return false, fmt.Errorf("erreur ajout fait dans alpha node: %w", err)
		}
	}

	return alreadyExists, nil
}

// propagateFactToChildren propage le fait aux enfants selon leur type.
func (an *AlphaNode) propagateFactToChildren(fact *Fact) error {
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
				Bindings: NewBindingChainWith(an.VariableName, fact),
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
