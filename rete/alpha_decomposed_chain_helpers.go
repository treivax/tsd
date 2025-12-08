// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"time"
)

// alpha_decomposed_chain_helpers.go contient des fonctions helper pour la construction
// de chaînes alpha décomposées. Ces fonctions ont été extraites de BuildDecomposedChain()
// pour améliorer la lisibilité et la maintenabilité.

// DecomposedChainBuildContext contient le contexte de construction d'une chaîne décomposée
type DecomposedChainBuildContext struct {
	StartTime       time.Time
	NodesCreated    int
	NodesReused     int
	HashesGenerated []string
	Chain           *AlphaChain
	CurrentParent   Node
}

// validateBuildDecomposedInputs valide les paramètres d'entrée pour BuildDecomposedChain
func validateBuildDecomposedInputs(
	conditions []DecomposedCondition,
	parentNode Node,
	network *ReteNetwork,
) error {
	if len(conditions) == 0 {
		return fmt.Errorf("impossible de construire une chaîne sans conditions")
	}

	if parentNode == nil {
		return fmt.Errorf("le nœud parent ne peut pas être nil")
	}

	if network.AlphaSharingManager == nil {
		return fmt.Errorf("AlphaSharingManager non initialisé dans le réseau")
	}

	if network.LifecycleManager == nil {
		return fmt.Errorf("LifecycleManager non initialisé dans le réseau")
	}

	return nil
}

// initializeDecomposedChainBuild initialise le contexte de construction d'une chaîne décomposée
func initializeDecomposedChainBuild(
	conditions []DecomposedCondition,
	parentNode Node,
	ruleID string,
) *DecomposedChainBuildContext {
	return &DecomposedChainBuildContext{
		StartTime:       time.Now(),
		NodesCreated:    0,
		NodesReused:     0,
		HashesGenerated: make([]string, 0, len(conditions)),
		Chain: &AlphaChain{
			Nodes:  make([]*AlphaNode, 0, len(conditions)),
			Hashes: make([]string, 0, len(conditions)),
			RuleID: ruleID,
		},
		CurrentParent: parentNode,
	}
}

// convertDecomposedConditionToMap convertit une DecomposedCondition en map
// pour la condition du nœud alpha
func convertDecomposedConditionToMap(decomposedCond DecomposedCondition) map[string]interface{} {
	return map[string]interface{}{
		"type":     decomposedCond.Type,
		"left":     decomposedCond.Left,
		"operator": decomposedCond.Operator,
		"right":    decomposedCond.Right,
	}
}

// configureNodeDecompositionMetadata configure les métadonnées de décomposition
// sur un nœud alpha
func configureNodeDecompositionMetadata(
	alphaNode *AlphaNode,
	decomposedCond DecomposedCondition,
) {
	alphaNode.ResultName = decomposedCond.ResultName
	alphaNode.IsAtomic = decomposedCond.IsAtomic
	alphaNode.Dependencies = decomposedCond.Dependencies
}

// addNodeToChain ajoute un nœud alpha et son hash à la chaîne en construction
func addNodeToChain(
	ctx *DecomposedChainBuildContext,
	alphaNode *AlphaNode,
	hash string,
) {
	ctx.Chain.Nodes = append(ctx.Chain.Nodes, alphaNode)
	ctx.Chain.Hashes = append(ctx.Chain.Hashes, hash)
	ctx.HashesGenerated = append(ctx.HashesGenerated, hash)
}

// handleReusedDecomposedNode gère le cas où un nœud alpha décomposé est réutilisé
func handleReusedDecomposedNode(
	acb *AlphaChainBuilder,
	alphaNode *AlphaNode,
	currentParent Node,
	ruleID string,
	conditionIndex int,
	totalConditions int,
) {
	fmt.Printf("♻️  [AlphaChainBuilder] Réutilisation du nœud alpha %s (decomposed: %s) pour la règle %s (condition %d/%d)",
		alphaNode.ID, alphaNode.ResultName, ruleID, conditionIndex+1, totalConditions)

	if !acb.isAlreadyConnectedCached(currentParent, alphaNode) {
		// Connecter au parent si pas déjà connecté
		currentParent.AddChild(alphaNode)
		fmt.Printf("🔗 [AlphaChainBuilder] Connexion du nœud réutilisé %s au parent %s",
			alphaNode.ID, currentParent.GetID())
	} else {
		fmt.Printf("✓  [AlphaChainBuilder] Nœud %s déjà connecté au parent %s",
			alphaNode.ID, currentParent.GetID())
	}
}

