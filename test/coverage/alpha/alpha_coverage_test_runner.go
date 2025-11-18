package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/treivax/tsd/constraint"
	"github.com/treivax/tsd/rete"
)

// AlphaTestResult stocke les résultats d'un test alpha
type AlphaTestResult struct {
	TestName          string
	Description       string
	ConstraintFile    string
	FactsFile         string
	Rules             []ParsedRule
	Facts             []*rete.Fact
	Network           *rete.ReteNetwork
	Actions           []ActionResult
	ValidationResults AlphaValidationResults
	ExecutionTime     time.Duration
	Success           bool
	ErrorMessage      string
}

// AlphaValidationResults contient les résultats de validation sémantique
type AlphaValidationResults struct {
	ExpectedActions []ExpectedActionResult
	ObservedActions []ActionResult
	SemanticScore   float64
	Errors          []string
}

// ExpectedActionResult représente une action attendue
type ExpectedActionResult struct {
	ActionName    string
	ExpectedCount int
	ExpectedFacts []string // IDs des faits qui devraient déclencher l'action
	Description   string
}

// ParsedRule représente une règle parsée
type ParsedRule struct {
	RuleNumber         int
	RuleText           string
	ActionName         string
	Condition          string      // Pour l'affichage
	OriginalConstraint interface{} // Pour l'évaluation
	IsNegation         bool
}

// ActionResult représente le résultat d'une action
type ActionResult struct {
	ActionName string
	Count      int
	Facts      []*rete.Fact
}

// NetworkNode représente un nœud du réseau RETE
type NetworkNode struct {
	ID         string
	Type       string
	Condition  interface{}
	FactsCount int
	Facts      []*rete.Fact
	Children   []string
}

func main() {
	fmt.Println("🔬 EXÉCUTION DES TESTS DE COUVERTURE ALPHA NODES")
	fmt.Println("================================================")

	testDir := "/home/resinsec/dev/tsd/test/coverage/alpha"
	resultsFile := "/home/resinsec/dev/tsd/ALPHA_NODES_COVERAGE_COMPLETE_RESULTS.md"

	// Découvrir tous les tests
	tests, err := discoverAlphaTests(testDir)
	if err != nil {
		fmt.Printf("❌ Erreur découverte tests: %v\n", err)
		return
	}

	fmt.Printf("📊 %d tests Alpha découverts\n\n", len(tests))

	// Exécuter tous les tests
	var allResults []AlphaTestResult
	for _, testName := range tests {
		fmt.Printf("🧪 Exécution test: %s\n", testName)
		result := executeAlphaTest(testDir, testName)
		allResults = append(allResults, result)

		if result.Success {
			fmt.Printf("✅ Succès (%v)\n", result.ExecutionTime)
		} else {
			fmt.Printf("❌ Échec: %s\n", result.ErrorMessage)
		}
		fmt.Println()
	}

	// Générer le rapport complet
	err = generateCompleteReport(allResults, resultsFile)
	if err != nil {
		fmt.Printf("❌ Erreur génération rapport: %v\n", err)
		return
	}

	// Résumé final
	successCount := 0
	for _, result := range allResults {
		if result.Success {
			successCount++
		}
	}

	fmt.Printf("🎯 RÉSUMÉ FINAL\n")
	fmt.Printf("==============\n")
	fmt.Printf("✅ Tests réussis: %d/%d\n", successCount, len(allResults))
	fmt.Printf("📄 Rapport complet: %s\n", resultsFile)
}

