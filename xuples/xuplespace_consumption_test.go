// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package xuples

import (
	"fmt"
	"testing"
	"time"

	"github.com/treivax/tsd/rete"
)

// TestRetrieveAutomaticallyMarksConsumed vérifie que Retrieve() marque
// automatiquement le xuple comme consommé (fix du bug critique).
func TestRetrieveAutomaticallyMarksConsumed(t *testing.T) {
	t.Log("═══════════════════════════════════════════════════════════════")
	t.Log("🐛 TEST FIX BUG: Retrieve() doit marquer automatiquement consommé")
	t.Log("═══════════════════════════════════════════════════════════════")

	config := XupleSpaceConfig{
		Name:              "test_once",
		SelectionPolicy:   NewFIFOSelectionPolicy(),
		ConsumptionPolicy: NewOnceConsumptionPolicy(),
		RetentionPolicy:   NewDurationRetentionPolicy(time.Hour),
		MaxSize:           0,
	}

	space := NewXupleSpace(config)

	// Créer un fait de test
	fact := &rete.Fact{
		Type: "TestFact",
		Fields: map[string]interface{}{
			"id":      "fact1",
			"message": "Test message",
		},
	}

	// Créer et insérer un xuple
	xuple := &Xuple{
		ID:              "xuple-test-001",
		Fact:            fact,
		TriggeringFacts: nil,
		CreatedAt:       time.Now(),
		Metadata: XupleMetadata{
			State:      XupleStateAvailable,
			ConsumedBy: make(map[string]time.Time),
		},
	}

	err := space.Insert(xuple)
	if err != nil {
		t.Fatalf("❌ Erreur Insert: %v", err)
	}

	t.Log("✅ Xuple inséré avec succès")
	t.Logf("   ID: %s", xuple.ID)
	t.Logf("   État initial: %s", xuple.Metadata.State)

	// Vérifier le compte avant récupération
	countBefore := space.Count()
	if countBefore != 1 {
		t.Errorf("❌ Count avant Retrieve devrait être 1, obtenu: %d", countBefore)
	} else {
		t.Logf("✅ Count avant Retrieve: %d", countBefore)
	}

	// PREMIÈRE RÉCUPÉRATION par agent1
	t.Log("")
	t.Log("───────────────────────────────────────────────────────────────")
	t.Log("TEST 1: Première récupération par agent1")
	t.Log("───────────────────────────────────────────────────────────────")

	retrieved1, err := space.Retrieve("agent1")
	if err != nil {
		t.Fatalf("❌ Erreur Retrieve agent1: %v", err)
	}

	if retrieved1.ID != xuple.ID {
		t.Errorf("❌ Mauvais xuple récupéré, attendu: %s, obtenu: %s", xuple.ID, retrieved1.ID)
	} else {
		t.Logf("✅ Xuple récupéré: %s", retrieved1.ID)
	}

	// VÉRIFICATION CRITIQUE: Le xuple DOIT être marqué comme consommé
	if retrieved1.Metadata.ConsumptionCount != 1 {
		t.Errorf("❌ BUG! ConsumptionCount devrait être 1, obtenu: %d", retrieved1.Metadata.ConsumptionCount)
	} else {
		t.Logf("✅ ConsumptionCount correctement incrémenté: %d", retrieved1.Metadata.ConsumptionCount)
	}

	if _, consumed := retrieved1.Metadata.ConsumedBy["agent1"]; !consumed {
		t.Errorf("❌ BUG! agent1 devrait être dans ConsumedBy")
	} else {
		t.Logf("✅ agent1 correctement enregistré dans ConsumedBy")
	}

	// VÉRIFICATION: Avec politique 'once', l'état devrait être XupleStateConsumed
	if retrieved1.Metadata.State != XupleStateConsumed {
		t.Errorf("❌ BUG! État devrait être XupleStateConsumed avec policy 'once', obtenu: %s", retrieved1.Metadata.State)
	} else {
		t.Logf("✅ État correctement changé à: %s", retrieved1.Metadata.State)
	}

	// Le count devrait maintenant être 0 car le xuple est consommé
	countAfter := space.Count()
	if countAfter != 0 {
		t.Errorf("❌ BUG! Count après Retrieve devrait être 0 (xuple consommé), obtenu: %d", countAfter)
	} else {
		t.Logf("✅ Count après Retrieve: %d (xuple consommé)", countAfter)
	}

	// DEUXIÈME TENTATIVE par le même agent
	t.Log("")
	t.Log("───────────────────────────────────────────────────────────────")
	t.Log("TEST 2: Deuxième tentative par le même agent (devrait échouer)")
	t.Log("───────────────────────────────────────────────────────────────")

	retrieved2, err := space.Retrieve("agent1")
	if err == nil {
		t.Errorf("❌ BUG! Deuxième Retrieve devrait échouer (policy 'once'), mais a retourné: %s", retrieved2.ID)
	} else if err != ErrNoAvailableXuple {
		t.Errorf("❌ Erreur inattendue: %v (attendu: ErrNoAvailableXuple)", err)
	} else {
		t.Logf("✅ Deuxième Retrieve a correctement échoué: %v", err)
	}

	// TENTATIVE par un autre agent
	t.Log("")
	t.Log("───────────────────────────────────────────────────────────────")
	t.Log("TEST 3: Tentative par agent2 (devrait aussi échouer)")
	t.Log("───────────────────────────────────────────────────────────────")

	retrieved3, err := space.Retrieve("agent2")
	if err == nil {
		t.Errorf("❌ BUG! Retrieve par agent2 devrait échouer (policy 'once' = consommé globalement), mais a retourné: %s", retrieved3.ID)
	} else if err != ErrNoAvailableXuple {
		t.Errorf("❌ Erreur inattendue: %v (attendu: ErrNoAvailableXuple)", err)
	} else {
		t.Logf("✅ Retrieve par agent2 a correctement échoué: %v", err)
	}

	t.Log("")
	t.Log("═══════════════════════════════════════════════════════════════")
	t.Log("✅ TOUS LES TESTS RÉUSSIS - BUG 'once' CORRIGÉ!")
	t.Log("═══════════════════════════════════════════════════════════════")
}

