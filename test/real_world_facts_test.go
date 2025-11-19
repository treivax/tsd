package testutil

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/treivax/tsd/constraint"
	"github.com/treivax/tsd/rete"
)

// Helper functions for test file management
func writeTestFile(filename, content string) error {
	return os.WriteFile(filename, []byte(content), 0644)
}

func removeTestFile(filename string) {
	os.Remove(filename)
}

func TestRealWorldGrammarFactsIntegration(t *testing.T) {
	fmt.Println("🧪 Test d'intégration avec fichiers réels: .constraint + .facts")

	// 1. Parser un fichier de contraintes existant
	constraintFile := filepath.Join("..", "beta_coverage_tests", "exists_simple.constraint")
	constraintResult, err := constraint.ParseConstraintFile(constraintFile)
	if err != nil {
		t.Fatalf("Erreur parsing fichier constraint: %v", err)
	}

	// 2. Parser le fichier .facts correspondant avec la grammaire
	factsFile := filepath.Join("..", "beta_coverage_tests", "exists_simple.facts")
	factsResult, err := constraint.ParseFactsFile(factsFile)
	if err != nil {
		t.Fatalf("Erreur parsing fichier .facts: %v", err)
	}

	// 3. Extraire les faits parsés
	parsedFacts, err := constraint.ExtractFactsFromProgram(factsResult)
	if err != nil {
		t.Fatalf("Erreur extraction faits: %v", err)
	}

	fmt.Printf("📊 Faits parsés du fichier .facts: %d\n", len(parsedFacts))
	for i, fact := range parsedFacts {
		fmt.Printf("  Fait %d: %v\n", i+1, fact)
	}

	// 4. Créer le réseau RETE et charger l'AST des contraintes
	storage := rete.NewMemoryStorage()
	network := rete.NewReteNetwork(storage)

	// Convertir les contraintes au format RETE
	constraintProgram, err := constraint.ConvertResultToProgram(constraintResult)
	if err != nil {
		t.Fatalf("Erreur conversion contraintes: %v", err)
	}

	reteProgram := constraint.ConvertToReteProgram(constraintProgram)
	err = network.LoadFromGenericAST(reteProgram)
	if err != nil {
		t.Fatalf("Erreur chargement contraintes dans réseau: %v", err)
	}

	// 5. Soumettre les faits parsés au réseau RETE
	err = network.SubmitFactsFromGrammar(parsedFacts)
	if err != nil {
		t.Fatalf("Erreur soumission faits au réseau: %v", err)
	}

	// 6. Vérifier que le réseau a traité les faits
	fmt.Printf("✅ Test intégration réussie avec %d faits parsés et %d nœuds terminaux\n",
		len(parsedFacts), len(network.TerminalNodes))

	// Vérifier qu'on a au moins quelques faits
	if len(parsedFacts) < 2 {
		t.Errorf("Attendu au moins 2 faits, obtenu %d", len(parsedFacts))
	}
}

func TestMixedConstraintFactsFile(t *testing.T) {
	fmt.Println("🧪 Test: Fichier .constraint avec types, contraintes ET faits mélangés")

	// Créer un fichier de test avec un mélange complet
	constraintContent := `// Test mélange types, faits et contraintes
type Person : <id: string, name: string, age: number>
type Order : <id: string, customer_id: string, amount: number, status: string>

// Quelques faits au début
Person(id:P001, name:Alice, age:25)
Order(id:O001, customer_id:P001, amount:100, status:pending)

// Règle 1: Personnes adultes
{p: Person} / p.age >= 18 ==> person_adult(p.id)

// Plus de faits au milieu
Person(id:P002, name:Bob, age:17)
Order(id:O002, customer_id:P002, amount:50, status:shipped)

// Règle 2: Commandes importantes
{o: Order} / o.amount > 75 ==> important_order(o.id)

// Encore des faits à la fin
Person(id:P003, name:Charlie, age:30)
Order(id:O003, customer_id:P003, amount:200, status:confirmed)
`

	// Créer le fichier temporaire
	tempFile := filepath.Join("..", "test", "mixed_test.constraint")
	err := writeTestFile(tempFile, constraintContent)
	if err != nil {
		t.Fatalf("Erreur création fichier test: %v", err)
	}
	defer removeTestFile(tempFile)

	// Parser le fichier mixte
	result, err := constraint.ParseConstraintFile(tempFile)
	if err != nil {
		t.Fatalf("Erreur parsing fichier mixte: %v", err)
	}

	// Valider le programme
	err = constraint.ValidateConstraintProgram(result)
	if err != nil {
		t.Fatalf("Erreur validation programme mixte: %v", err)
	}

	// Extraire les faits
	parsedFacts, err := constraint.ExtractFactsFromProgram(result)
	if err != nil {
		t.Fatalf("Erreur extraction faits du fichier mixte: %v", err)
	}

	// Convertir au format Program
	program, err := constraint.ConvertResultToProgram(result)
	if err != nil {
		t.Fatalf("Erreur conversion: %v", err)
	}

	fmt.Printf("📊 Analysé fichier mixte:\n")
	fmt.Printf("  - Types: %d\n", len(program.Types))
	fmt.Printf("  - Expressions: %d\n", len(program.Expressions))
	fmt.Printf("  - Faits: %d\n", len(program.Facts))

	// Vérifications
	if len(program.Types) != 2 {
		t.Errorf("Attendu 2 types, obtenu %d", len(program.Types))
	}
	if len(program.Expressions) != 2 {
		t.Errorf("Attendu 2 expressions, obtenu %d", len(program.Expressions))
	}
	if len(program.Facts) != 6 {
		t.Errorf("Attendu 6 faits, obtenu %d", len(program.Facts))
	}

	// Créer et tester le réseau
	storage := rete.NewMemoryStorage()
	network := rete.NewReteNetwork(storage)

	reteProgram := constraint.ConvertToReteProgram(program)
	err = network.LoadFromGenericAST(reteProgram)
	if err != nil {
		t.Fatalf("Erreur chargement réseau: %v", err)
	}

	err = network.SubmitFactsFromGrammar(parsedFacts)
	if err != nil {
		t.Fatalf("Erreur soumission faits: %v", err)
	}

	fmt.Printf("✅ Test fichier mixte réussi! Réseau avec %d nœuds terminaux\n",
		len(network.TerminalNodes))
}
