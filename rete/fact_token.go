// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"sync/atomic"
)

// FieldNameID est le nom du champ spécial pour l'identifiant du fait.
// Ce champ est accessible dans les expressions mais stocké dans Fact.ID, pas dans Fact.Fields.
const FieldNameID = "id"

// Fact représente un fait dans le réseau RETE
type Fact struct {
	// ID est l'identifiant unique du fait.
	// Il est soit généré à partir des clés primaires, soit calculé comme hash.
	// Format: "TypeName~value1_value2..." ou "TypeName~<hash>"
	// Accessible dans les expressions via le champ spécial 'id'.
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Fields     map[string]interface{} `json:"fields"`
	Attributes map[string]interface{} `json:"attributes,omitempty"` // Alias pour Fields (compatibilité)
}

// String retourne la représentation string d'un fait
func (f *Fact) String() string {
	return fmt.Sprintf("Fact{ID:%s, Type:%s, Fields:%v}", f.ID, f.Type, f.Fields)
}

// GetInternalID retourne l'identifiant interne unique (Type_ID)
func (f *Fact) GetInternalID() string {
	return fmt.Sprintf("%s_%s", f.Type, f.ID)
}

// GetField retourne la valeur d'un champ
func (f *Fact) GetField(fieldName string) (interface{}, bool) {
	value, exists := f.Fields[fieldName]
	return value, exists
}

// Clone crée une copie profonde d'un fait
func (f *Fact) Clone() *Fact {
	clone := &Fact{
		ID:     f.ID,
		Type:   f.Type,
		Fields: make(map[string]interface{}),
	}
	for k, v := range f.Fields {
		clone.Fields[k] = v
	}
	if f.Attributes != nil {
		clone.Attributes = make(map[string]interface{})
		for k, v := range f.Attributes {
			clone.Attributes[k] = v
		}
	}
	return clone
}

// MakeInternalID construit un identifiant interne à partir d'un type et d'un ID
func MakeInternalID(factType, factID string) string {
	return fmt.Sprintf("%s_%s", factType, factID)
}

// ParseInternalID décompose un identifiant interne en type et ID
// Retourne (type, id, true) si le format est valide, sinon ("", "", false)
func ParseInternalID(internalID string) (string, string, bool) {
	for i := 0; i < len(internalID); i++ {
		if internalID[i] == '_' {
			return internalID[:i], internalID[i+1:], true
		}
	}
	return "", "", false
}

// TokenMetadata contient les métadonnées d'un token pour traçage et debug.
type TokenMetadata struct {
	CreatedAt    string   `json:"created_at,omitempty"`    // Timestamp de création
	CreatedBy    string   `json:"created_by,omitempty"`    // ID du nœud créateur
	JoinLevel    int      `json:"join_level,omitempty"`    // Niveau de jointure (0 = fact initial, 1+ = jointures)
	ParentTokens []string `json:"parent_tokens,omitempty"` // IDs des tokens parents (pour jointures)
}

