// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
package rete

import (
	"testing"
)

// TestBuildChain_SingleCondition teste la construction d'une chaîne avec une seule condition
func TestBuildChain_SingleCondition(t *testing.T) {
	// Créer un réseau RETE avec storage
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	// Créer un nœud parent (TypeNode)
	typeDef := TypeDefinition{Type: "type", Name: "Person", Fields: []Field{}}
	parentNode := NewTypeNode("person", typeDef, storage)
	network.TypeNodes["Person"] = parentNode
	// Créer le builder
	builder := NewAlphaChainBuilder(network, storage)
	// Créer une condition simple
	conditions := []SimpleCondition{
		NewSimpleCondition("comparison", "p.age", ">", 18),
	}
	// Construire la chaîne
	chain, err := builder.BuildChain(conditions, "p", parentNode, "rule1")
	if err != nil {
		t.Fatalf("Erreur lors de la construction de la chaîne: %v", err)
	}
	// Vérifications
	if chain == nil {
		t.Fatal("La chaîne ne devrait pas être nil")
	}
	if len(chain.Nodes) != 1 {
		t.Errorf("Attendu 1 nœud, obtenu %d", len(chain.Nodes))
	}
	if len(chain.Hashes) != 1 {
		t.Errorf("Attendu 1 hash, obtenu %d", len(chain.Hashes))
	}
	if chain.FinalNode == nil {
		t.Fatal("Le nœud final ne devrait pas être nil")
	}
	if chain.FinalNode != chain.Nodes[0] {
		t.Error("Le nœud final devrait être le premier nœud")
	}
	if chain.RuleID != "rule1" {
		t.Errorf("Attendu RuleID 'rule1', obtenu '%s'", chain.RuleID)
	}
	// Vérifier que le nœud est dans le réseau
	if _, exists := network.AlphaNodes[chain.Nodes[0].ID]; !exists {
		t.Error("Le nœud alpha devrait être dans le réseau")
	}
	// Vérifier l'enregistrement dans le LifecycleManager
	lifecycle, exists := network.LifecycleManager.GetNodeLifecycle(chain.Nodes[0].ID)
	if !exists {
		t.Fatal("Le nœud devrait être enregistré dans le LifecycleManager")
	}
	if lifecycle.GetRefCount() != 1 {
		t.Errorf("Attendu 1 référence, obtenu %d", lifecycle.GetRefCount())
	}
}

// TestBuildChain_TwoConditions_New teste la construction d'une chaîne avec deux conditions nouvelles
func TestBuildChain_TwoConditions_New(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	typeDef := TypeDefinition{Type: "type", Name: "Person", Fields: []Field{}}
	parentNode := NewTypeNode("person", typeDef, storage)
	network.TypeNodes["Person"] = parentNode
	builder := NewAlphaChainBuilder(network, storage)
	conditions := []SimpleCondition{
		NewSimpleCondition("comparison", "p.age", ">", 18),
		NewSimpleCondition("comparison", "p.name", "==", "Alice"),
	}
	chain, err := builder.BuildChain(conditions, "p", parentNode, "rule1")
	if err != nil {
		t.Fatalf("Erreur lors de la construction de la chaîne: %v", err)
	}
	// Vérifications
	if len(chain.Nodes) != 2 {
		t.Errorf("Attendu 2 nœuds, obtenu %d", len(chain.Nodes))
	}
	if len(chain.Hashes) != 2 {
		t.Errorf("Attendu 2 hashes, obtenu %d", len(chain.Hashes))
	}
	if chain.FinalNode != chain.Nodes[1] {
		t.Error("Le nœud final devrait être le deuxième nœud")
	}
	// Vérifier que les deux nœuds sont dans le réseau
	for i, node := range chain.Nodes {
		if _, exists := network.AlphaNodes[node.ID]; !exists {
			t.Errorf("Le nœud %d devrait être dans le réseau", i)
		}
	}
	// Vérifier la connexion en chaîne
	// Le premier nœud devrait être enfant du parent
	if !isAlreadyConnected(parentNode, chain.Nodes[0]) {
		t.Error("Le premier nœud devrait être connecté au parent")
	}
	// Le deuxième nœud devrait être enfant du premier
	if !isAlreadyConnected(chain.Nodes[0], chain.Nodes[1]) {
		t.Error("Le deuxième nœud devrait être connecté au premier")
	}
	// Vérifier les références dans le LifecycleManager
	for i, node := range chain.Nodes {
		lifecycle, exists := network.LifecycleManager.GetNodeLifecycle(node.ID)
		if !exists {
			t.Errorf("Le nœud %d devrait être enregistré dans le LifecycleManager", i)
			continue
		}
		if lifecycle.GetRefCount() != 1 {
			t.Errorf("Nœud %d: attendu 1 référence, obtenu %d", i, lifecycle.GetRefCount())
		}
	}
}

