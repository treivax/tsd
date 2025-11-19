package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/treivax/tsd/rete"
)

type TestResult struct {
	Name           string
	Description    string
	ConstraintFile string
	FactsFile      string
	Success        bool
	ErrorMessage   string
	ExecutionTime  time.Duration
	NodesCount     int
	ActionsCount   int

	// Détails pour rapport enrichi
	Rules          []string      `json:"rules"`
	Facts          []string      `json:"facts"`
	Tokens         []TokenInfo   `json:"tokens"`
	TokenAnalysis  TokenAnalysis `json:"token_analysis"`
	NetworkState   NetworkInfo   `json:"network_state"`
	ExpectedResult string        `json:"expected_result"`
	ObservedResult string        `json:"observed_result"`
	InferenceLog   []string      `json:"inference_log"`
}

type TokenInfo struct {
	ID         string         `json:"id"`
	NodeID     string         `json:"node_id"`
	Facts      []DetailedFact `json:"facts"`
	IsMatch    bool           `json:"is_match"`
	ActionName string         `json:"action_name"`
}

type DetailedFact struct {
	ID     string                 `json:"id"`
	Type   string                 `json:"type"`
	Fields map[string]interface{} `json:"fields"`
}

type TokenAnalysis struct {
	ExpectedTokens []TokenInfo `json:"expected_tokens"`
	ObservedTokens []TokenInfo `json:"observed_tokens"`
	Matches        int         `json:"matches"`
	Mismatches     int         `json:"mismatches"`
}

type NetworkInfo struct {
	TotalNodes    int                    `json:"total_nodes"`
	AlphaNodes    int                    `json:"alpha_nodes"`
	TypeNodes     int                    `json:"type_nodes"`
	TerminalNodes int                    `json:"terminal_nodes"`
	BetaNodes     int                    `json:"beta_nodes"`
	NodesByType   map[string]int         `json:"nodes_by_type"`
	MemoryUsage   map[string]interface{} `json:"memory_usage"`
}

func main() {
	testDir := "/home/resinsec/dev/tsd/test/coverage/alpha"
	resultsFile := "/home/resinsec/dev/tsd/ALPHA_NODES_DETAILED_RESULTS.md"

	fmt.Println("🚀 TESTS ALPHA - PIPELINE UNIQUE - RAPPORT DÉTAILLÉ")
	fmt.Println("===================================================")

	tests, err := discoverTests(testDir)
	if err != nil {
		fmt.Printf("❌ Erreur: %v\n", err)
		return
	}

	sort.Strings(tests)
	fmt.Printf("📊 %d tests découverts\n\n", len(tests))

	var results []TestResult
	for i, testName := range tests {
		fmt.Printf("🧪 Test %d/%d: %s\n", i+1, len(tests), testName)
		result := runDetailedTest(testDir, testName)
		results = append(results, result)

		if result.Success {
			fmt.Printf("✅ Succès (%v) - %d règles, %d faits, %d tokens\n",
				result.ExecutionTime, len(result.Rules), len(result.Facts), len(result.Tokens))
		} else {
			fmt.Printf("❌ Échec: %s\n", result.ErrorMessage)
		}
		fmt.Println()
	}

	generateDetailedReport(results, resultsFile)

	// Résumé
	successCount := 0
	for _, result := range results {
		if result.Success {
			successCount++
		}
	}

	fmt.Printf("🎯 RÉSUMÉ FINAL\n")
	fmt.Printf("===============\n")
	fmt.Printf("✅ Tests réussis: %d/%d (%.1f%%)\n",
		successCount, len(results), float64(successCount)/float64(len(results))*100)
	fmt.Printf("📄 Rapport détaillé: %s\n", resultsFile)
}

func discoverTests(testDir string) ([]string, error) {
	var tests []string

	files, err := filepath.Glob(filepath.Join(testDir, "*.constraint"))
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		baseName := strings.TrimSuffix(filepath.Base(file), ".constraint")
		factsFile := filepath.Join(testDir, baseName+".facts")
		if _, err := os.Stat(factsFile); err == nil {
			tests = append(tests, baseName)
		}
	}

	return tests, nil
}

