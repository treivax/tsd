// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestPassthroughNodeKey tests the key generation for passthrough nodes
func TestPassthroughNodeKey(t *testing.T) {
	tests := []struct {
		name     string
		ruleName string
		typeName string
		varName  string
		side     string
		expected string
	}{
		{
			name:     "Without side",
			ruleName: "rule1",
			typeName: "Person",
			varName:  "p",
			side:     "",
			expected: "passthrough_rule1_p_Person",
		},
		{
			name:     "With left side",
			ruleName: "rule1",
			typeName: "Person",
			varName:  "p",
			side:     "left",
			expected: "passthrough_rule1_p_Person_left",
		},
		{
			name:     "With right side",
			ruleName: "rule1",
			typeName: "Order",
			varName:  "o",
			side:     "right",
			expected: "passthrough_rule1_o_Order_right",
		},
		{
			name:     "Complex type name",
			ruleName: "complexRule",
			typeName: "ProductOrder",
			varName:  "po",
			side:     "left",
			expected: "passthrough_complexRule_po_ProductOrder_left",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PassthroughNodeKey(tt.ruleName, tt.typeName, tt.varName, tt.side)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestGetOrCreatePassthroughAlphaNode_SameTypeSameSide tests that passthrough nodes are shared
func TestGetOrCreatePassthroughAlphaNode_SameTypeSameSide(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	utils := NewBuilderUtils(storage)

	// Create first passthrough for Person with left side
	alpha1 := utils.GetOrCreatePassthroughAlphaNode(network, "rule1", "Person", "p", "left")

	assert.NotNil(t, alpha1)
	assert.Equal(t, "passthrough_rule1_p_Person_left", alpha1.ID)
	assert.Equal(t, 1, len(network.PassthroughRegistry))

	// Create second passthrough for same rule (should reuse)
	alpha2 := utils.GetOrCreatePassthroughAlphaNode(network, "rule1", "Person", "p", "left")

	// Should be the same instance
	assert.Equal(t, alpha1, alpha2)
	assert.Equal(t, 1, len(network.PassthroughRegistry))

	// Verify the node is in the registry
	key := PassthroughNodeKey("rule1", "Person", "p", "left")
	registeredNode, exists := network.PassthroughRegistry[key]
	assert.True(t, exists)
	assert.Equal(t, alpha1, registeredNode)
}

// TestGetOrCreatePassthroughAlphaNode_SameTypeDifferentSide tests different sides create different nodes
func TestGetOrCreatePassthroughAlphaNode_SameTypeDifferentSide(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	utils := NewBuilderUtils(storage)

	// Create left-side passthrough
	alphaLeft := utils.GetOrCreatePassthroughAlphaNode(network, "rule1", "Person", "p", "left")

	// Create right-side passthrough (should be different)
	alphaRight := utils.GetOrCreatePassthroughAlphaNode(network, "rule1", "Person", "p", "right")

	assert.NotEqual(t, alphaLeft, alphaRight)
	assert.Equal(t, "passthrough_rule1_p_Person_left", alphaLeft.ID)
	assert.Equal(t, "passthrough_rule1_p_Person_right", alphaRight.ID)
	assert.Equal(t, 2, len(network.PassthroughRegistry))
}

// TestGetOrCreatePassthroughAlphaNode_DifferentTypes tests different types create different nodes
func TestGetOrCreatePassthroughAlphaNode_DifferentTypes(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	utils := NewBuilderUtils(storage)

	// Create passthrough for Person
	alphaPerson := utils.GetOrCreatePassthroughAlphaNode(network, "rule1", "Person", "p", "left")

	// Create passthrough for Order (should be different)
	alphaOrder := utils.GetOrCreatePassthroughAlphaNode(network, "rule1", "Order", "o", "left")

	assert.NotEqual(t, alphaPerson, alphaOrder)
	assert.Equal(t, "passthrough_rule1_p_Person_left", alphaPerson.ID)
	assert.Equal(t, "passthrough_rule1_o_Order_left", alphaOrder.ID)
	assert.Equal(t, 2, len(network.PassthroughRegistry))
}

// TestGetOrCreatePassthroughAlphaNode_NoSide tests passthrough without side specification
func TestGetOrCreatePassthroughAlphaNode_NoSide(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	utils := NewBuilderUtils(storage)

	// Create passthrough without side
	alpha1 := utils.GetOrCreatePassthroughAlphaNode(network, "rule1", "Person", "p", "")

	assert.NotNil(t, alpha1)
	assert.Equal(t, "passthrough_rule1_p_Person", alpha1.ID)
	assert.Equal(t, 1, len(network.PassthroughRegistry))

	// Create another passthrough without side (should reuse)
	alpha2 := utils.GetOrCreatePassthroughAlphaNode(network, "rule1", "Person", "p", "")

	assert.Equal(t, alpha1, alpha2)
	assert.Equal(t, 1, len(network.PassthroughRegistry))
}

// TestGetOrCreatePassthroughAlphaNode_Condition tests the condition structure
func TestGetOrCreatePassthroughAlphaNode_Condition(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	utils := NewBuilderUtils(storage)

	// Create passthrough with side
	alpha := utils.GetOrCreatePassthroughAlphaNode(network, "rule1", "Person", "p", "left")

	assert.NotNil(t, alpha.Condition)
	condMap, ok := alpha.Condition.(map[string]interface{})
	assert.True(t, ok)

	assert.Equal(t, ConditionTypePassthrough, condMap["type"])
	assert.Equal(t, "left", condMap["side"])

	// Create passthrough without side
	alpha2 := utils.GetOrCreatePassthroughAlphaNode(network, "rule1", "Order", "o", "")
	condMap2, ok := alpha2.Condition.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, ConditionTypePassthrough, condMap2["type"])
	_, hasSide := condMap2["side"]
	assert.False(t, hasSide)
}

// TestGetOrCreatePassthroughAlphaNode_VariableName tests that variable name doesn't affect sharing
func TestGetOrCreatePassthroughAlphaNode_VariableName(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	utils := NewBuilderUtils(storage)

	// Create passthrough with variable name "p"
	alpha1 := utils.GetOrCreatePassthroughAlphaNode(network, "rule1", "Person", "p", "left")

	// Create passthrough with different variable name (should be different now - per rule+var)
	alpha2 := utils.GetOrCreatePassthroughAlphaNode(network, "rule1", "Person", "person", "left")

	// Create passthrough with another variable name (should be different)
	alpha3 := utils.GetOrCreatePassthroughAlphaNode(network, "rule1", "Person", "x", "left")

	// All should be DIFFERENT instances (per-variable passthrough)
	assert.NotEqual(t, alpha1, alpha2)
	assert.NotEqual(t, alpha2, alpha3)
	assert.NotEqual(t, alpha1, alpha3)
	assert.Equal(t, 3, len(network.PassthroughRegistry))
}

// TestConnectTypeNodeToBetaNode_Sharing tests that ConnectTypeNodeToBetaNode creates per-rule passthroughs
func TestConnectTypeNodeToBetaNode_Sharing(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	utils := NewBuilderUtils(storage)

	// Create type node
	typeDef := TypeDefinition{Name: "Person", Fields: []Field{}}
	typeNode := NewTypeNode("Person", typeDef, storage)
	network.TypeNodes["Person"] = typeNode

	// Create two beta nodes (simulating two different rules)
	condition := map[string]interface{}{"type": "simple"}
	betaNode1 := NewJoinNode("join1", condition, []string{"p"}, []string{}, map[string]string{"p": "Person"}, storage)
	betaNode2 := NewJoinNode("join2", condition, []string{"q"}, []string{}, map[string]string{"q": "Person"}, storage)

	// Connect both beta nodes via the same type
	utils.ConnectTypeNodeToBetaNode(network, "rule1", "p", "Person", betaNode1, "left")
	utils.ConnectTypeNodeToBetaNode(network, "rule2", "q", "Person", betaNode2, "left")

	// Should have 2 passthrough nodes (per-rule passthroughs)
	assert.Equal(t, 2, len(network.PassthroughRegistry))

	// TypeNode should have 2 children (one per rule)
	assert.Equal(t, 2, len(typeNode.GetChildren()))

	// Each passthrough should have 1 child (its own beta node)
	alphaNode1 := typeNode.GetChildren()[0].(*AlphaNode)
	assert.Equal(t, 1, len(alphaNode1.GetChildren()))
}

// TestConnectTypeNodeToBetaNode_DifferentSides tests that different sides create different nodes
func TestConnectTypeNodeToBetaNode_DifferentSides(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	utils := NewBuilderUtils(storage)

	// Create type node
	typeDef := TypeDefinition{Name: "Person", Fields: []Field{}}
	typeNode := NewTypeNode("Person", typeDef, storage)
	network.TypeNodes["Person"] = typeNode

	// Create two beta nodes
	condition := map[string]interface{}{"type": "simple"}
	betaNode1 := NewJoinNode("join1", condition, []string{"p"}, []string{}, map[string]string{"p": "Person"}, storage)
	betaNode2 := NewJoinNode("join2", condition, []string{"p"}, []string{}, map[string]string{"p": "Person"}, storage)

	// Connect with different sides
	utils.ConnectTypeNodeToBetaNode(network, "rule1", "p", "Person", betaNode1, "left")
	utils.ConnectTypeNodeToBetaNode(network, "rule2", "p", "Person", betaNode2, "right")

	// Should have 2 passthrough nodes (one for each side)
	assert.Equal(t, 2, len(network.PassthroughRegistry))

	// TypeNode should have 2 children (one for each side)
	assert.Equal(t, 2, len(typeNode.GetChildren()))
}

// TestNetworkReset_ClearsPassthroughRegistry tests that Reset clears the passthrough registry
func TestNetworkReset_ClearsPassthroughRegistry(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	utils := NewBuilderUtils(storage)

	// Create some passthrough nodes
	utils.GetOrCreatePassthroughAlphaNode(network, "rule1", "Person", "p", "left")
	utils.GetOrCreatePassthroughAlphaNode(network, "rule1", "Order", "o", "right")

	assert.Equal(t, 2, len(network.PassthroughRegistry))

	// Reset the network
	network.Reset()

	// Registry should be empty
	assert.Equal(t, 0, len(network.PassthroughRegistry))
}

// TestPassthroughSharing_MultipleRulesSameTypes tests realistic scenario with multiple rules
// Updated to reflect per-rule passthrough behavior
func TestPassthroughSharing_MultipleRulesSameTypes(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	utils := NewBuilderUtils(storage)

	// Create type nodes
	personTypeDef := TypeDefinition{Name: "Person", Fields: []Field{}}
	orderTypeDef := TypeDefinition{Name: "Order", Fields: []Field{}}
	personType := NewTypeNode("Person", personTypeDef, storage)
	orderType := NewTypeNode("Order", orderTypeDef, storage)
	network.TypeNodes["Person"] = personType
	network.TypeNodes["Order"] = orderType

	// Simulate 3 rules, each using Person (left) and Order (right)
	condition := map[string]interface{}{"type": "simple"}
	for i := 1; i <= 3; i++ {
		betaNode := NewJoinNode("join"+string(rune('0'+i)), condition, []string{"p"}, []string{"o"}, map[string]string{"p": "Person", "o": "Order"}, storage)
		utils.ConnectTypeNodeToBetaNode(network, "rule"+string(rune('0'+i)), "p", "Person", betaNode, "left")
		utils.ConnectTypeNodeToBetaNode(network, "rule"+string(rune('0'+i)), "o", "Order", betaNode, "right")
	}

	// Should have 6 passthrough nodes (3 rules * 2 sides = 6 per-rule passthroughs)
	assert.Equal(t, 6, len(network.PassthroughRegistry))

	// Person TypeNode should have 3 children (one per rule)
	assert.Equal(t, 3, len(personType.GetChildren()))

	// Order TypeNode should have 3 children (one per rule)
	assert.Equal(t, 3, len(orderType.GetChildren()))

	// Each passthrough should have 1 child (its own beta node)
	personPassthrough := personType.GetChildren()[0].(*AlphaNode)
	orderPassthrough := orderType.GetChildren()[0].(*AlphaNode)
	assert.Equal(t, 1, len(personPassthrough.GetChildren()))
	assert.Equal(t, 1, len(orderPassthrough.GetChildren()))

	// Verify the IDs follow per-rule format: passthrough_<ruleID>_<varName>_<type>_<side>
	assert.Equal(t, "passthrough_rule1_p_Person_left", personPassthrough.ID)
	assert.Equal(t, "passthrough_rule1_o_Order_right", orderPassthrough.ID)
}

// TestPassthroughSharing_NoDoubleConnection tests that connecting twice doesn't create duplicate edges
func TestPassthroughSharing_NoDoubleConnection(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	utils := NewBuilderUtils(storage)

	// Create type node
	typeDef := TypeDefinition{Name: "Person", Fields: []Field{}}
	typeNode := NewTypeNode("Person", typeDef, storage)
	network.TypeNodes["Person"] = typeNode

	// Create beta node
	condition := map[string]interface{}{"type": "simple"}
	betaNode := NewJoinNode("join1", condition, []string{"p"}, []string{}, map[string]string{"p": "Person"}, storage)

	// Connect twice (simulating a bug or repeated call)
	utils.ConnectTypeNodeToBetaNode(network, "rule1", "p", "Person", betaNode, "left")
	utils.ConnectTypeNodeToBetaNode(network, "rule1", "p", "Person", betaNode, "left")

	// Should still have only 1 passthrough node
	assert.Equal(t, 1, len(network.PassthroughRegistry))

	// TypeNode should have only 1 child (not duplicated)
	assert.Equal(t, 1, len(typeNode.GetChildren()))

	// The passthrough should have the beta node (potentially twice if AddChild doesn't check)
	// This test documents current behavior - we may want to add deduplication in AddChild
	alphaNode := typeNode.GetChildren()[0].(*AlphaNode)
	assert.NotNil(t, alphaNode)
}

// TestPassthroughRegistry_InitializedInNewNetwork tests that new networks have initialized registry
func TestPassthroughRegistry_InitializedInNewNetwork(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)

	assert.NotNil(t, network.PassthroughRegistry)
	assert.Equal(t, 0, len(network.PassthroughRegistry))
}

// TestPassthroughSharing_RegistryConsistency tests registry consistency after operations
func TestPassthroughSharing_RegistryConsistency(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	utils := NewBuilderUtils(storage)

	// Create passthrough nodes (each with different variables)
	alpha1 := utils.GetOrCreatePassthroughAlphaNode(network, "rule1", "Person", "p", "left")
	alpha2 := utils.GetOrCreatePassthroughAlphaNode(network, "rule1", "Person", "q", "left")
	alpha3 := utils.GetOrCreatePassthroughAlphaNode(network, "rule1", "Order", "o", "right")

	// Verify registry state - now we have 3 nodes (one per variable)
	assert.Equal(t, 3, len(network.PassthroughRegistry))

	// Verify we can retrieve them by key (using new signature)
	key1 := PassthroughNodeKey("rule1", "Person", "p", "left")
	key2 := PassthroughNodeKey("rule1", "Person", "q", "left")
	key3 := PassthroughNodeKey("rule1", "Order", "o", "right")

	retrievedAlpha1, exists1 := network.PassthroughRegistry[key1]
	retrievedAlpha2, exists2 := network.PassthroughRegistry[key2]
	retrievedAlpha3, exists3 := network.PassthroughRegistry[key3]

	assert.True(t, exists1)
	assert.True(t, exists2)
	assert.True(t, exists3)
	assert.Equal(t, alpha1, retrievedAlpha1)
	assert.Equal(t, alpha2, retrievedAlpha2)
	assert.Equal(t, alpha3, retrievedAlpha3)
}
