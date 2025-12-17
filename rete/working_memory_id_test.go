// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"testing"
)

// TestWorkingMemory_AddFactWithPKSimple teste l'ajout de fait avec ID basé sur PK simple
func TestWorkingMemory_AddFactWithPKSimple(t *testing.T) {
	t.Log("🧪 TEST: Working Memory - Ajout fait avec PK simple")
	t.Log("====================================================")

	wm := &WorkingMemory{
		NodeID: "test_node",
		Facts:  make(map[string]*Fact),
		Tokens: make(map[string]*Token),
	}

	fact := &Fact{
		ID:   "Person~Alice",
		Type: "Person",
		Fields: map[string]interface{}{
			"nom": "Alice",
			"age": 30,
		},
	}

	// Ajouter le fait
	err := wm.AddFact(fact)
	if err != nil {
		t.Fatalf("❌ Erreur lors de l'ajout du fait: %v", err)
	}

	// Vérifier que le fait est indexé avec l'ID interne
	internalID := "Person_Person~Alice"
	storedFact, exists := wm.GetFact(internalID)
	if !exists {
		t.Fatalf("❌ Fait non trouvé avec ID interne %s", internalID)
	}

	if storedFact.ID != "Person~Alice" {
		t.Errorf("❌ ID du fait attendu 'Person~Alice', reçu '%s'", storedFact.ID)
	}

	if storedFact.Type != "Person" {
		t.Errorf("❌ Type du fait attendu 'Person', reçu '%s'", storedFact.Type)
	}

	t.Log("✅ Test réussi: Fait ajouté avec PK simple")
}

// TestWorkingMemory_AddFactWithPKComposite teste l'ajout de fait avec ID basé sur PK composite
func TestWorkingMemory_AddFactWithPKComposite(t *testing.T) {
	t.Log("🧪 TEST: Working Memory - Ajout fait avec PK composite")
	t.Log("=======================================================")

	wm := &WorkingMemory{
		NodeID: "test_node",
		Facts:  make(map[string]*Fact),
		Tokens: make(map[string]*Token),
	}

	fact := &Fact{
		ID:   "Person~Alice_Dupont",
		Type: "Person",
		Fields: map[string]interface{}{
			"prenom": "Alice",
			"nom":    "Dupont",
			"age":    30,
		},
	}

	// Ajouter le fait
	err := wm.AddFact(fact)
	if err != nil {
		t.Fatalf("❌ Erreur lors de l'ajout du fait: %v", err)
	}

	// Vérifier que le fait est indexé avec l'ID interne
	internalID := "Person_Person~Alice_Dupont"
	storedFact, exists := wm.GetFact(internalID)
	if !exists {
		t.Fatalf("❌ Fait non trouvé avec ID interne %s", internalID)
	}

	if storedFact.ID != "Person~Alice_Dupont" {
		t.Errorf("❌ ID du fait attendu 'Person~Alice_Dupont', reçu '%s'", storedFact.ID)
	}

	t.Log("✅ Test réussi: Fait ajouté avec PK composite")
}

// TestWorkingMemory_AddFactWithHashID teste l'ajout de fait avec ID basé sur hash
func TestWorkingMemory_AddFactWithHashID(t *testing.T) {
	t.Log("🧪 TEST: Working Memory - Ajout fait avec hash ID")
	t.Log("==================================================")

	wm := &WorkingMemory{
		NodeID: "test_node",
		Facts:  make(map[string]*Fact),
		Tokens: make(map[string]*Token),
	}

	fact := &Fact{
		ID:   "Event~a1b2c3d4e5f6g7h8",
		Type: "Event",
		Fields: map[string]interface{}{
			"timestamp": 1234567890,
			"message":   "test event",
		},
	}

	// Ajouter le fait
	err := wm.AddFact(fact)
	if err != nil {
		t.Fatalf("❌ Erreur lors de l'ajout du fait: %v", err)
	}

	// Vérifier que le fait est indexé avec l'ID interne
	internalID := "Event_Event~a1b2c3d4e5f6g7h8"
	storedFact, exists := wm.GetFact(internalID)
	if !exists {
		t.Fatalf("❌ Fait non trouvé avec ID interne %s", internalID)
	}

	if storedFact.ID != "Event~a1b2c3d4e5f6g7h8" {
		t.Errorf("❌ ID du fait attendu 'Event~a1b2c3d4e5f6g7h8', reçu '%s'", storedFact.ID)
	}

	t.Log("✅ Test réussi: Fait ajouté avec hash ID")
}