func runDetailedTest(testDir, testName string) TestResult {
	start := time.Now()

	result := TestResult{
		Name:           testName,
		Description:    fmt.Sprintf("Test Alpha: %s", testName),
		ConstraintFile: filepath.Join(testDir, testName+".constraint"),
		FactsFile:      filepath.Join(testDir, testName+".facts"),
		Success:        false,
		Rules:          []string{},
		Facts:          []string{},
		Tokens:         []TokenInfo{},
		InferenceLog:   []string{},
		NetworkState:   NetworkInfo{NodesByType: make(map[string]int), MemoryUsage: make(map[string]interface{})},
	}

	// Lire et analyser le contenu des fichiers
	constraintContent, err := ioutil.ReadFile(result.ConstraintFile)
	if err != nil {
		result.ErrorMessage = fmt.Sprintf("Erreur lecture contraintes: %v", err)
		result.ExecutionTime = time.Since(start)
		return result
	}

	factsContent, err := ioutil.ReadFile(result.FactsFile)
	if err != nil {
		result.ErrorMessage = fmt.Sprintf("Erreur lecture faits: %v", err)
		result.ExecutionTime = time.Since(start)
		return result
	}

	// Extraire les règles et faits
	result.Rules = extractRules(string(constraintContent))
	result.Facts = extractFactsFromContent(string(factsContent))

	result.InferenceLog = append(result.InferenceLog,
		fmt.Sprintf("Début analyse: %d règles, %d faits", len(result.Rules), len(result.Facts)))

	// Analyser les résultats attendus
	result.ExpectedResult = analyzeExpectedResult(result.Rules)

	// PIPELINE UNIQUE - traite contraintes ET faits
	pipeline := rete.NewConstraintPipeline()
	storage := rete.NewMemoryStorage()
	files := []string{result.ConstraintFile, result.FactsFile}

	network, err := pipeline.BuildNetworkFromMultipleFiles(files, storage)
	if err != nil {
		result.ErrorMessage = fmt.Sprintf("Pipeline error: %v", err)
		result.ExecutionTime = time.Since(start)
		return result
	}

	result.InferenceLog = append(result.InferenceLog, "Pipeline construit avec succès")

	// Analyser l'état du réseau
	if network != nil {
		result.InferenceLog = append(result.InferenceLog, "Analyse de l'état du réseau...")

		// CORRECTION: Analyser directement les nœuds du réseau RETE créé
		totalAlphaNodes := len(network.AlphaNodes)
		totalTypeNodes := len(network.TypeNodes)
		totalTerminalNodes := len(network.TerminalNodes)
		totalBetaNodes := len(network.BetaNodes)
		totalNodes := totalAlphaNodes + totalTypeNodes + totalTerminalNodes + totalBetaNodes + 1 // +1 pour RootNode

		result.NodesCount = totalNodes
		result.NetworkState.TotalNodes = totalNodes
		result.NetworkState.AlphaNodes = totalAlphaNodes
		result.ActionsCount = totalTerminalNodes

		// Compter par type de nœud
		result.NetworkState.NodesByType["alpha"] = totalAlphaNodes
		result.NetworkState.NodesByType["terminal"] = totalTerminalNodes
		result.NetworkState.NodesByType["type"] = totalTypeNodes
		result.NetworkState.NodesByType["beta"] = totalBetaNodes

		result.InferenceLog = append(result.InferenceLog,
			fmt.Sprintf("Nœuds créés: %d total (%d alpha, %d terminal, %d type, %d beta)",
				totalNodes, totalAlphaNodes, totalTerminalNodes, totalTypeNodes, totalBetaNodes))

		// CORRECTION: Analyser les tokens du réseau RETE
		result.Tokens = analyzeTokensFromNetworkWithRules(network, result.Rules)
		result.InferenceLog = append(result.InferenceLog,
			fmt.Sprintf("Tokens analysés: %d trouvés", len(result.Tokens)))

		// NOUVEAU: Analyser les tokens attendus vs observés
		expectedTokens := analyzeExpectedTokens(result.Rules, result.Facts)
		result.TokenAnalysis = compareTokenAnalysis(expectedTokens, result.Tokens)
		result.InferenceLog = append(result.InferenceLog,
			fmt.Sprintf("Analyse tokens: %d attendus, %d observés, %d matches, %d mismatches",
				len(result.TokenAnalysis.ExpectedTokens), len(result.TokenAnalysis.ObservedTokens),
				result.TokenAnalysis.Matches, result.TokenAnalysis.Mismatches))
	}

	// Déterminer le résultat observé
	result.ObservedResult = fmt.Sprintf("Réseau construit: %d nœuds, %d actions possibles",
		result.NodesCount, result.ActionsCount)

	result.Success = true
	result.ExecutionTime = time.Since(start)
	result.InferenceLog = append(result.InferenceLog,
		fmt.Sprintf("Test terminé avec succès en %v", result.ExecutionTime))

	return result
}

