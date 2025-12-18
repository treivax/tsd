// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"sync"
	"testing"
)

// TestActionObserver capture les exécutions pour assertions dans tests.
type TestActionObserver struct {
	t          *testing.T
	executions []ExecutionResult
	mu         sync.RWMutex
}

// NewTestActionObserver crée un observateur de test.
func NewTestActionObserver(t *testing.T) *TestActionObserver {
	return &TestActionObserver{
		t:          t,
		executions: make([]ExecutionResult, 0),
	}
}

// OnActionExecuted capture l'exécution.
func (o *TestActionObserver) OnActionExecuted(result ExecutionResult) {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.t.Logf("🎯 Action executed: %s (rule: %s, success: %v, duration: %v)",
		result.Context.ActionName,
		result.Context.RuleName,
		result.Success,
		result.Duration)

	o.executions = append(o.executions, result)
}

// GetExecutions retourne tous les résultats capturés.
func (o *TestActionObserver) GetExecutions() []ExecutionResult {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return append([]ExecutionResult{}, o.executions...)
}

// GetExecutionCount retourne le nombre d'exécutions.
func (o *TestActionObserver) GetExecutionCount() int {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return len(o.executions)
}

// AssertExecutionCount vérifie le nombre d'exécutions.
func (o *TestActionObserver) AssertExecutionCount(expected int) {
	o.mu.RLock()
	count := len(o.executions)
	o.mu.RUnlock()

	if count != expected {
		o.t.Errorf("❌ Expected %d executions, got %d", expected, count)
	}
}

// AssertActionExecuted vérifie qu'une action a été exécutée.
func (o *TestActionObserver) AssertActionExecuted(actionName string) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	for _, exec := range o.executions {
		if exec.Context.ActionName == actionName {
			return
		}
	}
	o.t.Errorf("❌ Action '%s' was not executed", actionName)
}

// AssertAllSuccessful vérifie que toutes les exécutions ont réussi.
func (o *TestActionObserver) AssertAllSuccessful() {
	o.mu.RLock()
	defer o.mu.RUnlock()

	for i, exec := range o.executions {
		if !exec.Success {
			o.t.Errorf("❌ Execution %d failed: %v", i, exec.Error)
		}
	}
}

// Reset réinitialise le collecteur.
func (o *TestActionObserver) Reset() {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.executions = make([]ExecutionResult, 0)
}

// GetSuccessfulExecutions retourne uniquement les exécutions réussies.
func (o *TestActionObserver) GetSuccessfulExecutions() []ExecutionResult {
	o.mu.RLock()
	defer o.mu.RUnlock()

	successful := make([]ExecutionResult, 0)
	for _, exec := range o.executions {
		if exec.Success {
			successful = append(successful, exec)
		}
	}
	return successful
}

// GetFailedExecutions retourne uniquement les exécutions échouées.
func (o *TestActionObserver) GetFailedExecutions() []ExecutionResult {
	o.mu.RLock()
	defer o.mu.RUnlock()

	failed := make([]ExecutionResult, 0)
	for _, exec := range o.executions {
		if !exec.Success {
			failed = append(failed, exec)
		}
	}
	return failed
}
