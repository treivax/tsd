// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

// Package delta implémente le système de propagation incrémentale (RETE-II/TREAT)
// pour optimiser les mises à jour de faits.
//
// Ce package fournit :
//   - Détection des changements de champs (FieldDelta, FactDelta)
//   - Indexation des dépendances (nœuds sensibles à chaque champ)
//   - Propagation sélective (uniquement vers nœuds affectés)
//
// Architecture :
//
//	Update(fact, {field: value})
//	    ↓
//	DetectDelta(oldFact, newFact) → FactDelta
//	    ↓
//	GetAffectedNodes(delta) → [nodes]
//	    ↓
//	PropagateSelective(delta, nodes)
//
// Performance :
//   - Propagation O(nœuds sensibles) au lieu de O(tous nœuds)
//   - Gain typique : 10-100x sur mises à jour partielles
//
// Compatibilité :
//   - Backward compatible (fallback Retract+Insert disponible)
//   - Feature flag pour activation progressive
package delta
