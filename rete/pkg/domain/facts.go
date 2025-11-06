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
