// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package servercmd

import (
	"fmt"
	"sync"

	"github.com/treivax/tsd/rete"
	"github.com/treivax/tsd/tsdio"
)

// ExecutionStatsCollector collecte les statistiques d'exécution des actions.
//
// Implémente ActionObserver pour capturer toutes les exécutions
// et les convertir en format tsdio.Activation pour l'API.
//
// Thread-Safety :
//   - Thread-safe grâce au mutex interne
//   - Peut être utilisé par plusieurs terminal nodes en parallèle
type ExecutionStatsCollector struct {
	executions []rete.ExecutionResult
	mu         sync.RWMutex
}

// NewExecutionStatsCollector crée un nouveau collecteur.
func NewExecutionStatsCollector() *ExecutionStatsCollector {
	return &ExecutionStatsCollector{
		executions: make([]rete.ExecutionResult, 0),
	}
}

// OnActionExecuted capture un résultat d'exécution.
// Implémente rete.ActionObserver.
func (c *ExecutionStatsCollector) OnActionExecuted(result rete.ExecutionResult) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.executions = append(c.executions, result)
}

// GetExecutions retourne une copie de tous les résultats capturés.
func (c *ExecutionStatsCollector) GetExecutions() []rete.ExecutionResult {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Retourner une copie pour éviter modifications concurrentes
	return append([]rete.ExecutionResult{}, c.executions...)
}

// GetActivations convertit les résultats en format tsdio.Activation.
func (c *ExecutionStatsCollector) GetActivations() []tsdio.Activation {
	c.mu.RLock()
	defer c.mu.RUnlock()

	activations := make([]tsdio.Activation, 0, len(c.executions))

	for _, exec := range c.executions {
		activation := tsdio.Activation{
			ActionName:      exec.Context.ActionName,
			Arguments:       formatArguments(exec.Arguments),
			TriggeringFacts: extractFacts(exec.Context.Token),
			BindingsCount:   len(exec.Context.Token.Facts),
		}
		activations = append(activations, activation)
	}

	return activations
}

// GetExecutionCount retourne le nombre d'exécutions capturées.
func (c *ExecutionStatsCollector) GetExecutionCount() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.executions)
}

// Reset réinitialise le collecteur.
func (c *ExecutionStatsCollector) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.executions = make([]rete.ExecutionResult, 0)
}

// Helper functions

func formatArguments(args []interface{}) []tsdio.ArgumentValue {
	result := make([]tsdio.ArgumentValue, 0, len(args))
	for _, arg := range args {
		result = append(result, tsdio.ArgumentValue{
			Type:  "value",
			Value: fmt.Sprintf("%v", arg),
		})
	}
	return result
}

func extractFacts(token *rete.Token) []tsdio.Fact {
	if token == nil {
		return []tsdio.Fact{}
	}

	facts := make([]tsdio.Fact, 0, len(token.Facts))
	for _, fact := range token.Facts {
		if fact == nil {
			continue
		}

		// Extraire les attributs du fait
		fields := make(map[string]interface{})
		if fact.Attributes != nil {
			for k, v := range fact.Attributes {
				fields[k] = v
			}
		}

		f := tsdio.Fact{
			Type:   fact.Type,
			Fields: fields,
		}
		f.SetInternalID(fact.ID)
		facts = append(facts, f)
	}

	return facts
}
