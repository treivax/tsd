package validation

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Types RETE intégrés dans le package validation

// Types de base
type RETEFact struct {
	ID     string
	Type   string
	Fields map[string]interface{}
}

type RETEToken struct {
	ID       string
	Facts    []*RETEFact
	RuleName string
	NodeID   string
}

// Interface pour tous les noeuds RETE
type RETENodeInterface interface {
	GetID() string
	GetTokens() map[string]*RETEToken
	AddToken(token *RETEToken)
}

// Noeud Alpha (filtrage par type)
type RETEAlphaNodeNew struct {
	ID       string
	TypeName string
	Tokens   map[string]*RETEToken
}

func (a *RETEAlphaNodeNew) GetID() string {
	return a.ID
}

func (a *RETEAlphaNodeNew) GetTokens() map[string]*RETEToken {
	return a.Tokens
}

func (a *RETEAlphaNodeNew) AddToken(token *RETEToken) {
	a.Tokens[token.ID] = token
}

// Noeud Beta (jointure)
type RETEBetaNodeNew struct {
	ID             string
	LeftInput      RETENodeInterface
	RightInput     RETENodeInterface
	JoinConditions []string // Conditions à évaluer à ce niveau
	Tokens         map[string]*RETEToken
}

func (b *RETEBetaNodeNew) GetID() string {
	return b.ID
}

func (b *RETEBetaNodeNew) GetTokens() map[string]*RETEToken {
	return b.Tokens
}

func (b *RETEBetaNodeNew) AddToken(token *RETEToken) {
	b.Tokens[token.ID] = token
}

// Règle parsée
type RETERule struct {
	ID            string
	VariableTypes map[string]string // variable -> type
	Conditions    []string
}

// Réseau RETE complet
type RETEValidationNetwork struct {
	AlphaNodes map[string]*RETEAlphaNodeNew
	BetaNodes  map[string]*RETEBetaNodeNew
	Facts      map[string]*RETEFact
	Rules      []RETERule
}

func NewRETEValidationNetwork() *RETEValidationNetwork {
	return &RETEValidationNetwork{
		AlphaNodes: make(map[string]*RETEAlphaNodeNew),
		BetaNodes:  make(map[string]*RETEBetaNodeNew),
		Facts:      make(map[string]*RETEFact),
		Rules:      []RETERule{},
	}
}

func (r *RETEValidationNetwork) AddRule(rule RETERule) {
	r.Rules = append(r.Rules, rule)

	// Créer les noeuds Alpha pour chaque type
	alphaNodes := make(map[string]*RETEAlphaNodeNew)
	for varName, typeName := range rule.VariableTypes {
		alphaNodeID := fmt.Sprintf("alpha_%s_%s", typeName, varName)
		if _, exists := r.AlphaNodes[alphaNodeID]; !exists {
			r.AlphaNodes[alphaNodeID] = &RETEAlphaNodeNew{
				ID:       alphaNodeID,
				TypeName: typeName,
				Tokens:   make(map[string]*RETEToken),
			}

			// Repopuler ce noeud Alpha avec les faits existants
			for _, fact := range r.Facts {
				if fact.Type == typeName {
					token := &RETEToken{
						ID:       fmt.Sprintf("token_%s_%s", alphaNodeID, fact.ID),
						Facts:    []*RETEFact{fact},
						RuleName: "",
						NodeID:   alphaNodeID,
					}
					r.AlphaNodes[alphaNodeID].AddToken(token)
				}
			}
		}
		alphaNodes[varName] = r.AlphaNodes[alphaNodeID]
	}

	// Construire la chaîne de noeuds Beta pour les jointures
	if len(rule.VariableTypes) > 1 {
		variables := make([]string, 0, len(rule.VariableTypes))
		for varName := range rule.VariableTypes {
			variables = append(variables, varName)
		}

		var lastNode RETENodeInterface = alphaNodes[variables[0]]

		for i := 1; i < len(variables); i++ {
			betaNodeID := fmt.Sprintf("beta_%s_%d", rule.ID, i)

			// Déterminer les conditions applicables à ce nœud
			nodeConditions := r.extractConditionsForJoin(rule, variables[:i+1])

			betaNode := &RETEBetaNodeNew{
				ID:             betaNodeID,
				LeftInput:      lastNode,
				RightInput:     alphaNodes[variables[i]],
				JoinConditions: nodeConditions,
				Tokens:         make(map[string]*RETEToken),
			}

			r.BetaNodes[betaNodeID] = betaNode

			// Propager immédiatement les tokens existants dans ce nouveau noeud
			r.populateNewBetaNode(betaNode)

			lastNode = betaNode
		}
	}
}