// TestRetrievePerAgentPolicy vérifie que Retrieve() fonctionne correctement
// avec la politique per-agent (plusieurs agents peuvent consommer le même xuple).
func TestRetrievePerAgentPolicy(t *testing.T) {
	t.Log("═══════════════════════════════════════════════════════════════")
	t.Log("🧪 TEST: Politique per-agent")
	t.Log("═══════════════════════════════════════════════════════════════")

	config := XupleSpaceConfig{
		Name:              "test_per_agent",
		SelectionPolicy:   NewFIFOSelectionPolicy(),
		ConsumptionPolicy: NewPerAgentConsumptionPolicy(),
		RetentionPolicy:   NewDurationRetentionPolicy(time.Hour),
		MaxSize:           0,
	}

	space := NewXupleSpace(config)

	// Créer et insérer un xuple
	fact := &rete.Fact{
		Type: "TestFact",
		Fields: map[string]interface{}{
			"id":      "fact2",
			"message": "Per-agent test",
		},
	}

	xuple := &Xuple{
		ID:              "xuple-per-agent-001",
		Fact:            fact,
		TriggeringFacts: nil,
		CreatedAt:       time.Now(),
		Metadata: XupleMetadata{
			State:      XupleStateAvailable,
			ConsumedBy: make(map[string]time.Time),
		},
	}

	err := space.Insert(xuple)
	if err != nil {
		t.Fatalf("❌ Erreur Insert: %v", err)
	}

	t.Logf("✅ Xuple inséré: %s", xuple.ID)

	// Agent1 récupère le xuple
	t.Log("")
	t.Log("───────────────────────────────────────────────────────────────")
	t.Log("TEST: Agent1 récupère le xuple")
	t.Log("───────────────────────────────────────────────────────────────")

	retrieved1, err := space.Retrieve("agent1")
	if err != nil {
		t.Fatalf("❌ Erreur Retrieve agent1: %v", err)
	}

	t.Logf("✅ Agent1 a récupéré: %s", retrieved1.ID)
	t.Logf("   ConsumptionCount: %d", retrieved1.Metadata.ConsumptionCount)
	t.Logf("   État: %s", retrieved1.Metadata.State)

	// Avec per-agent, l'état devrait rester Available (pas encore consommé globalement)
	if retrieved1.Metadata.State != XupleStateAvailable {
		t.Errorf("❌ État devrait être Available avec per-agent, obtenu: %s", retrieved1.Metadata.State)
	} else {
		t.Logf("✅ État reste Available (per-agent policy)")
	}

	// Agent2 récupère le même xuple (doit fonctionner avec per-agent)
	t.Log("")
	t.Log("───────────────────────────────────────────────────────────────")
	t.Log("TEST: Agent2 récupère le même xuple")
	t.Log("───────────────────────────────────────────────────────────────")

	retrieved2, err := space.Retrieve("agent2")
	if err != nil {
		t.Errorf("❌ Agent2 devrait pouvoir récupérer le xuple (per-agent), erreur: %v", err)
	} else {
		t.Logf("✅ Agent2 a récupéré: %s", retrieved2.ID)
		t.Logf("   ConsumptionCount: %d", retrieved2.Metadata.ConsumptionCount)
	}

	if retrieved2.ID != xuple.ID {
		t.Errorf("❌ Agent2 devrait obtenir le même xuple")
	}

	if retrieved2.Metadata.ConsumptionCount != 2 {
		t.Errorf("❌ ConsumptionCount devrait être 2, obtenu: %d", retrieved2.Metadata.ConsumptionCount)
	} else {
		t.Logf("✅ ConsumptionCount correctement incrémenté à: %d", retrieved2.Metadata.ConsumptionCount)
	}

	// Agent1 ne devrait pas pouvoir le récupérer à nouveau
	t.Log("")
	t.Log("───────────────────────────────────────────────────────────────")
	t.Log("TEST: Agent1 tente de récupérer à nouveau (devrait échouer)")
	t.Log("───────────────────────────────────────────────────────────────")

	retrieved3, err := space.Retrieve("agent1")
	if err == nil {
		t.Errorf("❌ Agent1 ne devrait pas pouvoir récupérer à nouveau, mais a obtenu: %s", retrieved3.ID)
	} else {
		t.Logf("✅ Agent1 ne peut pas récupérer à nouveau: %v", err)
	}

	t.Log("")
	t.Log("═══════════════════════════════════════════════════════════════")
	t.Log("✅ TEST per-agent RÉUSSI!")
	t.Log("═══════════════════════════════════════════════════════════════")
}

