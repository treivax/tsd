// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package xuples

import (
	"testing"
	"time"

	"github.com/treivax/tsd/rete"
)

// Test constants
const (
	// TestNumGoroutines nombre de goroutines pour tests concurrents
	TestNumGoroutines = 10

	// TestItemsPerGoroutine nombre d'items par goroutine
	TestItemsPerGoroutine = 10

	// TestNumAgents nombre d'agents pour tests concurrents
	TestNumAgents = 10

	// TestRetrievalsPerAgent nombre de récupérations par agent
	TestRetrievalsPerAgent = 5

	// TestNumXuples nombre de xuples à créer dans les tests
	TestNumXuples = 50

	// TestCleanupWaitDuration temps d'attente pour tests de cleanup
	TestCleanupWaitDuration = 100 * time.Millisecond

	// TestRetentionDuration durée de rétention pour tests
	TestRetentionDuration = 50 * time.Millisecond
)

// Test helpers

func createTestFact(id string) *rete.Fact {
	return &rete.Fact{
		ID:   id,
		Type: "TestFact",
	}
}

func createTestXuple(id string) *Xuple {
	return &Xuple{
		ID:              id,
		Fact:            createTestFact("f_" + id),
		TriggeringFacts: []*rete.Fact{createTestFact("t1_" + id), createTestFact("t2_" + id)},
		CreatedAt:       time.Now(),
		Metadata: XupleMetadata{
			State:      XupleStateAvailable,
			ConsumedBy: make(map[string]time.Time),
		},
	}
}

// Tests for Xuple

func TestXupleIsAvailable(t *testing.T) {
	t.Log("🧪 TEST XUPLE IS AVAILABLE")

	xuple := createTestXuple("x1")

	if !xuple.IsAvailable() {
		t.Error("❌ Xuple devrait être disponible")
	}

	xuple.Metadata.State = XupleStateConsumed
	if xuple.IsAvailable() {
		t.Error("❌ Xuple consumed ne devrait pas être disponible")
	}

	t.Log("✅ IsAvailable fonctionne correctement")
}

func TestXupleIsExpired(t *testing.T) {
	t.Log("🧪 TEST XUPLE IS EXPIRED")

	xuple := createTestXuple("x1")

	// Pas d'expiration définie
	if xuple.IsExpired() {
		t.Error("❌ Xuple sans expiration ne devrait pas être expiré")
	}

	// Expiration dans le futur
	futureTime := time.Now().Add(1 * time.Hour)
	xuple.Metadata.ExpiresAt = futureTime
	if xuple.IsExpired() {
		t.Error("❌ Xuple avec expiration future ne devrait pas être expiré")
	}

	// Expiration dans le passé
	pastTime := time.Now().Add(-1 * time.Hour)
	xuple.Metadata.ExpiresAt = pastTime
	if !xuple.IsExpired() {
		t.Error("❌ Xuple avec expiration passée devrait être expiré")
	}

	// Vérifier que IsExpired est read-only (ne modifie pas l'état)
	// L'état initial devrait toujours être Available
	if xuple.Metadata.State != XupleStateAvailable {
		t.Error("❌ IsExpired ne devrait pas modifier l'état (read-only)")
	}

	// Tester avec état déjà marqué comme expiré
	xuple.Metadata.State = XupleStateExpired
	if !xuple.IsExpired() {
		t.Error("❌ Xuple avec State=Expired devrait être considéré expiré")
	}

	t.Log("✅ IsExpired fonctionne correctement")
}

func TestXupleCanBeConsumedBy(t *testing.T) {
	t.Log("🧪 TEST XUPLE CAN BE CONSUMED BY")

	xuple := createTestXuple("x1")
	policy := NewOnceConsumptionPolicy()

	// Disponible et non expiré
	if !xuple.CanBeConsumedBy("agent1", policy) {
		t.Error("❌ Xuple disponible devrait pouvoir être consommé")
	}

	// Marquer comme consommé
	xuple.Metadata.State = XupleStateConsumed
	if xuple.CanBeConsumedBy("agent2", policy) {
		t.Error("❌ Xuple consumed ne devrait pas pouvoir être consommé")
	}

	t.Log("✅ CanBeConsumedBy fonctionne correctement")
}

