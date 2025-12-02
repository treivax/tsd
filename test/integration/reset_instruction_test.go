// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package main

import (
	"testing"
)

// TestResetInstruction_BasicReset tests that reset clears all types, rules and facts
// and allows defining new ones after the reset
func TestResetInstruction_BasicReset(t *testing.T) {
	t.Skip("Reset au milieu de fichier n'a plus de sens dans le mode incrémental - reset devrait être en début de fichier uniquement")

	helper := NewTestHelper()

	constraintFile := "../../constraint/test/integration/reset_integration_test.tsd"
	factsFile := "../../constraint/test/integration/reset_integration_test.tsd"

	// Build network from constraint and facts
	network, facts, storage := helper.BuildNetworkFromConstraintFileWithMassiveFacts(t, constraintFile, factsFile)

	// Validate basic structure
	if len(facts) == 0 {
		t.Fatal("No facts were parsed")
	}

	if len(network.TerminalNodes) == 0 {
		t.Fatal("No terminal nodes created")
	}

	// After reset, only types defined AFTER the reset should exist
	// Types before reset: User, Order, Product
	// Types after reset: Customer, Invoice, Payment
	expectedTypesAfterReset := []string{"Customer", "Invoice", "Payment"}
	unexpectedTypesBeforeReset := []string{"User", "Order", "Product"}

	// Verify that types from BEFORE reset do NOT exist
	for _, unexpectedType := range unexpectedTypesBeforeReset {
		if _, found := network.TypeNodes[unexpectedType]; found {
			t.Errorf("❌ Type '%s' should have been removed by reset but still exists", unexpectedType)
		}
	}

	// Verify that types from AFTER reset DO exist
	for _, expectedType := range expectedTypesAfterReset {
		if _, found := network.TypeNodes[expectedType]; !found {
			t.Errorf("❌ Type '%s' should exist after reset but was not found", expectedType)
		}
	}

	// Verify storage is not nil
	if storage == nil {
		t.Error("Storage should not be nil after reset")
	}

	// Verify facts were created for the types AFTER reset
	customerFactsCount := 0
	invoiceFactsCount := 0
	paymentFactsCount := 0

	for _, fact := range facts {
		switch fact.Type {
		case "Customer":
			customerFactsCount++
		case "Invoice":
			invoiceFactsCount++
		case "Payment":
			paymentFactsCount++
		case "User", "Order", "Product":
			t.Errorf("❌ Fact of type '%s' should not exist after reset", fact.Type)
		}
	}

	if customerFactsCount == 0 {
		t.Error("Expected at least one Customer fact after reset")
	}
	if invoiceFactsCount == 0 {
		t.Error("Expected at least one Invoice fact after reset")
	}
	if paymentFactsCount == 0 {
		t.Error("Expected at least one Payment fact after reset")
	}

	t.Logf("✅ Basic reset test passed: %d type nodes, %d terminal nodes, %d facts (%d Customer, %d Invoice, %d Payment)",
		len(network.TypeNodes), len(network.TerminalNodes), len(facts),
		customerFactsCount, invoiceFactsCount, paymentFactsCount)
}

