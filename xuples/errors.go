// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package xuples

import "errors"

// Erreurs du module xuples
var (
	// ErrNilXuple est retourné quand un xuple nil est fourni
	ErrNilXuple = errors.New("xuple cannot be nil")

	// ErrXupleNotFound est retourné quand un xuple n'existe pas
	ErrXupleNotFound = errors.New("xuple not found")

	// ErrXupleNotAvailable est retourné quand un xuple n'est pas disponible
	ErrXupleNotAvailable = errors.New("xuple not available for consumption")

	// ErrNoAvailableXuple est retourné quand aucun xuple n'est disponible
	ErrNoAvailableXuple = errors.New("no available xuple")

	// ErrXupleSpaceNotFound est retourné quand un xuple-space n'existe pas
	ErrXupleSpaceNotFound = errors.New("xuple-space not found")

	// ErrXupleSpaceExists est retourné lors d'une tentative de création d'un xuple-space existant
	ErrXupleSpaceExists = errors.New("xuple-space already exists")

	// ErrInvalidPolicy est retourné quand une politique est invalide
	ErrInvalidPolicy = errors.New("invalid policy")

	// ErrInvalidConfiguration est retourné quand une configuration est invalide
	ErrInvalidConfiguration = errors.New("invalid xuple-space configuration")

	// ErrEmptyAgentID est retourné quand l'ID d'agent est vide
	ErrEmptyAgentID = errors.New("agent ID cannot be empty")

	// ErrNilFact est retourné quand un fait nil est fourni
	ErrNilFact = errors.New("fact cannot be nil")

	// ErrXupleSpaceFull est retourné quand le xuple-space a atteint sa capacité maximale
	ErrXupleSpaceFull = errors.New("xuple-space is full")
)