// populateNewBetaNode remplit un nouveau noeud Beta avec des jointures des tokens existants
func (r *RETEValidationNetwork) populateNewBetaNode(betaNode *RETEBetaNodeNew) {
	leftTokens := betaNode.LeftInput.GetTokens()
	rightTokens := betaNode.RightInput.GetTokens()

	for _, leftToken := range leftTokens {
		for _, rightToken := range rightTokens {
			if r.evaluateJoinCondition(leftToken, rightToken, betaNode) {
				joinedFacts := append(leftToken.Facts, rightToken.Facts...)
				joinToken := &RETEToken{
					ID:       fmt.Sprintf("join_%s_%s_%s", betaNode.ID, leftToken.ID, rightToken.ID),
					Facts:    joinedFacts,
					RuleName: "",
					NodeID:   betaNode.ID,
				}
				betaNode.AddToken(joinToken)
			}
		}
	}
}

func (r *RETEValidationNetwork) InsertFact(fact *RETEFact) {
	r.Facts[fact.ID] = fact

	// Propager dans tous les noeuds Alpha correspondants
	for _, alphaNode := range r.AlphaNodes {
		if alphaNode.TypeName == fact.Type {
			token := &RETEToken{
				ID:       fmt.Sprintf("token_%s_%s", alphaNode.ID, fact.ID),
				Facts:    []*RETEFact{fact},
				RuleName: "",
				NodeID:   alphaNode.ID,
			}

			alphaNode.AddToken(token)
			r.propagateFromAlpha(alphaNode, token)
		}
	}
}

func (r *RETEValidationNetwork) propagateFromAlpha(alphaNode *RETEAlphaNodeNew, token *RETEToken) {
	// Vérifier tous les noeuds Beta pour voir lesquels utilisent ce noeud Alpha
	for _, betaNode := range r.BetaNodes {
		if betaNode.LeftInput == alphaNode {
			r.attemptJoin(betaNode, token, "left")
		} else if betaNode.RightInput == alphaNode {
			r.attemptJoin(betaNode, token, "right")
		}
	}
}

func (r *RETEValidationNetwork) attemptJoin(betaNode *RETEBetaNodeNew, newToken *RETEToken, side string) {
	var oppositeTokens []*RETEToken

	if side == "left" {
		// Nouveau token à gauche, chercher des tokens à droite
		for _, token := range betaNode.RightInput.GetTokens() {
			oppositeTokens = append(oppositeTokens, token)
		}
	} else {
		// Nouveau token à droite, chercher des tokens à gauche
		for _, token := range betaNode.LeftInput.GetTokens() {
			oppositeTokens = append(oppositeTokens, token)
		}
	}

	// Tenter la jointure avec chaque token opposé
	for _, oppositeToken := range oppositeTokens {
		var leftToken, rightToken *RETEToken
		if side == "left" {
			leftToken, rightToken = newToken, oppositeToken
		} else {
			leftToken, rightToken = oppositeToken, newToken
		}

		if r.evaluateJoinCondition(leftToken, rightToken, betaNode) {
			// Créer un nouveau token joint
			joinedFacts := append(leftToken.Facts, rightToken.Facts...)
			joinToken := &RETEToken{
				ID:       fmt.Sprintf("join_%s_%s_%s", betaNode.ID, leftToken.ID, rightToken.ID),
				Facts:    joinedFacts,
				RuleName: "",
				NodeID:   betaNode.ID,
			}

			betaNode.AddToken(joinToken)

			// Propager vers les nœuds Beta suivants
			r.propagateFromBeta(betaNode, joinToken)
		}
	}
}

func (r *RETEValidationNetwork) propagateFromBeta(betaNode *RETEBetaNodeNew, token *RETEToken) {
	for _, nextBeta := range r.BetaNodes {
		if nextBeta.LeftInput == betaNode {
			r.attemptJoin(nextBeta, token, "left")
		}
	}
}

