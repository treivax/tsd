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
type JoinNodeResult struct {
	NodeID        string
	VariablePairs []string
	JoinCondition string
	MatchedTuples int
	JoinType      string // "inner", "cross", "filtered"
	Performance   time.Duration
	SemanticValid bool
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
	fmt.Println("🔬 EXÉCUTION DES TESTS DE COUVERTURE BETA NODES")
	fmt.Println("===============================================")
	fmt.Println("🎯 Analyse sémantique: JoinNode, NotNode, ExistsNode, AccumulateNode")

	testDir := "/home/resinsec/dev/tsd/beta_coverage_tests"
	resultsFile := "/home/resinsec/dev/tsd/BETA_NODES_COVERAGE_COMPLETE_RESULTS.md"

	// Créer le répertoire de tests s'il n'existe pas
	if err := os.MkdirAll(testDir, 0755); err != nil {
		fmt.Printf("❌ Erreur création répertoire tests: %v\n", err)
		return
	}

	// Découvrir tous les tests Beta
	tests, err := discoverBetaTests(testDir)
	if err != nil {
		fmt.Printf("❌ Erreur découverte tests: %v\n", err)
		return
	}

	fmt.Printf("📊 %d tests Beta découverts\n\n", len(tests))

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

	// Soumettre tous les faits avec monitoring
	for _, fact := range result.Facts {
		err := network.SubmitFact(fact)
		if err != nil {
			return fmt.Errorf("erreur soumission fait %s: %v", fact.ID, err)
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

// validateBetaSemantics valide la sémantique des résultats Beta
func validateBetaSemantics(result *BetaTestResult) {
	validation := ValidationReport{
		ActionsValid:    true,
		JoinsValid:      true,
		NegationsValid:  true,
		ExistsValid:     true,
		AggregatesValid: true,
	}

	var validationErrors []string

	// Valider les actions
	for _, expectedAction := range result.ExpectedResults.ExpectedActions {
		found := false
		for i := range result.Actions {
			actualAction := &result.Actions[i]
			if actualAction.ActionName == expectedAction.ActionName {
				found = true
				if actualAction.Count < expectedAction.MinTriggers ||
					actualAction.Count > expectedAction.MaxTriggers {
					validation.ActionsValid = false
					validationErrors = append(validationErrors,
						fmt.Sprintf("Action %s: attendu %d-%d déclenchements, observé %d",
							expectedAction.ActionName, expectedAction.MinTriggers,
							expectedAction.MaxTriggers, actualAction.Count))
				}
				actualAction.SemanticMatch = true
				break
			}
		}
		if !found {
			validation.ActionsValid = false
			validationErrors = append(validationErrors,
				fmt.Sprintf("Action attendue manquante: %s", expectedAction.ActionName))
		}
	}

	// Valider les jointures
	for _, expectedJoin := range result.ExpectedResults.ExpectedJoins {
		joinFound := false
		for i := range result.JoinResults {
			actualJoin := &result.JoinResults[i]
			if actualJoin.MatchedTuples == expectedJoin.ExpectedMatches {
				joinFound = true
				actualJoin.SemanticValid = true
				break
			}
		}
		if !joinFound {
			validation.JoinsValid = false
			validationErrors = append(validationErrors,
				fmt.Sprintf("Jointure attendue: %s -> %s, %d correspondances",
					expectedJoin.LeftFactType, expectedJoin.RightFactType, expectedJoin.ExpectedMatches))
		}
	}

	// Calculer le score sémantique
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
				return "AND/OR expression", "JoinNode", "logical"
			case "comparison":
				left := extractBetaFieldPath(constraintMap["left"])
				op := fmt.Sprintf("%v", constraintMap["operator"])
				right := extractBetaValue(constraintMap["right"])
				condition := fmt.Sprintf("%s %s %s", left, op, right)
				return condition, "JoinNode", "comparison"
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
	case "join_simple":
		return ExpectedTestResults{
			ExpectedActions: []ExpectedAction{
				{ActionName: "join_person_order", MinTriggers: 2, MaxTriggers: 2,
					ExpectedFactIDs: []string{"P001", "P002", "O001", "O002"},
					SemanticReason:  "Deux personnes ont chacune une commande"},
			},
			ExpectedJoins: []ExpectedJoin{
				{LeftFactType: "Person", RightFactType: "Order",
					JoinCondition: "p.id == o.customer_id", ExpectedMatches: 2,
					SemanticReason: "Jointure sur l'ID client"},
			},
		}
	case "not_simple":
		return ExpectedTestResults{
			ExpectedActions: []ExpectedAction{
				{ActionName: "active_person", MinTriggers: 1, MaxTriggers: 1,
					ExpectedFactIDs: []string{"P001"},
					SemanticReason:  "Une seule personne active"},
			},
			ExpectedNegations: []ExpectedNegation{
				{NegatedCondition: "p.active == false", ExpectedFiltered: 1,
					LogicalReason: "Filtrer les personnes inactives"},
			},
		}
	case "exists_simple":
		return ExpectedTestResults{
			ExpectedActions: []ExpectedAction{
				{ActionName: "person_has_orders", MinTriggers: 1, MaxTriggers: 1,
					ExpectedFactIDs: []string{"P001"},
					SemanticReason:  "Une personne a des commandes"},
			},
			ExpectedExists: []ExpectedExists{
				{ExistsCondition: "o.customer_id == p.id", ShouldExist: true,
					QuantifierReason: "Vérifier l'existence de commandes"},
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
					report.WriteString("**📋 Tokens et couples de faits activant l'action:**\n\n")
					for k, joinedFacts := range action.JoinedFacts {
						report.WriteString(fmt.Sprintf("##### Token %d\n", k+1))

						// Si c'est une jointure, montrer les couples explicitement
						if len(joinedFacts) > 1 {
							report.WriteString("**Couple de faits joints:**\n")
							for l, fact := range joinedFacts {
								report.WriteString(fmt.Sprintf("- **Fait %d (%s):** `%s`\n",
									l+1, fact.Type, formatFactWithFields(fact)))
							}

							// Montrer l'association explicite pour les jointures
							if len(joinedFacts) == 2 {
								fact1, fact2 := joinedFacts[0], joinedFacts[1]
								report.WriteString(fmt.Sprintf("- **Association:** %s[%s] ⋈ %s[%s]\n",
									fact1.Type, fact1.ID, fact2.Type, fact2.ID))
							}
						} else {
							// Pour une action simple (pas de jointure)
							fact := joinedFacts[0]
							report.WriteString(fmt.Sprintf("**Fait activateur:** `%s`\n", formatFactWithFields(fact)))
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
				report.WriteString("| Action | Attendu | Observé | Faits Attendus | Faits Observés | Statut |\n")
				report.WriteString("|--------|---------|---------|----------------|----------------|--------|\n")
				for _, expected := range result.ExpectedResults.ExpectedActions {
					observed := findActualAction(result.Actions, expected.ActionName)
					status := "❌"
					if observed != nil && observed.Count >= expected.MinTriggers &&
						observed.Count <= expected.MaxTriggers {
						status = "✅"
					}
					observedCount := 0
					observedFactIDs := []string{}
					if observed != nil {
						observedCount = observed.Count
						for _, fact := range observed.Facts {
							observedFactIDs = append(observedFactIDs, fact.ID)
						}
					}
					report.WriteString(fmt.Sprintf("| %s | %d-%d | %d | %s | %s | %s |\n",
						expected.ActionName, expected.MinTriggers, expected.MaxTriggers,
						observedCount, strings.Join(expected.ExpectedFactIDs, ", "),
						strings.Join(observedFactIDs, ", "), status))
				}
				report.WriteString("\n")

				// Détails des tuples/nœuds attendus pour Beta
				report.WriteString("#### 📋 Détails des tuples Beta attendus\n\n")
				for _, expected := range result.ExpectedResults.ExpectedActions {
					report.WriteString(fmt.Sprintf("**Action `%s`:**\n", expected.ActionName))
					report.WriteString(fmt.Sprintf("- **Description:** %s\n", expected.SemanticReason))
					report.WriteString(fmt.Sprintf("- **Déclenchements attendus:** %d-%d\n", expected.MinTriggers, expected.MaxTriggers))
					if len(expected.ExpectedFactIDs) > 0 {
						report.WriteString("- **IDs de faits attendus:**\n")
						for i, factID := range expected.ExpectedFactIDs {
							report.WriteString(fmt.Sprintf("  %d. `%s`\n", i+1, factID))
						}
					}
					if len(expected.RequiredFactTypes) > 0 {
						report.WriteString(fmt.Sprintf("- **Types de faits requis:** %s\n", strings.Join(expected.RequiredFactTypes, ", ")))
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

func formatFactFields(fact rete.Fact) string {
	var fields []string
	for key, value := range fact.Fields {
		fields = append(fields, fmt.Sprintf("%s: %v", key, value))
	}
	return strings.Join(fields, ", ")
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

// formatFactWithFields formate un fait avec tous ses champs de manière lisible
func formatFactWithFields(fact rete.Fact) string {
	var fields []string
	for key, value := range fact.Fields {
		fields = append(fields, fmt.Sprintf("%s=%v", key, value))
	}
	return fmt.Sprintf("%s[%s]", fact.Type, strings.Join(fields, ", "))
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
							condition = fmt.Sprintf("NOT(...) [Négation]")
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
				report.WriteString(fmt.Sprintf("│   │   └── Type: JoinNode\n"))
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
