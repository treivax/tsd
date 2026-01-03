// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"fmt"
	"time"
)

// PropagationMode définit le mode de propagation à utiliser.
type PropagationMode int

const (
	// PropagationModeDelta utilise la propagation sélective par delta
	PropagationModeDelta PropagationMode = iota

	// PropagationModeClassic utilise Retract+Insert classique (fallback)
	PropagationModeClassic

	// PropagationModeAuto choisit automatiquement selon le contexte
	PropagationModeAuto
)

// String retourne la représentation string du PropagationMode
func (pm PropagationMode) String() string {
	switch pm {
	case PropagationModeDelta:
		return "Delta"
	case PropagationModeClassic:
		return "Classic"
	case PropagationModeAuto:
		return "Auto"
	default:
		return "Unknown"
	}
}

// PropagationConfig contient la configuration du DeltaPropagator.
//
// Cette configuration permet de contrôler le comportement de la propagation
// sélective et les critères de fallback vers la propagation classique.
type PropagationConfig struct {
	// Mode de propagation par défaut
	DefaultMode PropagationMode

	// EnableDeltaPropagation active/désactive la propagation delta
	// (master switch pour activation/désactivation globale)
	EnableDeltaPropagation bool

	// DeltaThreshold est le seuil de ratio de changement au-delà duquel
	// on bascule en mode classique (Retract+Insert).
	// Valeur entre 0.0 et 1.0.
	// Exemple : 0.3 → si > 30% des champs changent, utiliser mode classique
	// Default: 0.5
	DeltaThreshold float64

	// MinFieldsForDelta est le nombre minimum de champs dans un fait
	// pour que la propagation delta soit utilisée.
	// Si le fait a moins de champs, utiliser mode classique.
	// Default: 3
	MinFieldsForDelta int

	// MaxAffectedNodesForDelta est le nombre maximum de nœuds affectés
	// au-delà duquel on bascule en mode classique.
	// Rationale : si trop de nœuds affectés, overhead delta > bénéfice
	// Default: 100
	MaxAffectedNodesForDelta int

	// AllowPrimaryKeyChange indique si les modifications de clé primaire
	// sont autorisées en mode delta.
	// Si false, tout changement de PK force le mode classique.
	// Default: false (car changement PK = changement d'ID interne)
	AllowPrimaryKeyChange bool

	// PrimaryKeyFields liste les noms de champs considérés comme clés primaires.
	// Si vide, détection automatique depuis les TypeDefinitions.
	// Default: []
	PrimaryKeyFields []string

	// EnableMetrics active la collecte de métriques de propagation
	// Default: true
	EnableMetrics bool

	// PropagationTimeout est le timeout maximum pour une propagation delta.
	// Si dépassé, la propagation est annulée (protection deadlock).
	// Default: 30 secondes
	PropagationTimeout time.Duration

	// RetryOnError indique si une propagation delta échouée doit être
	// retentée en mode classique (fallback automatique).
	// Default: true
	RetryOnError bool

	// MaxConcurrentPropagations est le nombre maximum de propagations
	// delta simultanées autorisées (contrôle charge).
	// Default: 10
	MaxConcurrentPropagations int

	// EnableOptimisticPropagation active la propagation optimiste :
	// ne pas attendre la fin de la propagation pour retourner.
	// Default: false (attente synchrone)
	EnableOptimisticPropagation bool

	// LogPropagationDetails active le logging détaillé de chaque propagation
	// (utile pour debugging, overhead en production).
	// Default: false
	LogPropagationDetails bool
}

// DefaultPropagationConfig retourne une configuration par défaut.
func DefaultPropagationConfig() PropagationConfig {
	return PropagationConfig{
		DefaultMode:                 PropagationModeAuto,
		EnableDeltaPropagation:      true,
		DeltaThreshold:              0.5,
		MinFieldsForDelta:           3,
		MaxAffectedNodesForDelta:    100,
		AllowPrimaryKeyChange:       false,
		PrimaryKeyFields:            []string{},
		EnableMetrics:               true,
		PropagationTimeout:          30 * time.Second,
		RetryOnError:                true,
		MaxConcurrentPropagations:   10,
		EnableOptimisticPropagation: false,
		LogPropagationDetails:       false,
	}
}

