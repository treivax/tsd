// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"testing"
)

// TestDebugBetaJoinComplex trace la propagation des bindings pour le cas qui échoue
func TestDebugBetaJoinComplex(t *testing.T) {
	t.Log("🔍 DEBUG: Traçage de la règle r2 de beta_join_complex.tsd")
	t.Log("Règle: {u: User, o: Order, p: Product} / u.status == \"vip\" AND o.user_id == u.id AND p.id == o.product_id AND p.category == \"luxury\" ==> vip_luxury_purchase(u.id, p.name)")
	t.Log("")

	// Setup
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)

	// Créer les TypeNodes manuellement
	userTypeDef := TypeDefinition{Type: "type", Name: "User", Fields: []Field{
		{Name: "id", Type: "string"},
		{Name: "name", Type: "string"},
		{Name: "age", Type: "number"},
		{Name: "city", Type: "string"},
		{Name: "status", Type: "string"},
	}}
	network.TypeNodes["User"] = NewTypeNode("User", userTypeDef, storage)

	orderTypeDef := TypeDefinition{Type: "type", Name: "Order", Fields: []Field{
		{Name: "id", Type: "string"},
		{Name: "user_id", Type: "string"},
		{Name: "product_id", Type: "string"},
		{Name: "amount", Type: "number"},
		{Name: "date", Type: "string"},
	}}
	network.TypeNodes["Order"] = NewTypeNode("Order", orderTypeDef, storage)

	productTypeDef := TypeDefinition{Type: "type", Name: "Product", Fields: []Field{
		{Name: "id", Type: "string"},
		{Name: "name", Type: "string"},
		{Name: "category", Type: "string"},
		{Name: "price", Type: "number"},
	}}
	network.TypeNodes["Product"] = NewTypeNode("Product", productTypeDef, storage)

	// Connecter les TypeNodes au RootNode pour permettre la propagation des faits
	network.RootNode.AddChild(network.TypeNodes["User"])
	network.RootNode.AddChild(network.TypeNodes["Order"])
	network.RootNode.AddChild(network.TypeNodes["Product"])
	t.Log("✓ TypeNodes connectés au RootNode")

	// Construire la règle r2 manuellement avec debug activé
	ruleID := "r2"
	varTypesMap := map[string]string{
		"u": "User",
		"o": "Order",
		"p": "Product",
	}

	// Condition complète de la règle
	condition := map[string]interface{}{
		"AND": []interface{}{
			map[string]interface{}{"==": []interface{}{"u.status", "vip"}},
			map[string]interface{}{"==": []interface{}{"o.user_id", "u.id"}},
			map[string]interface{}{"==": []interface{}{"p.id", "o.product_id"}},
			map[string]interface{}{"==": []interface{}{"p.category", "luxury"}},
		},
	}

	// Créer les patterns de jointure
	patterns := []JoinPattern{
		{
			LeftVars:    []string{"u"},
			RightVars:   []string{"o"},
			AllVars:     []string{"u", "o"},
			VarTypes:    varTypesMap,
			Condition:   condition,
			Selectivity: 0.5,
		},
		{
			LeftVars:    []string{"u", "o"},
			RightVars:   []string{"p"},
			AllVars:     []string{"u", "o", "p"},
			VarTypes:    varTypesMap,
			Condition:   condition,
			Selectivity: 0.5,
		},
	}

	t.Log("📋 Patterns créés:")
	for i, p := range patterns {
		t.Logf("   Pattern %d: Left=%v, Right=%v, All=%v", i, p.LeftVars, p.RightVars, p.AllVars)
	}
	t.Log("")

	// Créer la chaîne de jointures avec BetaChainBuilder
	chain, err := network.BetaChainBuilder.BuildChain(patterns, ruleID)
	if err != nil {
		t.Fatalf("Erreur création chaîne: %v", err)
	}

	t.Logf("✅ Chaîne créée: %d JoinNodes", len(chain.Nodes))
	for i, node := range chain.Nodes {
		t.Logf("   JoinNode %d: ID=%s, LeftVars=%v, RightVars=%v, AllVars=%v",
			i, node.ID, node.LeftVariables, node.RightVariables, node.AllVariables)
	}
	t.Log("")

	// ACTIVER LE DEBUG sur tous les JoinNodes
	for _, node := range chain.Nodes {
		node.Debug = true
	}
	t.Log("🔍 Mode DEBUG activé sur tous les JoinNodes")
	t.Log("")

	// Créer le TerminalNode (avec une action fictive pour le test)
	action := &Action{
		Type: "action",
		Job: &JobCall{
			Type: "jobCall",
			Name: "vip_luxury_purchase",
			Args: []interface{}{
				map[string]interface{}{"type": "variable", "value": "u.id"},
				map[string]interface{}{"type": "variable", "value": "p.name"},
			},
		},
	}
	terminalNode := NewTerminalNode(ruleID+"_terminal", action, storage)
	network.TerminalNodes[terminalNode.ID] = terminalNode

	// Connecter la chaîne au réseau
	// JoinNode1 (u ⋈ o): User → left, Order → right
	// JoinNode2 (u,o ⋈ p): JoinNode1 → left, Product → right

	joinNode1 := chain.Nodes[0]
	joinNode2 := chain.Nodes[1]

	// Ajouter les JoinNodes au réseau
	network.BetaNodes[joinNode1.ID] = joinNode1
	network.BetaNodes[joinNode2.ID] = joinNode2

	// Utiliser BuilderUtils pour les connexions correctes
	utils := NewBuilderUtils(storage)

	// Connecter TypeNode(User) → JoinNode1.Left (via passthrough AlphaNode)
	utils.ConnectTypeNodeToBetaNode(network, ruleID, "u", "User", joinNode1, NodeSideLeft)
	t.Log("🔗 Connected: TypeNode(User) → PassthroughAlpha → JoinNode1.ActivateLeft")

	// Connecter TypeNode(Order) → JoinNode1.Right (via passthrough AlphaNode)
	utils.ConnectTypeNodeToBetaNode(network, ruleID, "o", "Order", joinNode1, NodeSideRight)
	t.Log("🔗 Connected: TypeNode(Order) → PassthroughAlpha → JoinNode1.ActivateRight")

	// Connecter JoinNode1 → JoinNode2.Left
	joinNode1.AddChild(joinNode2)
	t.Log("🔗 Connected: JoinNode1 → JoinNode2.ActivateLeft")

	// Connecter TypeNode(Product) → JoinNode2.Right (via passthrough AlphaNode)
	utils.ConnectTypeNodeToBetaNode(network, ruleID, "p", "Product", joinNode2, NodeSideRight)
	t.Log("🔗 Connected: TypeNode(Product) → PassthroughAlpha → JoinNode2.ActivateRight")

	// Connecter JoinNode2 → TerminalNode
	joinNode2.AddChild(terminalNode)
	t.Log("🔗 Connected: JoinNode2 → TerminalNode")
	t.Log("")

	// Vérifier les connexions
	t.Log("🔍 Vérification des connexions:")
	t.Logf("   TypeNode(User) a %d enfants", len(network.TypeNodes["User"].Children))
	for i, child := range network.TypeNodes["User"].Children {
		t.Logf("      Child %d: %s (type=%s)", i, child.GetID(), child.GetType())
	}
	t.Logf("   TypeNode(Order) a %d enfants", len(network.TypeNodes["Order"].Children))
	for i, child := range network.TypeNodes["Order"].Children {
		t.Logf("      Child %d: %s (type=%s)", i, child.GetID(), child.GetType())
	}
	t.Logf("   TypeNode(Product) a %d enfants", len(network.TypeNodes["Product"].Children))
	for i, child := range network.TypeNodes["Product"].Children {
		t.Logf("      Child %d: %s (type=%s)", i, child.GetID(), child.GetType())
	}
	t.Log("")

	// Créer les faits (NE PAS les ajouter au storage, SubmitFact le fera)
	t.Log("📥 Création des faits:")
	userFact := &Fact{
		ID:   "USER001",
		Type: "User",
		Fields: map[string]interface{}{
			"id":     "USER001",
			"name":   "Alice",
			"age":    float64(30),
			"city":   "Paris",
			"status": "vip",
		},
	}
	t.Log("   ✓ User(USER001, status=vip)")

	orderFact := &Fact{
		ID:   "ORD001",
		Type: "Order",
		Fields: map[string]interface{}{
			"id":         "ORD001",
			"user_id":    "USER001",
			"product_id": "PROD003",
			"amount":     float64(1),
			"date":       "2025-01-18",
		},
	}
	t.Log("   ✓ Order(ORD001, user_id=USER001, product_id=PROD003)")

	productFact := &Fact{
		ID:   "PROD003",
		Type: "Product",
		Fields: map[string]interface{}{
			"id":       "PROD003",
			"name":     "Luxury Watch",
			"category": "luxury",
			"price":    float64(2500.00),
		},
	}
	t.Log("   ✓ Product(PROD003, category=luxury)")
	t.Log("")

	// Soumettre les faits au réseau
	t.Log("=" + fmt.Sprintf("%80s", "="))
	t.Log("🚀 SOUMISSION User(USER001)")
	t.Log("=" + fmt.Sprintf("%80s", "="))
	if err := network.SubmitFact(userFact); err != nil {
		t.Errorf("Erreur soumission User: %v", err)
	}
	t.Log("")

	t.Log("=" + fmt.Sprintf("%80s", "="))
	t.Log("🚀 SOUMISSION Order(ORD001)")
	t.Log("=" + fmt.Sprintf("%80s", "="))
	if err := network.SubmitFact(orderFact); err != nil {
		t.Errorf("Erreur soumission Order: %v", err)
	}
	t.Log("")

	t.Log("=" + fmt.Sprintf("%80s", "="))
	t.Log("🚀 SOUMISSION Product(PROD003)")
	t.Log("=" + fmt.Sprintf("%80s", "="))
	if err := network.SubmitFact(productFact); err != nil {
		t.Errorf("Erreur soumission Product: %v", err)
	}
	t.Log("")

	// Vérifier l'état des mémoires
	t.Log("=" + fmt.Sprintf("%80s", "="))
	t.Log("📊 ÉTAT DES MÉMOIRES")
	t.Log("=" + fmt.Sprintf("%80s", "="))

	t.Logf("JoinNode1 (%s):", joinNode1.ID)
	t.Logf("   Left Memory:   %d tokens", len(joinNode1.LeftMemory.Tokens))
	for _, token := range joinNode1.LeftMemory.Tokens {
		t.Logf("      Token %s: Bindings=%v", token.ID, token.GetVariables())
	}
	t.Logf("   Right Memory:  %d tokens", len(joinNode1.RightMemory.Tokens))
	for _, token := range joinNode1.RightMemory.Tokens {
		t.Logf("      Token %s: Bindings=%v", token.ID, token.GetVariables())
	}
	t.Logf("   Result Memory: %d tokens", len(joinNode1.ResultMemory.Tokens))
	for _, token := range joinNode1.ResultMemory.Tokens {
		t.Logf("      Token %s: Bindings=%v", token.ID, token.GetVariables())
	}
	t.Log("")

	t.Logf("JoinNode2 (%s):", joinNode2.ID)
	t.Logf("   Left Memory:   %d tokens", len(joinNode2.LeftMemory.Tokens))
	for _, token := range joinNode2.LeftMemory.Tokens {
		t.Logf("      Token %s: Bindings=%v", token.ID, token.GetVariables())
		// Vérifier le contenu du BindingChain
		vars := token.Bindings.Variables()
		t.Logf("         Variables complètes: %v", vars)
		for _, v := range vars {
			fact := token.Bindings.Get(v)
			if fact != nil {
				t.Logf("         %s -> %s (%s)", v, fact.ID, fact.Type)
			}
		}
	}
	t.Logf("   Right Memory:  %d tokens", len(joinNode2.RightMemory.Tokens))
	for _, token := range joinNode2.RightMemory.Tokens {
		t.Logf("      Token %s: Bindings=%v", token.ID, token.GetVariables())
	}
	t.Logf("   Result Memory: %d tokens", len(joinNode2.ResultMemory.Tokens))
	for _, token := range joinNode2.ResultMemory.Tokens {
		t.Logf("      Token %s: Bindings=%v", token.ID, token.GetVariables())
		// Vérifier le contenu du BindingChain
		vars := token.Bindings.Variables()
		t.Logf("         Variables complètes: %v", vars)
		for _, v := range vars {
			fact := token.Bindings.Get(v)
			if fact != nil {
				t.Logf("         %s -> %s (%s)", v, fact.ID, fact.Type)
			}
		}
	}
	t.Log("")

	t.Logf("TerminalNode (%s):", terminalNode.ID)
	t.Logf("   Memory: %d tokens", len(terminalNode.Memory.Tokens))
	for _, token := range terminalNode.Memory.Tokens {
		t.Logf("      Token %s: Bindings=%v", token.ID, token.GetVariables())
		// Vérifier le contenu complet
		vars := token.Bindings.Variables()
		t.Logf("         Variables complètes: %v", vars)
		for _, v := range vars {
			fact := token.Bindings.Get(v)
			if fact != nil {
				t.Logf("         %s -> %s (%s)", v, fact.ID, fact.Type)
			}
		}
	}
	t.Log("")

	// Vérifications
	if len(joinNode1.ResultMemory.Tokens) == 0 {
		t.Error("❌ PROBLÈME: JoinNode1 n'a pas produit de token de jointure")
	} else {
		t.Log("✅ JoinNode1 a produit des tokens")
	}

	if len(joinNode2.LeftMemory.Tokens) == 0 {
		t.Error("❌ PROBLÈME: JoinNode2 n'a pas reçu de token du côté gauche")
	} else {
		t.Log("✅ JoinNode2 a reçu des tokens du côté gauche")
		// Vérifier que le token a TOUS les bindings
		for _, token := range joinNode2.LeftMemory.Tokens {
			vars := token.GetVariables()
			if len(vars) < 2 {
				t.Errorf("❌ PROBLÈME: Token dans JoinNode2.LeftMemory n'a que %d bindings: %v (attendu: [u, o])", len(vars), vars)
			} else {
				hasU := false
				hasO := false
				for _, v := range vars {
					if v == "u" {
						hasU = true
					}
					if v == "o" {
						hasO = true
					}
				}
				if !hasU || !hasO {
					t.Errorf("❌ PROBLÈME: Token manque des bindings. hasU=%v, hasO=%v, vars=%v", hasU, hasO, vars)
				} else {
					t.Logf("✅ Token dans JoinNode2.LeftMemory contient bien [u, o]")
				}
			}
		}
	}

	if len(joinNode2.ResultMemory.Tokens) == 0 {
		t.Error("❌ PROBLÈME: JoinNode2 n'a pas produit de token de jointure")
	} else {
		t.Log("✅ JoinNode2 a produit des tokens")
		// Vérifier que le token final a TOUS les bindings
		for _, token := range joinNode2.ResultMemory.Tokens {
			vars := token.GetVariables()
			if len(vars) < 3 {
				t.Errorf("❌ PROBLÈME: Token final n'a que %d bindings: %v (attendu: [u, o, p])", len(vars), vars)
			} else {
				hasU := false
				hasO := false
				hasP := false
				for _, v := range vars {
					if v == "u" {
						hasU = true
					}
					if v == "o" {
						hasO = true
					}
					if v == "p" {
						hasP = true
					}
				}
				if !hasU || !hasO || !hasP {
					t.Errorf("❌ PROBLÈME: Token final manque des bindings. hasU=%v, hasO=%v, hasP=%v, vars=%v", hasU, hasO, hasP, vars)
				} else {
					t.Logf("✅ Token final contient bien [u, o, p]")
				}
			}
		}
	}

	if len(terminalNode.Memory.Tokens) == 0 {
		t.Error("❌ PROBLÈME: TerminalNode n'a pas reçu de token")
	} else {
		t.Log("✅ TerminalNode a reçu des tokens")
	}
}
