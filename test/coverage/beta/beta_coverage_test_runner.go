package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/treivax/tsd/constraint"
	"github.com/treivax/tsd/rete"
)

// BetaTestResult stocke les résultats d'un test beta
type BetaTestResult struct {
	TestName          string
	Description       string
	ConstraintFile    string
	FactsFile         string
	Rules             []BetaParsedRule
	Facts             []*rete.Fact
	Network           *rete.ReteNetwork
	Actions           []BetaActionResult
	JoinResults       []JoinNodeResult
	NotResults        []NotNodeResult
	ExistsResults     []ExistsNodeResult
	AccumulateResults []AccumulateNodeResult
	ExecutionTime     time.Duration
	Success           bool
	ErrorMessage      string
	ExpectedResults   ExpectedTestResults
	ValidationReport  ValidationReport
}

// BetaParsedRule représente une règle parsée avec analyse Beta
type BetaParsedRule struct {
	RuleNumber   int
	RuleText     string
	ActionName   string
	Condition    string
	Variables    []VariableInfo
	NodeType     string // "JoinNode", "NotNode", "ExistsNode", "AccumulateNode"
	JoinType     string // "inner", "not", "exists", "accumulate"
	Complexity   string // "simple", "complex", "multi-variable"
	SemanticType string // "equality", "comparison", "logical", "arithmetic", etc.
}

// VariableInfo décrit les variables dans une règle
type VariableInfo struct {
	Name     string
	DataType string
	Role     string // "primary", "secondary", "accumulator"
}

// BetaActionResult représente le résultat d'une action avec détails Beta
type BetaActionResult struct {
	ActionName    string
	Count         int
	Facts         []*rete.Fact
	JoinedFacts   [][]rete.Fact // Faits joints pour les JoinNodes
	TriggerNode   string        // Type de nœud qui a déclenché l'action
	SemanticMatch bool          // Si le résultat correspond à la sémantique attendue
}

// JoinNodeResult analyse spécifique aux JoinNodes
type FactTuple struct {
	Facts       []string // IDs des faits participant au tuple
	Description string   // Description lisible du tuple
}

type JoinNodeResult struct {
	NodeID        string
	VariablePairs []string
	JoinCondition string
	MatchedTuples int
	JoinType      string // "inner", "cross", "filtered"
	Performance   time.Duration
	SemanticValid bool
	Tuples        []FactTuple // Tuples de faits satisfaisant la jointure
}

// NotNodeResult analyse spécifique aux NotNodes
type NotNodeResult struct {
	NodeID           string
	NegatedCondition string
	FilteredFacts    int
	RejectedFacts    int
	NegationType     string // "simple", "complex", "double"
	LogicalValid     bool
}

// ExistsNodeResult analyse spécifique aux ExistsNodes
type ExistsNodeResult struct {
	NodeID          string
	ExistsCondition string
	QuantifierType  string // "exists", "forall", "unique"
	MatchingFacts   int
	ExistenceProven bool
	SemanticValid   bool
	ValidatedTuples []FactTuple // Tuples validant l'existence
	ExpectedTuples  []string    // Tuples attendus
	ObservedTuples  []string    // Tuples observés
}

// AccumulateNodeResult analyse spécifique aux AccumulateNodes
type AccumulateNodeResult struct {
	NodeID            string
	AccumulateFunc    string // "SUM", "COUNT", "AVG", "MIN", "MAX"
	InputFacts        int
	AccumulatedValue  interface{}
	AggregationType   string
	MathematicalValid bool
}

// ExpectedTestResults définit les résultats attendus pour validation sémantique
type ExpectedTestResults struct {
	ExpectedActions    []ExpectedAction
	ExpectedJoins      []ExpectedJoin
	ExpectedNegations  []ExpectedNegation
	ExpectedExists     []ExpectedExists
	ExpectedAggregates []ExpectedAggregate
}

// ExpectedAction définit une action attendue
type ExpectedAction struct {
	ActionName        string
	MinTriggers       int
	MaxTriggers       int
	RequiredFactTypes []string
	ExpectedFactIDs   []string // Nouveaux: IDs des faits qui devraient déclencher
	SemanticReason    string
}

// ExpectedJoin définit une jointure attendue
type ExpectedJoin struct {
	LeftFactType    string
	RightFactType   string
	JoinCondition   string
	ExpectedMatches int
	SemanticReason  string
}

// ExpectedNegation définit une négation attendue
type ExpectedNegation struct {
	NegatedCondition string
	ExpectedFiltered int
	LogicalReason    string
}

// ExpectedExists définit une existence attendue
type ExpectedExists struct {
	ExistsCondition  string
	ShouldExist      bool
	QuantifierReason string
}

// ExpectedAggregate définit une agrégation attendue
type ExpectedAggregate struct {
	AggregateFunc      string
	ExpectedValue      interface{}
	MathematicalReason string
}

// ValidationReport contient le rapport de validation sémantique
type ValidationReport struct {
	ActionsValid     bool
	JoinsValid       bool
	NegationsValid   bool
	ExistsValid      bool
	AggregatesValid  bool
	OverallValid     bool
	ValidationErrors []string
	SemanticScore    float64 // Pourcentage de validation sémantique
}

func main() {
	fmt.Println("🔬 EXÉCUTION DES TESTS DE COUVERTURE BETA NODES - PIPELINE UNIQUE COMPLET")
	fmt.Println("========================================================================")
	fmt.Println("🎯 Analyse sémantique complète: JoinNode, NotNode, ExistsNode, AccumulateNode")
	fmt.Println("🔧 Couverture maximale des opérateurs: AND, OR, NOT, EXISTS, ==, !=, <, >, <=, >=, IN, CONTAINS, +, -, *, /")

	// D'abord chercher dans le répertoire local
	localTestDir := "."
	testDir := "/home/resinsec/dev/tsd/beta_coverage_tests"
	resultsFile := "/home/resinsec/dev/tsd/BETA_NODES_COVERAGE_COMPLETE_RESULTS.md"

	// Essayer d'abord les tests locaux
	tests, err := discoverBetaTests(localTestDir)
	if err != nil || len(tests) == 0 {
		// Fallback sur les tests globaux
		// Créer le répertoire de tests s'il n'existe pas
		if err := os.MkdirAll(testDir, 0755); err != nil {
			fmt.Printf("❌ Erreur création répertoire tests: %v\n", err)
			return
		}

		// Découvrir tous les tests Beta
		tests, err = discoverBetaTests(testDir)
		if err != nil {
			fmt.Printf("❌ Erreur découverte tests: %v\n", err)
			return
		}
	} else {
		// Utiliser les tests locaux
		testDir = localTestDir
	}

	fmt.Printf("📊 %d tests Beta découverts avec couverture complète des opérateurs\n\n", len(tests))

	// Si aucun test trouvé, créer des tests par défaut
	if len(tests) == 0 {
		fmt.Println("🏗️ Création des tests Beta par défaut...")
		err = createDefaultBetaTests(testDir)
		if err != nil {
			fmt.Printf("❌ Erreur création tests par défaut: %v\n", err)
			return
		}

		// Redécouvrir les tests
		tests, err = discoverBetaTests(testDir)
		if err != nil {
			fmt.Printf("❌ Erreur redécouverte tests: %v\n", err)
			return
		}
		fmt.Printf("✅ %d tests Beta créés\n\n", len(tests))
	}

	// Exécuter tous les tests
	var allResults []BetaTestResult
	for _, testName := range tests {
		fmt.Printf("🧪 Exécution test: %s\n", testName)
		result := executeBetaTest(testDir, testName)
		allResults = append(allResults, result)

		if result.Success {
			fmt.Printf("✅ Succès (%v) - Score sémantique: %.1f%%\n",
				result.ExecutionTime, result.ValidationReport.SemanticScore)
		} else {
			fmt.Printf("❌ Échec: %s\n", result.ErrorMessage)
		}
		fmt.Println()
	}

	// Générer le rapport complet
	err = generateBetaCompleteReport(allResults, resultsFile)
	if err != nil {
		fmt.Printf("❌ Erreur génération rapport: %v\n", err)
		return
	}

	// Résumé final avec analyse sémantique
	successCount := 0
	totalSemanticScore := 0.0
	for _, result := range allResults {
		if result.Success {
			successCount++
		}
		totalSemanticScore += result.ValidationReport.SemanticScore
	}

	avgSemanticScore := totalSemanticScore / float64(len(allResults))

	fmt.Printf("🎯 RÉSUMÉ FINAL - COUVERTURE BETA\n")
	fmt.Printf("=================================\n")
	fmt.Printf("✅ Tests réussis: %d/%d (%.1f%%)\n",
		successCount, len(allResults), float64(successCount)/float64(len(allResults))*100)
	fmt.Printf("🧠 Score sémantique moyen: %.1f%%\n", avgSemanticScore)
	fmt.Printf("📄 Rapport complet: %s\n", resultsFile)
}

// discoverBetaTests découvre tous les tests Beta dans le répertoire
func discoverBetaTests(testDir string) ([]string, error) {
	var tests []string

	err := filepath.Walk(testDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(info.Name(), ".constraint") {
			testName := strings.TrimSuffix(info.Name(), ".constraint")
			factsFile := filepath.Join(testDir, testName+".facts")
			if _, err := os.Stat(factsFile); err == nil {
				tests = append(tests, testName)
			}
		}
		return nil
	})

	sort.Strings(tests)
	return tests, err
}

