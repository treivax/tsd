// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"github.com/treivax/tsd/tsdio"
	"fmt"
)

type TypeNode struct {
	BaseNode
	TypeName       string         `json:"type_name"`
	TypeDefinition TypeDefinition `json:"type_definition"`
}

// NewTypeNode crée un nouveau nœud de type
func NewTypeNode(typeName string, typeDef TypeDefinition, storage Storage) *TypeNode {
	return &TypeNode{
		BaseNode: BaseNode{
			ID:       fmt.Sprintf("type_%s", typeName),
			Type:     "type",
			Memory:   &WorkingMemory{NodeID: fmt.Sprintf("type_%s", typeName), Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
			Children: make([]Node, 0),
			Storage:  storage,
		},
		TypeName:       typeName,
		TypeDefinition: typeDef,
	}
}

// ActivateLeft (non utilisé pour les nœuds de type)
func (tn *TypeNode) ActivateLeft(token *Token) error {
	return fmt.Errorf("les nœuds de type ne reçoivent pas de tokens")
}

// ActivateRetract retire le fait de la mémoire de type et propage aux enfants
// factID doit être l'identifiant interne (Type_ID)
func (tn *TypeNode) ActivateRetract(factID string) error {
	tn.mutex.Lock()
	_, exists := tn.Memory.GetFact(factID)
	if exists {
		tn.Memory.RemoveFact(factID)
	}
	tn.mutex.Unlock()
	if !exists {
		return nil
	}
	tsdio.Printf("🗑️  [TYPE_%s] Rétractation du fait: %s\n", tn.ID, factID)
	return tn.PropagateRetractToChildren(factID)
}

// ActivateRight filtre les faits par type et les propage
func (tn *TypeNode) ActivateRight(fact *Fact) error {
	// Vérifier si le fait correspond au type de ce nœud
	if fact.Type != tn.TypeName {
		return nil // Ignorer silencieusement les faits d'autres types
	}

	// Log désactivé pour les performances
	// fmt.Printf("[TYPE_%s] Reçu fait: %s\n", tn.TypeName, fact.String())

	// Valider les champs du fait
	if err := tn.validateFact(fact); err != nil {
		return fmt.Errorf("validation du fait échouée: %w", err)
	}

	tn.mutex.Lock()
	if err := tn.Memory.AddFact(fact); err != nil {
		tn.mutex.Unlock()
		return fmt.Errorf("erreur ajout fait dans type node: %w", err)
	}
	tn.mutex.Unlock()

	// Persistance désactivée pour les performances

	// Propager aux enfants (AlphaNodes)
	// Check if any child is a decomposed alpha chain and use context-aware activation
	for _, child := range tn.GetChildren() {
		if alphaNode, ok := child.(*AlphaNode); ok {
			// Check if this AlphaNode is part of a decomposed chain
			if alphaNode.IsAtomic || len(alphaNode.Dependencies) > 0 {
				// Create evaluation context for decomposed chain with cache if available
				var ctx *EvaluationContext
				if tn.network != nil && tn.network.ArithmeticResultCache != nil {
					ctx = NewEvaluationContextWithCache(fact, tn.network.ArithmeticResultCache)
				} else {
					ctx = NewEvaluationContext(fact)
				}
				if err := alphaNode.ActivateWithContext(fact, ctx); err != nil {
					return fmt.Errorf("error activating decomposed alpha chain: %w", err)
				}
			} else {
				// Standard activation for non-decomposed alpha nodes
				if err := alphaNode.ActivateRight(fact); err != nil {
					return fmt.Errorf("error activating alpha node: %w", err)
				}
			}
		} else {
			// Non-alpha child, use standard propagation
			if err := child.ActivateRight(fact); err != nil {
				return fmt.Errorf("error propagating to child: %w", err)
			}
		}
	}

	return nil
}

// validateFact valide qu'un fait respecte la définition de type
func (tn *TypeNode) validateFact(fact *Fact) error {
	for _, field := range tn.TypeDefinition.Fields {
		// Le champ "id" est stocké dans fact.ID, pas dans Fields
		if field.Name == "id" {
			if fact.ID == "" {
				return fmt.Errorf("champ manquant: %s", field.Name)
			}
			continue
		}

		value, exists := fact.Fields[field.Name]
		if !exists {
			return fmt.Errorf("champ manquant: %s", field.Name)
		}

		// Validation basique des types
		if !tn.isValidType(value, field.Type) {
			return fmt.Errorf("type invalide pour le champ %s: attendu %s", field.Name, field.Type)
		}
	}
	return nil
}

// isValidType vérifie si une valeur correspond au type attendu
func (tn *TypeNode) isValidType(value interface{}, expectedType string) bool {
	switch expectedType {
	case "string":
		_, ok := value.(string)
		return ok
	case "number":
		switch value.(type) {
		case int, int32, int64, float32, float64:
			return true
		}
		return false
	case "bool":
		_, ok := value.(bool)
		return ok
	default:
		return false
	}
}

// Clone crée une copie profonde du TypeNode
func (tn *TypeNode) Clone() *TypeNode {
	clone := &TypeNode{
		BaseNode: BaseNode{
			ID:       tn.ID,
			Type:     tn.Type,
			Memory:   tn.Memory.Clone(),
			Children: make([]Node, len(tn.Children)),
			Storage:  tn.Storage,
		},
		TypeName:       tn.TypeName,
		TypeDefinition: tn.TypeDefinition.Clone(),
	}

	// Copier les enfants
	copy(clone.Children, tn.Children)

	return clone
}