// TestResetInstruction_MultipleResets tests that multiple successive resets work correctly
// and that only the last set of definitions remains
func TestResetInstruction_MultipleResets(t *testing.T) {
	t.Skip("Reset au milieu de fichier n'a plus de sens dans le mode incrémental - reset devrait être en début de fichier uniquement")

	helper := NewTestHelper()

	constraintFile := "../../constraint/test/integration/multiple_resets_test.tsd"
	factsFile := "../../constraint/test/integration/multiple_resets_test.tsd"

	// Build network from constraint and facts
	network, facts, storage := helper.BuildNetworkFromConstraintFileWithMassiveFacts(t, constraintFile, factsFile)

	// Validate basic structure
	if len(facts) == 0 {
		t.Fatal("No facts were parsed")
	}

	if len(network.TerminalNodes) == 0 {
		t.Fatal("No terminal nodes created")
	}

	// After two resets, only types defined AFTER the SECOND reset should exist
	// Types after Phase 1 (before first reset): UserV1, PostV1
	// Types after Phase 2 (before second reset): ProductV2, SupplierV2
	// Types after Phase 3 (after second reset): Employee, Project, Assignment
	expectedTypesAfterResets := []string{"Employee", "Project", "Assignment"}
	unexpectedTypesPhase1 := []string{"UserV1", "PostV1"}
	unexpectedTypesPhase2 := []string{"ProductV2", "SupplierV2"}

	// Verify that types from Phase 1 do NOT exist
	for _, unexpectedType := range unexpectedTypesPhase1 {
		if _, found := network.TypeNodes[unexpectedType]; found {
			t.Errorf("❌ Type '%s' from Phase 1 should have been removed by first reset but still exists", unexpectedType)
		}
	}

	// Verify that types from Phase 2 do NOT exist
	for _, unexpectedType := range unexpectedTypesPhase2 {
		if _, found := network.TypeNodes[unexpectedType]; found {
			t.Errorf("❌ Type '%s' from Phase 2 should have been removed by second reset but still exists", unexpectedType)
		}
	}

	// Verify that types from Phase 3 (after both resets) DO exist
	for _, expectedType := range expectedTypesAfterResets {
		if _, found := network.TypeNodes[expectedType]; !found {
			t.Errorf("❌ Type '%s' should exist after resets but was not found", expectedType)
		}
	}

	// Verify storage is not nil
	if storage == nil {
		t.Error("Storage should not be nil after multiple resets")
	}

	// Verify facts were created for the types AFTER the second reset
	employeeFactsCount := 0
	projectFactsCount := 0
	assignmentFactsCount := 0

	for _, fact := range facts {
		switch fact.Type {
		case "Employee":
			employeeFactsCount++
		case "Project":
			projectFactsCount++
		case "Assignment":
			assignmentFactsCount++
		case "UserV1", "PostV1", "ProductV2", "SupplierV2":
			t.Errorf("❌ Fact of type '%s' should not exist after resets", fact.Type)
		}
	}

	if employeeFactsCount == 0 {
		t.Error("Expected at least one Employee fact after resets")
	}
	if projectFactsCount == 0 {
		t.Error("Expected at least one Project fact after resets")
	}
	if assignmentFactsCount == 0 {
		t.Error("Expected at least one Assignment fact after resets")
	}

	t.Logf("✅ Multiple resets test passed: %d type nodes, %d terminal nodes, %d facts (%d Employee, %d Project, %d Assignment)",
		len(network.TypeNodes), len(network.TerminalNodes), len(facts),
		employeeFactsCount, projectFactsCount, assignmentFactsCount)
}

// TestResetInstruction_NetworkIntegrity validates that after a reset,
// the RETE network is properly reconstructed with only post-reset definitions
func TestResetInstruction_NetworkIntegrity(t *testing.T) {
	t.Skip("Reset au milieu de fichier n'a plus de sens dans le mode incrémental - reset devrait être en début de fichier uniquement")

	helper := NewTestHelper()

	constraintFile := "../../constraint/test/integration/reset_integration_test.tsd"
	factsFile := "../../constraint/test/integration/reset_integration_test.tsd"

	network, _, storage := helper.BuildNetworkFromConstraintFileWithMassiveFacts(t, constraintFile, factsFile)

	// Verify network structure is valid after reset
	if network.RootNode == nil {
		t.Fatal("RootNode should not be nil after reset")
	}

	if len(network.TypeNodes) == 0 {
		t.Error("TypeNodes should not be empty after reset with new types")
	}

	if len(network.TerminalNodes) == 0 {
		t.Error("TerminalNodes should not be empty after reset with new rules")
	}

	// Verify storage reference is preserved
	if storage == nil {
		t.Error("Storage should not be nil")
	}

	// Verify that the network can process facts (rules are active)
	activatedRules := 0
	for _, terminal := range network.TerminalNodes {
		if len(terminal.Memory.Tokens) > 0 {
			activatedRules++
		}
	}

	if activatedRules == 0 {
		t.Error("Expected at least one rule to be activated after reset")
	}

	// Verify only types AFTER reset are present (Customer, Invoice, Payment)
	expectedTypeCount := 3
	if len(network.TypeNodes) != expectedTypeCount {
		t.Errorf("Expected %d type nodes after reset, got %d", expectedTypeCount, len(network.TypeNodes))
	}

	t.Logf("✅ Network integrity test passed: %d type nodes, %d alpha nodes, %d beta nodes, %d terminal nodes, %d rules activated",
		len(network.TypeNodes), len(network.AlphaNodes), len(network.BetaNodes),
		len(network.TerminalNodes), activatedRules)
}