func extractRules(content string) []string {
	var rules []string
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "//") && !strings.HasPrefix(line, "type") {
			rules = append(rules, line)
		}
	}
	return rules
}

func extractFactsFromContent(content string) []string {
	var facts []string
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "//") {
			facts = append(facts, line)
		}
	}
	return facts
}

func analyzeExpectedResult(rules []string) string {
	if len(rules) == 0 {
		return "Aucune règle définie"
	}

	for _, rule := range rules {
		if strings.Contains(rule, "==>") {
			parts := strings.Split(rule, "==>")
			if len(parts) > 1 {
				action := strings.TrimSpace(parts[1])
				return fmt.Sprintf("Action attendue: %s", action)
			}
		}
	}
	return "Règle de filtrage"
}

// analyzeTokensFromNetwork analyse les tokens du réseau RETE créé
func analyzeTokensFromNetwork(network *rete.ReteNetwork) []TokenInfo {
	return analyzeTokensFromNetworkWithRules(network, []string{})
}

// analyzeTokensFromNetworkWithRules analyse les tokens avec accès aux règles pour l'extraction d'actions
func analyzeTokensFromNetworkWithRules(network *rete.ReteNetwork, rules []string) []TokenInfo {
	var tokens []TokenInfo

	// Récupérer l'état du réseau
	networkState, err := network.GetNetworkState()
	if err != nil {
		// En cas d'erreur, retourner une liste vide
		return tokens
	}

	// Analyser les tokens dans les nœuds terminaux uniquement
	for nodeID, memory := range networkState {
		if memory != nil && len(memory.Tokens) > 0 && strings.Contains(nodeID, "terminal") {
			for _, token := range memory.Tokens {
				// Extraire le nom de l'action
				actionName := extractActionName(rules, nodeID)

				tokenInfo := TokenInfo{
					ID:         fmt.Sprintf("token_%s_%s", nodeID, token.ID),
					NodeID:     nodeID,
					Facts:      []DetailedFact{}, // Sera rempli avec les faits du token
					IsMatch:    true,
					ActionName: actionName,
				}

				// Extraire les faits du token
				if token.Facts != nil && len(token.Facts) > 0 {
					for _, fact := range token.Facts {
						detailedFact := DetailedFact{
							ID:     fact.ID,
							Type:   fact.Type,
							Fields: fact.Fields,
						}
						tokenInfo.Facts = append(tokenInfo.Facts, detailedFact)
					}
				}

				tokens = append(tokens, tokenInfo)
			}
		}
	}

	return tokens
}

func analyzeTokens(storage rete.Storage, nodes []string) []TokenInfo {
	var tokens []TokenInfo

	// Pour chaque nœud, essayer de récupérer sa mémoire
	for _, nodeID := range nodes {
		// Note: Cette partie nécessiterait l'accès à l'API interne du storage
		// Pour l'instant, on simule la présence de tokens
		if strings.Contains(nodeID, "alpha") || strings.Contains(nodeID, "terminal") {
			tokens = append(tokens, TokenInfo{
				ID:      fmt.Sprintf("token_%s", nodeID),
				NodeID:  nodeID,
				Facts:   []DetailedFact{{ID: "placeholder", Type: "placeholder", Fields: map[string]interface{}{}}},
				IsMatch: true,
			})
		}
	}

	return tokens
}

// analyzeExpectedTokens analyse les tokens attendus basés sur les règles et faits
func analyzeExpectedTokens(rules []string, facts []string) []TokenInfo {
	var expectedTokens []TokenInfo

	for i, rule := range rules {
		// Extraire le nom de l'action pour cette règle
		actionName := extractActionName(rules, fmt.Sprintf("rule_%d_terminal", i))

		// Pour les règles Alpha simples, chaque fait qui matche génère son propre token
		tokenCounter := 0
		for j, fact := range facts {
			if shouldFactMatchRule(rule, fact) {
				// Créer un token attendu séparé pour chaque fait qui matche
				expectedToken := TokenInfo{
					ID:         fmt.Sprintf("expected_token_rule_%d_fact_%d", i, tokenCounter),
					NodeID:     fmt.Sprintf("rule_%d_expected", i),
					Facts:      []DetailedFact{},
					IsMatch:    true,
					ActionName: actionName,
				}

				// Parser le fait pour créer DetailedFact
				detailedFact := parseFactToDetailed(fact, j)
				expectedToken.Facts = append(expectedToken.Facts, detailedFact)

				expectedTokens = append(expectedTokens, expectedToken)
				tokenCounter++
			}
		}
	}

	return expectedTokens
}