// TestWorkingMemory_RemoveFactWithNewIDFormat teste la suppression de fait avec nouveau format d'ID
func TestWorkingMemory_RemoveFactWithNewIDFormat(t *testing.T) {
	t.Log("🧪 TEST: Working Memory - Suppression fait avec nouveau format ID")
	t.Log("==================================================================")

	wm := &WorkingMemory{
		NodeID: "test_node",
		Facts:  make(map[string]*Fact),
		Tokens: make(map[string]*Token),
	}

	fact := &Fact{
		ID:   "Person~Alice",
		Type: "Person",
		Fields: map[string]interface{}{
			"nom": "Alice",
		},
	}

	// Ajouter le fait
	err := wm.AddFact(fact)
	if err != nil {
		t.Fatalf("❌ Erreur lors de l'ajout du fait: %v", err)
	}

	internalID := "Person_Person~Alice"

	// Vérifier que le fait existe
	if _, exists := wm.GetFact(internalID); !exists {
		t.Fatalf("❌ Fait devrait être présent avant suppression")
	}

	// Supprimer le fait
	wm.RemoveFact(internalID)

	// Vérifier que le fait n'existe plus
	if _, exists := wm.GetFact(internalID); exists {
		t.Errorf("❌ Fait devrait être supprimé")
	}

	t.Log("✅ Test réussi: Fait supprimé avec nouveau format ID")
}

// TestWorkingMemory_GetFactByTypeAndID_NewIDFormats teste la récupération par type et ID avec nouveaux formats
func TestWorkingMemory_GetFactByTypeAndID_NewIDFormats(t *testing.T) {
	t.Log("🧪 TEST: Working Memory - Récupération par type et ID (nouveaux formats)")
	t.Log("=========================================================================")

	wm := &WorkingMemory{
		NodeID: "test_node",
		Facts:  make(map[string]*Fact),
		Tokens: make(map[string]*Token),
	}

	fact := &Fact{
		ID:   "Person~Alice_Dupont",
		Type: "Person",
		Fields: map[string]interface{}{
			"prenom": "Alice",
			"nom":    "Dupont",
		},
	}

	err := wm.AddFact(fact)
	if err != nil {
		t.Fatalf("❌ Erreur lors de l'ajout du fait: %v", err)
	}

	// Récupérer par type et ID
	retrievedFact, exists := wm.GetFactByTypeAndID("Person", "Person~Alice_Dupont")
	if !exists {
		t.Fatalf("❌ Fait non trouvé par type et ID")
	}

	if retrievedFact.ID != "Person~Alice_Dupont" {
		t.Errorf("❌ ID attendu 'Person~Alice_Dupont', reçu '%s'", retrievedFact.ID)
	}

	t.Log("✅ Test réussi: Récupération par type et ID")
}