// Token représente un token dans le réseau RETE avec bindings immuables.
//
// Changement majeur: Bindings utilise maintenant BindingChain au lieu de map[string]*Fact
// pour garantir l'immutabilité et éviter la perte de bindings lors des jointures en cascade.
type Token struct {
	ID           string        `json:"id"`
	Facts        []*Fact       `json:"facts"`
	NodeID       string        `json:"node_id"`
	Parent       *Token        `json:"parent,omitempty"`
	Bindings     *BindingChain `json:"-"`                        // Chaîne immuable de bindings (non sérialisable)
	IsJoinResult bool          `json:"is_join_result,omitempty"` // Indique si c'est un token de jointure réussie
	Metadata     TokenMetadata `json:"metadata,omitempty"`       // Métadonnées pour traçage
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

	if existingFact, exists := wm.Facts[internalID]; exists {
		return fmt.Errorf("fait avec ID '%s' et type '%s' existe déjà dans la mémoire du nœud %s (champs existants: %v)",
			fact.ID, fact.Type, wm.NodeID, existingFact.Fields)
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

// GetFactsByVariable retourne les faits associés aux variables spécifiées.
// Si variables est vide ou nil, retourne tous les faits.
//
// Note: Cette implémentation retourne tous les faits car WorkingMemory ne maintient
// pas d'index par variable. L'appelant doit filtrer lui-même si nécessaire en utilisant
// les bindings des tokens.
func (wm *WorkingMemory) GetFactsByVariable(variables []string) []*Fact {
	// Si pas de filtre, retourner tous les faits
	if len(variables) == 0 {
		return wm.GetFacts()
	}

	// TODO: Implémentation complète nécessiterait:
	// 1. Soit maintenir un index variable->fait dans WorkingMemory
	// 2. Soit parcourir les tokens et extraire les faits liés aux variables
	//
	// Pour l'instant, retourne tous les faits car WorkingMemory ne contient
	// pas suffisamment d'information pour filtrer par variable.
	// L'appelant devra filtrer lui-même si nécessaire.
	return wm.GetFacts()
}

// GetTokensByVariable retourne les tokens contenant au moins une des variables spécifiées.
// Si variables est vide ou nil, retourne tous les tokens.
//
// Le filtrage est basé sur Token.Bindings.Has() pour vérifier la présence de chaque variable.
func (wm *WorkingMemory) GetTokensByVariable(variables []string) []*Token {
	// Si pas de filtre, retourner tous les tokens
	if len(variables) == 0 {
		return wm.GetTokens()
	}

	// Filtrer les tokens qui contiennent au moins une des variables
	result := make([]*Token, 0)
	for _, token := range wm.Tokens {
		if token.Bindings != nil {
			for _, varName := range variables {
				if token.Bindings.Has(varName) {
					result = append(result, token)
					break // Token déjà ajouté, passer au suivant
				}
			}
		}
	}

	return result
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

// Clone crée une copie profonde d'un token.
//
// Note: BindingChain est immuable donc pas besoin de cloner la chaîne elle-même.
// On réutilise la même référence (partage structurel).
func (t *Token) Clone() *Token {
	clone := &Token{
		ID:           t.ID,
		Facts:        make([]*Fact, len(t.Facts)),
		NodeID:       t.NodeID,
		Bindings:     t.Bindings, // Immuable, pas besoin de cloner
		IsJoinResult: t.IsJoinResult,
		Metadata:     t.Metadata, // Copie de la structure
	}

	// Copier les faits
	for i, fact := range t.Facts {
		clone.Facts[i] = fact.Clone()
	}

	// Copier les ParentTokens si présents
	if len(t.Metadata.ParentTokens) > 0 {
		clone.Metadata.ParentTokens = make([]string, len(t.Metadata.ParentTokens))
		copy(clone.Metadata.ParentTokens, t.Metadata.ParentTokens)
	}

	// Note: Parent n'est pas cloné pour éviter récursion infinie
	// Note: Bindings n'est pas cloné car BindingChain est immuable

	return clone
}

// GetBinding retourne le fait lié à une variable.
//
// Wrapper autour de Bindings.Get() pour compatibilité.
//
// Paramètres:
//   - variable: nom de la variable
//
// Retourne:
//   - *Fact: pointeur vers le fait si trouvé, nil sinon
func (t *Token) GetBinding(variable string) *Fact {
	if t.Bindings == nil {
		return nil
	}
	return t.Bindings.Get(variable)
}

// HasBinding vérifie si une variable est liée dans ce token.
//
// Wrapper autour de Bindings.Has() pour compatibilité.
//
// Paramètres:
//   - variable: nom de la variable
//
// Retourne:
//   - bool: true si la variable existe, false sinon
func (t *Token) HasBinding(variable string) bool {
	if t.Bindings == nil {
		return false
	}
	return t.Bindings.Has(variable)
}

// GetVariables retourne toutes les variables liées dans ce token.
//
// Wrapper autour de Bindings.Variables() pour compatibilité.
//
// Retourne:
//   - []string: liste des noms de variables
func (t *Token) GetVariables() []string {
	if t.Bindings == nil {
		return []string{}
	}
	return t.Bindings.Variables()
}

// generateTokenID génère un ID unique pour un token.
//
// Format: "token_<counter>"
// Cette fonction utilise un compteur atomique pour garantir l'unicité et la thread-safety.
var tokenCounter uint64

func generateTokenID() string {
	// Utiliser atomic.AddUint64 pour garantir thread-safety
	count := atomic.AddUint64(&tokenCounter, 1)
	return fmt.Sprintf("token_%d", count)
}

// NewTokenWithFact crée un nouveau token avec un seul binding.
//
// Fonction utilitaire pour créer un token initial avec un fait unique,
// typiquement utilisé lors de la première activation d'un JoinNode.
//
// Paramètres:
//   - fact: pointeur vers le fait à lier
//   - variable: nom de la variable à lier au fait
//   - nodeID: ID du nœud créateur du token
//
// Retourne:
//   - *Token: nouveau token avec le binding spécifié
//
// Exemple:
//
//	userFact := &Fact{ID: "u1", Type: "User", Fields: map[string]interface{}{"id": 1}}
//	token := NewTokenWithFact(userFact, "user", "type_node_user")
//	fmt.Println(token.HasBinding("user"))  // true
//	fmt.Println(token.GetBinding("user") == userFact)  // true
func NewTokenWithFact(fact *Fact, variable string, nodeID string) *Token {
	return &Token{
		ID:       generateTokenID(),
		Facts:    []*Fact{fact},
		NodeID:   nodeID,
		Bindings: NewBindingChainWith(variable, fact),
		Metadata: TokenMetadata{
			CreatedBy: nodeID,
			JoinLevel: 0,
		},
	}
}