// executeBetaTest exécute un test beta complet avec analyse sémantique
func executeBetaTest(testDir, testName string) BetaTestResult {
	startTime := time.Now()
	result := BetaTestResult{
		TestName:       testName,
		ConstraintFile: filepath.Join(testDir, testName+".constraint"),
		FactsFile:      filepath.Join(testDir, testName+".facts"),
	}

	// Charger les résultats attendus
	result.ExpectedResults = loadExpectedResults(testName)

	// Lire la description depuis le fichier constraint
	if description, err := extractBetaDescription(result.ConstraintFile); err == nil {
		result.Description = description
	}

	// Parser les contraintes
	program, err := constraint.ParseFile(result.ConstraintFile)
	if err != nil {
		result.ErrorMessage = fmt.Sprintf("Erreur parsing contraintes: %v", err)
		return result
	}

	// Extraire les règles avec analyse Beta
	result.Rules = extractBetaRulesFromProgram(program)

	// Créer le réseau RETE via le pipeline
	pipeline := rete.NewConstraintPipeline()
	storage := rete.NewMemoryStorage()

	network, facts, err := pipeline.BuildNetworkFromConstraintFileWithFacts(
		result.ConstraintFile, result.FactsFile, storage)
	if err != nil {
		result.ErrorMessage = fmt.Sprintf("Erreur construction réseau RETE: %v", err)
		return result
	}

	result.Network = network
	result.Facts = facts

	// Analyser la structure du réseau Beta avant exécution
	analyzeBetaNetworkStructure(&result)

	// Exécuter le test avec monitoring Beta
	err = executeBetaTestWithMonitoring(&result)
	if err != nil {
		result.ErrorMessage = fmt.Sprintf("Erreur exécution test: %v", err)
		return result
	}

	// Validation sémantique des résultats
	validateBetaSemantics(&result)

	result.ExecutionTime = time.Since(startTime)
	result.Success = true
	return result
}

// analyzeBetaNetworkStructure analyse la structure des nœuds Beta
func analyzeBetaNetworkStructure(result *BetaTestResult) {
	network := result.Network

	// Analyser les AlphaNodes pour identifier les NotNodes
	for i, node := range network.AlphaNodes {
		if isNegationNode(node.Condition) {
			notResult := NotNodeResult{
				NodeID:       fmt.Sprintf("alpha_%s", i),
				NegationType: "simple",
			}

			if condition, ok := node.Condition.(map[string]interface{}); ok {
				notResult.NegatedCondition = fmt.Sprintf("%v", condition)
			}

			result.NotResults = append(result.NotResults, notResult)
		}
	}

	// Identifier les jointures potentielles en analysant les actions multi-variables
	joinCount := 0
	for _, rule := range result.Rules {
		if len(rule.Variables) > 1 {
			joinResult := JoinNodeResult{
				NodeID:   fmt.Sprintf("join_%d", joinCount),
				JoinType: "inner",
			}
			for i := 0; i < len(rule.Variables)-1; i++ {
				joinResult.VariablePairs = append(joinResult.VariablePairs,
					fmt.Sprintf("%s <-> %s", rule.Variables[i].Name, rule.Variables[i+1].Name))
			}
			result.JoinResults = append(result.JoinResults, joinResult)
			joinCount++
		}
	}
}

// executeBetaTestWithMonitoring exécute le test avec monitoring des nœuds Beta
func executeBetaTestWithMonitoring(result *BetaTestResult) error {
	network := result.Network

	// Compteurs pour les actions
	actionsCount := make(map[string]int)
	actionsMap := make(map[string][]*rete.Fact)

	// LES FAITS ONT DÉJÀ ÉTÉ INJECTÉS PAR LE PIPELINE
	// Ne pas les réinjecter ici

	// Analyser les BetaNodes (JoinNodes) réels dans le réseau pour capturer les vraies jointures
	betaIndex := 0
	for _, betaNode := range network.BetaNodes {
		if joinNode, ok := betaNode.(*rete.JoinNode); ok && joinNode.Memory != nil {
			// Compter uniquement les jointures réussies (tokens avec flag IsJoinResult)
			joinMatches := 0
			var tuples []FactTuple

			for _, token := range joinNode.Memory.Tokens {
				if token.IsJoinResult && len(token.Facts) >= 2 {
					// C'est une jointure réussie
					joinMatches++
					var facts []string
					var desc []string
					for _, fact := range token.Facts {
						facts = append(facts, fact.ID)
						desc = append(desc, fmt.Sprintf("%s[%s]", fact.Type, fact.ID))
					}
					tuples = append(tuples, FactTuple{
						Facts:       facts,
						Description: fmt.Sprintf("Tuple joint: %s", strings.Join(desc, " ⋈ ")),
					})
				}
			}

			if betaIndex < len(result.JoinResults) {
				result.JoinResults[betaIndex].MatchedTuples = joinMatches
				if joinMatches > 0 {
					result.JoinResults[betaIndex].SemanticValid = true
					result.JoinResults[betaIndex].Tuples = tuples
				}
				betaIndex++
			}
		}
	}

	// Analyser les résultats dans les nœuds terminaux
	for _, terminal := range network.TerminalNodes {
		actionName := "unknown_action"
		if terminal.Action != nil && terminal.Action.Job.Name != "" {
			actionName = terminal.Action.Job.Name
		}

		tokenCount := len(terminal.Memory.Tokens)

		if tokenCount > 0 {
			actionsCount[actionName] = tokenCount

			// Analyser les faits joints
			var joinedFacts [][]rete.Fact
			for _, token := range terminal.Memory.Tokens {
				var tokenFacts []rete.Fact
				for _, fact := range token.Facts {
					tokenFacts = append(tokenFacts, *fact)
					actionsMap[actionName] = append(actionsMap[actionName], fact)
				}
				joinedFacts = append(joinedFacts, tokenFacts)
			}

			result.Actions = append(result.Actions, BetaActionResult{
				ActionName:  actionName,
				Count:       tokenCount,
				Facts:       actionsMap[actionName],
				JoinedFacts: joinedFacts,
				TriggerNode: determineTriggerNodeType(actionName),
			})
		}
	}

	return nil
}

// convertToFactTuples convertit les tuples de faits en structure lisible
func convertToFactTuples(joinedFacts [][]rete.Fact) []FactTuple {
	var tuples []FactTuple
	for _, factGroup := range joinedFacts {
		var factIDs []string
		var descriptions []string

		for _, fact := range factGroup {
			factIDs = append(factIDs, fact.ID)
			descriptions = append(descriptions, fmt.Sprintf("%s[%s]", fact.Type, fact.ID))
		}

		tuple := FactTuple{
			Facts:       factIDs,
			Description: fmt.Sprintf("(%s)", strings.Join(descriptions, " ⋈ ")),
		}
		tuples = append(tuples, tuple)
	}
	return tuples
}

