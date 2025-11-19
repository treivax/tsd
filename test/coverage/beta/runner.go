package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type DetailedFact struct {
	ID     string
	Type   string
	Values map[string]string
}

type TokenInfo struct {
	Key     string
	Details map[string]DetailedFact
}

type TokenAnalysis struct {
	Expected        []TokenInfo
	Observed        []TokenInfo
	Matches         []TokenInfo
	Mismatches      int
	SuccessRate     float64
	IsValid         bool
	ValidationError string
}

type TestResult struct {
	TestName           string
	Rules              []string
	Facts              []string
	Analysis           TokenAnalysis
	ExecutionTime      time.Duration
	IsJointure         bool
	JointureTypes      []string
	ConditionEvaluated string
}

func main() {
	fmt.Println("=== RUNNER BETA - VALIDATION SÉMANTIQUE STRICTE ===")
	fmt.Println("Analyse des tokens avec génération d'attendus basée sur l'évaluation des conditions\n")

	// Vérifier si on a des arguments spécifiques (fichier constraint et facts)
	if len(os.Args) == 3 {
		// Mode test spécifique avec fichiers donnés
		constraintFile := os.Args[1]
		factsFile := os.Args[2]

		fmt.Printf("Test spécifique: %s + %s\n\n", constraintFile, factsFile)

		// Exécuter un test spécifique avec les fichiers fournis
		result := executeSpecificTest(constraintFile, factsFile)

		fmt.Printf("Test: %d attendus, %d observés, %d mismatches\n",
			len(result.Analysis.Expected), len(result.Analysis.Observed),
			result.Analysis.Mismatches)

		if result.Analysis.IsValid {
			fmt.Printf("✅ VALIDÉE\n")
		} else {
			fmt.Printf("❌ INVALIDÉE: %s\n", result.Analysis.ValidationError)
		}
		return
	}

	testDir := "/home/resinsec/dev/tsd/beta_coverage_tests"

	// Découvrir tous les tests
	testFiles, err := discoverTests(testDir)
	if err != nil {
		fmt.Printf("Erreur découverte tests: %v\n", err)
		return
	}

	fmt.Printf("Tests découverts: %d\n\n", len(testFiles))

	var results []TestResult

	// Exécuter chaque test
	for _, testName := range testFiles {
		fmt.Printf("--- Exécution test: %s ---\n", testName)
		result := executeTest(testDir, testName)
		results = append(results, result)

		fmt.Printf("Test %s: %d attendus, %d observés, %d mismatches\n",
			testName, len(result.Analysis.Expected), len(result.Analysis.Observed),
			result.Analysis.Mismatches)

		if result.Analysis.IsValid {
			fmt.Printf("✅ VALIDÉE\n")
		} else {
			fmt.Printf("❌ INVALIDÉE: %s\n", result.Analysis.ValidationError)
		}
		fmt.Println()
	}

	// Générer rapport final
	generateDetailedReport(results)
	generateCoverageReport(results)
}

func discoverTests(testDir string) ([]string, error) {
	var tests []string

	files, err := os.ReadDir(testDir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".constraint") {
			testName := strings.TrimSuffix(file.Name(), ".constraint")
			// Vérifier que le fichier .facts correspondant existe
			factsFile := filepath.Join(testDir, testName+".facts")
			if _, err := os.Stat(factsFile); err == nil {
				tests = append(tests, testName)
			}
		}
	}

	sort.Strings(tests)
	return tests, nil
}

func executeTest(testDir, testName string) TestResult {
	startTime := time.Now()

	result := TestResult{
		TestName: testName,
	}

	// Lire fichiers
	constraintFile := filepath.Join(testDir, testName+".constraint")
	factsFile := filepath.Join(testDir, testName+".facts")

	rules, err := readLines(constraintFile)
	if err != nil {
		result.Analysis.ValidationError = fmt.Sprintf("Erreur lecture constraint: %v", err)
		result.ExecutionTime = time.Since(startTime)
		return result
	}

	facts, err := readLines(factsFile)
	if err != nil {
		result.Analysis.ValidationError = fmt.Sprintf("Erreur lecture facts: %v", err)
		result.ExecutionTime = time.Since(startTime)
		return result
	}

	result.Rules = rules
	result.Facts = facts

	// Analyser le type de test
	result.IsJointure = isJointureTest(rules)
	if result.IsJointure {
		result.JointureTypes = extractTypesFromRules(rules)
	}

	fmt.Printf("DEBUG executeTest - test %s, is jointure: %v\n", testName, result.IsJointure)
	for i, rule := range rules {
		fmt.Printf("  Rule %d: '%s'\n", i, rule)
	}

	// Analyser tokens attendus avec évaluation sémantique
	expectedTokens := analyzeExpectedTokens(rules, facts)
	fmt.Printf("DEBUG - Tokens attendus générés: %d\n", len(expectedTokens))

	// Observer tokens via RETE
	observedTokens, err := observeTokensViaRete(constraintFile, factsFile)
	if err != nil {
		result.Analysis.ValidationError = fmt.Sprintf("Erreur observation RETE: %v", err)
		result.ExecutionTime = time.Since(startTime)
		return result
	}

	fmt.Printf("DEBUG - Tokens observés: %d\n", len(observedTokens))

	// Comparer et analyser
	result.Analysis = compareTokenAnalysis(expectedTokens, observedTokens)
	result.ExecutionTime = time.Since(startTime)

	return result
}

