// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

// NewBetaChainBuilder crée un nouveau constructeur de chaînes beta avec des métriques neuves.
//
// Cette fonction initialise un builder avec un objet de métriques local. Pour partager
// les métriques entre plusieurs builders (recommandé), utilisez NewBetaChainBuilderWithMetrics.
//
// Le builder utilise le BetaSharingRegistry du réseau s'il existe, sinon désactive le partage.
//
// Paramètres:
//   - network: Réseau RETE auquel ajouter les nœuds
//   - storage: Backend de persistance pour les nœuds
//
// Retourne:
//   - Un nouveau builder prêt à l'emploi
//
// Exemple:
//
//	storage := NewMemoryStorage()
//	network := NewReteNetwork(storage)
//	builder := NewBetaChainBuilder(network, storage)
func NewBetaChainBuilder(network *ReteNetwork, storage Storage) *BetaChainBuilder {
	return NewBetaChainBuilderWithRegistry(network, storage, nil)
}

// NewBetaChainBuilderWithRegistry crée un builder avec un registry spécifique.
//
// Cette fonction permet de spécifier un BetaSharingRegistry personnalisé pour contrôler
// le partage de nœuds de jointure entre les règles.
//
// Paramètres:
//   - network: Réseau RETE auquel ajouter les nœuds
//   - storage: Backend de persistance pour les nœuds
//   - betaRegistry: Registry pour le partage de JoinNodes (peut être nil)
//
// Retourne:
//   - Un nouveau builder avec le registry spécifié
//
// Exemple:
//
//	registry := NewBetaSharingRegistry(config, lifecycle)
//	builder := NewBetaChainBuilderWithRegistry(network, storage, registry)
func NewBetaChainBuilderWithRegistry(network *ReteNetwork, storage Storage, betaRegistry BetaSharingRegistry) *BetaChainBuilder {
	return &BetaChainBuilder{
		network:             network,
		storage:             storage,
		betaSharingRegistry: betaRegistry,
		connectionCache:     make(map[string]bool),
		prefixCache:         make(map[string]*JoinNode),
		metrics:             NewBetaChainMetrics(),
		enableOptimization:  true, // Optimisation activée par défaut
		enablePrefixSharing: true, // Partage de préfixes activé par défaut
	}
}

// NewBetaChainBuilderWithMetrics crée un constructeur avec des métriques partagées.
//
// Recommandé quand le réseau RETE crée son propre builder, permettant de partager
// les métriques entre le builder et d'autres composants.
//
// Paramètres:
//   - network: Réseau RETE auquel ajouter les nœuds
//   - storage: Backend de persistance pour les nœuds
//   - metrics: Objet de métriques partagé (non nil)
//
// Retourne:
//   - Un nouveau builder utilisant les métriques fournies
//
// Exemple:
//
//	metrics := &BetaChainMetrics{}
//	builder := NewBetaChainBuilderWithMetrics(network, storage, metrics)
//	// Les métriques sont accessibles via builder.GetMetrics() et directement via 'metrics'
func NewBetaChainBuilderWithMetrics(network *ReteNetwork, storage Storage, metrics *BetaChainMetrics) *BetaChainBuilder {
	return NewBetaChainBuilderWithRegistryAndMetrics(network, storage, nil, metrics)
}

// NewBetaChainBuilderWithRegistryAndMetrics crée un builder avec registry et métriques.
//
// Cette fonction offre un contrôle total sur la configuration du builder en permettant
// de spécifier à la fois le registry de partage et l'objet de métriques.
//
// Paramètres:
//   - network: Réseau RETE auquel ajouter les nœuds
//   - storage: Backend de persistance pour les nœuds
//   - betaRegistry: Registry pour le partage de JoinNodes (peut être nil)
//   - metrics: Objet de métriques partagé (non nil)
//
// Retourne:
//   - Un nouveau builder avec registry et métriques spécifiés
//
// Exemple:
//
//	registry := NewBetaSharingRegistry(config, lifecycle)
//	metrics := &BetaChainMetrics{}
//	builder := NewBetaChainBuilderWithRegistryAndMetrics(network, storage, registry, metrics)
func NewBetaChainBuilderWithRegistryAndMetrics(network *ReteNetwork, storage Storage, betaRegistry BetaSharingRegistry, metrics *BetaChainMetrics) *BetaChainBuilder {
	return &BetaChainBuilder{
		network:             network,
		storage:             storage,
		betaSharingRegistry: betaRegistry,
		connectionCache:     make(map[string]bool),
		prefixCache:         make(map[string]*JoinNode),
		metrics:             metrics,
		enableOptimization:  true,
		enablePrefixSharing: true,
	}
}

// NewBetaChainBuilderWithComponents crée un builder avec tous les composants nécessaires.
//
// Cette fonction est utilisée lors de l'initialisation du ReteNetwork pour créer
// un builder complètement configuré avec registry de partage, lifecycle manager,
// et métriques partagées.
//
// Paramètres:
//   - network: Réseau RETE auquel ajouter les nœuds
//   - storage: Backend de persistance pour les nœuds
//   - betaRegistry: Registry pour le partage de JoinNodes (peut être nil)
//   - lifecycle: LifecycleManager pour la gestion du cycle de vie (peut être nil)
//
// Retourne:
//   - Un nouveau builder configuré avec tous les composants
//
// Exemple:
//
//	registry := NewBetaSharingRegistry(config, lifecycle)
//	builder := NewBetaChainBuilderWithComponents(network, storage, registry, lifecycle)
//
// Notes:
//   - Le lifecycle manager est géré par le registry s'il est fourni
//   - Si betaRegistry est nil, un builder basique est créé
//   - L'optimisation et le partage de préfixes sont activés par défaut
func NewBetaChainBuilderWithComponents(
	network *ReteNetwork,
	storage Storage,
	betaRegistry BetaSharingRegistry,
	lifecycle *LifecycleManager,
) *BetaChainBuilder {
	// Use the registry's lifecycle manager if available
	// Otherwise use the provided one
	if betaRegistry != nil {
		// Registry already has a lifecycle manager
		return &BetaChainBuilder{
			network:             network,
			storage:             storage,
			betaSharingRegistry: betaRegistry,
			connectionCache:     make(map[string]bool),
			prefixCache:         make(map[string]*JoinNode),
			metrics:             NewBetaChainMetrics(),
			enableOptimization:  true,
			enablePrefixSharing: true,
		}
	}

	// No registry, use basic builder
	return &BetaChainBuilder{
		network:             network,
		storage:             storage,
		betaSharingRegistry: nil,
		connectionCache:     make(map[string]bool),
		prefixCache:         make(map[string]*JoinNode),
		metrics:             NewBetaChainMetrics(),
		enableOptimization:  true,
		enablePrefixSharing: true,
	}
}
