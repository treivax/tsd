// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"sort"
)

// Constantes pour l'estimation de sélectivité des jointures
const (
	// DefaultSelectivity est la sélectivité par défaut si aucune information disponible
	DefaultSelectivity = 0.5

	// BinaryJoinSelectivity est la sélectivité estimée pour une jointure binaire simple (2 variables)
	BinaryJoinSelectivity = 0.3

	// BaseMultiVarSelectivity est la sélectivité de base pour les jointures multi-variables (3+)
	BaseMultiVarSelectivity = 0.4

	// SelectivityIncrementPerVar est l'incrément de sélectivité par variable supplémentaire
	SelectivityIncrementPerVar = 0.1

	// ConditionSelectivityFactor est le facteur de réduction par condition de jointure
	ConditionSelectivityFactor = 0.1

	// MinSelectivity est la sélectivité minimale autorisée
	MinSelectivity = 0.1

	// MinVariablesForBinaryJoin est le nombre de variables pour une jointure binaire
	MinVariablesForBinaryJoin = 2

	// MinReferencesForSharing est le nombre minimum de références pour considérer un nœud comme partagé
	MinReferencesForSharing = 1
)

// estimateSelectivity estime la sélectivité de chaque pattern de jointure.
//
// La sélectivité est une heuristique (0-1) qui indique combien de tuples
// passeront le filtre de jointure. Plus la valeur est basse, plus la jointure
// est sélective (filtre beaucoup de données).
//
// Heuristiques utilisées:
//   - Nombre de variables: 2 variables = plus sélectif que 3+
//   - Nombre de conditions: plus de conditions = plus sélectif
//   - Type d'opérateur: égalité > inégalité > range
//
// Cette estimation est ensuite utilisée par optimizeJoinOrder pour
// choisir l'ordre optimal d'évaluation des jointures.
//
// Paramètres:
//   - patterns: Liste de patterns à estimer (modifiés en place)
//
// Comportement:
//   - Ne fait rien si la sélectivité est déjà définie (> 0)
//   - Définit Selectivity et EstimatedCost pour chaque pattern
//
// Exemple:
//
//	patterns := []JoinPattern{
//	    {LeftVars: []string{"p"}, RightVars: []string{"o"}},
//	    {LeftVars: []string{"p","o"}, RightVars: []string{"pay"}},
//	}
//	builder.estimateSelectivity(patterns)
//	// patterns[0].Selectivity ≈ 0.3 (binaire)
//	// patterns[1].Selectivity ≈ 0.4 (cascade)
func (bcb *BetaChainBuilder) estimateSelectivity(patterns []JoinPattern) {
	for i := range patterns {
		pattern := &patterns[i]

		// Si déjà estimée, ne rien faire
		if pattern.Selectivity > 0 {
			continue
		}

		// Estimation par défaut
		selectivity := DefaultSelectivity

		// Ajuster selon le nombre de variables
		numVars := len(pattern.LeftVars) + len(pattern.RightVars)
		if numVars == MinVariablesForBinaryJoin {
			selectivity = BinaryJoinSelectivity // Jointure binaire simple (plus sélective)
		} else if numVars > MinVariablesForBinaryJoin {
			// Plus de variables = moins sélectif (plus de données passent)
			selectivity = BaseMultiVarSelectivity + (float64(numVars-MinVariablesForBinaryJoin) * SelectivityIncrementPerVar)
		}

		// Ajuster selon les conditions de jointure
		if len(pattern.JoinConditions) > 0 {
			// Plus de conditions = plus sélectif
			selectivity *= (1.0 - float64(len(pattern.JoinConditions))*ConditionSelectivityFactor)
			if selectivity < MinSelectivity {
				selectivity = MinSelectivity // Minimum de sélectivité
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
// Principe d'optimisation:
//   - Patterns sélectifs (selectivity faible) en premier → filtrent beaucoup
//   - Patterns moins sélectifs ensuite → traitent moins de données
//   - Résultat: réduction du coût total d'évaluation
//
// Note: Pour une optimisation plus avancée, on pourrait tenir compte des dépendances
// entre variables (un pattern ne peut être évalué que si ses variables dépendantes
// ont été produites par des patterns précédents).
//
// Paramètres:
//   - patterns: Liste originale de patterns
//
// Retourne:
//   - Une nouvelle slice avec les patterns réordonnés (copie)
//
// Exemple:
//
//	patterns := []JoinPattern{
//	    {Selectivity: 0.5}, // Moins sélectif
//	    {Selectivity: 0.2}, // Plus sélectif
//	    {Selectivity: 0.4},
//	}
//	optimized := builder.optimizeJoinOrder(patterns)
//	// Ordre résultant: [0.2, 0.4, 0.5]
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
//
// Comparaison structurelle:
//   - Même longueur
//   - Chaque pattern à la même position est égal (via patternEqual)
//
// Utilisé pour détecter si l'optimisation a réordonné les patterns.
//
// Paramètres:
//   - a, b: Listes de patterns à comparer
//
// Retourne:
//   - true si les listes sont identiques
//   - false sinon (longueurs différentes ou patterns différents)
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
//
// Comparaison basée sur les variables uniquement:
//   - LeftVars doivent être identiques (même ordre)
//   - RightVars doivent être identiques (même ordre)
//
// Note: Ne compare pas les conditions ou la sélectivité pour simplifier.
// Cela suffit pour détecter les réordonnancement dus à l'optimisation.
//
// Paramètres:
//   - a, b: Patterns à comparer
//
// Retourne:
//   - true si les patterns ont les mêmes variables
//   - false sinon
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
// Stratégie de recherche:
//   - Cherche du plus long au plus court (len-1 à 1)
//   - S'arrête au premier préfixe trouvé
//   - Utilise computePrefixKey pour générer les clés de cache
//
// Bénéfices:
//   - Réutilise des sous-chaînes existantes
//   - Évite reconstruction complète
//   - Améliore performance sur règles similaires
//
// Paramètres:
//   - patterns: Liste complète de patterns
//   - ruleID: ID de la règle (pour logging)
//
// Retourne:
//   - Le dernier nœud du préfixe réutilisable (ou nil si aucun)
//   - La longueur du préfixe (nombre de patterns, 0 si aucun)
//
// Thread-safety:
//   - Utilise RLock pour lecture concurrent-safe du cache
//
// Exemple:
//
//	// Règle 1 construit: p ⋈ o ⋈ pay
//	// Règle 2 veut: p ⋈ o ⋈ item
//	node, len := builder.findReusablePrefix(patterns, "rule2")
//	// Retourne: le JoinNode de (p ⋈ o), len=2
func (bcb *BetaChainBuilder) findReusablePrefix(patterns []JoinPattern, ruleID string) (*JoinNode, int) {
	bcb.mutex.RLock()
	defer bcb.mutex.RUnlock()

	// Chercher le plus long préfixe disponible (de len-1 à 1)
	for prefixLen := len(patterns) - 1; prefixLen >= 1; prefixLen-- {
		prefixKey := bcb.computePrefixKey(patterns[0:prefixLen], ruleID)
		if node, exists := bcb.prefixCache[prefixKey]; exists {
			return node, prefixLen
		}
	}

	return nil, 0
}

// computePrefixKey calcule une clé pour un préfixe de patterns.
//
// La clé est construite en concaténant le ruleID et les signatures des patterns.
// Chaque pattern contribue ses LeftVars et RightVars séparés par '|'.
//
// Format de clé:
//
//	ruleID::[left1,left2]|[right1,right2]|[left3]|[right3]|...
//
// Propriétés:
//   - Déterministe: même préfixe → même clé
//   - Unique: préfixes différents → clés différentes
//   - Ordre-dépendant: [p,o] ≠ [o,p]
//   - Règle-spécifique: préfixes de règles différentes ne sont pas partagés
//
// Paramètres:
//   - patterns: Sous-séquence de patterns (préfixe)
//   - ruleID: Identifiant de la règle pour éviter le partage entre règles
//
// Retourne:
//   - Clé de cache sous forme de string
//
// Exemple:
//
//	patterns := []JoinPattern{
//	    {LeftVars: []string{"p"}, RightVars: []string{"o"}},
//	    {LeftVars: []string{"p","o"}, RightVars: []string{"pay"}},
//	}
//	key := builder.computePrefixKey(patterns, "rule1")
//	// key = "rule1::[p]|[o]|[p o]|[pay]|"
func (bcb *BetaChainBuilder) computePrefixKey(patterns []JoinPattern, ruleID string) string {
	// Inclure ruleID pour éviter le partage de préfixes entre règles différentes
	key := ruleID + "::"
	for _, pattern := range patterns {
		// Utiliser les variables comme base de la clé
		key += fmt.Sprintf("%v|%v|", pattern.LeftVars, pattern.RightVars)
	}
	return key
}

// updatePrefixCache met à jour le cache de préfixes.
//
// Enregistre un JoinNode comme le dernier nœud d'un préfixe de chaîne
// pour réutilisation future par d'autres règles.
//
// Paramètres:
//   - key: Clé de cache (généré par computePrefixKey)
//   - node: JoinNode à associer à ce préfixe
//
// Thread-safety:
//   - Utilise Lock pour écriture exclusive
//
// Note:
//   - Le cache grandit indéfiniment
//   - Utiliser ClearPrefixCache() pour libérer mémoire
func (bcb *BetaChainBuilder) updatePrefixCache(key string, node *JoinNode) {
	bcb.mutex.Lock()
	defer bcb.mutex.Unlock()
	bcb.prefixCache[key] = node
}

// CountSharedNodes compte le nombre de nœuds partagés dans une chaîne.
//
// Un nœud est considéré comme partagé s'il est utilisé par plusieurs règles
// (RefCount > 1 dans le LifecycleManager).
//
// Cette métrique indique l'efficacité du partage de nœuds:
//   - 0 nœuds partagés: aucune réutilisation
//   - Tous partagés: excellente réutilisation
//
// Paramètres:
//   - chain: Chaîne beta à analyser
//
// Retourne:
//   - Nombre de nœuds avec RefCount > 1
//   - 0 si LifecycleManager non initialisé
//
// Exemple:
//
//	sharedCount := builder.CountSharedNodes(chain)
//	totalNodes := len(chain.Nodes)
//	ratio := float64(sharedCount) / float64(totalNodes)
//	fmt.Printf("Partage: %.1f%% (%d/%d)\n", ratio*100, sharedCount, totalNodes)
func (bcb *BetaChainBuilder) CountSharedNodes(chain *BetaChain) int {
	if bcb.network.LifecycleManager == nil {
		return 0
	}

	sharedCount := 0
	for _, node := range chain.Nodes {
		lifecycle, _ := bcb.network.LifecycleManager.GetNodeLifecycle(node.ID)
		if lifecycle != nil && lifecycle.GetRefCount() > MinReferencesForSharing {
			sharedCount++
		}
	}

	return sharedCount
}

// GetChainStats retourne des statistiques complètes sur une chaîne.
//
// Statistiques disponibles:
//   - total_nodes: nombre total de nœuds dans la chaîne
//   - shared_nodes: nombre de nœuds partagés (RefCount > 1)
//   - sharing_ratio: ratio de partage (0.0 à 1.0)
//   - average_refcount: refcount moyen des nœuds
//
// Utilisation:
//   - Monitoring de l'efficacité du partage
//   - Métriques de performance
//   - Débogage et optimisation
//   - Reporting
//
// Paramètres:
//   - chain: Chaîne beta à analyser
//
// Retourne:
//   - Map avec toutes les statistiques
//
// Exemple:
//
//	stats := builder.GetChainStats(chain)
//	fmt.Printf("Statistiques de la chaîne:\n")
//	fmt.Printf("  Nœuds totaux: %d\n", stats["total_nodes"])
//	fmt.Printf("  Nœuds partagés: %d\n", stats["shared_nodes"])
//	fmt.Printf("  Ratio de partage: %.2f\n", stats["sharing_ratio"])
//	fmt.Printf("  RefCount moyen: %.2f\n", stats["average_refcount"])
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

// determineJoinType détermine le type de jointure d'un pattern.
//
// Types supportés:
//   - "binary": jointure binaire simple (2 variables)
//   - "cascade": jointure en cascade (gauche multi-variables)
//   - "multi": jointure multi-variables complexe
//
// Classification:
//   - binary: exactement 2 variables au total (1 gauche + 1 droite)
//   - cascade: plus d'une variable à gauche (jointure cumulative)
//   - multi: autres cas (complexe)
//
// Paramètres:
//   - pattern: Pattern de jointure à classifier
//
// Retourne:
//   - "binary", "cascade", ou "multi"
//
// Exemple:
//
//	// Pattern {p} ⋈ {o}
//	joinType := builder.determineJoinType(pattern)
//	// Retourne: "binary"
//
//	// Pattern {p,o} ⋈ {pay}
//	joinType := builder.determineJoinType(pattern)
//	// Retourne: "cascade"
func (bcb *BetaChainBuilder) determineJoinType(pattern JoinPattern) string {
	numVars := len(pattern.LeftVars) + len(pattern.RightVars)
	if numVars == MinVariablesForBinaryJoin {
		return "binary"
	} else if len(pattern.LeftVars) > MinReferencesForSharing {
		return "cascade"
	} else {
		return "multi"
	}
}