// discoverAlphaTests découvre tous les tests dans le répertoire
func discoverAlphaTests(testDir string) ([]string, error) {
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

// executeAlphaTest exécute un test alpha complet
func executeAlphaTest(testDir, testName string) AlphaTestResult {
	startTime := time.Now()
	result := AlphaTestResult{
		TestName:       testName,
		ConstraintFile: filepath.Join(testDir, testName+".constraint"),
		FactsFile:      filepath.Join(testDir, testName+".facts"),
	}

	// Lire la description depuis le fichier constraint
	if description, err := extractDescription(result.ConstraintFile); err == nil {
		result.Description = description
	}

	// Parser les contraintes
	program, err := constraint.ParseFile(result.ConstraintFile)
	if err != nil {
		result.ErrorMessage = fmt.Sprintf("Erreur parsing contraintes: %v", err)
		return result
	}

	// Extraire les règles
	result.Rules = extractRulesFromProgram(program)

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

	// Exécuter le test
	actionsCount := make(map[string]int)
	actionsMap := make(map[string][]*rete.Fact)

	// Soumettre tous les faits au réseau
	for _, fact := range facts {
		err := network.SubmitFact(fact)
		if err != nil {
			result.ErrorMessage = fmt.Sprintf("Erreur soumission fait %s: %v", fact.ID, err)
			return result
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

			// Extraire les faits des tokens
			for _, token := range terminal.Memory.Tokens {
				for _, fact := range token.Facts {
					actionsMap[actionName] = append(actionsMap[actionName], fact)
				}
			}
		}
	}

	// Créer les résultats d'actions
	for actionName, count := range actionsCount {
		result.Actions = append(result.Actions, ActionResult{
			ActionName: actionName,
			Count:      count,
			Facts:      actionsMap[actionName],
		})
	}

	// Validation sémantique
	result.ValidationResults = performAlphaSemanticValidation(result.TestName, result.Actions, result.Facts, result.Rules)

	result.ExecutionTime = time.Since(startTime)
	result.Success = true
	return result
}

// performAlphaSemanticValidation effectue la validation sémantique des résultats Alpha
func performAlphaSemanticValidation(testName string, observedActions []ActionResult, facts []*rete.Fact, rules []ParsedRule) AlphaValidationResults {
	validation := AlphaValidationResults{
		ObservedActions: observedActions,
		Errors:          []string{},
	}

	// Définir les résultats attendus basés sur l'analyse des règles et faits
	expectedActions := generateExpectedActionsAlpha(testName, rules, facts)
	validation.ExpectedActions = expectedActions

	// Comparer attendu vs observé
	score := 0.0
	totalChecks := 0

	for _, expected := range expectedActions {
		totalChecks++
		found := false
		for _, observed := range observedActions {
			if expected.ActionName == observed.ActionName {
				found = true
				if expected.ExpectedCount == observed.Count {
					score += 1.0
				} else {
					validation.Errors = append(validation.Errors,
						fmt.Sprintf("Action %s: attendu %d déclenchements, observé %d",
							expected.ActionName, expected.ExpectedCount, observed.Count))
				}
				break
			}
		}
		if !found {
			validation.Errors = append(validation.Errors,
				fmt.Sprintf("Action attendue manquante: %s", expected.ActionName))
		}
	}

	if totalChecks > 0 {
		validation.SemanticScore = (score / float64(totalChecks)) * 100.0
	}

	return validation
}

// generateExpectedActionsAlpha génère les actions attendues pour les tests Alpha
func generateExpectedActionsAlpha(testName string, rules []ParsedRule, facts []*rete.Fact) []ExpectedActionResult {
	var expected []ExpectedActionResult

	// Logique basée sur le nom du test et l'analyse des règles/faits
	for _, rule := range rules {
		expectedAction := ExpectedActionResult{
			ActionName:    rule.ActionName,
			ExpectedCount: 0,
			ExpectedFacts: []string{},
			Description:   fmt.Sprintf("Action basée sur règle %d", rule.RuleNumber),
		}

		// Analyser les faits qui devraient déclencher cette règle
		for _, fact := range facts {
			shouldTrigger := shouldFactTriggerAlphaRule(rule, fact, testName)
			if shouldTrigger {
				expectedAction.ExpectedCount++
				expectedAction.ExpectedFacts = append(expectedAction.ExpectedFacts, fact.ID)
			}
		}

		expected = append(expected, expectedAction)
	}

	return expected
}

// shouldFactTriggerAlphaRule détermine si un fait devrait déclencher une règle Alpha
func shouldFactTriggerAlphaRule(rule ParsedRule, fact *rete.Fact, testName string) bool {
	// Évaluer la condition réelle de la règle
	// Pour les notConstraints, la négation est déjà intégrée dans l'évaluation
	// donc on ne passe pas rule.IsNegation
	return evaluateAlphaCondition(rule.OriginalConstraint, fact, false)
}

// evaluateAlphaCondition évalue une condition Alpha sur un fait
func evaluateAlphaCondition(condition interface{}, fact *rete.Fact, isNegation bool) bool {
	if condition == nil {
		return !isNegation // Si pas de condition, dépend de la négation
	}

	conditionMet := false

	if condMap, ok := condition.(map[string]interface{}); ok {
		// Identifier le type de condition
		if condType, hasType := condMap["type"]; hasType {
			switch condType {
			case "notConstraint":
				// Pour les contraintes NOT, évaluer l'expression interne
				if expr, hasExpr := condMap["expression"]; hasExpr {
					innerMet := evaluateAlphaCondition(expr, fact, false) // pas de négation sur l'évaluation interne
					conditionMet = !innerMet                              // appliquer la négation du NOT
				}
			case "comparison":
				conditionMet = evaluateConditionOperator(condMap, fact)
			default:
				// Fallback: chercher un opérateur direct
				if _, hasOp := condMap["operator"]; hasOp {
					conditionMet = evaluateConditionOperator(condMap, fact)
				}
			}
		} else if _, hasOp := condMap["operator"]; hasOp {
			conditionMet = evaluateConditionOperator(condMap, fact)
		}
	}

	// Appliquer la négation si nécessaire
	if isNegation {
		return !conditionMet
	}
	return conditionMet
}

// evaluateConditionOperator évalue un opérateur de condition
func evaluateConditionOperator(condMap map[string]interface{}, fact *rete.Fact) bool {
	operator, hasOp := condMap["operator"]
	if !hasOp {
		return false
	}

	left := getConditionValue(condMap["left"], fact)
	right := getConditionValue(condMap["right"], fact)

	switch operator {
	case "==":
		return compareValues(left, right, "==")
	case "!=":
		return compareValues(left, right, "!=")
	case ">":
		return compareValues(left, right, ">")
	case ">=":
		return compareValues(left, right, ">=")
	case "<":
		return compareValues(left, right, "<")
	case "<=":
		return compareValues(left, right, "<=")
	case "CONTAINS":
		return evaluateContains(left, right)
	case "LIKE":
		return evaluateLike(left, right)
	case "MATCHES":
		return evaluateMatches(left, right)
	case "IN":
		return evaluateIn(left, right)
	default:
		return false
	}
}

// getConditionValue extrait une valeur d'une condition
func getConditionValue(value interface{}, fact *rete.Fact) interface{} {
	if valueMap, ok := value.(map[string]interface{}); ok {
		switch valueMap["type"] {
		case "fieldAccess":
			// Accès à un champ du fait
			if field, hasField := valueMap["field"]; hasField {
				if fieldStr, ok := field.(string); ok {
					return fact.Fields[fieldStr]
				}
			}
		case "functionCall":
			// Appel de fonction
			if name, hasName := valueMap["name"]; hasName {
				return evaluateFunctionCall(name, valueMap["args"], fact)
			}
		case "string":
			return valueMap["value"]
		case "number":
			return valueMap["value"]
		case "boolean":
			return valueMap["value"]
		case "arrayLiteral":
			if elements, hasElems := valueMap["elements"]; hasElems {
				return elements
			}
		}
	}
	return value
}

// evaluateFunctionCall évalue un appel de fonction
func evaluateFunctionCall(name interface{}, args interface{}, fact *rete.Fact) interface{} {
	nameStr, ok := name.(string)
	if !ok {
		return nil
	}

	switch nameStr {
	case "ABS":
		if argsList, ok := args.([]interface{}); ok && len(argsList) > 0 {
			value := getConditionValue(argsList[0], fact)
			if num, ok := value.(float64); ok {
				result := num
				if num < 0 {
					result = -num
				}
				return result
			}
			if num, ok := value.(int); ok {
				result := num
				if num < 0 {
					result = -num
				}
				return result
			}
			// Essayer de convertir en float64
			numFloat, isValid := toFloat64(value)
			if isValid {
				result := numFloat
				if numFloat < 0 {
					result = -numFloat
				}
				return result
			}
		}
	case "LENGTH":
		if argsList, ok := args.([]interface{}); ok && len(argsList) > 0 {
			value := getConditionValue(argsList[0], fact)
			if str, ok := value.(string); ok {
				return float64(len(str))
			}
		}
	case "UPPER":
		if argsList, ok := args.([]interface{}); ok && len(argsList) > 0 {
			value := getConditionValue(argsList[0], fact)
			if str, ok := value.(string); ok {
				return strings.ToUpper(str)
			}
		}
	}
	return nil
}

// compareValues compare deux valeurs selon un opérateur
func compareValues(left, right interface{}, operator string) bool {
	// Convertir en float64 si possible pour les comparaisons numériques
	if leftFloat, leftOk := toFloat64(left); leftOk {
		if rightFloat, rightOk := toFloat64(right); rightOk {
			switch operator {
			case "==":
				return leftFloat == rightFloat
			case "!=":
				return leftFloat != rightFloat
			case ">":
				return leftFloat > rightFloat
			case ">=":
				return leftFloat >= rightFloat
			case "<":
				return leftFloat < rightFloat
			case "<=":
				return leftFloat <= rightFloat
			}
		}
	}

	// Comparaison de strings
	leftStr := fmt.Sprintf("%v", left)
	rightStr := fmt.Sprintf("%v", right)

	switch operator {
	case "==":
		return leftStr == rightStr
	case "!=":
		return leftStr != rightStr
	}

	return false
}

// toFloat64 convertit une valeur en float64
func toFloat64(value interface{}) (float64, bool) {
	switch v := value.(type) {
	case float64:
		return v, true
	case int:
		return float64(v), true
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f, true
		}
	}
	return 0, false
}