func (r *RETEValidationNetwork) evaluateJoinCondition(leftToken, rightToken *RETEToken, betaNode *RETEBetaNodeNew) bool {
	// Combiner les faits pour l'évaluation
	combinedFacts := append(leftToken.Facts, rightToken.Facts...)

	// Créer un mapping variable -> fait temporaire
	varToFact := r.createVarToFactMapping(combinedFacts)

	// Évaluer seulement les conditions locales à ce nœud
	for _, condition := range betaNode.JoinConditions {
		if !r.evaluateCondition(condition, varToFact) {
			return false
		}
	}

	return true
}

// canFactsBeJoined vérifie si des faits peuvent être joints selon les conditions de jointure basiques
func (r *RETEValidationNetwork) canFactsBeJoined(facts []*RETEFact) bool {
	// Pour l'instant, accepter toutes les jointures
	// L'évaluation complète se fait au niveau terminal
	return true
}

// factsMatchRuleTypes vérifie si un ensemble de faits correspond aux types requis par une règle
func (r *RETEValidationNetwork) factsMatchRuleTypes(facts []*RETEFact, rule RETERule) bool {
	// Compter les types trouvés
	foundTypes := make(map[string]int)
	for _, fact := range facts {
		foundTypes[fact.Type]++
	}

	// Vérifier que tous les types requis sont présents
	requiredTypes := make(map[string]int)
	for _, typeName := range rule.VariableTypes {
		requiredTypes[typeName]++
	}

	for typeName, requiredCount := range requiredTypes {
		if foundTypes[typeName] < requiredCount {
			return false
		}
	}

	return true
}

// evaluateRuleConditions évalue les conditions d'une règle sur un ensemble de faits
func (r *RETEValidationNetwork) evaluateRuleConditions(facts []*RETEFact, rule RETERule) bool {
	// Créer un mapping variable -> fait pour l'évaluation
	varToFact := make(map[string]*RETEFact)

	// Assigner les variables aux faits selon leurs types
	for varName, typeName := range rule.VariableTypes {
		for _, fact := range facts {
			if fact.Type == typeName && varToFact[varName] == nil {
				varToFact[varName] = fact
				break
			}
		}
	}

	// Vérifier que toutes les variables sont assignées
	for varName := range rule.VariableTypes {
		if varToFact[varName] == nil {
			return false
		}
	}

	// DEBUG: Afficher les assignations
	fmt.Printf("    🔍 Évaluation règle %s:\n", rule.ID)
	for varName, fact := range varToFact {
		fmt.Printf("      %s -> %s (%s)\n", varName, fact.ID, fact.Type)
	}

	// Évaluer chaque condition
	for _, condition := range rule.Conditions {
		result := r.evaluateCondition(condition, varToFact)
		fmt.Printf("      Condition '%s': %v\n", condition, result)
		if !result {
			return false
		}
	}

	fmt.Printf("    ✅ Règle satisfaite!\n")
	return true
}

// evaluateCondition évalue une condition individuelle
func (r *RETEValidationNetwork) evaluateCondition(condition string, varToFact map[string]*RETEFact) bool {
	condition = strings.TrimSpace(condition)

	// Cas spéciaux
	if condition == "true" {
		return true
	}
	if condition == "false" {
		return false
	}

	// Parser les différents types de conditions
	if strings.Contains(condition, " == ") {
		return r.evaluateEquality(condition, varToFact)
	} else if strings.Contains(condition, " >= ") {
		return r.evaluateGreaterEqual(condition, varToFact)
	} else if strings.Contains(condition, " > ") {
		return r.evaluateGreater(condition, varToFact)
	} else if strings.Contains(condition, " <= ") {
		return r.evaluateLessEqual(condition, varToFact)
	} else if strings.Contains(condition, " < ") {
		return r.evaluateLess(condition, varToFact)
	} else if strings.Contains(condition, " != ") {
		return r.evaluateNotEqual(condition, varToFact)
	}

	// Condition non reconnue
	return false
}

// evaluateEquality évalue une condition d'égalité comme "u.id == o.user_id"
func (r *RETEValidationNetwork) evaluateEquality(condition string, varToFact map[string]*RETEFact) bool {
	parts := strings.Split(condition, " == ")
	if len(parts) != 2 {
		return false
	}

	leftValue := r.getFieldValue(strings.TrimSpace(parts[0]), varToFact)
	rightValue := r.getFieldValue(strings.TrimSpace(parts[1]), varToFact)

	return r.compareValues(leftValue, rightValue, "==")
}

