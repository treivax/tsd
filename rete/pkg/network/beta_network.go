package network

import (
	"fmt"
	"sync"

	"github.com/treivax/tsd/rete/pkg/domain"
	"github.com/treivax/tsd/rete/pkg/nodes"
)

// BetaNetworkBuilder construit des réseaux de nœuds Beta pour les jointures multi-faits.
type BetaNetworkBuilder struct {
	logger    domain.Logger
	betaNodes map[string]domain.BetaNode
	mutex     sync.RWMutex
}

// NewBetaNetworkBuilder crée un nouveau constructeur de réseau Beta.
func NewBetaNetworkBuilder(logger domain.Logger) *BetaNetworkBuilder {
	return &BetaNetworkBuilder{
		logger:    logger,
		betaNodes: make(map[string]domain.BetaNode),
	}
}

// CreateJoinNode crée un nouveau nœud de jointure avec conditions.
func (bnb *BetaNetworkBuilder) CreateJoinNode(id string, conditions []domain.JoinCondition) domain.JoinNode {
	bnb.mutex.Lock()
	defer bnb.mutex.Unlock()

	joinNode := nodes.NewJoinNode(id, bnb.logger)
	joinNode.SetJoinConditions(conditions)

	bnb.betaNodes[id] = joinNode

	bnb.logger.Info("created join node", map[string]interface{}{
		"node_id":    id,
		"conditions": len(conditions),
	})

	return joinNode
}

// CreateBetaNode crée un nœud Beta simple (sans conditions).
func (bnb *BetaNetworkBuilder) CreateBetaNode(id string) domain.BetaNode {
	bnb.mutex.Lock()
	defer bnb.mutex.Unlock()

	betaNode := nodes.NewBaseBetaNode(id, "BetaNode", bnb.logger)
	bnb.betaNodes[id] = betaNode

	bnb.logger.Info("created beta node", map[string]interface{}{
		"node_id": id,
	})

	return betaNode
}

// ConnectNodes connecte un nœud parent à un nœud enfant.
func (bnb *BetaNetworkBuilder) ConnectNodes(parentID, childID string) error {
	bnb.mutex.RLock()
	parent, parentExists := bnb.betaNodes[parentID]
	child, childExists := bnb.betaNodes[childID]
	bnb.mutex.RUnlock()

	if !parentExists {
		return fmt.Errorf("parent node not found: %s", parentID)
	}
	if !childExists {
		return fmt.Errorf("child node not found: %s", childID)
	}

	parent.AddChild(child)

	bnb.logger.Info("connected nodes", map[string]interface{}{
		"parent_id": parentID,
		"child_id":  childID,
	})

	return nil
}

// GetBetaNode retourne un nœud Beta par son ID.
func (bnb *BetaNetworkBuilder) GetBetaNode(id string) (domain.BetaNode, bool) {
	bnb.mutex.RLock()
	defer bnb.mutex.RUnlock()

	node, exists := bnb.betaNodes[id]
	return node, exists
}

// ListBetaNodes retourne tous les nœuds Beta créés.
func (bnb *BetaNetworkBuilder) ListBetaNodes() map[string]domain.BetaNode {
	bnb.mutex.RLock()
	defer bnb.mutex.RUnlock()

	// Retourner une copie pour éviter les modifications concurrentes
	result := make(map[string]domain.BetaNode, len(bnb.betaNodes))
	for id, node := range bnb.betaNodes {
		result[id] = node
	}
	return result
}

// ClearNetwork vide tous les nœuds du réseau.
func (bnb *BetaNetworkBuilder) ClearNetwork() {
	bnb.mutex.Lock()
	defer bnb.mutex.Unlock()

	bnb.betaNodes = make(map[string]domain.BetaNode)

	bnb.logger.Info("cleared beta network", map[string]interface{}{})
}

// MultiJoinPattern représente un pattern de jointures multiples pour construire des règles complexes.
type MultiJoinPattern struct {
	PatternID   string              `json:"pattern_id"`
	JoinSpecs   []JoinSpecification `json:"join_specs"`
	FinalAction string              `json:"final_action"`
}

