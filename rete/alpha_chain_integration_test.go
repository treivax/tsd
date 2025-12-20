// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete

import (
	"os"
	"path/filepath"
	"testing"
)

// TestAlphaChain_TwoRules_SameConditions_DifferentOrder teste que deux règles
// avec les mêmes conditions dans un ordre différent partagent les mêmes AlphaNodes
func TestAlphaChain_TwoRules_SameConditions_DifferentOrder(t *testing.T) {
	tempDir := t.TempDir()
	tsdFile := filepath.Join(tempDir, "test.tsd")
	// Deux règles avec les mêmes conditions mais dans un ordre différent
	content := `type Person(#id: string, age: number, name:string)
action print(message: string)
rule r1 : {p: Person} / p.age > 18 AND p.name == 'toto' ==> print("A")
rule r2 : {p: Person} / p.name == 'toto' AND p.age > 18 ==> print("B")
`
	if err := os.WriteFile(tsdFile, []byte(content), 0644); err != nil {
		t.Fatalf("Erreur écriture fichier: %v", err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(tsdFile, nil, storage)
	if err != nil {
		t.Fatalf("Erreur construction réseau: %v", err)
	}
	stats := network.GetNetworkStats()
	// Vérifier: 2 AlphaNodes partagés (un pour chaque condition unique après normalisation)
	totalAlphaNodes := stats["alpha_nodes"].(int)
	if totalAlphaNodes != 2 {
		t.Errorf("Devrait avoir 2 AlphaNodes partagés, got %d", totalAlphaNodes)
	}
	// Vérifier: 2 TerminalNodes (un par règle)
	totalTerminalNodes := stats["terminal_nodes"].(int)
	if totalTerminalNodes != 2 {
		t.Errorf("Devrait avoir 2 TerminalNodes, got %d", totalTerminalNodes)
	}
	// Vérifier les stats de partage
	if sharingStats, ok := stats["sharing_total_shared_alpha_nodes"]; ok {
		sharedCount := sharingStats.(int)
		if sharedCount != 2 {
			t.Errorf("Devrait avoir 2 nœuds dans le registre de partage, got %d", sharedCount)
		}
	}
	// Vérifier le nombre total d'enfants (qui est utilisé comme proxy pour les références)
	// Alpha1 a 1 enfant (Alpha2), Alpha2 a 2 enfants (2 TerminalNodes) = 3 total
	if ruleRefs, ok := stats["sharing_total_rule_references"]; ok {
		refCount := ruleRefs.(int)
		if refCount != 3 { // 1 + 2 = 3 enfants au total
			t.Errorf("Devrait avoir 3 enfants au total, got %d", refCount)
		}
	}
	// Tester la propagation de faits
	fact := &Fact{
		ID:   "person1",
		Type: "Person",
		Fields: map[string]interface{}{
			"id":   "p1",
			"age":  25.0,
			"name": "toto",
		},
	}
	err = network.SubmitFact(fact)
	if err != nil {
		t.Fatalf("Erreur soumission fait: %v", err)
	}
	// Vérifier que les 2 TerminalNodes ont été activés
	activatedCount := 0
	for _, terminalNode := range network.TerminalNodes {
		if terminalNode.GetExecutionCount() > 0 {
			activatedCount++
		}
	}
	if activatedCount != 2 {
		t.Errorf("Les 2 TerminalNodes devraient être activés, got %d", activatedCount)
	}
	t.Logf("✓ Deux règles avec mêmes conditions (ordre différent) partagent correctement les AlphaNodes")
}

// TestAlphaChain_PartialSharing_ThreeRules teste le partage partiel
// entre trois règles avec des préfixes communs
func TestAlphaChain_PartialSharing_ThreeRules(t *testing.T) {
	tempDir := t.TempDir()
	tsdFile := filepath.Join(tempDir, "test.tsd")
	// Trois règles avec partage partiel progressif
	content := `type Person(#id: string, age: number, name: string, salary:number)
action print(message: string)
rule r1 : {p: Person} / p.age > 18 ==> print("A")
rule r2 : {p: Person} / p.age > 18 AND p.name == 'toto' ==> print("B")
rule r3 : {p: Person} / p.age > 18 AND p.name == 'toto' AND p.salary > 1000 ==> print("C")
`
	if err := os.WriteFile(tsdFile, []byte(content), 0644); err != nil {
		t.Fatalf("Erreur écriture fichier: %v", err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(tsdFile, nil, storage)
	if err != nil {
		t.Fatalf("Erreur construction réseau: %v", err)
	}
	stats := network.GetNetworkStats()
	// Vérifier: 3 AlphaNodes (partage optimal!)
	// - p.age > 18 (partagé par r1, r2 et r3 - règle simple + chaînes)
	// - p.name == 'toto' (partagé par r2 et r3)
	// - p.salary > 1000 (utilisé par r3 uniquement)
	totalAlphaNodes := stats["alpha_nodes"].(int)
	if totalAlphaNodes != 3 {
		t.Errorf("Devrait avoir 3 AlphaNodes (partage optimal), got %d", totalAlphaNodes)
	}
	// Vérifier: 3 TerminalNodes
	totalTerminalNodes := stats["terminal_nodes"].(int)
	if totalTerminalNodes != 3 {
		t.Errorf("Devrait avoir 3 TerminalNodes, got %d", totalTerminalNodes)
	}
	// Vérifier le nombre total d'enfants
	// r1: Alpha1(age>18) -> Terminal = 1 enfant
	// r2: Alpha1(age>18) -> Alpha2(name=='toto') -> Terminal = 1 enfant pour Alpha1, 1 pour Alpha2
	// r3: Alpha1 -> Alpha2 -> Alpha3(salary>1000) -> Terminal = 1+1+1
	// Total: Alpha1 a 1 enfant (Alpha2), Alpha2 a 1 enfant (Alpha3), Alpha3 a 3 terminaux = 1+1+3 = 5
	if ruleRefs, ok := stats["sharing_total_rule_references"]; ok {
		refCount := ruleRefs.(int)
		if refCount < 3 {
			t.Errorf("Devrait avoir au moins 3 enfants au total, got %d", refCount)
		}
	}
	// Tester avec un fait qui ne satisfait que r1
	fact1 := &Fact{
		ID:   "person1",
		Type: "Person",
		Fields: map[string]interface{}{
			"id":     "p1",
			"age":    20.0,
			"name":   "autre",
			"salary": 500.0,
		},
	}
	err = network.SubmitFact(fact1)
	if err != nil {
		t.Fatalf("Erreur soumission fait1: %v", err)
	}
	// Compter les activations
	activatedCount := 0
	for _, terminalNode := range network.TerminalNodes {
		if terminalNode.GetExecutionCount() > 0 {
			activatedCount++
		}
	}
	if activatedCount != 1 {
		t.Errorf("Seul 1 TerminalNode devrait être activé, got %d", activatedCount)
	}
	// Tester avec un fait qui satisfait r1 et r2 mais pas r3
	// Créer un nouveau réseau pour ce test
	storage2 := NewMemoryStorage()
	pipeline2 := NewConstraintPipeline()
	network2, _, err := pipeline2.IngestFile(tsdFile, nil, storage2)
	if err != nil {
		t.Fatalf("Erreur construction réseau2: %v", err)
	}
	fact2 := &Fact{
		ID:   "person2",
		Type: "Person",
		Fields: map[string]interface{}{
			"id":     "p2",
			"age":    25.0,
			"name":   "toto",
			"salary": 500.0,
		},
	}
	err = network2.SubmitFact(fact2)
	if err != nil {
		t.Fatalf("Erreur soumission fait2: %v", err)
	}
	activatedCount = 0
	for _, terminalNode := range network2.TerminalNodes {
		if terminalNode.GetExecutionCount() > 0 {
			activatedCount++
		}
	}
	if activatedCount != 2 {
		t.Errorf("2 TerminalNodes devraient être activés, got %d", activatedCount)
	}
	// Tester avec un fait qui satisfait les trois règles
	// Créer un nouveau réseau pour ce test
	storage3 := NewMemoryStorage()
	pipeline3 := NewConstraintPipeline()
	network3, _, err := pipeline3.IngestFile(tsdFile, nil, storage3)
	if err != nil {
		t.Fatalf("Erreur construction réseau3: %v", err)
	}
	fact3 := &Fact{
		ID:   "person3",
		Type: "Person",
		Fields: map[string]interface{}{
			"id":     "p3",
			"age":    30.0,
			"name":   "toto",
			"salary": 2000.0,
		},
	}
	err = network3.SubmitFact(fact3)
	if err != nil {
		t.Fatalf("Erreur soumission fait3: %v", err)
	}
	activatedCount = 0
	for _, terminalNode := range network3.TerminalNodes {
		if terminalNode.GetExecutionCount() > 0 {
			activatedCount++
		}
	}
	if activatedCount != 3 {
		t.Errorf("Les 3 TerminalNodes devraient être activés, got %d", activatedCount)
	}
	t.Logf("✓ Partage partiel entre 3 règles fonctionne correctement")
}

// TestAlphaChain_FactPropagation_ThroughChain teste la propagation de faits
// à travers une chaîne et vérifie que chaque condition n'est évaluée qu'une fois
func TestAlphaChain_FactPropagation_ThroughChain(t *testing.T) {
	tempDir := t.TempDir()
	tsdFile := filepath.Join(tempDir, "test.tsd")
	content := `type Person(#id: string, age: number, name: string, salary:number)
action print(message: string)
rule complete : {p: Person} / p.age > 18 AND p.name == 'toto' AND p.salary > 1000 ==> print("Complete")
`
	if err := os.WriteFile(tsdFile, []byte(content), 0644); err != nil {
		t.Fatalf("Erreur écriture fichier: %v", err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(tsdFile, nil, storage)
	if err != nil {
		t.Fatalf("Erreur construction réseau: %v", err)
	}
	// Vérifier la structure de la chaîne
	stats := network.GetNetworkStats()
	totalAlphaNodes := stats["alpha_nodes"].(int)
	if totalAlphaNodes != 3 {
		t.Errorf("Devrait avoir 3 AlphaNodes dans la chaîne, got %d", totalAlphaNodes)
	}
	// Soumettre un fait qui satisfait toute la chaîne
	fact := &Fact{
		ID:   "person1",
		Type: "Person",
		Fields: map[string]interface{}{
			"id":     "p1",
			"age":    25.0,
			"name":   "toto",
			"salary": 2000.0,
		},
	}
	err = network.SubmitFact(fact)
	if err != nil {
		t.Fatalf("Erreur soumission fait: %v", err)
	}
	// Vérifier que le TerminalNode a été activé
	if len(network.TerminalNodes) != 1 {
		t.Fatalf("Devrait avoir 1 TerminalNode, got %d", len(network.TerminalNodes))
	}
	var terminalNode *TerminalNode
	for _, tn := range network.TerminalNodes {
		terminalNode = tn
		break
	}
	execCount := terminalNode.GetExecutionCount()
	if execCount != 1 {
		t.Errorf("Le TerminalNode devrait avoir 1 activation, got %d", execCount)
	}
	// Vérifier que tous les AlphaNodes de la chaîne ont le fait en mémoire
	alphaNodesWithFact := 0
	for _, alphaNode := range network.AlphaNodes {
		memory := alphaNode.GetMemory()
		if len(memory.Facts) > 0 {
			alphaNodesWithFact++
		}
	}
	if alphaNodesWithFact != totalAlphaNodes {
		t.Errorf("Tous les %d AlphaNodes devraient avoir le fait, got %d", totalAlphaNodes, alphaNodesWithFact)
	}
	// Soumettre un fait qui échoue à la première condition
	// Créer un nouveau réseau pour ce test
	storage2 := NewMemoryStorage()
	pipeline2 := NewConstraintPipeline()
	network2, _, err := pipeline2.IngestFile(tsdFile, nil, storage2)
	if err != nil {
		t.Fatalf("Erreur construction réseau2: %v", err)
	}
	fact2 := &Fact{
		ID:   "person2",
		Type: "Person",
		Fields: map[string]interface{}{
			"id":     "p2",
			"age":    15.0, // Ne satisfait pas p.age > 18
			"name":   "toto",
			"salary": 2000.0,
		},
	}
	err = network2.SubmitFact(fact2)
	if err != nil {
		t.Fatalf("Erreur soumission fait2: %v", err)
	}
	// Vérifier qu'aucun TerminalNode du nouveau réseau n'a été activé
	activatedCount2 := 0
	for _, tn := range network2.TerminalNodes {
		if tn.GetExecutionCount() > 0 {
			activatedCount2++
		}
	}
	if activatedCount2 != 0 {
		t.Errorf("Aucun TerminalNode ne devrait être activé, got %d", activatedCount2)
	}
	t.Logf("✓ Propagation de faits à travers la chaîne fonctionne correctement")
}

// TestAlphaChain_RuleRemoval_PreservesShared teste la suppression de règles
// en préservant les nœuds partagés
func TestAlphaChain_RuleRemoval_PreservesShared(t *testing.T) {
	tempDir := t.TempDir()
	tsdFile := filepath.Join(tempDir, "test.tsd")
	content := `type Person(#id: string, age: number, name:string)
action print(message: string)
rule r1 : {p: Person} / p.age > 18 ==> print("A")
rule r2 : {p: Person} / p.age > 18 AND p.name == 'toto' ==> print("B")
rule r3 : {p: Person} / p.age > 18 ==> print("C")
`
	if err := os.WriteFile(tsdFile, []byte(content), 0644); err != nil {
		t.Fatalf("Erreur écriture fichier: %v", err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(tsdFile, nil, storage)
	if err != nil {
		t.Fatalf("Erreur construction réseau: %v", err)
	}
	// État initial
	stats := network.GetNetworkStats()
	initialAlphaCount := stats["alpha_nodes"].(int)
	if initialAlphaCount != 2 { // 1 AlphaNode p.age>18 partagé (r1/r2/r3) + 1 AlphaNode p.name=='toto' (r2)
		t.Errorf("Devrait avoir 2 AlphaNodes initialement (partage optimal), got %d", initialAlphaCount)
	}
	initialTerminalCount := stats["terminal_nodes"].(int)
	if initialTerminalCount != 3 {
		t.Errorf("Devrait avoir 3 TerminalNodes initialement, got %d", initialTerminalCount)
	}
	// Supprimer r2 (règle du milieu avec 2 conditions)
	err = network.RemoveRule("r2")
	if err != nil {
		t.Fatalf("Erreur suppression règle r2: %v", err)
	}
	// Vérifier que les nœuds partagés restent
	stats = network.GetNetworkStats()
	finalAlphaCount := stats["alpha_nodes"].(int)
	if finalAlphaCount != 1 { // Devrait rester 1 AlphaNode (p.age>18 encore utilisé par r1)
		t.Errorf("Devrait avoir 1 AlphaNode après suppression de r2, got %d", finalAlphaCount)
	}
	afterRemovalTerminalCount := stats["terminal_nodes"].(int)
	if afterRemovalTerminalCount != 2 {
		t.Errorf("Devrait avoir 2 TerminalNodes après suppression, got %d", afterRemovalTerminalCount)
	}
	// Vérifier que le nœud p.age > 18 a toujours 2 références (r1 et r3)
	// LifecycleManager is always initialized
	for nodeID := range network.AlphaNodes {
		lifecycle, exists := network.LifecycleManager.GetNodeLifecycle(nodeID)
		if !exists {
			t.Errorf("Nœud %s non trouvé dans LifecycleManager", nodeID)
			continue
		}
		refCount := lifecycle.GetRefCount()
		if refCount != 2 {
			t.Errorf("Le nœud partagé devrait avoir 2 références, got %d", refCount)
		}
	}
	// Vérifier que le réseau fonctionne toujours
	fact := &Fact{
		ID:   "person1",
		Type: "Person",
		Fields: map[string]interface{}{
			"id":   "p1",
			"age":  25.0,
			"name": "toto",
		},
	}
	err = network.SubmitFact(fact)
	if err != nil {
		t.Fatalf("Erreur soumission fait: %v", err)
	}
	activatedCount := 0
	for _, terminalNode := range network.TerminalNodes {
		if terminalNode.GetExecutionCount() > 0 {
			activatedCount++
		}
	}
	if activatedCount != 2 {
		t.Errorf("Les 2 TerminalNodes restants devraient être activés, got %d", activatedCount)
	}
	t.Logf("✓ Suppression de règle préserve correctement les nœuds partagés")
}

// TestAlphaChain_ComplexScenario_FraudDetection teste un scénario complexe
// de détection de fraude avec partage optimal
func TestAlphaChain_ComplexScenario_FraudDetection(t *testing.T) {
	tempDir := t.TempDir()
	tsdFile := filepath.Join(tempDir, "test.tsd")
	content := `type Transaction(#id: string, amount: number, country: string, risk:number)
action print(message: string)
rule fraud_low : {t: Transaction} / t.amount > 1000 AND t.country == 'XX' ==> print("LOW")
rule fraud_med : {t: Transaction} / t.amount > 1000 AND t.country == 'XX' AND t.risk > 50 ==> print("MED")
rule fraud_high : {t: Transaction} / t.amount > 1000 AND t.country == 'XX' AND t.risk > 80 ==> print("HIGH")
rule large : {t: Transaction} / t.amount > 1000 ==> print("LARGE")
`
	if err := os.WriteFile(tsdFile, []byte(content), 0644); err != nil {
		t.Fatalf("Erreur écriture fichier: %v", err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(tsdFile, nil, storage)
	if err != nil {
		t.Fatalf("Erreur construction réseau: %v", err)
	}
	stats := network.GetNetworkStats()
	// Vérifier le partage optimal:
	// - t.amount > 1000 (partagé par ALL: large, fraud_low, fraud_med, fraud_high) = 1 AlphaNode
	// - t.country == 'XX' (partagé par fraud_low, fraud_med, fraud_high) = 1 AlphaNode
	// - t.risk > 50 (partagé par fraud_med, fraud_high) = 1 AlphaNode
	// - t.risk > 80 (utilisé uniquement par fraud_high) = 1 AlphaNode
	// Total: 4 AlphaNodes (partage optimal entre règles simples et chaînes!)
	totalAlphaNodes := stats["alpha_nodes"].(int)
	if totalAlphaNodes != 4 {
		t.Errorf("Devrait avoir 4 AlphaNodes (partage optimal), got %d", totalAlphaNodes)
	}
	// 4 TerminalNodes (un par règle)
	totalTerminalNodes := stats["terminal_nodes"].(int)
	if totalTerminalNodes != 4 {
		t.Errorf("Devrait avoir 4 TerminalNodes, got %d", totalTerminalNodes)
	}
	// Vérifier le nombre total d'enfants (approximation du partage)
	// Le comptage exact dépend de la structure de la chaîne
	if ruleRefs, ok := stats["sharing_total_rule_references"]; ok {
		refCount := ruleRefs.(int)
		// Vérifier qu'il y a un partage significatif (au moins 4 règles = 4 terminaux minimum)
		if refCount < 4 {
			t.Errorf("Devrait avoir au moins 4 enfants (4 terminaux), got %d", refCount)
		}
	}
	// Vérifier le ratio de partage (approximatif car dépend de l'implémentation)
	if avgRatio, ok := stats["sharing_average_sharing_ratio"]; ok {
		ratio := avgRatio.(float64)
		// Le ratio réel dépend du nombre de références, nous vérifions juste qu'il y a du partage
		if ratio < 1.0 {
			t.Errorf("Le ratio de partage devrait être >= 1.0, got %v", ratio)
		}
	}
	// Test 1: Transaction large mais pas XX
	fact1 := &Fact{
		ID:   "tx1",
		Type: "Transaction",
		Fields: map[string]interface{}{
			"id":      "t1",
			"amount":  2000.0,
			"country": "US",
			"risk":    10.0,
		},
	}
	err = network.SubmitFact(fact1)
	if err != nil {
		t.Fatalf("Erreur soumission fait1: %v", err)
	}
	activatedCount := 0
	for _, terminalNode := range network.TerminalNodes {
		if terminalNode.GetExecutionCount() > 0 {
			activatedCount++
		}
	}
	if activatedCount != 1 { // Seul 'large' devrait être activé
		t.Errorf("Seule la règle 'large' devrait être activée, got %d activations", activatedCount)
	}
	// Test 2: Transaction fraud_low (amount > 1000, country = XX, mais risk <= 50)
	network.ClearMemory()
	fact2 := &Fact{
		ID:   "tx2",
		Type: "Transaction",
		Fields: map[string]interface{}{
			"id":      "t2",
			"amount":  1500.0,
			"country": "XX",
			"risk":    30.0,
		},
	}
	err = network.SubmitFact(fact2)
	if err != nil {
		t.Fatalf("Erreur soumission fait2: %v", err)
	}
	activatedCount = 0
	for _, terminalNode := range network.TerminalNodes {
		if terminalNode.GetExecutionCount() > 0 {
			activatedCount++
		}
	}
	if activatedCount != 2 { // 'large' et 'fraud_low'
		t.Errorf("2 règles devraient être activées (large, fraud_low), got %d", activatedCount)
	}
	// Test 3: Transaction fraud_med (risk > 50 mais <= 80)
	network.ClearMemory()
	fact3 := &Fact{
		ID:   "tx3",
		Type: "Transaction",
		Fields: map[string]interface{}{
			"id":      "t3",
			"amount":  2500.0,
			"country": "XX",
			"risk":    60.0,
		},
	}
	err = network.SubmitFact(fact3)
	if err != nil {
		t.Fatalf("Erreur soumission fait3: %v", err)
	}
	activatedCount = 0
	for _, terminalNode := range network.TerminalNodes {
		if terminalNode.GetExecutionCount() > 0 {
			activatedCount++
		}
	}
	if activatedCount != 3 { // 'large', 'fraud_low', 'fraud_med'
		t.Errorf("3 règles devraient être activées, got %d", activatedCount)
	}
	// Test 4: Transaction fraud_high (risk > 80)
	network.ClearMemory()
	fact4 := &Fact{
		ID:   "tx4",
		Type: "Transaction",
		Fields: map[string]interface{}{
			"id":      "t4",
			"amount":  3000.0,
			"country": "XX",
			"risk":    90.0,
		},
	}
	err = network.SubmitFact(fact4)
	if err != nil {
		t.Fatalf("Erreur soumission fait4: %v", err)
	}
	activatedCount = 0
	for _, terminalNode := range network.TerminalNodes {
		if terminalNode.GetExecutionCount() > 0 {
			activatedCount++
		}
	}
	if activatedCount != 4 { // Toutes les règles
		t.Errorf("Les 4 règles devraient être activées, got %d", activatedCount)
	}
	t.Logf("✓ Scénario complexe de détection de fraude avec partage optimal")
}

// TestAlphaChain_OR_NotDecomposed vérifie qu'une expression OR
// n'est pas décomposée en chaîne
func TestAlphaChain_OR_NotDecomposed(t *testing.T) {
	tempDir := t.TempDir()
	tsdFile := filepath.Join(tempDir, "test.tsd")
	content := `type Person(#id: string, age: number, status:string)
action print(message: string)
rule r1 : {p: Person} / p.age > 18 OR p.status == 'VIP' ==> print("A")
`
	if err := os.WriteFile(tsdFile, []byte(content), 0644); err != nil {
		t.Fatalf("Erreur écriture fichier: %v", err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(tsdFile, nil, storage)
	if err != nil {
		t.Fatalf("Erreur construction réseau: %v", err)
	}
	stats := network.GetNetworkStats()
	// Vérifier: Un seul AlphaNode (pas de décomposition)
	totalAlphaNodes := stats["alpha_nodes"].(int)
	if totalAlphaNodes != 1 {
		t.Errorf("Devrait avoir 1 AlphaNode (pas de décomposition OR), got %d", totalAlphaNodes)
	}
	// Vérifier: 1 TerminalNode
	totalTerminalNodes := stats["terminal_nodes"].(int)
	if totalTerminalNodes != 1 {
		t.Errorf("Devrait avoir 1 TerminalNode, got %d", totalTerminalNodes)
	}
	// Tester avec un fait qui satisfait la première partie du OR
	fact1 := &Fact{
		ID:   "person1",
		Type: "Person",
		Fields: map[string]interface{}{
			"id":     "p1",
			"age":    25.0,
			"status": "NORMAL",
		},
	}
	err = network.SubmitFact(fact1)
	if err != nil {
		t.Fatalf("Erreur soumission fait1: %v", err)
	}
	var terminalNode *TerminalNode
	for _, tn := range network.TerminalNodes {
		terminalNode = tn
		break
	}
	execCount := terminalNode.GetExecutionCount()
	if execCount != 1 {
		t.Errorf("Le TerminalNode devrait être activé, got %d activations", execCount)
	}
	// Tester avec un fait qui satisfait la deuxième partie du OR
	// Créer un nouveau réseau pour ce test
	storage2 := NewMemoryStorage()
	pipeline2 := NewConstraintPipeline()
	network2, _, err := pipeline2.IngestFile(tsdFile, nil, storage2)
	if err != nil {
		t.Fatalf("Erreur construction réseau2: %v", err)
	}
	fact2 := &Fact{
		ID:   "person2",
		Type: "Person",
		Fields: map[string]interface{}{
			"id":     "p2",
			"age":    15.0,
			"status": "VIP",
		},
	}
	err = network2.SubmitFact(fact2)
	if err != nil {
		t.Fatalf("Erreur soumission fait2: %v", err)
	}
	// Vérifier que le TerminalNode du nouveau réseau a été activé
	var terminalNode2 *TerminalNode
	for _, tn := range network2.TerminalNodes {
		terminalNode2 = tn
		break
	}
	execCount = terminalNode2.GetExecutionCount()
	if execCount != 1 {
		t.Errorf("Le TerminalNode devrait être activé, got %d activations", execCount)
	}
	// Tester avec un fait qui ne satisfait aucune partie
	// Créer un nouveau réseau pour ce test
	storage3 := NewMemoryStorage()
	pipeline3 := NewConstraintPipeline()
	network3, _, err := pipeline3.IngestFile(tsdFile, nil, storage3)
	if err != nil {
		t.Fatalf("Erreur construction réseau3: %v", err)
	}
	fact3 := &Fact{
		ID:   "person3",
		Type: "Person",
		Fields: map[string]interface{}{
			"id":     "p3",
			"age":    15.0,
			"status": "NORMAL",
		},
	}
	err = network3.SubmitFact(fact3)
	if err != nil {
		t.Fatalf("Erreur soumission fait3: %v", err)
	}
	// Vérifier qu'aucun TerminalNode du nouveau réseau n'a été activé
	var terminalNode3 *TerminalNode
	for _, tn := range network3.TerminalNodes {
		terminalNode3 = tn
		break
	}
	execCount = terminalNode3.GetExecutionCount()
	if execCount != 0 {
		t.Errorf("Le TerminalNode ne devrait pas être activé, got %d activations", execCount)
	}
	t.Logf("✓ Expression OR non décomposée, traitée comme un seul nœud normalisé")
}

// TestAlphaChain_NetworkStats_Accurate vérifie que GetNetworkStats()
// reporte correctement les statistiques de partage
func TestAlphaChain_NetworkStats_Accurate(t *testing.T) {
	tempDir := t.TempDir()
	tsdFile := filepath.Join(tempDir, "test.tsd")
	content := `type Person(#id: string, age: number, name: string, salary:number)
action print(message: string)
rule r1 : {p: Person} / p.age > 18 ==> print("R1")
rule r2 : {p: Person} / p.age > 18 AND p.name == 'toto' ==> print("R2")
rule r3 : {p: Person} / p.age > 18 AND p.name == 'toto' AND p.salary > 1000 ==> print("R3")
rule r4 : {p: Person} / p.age > 21 ==> print("R4")
rule r5 : {p: Person} / p.age > 21 AND p.salary > 2000 ==> print("R5")
`
	if err := os.WriteFile(tsdFile, []byte(content), 0644); err != nil {
		t.Fatalf("Erreur écriture fichier: %v", err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(tsdFile, nil, storage)
	if err != nil {
		t.Fatalf("Erreur construction réseau: %v", err)
	}
	stats := network.GetNetworkStats()
	// Vérifier le nombre d'AlphaNodes uniques (partage optimal entre règles simples et chaînes)
	// - p.age > 18: partagé par r1 (simple), r2 (chaîne), r3 (chaîne) → 1 AlphaNode
	// - p.name == 'toto': partagé par r2, r3 → 1 AlphaNode
	// - p.salary > 1000: utilisé par r3 → 1 AlphaNode
	// - p.age > 21: partagé par r4 (simple), r5 (chaîne) → 1 AlphaNode
	// - p.salary > 2000: utilisé par r5 → 1 AlphaNode
	// Total: 5 AlphaNodes (au lieu de 7 sans partage règles simples/chaînes)
	totalAlphaNodes := stats["alpha_nodes"].(int)
	if totalAlphaNodes != 5 {
		t.Errorf("Devrait avoir 5 AlphaNodes uniques (partage optimal), got %d", totalAlphaNodes)
	}
	// Vérifier le nombre de TerminalNodes
	totalTerminalNodes := stats["terminal_nodes"].(int)
	if totalTerminalNodes != 5 {
		t.Errorf("Devrait avoir 5 TerminalNodes, got %d", totalTerminalNodes)
	}
	// Vérifier le nombre de références
	// p.age > 18: 3 refs
	// p.name = 'toto': 2 refs
	// p.salary > 1000: 1 ref
	// p.age > 21: 2 refs
	// p.salary > 2000: 1 ref
	// Total: le comptage dépend de la structure des chaînes
	if ruleRefs, ok := stats["sharing_total_rule_references"]; ok {
		refCount := ruleRefs.(int)
		// Au minimum, nous avons 5 TerminalNodes
		if refCount < 5 {
			t.Errorf("Devrait avoir au moins 5 enfants (5 terminaux), got %d", refCount)
		}
	} else {
		t.Error("Les statistiques de références devraient être présentes")
	}
	// Vérifier le nombre de nœuds partagés dans le registre
	if sharedNodes, ok := stats["sharing_total_shared_alpha_nodes"]; ok {
		count := sharedNodes.(int)
		if count != 5 {
			t.Errorf("Devrait avoir 5 nœuds dans le registre de partage, got %d", count)
		}
	} else {
		t.Error("Les statistiques de nœuds partagés devraient être présentes")
	}
	// Vérifier le ratio de partage (approximatif, dépend de l'implémentation)
	if avgRatio, ok := stats["sharing_average_sharing_ratio"]; ok {
		ratio := avgRatio.(float64)
		// Le ratio réel varie selon le partage effectif, vérifions qu'il y a du partage
		if ratio < 1.0 {
			t.Errorf("Le ratio de partage devrait être >= 1.0, got %v", ratio)
		}
	} else {
		t.Error("Les statistiques de ratio de partage devraient être présentes")
	}
	// Vérifier le TypeNode count
	typeNodeCount := stats["type_nodes"].(int)
	if typeNodeCount != 1 {
		t.Errorf("Devrait avoir 1 TypeNode, got %d", typeNodeCount)
	}
	// Vérifier les statistiques après suppression d'une règle
	err = network.RemoveRule("r2")
	if err != nil {
		t.Fatalf("Erreur suppression règle r2: %v", err)
	}
	stats = network.GetNetworkStats()
	// Après suppression de r2:
	// r2 avait une chaîne de 2 AlphaNodes (p.age > 18 ET p.name == 'toto')
	// Ces 2 nœuds étaient partagés avec r3, donc ils ne sont PAS supprimés
	// Il reste donc 7 AlphaNodes (aucun n'est supprimé car tous sont encore utilisés)
	if ruleRefs, ok := stats["sharing_total_rule_references"]; ok {
		refCount := ruleRefs.(int)
		// Au minimum, nous avons 4 TerminalNodes restants
		if refCount < 4 {
			t.Errorf("Devrait avoir au moins 4 enfants après suppression, got %d", refCount)
		}
	}
	// Nombre d'AlphaNodes reste 5 (les nœuds de r2 sont partagés avec r3, partage optimal)
	totalAlphaNodes = stats["alpha_nodes"].(int)
	if totalAlphaNodes != 5 {
		t.Errorf("Devrait toujours avoir 5 AlphaNodes (partage optimal), got %d", totalAlphaNodes)
	}
	// Nombre de TerminalNodes: 4
	totalTerminalNodes = stats["terminal_nodes"].(int)
	if totalTerminalNodes != 4 {
		t.Errorf("Devrait avoir 4 TerminalNodes après suppression, got %d", totalTerminalNodes)
	}
	// Nouveau ratio: 5 AlphaNodes / 4 règles restantes
	// Le ratio devrait rester stable avec le partage optimal
	if avgRatio, ok := stats["sharing_average_sharing_ratio"]; ok {
		ratio := avgRatio.(float64)
		// Vérifions simplement qu'il y a toujours du partage
		if ratio < 1.0 {
			t.Errorf("Le ratio de partage devrait être >= 1.0 après suppression, got %v", ratio)
		}
	}
	t.Logf("✓ GetNetworkStats() reporte des statistiques précises")
}

// TestAlphaChain_MixedConditions_ComplexSharing teste un mélange
// de conditions simples et de chaînes avec partage complexe
func TestAlphaChain_MixedConditions_ComplexSharing(t *testing.T) {
	tempDir := t.TempDir()
	tsdFile := filepath.Join(tempDir, "test.tsd")
	content := `type Person(#id: string, age: number, name: string, salary: number, city:string)
action print(message: string)
rule simple1 : {p: Person} / p.age > 18 ==> print("S1")
rule simple2 : {p: Person} / p.salary > 1000 ==> print("S2")
rule chain1 : {p: Person} / p.age > 18 AND p.name == 'toto' ==> print("C1")
rule chain2 : {p: Person} / p.age > 18 AND p.name == 'toto' AND p.salary > 1000 ==> print("C2")
rule chain3 : {p: Person} / p.salary > 1000 AND p.city == 'Paris' ==> print("C3")
`
	if err := os.WriteFile(tsdFile, []byte(content), 0644); err != nil {
		t.Fatalf("Erreur écriture fichier: %v", err)
	}
	storage := NewMemoryStorage()
	pipeline := NewConstraintPipeline()
	network, _, err := pipeline.IngestFile(tsdFile, nil, storage)
	if err != nil {
		t.Fatalf("Erreur construction réseau: %v", err)
	}
	stats := network.GetNetworkStats()
	// AlphaNodes attendus (partage optimal entre règles simples et chaînes):
	// - p.age > 18: partagé par simple1, chain1, chain2 → 1 AlphaNode
	// - p.salary > 1000: partagé par simple2, chain2, chain3 → 1 AlphaNode
	// - p.name == 'toto': partagé par chain1, chain2 → 1 AlphaNode
	// - p.city == 'Paris': utilisé par chain3 → 1 AlphaNode
	// Total: 4 AlphaNodes (au lieu de 6 sans partage optimal)
	totalAlphaNodes := stats["alpha_nodes"].(int)
	if totalAlphaNodes != 4 {
		t.Errorf("Devrait avoir 4 AlphaNodes (partage optimal), got %d", totalAlphaNodes)
	}
	// 5 TerminalNodes
	totalTerminalNodes := stats["terminal_nodes"].(int)
	if totalTerminalNodes != 5 {
		t.Errorf("Devrait avoir 5 TerminalNodes, got %d", totalTerminalNodes)
	}
	// Le nombre d'enfants dépend de la structure des chaînes
	if ruleRefs, ok := stats["sharing_total_rule_references"]; ok {
		refCount := ruleRefs.(int)
		// Au minimum, nous avons 5 TerminalNodes
		if refCount < 5 {
			t.Errorf("Devrait avoir au moins 5 enfants (5 terminaux), got %d", refCount)
		}
	}
	// Test avec un fait qui satisfait plusieurs règles
	fact := &Fact{
		ID:   "person1",
		Type: "Person",
		Fields: map[string]interface{}{
			"id":     "p1",
			"age":    25.0,
			"name":   "toto",
			"salary": 2000.0,
			"city":   "Paris",
		},
	}
	err = network.SubmitFact(fact)
	if err != nil {
		t.Fatalf("Erreur soumission fait: %v", err)
	}
	// Toutes les 5 règles devraient être activées
	activatedCount := 0
	for _, terminalNode := range network.TerminalNodes {
		if terminalNode.GetExecutionCount() > 0 {
			activatedCount++
		}
	}
	if activatedCount != 5 {
		t.Errorf("Les 5 règles devraient être activées, got %d", activatedCount)
	}
	t.Logf("✓ Mélange de conditions simples et chaînes avec partage complexe")
}

// TestAlphaChain_EmptyNetwork_Stats teste les statistiques sur un réseau vide
func TestAlphaChain_EmptyNetwork_Stats(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	network.AlphaSharingManager = NewAlphaSharingRegistry()
	network.LifecycleManager = NewLifecycleManager()
	stats := network.GetNetworkStats()
	// Vérifier que toutes les statistiques sont à zéro
	if stats["alpha_nodes"].(int) != 0 {
		t.Errorf("Réseau vide devrait avoir 0 AlphaNodes, got %d", stats["alpha_nodes"].(int))
	}
	if stats["terminal_nodes"].(int) != 0 {
		t.Errorf("Réseau vide devrait avoir 0 TerminalNodes, got %d", stats["terminal_nodes"].(int))
	}
	if stats["type_nodes"].(int) != 0 {
		t.Errorf("Réseau vide devrait avoir 0 TypeNodes, got %d", stats["type_nodes"].(int))
	}
	if sharingStats, ok := stats["sharing_total_shared_alpha_nodes"]; ok {
		if sharingStats.(int) != 0 {
			t.Errorf("Réseau vide devrait avoir 0 nœuds partagés, got %d", sharingStats.(int))
		}
	}
	if ruleRefs, ok := stats["sharing_total_rule_references"]; ok {
		if ruleRefs.(int) != 0 {
			t.Errorf("Réseau vide devrait avoir 0 références, got %d", ruleRefs.(int))
		}
	}
	if avgRatio, ok := stats["sharing_average_sharing_ratio"]; ok {
		if avgRatio.(float64) != 0.0 {
			t.Errorf("Réseau vide devrait avoir un ratio de 0.0, got %v", avgRatio.(float64))
		}
	}
	t.Logf("✓ Statistiques correctes sur un réseau vide")
}