func TestXupleMarkConsumedByViaSpace(t *testing.T) {
	t.Log("🧪 TEST XUPLE MARK CONSUMED BY VIA SPACE")

	// Créer un xuple-space pour tester markConsumedBy de manière thread-safe
	manager := NewXupleManager()
	config := XupleSpaceConfig{
		Name:              "test",
		SelectionPolicy:   NewFIFOSelectionPolicy(),
		ConsumptionPolicy: NewPerAgentConsumptionPolicy(),
		RetentionPolicy:   NewUnlimitedRetentionPolicy(),
	}

	err := manager.CreateXupleSpace("test", config)
	if err != nil {
		t.Fatalf("❌ Erreur création space: %v", err)
	}

	space, err := manager.GetXupleSpace("test")
	if err != nil {
		t.Fatalf("❌ Erreur récupération space: %v", err)
	}

	// Créer et insérer un xuple
	err = manager.CreateXuple("test", createTestFact("f1"), []*rete.Fact{createTestFact("t1")})
	if err != nil {
		t.Fatalf("❌ Erreur création xuple: %v", err)
	}

	// Récupérer le xuple - cela le marque AUTOMATIQUEMENT comme consommé par agent1
	xuple, err := space.Retrieve("agent1")
	if err != nil {
		t.Fatalf("❌ Erreur récupération: %v", err)
	}

	// Vérifier que Retrieve() a automatiquement marqué comme consommé
	if xuple.Metadata.ConsumptionCount != 1 {
		t.Errorf("❌ ConsumptionCount devrait être 1 (marqué par Retrieve), reçu %d", xuple.Metadata.ConsumptionCount)
	}

	if _, consumed := xuple.Metadata.ConsumedBy["agent1"]; !consumed {
		t.Error("❌ agent1 devrait être dans ConsumedBy (marqué par Retrieve)")
	}

	// Tester MarkConsumed() avec un agent différent (per-agent policy)
	err = space.MarkConsumed(xuple.ID, "agent2")
	if err != nil {
		t.Fatalf("❌ Erreur MarkConsumed pour agent2: %v", err)
	}

	// Vérifier que agent2 est maintenant aussi enregistré
	if xuple.Metadata.ConsumptionCount != 2 {
		t.Errorf("❌ ConsumptionCount devrait être 2, reçu %d", xuple.Metadata.ConsumptionCount)
	}

	if _, consumed := xuple.Metadata.ConsumedBy["agent2"]; !consumed {
		t.Error("❌ agent2 devrait être dans ConsumedBy")
	}

	t.Log("✅ MarkConsumed via XupleSpace fonctionne correctement (thread-safe)")
}

// Tests for XupleManager

func TestNewXupleManager(t *testing.T) {
	t.Log("🧪 TEST NEW XUPLE MANAGER")

	manager := NewXupleManager()

	if manager == nil {
		t.Fatal("❌ Manager ne devrait pas être nil")
	}

	spaces := manager.ListXupleSpaces()
	if len(spaces) != 0 {
		t.Errorf("❌ Manager neuf devrait avoir 0 spaces, reçu %d", len(spaces))
	}

	t.Log("✅ NewXupleManager fonctionne correctement")
}