// TestBuildChain_TwoConditions_Reuse teste la réutilisation complète d'une chaîne
func TestBuildChain_TwoConditions_Reuse(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	typeDef := TypeDefinition{Type: "type", Name: "Person", Fields: []Field{}}
	parentNode := NewTypeNode("person", typeDef, storage)
	network.TypeNodes["Person"] = parentNode
	builder := NewAlphaChainBuilder(network, storage)
	conditions := []SimpleCondition{
		NewSimpleCondition("comparison", "p.age", ">", 18),
		NewSimpleCondition("comparison", "p.name", "==", "Alice"),
	}
	// Construire la première chaîne
	chain1, err := builder.BuildChain(conditions, "p", parentNode, "rule1")
	if err != nil {
		t.Fatalf("Erreur lors de la construction de la première chaîne: %v", err)
	}
	initialNodeCount := len(network.AlphaNodes)
	// Construire la deuxième chaîne avec les mêmes conditions
	chain2, err := builder.BuildChain(conditions, "p", parentNode, "rule2")
	if err != nil {
		t.Fatalf("Erreur lors de la construction de la deuxième chaîne: %v", err)
	}
	// Vérifications
	if len(chain2.Nodes) != 2 {
		t.Errorf("Attendu 2 nœuds dans la deuxième chaîne, obtenu %d", len(chain2.Nodes))
	}
	// Le nombre de nœuds dans le réseau ne devrait pas avoir changé (réutilisation)
	if len(network.AlphaNodes) != initialNodeCount {
		t.Errorf("Attendu %d nœuds dans le réseau (réutilisation), obtenu %d",
			initialNodeCount, len(network.AlphaNodes))
	}
	// Les nœuds devraient être identiques
	for i := 0; i < len(chain1.Nodes); i++ {
		if chain1.Nodes[i].ID != chain2.Nodes[i].ID {
			t.Errorf("Nœud %d: les IDs devraient être identiques (réutilisation)", i)
		}
		if chain1.Hashes[i] != chain2.Hashes[i] {
			t.Errorf("Nœud %d: les hashes devraient être identiques", i)
		}
	}
	// Vérifier que chaque nœud a maintenant 2 références
	for i, node := range chain1.Nodes {
		lifecycle, exists := network.LifecycleManager.GetNodeLifecycle(node.ID)
		if !exists {
			t.Errorf("Le nœud %d devrait être enregistré dans le LifecycleManager", i)
			continue
		}
		if lifecycle.GetRefCount() != 2 {
			t.Errorf("Nœud %d: attendu 2 références (rule1 + rule2), obtenu %d",
				i, lifecycle.GetRefCount())
		}
	}
}