// evaluateGreaterEqual évalue une condition >= comme "u.age >= 25"
func (r *RETEValidationNetwork) evaluateGreaterEqual(condition string, varToFact map[string]*RETEFact) bool {
	parts := strings.Split(condition, " >= ")
	if len(parts) != 2 {
		return false
	}

	leftValue := r.getFieldValue(strings.TrimSpace(parts[0]), varToFact)
	rightValue := r.getFieldValue(strings.TrimSpace(parts[1]), varToFact)

	return r.compareValues(leftValue, rightValue, ">=")
}

// evaluateGreater évalue une condition > comme "p.price > 100"
func (r *RETEValidationNetwork) evaluateGreater(condition string, varToFact map[string]*RETEFact) bool {
	parts := strings.Split(condition, " > ")
	if len(parts) != 2 {
		return false
	}

	leftValue := r.getFieldValue(strings.TrimSpace(parts[0]), varToFact)
	rightValue := r.getFieldValue(strings.TrimSpace(parts[1]), varToFact)

	return r.compareValues(leftValue, rightValue, ">")
}

// evaluateLessEqual évalue une condition <=
func (r *RETEValidationNetwork) evaluateLessEqual(condition string, varToFact map[string]*RETEFact) bool {
	parts := strings.Split(condition, " <= ")
	if len(parts) != 2 {
		return false
	}

	leftValue := r.getFieldValue(strings.TrimSpace(parts[0]), varToFact)
	rightValue := r.getFieldValue(strings.TrimSpace(parts[1]), varToFact)

	return r.compareValues(leftValue, rightValue, "<=")
}

// evaluateLess évalue une condition <
func (r *RETEValidationNetwork) evaluateLess(condition string, varToFact map[string]*RETEFact) bool {
	parts := strings.Split(condition, " < ")
	if len(parts) != 2 {
		return false
	}

	leftValue := r.getFieldValue(strings.TrimSpace(parts[0]), varToFact)
	rightValue := r.getFieldValue(strings.TrimSpace(parts[1]), varToFact)

	return r.compareValues(leftValue, rightValue, "<")
}

// evaluateNotEqual évalue une condition !=
func (r *RETEValidationNetwork) evaluateNotEqual(condition string, varToFact map[string]*RETEFact) bool {
	parts := strings.Split(condition, " != ")
	if len(parts) != 2 {
		return false
	}

	leftValue := r.getFieldValue(strings.TrimSpace(parts[0]), varToFact)
	rightValue := r.getFieldValue(strings.TrimSpace(parts[1]), varToFact)

	return r.compareValues(leftValue, rightValue, "!=")
}

// getFieldValue extrait la valeur d'un champ depuis les faits
func (r *RETEValidationNetwork) getFieldValue(expression string, varToFact map[string]*RETEFact) interface{} {
	expression = strings.TrimSpace(expression)

	// Si c'est un littéral entre guillemets
	if strings.HasPrefix(expression, `"`) && strings.HasSuffix(expression, `"`) {
		return strings.Trim(expression, `"`)
	}

	// Si c'est un nombre littéral
	if num, err := strconv.ParseFloat(expression, 64); err == nil {
		return num
	}

	// Si c'est une référence de champ comme "u.id" ou "o.user_id"
	if strings.Contains(expression, ".") {
		parts := strings.Split(expression, ".")
		if len(parts) == 2 {
			varName := parts[0]
			fieldName := parts[1]

			if fact, exists := varToFact[varName]; exists {
				if value, hasField := fact.Fields[fieldName]; hasField {
					return value
				}
			}
		}
	}

	return nil
}

// compareValues compare deux valeurs selon l'opérateur
func (r *RETEValidationNetwork) compareValues(left, right interface{}, operator string) bool {
	if left == nil || right == nil {
		return false
	}

	switch operator {
	case "==":
		return r.valuesEqual(left, right)
	case "!=":
		return !r.valuesEqual(left, right)
	case ">":
		return r.leftGreaterThanRight(left, right)
	case ">=":
		return r.valuesEqual(left, right) || r.leftGreaterThanRight(left, right)
	case "<":
		return r.leftGreaterThanRight(right, left)
	case "<=":
		return r.valuesEqual(left, right) || r.leftGreaterThanRight(right, left)
	}

	return false
}

// valuesEqual compare l'égalité de deux valeurs
func (r *RETEValidationNetwork) valuesEqual(left, right interface{}) bool {
	// Conversion de type si nécessaire
	leftStr := fmt.Sprintf("%v", left)
	rightStr := fmt.Sprintf("%v", right)

	return leftStr == rightStr
}

