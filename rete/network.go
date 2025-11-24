package rete

import (
	"fmt"
)

// ReteNetwork représente le réseau RETE complet
type ReteNetwork struct {
	RootNode      *RootNode                `json:"root_node"`
	TypeNodes     map[string]*TypeNode     `json:"type_nodes"`
	AlphaNodes    map[string]*AlphaNode    `json:"alpha_nodes"`
	BetaNodes     map[string]interface{}   `json:"beta_nodes"` // Nœuds Beta pour les jointures multi-faits
	TerminalNodes map[string]*TerminalNode `json:"terminal_nodes"`
	Storage       Storage                  `json:"-"`
	Types         []TypeDefinition         `json:"types"`
	BetaBuilder   interface{}              `json:"-"` // Constructeur de réseau Beta
}

// NewReteNetwork crée un nouveau réseau RETE
func NewReteNetwork(storage Storage) *ReteNetwork {
	rootNode := NewRootNode(storage)

	return &ReteNetwork{
		RootNode:      rootNode,
		TypeNodes:     make(map[string]*TypeNode),
		AlphaNodes:    make(map[string]*AlphaNode),
		BetaNodes:     make(map[string]interface{}),
		TerminalNodes: make(map[string]*TerminalNode),
		Storage:       storage,
		Types:         make([]TypeDefinition, 0),
		BetaBuilder:   nil, // Sera initialisé si nécessaire
	}
}

// SubmitFact soumet un nouveau fait au réseau
func (rn *ReteNetwork) SubmitFact(fact *Fact) error {
	fmt.Printf("🔥 Soumission d'un nouveau fait au réseau RETE: %s\n", fact.String())

	// Propager le fait depuis le nœud racine
	return rn.RootNode.ActivateRight(fact)
}

// SubmitFactsFromGrammar soumet plusieurs faits depuis la grammaire au réseau
func (rn *ReteNetwork) SubmitFactsFromGrammar(facts []map[string]interface{}) error {
	for i, factMap := range facts {
		// Convertir le map en Fact
		factID := fmt.Sprintf("fact_%d", i)
		if id, ok := factMap["id"].(string); ok {
			factID = id
		}

		factType := "unknown"
		if typ, ok := factMap["type"].(string); ok {
			factType = typ
		}

		fact := &Fact{
			ID:     factID,
			Type:   factType,
			Fields: make(map[string]interface{}),
		}

		// Copier tous les champs
		for key, value := range factMap {
			if key != "id" && key != "type" {
				fact.Fields[key] = value
			}
		}

		if err := rn.SubmitFact(fact); err != nil {
			return fmt.Errorf("erreur soumission fait %s: %w", fact.ID, err)
		}
	}
	return nil
}

// RetractFact retire un fait du réseau et propage la rétractation
func (rn *ReteNetwork) RetractFact(factID string) error {
	fmt.Printf("🗑️  Rétractation du fait: %s\n", factID)

	// Vérifier que le fait existe dans le réseau
	memory := rn.RootNode.GetMemory()
	if _, exists := memory.GetFact(factID); !exists {
		return fmt.Errorf("fait %s introuvable dans le réseau", factID)
	}

	// Propager la rétractation depuis le nœud racine
	return rn.RootNode.ActivateRetract(factID)
}
