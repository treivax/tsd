// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"strings"
	"testing"
)

func TestIntegration_ParseAndGenerateIDs(t *testing.T) {
	t.Log("🧪 TEST INTEGRATION - PARSING → VALIDATION → ID GENERATION")
	t.Log("===========================================================")

	tests := []struct {
		name       string
		input      string
		wantErr    bool
		checkIDsFn func(*testing.T, *Program)
	}{
		{
			name: "Programme complet avec PK simple",
			input: `
type Person(#nom: string, age: number)

Person(nom: "Alice", age: 30)
Person(nom: "Bob", age: 25)
`,
			wantErr: false,
			checkIDsFn: func(t *testing.T, prog *Program) {
				reteFacts, err := ConvertFactsToReteFormat(*prog)
				if err != nil {
					t.Fatalf("❌ ConvertFactsToReteFormat() error = %v", err)
				}

				if len(reteFacts) != 2 {
					t.Fatalf("❌ Expected 2 facts, got %d", len(reteFacts))
				}

				// Vérifier les IDs générés
				expectedIDs := map[string]bool{
					"Person~Alice": false,
					"Person~Bob":   false,
				}

				for _, fact := range reteFacts {
					factID, ok := fact["id"].(string)
					if !ok {
						t.Errorf("❌ Fact missing id field or wrong type")
						continue
					}

					if _, exists := expectedIDs[factID]; exists {
						expectedIDs[factID] = true
						t.Logf("✅ Found expected fact ID: %s", factID)
					} else {
						t.Errorf("❌ Unexpected fact ID: %s", factID)
					}
				}

				for id, found := range expectedIDs {
					if !found {
						t.Errorf("❌ Expected fact with ID %s not found", id)
					}
				}
			},
		},
		{
			name: "Programme complet avec PK composite",
			input: `
type Person(#prenom: string, #nom: string, age: number)

Person(prenom: "Alice", nom: "Dupont", age: 30)
Person(prenom: "Bob", nom: "Martin", age: 25)
`,
			wantErr: false,
			checkIDsFn: func(t *testing.T, prog *Program) {
				reteFacts, err := ConvertFactsToReteFormat(*prog)
				if err != nil {
					t.Fatalf("❌ ConvertFactsToReteFormat() error = %v", err)
				}

				expectedIDs := []string{
					"Person~Alice_Dupont",
					"Person~Bob_Martin",
				}

				if len(reteFacts) != len(expectedIDs) {
					t.Fatalf("❌ Expected %d facts, got %d", len(expectedIDs), len(reteFacts))
				}

				for i, fact := range reteFacts {
					factID, ok := fact["id"].(string)
					if !ok {
						t.Errorf("❌ Fact %d missing id field or wrong type", i)
						continue
					}

					if factID != expectedIDs[i] {
						t.Errorf("❌ Fact %d: expected ID %s, got %s", i, expectedIDs[i], factID)
					} else {
						t.Logf("✅ Fact %d: ID %s", i, factID)
					}
				}
			},
		},
		{
			name: "Programme avec type sans PK (hash)",
			input: `
type Event(timestamp: number, message: string)

Event(timestamp: 1234567890, message: "test1")
Event(timestamp: 1234567891, message: "test2")
`,
			wantErr: false,
			checkIDsFn: func(t *testing.T, prog *Program) {
				reteFacts, err := ConvertFactsToReteFormat(*prog)
				if err != nil {
					t.Fatalf("❌ ConvertFactsToReteFormat() error = %v", err)
				}

				if len(reteFacts) != 2 {
					t.Fatalf("❌ Expected 2 facts, got %d", len(reteFacts))
				}

				// Vérifier le format hash
				seenIDs := make(map[string]bool)
				for i, fact := range reteFacts {
					factID, ok := fact["id"].(string)
					if !ok {
						t.Errorf("❌ Fact %d missing id field or wrong type", i)
						continue
					}

					if !strings.HasPrefix(factID, "Event~") {
						t.Errorf("❌ Fact %d: expected ID to start with Event~, got %s", i, factID)
					}

					hashPart := strings.TrimPrefix(factID, "Event~")
					if len(hashPart) != 16 {
						t.Errorf("❌ Fact %d: expected hash of 16 chars, got %d: %s", i, len(hashPart), hashPart)
					} else {
						t.Logf("✅ Fact %d: hash ID %s", i, factID)
					}

					seenIDs[factID] = true
				}

				// Vérifier que les deux faits ont des IDs différents
				if len(seenIDs) != 2 {
					t.Errorf("❌ Expected 2 different IDs, got %d", len(seenIDs))
				}
			},
		},
		{
			name: "Rejet de id explicite dans assertion",
			input: `
type Person(#nom: string)

Person(id: "custom_id", nom: "Alice")
`,
			wantErr: true,
		},
		{
			name: "Caractères spéciaux dans PK",
			input: `
type Resource(#path: string)

Resource(path: "/home/user~test_file")
`,
			wantErr: false,
			checkIDsFn: func(t *testing.T, prog *Program) {
				reteFacts, err := ConvertFactsToReteFormat(*prog)
				if err != nil {
					t.Fatalf("❌ ConvertFactsToReteFormat() error = %v", err)
				}

				if len(reteFacts) != 1 {
					t.Fatalf("❌ Expected 1 fact, got %d", len(reteFacts))
				}

				factID, ok := reteFacts[0]["id"].(string)
				if !ok {
					t.Fatalf("❌ Fact missing id field or wrong type")
				}

				// Vérifier que les caractères spéciaux sont échappés
				if !strings.Contains(factID, "%") {
					t.Errorf("❌ Expected escaped characters in ID, got: %s", factID)
				}

				// Vérifier que l'ID commence par le type
				if !strings.HasPrefix(factID, "Resource~") {
					t.Errorf("❌ Expected ID to start with Resource~, got: %s", factID)
				}

				t.Logf("✅ Special chars escaped: %s", factID)
			},
		},
		{
			name: "Plusieurs types avec stratégies différentes",
			input: `
type Person(#nom: string, age: number)
type Event(timestamp: number, message: string)

Person(nom: "Alice", age: 30)
Event(timestamp: 1234567890, message: "User logged in")
`,
			wantErr: false,
			checkIDsFn: func(t *testing.T, prog *Program) {
				reteFacts, err := ConvertFactsToReteFormat(*prog)
				if err != nil {
					t.Fatalf("❌ ConvertFactsToReteFormat() error = %v", err)
				}

				if len(reteFacts) != 2 {
					t.Fatalf("❌ Expected 2 facts, got %d", len(reteFacts))
				}

				for _, fact := range reteFacts {
					factID, ok := fact["id"].(string)
					if !ok {
						t.Errorf("❌ Fact missing id field")
						continue
					}

					factType, ok := fact["reteType"].(string)
					if !ok {
						t.Errorf("❌ Fact missing reteType field")
						continue
					}

					if factType == "Person" {
						// Person doit avoir un ID basé sur PK
						if factID != "Person~Alice" {
							t.Errorf("❌ Person fact: expected ID Person~Alice, got %s", factID)
						} else {
							t.Logf("✅ Person with PK-based ID: %s", factID)
						}
					} else if factType == "Event" {
						// Event doit avoir un ID basé sur hash
						if !strings.HasPrefix(factID, "Event~") {
							t.Errorf("❌ Event fact: expected ID to start with Event~, got %s", factID)
						} else {
							hashPart := strings.TrimPrefix(factID, "Event~")
							if len(hashPart) == 16 {
								t.Logf("✅ Event with hash-based ID: %s", factID)
							} else {
								t.Errorf("❌ Event hash has wrong length: %d", len(hashPart))
							}
						}
					}
				}
			},
		},
		{
			name: "PK composite avec 3 champs",
			input: `
type Location(#country: string, #city: string, #street: string, population: number)

Location(country: "France", city: "Paris", street: "Rue de Rivoli", population: 1000000)
`,
			wantErr: false,
			checkIDsFn: func(t *testing.T, prog *Program) {
				reteFacts, err := ConvertFactsToReteFormat(*prog)
				if err != nil {
					t.Fatalf("❌ ConvertFactsToReteFormat() error = %v", err)
				}

				if len(reteFacts) != 1 {
					t.Fatalf("❌ Expected 1 fact, got %d", len(reteFacts))
				}

				factID, ok := reteFacts[0]["id"].(string)
				if !ok {
					t.Fatalf("❌ Fact missing id field")
				}

				expectedID := "Location~France_Paris_Rue%20de%20Rivoli"
				if factID != expectedID {
					t.Errorf("❌ Expected ID %s, got %s", expectedID, factID)
				} else {
					t.Logf("✅ Composite PK with 3 fields: %s", factID)
				}
			},
		},
		{
			name: "Type avec PK numérique",
			input: `
type Product(#code: number, name: string)

Product(code: 12345, name: "Widget")
`,
			wantErr: false,
			checkIDsFn: func(t *testing.T, prog *Program) {
				reteFacts, err := ConvertFactsToReteFormat(*prog)
				if err != nil {
					t.Fatalf("❌ ConvertFactsToReteFormat() error = %v", err)
				}

				if len(reteFacts) != 1 {
					t.Fatalf("❌ Expected 1 fact, got %d", len(reteFacts))
				}

				factID, ok := reteFacts[0]["id"].(string)
				if !ok {
					t.Fatalf("❌ Fact missing id field")
				}

				expectedID := "Product~12345"
				if factID != expectedID {
					t.Errorf("❌ Expected ID %s, got %s", expectedID, factID)
				} else {
					t.Logf("✅ Numeric PK: %s", factID)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("📋 Test: %s", tt.name)

			// Parser le programme
			result, err := Parse("test.tsd", []byte(tt.input))
			if err != nil {
				if !tt.wantErr {
					t.Fatalf("❌ Parse() error = %v, wantErr %v", err, tt.wantErr)
				}
				t.Logf("✅ Parse error as expected: %v", err)
				return
			}

			// Valider le programme
			err = ValidateConstraintProgram(result)
			if (err != nil) != tt.wantErr {
				t.Fatalf("❌ ValidateConstraintProgram() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err != nil {
				t.Logf("✅ Validation error as expected: %v", err)
				return
			}

			// Convertir en Program
			prog, err := ConvertResultToProgram(result)
			if err != nil {
				t.Fatalf("❌ ConvertResultToProgram() error = %v", err)
			}

			// Exécuter les vérifications d'ID si fourni
			if tt.checkIDsFn != nil {
				tt.checkIDsFn(t, prog)
			}

			t.Logf("✅ Test passed")
		})
	}
}

func TestIntegration_IDDeterminism(t *testing.T) {
	t.Log("🧪 TEST INTEGRATION - DÉTERMINISME DES IDS")
	t.Log("==========================================")

	input := `
type Person(#nom: string, age: number)
type Event(timestamp: number, message: string)

Person(nom: "Alice", age: 30)
Event(timestamp: 1234567890, message: "test")
`

	const numRuns = 5

	// Parser et convertir plusieurs fois
	var allIDs [][]string
	for i := 0; i < numRuns; i++ {
		result, err := Parse("test.tsd", []byte(input))
		if err != nil {
			t.Fatalf("❌ Parse() iteration %d error = %v", i, err)
		}

		err = ValidateConstraintProgram(result)
		if err != nil {
			t.Fatalf("❌ ValidateConstraintProgram() iteration %d error = %v", i, err)
		}

		prog, err := ConvertResultToProgram(result)
		if err != nil {
			t.Fatalf("❌ ConvertResultToProgram() iteration %d error = %v", i, err)
		}

		reteFacts, err := ConvertFactsToReteFormat(*prog)
		if err != nil {
			t.Fatalf("❌ ConvertFactsToReteFormat() iteration %d error = %v", i, err)
		}

		ids := make([]string, len(reteFacts))
		for j, fact := range reteFacts {
			factID, ok := fact["id"].(string)
			if !ok {
				t.Fatalf("❌ Fact %d missing id field in iteration %d", j, i)
			}
			ids[j] = factID
		}
		allIDs = append(allIDs, ids)

		if i == 0 {
			t.Logf("📋 Run %d IDs: %v", i+1, ids)
		}
	}

	// Vérifier que tous les runs ont produit les mêmes IDs
	firstRun := allIDs[0]
	for i, ids := range allIDs[1:] {
		for j, id := range ids {
			if id != firstRun[j] {
				t.Errorf("❌ Run %d, fact %d: ID mismatch: got %s, want %s", i+2, j, id, firstRun[j])
			}
		}
	}

	t.Logf("✅ All %d runs produced identical IDs", numRuns)
	t.Logf("✅ ID generation is deterministic")
}

func TestIntegration_BackwardCompatibility(t *testing.T) {
	t.Log("🧪 TEST INTEGRATION - RÉTROCOMPATIBILITÉ")
	t.Log("=========================================")

	// Programme sans clés primaires (ancien format)
	input := `
type Person(nom: string, age: number)

Person(nom: "Alice", age: 30)
Person(nom: "Bob", age: 25)
`

	result, err := Parse("test.tsd", []byte(input))
	if err != nil {
		t.Fatalf("❌ Parse() error = %v", err)
	}

	err = ValidateConstraintProgram(result)
	if err != nil {
		t.Fatalf("❌ ValidateConstraintProgram() error = %v", err)
	}

	prog, err := ConvertResultToProgram(result)
	if err != nil {
		t.Fatalf("❌ ConvertResultToProgram() error = %v", err)
	}

	// Vérifier que le type n'a pas de clé primaire
	if len(prog.Types) != 1 {
		t.Fatalf("❌ Expected 1 type, got %d", len(prog.Types))
	}

	if prog.Types[0].HasPrimaryKey() {
		t.Error("❌ Type should not have primary key")
	} else {
		t.Log("✅ Type without PK confirmed")
	}

	// Convertir les faits
	reteFacts, err := ConvertFactsToReteFormat(*prog)
	if err != nil {
		t.Fatalf("❌ ConvertFactsToReteFormat() error = %v", err)
	}

	// Vérifier que les IDs sont générés avec hash (pas de PK)
	for i, fact := range reteFacts {
		factID, ok := fact["id"].(string)
		if !ok {
			t.Errorf("❌ Fact %d missing id field", i)
			continue
		}

		if !strings.HasPrefix(factID, "Person~") {
			t.Errorf("❌ Fact %d: Expected ID to start with Person~, got %s", i, factID)
		}

		hashPart := strings.TrimPrefix(factID, "Person~")
		if len(hashPart) != 16 {
			t.Errorf("❌ Fact %d: Expected hash of 16 chars, got %d: %s", i, len(hashPart), hashPart)
		} else {
			t.Logf("✅ Fact %d: hash-based ID %s", i, factID)
		}
	}

	t.Log("✅ Backward compatibility preserved")
}