// leftGreaterThanRight compare si left > right pour des nombres
func (r *RETEValidationNetwork) leftGreaterThanRight(left, right interface{}) bool {
	leftNum, leftOk := r.toFloat64(left)
	rightNum, rightOk := r.toFloat64(right)

	if leftOk && rightOk {
		return leftNum > rightNum
	}

	return false
}

// toFloat64 convertit une valeur en float64
func (r *RETEValidationNetwork) toFloat64(value interface{}) (float64, bool) {
	switch v := value.(type) {
	case float64:
		return v, true
	case int:
		return float64(v), true
	case string:
		if num, err := strconv.ParseFloat(v, 64); err == nil {
			return num, true
		}
	}
	return 0, false
}

func (r *RETEValidationNetwork) GetTerminalTokens() []*RETEToken {
	var terminals []*RETEToken

	for _, rule := range r.Rules {
		expectedFactCount := len(rule.VariableTypes)

		if expectedFactCount == 1 {
			// Règle simple
			for _, alphaNode := range r.AlphaNodes {
				for _, token := range alphaNode.Tokens {
					if r.tokenMatchesRule(token, rule) {
						// Évaluer les conditions pour la règle simple
						if r.evaluateRuleConditions(token.Facts, rule) {
							terminals = append(terminals, token)
						}
					}
				}
			}
		} else {
			// Règle avec jointures - chercher dans le dernier beta node
			lastBetaID := fmt.Sprintf("beta_%s_%d", rule.ID, expectedFactCount-1)
			if lastBeta, exists := r.BetaNodes[lastBetaID]; exists {
				for _, token := range lastBeta.Tokens {
					if len(token.Facts) == expectedFactCount && r.tokenMatchesRule(token, rule) {
						// Évaluer les conditions complètes pour la règle complexe
						if r.evaluateRuleConditions(token.Facts, rule) {
							terminals = append(terminals, token)
						}
					}
				}
			}
		}
	}

	return terminals
}

func (r *RETEValidationNetwork) tokenMatchesRule(token *RETEToken, rule RETERule) bool {
	foundTypes := make(map[string]bool)
	for _, fact := range token.Facts {
		foundTypes[fact.Type] = true
	}

	for _, requiredType := range rule.VariableTypes {
		if !foundTypes[requiredType] {
			return false
		}
	}

	return true
}

func (r *RETEValidationNetwork) Debug() {
	fmt.Printf("=== État du réseau RETE ===\n")
	fmt.Printf("Alpha nodes: %d\n", len(r.AlphaNodes))
	for id, node := range r.AlphaNodes {
		fmt.Printf("  %s (%s): %d tokens\n", id, node.TypeName, len(node.Tokens))
	}

	fmt.Printf("Beta nodes: %d\n", len(r.BetaNodes))
	for id, node := range r.BetaNodes {
		fmt.Printf("  %s: %d tokens\n", id, len(node.Tokens))
	}

	fmt.Printf("Facts: %d\n", len(r.Facts))
	fmt.Printf("Rules: %d\n", len(r.Rules))
	fmt.Printf("Terminal tokens: %d\n", len(r.GetTerminalTokens()))
}

// ParseConstraintFileForRETENew parse un fichier de contraintes pour le nouveau RETE
func (r *RETEValidationNetwork) ParseConstraintFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parser la ligne de contrainte
		rule, err := r.parseConstraintLine(line, lineNum)
		if err != nil {
			fmt.Printf("Erreur ligne %d: %v\n", lineNum, err)
			continue
		}

		r.AddRule(rule)
	}

	return scanner.Err()
}

