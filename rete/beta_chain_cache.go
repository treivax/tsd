// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
)

// isAlreadyConnectedCached vérifie si un nœud enfant est déjà connecté à un parent
// en utilisant le cache de connexions pour améliorer les performances.
//
// Cette fonction utilise un cache à deux niveaux:
//  1. Vérifie d'abord le cache local (map thread-safe)
//  2. Si cache miss, appelle isAlreadyConnected() et met en cache le résultat
//
// Le cache évite les parcours répétés de la liste des enfants lors de la
// construction de plusieurs chaînes partageant des préfixes communs.
//
// Paramètres:
//   - parent: Le nœud parent à vérifier
//   - child: Le nœud enfant à rechercher
//
// Retourne:
//   - true si l'enfant est déjà dans la liste des enfants du parent
//   - false sinon (ou si parent/child est nil)
//
// Thread-safety:
//   - Utilise RWMutex pour accès concurrent sécurisé au cache
//   - Lecture optimiste avec lock partagé
//   - Écriture avec lock exclusif
//
// Exemple:
//
//	if builder.isAlreadyConnectedCached(typeNode, joinNode) {
//	    fmt.Println("Connexion déjà établie, skip")
//	} else {
//	    parent.AddChild(child)
//	}
func (bcb *BetaChainBuilder) isAlreadyConnectedCached(parent Node, child Node) bool {
	if parent == nil || child == nil {
		return false
	}

	parentID := parent.GetID()
	childID := child.GetID()
	cacheKey := fmt.Sprintf("%s_%s", parentID, childID)

	// Vérifier le cache (lecture optimiste)
	bcb.mutex.RLock()
	if connected, exists := bcb.connectionCache[cacheKey]; exists {
		bcb.mutex.RUnlock()
		return connected
	}
	bcb.mutex.RUnlock()

	// Cache miss - vérifier réellement dans la structure du nœud
	connected := isAlreadyConnected(parent, child)

	// Mettre à jour le cache avec le résultat
	bcb.updateConnectionCache(parentID, childID, connected)

	return connected
}

// updateConnectionCache met à jour le cache de connexion avec une nouvelle entrée.
//
// Cette fonction enregistre l'état de connexion entre deux nœuds dans le cache
// pour éviter des vérifications répétées coûteuses.
//
// Paramètres:
//   - parentID: ID du nœud parent
//   - childID: ID du nœud enfant
//   - connected: État de connexion (true si connectés)
//
// Thread-safety:
//   - Utilise mutex.Lock() pour écriture exclusive
//   - Garantit cohérence du cache en environnement concurrent
//
// Note:
//   - Le cache n'est jamais invalidé automatiquement
//   - Utiliser ClearConnectionCache() pour réinitialisation manuelle
//   - Taille du cache croît linéairement avec le nombre de connexions uniques
func (bcb *BetaChainBuilder) updateConnectionCache(parentID, childID string, connected bool) {
	cacheKey := fmt.Sprintf("%s_%s", parentID, childID)
	bcb.mutex.Lock()
	bcb.connectionCache[cacheKey] = connected
	bcb.mutex.Unlock()
}

// ClearConnectionCache vide complètement le cache de connexions.
//
// Utile pour libérer de la mémoire après suppression de nombreuses règles,
// ou pour forcer une réévaluation complète des connexions.
//
// Cas d'usage:
//   - Après suppression massive de règles (> 100 règles)
//   - Quand le cache grandit trop (> 10000 entrées)
//   - Avant reconstruction complète du réseau
//   - Tests unitaires nécessitant état propre
//
// Thread-safety:
//   - Utilise mutex.Lock() pour garantir atomicité
//   - Peut être appelé à tout moment sans risque
//
// Performance:
//   - O(1) - crée simplement une nouvelle map vide
//   - L'ancienne map est garbage collected automatiquement
//
// Exemple:
//
//	// Après suppression de beaucoup de règles
//	for _, ruleID := range oldRules {
//	    network.RemoveRule(ruleID)
//	}
//	builder.ClearConnectionCache() // Libérer mémoire du cache
func (bcb *BetaChainBuilder) ClearConnectionCache() {
	bcb.mutex.Lock()
	defer bcb.mutex.Unlock()
	bcb.connectionCache = make(map[string]bool)
	fmt.Printf("🧹 [BetaChainBuilder] Cache de connexions vidé\n")
}

// ClearPrefixCache vide complètement le cache de préfixes de chaînes.
//
// Le cache de préfixes stocke les JoinNodes réutilisables pour optimiser
// la construction de chaînes avec des préfixes communs. Vider ce cache
// force le builder à reconstruire les préfixes depuis zéro.
//
// Cas d'usage:
//   - Après modifications importantes du réseau
//   - Quand les patterns de règles changent radicalement
//   - Tests unitaires nécessitant état propre
//   - Débogage de problèmes de partage de préfixes
//
// Thread-safety:
//   - Utilise mutex.Lock() pour garantir atomicité
//   - Peut être appelé à tout moment sans risque
//
// Performance:
//   - O(1) - crée simplement une nouvelle map vide
//   - L'ancienne map est garbage collected automatiquement
//
// Impact:
//   - Construction de chaînes plus lente temporairement
//   - Pas d'impact sur les chaînes déjà construites
//   - Préfixes seront recalculés à la demande
//
// Exemple:
//
//	// Après modification majeure du réseau
//	network.RemoveRule("old_rule")
//	builder.ClearPrefixCache() // Invalider préfixes obsolètes
func (bcb *BetaChainBuilder) ClearPrefixCache() {
	bcb.mutex.Lock()
	defer bcb.mutex.Unlock()
	bcb.prefixCache = make(map[string]*JoinNode)
	fmt.Printf("🧹 [BetaChainBuilder] Cache de préfixes vidé\n")
}

// GetConnectionCacheSize retourne la taille actuelle du cache de connexions.
//
// Cette métrique indique combien de paires parent-enfant ont été vérifiées
// et mises en cache depuis la création du builder ou le dernier clear.
//
// Thread-safety:
//   - Utilise RLock pour lecture concurrent-safe
//   - Retourne snapshot instantané de la taille
//
// Utilisation:
//   - Monitoring de l'utilisation mémoire
//   - Décision de nettoyage du cache
//   - Métriques de performance
//   - Débogage
//
// Exemple:
//
//	size := builder.GetConnectionCacheSize()
//	fmt.Printf("Cache de connexions: %d entrées\n", size)
//	if size > 10000 {
//	    builder.ClearConnectionCache()
//	}
func (bcb *BetaChainBuilder) GetConnectionCacheSize() int {
	bcb.mutex.RLock()
	defer bcb.mutex.RUnlock()
	return len(bcb.connectionCache)
}

// GetPrefixCacheSize retourne la taille actuelle du cache de préfixes.
//
// Cette métrique indique combien de préfixes de chaînes ont été identifiés
// et mis en cache pour réutilisation future.
//
// Thread-safety:
//   - Utilise RLock pour lecture concurrent-safe
//   - Retourne snapshot instantané de la taille
//
// Utilisation:
//   - Monitoring de l'efficacité du partage de préfixes
//   - Métriques de performance
//   - Décision de nettoyage du cache
//   - Débogage
//
// Exemple:
//
//	size := builder.GetPrefixCacheSize()
//	fmt.Printf("Cache de préfixes: %d entrées\n", size)
func (bcb *BetaChainBuilder) GetPrefixCacheSize() int {
	bcb.mutex.RLock()
	defer bcb.mutex.RUnlock()
	return len(bcb.prefixCache)
}
