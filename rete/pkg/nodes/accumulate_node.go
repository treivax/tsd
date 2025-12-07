// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package nodes

import (
	"fmt"
	"sync"

	"github.com/treivax/tsd/rete/pkg/domain"
)

// AccumulateNodeImpl implements the AccumulateNode interface for aggregations
type AccumulateNodeImpl struct {
	*BaseBetaNode
	accumulator       domain.AccumulateFunction
	accumulatedValues map[string]interface{} // Stores aggregated values per token
	aggregationCalc   *AggregationCalculator
	mu                sync.RWMutex
}

// NewAccumulateNode creates a new accumulation node
func NewAccumulateNode(id string, accumulator domain.AccumulateFunction, logger domain.Logger) *AccumulateNodeImpl {
	baseBeta := NewBaseBetaNode(id, "AccumulateNode", logger)
	return &AccumulateNodeImpl{
		BaseBetaNode:      baseBeta,
		accumulator:       accumulator,
		accumulatedValues: make(map[string]interface{}),
		aggregationCalc:   NewAggregationCalculator(),
	}
}

// SetAccumulator sets the accumulation function
func (a *AccumulateNodeImpl) SetAccumulator(accumulator domain.AccumulateFunction) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.accumulator = accumulator
}

// GetAccumulator returns the accumulation function
func (a *AccumulateNodeImpl) GetAccumulator() domain.AccumulateFunction {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.accumulator
}

// ComputeAggregate calculates the aggregation for a given token
func (a *AccumulateNodeImpl) ComputeAggregate(token *domain.Token, facts []*domain.Fact) (interface{}, error) {
	a.mu.RLock()
	accumulator := a.accumulator
	a.mu.RUnlock()

	if accumulator.FunctionType == "" {
		return nil, fmt.Errorf("no accumulator function defined")
	}

	switch accumulator.FunctionType {
	case AggregateFunctionSum:
		return a.aggregationCalc.ComputeSum(facts, accumulator.Field)
	case AggregateFunctionCount:
		return a.aggregationCalc.ComputeCount(facts), nil
	case AggregateFunctionAvg:
		return a.aggregationCalc.ComputeAverage(facts, accumulator.Field)
	case AggregateFunctionMin:
		return a.aggregationCalc.ComputeMin(facts, accumulator.Field)
	case AggregateFunctionMax:
		return a.aggregationCalc.ComputeMax(facts, accumulator.Field)
	default:
		return nil, fmt.Errorf("unsupported accumulator function: %s", accumulator.FunctionType)
	}
}

// ProcessLeftToken processes a token coming from the left
func (a *AccumulateNodeImpl) ProcessLeftToken(token *domain.Token) error {
	a.logger.Debug("processing token in AccumulateNode", map[string]interface{}{
		"node_id":    a.id,
		"token_id":   token.ID,
		"node_type":  "AccumulateNode",
		"action":     "left_input",
		"fact_count": len(token.Facts),
	})

	// Store the token in left memory
	a.betaMemory.StoreToken(token)

	// Get all right facts for aggregation
	rightFacts := a.betaMemory.GetFacts()

	// Calculate the aggregation
	result, err := a.ComputeAggregate(token, rightFacts)
	if err != nil {
		a.logger.Error("failed to compute aggregate", err, map[string]interface{}{
			"node_id":  a.id,
			"token_id": token.ID,
		})
		return err
	}

	// Store the result
	a.mu.Lock()
	a.accumulatedValues[token.ID] = result
	a.mu.Unlock()

	// Propagate the enriched token with aggregation result
	return a.propagateAggregateToken(token, result)
}

// ProcessRightFact processes a fact coming from the right
func (a *AccumulateNodeImpl) ProcessRightFact(fact *domain.Fact) error {
	a.logger.Debug("processing fact in AccumulateNode", map[string]interface{}{
		"node_id":   a.id,
		"fact_id":   fact.ID,
		"fact_type": fact.Type,
		"node_type": "AccumulateNode",
		"action":    "right_input",
	})

	// Store the fact in right memory
	a.betaMemory.StoreFact(fact)

	// Recalculate the aggregation for all left tokens
	leftTokens := a.betaMemory.GetTokens()
	for _, token := range leftTokens {
		if err := a.recomputeAndPropagateAggregate(token); err != nil {
			a.logger.Error("failed to recompute aggregate for token", err, map[string]interface{}{
				"node_id":  a.id,
				"token_id": token.ID,
				"fact_id":  fact.ID,
			})
		}
	}

	return nil
}

// recomputeAndPropagateAggregate recalculates the aggregate for a token and propagates if changed
func (a *AccumulateNodeImpl) recomputeAndPropagateAggregate(token *domain.Token) error {
	// Get all right facts
	rightFacts := a.betaMemory.GetFacts()

	// Recalculate the aggregation
	result, err := a.ComputeAggregate(token, rightFacts)
	if err != nil {
		return fmt.Errorf("failed to compute aggregate: %w", err)
	}

	// Update the result
	a.mu.Lock()
	oldResult, existed := a.accumulatedValues[token.ID]
	a.accumulatedValues[token.ID] = result
	a.mu.Unlock()

	// If the result changed, propagate the update
	if !existed || oldResult != result {
		a.logger.Debug("aggregate result updated", map[string]interface{}{
			"node_id":    a.id,
			"token_id":   token.ID,
			"old_result": oldResult,
			"new_result": result,
		})

		return a.propagateAggregateToken(token, result)
	}

	return nil
}

// propagateAggregateToken creates and propagates a token with aggregate result
func (a *AccumulateNodeImpl) propagateAggregateToken(token *domain.Token, result interface{}) error {
	newToken := &domain.Token{
		ID: fmt.Sprintf("%s_agg", token.ID),
		Facts: append(token.Facts, &domain.Fact{
			ID:   fmt.Sprintf("agg_%s", token.ID),
			Type: "AggregateResult",
			Fields: map[string]interface{}{
				"function": a.accumulator.FunctionType,
				"value":    result,
			},
		}),
	}

	return a.propagateTokenToChildren(newToken)
}
