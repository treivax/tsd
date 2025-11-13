package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/treivax/tsd/rete"
)

// TestNegationRules teste spécifiquement les règles de négation (NotNode)
func TestNegationRules(t *testing.T) {
	fmt.Println("🔥 DÉMARRAGE TEST RÈGLES DE NÉGATION")
	fmt.Println("===================================")

	// Initialiser le helper avec le workspace TSD
	workspaceDir := "/home/resinsec/dev/tsd"
	helper := NewTestHelper()

	// Chemins vers les fichiers de négation
	constraintFile := filepath.Join(workspaceDir, "constraint", "test", "integration", "negation_rules.constraint")
	factsFile := filepath.Join(workspaceDir, "constraint", "test", "integration", "negation_rules.facts")

	// Vérification des fichiers
	if _, err := os.Stat(constraintFile); os.IsNotExist(err) {
		t.Fatalf("❌ Fichier contraintes négation introuvable: %s", constraintFile)
	}
	if _, err := os.Stat(factsFile); os.IsNotExist(err) {
		t.Fatalf("❌ Fichier faits négation introuvable: %s", factsFile)
	}

	fmt.Printf("✅ Fichier contraintes: %s\n", constraintFile)
	fmt.Printf("✅ Fichier faits: %s\n", factsFile)
	fmt.Println()

	// Traitement des contraintes de négation
	fmt.Println("🎯 TRAITEMENT CONTRAINTES DE NÉGATION")
	fmt.Println("====================================")

	network, facts, storage := helper.BuildNetworkFromConstraintFileWithMassiveFacts(t, constraintFile, factsFile)

	if network == nil {
		t.Fatal("❌ Réseau RETE non créé")
	}

	if len(facts) == 0 {
		t.Fatal("❌ Aucun fait chargé")
	}

	fmt.Printf("✅ %d faits négation chargés avec succès\n", len(facts))
	fmt.Printf("✅ Storage initialisé: %v\n", storage != nil)
	fmt.Println()

	// Analyse des règles de négation
	fmt.Println("🎯 ANALYSE RÈGLES DE NÉGATION")
	fmt.Println("=============================")

	// Lire le contenu des contraintes pour analyser les règles de négation
	constraintContent, err := os.ReadFile(constraintFile)
	if err != nil {
		t.Fatalf("❌ Erreur lecture fichier contraintes: %v", err)
	}

	content := string(constraintContent)

	// Compter les règles de négation
	notRules := strings.Count(content, "NOT (")
	totalRules := strings.Count(content, "==>")
	terminalNodes := len(network.TerminalNodes)

	fmt.Printf("📊 Règles totales: %d\n", totalRules)
	fmt.Printf("📊 Règles de négation (NOT): %d\n", notRules)
	fmt.Printf("📊 Nœuds terminaux: %d\n", terminalNodes)
	fmt.Printf("📊 Faits injectés: %d\n", len(facts))
	fmt.Println()

	// Test des résultats spécifiques aux négations
	fmt.Println("🧪 RÉSULTATS RÈGLES DE NÉGATION")
	fmt.Println("===============================")

	// Lister tous les nœuds terminaux disponibles
	fmt.Println("📋 Actions disponibles dans le réseau:")
	for actionName := range network.TerminalNodes {
		fmt.Printf("   - %s\n", actionName)
	}
	fmt.Println()

	// Test des négations avec analyse structurée par règle
	analyzeNegationRulesByRule(t, helper, network, facts, constraintFile)

	// Créer un fichier de résultats complet
	createCompleteResultsFile(t, helper, network, facts, constraintFile)

	fmt.Println()
	fmt.Printf("🎯 Test négation terminé: %d règles de négation analysées\n", notRules)
}

// NegationRule représente une règle de négation parsée
type NegationRule struct {
	RuleNumber   int
	TerminalName string
	RuleText     string
	ActionName   string
	Types        []string
	Condition    string
}

// analyzeNegationRulesByRule analyse les règles de négation une par une
func analyzeNegationRulesByRule(t *testing.T, helper *TestHelper, network *rete.ReteNetwork, facts []*rete.Fact, constraintFile string) {

	fmt.Println("🔍 ANALYSE PAR RÈGLE DE NÉGATION")
	fmt.Println("================================")

	// Parser les règles depuis le fichier constraint
	rules, err := parseNegationRules(constraintFile)
	if err != nil {
		t.Fatalf("❌ Erreur parsing règles: %v", err)
	}

	fmt.Printf("📊 %d règles de négation identifiées\n\n", len(rules))

	// Analyser chaque règle
	for _, rule := range rules {
		analyzeNegationRule(helper, network, facts, rule)
	}
}