// TestBuildChain_PartialReuse teste le partage partiel d'une chaîne
func TestBuildChain_PartialReuse(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	typeDef := TypeDefinition{Type: "type", Name: "Person", Fields: []Field{}}
	parentNode := NewTypeNode("person", typeDef, storage)
	network.TypeNodes["Person"] = parentNode
	builder := NewAlphaChainBuilder(network, storage)
	// Première règle: 2 conditions
	conditions1 := []SimpleCondition{
		NewSimpleCondition("comparison", "p.age", ">", 18),
		NewSimpleCondition("comparison", "p.name", "==", "Alice"),
	}
	chain1, err := builder.BuildChain(conditions1, "p", parentNode, "rule1")
	if err != nil {
		t.Fatalf("Erreur lors de la construction de la première chaîne: %v", err)
	}
	initialNodeCount := len(network.AlphaNodes)
	// Deuxième règle: partage la première condition, mais différente deuxième condition
	conditions2 := []SimpleCondition{
		NewSimpleCondition("comparison", "p.age", ">", 18),        // Même que rule1
		NewSimpleCondition("comparison", "p.city", "==", "Paris"), // Différent de rule1
	}
	chain2, err := builder.BuildChain(conditions2, "p", parentNode, "rule2")
	if err != nil {
		t.Fatalf("Erreur lors de la construction de la deuxième chaîne: %v", err)
	}
	// Vérifications
	if len(chain2.Nodes) != 2 {
		t.Errorf("Attendu 2 nœuds dans la deuxième chaîne, obtenu %d", len(chain2.Nodes))
	}
	// Un seul nouveau nœud devrait avoir été ajouté (la deuxième condition)
	expectedNodeCount := initialNodeCount + 1
	if len(network.AlphaNodes) != expectedNodeCount {
		t.Errorf("Attendu %d nœuds dans le réseau (1 nouveau), obtenu %d",
			expectedNodeCount, len(network.AlphaNodes))
	}
	// Le premier nœud devrait être partagé
	if chain1.Nodes[0].ID != chain2.Nodes[0].ID {
		t.Error("Le premier nœud devrait être partagé entre les deux chaînes")
	}
	// Le deuxième nœud devrait être différent
	if chain1.Nodes[1].ID == chain2.Nodes[1].ID {
		t.Error("Le deuxième nœud devrait être différent entre les deux chaînes")
	}
	// Vérifier les compteurs de références
	lifecycle0, _ := network.LifecycleManager.GetNodeLifecycle(chain2.Nodes[0].ID)
	if lifecycle0.GetRefCount() != 2 {
		t.Errorf("Premier nœud: attendu 2 références, obtenu %d", lifecycle0.GetRefCount())
	}
	lifecycle1Rule1, _ := network.LifecycleManager.GetNodeLifecycle(chain1.Nodes[1].ID)
	if lifecycle1Rule1.GetRefCount() != 1 {
		t.Errorf("Deuxième nœud de rule1: attendu 1 référence, obtenu %d",
			lifecycle1Rule1.GetRefCount())
	}
	lifecycle1Rule2, _ := network.LifecycleManager.GetNodeLifecycle(chain2.Nodes[1].ID)
	if lifecycle1Rule2.GetRefCount() != 1 {
		t.Errorf("Deuxième nœud de rule2: attendu 1 référence, obtenu %d",
			lifecycle1Rule2.GetRefCount())
	}
}

