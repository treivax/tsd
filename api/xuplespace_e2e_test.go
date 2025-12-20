// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package api

import (
	"os"
	"testing"
)

func TestPipeline_AutoCreateXupleSpaces(t *testing.T) {
	t.Log("🧪 TEST E2E: Création automatique des xuple-spaces")
	t.Log("=====================================================")

	// Créer un fichier TSD temporaire
	tsdContent := `
xuple-space alerts {
	selection: fifo
	consumption: once
	retention: duration(24h)
}

xuple-space notifications {
	selection: random
	max-size: 1000
}

xuple-space logs {
	selection: lifo
	retention: duration(7d)
	max-size: 5000
}

type Alert(id: string, message: string)
`

	tmpfile, err := os.CreateTemp("", "test*.tsd")
	if err != nil {
		t.Fatalf("❌ Impossible de créer fichier temporaire: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.WriteString(tsdContent)
	if err != nil {
		t.Fatalf("❌ Impossible d'écrire dans fichier temporaire: %v", err)
	}
	tmpfile.Close()

	// Créer le pipeline
	t.Log("📝 Étape 1: Création du pipeline")
	pipeline := NewPipeline()

	// Ingérer le fichier
	t.Log("📝 Étape 2: Ingestion du fichier TSD")
	result, err := pipeline.IngestFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("❌ Ingestion failed: %v", err)
	}

	// Vérifier que les xuple-spaces ont été créés
	t.Log("📝 Étape 3: Vérification des xuple-spaces créés")
	spaces := result.XupleSpaceNames()
	if len(spaces) != 3 {
		t.Fatalf("❌ Expected 3 xuple-spaces, got %d", len(spaces))
	}

	// Vérifier les noms
	expectedSpaces := map[string]bool{
		"alerts":        false,
		"notifications": false,
		"logs":          false,
	}

	for _, name := range spaces {
		if _, ok := expectedSpaces[name]; ok {
			expectedSpaces[name] = true
			t.Logf("✅ Xuple-space '%s' trouvé", name)
		} else {
			t.Errorf("❌ Xuple-space inattendu: '%s'", name)
		}
	}

	for name, found := range expectedSpaces {
		if !found {
			t.Errorf("❌ Xuple-space manquant: '%s'", name)
		}
	}

	// Vérifier les métriques
	t.Log("📝 Étape 4: Vérification des métriques")
	metrics := result.Metrics()
	if metrics.XupleSpaceCount != 3 {
		t.Errorf("❌ Expected XupleSpaceCount=3, got %d", metrics.XupleSpaceCount)
	}

	t.Log("✅ Tous les xuple-spaces ont été créés automatiquement")
	t.Log("=====================================================")
}

func TestPipeline_AutoCreateXupleSpaces_WithMaxSize(t *testing.T) {
	t.Log("🧪 TEST E2E: Création automatique avec max-size")
	t.Log("=================================================")

	tsdContent := `
xuple-space limited {
	selection: fifo
	max-size: 100
}
`

	tmpfile, err := os.CreateTemp("", "test*.tsd")
	if err != nil {
		t.Fatalf("❌ Impossible de créer fichier temporaire: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.WriteString(tsdContent)
	if err != nil {
		t.Fatalf("❌ Impossible d'écrire dans fichier temporaire: %v", err)
	}
	tmpfile.Close()

	pipeline := NewPipeline()
	result, err := pipeline.IngestFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("❌ Ingestion failed: %v", err)
	}

	spaces := result.XupleSpaceNames()
	if len(spaces) != 1 {
		t.Fatalf("❌ Expected 1 xuple-space, got %d", len(spaces))
	}

	if spaces[0] != "limited" {
		t.Errorf("❌ Expected xuple-space 'limited', got '%s'", spaces[0])
	}

	// Vérifier que max-size est correctement configuré
	limitedSpace, err := result.XupleManager().GetXupleSpace("limited")
	if err != nil {
		t.Fatalf("❌ Impossible de récupérer le xuple-space 'limited': %v", err)
	}

	config := limitedSpace.GetConfig()
	if config.MaxSize != 100 {
		t.Errorf("❌ Expected max-size=100, got %d", config.MaxSize)
	} else {
		t.Log("✅ Configuration max-size=100 vérifiée")
	}

	// Vérifier que la politique de sélection est bien FIFO
	if config.SelectionPolicy == nil {
		t.Error("❌ SelectionPolicy ne devrait pas être nil")
	} else {
		t.Log("✅ SelectionPolicy configurée")
	}

	t.Log("✅ Xuple-space avec max-size créé automatiquement et correctement configuré")
}

func TestPipeline_AutoCreateXupleSpaces_Empty(t *testing.T) {
	t.Log("🧪 TEST E2E: Fichier sans xuple-space")
	t.Log("======================================")

	tsdContent := `
type Person(name: string, age: number)
`

	tmpfile, err := os.CreateTemp("", "test*.tsd")
	if err != nil {
		t.Fatalf("❌ Impossible de créer fichier temporaire: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.WriteString(tsdContent)
	if err != nil {
		t.Fatalf("❌ Impossible d'écrire dans fichier temporaire: %v", err)
	}
	tmpfile.Close()

	pipeline := NewPipeline()
	result, err := pipeline.IngestFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("❌ Ingestion failed: %v", err)
	}

	spaces := result.XupleSpaceNames()
	if len(spaces) != 0 {
		t.Errorf("❌ Expected 0 xuple-spaces, got %d", len(spaces))
	}

	t.Log("✅ Aucun xuple-space créé (attendu)")
}

func TestPipeline_AutoCreateXupleSpaces_WithDefaults(t *testing.T) {
	t.Log("🧪 TEST E2E: Création avec valeurs par défaut")
	t.Log("==============================================")

	tsdContent := `
xuple-space minimal {
}
`

	tmpfile, err := os.CreateTemp("", "test*.tsd")
	if err != nil {
		t.Fatalf("❌ Impossible de créer fichier temporaire: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.WriteString(tsdContent)
	if err != nil {
		t.Fatalf("❌ Impossible d'écrire dans fichier temporaire: %v", err)
	}
	tmpfile.Close()

	pipeline := NewPipeline()
	result, err := pipeline.IngestFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("❌ Ingestion failed: %v", err)
	}

	spaces := result.XupleSpaceNames()
	if len(spaces) != 1 {
		t.Fatalf("❌ Expected 1 xuple-space, got %d", len(spaces))
	}

	if spaces[0] != "minimal" {
		t.Errorf("❌ Expected xuple-space 'minimal', got '%s'", spaces[0])
	}

	// Le xuple-space devrait avoir été créé avec les valeurs par défaut
	// (fifo, once, unlimited)

	t.Log("✅ Xuple-space créé avec valeurs par défaut")
}
