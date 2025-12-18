// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package actions

import (
	"bytes"
	"log"
	"testing"
	"time"

	"github.com/treivax/tsd/rete"
	"github.com/treivax/tsd/xuples"
)

func TestBuiltinActions_EndToEnd_DynamicFactOperations(t *testing.T) {
	t.Log("🧪 TEST End-to-End - Actions dynamiques Insert, Update, Retract")
	t.Log("==================================================================")

	// Setup
	storage := rete.NewMemoryStorage()
	network := rete.NewReteNetwork(storage)
	output := &bytes.Buffer{}
	logOutput := &bytes.Buffer{}
	logger := log.New(logOutput, "", 0)
	executor := NewBuiltinActionExecutor(network, nil, output, logger)

	// Scénario: Gestion du cycle de vie d'un utilisateur
	// 1. Insert : Créer un nouvel utilisateur
	// 2. Update : Promouvoir l'utilisateur
	// 3. Print : Afficher le statut
	// 4. Retract : Supprimer l'utilisateur

	t.Log("📝 Étape 1 : Insertion d'un nouvel utilisateur")
	newUser := &rete.Fact{
		ID:   "user001",
		Type: "User",
		Fields: map[string]interface{}{
			"name":   "Alice",
			"role":   "developer",
			"active": true,
		},
	}

	err := executor.Execute("Insert", []interface{}{newUser}, &rete.Token{})
	if err != nil {
		t.Fatalf("❌ Insert failed: %v", err)
	}

	// Vérifier l'insertion
	storedUser := storage.GetFact("User_user001")
	if storedUser == nil {
		t.Fatal("❌ User not found after insert")
	}
	if storedUser.Fields["role"] != "developer" {
		t.Errorf("❌ Expected role 'developer', got '%v'", storedUser.Fields["role"])
	}
	t.Log("✅ Utilisateur inséré avec succès")

	t.Log("📝 Étape 2 : Promotion de l'utilisateur (Update)")
	promotedUser := &rete.Fact{
		ID:   "user001",
		Type: "User",
		Fields: map[string]interface{}{
			"name":   "Alice",
			"role":   "senior_developer",
			"active": true,
		},
	}

	err = executor.Execute("Update", []interface{}{promotedUser}, &rete.Token{})
	if err != nil {
		t.Fatalf("❌ Update failed: %v", err)
	}

	// Vérifier la mise à jour
	storedUser = storage.GetFact("User_user001")
	if storedUser == nil {
		t.Fatal("❌ User not found after update")
	}
	if storedUser.Fields["role"] != "senior_developer" {
		t.Errorf("❌ Expected role 'senior_developer', got '%v'", storedUser.Fields["role"])
	}
	t.Log("✅ Utilisateur promu avec succès")

	t.Log("📝 Étape 3 : Affichage du statut (Print)")
	message := "User " + storedUser.Fields["name"].(string) + " is now " + storedUser.Fields["role"].(string)
	err = executor.Execute("Print", []interface{}{message}, &rete.Token{})
	if err != nil {
		t.Fatalf("❌ Print failed: %v", err)
	}

	// Vérifier l'affichage
	printedText := output.String()
	if printedText != message+"\n" {
		t.Errorf("❌ Expected '%s', got '%s'", message+"\n", printedText)
	}
	t.Log("✅ Statut affiché avec succès")

	t.Log("📝 Étape 4 : Logging de l'opération (Log)")
	logMessage := "User promotion completed: user001"
	err = executor.Execute("Log", []interface{}{logMessage}, &rete.Token{})
	if err != nil {
		t.Fatalf("❌ Log failed: %v", err)
	}

	// Vérifier le log
	loggedText := logOutput.String()
	if loggedText == "" || len(loggedText) < len(logMessage) {
		t.Errorf("❌ Log should contain '%s', got '%s'", logMessage, loggedText)
	}
	t.Log("✅ Opération loggée avec succès")

	t.Log("📝 Étape 5 : Suppression de l'utilisateur (Retract)")
	err = executor.Execute("Retract", []interface{}{"User_user001"}, &rete.Token{})
	if err != nil {
		t.Fatalf("❌ Retract failed: %v", err)
	}

	// Vérifier la suppression
	storedUser = storage.GetFact("User_user001")
	if storedUser != nil {
		t.Error("❌ User should have been removed")
	}
	t.Log("✅ Utilisateur supprimé avec succès")

	t.Log("🎉 Test end-to-end réussi - Toutes les actions dynamiques fonctionnent correctement")
}