// TestBuildChain_CompleteReuse teste la réutilisation complète par plusieurs règles
func TestBuildChain_CompleteReuse(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	typeDef := TypeDefinition{Type: "type", Name: "Person", Fields: []Field{}}
	parentNode := NewTypeNode("person", typeDef, storage)
	network.TypeNodes["Person"] = parentNode
	builder := NewAlphaChainBuilder(network, storage)
	conditions := []SimpleCondition{
		NewSimpleCondition("comparison", "p.age", ">", 18),
	}
	// Créer 5 règles avec la même condition
	ruleCount := 5
	chains := make([]*AlphaChain, ruleCount)
	for i := 0; i < ruleCount; i++ {
		ruleID := "rule" + string(rune('1'+i))
		chain, err := builder.BuildChain(conditions, "p", parentNode, ruleID)
		if err != nil {
			t.Fatalf("Erreur lors de la construction de la chaîne %d: %v", i, err)
		}
		chains[i] = chain
	}
	// Vérifications
	// Il ne devrait y avoir qu'un seul nœud alpha dans le réseau
	if len(network.AlphaNodes) != 1 {
		t.Errorf("Attendu 1 nœud alpha dans le réseau (réutilisation complète), obtenu %d",
			len(network.AlphaNodes))
	}
	// Toutes les chaînes devraient pointer vers le même nœud
	firstNodeID := chains[0].Nodes[0].ID
	for i := 1; i < ruleCount; i++ {
		if chains[i].Nodes[0].ID != firstNodeID {
			t.Errorf("Chaîne %d: le nœud devrait être le même que la première chaîne", i)
		}
	}
	// Le compteur de références devrait être égal au nombre de règles
	lifecycle, exists := network.LifecycleManager.GetNodeLifecycle(firstNodeID)
	if !exists {
		t.Fatal("Le nœud devrait être enregistré dans le LifecycleManager")
	}
	if lifecycle.GetRefCount() != ruleCount {
		t.Errorf("Attendu %d références, obtenu %d", ruleCount, lifecycle.GetRefCount())
	}
}

// TestBuildChain_MultipleRules_SharedSubchain teste plusieurs règles avec sous-chaînes partagées
func TestBuildChain_MultipleRules_SharedSubchain(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	typeDef := TypeDefinition{Type: "type", Name: "Person", Fields: []Field{}}
	parentNode := NewTypeNode("person", typeDef, storage)
	network.TypeNodes["Person"] = parentNode
	builder := NewAlphaChainBuilder(network, storage)
	// Rule 1: age > 18 AND name == "Alice"
	conditions1 := []SimpleCondition{
		NewSimpleCondition("comparison", "p.age", ">", 18),
		NewSimpleCondition("comparison", "p.name", "==", "Alice"),
	}
	// Rule 2: age > 18 AND city == "Paris"
	conditions2 := []SimpleCondition{
		NewSimpleCondition("comparison", "p.age", ">", 18),
		NewSimpleCondition("comparison", "p.city", "==", "Paris"),
	}
	// Rule 3: age > 18 AND name == "Alice" (identique à Rule 1)
	conditions3 := conditions1
	// Construire les trois chaînes
	chain1, err := builder.BuildChain(conditions1, "p", parentNode, "rule1")
	if err != nil {
		t.Fatalf("Erreur construction chain1: %v", err)
	}
	chain2, err := builder.BuildChain(conditions2, "p", parentNode, "rule2")
	if err != nil {
		t.Fatalf("Erreur construction chain2: %v", err)
	}
	chain3, err := builder.BuildChain(conditions3, "p", parentNode, "rule3")
	if err != nil {
		t.Fatalf("Erreur construction chain3: %v", err)
	}
	// Vérifications
	// Il devrait y avoir 3 nœuds alpha uniques:
	// - age > 18 (partagé par les 3 règles)
	// - name == "Alice" (partagé par rule1 et rule3)
	// - city == "Paris" (utilisé seulement par rule2)
	expectedNodeCount := 3
	if len(network.AlphaNodes) != expectedNodeCount {
		t.Errorf("Attendu %d nœuds alpha, obtenu %d", expectedNodeCount, len(network.AlphaNodes))
	}
	// Le premier nœud devrait être partagé entre toutes les chaînes
	if chain1.Nodes[0].ID != chain2.Nodes[0].ID || chain1.Nodes[0].ID != chain3.Nodes[0].ID {
		t.Error("Le premier nœud devrait être partagé entre toutes les chaînes")
	}
	// Chain1 et Chain3 devraient avoir exactement les mêmes nœuds
	if len(chain1.Nodes) != len(chain3.Nodes) {
		t.Error("Chain1 et Chain3 devraient avoir le même nombre de nœuds")
	}
	for i := 0; i < len(chain1.Nodes); i++ {
		if chain1.Nodes[i].ID != chain3.Nodes[i].ID {
			t.Errorf("Nœud %d: Chain1 et Chain3 devraient partager le même nœud", i)
		}
	}
	// Vérifier les compteurs de références
	lifecycle0, _ := network.LifecycleManager.GetNodeLifecycle(chain1.Nodes[0].ID)
	if lifecycle0.GetRefCount() != 3 {
		t.Errorf("Premier nœud: attendu 3 références, obtenu %d", lifecycle0.GetRefCount())
	}
	lifecycle1, _ := network.LifecycleManager.GetNodeLifecycle(chain1.Nodes[1].ID)
	if lifecycle1.GetRefCount() != 2 {
		t.Errorf("Deuxième nœud de chain1: attendu 2 références (rule1+rule3), obtenu %d",
			lifecycle1.GetRefCount())
	}
	lifecycle2, _ := network.LifecycleManager.GetNodeLifecycle(chain2.Nodes[1].ID)
	if lifecycle2.GetRefCount() != 1 {
		t.Errorf("Deuxième nœud de chain2: attendu 1 référence, obtenu %d",
			lifecycle2.GetRefCount())
	}
}

