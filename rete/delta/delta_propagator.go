// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// PropagateCallback est le callback pour propager un delta vers un nœud RETE.
type PropagateCallback func(nodeID string, delta *FactDelta) error

// ClassicPropagationCallback est le callback pour la propagation classique (Retract+Insert).
//
// Cette fonction est appelée lorsque la propagation delta n'est pas applicable
// (par exemple : PK modifiée, trop de champs changés, feature désactivée).
//
// Paramètres :
//   - factID : Identifiant interne du fait (format "Type~values")
//   - factType : Type du fait (ex: "Product", "Order")
//   - oldFact : Ancien état du fait (peut être nil si Insert pur)
//   - newFact : Nouvel état du fait (peut être nil si Retract pur)
//
// L'implémentation doit gérer :
//   - Retract de oldFact si non-nil
//   - Insert de newFact si non-nil
//   - Les deux opérations de manière atomique pour un Update
type ClassicPropagationCallback func(
	factID string,
	factType string,
	oldFact map[string]interface{},
	newFact map[string]interface{},
) error

// DeltaPropagator orchestre la propagation sélective des changements.
//
// Il coordonne la détection de delta, la recherche de nœuds affectés,
// et la propagation vers ces nœuds selon la stratégie configurée.
//
// Thread-safety : DeltaPropagator est safe pour utilisation concurrente.
type DeltaPropagator struct {
	detector          *DeltaDetector
	index             *DependencyIndex
	strategy          PropagationStrategy
	config            PropagationConfig
	metrics           *PropagationMetrics
	semaphore         chan struct{}
	mutex             sync.RWMutex
	onPropagate       PropagateCallback
	onClassicFallback ClassicPropagationCallback
}

// DeltaPropagatorBuilder construit un DeltaPropagator avec pattern builder.
type DeltaPropagatorBuilder struct {
	detector          *DeltaDetector
	index             *DependencyIndex
	strategy          PropagationStrategy
	config            PropagationConfig
	onPropagate       PropagateCallback
	onClassicFallback ClassicPropagationCallback
}

// NewDeltaPropagatorBuilder crée un nouveau builder.
func NewDeltaPropagatorBuilder() *DeltaPropagatorBuilder {
	return &DeltaPropagatorBuilder{
		config: DefaultPropagationConfig(),
	}
}

// WithDetector configure le détecteur de delta.
func (b *DeltaPropagatorBuilder) WithDetector(detector *DeltaDetector) *DeltaPropagatorBuilder {
	b.detector = detector
	return b
}

// WithIndex configure l'index de dépendances.
func (b *DeltaPropagatorBuilder) WithIndex(index *DependencyIndex) *DeltaPropagatorBuilder {
	b.index = index
	return b
}

// WithStrategy configure la stratégie de propagation.
func (b *DeltaPropagatorBuilder) WithStrategy(strategy PropagationStrategy) *DeltaPropagatorBuilder {
	b.strategy = strategy
	return b
}

// WithConfig configure la propagation.
func (b *DeltaPropagatorBuilder) WithConfig(config PropagationConfig) *DeltaPropagatorBuilder {
	b.config = config
	return b
}

// WithPropagateCallback configure le callback de propagation vers le réseau RETE.
func (b *DeltaPropagatorBuilder) WithPropagateCallback(
	callback PropagateCallback,
) *DeltaPropagatorBuilder {
	b.onPropagate = callback
	return b
}

// WithClassicPropagationCallback configure le callback pour la propagation classique (Retract+Insert).
//
// Ce callback est utilisé lorsque la propagation delta n'est pas applicable :
//   - Feature delta désactivée (EnableDeltaPropagation = false)
//   - Clé primaire modifiée
//   - Ratio de changement trop élevé (> DeltaThreshold)
//   - Erreur lors de la propagation delta et RetryOnError activé
//
// Si ce callback n'est pas configuré, les cas de fallback retourneront une erreur.
func (b *DeltaPropagatorBuilder) WithClassicPropagationCallback(
	callback ClassicPropagationCallback,
) *DeltaPropagatorBuilder {
	b.onClassicFallback = callback
	return b
}

// Build construit le DeltaPropagator.
func (b *DeltaPropagatorBuilder) Build() (*DeltaPropagator, error) {
	if err := b.config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	if b.index == nil {
		return nil, fmt.Errorf("dependency index is required")
	}

	if b.detector == nil {
		b.detector = NewDeltaDetector()
	}

	if b.strategy == nil {
		b.strategy = &SequentialStrategy{}
	}

	semaphore := make(chan struct{}, b.config.MaxConcurrentPropagations)

	return &DeltaPropagator{
		detector:          b.detector,
		index:             b.index,
		strategy:          b.strategy,
		config:            b.config,
		metrics:           &PropagationMetrics{},
		semaphore:         semaphore,
		onPropagate:       b.onPropagate,
		onClassicFallback: b.onClassicFallback,
	}, nil
}

// PropagateUpdate propage une mise à jour de fait.
//
// Cette méthode est le point d'entrée principal pour propager un Update.
//
// Paramètres :
//   - oldFact : fait avant modification
//   - newFact : fait après modification
//   - factID : identifiant interne du fait
//   - factType : type du fait
//
// Retourne une erreur si la propagation échoue.
func (dp *DeltaPropagator) PropagateUpdate(
	oldFact, newFact map[string]interface{},
	factID, factType string,
) error {
	return dp.PropagateUpdateWithContext(
		context.Background(),
		oldFact, newFact,
		factID, factType,
	)
}