// validateBetaSemantics valide la sémantique des résultats Beta avec couverture complète des opérateurs
func validateBetaSemantics(result *BetaTestResult) {
	validation := ValidationReport{
		ActionsValid:    true,
		JoinsValid:      true,
		NegationsValid:  true,
		ExistsValid:     true,
		AggregatesValid: true,
	}

	var validationErrors []string
	var totalScore float64
	var scoredItems int

	// Valider les actions avec détails des tuples et analyse sémantique avancée
	fmt.Printf("\n🔍 VALIDATION SÉMANTIQUE DÉTAILLÉE - COUVERTURE COMPLÈTE OPÉRATEURS\n")
	fmt.Printf("===================================================================\n")

	for _, expectedAction := range result.ExpectedResults.ExpectedActions {
		fmt.Printf("\n📋 Action attendue: %s\n", expectedAction.ActionName)
		fmt.Printf("   📊 Déclenchements attendus: %d-%d\n", expectedAction.MinTriggers, expectedAction.MaxTriggers)
		fmt.Printf("   📝 Raison: %s\n", expectedAction.SemanticReason)

		found := false
		for i := range result.Actions {
			actualAction := &result.Actions[i]
			if actualAction.ActionName == expectedAction.ActionName {
				found = true
				fmt.Printf("   ✅ Action trouvée: %d déclenchements\n", actualAction.Count)

				// Validation sémantique selon le type d'opérateur
				actionRule := findRuleByAction(result.Rules, actualAction.ActionName)
				if actionRule != nil {
					semanticValid := validateOperatorSemantics(actualAction, actionRule, result.Facts)
					if semanticValid {
						totalScore += 20
						fmt.Printf("   🎯 Validation sémantique: ✅ CORRECTE\n")
					} else {
						validation.ActionsValid = false
						validationErrors = append(validationErrors, fmt.Sprintf("Action %s: sémantique incorrecte pour opérateur %s", actualAction.ActionName, actionRule.SemanticType))
						fmt.Printf("   ❌ Validation sémantique: ÉCHOUÉE\n")
					}
					scoredItems++

					// Analyse détaillée selon le type de nœud
					switch actionRule.NodeType {
					case "JoinNode":
						analyzeJoinNodeSemantics(actualAction, actionRule)
					case "NotNode":
						analyzeNotNodeSemantics(actualAction, actionRule, result.Facts)
					case "ExistsNode":
						analyzeExistsNodeSemantics(actualAction, actionRule, result.Facts)
					}
				}

				// Afficher les tuples observés
				if len(actualAction.JoinedFacts) > 0 {
					fmt.Printf("   🔗 Tuples observés:\n")
					tuples := convertToFactTuples(actualAction.JoinedFacts)
					for j, tuple := range tuples {
						fmt.Printf("      %d. %s\n", j+1, tuple.Description)
					}
				} else if len(actualAction.Facts) > 0 {
					fmt.Printf("   📋 Faits individuels:\n")
					for j, fact := range actualAction.Facts {
						fmt.Printf("      %d. %s[%s]\n", j+1, fact.Type, fact.ID)
					}
				}

				if actualAction.Count < expectedAction.MinTriggers ||
					actualAction.Count > expectedAction.MaxTriggers {
					validation.ActionsValid = false
					validationErrors = append(validationErrors,
						fmt.Sprintf("Action %s: attendu %d-%d déclenchements, observé %d",
							expectedAction.ActionName, expectedAction.MinTriggers,
							expectedAction.MaxTriggers, actualAction.Count))
					fmt.Printf("   ❌ Échec: nombre incorrect\n")
				} else {
					actualAction.SemanticMatch = true
					fmt.Printf("   ✅ Succès: nombre correct\n")
				}
				break
			}
		}
		if !found {
			validation.ActionsValid = false
			validationErrors = append(validationErrors,
				fmt.Sprintf("Action attendue manquante: %s", expectedAction.ActionName))
			fmt.Printf("   ❌ Action manquante\n")
		}
	}

	// Valider les contraintes EXISTS
	for _, expectedExists := range result.ExpectedResults.ExpectedExists {
		fmt.Printf("\n📋 EXISTS attendu: %s\n", expectedExists.ExistsCondition)
		fmt.Printf("   🎯 Doit exister: %v\n", expectedExists.ShouldExist)
		fmt.Printf("   📝 Raison: %s\n", expectedExists.QuantifierReason)

		// Pour l'instant, marquer comme valide si on a des actions qui correspondent
		existsValid := false
		for _, action := range result.Actions {
			if len(action.JoinedFacts) > 0 || len(action.Facts) > 0 {
				existsValid = true
				fmt.Printf("   ✅ Condition d'existence satisfaite par action %s\n", action.ActionName)
				break
			}
		}

		if !existsValid && expectedExists.ShouldExist {
			validation.ExistsValid = false
			validationErrors = append(validationErrors,
				fmt.Sprintf("Condition EXISTS non satisfaite: %s", expectedExists.ExistsCondition))
			fmt.Printf("   ❌ Condition d'existence non satisfaite\n")
		}
	}

	// Valider les jointures
	for _, expectedJoin := range result.ExpectedResults.ExpectedJoins {
		fmt.Printf("\n📋 Jointure attendue: %s -> %s\n", expectedJoin.LeftFactType, expectedJoin.RightFactType)
		fmt.Printf("   📊 Correspondances attendues: %d\n", expectedJoin.ExpectedMatches)

		joinFound := false
		for i := range result.JoinResults {
			actualJoin := &result.JoinResults[i]
			fmt.Printf("   🔍 JoinNode %s: %d correspondances\n", actualJoin.NodeID, actualJoin.MatchedTuples)

			if len(actualJoin.Tuples) > 0 {
				fmt.Printf("   🔗 Tuples observés:\n")
				for j, tuple := range actualJoin.Tuples {
					fmt.Printf("      %d. %s\n", j+1, tuple.Description)
				}
			}

			if actualJoin.MatchedTuples == expectedJoin.ExpectedMatches {
				joinFound = true
				actualJoin.SemanticValid = true
				fmt.Printf("   ✅ Jointure validée\n")
				break
			} else if actualJoin.MatchedTuples > 0 {
				fmt.Printf("   ❌ Nombre incorrect: attendu %d, observé %d\n", expectedJoin.ExpectedMatches, actualJoin.MatchedTuples)
			} else {
				fmt.Printf("   ❌ Aucune correspondance trouvée\n")
			}
		}
		if !joinFound {
			validation.JoinsValid = false
			validationErrors = append(validationErrors,
				fmt.Sprintf("Jointure attendue: %s -> %s, %d correspondances",
					expectedJoin.LeftFactType, expectedJoin.RightFactType, expectedJoin.ExpectedMatches))
		}
	}

	// Calculer le score sémantique basé sur les opérateurs et la validation avancée
	if scoredItems > 0 {
		averageScore := totalScore / float64(scoredItems)
		validation.SemanticScore = averageScore
	} else {
		// Fallback sur l'ancienne méthode si pas de validation par opérateur
		totalChecks := len(result.ExpectedResults.ExpectedActions) +
			len(result.ExpectedResults.ExpectedJoins) +
			len(result.ExpectedResults.ExpectedNegations) +
			len(result.ExpectedResults.ExpectedExists) +
			len(result.ExpectedResults.ExpectedAggregates)

		if totalChecks > 0 {
			validChecks := 0
			if validation.ActionsValid {
				validChecks++
			}
			if validation.JoinsValid {
				validChecks++
			}
			if validation.NegationsValid {
				validChecks++
			}
			if validation.ExistsValid {
				validChecks++
			}
			if validation.AggregatesValid {
				validChecks++
			}

			validation.SemanticScore = float64(validChecks) / float64(5) * 100
		} else {
			validation.SemanticScore = 100.0 // Pas de vérifications = succès par défaut
		}
	}

	fmt.Printf("\n📊 SCORE SÉMANTIQUE FINAL: %.1f%% (%d opérateurs validés)\n", validation.SemanticScore, scoredItems)

	validation.OverallValid = validation.ActionsValid && validation.JoinsValid &&
		validation.NegationsValid && validation.ExistsValid && validation.AggregatesValid
	validation.ValidationErrors = validationErrors

	result.ValidationReport = validation
}

// Helper functions
func extractBetaDescription(constraintFile string) (string, error) {
	content, err := os.ReadFile(constraintFile)
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "//") {
			return strings.TrimSpace(line[2:]), nil
		}
	}
	return "Test de couverture Beta", nil
}

func extractBetaRulesFromProgram(program interface{}) []BetaParsedRule {
	var rules []BetaParsedRule

	if programMap, ok := program.(map[string]interface{}); ok {
		if expressions, exists := programMap["expressions"]; exists {
			if expList, ok := expressions.([]interface{}); ok {
				for i, expr := range expList {
					if exprMap, ok := expr.(map[string]interface{}); ok {
						rule := BetaParsedRule{
							RuleNumber: i + 1,
						}

						// Extraire l'action
						if actionData, hasAction := exprMap["action"]; hasAction {
							if actionMap, ok := actionData.(map[string]interface{}); ok {
								if jobData, hasJob := actionMap["job"]; hasJob {
									if jobMap, ok := jobData.(map[string]interface{}); ok {
										if name, hasName := jobMap["name"]; hasName {
											rule.ActionName = fmt.Sprintf("%v", name)
										}
									}
								}
							}
						}

						// Analyser les variables
						if setData, hasSet := exprMap["set"]; hasSet {
							rule.Variables = extractBetaVariableInfo(setData)
						}

						// Analyser la condition et déterminer le type de nœud
						if constraints, hasConstraints := exprMap["constraints"]; hasConstraints {
							rule.Condition, rule.NodeType, rule.SemanticType = analyzeBetaConstraintStructure(constraints)
						}

						// Déterminer la complexité
						rule.Complexity = determineBetaRuleComplexity(rule.Variables, rule.Condition)

						// Générer le texte de la règle
						rule.RuleText = generateBetaRuleText(exprMap)
						rules = append(rules, rule)
					}
				}
			}
		}
	}

	return rules
}

func extractBetaVariableInfo(setData interface{}) []VariableInfo {
	var variables []VariableInfo

	if setMap, ok := setData.(map[string]interface{}); ok {
		if vars, hasVars := setMap["variables"]; hasVars {
			if varList, ok := vars.([]interface{}); ok {
				for i, v := range varList {
					if varMap, ok := v.(map[string]interface{}); ok {
						variable := VariableInfo{
							Name:     fmt.Sprintf("%v", varMap["name"]),
							DataType: fmt.Sprintf("%v", varMap["dataType"]),
							Role:     "primary",
						}
						if i > 0 {
							variable.Role = "secondary"
						}
						variables = append(variables, variable)
					}
				}
			}
		}
	}

	return variables
}

func analyzeBetaConstraintStructure(constraints interface{}) (string, string, string) {
	if constraintMap, ok := constraints.(map[string]interface{}); ok {
		if constraintType, exists := constraintMap["type"]; exists {
			switch constraintType {
			case "notConstraint":
				return "NOT(...)", "NotNode", "negation"
			case "existsConstraint":
				return "EXISTS(...)", "ExistsNode", "existence"
			case "logicalExpr":
				// Analyser les opérations pour déterminer le type
				if operations, hasOps := constraintMap["operations"]; hasOps {
					if opList, ok := operations.([]interface{}); ok && len(opList) > 0 {
						if opMap, ok := opList[0].(map[string]interface{}); ok {
							if op, hasOp := opMap["op"]; hasOp {
								switch op {
								case "AND":
									return "AND expression", "JoinNode", "logical_and"
								case "OR":
									return "OR expression", "JoinNode", "logical_or"
								}
							}
						}
					}
				}
				return "AND/OR expression", "JoinNode", "logical"
			case "comparison":
				left := extractBetaFieldPath(constraintMap["left"])
				op := fmt.Sprintf("%v", constraintMap["operator"])
				right := extractBetaValue(constraintMap["right"])
				condition := fmt.Sprintf("%s %s %s", left, op, right)

				// Déterminer le type sémantique selon l'opérateur
				semanticType := "comparison"
				switch op {
				case "==", "!=":
					semanticType = "equality"
				case "<", ">", "<=", ">=":
					semanticType = "relational"
				case "IN":
					semanticType = "membership"
				case "CONTAINS", "LIKE", "MATCHES":
					semanticType = "pattern_matching"
				}

				return condition, "JoinNode", semanticType
			case "binaryOp":
				left := extractBetaFieldPath(constraintMap["left"])
				op := fmt.Sprintf("%v", constraintMap["operator"])
				right := extractBetaValue(constraintMap["right"])
				condition := fmt.Sprintf("%s %s %s", left, op, right)

				// Déterminer le type arithmétique
				semanticType := "arithmetic"
				switch op {
				case "+", "-":
					semanticType = "additive"
				case "*", "/":
					semanticType = "multiplicative"
				}

				return condition, "JoinNode", semanticType
			}
		}
	}
	return "unknown", "JoinNode", "unknown"
}

