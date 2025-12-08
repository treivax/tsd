// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

// fact_collection_helpers.go contient des fonctions helper pour la collection de faits
// depuis différents types de nœuds du réseau RETE.
// Ces fonctions ont été extraites de collectExistingFacts() pour réduire la complexité.

// collectFactsFromRootNode collecte les faits depuis le RootNode
func collectFactsFromRootNode(network *ReteNetwork, factMap map[string]*Fact) {
	if network.RootNode == nil || network.RootNode.Memory == nil {
		return
	}

	for _, fact := range network.RootNode.Memory.Facts {
		if fact != nil {
			factMap[fact.ID] = fact
		}
	}
}

// collectFactsFromTypeNodes collecte les faits depuis tous les TypeNodes
func collectFactsFromTypeNodes(network *ReteNetwork, factMap map[string]*Fact) {
	for _, typeNode := range network.TypeNodes {
		if typeNode == nil || typeNode.Memory == nil {
			continue
		}

		for _, token := range typeNode.Memory.Tokens {
			if token == nil {
				continue
			}
			for _, fact := range token.Facts {
				if fact != nil {
					factMap[fact.ID] = fact
				}
			}
		}
	}
}

// collectFactsFromAlphaNodes collecte les faits depuis tous les AlphaNodes
func collectFactsFromAlphaNodes(network *ReteNetwork, factMap map[string]*Fact) {
	for _, alphaNode := range network.AlphaNodes {
		if alphaNode == nil || alphaNode.Memory == nil {
			continue
		}

		for _, token := range alphaNode.Memory.Tokens {
			if token == nil {
				continue
			}
			for _, fact := range token.Facts {
				if fact != nil {
					factMap[fact.ID] = fact
				}
			}
		}
	}
}

// collectFactsFromJoinNode collecte les faits depuis un JoinNode
func collectFactsFromJoinNode(joinNode *JoinNode, factMap map[string]*Fact) {
	// Mémoire gauche
	if joinNode.LeftMemory != nil {
		for _, token := range joinNode.LeftMemory.Tokens {
			if token == nil {
				continue
			}

			// Faits dans le token actuel
			for _, fact := range token.Facts {
				if fact != nil {
					factMap[fact.ID] = fact
				}
			}

			// Collecter aussi les faits des parents dans les tokens de jointure
			for parent := token.Parent; parent != nil; parent = parent.Parent {
				for _, fact := range parent.Facts {
					if fact != nil {
						factMap[fact.ID] = fact
					}
				}
			}
		}
	}

	// Mémoire droite
	if joinNode.RightMemory != nil {
		for _, token := range joinNode.RightMemory.Tokens {
			if token == nil {
				continue
			}
			for _, fact := range token.Facts {
				if fact != nil {
					factMap[fact.ID] = fact
				}
			}
		}
	}
}

// collectFactsFromExistsNode collecte les faits depuis un ExistsNode
func collectFactsFromExistsNode(existsNode *ExistsNode, factMap map[string]*Fact) {
	// Mémoire principale
	if existsNode.MainMemory != nil {
		for _, token := range existsNode.MainMemory.Tokens {
			if token == nil {
				continue
			}
			for _, fact := range token.Facts {
				if fact != nil {
					factMap[fact.ID] = fact
				}
			}
		}
	}

	// Mémoire exists
	if existsNode.ExistsMemory != nil {
		for _, token := range existsNode.ExistsMemory.Tokens {
			if token == nil {
				continue
			}
			for _, fact := range token.Facts {
				if fact != nil {
					factMap[fact.ID] = fact
				}
			}
		}
	}
}

// collectFactsFromAccumulatorNode collecte les faits depuis un AccumulatorNode
func collectFactsFromAccumulatorNode(accNode *AccumulatorNode, factMap map[string]*Fact) {
	// Collecter depuis MainFacts
	for _, fact := range accNode.MainFacts {
		if fact != nil {
			factMap[fact.ID] = fact
		}
	}

	// Collecter depuis AllFacts
	for _, fact := range accNode.AllFacts {
		if fact != nil {
			factMap[fact.ID] = fact
		}
	}
}

// collectFactsFromBetaNodes collecte les faits depuis tous les BetaNodes
func collectFactsFromBetaNodes(network *ReteNetwork, factMap map[string]*Fact) {
	for _, betaNodeInterface := range network.BetaNodes {
		if betaNodeInterface == nil {
			continue
		}

		// Essayer de caster en JoinNode
		if joinNode, ok := betaNodeInterface.(*JoinNode); ok {
			collectFactsFromJoinNode(joinNode, factMap)
			continue
		}

		// Essayer de caster en ExistsNode
		if existsNode, ok := betaNodeInterface.(*ExistsNode); ok {
			collectFactsFromExistsNode(existsNode, factMap)
			continue
		}

		// Essayer de caster en AccumulatorNode
		if accNode, ok := betaNodeInterface.(*AccumulatorNode); ok {
			collectFactsFromAccumulatorNode(accNode, factMap)
			continue
		}
	}
}

// convertFactMapToSlice convertit une map de faits en slice
func convertFactMapToSlice(factMap map[string]*Fact) []*Fact {
	facts := make([]*Fact, 0, len(factMap))
	for _, fact := range factMap {
		facts = append(facts, fact)
	}
	return facts
}
