// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package xuples

import (
	"testing"
	"time"

	"github.com/treivax/tsd/rete"
)

// TestMaxSizeEnforcement teste que la limite MaxSize est respectée
func TestMaxSizeEnforcement(t *testing.T) {
	t.Log("🧪 TEST MAX SIZE ENFORCEMENT")
	t.Log("=============================")

	// Configuration avec MaxSize = 2
	config := XupleSpaceConfig{
		Name:              "test-space",
		SelectionPolicy:   NewFIFOSelectionPolicy(),
		ConsumptionPolicy: NewOnceConsumptionPolicy(),
		RetentionPolicy:   NewUnlimitedRetentionPolicy(),
		MaxSize:           2,
	}

	space := NewXupleSpace(config)

	// Créer 3 xuples
	xuple1 := &Xuple{
		ID:              "xuple-1",
		Fact:            &rete.Fact{},
		TriggeringFacts: []*rete.Fact{},
		CreatedAt:       time.Now(),
		Metadata: XupleMetadata{
			State:      XupleStateAvailable,
			ConsumedBy: make(map[string]time.Time),
		},
	}

	xuple2 := &Xuple{
		ID:              "xuple-2",
		Fact:            &rete.Fact{},
		TriggeringFacts: []*rete.Fact{},
		CreatedAt:       time.Now(),
		Metadata: XupleMetadata{
			State:      XupleStateAvailable,
			ConsumedBy: make(map[string]time.Time),
		},
	}

	xuple3 := &Xuple{
		ID:              "xuple-3",
		Fact:            &rete.Fact{},
		TriggeringFacts: []*rete.Fact{},
		CreatedAt:       time.Now(),
		Metadata: XupleMetadata{
			State:      XupleStateAvailable,
			ConsumedBy: make(map[string]time.Time),
		},
	}

	// Insert xuple1 - doit réussir
	err := space.Insert(xuple1)
	if err != nil {
		t.Fatalf("❌ Insert xuple1 échoué: %v", err)
	}
	t.Log("✅ Insert xuple1 réussi (1/2)")

	// Insert xuple2 - doit réussir
	err = space.Insert(xuple2)
	if err != nil {
		t.Fatalf("❌ Insert xuple2 échoué: %v", err)
	}
	t.Log("✅ Insert xuple2 réussi (2/2)")

	// Insert xuple3 - doit échouer (MaxSize atteint)
	err = space.Insert(xuple3)
	if err != ErrXupleSpaceFull {
		t.Fatalf("❌ Attendu ErrXupleSpaceFull, reçu: %v", err)
	}
	t.Log("✅ Insert xuple3 rejeté (MaxSize atteint)")

	// Vérifier le count
	count := space.Count()
	if count != 2 {
		t.Errorf("❌ Attendu count=2, reçu: %d", count)
	}

	t.Log("✅ MaxSize enforcement fonctionne correctement")
}

// TestMaxSizeZeroUnlimited teste que MaxSize=0 signifie illimité
func TestMaxSizeZeroUnlimited(t *testing.T) {
	t.Log("🧪 TEST MAX SIZE ZERO (UNLIMITED)")
	t.Log("==================================")

	config := XupleSpaceConfig{
		Name:              "unlimited-space",
		SelectionPolicy:   NewFIFOSelectionPolicy(),
		ConsumptionPolicy: NewOnceConsumptionPolicy(),
		RetentionPolicy:   NewUnlimitedRetentionPolicy(),
		MaxSize:           0, // Illimité
	}

	space := NewXupleSpace(config)

	// Insérer 100 xuples - tous doivent réussir
	const numXuples = 100
	for i := 0; i < numXuples; i++ {
		xuple := &Xuple{
			ID:              string(rune('a' + i)),
			Fact:            &rete.Fact{},
			TriggeringFacts: []*rete.Fact{},
			CreatedAt:       time.Now(),
			Metadata: XupleMetadata{
				State:      XupleStateAvailable,
				ConsumedBy: make(map[string]time.Time),
			},
		}

		err := space.Insert(xuple)
		if err != nil {
			t.Fatalf("❌ Insert xuple %d échoué: %v", i, err)
		}
	}

	count := space.Count()
	if count != numXuples {
		t.Errorf("❌ Attendu count=%d, reçu: %d", numXuples, count)
	}

	t.Logf("✅ Insertion de %d xuples avec MaxSize=0 réussie", numXuples)
}

