// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import "fmt"

// Command représente une opération réversible sur le réseau RETE
// Toutes les commandes doivent implémenter Execute() pour appliquer l'opération
// et Undo() pour l'annuler
type Command interface {
	// Execute applique l'opération sur le réseau
	Execute() error

	// Undo annule l'opération et restaure l'état précédent
	Undo() error

	// String retourne une description de la commande pour le debugging
	String() string
}

// CommandError représente une erreur lors de l'exécution ou de l'annulation d'une commande
type CommandError struct {
	CommandName string
	Operation   string // "Execute" ou "Undo"
	Err         error
}

func (e *CommandError) Error() string {
	return fmt.Sprintf("command %s failed during %s: %v", e.CommandName, e.Operation, e.Err)
}

func (e *CommandError) Unwrap() error {
	return e.Err
}

// NewCommandError crée une nouvelle erreur de commande
func NewCommandError(commandName, operation string, err error) *CommandError {
	return &CommandError{
		CommandName: commandName,
		Operation:   operation,
		Err:         err,
	}
}
