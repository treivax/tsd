// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"fmt"
	"sync"
	"time"
)

// NodeReference représente une référence à un nœud RETE.
//
// Cette structure stocke les informations nécessaires pour identifier
// et accéder à un nœud sans stocker directement le pointeur (évite cycles).
type NodeReference struct {
	// NodeID est l'identifiant unique du nœud
	NodeID string

	// NodeType indique le type de nœud ("alpha", "beta", "terminal")
	NodeType string

	// FactType est le type de fait concerné (ex: "Product")
	FactType string

	// Fields est la liste des champs utilisés par ce nœud
	Fields []string
}

// String retourne une représentation string de la référence
func (nr NodeReference) String() string {
	return fmt.Sprintf("%s[%s](%s)", nr.NodeType, nr.NodeID, nr.FactType)
}

// DependencyIndex est un index inversé permettant de trouver rapidement
// tous les nœuds RETE sensibles à un champ donné d'un type donné.
//
// Structure : factType → field → [nodeRefs]
//
// Exemple :
//
//	index.Get("Product", "price") → [alpha1, alpha2, beta3, terminal5]
//
// Thread-safety : toutes les opérations sont protégées par mutex.
type DependencyIndex struct {
	// alphaIndex indexe les nœuds alpha par (factType, field)
	// Structure : factType → field → [nodeIDs]
	alphaIndex map[string]map[string][]string

	// betaIndex indexe les nœuds beta par (factType, field)
	betaIndex map[string]map[string][]string

	// terminalIndex indexe les nœuds terminaux par (factType, field)
	terminalIndex map[string]map[string][]string

	// nodeReferences stocke les détails de chaque nœud indexé
	// Structure : nodeID → NodeReference
	nodeReferences map[string]NodeReference

	// metadata stocke des informations sur l'index
	builtAt    time.Time
	nodeCount  int
	fieldCount int

	// mutex protège l'accès concurrent
	mutex sync.RWMutex
}

// NewDependencyIndex crée un nouvel index de dépendances vide.
func NewDependencyIndex() *DependencyIndex {
	return &DependencyIndex{
		alphaIndex:     make(map[string]map[string][]string),
		betaIndex:      make(map[string]map[string][]string),
		terminalIndex:  make(map[string]map[string][]string),
		nodeReferences: make(map[string]NodeReference),
		builtAt:        time.Now(),
	}
}

// AddAlphaNode enregistre un nœud alpha pour un ensemble de champs.
//
// Paramètres :
//   - nodeID : identifiant unique du nœud
//   - factType : type de fait concerné
//   - fields : liste des champs testés par ce nœud
func (idx *DependencyIndex) AddAlphaNode(nodeID, factType string, fields []string) {
	idx.mutex.Lock()
	defer idx.mutex.Unlock()

	idx.addNodeToIndex(idx.alphaIndex, nodeID, factType, fields)
	idx.nodeReferences[nodeID] = NodeReference{
		NodeID:   nodeID,
		NodeType: NodeTypeAlpha,
		FactType: factType,
		Fields:   fields,
	}
	idx.nodeCount++
}

// AddBetaNode enregistre un nœud beta pour un ensemble de champs.
//
// Paramètres :
//   - nodeID : identifiant unique du nœud
//   - factType : type de fait concerné
//   - fields : liste des champs utilisés dans les jointures
func (idx *DependencyIndex) AddBetaNode(nodeID, factType string, fields []string) {
	idx.mutex.Lock()
	defer idx.mutex.Unlock()

	idx.addNodeToIndex(idx.betaIndex, nodeID, factType, fields)
	idx.nodeReferences[nodeID] = NodeReference{
		NodeID:   nodeID,
		NodeType: NodeTypeBeta,
		FactType: factType,
		Fields:   fields,
	}
	idx.nodeCount++
}

