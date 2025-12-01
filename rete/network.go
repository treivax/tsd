// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"log"
)

// ReteNetwork représente le réseau RETE complet
type ReteNetwork struct {
	RootNode            *RootNode                `json:"root_node"`
	TypeNodes           map[string]*TypeNode     `json:"type_nodes"`
	AlphaNodes          map[string]*AlphaNode    `json:"alpha_nodes"`
	BetaNodes           map[string]interface{}   `json:"beta_nodes"` // Nœuds Beta pour les jointures multi-faits
	TerminalNodes       map[string]*TerminalNode `json:"terminal_nodes"`
	Storage             Storage                  `json:"-"`
	Types               []TypeDefinition         `json:"types"`
	BetaBuilder         interface{}              `json:"-"` // Constructeur de réseau Beta (deprecated, use BetaChainBuilder)
	LifecycleManager    *LifecycleManager        `json:"-"` // Gestionnaire du cycle de vie des nœuds
	AlphaSharingManager *AlphaSharingRegistry    `json:"-"` // Gestionnaire du partage des AlphaNodes
	BetaSharingRegistry BetaSharingRegistry      `json:"-"` // Gestionnaire du partage des JoinNodes
	BetaChainBuilder    *BetaChainBuilder        `json:"-"` // Constructeur de chaînes beta avec partage
	ChainMetrics        *ChainBuildMetrics       `json:"-"` // Métriques de performance pour la construction des chaînes
	Config              *ChainPerformanceConfig  `json:"-"` // Configuration de performance
	ActionExecutor      *ActionExecutor          `json:"-"` // Exécuteur d'actions
}

// NewReteNetwork crée un nouveau réseau RETE avec la configuration par défaut
func NewReteNetwork(storage Storage) *ReteNetwork {
	return NewReteNetworkWithConfig(storage, DefaultChainPerformanceConfig())
}

// NewReteNetworkWithConfig crée un nouveau réseau RETE avec une configuration personnalisée
func NewReteNetworkWithConfig(storage Storage, config *ChainPerformanceConfig) *ReteNetwork {
	if config == nil {
		config = DefaultChainPerformanceConfig()
	}

	rootNode := NewRootNode(storage)
	metrics := NewChainBuildMetrics()
	lifecycleManager := NewLifecycleManager()

	// Initialize Beta sharing if enabled
	var betaSharingRegistry BetaSharingRegistry
	var betaChainBuilder *BetaChainBuilder

	if config.BetaSharingEnabled {
		betaSharingConfig := BetaSharingConfig{
			Enabled:                     true,
			HashCacheSize:               config.BetaHashCacheMaxSize,
			MaxSharedNodes:              10000, // Default limit
			EnableMetrics:               true,
			NormalizeOrder:              true,
			EnableAdvancedNormalization: false,
		}
		betaSharingRegistry = NewBetaSharingRegistry(betaSharingConfig, lifecycleManager)
	}

	network := &ReteNetwork{
		RootNode:            rootNode,
		TypeNodes:           make(map[string]*TypeNode),
		AlphaNodes:          make(map[string]*AlphaNode),
		BetaNodes:           make(map[string]interface{}),
		TerminalNodes:       make(map[string]*TerminalNode),
		Storage:             storage,
		Types:               make([]TypeDefinition, 0),
		BetaBuilder:         nil, // Deprecated field, kept for backward compatibility
		LifecycleManager:    lifecycleManager,
		AlphaSharingManager: NewAlphaSharingRegistryWithConfig(config, metrics),
		BetaSharingRegistry: betaSharingRegistry,
		BetaChainBuilder:    betaChainBuilder, // Will be initialized lazily if needed
		ChainMetrics:        metrics,
		Config:              config,
	}

	// Initialize action executor
	network.ActionExecutor = NewActionExecutor(network, log.Default())

	// Initialize BetaChainBuilder if Beta sharing is enabled
	if betaSharingRegistry != nil {
		betaChainBuilder = NewBetaChainBuilderWithComponents(
			network,
			storage,
			betaSharingRegistry,
			lifecycleManager,
		)
		betaChainBuilder.SetOptimizationEnabled(true)
		betaChainBuilder.SetPrefixSharingEnabled(true)
		network.BetaChainBuilder = betaChainBuilder
	}

	return network
}

// GetChainMetrics retourne les métriques de performance pour la construction des chaînes alpha
func (rn *ReteNetwork) GetChainMetrics() *ChainBuildMetrics {
	if rn.ChainMetrics == nil {
		rn.ChainMetrics = NewChainBuildMetrics()
	}
	return rn.ChainMetrics
}

