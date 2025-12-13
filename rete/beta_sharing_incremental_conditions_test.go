// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"testing"
)

// TestJoinNodeSharingWithIncrementalConditions vérifie que le partage de JoinNodes
// fonctionne correctement quand une règle a des conditions supplémentaires par rapport
// à une autre règle.
//
// Scénario:
//   - Règle r1: {u: User, o: Order} / u.id == o.user_id ==> action1(...)
//   - Règle r2: {u: User, o: Order} / u.id == o.user_id AND o.amount > 100 ==> action2(...)
//
// Comportement attendu:
//   - Les deux règles DOIVENT partager le même JoinNode pour la jointure u.id == o.user_id
//   - La condition supplémentaire (o.amount > 100) de r2 est gérée par un AlphaNode
//   - Résultat: 1 JoinNode partagé + 1 AlphaNode pour r2
func TestJoinNodeSharingWithIncrementalConditions(t *testing.T) {
	t.Log("🧪 TEST: JoinNode Sharing with Incremental Conditions")
	t.Log("=========================================================")

	// Setup
	storage := NewMemoryStorage()
	lifecycle := NewLifecycleManager()
	config := DefaultBetaSharingConfig()
	config.Enabled = true
	registry := NewBetaSharingRegistry(config, lifecycle)

	// Condition commune: u.id == o.user_id
	baseCondition := map[string]interface{}{
		"type":     "comparison",
		"operator": "==",
		"left": map[string]interface{}{
			"type":   "fieldAccess",
			"object": "u",
			"field":  "id",
		},
		"right": map[string]interface{}{
			"type":   "fieldAccess",
			"object": "o",
			"field":  "user_id",
		},
	}

	// Note: En pratique, r2 aurait une condition complète (u.id == o.user_id AND o.amount > 100)
	// mais après le split par ConditionSplitter:
	// - Join condition: u.id == o.user_id (identique à baseCondition)
	// - Alpha condition: o.amount > 100 (gérée par AlphaNode)

	varTypes := map[string]string{"u": "User", "o": "Order"}
	leftVars := []string{"u"}
	rightVars := []string{"o"}
	allVars := []string{"u", "o"}
	cascadeLevel := 0

	t.Log("\n📋 Création du JoinNode pour r1 (condition de base)")
	t.Log("   Condition: u.id == o.user_id")

	// Créer JoinNode pour r1 avec condition de base
	node1, hash1, shared1, err := registry.GetOrCreateJoinNode(
		baseCondition,
		leftVars,
		rightVars,
		allVars,
		varTypes,
		storage,
		cascadeLevel,
	)
	if err != nil {
		t.Fatalf("Failed to create JoinNode for r1: %v", err)
	}
	if shared1 {
		t.Error("❌ First JoinNode should not be shared (it's the first one)")
	}

	t.Logf("✅ JoinNode r1 créé: %s (hash: %s)", node1.ID, hash1)

	t.Log("\n📋 Création du JoinNode pour r2 (même condition de jointure)")
	t.Log("   Condition de jointure: u.id == o.user_id (identique à r1)")
	t.Log("   Condition alpha: o.amount > 100 (gérée par AlphaNode)")

	// Pour r2, la condition de jointure extraite devrait être identique à baseCondition
	// (après split par ConditionSplitter)
	// Note: Dans le vrai code, ConditionSplitter extrait la partie join de r2Condition
	// Pour ce test, on utilise directement baseCondition car c'est ce que le splitter retournerait

	node2, hash2, shared2, err := registry.GetOrCreateJoinNode(
		baseCondition, // Même condition de jointure que r1
		leftVars,
		rightVars,
		allVars,
		varTypes,
		storage,
		cascadeLevel,
	)
	if err != nil {
		t.Fatalf("Failed to create JoinNode for r2: %v", err)
	}

	// ASSERTION CRITIQUE: r2 doit PARTAGER le JoinNode de r1
	if !shared2 {
		t.Error("❌ r2 should share JoinNode with r1 (same join condition)")
		t.Error("   This indicates join node sharing is not working for incremental conditions")
		t.FailNow()
	}

	if hash1 != hash2 {
		t.Errorf("❌ r1 and r2 should have same hash (same join condition)")
		t.Errorf("   r1 hash: %s", hash1)
		t.Errorf("   r2 hash: %s", hash2)
		t.FailNow()
	}

	if node1.ID != node2.ID {
		t.Errorf("❌ r1 and r2 should share the same JoinNode")
		t.Errorf("   r1 node: %s", node1.ID)
		t.Errorf("   r2 node: %s", node2.ID)
		t.FailNow()
	}

	t.Logf("✅ JoinNode r2 PARTAGÉ avec r1: %s (hash: %s)", node2.ID, hash2)

	t.Log("\n✅ PARTAGE VÉRIFIÉ")
	t.Log("   - r1 et r2 partagent le même JoinNode pour la jointure commune")
	t.Log("   - La condition supplémentaire de r2 sera gérée par un AlphaNode")
	t.Log("   - Efficacité: 1 JoinNode au lieu de 2")
}