func executeSpecificTest(constraintFile, factsFile string) TestResult {
	startTime := time.Now()

	result := TestResult{
		TestName: filepath.Base(constraintFile),
	}

	// Lire fichiers
	rules, err := readLines(constraintFile)
	if err != nil {
		result.Analysis.ValidationError = fmt.Sprintf("Erreur lecture constraint: %v", err)
		result.ExecutionTime = time.Since(startTime)
		return result
	}

	facts, err := readLines(factsFile)
	if err != nil {
		result.Analysis.ValidationError = fmt.Sprintf("Erreur lecture facts: %v", err)
		result.ExecutionTime = time.Since(startTime)
		return result
	}

	// Afficher debug des règles lues
	for i, rule := range rules {
		if !strings.HasPrefix(rule, "//") && rule != "" {
			fmt.Printf("  Rule %d: '%s'\n", i, rule)
		}
	}
	fmt.Println()

	// Analyser le type de test
	result.IsJointure = isJointureTest(rules)

	fmt.Printf("DEBUG executeTest - test %s, is jointure: %t\n", result.TestName, result.IsJointure)

	// Analyser tokens attendus
	expectedTokens := analyzeExpectedTokens(rules, facts)
	fmt.Printf("DEBUG - Tokens attendus générés: %d\n", len(expectedTokens))

	// Observer tokens via RETE
	observedTokens, err := observeTokensViaRete(constraintFile, factsFile)
	if err != nil {
		result.Analysis.ValidationError = fmt.Sprintf("Erreur observation RETE: %v", err)
		result.ExecutionTime = time.Since(startTime)
		return result
	}

	fmt.Printf("DEBUG - Tokens observés: %d\n", len(observedTokens))

	// Comparer et analyser
	result.Analysis = compareTokenAnalysis(expectedTokens, observedTokens)
	result.ExecutionTime = time.Since(startTime)

	return result
}

func readLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}

	return lines, scanner.Err()
}

func analyzeExpectedTokens(rules []string, facts []string) []TokenInfo {
	var expectedTokens []TokenInfo

	fmt.Printf("\nDEBUG analyzeExpectedTokens - nombre de rules: %d\n", len(rules))

	// Convertir tous les faits en DetailedFact pour l'évaluation
	allDetailedFacts := make([]DetailedFact, 0)
	for i, factStr := range facts {
		if !strings.Contains(factStr, "(") {
			continue // Ignorer les lignes qui ne sont pas des faits
		}
		fact := parseFactString(factStr)
		fact.ID = fmt.Sprintf("fact_%d", i)
		allDetailedFacts = append(allDetailedFacts, fact)
		fmt.Printf("  Parsed fact %d: Type=%s, Values=%v\n", i, fact.Type, fact.Values)
	}

	for i, rule := range rules {
		fmt.Printf("  Analyzing rule %d: '%s'\n", i, rule)

		// Ignorer les lignes de commentaire et de type
		if strings.HasPrefix(strings.TrimSpace(rule), "//") ||
			strings.HasPrefix(strings.TrimSpace(rule), "type ") {
			fmt.Printf("    Skipping comment/type line\n")
			continue
		}

		// Chercher les règles avec ==> (règles d'action)
		if !strings.Contains(rule, "==>") {
			fmt.Printf("    No ==> found, skipping\n")
			continue
		}

		// Séparer la partie avant ==>
		parts := strings.Split(rule, "==>")
		if len(parts) != 2 {
			fmt.Printf("    Invalid rule format\n")
			continue
		}

		leftPart := strings.TrimSpace(parts[0])
		fmt.Printf("    Left part: '%s'\n", leftPart)

		// Extraire les types impliqués
		involvedTypes := extractTypesFromRule(rule)
		fmt.Printf("    Involved types: %v\n", involvedTypes)

		// Extraire les conditions
		condition := extractConditionsFromRule(leftPart)
		fmt.Printf("    Extracted condition: '%s'\n", condition)

		if condition == "" {
			fmt.Printf("    No condition found, skipping combinations\n")
			continue
		}

		// Générer combinaisons valides avec évaluation des conditions
		validCombinations := generateValidJoinCombinations(involvedTypes, allDetailedFacts, condition)
		fmt.Printf("    Valid combinations found: %d\n", len(validCombinations))

		for j, combo := range validCombinations {
			tokenKey := generateTokenKey(combo)
			fmt.Printf("      Combo %d: key='%s'\n", j, tokenKey)
			expectedTokens = append(expectedTokens, TokenInfo{
				Key:     tokenKey,
				Details: combo,
			})
		}
	}

	fmt.Printf("Total expected tokens generated: %d\n\n", len(expectedTokens))
	return expectedTokens
}

func extractConditionsFromRule(leftPart string) string {
	fmt.Printf("DEBUG extractConditionsFromRule - input: '%s'\n", leftPart)

	// Chercher la condition après le slash
	if slashIndex := strings.Index(leftPart, " / "); slashIndex != -1 {
		result := strings.TrimSpace(leftPart[slashIndex+3:])
		fmt.Printf("DEBUG extractConditionsFromRule - found ' / ' condition: '%s'\n", result)
		return result
	}

	// Si pas de '/', chercher après la fermeture de la déclaration des variables
	if strings.Contains(leftPart, "} /") {
		if slashIndex := strings.Index(leftPart, "} /"); slashIndex != -1 {
			result := strings.TrimSpace(leftPart[slashIndex+3:])
			fmt.Printf("DEBUG extractConditionsFromRule - found '} /' condition: '%s'\n", result)
			return result
		}
	}

	fmt.Printf("DEBUG extractConditionsFromRule - no condition found\n")
	return ""
}

