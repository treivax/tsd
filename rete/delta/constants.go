// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import "time"

// Node type constants - Types de nœuds RETE
const (
	NodeTypeAlpha    = "alpha"
	NodeTypeBeta     = "beta"
	NodeTypeTerminal = "terminal"
)

// AST node type constants - Types de nœuds AST du parser
const (
	ASTNodeTypeFieldAccess     = "fieldAccess"
	ASTNodeTypeBinaryOp        = "binaryOp"
	ASTNodeTypeComparison      = "comparison"
	ASTNodeTypeUpdateWithModif = "updateWithModifications"
	ASTNodeTypeFactCreation    = "factCreation"
	ASTFieldNameType           = "type"
	ASTFieldNameField          = "field"
	ASTFieldNameLeft           = "left"
	ASTFieldNameRight          = "right"
	ASTFieldNameModifications  = "modifications"
	ASTFieldNameFields         = "fields"
)

// Configuration defaults
const (
	// DefaultFloatEpsilon est la tolérance par défaut pour comparaison de floats
	DefaultFloatEpsilon = 1e-9

	// DefaultMaxDeltaAge est l'âge maximum d'un delta avant expiration (pour cache)
	DefaultMaxDeltaAge = 5 * time.Minute

	// DefaultMaxNestingLevel est la profondeur maximale pour comparaison récursive
	DefaultMaxNestingLevel = 10

	// DefaultCacheTTL est la durée de vie par défaut des entrées de cache
	DefaultCacheTTL = 1 * time.Minute

	// DefaultDeltaThreshold est le seuil de ratio pour fallback classique
	// Si >30% des champs changent, on utilise Retract+Insert au lieu de delta
	DefaultDeltaThreshold = 0.3

	// DefaultMaxConcurrentPropagations limite le nombre de propagations parallèles
	DefaultMaxConcurrentPropagations = 10
)

// Error messages - Messages d'erreur constants
const (
	ErrMsgPropagatorNotInit    = "propagator not initialized"
	ErrMsgCallbacksNotConfig   = "callbacks not configured"
	ErrMsgIndexNotInit         = "index not initialized"
	ErrMsgDeltaPropagationFail = "delta propagation failed"
	ErrMsgStorageUpdateFail    = "storage update failed"
	ErrMsgInvalidFactID        = "invalid fact ID"
	ErrMsgInvalidFactType      = "invalid fact type"
)

// String key prefixes - Préfixes pour clés de groupage
const (
	KeySeparator      = ":"
	MinAlphaPrefixLen = 5 // len("alpha")
	MinBetaPrefixLen  = 4 // len("beta")
	MinTermPrefixLen  = 8 // len("terminal")
)