func TestBuiltinActions_EndToEnd_ComplexScenario(t *testing.T) {
	t.Log("🧪 TEST End-to-End - Scénario complexe avec multiples actions")
	t.Log("=================================================================")

	// Setup
	storage := rete.NewMemoryStorage()
	network := rete.NewReteNetwork(storage)
	output := &bytes.Buffer{}
	executor := NewBuiltinActionExecutor(network, nil, output, nil)

	// Scénario: Système de gestion de commandes
	// 1. Insert : Créer une nouvelle commande
	// 2. Insert : Créer les items de la commande
	// 3. Update : Mettre à jour le statut de la commande
	// 4. Print : Afficher la confirmation
	// 5. Retract : Annuler un item
	// 6. Update : Recalculer le total

	t.Log("📝 Étape 1 : Création d'une commande")
	order := &rete.Fact{
		ID:   "order001",
		Type: "Order",
		Fields: map[string]interface{}{
			"customerId": "cust123",
			"status":     "pending",
			"total":      0.0,
		},
	}

	err := executor.Execute("Insert", []interface{}{order}, &rete.Token{})
	if err != nil {
		t.Fatalf("❌ Insert order failed: %v", err)
	}
	t.Log("✅ Commande créée")

	t.Log("📝 Étape 2 : Ajout d'items à la commande")
	item1 := &rete.Fact{
		ID:   "item001",
		Type: "OrderItem",
		Fields: map[string]interface{}{
			"orderId": "order001",
			"product": "Product A",
			"price":   29.99,
		},
	}
	item2 := &rete.Fact{
		ID:   "item002",
		Type: "OrderItem",
		Fields: map[string]interface{}{
			"orderId": "order001",
			"product": "Product B",
			"price":   49.99,
		},
	}

	err = executor.Execute("Insert", []interface{}{item1}, &rete.Token{})
	if err != nil {
		t.Fatalf("❌ Insert item1 failed: %v", err)
	}
	err = executor.Execute("Insert", []interface{}{item2}, &rete.Token{})
	if err != nil {
		t.Fatalf("❌ Insert item2 failed: %v", err)
	}
	t.Log("✅ Items ajoutés")

	// Vérifier que tous les faits sont dans le storage
	allFacts := storage.GetAllFacts()
	expectedFactCount := 3 // 1 order + 2 items
	if len(allFacts) != expectedFactCount {
		t.Errorf("❌ Expected %d facts, got %d", expectedFactCount, len(allFacts))
	}

	t.Log("📝 Étape 3 : Mise à jour du statut et total de la commande")
	updatedOrder := &rete.Fact{
		ID:   "order001",
		Type: "Order",
		Fields: map[string]interface{}{
			"customerId": "cust123",
			"status":     "processing",
			"total":      79.98, // 29.99 + 49.99
		},
	}

	err = executor.Execute("Update", []interface{}{updatedOrder}, &rete.Token{})
	if err != nil {
		t.Fatalf("❌ Update order failed: %v", err)
	}

	// Vérifier la mise à jour
	storedOrder := storage.GetFact("Order_order001")
	if storedOrder == nil {
		t.Fatal("❌ Order not found after update")
	}
	if storedOrder.Fields["status"] != "processing" {
		t.Errorf("❌ Expected status 'processing', got '%v'", storedOrder.Fields["status"])
	}
	if storedOrder.Fields["total"] != 79.98 {
		t.Errorf("❌ Expected total 79.98, got %v", storedOrder.Fields["total"])
	}
	t.Log("✅ Commande mise à jour")

	t.Log("📝 Étape 4 : Affichage de confirmation")
	err = executor.Execute("Print", []interface{}{"Order order001 is now processing"}, &rete.Token{})
	if err != nil {
		t.Fatalf("❌ Print failed: %v", err)
	}
	t.Log("✅ Confirmation affichée")

	t.Log("📝 Étape 5 : Annulation d'un item")
	err = executor.Execute("Retract", []interface{}{"OrderItem_item002"}, &rete.Token{})
	if err != nil {
		t.Fatalf("❌ Retract item failed: %v", err)
	}

	// Vérifier la suppression
	if storage.GetFact("OrderItem_item002") != nil {
		t.Error("❌ Item should have been removed")
	}
	t.Log("✅ Item annulé")

	t.Log("📝 Étape 6 : Recalcul du total")
	recalculatedOrder := &rete.Fact{
		ID:   "order001",
		Type: "Order",
		Fields: map[string]interface{}{
			"customerId": "cust123",
			"status":     "processing",
			"total":      29.99, // Seulement item1 restant
		},
	}

	err = executor.Execute("Update", []interface{}{recalculatedOrder}, &rete.Token{})
	if err != nil {
		t.Fatalf("❌ Update order total failed: %v", err)
	}

	// Vérification finale
	finalOrder := storage.GetFact("Order_order001")
	if finalOrder == nil {
		t.Fatal("❌ Order not found")
	}
	if finalOrder.Fields["total"] != 29.99 {
		t.Errorf("❌ Expected total 29.99, got %v", finalOrder.Fields["total"])
	}

	// Vérifier le nombre final de faits
	allFacts = storage.GetAllFacts()
	expectedFinalCount := 2 // 1 order + 1 item restant
	if len(allFacts) != expectedFinalCount {
		t.Errorf("❌ Expected %d facts, got %d", expectedFinalCount, len(allFacts))
	}

	t.Log("✅ Total recalculé")
	t.Log("🎉 Scénario complexe réussi - Toutes les opérations CRUD fonctionnent")
}

