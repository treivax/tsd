// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"sync"
)

type ExistsNode struct {
	BaseNode
	Condition       map[string]interface{} `json:"condition"`
	MainVariable    string                 `json:"main_variable"`    // Variable principale (p)
	ExistsVariable  string                 `json:"exists_variable"`  // Variable d'existence (o)
	VariableTypes   map[string]string      `json:"variable_types"`   // Mapping variable -> type
	ExistsCondition []JoinCondition        `json:"exists_condition"` // Condition d'existence (o.customer_id == p.id)
	mutex           sync.RWMutex
	// Mémoires pour architecture RETE
	MainMemory   *WorkingMemory // Faits de la variable principale
	ExistsMemory *WorkingMemory // Faits pour vérification d'existence
	ResultMemory *WorkingMemory // Tokens avec existence vérifiée
}

// NewExistsNode crée un nouveau nœud d'existence
func NewExistsNode(nodeID string, condition map[string]interface{}, mainVar string, existsVar string, varTypes map[string]string, storage Storage) *ExistsNode {
	return &ExistsNode{
		BaseNode: BaseNode{
			ID:       nodeID,
			Type:     "exists",
			Memory:   &WorkingMemory{NodeID: nodeID, Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
			Children: make([]Node, 0),
			Storage:  storage,
		},
		Condition:       condition,
		MainVariable:    mainVar,
		ExistsVariable:  existsVar,
		VariableTypes:   varTypes,
		ExistsCondition: extractJoinConditions(condition),
		// Initialiser les mémoires séparées
		MainMemory:   &WorkingMemory{NodeID: nodeID + "_main", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
		ExistsMemory: &WorkingMemory{NodeID: nodeID + "_exists", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
		ResultMemory: &WorkingMemory{NodeID: nodeID + "_result", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
	}
}

// ActivateLeft traite les faits de la variable principale
func (en *ExistsNode) ActivateLeft(token *Token) error {

	// Stocker le token dans la mémoire principale
	en.mutex.Lock()
	en.MainMemory.AddToken(token)
	en.mutex.Unlock()

	// Vérifier s'il existe des faits correspondants
	if en.checkExistence(token) {

		// Stocker le token avec existence vérifiée
		token.IsJoinResult = true // Marquer comme résultat validé
		en.mutex.Lock()
		en.ResultMemory.AddToken(token)
		en.Memory.AddToken(token) // Pour compatibilité avec le comptage
		en.mutex.Unlock()

		// Propager le token
		if err := en.PropagateToChildren(nil, token); err != nil {
			return err
		}
	} else {
	}

	return nil
}

// ActivateRetract retrait des tokens et faits contenant le fait rétracté
// factID doit être l'identifiant interne (Type_ID)
func (en *ExistsNode) ActivateRetract(factID string) error {
	en.mutex.Lock()
	var mainTokensToRemove []string
	for tokenID, token := range en.MainMemory.Tokens {
		for _, fact := range token.Facts {
			if fact.GetInternalID() == factID {
				mainTokensToRemove = append(mainTokensToRemove, tokenID)
				break
			}
		}
	}
	for _, tokenID := range mainTokensToRemove {
		delete(en.MainMemory.Tokens, tokenID)
	}
	_, existsInExistsMemory := en.ExistsMemory.GetFact(factID)
	if existsInExistsMemory {
		en.ExistsMemory.RemoveFact(factID)
	}
	var resultTokensToRemove []string
	for tokenID, token := range en.ResultMemory.Tokens {
		for _, fact := range token.Facts {
			if fact.GetInternalID() == factID {
				resultTokensToRemove = append(resultTokensToRemove, tokenID)
				break
			}
		}
	}
	for _, tokenID := range resultTokensToRemove {
		delete(en.ResultMemory.Tokens, tokenID)
		delete(en.Memory.Tokens, tokenID)
	}
	en.mutex.Unlock()
	totalRemoved := len(mainTokensToRemove) + len(resultTokensToRemove)
	if existsInExistsMemory {
		totalRemoved++
	}
	if totalRemoved > 0 {
		fmt.Printf("🗑️  [EXISTS_%s] Rétractation: %d éléments retirés (MAIN:%d EXISTS:%v RES:%d)\n", en.ID, totalRemoved, len(mainTokensToRemove), existsInExistsMemory, len(resultTokensToRemove))
	}
	return en.PropagateRetractToChildren(factID)
}

// ActivateRight traite les faits pour vérification d'existence
func (en *ExistsNode) ActivateRight(fact *Fact) error {

	// Stocker le fait dans la mémoire d'existence
	en.mutex.Lock()
	if err := en.ExistsMemory.AddFact(fact); err != nil {
		en.mutex.Unlock()
		return fmt.Errorf("erreur ajout fait dans exists node: %w", err)
	}
	en.mutex.Unlock()

	// Re-vérifier tous les tokens principaux avec ce nouveau fait
	mainTokens := en.MainMemory.GetTokens()
	for _, mainToken := range mainTokens {
		if en.checkExistence(mainToken) && !en.isAlreadyValidated(mainToken) {

			// Stocker le token avec existence vérifiée
			validatedToken := &Token{
				ID:           mainToken.ID + "_validated",
				Facts:        mainToken.Facts,
				NodeID:       en.ID,
				Bindings:     mainToken.Bindings,
				IsJoinResult: true,
			}

			en.mutex.Lock()
			en.ResultMemory.AddToken(validatedToken)
			en.Memory.AddToken(validatedToken)
			en.mutex.Unlock()

			// Propager le token
			if err := en.PropagateToChildren(nil, validatedToken); err != nil {
				return err
			}
		}
	}

	return nil
}

// checkExistence vérifie si un token principal a des faits correspondants
func (en *ExistsNode) checkExistence(mainToken *Token) bool {
	existsFacts := en.ExistsMemory.GetFacts()

	// Récupérer le fait principal du token
	if len(mainToken.Facts) == 0 {
		return false
	}
	mainFact := mainToken.Facts[0]

	// Vérifier les conditions d'existence
	for _, existsFact := range existsFacts {
		if en.evaluateExistsCondition(mainFact, existsFact) {
			return true
		}
	}

	return false
}

// evaluateExistsCondition évalue la condition d'existence entre deux faits
func (en *ExistsNode) evaluateExistsCondition(mainFact *Fact, existsFact *Fact) bool {

	for i, condition := range en.ExistsCondition {
		fmt.Printf("    Condition %d: %s.%s %s %s.%s\n", i,
			condition.LeftVar, condition.LeftField, condition.Operator,
			condition.RightVar, condition.RightField)

		// Déterminer quel fait correspond à quelle variable
		var leftFact, rightFact *Fact

		if condition.LeftVar == en.MainVariable {
			leftFact = mainFact
			rightFact = existsFact
			fmt.Printf("    → MainFact comme LeftVar (%s), ExistsFact comme RightVar (%s)\n", condition.LeftVar, condition.RightVar)
		} else if condition.LeftVar == en.ExistsVariable {
			leftFact = existsFact
			rightFact = mainFact
			fmt.Printf("    → ExistsFact comme LeftVar (%s), MainFact comme RightVar (%s)\n", condition.LeftVar, condition.RightVar)
		} else {
			fmt.Printf("    ❌ Variable %s non trouvée dans MainVariable:%s ou ExistsVariable:%s\n", condition.LeftVar, en.MainVariable, en.ExistsVariable)
			continue
		}

		leftValue := leftFact.Fields[condition.LeftField]
		rightValue := rightFact.Fields[condition.RightField]

		switch condition.Operator {
		case "==":
			if leftValue != rightValue {
				fmt.Printf("    ❌ Condition %d échoue: %v != %v\n", i, leftValue, rightValue)
				return false
			}
			fmt.Printf("    ✅ Condition %d réussie: %v == %v\n", i, leftValue, rightValue)
		case "!=":
			if leftValue == rightValue {
				fmt.Printf("    ❌ Condition %d échoue: %v == %v\n", i, leftValue, rightValue)
				return false
			}
			fmt.Printf("    ✅ Condition %d réussie: %v != %v\n", i, leftValue, rightValue)
		default:
			fmt.Printf("    ❌ Opérateur non supporté: %s\n", condition.Operator)
			return false
		}
	}

	fmt.Printf("  ✅ Toutes les conditions EXISTS satisfaites\n")
	return true
}

// isAlreadyValidated vérifie si un token a déjà été validé
func (en *ExistsNode) isAlreadyValidated(token *Token) bool {
	validatedTokens := en.ResultMemory.GetTokens()
	for _, validatedToken := range validatedTokens {
		if validatedToken.ID == token.ID+"_validated" || validatedToken.ID == token.ID {
			return true
		}
	}
	return false
}
