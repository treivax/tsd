// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"testing"
)

func TestResetInstruction(t *testing.T) {
	t.Log("🧪 TEST INSTRUCTION RESET")
	t.Log("========================")

	// Test 1: Parse a reset instruction
	t.Run("ParseResetInstruction", func(t *testing.T) {
		input := []byte("reset")
		result, err := Parse("test_reset.constraint", input)
		if err != nil {
			t.Fatalf("❌ Erreur de parsing: %v", err)
		}

		// Convert to program
		program, err := ConvertResultToProgram(result)
		if err != nil {
			t.Fatalf("❌ Erreur de conversion: %v", err)
		}

		// Verify we have a reset instruction
		if len(program.Resets) != 1 {
			t.Fatalf("❌ Attendu 1 instruction reset, reçu %d", len(program.Resets))
		}

		if program.Resets[0].Type != "reset" {
			t.Errorf("❌ Type incorrect: attendu 'reset', reçu '%s'", program.Resets[0].Type)
		}

		t.Log("✅ Instruction reset parsée avec succès")
	})

	// Test 2: Reset in a program with types
	t.Run("ResetInCompleteProgram", func(t *testing.T) {
		input := []byte(`
type User : <name: string, age: number>

reset
`)
		result, err := Parse("test_reset_complete.constraint", input)
		if err != nil {
			t.Fatalf("❌ Erreur de parsing: %v", err)
		}

		program, err := ConvertResultToProgram(result)
		if err != nil {
			t.Fatalf("❌ Erreur de conversion: %v", err)
		}

		// Verify all components are parsed
		if len(program.Types) != 1 {
			t.Errorf("❌ Attendu 1 type, reçu %d", len(program.Types))
		}
		if len(program.Resets) != 1 {
			t.Errorf("❌ Attendu 1 instruction reset, reçu %d", len(program.Resets))
		}

		t.Log("✅ Programme complet avec reset parsé avec succès")
	})

	// Test 3: Multiple reset instructions
	t.Run("MultipleResets", func(t *testing.T) {
		input := []byte(`
type User : <name: string>

reset

type Order : <id: number>

reset
`)
		result, err := Parse("test_multiple_resets.constraint", input)
		if err != nil {
			t.Fatalf("❌ Erreur de parsing: %v", err)
		}

		program, err := ConvertResultToProgram(result)
		if err != nil {
			t.Fatalf("❌ Erreur de conversion: %v", err)
		}

		if len(program.Resets) != 2 {
			t.Fatalf("❌ Attendu 2 instructions reset, reçu %d", len(program.Resets))
		}

		t.Log("✅ Multiples instructions reset parsées avec succès")
	})
}