func TestCreateXupleSpace(t *testing.T) {
	t.Log("🧪 TEST CREATE XUPLE SPACE")

	manager := NewXupleManager()
	config := XupleSpaceConfig{
		Name:              "test",
		SelectionPolicy:   NewFIFOSelectionPolicy(),
		ConsumptionPolicy: NewOnceConsumptionPolicy(),
		RetentionPolicy:   NewUnlimitedRetentionPolicy(),
	}

	err := manager.CreateXupleSpace("test", config)
	if err != nil {
		t.Fatalf("❌ Erreur création: %v", err)
	}

	// Vérifier que le space existe
	space, err := manager.GetXupleSpace("test")
	if err != nil {
		t.Fatalf("❌ Erreur récupération: %v", err)
	}

	if space.Name() != "test" {
		t.Errorf("❌ Nom devrait être 'test', reçu '%s'", space.Name())
	}

	t.Log("✅ CreateXupleSpace fonctionne correctement")
}

func TestCreateXupleSpaceDuplicate(t *testing.T) {
	t.Log("🧪 TEST CREATE XUPLE SPACE DUPLICATE")

	manager := NewXupleManager()
	config := XupleSpaceConfig{
		Name:              "test",
		SelectionPolicy:   NewFIFOSelectionPolicy(),
		ConsumptionPolicy: NewOnceConsumptionPolicy(),
		RetentionPolicy:   NewUnlimitedRetentionPolicy(),
	}

	err := manager.CreateXupleSpace("test", config)
	if err != nil {
		t.Fatalf("❌ Erreur première création: %v", err)
	}

	// Tenter de recréer
	err = manager.CreateXupleSpace("test", config)
	if err != ErrXupleSpaceExists {
		t.Errorf("❌ Devrait retourner ErrXupleSpaceExists, reçu %v", err)
	}

	t.Log("✅ Détection duplicate fonctionne correctement")
}

func TestGetXupleSpaceNotFound(t *testing.T) {
	t.Log("🧪 TEST GET XUPLE SPACE NOT FOUND")

	manager := NewXupleManager()

	_, err := manager.GetXupleSpace("inexistant")
	if err != ErrXupleSpaceNotFound {
		t.Errorf("❌ Devrait retourner ErrXupleSpaceNotFound, reçu %v", err)
	}

	t.Log("✅ Erreur space not found fonctionne correctement")
}

func TestCreateXuple(t *testing.T) {
	t.Log("🧪 TEST CREATE XUPLE")

	manager := NewXupleManager()
	config := XupleSpaceConfig{
		Name:              "test",
		SelectionPolicy:   NewFIFOSelectionPolicy(),
		ConsumptionPolicy: NewOnceConsumptionPolicy(),
		RetentionPolicy:   NewUnlimitedRetentionPolicy(),
	}

	err := manager.CreateXupleSpace("test", config)
	if err != nil {
		t.Fatalf("❌ Erreur création space: %v", err)
	}

	fact := createTestFact("f1")
	triggering := []*rete.Fact{createTestFact("t1")}

	err = manager.CreateXuple("test", fact, triggering)
	if err != nil {
		t.Fatalf("❌ Erreur création xuple: %v", err)
	}

	// Vérifier que le xuple a été ajouté
	space, _ := manager.GetXupleSpace("test")
	if space.Count() != 1 {
		t.Errorf("❌ Space devrait avoir 1 xuple, reçu %d", space.Count())
	}

	t.Log("✅ CreateXuple fonctionne correctement")
}

func TestListXupleSpaces(t *testing.T) {
	t.Log("🧪 TEST LIST XUPLE SPACES")

	manager := NewXupleManager()
	config := XupleSpaceConfig{
		SelectionPolicy:   NewFIFOSelectionPolicy(),
		ConsumptionPolicy: NewOnceConsumptionPolicy(),
		RetentionPolicy:   NewUnlimitedRetentionPolicy(),
	}

	if err := manager.CreateXupleSpace("space1", config); err != nil {
		t.Fatalf("❌ Erreur création space1: %v", err)
	}
	if err := manager.CreateXupleSpace("space2", config); err != nil {
		t.Fatalf("❌ Erreur création space2: %v", err)
	}
	if err := manager.CreateXupleSpace("space3", config); err != nil {
		t.Fatalf("❌ Erreur création space3: %v", err)
	}

	spaces := manager.ListXupleSpaces()
	if len(spaces) != 3 {
		t.Errorf("❌ Devrait avoir 3 spaces, reçu %d", len(spaces))
	}

	t.Log("✅ ListXupleSpaces fonctionne correctement")
}

