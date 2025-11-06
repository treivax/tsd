package main

import (
	"testing"
	
	parser "github.com/treivax/tsd/constraint/grammar"
)

// TestSimplePEGParsing teste le parsing avec le parseur corrigé
func TestSimplePEGParsing(t *testing.T) {
	
	t.Run("Parse_With_Comments", func(t *testing.T) {
		// Créer un fichier de test avec commentaires
		testContent := `// Commentaire de début
type User : <id: string, name: string>
// Commentaire intermédiaire  

{u: User} / u.id == "123" // Commentaire de fin`

		// Parser avec le nouveau parseur
		result, err := parser.Parse("test", []byte(testContent))
		
		if err != nil {
			t.Logf("❌ Parsing failed: %v", err)
			t.Fail()
		} else {
			t.Logf("✅ Parsing successful with comments!")
			
			if resultMap, ok := result.(map[string]interface{}); ok {
				t.Logf("📊 Result structure: %+v", resultMap)
				
				// Vérifier les types
				if types, hasTypes := resultMap["types"]; hasTypes {
					if typeList, ok := types.([]interface{}); ok {
						t.Logf("   📋 Types parsed: %d", len(typeList))
						for i, typeItem := range typeList {
							if typeMap, ok := typeItem.(map[string]interface{}); ok {
								t.Logf("     Type %d: %s", i+1, typeMap["name"])
							}
						}
					}
				}
				
				// Vérifier les expressions
				if exprs, hasExprs := resultMap["expressions"]; hasExprs {
					if exprList, ok := exprs.([]interface{}); ok {
						t.Logf("   🔍 Expressions parsed: %d", len(exprList))
					}
				}
			}
		}
	})
	
	t.Run("Parse_Multiline_Comments", func(t *testing.T) {
		testContent := `/* simple comment */
type Account : <id: string, balance: number>

{a: Account} / a.balance >= 0`

		_, err := parser.Parse("test_multiline", []byte(testContent))
		
		if err != nil {
			t.Logf("❌ Multiline comment parsing failed: %v", err) 
			t.Fail()
		} else {
			t.Logf("✅ Multiline comment parsing successful!")
		}
	})
}