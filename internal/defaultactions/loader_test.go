// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package defaultactions

import (
	"testing"

	"github.com/treivax/tsd/internal/actiondefs"
)

func TestLoadDefaultActions(t *testing.T) {
	t.Log("🧪 TEST CHARGEMENT ACTIONS PAR DÉFAUT")

	actions, err := LoadDefaultActions()
	if err != nil {
		t.Fatalf("❌ Erreur chargement: %v", err)
	}

	// Vérifier le nombre d'actions
	expectedCount := len(actiondefs.DefaultActionNames)
	if len(actions) != expectedCount {
		t.Errorf("❌ Attendu %d actions, reçu %d", expectedCount, len(actions))
	}

	// Vérifier que chaque action est marquée comme par défaut
	for _, action := range actions {
		if !action.IsDefault {
			t.Errorf("❌ Action '%s' devrait être marquée IsDefault", action.Name)
		}
	}

	// Vérifier que toutes les actions attendues sont présentes
	actionMap := make(map[string]bool)
	for _, action := range actions {
		actionMap[action.Name] = true
	}

	for _, name := range actiondefs.DefaultActionNames {
		if !actionMap[name] {
			t.Errorf("❌ Action par défaut manquante: %s", name)
		}
	}

	t.Log("✅ Toutes les actions par défaut chargées correctement")
}

func TestLoadDefaultActions_Signatures(t *testing.T) {
	t.Log("🧪 TEST SIGNATURES DES ACTIONS PAR DÉFAUT")

	actions, err := LoadDefaultActions()
	if err != nil {
		t.Fatalf("❌ Erreur chargement: %v", err)
	}

	// Définir les signatures attendues
	expectedSignatures := map[string]struct {
		paramCount int
		params     map[string]string // nom -> type
	}{
		"Print":   {1, map[string]string{"message": "string"}},
		"Log":     {1, map[string]string{"message": "string"}},
		"Update":  {2, map[string]string{"variable": "any", "modifications": "any"}},
		"Insert":  {1, map[string]string{"fact": "any"}},
		"Retract": {1, map[string]string{"fact": "any"}},
		"Xuple":   {2, map[string]string{"xuplespace": "string", "fact": "any"}},
	}

	for _, action := range actions {
		expected, exists := expectedSignatures[action.Name]
		if !exists {
			t.Errorf("❌ Action inattendue: %s", action.Name)
			continue
		}

		// Vérifier le nombre de paramètres
		if len(action.Parameters) != expected.paramCount {
			t.Errorf("❌ Action %s: attendu %d paramètres, reçu %d",
				action.Name, expected.paramCount, len(action.Parameters))
		}

		// Vérifier les types de paramètres
		for _, param := range action.Parameters {
			expectedType, exists := expected.params[param.Name]
			if !exists {
				t.Errorf("❌ Action %s: paramètre inattendu '%s'",
					action.Name, param.Name)
				continue
			}

			if param.Type != expectedType {
				t.Errorf("❌ Action %s, paramètre %s: attendu type '%s', reçu '%s'",
					action.Name, param.Name, expectedType, param.Type)
			}
		}
	}

	t.Log("✅ Toutes les signatures sont correctes")
}

func TestIsDefaultAction(t *testing.T) {
	t.Log("🧪 TEST IsDefaultAction")

	tests := []struct {
		name     string
		expected bool
	}{
		{"Print", true},
		{"Log", true},
		{"Update", true},
		{"Insert", true},
		{"Retract", true},
		{"Xuple", true},
		{"CustomAction", false},
		{"Unknown", false},
		{"print", false}, // case-sensitive
		{"PRINT", false}, // case-sensitive
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := actiondefs.IsDefaultAction(tt.name)
			if result != tt.expected {
				t.Errorf("❌ IsDefaultAction(%q) = %v, attendu %v", tt.name, result, tt.expected)
			}
		})
	}

	t.Log("✅ IsDefaultAction fonctionne correctement")
}

func TestDefaultActionNames_Complete(t *testing.T) {
	t.Log("🧪 TEST COMPLÉTUDE DefaultActionNames")

	// Vérifier que toutes les actions attendues sont présentes
	expectedNames := make(map[string]bool)
	for _, name := range actiondefs.DefaultActionNames {
		expectedNames[name] = true
	}

	// Vérifier que toutes les actions attendues sont dans la liste
	for _, name := range actiondefs.DefaultActionNames {
		if !expectedNames[name] {
			t.Errorf("❌ Action inattendue dans DefaultActionNames: %s", name)
		}
		delete(expectedNames, name)
	}

	// Vérifier qu'il ne manque aucune action
	for name := range expectedNames {
		t.Errorf("❌ Action manquante dans DefaultActionNames: %s", name)
	}

	// Vérifier qu'il n'y a pas de doublons
	seen := make(map[string]bool)
	for _, name := range actiondefs.DefaultActionNames {
		if seen[name] {
			t.Errorf("❌ Action en double dans DefaultActionNames: %s", name)
		}
		seen[name] = true
	}

	t.Log("✅ DefaultActionNames est complet et sans doublon")
}
