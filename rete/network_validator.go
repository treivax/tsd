// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import "fmt"

// ValidateNetwork vérifie la cohérence et l'intégrité du réseau RETE
func (rn *ReteNetwork) ValidateNetwork() error {
	if err := rn.validateStructure(); err != nil {
		return fmt.Errorf("erreur de validation de structure: %w", err)
	}

	if err := rn.validateNodeReferences(); err != nil {
		return fmt.Errorf("erreur de validation des références: %w", err)
	}

	if err := rn.validateLifecycle(); err != nil {
		return fmt.Errorf("erreur de validation du lifecycle: %w", err)
	}

	return nil
}

// validateStructure vérifie la structure de base du réseau
func (rn *ReteNetwork) validateStructure() error {
	// Vérifier que le RootNode existe
	if rn.RootNode == nil {
		return fmt.Errorf("RootNode est nil")
	}

	// Vérifier que le Storage existe
	if rn.Storage == nil {
		return fmt.Errorf("Storage est nil")
	}

	// Vérifier que le LifecycleManager existe
	if rn.LifecycleManager == nil {
		return fmt.Errorf("LifecycleManager est nil")
	}

	return nil
}

// validateNodeReferences vérifie que toutes les références entre nœuds sont valides
func (rn *ReteNetwork) validateNodeReferences() error {
	// Vérifier que tous les enfants du RootNode sont des TypeNodes
	for _, child := range rn.RootNode.GetChildren() {
		if child.GetType() != "type" {
			return fmt.Errorf("RootNode a un enfant qui n'est pas un TypeNode: %s", child.GetID())
		}

		// Vérifier que le TypeNode existe dans la map
		found := false
		for _, typeNode := range rn.TypeNodes {
			if typeNode.GetID() == child.GetID() {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("TypeNode %s référencé par RootNode mais absent de TypeNodes map", child.GetID())
		}
	}

	// Vérifier les enfants des TypeNodes (doivent être AlphaNodes ou BetaNodes)
	for typeName, typeNode := range rn.TypeNodes {
		if typeNode == nil {
			return fmt.Errorf("TypeNode '%s' est nil", typeName)
		}

		for _, child := range typeNode.GetChildren() {
			childType := child.GetType()
			if childType != "alpha" && childType != "join" && childType != "exists" && childType != "accumulate" {
				return fmt.Errorf("TypeNode %s a un enfant de type invalide: %s", typeNode.GetID(), childType)
			}
		}
	}

	// Vérifier les enfants des AlphaNodes (peuvent être d'autres AlphaNodes ou TerminalNodes)
	for alphaID, alphaNode := range rn.AlphaNodes {
		if alphaNode == nil {
			return fmt.Errorf("AlphaNode '%s' est nil", alphaID)
		}

		for _, child := range alphaNode.GetChildren() {
			childType := child.GetType()
			if childType != "alpha" && childType != "terminal" && childType != "join" {
				return fmt.Errorf("AlphaNode %s a un enfant de type invalide: %s", alphaNode.GetID(), childType)
			}
		}
	}

	return nil
}

// validateLifecycle vérifie la cohérence du LifecycleManager
func (rn *ReteNetwork) validateLifecycle() error {
	if rn.LifecycleManager == nil {
		return nil // Déjà vérifié dans validateStructure
	}

	// Vérifier que tous les TypeNodes sont enregistrés
	for _, typeNode := range rn.TypeNodes {
		lifecycle, exists := rn.LifecycleManager.GetNodeLifecycle(typeNode.GetID())
		if !exists {
			return fmt.Errorf("TypeNode %s non enregistré dans LifecycleManager", typeNode.GetID())
		}
		if lifecycle.NodeType != "type" {
			return fmt.Errorf("TypeNode %s a un type incorrect dans LifecycleManager: %s", typeNode.GetID(), lifecycle.NodeType)
		}
	}

	// Vérifier que tous les AlphaNodes sont enregistrés
	for _, alphaNode := range rn.AlphaNodes {
		lifecycle, exists := rn.LifecycleManager.GetNodeLifecycle(alphaNode.GetID())
		if !exists {
			return fmt.Errorf("AlphaNode %s non enregistré dans LifecycleManager", alphaNode.GetID())
		}
		if lifecycle.NodeType != "alpha" {
			return fmt.Errorf("AlphaNode %s a un type incorrect dans LifecycleManager: %s", alphaNode.GetID(), lifecycle.NodeType)
		}
	}

	// Vérifier que tous les BetaNodes sont enregistrés
	for betaID := range rn.BetaNodes {
		lifecycle, exists := rn.LifecycleManager.GetNodeLifecycle(betaID)
		if !exists {
			return fmt.Errorf("BetaNode %s non enregistré dans LifecycleManager", betaID)
		}
		// Les BetaNodes peuvent avoir différents types: join, exists, accumulate
		if lifecycle.NodeType != "join" && lifecycle.NodeType != "exists" && lifecycle.NodeType != "accumulate" {
			return fmt.Errorf("BetaNode %s a un type incorrect dans LifecycleManager: %s", betaID, lifecycle.NodeType)
		}
	}

	// Vérifier que tous les TerminalNodes sont enregistrés
	for _, terminalNode := range rn.TerminalNodes {
		lifecycle, exists := rn.LifecycleManager.GetNodeLifecycle(terminalNode.GetID())
		if !exists {
			return fmt.Errorf("TerminalNode %s non enregistré dans LifecycleManager", terminalNode.GetID())
		}
		if lifecycle.NodeType != "terminal" {
			return fmt.Errorf("TerminalNode %s a un type incorrect dans LifecycleManager: %s", terminalNode.GetID(), lifecycle.NodeType)
		}
	}

	return nil
}

// ValidateRule vérifie qu'une règle est correctement construite dans le réseau
func (rn *ReteNetwork) ValidateRule(ruleID string) error {
	if rn.LifecycleManager == nil {
		return fmt.Errorf("LifecycleManager non initialisé")
	}

	// Récupérer les nœuds de la règle
	nodeIDs := rn.LifecycleManager.GetNodesForRule(ruleID)
	if len(nodeIDs) == 0 {
		return fmt.Errorf("règle %s non trouvée ou aucun nœud associé", ruleID)
	}

	// Vérifier qu'il y a au moins un TerminalNode
	hasTerminal := false
	for _, nodeID := range nodeIDs {
		lifecycle, exists := rn.LifecycleManager.GetNodeLifecycle(nodeID)
		if !exists {
			return fmt.Errorf("nœud %s de la règle %s non trouvé dans LifecycleManager", nodeID, ruleID)
		}

		if lifecycle.NodeType == "terminal" {
			hasTerminal = true
		}

		// Vérifier que le nœud existe dans la map appropriée
		switch lifecycle.NodeType {
		case "type":
			found := false
			for _, typeNode := range rn.TypeNodes {
				if typeNode.GetID() == nodeID {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("TypeNode %s de la règle %s non trouvé dans TypeNodes map", nodeID, ruleID)
			}
		case "alpha":
			if _, exists := rn.AlphaNodes[nodeID]; !exists {
				return fmt.Errorf("AlphaNode %s de la règle %s non trouvé dans AlphaNodes map", nodeID, ruleID)
			}
		case "join", "exists", "accumulate":
			if _, exists := rn.BetaNodes[nodeID]; !exists {
				return fmt.Errorf("BetaNode %s de la règle %s non trouvé dans BetaNodes map", nodeID, ruleID)
			}
		case "terminal":
			if _, exists := rn.TerminalNodes[nodeID]; !exists {
				return fmt.Errorf("TerminalNode %s de la règle %s non trouvé dans TerminalNodes map", nodeID, ruleID)
			}
		}
	}

	if !hasTerminal {
		return fmt.Errorf("règle %s n'a pas de TerminalNode", ruleID)
	}

	return nil
}

// ValidateFactIntegrity vérifie l'intégrité d'un fait dans le réseau
func (rn *ReteNetwork) ValidateFactIntegrity(factID string) error {
	// Vérifier que le fait existe dans le storage
	fact := rn.Storage.GetFact(factID)
	if fact == nil {
		return fmt.Errorf("fait %s non trouvé dans le storage", factID)
	}

	// Vérifier que le type du fait existe dans le réseau
	if _, exists := rn.TypeNodes[fact.Type]; !exists {
		return fmt.Errorf("type %s du fait %s non trouvé dans le réseau", fact.Type, factID)
	}

	// Vérifier que le fait est dans la mémoire du RootNode
	rootMemory := rn.RootNode.GetMemory()
	if _, exists := rootMemory.GetFact(factID); !exists {
		return fmt.Errorf("fait %s absent de la mémoire du RootNode", factID)
	}

	return nil
}

// ValidateMemoryConsistency vérifie la cohérence des mémoires entre les nœuds
func (rn *ReteNetwork) ValidateMemoryConsistency() error {
	// Tous les faits du RootNode doivent être dans le Storage
	rootMemory := rn.RootNode.GetMemory()
	for factID := range rootMemory.Facts {
		if rn.Storage.GetFact(factID) == nil {
			return fmt.Errorf("fait %s dans RootNode mais absent du Storage", factID)
		}
	}

	// Tous les faits des TypeNodes doivent être dans le RootNode
	for typeName, typeNode := range rn.TypeNodes {
		typeMemory := typeNode.GetMemory()
		for factID := range typeMemory.Facts {
			if _, exists := rootMemory.GetFact(factID); !exists {
				return fmt.Errorf("fait %s dans TypeNode %s mais absent du RootNode", factID, typeName)
			}
		}
	}

	return nil
}