func TestProgramStateReset(t *testing.T) {
	t.Log("🧪 TEST RESET DU PROGRAM STATE")
	t.Log("==============================")

	// Test 1: Reset clears all state
	t.Run("ResetClearsAllState", func(t *testing.T) {
		ps := NewProgramState()

		// Add some data
		ps.Types["User"] = &TypeDefinition{
			Type: "typeDefinition",
			Name: "User",
			Fields: []Field{
				{Name: "name", Type: "string"},
			},
		}

		ps.Rules = append(ps.Rules, &Expression{
			Type: "expression",
		})

		ps.Facts = append(ps.Facts, &Fact{
			Type:     "fact",
			TypeName: "User",
		})

		ps.FilesParsed = append(ps.FilesParsed, "test.constraint")

		ps.Errors = append(ps.Errors, ValidationError{
			Type:    "test",
			Message: "test error",
		})

		// Verify state has data
		if len(ps.Types) == 0 {
			t.Fatal("❌ Types devrait contenir des données avant reset")
		}
		if len(ps.Rules) == 0 {
			t.Fatal("❌ Rules devrait contenir des données avant reset")
		}
		if len(ps.Facts) == 0 {
			t.Fatal("❌ Facts devrait contenir des données avant reset")
		}
		if len(ps.FilesParsed) == 0 {
			t.Fatal("❌ FilesParsed devrait contenir des données avant reset")
		}
		if len(ps.Errors) == 0 {
			t.Fatal("❌ Errors devrait contenir des données avant reset")
		}

		// Reset
		ps.Reset()

		// Verify everything is cleared
		if len(ps.Types) != 0 {
			t.Errorf("❌ Types devrait être vide après reset, reçu %d éléments", len(ps.Types))
		}
		if len(ps.Rules) != 0 {
			t.Errorf("❌ Rules devrait être vide après reset, reçu %d éléments", len(ps.Rules))
		}
		if len(ps.Facts) != 0 {
			t.Errorf("❌ Facts devrait être vide après reset, reçu %d éléments", len(ps.Facts))
		}
		if len(ps.FilesParsed) != 0 {
			t.Errorf("❌ FilesParsed devrait être vide après reset, reçu %d éléments", len(ps.FilesParsed))
		}
		if len(ps.Errors) != 0 {
			t.Errorf("❌ Errors devrait être vide après reset, reçu %d éléments", len(ps.Errors))
		}

		t.Log("✅ Reset a correctement vidé tout le state")
	})

	// Test 2: Reset can be called multiple times
	t.Run("ResetMultipleTimes", func(t *testing.T) {
		ps := NewProgramState()

		// Add data, reset, add data again, reset again
		ps.Types["User"] = &TypeDefinition{Name: "User"}
		ps.Reset()

		if len(ps.Types) != 0 {
			t.Fatal("❌ Premier reset a échoué")
		}

		ps.Types["Order"] = &TypeDefinition{Name: "Order"}
		ps.Reset()

		if len(ps.Types) != 0 {
			t.Fatal("❌ Deuxième reset a échoué")
		}

		t.Log("✅ Reset peut être appelé plusieurs fois")
	})

	// Test 3: Reset initializes empty structures (not nil)
	t.Run("ResetInitializesEmptyStructures", func(t *testing.T) {
		ps := NewProgramState()
		ps.Reset()

		// Verify we can add data after reset without nil errors
		ps.Types["Test"] = &TypeDefinition{Name: "Test"}
		ps.Rules = append(ps.Rules, &Expression{Type: "expression"})
		ps.Facts = append(ps.Facts, &Fact{Type: "fact"})
		ps.FilesParsed = append(ps.FilesParsed, "test.constraint")
		ps.Errors = append(ps.Errors, ValidationError{Type: "test"})

		if ps.Types == nil {
			t.Error("❌ Types ne devrait pas être nil après reset")
		}
		if ps.Rules == nil {
			t.Error("❌ Rules ne devrait pas être nil après reset")
		}
		if ps.Facts == nil {
			t.Error("❌ Facts ne devrait pas être nil après reset")
		}
		if ps.FilesParsed == nil {
			t.Error("❌ FilesParsed ne devrait pas être nil après reset")
		}
		if ps.Errors == nil {
			t.Error("❌ Errors ne devrait pas être nil après reset")
		}

		t.Log("✅ Reset initialise correctement des structures vides (non nil)")
	})
}

func TestIterativeParserReset(t *testing.T) {
	t.Log("🧪 TEST RESET DU ITERATIVE PARSER")
	t.Log("=================================")

	t.Run("ResetIterativeParser", func(t *testing.T) {
		parser := NewIterativeParser()

		// Parse some content
		content := `
type User : <name: string, age: number>
`
		err := parser.ParseContent(content, "test.constraint")
		if err != nil {
			t.Fatalf("❌ Erreur de parsing: %v", err)
		}

		// Verify we have data
		stats := parser.GetParsingStatistics()
		if stats.TypesCount == 0 {
			t.Fatal("❌ Parser devrait avoir des types avant reset")
		}

		// Reset
		parser.Reset()

		// Verify everything is cleared
		statsAfter := parser.GetParsingStatistics()
		if statsAfter.TypesCount != 0 {
			t.Errorf("❌ TypesCount devrait être 0 après reset, reçu %d", statsAfter.TypesCount)
		}

		if statsAfter.FactsCount != 0 {
			t.Errorf("❌ FactsCount devrait être 0 après reset, reçu %d", statsAfter.FactsCount)
		}
		if statsAfter.FilesParsedCount != 0 {
			t.Errorf("❌ FilesParsedCount devrait être 0 après reset, reçu %d", statsAfter.FilesParsedCount)
		}

		t.Log("✅ IterativeParser reset fonctionne correctement")
	})

	t.Run("CanParseAfterReset", func(t *testing.T) {
		parser := NewIterativeParser()

		// Parse, reset, parse again
		content1 := `type User : <name: string>`
		err := parser.ParseContent(content1, "test1.constraint")
		if err != nil {
			t.Fatalf("❌ Erreur de parsing initial: %v", err)
		}

		parser.Reset()

		content2 := `type Order : <id: number>`
		err = parser.ParseContent(content2, "test2.constraint")
		if err != nil {
			t.Fatalf("❌ Erreur de parsing après reset: %v", err)
		}

		// Verify we have the new data, not the old
		program := parser.GetProgram()
		if len(program.Types) != 1 {
			t.Fatalf("❌ Attendu 1 type, reçu %d", len(program.Types))
		}

		if program.Types[0].Name != "Order" {
			t.Errorf("❌ Attendu type 'Order', reçu '%s'", program.Types[0].Name)
		}

		t.Log("✅ Peut parser correctement après reset")
	})
}