// TestBuildChain_EmptyConditions teste le comportement avec une liste vide
func TestBuildChain_EmptyConditions(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	typeDef := TypeDefinition{Type: "type", Name: "Person", Fields: []Field{}}
	parentNode := NewTypeNode("person", typeDef, storage)
	builder := NewAlphaChainBuilder(network, storage)
	conditions := []SimpleCondition{}
	_, err := builder.BuildChain(conditions, "p", parentNode, "rule1")
	if err == nil {
		t.Error("Attendu une erreur avec une liste de conditions vide")
	}
}

// TestBuildChain_NilParent teste le comportement avec un parent nil
func TestBuildChain_NilParent(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	builder := NewAlphaChainBuilder(network, storage)
	conditions := []SimpleCondition{
		NewSimpleCondition("comparison", "p.age", ">", 18),
	}
	_, err := builder.BuildChain(conditions, "p", nil, "rule1")
	if err == nil {
		t.Error("Attendu une erreur avec un parent nil")
	}
}

// TestAlphaChain_ValidateChain teste la validation d'une chaîne
func TestAlphaChain_ValidateChain(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	typeDef := TypeDefinition{Type: "type", Name: "Person", Fields: []Field{}}
	parentNode := NewTypeNode("person", typeDef, storage)
	builder := NewAlphaChainBuilder(network, storage)
	conditions := []SimpleCondition{
		NewSimpleCondition("comparison", "p.age", ">", 18),
		NewSimpleCondition("comparison", "p.name", "==", "Alice"),
	}
	chain, err := builder.BuildChain(conditions, "p", parentNode, "rule1")
	if err != nil {
		t.Fatalf("Erreur construction chaîne: %v", err)
	}
	// La chaîne devrait être valide
	if err := chain.ValidateChain(); err != nil {
		t.Errorf("La chaîne devrait être valide: %v", err)
	}
}

// TestAlphaChain_GetChainInfo teste la récupération d'informations sur la chaîne
func TestAlphaChain_GetChainInfo(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	typeDef := TypeDefinition{Type: "type", Name: "Person", Fields: []Field{}}
	parentNode := NewTypeNode("person", typeDef, storage)
	builder := NewAlphaChainBuilder(network, storage)
	conditions := []SimpleCondition{
		NewSimpleCondition("comparison", "p.age", ">", 18),
	}
	chain, err := builder.BuildChain(conditions, "p", parentNode, "rule1")
	if err != nil {
		t.Fatalf("Erreur construction chaîne: %v", err)
	}
	info := chain.GetChainInfo()
	if info == nil {
		t.Fatal("GetChainInfo ne devrait pas retourner nil")
	}
	if info["rule_id"] != "rule1" {
		t.Errorf("Attendu rule_id 'rule1', obtenu '%v'", info["rule_id"])
	}
	if info["node_count"] != 1 {
		t.Errorf("Attendu node_count 1, obtenu %v", info["node_count"])
	}
	if info["final_node_id"] != chain.FinalNode.ID {
		t.Error("final_node_id incorrect")
	}
}

