package testutil

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/treivax/tsd/constraint"
	"github.com/treivax/tsd/rete"
)

func TestGrammarFactsIntegration(t *testing.T) {
	fmt.Println("🧪 Test d'intégration: Grammaire -> Faits -> Réseau RETE")

	// 1. Créer un fichier de contraintes avec des faits intégrés
	constraintContent := `// Test d'intégration avec faits
type Person : <id: string, name: string, age: number>
type Order : <customer_id: string, amount: number>

// Faits définis directement dans le fichier .constraint
Person(id:"P001", name:"Alice", age:25)
Person(id:"P002", name:"Bob", age:30)
Order(customer_id:"P001", amount:100)
Order(customer_id:"P001", amount:200)

// Règle pour tester
{p: Person} / p.age > 20 ==> person_adult(p.id)
`

	tempDir := os.TempDir()
	constraintFile := filepath.Join(tempDir, "test_constraint_with_facts.constraint")

	err := os.WriteFile(constraintFile, []byte(constraintContent), 0644)
	if err != nil {
		t.Fatalf("Erreur création fichier test: %v", err)
	}
	defer os.Remove(constraintFile)

	// 2. Parser le fichier avec la grammaire modifiée
	result, err := constraint.ParseConstraintFile(constraintFile)
	if err != nil {
		t.Fatalf("Erreur parsing fichier: %v", err)
	}

	// 3. Valider le programme parsé
	err = constraint.ValidateConstraintProgram(result)
	if err != nil {
		t.Fatalf("Erreur validation programme: %v", err)
	}

	// 4. Extraire les faits parsés
	parsedFacts, err := constraint.ExtractFactsFromProgram(result)
	if err != nil {
		t.Fatalf("Erreur extraction faits: %v", err)
	}

	fmt.Printf("📊 Faits extraits: %d\n", len(parsedFacts))
	for i, fact := range parsedFacts {
		fmt.Printf("  Fait %d: %v\n", i+1, fact)
	}

	// Vérifier qu'on a bien 4 faits (2 Person, 2 Order)
	if len(parsedFacts) != 4 {
		t.Errorf("Attendu 4 faits, obtenu %d", len(parsedFacts))
	}

	// 5. Créer le réseau RETE et charger l'AST
	storage := rete.NewMemoryStorage()
	network := rete.NewReteNetwork(storage)

	// Convertir le résultat en Program pour le réseau
	program, err := constraint.ConvertResultToProgram(result)
	if err != nil {
		t.Fatalf("Erreur conversion en Program: %v", err)
	}

	// Convertir au format RETE
	reteProgram := constraint.ConvertToReteProgram(program)

	err = network.LoadFromGenericAST(reteProgram)
	if err != nil {
		t.Fatalf("Erreur chargement AST dans réseau: %v", err)
	}

	// 6. Soumettre les faits parsés au réseau RETE
	err = network.SubmitFactsFromGrammar(parsedFacts)
	if err != nil {
		t.Fatalf("Erreur soumission faits au réseau: %v", err)
	}

	// 7. Vérifier les résultats
	fmt.Println("✅ Test d'intégration réussi!")
}

func TestPureFactsFileGrammar(t *testing.T) {
	fmt.Println("🧪 Test: Fichier .facts pur (sans contraintes)")

	// Créer un fichier .facts pur
	factsContent := `// Fichier contenant uniquement des faits
Person(id:"P001", name:"Alice", age:25)
Person(id:"P002", name:"Bob", age:30)
Order(customer_id:"P001", amount:100)
`

	tempDir := os.TempDir()
	factsFile := filepath.Join(tempDir, "test_pure.facts")

	err := os.WriteFile(factsFile, []byte(factsContent), 0644)
	if err != nil {
		t.Fatalf("Erreur création fichier: %v", err)
	}
	defer os.Remove(factsFile)

	// Parser le fichier .facts
	result, err := constraint.ParseFactsFile(factsFile)
	if err != nil {
		t.Fatalf("Erreur parsing fichier .facts: %v", err)
	}

	// Extraire les faits
	facts, err := constraint.ExtractFactsFromProgram(result)
	if err != nil {
		t.Fatalf("Erreur extraction faits: %v", err)
	}

	fmt.Printf("📊 Faits extraits: %d\n", len(facts))
	if len(facts) != 3 {
		t.Errorf("Attendu 3 faits, obtenu %d", len(facts))
	}

	fmt.Println("✅ Fichier .facts pur parsé avec succès!")
}

func TestInlineFactsInConstraintFile(t *testing.T) {
	fmt.Println("🧪 Test: Fichier .constraint avec types, faits et règles")

	constraintContent := `// Fichier complet avec types, faits et règles
type Person : <id: string, name: string, age: number>

// Faits inline
Person(id:"P001", name:"Alice", age:25)
Person(id:"P002", name:"Bob", age:30)

// Règle
{p: Person} / p.age > 20 ==> adult(p.id)
`

	tempDir := os.TempDir()
	constraintFile := filepath.Join(tempDir, "test_mixed.constraint")

	err := os.WriteFile(constraintFile, []byte(constraintContent), 0644)
	if err != nil {
		t.Fatalf("Erreur création fichier: %v", err)
	}
	defer os.Remove(constraintFile)

	// Parser
	result, err := constraint.ParseConstraintFile(constraintFile)
	if err != nil {
		t.Fatalf("Erreur parsing: %v", err)
	}

	// Valider
	err = constraint.ValidateConstraintProgram(result)
	if err != nil {
		t.Fatalf("Erreur validation: %v", err)
	}

	// Extraire faits
	facts, err := constraint.ExtractFactsFromProgram(result)
	if err != nil {
		t.Fatalf("Erreur extraction faits: %v", err)
	}

	if len(facts) != 2 {
		t.Errorf("Attendu 2 faits, obtenu %d", len(facts))
	}

	fmt.Println("✅ Fichier mixte parsé avec succès!")
}
