// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package xuples

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/treivax/tsd/rete"
)

// XupleState représente l'état d'un xuple dans son cycle de vie.
type XupleState int

const (
	// XupleStateAvailable indique qu'un xuple est disponible pour consommation
	XupleStateAvailable XupleState = iota

	// XupleStateConsumed indique qu'un xuple a été complètement consommé
	XupleStateConsumed

	// XupleStateExpired indique qu'un xuple a dépassé sa durée de vie
	XupleStateExpired
)

// String retourne la représentation textuelle de l'état.
func (s XupleState) String() string {
	switch s {
	case XupleStateAvailable:
		return "available"
	case XupleStateConsumed:
		return "consumed"
	case XupleStateExpired:
		return "expired"
	default:
		return "unknown"
	}
}

// XupleMetadata contient les métadonnées d'un xuple.
type XupleMetadata struct {
	// ConsumptionCount nombre total de consommations
	ConsumptionCount int

	// ConsumedBy agents ayant consommé (agent ID -> timestamp)
	ConsumedBy map[string]time.Time

	// ExpiresAt date d'expiration (zero time si illimitée)
	ExpiresAt time.Time

	// State état actuel du xuple
	State XupleState
}

// Xuple représente un tuple étendu avec fait principal et faits déclencheurs.
//
// Un xuple encapsule :
//   - Le fait principal résultant d'une activation
//   - Les faits déclencheurs qui ont conduit à cette activation
//   - Les métadonnées de tracking (consommation, expiration, état)
//
// Thread-Safety :
//   - Lecture des champs (ID, Fact, TriggeringFacts, CreatedAt) est thread-safe
//   - Modification des Metadata se fait UNIQUEMENT via XupleSpace avec lock approprié
//   - Ne jamais modifier directement les champs après création
type Xuple struct {
	// ID unique du xuple (format: "xuple_<counter>")
	ID string

	// Fait principal résultant de l'activation
	Fact *rete.Fact

	// Faits qui ont déclenché cette activation
	TriggeringFacts []*rete.Fact

	// Timestamp de création
	CreatedAt time.Time

	// Métadonnées (consommation, expiration, état)
	Metadata XupleMetadata
}

// IsAvailable retourne true si le xuple est disponible pour consommation.
func (x *Xuple) IsAvailable() bool {
	return x.Metadata.State == XupleStateAvailable
}

// IsExpired vérifie si le xuple a expiré.
// Note: La modification de l'état doit être faite par XupleSpace avec un lock.
// Cette méthode est read-only pour éviter les race conditions.
func (x *Xuple) IsExpired() bool {
	if x.Metadata.State == XupleStateExpired {
		return true
	}

	if !x.Metadata.ExpiresAt.IsZero() && time.Now().After(x.Metadata.ExpiresAt) {
		return true
	}

	return false
}

// CanBeConsumedBy vérifie si un agent peut consommer ce xuple.
func (x *Xuple) CanBeConsumedBy(agentID string, policy ConsumptionPolicy) bool {
	if !x.IsAvailable() || x.IsExpired() {
		return false
	}

	return policy.CanConsume(x, agentID)
}

// markConsumedBy marque le xuple comme consommé par un agent.
// Cette méthode est appelée uniquement depuis XupleSpace avec un lock approprié.
// Ne pas appeler directement - non thread-safe.
func (x *Xuple) markConsumedBy(agentID string) {
	if x.Metadata.ConsumedBy == nil {
		x.Metadata.ConsumedBy = make(map[string]time.Time)
	}

	x.Metadata.ConsumedBy[agentID] = time.Now()
	x.Metadata.ConsumptionCount++
}

// XupleSpaceConfig configure un xuple-space.
type XupleSpaceConfig struct {
	// Name nom du xuple-space
	Name string

	// SelectionPolicy politique de sélection
	SelectionPolicy SelectionPolicy

	// ConsumptionPolicy politique de consommation
	ConsumptionPolicy ConsumptionPolicy

	// RetentionPolicy politique de rétention
	RetentionPolicy RetentionPolicy

	// MaxSize taille maximale du xuple-space (0 = illimité)
	MaxSize int
}

