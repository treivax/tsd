// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"time"
)

// validateBuildChainInputs valide les paramètres d'entrée pour la construction d'une chaîne alpha
func validateBuildChainInputs(conditions []SimpleCondition, parentNode Node, network *ReteNetwork) error {
	if len(conditions) == 0 {
		return fmt.Errorf("impossible de construire une chaîne sans conditions")
	}

	if parentNode == nil {
		return fmt.Errorf("le nœud parent ne peut pas être nil")
	}

	return nil
}

// chainBuildMetrics contient les métriques de construction d'une chaîne alpha
type chainBuildMetrics struct {
	startTime       time.Time
	nodesCreated    int
	nodesReused     int
	hashesGenerated []string
}

// initializeChainMetrics initialise les métriques pour la construction d'une chaîne
func initializeChainMetrics(conditionsCount int) *chainBuildMetrics {
	return &chainBuildMetrics{
		startTime:       time.Now(),
		nodesCreated:    0,
		nodesReused:     0,
		hashesGenerated: make([]string, 0, conditionsCount),
	}
}

// alphaNodeBuildResult contient le résultat de la construction d'un nœud alpha
type alphaNodeBuildResult struct {
	node   *AlphaNode
	hash   string
	reused bool
}

// buildAndConnectAlphaNode construit ou réutilise un nœud alpha et le connecte à son parent
func (acb *AlphaChainBuilder) buildAndConnectAlphaNode(
	condition SimpleCondition,
	variableName string,
	currentParent Node,
	ruleID string,
	conditionIndex int,
	totalConditions int,
	metrics *chainBuildMetrics,
) (*alphaNodeBuildResult, error) {
	// Convertir SimpleCondition en map pour la condition du nœud alpha
	conditionMap := map[string]interface{}{
		"type":     condition.Type,
		"left":     condition.Left,
		"operator": condition.Operator,
		"right":    condition.Right,
	}

	// Obtenir ou créer l'AlphaNode via le gestionnaire de partage
	alphaNode, hash, reused, err := acb.network.AlphaSharingManager.GetOrCreateAlphaNode(
		conditionMap,
		variableName,
		acb.storage,
	)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la création/récupération du nœud alpha %d: %w", conditionIndex, err)
	}

	// Mettre à jour les métriques
	metrics.hashesGenerated = append(metrics.hashesGenerated, hash)
	if reused {
		metrics.nodesReused++
	} else {
		metrics.nodesCreated++
	}

	// Gérer la connexion du nœud
	if reused {
		acb.handleReusedNodeConnection(alphaNode, currentParent, ruleID, conditionIndex, totalConditions)
	} else {
		acb.handleNewNodeConnection(alphaNode, currentParent, ruleID, conditionIndex, totalConditions)
	}

	// Enregistrer le nœud dans le LifecycleManager avec la règle
	lifecycle := acb.network.LifecycleManager.RegisterNode(alphaNode.ID, "alpha")
	lifecycle.AddRuleReference(ruleID, "") // RuleName peut être ajouté plus tard si nécessaire

	if reused {
		fmt.Printf("📊 [AlphaChainBuilder] Nœud %s maintenant utilisé par %d règle(s)\n",
			alphaNode.ID, lifecycle.GetRefCount())
	}

	return &alphaNodeBuildResult{
		node:   alphaNode,
		hash:   hash,
		reused: reused,
	}, nil
}

// handleReusedNodeConnection gère la connexion d'un nœud alpha réutilisé
func (acb *AlphaChainBuilder) handleReusedNodeConnection(
	alphaNode *AlphaNode,
	currentParent Node,
	ruleID string,
	conditionIndex int,
	totalConditions int,
) {
	fmt.Printf("♻️  [AlphaChainBuilder] Réutilisation du nœud alpha %s pour la règle %s (condition %d/%d)\n",
		alphaNode.ID, ruleID, conditionIndex+1, totalConditions)

	if !acb.isAlreadyConnectedCached(currentParent, alphaNode) {
		// Connecter au parent si pas déjà connecté
		currentParent.AddChild(alphaNode)
		fmt.Printf("🔗 [AlphaChainBuilder] Connexion du nœud réutilisé %s au parent %s\n",
			alphaNode.ID, currentParent.GetID())
	} else {
		fmt.Printf("✓  [AlphaChainBuilder] Nœud %s déjà connecté au parent %s\n",
			alphaNode.ID, currentParent.GetID())
	}
}

// handleNewNodeConnection gère la connexion d'un nouveau nœud alpha
func (acb *AlphaChainBuilder) handleNewNodeConnection(
	alphaNode *AlphaNode,
	currentParent Node,
	ruleID string,
	conditionIndex int,
	totalConditions int,
) {
	// Connecter au parent et l'ajouter au réseau
	currentParent.AddChild(alphaNode)
	acb.network.AlphaNodes[alphaNode.ID] = alphaNode

	// Mettre à jour le cache de connexion
	acb.updateConnectionCache(currentParent.GetID(), alphaNode.ID, true)

	fmt.Printf("🆕 [AlphaChainBuilder] Nouveau nœud alpha %s créé pour la règle %s (condition %d/%d)\n",
		alphaNode.ID, ruleID, conditionIndex+1, totalConditions)
	fmt.Printf("🔗 [AlphaChainBuilder] Connexion du nœud %s au parent %s\n",
		alphaNode.ID, currentParent.GetID())
}

// recordChainMetrics enregistre les métriques finales de construction de chaîne
func (acb *AlphaChainBuilder) recordChainMetrics(
	ruleID string,
	chain *AlphaChain,
	metrics *chainBuildMetrics,
) {
	fmt.Printf("✅ [AlphaChainBuilder] Chaîne alpha complète construite pour la règle %s: %d nœud(s)\n",
		ruleID, len(chain.Nodes))

	// Enregistrer les métriques
	if acb.metrics != nil {
		buildTime := time.Since(metrics.startTime)
		detail := ChainMetricDetail{
			RuleID:          ruleID,
			ChainLength:     len(chain.Nodes),
			NodesCreated:    metrics.nodesCreated,
			NodesReused:     metrics.nodesReused,
			BuildTime:       buildTime,
			Timestamp:       time.Now(),
			HashesGenerated: metrics.hashesGenerated,
		}
		acb.metrics.RecordChainBuild(detail)
	}
}