// TestWorkingMemory_MultipleFactsDifferentTypes teste l'ajout de plusieurs faits de types différents
func TestWorkingMemory_MultipleFactsDifferentTypes(t *testing.T) {
	t.Log("🧪 TEST: Working Memory - Plusieurs faits de types différents")
	t.Log("==============================================================")

	wm := &WorkingMemory{
		NodeID: "test_node",
		Facts:  make(map[string]*Fact),
		Tokens: make(map[string]*Token),
	}

	facts := []*Fact{
		{
			ID:   "Person~Alice",
			Type: "Person",
			Fields: map[string]interface{}{
				"nom": "Alice",
			},
		},
		{
			ID:   "Order~Order123",
			Type: "Order",
			Fields: map[string]interface{}{
				"number": "Order123",
			},
		},
		{
			ID:   "Product~Laptop_15inch",
			Type: "Product",
			Fields: map[string]interface{}{
				"name": "Laptop",
				"size": "15inch",
			},
		},
	}

	// Ajouter tous les faits
	for _, fact := range facts {
		err := wm.AddFact(fact)
		if err != nil {
			t.Fatalf("❌ Erreur lors de l'ajout du fait %s: %v", fact.ID, err)
		}
	}

	// Vérifier que tous les faits sont présents
	expectedInternalIDs := []string{
		"Person_Person~Alice",
		"Order_Order~Order123",
		"Product_Product~Laptop_15inch",
	}

	for i, internalID := range expectedInternalIDs {
		storedFact, exists := wm.GetFact(internalID)
		if !exists {
			t.Errorf("❌ Fait non trouvé avec ID interne %s", internalID)
			continue
		}

		if storedFact.ID != facts[i].ID {
			t.Errorf("❌ ID attendu '%s', reçu '%s'", facts[i].ID, storedFact.ID)
		}

		if storedFact.Type != facts[i].Type {
			t.Errorf("❌ Type attendu '%s', reçu '%s'", facts[i].Type, storedFact.Type)
		}
	}

	// Vérifier le nombre total de faits
	allFacts := wm.GetFacts()
	if len(allFacts) != 3 {
		t.Errorf("❌ Attendu 3 faits, reçu %d", len(allFacts))
	}

	t.Log("✅ Test réussi: Plusieurs faits de types différents")
}

// TestWorkingMemory_DuplicateIDSameType teste que l'ajout d'un doublon échoue
func TestWorkingMemory_DuplicateIDSameType(t *testing.T) {
	t.Log("🧪 TEST: Working Memory - Rejet doublon même type")
	t.Log("=================================================")

	wm := &WorkingMemory{
		NodeID: "test_node",
		Facts:  make(map[string]*Fact),
		Tokens: make(map[string]*Token),
	}

	fact1 := &Fact{
		ID:   "Person~Alice",
		Type: "Person",
		Fields: map[string]interface{}{
			"nom": "Alice",
			"age": 30,
		},
	}

	fact2 := &Fact{
		ID:   "Person~Alice",
		Type: "Person",
		Fields: map[string]interface{}{
			"nom": "Alice",
			"age": 35, // Âge différent mais même ID
		},
	}

	// Ajouter le premier fait
	err := wm.AddFact(fact1)
	if err != nil {
		t.Fatalf("❌ Erreur lors de l'ajout du premier fait: %v", err)
	}

	// Tenter d'ajouter le doublon
	err = wm.AddFact(fact2)
	if err == nil {
		t.Fatalf("❌ L'ajout d'un doublon aurait dû échouer")
	}

	// Vérifier que seul le premier fait est présent
	internalID := "Person_Person~Alice"
	storedFact, exists := wm.GetFact(internalID)
	if !exists {
		t.Fatalf("❌ Fait non trouvé")
	}

	// Vérifier que c'est bien le premier fait (age = 30)
	if age, ok := storedFact.Fields["age"].(int); !ok || age != 30 {
		t.Errorf("❌ Le fait stocké devrait être le premier (age=30), reçu age=%v", storedFact.Fields["age"])
	}

	t.Log("✅ Test réussi: Doublon correctement rejeté")
}