// TestAlphaChainBuilder_CountSharedNodes teste le comptage des nœuds partagés
func TestAlphaChainBuilder_CountSharedNodes(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	typeDef := TypeDefinition{Type: "type", Name: "Person", Fields: []Field{}}
	parentNode := NewTypeNode("person", typeDef, storage)
	builder := NewAlphaChainBuilder(network, storage)
	conditions := []SimpleCondition{
		NewSimpleCondition("comparison", "p.age", ">", 18),
		NewSimpleCondition("comparison", "p.name", "==", "Alice"),
	}
	// Première chaîne - aucun nœud partagé
	chain1, _ := builder.BuildChain(conditions, "p", parentNode, "rule1")
	sharedCount1 := builder.CountSharedNodes(chain1)
	if sharedCount1 != 0 {
		t.Errorf("Première chaîne: attendu 0 nœud partagé, obtenu %d", sharedCount1)
	}
	// Deuxième chaîne - tous les nœuds sont partagés
	chain2, _ := builder.BuildChain(conditions, "p", parentNode, "rule2")
	sharedCount2 := builder.CountSharedNodes(chain2)
	if sharedCount2 != 2 {
		t.Errorf("Deuxième chaîne: attendu 2 nœuds partagés, obtenu %d", sharedCount2)
	}
}

// TestAlphaChainBuilder_GetChainStats teste la récupération de statistiques détaillées
func TestAlphaChainBuilder_GetChainStats(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	typeDef := TypeDefinition{Type: "type", Name: "Person", Fields: []Field{}}
	parentNode := NewTypeNode("person", typeDef, storage)
	builder := NewAlphaChainBuilder(network, storage)
	conditions := []SimpleCondition{
		NewSimpleCondition("comparison", "p.age", ">", 18),
		NewSimpleCondition("comparison", "p.name", "==", "Alice"),
	}
	chain1, _ := builder.BuildChain(conditions, "p", parentNode, "rule1")
	chain2, _ := builder.BuildChain(conditions, "p", parentNode, "rule2")
	stats := builder.GetChainStats(chain2)
	if stats["total_nodes"] != 2 {
		t.Errorf("Attendu total_nodes=2, obtenu %v", stats["total_nodes"])
	}
	if stats["shared_nodes"] != 2 {
		t.Errorf("Attendu shared_nodes=2, obtenu %v", stats["shared_nodes"])
	}
	if stats["new_nodes"] != 0 {
		t.Errorf("Attendu new_nodes=0, obtenu %v", stats["new_nodes"])
	}
	if stats["rule_id"] != "rule2" {
		t.Errorf("Attendu rule_id='rule2', obtenu '%v'", stats["rule_id"])
	}
	// Vérifier les détails des nœuds
	nodeDetails, ok := stats["node_details"].([]map[string]interface{})
	if !ok {
		t.Fatal("node_details devrait être un slice de maps")
	}
	if len(nodeDetails) != 2 {
		t.Errorf("Attendu 2 entrées dans node_details, obtenu %d", len(nodeDetails))
	}
	// Tous les nœuds devraient avoir is_shared=true et ref_count=2
	for i, detail := range nodeDetails {
		if detail["is_shared"] != true {
			t.Errorf("Nœud %d: is_shared devrait être true", i)
		}
		if detail["ref_count"] != 2 {
			t.Errorf("Nœud %d: ref_count devrait être 2, obtenu %v", i, detail["ref_count"])
		}
	}
	// Le dernier nœud devrait avoir is_final=true
	if nodeDetails[1]["is_final"] != true {
		t.Error("Le dernier nœud devrait avoir is_final=true")
	}
	// Ignorer chain1 pour éviter unused variable error
	_ = chain1
}