// XupleSpace représente un espace de xuples.
type XupleSpace interface {
	// Name retourne le nom du xuple-space
	Name() string

	// Insert insère un xuple dans le xuple-space
	Insert(xuple *Xuple) error

	// Retrieve récupère un xuple pour un agent selon les politiques
	Retrieve(agentID string) (*Xuple, error)

	// RetrieveMultiple récupère jusqu'à n xuples pour un agent selon les politiques
	// Si n <= 0, retourne une erreur
	// Si moins de n xuples disponibles, retourne tous les xuples disponibles sans erreur
	// Retourne une slice vide (non nil) si aucun xuple disponible
	RetrieveMultiple(agentID string, n int) ([]*Xuple, error)

	// MarkConsumed marque un xuple comme consommé par un agent
	MarkConsumed(xupleID string, agentID string) error

	// Count retourne le nombre de xuples disponibles
	Count() int

	// Cleanup nettoie les xuples expirés
	Cleanup() int

	// GetConfig retourne la configuration du xuple-space
	GetConfig() XupleSpaceConfig

	// ListAll retourne tous les xuples du xuple-space (pour inspection/tests)
	// Cette méthode n'applique pas les politiques de sélection/consommation
	ListAll() []*Xuple
}

// XupleManager gère les xuple-spaces.
type XupleManager interface {
	// CreateXupleSpace crée un nouveau xuple-space avec les politiques données
	CreateXupleSpace(name string, config XupleSpaceConfig) error

	// GetXupleSpace retourne un xuple-space par son nom
	GetXupleSpace(name string) (XupleSpace, error)

	// CreateXuple crée un xuple dans le xuple-space spécifié
	CreateXuple(xuplespace string, fact *rete.Fact, triggeringFacts []*rete.Fact) error

	// ListXupleSpaces retourne la liste des noms de xuple-spaces
	ListXupleSpaces() []string

	// Close ferme le manager et nettoie les ressources
	Close() error
}

// DefaultXupleManager implémente XupleManager.
type DefaultXupleManager struct {
	spaces map[string]XupleSpace
	mu     sync.RWMutex
}

// NewXupleManager crée un nouveau gestionnaire de xuple-spaces.
func NewXupleManager() XupleManager {
	return &DefaultXupleManager{
		spaces: make(map[string]XupleSpace),
	}
}

// generateXupleID génère un ID unique pour un xuple.
// Thread-safe via UUID v4.
func (m *DefaultXupleManager) generateXupleID() string {
	return uuid.New().String()
}

// CreateXupleSpace crée un nouveau xuple-space avec les politiques données.
func (m *DefaultXupleManager) CreateXupleSpace(name string, config XupleSpaceConfig) error {
	if name == "" {
		return ErrInvalidConfiguration
	}

	// Valider les politiques
	if config.SelectionPolicy == nil || config.ConsumptionPolicy == nil || config.RetentionPolicy == nil {
		return ErrInvalidPolicy
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.spaces[name]; exists {
		return ErrXupleSpaceExists
	}

	config.Name = name
	space := NewXupleSpace(config)
	m.spaces[name] = space

	return nil
}

// GetXupleSpace retourne un xuple-space par son nom.
func (m *DefaultXupleManager) GetXupleSpace(name string) (XupleSpace, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	space, exists := m.spaces[name]
	if !exists {
		return nil, ErrXupleSpaceNotFound
	}

	return space, nil
}

// CreateXuple crée un xuple dans le xuple-space spécifié.
func (m *DefaultXupleManager) CreateXuple(xuplespace string, fact *rete.Fact, triggeringFacts []*rete.Fact) error {
	// Validation des paramètres
	if xuplespace == "" {
		return ErrXupleSpaceNotFound
	}

	if fact == nil {
		return ErrNilFact
	}

	space, err := m.GetXupleSpace(xuplespace)
	if err != nil {
		return err
	}

	xuple := &Xuple{
		ID:              m.generateXupleID(),
		Fact:            fact,
		TriggeringFacts: triggeringFacts,
		CreatedAt:       time.Now(),
		Metadata: XupleMetadata{
			State:      XupleStateAvailable,
			ConsumedBy: make(map[string]time.Time),
		},
	}

	return space.Insert(xuple)
}

// ListXupleSpaces retourne la liste des noms de xuple-spaces.
func (m *DefaultXupleManager) ListXupleSpaces() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	names := make([]string, 0, len(m.spaces))
	for name := range m.spaces {
		names = append(names, name)
	}

	return names
}

// Close ferme le manager et nettoie les ressources.
func (m *DefaultXupleManager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Nettoyer tous les xuple-spaces
	for _, space := range m.spaces {
		space.Cleanup()
	}

	m.spaces = make(map[string]XupleSpace)
	return nil
}
