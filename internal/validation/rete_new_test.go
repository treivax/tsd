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