// Validate vérifie que la configuration est valide.
//
// Retourne une erreur si des paramètres sont incohérents.
func (pc *PropagationConfig) Validate() error {
	if pc.DeltaThreshold < 0.0 || pc.DeltaThreshold > 1.0 {
		return fmt.Errorf("DeltaThreshold must be between 0.0 and 1.0, got %v", pc.DeltaThreshold)
	}

	if pc.MinFieldsForDelta < 0 {
		return fmt.Errorf("MinFieldsForDelta must be >= 0, got %d", pc.MinFieldsForDelta)
	}

	if pc.MaxAffectedNodesForDelta < 1 {
		return fmt.Errorf("MaxAffectedNodesForDelta must be >= 1, got %d", pc.MaxAffectedNodesForDelta)
	}

	if pc.PropagationTimeout < 0 {
		return fmt.Errorf("PropagationTimeout must be >= 0, got %v", pc.PropagationTimeout)
	}

	if pc.MaxConcurrentPropagations < 1 {
		return fmt.Errorf("MaxConcurrentPropagations must be >= 1, got %d", pc.MaxConcurrentPropagations)
	}

	return nil
}

// ShouldUseDelta détermine si la propagation delta doit être utilisée
// pour un FactDelta donné.
//
// Cette méthode applique les heuristiques configurées pour décider
// du mode de propagation optimal.
//
// Paramètres :
//   - delta : le FactDelta à propager
//   - affectedNodesCount : nombre de nœuds qui seraient affectés
//
// Retourne true si la propagation delta doit être utilisée.
func (pc *PropagationConfig) ShouldUseDelta(delta *FactDelta, affectedNodesCount int) bool {
	// 1. Feature flag global
	if !pc.EnableDeltaPropagation {
		return false
	}

	// 2. Vérifier mode forcé
	if pc.DefaultMode == PropagationModeClassic {
		return false
	}
	if pc.DefaultMode == PropagationModeDelta {
		return true
	}

	// 3. Mode Auto : appliquer heuristiques

	// Heuristique 1 : Nombre de champs
	if delta.FieldCount < pc.MinFieldsForDelta {
		return false
	}

	// Heuristique 2 : Ratio de changement
	if delta.ChangeRatio() > pc.DeltaThreshold {
		return false
	}

	// Heuristique 3 : Nombre de nœuds affectés
	if affectedNodesCount > pc.MaxAffectedNodesForDelta {
		return false
	}

	// Heuristique 4 : Changement de clé primaire
	if !pc.AllowPrimaryKeyChange && pc.hasPrimaryKeyChange(delta) {
		return false
	}

	// Toutes les conditions passées : utiliser delta
	return true
}

// hasPrimaryKeyChange vérifie si le delta contient un changement de clé primaire.
func (pc *PropagationConfig) hasPrimaryKeyChange(delta *FactDelta) bool {
	if len(pc.PrimaryKeyFields) == 0 {
		// Pas de clés primaires configurées : autoriser
		return false
	}

	for _, pkField := range pc.PrimaryKeyFields {
		if _, changed := delta.Fields[pkField]; changed {
			return true
		}
	}

	return false
}

// Clone crée une copie de la configuration.
func (pc *PropagationConfig) Clone() PropagationConfig {
	pkFields := make([]string, len(pc.PrimaryKeyFields))
	copy(pkFields, pc.PrimaryKeyFields)

	return PropagationConfig{
		DefaultMode:                 pc.DefaultMode,
		EnableDeltaPropagation:      pc.EnableDeltaPropagation,
		DeltaThreshold:              pc.DeltaThreshold,
		MinFieldsForDelta:           pc.MinFieldsForDelta,
		MaxAffectedNodesForDelta:    pc.MaxAffectedNodesForDelta,
		AllowPrimaryKeyChange:       pc.AllowPrimaryKeyChange,
		PrimaryKeyFields:            pkFields,
		EnableMetrics:               pc.EnableMetrics,
		PropagationTimeout:          pc.PropagationTimeout,
		RetryOnError:                pc.RetryOnError,
		MaxConcurrentPropagations:   pc.MaxConcurrentPropagations,
		EnableOptimisticPropagation: pc.EnableOptimisticPropagation,
		LogPropagationDetails:       pc.LogPropagationDetails,
	}
}
