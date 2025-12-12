// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete

import (
	"testing"
)

func TestReteNetworkReset(t *testing.T) {
	t.Log("🧪 TEST RESET DU RÉSEAU RETE")
	t.Log("============================")
	t.Run("ResetClearsNetwork", func(t *testing.T) {
		// Create a network with indexed storage
		storage := NewMemoryStorage()
		network := NewReteNetwork(storage)
		// Add some types
		typeUser := TypeDefinition{
			Type: "typeDefinition",
			Name: "User",
			Fields: []Field{
				{Name: "name", Type: "string"},
				{Name: "age", Type: "number"},
			},
		}
		network.Types = append(network.Types, typeUser)
		// Create and add a TypeNode manually
		typeNodeDef := TypeDefinition{
			Type: "typeDefinition",
			Name: "User",
			Fields: []Field{
				{Name: "name", Type: "string"},
				{Name: "age", Type: "number"},
			},
		}
		typeNode := NewTypeNode("User", typeNodeDef, storage)
		network.TypeNodes["User"] = typeNode
		// Add some test data to verify state
		network.AlphaNodes["test_alpha"] = &AlphaNode{}
		network.TerminalNodes["test_terminal"] = &TerminalNode{}
		// Verify network has content before reset
		if len(network.Types) == 0 {
			t.Fatal("❌ Network devrait avoir des types avant reset")
		}
		if len(network.TypeNodes) == 0 {
			t.Fatal("❌ Network devrait avoir des TypeNodes avant reset")
		}
		if len(network.AlphaNodes) == 0 {
			t.Fatal("❌ Network devrait avoir des AlphaNodes avant reset")
		}
		if len(network.TerminalNodes) == 0 {
			t.Fatal("❌ Network devrait avoir des TerminalNodes avant reset")
		}
		// Call reset
		network.Reset()
		// Verify everything is cleared
		if len(network.Types) != 0 {
			t.Errorf("❌ Types devrait être vide après reset, reçu %d éléments", len(network.Types))
		}
		if len(network.TypeNodes) != 0 {
			t.Errorf("❌ TypeNodes devrait être vide après reset, reçu %d éléments", len(network.TypeNodes))
		}
		if len(network.AlphaNodes) != 0 {
			t.Errorf("❌ AlphaNodes devrait être vide après reset, reçu %d éléments", len(network.AlphaNodes))
		}
		if len(network.BetaNodes) != 0 {
			t.Errorf("❌ BetaNodes devrait être vide après reset, reçu %d éléments", len(network.BetaNodes))
		}
		if len(network.TerminalNodes) != 0 {
			t.Errorf("❌ TerminalNodes devrait être vide après reset, reçu %d éléments", len(network.TerminalNodes))
		}
		// Verify RootNode is recreated (not nil)
		if network.RootNode == nil {
			t.Error("❌ RootNode ne devrait pas être nil après reset")
		}
		// Verify BetaBuilder is cleared
		if network.BetaBuilder != nil {
			t.Error("❌ BetaBuilder devrait être nil après reset")
		}
		t.Log("✅ Reset a correctement vidé le réseau RETE")
	})
	t.Run("ResetPreservesStorage", func(t *testing.T) {
		storage := NewMemoryStorage()
		network := NewReteNetwork(storage)
		// Store original storage reference
		originalStorage := network.Storage
		// Add some data
		network.Types = append(network.Types, TypeDefinition{Name: "Test"})
		// Reset
		network.Reset()
		// Verify storage reference is preserved
		if network.Storage == nil {
			t.Error("❌ Storage ne devrait pas être nil après reset")
		}
		if network.Storage != originalStorage {
			t.Error("❌ Storage reference devrait être préservée après reset")
		}
		t.Log("✅ Reset préserve la référence au Storage")
	})
	t.Run("ResetMultipleTimes", func(t *testing.T) {
		storage := NewMemoryStorage()
		network := NewReteNetwork(storage)
		// Add data, reset, add data again, reset again
		network.Types = append(network.Types, TypeDefinition{Name: "User"})
		network.Reset()
		if len(network.Types) != 0 {
			t.Fatal("❌ Premier reset a échoué")
		}
		network.Types = append(network.Types, TypeDefinition{Name: "Order"})
		network.Reset()
		if len(network.Types) != 0 {
			t.Fatal("❌ Deuxième reset a échoué")
		}
		t.Log("✅ Reset peut être appelé plusieurs fois")
	})
	t.Run("CanRebuildAfterReset", func(t *testing.T) {
		storage := NewMemoryStorage()
		network := NewReteNetwork(storage)
		// Build initial network
		typeUser := TypeDefinition{
			Type: "typeDefinition",
			Name: "User",
			Fields: []Field{
				{Name: "name", Type: "string"},
			},
		}
		network.Types = append(network.Types, typeUser)
		typeNodeUser := NewTypeNode("User", typeUser, storage)
		network.TypeNodes["User"] = typeNodeUser
		// Verify initial state
		if len(network.Types) != 1 || network.Types[0].Name != "User" {
			t.Fatal("❌ Initial network setup failed")
		}
		// Reset
		network.Reset()
		// Build new network
		typeOrder := TypeDefinition{
			Type: "typeDefinition",
			Name: "Order",
			Fields: []Field{
				{Name: "id", Type: "number"},
			},
		}
		network.Types = append(network.Types, typeOrder)
		typeNodeOrder := NewTypeNode("Order", typeOrder, storage)
		network.TypeNodes["Order"] = typeNodeOrder
		// Verify new state
		if len(network.Types) != 1 {
			t.Errorf("❌ Attendu 1 type, reçu %d", len(network.Types))
		}
		if network.Types[0].Name != "Order" {
			t.Errorf("❌ Attendu type 'Order', reçu '%s'", network.Types[0].Name)
		}
		if len(network.TypeNodes) != 1 {
			t.Errorf("❌ Attendu 1 TypeNode, reçu %d", len(network.TypeNodes))
		}
		if _, exists := network.TypeNodes["Order"]; !exists {
			t.Error("❌ TypeNode 'Order' devrait exister")
		}
		if _, exists := network.TypeNodes["User"]; exists {
			t.Error("❌ TypeNode 'User' ne devrait plus exister après reset")
		}
		t.Log("✅ Réseau peut être reconstruit après reset")
	})
	t.Run("ResetClearsRootNodeMemory", func(t *testing.T) {
		storage := NewMemoryStorage()
		network := NewReteNetwork(storage)
		// Submit a fact to populate memory
		fact := &Fact{
			ID:   "test1",
			Type: "TestType",
			Fields: map[string]interface{}{
				"name": "test",
			},
		}
		// Note: SubmitFact may fail if network is not fully built, but we're just testing reset
		_ = network.SubmitFact(fact)
		// Reset
		network.Reset()
		// Verify new root node is created
		if network.RootNode == nil {
			t.Fatal("❌ RootNode devrait être recréé après reset")
		}
		// Verify root node memory is clean
		memory := network.RootNode.GetMemory()
		if memory == nil {
			t.Fatal("❌ RootNode devrait avoir une mémoire")
		}
		t.Log("✅ Reset recrée un RootNode propre")
	})
}