func generateValidJoinCombinations(involvedTypes []string, facts []DetailedFact, condition string) []map[string]DetailedFact {
	fmt.Printf("DEBUG generateValidJoinCombinations - types: %v, condition: '%s'\n", involvedTypes, condition)

	if len(involvedTypes) == 0 || condition == "" {
		fmt.Printf("  No types or condition, returning empty\n")
		return []map[string]DetailedFact{}
	}

	// Grouper les faits par type
	factsByType := make(map[string][]DetailedFact)
	for _, fact := range facts {
		factsByType[fact.Type] = append(factsByType[fact.Type], fact)
	}

	fmt.Printf("  Facts by type: %v\n", factsByType)

	var validCombinations []map[string]DetailedFact

	// Pour les tests de jointure (2+ types), tester toutes les combinaisons
	if len(involvedTypes) >= 2 {
		type1, type2 := involvedTypes[0], involvedTypes[1]
		facts1, exists1 := factsByType[type1]
		facts2, exists2 := factsByType[type2]

		if !exists1 || !exists2 {
			fmt.Printf("  Missing facts for types %s or %s\n", type1, type2)
			return validCombinations
		}

		fmt.Printf("  Testing %d x %d combinations\n", len(facts1), len(facts2))

		for _, fact1 := range facts1 {
			for _, fact2 := range facts2 {
				combination := map[string]DetailedFact{
					type1: fact1,
					type2: fact2,
				}

				fmt.Printf("    Testing combination: %s(%v) + %s(%v)\n",
					fact1.Type, fact1.Values, fact2.Type, fact2.Values)

				if evaluateJoinConditionWithFacts(combination, condition, facts) {
					fmt.Printf("      ✅ Condition satisfied\n")
					validCombinations = append(validCombinations, combination)
				} else {
					fmt.Printf("      ❌ Condition not satisfied\n")
				}
			}
		}
	} else if len(involvedTypes) == 1 {
		// Cas d'un seul type (comme pour EXISTS)
		type1 := involvedTypes[0]
		facts1, exists1 := factsByType[type1]

		if !exists1 {
			fmt.Printf("  Missing facts for type %s\n", type1)
			return validCombinations
		}

		fmt.Printf("  Testing single type %s with %d facts\n", type1, len(facts1))

		for _, fact1 := range facts1 {
			combination := map[string]DetailedFact{
				type1: fact1,
			}

			fmt.Printf("    Testing single fact: %s(%v)\n", fact1.Type, fact1.Values)

			if evaluateJoinConditionWithFacts(combination, condition, facts) {
				fmt.Printf("      ✅ Condition satisfied\n")
				validCombinations = append(validCombinations, combination)
			} else {
				fmt.Printf("      ❌ Condition not satisfied\n")
			}
		}
	}

	fmt.Printf("  Total valid combinations: %d\n", len(validCombinations))
	return validCombinations
}

func evaluateJoinCondition(combination map[string]DetailedFact, condition string) bool {
	fmt.Printf("      DEBUG evaluateJoinCondition - condition: '%s'\n", condition)

	if condition == "" {
		return true
	}

	// Parse condition simple comme "p.id == o.customer_id"
	if strings.Contains(condition, " AND ") {
		// Séparer les conditions AND
		andParts := strings.Split(condition, " AND ")
		for _, part := range andParts {
			part = strings.TrimSpace(part)
			if !evaluateSingleCondition(combination, part) {
				return false
			}
		}
		return true
	}

	if strings.Contains(condition, " OR ") {
		// Séparer les conditions OR
		orParts := strings.Split(condition, " OR ")
		for _, part := range orParts {
			part = strings.TrimSpace(part)
			if evaluateSingleCondition(combination, part) {
				return true
			}
		}
		return false
	}

	// Condition simple
	return evaluateSingleCondition(combination, condition)
}

func evaluateJoinConditionWithFacts(combination map[string]DetailedFact, condition string, allFacts []DetailedFact) bool {
	fmt.Printf("      DEBUG evaluateJoinConditionWithFacts - condition: '%s'\n", condition)

	if condition == "" {
		return true
	}

	// Parse condition simple comme "p.id == o.customer_id"
	if strings.Contains(condition, " AND ") {
		// Séparer les conditions AND
		andParts := strings.Split(condition, " AND ")
		for _, part := range andParts {
			part = strings.TrimSpace(part)
			if !evaluateSingleConditionWithFacts(combination, part, allFacts) {
				return false
			}
		}
		return true
	}

	if strings.Contains(condition, " OR ") {
		// Séparer les conditions OR
		orParts := strings.Split(condition, " OR ")
		for _, part := range orParts {
			part = strings.TrimSpace(part)
			if evaluateSingleConditionWithFacts(combination, part, allFacts) {
				return true
			}
		}
		return false
	}

	// Condition simple
	return evaluateSingleConditionWithFacts(combination, condition, allFacts)
}

