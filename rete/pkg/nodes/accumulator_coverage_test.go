package nodes

import (
	"testing"

	"github.com/treivax/tsd/rete/pkg/domain"
)

// TestAccumulateNode_ComputeSum_EdgeCases tests edge cases in sum computation
func TestAccumulateNode_ComputeSum_EdgeCases(t *testing.T) {
	logger := newMockLogger()

	tests := []struct {
		name      string
		facts     []*domain.Fact
		field     string
		wantValue float64
		wantErr   bool
	}{
		{
			name:      "empty facts",
			facts:     []*domain.Fact{},
			field:     "value",
			wantValue: 0,
			wantErr:   false,
		},
		{
			name: "all facts missing field",
			facts: []*domain.Fact{
				domain.NewFact("f1", "Test", map[string]interface{}{"other": 10}),
				domain.NewFact("f2", "Test", map[string]interface{}{"other": 20}),
			},
			field:     "value",
			wantValue: 0,
			wantErr:   false,
		},
		{
			name: "mixed numeric types",
			facts: []*domain.Fact{
				domain.NewFact("f1", "Test", map[string]interface{}{"value": 10}),
				domain.NewFact("f2", "Test", map[string]interface{}{"value": int64(20)}),
				domain.NewFact("f3", "Test", map[string]interface{}{"value": float32(5.5)}),
				domain.NewFact("f4", "Test", map[string]interface{}{"value": 3.5}),
			},
			field:     "value",
			wantValue: 39.0,
			wantErr:   false,
		},
		{
			name: "some facts with non-numeric values",
			facts: []*domain.Fact{
				domain.NewFact("f1", "Test", map[string]interface{}{"value": 10}),
				domain.NewFact("f2", "Test", map[string]interface{}{"value": "not a number"}),
				domain.NewFact("f3", "Test", map[string]interface{}{"value": 20}),
			},
			field:     "value",
			wantValue: 30,
			wantErr:   false,
		},
		{
			name: "negative values",
			facts: []*domain.Fact{
				domain.NewFact("f1", "Test", map[string]interface{}{"value": -10}),
				domain.NewFact("f2", "Test", map[string]interface{}{"value": -5}),
				domain.NewFact("f3", "Test", map[string]interface{}{"value": 15}),
			},
			field:     "value",
			wantValue: 0,
			wantErr:   false,
		},
		{
			name: "all negative values",
			facts: []*domain.Fact{
				domain.NewFact("f1", "Test", map[string]interface{}{"value": -10}),
				domain.NewFact("f2", "Test", map[string]interface{}{"value": -5}),
				domain.NewFact("f3", "Test", map[string]interface{}{"value": -3}),
			},
			field:     "value",
			wantValue: -18,
			wantErr:   false,
		},
		{
			name: "zeros",
			facts: []*domain.Fact{
				domain.NewFact("f1", "Test", map[string]interface{}{"value": 0}),
				domain.NewFact("f2", "Test", map[string]interface{}{"value": 0}),
				domain.NewFact("f3", "Test", map[string]interface{}{"value": 0}),
			},
			field:     "value",
			wantValue: 0,
			wantErr:   false,
		},
		{
			name: "single value",
			facts: []*domain.Fact{
				domain.NewFact("f1", "Test", map[string]interface{}{"value": 42}),
			},
			field:     "value",
			wantValue: 42,
			wantErr:   false,
		},
		{
			name: "very large values",
			facts: []*domain.Fact{
				domain.NewFact("f1", "Test", map[string]interface{}{"value": int64(1000000000000)}),
				domain.NewFact("f2", "Test", map[string]interface{}{"value": int64(2000000000000)}),
			},
			field:     "value",
			wantValue: 3000000000000.0,
			wantErr:   false,
		},
		{
			name: "fractional values",
			facts: []*domain.Fact{
				domain.NewFact("f1", "Test", map[string]interface{}{"value": 0.1}),
				domain.NewFact("f2", "Test", map[string]interface{}{"value": 0.2}),
				domain.NewFact("f3", "Test", map[string]interface{}{"value": 0.3}),
			},
			field:     "value",
			wantValue: 0.6,
			wantErr:   false,
		},
		{
			name: "some facts with field, some without",
			facts: []*domain.Fact{
				domain.NewFact("f1", "Test", map[string]interface{}{"value": 10}),
				domain.NewFact("f2", "Test", map[string]interface{}{"other": 5}),
				domain.NewFact("f3", "Test", map[string]interface{}{"value": 20}),
				domain.NewFact("f4", "Test", map[string]interface{}{"other": 3}),
				domain.NewFact("f5", "Test", map[string]interface{}{"value": 5}),
			},
			field:     "value",
			wantValue: 35,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accumulator := domain.AccumulateFunction{
				FunctionType: AggregateFunctionSum,
				Field:        tt.field,
			}
			node := NewAccumulateNode("acc1", accumulator, logger)

			got, err := node.computeSum(tt.facts, tt.field)

			if tt.wantErr {
				if err == nil {
					t.Error("computeSum() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("computeSum() unexpected error: %v", err)
				return
			}

			// Allow small floating point precision differences
			diff := got - tt.wantValue
			if diff < 0 {
				diff = -diff
			}
			if diff > 0.0001 {
				t.Errorf("computeSum() = %v, want %v", got, tt.wantValue)
			}
		})
	}
}

// TestAccumulateNode_ComputeAverage_EdgeCases tests edge cases in average computation
func TestAccumulateNode_ComputeAverage_EdgeCases(t *testing.T) {
	logger := newMockLogger()

	tests := []struct {
		name      string
		facts     []*domain.Fact
		field     string
		wantValue float64
		wantErr   bool
	}{
		{
			name:      "empty facts",
			facts:     []*domain.Fact{},
			field:     "value",
			wantValue: 0,
			wantErr:   false,
		},
		{
			name: "all facts missing field",
			facts: []*domain.Fact{
				domain.NewFact("f1", "Test", map[string]interface{}{"other": 10}),
				domain.NewFact("f2", "Test", map[string]interface{}{"other": 20}),
			},
			field:     "value",
			wantValue: 0,
			wantErr:   false,
		},
		{
			name: "single value",
			facts: []*domain.Fact{
				domain.NewFact("f1", "Test", map[string]interface{}{"value": 42}),
			},
			field:     "value",
			wantValue: 42,
			wantErr:   false,
		},
		{
			name: "even average",
			facts: []*domain.Fact{
				domain.NewFact("f1", "Test", map[string]interface{}{"value": 10}),
				domain.NewFact("f2", "Test", map[string]interface{}{"value": 20}),
				domain.NewFact("f3", "Test", map[string]interface{}{"value": 30}),
			},
			field:     "value",
			wantValue: 20,
			wantErr:   false,
		},
		{
			name: "fractional average",
			facts: []*domain.Fact{
				domain.NewFact("f1", "Test", map[string]interface{}{"value": 10}),
				domain.NewFact("f2", "Test", map[string]interface{}{"value": 20}),
			},
			field:     "value",
			wantValue: 15,
			wantErr:   false,
		},
		{
			name: "negative values",
			facts: []*domain.Fact{
				domain.NewFact("f1", "Test", map[string]interface{}{"value": -10}),
				domain.NewFact("f2", "Test", map[string]interface{}{"value": -20}),
			},
			field:     "value",
			wantValue: -15,
			wantErr:   false,
		},
		{
			name: "mixed positive and negative",
			facts: []*domain.Fact{
				domain.NewFact("f1", "Test", map[string]interface{}{"value": -10}),
				domain.NewFact("f2", "Test", map[string]interface{}{"value": 10}),
				domain.NewFact("f3", "Test", map[string]interface{}{"value": 5}),
			},
			field:     "value",
			wantValue: 5.0 / 3.0,
			wantErr:   false,
		},
		{
			name: "mixed numeric types",
			facts: []*domain.Fact{
				domain.NewFact("f1", "Test", map[string]interface{}{"value": 10}),
				domain.NewFact("f2", "Test", map[string]interface{}{"value": int64(20)}),
				domain.NewFact("f3", "Test", map[string]interface{}{"value": float32(15.0)}),
				domain.NewFact("f4", "Test", map[string]interface{}{"value": 15.0}),
			},
			field:     "value",
			wantValue: 15,
			wantErr:   false,
		},
		{
			name: "some non-numeric values skipped",
			facts: []*domain.Fact{
				domain.NewFact("f1", "Test", map[string]interface{}{"value": 10}),
				domain.NewFact("f2", "Test", map[string]interface{}{"value": "not a number"}),
				domain.NewFact("f3", "Test", map[string]interface{}{"value": 20}),
			},
			field:     "value",
			wantValue: 10, // sum=30 (10+0+20), count=3 (all have field), avg=10
			wantErr:   false,
		},
		{
			name: "some facts missing field",
			facts: []*domain.Fact{
				domain.NewFact("f1", "Test", map[string]interface{}{"value": 10}),
				domain.NewFact("f2", "Test", map[string]interface{}{"other": 100}),
				domain.NewFact("f3", "Test", map[string]interface{}{"value": 20}),
				domain.NewFact("f4", "Test", map[string]interface{}{"other": 200}),
				domain.NewFact("f5", "Test", map[string]interface{}{"value": 30}),
			},
			field:     "value",
			wantValue: 20,
			wantErr:   false,
		},
		{
			name: "zeros",
			facts: []*domain.Fact{
				domain.NewFact("f1", "Test", map[string]interface{}{"value": 0}),
				domain.NewFact("f2", "Test", map[string]interface{}{"value": 0}),
			},
			field:     "value",
			wantValue: 0,
			wantErr:   false,
		},
		{
			name: "very small values",
			facts: []*domain.Fact{
				domain.NewFact("f1", "Test", map[string]interface{}{"value": 0.001}),
				domain.NewFact("f2", "Test", map[string]interface{}{"value": 0.002}),
				domain.NewFact("f3", "Test", map[string]interface{}{"value": 0.003}),
			},
			field:     "value",
			wantValue: 0.002,
			wantErr:   false,
		},
		{
			name: "very large values",
			facts: []*domain.Fact{
				domain.NewFact("f1", "Test", map[string]interface{}{"value": int64(1000000000000)}),
				domain.NewFact("f2", "Test", map[string]interface{}{"value": int64(2000000000000)}),
				domain.NewFact("f3", "Test", map[string]interface{}{"value": int64(3000000000000)}),
			},
			field:     "value",
			wantValue: 2000000000000.0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accumulator := domain.AccumulateFunction{
				FunctionType: AggregateFunctionAvg,
				Field:        tt.field,
			}
			node := NewAccumulateNode("acc1", accumulator, logger)

			got, err := node.computeAverage(tt.facts, tt.field)

			if tt.wantErr {
				if err == nil {
					t.Error("computeAverage() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("computeAverage() unexpected error: %v", err)
				return
			}

			// Allow small floating point precision differences
			diff := got - tt.wantValue
			if diff < 0 {
				diff = -diff
			}
			if diff > 0.0001 {
				t.Errorf("computeAverage() = %v, want %v", got, tt.wantValue)
			}
		})
	}
}

// TestAccumulateNode_RecomputeAndPropagateAggregate tests recomputation edge cases
func TestAccumulateNode_RecomputeAndPropagateAggregate(t *testing.T) {
	logger := newMockLogger()

	tests := []struct {
		name        string
		accumulator domain.AccumulateFunction
		setupNode   func(*AccumulateNodeImpl)
		token       *domain.Token
		expectError bool
		description string
	}{
		{
			name: "recompute sum with no children",
			accumulator: domain.AccumulateFunction{
				FunctionType: AggregateFunctionSum,
				Field:        "value",
			},
			setupNode: func(node *AccumulateNodeImpl) {
				node.betaMemory.StoreFact(domain.NewFact("f1", "Test", map[string]interface{}{"value": 10}))
				node.betaMemory.StoreFact(domain.NewFact("f2", "Test", map[string]interface{}{"value": 20}))
			},
			token:       domain.NewToken("t1", "node1", []*domain.Fact{}),
			expectError: false,
			description: "should handle recomputation with no children",
		},
		{
			name: "recompute average",
			accumulator: domain.AccumulateFunction{
				FunctionType: AggregateFunctionAvg,
				Field:        "value",
			},
			setupNode: func(node *AccumulateNodeImpl) {
				node.betaMemory.StoreFact(domain.NewFact("f1", "Test", map[string]interface{}{"value": 10}))
				node.betaMemory.StoreFact(domain.NewFact("f2", "Test", map[string]interface{}{"value": 20}))
				node.betaMemory.StoreFact(domain.NewFact("f3", "Test", map[string]interface{}{"value": 30}))
			},
			token:       domain.NewToken("t1", "node1", []*domain.Fact{}),
			expectError: false,
			description: "should recompute average correctly",
		},
		{
			name: "recompute min",
			accumulator: domain.AccumulateFunction{
				FunctionType: AggregateFunctionMin,
				Field:        "value",
			},
			setupNode: func(node *AccumulateNodeImpl) {
				node.betaMemory.StoreFact(domain.NewFact("f1", "Test", map[string]interface{}{"value": 10}))
				node.betaMemory.StoreFact(domain.NewFact("f2", "Test", map[string]interface{}{"value": 5}))
				node.betaMemory.StoreFact(domain.NewFact("f3", "Test", map[string]interface{}{"value": 15}))
			},
			token:       domain.NewToken("t1", "node1", []*domain.Fact{}),
			expectError: false,
			description: "should recompute min correctly",
		},
		{
			name: "recompute max",
			accumulator: domain.AccumulateFunction{
				FunctionType: AggregateFunctionMax,
				Field:        "value",
			},
			setupNode: func(node *AccumulateNodeImpl) {
				node.betaMemory.StoreFact(domain.NewFact("f1", "Test", map[string]interface{}{"value": 10}))
				node.betaMemory.StoreFact(domain.NewFact("f2", "Test", map[string]interface{}{"value": 5}))
				node.betaMemory.StoreFact(domain.NewFact("f3", "Test", map[string]interface{}{"value": 15}))
			},
			token:       domain.NewToken("t1", "node1", []*domain.Fact{}),
			expectError: false,
			description: "should recompute max correctly",
		},
		{
			name: "recompute count",
			accumulator: domain.AccumulateFunction{
				FunctionType: AggregateFunctionCount,
				Field:        "value",
			},
			setupNode: func(node *AccumulateNodeImpl) {
				node.betaMemory.StoreFact(domain.NewFact("f1", "Test", map[string]interface{}{"value": 10}))
				node.betaMemory.StoreFact(domain.NewFact("f2", "Test", map[string]interface{}{"other": 20}))
				node.betaMemory.StoreFact(domain.NewFact("f3", "Test", map[string]interface{}{"value": 30}))
			},
			token:       domain.NewToken("t1", "node1", []*domain.Fact{}),
			expectError: false,
			description: "should recompute count correctly",
		},
		{
			name: "recompute with empty facts",
			accumulator: domain.AccumulateFunction{
				FunctionType: AggregateFunctionSum,
				Field:        "value",
			},
			setupNode: func(node *AccumulateNodeImpl) {
				// No facts added
			},
			token:       domain.NewToken("t1", "node1", []*domain.Fact{}),
			expectError: false,
			description: "should handle empty facts in recomputation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := NewAccumulateNode("acc1", tt.accumulator, logger)
			if tt.setupNode != nil {
				tt.setupNode(node)
			}

			err := node.recomputeAndPropagateAggregate(tt.token)

			if tt.expectError && err == nil {
				t.Errorf("recomputeAndPropagateAggregate() expected error but got none: %s", tt.description)
			}
			if !tt.expectError && err != nil {
				t.Errorf("recomputeAndPropagateAggregate() unexpected error: %v (%s)", err, tt.description)
			}
		})
	}
}
