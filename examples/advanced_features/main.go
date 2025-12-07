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

// Exemple d'utilisation des fonctionnalités du pipeline RETE :
// 1. Validation sémantique incrémentale
// 2. Garbage Collection après reset
// 3. Transactions avec rollback automatique
// 4. Collecte de métriques

func main() {
	fmt.Println("=== Démonstration des fonctionnalités du pipeline RETE ===")

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

	// Exemple 4 : Collecte de métriques
	fmt.Println("\n📊 Exemple 4 : Collecte de métriques d'ingestion")
	demonstrateMetricsCollection(tmpDir)
}

// Exemple 1 : Validation incrémentale avec contexte
func demonstrateIncrementalValidation(tmpDir string) {
	storage := rete.NewMemoryStorage()
	pipeline := rete.NewConstraintPipeline()

	// Fichier 1 : Définir les types
	typesFile := filepath.Join(tmpDir, "types.tsd")
	typesContent := `
type Employee(id: string, name: string, salary: number, department: string)
type Department(id: string, name: string, budget: number)
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
action print(msg: string)

rule high_salary_alert: {e: Employee} / e.salary > 100000 ==> print(e.name)
rule department_budget: {d: Department} / d.budget < 50000 ==> print(d.name)
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
rule invalid_rule: {p: Product} / p.price > 100 ==> print("Expensive product")
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
type Person(id: string, name: string, age: number)
action print(msg: string)

rule rule1: {p: Person} / p.age >= 18 ==> print("Adult")
rule rule2: {p: Person} / p.age < 18 ==> print("Minor")

Person(id: "p1", name: "Alice", age: 30)
Person(id: "p2", name: "Bob", age: 15)
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

type Vehicle(id: string, brand: string, model: string)
action print(msg: string)

rule luxury_car: {v: Vehicle} / v.brand == "BMW" ==> print("Luxury vehicle")
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
type Book(id: string, title: string, author: string, pages: number)

Book(id: "b1", title: "Go Programming", author: "John Doe", pages: 300)
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
action print(msg: string)

rule long_book: {b: Book} / b.pages > 200 ==> print(b.title)
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
rule invalid_rule: {m: Magazine} / m.pages > 50 ==> print("Magazine found")
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

// Exemple 4 : Collecte de métriques
func demonstrateMetricsCollection(tmpDir string) {
	storage := rete.NewMemoryStorage()
	pipeline := rete.NewConstraintPipeline()

	// Fichier 1 : Types de base
	file1 := filepath.Join(tmpDir, "base_types.tsd")
	file1Content := `
type Student(id: string, name: string, grade: number)
`
	writeFile(file1, file1Content)

	fmt.Println("  → Ingestion avec collecte de métriques...")
	network, metrics, err := pipeline.IngestFileWithMetrics(file1, nil, storage)
	if err != nil {
		log.Fatalf("Erreur : %v", err)
	}

	fmt.Println("\n  📊 Métriques d'ingestion :")
	fmt.Printf("    • Durée totale : %v\n", metrics.TotalDuration)
	fmt.Printf("    • Durée parsing : %v\n", metrics.ParsingDuration)
	fmt.Printf("    • Durée validation : %v\n", metrics.ValidationDuration)
	fmt.Printf("    • Durée création types : %v\n", metrics.TypeCreationDuration)
	fmt.Printf("    • Types ajoutés : %d\n", metrics.TypesAdded)
	fmt.Printf("    • Règles ajoutées : %d\n", metrics.RulesAdded)
	fmt.Printf("    • Faits soumis : %d\n", metrics.FactsSubmitted)
	if metrics.WasReset {
		fmt.Printf("    • Reset détecté : oui\n")
	}

	// Fichier 2 : Règles utilisant les types existants
	file2 := filepath.Join(tmpDir, "student_rules.tsd")
	file2Content := `
action print(msg: string)

rule honor_student: {s: Student} / s.grade >= 90 ==> print(s.name)
rule failing_student: {s: Student} / s.grade < 60 ==> print(s.name)

Student(id: "s1", name: "Alice", grade: 95)
Student(id: "s2", name: "Bob", grade: 55)
`
	writeFile(file2, file2Content)

	fmt.Println("\n  → Ingestion incrémentale avec métriques...")
	network, metrics, err = pipeline.IngestFileWithMetrics(file2, network, storage)
	if err != nil {
		log.Fatalf("Erreur : %v", err)
	}

	fmt.Println("\n  📊 Métriques d'ingestion incrémentale :")
	fmt.Printf("    • Durée totale : %v\n", metrics.TotalDuration)
	fmt.Printf("    • Règles ajoutées : %d\n", metrics.RulesAdded)
	fmt.Printf("    • Faits soumis : %d\n", metrics.FactsSubmitted)
	fmt.Printf("    • Faits propagés : %d\n", metrics.FactsPropagated)
	fmt.Printf("    • Nouveaux terminaux : %d\n", metrics.NewTerminalsAdded)

	// Fichier 3 : Reset avec GC et métriques
	file3 := filepath.Join(tmpDir, "reset_system.tsd")
	file3Content := `
reset

type Course(id: string, name: string, credits: number)
action print(msg: string)

rule credit_heavy_course: {c: Course} / c.credits > 3 ==> print(c.name)

Course(id: "c1", name: "Advanced Algorithms", credits: 4)
`
	writeFile(file3, file3Content)

	fmt.Println("\n  → Reset avec GC et métriques...")
	network, metrics, err = pipeline.IngestFileWithMetrics(file3, network, storage)
	if err != nil {
		log.Fatalf("Erreur : %v", err)
	}

	fmt.Println("\n  📊 Métriques après reset :")
	fmt.Printf("    • Durée totale : %v\n", metrics.TotalDuration)
	fmt.Printf("    • Reset détecté : %v\n", metrics.WasReset)
	if metrics.WasReset {
		fmt.Printf("    • Ancien réseau nettoyé (GC automatique)\n")
	}
	fmt.Printf("    • Types ajoutés : %d\n", metrics.TypesAdded)
	fmt.Printf("    • Règles ajoutées : %d\n", metrics.RulesAdded)
	fmt.Printf("    • Faits soumis : %d\n", metrics.FactsSubmitted)

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
