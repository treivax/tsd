// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"testing"
)

// TestXupleActionAutomaticRegistration vérifie que l'action Xuple est automatiquement
// enregistrée lorsqu'un handler est configuré.
func TestXupleActionAutomaticRegistration(t *testing.T) {
	t.Log("🧪 TEST: Enregistrement automatique de l'action Xuple")
	t.Log("=========================================================")

	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)

	// Mock XupleHandler
	network.SetXupleHandler(func(xuplespace string, fact *Fact, triggeringFacts []*Fact) error {
		t.Logf("   Handler Xuple appelé pour space: %s", xuplespace)
		return nil
	})

	// Créer l'ActionExecutor APRÈS avoir configuré le handler
	executor := NewActionExecutor(network, nil)
	network.ActionExecutor = executor

	// Vérifier que l'action Xuple est enregistrée
	registry := executor.GetRegistry()
	if !registry.Has("Xuple") {
		t.Fatal("❌ L'action Xuple n'a pas été automatiquement enregistrée")
	}

	t.Log("✅ Action Xuple automatiquement enregistrée lors de la création de l'executor")

	// Vérifier que l'action peut être récupérée
	handler := registry.Get("Xuple")
	if handler == nil {
		t.Fatal("❌ Impossible de récupérer le handler Xuple")
	}

	if handler.GetName() != "Xuple" {
		t.Errorf("❌ Nom du handler incorrect: attendu 'Xuple', reçu '%s'", handler.GetName())
	}

	t.Log("✅ Handler Xuple récupérable et correct")
}

// TestXupleActionLateRegistration vérifie que l'action Xuple peut être enregistrée
// après la création de l'executor via RegisterXupleActionIfNeeded.
func TestXupleActionLateRegistration(t *testing.T) {
	t.Log("🧪 TEST: Enregistrement tardif de l'action Xuple")
	t.Log("================================================")

	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)

	// Créer l'ActionExecutor SANS handler configuré
	executor := NewActionExecutor(network, nil)
	network.ActionExecutor = executor

	// Vérifier que l'action Xuple n'est PAS enregistrée
	registry := executor.GetRegistry()
	if registry.Has("Xuple") {
		t.Fatal("❌ L'action Xuple ne devrait pas être enregistrée sans handler")
	}

	t.Log("✅ Action Xuple non enregistrée (pas de handler)")

	// Configurer le handler APRÈS
	network.SetXupleHandler(func(xuplespace string, fact *Fact, triggeringFacts []*Fact) error {
		return nil
	})

	// Enregistrer l'action manuellement
	err := executor.RegisterXupleActionIfNeeded()
	if err != nil {
		t.Fatalf("❌ Erreur lors de l'enregistrement: %v", err)
	}

	// Vérifier que l'action est maintenant enregistrée
	if !registry.Has("Xuple") {
		t.Fatal("❌ L'action Xuple n'a pas été enregistrée via RegisterXupleActionIfNeeded")
	}

	t.Log("✅ Action Xuple enregistrée avec succès via RegisterXupleActionIfNeeded")

	// Appeler RegisterXupleActionIfNeeded une deuxième fois (devrait être idempotent)
	err = executor.RegisterXupleActionIfNeeded()
	if err != nil {
		t.Errorf("❌ RegisterXupleActionIfNeeded devrait être idempotent, erreur: %v", err)
	}

	t.Log("✅ RegisterXupleActionIfNeeded est idempotent")
}

