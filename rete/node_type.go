// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
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
	fmt.Printf("🗑️  [TYPE_%s] Rétractation du fait: %s\n", tn.ID, factID)
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
	return tn.PropagateToChildren(fact, nil)
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