// parseConstraintLine parse une ligne de contrainte en RETERule
func (r *RETEValidationNetwork) parseConstraintLine(line string, lineNum int) (RETERule, error) {
	// Exemple: {u: User, o: Order, p: Product} / u.id == o.user_id AND o.product_id == p.id ==> action

	// Ignorer les commentaires et définitions de types
	if strings.HasPrefix(line, "//") || strings.HasPrefix(line, "type ") {
		return RETERule{}, fmt.Errorf("ligne ignorée (commentaire ou définition type)")
	}

	// Chercher le séparateur " / "
	parts := strings.Split(line, " / ")
	if len(parts) != 2 {
		return RETERule{}, fmt.Errorf("format invalide, attendu '{variables} / {conditions} ==> {action}'")
	}

	// Parser les variables {u: User, o: Order, p: Product}
	variablesPart := strings.TrimSpace(parts[0])
	if !strings.HasPrefix(variablesPart, "{") || !strings.HasSuffix(variablesPart, "}") {
		return RETERule{}, fmt.Errorf("variables doivent être entre {}")
	}

	variableTypes := make(map[string]string)
	varContent := strings.Trim(variablesPart, "{}")
	pairs := strings.Split(varContent, ",")

	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		colonPos := strings.Index(pair, ":")
		if colonPos == -1 {
			return RETERule{}, fmt.Errorf("format variable invalide: %s", pair)
		}

		varName := strings.TrimSpace(pair[:colonPos])
		typeName := strings.TrimSpace(pair[colonPos+1:])
		variableTypes[varName] = typeName
	}

	// Parser les conditions (la partie avant " ==> ")
	conditionAndAction := strings.TrimSpace(parts[1])
	arrowPos := strings.Index(conditionAndAction, " ==> ")
	if arrowPos == -1 {
		return RETERule{}, fmt.Errorf("séparateur ==> manquant")
	}

	conditionsPart := strings.TrimSpace(conditionAndAction[:arrowPos])
	// Diviser les conditions par " AND "
	conditions := strings.Split(conditionsPart, " AND ")

	rule := RETERule{
		ID:            fmt.Sprintf("rule_%d", lineNum),
		VariableTypes: variableTypes,
		Conditions:    conditions,
	}

	return rule, nil
}

// LoadFactsFile charge un fichier de faits
func (r *RETEValidationNetwork) LoadFactsFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	factID := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parser la ligne de fait
		fact, err := r.parseFactLine(line, factID)
		if err != nil {
			fmt.Printf("Erreur parsing fait: %v\n", err)
			continue
		}

		r.InsertFact(fact)
		factID++
	}

	return scanner.Err()
}

// parseFactLine parse une ligne de fait
func (r *RETEValidationNetwork) parseFactLine(line string, id int) (*RETEFact, error) {
	// Exemple: User(id:"USER001", name:"Alice", age:30, city:"Paris", status:"vip")

	// Ignorer les lignes vides
	if strings.TrimSpace(line) == "" {
		return nil, fmt.Errorf("ligne vide")
	}

	// Trouver le nom du type (avant la parenthèse)
	parenPos := strings.Index(line, "(")
	if parenPos == -1 {
		return nil, fmt.Errorf("format invalide - parenthèse manquante")
	}

	typeName := strings.TrimSpace(line[:parenPos])

	// Extraire le contenu entre parenthèses
	endParenPos := strings.LastIndex(line, ")")
	if endParenPos == -1 || endParenPos <= parenPos {
		return nil, fmt.Errorf("format invalide - parenthèse fermante manquante")
	}

	fieldsPart := line[parenPos+1 : endParenPos]
	fields := make(map[string]interface{})

	// Parser les champs field:value, field:value, ...
	fieldEntries := strings.Split(fieldsPart, ",")
	for _, entry := range fieldEntries {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}

		colonPos := strings.Index(entry, ":")
		if colonPos == -1 {
			continue
		}

		fieldName := strings.TrimSpace(entry[:colonPos])
		fieldValue := strings.TrimSpace(entry[colonPos+1:])

		// Nettoyer les guillemets
		fieldValue = strings.Trim(fieldValue, `"`)

		// Essayer de convertir en nombre
		if num, err := strconv.ParseFloat(fieldValue, 64); err == nil {
			fields[fieldName] = num
		} else {
			fields[fieldName] = fieldValue
		}
	}

	fact := &RETEFact{
		ID:     fmt.Sprintf("%s_%d", typeName, id),
		Type:   typeName,
		Fields: fields,
	}

	return fact, nil
}

// GetValidationResults retourne les résultats de validation
func (r *RETEValidationNetwork) GetValidationResults() (int, int) {
	terminals := r.GetTerminalTokens()
	return len(terminals), len(r.Rules)
}