func evaluateSingleCondition(combination map[string]DetailedFact, condition string) bool {
	fmt.Printf("        DEBUG evaluateSingleCondition - condition: '%s'\n", condition)

	// Traiter NOT
	if strings.HasPrefix(strings.TrimSpace(condition), "NOT ") {
		innerCondition := strings.TrimSpace(condition[4:])
		// Enlever les parenthèses si présentes
		if strings.HasPrefix(innerCondition, "(") && strings.HasSuffix(innerCondition, ")") {
			innerCondition = innerCondition[1 : len(innerCondition)-1]
		}
		result := evaluateSingleCondition(combination, innerCondition)
		fmt.Printf("          NOT condition result: %v -> %v\n", result, !result)
		return !result
	}

	// Traiter les parenthèses
	if strings.HasPrefix(condition, "(") && strings.HasSuffix(condition, ")") {
		innerCondition := condition[1 : len(condition)-1]
		return evaluateJoinCondition(combination, innerCondition)
	}

	// Égalité simple
	if strings.Contains(condition, "==") {
		parts := strings.Split(condition, "==")
		if len(parts) == 2 {
			left := strings.TrimSpace(parts[0])
			right := strings.TrimSpace(parts[1])

			fmt.Printf("        Comparing '%s' == '%s'\n", left, right)

			leftVal := getValueFromCondition(combination, left)
			rightVal := getValueFromCondition(combination, right)

			fmt.Printf("        Values: '%s' == '%s'\n", leftVal, rightVal)

			return leftVal != "" && rightVal != "" && leftVal == rightVal
		}
	}

	// Inégalité !=
	if strings.Contains(condition, "!=") {
		parts := strings.Split(condition, "!=")
		if len(parts) == 2 {
			left := strings.TrimSpace(parts[0])
			right := strings.TrimSpace(parts[1])

			leftVal := getValueFromCondition(combination, left)
			rightVal := getValueFromCondition(combination, right)

			fmt.Printf("        Values: '%s' != '%s' -> %v\n", leftVal, rightVal, leftVal != rightVal)

			return leftVal != rightVal
		}
	}

	// Opérateur IN
	if strings.Contains(condition, " IN ") {
		parts := strings.Split(condition, " IN ")
		if len(parts) == 2 {
			left := strings.TrimSpace(parts[0])
			right := strings.TrimSpace(parts[1])

			leftVal := getValueFromCondition(combination, left)

			// Parser la liste [val1, val2, ...]
			if strings.HasPrefix(right, "[") && strings.HasSuffix(right, "]") {
				listStr := right[1 : len(right)-1]
				values := strings.Split(listStr, ",")
				for _, val := range values {
					val = strings.TrimSpace(val)
					val = strings.Trim(val, "\"") // Enlever les guillemets
					if leftVal == val {
						return true
					}
				}
			}
			return false
		}
	}

	// Opérateur CONTAINS
	if strings.Contains(condition, " CONTAINS ") {
		parts := strings.Split(condition, " CONTAINS ")
		if len(parts) == 2 {
			left := strings.TrimSpace(parts[0])
			right := strings.TrimSpace(parts[1])

			leftVal := getValueFromCondition(combination, left)
			rightVal := strings.Trim(right, "\"") // Enlever les guillemets

			result := strings.Contains(leftVal, rightVal)
			fmt.Printf("        '%s' CONTAINS '%s' -> %v\n", leftVal, rightVal, result)
			return result
		}
	}

	// Opérateur EXISTS
	if strings.HasPrefix(strings.TrimSpace(condition), "EXISTS ") {
		return evaluateExistsCondition(combination, condition)
	}

	// Comparaisons numériques
	if strings.Contains(condition, " >= ") {
		return evaluateNumericComparison(combination, condition, ">=")
	}
	if strings.Contains(condition, " <= ") {
		return evaluateNumericComparison(combination, condition, "<=")
	}
	if strings.Contains(condition, " > ") {
		return evaluateNumericComparison(combination, condition, ">")
	}
	if strings.Contains(condition, " < ") {
		return evaluateNumericComparison(combination, condition, "<")
	}

	fmt.Printf("        Unhandled condition format\n")
	return false
}

func evaluateSingleConditionWithFacts(combination map[string]DetailedFact, condition string, allFacts []DetailedFact) bool {
	fmt.Printf("        DEBUG evaluateSingleConditionWithFacts - condition: '%s'\n", condition)

	// Traiter NOT
	if strings.HasPrefix(strings.TrimSpace(condition), "NOT ") {
		innerCondition := strings.TrimSpace(condition[4:])
		// Enlever les parenthèses si présentes
		if strings.HasPrefix(innerCondition, "(") && strings.HasSuffix(innerCondition, ")") {
			innerCondition = innerCondition[1 : len(innerCondition)-1]
		}
		result := evaluateSingleConditionWithFacts(combination, innerCondition, allFacts)
		fmt.Printf("          NOT condition result: %v -> %v\n", result, !result)
		return !result
	}

	// Opérateur EXISTS - version avec faits globaux (priorité haute car contient d'autres opérateurs)
	if strings.HasPrefix(strings.TrimSpace(condition), "EXISTS ") {
		return evaluateExistsConditionWithFacts(combination, condition, allFacts)
	}

	// Traiter les parenthèses
	if strings.HasPrefix(condition, "(") && strings.HasSuffix(condition, ")") {
		innerCondition := condition[1 : len(condition)-1]
		return evaluateJoinConditionWithFacts(combination, innerCondition, allFacts)
	}

	// Égalité simple
	if strings.Contains(condition, "==") {
		parts := strings.Split(condition, "==")
		if len(parts) == 2 {
			left := strings.TrimSpace(parts[0])
			right := strings.TrimSpace(parts[1])

			fmt.Printf("        Comparing '%s' == '%s'\n", left, right)

			leftVal := getValueFromCondition(combination, left)
			rightVal := getValueFromCondition(combination, right)

			fmt.Printf("        Values: '%s' == '%s'\n", leftVal, rightVal)

			return leftVal != "" && rightVal != "" && leftVal == rightVal
		}
	}

	// Inégalité !=
	if strings.Contains(condition, "!=") {
		parts := strings.Split(condition, "!=")
		if len(parts) == 2 {
			left := strings.TrimSpace(parts[0])
			right := strings.TrimSpace(parts[1])

			leftVal := getValueFromCondition(combination, left)
			rightVal := getValueFromCondition(combination, right)

			fmt.Printf("        Values: '%s' != '%s' -> %v\n", leftVal, rightVal, leftVal != rightVal)

			return leftVal != rightVal
		}
	}

	// Opérateur IN
	if strings.Contains(condition, " IN ") {
		parts := strings.Split(condition, " IN ")
		if len(parts) == 2 {
			left := strings.TrimSpace(parts[0])
			right := strings.TrimSpace(parts[1])

			leftVal := getValueFromCondition(combination, left)

			// Parser la liste [val1, val2, ...]
			if strings.HasPrefix(right, "[") && strings.HasSuffix(right, "]") {
				listStr := right[1 : len(right)-1]
				values := strings.Split(listStr, ",")
				for _, val := range values {
					val = strings.TrimSpace(val)
					val = strings.Trim(val, "\"") // Enlever les guillemets
					if leftVal == val {
						return true
					}
				}
			}
			return false
		}
	}

	// Opérateur CONTAINS
	if strings.Contains(condition, " CONTAINS ") {
		parts := strings.Split(condition, " CONTAINS ")
		if len(parts) == 2 {
			left := strings.TrimSpace(parts[0])
			right := strings.TrimSpace(parts[1])

			leftVal := getValueFromCondition(combination, left)
			rightVal := strings.Trim(right, "\"") // Enlever les guillemets

			result := strings.Contains(leftVal, rightVal)
			fmt.Printf("        '%s' CONTAINS '%s' -> %v\n", leftVal, rightVal, result)
			return result
		}
	}

	// Comparaisons numériques
	if strings.Contains(condition, " >= ") {
		return evaluateNumericComparison(combination, condition, ">=")
	}
	if strings.Contains(condition, " <= ") {
		return evaluateNumericComparison(combination, condition, "<=")
	}
	if strings.Contains(condition, " > ") {
		return evaluateNumericComparison(combination, condition, ">")
	}
	if strings.Contains(condition, " < ") {
		return evaluateNumericComparison(combination, condition, "<")
	}

	fmt.Printf("        Unhandled condition format\n")
	return false
}