// AddTerminalNode enregistre un nœud terminal pour un ensemble de champs.
//
// Paramètres :
//   - nodeID : identifiant unique du nœud
//   - factType : type de fait concerné
//   - fields : liste des champs utilisés dans les actions
func (idx *DependencyIndex) AddTerminalNode(nodeID, factType string, fields []string) {
	idx.mutex.Lock()
	defer idx.mutex.Unlock()

	idx.addNodeToIndex(idx.terminalIndex, nodeID, factType, fields)
	idx.nodeReferences[nodeID] = NodeReference{
		NodeID:   nodeID,
		NodeType: NodeTypeTerminal,
		FactType: factType,
		Fields:   fields,
	}
	idx.nodeCount++
}

// addNodeToIndex est une fonction helper privée pour ajouter un nœud à un index.
// ATTENTION : doit être appelée avec mutex déjà acquis.
func (idx *DependencyIndex) addNodeToIndex(
	index map[string]map[string][]string,
	nodeID, factType string,
	fields []string,
) {
	// Initialiser la map pour ce factType si nécessaire
	if index[factType] == nil {
		index[factType] = make(map[string][]string)
	}

	// Ajouter le nœud pour chaque champ
	for _, field := range fields {
		// Vérifier si le nœud n'est pas déjà indexé pour ce champ
		nodes := index[factType][field]
		alreadyIndexed := false
		for _, existingNodeID := range nodes {
			if existingNodeID == nodeID {
				alreadyIndexed = true
				break
			}
		}

		if !alreadyIndexed {
			index[factType][field] = append(index[factType][field], nodeID)
			idx.fieldCount++
		}
	}
}

// GetAffectedNodes retourne tous les nœuds affectés par un changement de champ.
//
// Paramètres :
//   - factType : type de fait
//   - field : nom du champ modifié
//
// Retourne la liste des NodeReferences des nœuds sensibles à ce champ.
func (idx *DependencyIndex) GetAffectedNodes(factType, field string) []NodeReference {
	idx.mutex.RLock()
	defer idx.mutex.RUnlock()

	affectedNodeIDs := make(map[string]bool)

	// Collecter les nœuds alpha
	if alphaFields := idx.alphaIndex[factType]; alphaFields != nil {
		if nodes := alphaFields[field]; nodes != nil {
			for _, nodeID := range nodes {
				affectedNodeIDs[nodeID] = true
			}
		}
	}

	// Collecter les nœuds beta
	if betaFields := idx.betaIndex[factType]; betaFields != nil {
		if nodes := betaFields[field]; nodes != nil {
			for _, nodeID := range nodes {
				affectedNodeIDs[nodeID] = true
			}
		}
	}

	// Collecter les nœuds terminaux
	if terminalFields := idx.terminalIndex[factType]; terminalFields != nil {
		if nodes := terminalFields[field]; nodes != nil {
			for _, nodeID := range nodes {
				affectedNodeIDs[nodeID] = true
			}
		}
	}

	// Convertir en NodeReferences
	result := make([]NodeReference, 0, len(affectedNodeIDs))
	for nodeID := range affectedNodeIDs {
		if ref, exists := idx.nodeReferences[nodeID]; exists {
			result = append(result, ref)
		}
	}

	return result
}

// GetAffectedNodesForDelta retourne tous les nœuds affectés par un FactDelta.
//
// Paramètres :
//   - delta : le FactDelta contenant les champs modifiés
//
// Retourne la liste des NodeReferences des nœuds affectés (dédupliqués).
func (idx *DependencyIndex) GetAffectedNodesForDelta(delta *FactDelta) []NodeReference {
	idx.mutex.RLock()
	defer idx.mutex.RUnlock()

	affectedNodeIDs := make(map[string]bool)

	// Pour chaque champ modifié
	for fieldName := range delta.Fields {
		// Alpha nodes
		if alphaFields := idx.alphaIndex[delta.FactType]; alphaFields != nil {
			if nodes := alphaFields[fieldName]; nodes != nil {
				for _, nodeID := range nodes {
					affectedNodeIDs[nodeID] = true
				}
			}
		}

		// Beta nodes
		if betaFields := idx.betaIndex[delta.FactType]; betaFields != nil {
			if nodes := betaFields[fieldName]; nodes != nil {
				for _, nodeID := range nodes {
					affectedNodeIDs[nodeID] = true
				}
			}
		}

		// Terminal nodes
		if terminalFields := idx.terminalIndex[delta.FactType]; terminalFields != nil {
			if nodes := terminalFields[fieldName]; nodes != nil {
				for _, nodeID := range nodes {
					affectedNodeIDs[nodeID] = true
				}
			}
		}
	}

	// Convertir en NodeReferences
	result := make([]NodeReference, 0, len(affectedNodeIDs))
	for nodeID := range affectedNodeIDs {
		if ref, exists := idx.nodeReferences[nodeID]; exists {
			result = append(result, ref)
		}
	}

	return result
}

