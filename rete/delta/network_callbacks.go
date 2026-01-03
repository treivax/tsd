// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

// NetworkCallbacks définit les callbacks pour interagir avec le réseau RETE.
//
// Cette interface découple le package delta du package rete principal,
// évitant ainsi les dépendances circulaires.
//
// Les implémentations de cette interface doivent être thread-safe car
// elles peuvent être appelées concurremment par plusieurs goroutines.
type NetworkCallbacks interface {
	// PropagateToNode propage un delta vers un nœud RETE.
	// Le type de nœud (alpha, beta, terminal) est déterminé automatiquement.
	//
	// Paramètres :
	//   - nodeID : identifiant du nœud cible
	//   - delta : delta de fait à propager
	//
	// Retourne une erreur si la propagation échoue.
	PropagateToNode(nodeID string, delta *FactDelta) error

	// GetNode récupère un nœud par son ID.
	//
	// Paramètres :
	//   - nodeID : identifiant du nœud
	//
	// Retourne le nœud (interface{}) et nil, ou nil et une erreur si non trouvé.
	GetNode(nodeID string) (interface{}, error)

	// UpdateStorage met à jour le storage avec le fait modifié.
	//
	// Paramètres :
	//   - factID : identifiant interne du fait
	//   - newFact : nouveau contenu du fait (champs)
	//
	// Retourne une erreur si la mise à jour échoue.
	UpdateStorage(factID string, newFact map[string]interface{}) error
}

// DefaultNetworkCallbacks est une implémentation no-op pour tests.
//
// Cette implémentation ne fait rien et retourne toujours nil.
// Utile pour les tests unitaires du package delta qui n'ont pas besoin
// d'un réseau RETE réel.
type DefaultNetworkCallbacks struct{}

// PropagateToNode ne fait rien.
func (dnc *DefaultNetworkCallbacks) PropagateToNode(nodeID string, delta *FactDelta) error {
	return nil
}

// GetNode retourne toujours nil.
func (dnc *DefaultNetworkCallbacks) GetNode(nodeID string) (interface{}, error) {
	return nil, nil
}

// UpdateStorage ne fait rien.
func (dnc *DefaultNetworkCallbacks) UpdateStorage(factID string, newFact map[string]interface{}) error {
	return nil
}
