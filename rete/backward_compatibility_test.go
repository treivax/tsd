// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
package rete

import (
	"os"
	"path/filepath"
	"testing"
)

// TestBackwardCompatibility_SimpleRules vérifie que les règles simples
// fonctionnent toujours comme avant l'introduction des AlphaChains
func TestBackwardCompatibility_SimpleRules(t *testing.T) {
	tempDir := t.TempDir()
	tsdFile := filepath.Join(tempDir, "simple_rules.tsd")
	// Règles simples basiques qui doivent toujours fonctionner
	content := `type Person(#id: string, age: number, name:string)
action print(message: string)
action print(message: string)
rule adult : {p: Person} / p.age >= 18 ==> print("Adult detected")
rule senior : {p: Person} / p.age >= 65 ==> print("Senior detected")
rule young : {p: Person} / p.age < 18 ==> print("Young person")
`
	if err := os.WriteFile(tsdFile, []byte(content), 0644); err != nil {
		t.Fatalf("Erreur création fichier: %v", err)
	}
	// Construire le réseau
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(tsdFile, nil, storage)
	if err != nil {
		t.Fatalf("Erreur construction réseau: %v", err)
	}
	// Vérifier la structure de base
	stats := network.GetNetworkStats()
	if stats["type_nodes"].(int) != 1 {
		t.Errorf("Attendu 1 TypeNode, obtenu %d", stats["type_nodes"].(int))
	}
	if stats["terminal_nodes"].(int) != 3 {
		t.Errorf("Attendu 3 TerminalNodes, obtenu %d", stats["terminal_nodes"].(int))
	}
	// Soumettre des faits et vérifier le comportement
	// Les IDs sont fournis directement car c'est le test du réseau RETE
	// En production, les IDs seraient générés automatiquement par le pipeline d'ingestion
	facts := []Fact{
		{ID: "Alice", Type: "Person", Fields: map[string]interface{}{"age": 25.0, "name": "Alice"}},
		{ID: "Bob", Type: "Person", Fields: map[string]interface{}{"age": 70.0, "name": "Bob"}},
		{ID: "Charlie", Type: "Person", Fields: map[string]interface{}{"age": 15.0, "name": "Charlie"}},
	}
	for i, fact := range facts {
		if err := network.SubmitFact(&fact); err != nil {
			t.Errorf("Erreur ajout fait %d: %v", i, err)
		}
	}
	// Vérifier les activations
	activatedCount := 0
	for _, terminalNode := range network.TerminalNodes {
		memory := terminalNode.GetMemory()
		activatedCount += len(memory.Tokens)
	}
	// Alice (25 ans) devrait activer 'adult'
	// Bob (70 ans) devrait activer 'adult' et 'senior'
	// Charlie (15 ans) devrait activer 'young'
	// Total: 4 activations
	if activatedCount != 4 {
		t.Errorf("Attendu 4 activations, obtenu %d", activatedCount)
	}
	t.Logf("✅ Règles simples: backward compatible")
}