// evaluateContains évalue l'opérateur CONTAINS
func evaluateContains(left, right interface{}) bool {
	leftStr := fmt.Sprintf("%v", left)
	rightStr := fmt.Sprintf("%v", right)
	return strings.Contains(leftStr, rightStr)
}

// evaluateLike évalue l'opérateur LIKE
func evaluateLike(left, right interface{}) bool {
	leftStr := fmt.Sprintf("%v", left)
	rightStr := fmt.Sprintf("%v", right)

	// Convertir pattern LIKE en regex simple
	pattern := strings.ReplaceAll(rightStr, "%", ".*")
	pattern = "^" + pattern + "$"

	matched, err := regexp.MatchString(pattern, leftStr)
	return err == nil && matched
}

// evaluateMatches évalue l'opérateur MATCHES
func evaluateMatches(left, right interface{}) bool {
	leftStr := fmt.Sprintf("%v", left)
	rightStr := fmt.Sprintf("%v", right)

	matched, err := regexp.MatchString(rightStr, leftStr)
	return err == nil && matched
}

// evaluateIn évalue l'opérateur IN
func evaluateIn(left, right interface{}) bool {
	leftStr := fmt.Sprintf("%v", left)

	if rightList, ok := right.([]interface{}); ok {
		for _, item := range rightList {
			if itemMap, ok := item.(map[string]interface{}); ok {
				if itemValue, hasValue := itemMap["value"]; hasValue {
					if fmt.Sprintf("%v", itemValue) == leftStr {
						return true
					}
				}
			}
		}
	}

	return false
}

