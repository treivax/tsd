// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"fmt"
)

// IntegrationHelper facilite l'intégration du système delta avec RETE.
//
// Cette structure coordonne les composants du système delta (détecteur,
// index, propagateur) et fournit une API simplifiée pour l'intégration
// avec le réseau RETE.
//
// Thread-safety : IntegrationHelper est safe pour utilisation concurrente.
type IntegrationHelper struct {
	propagator *DeltaPropagator
	index      *DependencyIndex
	callbacks  NetworkCallbacks
	network    interface{} // Référence au ReteNetwork (interface{} pour éviter dépendance circulaire)
	builder    *IndexBuilder
}

// NewIntegrationHelper crée un nouveau helper d'intégration.
//
// Paramètres :
//   - propagator : propagateur delta configuré
//   - index : index de dépendances construit
//   - callbacks : callbacks vers le réseau RETE
//
// Retourne un nouveau helper prêt à l'emploi.
func NewIntegrationHelper(
	propagator *DeltaPropagator,
	index *DependencyIndex,
	callbacks NetworkCallbacks,
) *IntegrationHelper {
	return &IntegrationHelper{
		propagator: propagator,
		index:      index,
		callbacks:  callbacks,
		builder:    NewIndexBuilder(),
	}
}

// SetNetwork configure la référence au réseau RETE.
//
// Cette méthode doit être appelée pour permettre la reconstruction
// automatique de l'index lors de modifications du réseau.
//
// Paramètres :
//   - network : référence au ReteNetwork (typiquement *rete.ReteNetwork)
func (ih *IntegrationHelper) SetNetwork(network interface{}) {
	ih.network = network
}

// ProcessUpdate traite une mise à jour de fait de bout en bout.
//
// Cette méthode coordonne :
//  1. Détection du delta entre ancien et nouveau fait
//  2. Recherche des nœuds affectés via l'index
//  3. Décision delta vs classique (selon config et ratio)
//  4. Propagation vers les nœuds affectés
//  5. Mise à jour du storage
//
// Paramètres :
//   - oldFact : fait avant modification
//   - newFact : fait après modification
//   - factID : identifiant interne du fait
//   - factType : type du fait
//
// Retourne une erreur si le traitement échoue.
func (ih *IntegrationHelper) ProcessUpdate(
	oldFact, newFact map[string]interface{},
	factID, factType string,
) error {
	if ih.propagator == nil {
		return newComponentError("propagator", "ProcessUpdate", ErrMsgPropagatorNotInit)
	}

	if ih.callbacks == nil {
		return newComponentError("callbacks", "ProcessUpdate", ErrMsgCallbacksNotConfig)
	}

	// Déléguer à la propagation delta
	err := ih.propagator.PropagateUpdate(oldFact, newFact, factID, factType)
	if err != nil {
		return fmt.Errorf("%s: %w", ErrMsgDeltaPropagationFail, err)
	}

	// Mettre à jour le storage
	if err := ih.callbacks.UpdateStorage(factID, newFact); err != nil {
		return fmt.Errorf("%s: %w", ErrMsgStorageUpdateFail, err)
	}

	return nil
}

// RebuildIndex reconstruit l'index de dépendances depuis le réseau.
//
// Cette méthode doit être appelée si le réseau RETE est modifié
// dynamiquement (ajout/suppression de règles).
//
// Le réseau doit être configuré via SetNetwork() avant l'appel à cette méthode.
//
// Retourne une erreur si :
//   - L'index n'est pas initialisé
//   - Le réseau n'est pas configuré
//   - La reconstruction échoue
func (ih *IntegrationHelper) RebuildIndex() error {
	if ih.index == nil {
		return newComponentError("index", "RebuildIndex", ErrMsgIndexNotInit)
	}

	if ih.network == nil {
		return newComponentError("network", "RebuildIndex", "network not configured - use SetNetwork()")
	}

	// Utiliser le builder pour reconstruire depuis le réseau
	newIndex, err := ih.builder.BuildFromNetwork(ih.network)
	if err != nil {
		return fmt.Errorf("failed to rebuild index from network: %w", err)
	}

	// Remplacer l'index actuel par le nouvel index
	ih.index = newIndex

	return nil
}

// GetMetrics retourne les métriques du propagateur.
//
// Utile pour monitoring et diagnostic des performances de propagation.
func (ih *IntegrationHelper) GetMetrics() PropagationMetrics {
	if ih.propagator == nil {
		return PropagationMetrics{}
	}
	return ih.propagator.GetMetrics()
}

// GetIndexMetrics retourne les statistiques de l'index.
//
// Note: Retourne IndexStats qui contient les statistiques de structure,
// pas IndexMetrics qui concerne les performances d'utilisation.
func (ih *IntegrationHelper) GetIndexMetrics() IndexStats {
	if ih.index == nil {
		return IndexStats{}
	}
	return ih.index.GetStats()
}

// EnableDiagnostics active les diagnostics pour tous les composants.
//
// En mode diagnostic, les composants collectent des informations
// supplémentaires utiles pour le debugging.
//
// Note: Actuellement, seul le propagateur supporte les diagnostics.
// L'index n'a pas de mode diagnostic distinct.
func (ih *IntegrationHelper) EnableDiagnostics() {
	// L'index n'a pas de mode diagnostics à activer
	// Les stats sont toujours collectées via GetStats()
}

// DisableDiagnostics désactive les diagnostics.
func (ih *IntegrationHelper) DisableDiagnostics() {
	// L'index n'a pas de mode diagnostics à désactiver
}