// TestIsAlreadyConnected teste la fonction helper isAlreadyConnected
func TestIsAlreadyConnected(t *testing.T) {
	storage := NewMemoryStorage()
	typeDef := TypeDefinition{Type: "type", Name: "Person", Fields: []Field{}}
	parent := NewTypeNode("person", typeDef, storage)
	child1 := NewAlphaNode("alpha1", nil, "p", storage)
	child2 := NewAlphaNode("alpha2", nil, "p", storage)
	// Initialement, child1 n'est pas connecté
	if isAlreadyConnected(parent, child1) {
		t.Error("child1 ne devrait pas être connecté initialement")
	}
	// Connecter child1
	parent.AddChild(child1)
	// Maintenant child1 devrait être connecté
	if !isAlreadyConnected(parent, child1) {
		t.Error("child1 devrait être connecté après AddChild")
	}
	// child2 ne devrait toujours pas être connecté
	if isAlreadyConnected(parent, child2) {
		t.Error("child2 ne devrait pas être connecté")
	}
	// Tester avec parent nil
	if isAlreadyConnected(nil, child1) {
		t.Error("parent nil devrait retourner false")
	}
	// Tester avec child nil
	if isAlreadyConnected(parent, nil) {
		t.Error("child nil devrait retourner false")
	}
}

// TestAlphaChainBuilder_BuildDecomposedChain tests decomposed chain building with metadata
func TestAlphaChainBuilder_BuildDecomposedChain(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	builder := NewAlphaChainBuilder(network, storage)
	// Create a type node as parent
	typeDef := TypeDefinition{Fields: []Field{}}
	typeNode := NewTypeNode("Person", typeDef, storage)
	network.TypeNodes["Person"] = typeNode
	// Create decomposed conditions with metadata
	conditions := []DecomposedCondition{
		{
			SimpleCondition: NewSimpleCondition(
				"binaryOp",
				map[string]interface{}{"type": "fieldAccess", "field": "qte"},
				"*",
				map[string]interface{}{"type": "number", "value": 23},
			),
			ResultName:   "temp_1",
			IsAtomic:     true,
			Dependencies: []string{},
		},
		{
			SimpleCondition: NewSimpleCondition(
				"binaryOp",
				map[string]interface{}{"type": "tempResult", "step_name": "temp_1"},
				"-",
				map[string]interface{}{"type": "number", "value": 10},
			),
			ResultName:   "temp_2",
			IsAtomic:     true,
			Dependencies: []string{"temp_1"},
		},
		{
			SimpleCondition: NewSimpleCondition(
				"comparison",
				map[string]interface{}{"type": "tempResult", "step_name": "temp_2"},
				">",
				map[string]interface{}{"type": "number", "value": 0},
			),
			ResultName:   "temp_3",
			IsAtomic:     true,
			Dependencies: []string{"temp_2"},
		},
	}
	// Build the decomposed chain
	chain, err := builder.BuildDecomposedChain(conditions, "p", typeNode, "test_rule")
	if err != nil {
		t.Fatalf("BuildDecomposedChain failed: %v", err)
	}
	// Verify chain structure
	if len(chain.Nodes) != 3 {
		t.Errorf("Expected 3 nodes, got %d", len(chain.Nodes))
	}
	// Verify first node metadata
	node1 := chain.Nodes[0]
	if node1.ResultName != "temp_1" {
		t.Errorf("Node 1: Expected ResultName 'temp_1', got '%s'", node1.ResultName)
	}
	if !node1.IsAtomic {
		t.Error("Node 1: Expected IsAtomic=true")
	}
	if len(node1.Dependencies) != 0 {
		t.Errorf("Node 1: Expected no dependencies, got %v", node1.Dependencies)
	}
	// Verify second node metadata
	node2 := chain.Nodes[1]
	if node2.ResultName != "temp_2" {
		t.Errorf("Node 2: Expected ResultName 'temp_2', got '%s'", node2.ResultName)
	}
	if !node2.IsAtomic {
		t.Error("Node 2: Expected IsAtomic=true")
	}
	if len(node2.Dependencies) != 1 || node2.Dependencies[0] != "temp_1" {
		t.Errorf("Node 2: Expected Dependencies=['temp_1'], got %v", node2.Dependencies)
	}
	// Verify third node metadata
	node3 := chain.Nodes[2]
	if node3.ResultName != "temp_3" {
		t.Errorf("Node 3: Expected ResultName 'temp_3', got '%s'", node3.ResultName)
	}
	if !node3.IsAtomic {
		t.Error("Node 3: Expected IsAtomic=true")
	}
	if len(node3.Dependencies) != 1 || node3.Dependencies[0] != "temp_2" {
		t.Errorf("Node 3: Expected Dependencies=['temp_2'], got %v", node3.Dependencies)
	}
	// Verify chain is connected
	if chain.FinalNode != node3 {
		t.Error("Expected FinalNode to be the last node")
	}
}