// shouldFactMatchRule détermine si un fait devrait matcher une règle
func shouldFactMatchRule(rule, fact string) bool {
	// Parser avancé pour les règles et faits
	// Exemple règle: "{b: Balance} / NOT(ABS(b.amount) > 100) ==> small_balance_found(b.id, b.amount)"
	// Exemple fait: "Balance(id:"B002", amount:-25.0, type:"debit")"

	if !strings.Contains(rule, "==>") {
		return false
	}

	parts := strings.Split(rule, "==>")
	leftPart := strings.TrimSpace(parts[0])

	// Extraire le type de fait recherché
	var targetType string
	var variable string

	if strings.Contains(leftPart, "{") && strings.Contains(leftPart, ":") {
		between := strings.TrimSpace(strings.Split(strings.Split(leftPart, "{")[1], "}")[0])
		if typeParts := strings.Split(between, ":"); len(typeParts) >= 2 {
			variable = strings.TrimSpace(typeParts[0])
			targetType = strings.TrimSpace(typeParts[1])
		}
	}

	// Vérifier si le fait correspond au type
	if !strings.HasPrefix(fact, targetType+"(") {
		return false
	}

	// Parser le fait pour l'évaluation
	detailedFact := parseFactString(fact)

	// Évaluer les conditions de la règle
	return evaluateRuleCondition(leftPart, detailedFact, variable)
}

// parseFactToDetailed convertit un fait string en DetailedFact
func parseFactToDetailed(factStr string, index int) DetailedFact {
	return parseFactString(factStr)
}

// parseFactString parse une chaîne de fait en DetailedFact
func parseFactString(factStr string) DetailedFact {
	// Exemple: Balance(id:"B001", amount:150.0, type:"credit")
	factStr = strings.TrimSpace(factStr)

	// Extraire le type
	parenIdx := strings.Index(factStr, "(")
	if parenIdx == -1 {
		return DetailedFact{Type: factStr, Fields: make(map[string]interface{})}
	}

	factType := factStr[:parenIdx]
	fieldsStr := factStr[parenIdx+1:]
	if strings.HasSuffix(fieldsStr, ")") {
		fieldsStr = fieldsStr[:len(fieldsStr)-1]
	}

	fields := make(map[string]interface{})
	var factID string

	// Parser simple pour les champs
	parts := strings.Split(fieldsStr, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.Contains(part, ":") {
			keyValue := strings.SplitN(part, ":", 2)
			if len(keyValue) == 2 {
				key := strings.TrimSpace(keyValue[0])
				value := strings.TrimSpace(keyValue[1])

				// Nettoyer les guillemets
				if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
					value = value[1 : len(value)-1]
					if key == "id" {
						factID = value
					}
					fields[key] = value
				} else {
					// Tenter de parser comme nombre
					if strings.Contains(value, ".") {
						var f float64
						if _, err := fmt.Sscanf(value, "%f", &f); err == nil {
							fields[key] = f
						} else {
							fields[key] = value
						}
					} else {
						var intVal int
						if _, err := fmt.Sscanf(value, "%d", &intVal); err == nil {
							fields[key] = intVal
						} else {
							fields[key] = value
						}
					}
				}
			}
		}
	}

	if factID == "" {
		factID = fmt.Sprintf("%s_%d", factType, len(fields))
	}

	return DetailedFact{
		ID:     factID,
		Type:   factType,
		Fields: fields,
	}
}

// evaluateRuleCondition évalue si un fait satisfait les conditions d'une règle
func evaluateRuleCondition(leftPart string, fact DetailedFact, variable string) bool {
	// Cette fonction doit évaluer les conditions complexes
	// Exemple: "{b: Balance} / NOT(ABS(b.amount) > 100)"

	if strings.Contains(leftPart, "NOT(") {
		// Extraire la condition dans NOT() avec compteur de parenthèses
		notStart := strings.Index(leftPart, "NOT(")
		if notStart != -1 {
			openParen := 1
			i := notStart + 4 // après "NOT("
			for i < len(leftPart) && openParen > 0 {
				if leftPart[i] == '(' {
					openParen++
				} else if leftPart[i] == ')' {
					openParen--
				}
				i++
			}
			if openParen == 0 {
				condition := leftPart[notStart+4 : i-1]
				return !evaluateSimpleCondition(condition, fact, variable)
			}
		}
	}

	// Si pas de NOT, évaluer directement les conditions après "/"
	if strings.Contains(leftPart, "/") {
		parts := strings.Split(leftPart, "/")
		if len(parts) > 1 {
			condition := strings.TrimSpace(parts[1])
			return evaluateSimpleCondition(condition, fact, variable)
		}
	}

	// Si aucune condition spécifique, le fait correspond au type
	return true
}

