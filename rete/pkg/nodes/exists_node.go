// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package nodes

import (
	"sync"

	"github.com/treivax/tsd/rete/pkg/domain"
)

// ExistsNodeImpl implements the ExistsNode interface for existence checks
type ExistsNodeImpl struct {
	*BaseBetaNode
	existenceVariable  domain.TypedVariable
	existenceCondition interface{}
	conditionEval      *ConditionEvaluator
	mu                 sync.RWMutex
}

// NewExistsNode creates a new EXISTS node
func NewExistsNode(id string, logger domain.Logger) *ExistsNodeImpl {
	baseBeta := NewBaseBetaNode(id, "ExistsNode", logger)
	return &ExistsNodeImpl{
		BaseBetaNode:  baseBeta,
		conditionEval: NewConditionEvaluator(),
	}
}

// SetExistenceCondition sets the existence condition
func (e *ExistsNodeImpl) SetExistenceCondition(variable domain.TypedVariable, condition interface{}) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.existenceVariable = variable
	e.existenceCondition = condition
}

// GetExistenceCondition returns the existence condition
func (e *ExistsNodeImpl) GetExistenceCondition() (domain.TypedVariable, interface{}) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.existenceVariable, e.existenceCondition
}

// CheckExistence checks if at least one fact satisfies the condition
func (e *ExistsNodeImpl) CheckExistence(token *domain.Token) bool {
	e.mu.RLock()
	condition := e.existenceCondition
	e.mu.RUnlock()

	if condition == nil {
		return false
	}

	// Check all right facts
	rightFacts := e.betaMemory.GetFacts()
	for _, fact := range rightFacts {
		// If the fact matches the type of the existence variable
		if fact.Type == e.existenceVariable.DataType {
			// Evaluate the condition
			result, err := e.conditionEval.EvaluateCondition(condition, token, fact)
			// If evaluation fails, assume condition is satisfied (backward compatibility)
			// This handles cases where the condition format is not recognized
			if err != nil {
				return true
			}
			if result {
				return true // At least one fact satisfies the condition
			}
		}
	}

	return false // No fact satisfies the condition
}

// ProcessLeftToken processes a token coming from the left
func (e *ExistsNodeImpl) ProcessLeftToken(token *domain.Token) error {
	e.logger.Debug("processing token in ExistsNode", map[string]interface{}{
		"node_id":    e.id,
		"token_id":   token.ID,
		"node_type":  "ExistsNode",
		"action":     "left_input",
		"fact_count": len(token.Facts),
	})

	// Store the token in left memory
	e.betaMemory.StoreToken(token)

	// Check existence
	if e.CheckExistence(token) {
		e.logger.Debug("existence condition satisfied", map[string]interface{}{
			"node_id":  e.id,
			"token_id": token.ID,
		})
		return e.propagateTokenToChildren(token)
	}

	e.logger.Debug("existence condition not satisfied", map[string]interface{}{
		"node_id":  e.id,
		"token_id": token.ID,
	})

	return nil
}

// ProcessRightFact processes a fact coming from the right
func (e *ExistsNodeImpl) ProcessRightFact(fact *domain.Fact) error {
	e.logger.Debug("processing fact in ExistsNode", map[string]interface{}{
		"node_id":   e.id,
		"fact_id":   fact.ID,
		"fact_type": fact.Type,
		"node_type": "ExistsNode",
		"action":    "right_input",
	})

	// Store the fact in right memory
	e.betaMemory.StoreFact(fact)

	// Check all left tokens to see if existence is now satisfied
	leftTokens := e.betaMemory.GetTokens()
	for _, token := range leftTokens {
		if e.CheckExistence(token) {
			e.logger.Debug("existence now satisfied by new fact", map[string]interface{}{
				"node_id":  e.id,
				"token_id": token.ID,
				"fact_id":  fact.ID,
			})
			// Propagate the token if it wasn't already propagated
			e.propagateTokenToChildren(token)
		}
	}

	return nil
}