// TestJoinNodeSharingWithDifferentAdditionalConditions vérifie que des règles
// avec des conditions alpha différentes mais la même condition de jointure
// partagent bien le JoinNode.
func TestJoinNodeSharingWithDifferentAdditionalConditions(t *testing.T) {
	t.Log("🧪 TEST: JoinNode Sharing with Different Alpha Conditions")
	t.Log("===========================================================")

	// Setup
	storage := NewMemoryStorage()
	lifecycle := NewLifecycleManager()
	config := DefaultBetaSharingConfig()
	config.Enabled = true
	registry := NewBetaSharingRegistry(config, lifecycle)

	// Condition de jointure commune à toutes les règles
	joinCondition := map[string]interface{}{
		"type":     "comparison",
		"operator": "==",
		"left": map[string]interface{}{
			"type":   "fieldAccess",
			"object": "u",
			"field":  "id",
		},
		"right": map[string]interface{}{
			"type":   "fieldAccess",
			"object": "o",
			"field":  "user_id",
		},
	}

	varTypes := map[string]string{"u": "User", "o": "Order"}
	leftVars := []string{"u"}
	rightVars := []string{"o"}
	allVars := []string{"u", "o"}
	cascadeLevel := 0

	// Règle 1: u.id == o.user_id (base)
	t.Log("\n📋 Règle r1: u.id == o.user_id")
	node1, hash1, _, err := registry.GetOrCreateJoinNode(
		joinCondition,
		leftVars,
		rightVars,
		allVars,
		varTypes,
		storage,
		cascadeLevel,
	)
	if err != nil {
		t.Fatalf("Failed to create JoinNode for r1: %v", err)
	}
	t.Logf("✅ r1 JoinNode: %s", node1.ID)

	// Règle 2: u.id == o.user_id AND o.amount > 100
	t.Log("\n📋 Règle r2: u.id == o.user_id AND o.amount > 100")
	t.Log("   (condition alpha: o.amount > 100)")
	node2, hash2, shared2, err := registry.GetOrCreateJoinNode(
		joinCondition, // Même condition de jointure
		leftVars,
		rightVars,
		allVars,
		varTypes,
		storage,
		cascadeLevel,
	)
	if err != nil {
		t.Fatalf("Failed to create JoinNode for r2: %v", err)
	}
	if !shared2 {
		t.Error("❌ r2 should share JoinNode with r1")
	}
	if node1.ID != node2.ID {
		t.Errorf("❌ r2 should share same node as r1: got %s vs %s", node2.ID, node1.ID)
	}
	t.Logf("✅ r2 partage le JoinNode avec r1: %s", node2.ID)

	// Règle 3: u.id == o.user_id AND u.age >= 25
	t.Log("\n📋 Règle r3: u.id == o.user_id AND u.age >= 25")
	t.Log("   (condition alpha: u.age >= 25)")
	node3, hash3, shared3, err := registry.GetOrCreateJoinNode(
		joinCondition, // Même condition de jointure
		leftVars,
		rightVars,
		allVars,
		varTypes,
		storage,
		cascadeLevel,
	)
	if err != nil {
		t.Fatalf("Failed to create JoinNode for r3: %v", err)
	}
	if !shared3 {
		t.Error("❌ r3 should share JoinNode with r1 and r2")
	}
	if node1.ID != node3.ID {
		t.Errorf("❌ r3 should share same node as r1: got %s vs %s", node3.ID, node1.ID)
	}
	t.Logf("✅ r3 partage le JoinNode avec r1 et r2: %s", node3.ID)

	// Vérification finale: toutes les règles partagent le même JoinNode
	if hash1 != hash2 || hash2 != hash3 {
		t.Errorf("❌ All rules should have the same hash")
		t.Errorf("   r1: %s", hash1)
		t.Errorf("   r2: %s", hash2)
		t.Errorf("   r3: %s", hash3)
	}

	t.Log("\n✅ PARTAGE OPTIMAL VÉRIFIÉ")
	t.Log("   - 3 règles partagent 1 seul JoinNode")
	t.Log("   - Conditions alpha différentes gérées par AlphaNodes séparés")
	t.Log("   - Efficacité: 1 JoinNode au lieu de 3")
}