// PropagateUpdateWithContext propage avec un contexte (timeout, cancellation).
func (dp *DeltaPropagator) PropagateUpdateWithContext(
	ctx context.Context,
	oldFact, newFact map[string]interface{},
	factID, factType string,
) error {
	if !dp.config.EnableDeltaPropagation {
		return dp.classicPropagation(ctx, oldFact, newFact, factID, factType)
	}

	start := time.Now()

	delta, err := dp.detector.DetectDelta(oldFact, newFact, factID, factType)
	if err != nil {
		dp.metrics.RecordFailedPropagation()
		return fmt.Errorf("failed to detect delta: %w", err)
	}

	affectedNodes := dp.index.GetAffectedNodesForDelta(delta)

	shouldUseDelta := dp.config.ShouldUseDelta(delta, len(affectedNodes))

	if !shouldUseDelta {
		dp.recordFallbackReason(delta, len(affectedNodes))
		return dp.classicPropagation(ctx, oldFact, newFact, factID, factType)
	}

	if !dp.strategy.ShouldPropagate(delta, affectedNodes) {
		return nil
	}

	if err := dp.executeDeltaPropagation(ctx, delta, affectedNodes); err != nil {
		if dp.config.RetryOnError {
			dp.metrics.RecordFallback("error")
			return dp.classicPropagation(ctx, oldFact, newFact, factID, factType)
		}
		return err
	}

	if dp.config.EnableMetrics {
		duration := time.Since(start)
		dp.metrics.RecordDeltaPropagation(duration, len(affectedNodes), len(delta.Fields))
	}

	return nil
}

// executeDeltaPropagation exécute la propagation delta vers les nœuds affectés.
//
// Cette méthode utilise le batch processing pour grouper les nœuds par type
// et optimiser la propagation séquentielle.
func (dp *DeltaPropagator) executeDeltaPropagation(
	ctx context.Context,
	delta *FactDelta,
	affectedNodes []NodeReference,
) error {
	// Utiliser batch processing pour grouper par type de nœud
	batch := NewBatchNodeReferences(len(affectedNodes))
	for i := range affectedNodes {
		batch.Add(affectedNodes[i])
	}

	// Propager dans l'ordre : Alpha -> Beta -> Terminal
	err := batch.ProcessInOrder(func(nodeRef NodeReference) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case dp.semaphore <- struct{}{}:
			defer func() { <-dp.semaphore }()

			if err := dp.propagateToNode(nodeRef, delta); err != nil {
				return fmt.Errorf("failed to propagate to node %s: %w", nodeRef.NodeID, err)
			}
			return nil
		}
	})

	return err
}

// propagateToNode propage le delta vers un nœud spécifique.
func (dp *DeltaPropagator) propagateToNode(nodeRef NodeReference, delta *FactDelta) error {
	if dp.onPropagate == nil {
		return fmt.Errorf("no propagation callback configured")
	}

	return dp.onPropagate(nodeRef.NodeID, delta)
}

// classicPropagation effectue une propagation classique (Retract+Insert).
//
// Cette méthode délègue à onClassicFallback si configuré, sinon retourne une erreur.
func (dp *DeltaPropagator) classicPropagation(
	ctx context.Context,
	oldFact, newFact map[string]interface{},
	factID, factType string,
) error {
	start := time.Now()

	// Vérifier que le callback est configuré
	if dp.onClassicFallback == nil {
		return fmt.Errorf("classic propagation callback not configured - use WithClassicPropagationCallback()")
	}

	// Appeler le callback pour effectuer Retract+Insert
	err := dp.onClassicFallback(factID, factType, oldFact, newFact)

	// Enregistrer les métriques
	if dp.config.EnableMetrics {
		duration := time.Since(start)
		// Estimer nombre total de nœuds (à obtenir depuis l'index)
		totalNodes := dp.index.GetTotalNodeCount()
		dp.metrics.RecordClassicPropagation(duration, totalNodes)
	}

	return err
}

// recordFallbackReason enregistre la raison du fallback vers mode classique.
func (dp *DeltaPropagator) recordFallbackReason(delta *FactDelta, affectedCount int) {
	if delta.FieldCount < dp.config.MinFieldsForDelta {
		dp.metrics.RecordFallback("fields")
		return
	}

	if delta.ChangeRatio() > dp.config.DeltaThreshold {
		dp.metrics.RecordFallback("ratio")
		return
	}

	if affectedCount > dp.config.MaxAffectedNodesForDelta {
		dp.metrics.RecordFallback("nodes")
		return
	}

	if dp.config.hasPrimaryKeyChange(delta) {
		dp.metrics.RecordFallback("pk")
		return
	}
}

// GetMetrics retourne un snapshot des métriques actuelles.
func (dp *DeltaPropagator) GetMetrics() PropagationMetrics {
	return dp.metrics.GetSnapshot()
}

// ResetMetrics réinitialise les métriques.
func (dp *DeltaPropagator) ResetMetrics() {
	dp.metrics.Reset()
}

// GetConfig retourne la configuration actuelle.
func (dp *DeltaPropagator) GetConfig() PropagationConfig {
	dp.mutex.RLock()
	defer dp.mutex.RUnlock()
	return dp.config.Clone()
}

// UpdateConfig met à jour la configuration.
func (dp *DeltaPropagator) UpdateConfig(config PropagationConfig) error {
	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	dp.mutex.Lock()
	defer dp.mutex.Unlock()
	dp.config = config

	semaphore := make(chan struct{}, config.MaxConcurrentPropagations)
	dp.semaphore = semaphore

	return nil
}
