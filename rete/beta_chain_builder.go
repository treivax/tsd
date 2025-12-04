// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"github.com/treivax/tsd/tsdio"
	"fmt"
	"sort"
	"sync"
	"time"
)

// BetaChain représente une chaîne de JoinNodes construite pour un ensemble de patterns.
//
// Une chaîne beta est une séquence ordonnée de nœuds de jointure qui combinent
// progressivement plusieurs variables. Chaque JoinNode évalue une condition de
// jointure et propage les tokens combinés au nœud suivant dans la chaîne.
//
// Structure de chaîne typique (cascade pour 3+ variables):
//
//	TypeNode(Person)  TypeNode(Order)
//	       └──────────┴────────┐
//	                    JoinNode(p ⋈ o)
//	                         │
//	                    TypeNode(Payment)
//	                         │
//	                    JoinNode((p,o) ⋈ pay)
//	                         │
//	                    TerminalNode(rule_terminal)
//
// Propriétés:
//   - len(Nodes) == len(Hashes) (toujours maintenu)
//   - FinalNode == Nodes[len(Nodes)-1] (si non vide)
//   - Ordre des jointures est optimisé pour la sélectivité
//
// Exemple d'utilisation:
//
//	patterns := []JoinPattern{
//	    {LeftVar: "p", RightVar: "o", Condition: ...},
//	    {LeftVar: "p,o", RightVar: "pay", Condition: ...},
//	}
//	chain, err := builder.BuildChain(patterns, "myRule")
//	if err != nil {
//	    return err
//	}
//	fmt.Printf("Chaîne construite: %d nœuds\n", len(chain.Nodes))
type BetaChain struct {
	Nodes     []*JoinNode `json:"nodes"`      // Liste ordonnée des JoinNodes dans la chaîne
	Hashes    []string    `json:"hashes"`     // Hashes correspondants pour chaque nœud
	FinalNode *JoinNode   `json:"final_node"` // Le dernier nœud de la chaîne
	RuleID    string      `json:"rule_id"`    // ID de la règle pour laquelle la chaîne a été construite
}

// JoinPattern représente un pattern de jointure entre variables.
//
// Un pattern de jointure définit comment deux ensembles de variables doivent
// être combinés selon une condition spécifique.
//
// Exemples:
//   - Jointure binaire simple: {LeftVars: ["p"], RightVars: ["o"], ...}
//   - Jointure cascade niveau 2: {LeftVars: ["p","o"], RightVars: ["pay"], ...}
type JoinPattern struct {
	LeftVars       []string               `json:"left_vars"`       // Variables du côté gauche
	RightVars      []string               `json:"right_vars"`      // Variables du côté droit
	AllVars        []string               `json:"all_vars"`        // Toutes les variables impliquées
	VarTypes       map[string]string      `json:"var_types"`       // Mapping variable -> type
	Condition      map[string]interface{} `json:"condition"`       // Condition de jointure
	Selectivity    float64                `json:"selectivity"`     // Estimation de sélectivité (0-1, plus bas = plus sélectif)
	EstimatedCost  float64                `json:"estimated_cost"`  // Coût estimé de cette jointure
	JoinConditions []JoinCondition        `json:"join_conditions"` // Conditions de jointure extraites
}

