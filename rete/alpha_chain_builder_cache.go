// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

// Package rete fournit l'implémentation du réseau RETE pour l'évaluation de règles.
// Ce fichier contient les fonctions de gestion du cache de connexions pour l'AlphaChainBuilder,
// permettant d'éviter les vérifications coûteuses de connexions déjà établies.
package rete

import (
	"fmt"
)

// isAlreadyConnectedCached vérifie si un nœud enfant est déjà connecté à un nœud parent avec cache
func (acb *AlphaChainBuilder) isAlreadyConnectedCached(parent Node, child Node) bool {
	if parent == nil || child == nil {
		return false
	}

	parentID := parent.GetID()
	childID := child.GetID()
	cacheKey := fmt.Sprintf("%s_%s", parentID, childID)

	// Vérifier le cache
	acb.mutex.RLock()
	if connected, exists := acb.connectionCache[cacheKey]; exists {
		acb.mutex.RUnlock()
		if acb.metrics != nil {
			acb.metrics.RecordConnectionCacheHit()
		}
		return connected
	}
	acb.mutex.RUnlock()

	// Cache miss - vérifier réellement
	if acb.metrics != nil {
		acb.metrics.RecordConnectionCacheMiss()
	}

	connected := isAlreadyConnected(parent, child)

	// Mettre à jour le cache
	acb.updateConnectionCache(parentID, childID, connected)

	return connected
}

// updateConnectionCache met à jour le cache de connexion
func (acb *AlphaChainBuilder) updateConnectionCache(parentID, childID string, connected bool) {
	cacheKey := fmt.Sprintf("%s_%s", parentID, childID)
	acb.mutex.Lock()
	acb.connectionCache[cacheKey] = connected
	acb.mutex.Unlock()
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
func (acb *AlphaChainBuilder) ClearConnectionCache() {
	acb.mutex.Lock()
	defer acb.mutex.Unlock()
	acb.connectionCache = make(map[string]bool)
}

// GetConnectionCacheSize retourne la taille actuelle du cache de connexions.
//
// Utile pour monitoring et diagnostic de l'utilisation mémoire.
//
// Retourne:
//   - Nombre d'entrées dans le cache (une par connexion parent→child unique)
//
// Exemple:
//
//	size := builder.GetConnectionCacheSize()
//	fmt.Printf("Cache de connexions: %d entrées\n", size)
//	if size > 10000 {
//	    builder.ClearConnectionCache() // Nettoyage si trop grand
//	}
func (acb *AlphaChainBuilder) GetConnectionCacheSize() int {
	acb.mutex.RLock()
	defer acb.mutex.RUnlock()
	return len(acb.connectionCache)
}

// isAlreadyConnected vérifie si un nœud enfant est déjà connecté à un nœud parent
func isAlreadyConnected(parent Node, child Node) bool {
	if parent == nil || child == nil {
		return false
	}

	children := parent.GetChildren()
	childID := child.GetID()

	for _, c := range children {
		if c.GetID() == childID {
			return true
		}
	}

	return false
}