func TestCloseManager(t *testing.T) {
	t.Log("🧪 TEST CLOSE MANAGER")

	manager := NewXupleManager()
	config := XupleSpaceConfig{
		SelectionPolicy:   NewFIFOSelectionPolicy(),
		ConsumptionPolicy: NewOnceConsumptionPolicy(),
		RetentionPolicy:   NewUnlimitedRetentionPolicy(),
	}

	if err := manager.CreateXupleSpace("test", config); err != nil {
		t.Fatalf("❌ Erreur création test space: %v", err)
	}

	err := manager.Close()
	if err != nil {
		t.Errorf("❌ Erreur close: %v", err)
	}

	// Vérifier que les spaces ont été nettoyés
	spaces := manager.ListXupleSpaces()
	if len(spaces) != 0 {
		t.Errorf("❌ Devrait avoir 0 spaces après close, reçu %d", len(spaces))
	}

	t.Log("✅ Close fonctionne correctement")
}

// Tests for XupleSpace

func TestInsertXuple(t *testing.T) {
	t.Log("🧪 TEST INSERT XUPLE")

	space := NewXupleSpace(XupleSpaceConfig{
		Name:              "test",
		SelectionPolicy:   NewFIFOSelectionPolicy(),
		ConsumptionPolicy: NewOnceConsumptionPolicy(),
		RetentionPolicy:   NewUnlimitedRetentionPolicy(),
	})

	xuple := createTestXuple("x1")
	err := space.Insert(xuple)
	if err != nil {
		t.Fatalf("❌ Erreur insert: %v", err)
	}

	if space.Count() != 1 {
		t.Errorf("❌ Count devrait être 1, reçu %d", space.Count())
	}

	t.Log("✅ Insert fonctionne correctement")
}

func TestInsertNilXuple(t *testing.T) {
	t.Log("🧪 TEST INSERT NIL XUPLE")

	space := NewXupleSpace(XupleSpaceConfig{
		Name:              "test",
		SelectionPolicy:   NewFIFOSelectionPolicy(),
		ConsumptionPolicy: NewOnceConsumptionPolicy(),
		RetentionPolicy:   NewUnlimitedRetentionPolicy(),
	})

	err := space.Insert(nil)
	if err != ErrNilXuple {
		t.Errorf("❌ Devrait retourner ErrNilXuple, reçu %v", err)
	}

	t.Log("✅ Validation nil xuple fonctionne correctement")
}

func TestRetrieveXuple(t *testing.T) {
	t.Log("🧪 TEST RETRIEVE XUPLE")

	space := NewXupleSpace(XupleSpaceConfig{
		Name:              "test",
		SelectionPolicy:   NewFIFOSelectionPolicy(),
		ConsumptionPolicy: NewOnceConsumptionPolicy(),
		RetentionPolicy:   NewUnlimitedRetentionPolicy(),
	})

	xuple := createTestXuple("x1")
	if err := space.Insert(xuple); err != nil {
		t.Fatalf("❌ Erreur insertion xuple: %v", err)
	}

	retrieved, err := space.Retrieve("agent1")
	if err != nil {
		t.Fatalf("❌ Erreur retrieve: %v", err)
	}

	if retrieved == nil {
		t.Fatal("❌ Xuple retrieved ne devrait pas être nil")
	}

	if retrieved.ID != xuple.ID {
		t.Errorf("❌ ID incorrect, attendu %s, reçu %s", xuple.ID, retrieved.ID)
	}

	t.Log("✅ Retrieve fonctionne correctement")
}

