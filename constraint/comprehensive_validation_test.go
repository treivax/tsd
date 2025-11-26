// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"os"
	"path/filepath"
	"testing"
)

// TestComprehensiveNonBlockingValidation teste tous les aspects de la validation non-bloquante
func TestComprehensiveNonBlockingValidation(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "comprehensive_validation_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Fichier 1: Définitions de types
	typesFile := filepath.Join(tempDir, "types.tsd")
	typesContent := `type Person : <id: string, name: string, age: number>
type Company : <id: string, name: string, employees: number>`

	err = os.WriteFile(typesFile, []byte(typesContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write types file: %v", err)
	}

	// Fichier 2: Mix de faits valides et invalides
	factsFile := filepath.Join(tempDir, "facts.tsd")
	factsContent := `Person(id: "P001", name: "Alice", age: 30)
UnknownType(id: "U001", field: "invalid")
Person(id: "P002", name: "Bob", age: 25)
Person(id: "P003", invalidField: "test", age: 35)
Company(id: "C001", name: "TechCorp", employees: 100)`

	err = os.WriteFile(factsFile, []byte(factsContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write facts file: %v", err)
	}

	// Test du parsing incrémental
	ps := NewProgramState()

	// Étape 1: Parser les types
	err = ps.ParseAndMerge(typesFile)
	if err != nil {
		t.Fatalf("Failed to parse types: %v", err)
	}

	if len(ps.Types) != 2 {
		t.Errorf("Expected 2 types, got %d", len(ps.Types))
	}

	// Étape 2: Parser les faits (contient des erreurs)
	err = ps.ParseAndMerge(factsFile)
	if err != nil {
		t.Fatalf("ParseAndMerge should not fail with validation errors: %v", err)
	}

	// Vérification: Les erreurs doivent être non-bloquantes
	if !ps.HasErrors() {
		t.Error("Expected validation errors to be recorded")
	}

	errorCount := ps.GetErrorCount()
	if errorCount < 2 {
		t.Errorf("Expected at least 2 errors (UnknownType + invalidField), got %d", errorCount)
	}

	// Vérification: Les faits valides doivent être parsés
	validFactCount := len(ps.Facts)
	if validFactCount < 3 {
		t.Errorf("Expected at least 3 valid facts, got %d", validFactCount)
	}

	// Vérification: Les types doivent être préservés
	if len(ps.Types) != 2 {
		t.Errorf("Expected 2 types after fact parsing, got %d", len(ps.Types))
	}

	// Vérification des erreurs spécifiques
	errors := ps.GetErrors()
	hasUndefinedTypeError := false
	hasInvalidFieldError := false

	for _, verr := range errors {
		if verr.Type == "fact" {
			t.Logf("Error recorded: %s", verr.Message)
			if contains(verr.Message, "undefined type") {
				hasUndefinedTypeError = true
			}
			if contains(verr.Message, "undefined field") || contains(verr.Message, "invalidField") {
				hasInvalidFieldError = true
			}
		}
	}

	if !hasUndefinedTypeError {
		t.Error("Expected error for undefined type")
	}

	if !hasInvalidFieldError {
		t.Error("Expected error for invalid field")
	}

	t.Logf("✅ Comprehensive validation: %d types, %d valid facts, %d errors",
		len(ps.Types), len(ps.Facts), errorCount)
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