// TestResetInstruction_RulesAfterReset verifies that rules defined AFTER reset
// are correctly executed and rules from BEFORE reset are no longer active
func TestResetInstruction_RulesAfterReset(t *testing.T) {
	t.Skip("Reset au milieu de fichier n'a plus de sens dans le mode incrémental - reset devrait être en début de fichier uniquement")

	helper := NewTestHelper()

	constraintFile := "../../constraint/test/integration/reset_integration_test.tsd"
	factsFile := "../../constraint/test/integration/reset_integration_test.tsd"

	network, facts, _ := helper.BuildNetworkFromConstraintFileWithMassiveFacts(t, constraintFile, factsFile)

	// Count how many terminal nodes (rules) exist
	ruleCount := len(network.TerminalNodes)

	// Expected rules AFTER reset (from reset_integration_test.constraint):
	// 1. Customer(?c) where ?c.vip == true => log("Client VIP: " + ?c.name)
	// 2. Invoice(?i) where ?i.total > 5000 => log("Facture élevée: " + ?i.id)
	// 3. Payment(?p) & Invoice(?i) where ?p.invoiceId == ?i.id => log("Paiement pour facture: " + ?i.id)
	expectedRuleCount := 3

	if ruleCount != expectedRuleCount {
		t.Errorf("❌ Expected %d rules after reset, but found %d", expectedRuleCount, ruleCount)
	}

	// Verify that at least some rules have been activated
	totalActivations := 0
	for _, terminal := range network.TerminalNodes {
		activations := len(terminal.Memory.Tokens)
		totalActivations += activations

		if terminal.Action != nil {
			t.Logf("Rule activated %d times: %v", activations, terminal.Action.Job.Name)
		}
	}

	if totalActivations == 0 {
		t.Error("Expected at least some rule activations after reset")
	}

	// Verify VIP customer rule is activated for VIP customers
	vipCustomerCount := 0
	for _, fact := range facts {
		if fact.Type == "Customer" {
			if vip, ok := fact.Fields["vip"].(bool); ok && vip {
				vipCustomerCount++
			}
		}
	}

	if vipCustomerCount == 0 {
		t.Error("Test setup error: expected at least one VIP customer in facts")
	}

	t.Logf("✅ Rules after reset test passed: %d rules defined, %d total activations, %d VIP customers",
		ruleCount, totalActivations, vipCustomerCount)
}

// TestResetInstruction_StoragePreservation verifies that the Storage reference
// is preserved across resets (important for performance)
func TestResetInstruction_StoragePreservation(t *testing.T) {
	helper := NewTestHelper()

	constraintFile := "../../constraint/test/integration/reset_integration_test.tsd"
	factsFile := "../../constraint/test/integration/reset_integration_test.tsd"

	network, _, storage := helper.BuildNetworkFromConstraintFileWithMassiveFacts(t, constraintFile, factsFile)

	// Verify storage is the same instance (preserved)
	if storage == nil {
		t.Fatal("Storage should not be nil")
	}

	// Verify network has a valid storage reference
	if network.Storage == nil {
		t.Error("Network storage should not be nil after reset")
	}

	// Both should point to the same storage
	if network.Storage != storage {
		t.Error("Network storage and returned storage should be the same instance")
	}

	t.Logf("✅ Storage preservation test passed: storage correctly preserved across reset")
}

// TestResetInstruction_ParsingOnly tests that files containing reset instruction
// can be parsed and that only post-reset definitions are in the network
func TestResetInstruction_ParsingOnly(t *testing.T) {
	t.Skip("Reset au milieu de fichier n'a plus de sens dans le mode incrémental - reset devrait être en début de fichier uniquement")

	helper := NewTestHelper()

	// Test that a file with reset is parsed with correct reset semantics
	constraintFile := "../../constraint/test/integration/reset_integration_test.tsd"
	network, _ := helper.BuildNetworkFromConstraintFile(t, constraintFile)

	// Verify network was built with only post-reset types
	expectedTypeCount := 3 // Customer, Invoice, Payment (not User, Order, Product)
	if len(network.TypeNodes) != expectedTypeCount {
		t.Errorf("Expected %d type nodes after reset, got %d", expectedTypeCount, len(network.TypeNodes))
	}

	if network.RootNode == nil {
		t.Error("Expected RootNode to exist")
	}

	t.Logf("✅ File with reset instruction parsed successfully with correct reset semantics")
	t.Logf("   Network contains %d type nodes (only types after reset, as expected)", len(network.TypeNodes))
}