// Clear vide complètement l'index.
func (idx *DependencyIndex) Clear() {
	idx.mutex.Lock()
	defer idx.mutex.Unlock()

	idx.alphaIndex = make(map[string]map[string][]string)
	idx.betaIndex = make(map[string]map[string][]string)
	idx.terminalIndex = make(map[string]map[string][]string)
	idx.nodeReferences = make(map[string]NodeReference)
	idx.nodeCount = 0
	idx.fieldCount = 0
	idx.builtAt = time.Now()
}

// IndexStats retourne des statistiques sur l'index.
type IndexStats struct {
	NodeCount      int
	FieldCount     int
	AlphaNodeCount int
	BetaNodeCount  int
	TerminalCount  int
	FactTypes      []string
	BuiltAt        time.Time
	MemoryEstimate int64 // Estimation mémoire en bytes
}

// GetStats retourne les statistiques de l'index.
func (idx *DependencyIndex) GetStats() IndexStats {
	idx.mutex.RLock()
	defer idx.mutex.RUnlock()

	stats := IndexStats{
		NodeCount:  idx.nodeCount,
		FieldCount: idx.fieldCount,
		BuiltAt:    idx.builtAt,
	}

	// Compter par type de nœud
	for _, ref := range idx.nodeReferences {
		switch ref.NodeType {
		case "alpha":
			stats.AlphaNodeCount++
		case "beta":
			stats.BetaNodeCount++
		case "terminal":
			stats.TerminalCount++
		}
	}

	// Collecter les types de faits
	factTypeSet := make(map[string]bool)
	for factType := range idx.alphaIndex {
		factTypeSet[factType] = true
	}
	for factType := range idx.betaIndex {
		factTypeSet[factType] = true
	}
	for factType := range idx.terminalIndex {
		factTypeSet[factType] = true
	}

	stats.FactTypes = make([]string, 0, len(factTypeSet))
	for factType := range factTypeSet {
		stats.FactTypes = append(stats.FactTypes, factType)
	}

	// Estimation mémoire approximative
	const (
		nodeIDSize      = 50
		refSize         = 100
		sliceOverhead   = 24
		entrySize       = 74
		estimatePerNode = nodeIDSize + refSize + sliceOverhead
	)
	stats.MemoryEstimate = int64(idx.nodeCount * estimatePerNode)
	stats.MemoryEstimate += int64(idx.fieldCount * entrySize)

	return stats
}

// GetTotalNodeCount retourne le nombre total de nœuds indexés.
func (idx *DependencyIndex) GetTotalNodeCount() int {
	idx.mutex.RLock()
	defer idx.mutex.RUnlock()
	return idx.nodeCount
}

// String retourne une représentation string de l'index.
func (idx *DependencyIndex) String() string {
	stats := idx.GetStats()
	return fmt.Sprintf(
		"DependencyIndex[nodes=%d, fields=%d, alpha=%d, beta=%d, terminal=%d, types=%d]",
		stats.NodeCount, stats.FieldCount,
		stats.AlphaNodeCount, stats.BetaNodeCount, stats.TerminalCount,
		len(stats.FactTypes),
	)
}