// GetBetaSharingStats retourne les statistiques de partage des JoinNodes
func (rn *ReteNetwork) GetBetaSharingStats() *BetaSharingStats {
	if rn.BetaSharingRegistry == nil {
		return nil
	}
	return rn.BetaSharingRegistry.GetSharingStats()
}

// GetBetaChainMetrics retourne les métriques de construction des chaînes beta
func (rn *ReteNetwork) GetBetaChainMetrics() *BetaChainMetrics {
	if rn.BetaChainBuilder == nil {
		return nil
	}
	return rn.BetaChainBuilder.GetMetrics()
}

// GetConfig retourne la configuration de performance
func (rn *ReteNetwork) GetConfig() *ChainPerformanceConfig {
	if rn.Config == nil {
		rn.Config = DefaultChainPerformanceConfig()
	}
	return rn.Config
}

// ResetChainMetrics réinitialise toutes les métriques de performance
func (rn *ReteNetwork) ResetChainMetrics() {
	if rn.ChainMetrics != nil {
		rn.ChainMetrics.Reset()
	}
	if rn.BetaChainBuilder != nil {
		rn.BetaChainBuilder.ResetMetrics()
	}
}

// SubmitFact soumet un nouveau fait au réseau
func (rn *ReteNetwork) SubmitFact(fact *Fact) error {
	fmt.Printf("🔥 Soumission d'un nouveau fait au réseau RETE: %s\n", fact.String())

	// Propager le fait depuis le nœud racine
	return rn.RootNode.ActivateRight(fact)
}

// SubmitFactsFromGrammar soumet plusieurs faits depuis la grammaire au réseau
func (rn *ReteNetwork) SubmitFactsFromGrammar(facts []map[string]interface{}) error {
	for i, factMap := range facts {
		// Convertir le map en Fact
		factID := fmt.Sprintf("fact_%d", i)
		if id, ok := factMap["id"].(string); ok {
			factID = id
		}

		factType := "unknown"
		if typ, ok := factMap["type"].(string); ok {
			factType = typ
		}

		fact := &Fact{
			ID:     factID,
			Type:   factType,
			Fields: make(map[string]interface{}),
		}

		// Copier tous les champs
		for key, value := range factMap {
			if key != "id" && key != "type" {
				fact.Fields[key] = value
			}
		}

		if err := rn.SubmitFact(fact); err != nil {
			return fmt.Errorf("erreur soumission fait %s: %w", fact.ID, err)
		}
	}
	return nil
}

// RetractFact retire un fait du réseau et propage la rétractation
// factID doit être l'identifiant interne (Type_ID)
func (rn *ReteNetwork) RetractFact(factID string) error {
	fmt.Printf("🗑️  Rétractation du fait: %s\n", factID)

	// Vérifier que le fait existe dans le réseau
	memory := rn.RootNode.GetMemory()
	if _, exists := memory.GetFact(factID); !exists {
		return fmt.Errorf("fait %s introuvable dans le réseau", factID)
	}

	// Propager la rétractation depuis le nœud racine
	return rn.RootNode.ActivateRetract(factID)
}

// Reset clears the entire RETE network and resets it to an empty state.
// This removes all facts, rules, types, and network nodes.
// After calling Reset, the network is ready to accept new definitions from scratch.
func (rn *ReteNetwork) Reset() {
	fmt.Println("🧹 Réinitialisation complète du réseau RETE")

	// Clear all node collections
	rn.TypeNodes = make(map[string]*TypeNode)
	rn.AlphaNodes = make(map[string]*AlphaNode)
	rn.BetaNodes = make(map[string]interface{})
	rn.TerminalNodes = make(map[string]*TerminalNode)
	rn.Types = make([]TypeDefinition, 0)
	rn.BetaBuilder = nil

	// Reset lifecycle manager
	if rn.LifecycleManager != nil {
		rn.LifecycleManager.Reset()
	} else {
		rn.LifecycleManager = NewLifecycleManager()
	}

	// Reset alpha sharing manager
	if rn.AlphaSharingManager != nil {
		rn.AlphaSharingManager.Reset()
	} else {
		rn.AlphaSharingManager = NewAlphaSharingRegistry()
	}

	// Recreate a fresh root node with the existing storage
	rn.RootNode = NewRootNode(rn.Storage)

	fmt.Println("✅ Réseau RETE réinitialisé avec succès")
}