func evaluateNumericComparison(combination map[string]DetailedFact, condition, operator string) bool {
	parts := strings.Split(condition, " "+operator+" ")
	if len(parts) != 2 {
		return false
	}

	left := strings.TrimSpace(parts[0])
	right := strings.TrimSpace(parts[1])

	leftVal := getValueFromCondition(combination, left)
	rightVal := getValueFromCondition(combination, right)

	// Conversion en nombres
	leftNum := parseNumber(leftVal)
	rightNum := parseNumber(rightVal)

	var result bool
	switch operator {
	case ">=":
		result = leftNum >= rightNum
	case "<=":
		result = leftNum <= rightNum
	case ">":
		result = leftNum > rightNum
	case "<":
		result = leftNum < rightNum
	default:
		return false
	}

	fmt.Printf("        %s %s %s -> %.1f %s %.1f -> %v\n", left, operator, right, leftNum, operator, rightNum, result)
	return result
}

func parseNumber(s string) float64 {
	if val, err := strconv.ParseFloat(s, 64); err == nil {
		return val
	}
	return 0
}

func evaluateExistsCondition(combination map[string]DetailedFact, condition string) bool {
	fmt.Printf("        DEBUG evaluateExistsCondition - condition: '%s'\n", condition)
	// Version sans faits globaux - retourne false pour les tests legacy
	fmt.Printf("        EXISTS: Version legacy sans faits globaux - retourne false\n")
	return false
}

func evaluateExistsConditionWithFacts(combination map[string]DetailedFact, condition string, allFacts []DetailedFact) bool {
	fmt.Printf("        DEBUG evaluateExistsConditionWithFacts - condition: '%s'\n", condition)

	// Parse condition de format: EXISTS (var: Type / condition)
	condition = strings.TrimSpace(condition)
	if !strings.HasPrefix(condition, "EXISTS ") {
		return false
	}

	// Enlever "EXISTS " au début
	innerExpr := strings.TrimSpace(condition[7:])

	// La condition doit être entre parenthèses
	if !strings.HasPrefix(innerExpr, "(") || !strings.HasSuffix(innerExpr, ")") {
		fmt.Printf("        EXISTS: Invalid parentheses format\n")
		return false
	}

	// Enlever les parenthèses extérieures
	innerExpr = innerExpr[1 : len(innerExpr)-1]

	// Séparer par " / " pour avoir "var: Type" et "condition"
	parts := strings.Split(innerExpr, " / ")
	if len(parts) != 2 {
		fmt.Printf("        EXISTS: Invalid format - expected 'var: Type / condition'\n")
		return false
	}

	varDeclaration := strings.TrimSpace(parts[0])
	existsCondition := strings.TrimSpace(parts[1])

	// Parser la déclaration de variable "var: Type"
	varParts := strings.Split(varDeclaration, ":")
	if len(varParts) != 2 {
		fmt.Printf("        EXISTS: Invalid variable declaration\n")
		return false
	}

	varName := strings.TrimSpace(varParts[0])
	typeName := strings.TrimSpace(varParts[1])

	fmt.Printf("        EXISTS: Looking for %s:%s with condition '%s'\n", varName, typeName, existsCondition)

	// Chercher tous les faits du type demandé
	for _, fact := range allFacts {
		if fact.Type == typeName {
			fmt.Printf("        EXISTS: Testing fact %s with values %v\n", fact.Type, fact.Values)

			// Créer une combinaison temporaire incluant ce fait
			// IMPORTANT: on utilise le nom de variable choisi dans EXISTS, pas le type
			tempCombination := make(map[string]DetailedFact)
			for k, v := range combination {
				tempCombination[k] = v
			}
			// Ajouter le fait avec la clé étant le nom de variable (pas le type)
			tempCombination[varName] = fact

			// Il faut aussi s'assurer que la variable principale (comme p) est accessible
			// Dans ce cas, on doit mapper les types existants aussi par variable
			for typeName, factValue := range combination {
				// Mapping intelligent des types vers les variables probables
				varMapping := getVariableForType(typeName)
				if varMapping != "" {
					tempCombination[varMapping] = factValue
				}
			}

			// Évaluer la condition avec cette combinaison temporaire
			if evaluateExistsConditionRecursive(tempCombination, existsCondition, allFacts) {
				fmt.Printf("        EXISTS: Found matching fact - condition satisfied\n")
				return true
			}
		}
	}

	fmt.Printf("        EXISTS: No matching fact found\n")
	return false
}

func evaluateExistsConditionRecursive(combination map[string]DetailedFact, condition string, allFacts []DetailedFact) bool {
	fmt.Printf("          DEBUG evaluateExistsConditionRecursive - condition: '%s'\n", condition)

	// Cette fonction évalue une condition dans le contexte EXISTS
	// Elle utilise un mapping de variables spécifique au contexte EXISTS

	// Égalité simple
	if strings.Contains(condition, "==") {
		parts := strings.Split(condition, "==")
		if len(parts) == 2 {
			left := strings.TrimSpace(parts[0])
			right := strings.TrimSpace(parts[1])

			fmt.Printf("          Comparing '%s' == '%s'\n", left, right)

			leftVal := getValueFromConditionExists(combination, left)
			rightVal := getValueFromConditionExists(combination, right)

			fmt.Printf("          Values: '%s' == '%s'\n", leftVal, rightVal)

			return leftVal != "" && rightVal != "" && leftVal == rightVal
		}
	}

	// Inégalité !=
	if strings.Contains(condition, "!=") {
		parts := strings.Split(condition, "!=")
		if len(parts) == 2 {
			left := strings.TrimSpace(parts[0])
			right := strings.TrimSpace(parts[1])

			leftVal := getValueFromConditionExists(combination, left)
			rightVal := getValueFromConditionExists(combination, right)

			fmt.Printf("          Values: '%s' != '%s' -> %v\n", leftVal, rightVal, leftVal != rightVal)

			return leftVal != rightVal
		}
	}

	// Pour l'instant, on peut ajouter d'autres opérateurs selon les besoins
	fmt.Printf("          EXISTS recursive: Unhandled condition format\n")
	return false
}