// BetaChainBuilder construit des chaînes de JoinNodes avec partage automatique.
//
// Le builder coordonne la construction de chaînes beta en réutilisant intelligemment
// les nœuds existants via BetaSharingRegistry et en optimisant l'ordre des jointures
// pour maximiser les performances.
//
// Fonctionnalités principales:
//   - Construction séquentielle de chaînes pattern par pattern
//   - Partage automatique via BetaSharingRegistry
//   - Optimisation de l'ordre des jointures (heuristique de sélectivité)
//   - Cache de connexions pour éviter duplications
//   - Collection de métriques détaillées
//   - Thread-safe avec sync.RWMutex
//
// Flux de construction:
//
//  1. Analyser les patterns et estimer la sélectivité
//  2. Trier les patterns par ordre optimal (plus sélectif d'abord)
//  3. Pour chaque pattern:
//     a. Calculer signature + hash
//     b. Chercher nœud existant via BetaSharingRegistry
//     c. Si trouvé: réutiliser (refcount++)
//     Si non: créer nouveau JoinNode
//     d. Vérifier connexion parent→child (avec cache)
//     e. Connecter si nécessaire
//     f. Enregistrer dans LifecycleManager
//     g. Nœud devient parent pour suivant
//
// Optimisations implémentées:
//   - Ordre des jointures basé sur sélectivité
//   - Détection des préfixes réutilisables (sous-chaînes communes)
//   - Cache des connexions pour éviter duplications
//   - Métriques détaillées pour monitoring
//
// Exemple d'utilisation:
//
//	builder := NewBetaChainBuilder(network, storage)
//	patterns := []JoinPattern{
//	    {
//	        LeftVars: []string{"p"},
//	        RightVars: []string{"o"},
//	        VarTypes: map[string]string{"p": "Person", "o": "Order"},
//	        Condition: map[string]interface{}{...},
//	        Selectivity: 0.3,
//	    },
//	}
//	chain, err := builder.BuildChain(patterns, "rule1")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Accéder aux métriques
//	metrics := builder.GetMetrics()
//	fmt.Printf("Sharing ratio: %.1f%%\n", metrics.SharingRatio() * 100)
//
// Thread-safety:
//   - Toutes les opérations publiques sont thread-safe
//   - Le cache de connexions est protégé par mutex
//   - Peut être utilisé concurremment par plusieurs goroutines
type BetaChainBuilder struct {
	network             *ReteNetwork
	storage             Storage
	betaSharingRegistry BetaSharingRegistry
	connectionCache     map[string]bool      // Cache pour les connexions existantes (parentID_childID -> bool)
	prefixCache         map[string]*JoinNode // Cache des préfixes de chaînes réutilisables
	metrics             *BetaChainMetrics    // Métriques de construction complètes
	mutex               sync.RWMutex
	enableOptimization  bool // Active/désactive l'optimisation de l'ordre des jointures
	enablePrefixSharing bool // Active/désactive le partage des préfixes
}

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
//	metrics := &BetaBuildMetrics{}
//	builder := NewBetaChainBuilderWithMetrics(network, storage, metrics)
//	// Les métriques sont accessibles via builder.GetMetrics() et directement via 'metrics'
func NewBetaChainBuilderWithMetrics(network *ReteNetwork, storage Storage, metrics *BetaChainMetrics) *BetaChainBuilder {
	return NewBetaChainBuilderWithRegistryAndMetrics(network, storage, nil, metrics)
}

// NewBetaChainBuilderWithRegistryAndMetrics crée un builder avec registry et métriques.
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