// parseNegationRules parse le fichier constraint pour extraire les règles de négation
func parseNegationRules(constraintFile string) ([]NegationRule, error) {
	content, err := os.ReadFile(constraintFile)
	if err != nil {
		return nil, err
	}

	var rules []NegationRule
	lines := strings.Split(string(content), "\n")
	ruleNumber := 0

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Chercher les lignes contenant des règles avec NOT ou les règles positives
		if strings.Contains(line, "==>") && !strings.HasPrefix(line, "//") {
			rule := NegationRule{}
			rule.RuleNumber = ruleNumber
			rule.TerminalName = fmt.Sprintf("rule_%d_terminal", ruleNumber)
			rule.RuleText = line

			// Extraire l'action
			if actionMatch := regexp.MustCompile(`==>\s*(\w+)`).FindStringSubmatch(line); actionMatch != nil {
				rule.ActionName = actionMatch[1]
			}

			// Extraire les types
			if typesMatch := regexp.MustCompile(`\{([^}]+)\}`).FindStringSubmatch(line); typesMatch != nil {
				typesPart := typesMatch[1]
				for _, part := range strings.Split(typesPart, ",") {
					if colonIdx := strings.Index(part, ":"); colonIdx != -1 {
						typeName := strings.TrimSpace(part[colonIdx+1:])
						rule.Types = append(rule.Types, typeName)
					}
				}
			}

			// Extraire la condition
			if condMatch := regexp.MustCompile(`/\s*(.+?)\s*==>`).FindStringSubmatch(line); condMatch != nil {
				rule.Condition = strings.TrimSpace(condMatch[1])
			}

			rules = append(rules, rule)
			ruleNumber++
		}
	}

	return rules, nil
}

// analyzeNegationRule analyse une règle spécifique
func analyzeNegationRule(helper *TestHelper, network *rete.ReteNetwork, facts []*rete.Fact, rule NegationRule) {

	fmt.Printf("🎯 RÈGLE %d: %s\n", rule.RuleNumber, rule.ActionName)
	fmt.Printf("   Condition: %s\n", rule.Condition)
	fmt.Printf("   Types concernés: %v\n", rule.Types)
	fmt.Println("   " + strings.Repeat("-", 80))

	// Trouver le terminal correspondant
	terminal, exists := network.TerminalNodes[rule.TerminalName]
	if !exists {
		fmt.Printf("   ❌ Terminal %s introuvable\n\n", rule.TerminalName)
		return
	}

	// Afficher les faits soumis concernés par cette règle
	fmt.Println("   📥 FAITS SOUMIS (types concernés):")
	relevantFacts := getRelevantFacts(facts, rule.Types)
	if len(relevantFacts) == 0 {
		fmt.Println("      Aucun fait correspondant")
	} else {
		for i, fact := range relevantFacts {
			fmt.Printf("      - %s\n", helper.ShowFactDetails(fact, i+1))
		}
	}

	fmt.Printf("   📊 Total: %d faits soumis\n\n", len(relevantFacts))

	// Afficher les résultats dans le nœud terminal
	tokenCount := len(terminal.Memory.Tokens)
	fmt.Printf("   📤 RÉSULTATS TERMINAL (%s):\n", rule.TerminalName)

	if tokenCount == 0 {
		fmt.Println("      Aucun résultat (règle non déclenchée)")
	} else {
		fmt.Printf("      %d résultats obtenus\n", tokenCount)

		// Afficher tous les résultats
		count := 0
		for _, token := range terminal.Memory.Tokens {
			if len(token.Facts) > 0 {
				for j, fact := range token.Facts {
					fmt.Printf("      - Résultat %d: %s\n", count+1, helper.ShowFactDetails(fact, j+1))
				}
			}
			count++
		}
	}

	fmt.Printf("   📊 Taux de déclenchement: %d/%d (%.1f%%)\n", tokenCount, len(relevantFacts),
		float64(tokenCount)/float64(len(relevantFacts))*100)

	fmt.Println()
	fmt.Println()
}

// getRelevantFacts filtre les faits par types concernés
func getRelevantFacts(facts []*rete.Fact, types []string) []*rete.Fact {
	var relevant []*rete.Fact

	typeSet := make(map[string]bool)
	for _, t := range types {
		typeSet[t] = true
	}

	for _, fact := range facts {
		if typeSet[fact.Type] {
			relevant = append(relevant, fact)
		}
	}

	return relevant
}

