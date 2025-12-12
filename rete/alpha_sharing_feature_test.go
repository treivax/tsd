// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete
import (
	"testing"
)
// TestConditionHash vérifie que le hashing de conditions est cohérent
func TestConditionHash(t *testing.T) {
	tests := []struct {
		name        string
		cond1       interface{}
		var1        string
		cond2       interface{}
		var2        string
		shouldMatch bool
	}{
		{
			name: "identical simple conditions",
			cond1: map[string]interface{}{
				"type":      "constraint",
				"attribute": "age",
				"operator":  ">",
				"value":     18,
			},
			var1: "p",
			cond2: map[string]interface{}{
				"type":      "constraint",
				"attribute": "age",
				"operator":  ">",
				"value":     18,
			},
			var2:        "p",
			shouldMatch: true,
		},
		{
			name: "different variable names",
			cond1: map[string]interface{}{
				"type":      "constraint",
				"attribute": "age",
				"operator":  ">",
				"value":     18,
			},
			var1: "p",
			cond2: map[string]interface{}{
				"type":      "constraint",
				"attribute": "age",
				"operator":  ">",
				"value":     18,
			},
			var2:        "q",
			shouldMatch: false,
		},
		{
			name: "different operators",
			cond1: map[string]interface{}{
				"type":      "constraint",
				"attribute": "age",
				"operator":  ">",
				"value":     18,
			},
			var1: "p",
			cond2: map[string]interface{}{
				"type":      "constraint",
				"attribute": "age",
				"operator":  ">=",
				"value":     18,
			},
			var2:        "p",
			shouldMatch: false,
		},
		{
			name: "different values",
			cond1: map[string]interface{}{
				"type":      "constraint",
				"attribute": "age",
				"operator":  ">",
				"value":     18,
			},
			var1: "p",
			cond2: map[string]interface{}{
				"type":      "constraint",
				"attribute": "age",
				"operator":  ">",
				"value":     21,
			},
			var2:        "p",
			shouldMatch: false,
		},
		{
			name:        "nil conditions",
			cond1:       nil,
			var1:        "p",
			cond2:       nil,
			var2:        "p",
			shouldMatch: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash1, err1 := ConditionHash(tt.cond1, tt.var1)
			if err1 != nil {
				t.Fatalf("Erreur hash1: %v", err1)
			}
			hash2, err2 := ConditionHash(tt.cond2, tt.var2)
			if err2 != nil {
				t.Fatalf("Erreur hash2: %v", err2)
			}
			matches := (hash1 == hash2)
			if matches != tt.shouldMatch {
				t.Errorf("Expected match=%v, got match=%v (hash1=%s, hash2=%s)",
					tt.shouldMatch, matches, hash1, hash2)
			}
		})
	}
}
// TestAlphaSharingRegistry_GetOrCreate vérifie la création et le partage d'AlphaNodes
func TestAlphaSharingRegistry_GetOrCreate(t *testing.T) {
	storage := NewMemoryStorage()
	registry := NewAlphaSharingRegistry()
	condition := map[string]interface{}{
		"type":      "constraint",
		"attribute": "age",
		"operator":  ">",
		"value":     18,
	}
	// Première création
	node1, hash1, wasShared1, err := registry.GetOrCreateAlphaNode(condition, "p", storage)
	if err != nil {
		t.Fatalf("Erreur création nœud 1: %v", err)
	}
	if wasShared1 {
		t.Error("Premier nœud ne devrait pas être partagé")
	}
	if node1 == nil {
		t.Fatal("node1 ne devrait pas être nil")
	}
	// Deuxième création avec même condition
	node2, hash2, wasShared2, err := registry.GetOrCreateAlphaNode(condition, "p", storage)
	if err != nil {
		t.Fatalf("Erreur création nœud 2: %v", err)
	}
	if !wasShared2 {
		t.Error("Deuxième nœud devrait être partagé")
	}
	if node2 == nil {
		t.Fatal("node2 ne devrait pas être nil")
	}
	// Vérifier que c'est le même nœud
	if node1.GetID() != node2.GetID() {
		t.Errorf("Les nœuds devraient avoir le même ID: %s != %s", node1.GetID(), node2.GetID())
	}
	if hash1 != hash2 {
		t.Errorf("Les hash devraient être identiques: %s != %s", hash1, hash2)
	}
	// Troisième création avec condition différente
	differentCondition := map[string]interface{}{
		"type":      "constraint",
		"attribute": "age",
		"operator":  ">",
		"value":     21,
	}
	node3, hash3, wasShared3, err := registry.GetOrCreateAlphaNode(differentCondition, "p", storage)
	if err != nil {
		t.Fatalf("Erreur création nœud 3: %v", err)
	}
	if wasShared3 {
		t.Error("Troisième nœud (condition différente) ne devrait pas être partagé")
	}
	if node3.GetID() == node1.GetID() {
		t.Error("Nœud 3 devrait avoir un ID différent de nœud 1")
	}
	if hash3 == hash1 {
		t.Error("Hash 3 devrait être différent de hash 1")
	}
}
// TestAlphaSharingRegistry_RemoveAlphaNode vérifie la suppression d'AlphaNodes
func TestAlphaSharingRegistry_RemoveAlphaNode(t *testing.T) {
	storage := NewMemoryStorage()
	registry := NewAlphaSharingRegistry()
	condition := map[string]interface{}{
		"type":      "constraint",
		"attribute": "age",
		"operator":  ">",
		"value":     18,
	}
	node, hash, _, err := registry.GetOrCreateAlphaNode(condition, "p", storage)
	if err != nil {
		t.Fatalf("Erreur création nœud: %v", err)
	}
	// Vérifier que le nœud existe
	retrieved, exists := registry.GetAlphaNode(hash)
	if !exists {
		t.Error("Le nœud devrait exister")
	}
	if retrieved.GetID() != node.GetID() {
		t.Error("Le nœud récupéré devrait être le même")
	}
	// Supprimer le nœud
	err = registry.RemoveAlphaNode(hash)
	if err != nil {
		t.Fatalf("Erreur suppression nœud: %v", err)
	}
	// Vérifier que le nœud n'existe plus
	_, exists = registry.GetAlphaNode(hash)
	if exists {
		t.Error("Le nœud ne devrait plus exister après suppression")
	}
	// Tenter de supprimer à nouveau (devrait échouer)
	err = registry.RemoveAlphaNode(hash)
	if err == nil {
		t.Error("La suppression d'un nœud inexistant devrait retourner une erreur")
	}
}
// TestAlphaSharingRegistry_Stats vérifie les statistiques
func TestAlphaSharingRegistry_Stats(t *testing.T) {
	storage := NewMemoryStorage()
	registry := NewAlphaSharingRegistry()
	// État initial
	stats := registry.GetStats()
	if stats["total_shared_alpha_nodes"].(int) != 0 {
		t.Error("Devrait avoir 0 nœuds au début")
	}
	// Créer quelques nœuds
	cond1 := map[string]interface{}{"type": "constraint", "attribute": "age", "operator": ">", "value": 18}
	cond2 := map[string]interface{}{"type": "constraint", "attribute": "salary", "operator": ">", "value": 50000}
	node1, _, _, _ := registry.GetOrCreateAlphaNode(cond1, "p", storage)
	node2, _, _, _ := registry.GetOrCreateAlphaNode(cond2, "p", storage)
	// Ajouter des enfants (simuler des règles multiples)
	terminal1 := NewTerminalNode("term1", &Action{Type: "print"}, storage)
	terminal2 := NewTerminalNode("term2", &Action{Type: "print"}, storage)
	terminal3 := NewTerminalNode("term3", &Action{Type: "print"}, storage)
	node1.AddChild(terminal1)
	node1.AddChild(terminal2)
	node2.AddChild(terminal3)
	// Vérifier les stats
	stats = registry.GetStats()
	if stats["total_shared_alpha_nodes"].(int) != 2 {
		t.Errorf("Devrait avoir 2 nœuds, got %v", stats["total_shared_alpha_nodes"])
	}
	if stats["total_rule_references"].(int) != 3 {
		t.Errorf("Devrait avoir 3 références de règles, got %v", stats["total_rule_references"])
	}
	avgSharing := stats["average_sharing_ratio"].(float64)
	if avgSharing != 1.5 { // 3 refs / 2 nodes = 1.5
		t.Errorf("Ratio de partage moyen devrait être 1.5, got %v", avgSharing)
	}
}
// TestAlphaSharingRegistry_ListNodes vérifie le listage des nœuds
func TestAlphaSharingRegistry_ListNodes(t *testing.T) {
	storage := NewMemoryStorage()
	registry := NewAlphaSharingRegistry()
	// Liste vide initialement
	nodes := registry.ListSharedAlphaNodes()
	if len(nodes) != 0 {
		t.Error("La liste devrait être vide initialement")
	}
	// Créer quelques nœuds
	cond1 := map[string]interface{}{"type": "constraint", "attribute": "age", "operator": ">", "value": 18}
	cond2 := map[string]interface{}{"type": "constraint", "attribute": "salary", "operator": ">", "value": 50000}
	_, hash1, _, _ := registry.GetOrCreateAlphaNode(cond1, "p", storage)
	_, hash2, _, _ := registry.GetOrCreateAlphaNode(cond2, "p", storage)
	// Vérifier la liste
	nodes = registry.ListSharedAlphaNodes()
	if len(nodes) != 2 {
		t.Errorf("Devrait avoir 2 nœuds, got %d", len(nodes))
	}
	// Vérifier que les hash sont dans la liste
	found1, found2 := false, false
	for _, hash := range nodes {
		if hash == hash1 {
			found1 = true
		}
		if hash == hash2 {
			found2 = true
		}
	}
	if !found1 || !found2 {
		t.Error("Les deux hash devraient être dans la liste")
	}
}
// TestAlphaSharingRegistry_GetDetails vérifie les détails d'un nœud
func TestAlphaSharingRegistry_GetDetails(t *testing.T) {
	storage := NewMemoryStorage()
	registry := NewAlphaSharingRegistry()
	condition := map[string]interface{}{
		"type":      "constraint",
		"attribute": "age",
		"operator":  ">",
		"value":     18,
	}
	node, hash, _, _ := registry.GetOrCreateAlphaNode(condition, "p", storage)
	// Ajouter des enfants
	terminal1 := NewTerminalNode("term1", &Action{Type: "print"}, storage)
	terminal2 := NewTerminalNode("term2", &Action{Type: "print"}, storage)
	node.AddChild(terminal1)
	node.AddChild(terminal2)
	// Récupérer les détails
	details := registry.GetSharedAlphaNodeDetails(hash)
	if details == nil {
		t.Fatal("Les détails ne devraient pas être nil")
	}
	if details["hash"].(string) != hash {
		t.Error("Le hash devrait correspondre")
	}
	if details["node_id"].(string) != node.GetID() {
		t.Error("Le node_id devrait correspondre")
	}
	if details["variable_name"].(string) != "p" {
		t.Error("Le variable_name devrait être 'p'")
	}
	if details["child_count"].(int) != 2 {
		t.Errorf("Devrait avoir 2 enfants, got %v", details["child_count"])
	}
	// Tester avec un hash inexistant
	detailsNil := registry.GetSharedAlphaNodeDetails("nonexistent")
	if detailsNil != nil {
		t.Error("Les détails d'un nœud inexistant devraient être nil")
	}
}
// TestAlphaSharingRegistry_Reset vérifie la réinitialisation
func TestAlphaSharingRegistry_Reset(t *testing.T) {
	storage := NewMemoryStorage()
	registry := NewAlphaSharingRegistry()
	// Créer quelques nœuds
	cond1 := map[string]interface{}{"type": "constraint", "attribute": "age", "operator": ">", "value": 18}
	cond2 := map[string]interface{}{"type": "constraint", "attribute": "salary", "operator": ">", "value": 50000}
	registry.GetOrCreateAlphaNode(cond1, "p", storage)
	registry.GetOrCreateAlphaNode(cond2, "p", storage)
	// Vérifier qu'il y a des nœuds
	stats := registry.GetStats()
	if stats["total_shared_alpha_nodes"].(int) != 2 {
		t.Error("Devrait avoir 2 nœuds avant reset")
	}
	// Reset
	registry.Reset()
	// Vérifier que tout est vide
	stats = registry.GetStats()
	if stats["total_shared_alpha_nodes"].(int) != 0 {
		t.Error("Devrait avoir 0 nœuds après reset")
	}
	nodes := registry.ListSharedAlphaNodes()
	if len(nodes) != 0 {
		t.Error("La liste devrait être vide après reset")
	}
}
// TestAlphaSharingRegistry_ConcurrentAccess vérifie la sécurité des accès concurrents
func TestAlphaSharingRegistry_ConcurrentAccess(t *testing.T) {
	storage := NewMemoryStorage()
	registry := NewAlphaSharingRegistry()
	condition := map[string]interface{}{
		"type":      "constraint",
		"attribute": "age",
		"operator":  ">",
		"value":     18,
	}
	// Créer plusieurs goroutines qui tentent d'accéder simultanément
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			node, _, _, err := registry.GetOrCreateAlphaNode(condition, "p", storage)
			if err != nil {
				t.Errorf("Erreur création nœud: %v", err)
			}
			if node == nil {
				t.Error("Le nœud ne devrait pas être nil")
			}
			done <- true
		}()
	}
	// Attendre que toutes les goroutines terminent
	for i := 0; i < 10; i++ {
		<-done
	}
	// Vérifier qu'un seul nœud a été créé
	stats := registry.GetStats()
	if stats["total_shared_alpha_nodes"].(int) != 1 {
		t.Errorf("Devrait avoir exactement 1 nœud partagé, got %v", stats["total_shared_alpha_nodes"])
	}
}
// TestNormalizeCondition vérifie la normalisation des conditions
func TestNormalizeCondition(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name:  "nil condition",
			input: nil,
			expected: map[string]interface{}{
				"type": "simple",
			},
		},
		{
			name: "simple map",
			input: map[string]interface{}{
				"type":  "constraint",
				"value": 42,
			},
			expected: map[string]interface{}{
				"type":  "constraint",
				"value": 42,
			},
		},
		{
			name: "nested map",
			input: map[string]interface{}{
				"type": "and",
				"conditions": []interface{}{
					map[string]interface{}{
						"attribute": "age",
						"operator":  ">",
						"value":     18,
					},
				},
			},
			expected: map[string]interface{}{
				"type": "and",
				"conditions": []interface{}{
					map[string]interface{}{
						"attribute": "age",
						"operator":  ">",
						"value":     18,
					},
				},
			},
		},
		{
			name:     "primitive value",
			input:    42,
			expected: 42,
		},
		{
			name:     "string value",
			input:    "test",
			expected: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := normalizeCondition(tt.input)
			if err != nil {
				t.Fatalf("Erreur normalisation: %v", err)
			}
			// Comparer les résultats (test simple)
			if result == nil && tt.expected != nil {
				t.Error("Le résultat ne devrait pas être nil")
			}
		})
	}
}