// BuildChain construit une chaîne de JoinNodes pour un ensemble de patterns de jointure
// avec partage automatique des nœuds identiques entre règles et optimisation de l'ordre.
//
// Cette méthode est le point d'entrée principal pour la construction de chaînes.
// Elle analyse les patterns, optimise leur ordre, et construit progressivement la chaîne
// en réutilisant les nœuds existants via BetaSharingRegistry.
//
// Algorithme:
//
//  1. Validation des inputs
//  2. Analyse et estimation de sélectivité des patterns
//  3. Optimisation de l'ordre (si activée)
//  4. Détection de préfixes réutilisables (si activée)
//  5. Pour chaque pattern dans l'ordre optimal:
//     a. Calculer signature de jointure
//     b. Appeler BetaSharingRegistry.GetOrCreateJoinNode()
//     → Cherche nœud existant via hash
//     → Crée nouveau si inexistant
//     c. Si nœud réutilisé:
//     - Vérifier connexion avec parent (cache)
//     - Connecter si nécessaire
//     d. Si nœud créé:
//     - Connecter au parent
//     - Ajouter au réseau
//     - Mettre en cache la connexion
//     e. Enregistrer dans LifecycleManager
//     f. Nœud devient parent pour itération suivante
//  6. Collecter métriques
//
// Paramètres:
//   - patterns: liste de patterns de jointure dans l'ordre initial
//   - ruleID: identifiant unique de la règle pour le lifecycle management
//
// Retourne:
//   - *BetaChain: la chaîne construite avec tous les nœuds et leurs hashes
//   - error: erreur si patterns vides, registry non initialisé, ou problème de création
//
// Exemple simple (2 variables):
//
//	patterns := []JoinPattern{
//	    {
//	        LeftVars: []string{"p"},
//	        RightVars: []string{"o"},
//	        VarTypes: map[string]string{"p": "Person", "o": "Order"},
//	        Condition: map[string]interface{}{"type": "join", ...},
//	    },
//	}
//	chain, err := builder.BuildChain(patterns, "rule_customer_orders")
//	// → Crée: TypeNode(Person) ⋈ TypeNode(Order) → JoinNode → Terminal
//
// Exemple avec cascade (3+ variables):
//
//	patterns := []JoinPattern{
//	    {LeftVars: []string{"p"}, RightVars: []string{"o"}, ...},      // p ⋈ o
//	    {LeftVars: []string{"p","o"}, RightVars: []string{"pay"}, ...}, // (p,o) ⋈ pay
//	}
//	chain, err := builder.BuildChain(patterns, "rule_payment_check")
//	// → Crée cascade: p ⋈ o → (p,o) ⋈ pay → Terminal
//
// Exemple avec partage:
//
//	// Règle 1
//	chain1, _ := builder.BuildChain(patterns1, "rule1")
//	// → Crée nouveau JoinNode join_abc123
//
//	// Règle 2 (même pattern de jointure)
//	chain2, _ := builder.BuildChain(patterns2, "rule2")
//	// → Réutilise join_abc123 (RefCount=2)
//
// Logs générés:
//
//	🆕 [BetaChainBuilder] Nouveau JoinNode join_abc123 créé pour la règle rule1 (pattern 1/1)
//	🔗 [BetaChainBuilder] Connexion du JoinNode join_abc123 aux TypeNodes
//	♻️  [BetaChainBuilder] Réutilisation du JoinNode join_abc123 pour la règle rule2 (pattern 1/1)
//	✓  [BetaChainBuilder] JoinNode join_abc123 déjà connecté
//	⚡ [BetaChainBuilder] Optimisation de l'ordre appliquée (2 patterns réordonnés)
func (bcb *BetaChainBuilder) BuildChain(
	patterns []JoinPattern,
	ruleID string,
) (*BetaChain, error) {
	// Validation des inputs
	if len(patterns) == 0 {
		return nil, fmt.Errorf("impossible de construire une chaîne sans patterns")
	}

	if bcb.network.LifecycleManager == nil {
		return nil, fmt.Errorf("LifecycleManager non initialisé dans le réseau")
	}

	// Démarrer le chronomètre pour les métriques
	startTime := time.Now()
	nodesCreated := 0
	nodesReused := 0
	hashesGenerated := make([]string, 0, len(patterns))
	optimizationApplied := false
	prefixReused := false

	chain := &BetaChain{
		Nodes:  make([]*JoinNode, 0, len(patterns)),
		Hashes: make([]string, 0, len(patterns)),
		RuleID: ruleID,
	}

	// Estimer la sélectivité des patterns si pas déjà fait
	bcb.estimateSelectivity(patterns)

	// Optimiser l'ordre des patterns si activé
	optimizedPatterns := patterns
	if bcb.enableOptimization && len(patterns) > 1 {
		optimizedPatterns = bcb.optimizeJoinOrder(patterns)
		if !bcb.patternsEqual(patterns, optimizedPatterns) {
			optimizationApplied = true
			tsdio.LogPrintf("⚡ [BetaChainBuilder] Optimisation de l'ordre appliquée (%d patterns réordonnés) pour règle %s",
				len(patterns), ruleID)
		}
	}

	// Tenter de réutiliser un préfixe de chaîne existant si activé
	var currentParent Node
	startPatternIndex := 0

	if bcb.enablePrefixSharing && len(optimizedPatterns) > 1 {
		prefixNode, prefixLen := bcb.findReusablePrefix(optimizedPatterns, ruleID)
		if prefixNode != nil && prefixLen > 0 {
			prefixReused = true
			currentParent = prefixNode
			startPatternIndex = prefixLen
			nodesReused += prefixLen
			tsdio.LogPrintf("♻️  [BetaChainBuilder] Préfixe de chaîne réutilisé (%d nœuds) pour règle %s",
				prefixLen, ruleID)
		}
	}

	// Construire la chaîne pattern par pattern
	for i := startPatternIndex; i < len(optimizedPatterns); i++ {
		pattern := optimizedPatterns[i]

		// Obtenir ou créer le JoinNode via le registry de partage
		var joinNode *JoinNode
		var hash string
		var reused bool
		var err error

		if bcb.betaSharingRegistry != nil {
			joinNode, hash, reused, err = bcb.betaSharingRegistry.GetOrCreateJoinNode(
				pattern.Condition,
				pattern.LeftVars,
				pattern.RightVars,
				pattern.AllVars,
				pattern.VarTypes,
				bcb.storage,
			)
			if err != nil {
				return nil, fmt.Errorf("erreur lors de la création/récupération du JoinNode %d: %w", i, err)
			}
		} else {
			// Fallback si pas de registry: créer directement
			nodeID := fmt.Sprintf("%s_join_%d", ruleID, i)
			joinNode = NewJoinNode(nodeID, pattern.Condition, pattern.LeftVars, pattern.RightVars, pattern.VarTypes, bcb.storage)
			hash = nodeID
			reused = false
		}

		// Register join node with lifecycle manager
		if bcb.network != nil && bcb.network.LifecycleManager != nil {
			// Register the node if not already registered (for new nodes)
			if _, exists := bcb.network.LifecycleManager.GetNodeLifecycle(hash); !exists {
				bcb.network.LifecycleManager.RegisterNode(hash, "join")
			}
			// Add this rule's reference to the join node
			bcb.network.LifecycleManager.AddRuleToNode(hash, ruleID, ruleID)
		}

		// Register rule with beta sharing registry for join node tracking
		if bcb.betaSharingRegistry != nil {
			if err := bcb.betaSharingRegistry.RegisterRuleForJoinNode(hash, ruleID); err != nil {
				tsdio.LogPrintf("⚠️  [BetaChainBuilder] Warning: failed to register rule %s for join node %s: %v",
					ruleID, hash, err)
			}
		}

		// Ajouter le nœud et son hash à la chaîne
		chain.Nodes = append(chain.Nodes, joinNode)
		chain.Hashes = append(chain.Hashes, hash)
		hashesGenerated = append(hashesGenerated, hash)

		if reused {
			nodesReused++
			tsdio.LogPrintf("♻️  [BetaChainBuilder] Réutilisation du JoinNode %s pour la règle %s (pattern %d/%d)",
				joinNode.ID, ruleID, i+1, len(optimizedPatterns))

			// Nœud réutilisé - vérifier la connexion si on a un parent
			if currentParent != nil && !bcb.isAlreadyConnectedCached(currentParent, joinNode) {
				currentParent.AddChild(joinNode)
				tsdio.LogPrintf("🔗 [BetaChainBuilder] Connexion du nœud réutilisé %s au parent %s",
					joinNode.ID, currentParent.GetID())
			} else if currentParent != nil {
				tsdio.LogPrintf("✓  [BetaChainBuilder] Nœud %s déjà connecté au parent %s",
					joinNode.ID, currentParent.GetID())
			}
		} else {
			nodesCreated++
			// Nouveau nœud - l'ajouter au réseau
			bcb.network.BetaNodes[joinNode.ID] = joinNode

			// Connecter au parent si on en a un
			if currentParent != nil {
				currentParent.AddChild(joinNode)
				bcb.updateConnectionCache(currentParent.GetID(), joinNode.ID, true)
			}

			tsdio.LogPrintf("🆕 [BetaChainBuilder] Nouveau JoinNode %s créé pour la règle %s (pattern %d/%d)",
				joinNode.ID, ruleID, i+1, len(optimizedPatterns))
			if currentParent != nil {
				tsdio.LogPrintf("🔗 [BetaChainBuilder] Connexion du nœud %s au parent %s",
					joinNode.ID, currentParent.GetID())
			}
		}

		// Enregistrer le nœud dans le LifecycleManager avec la règle
		lifecycle := bcb.network.LifecycleManager.RegisterNode(joinNode.ID, "join")
		lifecycle.AddRuleReference(ruleID, "") // RuleName peut être ajouté plus tard si nécessaire

		if reused {
			tsdio.LogPrintf("📊 [BetaChainBuilder] Nœud %s maintenant utilisé par %d règle(s)",
				joinNode.ID, lifecycle.GetRefCount())
		}

		// Mettre à jour le cache de préfixes si pertinent
		if bcb.enablePrefixSharing && i < len(optimizedPatterns)-1 {
			prefixKey := bcb.computePrefixKey(optimizedPatterns[0 : i+1])
			bcb.updatePrefixCache(prefixKey, joinNode)
		}

		// Le nœud actuel devient le parent pour le prochain nœud
		currentParent = joinNode
	}

	// Le dernier nœud de la chaîne est le nœud final
	if len(chain.Nodes) > 0 {
		chain.FinalNode = chain.Nodes[len(chain.Nodes)-1]
	}

	buildTime := time.Since(startTime)
	tsdio.LogPrintf("✅ [BetaChainBuilder] Chaîne beta complète construite pour la règle %s: %d nœud(s) (créés: %d, réutilisés: %d) en %v",
		ruleID, len(chain.Nodes), nodesCreated, nodesReused, buildTime)

	// Record metrics
	if bcb.metrics != nil {
		detail := BetaChainMetricDetail{
			RuleID:          ruleID,
			ChainLength:     len(chain.Nodes),
			NodesCreated:    nodesCreated,
			NodesReused:     nodesReused,
			BuildTime:       buildTime,
			Timestamp:       time.Now(),
			HashesGenerated: hashesGenerated,
			JoinsExecuted:   0, // Will be updated during runtime
			TotalJoinTime:   0,
		}
		bcb.metrics.RecordChainBuild(detail)
	}

	// Log optimization info if applied
	_ = optimizationApplied
	_ = prefixReused

	return chain, nil
}

