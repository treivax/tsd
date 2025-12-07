// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

// collectExistingFacts parcourt tous les nœuds du réseau pour collecter les faits existants.
// Cette fonction agit comme orchestrateur, déléguant la collection à des fonctions spécialisées
// pour réduire la complexité et améliorer la testabilité.
func (cp *ConstraintPipeline) collectExistingFacts(network *ReteNetwork) []*Fact {
	factMap := make(map[string]*Fact)

	// Étape 1: Collecter depuis le RootNode
	collectFactsFromRootNode(network, factMap)

	// Étape 2: Collecter depuis les TypeNodes
	collectFactsFromTypeNodes(network, factMap)

	// Étape 3: Collecter depuis les AlphaNodes
	collectFactsFromAlphaNodes(network, factMap)

	// Étape 4: Collecter depuis les BetaNodes (JoinNodes, ExistsNodes, AccumulatorNodes)
	collectFactsFromBetaNodes(network, factMap)

	// Étape 5: Convertir la map en slice
	return convertFactMapToSlice(factMap)
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
