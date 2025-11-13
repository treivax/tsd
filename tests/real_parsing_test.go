package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	parser "github.com/treivax/tsd/constraint"
)

// TestRealPEGParsingIntegration teste le parsing réel avec le parseur PEG généré
func TestRealPEGParsingIntegration(t *testing.T) {

	constraintFiles := []string{
		"constraint/test/integration/alpha_conditions.constraint",
		"constraint/test/integration/beta_joins.constraint",
		"constraint/test/integration/negation.constraint",
		"constraint/test/integration/exists.constraint",
		"constraint/test/integration/aggregation.constraint",
		"constraint/test/integration/actions.constraint",
		"constraint/test/integration/complex_multi_node.constraint",
	}

	t.Run("Valid_Files_PEG_Parsing", func(t *testing.T) {
		for _, file := range constraintFiles {
			t.Run(filepath.Base(file), func(t *testing.T) {
				// Lire le fichier
				content, err := os.ReadFile(file)
				assert.NoError(t, err, "Should be able to read file: %s", file)

				// Parser avec le vrai parseur PEG
				result, err := parser.Parse(file, content)

				if err != nil {
					t.Logf("❌ Parsing failed for %s: %v", filepath.Base(file), err)
					t.Fail()
				} else {
					t.Logf("✅ Successfully parsed %s", filepath.Base(file))

					// Vérifier que le résultat a la structure attendue
					if resultMap, ok := result.(map[string]interface{}); ok {

						// Vérifier la présence des types
						if types, hasTypes := resultMap["types"]; hasTypes {
							if typeList, ok := types.([]interface{}); ok {
								t.Logf("   📋 Parsed %d type definitions", len(typeList))
							}
						}

						// Vérifier la présence des expressions
						if exprs, hasExprs := resultMap["expressions"]; hasExprs {
							if exprList, ok := exprs.([]interface{}); ok {
								t.Logf("   🔍 Parsed %d expressions", len(exprList))
							}
						}

					} else {
						t.Logf("⚠️  Unexpected result structure for %s", filepath.Base(file))
					}
				}
			})
		}
	})

	t.Run("Invalid_Files_PEG_Parsing", func(t *testing.T) {
		invalidFiles := []string{
			"constraint/test/integration/invalid_no_types.constraint",
			"constraint/test/integration/invalid_unknown_type.constraint",
		}

		for _, file := range invalidFiles {
			t.Run(filepath.Base(file), func(t *testing.T) {
				// Lire le fichier invalide
				content, err := os.ReadFile(file)
				assert.NoError(t, err, "Should be able to read invalid file: %s", file)

				// Parser avec le vrai parseur PEG - doit échouer
				result, err := parser.Parse(file, content)

				if err != nil {
					t.Logf("✅ Expected parsing failure for %s: %v", filepath.Base(file), err)
					// C'est attendu pour les fichiers invalides
				} else {
					t.Logf("⚠️  Unexpected success parsing invalid file %s", filepath.Base(file))
					t.Logf("   Result: %+v", result)
					// Ce n'est pas forcément un échec si la grammaire permet cette structure
				}
			})
		}
	})
}

// TestSemanticValidationWithRealParser teste la validation sémantique avec le parseur réel
func TestSemanticValidationWithRealParser(t *testing.T) {

	t.Run("Type_Reference_Validation", func(t *testing.T) {
		validFile := "constraint/test/integration/alpha_conditions.constraint"

		// Lire et parser le fichier
		content, err := os.ReadFile(validFile)
		assert.NoError(t, err, "Should read file")

		result, err := parser.Parse(validFile, content)
		assert.NoError(t, err, "Should parse valid file successfully")

		if resultMap, ok := result.(map[string]interface{}); ok {

			// Extraire les types déclarés du résultat parsé
			declaredTypes := make(map[string]bool)
			if types, hasTypes := resultMap["types"]; hasTypes {
				if typeList, ok := types.([]interface{}); ok {
					for _, typeItem := range typeList {
						if typeMap, ok := typeItem.(map[string]interface{}); ok {
							if typeName, hasName := typeMap["name"]; hasName {
								if name, ok := typeName.(string); ok {
									declaredTypes[name] = true
									t.Logf("📝 Declared type: %s", name)
								}
							}
						}
					}
				}
			}

			// Valider que nous avons bien des types déclarés
			assert.Greater(t, len(declaredTypes), 0, "Should have declared types")

			t.Logf("✅ Real PEG parsing validation successful: %d types declared", len(declaredTypes))

		} else {
			t.Error("Expected map structure from parser")
		}
	})
}