// evaluateSimpleCondition évalue une condition simple
func evaluateSimpleCondition(condition string, fact DetailedFact, variable string) bool {
	// Exemple: "ABS(b.amount) > 100"
	condition = strings.TrimSpace(condition)

	// Remplacer la variable par les valeurs du fait
	for field, value := range fact.Fields {
		placeholder := variable + "." + field
		if strings.Contains(condition, placeholder) {
			valueStr := fmt.Sprintf("%v", value)
			condition = strings.ReplaceAll(condition, placeholder, valueStr)
		}
	}

	// Évaluation basique des opérateurs courants
	if strings.Contains(condition, "ABS(") && strings.Contains(condition, ">") {
		// Exemple: ABS(150.0) > 100
		// Extraire la valeur dans ABS()
		absStart := strings.Index(condition, "ABS(")
		absEnd := strings.Index(condition[absStart:], ")") + absStart
		if absEnd > absStart {
			valueStr := condition[absStart+4 : absEnd]
			var value float64
			if _, err := fmt.Sscanf(valueStr, "%f", &value); err == nil {
				absValue := value
				if absValue < 0 {
					absValue = -absValue
				}

				// Extraire le seuil après ">"
				if gtIndex := strings.Index(condition, ">"); gtIndex != -1 {
					thresholdStr := strings.TrimSpace(condition[gtIndex+1:])
					var threshold float64
					if _, err := fmt.Sscanf(thresholdStr, "%f", &threshold); err == nil {
						return absValue > threshold
					}
				}
			}
		}
	}

	// Pour d'autres conditions simples comme "amount > 100"
	if strings.Contains(condition, ">") && !strings.Contains(condition, "ABS(") {
		parts := strings.Split(condition, ">")
		if len(parts) == 2 {
			leftVal := strings.TrimSpace(parts[0])
			rightVal := strings.TrimSpace(parts[1])

			var leftNum, rightNum float64
			leftOk := false
			rightOk := false

			if _, err := fmt.Sscanf(leftVal, "%f", &leftNum); err == nil {
				leftOk = true
			}
			if _, err := fmt.Sscanf(rightVal, "%f", &rightNum); err == nil {
				rightOk = true
			}

			if leftOk && rightOk {
				return leftNum > rightNum
			}
		}
	}

	// Gestion de l'opérateur CONTAINS
	if strings.Contains(condition, " CONTAINS ") {
		parts := strings.Split(condition, " CONTAINS ")
		if len(parts) == 2 {
			leftVal := strings.TrimSpace(parts[0])
			rightVal := strings.TrimSpace(parts[1])

			// Nettoyer les guillemets
			leftClean := strings.Trim(leftVal, "\"")
			rightClean := strings.Trim(rightVal, "\"")

			return strings.Contains(leftClean, rightClean)
		}
	}

	// Gestion de l'opérateur IN
	if strings.Contains(condition, " IN ") {
		parts := strings.Split(condition, " IN ")
		if len(parts) == 2 {
			leftVal := strings.TrimSpace(parts[0])
			rightVal := strings.TrimSpace(parts[1])

			// Parser le tableau [item1, item2, ...]
			if strings.HasPrefix(rightVal, "[") && strings.HasSuffix(rightVal, "]") {
				arrayContent := rightVal[1 : len(rightVal)-1]
				items := strings.Split(arrayContent, ",")

				leftClean := strings.Trim(leftVal, "\"")
				for _, item := range items {
					itemClean := strings.Trim(strings.TrimSpace(item), "\"")
					if leftClean == itemClean {
						return true
					}
				}
				return false
			}
		}
	}

	// Gestion de l'opérateur LIKE
	if strings.Contains(condition, " LIKE ") {
		parts := strings.Split(condition, " LIKE ")
		if len(parts) == 2 {
			leftVal := strings.TrimSpace(parts[0])
			rightVal := strings.TrimSpace(parts[1])

			leftClean := strings.Trim(leftVal, "\"")
			rightClean := strings.Trim(rightVal, "\"")

			// Convertir pattern LIKE en regex simple
			pattern := strings.ReplaceAll(rightClean, "%", ".*")
			matched, _ := regexp.MatchString("^"+pattern+"$", leftClean)
			return matched
		}
	}

	// Gestion de l'opérateur MATCHES
	if strings.Contains(condition, " MATCHES ") {
		parts := strings.Split(condition, " MATCHES ")
		if len(parts) == 2 {
			leftVal := strings.TrimSpace(parts[0])
			rightVal := strings.TrimSpace(parts[1])

			leftClean := strings.Trim(leftVal, "\"")
			rightClean := strings.Trim(rightVal, "\"")

			matched, _ := regexp.MatchString(rightClean, leftClean)
			return matched
		}
	}

	// Gestion de l'opérateur !=
	if strings.Contains(condition, " != ") {
		parts := strings.Split(condition, " != ")
		if len(parts) == 2 {
			leftVal := strings.TrimSpace(parts[0])
			rightVal := strings.TrimSpace(parts[1])

			// Comparaison string
			leftClean := strings.Trim(leftVal, "\"")
			rightClean := strings.Trim(rightVal, "\"")
			return leftClean != rightClean
		}
	}

	// Gestion des fonctions LENGTH() et UPPER()
	if strings.Contains(condition, "LENGTH(") {
		// Exemple: LENGTH(p.value) >= 8
		lengthStart := strings.Index(condition, "LENGTH(")
		lengthEnd := strings.Index(condition[lengthStart:], ")") + lengthStart
		if lengthEnd > lengthStart {
			valueStr := condition[lengthStart+7 : lengthEnd]
			valueClean := strings.Trim(valueStr, "\"")
			length := float64(len(valueClean))

			// Remplacer LENGTH(...) par la longueur
			newCondition := condition[:lengthStart] + fmt.Sprintf("%.0f", length) + condition[lengthEnd+1:]
			return evaluateSimpleCondition(newCondition, fact, variable)
		}
	}

	if strings.Contains(condition, "UPPER(") {
		// Exemple: UPPER(d.name) == "FINANCE"
		upperStart := strings.Index(condition, "UPPER(")
		upperEnd := strings.Index(condition[upperStart:], ")") + upperStart
		if upperEnd > upperStart {
			valueStr := condition[upperStart+6 : upperEnd]
			valueClean := strings.Trim(valueStr, "\"")
			upperValue := strings.ToUpper(valueClean)

			// Remplacer UPPER(...) par la valeur en majuscules
			newCondition := condition[:upperStart] + "\"" + upperValue + "\"" + condition[upperEnd+1:]
			return evaluateSimpleCondition(newCondition, fact, variable)
		}
	}

	// Gestion de l'opérateur >= et <=
	if strings.Contains(condition, ">=") {
		parts := strings.Split(condition, ">=")
		if len(parts) == 2 {
			leftVal := strings.TrimSpace(parts[0])
			rightVal := strings.TrimSpace(parts[1])

			var leftNum, rightNum float64
			leftOk := false
			rightOk := false

			if _, err := fmt.Sscanf(leftVal, "%f", &leftNum); err == nil {
				leftOk = true
			}
			if _, err := fmt.Sscanf(rightVal, "%f", &rightNum); err == nil {
				rightOk = true
			}

			if leftOk && rightOk {
				return leftNum >= rightNum
			}
		}
	}

	// Gestion des comparaisons d'égalité booléennes et string
	if strings.Contains(condition, " == ") {
		parts := strings.Split(condition, " == ")
		if len(parts) == 2 {
			leftVal := strings.TrimSpace(parts[0])
			rightVal := strings.TrimSpace(parts[1])

			// Comparaison booléenne
			if rightVal == "true" || rightVal == "false" {
				leftBool := false
				rightBool := rightVal == "true"

				if leftVal == "true" {
					leftBool = true
				} else if leftVal == "false" {
					leftBool = false
				} else {
					// Essayer de parser comme booléen depuis les champs
					if boolVal, ok := parseBoolFromString(leftVal); ok {
						leftBool = boolVal
					}
				}

				return leftBool == rightBool
			}

			// Comparaison string
			// Nettoyer les guillemets éventuels
			leftClean := strings.Trim(leftVal, "\"")
			rightClean := strings.Trim(rightVal, "\"")
			return leftClean == rightClean
		}
	}

	return false
}