// TestWorkingMemory_SameIDDifferentTypes teste que le même ID peut exister pour des types différents
func TestWorkingMemory_SameIDDifferentTypes(t *testing.T) {
	t.Log("🧪 TEST: Working Memory - Même ID pour types différents")
	t.Log("=========================================================")

	wm := &WorkingMemory{
		NodeID: "test_node",
		Facts:  make(map[string]*Fact),
		Tokens: make(map[string]*Token),
	}

	// Deux faits avec le même ID mais des types différents
	fact1 := &Fact{
		ID:   "Entity~123",
		Type: "Person",
		Fields: map[string]interface{}{
			"id": 123,
		},
	}

	fact2 := &Fact{
		ID:   "Entity~123",
		Type: "Company",
		Fields: map[string]interface{}{
			"id": 123,
		},
	}

	// Ajouter les deux faits
	err := wm.AddFact(fact1)
	if err != nil {
		t.Fatalf("❌ Erreur lors de l'ajout du premier fait: %v", err)
	}

	err = wm.AddFact(fact2)
	if err != nil {
		t.Fatalf("❌ Erreur lors de l'ajout du second fait: %v", err)
	}

	// Vérifier que les deux faits sont présents avec des IDs internes différents
	internalID1 := "Person_Entity~123"
	internalID2 := "Company_Entity~123"

	storedFact1, exists1 := wm.GetFact(internalID1)
	if !exists1 {
		t.Errorf("❌ Fait Person non trouvé")
	} else if storedFact1.Type != "Person" {
		t.Errorf("❌ Type attendu 'Person', reçu '%s'", storedFact1.Type)
	}

	storedFact2, exists2 := wm.GetFact(internalID2)
	if !exists2 {
		t.Errorf("❌ Fait Company non trouvé")
	} else if storedFact2.Type != "Company" {
		t.Errorf("❌ Type attendu 'Company', reçu '%s'", storedFact2.Type)
	}

	t.Log("✅ Test réussi: Même ID accepté pour types différents")
}

// TestWorkingMemory_ParseInternalID teste la décomposition des IDs internes
func TestWorkingMemory_ParseInternalID(t *testing.T) {
	t.Log("🧪 TEST: Working Memory - Parse Internal ID")
	t.Log("============================================")

	tests := []struct {
		name       string
		internalID string
		wantType   string
		wantID     string
		wantOK     bool
	}{
		{
			name:       "ID simple",
			internalID: "Person_Person~Alice",
			wantType:   "Person",
			wantID:     "Person~Alice",
			wantOK:     true,
		},
		{
			name:       "ID composite",
			internalID: "Person_Person~Alice_Dupont",
			wantType:   "Person",
			wantID:     "Person~Alice_Dupont",
			wantOK:     true,
		},
		{
			name:       "ID hash",
			internalID: "Event_Event~a1b2c3d4e5f6g7h8",
			wantType:   "Event",
			wantID:     "Event~a1b2c3d4e5f6g7h8",
			wantOK:     true,
		},
		{
			name:       "Format invalide",
			internalID: "InvalidFormat",
			wantType:   "",
			wantID:     "",
			wantOK:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotType, gotID, gotOK := ParseInternalID(tt.internalID)

			if gotOK != tt.wantOK {
				t.Errorf("❌ OK = %v, attendu %v", gotOK, tt.wantOK)
			}

			if gotType != tt.wantType {
				t.Errorf("❌ Type = '%s', attendu '%s'", gotType, tt.wantType)
			}

			if gotID != tt.wantID {
				t.Errorf("❌ ID = '%s', attendu '%s'", gotID, tt.wantID)
			}

			t.Log("✅ Test réussi")
		})
	}

	t.Log("✅ Test complet: Parse Internal ID")
}

// TestWorkingMemory_MakeInternalID teste la construction des IDs internes
func TestWorkingMemory_MakeInternalID(t *testing.T) {
	t.Log("🧪 TEST: Working Memory - Make Internal ID")
	t.Log("===========================================")

	tests := []struct {
		name     string
		factType string
		factID   string
		want     string
	}{
		{
			name:     "ID simple",
			factType: "Person",
			factID:   "Person~Alice",
			want:     "Person_Person~Alice",
		},
		{
			name:     "ID composite",
			factType: "Person",
			factID:   "Person~Alice_Dupont",
			want:     "Person_Person~Alice_Dupont",
		},
		{
			name:     "ID hash",
			factType: "Event",
			factID:   "Event~a1b2c3d4e5f6g7h8",
			want:     "Event_Event~a1b2c3d4e5f6g7h8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MakeInternalID(tt.factType, tt.factID)

			if got != tt.want {
				t.Errorf("❌ MakeInternalID() = '%s', attendu '%s'", got, tt.want)
			}

			t.Log("✅ Test réussi")
		})
	}

	t.Log("✅ Test complet: Make Internal ID")
}
