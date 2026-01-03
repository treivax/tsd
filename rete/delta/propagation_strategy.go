// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

// PropagationStrategy définit une stratégie de propagation.
//
// Cette interface permet d'implémenter différentes stratégies
// de propagation pour s'adapter à différents scénarios.
type PropagationStrategy interface {
	// GetName retourne le nom de la stratégie
	GetName() string

	// ShouldPropagate détermine si la propagation doit avoir lieu
	ShouldPropagate(delta *FactDelta, affectedNodes []NodeReference) bool

	// GetPropagationOrder retourne l'ordre de propagation des nœuds
	GetPropagationOrder(nodes []NodeReference) []NodeReference
}

// SequentialStrategy propage vers les nœuds dans l'ordre séquentiel.
//
// Cette stratégie est simple et prévisible : alpha → beta → terminal.
type SequentialStrategy struct{}

// GetName retourne "Sequential"
func (s *SequentialStrategy) GetName() string {
	return "Sequential"
}

// ShouldPropagate retourne toujours true (propage toujours)
func (s *SequentialStrategy) ShouldPropagate(delta *FactDelta, affectedNodes []NodeReference) bool {
	return len(affectedNodes) > 0
}

// GetPropagationOrder trie les nœuds par type : alpha, puis beta, puis terminal
func (s *SequentialStrategy) GetPropagationOrder(nodes []NodeReference) []NodeReference {
	var alphaNodes, betaNodes, terminalNodes []NodeReference

	for _, node := range nodes {
		switch node.NodeType {
		case NodeTypeAlpha:
			alphaNodes = append(alphaNodes, node)
		case NodeTypeBeta:
			betaNodes = append(betaNodes, node)
		case NodeTypeTerminal:
			terminalNodes = append(terminalNodes, node)
		}
	}

	ordered := make([]NodeReference, 0, len(nodes))
	ordered = append(ordered, alphaNodes...)
	ordered = append(ordered, betaNodes...)
	ordered = append(ordered, terminalNodes...)

	return ordered
}

// TopologicalStrategy propage en respectant les dépendances topologiques.
//
// Cette stratégie garantit qu'un nœud parent est toujours traité avant
// ses nœuds enfants (ordre topologique du graphe RETE).
type TopologicalStrategy struct {
	nodeDepths map[string]int
}

// NewTopologicalStrategy crée une nouvelle stratégie topologique
func NewTopologicalStrategy() *TopologicalStrategy {
	return &TopologicalStrategy{
		nodeDepths: make(map[string]int),
	}
}

// GetName retourne "Topological"
func (ts *TopologicalStrategy) GetName() string {
	return "Topological"
}

// ShouldPropagate retourne true si au moins un nœud est affecté
func (ts *TopologicalStrategy) ShouldPropagate(delta *FactDelta, affectedNodes []NodeReference) bool {
	return len(affectedNodes) > 0
}

// GetPropagationOrder trie les nœuds par profondeur topologique
func (ts *TopologicalStrategy) GetPropagationOrder(nodes []NodeReference) []NodeReference {
	if len(ts.nodeDepths) == 0 {
		sequential := &SequentialStrategy{}
		return sequential.GetPropagationOrder(nodes)
	}

	ordered := make([]NodeReference, len(nodes))
	copy(ordered, nodes)

	for i := 1; i < len(ordered); i++ {
		key := ordered[i]
		keyDepth := ts.getDepth(key.NodeID)
		j := i - 1

		for j >= 0 && ts.getDepth(ordered[j].NodeID) > keyDepth {
			ordered[j+1] = ordered[j]
			j--
		}
		ordered[j+1] = key
	}

	return ordered
}

// SetNodeDepth enregistre la profondeur d'un nœud
func (ts *TopologicalStrategy) SetNodeDepth(nodeID string, depth int) {
	ts.nodeDepths[nodeID] = depth
}

// getDepth retourne la profondeur d'un nœud (0 si inconnu)
func (ts *TopologicalStrategy) getDepth(nodeID string) int {
	if depth, exists := ts.nodeDepths[nodeID]; exists {
		return depth
	}
	return 0
}

// OptimizedStrategy est une stratégie hybride qui optimise selon le contexte.
//
// Elle combine plusieurs heuristiques :
// - Trier par type (alpha → beta → terminal)
// - Grouper par factType (meilleure localité cache)
// - Prioriser les nœuds avec moins de dépendances
type OptimizedStrategy struct{}

// GetName retourne "Optimized"
func (os *OptimizedStrategy) GetName() string {
	return "Optimized"
}

// ShouldPropagate retourne true si propagation justifiée
func (os *OptimizedStrategy) ShouldPropagate(delta *FactDelta, affectedNodes []NodeReference) bool {
	if len(affectedNodes) == 0 {
		return false
	}

	if delta.IsEmpty() {
		return false
	}

	return true
}

// GetPropagationOrder optimise l'ordre de propagation
func (os *OptimizedStrategy) GetPropagationOrder(nodes []NodeReference) []NodeReference {
	if len(nodes) == 0 {
		return nodes
	}

	// Grouper par type et factType pour meilleure localité cache
	groups := groupNodesByTypeAndFactType(nodes)

	// Ordre de propagation : alpha → beta → terminal
	ordered := make([]NodeReference, 0, len(nodes))
	ordered = appendGroupsWithPrefix(ordered, groups, NodeTypeAlpha)
	ordered = appendGroupsWithPrefix(ordered, groups, NodeTypeBeta)
	ordered = appendGroupsWithPrefix(ordered, groups, NodeTypeTerminal)

	return ordered
}

// groupNodesByTypeAndFactType regroupe les nœuds par clé "type:factType"
func groupNodesByTypeAndFactType(nodes []NodeReference) map[string][]NodeReference {
	groups := make(map[string][]NodeReference)
	for _, node := range nodes {
		key := node.NodeType + KeySeparator + node.FactType
		groups[key] = append(groups[key], node)
	}
	return groups
}

// appendGroupsWithPrefix ajoute tous les groupes commençant par prefix
func appendGroupsWithPrefix(
	ordered []NodeReference,
	groups map[string][]NodeReference,
	prefix string,
) []NodeReference {
	prefixWithSep := prefix + KeySeparator
	for key, group := range groups {
		if len(key) >= len(prefixWithSep) && key[:len(prefixWithSep)] == prefixWithSep {
			ordered = append(ordered, group...)
		}
	}
	return ordered
}
