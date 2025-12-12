// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"testing"
)

// TestBetaChainBuilder_BuildCascade3Variables tests the complete cascade building for 3 variables
func TestBetaChainBuilder_BuildCascade3Variables(t *testing.T) {
	t.Log("🧪 TEST BetaChainBuilder - Cascade 3 variables")
	t.Log("============================================")

	// Setup
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)

	// Create type definitions
	userType := TypeDefinition{Name: "User", Type: "type", Fields: []Field{{Name: "id", Type: "string"}}}
	orderType := TypeDefinition{Name: "Order", Type: "type", Fields: []Field{{Name: "id", Type: "string"}, {Name: "user_id", Type: "string"}}}
	productType := TypeDefinition{Name: "Product", Type: "type", Fields: []Field{{Name: "id", Type: "string"}, {Name: "order_id", Type: "string"}}}

	network.Types = []TypeDefinition{userType, orderType, productType}

	// Create TypeNodes
	network.TypeNodes = make(map[string]*TypeNode)
	network.TypeNodes["User"] = NewTypeNode("User", userType, storage)
	network.TypeNodes["Order"] = NewTypeNode("Order", orderType, storage)
	network.TypeNodes["Product"] = NewTypeNode("Product", productType, storage)

	// Connect to RootNode
	network.RootNode.AddChild(network.TypeNodes["User"])
	network.RootNode.AddChild(network.TypeNodes["Order"])
	network.RootNode.AddChild(network.TypeNodes["Product"])

	// Build patterns for 3 variables: u (User), o (Order), p (Product)
	variableNames := []string{"u", "o", "p"}
	variableTypes := []string{"User", "Order", "Product"}
	varTypes := BuildVarTypesMap(variableNames, variableTypes)

	patterns := []JoinPattern{
		// Pattern 1: u ⋈ o
		{
			LeftVars:    []string{"u"},
			RightVars:   []string{"o"},
			AllVars:     []string{"u", "o"},
			VarTypes:    varTypes,
			Condition:   nil,
			Selectivity: 0.5,
		},
		// Pattern 2: (u, o) ⋈ p
		{
			LeftVars:    []string{"u", "o"},
			RightVars:   []string{"p"},
			AllVars:     []string{"u", "o", "p"},
			VarTypes:    varTypes,
			Condition:   nil,
			Selectivity: 0.5,
		},
	}

	// Build the chain
	chain, err := network.BetaChainBuilder.BuildChain(patterns, "test_rule_3vars")
	if err != nil {
		t.Fatalf("❌ Erreur construction: %v", err)
	}

	// Verify chain structure
	if len(chain.Nodes) != 2 {
		t.Fatalf("❌ Attendu 2 JoinNodes, got %d", len(chain.Nodes))
	}

	// Verify JoinNode 1: u ⋈ o
	join1 := chain.Nodes[0]
	t.Logf("📍 JoinNode 1 (ID: %s)", join1.ID)
	t.Logf("   LeftVariables:  %v", join1.LeftVariables)
	t.Logf("   RightVariables: %v", join1.RightVariables)
	t.Logf("   AllVariables:   %v", join1.AllVariables)

	if !slicesEqual(join1.LeftVariables, []string{"u"}) {
		t.Errorf("❌ JoinNode1 LeftVariables incorrect: %v (attendu [u])", join1.LeftVariables)
	}
	if !slicesEqual(join1.RightVariables, []string{"o"}) {
		t.Errorf("❌ JoinNode1 RightVariables incorrect: %v (attendu [o])", join1.RightVariables)
	}
	if !slicesEqual(join1.AllVariables, []string{"u", "o"}) {
		t.Errorf("❌ JoinNode1 AllVariables incorrect: %v (attendu [u, o])", join1.AllVariables)
	}
	if len(join1.VariableTypes) != 3 {
		t.Errorf("❌ JoinNode1 VariableTypes devrait contenir 3 types, got %d: %v", len(join1.VariableTypes), join1.VariableTypes)
	}

	// Verify JoinNode 2: (u, o) ⋈ p
	join2 := chain.Nodes[1]
	t.Logf("📍 JoinNode 2 (ID: %s)", join2.ID)
	t.Logf("   LeftVariables:  %v", join2.LeftVariables)
	t.Logf("   RightVariables: %v", join2.RightVariables)
	t.Logf("   AllVariables:   %v", join2.AllVariables)

	if !slicesEqual(join2.LeftVariables, []string{"u", "o"}) {
		t.Errorf("❌ JoinNode2 LeftVariables incorrect: %v (attendu [u, o])", join2.LeftVariables)
	}
	if !slicesEqual(join2.RightVariables, []string{"p"}) {
		t.Errorf("❌ JoinNode2 RightVariables incorrect: %v (attendu [p])", join2.RightVariables)
	}
	if !slicesEqual(join2.AllVariables, []string{"u", "o", "p"}) {
		t.Errorf("❌ JoinNode2 AllVariables incorrect: %v (attendu [u, o, p])", join2.AllVariables)
	}
	if len(join2.VariableTypes) != 3 {
		t.Errorf("❌ JoinNode2 VariableTypes devrait contenir 3 types, got %d: %v", len(join2.VariableTypes), join2.VariableTypes)
	}

	// Verify VariableTypes content
	expectedTypes := map[string]string{"u": "User", "o": "Order", "p": "Product"}
	for varName, expectedType := range expectedTypes {
		if actualType, exists := join1.VariableTypes[varName]; !exists {
			t.Errorf("❌ JoinNode1 VariableTypes manque '%s'", varName)
		} else if actualType != expectedType {
			t.Errorf("❌ JoinNode1 VariableTypes['%s'] = '%s', attendu '%s'", varName, actualType, expectedType)
		}

		if actualType, exists := join2.VariableTypes[varName]; !exists {
			t.Errorf("❌ JoinNode2 VariableTypes manque '%s'", varName)
		} else if actualType != expectedType {
			t.Errorf("❌ JoinNode2 VariableTypes['%s'] = '%s', attendu '%s'", varName, actualType, expectedType)
		}
	}

	t.Log("✅ Cascade 3 variables construite correctement")
}

