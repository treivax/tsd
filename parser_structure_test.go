package main

import (
	"encoding/json"
	"testing"

	parser "github.com/treivax/tsd/constraint"
)

// TestParserStructureAnalysis analyse la structure exacte du parseur existant
func TestParserStructureAnalysis(t *testing.T) {

	content := `type User : <id: string, name: string>

{u: User} / u.name == "test"`

	result, err := parser.Parse("test", []byte(content))

	if err != nil {
		t.Fatalf("Parsing failed: %v", err)
	}

	// Afficher la structure complète en JSON pour analyse
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	t.Logf("📊 Complete parsed structure:\n%s", string(jsonData))

	// Analyser la structure
	if resultMap, ok := result.(map[string]interface{}); ok {
		t.Logf("\n🔍 Structure analysis:")

		for key, value := range resultMap {
			t.Logf("  Top-level key: %s", key)
			analyzeValue(t, value, "    ")
		}
	}
}

func analyzeValue(t *testing.T, value interface{}, indent string) {
	switch v := value.(type) {
	case map[string]interface{}:
		t.Logf("%smap[string]interface{} with keys:", indent)
		for key, subValue := range v {
			t.Logf("%s  %s:", indent, key)
			analyzeValue(t, subValue, indent+"    ")
		}
	case []interface{}:
		t.Logf("%s[]interface{} with %d elements:", indent, len(v))
		for i, item := range v {
			t.Logf("%s  [%d]:", indent, i)
			analyzeValue(t, item, indent+"    ")
		}
	case string:
		t.Logf("%sstring: %q", indent, v)
	case int, int64, float64:
		t.Logf("%snumber: %v", indent, v)
	case bool:
		t.Logf("%sbool: %v", indent, v)
	default:
		t.Logf("%sunknown type: %T = %v", indent, v, v)
	}
}
