// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"sync"
)

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
	// Delegate to orchestrated version
	return bcb.BuildChainOrchestrated(patterns, ruleID)
}

// GetMetrics retourne les métriques de construction de chaînes.
//
// Thread-safe: Utilise RLock pour lecture concurrent-safe.
//
// Retourne:
//   - Pointeur vers l'objet de métriques partagé
//
// Exemple:
//
//	metrics := builder.GetMetrics()
//	fmt.Printf("Nœuds créés: %d\n", metrics.UniqueJoinNodesCreated)
//	fmt.Printf("Nœuds réutilisés: %d\n", metrics.JoinNodesReused)
//	fmt.Printf("Ratio de partage: %.2f%%\n", metrics.SharingRatio()*100)
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
// Quand activée, BuildChain réordonne les patterns par sélectivité croissante
// pour optimiser les performances d'évaluation.
//
// Thread-safe: Protégé par mutex.
//
// Paramètres:
//   - enabled: true pour activer, false pour désactiver
//
// Exemple:
//
//	builder.SetOptimizationEnabled(false) // Désactiver l'optimisation
//	chain, _ := builder.BuildChain(patterns, "rule1") // Ordre original préservé
func (bcb *BetaChainBuilder) SetOptimizationEnabled(enabled bool) {
	bcb.mutex.Lock()
	defer bcb.mutex.Unlock()
	bcb.enableOptimization = enabled
	fmt.Printf("⚙️  [BetaChainBuilder] Optimisation de l'ordre: %v\n", enabled)
}

// SetPrefixSharingEnabled active/désactive le partage de préfixes.
//
// Quand activé, BuildChain cherche des préfixes de chaînes réutilisables
// pour éviter de reconstruire des sous-séquences de patterns identiques.
//
// Thread-safe: Protégé par mutex.
//
// Paramètres:
//   - enabled: true pour activer, false pour désactiver
//
// Exemple:
//
//	builder.SetPrefixSharingEnabled(false) // Désactiver le partage de préfixes
//	chain, _ := builder.BuildChain(patterns, "rule1") // Pas de réutilisation de préfixes
func (bcb *BetaChainBuilder) SetPrefixSharingEnabled(enabled bool) {
	bcb.mutex.Lock()
	defer bcb.mutex.Unlock()
	bcb.enablePrefixSharing = enabled
	fmt.Printf("⚙️  [BetaChainBuilder] Partage de préfixes: %v\n", enabled)
}
