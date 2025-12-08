// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"sync"
)

// AlphaChain représente une chaîne d'AlphaNodes construite pour un ensemble de conditions.
//
// Une chaîne alpha est une séquence ordonnée de nœuds alpha qui évaluent des conditions
// successives sur une même variable. Chaque nœud évalue une condition et propage les faits
// correspondants au nœud suivant dans la chaîne.
//
// Structure de chaîne typique:
//
//	TypeNode(Person)
//	  └── AlphaNode(p.age >= 18)
//	       └── AlphaNode(p.city == "Paris")
//	            └── TerminalNode(rule_terminal)
//
// Propriétés:
//   - len(Nodes) == len(Hashes) (toujours maintenu)
//   - FinalNode == Nodes[len(Nodes)-1] (si non vide)
//   - Ordre des nœuds correspond à l'ordre des conditions dans la règle
//
// Exemple d'utilisation:
//
//	chain, err := builder.BuildChain(conditions, "p", parentNode, "myRule")
//	if err != nil {
//	    return err
//	}
//	fmt.Printf("Chaîne construite: %d nœuds\n", len(chain.Nodes))
//	stats := builder.GetChainStats(chain)
//	fmt.Printf("Nœuds partagés: %d/%d\n", stats["shared_nodes"], len(chain.Nodes))
type AlphaChain struct {
	Nodes     []*AlphaNode `json:"nodes"`      // Liste ordonnée des nœuds alpha dans la chaîne
	Hashes    []string     `json:"hashes"`     // Hashes correspondants pour chaque nœud
	FinalNode *AlphaNode   `json:"final_node"` // Le dernier nœud de la chaîne
	RuleID    string       `json:"rule_id"`    // ID de la règle pour laquelle la chaîne a été construite
}

// AlphaChainBuilder construit des chaînes d'AlphaNodes avec partage automatique.
//
// Le builder coordonne la construction de chaînes alpha en réutilisant intelligemment
// les nœuds existants et en maintenant un cache des connexions parent→child pour
// éviter les duplications.
//
// Fonctionnalités principales:
//   - Construction séquentielle de chaînes condition par condition
//   - Partage automatique via AlphaSharingRegistry
//   - Cache de connexions pour éviter duplications
//   - Collection de métriques détaillées
//   - Thread-safe avec sync.RWMutex
//
// Flux de construction:
//
//	Pour chaque condition:
//	  1. Calculer hash (avec cache LRU)
//	  2. Chercher nœud existant via hash
//	  3. Si trouvé: réutiliser (refcount++)
//	     Si non: créer nouveau nœud
//	  4. Vérifier connexion parent→child (avec cache)
//	  5. Connecter si nécessaire
//	  6. Enregistrer dans LifecycleManager
//	  7. Nœud devient parent pour suivant
//
// Exemple d'utilisation:
//
//	builder := NewAlphaChainBuilder(network, storage)
//	conditions := []SimpleCondition{
//	    NewSimpleCondition("binaryOperation", "p.age", ">", 18),
//	    NewSimpleCondition("binaryOperation", "p.name", "==", "Alice"),
//	}
//	chain, err := builder.BuildChain(conditions, "p", typeNode, "rule1")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Accéder aux métriques
//	metrics := builder.GetMetrics()
//	fmt.Printf("Sharing ratio: %.1f%%\n", metrics.SharingRatio * 100)
//
// Thread-safety:
//   - Toutes les opérations publiques sont thread-safe
//   - Le cache de connexions est protégé par mutex
//   - Peut être utilisé concurremment par plusieurs goroutines
type AlphaChainBuilder struct {
	network         *ReteNetwork
	storage         Storage
	connectionCache map[string]bool // Cache pour les connexions existantes (parentID_childID -> bool)
	metrics         *ChainBuildMetrics
	mutex           sync.RWMutex
}