func TestBuiltinActions_ErrorHandling(t *testing.T) {
	t.Log("🧪 TEST End-to-End - Gestion des erreurs")
	t.Log("==========================================")

	storage := rete.NewMemoryStorage()
	network := rete.NewReteNetwork(storage)
	executor := NewBuiltinActionExecutor(network, nil, nil, nil)

	t.Log("📝 Test 1 : Update sur un fait inexistant")
	nonExistentFact := &rete.Fact{
		ID:   "ghost",
		Type: "User",
		Fields: map[string]interface{}{
			"name": "Ghost User",
		},
	}

	err := executor.Execute("Update", []interface{}{nonExistentFact}, &rete.Token{})
	if err == nil {
		t.Error("❌ Expected error when updating non-existent fact")
	}
	t.Logf("✅ Erreur attendue reçue: %v", err)

	t.Log("📝 Test 2 : Insert d'un fait déjà existant")
	existingFact := &rete.Fact{
		ID:   "user001",
		Type: "User",
		Fields: map[string]interface{}{
			"name": "Alice",
		},
	}

	// Premier insert (doit réussir)
	err = executor.Execute("Insert", []interface{}{existingFact}, &rete.Token{})
	if err != nil {
		t.Fatalf("❌ First insert failed: %v", err)
	}

	// Second insert (doit échouer)
	err = executor.Execute("Insert", []interface{}{existingFact}, &rete.Token{})
	if err == nil {
		t.Error("❌ Expected error when inserting duplicate fact")
	}
	t.Logf("✅ Erreur attendue reçue: %v", err)

	t.Log("📝 Test 3 : Retract sur un fait inexistant")
	err = executor.Execute("Retract", []interface{}{"User_ghost"}, &rete.Token{})
	if err == nil {
		t.Error("❌ Expected error when retracting non-existent fact")
	}
	t.Logf("✅ Erreur attendue reçue: %v", err)

	t.Log("🎉 Gestion des erreurs validée")
}

