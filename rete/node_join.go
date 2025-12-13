// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

type JoinNode struct {
	BaseNode
	Condition      map[string]interface{} `json:"condition"`
	LeftVariables  []string               `json:"left_variables"`
	RightVariables []string               `json:"right_variables"`
	AllVariables   []string               `json:"all_variables"`
	VariableTypes  map[string]string      `json:"variable_types"` // Nouveau: mapping variable -> type
	JoinConditions []JoinCondition        `json:"join_conditions"`
	mutex          sync.RWMutex
	// Mémoires séparées pour architecture RETE propre
	LeftMemory   *WorkingMemory // Tokens venant de la gauche
	RightMemory  *WorkingMemory // Tokens venant de la droite
	ResultMemory *WorkingMemory // Tokens de jointure réussie
}

// JoinCondition représente une condition de jointure entre variables
type JoinCondition struct {
	LeftField  string `json:"left_field"`  // p.id
	RightField string `json:"right_field"` // o.customer_id
	LeftVar    string `json:"left_var"`    // p
	RightVar   string `json:"right_var"`   // o
	Operator   string `json:"operator"`    // ==
}

// NewJoinNode crée un nouveau nœud de jointure
func NewJoinNode(nodeID string, condition map[string]interface{}, leftVars []string, rightVars []string, varTypes map[string]string, storage Storage) *JoinNode {
	allVars := append(leftVars, rightVars...)

	// Extract beta condition from composite condition if present
	// Composite conditions are in format: {"beta": ..., "alpha": ...}
	conditionForExtraction := condition
	if betaCond, hasBeta := condition["beta"]; hasBeta {
		if betaMap, ok := betaCond.(map[string]interface{}); ok {
			conditionForExtraction = betaMap
		}
	}

	return &JoinNode{
		BaseNode: BaseNode{
			ID:          nodeID,
			Type:        "join",
			Memory:      &WorkingMemory{NodeID: nodeID, Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
			Children:    make([]Node, 0),
			Storage:     storage,
			createdAt:   time.Now(),
		},
		Condition:      condition,
		LeftVariables:  leftVars,
		RightVariables: rightVars,
		AllVariables:   allVars,
		VariableTypes:  varTypes,
		JoinConditions: extractJoinConditions(conditionForExtraction),
		// Initialiser les mémoires séparées
		LeftMemory:   &WorkingMemory{NodeID: nodeID + "_left", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
		RightMemory:  &WorkingMemory{NodeID: nodeID + "_right", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
		ResultMemory: &WorkingMemory{NodeID: nodeID + "_result", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
	}
}

// ActivateLeft traite les tokens de la gauche (généralement des AlphaNodes)
func (jn *JoinNode) ActivateLeft(token *Token) error {
	// Enregistrer l'activation
	jn.recordActivation()
	
	logger := GetDebugLogger()

	logger.Log("[JOIN_%s] ActivateLeft: token vars=%v", jn.ID, token.GetVariables())
	logger.LogBindings(fmt.Sprintf("JOIN_%s ActivateLeft", jn.ID), token.Bindings)

	// Stocker le token dans la mémoire gauche
	jn.mutex.Lock()
	jn.LeftMemory.AddToken(token)
	leftSize := len(jn.LeftMemory.Tokens)
	rightSize := len(jn.RightMemory.Tokens)
	jn.mutex.Unlock()

	logger.Log("[JOIN_%s] After adding to LeftMemory: left=%d, right=%d", jn.ID, leftSize, rightSize)

	// Essayer de joindre avec tous les tokens de la mémoire droite
	rightTokens := jn.RightMemory.GetTokens()
	logger.Log("[JOIN_%s] Attempting join with %d right tokens", jn.ID, len(rightTokens))

	for _, rightToken := range rightTokens {
		logger.Log("[JOIN_%s] Trying join: left_vars=%v + right_vars=%v", jn.ID, token.GetVariables(), rightToken.GetVariables())
		if joinedToken := jn.performJoinWithTokens(token, rightToken); joinedToken != nil {

			// Stocker uniquement les tokens de jointure réussie
			joinedToken.IsJoinResult = true
			jn.mutex.Lock()
			jn.ResultMemory.AddToken(joinedToken)
			jn.Memory.AddToken(joinedToken) // Pour compatibilité avec le comptage
			jn.mutex.Unlock()

			logger.Log("[JOIN_%s] ✓ Join successful, propagating token with vars=%v", jn.ID, joinedToken.GetVariables())

			if err := jn.PropagateToChildren(nil, joinedToken); err != nil {
				return err
			}
		}
	}

	logger.LogMemorySizes(jn.ID, len(jn.LeftMemory.Tokens), len(jn.RightMemory.Tokens), len(jn.ResultMemory.Tokens))
	return nil
}

// ActivateRetract retrait des tokens contenant le fait rétracté des 3 mémoires.
// factID doit être l'identifiant interne (Type_ID).
// Refactorisé pour réduire la complexité et améliorer la lisibilité.
func (jn *JoinNode) ActivateRetract(factID string) error {
	jn.mutex.Lock()

	// Retirer des 3 mémoires
	leftRemoved := jn.retractFromMemory(jn.LeftMemory, factID)
	rightRemoved := jn.retractFromMemory(jn.RightMemory, factID)
	resultRemoved := jn.retractFromResultMemory(factID)

	jn.mutex.Unlock()

	// Log si des tokens ont été retirés
	totalRemoved := len(leftRemoved) + len(rightRemoved) + len(resultRemoved)
	if totalRemoved > 0 {
		fmt.Printf("🗑️  [JOIN_%s] Rétractation: %d tokens retirés (L:%d R:%d RES:%d)\n",
			jn.ID, totalRemoved, len(leftRemoved), len(rightRemoved), len(resultRemoved))
	}

	return jn.PropagateRetractToChildren(factID)
}

// retractFromMemory retire les tokens contenant le fait spécifié d'une mémoire.
// Retourne la liste des IDs des tokens retirés.
func (jn *JoinNode) retractFromMemory(memory *WorkingMemory, factID string) []string {
	var tokensToRemove []string

	for tokenID, token := range memory.Tokens {
		if jn.tokenContainsFact(token, factID) {
			tokensToRemove = append(tokensToRemove, tokenID)
		}
	}

	for _, tokenID := range tokensToRemove {
		delete(memory.Tokens, tokenID)
	}

	return tokensToRemove
}

// retractFromResultMemory retire les tokens de la mémoire de résultats.
// Met aussi à jour la mémoire principale pour compatibilité.
func (jn *JoinNode) retractFromResultMemory(factID string) []string {
	var tokensToRemove []string

	for tokenID, token := range jn.ResultMemory.Tokens {
		if jn.tokenContainsFact(token, factID) {
			tokensToRemove = append(tokensToRemove, tokenID)
		}
	}

	for _, tokenID := range tokensToRemove {
		delete(jn.ResultMemory.Tokens, tokenID)
		delete(jn.Memory.Tokens, tokenID) // Synchroniser avec mémoire principale
	}

	return tokensToRemove
}

// tokenContainsFact vérifie si un token contient le fait spécifié.
func (jn *JoinNode) tokenContainsFact(token *Token, factID string) bool {
	for _, fact := range token.Facts {
		if fact.GetInternalID() == factID {
			return true
		}
	}
	return false
}

// ActivateRight traite les faits de la droite (nouveau fait injecté via AlphaNode)
func (jn *JoinNode) ActivateRight(fact *Fact) error {
	// Enregistrer l'activation
	jn.recordActivation()
	
	logger := GetDebugLogger()

	logger.Log("[JOIN_%s] ActivateRight: fact type=%s, id=%s", jn.ID, fact.Type, fact.ID)

	// Convertir le fait en token pour la mémoire droite
	factVar := jn.getVariableForFact(fact)
	if factVar == "" {
		logger.Log("[JOIN_%s] Fact type %s not applicable to this JoinNode (expected RightVars=%v)",
			jn.ID, fact.Type, jn.RightVariables)
		return nil // Fait non applicable à ce JoinNode
	}

	logger.Log("[JOIN_%s] Fact mapped to variable: %s", jn.ID, factVar)

	factToken := &Token{
		ID:       fmt.Sprintf("right_token_%s_%s", jn.ID, fact.ID),
		Facts:    []*Fact{fact},
		NodeID:   jn.ID,
		Bindings: NewBindingChainWith(factVar, fact),
	}

	// Stocker le token dans la mémoire droite
	jn.mutex.Lock()
	jn.RightMemory.AddToken(factToken)
	leftSize := len(jn.LeftMemory.Tokens)
	rightSize := len(jn.RightMemory.Tokens)
	jn.mutex.Unlock()

	logger.Log("[JOIN_%s] After adding to RightMemory: left=%d, right=%d", jn.ID, leftSize, rightSize)

	// Essayer de joindre avec tous les tokens de la mémoire gauche
	leftTokens := jn.LeftMemory.GetTokens()
	logger.Log("[JOIN_%s] Attempting join with %d left tokens", jn.ID, len(leftTokens))

	for _, leftToken := range leftTokens {
		logger.Log("[JOIN_%s] Trying join: left_vars=%v + right_vars=%v", jn.ID, leftToken.GetVariables(), factToken.GetVariables())
		if joinedToken := jn.performJoinWithTokens(leftToken, factToken); joinedToken != nil {

			// Stocker uniquement les tokens de jointure réussie
			joinedToken.IsJoinResult = true
			jn.mutex.Lock()
			jn.ResultMemory.AddToken(joinedToken)
			jn.Memory.AddToken(joinedToken) // Pour compatibilité avec le comptage
			jn.mutex.Unlock()

			logger.Log("[JOIN_%s] ✓ Join successful, propagating token with vars=%v", jn.ID, joinedToken.GetVariables())

			if err := jn.PropagateToChildren(nil, joinedToken); err != nil {
				return err
			}
		}
	}

	logger.LogMemorySizes(jn.ID, len(jn.LeftMemory.Tokens), len(jn.RightMemory.Tokens), len(jn.ResultMemory.Tokens))
	return nil
}

// performJoinWithTokens effectue la jointure entre deux tokens avec BindingChain immuable.
//
// IMPORTANT: Cette fonction utilise maintenant BindingChain.Merge() pour combiner
// les bindings de manière immuable, garantissant qu'aucun binding n'est perdu.
func (jn *JoinNode) performJoinWithTokens(token1 *Token, token2 *Token) *Token {
	logger := GetDebugLogger()

	// Vérifier que les tokens ont des variables différentes
	if !jn.tokensHaveDifferentVariables(token1, token2) {
		logger.Log("[JOIN_%s] Tokens have same variables, skipping join", jn.ID)
		return nil
	}

	// Combiner les bindings de manière immuable
	// Merge garantit que tous les bindings des deux tokens sont préservés
	combinedBindings := token1.Bindings.Merge(token2.Bindings)

	logger.LogJoinNode(jn.ID, "performJoinWithTokens", map[string]interface{}{
		"token1_vars":   token1.GetVariables(),
		"token2_vars":   token2.GetVariables(),
		"combined_vars": combinedBindings.Variables(),
	})

	// Valider les conditions de jointure
	if !jn.evaluateJoinConditions(combinedBindings) {
		logger.Log("[JOIN_%s] Join conditions failed", jn.ID)
		return nil // Jointure échoue
	}

	logger.Log("[JOIN_%s] ✓ Join successful, combined vars: %v", jn.ID, combinedBindings.Variables())

	// Créer et retourner le token joint avec la chaîne combinée
	return &Token{
		ID:       fmt.Sprintf("%s_JOIN_%s", token1.ID, token2.ID),
		Bindings: combinedBindings,
		NodeID:   jn.ID,
		Facts:    append(token1.Facts, token2.Facts...),
	}
}

// tokensHaveDifferentVariables vérifie que les tokens représentent des variables différentes
func (jn *JoinNode) tokensHaveDifferentVariables(token1 *Token, token2 *Token) bool {
	vars1 := token1.GetVariables()
	vars2 := token2.GetVariables()

	for _, var1 := range vars1 {
		for _, var2 := range vars2 {
			if var1 == var2 {
				return false // Même variable = pas de jointure possible
			}
		}
	}
	return true
}

// getVariableForFact détermine la variable associée à un fait basé sur son type
func (jn *JoinNode) getVariableForFact(fact *Fact) string {
	// Chercher uniquement dans RightVariables, pas dans toutes les variables
	// Car ce fait arrive par le côté droit du JoinNode
	for _, varName := range jn.RightVariables {
		if expectedType, exists := jn.VariableTypes[varName]; exists {
			if expectedType == fact.Type {
				return varName
			}
		}
	}

	// Si pas trouvé dans RightVariables, chercher dans AllVariables (fallback)
	for _, varName := range jn.AllVariables {
		if expectedType, exists := jn.VariableTypes[varName]; exists {
			if expectedType == fact.Type {
				// Vérifier que cette variable n'est pas déjà dans une autre catégorie
				found := false
				for _, rv := range jn.RightVariables {
					if rv == varName {
						found = true
						break
					}
				}
				if !found {
					return varName
				}
			}
		}
	}

	return ""
}

// evaluateJoinConditions vérifie si toutes les conditions de jointure sont respectées.
// Accepte maintenant BindingChain au lieu de map[string]*Fact.
// Refactorisé pour réduire la complexité cyclomatique de 21 à <10.
func (jn *JoinNode) evaluateJoinConditions(bindings *BindingChain) bool {
	// Validation initiale
	if !jn.validateBindingsForJoin(bindings) {
		return false
	}

	// Étape 1: Évaluer les conditions de jointure simples (field-to-field)
	if !jn.evaluateSimpleConditions(bindings) {
		return false
	}

	// Étape 2: Évaluer les conditions complètes si présentes
	return jn.evaluateComplexConditions(bindings)
}

// validateBindingsForJoin vérifie que les bindings sont suffisants pour une jointure.
func (jn *JoinNode) validateBindingsForJoin(bindings *BindingChain) bool {
	return bindings != nil && bindings.Len() >= 2
}

// evaluateSimpleConditions évalue les conditions de jointure simples.
func (jn *JoinNode) evaluateSimpleConditions(bindings *BindingChain) bool {
	if len(jn.JoinConditions) == 0 {
		return true
	}
	return jn.evaluateSimpleJoinConditions(bindings)
}

// evaluateComplexConditions évalue les conditions complètes avec contraintes additionnelles.
func (jn *JoinNode) evaluateComplexConditions(bindings *BindingChain) bool {
	if jn.Condition == nil {
		return true
	}

	// Unwrap composite condition si présent
	actualCondition := jn.unwrapCompositeCondition()

	condType, exists := actualCondition["type"].(string)
	if !exists {
		return true
	}

	// Déléguer selon le type de condition
	switch condType {
	case "constraint":
		return jn.evaluateConstraintCondition(actualCondition)
	case "comparison":
		// Déjà validé par evaluateSimpleJoinConditions
		return true
	case "logicalExpr":
		return jn.evaluateLogicalExprCondition(actualCondition, bindings)
	default:
		return true
	}
}

// unwrapCompositeCondition décompose une condition composite (beta + alpha).
func (jn *JoinNode) unwrapCompositeCondition() map[string]interface{} {
	actualCondition := jn.Condition

	// Extract beta condition if composite
	if betaCond, isBeta := jn.Condition["beta"]; isBeta {
		if betaMap, ok := betaCond.(map[string]interface{}); ok {
			actualCondition = betaMap
		}
	}

	// Unwrap constraint wrapper if present
	if condType, exists := actualCondition["type"].(string); exists && condType == "constraint" {
		if constraint, ok := actualCondition["constraint"].(map[string]interface{}); ok {
			actualCondition = constraint
		}
	}

	return actualCondition
}

// evaluateConstraintCondition évalue une condition de type "constraint".
func (jn *JoinNode) evaluateConstraintCondition(condition map[string]interface{}) bool {
	// Les contraintes sont unwrappées dans unwrapCompositeCondition
	return true
}

// evaluateLogicalExprCondition évalue une condition de type "logicalExpr".
func (jn *JoinNode) evaluateLogicalExprCondition(condition map[string]interface{}, bindings *BindingChain) bool {
	alphaConditions := jn.extractAlphaConditions(condition)
	if len(alphaConditions) == 0 {
		// Pas de contraintes alpha, seulement des joins déjà validés
		return true
	}

	// Évaluer chaque contrainte alpha
	return jn.evaluateAlphaConditions(alphaConditions, bindings)
}

// evaluateAlphaConditions évalue toutes les contraintes alpha.
func (jn *JoinNode) evaluateAlphaConditions(alphaConditions []map[string]interface{}, bindings *BindingChain) bool {
	evaluator := NewAlphaConditionEvaluator()
	evaluator.SetPartialEvalMode(true)

	// Lier toutes les variables aux faits
	jn.bindVariablesToEvaluator(evaluator, bindings)

	// Évaluer chaque contrainte
	for _, alphaCond := range alphaConditions {
		result, err := evaluator.evaluateExpression(alphaCond)
		if err != nil {
			// Erreur d'évaluation - accepter par défaut
			continue
		}
		if !result {
			return false
		}
	}

	return true
}

// bindVariablesToEvaluator lie toutes les variables de bindings à l'évaluateur.
func (jn *JoinNode) bindVariablesToEvaluator(evaluator *AlphaConditionEvaluator, bindings *BindingChain) {
	vars := bindings.Variables()
	for _, varName := range vars {
		fact := bindings.Get(varName)
		if fact != nil {
			evaluator.variableBindings[varName] = fact
		}
	}
}

// extractAlphaConditions extrait les conditions alpha (non-join) d'une logicalExpr.
// Refactorisé pour améliorer la lisibilité et réduire la complexité.
func (jn *JoinNode) extractAlphaConditions(condition map[string]interface{}) []map[string]interface{} {
	var alphaConditions []map[string]interface{}

	// Extraire de la partie gauche
	if left, ok := condition["left"].(map[string]interface{}); ok {
		if isAlphaCondition(left) {
			alphaConditions = append(alphaConditions, left)
		}
	}

	// Extraire des opérations
	alphaFromOps := jn.extractAlphaFromOperations(condition)
	alphaConditions = append(alphaConditions, alphaFromOps...)

	return alphaConditions
}

// extractAlphaFromOperations extrait les conditions alpha depuis la liste d'opérations.
func (jn *JoinNode) extractAlphaFromOperations(condition map[string]interface{}) []map[string]interface{} {
	var alphaConditions []map[string]interface{}

	operationsRaw, exists := condition["operations"]
	if !exists {
		return alphaConditions
	}

	// Essayer []interface{} en premier
	if operations, ok := operationsRaw.([]interface{}); ok {
		for _, op := range operations {
			if opMap, ok := op.(map[string]interface{}); ok {
				if right, ok := opMap["right"].(map[string]interface{}); ok {
					if isAlphaCondition(right) {
						alphaConditions = append(alphaConditions, right)
					}
				}
			}
		}
		return alphaConditions
	}

	// Essayer []map[string]interface{} en fallback
	if operations, ok := operationsRaw.([]map[string]interface{}); ok {
		for _, opMap := range operations {
			if right, ok := opMap["right"].(map[string]interface{}); ok {
				if isAlphaCondition(right) {
					alphaConditions = append(alphaConditions, right)
				}
			}
		}
	}

	return alphaConditions
}

// isAlphaCondition détermine si une condition est une contrainte alpha (pas une jointure)
func isAlphaCondition(condition map[string]interface{}) bool {
	if condType, exists := condition["type"].(string); exists && condType == "comparison" {
		// Vérifier si c'est une comparaison field-to-constant (alpha) ou field-to-field (join)
		left, leftOk := condition["left"].(map[string]interface{})
		right, rightOk := condition["right"].(map[string]interface{})

		if !leftOk || !rightOk {
			return false
		}

		leftType, _ := left["type"].(string)
		rightType, _ := right["type"].(string)

		// Si les deux côtés sont des fieldAccess, c'est une condition de jointure
		if leftType == "fieldAccess" && rightType == "fieldAccess" {
			return false
		}

		// Sinon, c'est une condition alpha
		return true
	}

	return false
}

// evaluateSimpleJoinConditions évalue les conditions de jointure simples (champ à champ).
//
// Accepte maintenant BindingChain au lieu de map[string]*Fact.
// Refactorisé pour réduire la complexité cyclomatique de 26 à <10.
func (jn *JoinNode) evaluateSimpleJoinConditions(bindings *BindingChain) bool {
	logger := GetDebugLogger()

	logger.Log("[JOIN_%s] Evaluating %d join conditions", jn.ID, len(jn.JoinConditions))
	logger.LogBindings(fmt.Sprintf("JOIN_%s", jn.ID), bindings)

	for i, joinCondition := range jn.JoinConditions {
		if !jn.evaluateSingleJoinCondition(bindings, joinCondition, i, logger) {
			return false
		}
	}

	logger.Log("[JOIN_%s] ✓ All join conditions passed", jn.ID)
	return true
}

// evaluateSingleJoinCondition évalue une seule condition de jointure.
// Complexité réduite en extrayant la logique de chaque condition.
func (jn *JoinNode) evaluateSingleJoinCondition(bindings *BindingChain, cond JoinCondition, index int, logger *DebugLogger) bool {
	// Étape 1: Récupérer les faits
	leftFact, rightFact := jn.getJoinFacts(bindings, cond, index, logger)
	if leftFact == nil || rightFact == nil {
		// Skip si variables non disponibles (cascade joins)
		return true
	}

	// Étape 2: Récupérer les valeurs des champs
	leftValue, rightValue, ok := jn.getFieldValues(leftFact, rightFact, cond, index, logger)
	if !ok {
		return false
	}

	// Étape 3: Évaluer l'opérateur
	if !jn.evaluateOperator(cond.Operator, leftValue, rightValue, cond, index, logger) {
		return false
	}

	logger.Log("[JOIN_%s] Condition[%d] PASS: %s.%s %s %s.%s",
		jn.ID, index, cond.LeftVar, cond.LeftField,
		cond.Operator,
		cond.RightVar, cond.RightField)
	return true
}

// getJoinFacts récupère les faits gauche et droit pour une condition de jointure.
// Retourne (leftFact, rightFact) ou (nil, nil) si skip nécessaire.
func (jn *JoinNode) getJoinFacts(bindings *BindingChain, cond JoinCondition, index int, logger *DebugLogger) (*Fact, *Fact) {
	leftFact := bindings.Get(cond.LeftVar)
	rightFact := bindings.Get(cond.RightVar)

	// Skip conditions that reference variables not available at this join level
	if leftFact == nil || rightFact == nil {
		logger.Log("[JOIN_%s] Condition[%d] SKIP: leftVar=%s (found=%v), rightVar=%s (found=%v)",
			jn.ID, index, cond.LeftVar, leftFact != nil, cond.RightVar, rightFact != nil)
		return nil, nil
	}

	return leftFact, rightFact
}

// getFieldValues extrait les valeurs des champs depuis les faits.
// Retourne (leftValue, rightValue, true) si succès, (nil, nil, false) si échec.
func (jn *JoinNode) getFieldValues(leftFact, rightFact *Fact, cond JoinCondition, index int, logger *DebugLogger) (interface{}, interface{}, bool) {
	leftValue, leftExists := leftFact.Fields[cond.LeftField]
	rightValue, rightExists := rightFact.Fields[cond.RightField]

	if !leftExists || !rightExists {
		logger.Log("[JOIN_%s] Condition[%d] FAIL: field not found - %s.%s (exists=%v), %s.%s (exists=%v)",
			jn.ID, index, cond.LeftVar, cond.LeftField, leftExists,
			cond.RightVar, cond.RightField, rightExists)
		return nil, nil, false
	}

	return leftValue, rightValue, true
}

// evaluateOperator évalue un opérateur de comparaison.
// Complexité réduite en déléguant les comparaisons numériques.
func (jn *JoinNode) evaluateOperator(operator string, leftValue, rightValue interface{}, cond JoinCondition, index int, logger *DebugLogger) bool {
	switch operator {
	case "==":
		return jn.evaluateEquality(leftValue, rightValue, cond, index, logger)
	case "!=":
		return jn.evaluateInequality(leftValue, rightValue, cond, index, logger)
	case "<", ">", "<=", ">=":
		return jn.evaluateNumericComparison(operator, leftValue, rightValue)
	default:
		logger.Log("[JOIN_%s] Condition[%d] FAIL: unknown operator %s", jn.ID, index, operator)
		return false
	}
}

// evaluateEquality évalue l'opérateur ==
func (jn *JoinNode) evaluateEquality(leftValue, rightValue interface{}, cond JoinCondition, index int, logger *DebugLogger) bool {
	if leftValue != rightValue {
		logger.Log("[JOIN_%s] Condition[%d] FAIL: %s.%s (%v) == %s.%s (%v)",
			jn.ID, index, cond.LeftVar, cond.LeftField, leftValue,
			cond.RightVar, cond.RightField, rightValue)
		return false
	}
	return true
}

// evaluateInequality évalue l'opérateur !=
func (jn *JoinNode) evaluateInequality(leftValue, rightValue interface{}, cond JoinCondition, index int, logger *DebugLogger) bool {
	if leftValue == rightValue {
		logger.Log("[JOIN_%s] Condition[%d] FAIL: %s.%s (%v) != %s.%s (%v)",
			jn.ID, index, cond.LeftVar, cond.LeftField, leftValue,
			cond.RightVar, cond.RightField, rightValue)
		return false
	}
	return true
}

// evaluateNumericComparison évalue les opérateurs de comparaison numérique (<, >, <=, >=).
func (jn *JoinNode) evaluateNumericComparison(operator string, leftValue, rightValue interface{}) bool {
	leftFloat, leftOk := convertToFloat64(leftValue)
	if !leftOk {
		return false
	}

	rightFloat, rightOk := convertToFloat64(rightValue)
	if !rightOk {
		return false
	}

	switch operator {
	case "<":
		return leftFloat < rightFloat
	case ">":
		return leftFloat > rightFloat
	case "<=":
		return leftFloat <= rightFloat
	case ">=":
		return leftFloat >= rightFloat
	default:
		return false
	}
}

// convertToFloat64 tente de convertir une valeur en float64
func convertToFloat64(value interface{}) (float64, bool) {
	switch v := value.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f, true
		}
		return 0, false
	default:
		return 0, false
	}
}

// extractJoinConditions extrait les conditions de jointure d'une condition complexe.
// Refactorisé pour réduire la complexité cyclomatique de 22 à <10.
func extractJoinConditions(condition map[string]interface{}) []JoinCondition {
	for key, value := range condition {
		fmt.Printf("    %s: %v (type: %T)\n", key, value, value)
	}

	conditionType, _ := condition["type"].(string)

	switch conditionType {
	case "constraint":
		return extractConstraintJoinConditions(condition)
	case "exists":
		return extractExistsJoinConditions(condition)
	case "comparison":
		return extractComparisonJoinConditions(condition)
	case "logicalExpr":
		return extractLogicalExprJoinConditions(condition)
	default:
		return []JoinCondition{}
	}
}

// extractConstraintJoinConditions extrait les conditions depuis un type "constraint".
func extractConstraintJoinConditions(condition map[string]interface{}) []JoinCondition {
	if innerCondition, ok := condition["constraint"].(map[string]interface{}); ok {
		fmt.Printf("  ✅ Sous-condition extraite, analyse récursive\n")
		return extractJoinConditions(innerCondition)
	}
	return []JoinCondition{}
}

// extractExistsJoinConditions extrait les conditions depuis un type "exists".
func extractExistsJoinConditions(condition map[string]interface{}) []JoinCondition {
	var joinConditions []JoinCondition

	if conditionsData, ok := condition["conditions"].([]map[string]interface{}); ok {
		fmt.Printf("  ✅ Array de conditions EXISTS trouvé: %d conditions\n", len(conditionsData))
		for _, subCondition := range conditionsData {
			subJoinConditions := extractJoinConditions(subCondition)
			joinConditions = append(joinConditions, subJoinConditions...)
		}
	}

	return joinConditions
}

// extractComparisonJoinConditions extrait une condition de jointure depuis un type "comparison".
func extractComparisonJoinConditions(condition map[string]interface{}) []JoinCondition {
	fmt.Printf("  ✅ Condition de comparaison détectée\n")

	left, leftOk := condition["left"].(map[string]interface{})
	right, rightOk := condition["right"].(map[string]interface{})

	if !leftOk || !rightOk {
		return []JoinCondition{}
	}

	fmt.Printf("  ✅ Left et Right extraits\n")

	// Vérifier si c'est une jointure field-to-field
	leftType, _ := left["type"].(string)
	rightType, _ := right["type"].(string)

	if leftType != "fieldAccess" || rightType != "fieldAccess" {
		return []JoinCondition{}
	}

	// Extraire les détails de la condition de jointure
	return extractFieldAccessJoinCondition(left, right, condition)
}

// extractFieldAccessJoinCondition crée une JoinCondition depuis des fieldAccess.
func extractFieldAccessJoinCondition(left, right, condition map[string]interface{}) []JoinCondition {
	fmt.Printf("  ✅ Condition de jointure fieldAccess détectée\n")

	leftObj, _ := left["object"].(string)
	leftField, _ := left["field"].(string)
	rightObj, _ := right["object"].(string)
	rightField, _ := right["field"].(string)
	operator, _ := condition["operator"].(string)

	fmt.Printf("    📌 %s.%s %s %s.%s\n", leftObj, leftField, operator, rightObj, rightField)

	return []JoinCondition{{
		LeftField:  leftField,
		RightField: rightField,
		LeftVar:    leftObj,
		RightVar:   rightObj,
		Operator:   operator,
	}}
}

// extractLogicalExprJoinConditions extrait les conditions depuis un type "logicalExpr".
func extractLogicalExprJoinConditions(condition map[string]interface{}) []JoinCondition {
	fmt.Printf("  ✅ LogicalExpr détectée, extraction des conditions\n")

	var joinConditions []JoinCondition

	// Extraire la partie gauche
	if left, ok := condition["left"].(map[string]interface{}); ok {
		leftJoinConditions := extractJoinConditions(left)
		joinConditions = append(joinConditions, leftJoinConditions...)
	}

	// Extraire les opérations
	if operations, ok := condition["operations"].([]interface{}); ok {
		for _, op := range operations {
			if opMap, ok := op.(map[string]interface{}); ok {
				if right, ok := opMap["right"].(map[string]interface{}); ok {
					rightJoinConditions := extractJoinConditions(right)
					joinConditions = append(joinConditions, rightJoinConditions...)
				}
			}
		}
	}

	return joinConditions
}