// TestUnlimitedRetentionCleansConsumed teste que UnlimitedRetentionPolicy nettoie les xuples consommés
func TestUnlimitedRetentionCleansConsumed(t *testing.T) {
	t.Log("🧪 TEST UNLIMITED RETENTION CLEANS CONSUMED")
	t.Log("============================================")

	config := XupleSpaceConfig{
		Name:              "cleanup-space",
		SelectionPolicy:   NewFIFOSelectionPolicy(),
		ConsumptionPolicy: NewOnceConsumptionPolicy(),
		RetentionPolicy:   NewUnlimitedRetentionPolicy(),
		MaxSize:           0,
	}

	space := NewXupleSpace(config)

	// Créer et insérer un xuple
	xuple := &Xuple{
		ID:              "xuple-consumed",
		Fact:            &rete.Fact{},
		TriggeringFacts: []*rete.Fact{},
		CreatedAt:       time.Now(),
		Metadata: XupleMetadata{
			State:      XupleStateAvailable,
			ConsumedBy: make(map[string]time.Time),
		},
	}

	err := space.Insert(xuple)
	if err != nil {
		t.Fatalf("❌ Insert échoué: %v", err)
	}

	// Marquer comme consommé
	err = space.MarkConsumed(xuple.ID, "agent-1")
	if err != nil {
		t.Fatalf("❌ MarkConsumed échoué: %v", err)
	}

	// Vérifier que le xuple est marqué comme consommé
	if xuple.Metadata.State != XupleStateConsumed {
		t.Errorf("❌ Xuple devrait être consommé, état: %s", xuple.Metadata.State)
	}
	t.Log("✅ Xuple marqué comme consommé")

	// Cleanup devrait retirer le xuple consommé
	cleaned := space.Cleanup()
	if cleaned != 1 {
		t.Errorf("❌ Attendu 1 xuple nettoyé, reçu: %d", cleaned)
	}

	count := space.Count()
	if count != 0 {
		t.Errorf("❌ Attendu count=0 après cleanup, reçu: %d", count)
	}

	t.Log("✅ Unlimited retention policy nettoie les xuples consommés")
}

// TestInsertWithoutID teste que Insert rejette les xuples sans ID
func TestInsertWithoutID(t *testing.T) {
	t.Log("🧪 TEST INSERT WITHOUT ID")
	t.Log("=========================")

	config := XupleSpaceConfig{
		Name:              "test-space",
		SelectionPolicy:   NewFIFOSelectionPolicy(),
		ConsumptionPolicy: NewOnceConsumptionPolicy(),
		RetentionPolicy:   NewUnlimitedRetentionPolicy(),
		MaxSize:           0,
	}

	space := NewXupleSpace(config)

	// Xuple sans ID
	xuple := &Xuple{
		ID:              "", // Vide
		Fact:            &rete.Fact{},
		TriggeringFacts: []*rete.Fact{},
		CreatedAt:       time.Now(),
		Metadata: XupleMetadata{
			State:      XupleStateAvailable,
			ConsumedBy: make(map[string]time.Time),
		},
	}

	err := space.Insert(xuple)
	if err != ErrInvalidConfiguration {
		t.Fatalf("❌ Attendu ErrInvalidConfiguration, reçu: %v", err)
	}

	t.Log("✅ Insert rejette correctement les xuples sans ID")
}
