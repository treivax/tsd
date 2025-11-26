// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
)

type TerminalNode struct {
	BaseNode
	Action *Action `json:"action"`
}

// NewTerminalNode crée un nouveau nœud terminal
func NewTerminalNode(nodeID string, action *Action, storage Storage) *TerminalNode {
	return &TerminalNode{
		BaseNode: BaseNode{
			ID:       nodeID,
			Type:     "terminal",
			Memory:   &WorkingMemory{NodeID: nodeID, Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
			Children: make([]Node, 0), // Les nœuds terminaux n'ont pas d'enfants
			Storage:  storage,
		},
		Action: action,
	}
}

// ActivateLeft déclenche l'action
func (tn *TerminalNode) ActivateLeft(token *Token) error {
	// Log désactivé pour les performances
	// fmt.Printf("[TERMINAL_%s] Déclenchement action avec token: %s\n", tn.ID, token.ID)

	// Stocker le token
	tn.mutex.Lock()
	if tn.Memory.Tokens == nil {
		tn.Memory.Tokens = make(map[string]*Token)
	}
	tn.Memory.Tokens[token.ID] = token
	tn.mutex.Unlock()

	// Persistance désactivée pour les performances

	// Déclencher l'action
	return tn.executeAction(token)
}

// ActivateRetract retrait des tokens contenant le fait rétracté
// factID doit être l'identifiant interne (Type_ID)
func (tn *TerminalNode) ActivateRetract(factID string) error {
	tn.mutex.Lock()
	var tokensToRemove []string
	for tokenID, token := range tn.Memory.Tokens {
		for _, fact := range token.Facts {
			if fact.GetInternalID() == factID {
				tokensToRemove = append(tokensToRemove, tokenID)
				break
			}
		}
	}
	for _, tokenID := range tokensToRemove {
		delete(tn.Memory.Tokens, tokenID)
	}
	tn.mutex.Unlock()
	if len(tokensToRemove) > 0 {
		fmt.Printf("🗑️  [TERMINAL_%s] Rétractation: %d tokens retirés\n", tn.ID, len(tokensToRemove))
	}
	return nil
}

// GetTriggeredActions retourne les actions déclenchées (pour les tests)
func (tn *TerminalNode) GetTriggeredActions() []*Action {
	tn.mutex.RLock()
	defer tn.mutex.RUnlock()

	actions := make([]*Action, 0, len(tn.Memory.Tokens))
	for range tn.Memory.Tokens {
		actions = append(actions, tn.Action)
	}
	return actions
}

// ActivateRight (non utilisé pour les nœuds terminaux)
func (tn *TerminalNode) ActivateRight(fact *Fact) error {
	return fmt.Errorf("les nœuds terminaux ne reçoivent pas de faits directement")
}

// executeAction affiche l'action déclenchée avec les faits déclencheurs (version tuple-space)
func (tn *TerminalNode) executeAction(token *Token) error {
	// Les actions sont maintenant obligatoires dans la grammaire
	// Mais nous gardons cette vérification par sécurité
	if tn.Action == nil {
		return fmt.Errorf("aucune action définie pour le nœud %s", tn.ID)
	}

	// === VERSION TUPLE-SPACE ===
	// Au lieu d'exécuter l'action, on l'affiche avec les faits déclencheurs
	// Les agents du tuple-space viendront "prendre" ces tuples plus tard

	actionName := tn.Action.Job.Name
	fmt.Printf("🎯 ACTION DISPONIBLE DANS TUPLE-SPACE: %s", actionName)

	// Afficher les faits déclencheurs entre parenthèses
	if len(token.Facts) > 0 {
		fmt.Print(" (")
		for i, fact := range token.Facts {
			if i > 0 {
				fmt.Print(", ")
			}
			// Format compact : Type(id:value, field:value, ...)
			fmt.Printf("%s(", fact.Type)
			fieldCount := 0
			for key, value := range fact.Fields {
				if fieldCount > 0 {
					fmt.Print(", ")
				}
				fmt.Printf("%s:%v", key, value)
				fieldCount++
			}
			fmt.Print(")")
		}
		fmt.Print(")")
	}
	fmt.Println()

	return nil
}