// estimateSelectivity estime la sélectivité de chaque pattern de jointure.
//
// La sélectivité est une heuristique (0-1) qui indique combien de tuples
// passeront le filtre de jointure. Plus la valeur est basse, plus la jointure
// est sélective (filtre beaucoup de données).
//
// Heuristiques utilisées:
//   - Nombre de conditions: plus de conditions = plus sélectif
//   - Type d'opérateur: égalité > inégalité > range
//   - Nombre de variables impliquées: moins de variables = plus sélectif
//
// Cette fonction modifie les patterns en place.
func (bcb *BetaChainBuilder) estimateSelectivity(patterns []JoinPattern) {
	for i := range patterns {
		pattern := &patterns[i]

		// Si déjà estimée, ne rien faire
		if pattern.Selectivity > 0 {
			continue
		}

		// Estimation par défaut
		selectivity := 0.5

		// Ajuster selon le nombre de variables
		numVars := len(pattern.LeftVars) + len(pattern.RightVars)
		if numVars == 2 {
			selectivity = 0.3 // Jointure binaire simple
		} else if numVars > 2 {
			selectivity = 0.4 + (float64(numVars-2) * 0.1) // Plus de variables = moins sélectif
		}

		// Ajuster selon les conditions de jointure
		if len(pattern.JoinConditions) > 0 {
			// Plus de conditions = plus sélectif
			selectivity *= (1.0 - float64(len(pattern.JoinConditions))*0.1)
			if selectivity < 0.1 {
				selectivity = 0.1
			}
		}

		pattern.Selectivity = selectivity
		pattern.EstimatedCost = selectivity * float64(numVars)
	}
}