// ClearMemory efface uniquement les mémoires (faits et tokens) de tous les nœuds
// sans détruire la structure du réseau
func (rn *ReteNetwork) ClearMemory() {
	fmt.Println("🧹 Nettoyage de la mémoire du réseau RETE")

	// Clear TypeNode memories
	for _, typeNode := range rn.TypeNodes {
		typeNode.mutex.Lock()
		typeNode.Memory.Facts = make(map[string]*Fact)
		typeNode.Memory.Tokens = make(map[string]*Token)
		typeNode.mutex.Unlock()
	}

	// Clear AlphaNode memories
	for _, alphaNode := range rn.AlphaNodes {
		alphaNode.mutex.Lock()
		alphaNode.Memory.Facts = make(map[string]*Fact)
		alphaNode.Memory.Tokens = make(map[string]*Token)
		alphaNode.mutex.Unlock()
	}

	// Clear BetaNode memories (JoinNodes, etc.)
	for _, betaNode := range rn.BetaNodes {
		if node, ok := betaNode.(Node); ok {
			node.GetMemory().Facts = make(map[string]*Fact)
			node.GetMemory().Tokens = make(map[string]*Token)
		}
	}

	// Clear TerminalNode memories
	for _, terminalNode := range rn.TerminalNodes {
		terminalNode.mutex.Lock()
		terminalNode.Memory.Facts = make(map[string]*Fact)
		terminalNode.Memory.Tokens = make(map[string]*Token)
		terminalNode.mutex.Unlock()
	}

	fmt.Println("✅ Mémoire du réseau RETE nettoyée avec succès")
}

// RemoveRule supprime une règle et tous ses nœuds qui ne sont plus utilisés
func (rn *ReteNetwork) RemoveRule(ruleID string) error {
	fmt.Printf("🗑️  Suppression de la règle: %s\n", ruleID)

	if rn.LifecycleManager == nil {
		return fmt.Errorf("LifecycleManager non initialisé")
	}

	// Récupérer tous les nœuds utilisés par cette règle
	nodeIDs := rn.LifecycleManager.GetNodesForRule(ruleID)
	if len(nodeIDs) == 0 {
		return fmt.Errorf("règle %s non trouvée ou aucun nœud associé", ruleID)
	}

	fmt.Printf("   📊 Nœuds associés à la règle: %d\n", len(nodeIDs))

	// Detect rule type and use appropriate removal strategy
	hasChain := false
	hasJoinNodes := false

	for _, nodeID := range nodeIDs {
		if rn.isPartOfChain(nodeID) {
			hasChain = true
		}
		if rn.isJoinNode(nodeID) {
			hasJoinNodes = true
		}
	}

	// Utiliser la suppression optimisée pour les chaînes avec joins
	if hasJoinNodes {
		fmt.Printf("   🔗 JoinNodes détectés, utilisation de la suppression avec lifecycle\n")
		return rn.removeRuleWithJoins(ruleID, nodeIDs)
	}

	// Utiliser la suppression optimisée pour les chaînes alpha
	if hasChain {
		fmt.Printf("   🔗 Chaîne d'AlphaNodes détectée, utilisation de la suppression optimisée\n")
		return rn.removeAlphaChain(ruleID)
	}

	// Comportement classique pour les règles simples
	return rn.removeSimpleRule(ruleID, nodeIDs)
}

// removeSimpleRule supprime une règle simple (sans chaîne)
func (rn *ReteNetwork) removeSimpleRule(ruleID string, nodeIDs []string) error {
	// Parcourir chaque nœud et retirer la référence à la règle
	nodesToDelete := make([]string, 0)
	for _, nodeID := range nodeIDs {
		shouldDelete, err := rn.LifecycleManager.RemoveRuleFromNode(nodeID, ruleID)
		if err != nil {
			fmt.Printf("   ⚠️  Erreur lors de la suppression de la règle du nœud %s: %v\n", nodeID, err)
			continue
		}

		if shouldDelete {
			nodesToDelete = append(nodesToDelete, nodeID)
			fmt.Printf("   ✓ Nœud %s marqué pour suppression (plus de références)\n", nodeID)
		} else {
			lifecycle, _ := rn.LifecycleManager.GetNodeLifecycle(nodeID)
			fmt.Printf("   ✓ Nœud %s conservé (%d référence(s) restante(s))\n", nodeID, lifecycle.GetRefCount())
		}
	}

	// Supprimer les nœuds qui n'ont plus de références
	for _, nodeID := range nodesToDelete {
		if err := rn.removeNodeFromNetwork(nodeID); err != nil {
			fmt.Printf("   ⚠️  Erreur lors de la suppression du nœud %s: %v\n", nodeID, err)
		} else {
			fmt.Printf("   🗑️  Nœud %s supprimé du réseau\n", nodeID)
		}
	}

	fmt.Printf("✅ Règle %s supprimée avec succès (%d nœud(s) supprimé(s))\n", ruleID, len(nodesToDelete))
	return nil
}