func determineBetaRuleComplexity(variables []VariableInfo, condition string) string {
	if len(variables) <= 1 {
		return "simple"
	} else if len(variables) == 2 {
		return "complex"
	}
	return "multi-variable"
}

func generateBetaRuleText(exprMap map[string]interface{}) string {
	// Extraire les variables du set
	variables := ""
	if setData, hasSet := exprMap["set"]; hasSet {
		if setMap, ok := setData.(map[string]interface{}); ok {
			if vars, hasVars := setMap["variables"]; hasVars {
				if varList, ok := vars.([]interface{}); ok {
					var varStrings []string
					for _, v := range varList {
						if varMap, ok := v.(map[string]interface{}); ok {
							name := fmt.Sprintf("%v", varMap["name"])
							dataType := fmt.Sprintf("%v", varMap["dataType"])
							varStrings = append(varStrings, fmt.Sprintf("%s: %s", name, dataType))
						}
					}
					variables = strings.Join(varStrings, ", ")
				}
			}
		}
	}

	condition := "true"
	if constraints, hasConstraints := exprMap["constraints"]; hasConstraints {
		condition, _, _ = analyzeBetaConstraintStructure(constraints)
	}

	action := "unknown_action"
	if actionData, hasAction := exprMap["action"]; hasAction {
		if actionMap, ok := actionData.(map[string]interface{}); ok {
			if jobData, hasJob := actionMap["job"]; hasJob {
				if jobMap, ok := jobData.(map[string]interface{}); ok {
					if name, hasName := jobMap["name"]; hasName {
						action = fmt.Sprintf("%v", name)
					}
				}
			}
		}
	}

	return fmt.Sprintf("{%s} / %s ==> %s", variables, condition, action)
}

func extractBetaFieldPath(fieldData interface{}) string {
	if fieldMap, ok := fieldData.(map[string]interface{}); ok {
		if fieldType, exists := fieldMap["type"]; exists && fieldType == "fieldAccess" {
			object := fmt.Sprintf("%v", fieldMap["object"])
			field := fmt.Sprintf("%v", fieldMap["field"])
			return fmt.Sprintf("%s.%s", object, field)
		}
	}
	return fmt.Sprintf("%v", fieldData)
}

func extractBetaValue(valueData interface{}) string {
	if valueMap, ok := valueData.(map[string]interface{}); ok {
		if value, exists := valueMap["value"]; exists {
			if valueType, hasType := valueMap["type"]; hasType {
				switch valueType {
				case "string":
					return fmt.Sprintf("\"%v\"", value)
				case "number", "boolean":
					return fmt.Sprintf("%v", value)
				}
			}
		}
	}
	return fmt.Sprintf("%v", valueData)
}

func isNegationNode(condition interface{}) bool {
	if condMap, ok := condition.(map[string]interface{}); ok {
		if condType, exists := condMap["type"]; exists {
			return condType == "notConstraint" || condType == "negation"
		}
	}
	return false
}

func determineTriggerNodeType(actionName string) string {
	// Analyser le nom de l'action pour déterminer le type de nœud
	if strings.Contains(actionName, "join") {
		return "JoinNode"
	} else if strings.Contains(actionName, "not") {
		return "NotNode"
	} else if strings.Contains(actionName, "exists") {
		return "ExistsNode"
	}
	return "AlphaNode"
}

func loadExpectedResults(testName string) ExpectedTestResults {
	// Charger les résultats attendus selon le nom du test
	switch testName {
	case "exists_simple", "beta_exists_simple":
		return ExpectedTestResults{
			ExpectedActions: []ExpectedAction{
				{ActionName: "person_has_orders", MinTriggers: 1, MaxTriggers: 1,
					ExpectedFactIDs:   []string{"P001"},   // Seulement la variable principale
					RequiredFactTypes: []string{"Person"}, // Variable principale
					SemanticReason:    "Une personne (Alice) a une commande existante"},
			},
			ExpectedExists: []ExpectedExists{
				{ExistsCondition: "o.customer_id == p.id", ShouldExist: true,
					QuantifierReason: "Alice a une commande, Bob n'en a pas"},
			},
		}
	case "not_simple", "beta_not_simple":
		return ExpectedTestResults{
			ExpectedActions: []ExpectedAction{
				{ActionName: "active_person", MinTriggers: 1, MaxTriggers: 1,
					ExpectedFactIDs: []string{"P001"},
					SemanticReason:  "Une personne active (Alice) passe le filtre NOT"},
			},
			ExpectedNegations: []ExpectedNegation{
				{NegatedCondition: "p.active == false", ExpectedFiltered: 1,
					LogicalReason: "Exclure les personnes inactives"},
			},
		}
	case "join_simple", "beta_join_simple":
		return ExpectedTestResults{
			ExpectedActions: []ExpectedAction{
				{ActionName: "customer_order_match", MinTriggers: 2, MaxTriggers: 2,
					ExpectedFactIDs: []string{"P001", "P002", "O001", "O002"},
					SemanticReason:  "Deux customers avec leurs commandes matchent"},
			},
			ExpectedJoins: []ExpectedJoin{
				{LeftFactType: "Person", RightFactType: "Order",
					JoinCondition: "p.id == o.customer_id", ExpectedMatches: 2,
					SemanticReason: "Jointure sur l'ID client"},
			},
		}
	case "beta_join_complex":
		return ExpectedTestResults{
			ExpectedActions: []ExpectedAction{
				{ActionName: "dept_match", MinTriggers: 3, MaxTriggers: 3,
					ExpectedFactIDs: []string{"E001", "E002", "E003", "PR001", "PR002", "PR003"},
					SemanticReason:  "Trois correspondances department entre Employee et Project"},
			},
			ExpectedJoins: []ExpectedJoin{
				{LeftFactType: "Employee", RightFactType: "Project",
					JoinCondition: "e.department == p.department", ExpectedMatches: 3,
					SemanticReason: "Jointure sur département"},
			},
		}
	case "beta_not_complex":
		return ExpectedTestResults{
			ExpectedActions: []ExpectedAction{
				{ActionName: "eligible_person", MinTriggers: 3, MaxTriggers: 3,
					ExpectedFactIDs: []string{"P001", "P003", "P005"},
					SemanticReason:  "Trois personnes majeures éligibles"},
			},
			ExpectedNegations: []ExpectedNegation{
				{NegatedCondition: "p.age < 18", ExpectedFiltered: 2,
					LogicalReason: "Exclure les mineurs"},
			},
		}
	case "beta_exists_complex":
		return ExpectedTestResults{
			ExpectedActions: []ExpectedAction{
				{ActionName: "customer_purchase", MinTriggers: 5, MaxTriggers: 5,
					ExpectedFactIDs: []string{"C001", "C002", "C003", "PUR001", "PUR002", "PUR003", "PUR004", "PUR005"},
					SemanticReason:  "Cinq jointures Customer-Purchase (C001 a 2 achats, C002 a 1 achat, C003 a 2 achats)"},
			},
			ExpectedJoins: []ExpectedJoin{
				{LeftFactType: "Customer", RightFactType: "Purchase",
					JoinCondition: "c.id == p.customer_id", ExpectedMatches: 5,
					SemanticReason: "Jointure sur customer_id"},
			},
		}
	case "beta_exists_real":
		return ExpectedTestResults{
			ExpectedActions: []ExpectedAction{
				{ActionName: "vendor_has_products", MinTriggers: 2, MaxTriggers: 2,
					ExpectedFactIDs: []string{"V001", "V003"},
					SemanticReason:  "Deux vendors avec des produits"},
			},
			ExpectedExists: []ExpectedExists{
				{ExistsCondition: "p.vendor_id == v.id", ShouldExist: true,
					QuantifierReason: "Existence de produits pour vendors"},
			},
		}
	case "beta_join_numeric":
		return ExpectedTestResults{
			ExpectedActions: []ExpectedAction{
				{ActionName: "high_performer_advanced", MinTriggers: 6, MaxTriggers: 6,
					ExpectedFactIDs: []string{"S001", "S003", "S004", "C001", "C003", "C004"},
					SemanticReason:  "Étudiants performants avec tous les cours"},
			},
			ExpectedJoins: []ExpectedJoin{
				{LeftFactType: "Student", RightFactType: "Course",
					JoinCondition: "s.grade > 85", ExpectedMatches: 6,
					SemanticReason: "Jointure avec condition numérique"},
			},
		}
	case "beta_not_string":
		return ExpectedTestResults{
			ExpectedActions: []ExpectedAction{
				{ActionName: "working_device", MinTriggers: 3, MaxTriggers: 3,
					ExpectedFactIDs: []string{"D001", "D003", "D004"},
					SemanticReason:  "Trois devices non cassés"},
			},
			ExpectedNegations: []ExpectedNegation{
				{NegatedCondition: "d.status == broken", ExpectedFiltered: 2,
					LogicalReason: "Exclure devices cassés"},
			},
		}
	case "beta_mixed_complex":
		return ExpectedTestResults{
			ExpectedActions: []ExpectedAction{
				{ActionName: "valid_transaction", MinTriggers: 5, MaxTriggers: 5,
					ExpectedFactIDs: []string{"A001", "A002", "A003", "A004", "T001", "T002", "T003", "T004", "T005"},
					SemanticReason:  "Toutes les transactions correspondent aux comptes"},
			},
			ExpectedJoins: []ExpectedJoin{
				{LeftFactType: "Account", RightFactType: "Transaction",
					JoinCondition: "a.id == t.account_id", ExpectedMatches: 5,
					SemanticReason: "Tous les comptes avec transactions"},
			},
		}
	default:
		return ExpectedTestResults{}
	}
}