func TestRetrieveNoAvailable(t *testing.T) {
	t.Log("🧪 TEST RETRIEVE NO AVAILABLE")

	space := NewXupleSpace(XupleSpaceConfig{
		Name:              "test",
		SelectionPolicy:   NewFIFOSelectionPolicy(),
		ConsumptionPolicy: NewOnceConsumptionPolicy(),
		RetentionPolicy:   NewUnlimitedRetentionPolicy(),
	})

	_, err := space.Retrieve("agent1")
	if err != ErrNoAvailableXuple {
		t.Errorf("❌ Devrait retourner ErrNoAvailableXuple, reçu %v", err)
	}

	t.Log("✅ No available xuple fonctionne correctement")
}

func TestMarkConsumed(t *testing.T) {
	t.Log("🧪 TEST MARK CONSUMED")

	space := NewXupleSpace(XupleSpaceConfig{
		Name:              "test",
		SelectionPolicy:   NewFIFOSelectionPolicy(),
		ConsumptionPolicy: NewOnceConsumptionPolicy(),
		RetentionPolicy:   NewUnlimitedRetentionPolicy(),
	})

	xuple := createTestXuple("x1")
	if err := space.Insert(xuple); err != nil {
		t.Fatalf("❌ Erreur insertion xuple: %v", err)
	}

	err := space.MarkConsumed(xuple.ID, "agent1")
	if err != nil {
		t.Fatalf("❌ Erreur mark consumed: %v", err)
	}

	// Vérifier que le xuple a été marqué
	if xuple.Metadata.ConsumptionCount != 1 {
		t.Errorf("❌ ConsumptionCount devrait être 1, reçu %d", xuple.Metadata.ConsumptionCount)
	}

	// Avec politique Once, devrait être consumed
	if xuple.Metadata.State != XupleStateConsumed {
		t.Errorf("❌ State devrait être Consumed, reçu %s", xuple.Metadata.State)
	}

	t.Log("✅ MarkConsumed fonctionne correctement")
}

func TestCleanup(t *testing.T) {
	t.Log("🧪 TEST CLEANUP")

	space := NewXupleSpace(XupleSpaceConfig{
		Name:              "test",
		SelectionPolicy:   NewFIFOSelectionPolicy(),
		ConsumptionPolicy: NewOnceConsumptionPolicy(),
		RetentionPolicy:   NewDurationRetentionPolicy(TestRetentionDuration),
	})

	xuple := createTestXuple("x1")
	if err := space.Insert(xuple); err != nil {
		t.Fatalf("❌ Erreur insertion xuple: %v", err)
	}

	// Attendre l'expiration
	time.Sleep(TestCleanupWaitDuration)

	cleaned := space.Cleanup()
	if cleaned != 1 {
		t.Errorf("❌ Devrait nettoyer 1 xuple, reçu %d", cleaned)
	}

	t.Log("✅ Cleanup fonctionne correctement")
}

// Tests de concurrence

func TestConcurrentInsert(t *testing.T) {
	t.Log("🧪 TEST CONCURRENT INSERT")

	space := NewXupleSpace(XupleSpaceConfig{
		Name:              "test",
		SelectionPolicy:   NewFIFOSelectionPolicy(),
		ConsumptionPolicy: NewOnceConsumptionPolicy(),
		RetentionPolicy:   NewUnlimitedRetentionPolicy(),
	})

	done := make(chan bool)

	for i := 0; i < TestNumGoroutines; i++ {
		go func(id int) {
			for j := 0; j < TestItemsPerGoroutine; j++ {
				xuple := createTestXuple(string(rune(id*1000 + j)))
				if err := space.Insert(xuple); err != nil {
					t.Errorf("❌ Erreur insertion concurrent: %v", err)
				}
			}
			done <- true
		}(i)
	}

	// Attendre toutes les goroutines
	for i := 0; i < TestNumGoroutines; i++ {
		<-done
	}

	expectedCount := TestNumGoroutines * TestItemsPerGoroutine
	if space.Count() != expectedCount {
		t.Errorf("❌ Count devrait être %d, reçu %d", expectedCount, space.Count())
	}

	t.Log("✅ Insertion concurrente fonctionne correctement")
}

