// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package xuples

import "sync"

// DefaultXupleSpace implémente XupleSpace.
type DefaultXupleSpace struct {
	name   string
	config XupleSpaceConfig
	xuples map[string]*Xuple // xupleID -> Xuple
	mu     sync.RWMutex
}

// NewXupleSpace crée un nouveau xuple-space.
func NewXupleSpace(config XupleSpaceConfig) XupleSpace {
	return &DefaultXupleSpace{
		name:   config.Name,
		config: config,
		xuples: make(map[string]*Xuple),
	}
}

// Name retourne le nom du xuple-space.
func (xs *DefaultXupleSpace) Name() string {
	return xs.name
}

// Insert insère un xuple dans le xuple-space.
//
// Validation :
//   - Le xuple ne peut pas être nil
//   - Le xuple doit avoir un ID (généré par XupleManager)
//   - Si MaxSize > 0, vérifie que la capacité n'est pas atteinte
//
// Side-effects :
//   - Applique la politique de rétention (calcul ExpiresAt)
//   - Modifie xuple.Metadata.ExpiresAt
//
// Thread-Safety :
//   - Méthode thread-safe (protégée par mutex)
func (xs *DefaultXupleSpace) Insert(xuple *Xuple) error {
	if xuple == nil {
		return ErrNilXuple
	}

	// L'ID doit être généré par le XupleManager
	if xuple.ID == "" {
		return ErrInvalidConfiguration
	}

	xs.mu.Lock()
	defer xs.mu.Unlock()

	// Vérifier la capacité maximale si définie
	if xs.config.MaxSize > 0 && len(xs.xuples) >= xs.config.MaxSize {
		return ErrXupleSpaceFull
	}

	// Appliquer la politique de rétention
	xuple.Metadata.ExpiresAt = xs.config.RetentionPolicy.ComputeExpiration(xuple.CreatedAt)

	xs.xuples[xuple.ID] = xuple
	return nil
}

// Retrieve récupère un xuple pour un agent selon les politiques.
//
// Process :
//  1. Marque les xuples expirés (modification de Metadata.State si nécessaire)
//  2. Filtre les xuples disponibles pour cet agent
//  3. Sélectionne un xuple selon la SelectionPolicy
//  4. MARQUE AUTOMATIQUEMENT LE XUPLE COMME CONSOMMÉ par cet agent
//
// Side-effects :
//   - Peut modifier Metadata.State pour marquer les xuples comme expirés
//   - MARQUE LE XUPLE SÉLECTIONNÉ COMME CONSOMMÉ (appelle markConsumedBy)
//   - Peut changer l'état du xuple à XupleStateConsumed selon la ConsumptionPolicy
//   - Toutes les modifications sont thread-safe (protégées par mutex)
//
// Thread-Safety :
//   - Méthode thread-safe (protégée par mutex)
//
// Note : Cette méthode implémente une sémantique "retrieve-and-consume" atomique.
//
//	Pour marquer manuellement un xuple comme consommé, utilisez MarkConsumed().
func (xs *DefaultXupleSpace) Retrieve(agentID string) (*Xuple, error) {
	if agentID == "" {
		return nil, ErrEmptyAgentID
	}

	xs.mu.Lock()
	defer xs.mu.Unlock()

	// Collecter les xuples disponibles pour cet agent
	available := make([]*Xuple, 0)
	for _, xuple := range xs.xuples {
		// Marquer comme expiré si nécessaire (avec lock)
		if xuple.IsExpired() && xuple.Metadata.State != XupleStateExpired {
			xuple.Metadata.State = XupleStateExpired
		}

		if xuple.CanBeConsumedBy(agentID, xs.config.ConsumptionPolicy) {
			available = append(available, xuple)
		}
	}

	if len(available) == 0 {
		return nil, ErrNoAvailableXuple
	}

	// Sélectionner selon la politique
	selected := xs.config.SelectionPolicy.Select(available)
	if selected == nil {
		return nil, ErrNoAvailableXuple
	}

	// CORRECTION DU BUG CRITIQUE : Marquer automatiquement comme consommé
	// Cela évite que l'appelant oublie d'appeler MarkConsumed() et garantit
	// que la politique de consommation 'once' fonctionne correctement
	selected.markConsumedBy(agentID)

	// Vérifier si le xuple doit être marqué comme complètement consommé
	if xs.config.ConsumptionPolicy.OnConsumed(selected, agentID) {
		selected.Metadata.State = XupleStateConsumed
	}

	return selected, nil
}

