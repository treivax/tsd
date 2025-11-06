package main

import (
	"os"
	"strings"
	"testing"
)

// TestSemanticValidation teste la validation sémantique des types
func TestSemanticValidation(t *testing.T) {

	t.Run("Valid_Type_References", func(t *testing.T) {
		// Simuler un fichier avec des références de types valides
		validFile := "constraint/test/integration/alpha_conditions.constraint"

		content, err := os.ReadFile(validFile)
		if err != nil {
			t.Fatalf("Cannot read test file: %v", err)
		}

		// Parser le contenu pour extraire les types déclarés et référencés
		declaredTypes := extractDeclaredTypes(string(content))
		referencedTypes := extractReferencedTypes(string(content))

		t.Logf("Declared types: %v", declaredTypes)
		t.Logf("Referenced types: %v", referencedTypes)

		// Vérifier que tous les types référencés sont déclarés
		for _, refType := range referencedTypes {
			found := false
			for _, declType := range declaredTypes {
				if refType == declType {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Type referenced but not declared: %s", refType)
			}
		}

		t.Log("✅ All referenced types are properly declared")
	})

	t.Run("Invalid_Type_References", func(t *testing.T) {
		// Test avec le fichier contenant une référence invalide
		invalidFile := "constraint/test/integration/invalid_unknown_type.constraint"

		content, err := os.ReadFile(invalidFile)
		if err != nil {
			t.Fatalf("Cannot read invalid test file: %v", err)
		}

		// Parser le contenu
		declaredTypes := extractDeclaredTypes(string(content))
		referencedTypes := extractReferencedTypes(string(content))

		t.Logf("Declared types: %v", declaredTypes)
		t.Logf("Referenced types: %v", referencedTypes)

		// Ce fichier DOIT avoir des références invalides
		invalidRefs := 0
		for _, refType := range referencedTypes {
			found := false
			for _, declType := range declaredTypes {
				if refType == declType {
					found = true
					break
				}
			}
			if !found {
				invalidRefs++
				t.Logf("❌ Invalid type reference detected: %s", refType)
			}
		}

		if invalidRefs == 0 {
			t.Error("Expected invalid type references but found none")
		} else {
			t.Logf("✅ Semantic validation working: %d invalid type references detected", invalidRefs)
		}
	})
}

// extractDeclaredTypes extrait les types déclarés d'un contenu de fichier
func extractDeclaredTypes(content string) []string {
	var types []string
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "type ") && strings.Contains(line, ":") {
			// Extraire le nom du type
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				typeName := parts[1]
				types = append(types, typeName)
			}
		}
	}

	return types
}

// extractReferencedTypes extrait les types référencés dans les expressions
func extractReferencedTypes(content string) []string {
	var types []string
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "{") && strings.Contains(line, ":") && strings.Contains(line, "}") {
			// Extraire les types des variables : {t: Transaction, a: Account}
			start := strings.Index(line, "{") + 1
			end := strings.Index(line, "}")
			if end > start {
				varSection := line[start:end]
				variables := strings.Split(varSection, ",")

				for _, variable := range variables {
					if strings.Contains(variable, ":") {
						parts := strings.Split(variable, ":")
						if len(parts) >= 2 {
							typeName := strings.TrimSpace(parts[1])
							types = append(types, typeName)
						}
					}
				}
			}
		}
	}

	return types
}