// removeAlphaChain supprime une règle avec une chaîne d'AlphaNodes
// Remonte la chaîne en ordre inverse depuis le terminal pour supprimer les nœuds
func (rn *ReteNetwork) removeAlphaChain(ruleID string) error {
	// Récupérer tous les nœuds de la règle
	nodeIDs := rn.LifecycleManager.GetNodesForRule(ruleID)

	// Séparer les nœuds par type
	var terminalID string
	alphaNodes := make([]string, 0)
	otherNodes := make([]string, 0)

	for _, nodeID := range nodeIDs {
		lifecycle, exists := rn.LifecycleManager.GetNodeLifecycle(nodeID)
		if !exists {
			continue
		}

		switch lifecycle.NodeType {
		case "terminal":
			terminalID = nodeID
		case "alpha":
			alphaNodes = append(alphaNodes, nodeID)
		default:
			otherNodes = append(otherNodes, nodeID)
		}
	}

	// Supprimer le terminal en premier
	deletedCount := 0
	if terminalID != "" {
		if err := rn.removeNodeWithCheck(terminalID, ruleID); err == nil {
			deletedCount++
			fmt.Printf("   🗑️  TerminalNode %s supprimé\n", terminalID)
		}
	}

	// Ordonner les AlphaNodes dans l'ordre inverse de la chaîne (du terminal vers le TypeNode)
	orderedAlphaNodes := rn.orderAlphaNodesReverse(alphaNodes)

	// Parcourir les AlphaNodes en ordre inverse
	stopDeletion := false
	for i, nodeID := range orderedAlphaNodes {
		lifecycle, exists := rn.LifecycleManager.GetNodeLifecycle(nodeID)
		if !exists {
			continue
		}

		// Décrémenter RefCount pour tous les nœuds
		shouldDelete, err := rn.LifecycleManager.RemoveRuleFromNode(nodeID, ruleID)
		if err != nil {
			fmt.Printf("   ⚠️  Erreur lors de la suppression de la règle du nœud %s: %v\n", nodeID, err)
			continue
		}

		if !stopDeletion && shouldDelete {
			// RefCount == 0, on peut supprimer
			if err := rn.removeNodeFromNetwork(nodeID); err != nil {
				fmt.Printf("   ⚠️  Erreur suppression nœud %s: %v\n", nodeID, err)
			} else {
				deletedCount++
				fmt.Printf("   🗑️  AlphaNode %s supprimé (position %d dans la chaîne)\n", nodeID, len(orderedAlphaNodes)-i)
			}
		} else if !shouldDelete && !stopDeletion {
			// Premier nœud partagé rencontré - on arrête la suppression mais on continue à décrémenter
			refCount := lifecycle.GetRefCount()
			fmt.Printf("   ♻️  AlphaNode %s conservé (%d référence(s) restante(s)) - arrêt des suppressions\n", nodeID, refCount)
			fmt.Printf("   ℹ️  Décrémentation du RefCount des nœuds parents partagés\n")
			stopDeletion = true
		} else if stopDeletion {
			// Nœuds parents - juste décrémenter le RefCount
			refCount := lifecycle.GetRefCount()
			fmt.Printf("   ♻️  AlphaNode %s: RefCount décrémenté (%d référence(s) restante(s))\n", nodeID, refCount)
		}
	}

	// Supprimer les autres nœuds (TypeNodes, JoinNodes, etc.)
	for _, nodeID := range otherNodes {
		if err := rn.removeNodeWithCheck(nodeID, ruleID); err == nil {
			deletedCount++
			lifecycle, _ := rn.LifecycleManager.GetNodeLifecycle(nodeID)
			fmt.Printf("   🗑️  %s %s supprimé\n", lifecycle.NodeType, nodeID)
		}
	}

	fmt.Printf("✅ Règle %s avec chaîne supprimée avec succès (%d nœud(s) supprimé(s))\n", ruleID, deletedCount)
	return nil
}

