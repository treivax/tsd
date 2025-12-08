// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"

	"github.com/treivax/tsd/rete/pkg/domain"
)

// Fact est un alias vers domain.Fact pour compatibilité avec le code existant.
// Toutes les méthodes et fonctions sont définies dans rete/pkg/domain/facts.go.
type Fact = domain.Fact

// Réexporter les fonctions helper pour compatibilité
var (
	MakeInternalID  = domain.MakeInternalID
	ParseInternalID = domain.ParseInternalID
)

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

// AddFact ajoute un fait à la mémoire en utilisant un identifiant interne unique (Type_ID)
// Retourne une erreur si un fait avec le même type et ID existe déjà
func (wm *WorkingMemory) AddFact(fact *Fact) error {
	if wm.Facts == nil {
		wm.Facts = make(map[string]*Fact)
	}

	// Utiliser l'identifiant interne (Type_ID) pour garantir l'unicité par type
	internalID := fact.GetInternalID()

	if _, exists := wm.Facts[internalID]; exists {
		return fmt.Errorf("fait avec ID '%s' et type '%s' existe déjà dans la mémoire", fact.ID, fact.Type)
	}

	wm.Facts[internalID] = fact
	return nil
}

// RemoveFact supprime un fait de la mémoire
// factID doit être l'identifiant interne (Type_ID)
func (wm *WorkingMemory) RemoveFact(factID string) {
	delete(wm.Facts, factID)
}

// GetFact récupère un fait par son identifiant interne (Type_ID)
// Pour rechercher par type et ID séparément, utiliser GetFactByTypeAndID
func (wm *WorkingMemory) GetFact(internalID string) (*Fact, bool) {
	fact, exists := wm.Facts[internalID]
	return fact, exists
}

// GetFactByInternalID récupère un fait uniquement par son identifiant interne
func (wm *WorkingMemory) GetFactByInternalID(internalID string) (*Fact, bool) {
	fact, exists := wm.Facts[internalID]
	return fact, exists
}

// GetFactByTypeAndID récupère un fait par son type et son ID
func (wm *WorkingMemory) GetFactByTypeAndID(factType, factID string) (*Fact, bool) {
	internalID := MakeInternalID(factType, factID)
	return wm.GetFactByInternalID(internalID)
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

// Clone crée une copie profonde de WorkingMemory
func (wm *WorkingMemory) Clone() *WorkingMemory {
	clone := &WorkingMemory{
		NodeID: wm.NodeID,
		Facts:  make(map[string]*Fact),
		Tokens: make(map[string]*Token),
	}

	// Copier les faits
	for id, fact := range wm.Facts {
		clone.Facts[id] = fact.Clone()
	}

	// Copier les tokens
	for id, token := range wm.Tokens {
		clone.Tokens[id] = token.Clone()
	}

	return clone
}

// Clone crée une copie profonde d'un token
func (t *Token) Clone() *Token {
	clone := &Token{
		ID:           t.ID,
		Facts:        make([]*Fact, len(t.Facts)),
		NodeID:       t.NodeID,
		Bindings:     make(map[string]*Fact),
		IsJoinResult: t.IsJoinResult,
	}

	// Copier les faits
	for i, fact := range t.Facts {
		clone.Facts[i] = fact.Clone()
	}

	// Copier les bindings
	for k, v := range t.Bindings {
		clone.Bindings[k] = v.Clone()
	}

	// Note: Parent n'est pas cloné pour éviter récursion infinie

	return clone
}