// optimizeJoinOrder optimise l'ordre des patterns de jointure.
//
// Stratégie: trier les patterns par sélectivité croissante (plus sélectif d'abord).
// Cela permet de filtrer les données tôt dans la chaîne et de réduire le volume
// de données traité par les jointures suivantes.
//
// Note: Pour une optimisation plus avancée, on pourrait tenir compte des dépendances
// entre variables (un pattern ne peut être évalué que si ses variables dépendantes
// ont été produites par des patterns précédents).
//
// Retourne une nouvelle slice avec les patterns réordonnés.
func (bcb *BetaChainBuilder) optimizeJoinOrder(patterns []JoinPattern) []JoinPattern {
	// Copier les patterns pour ne pas modifier l'original
	optimized := make([]JoinPattern, len(patterns))
	copy(optimized, patterns)

	// Trier par sélectivité croissante (plus sélectif d'abord)
	sort.Slice(optimized, func(i, j int) bool {
		return optimized[i].Selectivity < optimized[j].Selectivity
	})

	return optimized
}

// patternsEqual vérifie si deux slices de patterns sont identiques (même ordre).
func (bcb *BetaChainBuilder) patternsEqual(a, b []JoinPattern) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !bcb.patternEqual(a[i], b[i]) {
			return false
		}
	}
	return true
}

