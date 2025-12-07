// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package nodes

import (
	"sync"

	"github.com/treivax/tsd/rete/pkg/domain"
)

// NotNodeImpl implements the NotNode interface for negation
type NotNodeImpl struct {
	*BaseBetaNode
	negationCondition interface{}
	conditionEval     *ConditionEvaluator
	mu                sync.RWMutex
}

// NewNotNode creates a new negation node
func NewNotNode(id string, logger domain.Logger) *NotNodeImpl {
	baseBeta := NewBaseBetaNode(id, "NotNode", logger)
	return &NotNodeImpl{
		BaseBetaNode:  baseBeta,
		conditionEval: NewConditionEvaluator(),
	}
}

// SetNegationCondition sets the negation condition
func (n *NotNodeImpl) SetNegationCondition(condition interface{}) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.negationCondition = condition
}

// GetNegationCondition returns the negation condition
func (n *NotNodeImpl) GetNegationCondition() interface{} {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.negationCondition
}

// ProcessNegation evaluates the negation of a condition
func (n *NotNodeImpl) ProcessNegation(token *domain.Token, fact *domain.Fact) bool {
	n.mu.RLock()
	condition := n.negationCondition
	n.mu.RUnlock()

	if condition == nil {
		return false
	}

	// Evaluate the condition and return its negation
	result, err := n.conditionEval.EvaluateCondition(condition, token, fact)
	if err != nil {
		n.logger.Error("Erreur évaluation condition négation", err, map[string]interface{}{
			"node_id": n.id,
			"token":   token.ID,
			"fact":    fact.ID,
		})
		return false
	}

	return !result // Negation of the result
}

// ProcessLeftToken processes a token coming from the left
func (n *NotNodeImpl) ProcessLeftToken(token *domain.Token) error {
	n.logger.Debug("processing token in NotNode", map[string]interface{}{
		"node_id":    n.id,
		"token_id":   token.ID,
		"node_type":  "NotNode",
		"action":     "left_input",
		"fact_count": len(token.Facts),
	})

	// Store the token in left memory
	n.betaMemory.StoreToken(token)

	// Check negation against all right facts
	rightFacts := n.betaMemory.GetFacts()
	shouldPropagate := true

	for _, fact := range rightFacts {
		if n.ProcessNegation(token, fact) {
			// If negation is true (original condition false), continue
			continue
		} else {
			// If negation is false (original condition true), block propagation
			shouldPropagate = false
			break
		}
	}

	// If no right fact satisfies the condition, propagate the token (negation succeeded)
	if shouldPropagate && len(rightFacts) > 0 {
		return n.propagateTokenToChildren(token)
	}

	// If no right facts, also propagate (default negation)
	if len(rightFacts) == 0 {
		return n.propagateTokenToChildren(token)
	}

	return nil
}

// ProcessRightFact processes a fact coming from the right
func (n *NotNodeImpl) ProcessRightFact(fact *domain.Fact) error {
	n.logger.Debug("processing fact in NotNode", map[string]interface{}{
		"node_id":   n.id,
		"fact_id":   fact.ID,
		"fact_type": fact.Type,
		"node_type": "NotNode",
		"action":    "right_input",
	})

	// Store the fact in right memory
	n.betaMemory.StoreFact(fact)

	// Check all left tokens
	leftTokens := n.betaMemory.GetTokens()
	for _, token := range leftTokens {
		if !n.ProcessNegation(token, fact) {
			// If negation fails (condition true), remove the token if it was propagated
			n.logger.Debug("negation failed, blocking token", map[string]interface{}{
				"node_id":  n.id,
				"token_id": token.ID,
				"fact_id":  fact.ID,
			})
		}
	}

	return nil
}