// createCompleteResultsFile crée un fichier avec tous les résultats détaillés
func createCompleteResultsFile(t *testing.T, helper *TestHelper, network *rete.ReteNetwork, facts []*rete.Fact, constraintFile string) {

	// Parser les règles
	rules, err := parseNegationRules(constraintFile)
	if err != nil {
		t.Logf("❌ Erreur parsing règles pour fichier résultats: %v", err)
		return
	}

	// Créer le contenu du fichier
	var content strings.Builder
	content.WriteString("# RÉSULTATS COMPLETS - ANALYSE RÈGLES DE NÉGATION TSD\n")
	content.WriteString("=====================================================\n\n")
	content.WriteString(fmt.Sprintf("**Date d'exécution**: %s\n", "13 novembre 2025"))
	content.WriteString(fmt.Sprintf("**Fichier contraintes**: %s\n", constraintFile))
	content.WriteString(fmt.Sprintf("**Nombre de règles**: %d\n", len(rules)))
	content.WriteString(fmt.Sprintf("**Nombre de faits**: %d\n\n", len(facts)))

	// Analyser chaque règle et ajouter au contenu
	for _, rule := range rules {
		content.WriteString(fmt.Sprintf("## 🎯 RÈGLE %d: %s\n\n", rule.RuleNumber, rule.ActionName))
		content.WriteString(fmt.Sprintf("**Condition**: `%s`\n", rule.Condition))
		content.WriteString(fmt.Sprintf("**Types concernés**: %v\n", rule.Types))
		content.WriteString(fmt.Sprintf("**Terminal**: %s\n\n", rule.TerminalName))

		// Trouver le terminal correspondant
		terminal, exists := network.TerminalNodes[rule.TerminalName]
		if !exists {
			content.WriteString("❌ Terminal introuvable\n\n")
			continue
		}

		// Faits soumis
		relevantFacts := getRelevantFacts(facts, rule.Types)
		content.WriteString("### 📥 FAITS SOUMIS\n\n")
		if len(relevantFacts) == 0 {
			content.WriteString("Aucun fait correspondant\n\n")
		} else {
			for i, fact := range relevantFacts {
				content.WriteString(fmt.Sprintf("%d. %s\n", i+1, helper.ShowFactDetails(fact, 1)))
			}
			content.WriteString(fmt.Sprintf("\n**Total**: %d faits soumis\n\n", len(relevantFacts)))
		}

		// Résultats terminal
		tokenCount := len(terminal.Memory.Tokens)
		content.WriteString("### 📤 RÉSULTATS TERMINAL\n\n")

		if tokenCount == 0 {
			content.WriteString("Aucun résultat (règle non déclenchée)\n\n")
		} else {
			content.WriteString(fmt.Sprintf("**%d résultats obtenus**:\n\n", tokenCount))

			count := 0
			for _, token := range terminal.Memory.Tokens {
				if len(token.Facts) > 0 {
					count++
					content.WriteString(fmt.Sprintf("%d. **Token %d**:\n", count, count))
					for j, fact := range token.Facts {
						content.WriteString(fmt.Sprintf("   - Fait %d: %s\n", j+1, helper.ShowFactDetails(fact, 1)))
					}
					content.WriteString("\n")
				}
			}
		}

		// Statistiques
		if len(relevantFacts) > 0 {
			percentage := float64(tokenCount) / float64(len(relevantFacts)) * 100
			content.WriteString(fmt.Sprintf("### 📊 STATISTIQUES\n\n"))
			content.WriteString(fmt.Sprintf("- **Taux de déclenchement**: %d/%d (%.1f%%)\n", tokenCount, len(relevantFacts), percentage))
			content.WriteString(fmt.Sprintf("- **Efficacité**: %s\n", getEfficiencyLabel(percentage)))
		}

		content.WriteString("\n---\n\n")
	}

	// Ajouter un résumé global
	content.WriteString("## 📊 RÉSUMÉ GLOBAL\n\n")

	totalTerminals := len(network.TerminalNodes)
	activeTerminals := 0
	totalTokens := 0

	for _, terminal := range network.TerminalNodes {
		tokenCount := len(terminal.Memory.Tokens)
		totalTokens += tokenCount
		if tokenCount > 0 {
			activeTerminals++
		}
	}

	content.WriteString(fmt.Sprintf("- **Terminaux totaux**: %d\n", totalTerminals))
	content.WriteString(fmt.Sprintf("- **Terminaux actifs**: %d (%.1f%%)\n", activeTerminals, float64(activeTerminals)/float64(totalTerminals)*100))
	content.WriteString(fmt.Sprintf("- **Tokens générés**: %d\n", totalTokens))
	content.WriteString(fmt.Sprintf("- **Faits traités**: %d\n", len(facts)))

	// Créer le fichier
	outputPath := filepath.Join("/home/resinsec/dev/tsd/constraint/test/integration", "NEGATION_RESULTS_COMPLETE.md")
	err = os.WriteFile(outputPath, []byte(content.String()), 0644)
	if err != nil {
		t.Logf("❌ Erreur création fichier résultats: %v", err)
		return
	}

	fmt.Printf("✅ Fichier résultats complet créé: %s\n", outputPath)
}

// getEfficiencyLabel retourne un label d'efficacité selon le pourcentage
func getEfficiencyLabel(percentage float64) string {
	switch {
	case percentage >= 90:
		return "🟢 Très élevée"
	case percentage >= 70:
		return "🟡 Élevée"
	case percentage >= 50:
		return "🟠 Moyenne"
	case percentage >= 30:
		return "🔴 Faible"
	default:
		return "⚫ Très faible"
	}
}
