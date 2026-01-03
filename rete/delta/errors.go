// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import "fmt"

// ComponentNotInitializedError représente une erreur de composant non initialisé.
type ComponentNotInitializedError struct {
	Component string // Nom du composant (ex: "propagator", "index")
	Function  string // Nom de la fonction appelante
	Message   string // Message d'erreur détaillé
}

// Error implémente l'interface error
func (e *ComponentNotInitializedError) Error() string {
	return fmt.Sprintf("%s not initialized in %s: %s", e.Component, e.Function, e.Message)
}

// newComponentError crée une nouvelle erreur de composant non initialisé
func newComponentError(component, function, message string) error {
	return &ComponentNotInitializedError{
		Component: component,
		Function:  function,
		Message:   message,
	}
}

// InvalidConfigError représente une erreur de configuration invalide.
type InvalidConfigError struct {
	Field  string // Nom du champ invalide
	Reason string // Raison de l'invalidité
}

// Error implémente l'interface error
func (e *InvalidConfigError) Error() string {
	return fmt.Sprintf("invalid detector config [%s]: %s", e.Field, e.Reason)
}

// newInvalidConfigError crée une nouvelle erreur de configuration invalide
func newInvalidConfigError(field, reason string) error {
	return &InvalidConfigError{
		Field:  field,
		Reason: reason,
	}
}

// InvalidFactError représente une erreur de fait invalide.
type InvalidFactError struct {
	FactID   string // ID du fait invalide
	FactType string // Type du fait
	Reason   string // Raison de l'invalidité
}

// Error implémente l'interface error
func (e *InvalidFactError) Error() string {
	if e.FactID != "" {
		return fmt.Sprintf("invalid fact [%s:%s]: %s", e.FactType, e.FactID, e.Reason)
	}
	return fmt.Sprintf("invalid fact [%s]: %s", e.FactType, e.Reason)
}