// NewAlphaChainBuilder crée un nouveau constructeur de chaînes alpha avec des métriques neuves.
//
// Cette fonction initialise un builder avec un objet de métriques local. Pour partager
// les métriques entre plusieurs builders (recommandé), utilisez NewAlphaChainBuilderWithMetrics.
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
//	builder := NewAlphaChainBuilder(network, storage)
func NewAlphaChainBuilder(network *ReteNetwork, storage Storage) *AlphaChainBuilder {
	return &AlphaChainBuilder{
		network:         network,
		storage:         storage,
		connectionCache: make(map[string]bool),
		metrics:         NewChainBuildMetrics(),
	}
}

// NewAlphaChainBuilderWithMetrics crée un constructeur avec des métriques partagées.
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
//	metrics := NewChainBuildMetrics()
//	builder := NewAlphaChainBuilderWithMetrics(network, storage, metrics)
//	// Les métriques sont accessibles via builder.GetMetrics() et directement via 'metrics'
func NewAlphaChainBuilderWithMetrics(network *ReteNetwork, storage Storage, metrics *ChainBuildMetrics) *AlphaChainBuilder {
	return &AlphaChainBuilder{
		network:         network,
		storage:         storage,
		connectionCache: make(map[string]bool),
		metrics:         metrics,
	}
}

// BuildChain construit une chaîne d'AlphaNodes pour un ensemble de conditions
// avec partage automatique des nœuds identiques entre règles.
//
// Cette méthode est le point d'entrée principal pour la construction de chaînes.
// Elle itère sur chaque condition, tente de réutiliser un nœud existant, sinon
// en crée un nouveau, et maintient les connexions parent→child appropriées.
//
// Algorithme:
//
//	Pour chaque condition dans la liste:
//	  1. Convertir SimpleCondition en map
//	  2. Appeler AlphaSharingRegistry.GetOrCreateAlphaNode()
//	     → Calcule hash (avec cache LRU)
//	     → Cherche nœud existant
//	     → Crée nouveau si inexistant
//	  3. Si nœud réutilisé:
//	     - Vérifier connexion avec parent (cache)
//	     - Connecter si nécessaire
//	  4. Si nœud créé:
//	     - Connecter au parent
//	     - Ajouter au réseau
//	     - Mettre en cache la connexion
//	  5. Enregistrer dans LifecycleManager
//	  6. Nœud devient parent pour itération suivante
//
// Paramètres:
//   - conditions: liste de conditions simples dans l'ordre d'évaluation
//   - variableName: nom de la variable (ex: "p", "u") - utilisé dans le hash
//   - parentNode: nœud parent (généralement TypeNode) auquel connecter le premier nœud
//   - ruleID: identifiant unique de la règle pour le lifecycle management
//
// Retourne:
//   - *AlphaChain: la chaîne construite avec tous les nœuds et leurs hashes
//   - error: erreur si conditions vides, parent nil, ou problème de création
//
// Exemple simple:
//
//	conditions := []SimpleCondition{
//	    NewSimpleCondition("binaryOperation", "p.age", ">", 18),
//	}
//	chain, err := builder.BuildChain(conditions, "p", typeNode, "rule_adult")
//	// → Crée: TypeNode → AlphaNode(p.age>18) → Terminal
//
// Exemple avec partage:
//
//	// Règle 1
//	chain1, _ := builder.BuildChain(
//	    []SimpleCondition{NewSimpleCondition("binaryOperation", "p.age", ">", 18)},
//	    "p", typeNode, "rule1")
//	// → Crée nouveau nœud alpha_abc123
//
//	// Règle 2 (même condition)
//	chain2, _ := builder.BuildChain(
//	    []SimpleCondition{NewSimpleCondition("binaryOperation", "p.age", ">", 18)},
//	    "p", typeNode, "rule2")
//	// → Réutilise alpha_abc123 (RefCount=2)
//
// Logs générés:
//
//	🆕 [AlphaChainBuilder] Nouveau nœud alpha alpha_abc123 créé pour la règle rule1 (condition 1/1)
//	🔗 [AlphaChainBuilder] Connexion du nœud alpha_abc123 au parent type_person
//	♻️  [AlphaChainBuilder] Réutilisation du nœud alpha alpha_abc123 pour la règle rule2 (condition 1/1)
//	✓  [AlphaChainBuilder] Nœud alpha_abc123 déjà connecté au parent type_person
func (acb *AlphaChainBuilder) BuildChain(
	conditions []SimpleCondition,
	variableName string,
	parentNode Node,
	ruleID string,
) (*AlphaChain, error) {
	// Validation des entrées
	if err := validateBuildChainInputs(conditions, parentNode, acb.network); err != nil {
		return nil, err
	}

	// Initialisation des métriques et de la chaîne
	metrics := initializeChainMetrics(len(conditions))
	chain := &AlphaChain{
		Nodes:  make([]*AlphaNode, 0, len(conditions)),
		Hashes: make([]string, 0, len(conditions)),
		RuleID: ruleID,
	}
	currentParent := parentNode

	// Construire la chaîne condition par condition
	for i, condition := range conditions {
		result, err := acb.buildAndConnectAlphaNode(
			condition, variableName, currentParent, ruleID,
			i, len(conditions), metrics,
		)
		if err != nil {
			return nil, err
		}

		chain.Nodes = append(chain.Nodes, result.node)
		chain.Hashes = append(chain.Hashes, result.hash)
		currentParent = result.node
	}

	// Finalisation de la chaîne
	chain.FinalNode = chain.Nodes[len(chain.Nodes)-1]
	acb.recordChainMetrics(ruleID, chain, metrics)

	return chain, nil
}