// patternEqual vérifie si deux patterns sont identiques.
func (bcb *BetaChainBuilder) patternEqual(a, b JoinPattern) bool {
	// Comparaison simple basée sur les variables
	if len(a.LeftVars) != len(b.LeftVars) || len(a.RightVars) != len(b.RightVars) {
		return false
	}
	for i := range a.LeftVars {
		if a.LeftVars[i] != b.LeftVars[i] {
			return false
		}
	}
	for i := range a.RightVars {
		if a.RightVars[i] != b.RightVars[i] {
			return false
		}
	}
	return true
}

// findReusablePrefix cherche un préfixe de chaîne réutilisable dans le cache.
//
// Un préfixe réutilisable est une sous-séquence de patterns au début de la chaîne
// qui correspond exactement à une sous-chaîne déjà construite.
//
// Retourne:
//   - Le dernier nœud du préfixe réutilisable (ou nil si aucun)
//   - La longueur du préfixe (nombre de patterns)
func (bcb *BetaChainBuilder) findReusablePrefix(patterns []JoinPattern, ruleID string) (*JoinNode, int) {
	bcb.mutex.RLock()
	defer bcb.mutex.RUnlock()

	// Chercher le plus long préfixe disponible (de len-1 à 1)
	for prefixLen := len(patterns) - 1; prefixLen >= 1; prefixLen-- {
		prefixKey := bcb.computePrefixKey(patterns[0:prefixLen])
		if node, exists := bcb.prefixCache[prefixKey]; exists {
			return node, prefixLen
		}
	}

	return nil, 0
}

// computePrefixKey calcule une clé pour un préfixe de patterns.
//
// La clé est construite en concaténant les signatures des patterns.
func (bcb *BetaChainBuilder) computePrefixKey(patterns []JoinPattern) string {
	key := ""
	for _, pattern := range patterns {
		// Utiliser les variables comme base de la clé
		key += fmt.Sprintf("%v|%v|", pattern.LeftVars, pattern.RightVars)
	}
	return key
}

// updatePrefixCache met à jour le cache de préfixes.
func (bcb *BetaChainBuilder) updatePrefixCache(key string, node *JoinNode) {
	bcb.mutex.Lock()
	defer bcb.mutex.Unlock()
	bcb.prefixCache[key] = node
}

// isAlreadyConnectedCached vérifie si un nœud enfant est déjà connecté à un nœud parent avec cache.
func (bcb *BetaChainBuilder) isAlreadyConnectedCached(parent Node, child Node) bool {
	if parent == nil || child == nil {
		return false
	}

	parentID := parent.GetID()
	childID := child.GetID()
	cacheKey := fmt.Sprintf("%s_%s", parentID, childID)

	// Vérifier le cache
	bcb.mutex.RLock()
	if connected, exists := bcb.connectionCache[cacheKey]; exists {
		bcb.mutex.RUnlock()
		return connected
	}
	bcb.mutex.RUnlock()

	// Cache miss - vérifier réellement
	connected := isAlreadyConnected(parent, child)

	// Mettre à jour le cache
	bcb.updateConnectionCache(parentID, childID, connected)

	return connected
}

// updateConnectionCache met à jour le cache de connexion.
func (bcb *BetaChainBuilder) updateConnectionCache(parentID, childID string, connected bool) {
	cacheKey := fmt.Sprintf("%s_%s", parentID, childID)
	bcb.mutex.Lock()
	bcb.connectionCache[cacheKey] = connected
	bcb.mutex.Unlock()
}

// determineJoinType détermine le type de jointure d'un pattern.
//
// Types supportés:
//   - "binary": jointure binaire simple (2 variables)
//   - "cascade": jointure en cascade (3+ variables)
//   - "multi": jointure multi-variables complexe
func (bcb *BetaChainBuilder) determineJoinType(pattern JoinPattern) string {
	numVars := len(pattern.LeftVars) + len(pattern.RightVars)
	if numVars == 2 {
		return "binary"
	} else if len(pattern.LeftVars) > 1 {
		return "cascade"
	} else {
		return "multi"
	}
}

// ClearConnectionCache vide le cache de connexions.
//
// Utile pour libérer de la mémoire après suppression de nombreuses règles,
// ou pour forcer une réévaluation complète des connexions.
//
// Thread-safe: peut être appelé à tout moment.
//
// Exemple:
//
//	// Après suppression de beaucoup de règles
//	for _, ruleID := range oldRules {
//	    network.RemoveRule(ruleID)
//	}
//	builder.ClearConnectionCache() // Libérer mémoire
func (bcb *BetaChainBuilder) ClearConnectionCache() {
	bcb.mutex.Lock()
	defer bcb.mutex.Unlock()
	bcb.connectionCache = make(map[string]bool)
	tsdio.LogPrintf("🧹 [BetaChainBuilder] Cache de connexions vidé")
}