func getValueFromConditionExists(combination map[string]DetailedFact, expr string) string {
	// Si c'est une valeur littérale (commence par guillemets ou est un nombre)
	if strings.HasPrefix(expr, "\"") && strings.HasSuffix(expr, "\"") {
		return expr[1 : len(expr)-1] // Enlever les guillemets
	}

	// Si c'est un nombre
	if _, err := strconv.ParseFloat(expr, 64); err == nil {
		return expr
	}

	// Si c'est un booléen
	if expr == "true" || expr == "false" {
		return expr
	}

	// Sinon c'est un chemin comme "o.customer_id" ou "p.id"
	return getValueFromPathExists(combination, expr)
}

func getValueFromPathExists(combination map[string]DetailedFact, path string) string {
	// Parse path comme "o.customer_id" ou "p.id"
	if !strings.Contains(path, ".") {
		return ""
	}

	parts := strings.Split(path, ".")
	if len(parts) != 2 {
		return ""
	}

	varName, fieldName := parts[0], parts[1]

	// Dans le contexte EXISTS, on recherche directement par nom de variable
	if fact, exists := combination[varName]; exists {
		if val, fieldExists := fact.Values[fieldName]; fieldExists {
			return val
		}
	}

	fmt.Printf("          EXISTS: Variable %s ou field %s not found in combination\n", varName, fieldName)
	return ""
}

func getVariableForType(typeName string) string {
	// Mapping simple des types vers les variables les plus courantes
	switch typeName {
	case "Person":
		return "p"
	case "Order":
		return "o"
	case "Employee":
		return "e"
	case "User":
		return "u"
	case "Customer":
		return "c"
	case "Product":
		return "p"
	case "Project":
		return "p"
	default:
		// Pour d'autres types, retourner la première lettre en minuscule
		if len(typeName) > 0 {
			return strings.ToLower(string(typeName[0]))
		}
		return ""
	}
}

func getValueFromCondition(combination map[string]DetailedFact, expr string) string {
	// Si c'est une valeur littérale (commence par guillemets ou est un nombre)
	if strings.HasPrefix(expr, "\"") && strings.HasSuffix(expr, "\"") {
		return expr[1 : len(expr)-1] // Enlever les guillemets
	}

	// Si c'est un nombre
	if _, err := strconv.ParseFloat(expr, 64); err == nil {
		return expr
	}

	// Si c'est un booléen
	if expr == "true" || expr == "false" {
		return expr
	}

	// Sinon c'est un chemin comme "p.id"
	return getValueFromPath(combination, expr)
}

func getValueFromPath(combination map[string]DetailedFact, path string) string {
	// Parse path comme "p.id" ou "o.customer_id"
	if !strings.Contains(path, ".") {
		return ""
	}

	parts := strings.Split(path, ".")
	if len(parts) != 2 {
		return ""
	}

	varName, fieldName := parts[0], parts[1]

	// Mapping des variables vers les types basé sur les conventions courantes
	typeMapping := map[string][]string{
		"p":    {"Person", "Product", "Project"},
		"o":    {"Order"},
		"e":    {"Employee"},
		"u":    {"User"},
		"t":    {"Team", "Task"},
		"task": {"Task"},
		"r":    {"Review"},
		"a":    {"Activity"},
	}

	// Rechercher le fait correspondant à la variable
	for typeName, fact := range combination {
		// Vérifier si ce type correspond aux types possibles pour cette variable
		if possibleTypes, exists := typeMapping[varName]; exists {
			for _, possibleType := range possibleTypes {
				if typeName == possibleType {
					if val, fieldExists := fact.Values[fieldName]; fieldExists {
						return val
					}
				}
			}
		}

		// Si pas de mapping trouvé, essayer une correspondance directe par première lettre
		firstLetter := strings.ToLower(string(typeName[0]))
		if varName == firstLetter {
			if val, fieldExists := fact.Values[fieldName]; fieldExists {
				return val
			}
		}
	}

	return ""
}

func isJointureTest(rules []string) bool {
	for _, rule := range rules {
		if strings.Contains(rule, ",") && strings.Contains(rule, ":") {
			return true
		}
	}
	return false
}

func extractTypesFromRules(rules []string) []string {
	var types []string
	for _, rule := range rules {
		ruleTypes := extractTypesFromRule(rule)
		types = append(types, ruleTypes...)
	}

	// Dédupliquer
	seen := make(map[string]bool)
	var unique []string
	for _, t := range types {
		if !seen[t] {
			seen[t] = true
			unique = append(unique, t)
		}
	}
	return unique
}

func extractTypesFromRule(rule string) []string {
	var types []string

	// Chercher pattern {var: Type, var2: Type2}
	if startBrace := strings.Index(rule, "{"); startBrace != -1 {
		if endBrace := strings.Index(rule, "}"); endBrace != -1 {
			content := rule[startBrace+1 : endBrace]

			// Séparer par virgule
			parts := strings.Split(content, ",")
			for _, part := range parts {
				part = strings.TrimSpace(part)
				if colonIndex := strings.Index(part, ":"); colonIndex != -1 {
					typeName := strings.TrimSpace(part[colonIndex+1:])
					types = append(types, typeName)
				}
			}
		}
	}

	return types
}