// BuildDecomposedChain constructs an alpha chain from decomposed conditions with full metadata.
// This method sets ResultName, IsAtomic, and Dependencies on each AlphaNode for
// context-aware evaluation with intermediate result propagation.
func (acb *AlphaChainBuilder) BuildDecomposedChain(
	conditions []DecomposedCondition,
	variableName string,
	parentNode Node,
	ruleID string,
) (*AlphaChain, error) {
	// Phase 1: Valider les entrées
	if err := validateBuildDecomposedInputs(conditions, parentNode, acb.network); err != nil {
		return nil, err
	}

	// Phase 2: Initialiser le contexte de construction
	ctx := initializeDecomposedChainBuild(conditions, parentNode, ruleID)

	// Phase 3: Construire la chaîne condition par condition
	for i, decomposedCond := range conditions {
		if err := processDecomposedCondition(
			acb,
			ctx,
			decomposedCond,
			variableName,
			i,
			len(conditions),
			ruleID,
		); err != nil {
			return nil, err
		}
	}

	// Phase 4: Finaliser la chaîne et enregistrer les métriques
	finalizeDecomposedChain(ctx, acb.metrics, ruleID)

	return ctx.Chain, nil
}

// GetMetrics retourne les métriques de performance du builder.
//
// Les métriques incluent:
//   - Nombre de chaînes construites
//   - Nœuds créés vs réutilisés
//   - Ratio de partage
//   - Hit rate du cache de hash
//   - Temps moyen de construction
//
// Retourne:
//   - Pointeur vers l'objet de métriques (non nil)
//
// Exemple:
//
//	metrics := builder.GetMetrics()
//	fmt.Printf("Chaînes construites: %d\n", metrics.TotalChainsBuilt)
//	fmt.Printf("Ratio de partage: %.1f%%\n", metrics.SharingRatio * 100)
//	fmt.Printf("Cache hit rate: %.1f%%\n",
//	    float64(metrics.HashCacheHits) /
//	    float64(metrics.HashCacheHits + metrics.HashCacheMisses) * 100)
func (acb *AlphaChainBuilder) GetMetrics() *ChainBuildMetrics {
	return acb.metrics
}