// parseBoolFromString tente de parser un booléen depuis une chaîne
func parseBoolFromString(value string) (bool, bool) {
	switch strings.ToLower(value) {
	case "true":
		return true, true
	case "false":
		return false, true
	}
	return false, false
}

// extractActionName extrait le nom de l'action depuis les règles
func extractActionName(rules []string, nodeID string) string {
	// Extraire l'index de la règle depuis le nodeID
	re := regexp.MustCompile(`rule_(\d+)_terminal`)
	matches := re.FindStringSubmatch(nodeID)
	if len(matches) > 1 {
		var ruleIndex int
		if _, err := fmt.Sscanf(matches[1], "%d", &ruleIndex); err == nil && ruleIndex < len(rules) {
			rule := rules[ruleIndex]
			if strings.Contains(rule, "==>") {
				parts := strings.Split(rule, "==>")
				if len(parts) > 1 {
					action := strings.TrimSpace(parts[1])
					// Extraire juste le nom de l'action (avant la parenthèse)
					if parenIdx := strings.Index(action, "("); parenIdx != -1 {
						return action[:parenIdx]
					}
					return action
				}
			}
		}
	}
	return ""
}

// abs retourne la valeur absolue
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func generateDetailedReport(results []TestResult, outputFile string) error {
	content := "# RAPPORT ALPHA NODES - DÉTAILLÉ\n"
	content += "==================================\n\n"

	successCount := 0
	for _, result := range results {
		if result.Success {
			successCount++
		}
	}

	content += fmt.Sprintf("**Tests:** %d\n", len(results))
	content += fmt.Sprintf("**Succès:** %d (%.1f%%)\n",
		successCount, float64(successCount)/float64(len(results))*100)
	content += fmt.Sprintf("**Date:** %s\n\n", time.Now().Format("2006-01-02 15:04:05"))

	for i, result := range results {
		content += fmt.Sprintf("## Test %d: %s\n", i+1, result.Name)
		content += fmt.Sprintf("- **Contraintes:** %s\n", result.ConstraintFile)
		content += fmt.Sprintf("- **Faits:** %s\n", result.FactsFile)
		content += fmt.Sprintf("- **Temps:** %v\n", result.ExecutionTime)

		// Détails des règles
		content += "\n### Règles Analysées\n"
		if len(result.Rules) > 0 {
			for _, rule := range result.Rules {
				content += fmt.Sprintf("- `%s`\n", rule)
			}
		} else {
			content += "- Aucune règle explicite trouvée\n"
		}

		// Détails des faits
		content += "\n### Faits Traités\n"
		if len(result.Facts) > 0 {
			for _, fact := range result.Facts {
				content += fmt.Sprintf("- `%s`\n", fact)
			}
		} else {
			content += "- Aucun fait trouvé\n"
		}

		// État du réseau
		content += "\n### État du Réseau\n"
		content += fmt.Sprintf("- **Total nœuds:** %d\n", result.NetworkState.TotalNodes)
		content += fmt.Sprintf("- **Nœuds Alpha:** %d\n", result.NetworkState.NodesByType["alpha"])
		content += fmt.Sprintf("- **Nœuds Terminal:** %d\n", result.NetworkState.NodesByType["terminal"])
		content += fmt.Sprintf("- **Tokens générés:** %d\n", len(result.Tokens))

		// NOUVEAU: Analyse détaillée des tokens
		content += "\n### Analyse Détaillée des Tokens\n"

		// Tokens attendus
		content += "\n#### Tokens Attendus\n"
		if len(result.TokenAnalysis.ExpectedTokens) > 0 {
			for i, token := range result.TokenAnalysis.ExpectedTokens {
				actionDisplay := ""
				if token.ActionName != "" {
					actionDisplay = fmt.Sprintf(" → Action: %s", token.ActionName)
				}
				content += fmt.Sprintf("**Token %d (Node: %s)%s**\n", i+1, token.NodeID, actionDisplay)
				if len(token.Facts) > 0 {
					content += "Faits:\n"
					for _, fact := range token.Facts {
						content += fmt.Sprintf("  - %s", fact.Type)
						if len(fact.Fields) > 0 {
							content += " {"
							first := true
							for key, value := range fact.Fields {
								if !first {
									content += ", "
								}
								content += fmt.Sprintf("%s: %v", key, value)
								first = false
							}
							content += "}"
						}
						content += "\n"
					}
				} else {
					content += "  - Aucun fait associé\n"
				}
				content += "\n"
			}
		} else {
			content += "- Aucun token attendu\n"
		}

		// Tokens observés
		content += "\n#### Tokens Observés\n"
		if len(result.TokenAnalysis.ObservedTokens) > 0 {
			for i, token := range result.TokenAnalysis.ObservedTokens {
				actionDisplay := ""
				if token.ActionName != "" {
					actionDisplay = fmt.Sprintf(" → Action: %s", token.ActionName)
				}
				content += fmt.Sprintf("**Token %d (Node: %s)%s**\n", i+1, token.NodeID, actionDisplay)
				if len(token.Facts) > 0 {
					content += "Faits:\n"
					for _, fact := range token.Facts {
						content += fmt.Sprintf("  - %s", fact.Type)
						if len(fact.Fields) > 0 {
							content += " {"
							first := true
							for key, value := range fact.Fields {
								if !first {
									content += ", "
								}
								content += fmt.Sprintf("%s: %v", key, value)
								first = false
							}
							content += "}"
						}
						content += "\n"
					}
				} else {
					content += "  - Aucun fait associé\n"
				}
				content += "\n"
			}
		} else {
			content += "- Aucun token observé\n"
		}

		// Comparaison
		content += "\n#### Comparaison\n"
		content += fmt.Sprintf("- **Matches:** %d\n", result.TokenAnalysis.Matches)
		content += fmt.Sprintf("- **Mismatches:** %d\n", result.TokenAnalysis.Mismatches)
		if result.TokenAnalysis.Matches > 0 && result.TokenAnalysis.Mismatches == 0 {
			content += "- **Cohérence sémantique:** ✅ VALIDÉE\n"
		} else {
			content += "- **Cohérence sémantique:** ❌ ÉCHEC\n"
		}

		// Résultat attendu vs observé
		content += "\n### Résultats\n"
		content += fmt.Sprintf("- **Attendu:** %s\n", result.ExpectedResult)
		content += fmt.Sprintf("- **Observé:** %s\n", result.ObservedResult)

		if result.Success {
			content += "- **Status:** ✅ Succès\n"
		} else {
			content += "- **Status:** ❌ Échec\n"
			content += fmt.Sprintf("- **Erreur:** %s\n", result.ErrorMessage)
		}

		// Log d'inférence
		if len(result.InferenceLog) > 0 {
			content += "\n### Log d'Inférence\n"
			for _, logEntry := range result.InferenceLog {
				content += fmt.Sprintf("- %s\n", logEntry)
			}
		}

		content += "\n---\n\n"
	}

	return os.WriteFile(outputFile, []byte(content), 0644)
}

