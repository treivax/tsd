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

// LoadFromAST construit le réseau RETE à partir d'un AST
func (rn *ReteNetwork) LoadFromAST(program *Program) error {
	fmt.Printf("🏗️  Construction du réseau RETE à partir de l'AST\n")

	// Charger les types
	rn.Types = program.Types
	fmt.Printf("   Types définis: %d\n", len(program.Types))

	// Créer les nœuds de type
	for _, typeDef := range program.Types {
		typeNode := NewTypeNode(typeDef.Name, typeDef, rn.Storage)
		rn.TypeNodes[typeDef.Name] = typeNode
		rn.RootNode.AddChild(typeNode)
		fmt.Printf("   ✓ Créé TypeNode: %s\n", typeDef.Name)
	}

	// Créer les nœuds pour chaque expression (règle)
	for i, expr := range program.Expressions {
		ruleID := fmt.Sprintf("rule_%d", i)
		fmt.Printf("   📋 Traitement de la règle: %s\n", ruleID)

		// Créer les nœuds alpha pour les conditions
		alphaNodes, err := rn.createAlphaNodes(expr, ruleID)
		if err != nil {
			return fmt.Errorf("erreur création nœuds alpha: %w", err)
		}

		// Créer le nœud terminal pour l'action
		terminalNode := NewTerminalNode(ruleID+"_terminal", expr.Action, rn.Storage)
		rn.TerminalNodes[terminalNode.ID] = terminalNode

		// Connecter les nœuds alpha au nœud terminal
		for _, alphaNode := range alphaNodes {
			alphaNode.AddChild(terminalNode)
			fmt.Printf("     ✓ Connecté AlphaNode %s -> TerminalNode %s\n", alphaNode.ID, terminalNode.ID)
		}
	}

	fmt.Printf("🎯 Réseau RETE construit avec succès!\n")
	fmt.Printf("   - %d TypeNodes\n", len(rn.TypeNodes))
	fmt.Printf("   - %d AlphaNodes\n", len(rn.AlphaNodes))
	fmt.Printf("   - %d BetaNodes\n", len(rn.BetaNodes))
	fmt.Printf("   - %d TerminalNodes\n", len(rn.TerminalNodes))

	return nil
}

// createAlphaNodes crée les nœuds alpha pour une expression
func (rn *ReteNetwork) createAlphaNodes(expr Expression, ruleID string) ([]*AlphaNode, error) {
	var alphaNodes []*AlphaNode

	// Identifier le type de la variable dans le set
	for _, variable := range expr.Set.Variables {
		typeNode, exists := rn.TypeNodes[variable.DataType]
		if !exists {
			return nil, fmt.Errorf("type non trouvé: %s", variable.DataType)
		}

		// Créer un nœud alpha pour cette condition
		alphaNodeID := fmt.Sprintf("%s_alpha_%s", ruleID, variable.Name)
		alphaNode := NewAlphaNode(alphaNodeID, expr.Constraints, variable.Name, rn.Storage)
		rn.AlphaNodes[alphaNodeID] = alphaNode

		// Connecter le nœud de type au nœud alpha
		typeNode.AddChild(alphaNode)

		alphaNodes = append(alphaNodes, alphaNode)
		fmt.Printf("     ✓ Créé AlphaNode: %s pour variable: %s\n", alphaNodeID, variable.Name)
	}

	return alphaNodes, nil
}

// GetNetworkState retourne l'état complet du réseau
func (rn *ReteNetwork) GetNetworkState() (map[string]*WorkingMemory, error) {
	state := make(map[string]*WorkingMemory)

	// Récupérer l'état de tous les nœuds
	nodes := []Node{rn.RootNode}

	for _, typeNode := range rn.TypeNodes {
		nodes = append(nodes, typeNode)
	}
	for _, alphaNode := range rn.AlphaNodes {
		nodes = append(nodes, alphaNode)
	}
	for _, terminalNode := range rn.TerminalNodes {
		nodes = append(nodes, terminalNode)
	}

	for _, node := range nodes {
		memory, err := rn.Storage.LoadMemory(node.GetID())
		if err != nil {
			// Si pas de mémoire sauvée, utiliser la mémoire courante
			memory = node.GetMemory()
		}
		state[node.GetID()] = memory
	}

	return state, nil
}

// PrintNetworkStructure affiche la structure du réseau
func (rn *ReteNetwork) PrintNetworkStructure() {
	fmt.Printf("\n📊 STRUCTURE DU RÉSEAU RETE:\n")
	fmt.Printf("Root: %s\n", rn.RootNode.GetID())

	for typeName, typeNode := range rn.TypeNodes {
		fmt.Printf("├── TypeNode[%s]: %s\n", typeName, typeNode.GetID())

		for _, child := range typeNode.GetChildren() {
			if alphaNode, ok := child.(*AlphaNode); ok {
				fmt.Printf("│   ├── AlphaNode: %s\n", alphaNode.GetID())

				for _, grandChild := range alphaNode.GetChildren() {
					if terminalNode, ok := grandChild.(*TerminalNode); ok {
						fmt.Printf("│   │   └── TerminalNode: %s\n", terminalNode.GetID())
					}
				}
			}
		}
	}

	// Afficher les nœuds Beta si présents
	if len(rn.BetaNodes) > 0 {
		fmt.Printf("Beta Nodes:\n")
		for nodeID := range rn.BetaNodes {
			fmt.Printf("├── BetaNode: %s\n", nodeID)
		}
	}

	fmt.Printf("\n")
}

// EnableBetaNodes active le support des nœuds Beta dans le réseau
// Cette méthode doit être appelée avant de créer des jointures multi-faits
func (rn *ReteNetwork) EnableBetaNodes() error {
	// Note: Cette implémentation utilise des interfaces génériques pour éviter
	// les dépendances circulaires. Dans une vraie implémentation, on utiliserait
	// directement les types du package network.
	fmt.Printf("🔗 Activation du support des nœuds Beta\n")

	// Placeholder pour l'initialisation du BetaNetworkBuilder
	// Dans la vraie implémentation, on ferait:
	// rn.BetaBuilder = network.NewBetaNetworkBuilder(logger)

	return nil
}

// CreateBetaJoin crée une jointure Beta entre deux sources de données
// Ceci est une méthode d'exemple montrant comment intégrer les nœuds Beta
func (rn *ReteNetwork) CreateBetaJoin(leftSource, rightSource, joinID string, conditions []interface{}) error {
	fmt.Printf("🔗 Création d'une jointure Beta: %s\n", joinID)
	fmt.Printf("   Sources: %s ↔ %s\n", leftSource, rightSource)
	fmt.Printf("   Conditions: %d\n", len(conditions))

	// Placeholder pour la création d'un nœud de jointure
	// Dans la vraie implémentation, on utiliserait le BetaBuilder
	rn.BetaNodes[joinID] = map[string]interface{}{
		"type":        "JoinNode",
		"id":          joinID,
		"leftSource":  leftSource,
		"rightSource": rightSource,
		"conditions":  conditions,
	}

	fmt.Printf("   ✓ Nœud Beta créé: %s\n", joinID)
	return nil
}

// GetBetaNodeStatistics retourne les statistiques des nœuds Beta
func (rn *ReteNetwork) GetBetaNodeStatistics() map[string]interface{} {
	stats := map[string]interface{}{
		"totalBetaNodes": len(rn.BetaNodes),
		"betaEnabled":    rn.BetaBuilder != nil,
		"nodes":          make(map[string]interface{}),
	}

	for nodeID, node := range rn.BetaNodes {
		stats["nodes"].(map[string]interface{})[nodeID] = node
	}

	return stats
}