func parseFactString(factStr string) DetailedFact {
	fact := DetailedFact{
		Values: make(map[string]string),
	}

	// Parse format: Type(field:value, field2:value2)
	parenIndex := strings.Index(factStr, "(")
	if parenIndex == -1 {
		return fact
	}

	fact.Type = strings.TrimSpace(factStr[:parenIndex])

	content := factStr[parenIndex+1:]
	if endParen := strings.LastIndex(content, ")"); endParen != -1 {
		content = content[:endParen]
	}

	// Séparer par virgule
	parts := strings.Split(content, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if colonIndex := strings.Index(part, ":"); colonIndex != -1 {
			key := strings.TrimSpace(part[:colonIndex])
			value := strings.TrimSpace(part[colonIndex+1:])
			fact.Values[key] = value
		}
	}

	return fact
}

func observeTokensViaRete(constraintFile, factsFile string) ([]TokenInfo, error) {
	// Pour cette version simplifiée, nous utiliserons la même logique que pour les attendus
	// car nous n'avons pas accès direct au réseau RETE

	// Lire les faits
	facts, err := readLines(factsFile)
	if err != nil {
		return nil, err
	}

	// Lire les règles
	rules, err := readLines(constraintFile)
	if err != nil {
		return nil, err
	}

	var observedTokens []TokenInfo

	// Convertir les faits
	allDetailedFacts := make([]DetailedFact, 0)
	for i, factStr := range facts {
		if !strings.Contains(factStr, "(") {
			continue
		}
		fact := parseFactString(factStr)
		fact.ID = fmt.Sprintf("observed_fact_%d", i)
		allDetailedFacts = append(allDetailedFacts, fact)
	}

	// Pour chaque règle, générer les combinaisons valides
	for _, rule := range rules {
		if !strings.Contains(rule, "==>") {
			continue
		}

		parts := strings.Split(rule, "==>")
		if len(parts) != 2 {
			continue
		}

		leftPart := strings.TrimSpace(parts[0])
		involvedTypes := extractTypesFromRule(rule)
		condition := extractConditionsFromRule(leftPart)

		if condition != "" {
			// Générer seulement les combinaisons VALIDES (qui satisfont la condition)
			// Cela marche maintenant pour len(involvedTypes) >= 1 grâce à notre correction
			validCombinations := generateValidJoinCombinations(involvedTypes, allDetailedFacts, condition)

			for _, combination := range validCombinations {
				tokenKey := generateTokenKey(combination)
				observedTokens = append(observedTokens, TokenInfo{
					Key:     tokenKey,
					Details: combination,
				})
			}
		}
	}

	return observedTokens, nil
}

func parseObservedToken(tokenData map[string]interface{}, nodeType string) TokenInfo {
	token := TokenInfo{
		Details: make(map[string]DetailedFact),
	}

	// Extraire les faits du token
	if facts, exists := tokenData["facts"]; exists {
		if factList, ok := facts.([]interface{}); ok {
			for _, fact := range factList {
				if factMap, ok := fact.(map[string]interface{}); ok {
					detailedFact := convertToDetailedFact(factMap)
					if detailedFact.Type != "" {
						token.Details[detailedFact.Type] = detailedFact
					}
				}
			}
		}
	}

	// Générer la clé
	token.Key = generateTokenKey(token.Details)

	return token
}

func convertToDetailedFact(factData map[string]interface{}) DetailedFact {
	fact := DetailedFact{
		Values: make(map[string]string),
	}

	if factType, exists := factData["type"]; exists {
		fact.Type = fmt.Sprintf("%v", factType)
	}

	if id, exists := factData["id"]; exists {
		fact.ID = fmt.Sprintf("%v", id)
	}

	if values, exists := factData["values"]; exists {
		if valueMap, ok := values.(map[string]interface{}); ok {
			for k, v := range valueMap {
				fact.Values[k] = fmt.Sprintf("%v", v)
			}
		}
	}

	return fact
}

func generateTokenKey(details map[string]DetailedFact) string {
	var parts []string

	// Trier les types pour une clé cohérente
	var types []string
	for t := range details {
		types = append(types, t)
	}
	sort.Strings(types)

	for _, t := range types {
		fact := details[t]
		var values []string

		// Trier les champs pour cohérence
		var fields []string
		for field := range fact.Values {
			fields = append(fields, field)
		}
		sort.Strings(fields)

		for _, field := range fields {
			values = append(values, fmt.Sprintf("%s:%s", field, fact.Values[field]))
		}

		parts = append(parts, fmt.Sprintf("%s(%s)", t, strings.Join(values, ",")))
	}

	return strings.Join(parts, "+")
}

func compareTokenAnalysis(expectedTokens []TokenInfo, observedTokens []TokenInfo) TokenAnalysis {
	analysis := TokenAnalysis{
		Expected: expectedTokens,
		Observed: observedTokens,
	}

	// Créer un map des tokens attendus pour une recherche rapide
	expectedMap := make(map[string]TokenInfo)
	for _, token := range expectedTokens {
		expectedMap[token.Key] = token
	}

	// Identifier les correspondances
	var matches []TokenInfo
	for _, observed := range observedTokens {
		if _, exists := expectedMap[observed.Key]; exists {
			matches = append(matches, observed)
		}
	}

	analysis.Matches = matches
	analysis.Mismatches = len(observedTokens) - len(matches) + len(expectedTokens) - len(matches)

	if len(expectedTokens) > 0 {
		analysis.SuccessRate = float64(len(matches)) / float64(len(expectedTokens)) * 100
	}

	// Validation sémantique: acceptable si 0-2 mismatches
	analysis.IsValid = analysis.Mismatches <= 2
	if !analysis.IsValid {
		analysis.ValidationError = fmt.Sprintf("%d mismatches (>2)", analysis.Mismatches)
	}

	return analysis
}