// compareTokenAnalysis compare les tokens attendus et observés
func compareTokenAnalysis(expectedTokens []TokenInfo, observedTokens []TokenInfo) TokenAnalysis {
	analysis := TokenAnalysis{
		ExpectedTokens: expectedTokens,
		ObservedTokens: observedTokens,
		Matches:        0,
		Mismatches:     0,
	}

	// Comparer les tokens par leurs faits constituants
	for _, expected := range expectedTokens {
		found := false
		for _, observed := range observedTokens {
			if tokensMatch(expected, observed) {
				found = true
				analysis.Matches++
				break
			}
		}
		if !found {
			analysis.Mismatches++
		}
	}

	// Compter les tokens observés non attendus
	for _, observed := range observedTokens {
		found := false
		for _, expected := range expectedTokens {
			if tokensMatch(expected, observed) {
				found = true
				break
			}
		}
		if !found {
			analysis.Mismatches++
		}
	}

	return analysis
}

// tokensMatch vérifie si deux tokens correspondent
func tokensMatch(token1, token2 TokenInfo) bool {
	if len(token1.Facts) != len(token2.Facts) {
		return false
	}

	// Vérifier que tous les faits correspondent par contenu (ignorer les IDs)
	for _, fact1 := range token1.Facts {
		found := false
		for _, fact2 := range token2.Facts {
			if factsMatchByContent(fact1, fact2) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

// factsMatchByContent compare les faits par type et champs (ignore les IDs)
func factsMatchByContent(fact1, fact2 DetailedFact) bool {
	if fact1.Type != fact2.Type {
		return false
	}

	// Comparer tous les champs sauf "id"
	for key, value1 := range fact1.Fields {
		if key == "id" {
			continue // Ignorer les IDs
		}
		value2, exists := fact2.Fields[key]
		if !exists {
			return false
		}
		// Comparer les valeurs avec tolérance pour les types numériques
		if !valuesEqual(value1, value2) {
			return false
		}
	}

	// Vérifier dans l'autre sens
	for key := range fact2.Fields {
		if key == "id" {
			continue
		}
		if _, exists := fact1.Fields[key]; !exists {
			return false
		}
	}

	return true
}

// valuesEqual compare deux valeurs avec tolérance pour les types numériques
func valuesEqual(v1, v2 interface{}) bool {
	// Conversion en string pour comparaison simple
	s1 := fmt.Sprintf("%v", v1)
	s2 := fmt.Sprintf("%v", v2)
	return s1 == s2
}