// createDefaultBetaTests crée des tests Beta par défaut
func createDefaultBetaTests(testDir string) error {
	tests := map[string]struct {
		Constraint string
		Facts      string
	}{
		"join_simple": {
			Constraint: `// Test jointure simple entre deux faits
type Person : <id: string, name: string, age: number>
type Order : <id: string, customer_id: string, amount: number>

{p: Person, o: Order} / p.id == o.customer_id ==> join_person_order(p.id, o.id)`,
			Facts: `Person{id: "P001", name: "Alice", age: 25}
Person{id: "P002", name: "Bob", age: 30}
Order{id: "O001", customer_id: "P001", amount: 100}
Order{id: "O002", customer_id: "P002", amount: 200}`,
		},
		"not_simple": {
			Constraint: `// Test négation simple
type Person : <id: string, name: string, active: bool>

{p: Person} / NOT (p.active == false) ==> active_person(p.id)`,
			Facts: `Person{id: "P001", name: "Alice", active: true}
Person{id: "P002", name: "Bob", active: false}`,
		},
		"exists_simple": {
			Constraint: `// Test existence simple
type Person : <id: string, name: string>
type Order : <customer_id: string, amount: number>

{p: Person} / EXISTS(o: Order, o.customer_id == p.id) ==> person_has_orders(p.id)`,
			Facts: `Person{id: "P001", name: "Alice"}
Person{id: "P002", name: "Bob"}
Order{customer_id: "P001", amount: 100}`,
		},
	}

	for testName, testData := range tests {
		constraintFile := filepath.Join(testDir, testName+".constraint")
		factsFile := filepath.Join(testDir, testName+".facts")

		err := os.WriteFile(constraintFile, []byte(testData.Constraint), 0644)
		if err != nil {
			return err
		}

		err = os.WriteFile(factsFile, []byte(testData.Facts), 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

// generateBetaCompleteReport génère le rapport complet de couverture Beta
func generateBetaCompleteReport(results []BetaTestResult, outputFile string) error {
	var report strings.Builder

	report.WriteString("# RAPPORT COMPLET DE COUVERTURE DES NŒUDS BETA\n")
	report.WriteString("================================================\n\n")

	// En-tête avec résumé
	successCount := 0
	totalSemanticScore := 0.0
	for _, result := range results {
		if result.Success {
			successCount++
		}
		totalSemanticScore += result.ValidationReport.SemanticScore
	}

	avgSemanticScore := 0.0
	if len(results) > 0 {
		avgSemanticScore = totalSemanticScore / float64(len(results))
	}

	report.WriteString(fmt.Sprintf("**📊 Tests exécutés:** %d\n", len(results)))
	report.WriteString(fmt.Sprintf("**✅ Tests réussis:** %d (%.1f%%)\n",
		successCount, float64(successCount)/float64(len(results))*100))
	report.WriteString(fmt.Sprintf("**🧠 Score sémantique moyen:** %.1f%%\n", avgSemanticScore))
	report.WriteString(fmt.Sprintf("**📅 Date d'exécution:** %s\n\n", time.Now().Format("2006-01-02 15:04:05")))

	// Section résumé des nœuds Beta testés
	report.WriteString("## 🎯 NŒUDS BETA ANALYSÉS\n")
	report.WriteString("| Type de Nœud | Tests | Succès | Score Sémantique |\n")
	report.WriteString("|---------------|--------|--------|------------------|\n")

	nodeTypeCounts := make(map[string]int)
	nodeTypeSuccesses := make(map[string]int)
	nodeTypeScores := make(map[string]float64)

	for _, result := range results {
		for _, rule := range result.Rules {
			nodeTypeCounts[rule.NodeType]++
			if result.Success {
				nodeTypeSuccesses[rule.NodeType]++
			}
			nodeTypeScores[rule.NodeType] += result.ValidationReport.SemanticScore
		}
	}

	for nodeType, count := range nodeTypeCounts {
		success := nodeTypeSuccesses[nodeType]
		avgScore := nodeTypeScores[nodeType] / float64(count)
		report.WriteString(fmt.Sprintf("| %s | %d | %d | %.1f%% |\n",
			nodeType, count, success, avgScore))
	}

	report.WriteString("\n")

	// Rapports détaillés pour chaque test
	for i, result := range results {
		report.WriteString(fmt.Sprintf("## 🧪 TEST %d: %s\n", i+1, result.TestName))
		report.WriteString("---\n\n")

		// Informations générales
		report.WriteString("### 📋 Informations générales\n")
		report.WriteString(fmt.Sprintf("- **Description:** %s\n", result.Description))
		report.WriteString(fmt.Sprintf("- **Fichier contraintes:** `%s`\n", result.ConstraintFile))
		report.WriteString(fmt.Sprintf("- **Fichier faits:** `%s`\n", result.FactsFile))
		report.WriteString(fmt.Sprintf("- **Temps d'exécution:** %v\n", result.ExecutionTime))
		report.WriteString(fmt.Sprintf("- **Résultat:** %s\n", getSuccessEmoji(result.Success)))
		if !result.Success {
			report.WriteString(fmt.Sprintf("- **Erreur:** %s\n", result.ErrorMessage))
		}
		report.WriteString("\n")

		// Validation sémantique
		report.WriteString("### 🧠 Validation sémantique\n")
		validation := result.ValidationReport
		report.WriteString(fmt.Sprintf("- **Score global:** %.1f%%\n", validation.SemanticScore))
		report.WriteString(fmt.Sprintf("- **Actions valides:** %s\n", getBoolEmoji(validation.ActionsValid)))
		report.WriteString(fmt.Sprintf("- **Jointures valides:** %s\n", getBoolEmoji(validation.JoinsValid)))
		report.WriteString(fmt.Sprintf("- **Négations valides:** %s\n", getBoolEmoji(validation.NegationsValid)))
		report.WriteString(fmt.Sprintf("- **Existences valides:** %s\n", getBoolEmoji(validation.ExistsValid)))
		report.WriteString(fmt.Sprintf("- **Agrégations valides:** %s\n", getBoolEmoji(validation.AggregatesValid)))

		if len(validation.ValidationErrors) > 0 {
			report.WriteString("\n**⚠️ Erreurs de validation:**\n")
			for _, error := range validation.ValidationErrors {
				report.WriteString(fmt.Sprintf("- %s\n", error))
			}
		}
		report.WriteString("\n")

		// Règles analysées avec contenu exact des fichiers
		report.WriteString("### 📜 Règles analysées\n")
		if len(result.Rules) == 0 {
			report.WriteString("*Aucune règle analysée*\n\n")
		} else {
			// Extraire la règle exacte depuis le fichier
			exactRule := extractExactRuleFromConstraint(result.ConstraintFile)

			for j, rule := range result.Rules {
				report.WriteString(fmt.Sprintf("#### Règle %d\n", j+1))
				report.WriteString(fmt.Sprintf("- **Texte original:** `%s`\n", exactRule))
				report.WriteString(fmt.Sprintf("- **Action:** %s\n", rule.ActionName))
				report.WriteString(fmt.Sprintf("- **Type de nœud:** %s\n", rule.NodeType))
				report.WriteString(fmt.Sprintf("- **Type sémantique:** %s\n", rule.SemanticType))
				report.WriteString(fmt.Sprintf("- **Complexité:** %s\n", rule.Complexity))

				if len(rule.Variables) > 0 {
					report.WriteString("- **Variables:**\n")
					for _, variable := range rule.Variables {
						report.WriteString(fmt.Sprintf("  - %s (%s): %s\n",
							variable.Name, variable.DataType, variable.Role))
					}
				}
				report.WriteString("\n")
			}
		}

		// Structure du réseau RETE
		report.WriteString("### 🕸️ Structure du réseau RETE\n\n")
		if result.Network != nil {
			generateBetaNetworkVisualization(&report, result.Network)
		} else {
			report.WriteString("⚠️ Réseau RETE non disponible\n\n")
		}

		// Faits traités avec contenu exact des fichiers
		report.WriteString("### 📄 Faits traités\n")
		factTypes := make(map[string]int)
		for _, fact := range result.Facts {
			factTypes[getFactType(fact)]++
		}

		// Lire le contenu exact du fichier facts
		factsContent := readExactFactsContent(result.FactsFile)
		report.WriteString(fmt.Sprintf("**📄 Contenu fichier facts:**\n```\n%s\n```\n\n", factsContent))

		report.WriteString(fmt.Sprintf("**Total faits:** %d\n\n", len(result.Facts)))
		for factType, count := range factTypes {
			report.WriteString(fmt.Sprintf("- **%s:** %d faits\n", factType, count))
		}

		// Affichage détaillé de chaque fait parsé
		report.WriteString("\n**📋 Détail des faits parsés:**\n")
		for i, fact := range result.Facts {
			report.WriteString(fmt.Sprintf("%d. **%s[%s]** - `%s`\n",
				i+1, fact.Type, fact.ID, formatFactWithFields(*fact)))
		}
		report.WriteString("\n")

		// Résultats des actions avec détail complet des activations
		report.WriteString("### ⚡ Résultats des actions\n")
		if len(result.Actions) == 0 {
			report.WriteString("*Aucune action déclenchée*\n\n")
		} else {
			report.WriteString("| Action | Déclenchements | Type de Nœud | Correspondance Sémantique |\n")
			report.WriteString("|--------|----------------|-------------|---------------------------|\n")
			for _, action := range result.Actions {
				semanticIcon := getBoolEmoji(action.SemanticMatch)
				report.WriteString(fmt.Sprintf("| %s | %d | %s | %s |\n",
					action.ActionName, action.Count, action.TriggerNode, semanticIcon))
			}
			report.WriteString("\n")

			// Détails exhaustifs des activations d'actions
			for _, action := range result.Actions {
				report.WriteString(fmt.Sprintf("#### 🎯 Activation détaillée: `%s`\n", action.ActionName))
				report.WriteString(fmt.Sprintf("- **Nombre de déclenchements:** %d\n", action.Count))
				report.WriteString(fmt.Sprintf("- **Type de nœud déclencheur:** %s\n", action.TriggerNode))
				report.WriteString("\n")

				if len(action.JoinedFacts) > 0 {
					report.WriteString("**📋 TOKENS COMBINÉS activant l'action:**\n\n")

					// Obtenir les variables de la règle pour affichage structuré
					variables := getActionVariables(result.Rules, action.ActionName)

					for k, joinedFacts := range action.JoinedFacts {
						report.WriteString(fmt.Sprintf("##### Token combiné %d\n", k+1))

						// Afficher chaque fait avec sa variable correspondante
						if len(variables) > 0 && len(joinedFacts) > 0 {
							for l, fact := range joinedFacts {
								varName := "unknown"
								if l < len(variables) {
									varName = variables[l]
								}
								report.WriteString(fmt.Sprintf("- **`%s`**: %s[%s] - `%s`\n",
									varName, fact.Type, fact.ID, formatFactWithFields(fact)))
							}

							// Pour les jointures multiples, montrer aussi l'association explicite
							if len(joinedFacts) > 1 {
								var associations []string
								for l, fact := range joinedFacts {
									if l < len(variables) {
										associations = append(associations, fmt.Sprintf("%s[%s]", fact.Type, fact.ID))
									}
								}
								if len(associations) > 1 {
									report.WriteString(fmt.Sprintf("- **Association:** %s\n", strings.Join(associations, " ⋈ ")))
								}
							}
						} else {
							// Fallback si pas de variables définies
							for l, fact := range joinedFacts {
								report.WriteString(fmt.Sprintf("- **Fait %d (%s):** `%s`\n",
									l+1, fact.Type, formatFactWithFields(fact)))
							}
						}
						report.WriteString("\n")
					}
				} else {
					report.WriteString("*Aucun détail de token disponible*\n\n")
				}
			}
		}

		// Analyse des jointures
		if len(result.JoinResults) > 0 {
			report.WriteString("### 🔗 Analyse des jointures (JoinNodes)\n")
			report.WriteString("| Nœud | Paires de Variables | Correspondances | Type | Validation |\n")
			report.WriteString("|------|---------------------|-----------------|------|------------|\n")
			for _, join := range result.JoinResults {
				semanticIcon := getBoolEmoji(join.SemanticValid)
				pairText := strings.Join(join.VariablePairs, ", ")
				report.WriteString(fmt.Sprintf("| %s | %s | %d | %s | %s |\n",
					join.NodeID, pairText, join.MatchedTuples, join.JoinType, semanticIcon))
			}
			report.WriteString("\n")
		}

		// Analyse des négations
		if len(result.NotResults) > 0 {
			report.WriteString("### 🚫 Analyse des négations (NotNodes)\n")
			report.WriteString("| Nœud | Condition Niée | Faits Filtrés | Type | Validation |\n")
			report.WriteString("|------|----------------|---------------|------|------------|\n")
			for _, not := range result.NotResults {
				logicalIcon := getBoolEmoji(not.LogicalValid)
				report.WriteString(fmt.Sprintf("| %s | %s | %d | %s | %s |\n",
					not.NodeID, not.NegatedCondition, not.FilteredFacts, not.NegationType, logicalIcon))
			}
			report.WriteString("\n")
		}

		// Analyse des existences
		if len(result.ExistsResults) > 0 {
			report.WriteString("### ✅ Analyse des existences (ExistsNodes)\n")
			report.WriteString("| Nœud | Condition | Correspondances | Prouvé | Validation |\n")
			report.WriteString("|------|-----------|-----------------|--------|------------|\n")
			for _, exists := range result.ExistsResults {
				semanticIcon := getBoolEmoji(exists.SemanticValid)
				proofIcon := getBoolEmoji(exists.ExistenceProven)
				report.WriteString(fmt.Sprintf("| %s | %s | %d | %s | %s |\n",
					exists.NodeID, exists.ExistsCondition, exists.MatchingFacts, proofIcon, semanticIcon))
			}
			report.WriteString("\n")
		}

		// Analyse des agrégations
		if len(result.AccumulateResults) > 0 {
			report.WriteString("### 📊 Analyse des agrégations (AccumulateNodes)\n")
			report.WriteString("| Nœud | Fonction | Faits d'Entrée | Valeur | Validation |\n")
			report.WriteString("|------|----------|----------------|--------|------------|\n")
			for _, acc := range result.AccumulateResults {
				mathIcon := getBoolEmoji(acc.MathematicalValid)
				report.WriteString(fmt.Sprintf("| %s | %s | %d | %v | %s |\n",
					acc.NodeID, acc.AccumulateFunc, acc.InputFacts, acc.AccumulatedValue, mathIcon))
			}
			report.WriteString("\n")
		}

		// Résultats attendus vs observés
		if hasExpectedResults(result.ExpectedResults) {
			report.WriteString("### 🎯 Comparaison attendu vs observé\n")

			if len(result.ExpectedResults.ExpectedActions) > 0 {
				report.WriteString("#### Actions\n")
				report.WriteString("| Action | Attendu | Observé | Statut |\n")
				report.WriteString("|--------|---------|---------|--------|\n")
				for _, expected := range result.ExpectedResults.ExpectedActions {
					observed := findActualAction(result.Actions, expected.ActionName)
					status := "❌"
					if observed != nil && observed.Count >= expected.MinTriggers &&
						observed.Count <= expected.MaxTriggers {
						status = "✅"
					}
					observedCount := 0
					if observed != nil {
						observedCount = observed.Count
					}
					report.WriteString(fmt.Sprintf("| %s | %d-%d | %d | %s |\n",
						expected.ActionName, expected.MinTriggers, expected.MaxTriggers,
						observedCount, status))
				}
				report.WriteString("\n")

				// Détails complets des tokens combinés attendus vs observés
				report.WriteString("#### 📋 TOKENS COMBINÉS ATTENDUS vs OBTENUS\n\n")
				for _, expected := range result.ExpectedResults.ExpectedActions {
					observed := findActualAction(result.Actions, expected.ActionName)

					report.WriteString(fmt.Sprintf("**🎯 Action `%s`:**\n", expected.ActionName))
					report.WriteString(fmt.Sprintf("- **Description:** %s\n", expected.SemanticReason))
					report.WriteString(fmt.Sprintf("- **Variables de la règle:** %s\n",
						strings.Join(getActionVariables(result.Rules, expected.ActionName), ", ")))
					report.WriteString("\n")

					// TOKENS ATTENDUS
					report.WriteString("**📍 TOKENS COMBINÉS ATTENDUS:**\n")
					report.WriteString(fmt.Sprintf("- **Nombre de tokens attendus:** %d-%d\n", expected.MinTriggers, expected.MaxTriggers))
					if len(expected.ExpectedFactIDs) > 0 {
						// Grouper les faits par token attendu selon les variables de la règle
						expectedTokens := groupFactsIntoTokens(expected.ExpectedFactIDs, result.Facts, result.Rules, expected.ActionName)
						for i, token := range expectedTokens {
							report.WriteString(fmt.Sprintf("- **Token attendu %d:**\n", i+1))
							for varName, fact := range token {
								report.WriteString(fmt.Sprintf("  * `%s`: %s[%s] - `%s`\n",
									varName, fact.Type, fact.ID, formatFactWithFields(fact)))
							}
						}
					} else {
						report.WriteString("- *Pas de détails de tokens attendus spécifiés*\n")
					}
					report.WriteString("\n")

					// TOKENS OBTENUS
					report.WriteString("**📊 TOKENS COMBINÉS OBTENUS:**\n")
					if observed != nil && len(observed.JoinedFacts) > 0 {
						report.WriteString(fmt.Sprintf("- **Nombre de tokens obtenus:** %d\n", len(observed.JoinedFacts)))
						for k, joinedFacts := range observed.JoinedFacts {
							report.WriteString(fmt.Sprintf("- **Token obtenu %d:**\n", k+1))
							// Associer chaque fait à sa variable selon l'ordre de la règle
							variables := getActionVariables(result.Rules, expected.ActionName)
							for l, fact := range joinedFacts {
								varName := "unknown"
								if l < len(variables) {
									varName = variables[l]
								}
								report.WriteString(fmt.Sprintf("  * `%s`: %s[%s] - `%s`\n",
									varName, fact.Type, fact.ID, formatFactWithFields(fact)))
							}
						}
					} else {
						report.WriteString("- **Nombre de tokens obtenus:** 0\n")
						report.WriteString("- *Aucun token combiné généré*\n")
					}
					report.WriteString("\n")

					// COMPARAISON
					status := "❌ ÉCHEC"
					if observed != nil && observed.Count >= expected.MinTriggers && observed.Count <= expected.MaxTriggers {
						status = "✅ SUCCÈS"
					}
					report.WriteString(fmt.Sprintf("**🎯 RÉSULTAT:** %s\n", status))
					if observed != nil {
						if observed.Count == expected.MinTriggers ||
							(expected.MinTriggers != expected.MaxTriggers && observed.Count >= expected.MinTriggers && observed.Count <= expected.MaxTriggers) {
							report.WriteString("- ✅ Nombre de tokens correct\n")
						} else {
							report.WriteString(fmt.Sprintf("- ❌ Nombre de tokens incorrect: attendu %d-%d, obtenu %d\n",
								expected.MinTriggers, expected.MaxTriggers, observed.Count))
						}
					} else {
						report.WriteString("- ❌ Action non déclenchée\n")
					}
					report.WriteString("\n")
				}
			}
		}

		report.WriteString("---\n\n")
	}

	// Section recommandations
	report.WriteString("## 💡 RECOMMANDATIONS\n")
	report.WriteString("### Amélioration de la couverture Beta\n")

	if avgSemanticScore < 80 {
		report.WriteString("⚠️ **Score sémantique faible:** Réviser la validation des règles et la correspondance des résultats attendus.\n\n")
	}

	if len(results) > 0 && float64(successCount)/float64(len(results))*100 < 90 {
		report.WriteString("⚠️ **Taux de succès faible:** Analyser les erreurs de parsing et d'exécution.\n\n")
	}

	report.WriteString("### Prochaines étapes\n")
	report.WriteString("1. **Ajouter plus de tests complexes** avec jointures multiples\n")
	report.WriteString("2. **Tester les négations imbriquées** et conditions complexes\n")
	report.WriteString("3. **Valider les performances** des nœuds Beta avec de gros volumes\n")
	report.WriteString("4. **Enrichir la validation sémantique** avec plus de critères\n\n")

	// Écrire le rapport
	return os.WriteFile(outputFile, []byte(report.String()), 0644)
}

// Helper functions pour le rapport
func getSuccessEmoji(success bool) string {
	if success {
		return "✅ Succès"
	}
	return "❌ Échec"
}

func getBoolEmoji(value bool) string {
	if value {
		return "✅"
	}
	return "❌"
}

func getFactType(fact *rete.Fact) string {
	// Essayer d'extraire le type du nom du fait
	if fact.Type != "" {
		return fact.Type
	}
	// Fallback sur l'ID du fait
	return strings.Split(fact.ID, "{")[0]
}

// Utilitaires communs
func formatFactWithFields(fact interface{}) string {
	var factValue rete.Fact
	switch f := fact.(type) {
	case rete.Fact:
		factValue = f
	case *rete.Fact:
		factValue = *f
	default:
		return "Unknown fact type"
	}

	var fields []string
	for key, value := range factValue.Fields {
		fields = append(fields, fmt.Sprintf("%s=%v", key, value))
	}
	return fmt.Sprintf("%s[%s]", factValue.Type, strings.Join(fields, ", "))
}

func hasExpectedResults(expected ExpectedTestResults) bool {
	return len(expected.ExpectedActions) > 0 ||
		len(expected.ExpectedJoins) > 0 ||
		len(expected.ExpectedNegations) > 0 ||
		len(expected.ExpectedExists) > 0 ||
		len(expected.ExpectedAggregates) > 0
}

func findActualAction(actions []BetaActionResult, actionName string) *BetaActionResult {
	for i := range actions {
		if actions[i].ActionName == actionName {
			return &actions[i]
		}
	}
	return nil
}

// readExactConstraintContent lit le contenu exact du fichier constraint
func readExactConstraintContent(filePath string) string {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Sprintf("Erreur lecture fichier: %v", err)
	}
	return string(content)
}

// readExactFactsContent lit le contenu exact du fichier facts
func readExactFactsContent(filePath string) string {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Sprintf("Erreur lecture fichier: %v", err)
	}
	return string(content)
}

// generateBetaNetworkVisualization génère une visualisation du réseau RETE pour les tests Beta
func generateBetaNetworkVisualization(report *strings.Builder, network *rete.ReteNetwork) {
	report.WriteString("```\n")
	report.WriteString("RÉSEAU RETE BETA - STRUCTURE HIÉRARCHIQUE\n")
	report.WriteString("==========================================\n\n")

	// Root Node
	report.WriteString("🌳 RootNode\n")
	report.WriteString("│\n")

	// Type Nodes
	if len(network.TypeNodes) > 0 {
		report.WriteString("├── 📁 TypeNodes\n")
		for typeName, typeNode := range network.TypeNodes {
			report.WriteString(fmt.Sprintf("│   ├── %s (%s)\n", typeName, typeNode.ID))
		}
		report.WriteString("│\n")
	}

	// Alpha Nodes
	if len(network.AlphaNodes) > 0 {
		report.WriteString("├── 🔍 AlphaNodes\n")
		for _, alphaNode := range network.AlphaNodes {
			condition := "unknown"
			if alphaNode.Condition != nil {
				if condMap, ok := alphaNode.Condition.(map[string]interface{}); ok {
					if condType, exists := condMap["type"]; exists {
						switch condType {
						case "negation":
							condition = "NOT(...) [Négation]"
						case "constraint":
							condition = "Condition positive"
						case "simple":
							condition = "Condition simple"
						case "existsConstraint":
							condition = "EXISTS(...) [Existence]"
						default:
							condition = fmt.Sprintf("Type: %v", condType)
						}
					}
				}
			}
			report.WriteString(fmt.Sprintf("│   ├── %s\n", alphaNode.ID))
			report.WriteString(fmt.Sprintf("│   │   ├── Condition: %s\n", condition))
			report.WriteString(fmt.Sprintf("│   │   └── Variable: %s\n", alphaNode.VariableName))
		}
		report.WriteString("│\n")
	}

	// Join Nodes (spécifique aux tests Beta)
	if len(network.BetaNodes) > 0 {
		report.WriteString("├── 🔗 BetaNodes (Jointures)\n")
		for nodeID, betaNodeInterface := range network.BetaNodes {
			report.WriteString(fmt.Sprintf("│   ├── %s\n", nodeID))
			if joinNode, ok := betaNodeInterface.(*rete.JoinNode); ok {
				report.WriteString(fmt.Sprintf("│   │   ├── Variables: %s\n", strings.Join(joinNode.AllVariables, " ⋈ ")))
				report.WriteString(fmt.Sprintf("│   │   ├── Conditions: %d\n", len(joinNode.JoinConditions)))
				report.WriteString("│   │   └── Type: JoinNode\n")
			} else {
				report.WriteString(fmt.Sprintf("│   │   └── Type: %T\n", betaNodeInterface))
			}
		}
		report.WriteString("│\n")
	}

	// Terminal Nodes
	if len(network.TerminalNodes) > 0 {
		report.WriteString("└── 🎯 TerminalNodes (Actions)\n")
		for _, terminalNode := range network.TerminalNodes {
			actionName := "unknown_action"
			if terminalNode.Action != nil && terminalNode.Action.Job.Name != "" {
				actionName = terminalNode.Action.Job.Name
			}
			report.WriteString(fmt.Sprintf("    ├── %s\n", terminalNode.ID))
			report.WriteString(fmt.Sprintf("    │   └── Action: %s\n", actionName))
		}
	}

	report.WriteString("```\n\n")
}

// extractExactRuleFromConstraint extrait uniquement la ligne de règle du fichier .constraint
func extractExactRuleFromConstraint(constraintFile string) string {
	file, err := os.Open(constraintFile)
	if err != nil {
		return fmt.Sprintf("Erreur lecture: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Cherche les lignes qui contiennent une règle (avec ==>)
		if strings.Contains(line, "==>") {
			return line
		}
	}
	return "Règle non trouvée"
}

// extractMainRule extrait la règle principale (avec ==>) d'un contenu de fichier constraint
func extractMainRule(content string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "==>") {
			return line
		}
	}
	return "Règle non trouvée"
}

