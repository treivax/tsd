// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Transaction gère le cycle de vie d'une transaction sur le réseau RETE
// Utilise le Command Pattern pour enregistrer et annuler les opérations
// Cette implémentation remplace l'ancienne approche par snapshot qui doublait la mémoire
type Transaction struct {
	ID           string
	Network      *ReteNetwork
	Commands     []Command           // Log des commandes exécutées (pour rollback)
	Options      *TransactionOptions // Configuration de la transaction (Strong mode)
	IsActive     bool
	IsCommitted  bool
	IsRolledBack bool
	StartTime    time.Time
	mutex        sync.RWMutex
}

// BeginTransaction démarre une nouvelle transaction sur le réseau avec options par défaut
// Cette opération est O(1) en temps et mémoire (vs O(N) avec snapshot)
func (network *ReteNetwork) BeginTransaction() *Transaction {
	return network.BeginTransactionWithOptions(nil)
}

// BeginTransactionWithOptions démarre une nouvelle transaction avec options personnalisées
func (network *ReteNetwork) BeginTransactionWithOptions(opts *TransactionOptions) *Transaction {
	if opts == nil {
		opts = DefaultTransactionOptions()
	}

	return &Transaction{
		ID:           uuid.New().String(),
		Network:      network,
		Commands:     make([]Command, 0, 16), // Pré-allocation raisonnable
		Options:      opts,
		IsActive:     true,
		IsCommitted:  false,
		IsRolledBack: false,
		StartTime:    time.Now(),
	}
}

// RecordAndExecute enregistre et exécute une commande dans la transaction
// La commande est exécutée immédiatement et enregistrée pour un potentiel rollback
func (tx *Transaction) RecordAndExecute(cmd Command) error {
	tx.mutex.Lock()
	defer tx.mutex.Unlock()

	if !tx.IsActive {
		return fmt.Errorf("transaction %s is not active", tx.ID)
	}

	// Exécuter la commande
	if err := cmd.Execute(); err != nil {
		return fmt.Errorf("command execution failed: %w", err)
	}

	// Enregistrer pour potentiel rollback
	tx.Commands = append(tx.Commands, cmd)

	return nil
}

// Commit valide les modifications de la transaction
// Les commandes sont déjà exécutées, donc commit ne fait que marquer la transaction
// comme terminée et libère le log des commandes
func (tx *Transaction) Commit() error {
	tx.mutex.Lock()
	defer tx.mutex.Unlock()

	if !tx.IsActive {
		return fmt.Errorf("transaction %s is not active", tx.ID)
	}

	if tx.IsCommitted {
		return fmt.Errorf("transaction %s already committed", tx.ID)
	}

	if tx.IsRolledBack {
		return fmt.Errorf("transaction %s already rolled back", tx.ID)
	}

	// Commit : rien à faire, les commandes sont déjà exécutées
	tx.IsActive = false
	tx.IsCommitted = true

	// Libérer le log des commandes (plus besoin pour rollback)
	tx.Commands = nil

	return nil
}

// Rollback annule les modifications en rejouant les commandes à l'envers
// Utilise le pattern de "rejeu inversé" : chaque commande est Undo() dans l'ordre inverse
// Complexité : O(k) où k = nombre de commandes (vs O(N) pour restaurer snapshot complet)
func (tx *Transaction) Rollback() error {
	tx.mutex.Lock()
	defer tx.mutex.Unlock()

	if !tx.IsActive {
		return fmt.Errorf("transaction %s is not active", tx.ID)
	}

	if tx.IsCommitted {
		return fmt.Errorf("transaction %s already committed, cannot rollback", tx.ID)
	}

	if tx.IsRolledBack {
		return fmt.Errorf("transaction %s already rolled back", tx.ID)
	}

	// Rejouer les commandes EN ORDRE INVERSE (rejeu inversé)
	// Si une commande échoue, on arrête le rollback et on retourne l'erreur
	for i := len(tx.Commands) - 1; i >= 0; i-- {
		if err := tx.Commands[i].Undo(); err != nil {
			// CRITIQUE : Rollback partiel échoué
			// Le réseau peut être dans un état inconsistant
			return fmt.Errorf("rollback failed at command %d (%s): %w",
				i, tx.Commands[i].String(), err)
		}
	}

	tx.IsActive = false
	tx.IsRolledBack = true
	tx.Commands = nil

	return nil
}

// GetCommandCount retourne le nombre de commandes exécutées dans la transaction
func (tx *Transaction) GetCommandCount() int {
	tx.mutex.RLock()
	defer tx.mutex.RUnlock()
	return len(tx.Commands)
}

// GetDuration retourne la durée de la transaction depuis sa création
func (tx *Transaction) GetDuration() time.Duration {
	return time.Since(tx.StartTime)
}

// GetCommands retourne une copie de la liste des commandes (pour debugging/logging)
func (tx *Transaction) GetCommands() []Command {
	tx.mutex.RLock()
	defer tx.mutex.RUnlock()

	if tx.Commands == nil {
		return nil
	}

	commands := make([]Command, len(tx.Commands))
	copy(commands, tx.Commands)
	return commands
}

// String retourne une représentation textuelle de la transaction
func (tx *Transaction) String() string {
	tx.mutex.RLock()
	defer tx.mutex.RUnlock()

	status := "active"
	if tx.IsCommitted {
		status = "committed"
	} else if tx.IsRolledBack {
		status = "rolled back"
	}

	return fmt.Sprintf("Transaction{ID: %s, Status: %s, Commands: %d, Duration: %v}",
		tx.ID, status, len(tx.Commands), tx.GetDuration())
}

// GetMemoryFootprint retourne une estimation de l'empreinte mémoire de la transaction
// Utile pour le monitoring et les métriques
func (tx *Transaction) GetMemoryFootprint() int64 {
	tx.mutex.RLock()
	defer tx.mutex.RUnlock()

	// Estimation : chaque commande occupe environ 200 bytes
	// (pointeurs + métadonnées + petit overhead)
	return int64(len(tx.Commands) * 200)
}