func TestBuiltinActions_EndToEnd_XupleAction(t *testing.T) {
	t.Log("🧪 TEST End-to-End - Action Xuple avec xuple-spaces")
	t.Log("======================================================")

	// Setup
	storage := rete.NewMemoryStorage()
	network := rete.NewReteNetwork(storage)
	xupleManager := xuples.NewXupleManager()
	executor := NewBuiltinActionExecutor(network, xupleManager, nil, nil)

	// Créer les xuple-spaces
	t.Log("📝 Étape 1 : Création des xuple-spaces")

	// Xuple-space pour alertes critiques (LIFO, per-agent)
	criticalConfig := xuples.XupleSpaceConfig{
		Name:              "critical-alerts",
		SelectionPolicy:   xuples.NewLIFOSelectionPolicy(),
		ConsumptionPolicy: xuples.NewPerAgentConsumptionPolicy(),
		RetentionPolicy:   xuples.NewDurationRetentionPolicy(10 * time.Minute),
		MaxSize:           0,
	}
	err := xupleManager.CreateXupleSpace("critical-alerts", criticalConfig)
	if err != nil {
		t.Fatalf("❌ Failed to create critical-alerts space: %v", err)
	}

	// Xuple-space pour commandes (FIFO, once)
	commandConfig := xuples.XupleSpaceConfig{
		Name:              "command-queue",
		SelectionPolicy:   xuples.NewFIFOSelectionPolicy(),
		ConsumptionPolicy: xuples.NewOnceConsumptionPolicy(),
		RetentionPolicy:   xuples.NewDurationRetentionPolicy(1 * time.Hour),
		MaxSize:           0,
	}
	err = xupleManager.CreateXupleSpace("command-queue", commandConfig)
	if err != nil {
		t.Fatalf("❌ Failed to create command-queue space: %v", err)
	}

	t.Log("✅ Xuple-spaces créés avec succès")

	// Scénario: Système de monitoring de capteurs
	t.Log("📝 Étape 2 : Création de xuples via l'action Xuple")

	// Créer des faits déclencheurs (sensors)
	sensor1 := &rete.Fact{
		ID:   "S001",
		Type: "Sensor",
		Fields: map[string]interface{}{
			"location":    "Room-A",
			"temperature": 45.0,
			"humidity":    60.0,
		},
	}

	sensor2 := &rete.Fact{
		ID:   "S002",
		Type: "Sensor",
		Fields: map[string]interface{}{
			"location":    "Room-B",
			"temperature": 35.0,
			"humidity":    50.0,
		},
	}

	// Token avec faits déclencheurs
	token := &rete.Token{
		Facts: []*rete.Fact{sensor1, sensor2},
	}

	// Créer une alerte critique via Xuple
	alert1 := &rete.Fact{
		ID:   "A001",
		Type: "Alert",
		Fields: map[string]interface{}{
			"level":    "CRITICAL",
			"message":  "Temperature critical at Room-A",
			"sensorId": "S001",
		},
	}

	err = executor.Execute("Xuple", []interface{}{"critical-alerts", alert1}, token)
	if err != nil {
		t.Fatalf("❌ Failed to create xuple in critical-alerts: %v", err)
	}
	t.Log("✅ Alerte critique créée dans critical-alerts")

	// Créer une deuxième alerte critique
	alert2 := &rete.Fact{
		ID:   "A002",
		Type: "Alert",
		Fields: map[string]interface{}{
			"level":    "CRITICAL",
			"message":  "System overload detected",
			"sensorId": "S002",
		},
	}

	err = executor.Execute("Xuple", []interface{}{"critical-alerts", alert2}, token)
	if err != nil {
		t.Fatalf("❌ Failed to create second xuple in critical-alerts: %v", err)
	}
	t.Log("✅ Deuxième alerte critique créée dans critical-alerts")

	// Créer des commandes via Xuple
	command1 := &rete.Fact{
		ID:   "C001",
		Type: "Command",
		Fields: map[string]interface{}{
			"action":   "activate_cooling",
			"target":   "Room-A",
			"priority": 10,
		},
	}

	err = executor.Execute("Xuple", []interface{}{"command-queue", command1}, token)
	if err != nil {
		t.Fatalf("❌ Failed to create xuple in command-queue: %v", err)
	}
	t.Log("✅ Commande créée dans command-queue")

	command2 := &rete.Fact{
		ID:   "C002",
		Type: "Command",
		Fields: map[string]interface{}{
			"action":   "send_notification",
			"target":   "admin",
			"priority": 5,
		},
	}

	err = executor.Execute("Xuple", []interface{}{"command-queue", command2}, token)
	if err != nil {
		t.Fatalf("❌ Failed to create second xuple in command-queue: %v", err)
	}
	t.Log("✅ Deuxième commande créée dans command-queue")

	// Vérifier le contenu des xuple-spaces
	t.Log("📝 Étape 3 : Vérification du contenu des xuple-spaces")

	// Vérifier critical-alerts
	criticalSpace, err := xupleManager.GetXupleSpace("critical-alerts")
	if err != nil {
		t.Fatalf("❌ Failed to get critical-alerts space: %v", err)
	}

	criticalXuples := criticalSpace.ListAll()
	if len(criticalXuples) != 2 {
		t.Errorf("❌ Expected 2 xuples in critical-alerts, got %d", len(criticalXuples))
	} else {
		t.Logf("✅ critical-alerts contient %d xuples", len(criticalXuples))

		// Vérifier les détails des xuples
		for i, xuple := range criticalXuples {
			t.Logf("   Xuple %d: ID=%s, Type=%s, State=%s",
				i+1, xuple.Fact.ID, xuple.Fact.Type, xuple.Metadata.State)

			if xuple.Fact.Type != "Alert" {
				t.Errorf("❌ Expected Alert type, got %s", xuple.Fact.Type)
			}

			// Vérifier que les faits déclencheurs sont présents
			if len(xuple.TriggeringFacts) != 2 {
				t.Errorf("❌ Expected 2 triggering facts, got %d", len(xuple.TriggeringFacts))
			}
		}
	}

	// Vérifier command-queue
	commandSpace, err := xupleManager.GetXupleSpace("command-queue")
	if err != nil {
		t.Fatalf("❌ Failed to get command-queue space: %v", err)
	}

	commandXuples := commandSpace.ListAll()
	if len(commandXuples) != 2 {
		t.Errorf("❌ Expected 2 xuples in command-queue, got %d", len(commandXuples))
	} else {
		t.Logf("✅ command-queue contient %d xuples", len(commandXuples))

		// Vérifier les détails
		for i, xuple := range commandXuples {
			t.Logf("   Xuple %d: ID=%s, Type=%s, Action=%s, Priority=%v",
				i+1, xuple.Fact.ID, xuple.Fact.Type,
				xuple.Fact.Fields["action"], xuple.Fact.Fields["priority"])

			if xuple.Fact.Type != "Command" {
				t.Errorf("❌ Expected Command type, got %s", xuple.Fact.Type)
			}
		}
	}

	// Tester la récupération avec politiques
	t.Log("📝 Étape 4 : Test de récupération avec politiques")

	// Récupérer depuis critical-alerts (LIFO + per-agent)
	retrievedAlert, err := criticalSpace.Retrieve("agent1")
	if err != nil {
		t.Errorf("❌ Failed to retrieve from critical-alerts: %v", err)
	} else {
		t.Logf("✅ Agent1 a récupéré alerte: %s (LIFO: devrait être la dernière créée)", retrievedAlert.Fact.ID)
		// LIFO devrait retourner A002 (dernière créée)
		if retrievedAlert.Fact.ID != "A002" {
			t.Logf("⚠️  LIFO policy: expected A002, got %s", retrievedAlert.Fact.ID)
		}
	}

	// Marquer comme consommée
	if retrievedAlert != nil {
		err = criticalSpace.MarkConsumed(retrievedAlert.ID, "agent1")
		if err != nil {
			t.Errorf("❌ Failed to mark consumed: %v", err)
		} else {
			t.Log("✅ Alerte marquée comme consommée par agent1")
		}

		// Per-agent policy: un autre agent devrait pouvoir récupérer le même xuple
		retrievedAlert2, err := criticalSpace.Retrieve("agent2")
		if err != nil {
			t.Errorf("❌ Failed to retrieve for agent2: %v", err)
		} else {
			t.Logf("✅ Agent2 a récupéré alerte: %s (per-agent policy fonctionne)", retrievedAlert2.Fact.ID)
		}
	}

	// Récupérer depuis command-queue (FIFO + once)
	retrievedCmd, err := commandSpace.Retrieve("agent1")
	if err != nil {
		t.Errorf("❌ Failed to retrieve from command-queue: %v", err)
	} else {
		t.Logf("✅ Agent1 a récupéré commande: %s (FIFO: devrait être la première créée)", retrievedCmd.Fact.ID)
		// FIFO devrait retourner C001 (première créée)
		if retrievedCmd.Fact.ID != "C001" {
			t.Logf("⚠️  FIFO policy: expected C001, got %s", retrievedCmd.Fact.ID)
		}
	}

	// Test d'erreur: xuple-space inexistant
	t.Log("📝 Étape 5 : Test de gestion d'erreurs")

	fakeFact := &rete.Fact{ID: "F001", Type: "Fake"}
	err = executor.Execute("Xuple", []interface{}{"non-existent-space", fakeFact}, token)
	if err == nil {
		t.Error("❌ Expected error for non-existent xuple-space")
	} else {
		t.Logf("✅ Erreur attendue pour xuple-space inexistant: %v", err)
	}

	// Test statistiques
	t.Log("📝 Étape 6 : Statistiques des xuple-spaces")

	criticalCount := criticalSpace.Count()
	commandCount := commandSpace.Count()

	t.Logf("✅ critical-alerts: %d xuples disponibles", criticalCount)
	t.Logf("✅ command-queue: %d xuples disponibles", commandCount)

	spaces := xupleManager.ListXupleSpaces()
	t.Logf("✅ Nombre total de xuple-spaces: %d", len(spaces))
	for _, name := range spaces {
		t.Logf("   - %s", name)
	}

	t.Log("🎉 Tests de l'action Xuple validés avec succès!")
}