// removeNodeWithCheck supprime un nœud seulement si RefCount == 0
func (rn *ReteNetwork) removeNodeWithCheck(nodeID, ruleID string) error {
	shouldDelete, err := rn.LifecycleManager.RemoveRuleFromNode(nodeID, ruleID)
	if err != nil {
		return err
	}

	if shouldDelete {
		return rn.removeNodeFromNetwork(nodeID)
	}

	return fmt.Errorf("nœud %s encore référencé", nodeID)
}

// orderAlphaNodesReverse ordonne les AlphaNodes en ordre inverse de la chaîne
// (du nœud le plus éloigné du TypeNode vers le TypeNode)
func (rn *ReteNetwork) orderAlphaNodesReverse(alphaNodeIDs []string) []string {
	if len(alphaNodeIDs) <= 1 {
		return alphaNodeIDs
	}

	// Construire un graphe parent->enfants pour trouver l'ordre
	childToParent := make(map[string]string)
	hasParent := make(map[string]bool)

	for _, nodeID := range alphaNodeIDs {
		alphaNode, exists := rn.AlphaNodes[nodeID]
		if !exists {
			continue
		}

		parent := rn.getChainParent(alphaNode)
		if parent != nil {
			parentID := parent.GetID()
			// Vérifier si le parent est aussi un AlphaNode de notre liste
			for _, candidateID := range alphaNodeIDs {
				if candidateID == parentID {
					childToParent[nodeID] = parentID
					hasParent[nodeID] = true
					break
				}
			}
		}
	}

	// Trouver le nœud terminal de la chaîne (celui qui n'est parent de personne)
	var terminalNode string
	for _, nodeID := range alphaNodeIDs {
		isParent := false
		for _, parentID := range childToParent {
			if parentID == nodeID {
				isParent = true
				break
			}
		}
		if !isParent {
			terminalNode = nodeID
			break
		}
	}

	// Si pas de structure de chaîne détectée, retourner l'ordre original
	if terminalNode == "" {
		return alphaNodeIDs
	}

	// Remonter la chaîne depuis le terminal
	ordered := make([]string, 0, len(alphaNodeIDs))
	current := terminalNode
	visited := make(map[string]bool)

	for current != "" && !visited[current] {
		ordered = append(ordered, current)
		visited[current] = true
		current = childToParent[current]
	}

	// Ajouter les nœuds non visités (au cas où)
	for _, nodeID := range alphaNodeIDs {
		if !visited[nodeID] {
			ordered = append(ordered, nodeID)
		}
	}

	return ordered
}

// isPartOfChain détecte si un nœud fait partie d'une chaîne d'AlphaNodes
func (rn *ReteNetwork) isPartOfChain(nodeID string) bool {
	lifecycle, exists := rn.LifecycleManager.GetNodeLifecycle(nodeID)
	if !exists || lifecycle.NodeType != "alpha" {
		return false
	}

	alphaNode, exists := rn.AlphaNodes[nodeID]
	if !exists {
		return false
	}

	// Un AlphaNode fait partie d'une chaîne si:
	// 1. Son parent est un autre AlphaNode, OU
	// 2. Un de ses enfants est un autre AlphaNode

	parent := rn.getChainParent(alphaNode)
	if parent != nil && parent.GetType() == "alpha" {
		return true
	}

	children := alphaNode.GetChildren()
	for _, child := range children {
		if child.GetType() == "alpha" {
			return true
		}
	}

	return false
}

// getChainParent récupère le nœud parent d'un AlphaNode dans une chaîne
func (rn *ReteNetwork) getChainParent(alphaNode *AlphaNode) Node {
	if alphaNode == nil {
		return nil
	}

	alphaID := alphaNode.GetID()

	// Chercher dans les TypeNodes
	for _, typeNode := range rn.TypeNodes {
		for _, child := range typeNode.GetChildren() {
			if child.GetID() == alphaID {
				return typeNode
			}
		}
	}

	// Chercher dans les autres AlphaNodes
	for _, parentAlpha := range rn.AlphaNodes {
		if parentAlpha.GetID() == alphaID {
			continue
		}
		for _, child := range parentAlpha.GetChildren() {
			if child.GetID() == alphaID {
				return parentAlpha
			}
		}
	}

	return nil
}

