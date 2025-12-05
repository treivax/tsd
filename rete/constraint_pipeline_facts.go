// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

// collectExistingFacts parcourt tous les nœuds du réseau pour collecter les faits existants
func (cp *ConstraintPipeline) collectExistingFacts(network *ReteNetwork) []*Fact {
	factMap := make(map[string]*Fact)

	// Collecter depuis le RootNode
	if network.RootNode != nil && network.RootNode.Memory != nil {
		for _, fact := range network.RootNode.Memory.Facts {
			if fact != nil {
				factMap[fact.ID] = fact
			}
		}
	}

	// Collecter depuis les TypeNodes
	for _, typeNode := range network.TypeNodes {
		for _, token := range typeNode.Memory.Tokens {
			for _, fact := range token.Facts {
				if fact != nil {
					factMap[fact.ID] = fact
				}
			}
		}
	}

	// Collecter depuis les AlphaNodes
	for _, alphaNode := range network.AlphaNodes {
		for _, token := range alphaNode.Memory.Tokens {
			for _, fact := range token.Facts {
				if fact != nil {
					factMap[fact.ID] = fact
				}
			}
		}
	}

	// Collecter depuis les BetaNodes (JoinNodes, ExistsNodes, AccumulatorNodes, etc.)
	for _, betaNodeInterface := range network.BetaNodes {
		// Essayer de caster en JoinNode
		if joinNode, ok := betaNodeInterface.(*JoinNode); ok {
			// Mémoire gauche
			for _, token := range joinNode.LeftMemory.Tokens {
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
			// Mémoire droite
			for _, token := range joinNode.RightMemory.Tokens {
				for _, fact := range token.Facts {
					if fact != nil {
						factMap[fact.ID] = fact
					}
				}
			}
		}
		// Essayer de caster en ExistsNode
		if existsNode, ok := betaNodeInterface.(*ExistsNode); ok {
			for _, token := range existsNode.MainMemory.Tokens {
				for _, fact := range token.Facts {
					if fact != nil {
						factMap[fact.ID] = fact
					}
				}
			}
			for _, token := range existsNode.ExistsMemory.Tokens {
				for _, fact := range token.Facts {
					if fact != nil {
						factMap[fact.ID] = fact
					}
				}
			}
		}
		// Essayer de caster en AccumulatorNode
		if accNode, ok := betaNodeInterface.(*AccumulatorNode); ok {
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
	}

	// Convertir la map en slice
	facts := make([]*Fact, 0, len(factMap))
	for _, fact := range factMap {
		facts = append(facts, fact)
	}

	return facts
}

// organizeFactsByType organise les faits par type pour une propagation ciblée
func (cp *ConstraintPipeline) organizeFactsByType(facts []*Fact) map[string][]*Fact {
	factsByType := make(map[string][]*Fact)
	for _, fact := range facts {
		if fact != nil {
			factsByType[fact.Type] = append(factsByType[fact.Type], fact)
		}
	}
	return factsByType
}