// JoinSpecification spécifie une jointure dans un pattern multi-jointures.
type JoinSpecification struct {
	LeftType   string                 `json:"left_type"`  // Type des faits du côté gauche
	RightType  string                 `json:"right_type"` // Type des faits du côté droit
	Conditions []domain.JoinCondition `json:"conditions"` // Conditions de jointure
	NodeID     string                 `json:"node_id"`    // ID unique du nœud de jointure
}

// BuildMultiJoinNetwork construit un réseau de jointures multiples à partir d'un pattern.
func (bnb *BetaNetworkBuilder) BuildMultiJoinNetwork(pattern MultiJoinPattern) ([]domain.BetaNode, error) {
	createdNodes := make([]domain.BetaNode, 0, len(pattern.JoinSpecs))

	bnb.logger.Info("building multi-join network", map[string]interface{}{
		"pattern_id": pattern.PatternID,
		"join_count": len(pattern.JoinSpecs),
	})

	// Créer les nœuds de jointure en cascade
	var previousNode domain.BetaNode

	for i, spec := range pattern.JoinSpecs {
		nodeID := spec.NodeID
		if nodeID == "" {
			nodeID = fmt.Sprintf("%s_join_%d", pattern.PatternID, i)
		}

		// Créer le nœud de jointure avec ses conditions
		joinNode := bnb.CreateJoinNode(nodeID, spec.Conditions)
		createdNodes = append(createdNodes, joinNode)

		// Connecter au nœud précédent si il existe
		if previousNode != nil {
			previousNode.AddChild(joinNode)

			bnb.logger.Debug("connected join nodes", map[string]interface{}{
				"parent_id": previousNode.ID(),
				"child_id":  joinNode.ID(),
			})
		}

		previousNode = joinNode

		bnb.logger.Debug("created join node in pattern", map[string]interface{}{
			"pattern_id": pattern.PatternID,
			"node_id":    nodeID,
			"left_type":  spec.LeftType,
			"right_type": spec.RightType,
			"conditions": len(spec.Conditions),
		})
	}

	bnb.logger.Info("multi-join network built successfully", map[string]interface{}{
		"pattern_id":    pattern.PatternID,
		"nodes_created": len(createdNodes),
	})

	return createdNodes, nil
}

// NetworkStatistics retourne les statistiques du réseau Beta.
func (bnb *BetaNetworkBuilder) NetworkStatistics() NetworkStats {
	bnb.mutex.RLock()
	defer bnb.mutex.RUnlock()

	stats := NetworkStats{
		TotalNodes:      len(bnb.betaNodes),
		JoinNodes:       0,
		SimpleBetaNodes: 0,
		MemoryStats:     make(map[string]MemoryStat),
	}

	for nodeID, node := range bnb.betaNodes {
		// Compter le type de nœud
		if _, isJoinNode := node.(domain.JoinNode); isJoinNode {
			stats.JoinNodes++
		} else {
			stats.SimpleBetaNodes++
		}

		// Collecter les statistiques mémoire
		tokens, facts := 0, 0
		if leftMemory := node.GetLeftMemory(); leftMemory != nil {
			tokens = len(leftMemory)
		}
		if rightMemory := node.GetRightMemory(); rightMemory != nil {
			facts = len(rightMemory)
		}

		stats.MemoryStats[nodeID] = MemoryStat{
			TokenCount: tokens,
			FactCount:  facts,
		}

		stats.TotalTokens += tokens
		stats.TotalFacts += facts
	}

	return stats
}

// NetworkStats représente les statistiques du réseau Beta.
type NetworkStats struct {
	TotalNodes      int                   `json:"total_nodes"`
	JoinNodes       int                   `json:"join_nodes"`
	SimpleBetaNodes int                   `json:"simple_beta_nodes"`
	TotalTokens     int                   `json:"total_tokens"`
	TotalFacts      int                   `json:"total_facts"`
	MemoryStats     map[string]MemoryStat `json:"memory_stats"`
}

// MemoryStat représente les statistiques mémoire d'un nœud.
type MemoryStat struct {
	TokenCount int `json:"token_count"`
	FactCount  int `json:"fact_count"`
}