// removeNodeFromNetwork supprime un nœud du réseau RETE
// Ne supprime que si RefCount == 0
func (rn *ReteNetwork) removeNodeFromNetwork(nodeID string) error {
	// Vérifier que le nœud existe et peut être supprimé
	lifecycle, exists := rn.LifecycleManager.GetNodeLifecycle(nodeID)
	if !exists {
		return fmt.Errorf("nœud %s non trouvé dans le LifecycleManager", nodeID)
	}

	// Ne pas supprimer si le nœud a encore des références
	if lifecycle.HasReferences() {
		return fmt.Errorf("impossible de supprimer le nœud %s: encore %d référence(s)",
			nodeID, lifecycle.GetRefCount())
	}

	// Déterminer le type de nœud et le supprimer de la map appropriée
	switch lifecycle.NodeType {
	case "type":
		// Trouver et supprimer le TypeNode
		for typeName, typeNode := range rn.TypeNodes {
			if typeNode.GetID() == nodeID {
				// Déconnecter du RootNode
				rn.removeChildFromNode(rn.RootNode, typeNode)
				delete(rn.TypeNodes, typeName)
				return rn.LifecycleManager.RemoveNode(nodeID)
			}
		}

	case "alpha":
		// Supprimer l'AlphaNode
		if alphaNode, exists := rn.AlphaNodes[nodeID]; exists {
			// Déconnecter des parents (TypeNodes ou autres AlphaNodes)
			parent := rn.getChainParent(alphaNode)
			if parent != nil {
				rn.removeChildFromNode(parent, alphaNode)
				fmt.Printf("   🔗 AlphaNode %s déconnecté de son parent %s\n", nodeID, parent.GetID())
			}

			delete(rn.AlphaNodes, nodeID)

			// Supprimer du registre de partage AlphaSharingManager
			if rn.AlphaSharingManager != nil {
				// Vérifier si c'est un nœud partagé (les nœuds partagés ont un ID qui commence par "alpha_")
				if len(nodeID) > 6 && nodeID[:6] == "alpha_" {
					if err := rn.AlphaSharingManager.RemoveAlphaNode(nodeID); err != nil {
						fmt.Printf("   ⚠️  Erreur suppression AlphaNode du registre de partage: %v\n", err)
					} else {
						fmt.Printf("   ✓ AlphaNode %s supprimé du AlphaSharingManager\n", nodeID)
					}
				}
			}

			return rn.LifecycleManager.RemoveNode(nodeID)
		}

	case "terminal":
		// Supprimer le TerminalNode
		if terminalNode, exists := rn.TerminalNodes[nodeID]; exists {
			// Déconnecter des parents (AlphaNodes ou JoinNodes)
			for _, alphaNode := range rn.AlphaNodes {
				rn.removeChildFromNode(alphaNode, terminalNode)
			}
			// Aussi déconnecter des BetaNodes si nécessaire
			for _, betaNode := range rn.BetaNodes {
				if node, ok := betaNode.(Node); ok {
					rn.removeChildFromNode(node, terminalNode)
				}
			}
			delete(rn.TerminalNodes, nodeID)
			return rn.LifecycleManager.RemoveNode(nodeID)
		}

	case "join", "exists", "accumulate":
		// Supprimer le BetaNode
		if betaNode, exists := rn.BetaNodes[nodeID]; exists {
			// Déconnecter des parents
			for _, typeNode := range rn.TypeNodes {
				if node, ok := betaNode.(Node); ok {
					rn.removeChildFromNode(typeNode, node)
				}
			}
			delete(rn.BetaNodes, nodeID)
			return rn.LifecycleManager.RemoveNode(nodeID)
		}
	}

	return fmt.Errorf("nœud %s non trouvé dans le réseau", nodeID)
}

// removeChildFromNode retire un nœud enfant d'un nœud parent
func (rn *ReteNetwork) removeChildFromNode(parent Node, child Node) {
	if parent == nil || child == nil {
		return
	}

	children := parent.GetChildren()
	newChildren := make([]Node, 0, len(children))
	for _, c := range children {
		if c.GetID() != child.GetID() {
			newChildren = append(newChildren, c)
		}
	}

	// Mettre à jour les enfants (nécessite un cast vers le type concret)
	switch p := parent.(type) {
	case *RootNode:
		p.Children = newChildren
	case *TypeNode:
		p.Children = newChildren
	case *AlphaNode:
		p.Children = newChildren
	case *JoinNode:
		p.Children = newChildren
	case *ExistsNode:
		p.Children = newChildren
	}
}

// GetRuleInfo retourne les informations d'une règle
func (rn *ReteNetwork) GetRuleInfo(ruleID string) *RuleInfo {
	if rn.LifecycleManager == nil {
		return &RuleInfo{
			RuleID:    ruleID,
			NodeIDs:   []string{},
			NodeCount: 0,
		}
	}
	return rn.LifecycleManager.GetRuleInfo(ruleID)
}