// TestBackwardCompatibility_ExistingBehavior vérifie que le comportement
// existant du réseau RETE n'a pas changé
func TestBackwardCompatibility_ExistingBehavior(t *testing.T) {
	tempDir := t.TempDir()
	tsdFile := filepath.Join(tempDir, "behavior_test.tsd")
	// Test simplifié avec des règles alpha uniquement (pas de join)
	content := `type Order(#id: string, amount:number)
type Customer(#id: string, name: string, vip:number)
action print(message: string)
rule large_order : {o: Order} / o.amount > 1000 ==> print("Large order")
rule vip_customer : {c: Customer} / c.vip == 1 ==> print("VIP customer")
`
	if err := os.WriteFile(tsdFile, []byte(content), 0644); err != nil {
		t.Fatalf("Erreur création fichier: %v", err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(tsdFile, nil, storage)
	if err != nil {
		t.Fatalf("Erreur construction réseau: %v", err)
	}
	// Test 1: TypeNode sharing fonctionne toujours
	stats := network.GetNetworkStats()
	if stats["type_nodes"].(int) != 2 {
		t.Errorf("TypeNode sharing: attendu 2 TypeNodes, obtenu %d", stats["type_nodes"].(int))
	}
	// Test 2: Ajout de faits
	customer := Fact{
		ID:   "Alice",
		Type: "Customer",
		Fields: map[string]interface{}{
			"name": "Alice",
			"vip":  1.0,
		},
	}
	order := Fact{
		ID:   "order1",
		Type: "Order",
		Fields: map[string]interface{}{
			"amount": 1500.0,
		},
	}
	if err := network.SubmitFact(&customer); err != nil {
		t.Fatalf("Erreur ajout customer: %v", err)
	}
	if err := network.SubmitFact(&order); err != nil {
		t.Fatalf("Erreur ajout order: %v", err)
	}
	// Test 3: Vérifier les activations
	activatedCount := 0
	activatedRules := []string{}
	for ruleName, terminalNode := range network.TerminalNodes {
		memory := terminalNode.GetMemory()
		if len(memory.Tokens) > 0 {
			activatedCount += len(memory.Tokens)
			activatedRules = append(activatedRules, ruleName)
		}
	}
	if activatedCount != 2 {
		t.Errorf("Attendu 2 activations (large_order, vip_customer), obtenu %d", activatedCount)
		for _, rule := range activatedRules {
			t.Logf("  Activation: %s", rule)
		}
	}
	// Test 4: Suppression de fait
	// L'identifiant interne est: Type_ID
	if err := network.RetractFact("Order_order1"); err != nil {
		t.Fatalf("Erreur suppression fait: %v", err)
	}
	activatedCount = 0
	for _, terminalNode := range network.TerminalNodes {
		memory := terminalNode.GetMemory()
		activatedCount += len(memory.Tokens)
	}
	if activatedCount != 1 {
		t.Errorf("Après suppression: attendu 1 activation (vip_customer), obtenu %d", activatedCount)
	}
	t.Logf("✅ Comportement existant: backward compatible")
}

// TestNoRegression_AllPreviousTests exécute un ensemble de tests
// qui valident qu'il n'y a pas de régression
func TestNoRegression_AllPreviousTests(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		factCount   int
		activations int
	}{
		{
			name: "Single condition",
			content: `type Person(#id: string, age:number)
action print(message: string)
rule adult : {p: Person} / p.age >= 18 ==> print("Adult")`,
			factCount:   2,
			activations: 1,
		},
		{
			name: "Multiple conditions AND",
			content: `type Person(#id: string, age: number, name:string)
action print(message: string)
rule specific : {p: Person} / p.age > 18 AND p.name == 'Alice' ==> print("Found")`,
			factCount:   2,
			activations: 1,
		},
		{
			name: "Multiple conditions OR",
			content: `type Person(#id: string, age:number)
action print(message: string)
rule young_or_old : {p: Person} / p.age < 18 OR p.age > 65 ==> print("Young or old")`,
			factCount:   3,
			activations: 2,
		},
		{
			name: "Numeric comparisons",
			content: `type Product(#id: string, price:number)
action print(message: string)
rule expensive : {p: Product} / p.price > 100 ==> print("Expensive")
rule cheap : {p: Product} / p.price <= 50 ==> print("Cheap")`,
			factCount:   3,
			activations: 2,
		},
		{
			name: "String equality",
			content: `type User(#id: string, role:string)
action print(message: string)
rule admin : {u: User} / u.role == 'admin' ==> print("Admin user")`,
			factCount:   2,
			activations: 1,
		},
		{
			name: "Boolean conditions",
			content: `type Account(#id: string, active:number)
action print(message: string)
rule active_account : {a: Account} / a.active == 1 ==> print("Active")`,
			factCount:   2,
			activations: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			tsdFile := filepath.Join(tempDir, "test.tsd")
			if err := os.WriteFile(tsdFile, []byte(tt.content), 0644); err != nil {
				t.Fatalf("Erreur création fichier: %v", err)
			}
			storage := NewMemoryStorage()
			pipeline := NewConstraintPipeline()
			network, _, err := pipeline.IngestFile(tsdFile, nil, storage)
			if err != nil {
				t.Fatalf("Erreur construction réseau: %v", err)
			}
			// Vérifier que le réseau est construit correctement
			stats := network.GetNetworkStats()
			if stats["terminal_nodes"].(int) < 1 {
				t.Errorf("Aucun TerminalNode créé")
			}
			// Ajouter des faits de test appropriés
			// Les IDs sont fournis directement pour les tests unitaires
			switch tt.name {
			case "Single condition":
				network.SubmitFact(&Fact{ID: "p1", Type: "Person", Fields: map[string]interface{}{"age": 25.0}})
				network.SubmitFact(&Fact{ID: "p2", Type: "Person", Fields: map[string]interface{}{"age": 15.0}})
			case "Multiple conditions AND":
				network.SubmitFact(&Fact{ID: "p1", Type: "Person", Fields: map[string]interface{}{"age": 25.0, "name": "Alice"}})
				network.SubmitFact(&Fact{ID: "p2", Type: "Person", Fields: map[string]interface{}{"age": 25.0, "name": "Bob"}})
			case "Multiple conditions OR":
				network.SubmitFact(&Fact{ID: "p1", Type: "Person", Fields: map[string]interface{}{"age": 15.0}})
				network.SubmitFact(&Fact{ID: "p2", Type: "Person", Fields: map[string]interface{}{"age": 70.0}})
				network.SubmitFact(&Fact{ID: "p3", Type: "Person", Fields: map[string]interface{}{"age": 30.0}})
			case "Numeric comparisons":
				network.SubmitFact(&Fact{ID: "prod1", Type: "Product", Fields: map[string]interface{}{"price": 150.0}})
				network.SubmitFact(&Fact{ID: "prod2", Type: "Product", Fields: map[string]interface{}{"price": 30.0}})
				network.SubmitFact(&Fact{ID: "prod3", Type: "Product", Fields: map[string]interface{}{"price": 75.0}})
			case "String equality":
				network.SubmitFact(&Fact{ID: "u1", Type: "User", Fields: map[string]interface{}{"role": "admin"}})
				network.SubmitFact(&Fact{ID: "u2", Type: "User", Fields: map[string]interface{}{"role": "user"}})
			case "Boolean conditions":
				network.SubmitFact(&Fact{ID: "a1", Type: "Account", Fields: map[string]interface{}{"active": 1.0}})
				network.SubmitFact(&Fact{ID: "a2", Type: "Account", Fields: map[string]interface{}{"active": 0.0}})
			}
			// Vérifier le nombre d'activations
			activatedCount := 0
			for _, terminalNode := range network.TerminalNodes {
				memory := terminalNode.GetMemory()
				activatedCount += len(memory.Tokens)
			}
			if activatedCount != tt.activations {
				t.Errorf("Attendu %d activations, obtenu %d", tt.activations, activatedCount)
			}
			t.Logf("✅ Test '%s': pas de régression", tt.name)
		})
	}
}