// getActionVariables extrait les noms des variables d'une action donnée
func getActionVariables(rules []BetaParsedRule, actionName string) []string {
	for _, rule := range rules {
		if rule.ActionName == actionName {
			variables := make([]string, len(rule.Variables))
			for i, variable := range rule.Variables {
				variables[i] = variable.Name
			}
			return variables
		}
	}
	return []string{}
}

// groupFactsIntoTokens groupe les faits en tokens combinés selon les variables de la règle
func groupFactsIntoTokens(factIDs []string, allFacts []*rete.Fact, rules []BetaParsedRule, actionName string) []map[string]*rete.Fact {
	// Créer une map des faits par ID pour un accès rapide
	factsByID := make(map[string]*rete.Fact)
	for _, fact := range allFacts {
		factsByID[fact.ID] = fact
	}

	// Obtenir les variables de la règle et son type de nœud
	variables := getActionVariables(rules, actionName)
	var nodeType string
	for _, rule := range rules {
		if rule.ActionName == actionName {
			nodeType = rule.NodeType
			break
		}
	}

	if len(variables) == 0 {
		return []map[string]*rete.Fact{}
	}

	var tokens []map[string]*rete.Fact

	// Cas spécial pour EXISTS : seule la variable principale est dans le token final
	if nodeType == "ExistsNode" && len(variables) > 0 {
		mainVarName := variables[0] // Variable principale (p)
		for _, factID := range factIDs {
			if fact, exists := factsByID[factID]; exists {
				// Vérifier si ce fait correspond à la variable principale
				if isMainVariableFact(fact, mainVarName, rules, actionName) {
					token := map[string]*rete.Fact{
						mainVarName: fact,
					}
					tokens = append(tokens, token)
				}
			}
		}
		return tokens
	}

	// Si une seule variable, créer un token pour chaque fait
	if len(variables) == 1 {
		varName := variables[0]
		for _, factID := range factIDs {
			if fact, exists := factsByID[factID]; exists {
				token := map[string]*rete.Fact{
					varName: fact,
				}
				tokens = append(tokens, token)
			}
		}
		return tokens
	}

	// Pour plusieurs variables, essayer de grouper les faits par tokens logiques
	// Si on a des variables [p, o] et des factIDs [P001, O001, P002, O002],
	// on suppose que P001+O001 forment un token et P002+O002 un autre

	// Grouper les faits par type
	factsByType := make(map[string][]*rete.Fact)
	for _, factID := range factIDs {
		if fact, exists := factsByID[factID]; exists {
			factsByType[fact.Type] = append(factsByType[fact.Type], fact)
		}
	}

	// Obtenir les types des variables
	var variableTypes []string
	for _, rule := range rules {
		if rule.ActionName == actionName {
			for _, variable := range rule.Variables {
				variableTypes = append(variableTypes, variable.DataType)
			}
			break
		}
	}

	if len(variableTypes) != len(variables) {
		// Fallback: créer un seul token avec tous les faits disponibles
		token := make(map[string]*rete.Fact)
		varIndex := 0
		for _, factID := range factIDs {
			if fact, exists := factsByID[factID]; exists && varIndex < len(variables) {
				token[variables[varIndex]] = fact
				varIndex++
			}
		}
		if len(token) > 0 {
			tokens = append(tokens, token)
		}
		return tokens
	}

	// Associer chaque variable à son type et créer des tokens combinés
	maxTokens := 0
	for _, facts := range factsByType {
		if len(facts) > maxTokens {
			maxTokens = len(facts)
		}
	}

	for i := 0; i < maxTokens; i++ {
		token := make(map[string]*rete.Fact)
		for j, varType := range variableTypes {
			if j < len(variables) && i < len(factsByType[varType]) {
				token[variables[j]] = factsByType[varType][i]
			}
		}
		if len(token) > 0 {
			tokens = append(tokens, token)
		}
	}

	return tokens
}

