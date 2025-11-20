package validation

import (
	"testing"
)

func TestRETENewBasic(t *testing.T) {
	// Créer un nouveau réseau
	network := NewRETEValidationNetwork()

	// Ajouter une règle simple avec un seul type
	simpleRule := RETERule{
		ID: "simple_rule",
		VariableTypes: map[string]string{
			"u": "User",
		},
		Conditions: []string{"true"},
	}
	network.AddRule(simpleRule)

	// Ajouter un fait
	fact := &RETEFact{
		ID:   "user1",
		Type: "User",
		Fields: map[string]interface{}{
			"id":   "1",
			"name": "Alice",
		},
	}
	network.InsertFact(fact)

	// Vérifier les tokens terminaux
	terminals := network.GetTerminalTokens()
	if len(terminals) != 1 {
		t.Errorf("Attendu 1 token terminal, reçu %d", len(terminals))
	}
	if len(terminals) > 0 && len(terminals[0].Facts) != 1 {
		t.Errorf("Attendu 1 fait dans le token, reçu %d", len(terminals[0].Facts))
	}
}

func TestRETENewJointure(t *testing.T) {
	// Créer un nouveau réseau
	network := NewRETEValidationNetwork()

	// Ajouter une règle avec jointure (2 types)
	joinRule := RETERule{
		ID: "join_rule",
		VariableTypes: map[string]string{
			"u": "User",
			"o": "Order",
		},
		Conditions: []string{"u.id == o.user_id"},
	}
	network.AddRule(joinRule)

	// Ajouter des faits
	userFact := &RETEFact{
		ID:   "user1",
		Type: "User",
		Fields: map[string]interface{}{
			"id":   "1",
			"name": "Alice",
		},
	}
	orderFact := &RETEFact{
		ID:   "order1",
		Type: "Order",
		Fields: map[string]interface{}{
			"id":      "100",
			"user_id": "1",
		},
	}

	network.InsertFact(userFact)
	network.InsertFact(orderFact)

	// Debug
	network.Debug()

	// Vérifier les tokens terminaux
	terminals := network.GetTerminalTokens()
	if len(terminals) != 1 {
		t.Errorf("Attendu 1 token terminal, reçu %d", len(terminals))
	}
	if len(terminals) > 0 && len(terminals[0].Facts) != 2 {
		t.Errorf("Attendu 2 faits dans le token, reçu %d", len(terminals[0].Facts))
	}
}

func TestRETEIncrementalPropagation(t *testing.T) {
	// Test spécifique de la propagation incrémentale
	network := NewRETEValidationNetwork()

	// Règle avec 3 niveaux et conditions étagées
	rule := RETERule{
		ID: "incremental_test",
		VariableTypes: map[string]string{
			"u": "User",
			"o": "Order",
			"p": "Product",
		},
		Conditions: []string{
			"u.id == o.user_id",    // Niveau 1: u+o
			"o.product_id == p.id", // Niveau 2: résultat niveau 1 + p
			"u.age >= 18",          // Condition finale
		},
	}
	network.AddRule(rule)

	// Ajouter faits un par un pour observer la propagation
	userFact := &RETEFact{
		ID:   "user1",
		Type: "User",
		Fields: map[string]interface{}{
			"id":  "U1",
			"age": 25,
		},
	}

	// 1. Ajouter User - doit créer token alpha
	network.InsertFact(userFact)

	// Vérifier alpha nodes
	alphaTokenCount := 0
	for _, node := range network.AlphaNodes {
		if node.TypeName == "User" {
			alphaTokenCount = len(node.Tokens)
		}
	}
	if alphaTokenCount != 1 {
		t.Errorf("Attendu 1 token alpha User, reçu %d", alphaTokenCount)
	}

	// 2. Ajouter Order - doit déclencher jointure si condition respectée
	orderFact := &RETEFact{
		ID:   "order1",
		Type: "Order",
		Fields: map[string]interface{}{
			"id":         "O1",
			"user_id":    "U1", // Match avec user
			"product_id": "P1",
		},
	}
	network.InsertFact(orderFact)

	// Vérifier qu'il y a un token dans le beta niveau 1
	beta1Count := 0
	for id, node := range network.BetaNodes {
		if id == "beta_incremental_test_1" {
			beta1Count = len(node.Tokens)
		}
	}
	if beta1Count != 1 {
		t.Errorf("Attendu 1 token dans beta niveau 1, reçu %d", beta1Count)
	}

	// 3. Ajouter Product - doit compléter la chaîne
	productFact := &RETEFact{
		ID:   "product1",
		Type: "Product",
		Fields: map[string]interface{}{
			"id":   "P1", // Match avec order
			"name": "Test",
		},
	}
	network.InsertFact(productFact)

	// Vérifier les tokens terminaux
	terminals := network.GetTerminalTokens()
	if len(terminals) != 1 {
		t.Errorf("Attendu 1 token terminal après propagation complète, reçu %d", len(terminals))
	}

	// 4. Test avec fait qui ne match pas - doit être filtré par conditions
	badOrderFact := &RETEFact{
		ID:   "order2",
		Type: "Order",
		Fields: map[string]interface{}{
			"id":         "O2",
			"user_id":    "U999", // Ne match pas avec user
			"product_id": "P1",
		},
	}
	network.InsertFact(badOrderFact)

	// Le nombre de terminaux ne doit pas changer
	terminalsAfterBad := network.GetTerminalTokens()
	if len(terminalsAfterBad) != 1 {
		t.Errorf("Les mauvais faits ne doivent pas créer de tokens terminaux, reçu %d", len(terminalsAfterBad))
	}

	t.Logf("✅ Propagation incrémentale validée avec filtrage par conditions")
}
