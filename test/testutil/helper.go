// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

// Package testutil provides centralized testing utilities for the TSD project
package testutil

import (
	"testing"

	"github.com/treivax/tsd/rete"
)

// TestHelper provides common testing utilities and helpers
type TestHelper struct {
	pipeline *rete.ConstraintPipeline
}

// NewTestHelper creates a new test helper instance
func NewTestHelper() *TestHelper {
	return &TestHelper{
		pipeline: rete.NewConstraintPipeline(),
	}
}

// BuildNetworkFromConstraintFile builds a RETE network from constraint file
// This function MUST be used by ALL tests using .constraint files
func (th *TestHelper) BuildNetworkFromConstraintFile(t *testing.T, constraintFile string) (*rete.ReteNetwork, rete.Storage) {
	storage := rete.NewMemoryStorage()
	network, err := th.pipeline.BuildNetworkFromConstraintFile(constraintFile, storage)
	if err != nil {
		t.Fatalf("Failed to build network: %v", err)
	}
	return network, storage
}

// BuildNetworkFromConstraintFileWithFacts builds a RETE network and loads facts
func (th *TestHelper) BuildNetworkFromConstraintFileWithFacts(t *testing.T, constraintFile, factsFile string) (*rete.ReteNetwork, []*rete.Fact, rete.Storage) {
	storage := rete.NewMemoryStorage()
	network, facts, err := th.pipeline.BuildNetworkFromConstraintFileWithFacts(constraintFile, factsFile, storage)
	if err != nil {
		t.Fatalf("Failed to build network with facts: %v", err)
	}
	return network, facts, storage
}

// CreateUserFact creates a test user fact (standard format)
func (th *TestHelper) CreateUserFact(id, name, firstName string, age float64) *rete.Fact {
	return &rete.Fact{
		ID:   id,
		Type: "User",
		Fields: map[string]interface{}{
			"id":        id,
			"name":      name,
			"firstName": firstName,
			"age":       age,
		},
	}
}

// CreateAddressFact creates a test address fact (standard format)
func (th *TestHelper) CreateAddressFact(userID, street, city string) *rete.Fact {
	return &rete.Fact{
		ID:   userID + "_address",
		Type: "Address",
		Fields: map[string]interface{}{
			"userID": userID,
			"street": street,
			"city":   city,
		},
	}
}

// CreateCustomerFact creates a test customer fact (standard format)
func (th *TestHelper) CreateCustomerFact(id string, age float64, isVIP bool) *rete.Fact {
	return &rete.Fact{
		ID:   id,
		Type: "Customer",
		Fields: map[string]interface{}{
			"id":    id,
			"age":   age,
			"isVIP": isVIP,
		},
	}
}

// SubmitFactsAndAnalyze submits facts to network and returns action count
func (th *TestHelper) SubmitFactsAndAnalyze(t *testing.T, network *rete.ReteNetwork, facts []*rete.Fact) int {
	totalActions := 0

	// Submit all facts
	for _, fact := range facts {
		err := network.SubmitFact(fact)
		if err != nil {
			t.Logf("Warning: Failed to submit fact %s: %v", fact.ID, err)
		}
	}

	// Analyze terminal nodes for triggered actions
	for terminalID, terminal := range network.TerminalNodes {
		tokenCount := len(terminal.Memory.Tokens)
		totalActions += tokenCount
		t.Logf("Terminal %s: %d tokens stored", terminalID, tokenCount)
	}

	return totalActions
}
