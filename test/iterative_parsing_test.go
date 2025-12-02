// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package testutil

import (
	"fmt"
	"os"
	"testing"

	"github.com/treivax/tsd/constraint"
	"github.com/treivax/tsd/rete"
)

// Test complet du parsing itératif avec construction du réseau RETE
func TestIterativeParsingWithReteNetwork(t *testing.T) {
	fmt.Println("🚀 Test de parsing itératif avec réseau RETE")

	// Créer un parser itératif
	parser := constraint.NewIterativeParser()

	// Données de test avec types, règles et faits répartis
	typeContent := `
		// Types de test
		type Person(name: string, age:number)
		type Company(name: string, sector:string)
	`

	ruleContent := `
		// Règle de validation
		rule r1 : {p: Person} / p.age >= 18 ==> adult_status(p.name)
	`

	factContent := `
		// Faits de test
		Person(name:Alice, age:25)
		Person(name:Bob, age:17)
		Company(name:TechCorp, sector:IT)
	`

	// Parser les contenus de manière itérative
	fmt.Printf("📋 Parsing des types...\n")
	err := parser.ParseContent(typeContent, "types.tsd")
	if err != nil {
		t.Fatalf("Erreur parsing types: %v", err)
	}

	fmt.Printf("📋 Parsing des règles...\n")
	err = parser.ParseContent(ruleContent, "rules.tsd")
	if err != nil {
		t.Fatalf("Erreur parsing règles: %v", err)
	}

	fmt.Printf("📋 Parsing des faits...\n")
	err = parser.ParseContent(factContent, "facts.tsd")
	if err != nil {
		t.Fatalf("Erreur parsing faits: %v", err)
	}

	// Obtenir les statistiques de parsing
	stats := parser.GetParsingStatistics()
	fmt.Printf("📊 Statistiques finales: %+v\n", stats)

	// Vérifications
	if stats.TypesCount != 2 {
		t.Errorf("Attendu 2 types, obtenu %d", stats.TypesCount)
	}
	if stats.RulesCount != 1 {
		t.Errorf("Attendu 1 règle, obtenu %d", stats.RulesCount)
	}
	if stats.FactsCount != 3 {
		t.Errorf("Attendu 3 faits, obtenu %d", stats.FactsCount)
	}

	// Créer un pipeline et construire le réseau RETE
	fmt.Printf("🏗️ Construction du réseau RETE...\n")
	pipeline := rete.NewConstraintPipeline()
	storage := rete.NewMemoryStorage()

	network, err := pipeline.BuildNetworkFromIterativeParser(parser, storage)
	if err != nil {
		t.Fatalf("Erreur construction réseau: %v", err)
	}

	// Vérifications du réseau
	if len(network.TypeNodes) == 0 {
		t.Errorf("Aucun TypeNode créé")
	}
	if len(network.TerminalNodes) == 0 {
		t.Errorf("Aucun TerminalNode créé")
	}

	fmt.Printf("✅ Test réussi - Réseau créé avec %d TypeNodes et %d TerminalNodes\n",
		len(network.TypeNodes), len(network.TerminalNodes))
}

// Test du parsing multi-fichiers
func TestMultiFileParsing(t *testing.T) {
	fmt.Println("🚀 Test de parsing multi-fichiers")

	// Créer des fichiers temporaires
	files := []string{
		"/tmp/types.tsd",
		"/tmp/rules.tsd",
		"/tmp/facts.tsd",
	}

	// Contenu des fichiers
	contents := []string{
		`type Person(name: string, age:number)`,
		`rule r1 : {p: Person} / p.age >= 18 ==> adult_status(p.name)`,
		`Person(name:John, age:30)
Person(name:Jane, age:16)`,
	}

	// Créer les fichiers temporaires
	for i, content := range contents {
		err := createTempFile(files[i], content)
		if err != nil {
			t.Fatalf("Erreur création fichier %s: %v", files[i], err)
		}
		defer removeTempFile(files[i])
	}

	// Test du pipeline multi-fichiers
	pipeline := rete.NewConstraintPipeline()
	storage := rete.NewMemoryStorage()

	// Ingest files sequentially
	var network *rete.ReteNetwork
	for _, file := range files {
		network, err = pipeline.IngestFile(file, network, storage)
		if err != nil {
			t.Fatalf("Erreur ingestion fichier %s: %v", file, err)
		}
	}

	// Vérifications
	if len(network.TypeNodes) == 0 {
		t.Errorf("Aucun TypeNode créé")
	}

	fmt.Printf("✅ Test multi-fichiers réussi\n")
}

// Fonctions utilitaires pour les tests
func createTempFile(filename, content string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}

func removeTempFile(filename string) {
	os.Remove(filename)
}