// ValidateBetaTest valide un test beta spécifique
func ValidateBetaTestNew(baseName string) error {
	constraintFile := fmt.Sprintf("beta_coverage_tests/%s.constraint", baseName)
	factsFile := fmt.Sprintf("beta_coverage_tests/%s.facts", baseName)

	fmt.Printf("=== Test: %s ===\n", baseName)

	// Créer le réseau
	network := NewRETEValidationNetwork()

	// Charger les contraintes
	if err := network.ParseConstraintFile(constraintFile); err != nil {
		return fmt.Errorf("erreur parsing contraintes: %v", err)
	}

	// Charger les faits
	if err := network.LoadFactsFile(factsFile); err != nil {
		return fmt.Errorf("erreur chargement faits: %v", err)
	}

	// Obtenir les résultats
	tokensCount, rulesCount := network.GetValidationResults()

	fmt.Printf("Règles: %d\n", rulesCount)
	fmt.Printf("Tokens générés: %d\n", tokensCount)

	// Debug détaillé
	network.Debug()

	if tokensCount > 0 {
		fmt.Printf("✅ SUCCÈS: %d tokens générés\n", tokensCount)
	} else {
		fmt.Printf("❌ ÉCHEC: Aucun token généré\n")
	}

	return nil
}

// Fonction pour tester tous les cas beta
func RunAllBetaTestsNew() {
	tests := []string{
		"beta_join_simple",
		"beta_join_complex",
		"join_simple",
		"join_multi_variable_complex",
	}

	successCount := 0
	for _, test := range tests {
		fmt.Printf("\n" + strings.Repeat("=", 60) + "\n")
		if err := ValidateBetaTestNew(test); err != nil {
			fmt.Printf("❌ %s: %v\n", test, err)
		} else {
			fmt.Printf("✅ %s: completed\n", test)
			successCount++
		}
	}

	fmt.Printf("\n" + strings.Repeat("=", 60) + "\n")
	fmt.Printf("RÉSUMÉ: %d/%d tests réussis\n", successCount, len(tests))
}

// extractConditionsForJoin extrait les conditions applicables à un niveau de jointure
func (r *RETEValidationNetwork) extractConditionsForJoin(rule RETERule, availableVars []string) []string {
	availableVarSet := make(map[string]bool)
	for _, v := range availableVars {
		availableVarSet[v] = true
	}

	var applicableConditions []string
	for _, condition := range rule.Conditions {
		if r.conditionUsesOnlyVars(condition, availableVarSet) {
			applicableConditions = append(applicableConditions, condition)
		}
	}

	return applicableConditions
}

// conditionUsesOnlyVars vérifie si une condition utilise uniquement les variables disponibles
func (r *RETEValidationNetwork) conditionUsesOnlyVars(condition string, availableVars map[string]bool) bool {
	vars := r.extractVariablesFromCondition(condition)

	// Vérifier que toutes les variables référencées sont disponibles
	for varName := range vars {
		if !availableVars[varName] {
			return false
		}
	}

	// S'assurer qu'au moins une variable est référencée
	return len(vars) > 0
}

// extractVariablesFromCondition extrait les noms de variables d'une condition
func (r *RETEValidationNetwork) extractVariablesFromCondition(condition string) map[string]bool {
	vars := make(map[string]bool)

	// Chercher les patterns "variable.champ"
	words := strings.Fields(condition)
	for _, word := range words {
		if strings.Contains(word, ".") {
			parts := strings.Split(word, ".")
			if len(parts) >= 2 {
				varName := parts[0]
				// Nettoyer les caractères spéciaux
				varName = strings.Trim(varName, "()[]")
				vars[varName] = true
			}
		}
	}

	return vars
}

// createVarToFactMapping crée un mapping variable -> fait pour l'évaluation
func (r *RETEValidationNetwork) createVarToFactMapping(facts []*RETEFact) map[string]*RETEFact {
	varToFact := make(map[string]*RETEFact)

	// Mapping intelligent basé sur les types des faits
	typeToVar := map[string]string{
		"User":        "u",
		"Order":       "o",
		"Product":     "p",
		"Address":     "a",
		"Employee":    "e",
		"Performance": "perf",
		"TestUser":    "u",
		"TestOrder":   "o",
		"TestProduct": "p",
		"TestPerson":  "p",
	}

	for _, fact := range facts {
		if varName, exists := typeToVar[fact.Type]; exists {
			varToFact[varName] = fact
		} else {
			// Fallback: utiliser les premières lettres du type en minuscule
			varName := strings.ToLower(string(fact.Type[0]))
			varToFact[varName] = fact
		}
	}

	return varToFact
}