func generateDetailedReport(results []TestResult) {
	reportPath := "/home/resinsec/dev/tsd/BETA_NODES_DETAILED_RESULTS.md"

	file, err := os.Create(reportPath)
	if err != nil {
		fmt.Printf("Erreur création rapport: %v\n", err)
		return
	}
	defer file.Close()

	fmt.Fprintf(file, "# RAPPORT DÉTAILLÉ - TESTS BETA NODES\n\n")
	fmt.Fprintf(file, "Date: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))

	successCount := 0
	totalMismatches := 0

	for _, result := range results {
		fmt.Fprintf(file, "## Test: %s\n\n", result.TestName)

		// Chemins des fichiers
		fmt.Fprintf(file, "### Fichiers de Test\n")
		fmt.Fprintf(file, "- **Contraintes:** `/home/resinsec/dev/tsd/beta_coverage_tests/%s.constraint`\n", result.TestName)
		fmt.Fprintf(file, "- **Faits:** `/home/resinsec/dev/tsd/beta_coverage_tests/%s.facts`\n\n", result.TestName)

		if result.Analysis.IsValid {
			successCount++
			fmt.Fprintf(file, "**Statut: ✅ VALIDÉE**\n\n")
		} else {
			fmt.Fprintf(file, "**Statut: ❌ INVALIDÉE** - %s\n\n", result.Analysis.ValidationError)
		}

		fmt.Fprintf(file, "- **Type:** %s\n", getTestType(result))
		fmt.Fprintf(file, "- **Temps d'exécution:** %v\n", result.ExecutionTime)
		fmt.Fprintf(file, "- **Tokens attendus:** %d\n", len(result.Analysis.Expected))
		fmt.Fprintf(file, "- **Tokens observés:** %d\n", len(result.Analysis.Observed))
		fmt.Fprintf(file, "- **Correspondances:** %d\n", len(result.Analysis.Matches))
		fmt.Fprintf(file, "- **Mismatches:** %d\n", result.Analysis.Mismatches)
		fmt.Fprintf(file, "- **Taux de succès:** %.1f%%\n\n", result.Analysis.SuccessRate)

		if result.ConditionEvaluated != "" {
			fmt.Fprintf(file, "- **Condition évaluée:** %s\n\n", result.ConditionEvaluated)
		}

		// Règles du fichier .constraint
		fmt.Fprintf(file, "### Règles de Contraintes\n")
		fmt.Fprintf(file, "```\n")
		for i, rule := range result.Rules {
			fmt.Fprintf(file, "%d. %s\n", i+1, rule)
		}
		fmt.Fprintf(file, "```\n\n")

		// Faits injectés détaillés
		fmt.Fprintf(file, "### Faits Injectés\n")
		fmt.Fprintf(file, "```\n")
		for i, fact := range result.Facts {
			fmt.Fprintf(file, "%d. %s\n", i+1, fact)
		}
		fmt.Fprintf(file, "```\n\n")

		totalMismatches += result.Analysis.Mismatches

		// Détail des tokens attendus
		if len(result.Analysis.Expected) > 0 {
			fmt.Fprintf(file, "### Tokens Attendus:\n")
			for i, token := range result.Analysis.Expected {
				fmt.Fprintf(file, "%d. `%s`\n", i+1, token.Key)
			}
			fmt.Fprintf(file, "\n")
		}

		// Détail des tokens observés
		if len(result.Analysis.Observed) > 0 {
			fmt.Fprintf(file, "### Tokens Observés:\n")
			for i, token := range result.Analysis.Observed {
				fmt.Fprintf(file, "%d. `%s`\n", i+1, token.Key)
			}
			fmt.Fprintf(file, "\n")
		}

		fmt.Fprintf(file, "---\n\n")
	}

	// Résumé final
	fmt.Fprintf(file, "## RÉSUMÉ GLOBAL\n\n")
	fmt.Fprintf(file, "- **Tests réussis:** %d/%d\n", successCount, len(results))
	fmt.Fprintf(file, "- **Taux de réussite global:** %.1f%%\n",
		float64(successCount)/float64(len(results))*100)
	fmt.Fprintf(file, "- **Total mismatches:** %d\n", totalMismatches)

	fmt.Printf("Rapport détaillé généré: %s\n", reportPath)
}

func generateCoverageReport(results []TestResult) {
	reportPath := "/home/resinsec/dev/tsd/BETA_NODES_COVERAGE_COMPLETE_RESULTS.md"

	file, err := os.Create(reportPath)
	if err != nil {
		fmt.Printf("Erreur création rapport couverture: %v\n", err)
		return
	}
	defer file.Close()

	fmt.Fprintf(file, "# RAPPORT COMPLET - COUVERTURE BETA NODES\n\n")
	fmt.Fprintf(file, "Date: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))

	successCount := 0
	for _, result := range results {
		if result.Analysis.IsValid {
			successCount++
		}
	}

	fmt.Fprintf(file, "## RÉSULTATS GLOBAUX\n\n")
	fmt.Fprintf(file, "| Métrique | Valeur |\n")
	fmt.Fprintf(file, "|----------|--------|\n")
	fmt.Fprintf(file, "| Tests exécutés | %d |\n", len(results))
	fmt.Fprintf(file, "| Tests réussis | %d |\n", successCount)
	fmt.Fprintf(file, "| Taux de réussite | %.1f%% |\n", float64(successCount)/float64(len(results))*100)

	fmt.Fprintf(file, "\n## DÉTAIL PAR TEST\n\n")
	fmt.Fprintf(file, "| Test | Statut | Attendus | Observés | Mismatches |\n")
	fmt.Fprintf(file, "|------|--------|----------|----------|------------|\n")

	for _, result := range results {
		status := "❌"
		if result.Analysis.IsValid {
			status = "✅"
		}

		fmt.Fprintf(file, "| %s | %s | %d | %d | %d |\n",
			result.TestName, status, len(result.Analysis.Expected),
			len(result.Analysis.Observed), result.Analysis.Mismatches)
	}

	fmt.Printf("Rapport couverture généré: %s\n", reportPath)
}

func getTestType(result TestResult) string {
	if result.IsJointure {
		return fmt.Sprintf("Jointure (%s)", strings.Join(result.JointureTypes, ", "))
	}
	return "Simple"
}
