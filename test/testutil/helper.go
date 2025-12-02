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
// This function uses the incremental IngestFile method
func (th *TestHelper) BuildNetworkFromConstraintFile(t *testing.T, constraintFile string) (*rete.ReteNetwork, rete.Storage) {
	storage := rete.NewMemoryStorage()
	network, err := th.pipeline.IngestFile(constraintFile, nil, storage)
	if err != nil {
		t.Fatalf("Failed to build network: %v", err)
	}
	return network, storage
}

// BuildNetworkFromConstraintFileWithMetrics builds a RETE network and returns metrics
func (th *TestHelper) BuildNetworkFromConstraintFileWithMetrics(t *testing.T, constraintFile string) (*rete.ReteNetwork, rete.Storage, *rete.IngestionMetrics) {
	storage := rete.NewMemoryStorage()
	network, metrics, err := th.pipeline.IngestFileWithMetrics(constraintFile, nil, storage)
	if err != nil {
		t.Fatalf("Failed to build network: %v", err)
	}
	return network, storage, metrics
}

// IngestFile incrementally adds types, rules, and facts from a file to an existing network
// If network is nil, creates a new network
func (th *TestHelper) IngestFile(t *testing.T, filename string, network *rete.ReteNetwork, storage rete.Storage) *rete.ReteNetwork {
	var err error
	network, err = th.pipeline.IngestFile(filename, network, storage)
	if err != nil {
		t.Fatalf("Failed to ingest file %s: %v", filename, err)
	}
	return network
}

// IngestFileWithMetrics incrementally adds types, rules, and facts with metrics collection
func (th *TestHelper) IngestFileWithMetrics(t *testing.T, filename string, network *rete.ReteNetwork, storage rete.Storage) (*rete.ReteNetwork, *rete.IngestionMetrics) {
	var err error
	var metrics *rete.IngestionMetrics
	network, metrics, err = th.pipeline.IngestFileWithMetrics(filename, network, storage)
	if err != nil {
		t.Fatalf("Failed to ingest file %s: %v", filename, err)
	}
	return network, metrics
}

// BuildNetworkFromConstraintFileWithFacts builds a RETE network and loads facts
// This uses IngestFile incrementally for both constraint and facts files
func (th *TestHelper) BuildNetworkFromConstraintFileWithFacts(t *testing.T, constraintFile, factsFile string) (*rete.ReteNetwork, []*rete.Fact, rete.Storage) {
	storage := rete.NewMemoryStorage()

	// If the same file is used for both constraints and facts, ingest it only once
	if constraintFile == factsFile {
		network, err := th.pipeline.IngestFile(constraintFile, nil, storage)
		if err != nil {
			t.Fatalf("Failed to ingest file: %v", err)
		}

		// Collect submitted facts from the network
		facts := th.collectAllFactsFromNetwork(network)
		return network, facts, storage
	}

	// First ingest the constraint file (types and rules)
	network, err := th.pipeline.IngestFile(constraintFile, nil, storage)
	if err != nil {
		t.Fatalf("Failed to ingest constraint file: %v", err)
	}

	// Then ingest the facts file
	network, err = th.pipeline.IngestFile(factsFile, network, storage)
	if err != nil {
		t.Fatalf("Failed to ingest facts file: %v", err)
	}

	// Collect submitted facts from the network
	facts := th.collectAllFactsFromNetwork(network)
	return network, facts, storage
}

// collectAllFactsFromNetwork collects all unique facts from the network
func (th *TestHelper) collectAllFactsFromNetwork(network *rete.ReteNetwork) []*rete.Fact {
	factMap := make(map[string]*rete.Fact)

	// Collect from TypeNodes
	for _, typeNode := range network.TypeNodes {
		for _, token := range typeNode.Memory.Tokens {
			for _, fact := range token.Facts {
				if fact != nil {
					factMap[fact.ID] = fact
				}
			}
		}
	}

	// Collect from AlphaNodes
	for _, alphaNode := range network.AlphaNodes {
		for _, token := range alphaNode.Memory.Tokens {
			for _, fact := range token.Facts {
				if fact != nil {
					factMap[fact.ID] = fact
				}
			}
		}
	}

	// Collect from TerminalNodes
	for _, terminalNode := range network.TerminalNodes {
		for _, token := range terminalNode.Memory.Tokens {
			for _, fact := range token.Facts {
				if fact != nil {
					factMap[fact.ID] = fact
				}
			}
		}
	}

	// Convert map to slice
	facts := make([]*rete.Fact, 0, len(factMap))
	for _, fact := range factMap {
		facts = append(facts, fact)
	}

	return facts
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