// handleNewDecomposedNode gère le cas où un nouveau nœud alpha décomposé est créé
func handleNewDecomposedNode(
	acb *AlphaChainBuilder,
	alphaNode *AlphaNode,
	currentParent Node,
	ruleID string,
	conditionIndex int,
	totalConditions int,
) {
	// Connecter le nouveau nœud au parent
	currentParent.AddChild(alphaNode)

	// Ajouter le nœud au réseau
	acb.network.AlphaNodes[alphaNode.ID] = alphaNode

	// Mettre à jour le cache de connexion
	acb.updateConnectionCache(currentParent.GetID(), alphaNode.ID, true)

	fmt.Printf("🆕 [AlphaChainBuilder] Nouveau nœud alpha %s créé (decomposed: %s, deps: %v) pour la règle %s (condition %d/%d)",
		alphaNode.ID, alphaNode.ResultName, alphaNode.Dependencies, ruleID, conditionIndex+1, totalConditions)
	fmt.Printf("🔗 [AlphaChainBuilder] Connexion du nœud %s au parent %s",
		alphaNode.ID, currentParent.GetID())
}

// registerDecomposedNodeInLifecycle enregistre un nœud alpha décomposé
// dans le LifecycleManager
func registerDecomposedNodeInLifecycle(
	network *ReteNetwork,
	alphaNode *AlphaNode,
	ruleID string,
	reused bool,
) {
	lifecycle := network.LifecycleManager.RegisterNode(alphaNode.ID, "alpha")
	lifecycle.AddRuleReference(ruleID, "") // RuleName peut être ajouté plus tard si nécessaire

	if reused {
		fmt.Printf("📊 [AlphaChainBuilder] Nœud %s maintenant utilisé par %d règle(s)",
			alphaNode.ID, lifecycle.GetRefCount())
	}
}

// finalizeDecomposedChain finalise la construction d'une chaîne décomposée
// et enregistre les métriques
func finalizeDecomposedChain(
	ctx *DecomposedChainBuildContext,
	metrics *ChainBuildMetrics,
	ruleID string,
) {
	// Le dernier nœud de la chaîne est le nœud final
	ctx.Chain.FinalNode = ctx.Chain.Nodes[len(ctx.Chain.Nodes)-1]

	fmt.Printf("✅ [AlphaChainBuilder] Chaîne alpha décomposée complète construite pour la règle %s: %d nœud(s) atomiques",
		ruleID, len(ctx.Chain.Nodes))

	// Enregistrer les métriques si disponibles
	if metrics != nil {
		buildTime := time.Since(ctx.StartTime)
		detail := ChainMetricDetail{
			RuleID:          ruleID,
			ChainLength:     len(ctx.Chain.Nodes),
			NodesCreated:    ctx.NodesCreated,
			NodesReused:     ctx.NodesReused,
			BuildTime:       buildTime,
			Timestamp:       time.Now(),
			HashesGenerated: ctx.HashesGenerated,
		}
		metrics.RecordChainBuild(detail)
	}
}

// processDecomposedCondition traite une condition décomposée et l'ajoute à la chaîne
func processDecomposedCondition(
	acb *AlphaChainBuilder,
	ctx *DecomposedChainBuildContext,
	decomposedCond DecomposedCondition,
	variableName string,
	conditionIndex int,
	totalConditions int,
	ruleID string,
) error {
	// 1. Convertir DecomposedCondition en map
	conditionMap := convertDecomposedConditionToMap(decomposedCond)

	// 2. Obtenir ou créer l'AlphaNode via le gestionnaire de partage
	alphaNode, hash, reused, err := acb.network.AlphaSharingManager.GetOrCreateAlphaNode(
		conditionMap,
		variableName,
		acb.storage,
	)
	if err != nil {
		return fmt.Errorf("erreur lors de la création/récupération du nœud alpha %d: %w", conditionIndex, err)
	}

	// 3. Configurer les métadonnées de décomposition
	configureNodeDecompositionMetadata(alphaNode, decomposedCond)

	// 4. Ajouter le nœud et son hash à la chaîne
	addNodeToChain(ctx, alphaNode, hash)

	// 5. Traiter selon que le nœud est réutilisé ou nouveau
	if reused {
		ctx.NodesReused++
		handleReusedDecomposedNode(acb, alphaNode, ctx.CurrentParent, ruleID, conditionIndex, totalConditions)
	} else {
		ctx.NodesCreated++
		handleNewDecomposedNode(acb, alphaNode, ctx.CurrentParent, ruleID, conditionIndex, totalConditions)
	}

	// 6. Enregistrer le nœud dans le LifecycleManager
	registerDecomposedNodeInLifecycle(acb.network, alphaNode, ruleID, reused)

	// 7. Le nœud actuel devient le parent pour le prochain nœud
	ctx.CurrentParent = alphaNode

	return nil
}