// RetrieveMultiple récupère jusqu'à n xuples pour un agent selon les politiques.
//
// Cette méthode implémente une récupération batch efficace en une seule opération atomique.
//
// Paramètres:
//   - agentID: Identifiant de l'agent récupérant les xuples
//   - n: Nombre maximum de xuples à récupérer (doit être > 0)
//
// Retour:
//   - []*Xuple: Slice de xuples récupérés (vide si aucun disponible, jamais nil)
//   - error: Erreur si agentID vide ou n <= 0
//
// Comportement:
//   - Si n <= 0, retourne ErrInvalidConfiguration
//   - Si agentID vide, retourne ErrEmptyAgentID
//   - Récupère jusqu'à n xuples disponibles selon SelectionPolicy
//   - Si moins de n xuples disponibles, retourne tous les disponibles (pas d'erreur)
//   - Marque automatiquement tous les xuples récupérés comme consommés
//   - Met à jour l'état selon ConsumptionPolicy
//
// Thread-Safety:
//   - Opération atomique sous mutex (tous les xuples récupérés en une transaction)
//
// Exemple:
//
//	xuples, err := space.RetrieveMultiple("agent1", 5)
//	// Récupère jusqu'à 5 xuples pour agent1
func (xs *DefaultXupleSpace) RetrieveMultiple(agentID string, n int) ([]*Xuple, error) {
	if agentID == "" {
		return nil, ErrEmptyAgentID
	}

	if n < 0 {
		return nil, ErrInvalidConfiguration
	}

	// Si n == 0, retourner slice vide (comportement valide)
	if n == 0 {
		return make([]*Xuple, 0), nil
	}

	xs.mu.Lock()
	defer xs.mu.Unlock()

	// Collecter les xuples disponibles pour cet agent
	available := make([]*Xuple, 0)
	for _, xuple := range xs.xuples {
		// Marquer comme expiré si nécessaire
		if xuple.IsExpired() && xuple.Metadata.State != XupleStateExpired {
			xuple.Metadata.State = XupleStateExpired
		}

		if xuple.CanBeConsumedBy(agentID, xs.config.ConsumptionPolicy) {
			available = append(available, xuple)
		}
	}

	if len(available) == 0 {
		// Retourner slice vide (non nil) si aucun xuple disponible
		return make([]*Xuple, 0), nil
	}

	// Limiter au nombre demandé
	count := n
	if count > len(available) {
		count = len(available)
	}

	// Sélectionner les xuples selon la politique
	// Note: On doit appeler Select() count fois pour respecter la politique
	// (ex: FIFO prend le premier, puis le suivant, etc.)
	selected := make([]*Xuple, 0, count)
	remaining := available

	for i := 0; i < count; i++ {
		if len(remaining) == 0 {
			break
		}

		xuple := xs.config.SelectionPolicy.Select(remaining)
		if xuple == nil {
			break
		}

		// Marquer comme consommé
		xuple.markConsumedBy(agentID)

		// Vérifier si le xuple doit être marqué comme complètement consommé
		if xs.config.ConsumptionPolicy.OnConsumed(xuple, agentID) {
			xuple.Metadata.State = XupleStateConsumed
		}

		selected = append(selected, xuple)

		// Retirer le xuple sélectionné de la liste des disponibles
		newRemaining := make([]*Xuple, 0, len(remaining)-1)
		for _, x := range remaining {
			if x.ID != xuple.ID {
				newRemaining = append(newRemaining, x)
			}
		}
		remaining = newRemaining
	}

	return selected, nil
}

// MarkConsumed marque un xuple comme consommé par un agent.
func (xs *DefaultXupleSpace) MarkConsumed(xupleID string, agentID string) error {
	if agentID == "" {
		return ErrEmptyAgentID
	}

	xs.mu.Lock()
	defer xs.mu.Unlock()

	xuple, exists := xs.xuples[xupleID]
	if !exists {
		return ErrXupleNotFound
	}

	if !xuple.CanBeConsumedBy(agentID, xs.config.ConsumptionPolicy) {
		return ErrXupleNotAvailable
	}

	// Marquer comme consommé (thread-safe car nous avons le lock)
	xuple.markConsumedBy(agentID)

	// Vérifier si le xuple doit être marqué comme complètement consommé
	if xs.config.ConsumptionPolicy.OnConsumed(xuple, agentID) {
		xuple.Metadata.State = XupleStateConsumed
	}

	return nil
}

// Count retourne le nombre de xuples disponibles.
func (xs *DefaultXupleSpace) Count() int {
	xs.mu.RLock()
	defer xs.mu.RUnlock()

	count := 0
	for _, xuple := range xs.xuples {
		if xuple.IsAvailable() && !xuple.IsExpired() {
			count++
		}
	}

	return count
}

// Cleanup nettoie les xuples expirés.
func (xs *DefaultXupleSpace) Cleanup() int {
	xs.mu.Lock()
	defer xs.mu.Unlock()

	cleaned := 0
	for id, xuple := range xs.xuples {
		if !xs.config.RetentionPolicy.ShouldRetain(xuple) || xuple.IsExpired() {
			delete(xs.xuples, id)
			cleaned++
		}
	}

	return cleaned
}

// GetConfig retourne la configuration du xuple-space.
func (xs *DefaultXupleSpace) GetConfig() XupleSpaceConfig {
	return xs.config
}

// ListAll retourne tous les xuples du xuple-space.
//
// Cette méthode est principalement destinée à l'inspection et aux tests.
// Elle ne respecte PAS les politiques de sélection ou de consommation.
//
// Retourne:
//   - []*Xuple: liste de tous les xuples (disponibles, consommés, expirés)
//
// Thread-Safety:
//   - Méthode thread-safe (protégée par mutex en lecture)
func (xs *DefaultXupleSpace) ListAll() []*Xuple {
	xs.mu.RLock()
	defer xs.mu.RUnlock()

	xuples := make([]*Xuple, 0, len(xs.xuples))
	for _, xuple := range xs.xuples {
		xuples = append(xuples, xuple)
	}

	return xuples
}