// TestRetrieveLimitedPolicy vérifie que Retrieve() fonctionne correctement
// avec la politique limited(n).
func TestRetrieveLimitedPolicy(t *testing.T) {
	t.Log("═══════════════════════════════════════════════════════════════")
	t.Log("🧪 TEST: Politique limited(3)")
	t.Log("═══════════════════════════════════════════════════════════════")

	config := XupleSpaceConfig{
		Name:              "test_limited",
		SelectionPolicy:   NewFIFOSelectionPolicy(),
		ConsumptionPolicy: NewLimitedConsumptionPolicy(3), // Limite à 3 consommations
		RetentionPolicy:   NewDurationRetentionPolicy(time.Hour),
		MaxSize:           0,
	}

	space := NewXupleSpace(config)

	// Créer et insérer un xuple
	fact := &rete.Fact{
		Type: "TestFact",
		Fields: map[string]interface{}{
			"id":      "fact3",
			"message": "Limited test",
		},
	}

	xuple := &Xuple{
		ID:              "xuple-limited-001",
		Fact:            fact,
		TriggeringFacts: nil,
		CreatedAt:       time.Now(),
		Metadata: XupleMetadata{
			State:      XupleStateAvailable,
			ConsumedBy: make(map[string]time.Time),
		},
	}

	err := space.Insert(xuple)
	if err != nil {
		t.Fatalf("❌ Erreur Insert: %v", err)
	}

	t.Logf("✅ Xuple inséré: %s", xuple.ID)

	// Consommer 3 fois (la limite)
	agents := []string{"agent1", "agent2", "agent3"}
	for i, agentID := range agents {
		t.Logf("")
		t.Logf("───────────────────────────────────────────────────────────────")
		t.Logf("TEST: Consommation %d/%d par %s", i+1, 3, agentID)
		t.Logf("───────────────────────────────────────────────────────────────")

		retrieved, err := space.Retrieve(agentID)
		if err != nil {
			t.Fatalf("❌ Erreur Retrieve %s: %v", agentID, err)
		}

		t.Logf("✅ %s a récupéré: %s", agentID, retrieved.ID)
		t.Logf("   ConsumptionCount: %d", retrieved.Metadata.ConsumptionCount)
		t.Logf("   État: %s", retrieved.Metadata.State)

		expectedCount := i + 1
		if retrieved.Metadata.ConsumptionCount != expectedCount {
			t.Errorf("❌ ConsumptionCount devrait être %d, obtenu: %d", expectedCount, retrieved.Metadata.ConsumptionCount)
		}

		// Après la 3ème consommation, l'état devrait être Consumed
		if i == 2 { // 3ème itération (index 2)
			if retrieved.Metadata.State != XupleStateConsumed {
				t.Errorf("❌ Après 3 consommations, l'état devrait être Consumed, obtenu: %s", retrieved.Metadata.State)
			} else {
				t.Logf("✅ État correctement changé à Consumed après limite atteinte")
			}
		}
	}

	// Tentative de 4ème consommation (devrait échouer)
	t.Log("")
	t.Log("───────────────────────────────────────────────────────────────")
	t.Log("TEST: Tentative de 4ème consommation (devrait échouer)")
	t.Log("───────────────────────────────────────────────────────────────")

	retrieved4, err := space.Retrieve("agent4")
	if err == nil {
		t.Errorf("❌ 4ème consommation devrait échouer (limite=3), mais a obtenu: %s", retrieved4.ID)
	} else if err != ErrNoAvailableXuple {
		t.Errorf("❌ Erreur inattendue: %v (attendu: ErrNoAvailableXuple)", err)
	} else {
		t.Logf("✅ 4ème consommation a correctement échoué: %v", err)
	}

	t.Log("")
	t.Log("═══════════════════════════════════════════════════════════════")
	t.Log("✅ TEST limited(3) RÉUSSI!")
	t.Log("═══════════════════════════════════════════════════════════════")
}

