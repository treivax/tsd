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
Person(id:P001, name:Alice, age:25)
Person(id:P002, name:Bob, age:30)
Order(customer_id:P001, amount:100)
Order(customer_id:P001, amount:200)

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

	// 7. Vérifier que le réseau a traité les faits
	// On devrait avoir des activations dans les nœuds terminaux
	terminalCount := len(network.TerminalNodes)
	if terminalCount == 0 {
		t.Error("Aucun nœud terminal créé")
	}

	fmt.Printf("✅ Test d'intégration réussi! Réseau avec %d nœuds terminaux\n", terminalCount)
}

func TestPureFactsFile(t *testing.T) {
	fmt.Println("🧪 Test: Parser un fichier .facts pur avec la grammaire")

	// 1. Créer un fichier .facts pur
	factsContent := `Person(id:P001, name:Alice, age:25)
Person(id:P002, name:Bob, age:30)
Order(customer_id:P001, amount:100)
Order(customer_id:P002, amount:50)
`

	tempDir := os.TempDir()
	factsFile := filepath.Join(tempDir, "test_pure.facts")

	err := os.WriteFile(factsFile, []byte(factsContent), 0644)
	if err != nil {
		t.Fatalf("Erreur création fichier test: %v", err)
	}
	defer os.Remove(factsFile)

	// 2. Parser le fichier .facts avec la grammaire
	result, err := constraint.ParseFactsFile(factsFile)
	if err != nil {
		t.Fatalf("Erreur parsing fichier .facts: %v", err)
	}

	// 3. Extraire les faits parsés
	parsedFacts, err := constraint.ExtractFactsFromProgram(result)
	if err != nil {
		t.Fatalf("Erreur extraction faits: %v", err)
	}

	fmt.Printf("📊 Faits extraits du fichier .facts: %d\n", len(parsedFacts))
	for i, fact := range parsedFacts {
		fmt.Printf("  Fait %d: %v\n", i+1, fact)
	}

	// Vérifier qu'on a bien 4 faits
	if len(parsedFacts) != 4 {
		t.Errorf("Attendu 4 faits, obtenu %d", len(parsedFacts))
	}

	// Vérifier le contenu des faits
	personCount := 0
	orderCount := 0
	for _, fact := range parsedFacts {
		switch fact["type"] {
		case "Person":
			personCount++
		case "Order":
			orderCount++
		}
	}

	if personCount != 2 {
		t.Errorf("Attendu 2 faits Person, obtenu %d", personCount)
	}
	if orderCount != 2 {
		t.Errorf("Attendu 2 faits Order, obtenu %d", orderCount)
	}

	fmt.Printf("✅ Test fichier .facts pur réussi! %d Person, %d Order\n", personCount, orderCount)
}
