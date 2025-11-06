// Package domain contient les types métier fondamentaux du système RETE
package domain

import (
	"time"
)

// Fact représente un fait dans le système RETE.
// Un fait est une donnée d'entrée qui sera propagée à travers le réseau.
type Fact struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Fields    map[string]interface{} `json:"fields"`
	Timestamp time.Time              `json:"timestamp"`
}

// NewFact crée un nouveau fait avec timestamp automatique.
func NewFact(id, factType string, fields map[string]interface{}) *Fact {
	return &Fact{
		ID:        id,
		Type:      factType,
		Fields:    fields,
		Timestamp: time.Now(),
	}
}

// String retourne la représentation string d'un fait.
func (f *Fact) String() string {
	return f.ID + ":" + f.Type
}

// GetField retourne la valeur d'un champ et indique s'il existe.
func (f *Fact) GetField(fieldName string) (interface{}, bool) {
	value, exists := f.Fields[fieldName]
	return value, exists
}

// Token représente un token dans le réseau RETE.
// Un token contient un ou plusieurs faits et circule entre les nœuds.
type Token struct {
	ID     string  `json:"id"`
	Facts  []*Fact `json:"facts"`
	NodeID string  `json:"node_id"`
	Parent *Token  `json:"parent,omitempty"`
}

// NewToken crée un nouveau token.
func NewToken(id, nodeID string, facts []*Fact) *Token {
	return &Token{
		ID:     id,
		Facts:  facts,
		NodeID: nodeID,
	}
}

// WorkingMemory représente la mémoire de travail d'un nœud.
// Elle stocke les faits et tokens actuellement traités par le nœud.
type WorkingMemory struct {
	NodeID string            `json:"node_id"`
	Facts  map[string]*Fact  `json:"facts"`
	Tokens map[string]*Token `json:"tokens"`
}

// NewWorkingMemory crée une nouvelle mémoire de travail.
func NewWorkingMemory(nodeID string) *WorkingMemory {
	return &WorkingMemory{
		NodeID: nodeID,
		Facts:  make(map[string]*Fact),
		Tokens: make(map[string]*Token),
	}
}

// AddFact ajoute un fait à la mémoire.
func (wm *WorkingMemory) AddFact(fact *Fact) {
	if wm.Facts == nil {
		wm.Facts = make(map[string]*Fact)
	}
	wm.Facts[fact.ID] = fact
}

// RemoveFact supprime un fait de la mémoire.
func (wm *WorkingMemory) RemoveFact(factID string) {
	delete(wm.Facts, factID)
}

// GetFacts retourne tous les faits de la mémoire.
func (wm *WorkingMemory) GetFacts() []*Fact {
	facts := make([]*Fact, 0, len(wm.Facts))
	for _, fact := range wm.Facts {
		facts = append(facts, fact)
	}
	return facts
}

// AddToken ajoute un token à la mémoire.
func (wm *WorkingMemory) AddToken(token *Token) {
	if wm.Tokens == nil {
		wm.Tokens = make(map[string]*Token)
	}
	wm.Tokens[token.ID] = token
}

// RemoveToken supprime un token de la mémoire.
func (wm *WorkingMemory) RemoveToken(tokenID string) {
	delete(wm.Tokens, tokenID)
}

// GetTokens retourne tous les tokens de la mémoire.
func (wm *WorkingMemory) GetTokens() []*Token {
	tokens := make([]*Token, 0, len(wm.Tokens))
	for _, token := range wm.Tokens {
		tokens = append(tokens, token)
	}
	return tokens
}

// BasicJoinCondition implémente une condition de jointure simple.
type BasicJoinCondition struct {
	LeftField  string `json:"left_field"`
	RightField string `json:"right_field"`
	Operator   string `json:"operator"`
}

// NewBasicJoinCondition crée une nouvelle condition de jointure.
func NewBasicJoinCondition(leftField, rightField, operator string) *BasicJoinCondition {
	return &BasicJoinCondition{
		LeftField:  leftField,
		RightField: rightField,
		Operator:   operator,
	}
}

// Evaluate évalue la condition de jointure entre un token et un fait.
func (bjc *BasicJoinCondition) Evaluate(token *Token, fact *Fact) bool {
	if len(token.Facts) == 0 {
		return false
	}

	// Prendre le dernier fait du token pour la comparaison
	leftFact := token.Facts[len(token.Facts)-1]

	leftValue, leftExists := leftFact.GetField(bjc.LeftField)
	rightValue, rightExists := fact.GetField(bjc.RightField)

	if !leftExists || !rightExists {
		return false
	}

	return bjc.evaluateOperator(leftValue, rightValue)
}

// GetLeftField retourne le champ gauche de la condition.
func (bjc *BasicJoinCondition) GetLeftField() string {
	return bjc.LeftField
}

// GetRightField retourne le champ droit de la condition.
func (bjc *BasicJoinCondition) GetRightField() string {
	return bjc.RightField
}

// GetOperator retourne l'opérateur de la condition.
func (bjc *BasicJoinCondition) GetOperator() string {
	return bjc.Operator
}

// evaluateOperator évalue l'opérateur entre deux valeurs.
func (bjc *BasicJoinCondition) evaluateOperator(left, right interface{}) bool {
	switch bjc.Operator {
	case "==", "=":
		return left == right
	case "!=":
		return left != right
	case "<":
		return bjc.compareValues(left, right) < 0
	case "<=":
		return bjc.compareValues(left, right) <= 0
	case ">":
		return bjc.compareValues(left, right) > 0
	case ">=":
		return bjc.compareValues(left, right) >= 0
	default:
		return false
	}
}

// compareValues compare deux valeurs pour les opérateurs relationnels.
func (bjc *BasicJoinCondition) compareValues(left, right interface{}) int {
	// Implémentation simplifiée pour les types de base
	switch lv := left.(type) {
	case int:
		if rv, ok := right.(int); ok {
			if lv < rv {
				return -1
			} else if lv > rv {
				return 1
			}
			return 0
		}
	case float64:
		if rv, ok := right.(float64); ok {
			if lv < rv {
				return -1
			} else if lv > rv {
				return 1
			}
			return 0
		}
	case string:
		if rv, ok := right.(string); ok {
			if lv < rv {
				return -1
			} else if lv > rv {
				return 1
			}
			return 0
		}
	}
	return 0
}
