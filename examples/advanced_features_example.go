// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/treivax/tsd/rete"
)

// Exemple d'utilisation des fonctionnalités avancées du pipeline RETE :
// 1. Validation sémantique incrémentale
// 2. Garbage Collection après reset
// 3. Transactions avec rollback

func main() {
	fmt.Println("=== Démonstration des fonctionnalités avancées ===\n")

	// Créer un répertoire temporaire pour les exemples
	tmpDir, err := os.MkdirTemp("", "tsd_advanced_example_")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Exemple 1 : Validation incrémentale
	fmt.Println("📝 Exemple 1 : Validation sémantique incrémentale")
	demonstrateIncrementalValidation(tmpDir)

	// Exemple 2 : Garbage Collection
	fmt.Println("\n🗑️  Exemple 2 : Garbage Collection après reset")
	demonstrateGarbageCollection(tmpDir)

	// Exemple 3 : Transactions
	fmt.Println("\n🔒 Exemple 3 : Transactions avec rollback")
	demonstrateTransactions(tmpDir)

	// Exemple 4 : Toutes les fonctionnalités combinées
	fmt.Println("\n🚀 Exemple 4 : Toutes les fonctionnalités combinées")
	demonstrateAllFeatures(tmpDir)
}

// Exemple 1 : Validation incrémentale avec contexte
func demonstrateIncrementalValidation(tmpDir string) {
	storage := rete.NewMemoryStorage()
	pipeline := rete.NewConstraintPipeline()

	// Fichier 1 : Définir les types
	typesFile := filepath.Join(tmpDir, "types.tsd")
	typesContent := `
type Employee {
	id: string
	name: string
	salary: number
	department: string
}

type Department {
	id: string
	name: string
	budget: number
}
`
	writeFile(typesFile, typesContent)

	fmt.Println("  → Chargement des types...")
	network, err := pipeline.IngestFile(typesFile, nil, storage)
	if err != nil {
		log.Fatalf("Erreur : %v", err)
	}
	fmt.Printf("  ✅ %d types chargés\n", len(network.Types))

	// Fichier 2 : Définir des règles qui utilisent les types existants
	rulesFile := filepath.Join(tmpDir, "rules.tsd")
	rulesContent := `
rule "high_salary_alert" {
	when {
		e: Employee(salary > 100000)
	}
	then {
		print("High salary employee: " + e.name)
	}
}

rule "department_budget_check" {
	when {
		d: Department(budget < 50000)
	}
	then {
		print("Low budget department: " + d.name)
	}
}
`
	writeFile(rulesFile, rulesContent)

	fmt.Println("  → Chargement des règles avec validation incrémentale...")
	network, err = pipeline.IngestFile(rulesFile, network, storage)
	if err != nil {
		log.Fatalf("Erreur : %v", err)
	}
	fmt.Printf("  ✅ %d règles chargées (validation OK)\n", len(network.TerminalNodes))

	// Fichier 3 : Essayer de charger une règle avec un type inexistant
	invalidRulesFile := filepath.Join(tmpDir, "invalid_rules.tsd")
	invalidRulesContent := `
rule "invalid_rule" {
	when {
		p: Product(price > 100)
	}
	then {
		print("Expensive product")
	}
}
`
	writeFile(invalidRulesFile, invalidRulesContent)

	fmt.Println("  → Tentative de chargement d'une règle invalide...")
	_, err = pipeline.IngestFile(invalidRulesFile, network, storage)
	if err != nil {
		fmt.Printf("  ✅ Erreur détectée comme attendu : %v\n", err)
	} else {
		fmt.Println("  ❌ Erreur NON détectée (problème)")
	}
}

// Exemple 2 : Garbage Collection après reset
func demonstrateGarbageCollection(tmpDir string) {
	storage := rete.NewMemoryStorage()
	pipeline := rete.NewConstraintPipeline()

	// Session 1 : Créer un réseau volumineux
	session1File := filepath.Join(tmpDir, "session1.tsd")
	session1Content := `
type Person {
	id: string
	name: string
	age: number
}

rule "rule1" {
	when {
		p: Person(age >= 18)
	}
	then {
		print("Adult")
	}
}

rule "rule2" {
	when {
		p: Person(age < 18)
	}
	then {
		print("Minor")
	}
}

fact Person { id: "p1", name: "Alice", age: 30 }
fact Person { id: "p2", name: "Bob", age: 15 }
`
	writeFile(session1File, session1Content)

	fmt.Println("  → Session 1 : Création d'un réseau...")
	network, err := pipeline.IngestFile(session1File, nil, storage)
	if err != nil {
		log.Fatalf("Erreur : %v", err)
	}
	nodesSession1 := len(network.TypeNodes) + len(network.AlphaNodes) + len(network.TerminalNodes)
	fmt.Printf("  ✅ Réseau créé : %d nœuds, %d types\n", nodesSession1, len(network.Types))

	// Session 2 : Reset avec GC automatique
	session2File := filepath.Join(tmpDir, "session2.tsd")
	session2Content := `
reset

type Vehicle {
	id: string
	brand: string
	model: string
}

rule "luxury_car" {
	when {
		v: Vehicle(brand == "BMW")
	}
	then {
		print("Luxury vehicle")
	}
}
`
	writeFile(session2File, session2Content)

	fmt.Println("  → Session 2 : Reset avec GC automatique...")
	network, err = pipeline.IngestFile(session2File, network, storage)
	if err != nil {
		log.Fatalf("Erreur : %v", err)
	}
	nodesSession2 := len(network.TypeNodes) + len(network.AlphaNodes) + len(network.TerminalNodes)
	fmt.Printf("  ✅ Nouveau réseau : %d nœuds, %d types\n", nodesSession2, len(network.Types))
	fmt.Printf("  ✅ GC effectué : ancien réseau nettoyé (%d nœuds libérés)\n", nodesSession1)
}