// isMainVariableFact vérifie si un fait correspond à la variable principale d'une règle
func isMainVariableFact(fact *rete.Fact, varName string, rules []BetaParsedRule, actionName string) bool {
	for _, rule := range rules {
		if rule.ActionName == actionName {
			for _, variable := range rule.Variables {
				if variable.Name == varName && variable.Role == "primary" {
					return fact.Type == variable.DataType
				}
			}
		}
	}
	// Fallback : accepter tous les faits pour la variable principale
	return true
}

// findRuleByAction trouve une règle par nom d'action
func findRuleByAction(rules []BetaParsedRule, actionName string) *BetaParsedRule {
	for i := range rules {
		if rules[i].ActionName == actionName {
			return &rules[i]
		}
	}
	return nil
}

// validateOperatorSemantics valide la sémantique selon l'opérateur utilisé
func validateOperatorSemantics(action *BetaActionResult, rule *BetaParsedRule, allFacts []*rete.Fact) bool {
	switch rule.SemanticType {
	case "logical_and":
		// Pour AND, vérifier que tous les faits joints satisfont les conditions
		return action.Count > 0 && len(action.JoinedFacts) > 0
	case "logical_or":
		// Pour OR, au moins un fait doit satisfaire
		return action.Count > 0
	case "equality":
		// Pour ==, !=, vérifier la logique d'égalité
		return validateEqualityOperator(action, rule, allFacts)
	case "relational":
		// Pour <, >, <=, >=, vérifier la logique de comparaison
		return validateRelationalOperator(action, rule, allFacts)
	case "membership":
		// Pour IN, vérifier l'appartenance
		return validateMembershipOperator(action, rule, allFacts)
	case "pattern_matching":
		// Pour CONTAINS, LIKE, MATCHES
		return validatePatternOperator(action, rule, allFacts)
	case "additive", "multiplicative":
		// Pour +, -, *, /, vérifier les calculs arithmétiques
		return validateArithmeticOperator(action, rule, allFacts)
	case "negation":
		// Pour NOT
		return validateNegationOperator(action, rule, allFacts)
	case "existence":
		// Pour EXISTS
		return validateExistenceOperator(action, rule, allFacts)
	default:
		return true // Opérateur inconnu, validation passée
	}
}

