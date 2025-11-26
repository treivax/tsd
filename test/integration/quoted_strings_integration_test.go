// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package main

import (
	"testing"
)

// TestQuotedStringsIntegration tests the full pipeline with quoted strings
func TestQuotedStringsIntegration(t *testing.T) {
	// Use helper like other integration tests
	helper := NewTestHelper()

	constraintFile := "../../constraint/test/integration/quoted_strings_integration.tsd"
	factsFile := "../../constraint/test/integration/quoted_strings_integration.tsd"

	network, facts, _ := helper.BuildNetworkFromConstraintFileWithMassiveFacts(t, constraintFile, factsFile)

	if network == nil {
		t.Fatal("Network should not be nil")
	}

	if facts == nil {
		t.Fatal("Facts should not be nil")
	}

	// Verify types were created
	if len(network.TypeNodes) != 2 {
		t.Errorf("Expected 2 type nodes, got %d", len(network.TypeNodes))
	}

	if _, exists := network.TypeNodes["Person"]; !exists {
		t.Error("Expected Person type node to exist")
	}

	if _, exists := network.TypeNodes["Message"]; !exists {
		t.Error("Expected Message type node to exist")
	}

	// Count facts by type
	personFacts := 0
	messageFacts := 0

	for _, fact := range facts {
		switch fact.Type {
		case "Person":
			personFacts++
			// Verify field values with quotes were parsed correctly
			if fact.ID == "p1" {
				if name, ok := fact.Fields["name"].(string); ok {
					if name != "Alice Smith" {
						t.Errorf("Expected name 'Alice Smith', got '%s'", name)
					}
				} else {
					t.Error("Name field should be a string")
				}
				if city, ok := fact.Fields["city"].(string); ok {
					if city != "New York" {
						t.Errorf("Expected city 'New York', got '%s'", city)
					}
				}
			}
			if fact.ID == "p2" {
				if name, ok := fact.Fields["name"].(string); ok {
					if name != "Bob Jones" {
						t.Errorf("Expected name 'Bob Jones', got '%s'", name)
					}
				}
			}
		case "Message":
			messageFacts++
			if fact.ID == "m1" {
				if text, ok := fact.Fields["text"].(string); ok {
					if text != "Hello, World!" {
						t.Errorf("Expected text 'Hello, World!', got '%s'", text)
					}
				}
			}
		}
	}

	if personFacts != 4 {
		t.Errorf("Expected 4 Person facts, got %d", personFacts)
	}

	if messageFacts != 4 {
		t.Errorf("Expected 4 Message facts, got %d", messageFacts)
	}

	// Count terminal nodes (rules)
	terminalCount := len(network.TerminalNodes)
	expectedRules := 4
	if terminalCount != expectedRules {
		t.Errorf("Expected %d terminal nodes (rules), got %d", expectedRules, terminalCount)
	}

	// Count activations for each terminal
	totalActivations := 0
	for _, terminal := range network.TerminalNodes {
		tokenCount := len(terminal.Memory.Tokens)
		totalActivations += tokenCount
		t.Logf("Terminal: %d activation(s)", tokenCount)
	}

	// We expect at least some activations
	if totalActivations == 0 {
		t.Error("Expected at least one activation from the rules")
	}

	t.Logf("✓ Quoted strings integration test passed with %d total activations", totalActivations)
}

// TestQuotedStringsEscapeSequences tests various escape sequences in strings
func TestQuotedStringsEscapeSequences(t *testing.T) {
	// This test is covered by the unit tests in constraint package
	// and the integration test above
	t.Skip("Escape sequences are tested in constraint/quoted_strings_test.go")

}