// TestAlphaChainBuilder_DecomposedChainSharing tests node sharing with decomposed chains
func TestAlphaChainBuilder_DecomposedChainSharing(t *testing.T) {
	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)
	builder := NewAlphaChainBuilder(network, storage)
	typeDef := TypeDefinition{Fields: []Field{}}
	typeNode := NewTypeNode("Order", typeDef, storage)
	network.TypeNodes["Order"] = typeNode
	// Condition that will be shared: c.qte * 23
	sharedCondition := DecomposedCondition{
		SimpleCondition: NewSimpleCondition(
			"binaryOp",
			map[string]interface{}{"type": "fieldAccess", "field": "qte"},
			"*",
			map[string]interface{}{"type": "number", "value": 23},
		),
		ResultName:   "temp_1",
		IsAtomic:     true,
		Dependencies: []string{},
	}
	// Rule 1: c.qte * 23 > 100
	conditions1 := []DecomposedCondition{
		sharedCondition,
		{
			SimpleCondition: NewSimpleCondition(
				"comparison",
				map[string]interface{}{"type": "tempResult", "step_name": "temp_1"},
				">",
				map[string]interface{}{"type": "number", "value": 100},
			),
			ResultName:   "temp_2",
			IsAtomic:     true,
			Dependencies: []string{"temp_1"},
		},
	}
	// Rule 2: c.qte * 23 > 50 (shares first step with rule 1)
	conditions2 := []DecomposedCondition{
		sharedCondition,
		{
			SimpleCondition: NewSimpleCondition(
				"comparison",
				map[string]interface{}{"type": "tempResult", "step_name": "temp_1"},
				">",
				map[string]interface{}{"type": "number", "value": 50},
			),
			ResultName:   "temp_2",
			IsAtomic:     true,
			Dependencies: []string{"temp_1"},
		},
	}
	// Build first chain
	chain1, err := builder.BuildDecomposedChain(conditions1, "c", typeNode, "rule_1")
	if err != nil {
		t.Fatalf("BuildDecomposedChain rule_1 failed: %v", err)
	}
	// Build second chain
	chain2, err := builder.BuildDecomposedChain(conditions2, "c", typeNode, "rule_2")
	if err != nil {
		t.Fatalf("BuildDecomposedChain rule_2 failed: %v", err)
	}
	// Verify that first nodes are shared (same ID)
	if chain1.Nodes[0].ID != chain2.Nodes[0].ID {
		t.Errorf("Expected first nodes to be shared, got IDs %s and %s",
			chain1.Nodes[0].ID, chain2.Nodes[0].ID)
	}
	t.Logf("✅ Node sharing verified: %s is shared between rule_1 and rule_2",
		chain1.Nodes[0].ID)
	// Verify that second nodes are different (different comparison values)
	if chain1.Nodes[1].ID == chain2.Nodes[1].ID {
		t.Error("Expected second nodes to be different (different comparison values)")
	}
}