// GetNetworkStats retourne des statistiques sur le réseau
func (rn *ReteNetwork) GetNetworkStats() map[string]interface{} {
	stats := map[string]interface{}{
		"type_nodes":     len(rn.TypeNodes),
		"alpha_nodes":    len(rn.AlphaNodes),
		"beta_nodes":     len(rn.BetaNodes),
		"terminal_nodes": len(rn.TerminalNodes),
	}

	if rn.LifecycleManager != nil {
		lifecycleStats := rn.LifecycleManager.GetStats()
		for k, v := range lifecycleStats {
			stats["lifecycle_"+k] = v
		}
	}

	if rn.AlphaSharingManager != nil {
		alphaStats := rn.AlphaSharingManager.GetStats()
		for k, v := range alphaStats {
			stats["sharing_"+k] = v
		}
	}

	return stats
}

// isJoinNode checks if a node ID corresponds to a JoinNode
func (rn *ReteNetwork) isJoinNode(nodeID string) bool {
	_, exists := rn.BetaNodes[nodeID]
	return exists
}

// removeRuleWithJoins removes a rule that contains join nodes
func (rn *ReteNetwork) removeRuleWithJoins(ruleID string, nodeIDs []string) error {
	fmt.Printf("   🔗 Removing rule with join nodes: %s\n", ruleID)

	// Separate nodes by type
	var terminalNodes []string
	var joinNodes []string
	var alphaNodes []string
	var typeNodes []string

	for _, nodeID := range nodeIDs {
		lifecycle, exists := rn.LifecycleManager.GetNodeLifecycle(nodeID)
		if !exists {
			continue
		}

		switch lifecycle.NodeType {
		case "terminal":
			terminalNodes = append(terminalNodes, nodeID)
		case "join":
			joinNodes = append(joinNodes, nodeID)
		case "alpha":
			alphaNodes = append(alphaNodes, nodeID)
		case "type":
			typeNodes = append(typeNodes, nodeID)
		}
	}

	deletedCount := 0

	// Step 1: Remove terminal nodes first
	for _, nodeID := range terminalNodes {
		if err := rn.removeNodeWithCheck(nodeID, ruleID); err == nil {
			deletedCount++
			fmt.Printf("   🗑️  TerminalNode %s removed\n", nodeID)
		}
	}

	// Step 2: Remove join nodes with reference counting
	for _, nodeID := range joinNodes {
		// Remove rule reference from join node
		if rn.BetaSharingRegistry != nil {
			canDelete, err := rn.BetaSharingRegistry.RemoveRuleFromJoinNode(nodeID, ruleID)
			if err != nil {
				fmt.Printf("   ⚠️  Error removing rule from join node %s: %v\n", nodeID, err)
				continue
			}

			if canDelete {
				// No more rules reference this join node - safe to delete
				if err := rn.removeJoinNodeFromNetwork(nodeID); err == nil {
					deletedCount++
					fmt.Printf("   🗑️  JoinNode %s removed (no more references)\n", nodeID)
				}
			} else {
				// Join node is still shared by other rules
				refCount := rn.BetaSharingRegistry.GetJoinNodeRefCount(nodeID)
				fmt.Printf("   ✓ JoinNode %s preserved (%d rule(s) remaining)\n", nodeID, refCount)
			}
		} else {
			// No sharing registry - use lifecycle manager
			if err := rn.removeNodeWithCheck(nodeID, ruleID); err == nil {
				deletedCount++
				fmt.Printf("   🗑️  JoinNode %s removed\n", nodeID)
			}
		}
	}

	// Step 3: Remove alpha nodes
	for _, nodeID := range alphaNodes {
		if err := rn.removeNodeWithCheck(nodeID, ruleID); err == nil {
			deletedCount++
			fmt.Printf("   🗑️  AlphaNode %s removed\n", nodeID)
		} else {
			lifecycle, _ := rn.LifecycleManager.GetNodeLifecycle(nodeID)
			if lifecycle != nil && lifecycle.HasReferences() {
				fmt.Printf("   ✓ AlphaNode %s preserved (%d reference(s))\n", nodeID, lifecycle.GetRefCount())
			}
		}
	}

	// Step 4: Type nodes are usually shared - only remove if no references
	for _, nodeID := range typeNodes {
		lifecycle, exists := rn.LifecycleManager.GetNodeLifecycle(nodeID)
		if !exists {
			continue
		}

		shouldDelete, err := rn.LifecycleManager.RemoveRuleFromNode(nodeID, ruleID)
		if err != nil {
			fmt.Printf("   ⚠️  Error removing rule from type node %s: %v\n", nodeID, err)
			continue
		}

		if shouldDelete {
			if err := rn.removeNodeFromNetwork(nodeID); err == nil {
				deletedCount++
				fmt.Printf("   🗑️  TypeNode %s removed\n", nodeID)
			}
		} else {
			fmt.Printf("   ✓ TypeNode %s preserved (%d reference(s))\n", nodeID, lifecycle.GetRefCount())
		}
	}

	fmt.Printf("✅ Rule %s removed successfully (%d node(s) deleted)\n", ruleID, deletedCount)
	return nil
}

