// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"strconv"
	"sync"
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

	return &JoinNode{
		BaseNode: BaseNode{
			ID:       nodeID,
			Type:     "join",
			Memory:   &WorkingMemory{NodeID: nodeID, Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
			Children: make([]Node, 0),
			Storage:  storage,
		},
		Condition:      condition,
		LeftVariables:  leftVars,
		RightVariables: rightVars,
		AllVariables:   allVars,
		VariableTypes:  varTypes,
		JoinConditions: extractJoinConditions(condition),
		// Initialiser les mémoires séparées
		LeftMemory:   &WorkingMemory{NodeID: nodeID + "_left", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
		RightMemory:  &WorkingMemory{NodeID: nodeID + "_right", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
		ResultMemory: &WorkingMemory{NodeID: nodeID + "_result", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
	}
}

// ActivateLeft traite les tokens de la gauche (généralement des AlphaNodes)
func (jn *JoinNode) ActivateLeft(token *Token) error {
	// Stocker le token dans la mémoire gauche
	jn.mutex.Lock()
	jn.LeftMemory.AddToken(token)
	jn.mutex.Unlock()

	// Essayer de joindre avec tous les tokens de la mémoire droite
	rightTokens := jn.RightMemory.GetTokens()

	for _, rightToken := range rightTokens {
		if joinedToken := jn.performJoinWithTokens(token, rightToken); joinedToken != nil {

			// Stocker uniquement les tokens de jointure réussie
			joinedToken.IsJoinResult = true
			jn.mutex.Lock()
			jn.ResultMemory.AddToken(joinedToken)
			jn.Memory.AddToken(joinedToken) // Pour compatibilité avec le comptage
			jn.mutex.Unlock()

			if err := jn.PropagateToChildren(nil, joinedToken); err != nil {
				return err
			}
		}
	}
	return nil
}

// ActivateRetract retrait des tokens contenant le fait rétracté des 3 mémoires
// factID doit être l'identifiant interne (Type_ID)
func (jn *JoinNode) ActivateRetract(factID string) error {
	jn.mutex.Lock()
	var leftTokensToRemove []string
	for tokenID, token := range jn.LeftMemory.Tokens {
		for _, fact := range token.Facts {
			if fact.GetInternalID() == factID {
				leftTokensToRemove = append(leftTokensToRemove, tokenID)
				break
			}
		}
	}
	for _, tokenID := range leftTokensToRemove {
		delete(jn.LeftMemory.Tokens, tokenID)
	}
	var rightTokensToRemove []string
	for tokenID, token := range jn.RightMemory.Tokens {
		for _, fact := range token.Facts {
			if fact.GetInternalID() == factID {
				rightTokensToRemove = append(rightTokensToRemove, tokenID)
				break
			}
		}
	}
	for _, tokenID := range rightTokensToRemove {
		delete(jn.RightMemory.Tokens, tokenID)
	}
	var resultTokensToRemove []string
	for tokenID, token := range jn.ResultMemory.Tokens {
		for _, fact := range token.Facts {
			if fact.GetInternalID() == factID {
				resultTokensToRemove = append(resultTokensToRemove, tokenID)
				break
			}
		}
	}
	for _, tokenID := range resultTokensToRemove {
		delete(jn.ResultMemory.Tokens, tokenID)
		delete(jn.Memory.Tokens, tokenID)
	}
	jn.mutex.Unlock()
	totalRemoved := len(leftTokensToRemove) + len(rightTokensToRemove) + len(resultTokensToRemove)
	if totalRemoved > 0 {
		fmt.Printf("🗑️  [JOIN_%s] Rétractation: %d tokens retirés (L:%d R:%d RES:%d)\n", jn.ID, totalRemoved, len(leftTokensToRemove), len(rightTokensToRemove), len(resultTokensToRemove))
	}
	return jn.PropagateRetractToChildren(factID)
}

// ActivateRight traite les faits de la droite (nouveau fait injecté via AlphaNode)
func (jn *JoinNode) ActivateRight(fact *Fact) error {
	// Convertir le fait en token pour la mémoire droite
	factVar := jn.getVariableForFact(fact)
	if factVar == "" {
		return nil // Fait non applicable à ce JoinNode
	}

	factToken := &Token{
		ID:       fmt.Sprintf("right_token_%s_%s", jn.ID, fact.ID),
		Facts:    []*Fact{fact},
		NodeID:   jn.ID,
		Bindings: map[string]*Fact{factVar: fact},
	}

	// Stocker le token dans la mémoire droite
	jn.mutex.Lock()
	jn.RightMemory.AddToken(factToken)
	jn.mutex.Unlock()

	// Essayer de joindre avec tous les tokens de la mémoire gauche
	leftTokens := jn.LeftMemory.GetTokens()

	for _, leftToken := range leftTokens {
		if joinedToken := jn.performJoinWithTokens(leftToken, factToken); joinedToken != nil {

			// Stocker uniquement les tokens de jointure réussie
			joinedToken.IsJoinResult = true
			jn.mutex.Lock()
			jn.ResultMemory.AddToken(joinedToken)
			jn.Memory.AddToken(joinedToken) // Pour compatibilité avec le comptage
			jn.mutex.Unlock()

			if err := jn.PropagateToChildren(nil, joinedToken); err != nil {
				return err
			}
		}
	}
	return nil
}

// performJoinWithTokens effectue la jointure entre deux tokens
func (jn *JoinNode) performJoinWithTokens(token1 *Token, token2 *Token) *Token {
	// Vérifier que les tokens ont des variables différentes
	if !jn.tokensHaveDifferentVariables(token1, token2) {
		return nil
	}

	// Combiner les bindings des deux tokens
	combinedBindings := make(map[string]*Fact)

	// Copier les bindings du premier token
	for varName, varFact := range token1.Bindings {
		combinedBindings[varName] = varFact
	}

	// Copier les bindings du second token
	for varName, varFact := range token2.Bindings {
		combinedBindings[varName] = varFact
	}

	// Valider les conditions de jointure
	if !jn.evaluateJoinConditions(combinedBindings) {
		return nil // Jointure échoue
	}

	// Créer et retourner le token joint
	return &Token{
		ID:       fmt.Sprintf("%s_JOIN_%s", token1.ID, token2.ID),
		Bindings: combinedBindings,
		NodeID:   jn.ID,
		Facts:    append(token1.Facts, token2.Facts...),
	}
}

// tokensHaveDifferentVariables vérifie que les tokens représentent des variables différentes
func (jn *JoinNode) tokensHaveDifferentVariables(token1 *Token, token2 *Token) bool {
	for var1 := range token1.Bindings {
		for var2 := range token2.Bindings {
			if var1 == var2 {
				return false // Même variable = pas de jointure possible
			}
		}
	}
	return true
}

// getVariableForFact détermine la variable associée à un fait basé sur son type
func (jn *JoinNode) getVariableForFact(fact *Fact) string {
	// Utiliser le mapping variable -> type du JoinNode
	for _, varName := range jn.AllVariables {
		if expectedType, exists := jn.VariableTypes[varName]; exists {
			if expectedType == fact.Type {
				return varName
			}
		}
	}

	fmt.Printf("❌ JOINNODE[%s]: Aucune variable trouvée pour fait %s (type: %s)\n", jn.ID, fact.ID, fact.Type)
	fmt.Printf("   Variables disponibles: %v\n", jn.AllVariables)
	fmt.Printf("   Types attendus: %v\n", jn.VariableTypes)
	return ""
}

// evaluateJoinConditions vérifie si toutes les conditions de jointure sont respectées
func (jn *JoinNode) evaluateJoinConditions(bindings map[string]*Fact) bool {
	// Vérifier qu'on a au moins 2 variables différentes
	if len(bindings) < 2 {
		return false
	}

	// Étape 1: Évaluer les conditions de jointure simples (field-to-field)
	// Ces conditions sont extraites et toujours présentes pour les joins
	if len(jn.JoinConditions) > 0 {
		if !jn.evaluateSimpleJoinConditions(bindings) {
			return false
		}
	}

	// Étape 2: Si on a une condition complète avec des contraintes alpha additionnelles,
	// l'évaluer pour vérifier les conditions non-join (ex: o.amount > 100)
	if jn.Condition != nil {
		// Unwrap the constraint wrapper if present
		actualCondition := jn.Condition
		if condType, exists := jn.Condition["type"].(string); exists && condType == "constraint" {
			if constraint, ok := jn.Condition["constraint"].(map[string]interface{}); ok {
				actualCondition = constraint
			}
		}

		condType, exists := actualCondition["type"].(string)

		// Si c'est une simple comparison (join pur), on a déjà validé avec JoinConditions
		if exists && condType == "comparison" {
			// Déjà validé par evaluateSimpleJoinConditions
			return true
		}

		// Si c'est un logicalExpr, extraire et évaluer seulement les contraintes alpha (non-join)
		if exists && condType == "logicalExpr" {
			alphaConditions := jn.extractAlphaConditions(actualCondition)
			if len(alphaConditions) == 0 {
				// Pas de contraintes alpha, seulement des joins déjà validés
				return true
			}

			// Évaluer chaque contrainte alpha
			evaluator := NewAlphaConditionEvaluator()
			evaluator.SetPartialEvalMode(true)

			// Lier toutes les variables aux faits
			for varName, fact := range bindings {
				evaluator.variableBindings[varName] = fact
			}

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
	}

	return true
}

// extractAlphaConditions extrait les conditions alpha (non-join) d'une logicalExpr
func (jn *JoinNode) extractAlphaConditions(condition map[string]interface{}) []map[string]interface{} {
	var alphaConditions []map[string]interface{}

	// Vérifier la partie gauche
	if left, ok := condition["left"].(map[string]interface{}); ok {
		if isAlphaCondition(left) {
			alphaConditions = append(alphaConditions, left)
		}
	}

	// Vérifier les opérations
	if operations, ok := condition["operations"].([]interface{}); ok {
		for _, op := range operations {
			if opMap, ok := op.(map[string]interface{}); ok {
				if right, ok := opMap["right"].(map[string]interface{}); ok {
					if isAlphaCondition(right) {
						alphaConditions = append(alphaConditions, right)
					}
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

// evaluateSimpleJoinConditions évalue les conditions de jointure simples (champ à champ)
func (jn *JoinNode) evaluateSimpleJoinConditions(bindings map[string]*Fact) bool {
	for _, joinCondition := range jn.JoinConditions {
		leftFact := bindings[joinCondition.LeftVar]
		rightFact := bindings[joinCondition.RightVar]

		if leftFact == nil || rightFact == nil {
			return false // Une variable manque
		}

		// Récupérer les valeurs des champs
		leftValue := leftFact.Fields[joinCondition.LeftField]
		rightValue := rightFact.Fields[joinCondition.RightField]

		// Évaluer l'opérateur
		switch joinCondition.Operator {
		case "==":
			if leftValue != rightValue {
				return false
			}
		case "!=":
			if leftValue == rightValue {
				return false
			}
		case "<":
			if leftFloat, leftOk := convertToFloat64(leftValue); leftOk {
				if rightFloat, rightOk := convertToFloat64(rightValue); rightOk {
					if leftFloat >= rightFloat {
						return false
					}
				} else {
					return false // Comparaison numérique impossible
				}
			} else {
				return false
			}
		case ">":
			if leftFloat, leftOk := convertToFloat64(leftValue); leftOk {
				if rightFloat, rightOk := convertToFloat64(rightValue); rightOk {
					if leftFloat <= rightFloat {
						return false
					}
				} else {
					return false
				}
			} else {
				return false
			}
		case "<=":
			if leftFloat, leftOk := convertToFloat64(leftValue); leftOk {
				if rightFloat, rightOk := convertToFloat64(rightValue); rightOk {
					if leftFloat > rightFloat {
						return false
					}
				} else {
					return false
				}
			} else {
				return false
			}
		case ">=":
			if leftFloat, leftOk := convertToFloat64(leftValue); leftOk {
				if rightFloat, rightOk := convertToFloat64(rightValue); rightOk {
					if leftFloat < rightFloat {
						return false
					}
				} else {
					return false
				}
			} else {
				return false
			}
		default:
			return false
		}
	}

	return true
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

// extractJoinConditions extrait les conditions de jointure d'une condition complexe
func extractJoinConditions(condition map[string]interface{}) []JoinCondition {
	for key, value := range condition {
		fmt.Printf("    %s: %v (type: %T)\n", key, value, value)
	}

	var joinConditions []JoinCondition

	// Cas 1: condition wrappée dans un type "constraint"
	if conditionType, exists := condition["type"].(string); exists && conditionType == "constraint" {
		if innerCondition, ok := condition["constraint"].(map[string]interface{}); ok {
			fmt.Printf("  ✅ Sous-condition extraite, analyse récursive\n")
			return extractJoinConditions(innerCondition)
		}
	}

	// Cas 2: condition EXISTS avec array de conditions
	if conditionType, exists := condition["type"].(string); exists && conditionType == "exists" {
		if conditionsData, ok := condition["conditions"].([]map[string]interface{}); ok {
			fmt.Printf("  ✅ Array de conditions EXISTS trouvé: %d conditions\n", len(conditionsData))
			for _, subCondition := range conditionsData {
				subJoinConditions := extractJoinConditions(subCondition)
				joinConditions = append(joinConditions, subJoinConditions...)
			}
			return joinConditions
		}
	}

	// Cas 3: condition directe de comparaison
	if conditionType, exists := condition["type"].(string); exists && conditionType == "comparison" {
		fmt.Printf("  ✅ Condition de comparaison détectée\n")
		if left, leftOk := condition["left"].(map[string]interface{}); leftOk {
			if right, rightOk := condition["right"].(map[string]interface{}); rightOk {
				fmt.Printf("  ✅ Left et Right extraits\n")
				if leftType, _ := left["type"].(string); leftType == "fieldAccess" {
					if rightType, _ := right["type"].(string); rightType == "fieldAccess" {
						// Condition de jointure détectée
						fmt.Printf("  ✅ Condition de jointure fieldAccess détectée\n")
						leftObj, _ := left["object"].(string)
						leftField, _ := left["field"].(string)
						rightObj, _ := right["object"].(string)
						rightField, _ := right["field"].(string)
						operator, _ := condition["operator"].(string)

						fmt.Printf("    📌 %s.%s %s %s.%s\n", leftObj, leftField, operator, rightObj, rightField)

						joinConditions = append(joinConditions, JoinCondition{
							LeftField:  leftField,
							RightField: rightField,
							LeftVar:    leftObj,
							RightVar:   rightObj,
							Operator:   operator,
						})
					}
				}
			}
		}
	}

	// Cas 4: logicalExpr avec opérations AND/OR
	if conditionType, exists := condition["type"].(string); exists && conditionType == "logicalExpr" {
		fmt.Printf("  ✅ LogicalExpr détectée, extraction des conditions\n")

		// Extraire les conditions de la partie gauche
		if left, ok := condition["left"].(map[string]interface{}); ok {
			leftJoinConditions := extractJoinConditions(left)
			joinConditions = append(joinConditions, leftJoinConditions...)
		}

		// Extraire les conditions des opérations
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

	return joinConditions
}