// extractDescription extrait la description du fichier constraint
func extractDescription(constraintFile string) (string, error) {
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
	return "", fmt.Errorf("pas de description trouvée")
}

// extractRulesFromProgram extrait les règles parsées
func extractRulesFromProgram(program interface{}) []ParsedRule {
	var rules []ParsedRule

	// Tenter de convertir en map pour accéder aux expressions
	if programMap, ok := program.(map[string]interface{}); ok {
		if expressions, exists := programMap["expressions"]; exists {
			if expList, ok := expressions.([]interface{}); ok {
				for i, expr := range expList {
					if exprMap, ok := expr.(map[string]interface{}); ok {
						rule := ParsedRule{
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

						// Analyser les contraintes pour détecter la négation
						if constraints, hasConstraints := exprMap["constraints"]; hasConstraints {
							rule.Condition, rule.IsNegation = analyzeConstraintStructure(constraints)
							rule.OriginalConstraint = constraints // Stocker la contrainte originale
						}

						// Générer le texte de la règle
						rule.RuleText = generateRuleText(exprMap)
						rules = append(rules, rule)
					}
				}
			}
		}
	}

	return rules
}

// analyzeConstraintStructure analyse la structure des contraintes
func analyzeConstraintStructure(constraints interface{}) (string, bool) {
	if constraintMap, ok := constraints.(map[string]interface{}); ok {
		if constraintType, exists := constraintMap["type"]; exists {
			switch constraintType {
			case "notConstraint":
				if expr, hasExpr := constraintMap["expression"]; hasExpr {
					innerCondition, _ := analyzeConstraintStructure(expr)
					return fmt.Sprintf("NOT(%s)", innerCondition), true
				}
			case "comparison":
				left := extractFieldPath(constraintMap["left"])
				op := fmt.Sprintf("%v", constraintMap["operator"])
				right := extractValue(constraintMap["right"])
				return fmt.Sprintf("%s %s %s", left, op, right), false
			case "logicalExpr":
				left, _ := analyzeConstraintStructure(constraintMap["left"])
				condition := left
				if operations, hasOps := constraintMap["operations"]; hasOps {
					if opList, ok := operations.([]interface{}); ok {
						for _, op := range opList {
							if opMap, ok := op.(map[string]interface{}); ok {
								operator := fmt.Sprintf("%v", opMap["op"])
								right, _ := analyzeConstraintStructure(opMap["right"])
								condition = fmt.Sprintf("%s %s %s", condition, operator, right)
							}
						}
					}
				}
				return condition, false
			}
		}
	}
	return "unknown_condition", false
}

// extractFieldPath extrait le chemin d'accès au champ
func extractFieldPath(fieldData interface{}) string {
	if fieldMap, ok := fieldData.(map[string]interface{}); ok {
		if fieldType, exists := fieldMap["type"]; exists && fieldType == "fieldAccess" {
			object := fmt.Sprintf("%v", fieldMap["object"])
			field := fmt.Sprintf("%v", fieldMap["field"])
			return fmt.Sprintf("%s.%s", object, field)
		}
	}
	return fmt.Sprintf("%v", fieldData)
}

// extractValue extrait une valeur
func extractValue(valueData interface{}) string {
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

// generateRuleText génère le texte de la règle
func generateRuleText(exprMap map[string]interface{}) string {
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

	// Extraire la condition
	condition := "true"
	if constraints, hasConstraints := exprMap["constraints"]; hasConstraints {
		condition, _ = analyzeConstraintStructure(constraints)
	}

	// Extraire l'action
	action := "unknown_action"
	if actionData, hasAction := exprMap["action"]; hasAction {
		if actionMap, ok := actionData.(map[string]interface{}); ok {
			if jobData, hasJob := actionMap["job"]; hasJob {
				if jobMap, ok := jobData.(map[string]interface{}); ok {
					if name, hasName := jobMap["name"]; hasName {
						action = fmt.Sprintf("%v", name)
						// Ajouter les arguments si disponibles
						if args, hasArgs := jobMap["args"]; hasArgs {
							if argList, ok := args.([]interface{}); ok && len(argList) > 0 {
								var argStrings []string
								for _, arg := range argList {
									argStrings = append(argStrings, extractFieldPath(arg))
								}
								action = fmt.Sprintf("%s(%s)", action, strings.Join(argStrings, ", "))
							}
						}
					}
				}
			}
		}
	}

	return fmt.Sprintf("{%s} / %s ==> %s", variables, condition, action)
}

// generateCompleteReport génère le rapport complet
func generateCompleteReport(results []AlphaTestResult, outputFile string) error {
	var report strings.Builder

	// En-tête du rapport
	report.WriteString("# 📊 RAPPORT COMPLET - TESTS DE COUVERTURE ALPHA NODES\n\n")
	report.WriteString(fmt.Sprintf("**Date d'exécution:** %s\n", time.Now().Format("2006-01-02 15:04:05")))
	report.WriteString(fmt.Sprintf("**Nombre de tests:** %d\n\n", len(results)))

	// Résumé exécutif
	successCount := 0
	totalActions := 0
	for _, result := range results {
		if result.Success {
			successCount++
		}
		totalActions += len(result.Actions)
	}

	report.WriteString("## 🎯 RÉSUMÉ EXÉCUTIF\n\n")
	report.WriteString(fmt.Sprintf("- ✅ **Tests réussis:** %d/%d (%.1f%%)\n",
		successCount, len(results), float64(successCount)/float64(len(results))*100))
	report.WriteString(fmt.Sprintf("- 🎬 **Actions déclenchées:** %d\n", totalActions))
	report.WriteString(fmt.Sprintf("- ⚡ **Couverture:** Nœuds Alpha positifs et négatifs\n\n"))

	// Détail de chaque test
	for i, result := range results {
		report.WriteString(fmt.Sprintf("## 🧪 TEST %d: %s\n\n", i+1, result.TestName))

		// Informations générales
		report.WriteString("### 📋 Informations générales\n\n")
		report.WriteString(fmt.Sprintf("- **Description:** %s\n", result.Description))
		report.WriteString(fmt.Sprintf("- **Fichier contraintes:** `%s`\n", result.ConstraintFile))
		report.WriteString(fmt.Sprintf("- **Fichier faits:** `%s`\n", result.FactsFile))
		report.WriteString(fmt.Sprintf("- **Temps d'exécution:** %v\n", result.ExecutionTime))
		report.WriteString(fmt.Sprintf("- **Statut:** %s\n\n", getStatusEmoji(result.Success)))

		if !result.Success {
			report.WriteString(fmt.Sprintf("**❌ Erreur:** %s\n\n", result.ErrorMessage))
			continue
		}

		// Règles du test
		report.WriteString("### 📏 Règles du test\n\n")

		for _, rule := range result.Rules {
			negationIcon := ""
			if rule.IsNegation {
				negationIcon = " 🚫"
			}
			// Lire la règle exacte depuis le fichier .constraint
			exactRule := extractExactRuleFromConstraint(result.ConstraintFile)
			report.WriteString(fmt.Sprintf("**Règle %d%s:**\n", rule.RuleNumber, negationIcon))
			report.WriteString("```constraint\n")
			report.WriteString(exactRule)
			report.WriteString("\n```\n")
			report.WriteString(fmt.Sprintf("- **Action:** `%s`\n", rule.ActionName))
			report.WriteString(fmt.Sprintf("- **Condition:** `%s`\n", rule.Condition))
			report.WriteString(fmt.Sprintf("- **Type:** %s\n\n", getConditionType(rule.IsNegation)))
		}

		// Faits du test
		report.WriteString("### 📦 Faits du test\n\n")
		report.WriteString(fmt.Sprintf("**Nombre total:** %d faits\n\n", len(result.Facts)))
		for j, fact := range result.Facts {
			report.WriteString(fmt.Sprintf("**Fait %d:** `%s`\n", j+1, fact.ID))
			report.WriteString("```json\n")
			report.WriteString(fmt.Sprintf("Type: %s\n", fact.Type))
			report.WriteString("Champs:\n")
			for field, value := range fact.Fields {
				report.WriteString(fmt.Sprintf("  %s: %v\n", field, value))
			}
			report.WriteString("```\n\n")
		}

		// Structure du réseau RETE
		report.WriteString("### 🕸️ Structure du réseau RETE\n\n")
		generateNetworkVisualization(&report, result.Network)

		// Résultats d'exécution
		report.WriteString("### ⚡ Résultats d'exécution\n\n")
		if len(result.Actions) == 0 {
			report.WriteString("**Aucune action déclenchée**\n\n")
		} else {
			report.WriteString(fmt.Sprintf("**%d actions déclenchées:**\n\n", len(result.Actions)))
			for _, action := range result.Actions {
				report.WriteString(fmt.Sprintf("#### 🎯 Action: `%s`\n", action.ActionName))
				report.WriteString(fmt.Sprintf("- **Nombre de déclenchements:** %d\n", action.Count))
				report.WriteString("- **Faits concernés:**\n")
				for k, fact := range action.Facts {
					report.WriteString(fmt.Sprintf("  %d. `%s` (Type: %s)\n", k+1, fact.ID, fact.Type))
				}
				report.WriteString("\n")
			}
		}

		// Validation sémantique
		report.WriteString("### 🧠 Validation sémantique\n\n")
		validation := result.ValidationResults
		report.WriteString(fmt.Sprintf("- **Score de validation:** %.1f%%\n", validation.SemanticScore))

		if len(validation.Errors) > 0 {
			report.WriteString("- **Erreurs détectées:**\n")
			for _, err := range validation.Errors {
				report.WriteString(fmt.Sprintf("  ❌ %s\n", err))
			}
		} else {
			report.WriteString("- **✅ Validation parfaite**\n")
		}
		report.WriteString("\n")

		// Comparaison Attendu vs Observé
		report.WriteString("### 📊 Comparaison Attendu vs Observé\n\n")
		report.WriteString("| Action | Attendu | Observé | Faits Attendus | Faits Observés | Statut |\n")
		report.WriteString("|--------|---------|---------|----------------|----------------|--------|\n")

		for _, expected := range validation.ExpectedActions {
			observed := findObservedAction(expected.ActionName, validation.ObservedActions)
			observedCount := 0
			observedFactIDs := []string{}

			if observed != nil {
				observedCount = observed.Count
				for _, fact := range observed.Facts {
					observedFactIDs = append(observedFactIDs, fact.ID)
				}
			}

			status := "✅"
			if expected.ExpectedCount != observedCount {
				status = "❌"
			}

			report.WriteString(fmt.Sprintf("| %s | %d | %d | %s | %s | %s |\n",
				expected.ActionName,
				expected.ExpectedCount,
				observedCount,
				strings.Join(expected.ExpectedFacts, ", "),
				strings.Join(observedFactIDs, ", "),
				status))
		}
		report.WriteString("\n")

		// Détails des tuples/nœuds attendus
		report.WriteString("#### 📋 Détails des tuples attendus\n\n")
		for _, expected := range validation.ExpectedActions {
			report.WriteString(fmt.Sprintf("**Action `%s`:**\n", expected.ActionName))
			report.WriteString(fmt.Sprintf("- **Description:** %s\n", expected.Description))
			report.WriteString(fmt.Sprintf("- **Faits attendus:** %d\n", expected.ExpectedCount))
			if len(expected.ExpectedFacts) > 0 {
				report.WriteString("- **IDs attendus:**\n")
				for i, factID := range expected.ExpectedFacts {
					report.WriteString(fmt.Sprintf("  %d. `%s`\n", i+1, factID))
				}
			}
			report.WriteString("\n")
		}
		report.WriteString("\n")

		report.WriteString("---\n\n")
	}

	// Écrire le fichier
	return os.WriteFile(outputFile, []byte(report.String()), 0644)
}

// generateNetworkVisualization génère une visualisation du réseau RETE
func generateNetworkVisualization(report *strings.Builder, network *rete.ReteNetwork) {
	report.WriteString("```\n")
	report.WriteString("RÉSEAU RETE - STRUCTURE HIÉRARCHIQUE\n")
	report.WriteString("=====================================\n\n")

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

// Helper functions
func getStatusEmoji(success bool) string {
	if success {
		return "✅ Succès"
	}
	return "❌ Échec"
}

func getConditionType(isNegation bool) string {
	if isNegation {
		return "Condition négative (NOT)"
	}
	return "Condition positive"
}

// findObservedAction trouve une action observée par nom
func findObservedAction(actionName string, observedActions []ActionResult) *ActionResult {
	for i := range observedActions {
		if observedActions[i].ActionName == actionName {
			return &observedActions[i]
		}
	}
	return nil
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

// readExactRuleFromConstraint lit la règle exacte depuis le fichier .constraint
func readExactRuleFromConstraint(constraintFile string) string {
	content, err := os.ReadFile(constraintFile)
	if err != nil {
		return fmt.Sprintf("Erreur lecture fichier: %v", err)
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "==>") {
			return line
		}
	}
	return "Règle non trouvée"
}