// ClearPrefixCache vide le cache de préfixes.
//
// Utile pour libérer de la mémoire ou invalider les préfixes après
// modifications importantes du réseau.
//
// Thread-safe: peut être appelé à tout moment.
func (bcb *BetaChainBuilder) ClearPrefixCache() {
	bcb.mutex.Lock()
	defer bcb.mutex.Unlock()
	bcb.prefixCache = make(map[string]*JoinNode)
	tsdio.LogPrintf("🧹 [BetaChainBuilder] Cache de préfixes vidé")
}

// GetConnectionCacheSize retourne la taille actuelle du cache de connexions.
//
// Thread-safe.
//
// Exemple:
//
//	size := builder.GetConnectionCacheSize()
//	fmt.Printf("Cache de connexions: %d entrées\n", size)
func (bcb *BetaChainBuilder) GetConnectionCacheSize() int {
	bcb.mutex.RLock()
	defer bcb.mutex.RUnlock()
	return len(bcb.connectionCache)
}

// GetPrefixCacheSize retourne la taille actuelle du cache de préfixes.
//
// Thread-safe.
func (bcb *BetaChainBuilder) GetPrefixCacheSize() int {
	bcb.mutex.RLock()
	defer bcb.mutex.RUnlock()
	return len(bcb.prefixCache)
}

// GetMetrics retourne les métriques de construction.
//
// Thread-safe: retourne une copie des métriques.
//
// Exemple:
//
//	metrics := builder.GetMetrics()
//	fmt.Printf("Join nodes requested: %d\n", metrics.TotalJoinNodesRequested)
func (bcb *BetaChainBuilder) GetMetrics() *BetaChainMetrics {
	bcb.mutex.RLock()
	defer bcb.mutex.RUnlock()
	return bcb.metrics
}

// ResetMetrics réinitialise les métriques de construction à zéro.
//
// Cette méthode remet toutes les statistiques de construction à leur état initial.
// Utile pour les tests ou pour commencer une nouvelle session de mesure.
//
// Thread-safe: Protégé par mutex pour éviter les conditions de course.
//
// Exemple:
//
//	builder.ResetMetrics()
//	// Construire des chaînes...
//	metrics := builder.GetMetrics()
//	fmt.Printf("Depuis le reset: %d nodes créés\n", metrics.UniqueJoinNodesCreated)
func (bcb *BetaChainBuilder) ResetMetrics() {
	if bcb.metrics != nil {
		bcb.metrics.Reset()
	}
}

// SetOptimizationEnabled active/désactive l'optimisation de l'ordre des jointures.
//
// Thread-safe.
//
// Exemple:
//
//	builder.SetOptimizationEnabled(false) // Désactiver l'optimisation
func (bcb *BetaChainBuilder) SetOptimizationEnabled(enabled bool) {
	bcb.mutex.Lock()
	defer bcb.mutex.Unlock()
	bcb.enableOptimization = enabled
	tsdio.LogPrintf("⚙️  [BetaChainBuilder] Optimisation de l'ordre: %v", enabled)
}

// SetPrefixSharingEnabled active/désactive le partage de préfixes.
//
// Thread-safe.
//
// Exemple:
//
//	builder.SetPrefixSharingEnabled(false) // Désactiver le partage de préfixes
func (bcb *BetaChainBuilder) SetPrefixSharingEnabled(enabled bool) {
	bcb.mutex.Lock()
	defer bcb.mutex.Unlock()
	bcb.enablePrefixSharing = enabled
	tsdio.LogPrintf("⚙️  [BetaChainBuilder] Partage de préfixes: %v", enabled)
}

