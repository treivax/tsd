// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import "fmt"

// AddFactCommand représente l'ajout d'un fait au réseau
type AddFactCommand struct {
	storage Storage
	fact    *Fact
	factID  string
}

// NewAddFactCommand crée une nouvelle commande d'ajout de fait
func NewAddFactCommand(storage Storage, fact *Fact) *AddFactCommand {
	return &AddFactCommand{
		storage: storage,
		fact:    fact,
		factID:  fact.GetInternalID(),
	}
}

// Execute ajoute le fait au storage
func (c *AddFactCommand) Execute() error {
	if err := c.storage.AddFact(c.fact); err != nil {
		return NewCommandError("AddFact", "Execute", err)
	}
	return nil
}

// Undo supprime le fait du storage
func (c *AddFactCommand) Undo() error {
	if err := c.storage.RemoveFact(c.factID); err != nil {
		return NewCommandError("AddFact", "Undo", err)
	}
	return nil
}

// String retourne une description de la commande
func (c *AddFactCommand) String() string {
	return fmt.Sprintf("AddFact(id=%s, type=%s)", c.fact.ID, c.fact.Type)
}

// RemoveFactCommand représente la suppression d'un fait du réseau
type RemoveFactCommand struct {
	storage     Storage
	factID      string
	removedFact *Fact // Sauvegardé pour restauration lors du Undo
}

// NewRemoveFactCommand crée une nouvelle commande de suppression de fait
func NewRemoveFactCommand(storage Storage, factID string) *RemoveFactCommand {
	return &RemoveFactCommand{
		storage: storage,
		factID:  factID,
	}
}

// Execute supprime le fait du storage (après l'avoir sauvegardé)
func (c *RemoveFactCommand) Execute() error {
	// Sauvegarder le fait avant de le supprimer
	fact := c.storage.GetFact(c.factID)
	if fact == nil {
		return NewCommandError("RemoveFact", "Execute",
			fmt.Errorf("fact %s not found", c.factID))
	}

	c.removedFact = fact.Clone()

	// Supprimer le fait
	if err := c.storage.RemoveFact(c.factID); err != nil {
		return NewCommandError("RemoveFact", "Execute", err)
	}

	return nil
}

// Undo restaure le fait supprimé
func (c *RemoveFactCommand) Undo() error {
	if c.removedFact == nil {
		return NewCommandError("RemoveFact", "Undo",
			fmt.Errorf("no fact to restore"))
	}

	if err := c.storage.AddFact(c.removedFact); err != nil {
		return NewCommandError("RemoveFact", "Undo", err)
	}

	return nil
}

// String retourne une description de la commande
func (c *RemoveFactCommand) String() string {
	if c.removedFact != nil {
		return fmt.Sprintf("RemoveFact(id=%s, type=%s)", c.removedFact.ID, c.removedFact.Type)
	}
	return fmt.Sprintf("RemoveFact(id=%s)", c.factID)
}