// TestMultipleXuplesWithOncePolicy vérifie le comportement avec plusieurs xuples
// et la politique 'once'.
func TestMultipleXuplesWithOncePolicy(t *testing.T) {
	t.Log("═══════════════════════════════════════════════════════════════")
	t.Log("🧪 TEST: Plusieurs xuples avec politique 'once'")
	t.Log("═══════════════════════════════════════════════════════════════")

	config := XupleSpaceConfig{
		Name:              "test_multiple_once",
		SelectionPolicy:   NewFIFOSelectionPolicy(),
		ConsumptionPolicy: NewOnceConsumptionPolicy(),
		RetentionPolicy:   NewDurationRetentionPolicy(time.Hour),
		MaxSize:           0,
	}

	space := NewXupleSpace(config)

	// Insérer 5 xuples
	numXuples := 5
	for i := 1; i <= numXuples; i++ {
		fact := &rete.Fact{
			Type: "TestFact",
			Fields: map[string]interface{}{
				"id":      fmt.Sprintf("fact-%d", i),
				"message": fmt.Sprintf("Message %d", i),
			},
		}

		xuple := &Xuple{
			ID:              fmt.Sprintf("xuple-%03d", i),
			Fact:            fact,
			TriggeringFacts: nil,
			CreatedAt:       time.Now(),
			Metadata: XupleMetadata{
				State:      XupleStateAvailable,
				ConsumedBy: make(map[string]time.Time),
			},
		}

		err := space.Insert(xuple)
		if err != nil {
			t.Fatalf("❌ Erreur Insert xuple %d: %v", i, err)
		}
	}

	t.Logf("✅ %d xuples insérés", numXuples)

	initialCount := space.Count()
	if initialCount != numXuples {
		t.Errorf("❌ Count initial devrait être %d, obtenu: %d", numXuples, initialCount)
	} else {
		t.Logf("✅ Count initial: %d", initialCount)
	}

	// Un agent consomme tous les xuples un par un
	t.Log("")
	t.Log("───────────────────────────────────────────────────────────────")
	t.Log("TEST: Agent1 consomme tous les xuples")
	t.Log("───────────────────────────────────────────────────────────────")

	retrievedIDs := make(map[string]bool)
	for i := 1; i <= numXuples; i++ {
		retrieved, err := space.Retrieve("agent1")
		if err != nil {
			t.Fatalf("❌ Erreur Retrieve %d: %v", i, err)
		}

		t.Logf("   Récupération %d: %s", i, retrieved.ID)

		// Vérifier qu'on ne récupère pas le même xuple deux fois
		if retrievedIDs[retrieved.ID] {
			t.Errorf("❌ BUG! Xuple %s récupéré plusieurs fois!", retrieved.ID)
		}
		retrievedIDs[retrieved.ID] = true

		// Vérifier que le xuple est marqué comme consommé
		if retrieved.Metadata.State != XupleStateConsumed {
			t.Errorf("❌ Xuple %s devrait être Consumed, obtenu: %s", retrieved.ID, retrieved.Metadata.State)
		}

		// Vérifier le count décrémente
		expectedCount := numXuples - i
		actualCount := space.Count()
		if actualCount != expectedCount {
			t.Errorf("❌ Après %d consommations, Count devrait être %d, obtenu: %d", i, expectedCount, actualCount)
		}
	}

	t.Logf("✅ Tous les xuples consommés (IDs uniques: %d)", len(retrievedIDs))

	// Le count final devrait être 0
	finalCount := space.Count()
	if finalCount != 0 {
		t.Errorf("❌ Count final devrait être 0, obtenu: %d", finalCount)
	} else {
		t.Logf("✅ Count final: %d", finalCount)
	}

	// Tentative de récupération sur espace vide
	t.Log("")
	t.Log("───────────────────────────────────────────────────────────────")
	t.Log("TEST: Tentative sur xuple-space vide (devrait échouer)")
	t.Log("───────────────────────────────────────────────────────────────")

	retrieved, err := space.Retrieve("agent1")
	if err == nil {
		t.Errorf("❌ Retrieve devrait échouer sur espace vide, mais a obtenu: %s", retrieved.ID)
	} else if err != ErrNoAvailableXuple {
		t.Errorf("❌ Erreur inattendue: %v (attendu: ErrNoAvailableXuple)", err)
	} else {
		t.Logf("✅ Retrieve a correctement échoué: %v", err)
	}

	t.Log("")
	t.Log("═══════════════════════════════════════════════════════════════")
	t.Log("✅ TEST PLUSIEURS XUPLES RÉUSSI!")
	t.Log("═══════════════════════════════════════════════════════════════")
}