// removeJoinNodeFromNetwork removes a join node and all its dependent nodes from the network.
// This should only be called when the join node has no remaining rule references.
func (rn *ReteNetwork) removeJoinNodeFromNetwork(nodeID string) error {
	// Get the join node
	joinNode, exists := rn.BetaNodes[nodeID]
	if !exists {
		return fmt.Errorf("join node %s not found in network", nodeID)
	}

	fmt.Printf("   🗑️  Removing join node %s from network\n", nodeID)

	// Convert join node to proper type first
	var node Node
	var jn *JoinNode
	var ok bool
	if jn, ok = joinNode.(*JoinNode); !ok {
		return fmt.Errorf("beta node %s is not a JoinNode", nodeID)
	}
	node = jn

	// Step 1: Find and remove all terminal nodes that depend on this join node
	// Check if any terminal nodes are children of this join node
	for terminalID := range rn.TerminalNodes {
		// Check if this terminal is in the join node's children list
		isChild := false
		for _, child := range jn.GetChildren() {
			if child.GetID() == terminalID {
				isChild = true
				break
			}
		}

		if isChild {
			delete(rn.TerminalNodes, terminalID)
			fmt.Printf("   🗑️  Removed terminal node %s (child of join node)\n", terminalID)

			// Remove from lifecycle manager
			if rn.LifecycleManager != nil {
				rn.LifecycleManager.RemoveNode(terminalID)
			}
		}
	}

	// Step 2: Disconnect from parent nodes using the disconnectChild helper

	// Join nodes can have alpha nodes as parents
	for _, alphaNode := range rn.AlphaNodes {
		rn.disconnectChild(alphaNode, node)
	}

	// Check all other beta nodes (for cascading joins)
	for betaNodeID, betaNode := range rn.BetaNodes {
		if betaNodeID != nodeID {
			if bn, ok := betaNode.(*JoinNode); ok {
				rn.disconnectChild(bn, node)
			}
		}
	}

	// Also check type nodes (join nodes can connect directly to type nodes)
	for _, typeNode := range rn.TypeNodes {
		rn.disconnectChild(typeNode, node)
	}

	// Step 3: Remove from beta nodes map
	delete(rn.BetaNodes, nodeID)

	// Step 4: Remove from lifecycle manager
	if rn.LifecycleManager != nil {
		if err := rn.LifecycleManager.RemoveNode(nodeID); err != nil {
			fmt.Printf("   ⚠️  Warning: failed to remove join node %s from lifecycle manager: %v\n", nodeID, err)
		}
	}

	// Step 5: Remove from beta sharing registry
	if rn.BetaSharingRegistry != nil {
		if err := rn.BetaSharingRegistry.UnregisterJoinNode(nodeID); err != nil {
			fmt.Printf("   ⚠️  Warning: failed to unregister join node %s from beta sharing: %v\n", nodeID, err)
		}
	}

	fmt.Printf("   ✅ Join node %s successfully removed from network\n", nodeID)
	return nil
}

// disconnectChild removes a child from a node's children list
func (rn *ReteNetwork) disconnectChild(parent Node, child Node) {
	if parent == nil || child == nil {
		return
	}

	children := parent.GetChildren()
	newChildren := make([]Node, 0, len(children))
	for _, c := range children {
		if c.GetID() != child.GetID() {
			newChildren = append(newChildren, c)
		}
	}

	// Update parent's children list (this assumes BaseNode structure)
	if baseNode, ok := parent.(interface{ SetChildren([]Node) }); ok {
		baseNode.SetChildren(newChildren)
	}
}

// GetTypeDefinition retourne la définition d'un type par son nom
func (rn *ReteNetwork) GetTypeDefinition(typeName string) *TypeDefinition {
	for i := range rn.Types {
		if rn.Types[i].Name == typeName {
			return &rn.Types[i]
		}
	}
	return nil
}