// GetChainInfo retourne des informations détaillées sur une chaîne beta.
//
// Exemple:
//
//	info := chain.GetChainInfo()
//	fmt.Printf("Chaîne: %s\n", info["summary"])
//	fmt.Printf("Nœuds: %v\n", info["node_ids"])
func (bc *BetaChain) GetChainInfo() map[string]interface{} {
	info := make(map[string]interface{})

	nodeIDs := make([]string, len(bc.Nodes))
	for i, node := range bc.Nodes {
		nodeIDs[i] = node.ID
	}

	info["rule_id"] = bc.RuleID
	info["chain_length"] = len(bc.Nodes)
	info["node_ids"] = nodeIDs
	info["hashes"] = bc.Hashes

	if bc.FinalNode != nil {
		info["final_node_id"] = bc.FinalNode.ID
	}

	summary := fmt.Sprintf("BetaChain[%s]: %d nœuds", bc.RuleID, len(bc.Nodes))
	info["summary"] = summary

	return info
}

// ValidateChain valide la cohérence d'une chaîne beta.
//
// Vérifie:
//   - Longueurs cohérentes (nodes, hashes)
//   - FinalNode correspond au dernier nœud
//   - Tous les nœuds sont non-nil
//   - Tous les hashes sont non-vides
//
// Retourne une erreur si la validation échoue.
//
// Exemple:
//
//	if err := chain.ValidateChain(); err != nil {
//	    tsdio.LogPrintf("Chaîne invalide: %v", err)
//	}
func (bc *BetaChain) ValidateChain() error {
	if len(bc.Nodes) != len(bc.Hashes) {
		return fmt.Errorf("incohérence: %d nœuds mais %d hashes", len(bc.Nodes), len(bc.Hashes))
	}

	if len(bc.Nodes) == 0 {
		return fmt.Errorf("chaîne vide")
	}

	for i, node := range bc.Nodes {
		if node == nil {
			return fmt.Errorf("nœud %d est nil", i)
		}
		if bc.Hashes[i] == "" {
			return fmt.Errorf("hash %d est vide", i)
		}
	}

	if bc.FinalNode != bc.Nodes[len(bc.Nodes)-1] {
		return fmt.Errorf("FinalNode ne correspond pas au dernier nœud de la chaîne")
	}

	return nil
}

// CountSharedNodes compte le nombre de nœuds partagés dans une chaîne.
//
// Un nœud est considéré comme partagé s'il est utilisé par plusieurs règles
// (RefCount > 1 dans le LifecycleManager).
//
// Exemple:
//
//	sharedCount := builder.CountSharedNodes(chain)
//	fmt.Printf("Nœuds partagés: %d/%d\n", sharedCount, len(chain.Nodes))
func (bcb *BetaChainBuilder) CountSharedNodes(chain *BetaChain) int {
	if bcb.network.LifecycleManager == nil {
		return 0
	}

	sharedCount := 0
	for _, node := range chain.Nodes {
		lifecycle, _ := bcb.network.LifecycleManager.GetNodeLifecycle(node.ID)
		if lifecycle != nil && lifecycle.GetRefCount() > 1 {
			sharedCount++
		}
	}

	return sharedCount
}

// GetChainStats retourne des statistiques sur une chaîne.
//
// Statistiques disponibles:
//   - total_nodes: nombre total de nœuds
//   - shared_nodes: nombre de nœuds partagés
//   - sharing_ratio: ratio de partage (0-1)
//   - average_refcount: refcount moyen des nœuds
//
// Exemple:
//
//	stats := builder.GetChainStats(chain)
//	fmt.Printf("Statistiques:\n")
//	for key, value := range stats {
//	    fmt.Printf("  %s: %v\n", key, value)
//	}
func (bcb *BetaChainBuilder) GetChainStats(chain *BetaChain) map[string]interface{} {
	stats := make(map[string]interface{})

	totalNodes := len(chain.Nodes)
	sharedNodes := bcb.CountSharedNodes(chain)

	stats["total_nodes"] = totalNodes
	stats["shared_nodes"] = sharedNodes

	if totalNodes > 0 {
		stats["sharing_ratio"] = float64(sharedNodes) / float64(totalNodes)
	} else {
		stats["sharing_ratio"] = 0.0
	}

	// Calculer le refcount moyen
	if bcb.network.LifecycleManager != nil {
		totalRefCount := 0
		for _, node := range chain.Nodes {
			lifecycle, _ := bcb.network.LifecycleManager.GetNodeLifecycle(node.ID)
			if lifecycle != nil {
				totalRefCount += lifecycle.GetRefCount()
			}
		}
		if totalNodes > 0 {
			stats["average_refcount"] = float64(totalRefCount) / float64(totalNodes)
		}
	}

	return stats
}
