// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
)

// identifyNewTerminals identifie les nœuds terminaux qui viennent d'être ajoutés
func (cp *ConstraintPipeline) identifyNewTerminals(network *ReteNetwork, existingTerminals map[string]bool) []*TerminalNode {
	var newTerminals []*TerminalNode
	for terminalID, terminal := range network.TerminalNodes {
		if !existingTerminals[terminalID] {
			newTerminals = append(newTerminals, terminal)
		}
	}
	return newTerminals
}

// propagateToNewTerminals propage les faits existants uniquement vers les nouvelles chaînes de règles
func (cp *ConstraintPipeline) propagateToNewTerminals(
	network *ReteNetwork,
	newTerminals []*TerminalNode,
	factsByType map[string][]*Fact,
) int {
	propagatedCount := 0

	// Pour chaque nouveau terminal, identifier les types de faits qu'il attend
	for _, terminal := range newTerminals {
		// Identifier les types de faits attendus par cette règle
		expectedTypes := cp.identifyExpectedTypesForTerminal(network, terminal)

		// Propager uniquement les faits des types attendus
		for _, typeName := range expectedTypes {
			if facts, exists := factsByType[typeName]; exists {
				for _, fact := range facts {
					// Propager le fait via le TypeNode correspondant
					if typeNode, exists := network.TypeNodes[typeName]; exists {
						// Créer un token pour ce fait
						token := &Token{
							ID:     fmt.Sprintf("retro_%s_%s", typeName, fact.ID),
							NodeID: typeNode.GetID(),
							Facts:  []*Fact{fact},
						}

						// Propager aux enfants du TypeNode
						err := typeNode.PropagateToChildren(fact, token)
						if err == nil {
							propagatedCount++
						}
					}
				}
			}
		}
	}

	return propagatedCount
}

// identifyExpectedTypesForTerminal identifie les types de faits attendus par un terminal
func (cp *ConstraintPipeline) identifyExpectedTypesForTerminal(network *ReteNetwork, terminal *TerminalNode) []string {
	expectedTypes := make(map[string]bool)

	// Parcourir les TypeNodes pour trouver ceux qui ont ce terminal comme descendant
	for typeName, typeNode := range network.TypeNodes {
		if cp.isTerminalReachableFrom(typeNode, terminal.GetID()) {
			expectedTypes[typeName] = true
		}
	}

	// Convertir en slice
	types := make([]string, 0, len(expectedTypes))
	for typeName := range expectedTypes {
		types = append(types, typeName)
	}

	return types
}

// isTerminalReachableFrom vérifie si un terminal est accessible depuis un nœud donné
func (cp *ConstraintPipeline) isTerminalReachableFrom(node Node, terminalID string) bool {
	// Vérification directe
	if node.GetID() == terminalID {
		return true
	}

	// Vérification récursive dans les enfants
	for _, child := range node.GetChildren() {
		if cp.isTerminalReachableFrom(child, terminalID) {
			return true
		}
	}

	return false
}

// processRuleRemovals traite les commandes de suppression de règles
func (cp *ConstraintPipeline) processRuleRemovals(network *ReteNetwork, resultMap map[string]interface{}) error {
	// Vérifier si des suppressions de règles sont présentes
	ruleRemovalsData, exists := resultMap["ruleRemovals"]
	if !exists {
		return nil // Pas de suppressions de règles
	}

	ruleRemovals, ok := ruleRemovalsData.([]interface{})
	if !ok || len(ruleRemovals) == 0 {
		return nil // Pas de suppressions de règles
	}

	cp.GetLogger().Info("🗑️  Traitement de %d suppression(s) de règles", len(ruleRemovals))

	// Traiter chaque suppression de règle
	for _, removalData := range ruleRemovals {
		removalMap, ok := removalData.(map[string]interface{})
		if !ok {
			cp.GetLogger().Warn("⚠️  Format de suppression invalide: %v", removalData)
			continue
		}

		ruleID, ok := removalMap["ruleID"].(string)
		if !ok || ruleID == "" {
			cp.GetLogger().Warn("⚠️  Identifiant de règle manquant ou invalide: %v", removalMap)
			continue
		}

		// Supprimer la règle du réseau
		cp.GetLogger().Info("🗑️  Suppression de la règle: %s", ruleID)
		err := network.RemoveRule(ruleID)
		if err != nil {
			// Logger l'erreur mais continuer avec les autres suppressions
			cp.GetLogger().Warn("⚠️  Erreur lors de la suppression de la règle %s: %v", ruleID, err)
			continue
		}

		cp.GetLogger().Info("✅ Règle %s supprimée avec succès", ruleID)
	}

	return nil
}