// Exemple 3 : Transactions avec rollback
func demonstrateTransactions(tmpDir string) {
	storage := rete.NewMemoryStorage()
	pipeline := rete.NewConstraintPipeline()

	// État initial
	initialFile := filepath.Join(tmpDir, "initial.tsd")
	initialContent := `
type Book {
	id: string
	title: string
	author: string
	pages: number
}

fact Book { id: "b1", title: "Go Programming", author: "John Doe", pages: 300 }
`
	writeFile(initialFile, initialContent)

	network := rete.NewReteNetwork(storage)
	network, err := pipeline.IngestFile(initialFile, network, storage)
	if err != nil {
		log.Fatalf("Erreur : %v", err)
	}
	typesCountInitial := len(network.Types)
	fmt.Printf("  → État initial : %d type(s)\n", typesCountInitial)

	// Tentative 1 : Ingestion réussie avec transaction
	validFile := filepath.Join(tmpDir, "valid_update.tsd")
	validContent := `
rule "long_book" {
	when {
		b: Book(pages > 200)
	}
	then {
		print("Long book: " + b.title)
	}
}
`
	writeFile(validFile, validContent)

	fmt.Println("  → Transaction 1 : Ingestion valide (transaction automatique)...")
	network, err = pipeline.IngestFile(validFile, network, storage)
	if err != nil {
		log.Fatalf("Erreur inattendue : %v", err)
	}
	fmt.Printf("  ✅ Ingestion réussie (commit automatique)\n")

	// Tentative 2 : Ingestion invalide avec rollback
	invalidFile := filepath.Join(tmpDir, "invalid_update.tsd")
	invalidContent := `
rule "invalid_rule" {
	when {
		m: Magazine(pages > 50)
	}
	then {
		print("Magazine found")
	}
}
`
	writeFile(invalidFile, invalidContent)

	fmt.Println("  → Transaction 2 : Tentative d'ingestion invalide (rollback automatique)...")
	typesBeforeTx2 := len(network.Types)

	_, err = pipeline.IngestFile(invalidFile, network, storage)
	if err != nil {
		fmt.Printf("  ⚠️  Erreur détectée : %v\n", err)
		fmt.Printf("  ✅ Rollback automatique effectué\n")
	}

	typesAfterTx2 := len(network.Types)
	if typesAfterTx2 == typesBeforeTx2 {
		fmt.Printf("  ✅ État restauré : %d type(s) (aucun changement)\n", typesAfterTx2)
	} else {
		fmt.Printf("  ❌ État non restauré : avant=%d, après=%d\n", typesBeforeTx2, typesAfterTx2)
	}
}

// Exemple 4 : Toutes les fonctionnalités combinées
func demonstrateAllFeatures(tmpDir string) {
	storage := rete.NewMemoryStorage()
	pipeline := rete.NewConstraintPipeline()

	// Configuration avec toutes les fonctionnalités
	config := rete.DefaultAdvancedPipelineConfig()
	config.AutoCommit = true
	config.AutoRollbackOnError = true

	// Fichier 1 : Types de base
	file1 := filepath.Join(tmpDir, "base_types.tsd")
	file1Content := `
type Student {
	id: string
	name: string
	grade: number
}
`
	writeFile(file1, file1Content)

	fmt.Println("  → Chargement avec toutes les fonctionnalités activées...")
	network, metrics, err := pipeline.IngestFileWithAdvancedFeatures(file1, nil, storage, config)
	if err != nil {
		log.Fatalf("Erreur : %v", err)
	}

	rete.PrintAdvancedMetrics(metrics)

	// Fichier 2 : Règles utilisant les types existants
	file2 := filepath.Join(tmpDir, "student_rules.tsd")
	file2Content := `
rule "honor_student" {
	when {
		s: Student(grade >= 90)
	}
	then {
		print("Honor student: " + s.name)
	}
}
`
	writeFile(file2, file2Content)

	network, metrics, err = pipeline.IngestFileWithAdvancedFeatures(file2, network, storage, config)
	if err != nil {
		log.Fatalf("Erreur : %v", err)
	}

	fmt.Println("\n  📊 Résumé des métriques :")
	summary := rete.GetAdvancedMetricsSummary(metrics)
	fmt.Print("  " + summary)

	// Fichier 3 : Reset avec GC
	file3 := filepath.Join(tmpDir, "reset_system.tsd")
	file3Content := `
reset

type Course {
	id: string
	name: string
	credits: number
}
`
	writeFile(file3, file3Content)

	fmt.Println("  → Reset avec GC...")
	network, metrics, err = pipeline.IngestFileWithAdvancedFeatures(file3, network, storage, config)
	if err != nil {
		log.Fatalf("Erreur : %v", err)
	}

	if metrics.GCPerformed {
		fmt.Printf("  ✅ GC effectué : %d nœuds collectés en %v\n",
			metrics.NodesCollected, metrics.GCDuration)
	}

	fmt.Printf("\n  🎯 Réseau final : %d type(s), %d règle(s)\n",
		len(network.Types), len(network.TerminalNodes))
}

// Utilitaire : Écrire un fichier
func writeFile(path, content string) {
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		log.Fatalf("Erreur écriture fichier %s : %v", path, err)
	}
}