func TestConcurrentRetrieve(t *testing.T) {
	t.Log("🧪 TEST CONCURRENT RETRIEVE")

	space := NewXupleSpace(XupleSpaceConfig{
		Name:              "test",
		SelectionPolicy:   NewFIFOSelectionPolicy(),
		ConsumptionPolicy: NewPerAgentConsumptionPolicy(),
		RetentionPolicy:   NewUnlimitedRetentionPolicy(),
	})

	// Insérer plusieurs xuples
	for i := 0; i < TestNumXuples; i++ {
		if err := space.Insert(createTestXuple(string(rune(i)))); err != nil {
			t.Fatalf("❌ Erreur insertion xuple %d: %v", i, err)
		}
	}

	done := make(chan bool)

	for i := 0; i < TestNumAgents; i++ {
		go func(agentNum int) {
			agentID := string(rune('A' + agentNum))
			for j := 0; j < TestRetrievalsPerAgent; j++ {
				_, err := space.Retrieve(agentID)
				if err != nil && err != ErrNoAvailableXuple {
					t.Errorf("❌ Erreur retrieve agent %s: %v", agentID, err)
				}
			}
			done <- true
		}(i)
	}

	// Attendre tous les agents
	for i := 0; i < TestNumAgents; i++ {
		<-done
	}

	t.Log("✅ Retrieve concurrent fonctionne correctement")
}

// Tests pour fonctions utilitaires

func TestXupleState_String(t *testing.T) {
	t.Log("🧪 TEST XUPLE STATE STRING")

	tests := []struct {
		state    XupleState
		expected string
	}{
		{XupleStateAvailable, "available"},
		{XupleStateConsumed, "consumed"},
		{XupleStateExpired, "expired"},
		{XupleState(999), "unknown"},
	}

	for _, tt := range tests {
		result := tt.state.String()
		if result != tt.expected {
			t.Errorf("❌ State %d: attendu '%s', reçu '%s'", tt.state, tt.expected, result)
		}
	}

	t.Log("✅ XupleState.String() fonctionne correctement")
}

func TestGetConfig(t *testing.T) {
	t.Log("🧪 TEST GET CONFIG")

	config := XupleSpaceConfig{
		Name:              "test-config",
		SelectionPolicy:   NewFIFOSelectionPolicy(),
		ConsumptionPolicy: NewOnceConsumptionPolicy(),
		RetentionPolicy:   NewUnlimitedRetentionPolicy(),
	}

	space := NewXupleSpace(config)

	retrievedConfig := space.GetConfig()

	if retrievedConfig.Name != "test-config" {
		t.Errorf("❌ Name incorrect: attendu 'test-config', reçu '%s'", retrievedConfig.Name)
	}

	if retrievedConfig.SelectionPolicy.Name() != "fifo" {
		t.Errorf("❌ SelectionPolicy incorrect: attendu 'fifo', reçu '%s'", retrievedConfig.SelectionPolicy.Name())
	}

	if retrievedConfig.ConsumptionPolicy.Name() != "once" {
		t.Errorf("❌ ConsumptionPolicy incorrect: attendu 'once', reçu '%s'", retrievedConfig.ConsumptionPolicy.Name())
	}

	if retrievedConfig.RetentionPolicy.Name() != "unlimited" {
		t.Errorf("❌ RetentionPolicy incorrect: attendu 'unlimited', reçu '%s'", retrievedConfig.RetentionPolicy.Name())
	}

	t.Log("✅ GetConfig fonctionne correctement")
}