// TestBetaChainBuilder_BuildCascadeNVariables tests cascade building for N variables
func TestBetaChainBuilder_BuildCascadeNVariables(t *testing.T) {
	t.Log("🧪 TEST BetaChainBuilder - Cascade N variables")
	t.Log("============================================")

	for n := 2; n <= 5; n++ {
		t.Run(fmt.Sprintf("n=%d", n), func(t *testing.T) {
			// Setup
			storage := NewMemoryStorage()
			network := NewReteNetwork(storage)

			// Create variable names and types
			varNames := make([]string, n)
			varTypes := make([]string, n)
			typeMap := make(map[string]string)
			for i := 0; i < n; i++ {
				varNames[i] = fmt.Sprintf("v%d", i)
				varTypes[i] = fmt.Sprintf("Type%d", i)
				typeMap[varNames[i]] = varTypes[i]
			}

			// Build patterns
			patterns := make([]JoinPattern, n-1)
			for i := 0; i < n-1; i++ {
				leftVars := make([]string, i+1)
				copy(leftVars, varNames[0:i+1])

				allVars := make([]string, i+2)
				copy(allVars, varNames[0:i+2])

				patterns[i] = JoinPattern{
					LeftVars:    leftVars,
					RightVars:   []string{varNames[i+1]},
					AllVars:     allVars,
					VarTypes:    typeMap,
					Condition:   nil,
					Selectivity: 0.5,
				}
			}

			// Build chain
			chain, err := network.BetaChainBuilder.BuildChain(patterns, fmt.Sprintf("test_rule_%dvars", n))
			if err != nil {
				t.Fatalf("❌ Erreur construction pour %d variables: %v", n, err)
			}

			// Verify chain length
			if len(chain.Nodes) != n-1 {
				t.Errorf("❌ Pour %d variables, attendu %d JoinNodes, got %d", n, n-1, len(chain.Nodes))
			}

			// Verify each JoinNode
			for i, joinNode := range chain.Nodes {
				expectedLeftVars := i + 1
				expectedAllVars := i + 2

				if len(joinNode.LeftVariables) != expectedLeftVars {
					t.Errorf("❌ JoinNode %d: LeftVariables attendu %d, got %d: %v",
						i, expectedLeftVars, len(joinNode.LeftVariables), joinNode.LeftVariables)
				}

				if len(joinNode.RightVariables) != 1 {
					t.Errorf("❌ JoinNode %d: RightVariables attendu 1, got %d: %v",
						i, len(joinNode.RightVariables), joinNode.RightVariables)
				}

				if len(joinNode.AllVariables) != expectedAllVars {
					t.Errorf("❌ JoinNode %d: AllVariables attendu %d, got %d: %v",
						i, expectedAllVars, len(joinNode.AllVariables), joinNode.AllVariables)
				}

				if len(joinNode.VariableTypes) != n {
					t.Errorf("❌ JoinNode %d: VariableTypes devrait contenir %d types, got %d",
						i, n, len(joinNode.VariableTypes))
				}
			}

			t.Logf("✅ Cascade %d variables construite correctement", n)
		})
	}
}
