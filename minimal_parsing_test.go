package main

import (
	"testing"
	"os"
	
	"github.com/stretchr/testify/assert"
	parser "github.com/treivax/tsd/constraint"
)

// TestMinimalParsing teste le parsing avec un fichier très simple
func TestMinimalParsing(t *testing.T) {
	
	t.Run("Parse_Minimal_File", func(t *testing.T) {
		file := "constraint/test/integration/minimal_test.constraint"
		
		// Lire le fichier
		content, err := os.ReadFile(file)
		assert.NoError(t, err, "Should be able to read minimal file")
		
		t.Logf("📄 File content:\n%s", string(content))
		
		// Tenter le parsing avec le vrai parseur PEG
		result, err := parser.Parse(file, content)
		
		if err != nil {
			t.Logf("❌ Parsing failed: %v", err)
		} else {
			t.Logf("✅ Parsing successful!")
			
			if resultMap, ok := result.(map[string]interface{}); ok {
				t.Logf("📊 Result structure: %+v", resultMap)
				
				// Analyser la structure
				if types, hasTypes := resultMap["types"]; hasTypes {
					t.Logf("   📋 Types found: %+v", types)
				}
				if exprs, hasExprs := resultMap["expressions"]; hasExprs {
					t.Logf("   🔍 Expressions found: %+v", exprs)
				}
			}
		}
	})
}