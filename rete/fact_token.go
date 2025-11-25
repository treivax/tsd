package rete

import (
	"fmt"
	"time"
)

type Fact struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Fields    map[string]interface{} `json:"fields"`
	Timestamp time.Time              `json:"timestamp"`
}

// String retourne la représentation string d'un fait
func (f *Fact) String() string {
	return fmt.Sprintf("Fact{ID:%s, Type:%s, Fields:%v}", f.ID, f.Type, f.Fields)
}

// GetField retourne la valeur d'un champ
func (f *Fact) GetField(fieldName string) (interface{}, bool) {
	value, exists := f.Fields[fieldName]
	return value, exists
}

// Token représente un token dans le réseau RETE
type Token struct {
	ID           string           `json:"id"`
	Facts        []*Fact          `json:"facts"`
	NodeID       string           `json:"node_id"`
	Parent       *Token           `json:"parent,omitempty"`
	Bindings     map[string]*Fact `json:"bindings,omitempty"`       // Nouveau: bindings pour jointures
	IsJoinResult bool             `json:"is_join_result,omitempty"` // Indique si c'est un token de jointure réussie
}

// WorkingMemory représente la mémoire de travail d'un nœud
type WorkingMemory struct {
	NodeID string            `json:"node_id"`
	Facts  map[string]*Fact  `json:"facts"`
	Tokens map[string]*Token `json:"tokens"`
}

// AddFact ajoute un fait à la mémoire
func (wm *WorkingMemory) AddFact(fact *Fact) {
	if wm.Facts == nil {
		wm.Facts = make(map[string]*Fact)
	}
	wm.Facts[fact.ID] = fact
}

// RemoveFact supprime un fait de la mémoire
func (wm *WorkingMemory) RemoveFact(factID string) {
	delete(wm.Facts, factID)
}

// GetFact récupère un fait par son ID
func (wm *WorkingMemory) GetFact(factID string) (*Fact, bool) {
	fact, exists := wm.Facts[factID]
	return fact, exists
}

// GetFacts retourne tous les faits de la mémoire
func (wm *WorkingMemory) GetFacts() []*Fact {
	facts := make([]*Fact, 0, len(wm.Facts))
	for _, fact := range wm.Facts {
		facts = append(facts, fact)
	}
	return facts
}

// AddToken ajoute un token à la mémoire
func (wm *WorkingMemory) AddToken(token *Token) {
	if wm.Tokens == nil {
		wm.Tokens = make(map[string]*Token)
	}
	wm.Tokens[token.ID] = token
}

// RemoveToken supprime un token de la mémoire
func (wm *WorkingMemory) RemoveToken(tokenID string) {
	delete(wm.Tokens, tokenID)
}

// GetTokens retourne tous les tokens de la mémoire
func (wm *WorkingMemory) GetTokens() []*Token {
	tokens := make([]*Token, 0, len(wm.Tokens))
	for _, token := range wm.Tokens {
		tokens = append(tokens, token)
	}
	return tokens
}

// GetFactsByVariable retourne les faits associés aux variables spécifiées
func (wm *WorkingMemory) GetFactsByVariable(variables []string) []*Fact {
	// Pour l'instant, retourne tous les faits (implémentation simplifiée)
	return wm.GetFacts()
}

// GetTokensByVariable retourne les tokens associés aux variables spécifiées
func (wm *WorkingMemory) GetTokensByVariable(variables []string) []*Token {
	// Pour l'instant, retourne tous les tokens (implémentation simplifiée)
	return wm.GetTokens()
}