// TestBackwardCompatibility_TypeNodeSharing vérifie que le partage
// des TypeNodes fonctionne toujours correctement
func TestBackwardCompatibility_TypeNodeSharing(t *testing.T) {
	tempDir := t.TempDir()
	tsdFile := filepath.Join(tempDir, "sharing.tsd")
	// Plusieurs règles sur le même type
	content := `type Person(#id: string, age: number, name:string)
action print(message: string)
rule r1 : {p: Person} / p.age > 18 ==> print("R1")
rule r2 : {p: Person} / p.age > 30 ==> print("R2")
rule r3 : {p: Person} / p.age > 50 ==> print("R3")
rule r4 : {p: Person} / p.name == 'Alice' ==> print("R4")
`
	if err := os.WriteFile(tsdFile, []byte(content), 0644); err != nil {
		t.Fatalf("Erreur création fichier: %v", err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(tsdFile, nil, storage)
	if err != nil {
		t.Fatalf("Erreur construction réseau: %v", err)
	}
	stats := network.GetNetworkStats()
	// Vérifier qu'un seul TypeNode est créé pour Person
	if stats["type_nodes"].(int) != 1 {
		t.Errorf("TypeNode sharing: attendu 1 TypeNode, obtenu %d", stats["type_nodes"].(int))
	}
	// Vérifier que 4 règles sont créées
	if stats["terminal_nodes"].(int) != 4 {
		t.Errorf("Attendu 4 TerminalNodes, obtenu %d", stats["terminal_nodes"].(int))
	}
	// Soumettre un fait et vérifier qu'il est propagé à toutes les règles
	fact := Fact{
		ID:   "Alice",
		Type: "Person",
		Fields: map[string]interface{}{
			"age":  55.0,
			"name": "Alice",
		},
	}
	if err := network.SubmitFact(&fact); err != nil {
		t.Fatalf("Erreur ajout fait: %v", err)
	}
	// Ce fait devrait activer r1, r2, r3 et r4 (4 activations)
	activatedCount := 0
	activatedRules := []string{}
	for ruleName, terminalNode := range network.TerminalNodes {
		memory := terminalNode.GetMemory()
		if len(memory.Tokens) > 0 {
			activatedCount += len(memory.Tokens)
			activatedRules = append(activatedRules, ruleName)
		}
	}
	if activatedCount != 4 {
		t.Errorf("Attendu 4 activations, obtenu %d", activatedCount)
		for _, rule := range activatedRules {
			t.Logf("  Activation: %s", rule)
		}
	}
	t.Logf("✅ TypeNode sharing: backward compatible")
}

// TestBackwardCompatibility_LifecycleManagement vérifie que la gestion
// du lifecycle des nœuds fonctionne toujours correctement
func TestBackwardCompatibility_LifecycleManagement(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	// Créer un TypeNode
	typeDef := TypeDefinition{Type: "type", Name: "Person", Fields: []Field{}}
	parentNode := NewTypeNode("person", typeDef, storage)
	network.TypeNodes["Person"] = parentNode
	// Créer le builder
	builder := NewAlphaChainBuilder(network, storage)
	// Règle 1: age > 18
	conditions := []SimpleCondition{
		NewSimpleCondition("comparison", "p.age", ">", 18),
	}
	chain1, err := builder.BuildChain(conditions, "p", parentNode, "rule1")
	if err != nil {
		t.Fatalf("Erreur construction chain1: %v", err)
	}
	// Règle 2: age > 18 (même condition, devrait réutiliser)
	chain2, err := builder.BuildChain(conditions, "p", parentNode, "rule2")
	if err != nil {
		t.Fatalf("Erreur construction chain2: %v", err)
	}
	// Vérifier que le même nœud est réutilisé
	if chain1.Nodes[0].ID != chain2.Nodes[0].ID {
		t.Errorf("Le même nœud devrait être réutilisé")
	}
	// Vérifier le compteur de références dans le LifecycleManager
	lifecycle, exists := network.LifecycleManager.GetNodeLifecycle(chain1.Nodes[0].ID)
	if !exists {
		t.Fatal("Le nœud devrait être dans le LifecycleManager")
	}
	if lifecycle.GetRefCount() != 2 {
		t.Errorf("Attendu 2 références, obtenu %d", lifecycle.GetRefCount())
	}
	t.Logf("✅ Lifecycle management: backward compatible")
}

// TestBackwardCompatibility_RuleRemoval vérifie que la suppression
// de règles fonctionne toujours correctement
func TestBackwardCompatibility_RuleRemoval(t *testing.T) {
	tempDir := t.TempDir()
	tsdFile := filepath.Join(tempDir, "removal.tsd")
	content := `type Person(#id: string, age:number)
action print(message: string)
action print(message: string)
rule adult : {p: Person} / p.age >= 18 ==> print("Adult")
rule senior : {p: Person} / p.age >= 65 ==> print("Senior")
rule teenager : {p: Person} / p.age >= 13 AND p.age < 18 ==> print("Teen")
`
	if err := os.WriteFile(tsdFile, []byte(content), 0644); err != nil {
		t.Fatalf("Erreur création fichier: %v", err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(tsdFile, nil, storage)
	if err != nil {
		t.Fatalf("Erreur construction réseau: %v", err)
	}
	// Vérifier l'état initial
	statsInitial := network.GetNetworkStats()
	initialTerminalCount := statsInitial["terminal_nodes"].(int)
	if initialTerminalCount != 3 {
		t.Errorf("Attendu 3 TerminalNodes initialement, obtenu %d", initialTerminalCount)
	}
	// Supprimer la règle 'adult'
	if err := network.RemoveRule("adult"); err != nil {
		t.Fatalf("Erreur suppression règle: %v", err)
	}
	// Vérifier qu'il reste 2 TerminalNodes
	statsAfter := network.GetNetworkStats()
	afterTerminalCount := statsAfter["terminal_nodes"].(int)
	if afterTerminalCount != 2 {
		t.Errorf("Attendu 2 TerminalNodes après suppression, obtenu %d", afterTerminalCount)
	}
	// Vérifier que les règles restantes fonctionnent toujours
	fact1 := Fact{ID: "senior1", Type: "Person", Fields: map[string]interface{}{"age": 70.0}}
	fact2 := Fact{ID: "teen1", Type: "Person", Fields: map[string]interface{}{"age": 15.0}}
	network.SubmitFact(&fact1)
	network.SubmitFact(&fact2)
	activatedCount := 0
	for _, terminalNode := range network.TerminalNodes {
		memory := terminalNode.GetMemory()
		activatedCount += len(memory.Tokens)
	}
	// senior1 devrait activer 'senior', teen1 devrait activer 'teenager'
	if activatedCount != 2 {
		t.Errorf("Attendu 2 activations, obtenu %d", activatedCount)
	}
	t.Logf("✅ Rule removal: backward compatible")
}

// TestBackwardCompatibility_PerformanceCharacteristics vérifie que
// les caractéristiques de performance sont maintenues ou améliorées
func TestBackwardCompatibility_PerformanceCharacteristics(t *testing.T) {
	tempDir := t.TempDir()
	tsdFile := filepath.Join(tempDir, "perf.tsd")
	// Créer un ensemble de règles avec beaucoup de conditions partagées
	content := `type Person(#id: string, age: number, name: string, country:string)
action print(message: string)
rule r1 : {p: Person} / p.age > 18 ==> print("R1")
rule r2 : {p: Person} / p.age > 18 AND p.country == 'USA' ==> print("R2")
rule r3 : {p: Person} / p.age > 18 AND p.country == 'Canada' ==> print("R3")
rule r4 : {p: Person} / p.age > 30 ==> print("R4")
rule r5 : {p: Person} / p.age > 30 AND p.name == 'Alice' ==> print("R5")
`
	if err := os.WriteFile(tsdFile, []byte(content), 0644); err != nil {
		t.Fatalf("Erreur création fichier: %v", err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(tsdFile, nil, storage)
	if err != nil {
		t.Fatalf("Erreur construction réseau: %v", err)
	}
	stats := network.GetNetworkStats()
	// Avec le partage, on devrait avoir moins d'AlphaNodes que de conditions
	// (5 règles avec des conditions partagées)
	alphaCount := stats["alpha_nodes"].(int)
	if alphaCount > 7 {
		t.Errorf("Trop d'AlphaNodes créés: %d (le partage devrait réduire ce nombre)", alphaCount)
	}
	// Vérifier que toutes les règles sont créées
	terminalCount := stats["terminal_nodes"].(int)
	if terminalCount != 5 {
		t.Errorf("Attendu 5 TerminalNodes, obtenu %d", terminalCount)
	}
	// Ajouter un fait et vérifier que le traitement est efficace
	fact := Fact{
		ID:   "person1",
		Type: "Person",
		Fields: map[string]interface{}{
			"age":     35.0,
			"name":    "Alice",
			"country": "USA",
		},
	}
	if err := network.SubmitFact(&fact); err != nil {
		t.Fatalf("Erreur ajout fait: %v", err)
	}
	// Ce fait devrait activer r1, r2, r4, et r5 (4 activations)
	activatedCount := 0
	for _, terminalNode := range network.TerminalNodes {
		memory := terminalNode.GetMemory()
		activatedCount += len(memory.Tokens)
	}
	if activatedCount != 4 {
		t.Errorf("Attendu 4 activations, obtenu %d", activatedCount)
	}
	t.Logf("✅ Performance: %d AlphaNodes pour 5 règles (partage efficace)", alphaCount)
}