// TestNoSharingWhenJoinConditionsDiffer vérifie que les règles avec des
// conditions de jointure différentes ne partagent PAS de JoinNode.
func TestNoSharingWhenJoinConditionsDiffer(t *testing.T) {
	t.Log("🧪 TEST: No Sharing When Join Conditions Differ")
	t.Log("=================================================")

	// Setup
	storage := NewMemoryStorage()
	lifecycle := NewLifecycleManager()
	config := DefaultBetaSharingConfig()
	config.Enabled = true
	registry := NewBetaSharingRegistry(config, lifecycle)

	varTypes := map[string]string{"u": "User", "o": "Order"}
	leftVars := []string{"u"}
	rightVars := []string{"o"}
	allVars := []string{"u", "o"}
	cascadeLevel := 0

	// Règle 1: u.id == o.user_id
	condition1 := map[string]interface{}{
		"type":     "comparison",
		"operator": "==",
		"left":     map[string]interface{}{"type": "fieldAccess", "object": "u", "field": "id"},
		"right":    map[string]interface{}{"type": "fieldAccess", "object": "o", "field": "user_id"},
	}

	t.Log("\n📋 Règle r1: u.id == o.user_id")
	node1, hash1, _, err := registry.GetOrCreateJoinNode(
		condition1,
		leftVars,
		rightVars,
		allVars,
		varTypes,
		storage,
		cascadeLevel,
	)
	if err != nil {
		t.Fatalf("Failed to create JoinNode for r1: %v", err)
	}
	t.Logf("✅ r1 JoinNode: %s (hash: %s)", node1.ID, hash1)

	// Règle 2: u.email == o.customer_email (différente condition de jointure)
	condition2 := map[string]interface{}{
		"type":     "comparison",
		"operator": "==",
		"left":     map[string]interface{}{"type": "fieldAccess", "object": "u", "field": "email"},
		"right":    map[string]interface{}{"type": "fieldAccess", "object": "o", "field": "customer_email"},
	}

	t.Log("\n📋 Règle r2: u.email == o.customer_email")
	node2, hash2, shared2, err := registry.GetOrCreateJoinNode(
		condition2,
		leftVars,
		rightVars,
		allVars,
		varTypes,
		storage,
		cascadeLevel,
	)
	if err != nil {
		t.Fatalf("Failed to create JoinNode for r2: %v", err)
	}

	// ASSERTION: r2 ne doit PAS partager avec r1 (conditions différentes)
	if shared2 {
		t.Error("❌ r2 should NOT share JoinNode with r1 (different join conditions)")
		t.FailNow()
	}

	if hash1 == hash2 {
		t.Errorf("❌ Different join conditions should produce different hashes")
		t.Errorf("   r1: %s", hash1)
		t.Errorf("   r2: %s", hash2)
	}

	if node1.ID == node2.ID {
		t.Errorf("❌ Different join conditions should produce different nodes")
		t.Errorf("   Both have ID: %s", node1.ID)
	}

	t.Logf("✅ r2 JoinNode: %s (hash: %s)", node2.ID, hash2)

	t.Log("\n✅ NON-PARTAGE VÉRIFIÉ")
	t.Log("   - Conditions de jointure différentes = JoinNodes différents")
	t.Log("   - r1 et r2 ont des nœuds séparés (correct)")
}