// TestXupleActionWithoutHandler vérifie le comportement quand aucun handler n'est configuré.
func TestXupleActionWithoutHandler(t *testing.T) {
	t.Log("🧪 TEST: Action Xuple sans handler")
	t.Log("===================================")

	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)

	// Créer l'ActionExecutor sans handler
	executor := NewActionExecutor(network, nil)
	network.ActionExecutor = executor

	// Vérifier qu'aucune action Xuple n'est enregistrée
	registry := executor.GetRegistry()
	if registry.Has("Xuple") {
		t.Fatal("❌ L'action Xuple ne devrait pas être enregistrée sans handler")
	}

	t.Log("✅ Aucune action Xuple enregistrée sans handler (comportement attendu)")

	// Appeler RegisterXupleActionIfNeeded sans handler (devrait ne rien faire)
	err := executor.RegisterXupleActionIfNeeded()
	if err != nil {
		t.Errorf("❌ RegisterXupleActionIfNeeded devrait retourner nil sans handler, erreur: %v", err)
	}

	if registry.Has("Xuple") {
		t.Fatal("❌ L'action Xuple ne devrait toujours pas être enregistrée")
	}

	t.Log("✅ RegisterXupleActionIfNeeded ne fait rien sans handler (comportement attendu)")
}

// TestXupleActionValidation vérifie la validation des arguments.
func TestXupleActionValidation(t *testing.T) {
	t.Log("🧪 TEST: Validation des arguments de l'action Xuple")
	t.Log("==================================================")

	storage := NewMemoryStorage()
	network := NewReteNetwork(storage)

	network.SetXupleHandler(func(xuplespace string, fact *Fact, triggeringFacts []*Fact) error {
		return nil
	})

	executor := NewActionExecutor(network, nil)
	network.ActionExecutor = executor

	handler := executor.GetRegistry().Get("Xuple")
	if handler == nil {
		t.Fatal("❌ Handler Xuple non trouvé")
	}

	// Test 1: Nombre d'arguments incorrect
	t.Run("Arguments insuffisants", func(t *testing.T) {
		err := handler.Validate([]interface{}{"space"})
		if err == nil {
			t.Error("❌ La validation devrait échouer avec 1 argument")
		}
		t.Logf("✅ Validation échoue correctement: %v", err)
	})

	// Test 2: Premier argument non-string
	t.Run("Premier argument invalide", func(t *testing.T) {
		fact := &Fact{ID: "f1", Type: "Test", Fields: map[string]interface{}{}}
		err := handler.Validate([]interface{}{123, fact})
		if err == nil {
			t.Error("❌ La validation devrait échouer avec premier argument non-string")
		}
		t.Logf("✅ Validation échoue correctement: %v", err)
	})

	// Test 3: Nom de xuplespace vide
	t.Run("Nom de xuplespace vide", func(t *testing.T) {
		fact := &Fact{ID: "f1", Type: "Test", Fields: map[string]interface{}{}}
		err := handler.Validate([]interface{}{"", fact})
		if err == nil {
			t.Error("❌ La validation devrait échouer avec nom vide")
		}
		t.Logf("✅ Validation échoue correctement: %v", err)
	})

	// Test 4: Second argument non-Fact
	t.Run("Second argument invalide", func(t *testing.T) {
		err := handler.Validate([]interface{}{"space", "not a fact"})
		if err == nil {
			t.Error("❌ La validation devrait échouer avec second argument non-Fact")
		}
		t.Logf("✅ Validation échoue correctement: %v", err)
	})

	// Test 5: Fait nil
	t.Run("Fait nil", func(t *testing.T) {
		var nilFact *Fact = nil
		err := handler.Validate([]interface{}{"space", nilFact})
		if err == nil {
			t.Error("❌ La validation devrait échouer avec fait nil")
		}
		t.Logf("✅ Validation échoue correctement: %v", err)
	})

	// Test 6: Arguments valides
	t.Run("Arguments valides", func(t *testing.T) {
		fact := &Fact{ID: "f1", Type: "Test", Fields: map[string]interface{}{"name": "test"}}
		err := handler.Validate([]interface{}{"space", fact})
		if err != nil {
			t.Errorf("❌ La validation devrait réussir avec arguments valides: %v", err)
		} else {
			t.Log("✅ Validation réussit avec arguments valides")
		}
	})
}
