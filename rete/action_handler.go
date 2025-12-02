// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"sync"
)

// ActionHandler définit l'interface pour les gestionnaires d'actions personnalisées.
// Chaque action peut avoir son propre handler qui définit son comportement.
type ActionHandler interface {
	// Execute exécute l'action avec les arguments évalués fournis.
	// Retourne une erreur si l'exécution échoue.
	Execute(args []interface{}, ctx *ExecutionContext) error

	// GetName retourne le nom de l'action gérée par ce handler.
	GetName() string

	// Validate valide que les arguments sont corrects pour cette action.
	// Cette validation est optionnelle et peut retourner nil si aucune validation spécifique n'est nécessaire.
	Validate(args []interface{}) error
}

// ActionRegistry gère l'enregistrement et la récupération des handlers d'actions.
type ActionRegistry struct {
	handlers map[string]ActionHandler
	mu       sync.RWMutex
}

// NewActionRegistry crée un nouveau registry d'actions.
func NewActionRegistry() *ActionRegistry {
	return &ActionRegistry{
		handlers: make(map[string]ActionHandler),
	}
}

// Register enregistre un handler d'action.
// Si un handler existe déjà pour ce nom, il est remplacé.
func (ar *ActionRegistry) Register(handler ActionHandler) error {
	if handler == nil {
		return fmt.Errorf("handler cannot be nil")
	}

	name := handler.GetName()
	if name == "" {
		return fmt.Errorf("handler name cannot be empty")
	}

	ar.mu.Lock()
	defer ar.mu.Unlock()

	ar.handlers[name] = handler
	return nil
}

// Unregister supprime un handler d'action du registry.
func (ar *ActionRegistry) Unregister(actionName string) {
	ar.mu.Lock()
	defer ar.mu.Unlock()

	delete(ar.handlers, actionName)
}

// Get récupère un handler d'action par son nom.
// Retourne nil si aucun handler n'est enregistré pour ce nom.
func (ar *ActionRegistry) Get(actionName string) ActionHandler {
	ar.mu.RLock()
	defer ar.mu.RUnlock()

	return ar.handlers[actionName]
}

// Has vérifie si un handler est enregistré pour une action donnée.
func (ar *ActionRegistry) Has(actionName string) bool {
	ar.mu.RLock()
	defer ar.mu.RUnlock()

	_, exists := ar.handlers[actionName]
	return exists
}

// GetAll retourne une copie de tous les handlers enregistrés.
func (ar *ActionRegistry) GetAll() map[string]ActionHandler {
	ar.mu.RLock()
	defer ar.mu.RUnlock()

	copy := make(map[string]ActionHandler, len(ar.handlers))
	for k, v := range ar.handlers {
		copy[k] = v
	}
	return copy
}

// Count retourne le nombre de handlers enregistrés.
func (ar *ActionRegistry) Count() int {
	ar.mu.RLock()
	defer ar.mu.RUnlock()

	return len(ar.handlers)
}

// Clear supprime tous les handlers du registry.
func (ar *ActionRegistry) Clear() {
	ar.mu.Lock()
	defer ar.mu.Unlock()

	ar.handlers = make(map[string]ActionHandler)
}

// RegisterMultiple enregistre plusieurs handlers en une seule opération.
// Si un handler échoue, les précédents restent enregistrés.
func (ar *ActionRegistry) RegisterMultiple(handlers []ActionHandler) error {
	for i, handler := range handlers {
		if err := ar.Register(handler); err != nil {
			return fmt.Errorf("failed to register handler at index %d: %w", i, err)
		}
	}
	return nil
}

// GetRegisteredNames retourne la liste des noms d'actions enregistrées.
func (ar *ActionRegistry) GetRegisteredNames() []string {
	ar.mu.RLock()
	defer ar.mu.RUnlock()

	names := make([]string, 0, len(ar.handlers))
	for name := range ar.handlers {
		names = append(names, name)
	}
	return names
}
