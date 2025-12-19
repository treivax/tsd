// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package api

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIngestFile_ReallySimple(t *testing.T) {
	t.Log("🧪 TEST INGEST FILE REALLY SIMPLE")

	tmpDir := t.TempDir()
	tsdFile := filepath.Join(tmpDir, "simple.tsd")

	// Using actual working TSD file content
	program := `// Simple test
type Product(code: string, price: number, category: string)

action expensive_product(arg1: string, arg2: number)

rule r1 : {prod: Product} / prod.price > 100 ==> expensive_product(prod.code, prod.price)

Product(code:PROD001, price:150, category:electronics)
`

	if err := os.WriteFile(tsdFile, []byte(program), 0644); err != nil {
		t.Fatal("❌ Impossible d'écrire le fichier:", err)
	}

	pipeline := NewPipeline()
	result, err := pipeline.IngestFile(tsdFile)
	if err != nil {
		t.Fatal("❌ Erreur d'ingestion:", err)
	}

	if result == nil {
		t.Fatal("❌ Result ne devrait pas être nil")
	}

	if result.TypeCount() != 1 {
		t.Errorf("❌ TypeCount attendu: 1, reçu: %d", result.TypeCount())
	}

	if result.RuleCount() != 1 {
		t.Errorf("❌ RuleCount attendu: 1, reçu: %d", result.RuleCount())
	}

	if result.FactCount() != 1 {
		t.Errorf("❌ FactCount attendu: 1, reçu: %d", result.FactCount())
	}

	t.Log("✅ Ingestion simple réussie")
	t.Logf("   Types: %d, Règles: %d, Faits: %d",
		result.TypeCount(), result.RuleCount(), result.FactCount())
}
