package main

import (
	"fmt"
	"os"

	"github.com/treivax/tsd/internal/validation"
)

func main() {
	// Changer vers le répertoire de travail
	if err := os.Chdir("/home/resinsec/dev/tsd"); err != nil {
		fmt.Printf("Erreur changement répertoire: %v\n", err)
		return
	}

	fmt.Println("=== Test du nouveau système RETE avec jointures binaires en cascade ===")

	// Test spécifique pour beta_join_complex
	fmt.Println("\n=== Test focus: beta_join_complex ===")
	if err := testSpecificCase("beta_join_complex"); err != nil {
		fmt.Printf("❌ Erreur: %v\n", err)
	}
}

func testSpecificCase(baseName string) error {
	constraintFile := fmt.Sprintf("beta_coverage_tests/%s.constraint", baseName)
	factsFile := fmt.Sprintf("beta_coverage_tests/%s.facts", baseName)

	// Vérifier que les fichiers existent
	if _, err := os.Stat(constraintFile); os.IsNotExist(err) {
		return fmt.Errorf("fichier constraint manquant: %s", constraintFile)
	}
	if _, err := os.Stat(factsFile); os.IsNotExist(err) {
		return fmt.Errorf("fichier facts manquant: %s", factsFile)
	}

	fmt.Printf("📝 Contraintes: %s\n", constraintFile)
	fmt.Printf("📊 Faits: %s\n", factsFile)

	// Créer le réseau RETE
	network := validation.NewRETEValidationNetwork()

	// Charger les contraintes
	if err := network.ParseConstraintFile(constraintFile); err != nil {
		return fmt.Errorf("erreur parsing contraintes: %v", err)
	}

	// Charger les faits
	if err := network.LoadFactsFile(factsFile); err != nil {
		return fmt.Errorf("erreur chargement faits: %v", err)
	}

	// Debug du réseau
	network.Debug()

	// Obtenir les tokens terminaux
	terminals := network.GetTerminalTokens()

	fmt.Printf("🎯 Tokens terminaux générés: %d\n", len(terminals))

	// Afficher les détails
	for i, token := range terminals {
		fmt.Printf("  Token %d: %d faits (%s)\n", i+1, len(token.Facts), token.NodeID)
		for j, fact := range token.Facts {
			fmt.Printf("    %d. %s (%s) %+v\n", j+1, fact.ID, fact.Type, fact.Fields)
		}
	}

	if len(terminals) > 0 {
		fmt.Printf("✅ Succès: %d tokens générés\n", len(terminals))
		return nil
	} else {
		return fmt.Errorf("aucun token terminal généré")
	}
}