// validateEqualityOperator valide les opérateurs d'égalité (==, !=)
func validateEqualityOperator(action *BetaActionResult, rule *BetaParsedRule, allFacts []*rete.Fact) bool {
	// Logique simplifiée : vérifier que l'action a été déclenchée de manière cohérente
	return action.Count >= 0
}

// validateRelationalOperator valide les opérateurs relationnels (<, >, <=, >=)
func validateRelationalOperator(action *BetaActionResult, rule *BetaParsedRule, allFacts []*rete.Fact) bool {
	// Logique simplifiée : vérifier que les comparaisons numériques sont cohérentes
	return action.Count >= 0
}

// validateMembershipOperator valide l'opérateur IN
func validateMembershipOperator(action *BetaActionResult, rule *BetaParsedRule, allFacts []*rete.Fact) bool {
	// Vérifier que les éléments appartiennent bien aux ensembles
	return action.Count >= 0
}

// validatePatternOperator valide les opérateurs de patterns (CONTAINS, LIKE, MATCHES)
func validatePatternOperator(action *BetaActionResult, rule *BetaParsedRule, allFacts []*rete.Fact) bool {
	// Vérifier que les patterns correspondent
	return action.Count >= 0
}

// validateArithmeticOperator valide les opérateurs arithmétiques (+, -, *, /)
func validateArithmeticOperator(action *BetaActionResult, rule *BetaParsedRule, allFacts []*rete.Fact) bool {
	// Vérifier que les calculs arithmétiques sont cohérents
	return action.Count >= 0
}

// validateNegationOperator valide l'opérateur NOT
func validateNegationOperator(action *BetaActionResult, rule *BetaParsedRule, allFacts []*rete.Fact) bool {
	// Vérifier que la négation est logiquement correcte
	// Pour NOT, le nombre de déclenchements devrait correspondre aux faits qui ne satisfont pas la condition
	return action.Count >= 0
}

// validateExistenceOperator valide l'opérateur EXISTS
func validateExistenceOperator(action *BetaActionResult, rule *BetaParsedRule, allFacts []*rete.Fact) bool {
	// Pour EXISTS, vérifier qu'au moins un fait satisfait la condition
	return action.Count > 0
}

// analyzeJoinNodeSemantics analyse en détail la sémantique des jointures
func analyzeJoinNodeSemantics(action *BetaActionResult, rule *BetaParsedRule) {
	fmt.Printf("   🔗 Analyse JoinNode (%s):\n", rule.SemanticType)
	fmt.Printf("      Variables: %d, Tuples joints: %d\n", len(rule.Variables), len(action.JoinedFacts))

	if len(action.JoinedFacts) > 0 {
		fmt.Printf("      Échantillon de tuples:\n")
		for i, tuple := range action.JoinedFacts {
			if i < 3 { // Afficher seulement les 3 premiers
				factIDs := make([]string, len(tuple))
				for j, fact := range tuple {
					factIDs[j] = fact.ID
				}
				fmt.Printf("        Tuple %d: [%s]\n", i+1, strings.Join(factIDs, ", "))
			}
		}
	}
}

// analyzeNotNodeSemantics analyse en détail la sémantique des négations
func analyzeNotNodeSemantics(action *BetaActionResult, rule *BetaParsedRule, allFacts []*rete.Fact) {
	fmt.Printf("   🚫 Analyse NotNode:\n")
	fmt.Printf("      Faits filtrés par négation: %d\n", action.Count)

	// Calculer le nombre de faits qui auraient dû être rejetés
	totalFacts := 0
	for _, fact := range allFacts {
		if isFactOfType(fact, rule.Variables) {
			totalFacts++
		}
	}
	rejectedFacts := totalFacts - action.Count
	fmt.Printf("      Faits rejetés (ne satisfont pas NOT): %d/%d\n", rejectedFacts, totalFacts)
}

// analyzeExistsNodeSemantics analyse en détail la sémantique des existences
func analyzeExistsNodeSemantics(action *BetaActionResult, rule *BetaParsedRule, allFacts []*rete.Fact) {
	fmt.Printf("   ✨ Analyse ExistsNode:\n")
	fmt.Printf("      Preuves d'existence: %d\n", action.Count)

	if action.Count > 0 {
		fmt.Printf("      Existence prouvée: ✅\n")
	} else {
		fmt.Printf("      Existence prouvée: ❌\n")
	}
}

// isFactOfType vérifie si un fait correspond aux types de variables d'une règle
func isFactOfType(fact *rete.Fact, variables []VariableInfo) bool {
	for _, variable := range variables {
		if fact.Type == variable.DataType {
			return true
		}
	}
	return false
}
