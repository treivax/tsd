package rete

import (
	"fmt"
)

type AlphaNode struct {
	BaseNode
	Condition    interface{} `json:"condition"`
	VariableName string      `json:"variable_name"`
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
	}
}

// ActivateLeft (non utilisé pour les nœuds alpha)
func (an *AlphaNode) ActivateLeft(token *Token) error {
	return fmt.Errorf("les nœuds alpha ne reçoivent pas de tokens")
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
			an.mutex.Unlock()				// Créer un token pour le fait avec la variable correspondante
				token := &Token{
					ID:       fmt.Sprintf("alpha_token_%s_%s", an.ID, fact.ID),
					Facts:    []*Fact{fact},
					NodeID:   an.ID,
					Bindings: map[string]*Fact{an.VariableName: fact},
				}

				// Déterminer le côté et propager selon l'architecture RETE
				side, sideExists := condMap["side"].(string)
				if sideExists && side == "left" {
					fmt.Printf("🔗 ALPHA PASSTHROUGH[%s]: Propagation LEFT pour JoinNode\n", an.ID)
					return an.PropagateToChildren(nil, token) // ActivateLeft
				} else {
					fmt.Printf("🔗 ALPHA PASSTHROUGH[%s]: Propagation RIGHT pour JoinNode\n", an.ID)
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

	an.mutex.Lock()
	if err := an.Memory.AddFact(fact); err != nil {
		an.mutex.Unlock()
		return fmt.Errorf("erreur ajout fait dans alpha node: %w", err)
	}
	an.mutex.Unlock()

	// Persistance désactivée pour les performances

	// Créer un token et le propager
	token := &Token{
		ID:     fmt.Sprintf("token_%s_%s", an.ID, fact.ID),
		Facts:  []*Fact{fact},
		NodeID: an.ID,
	}

	return an.PropagateToChildren(nil, token)
}
