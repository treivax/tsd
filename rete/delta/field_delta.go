// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"fmt"
	"time"
)

// FieldDelta représente le changement d'un champ spécifique.
//
// Il capture l'ancienne et la nouvelle valeur, le type de changement,
// et le type de la valeur pour validation et comparaison.
type FieldDelta struct {
	// FieldName est le nom du champ modifié
	FieldName string

	// OldValue est l'ancienne valeur (nil si champ ajouté)
	OldValue interface{}

	// NewValue est la nouvelle valeur (nil si champ supprimé)
	NewValue interface{}

	// ChangeType indique le type de changement (Modified, Added, Removed)
	ChangeType ChangeType

	// ValueType est le type de la nouvelle valeur
	ValueType ValueType
}

// NewFieldDelta crée un nouveau FieldDelta.
//
// Paramètres :
//   - fieldName : nom du champ
//   - oldValue : ancienne valeur
//   - newValue : nouvelle valeur
//
// Retourne un FieldDelta avec ChangeType et ValueType automatiquement déduits.
func NewFieldDelta(fieldName string, oldValue, newValue interface{}) FieldDelta {
	changeType := ChangeTypeModified
	if oldValue == nil && newValue != nil {
		changeType = ChangeTypeAdded
	} else if oldValue != nil && newValue == nil {
		changeType = ChangeTypeRemoved
	}

	valueType := inferValueType(newValue)
	if newValue == nil {
		valueType = inferValueType(oldValue)
	}

	return FieldDelta{
		FieldName:  fieldName,
		OldValue:   oldValue,
		NewValue:   newValue,
		ChangeType: changeType,
		ValueType:  valueType,
	}
}

// String retourne une représentation string du FieldDelta
func (fd FieldDelta) String() string {
	switch fd.ChangeType {
	case ChangeTypeAdded:
		return fmt.Sprintf("%s: (nil → %v)", fd.FieldName, fd.NewValue)
	case ChangeTypeRemoved:
		return fmt.Sprintf("%s: (%v → nil)", fd.FieldName, fd.OldValue)
	default:
		return fmt.Sprintf("%s: (%v → %v)", fd.FieldName, fd.OldValue, fd.NewValue)
	}
}

// FactDelta représente l'ensemble des changements sur un fait.
//
// Il contient tous les champs modifiés, ajoutés ou supprimés,
// ainsi que des métadonnées pour tracking et cache.
type FactDelta struct {
	// FactID est l'identifiant interne du fait (ex: "Product~123")
	FactID string

	// FactType est le type du fait (ex: "Product")
	FactType string

	// Fields contient les changements par nom de champ
	Fields map[string]FieldDelta

	// Timestamp est le moment de création du delta
	Timestamp time.Time

	// FieldCount est le nombre total de champs dans le fait complet
	// (utilisé pour calculer le ratio de changement)
	FieldCount int
}

// NewFactDelta crée un nouveau FactDelta.
//
// Paramètres :
//   - factID : identifiant interne du fait
//   - factType : type du fait
//
// Retourne un FactDelta initialisé avec timestamp actuel.
func NewFactDelta(factID, factType string) *FactDelta {
	return &FactDelta{
		FactID:    factID,
		FactType:  factType,
		Fields:    make(map[string]FieldDelta),
		Timestamp: time.Now(),
	}
}

// AddFieldChange enregistre un changement de champ.
//
// Paramètres :
//   - fieldName : nom du champ
//   - oldValue : ancienne valeur
//   - newValue : nouvelle valeur
func (fd *FactDelta) AddFieldChange(fieldName string, oldValue, newValue interface{}) {
	fieldDelta := NewFieldDelta(fieldName, oldValue, newValue)
	fd.Fields[fieldName] = fieldDelta
}

// IsEmpty retourne true si aucun champ n'a changé.
func (fd *FactDelta) IsEmpty() bool {
	return len(fd.Fields) == 0
}

// FieldsChanged retourne la liste des noms de champs modifiés.
func (fd *FactDelta) FieldsChanged() []string {
	fields := make([]string, 0, len(fd.Fields))
	for fieldName := range fd.Fields {
		fields = append(fields, fieldName)
	}
	return fields
}

// ChangeRatio retourne le ratio de champs modifiés (entre 0.0 et 1.0).
//
// Exemple : 2 champs modifiés sur 10 → 0.2
//
// Ce ratio est utilisé pour décider si la propagation delta est pertinente
// (seuil typique : < 0.3 → delta, >= 0.3 → Retract+Insert classique).
func (fd *FactDelta) ChangeRatio() float64 {
	if fd.FieldCount == 0 {
		return 0.0
	}
	return float64(len(fd.Fields)) / float64(fd.FieldCount)
}

// String retourne une représentation string du FactDelta
func (fd *FactDelta) String() string {
	if fd.IsEmpty() {
		return fmt.Sprintf("FactDelta[%s:%s] (no changes)", fd.FactType, fd.FactID)
	}

	fieldsStr := ""
	for _, delta := range fd.Fields {
		fieldsStr += delta.String() + ", "
	}
	if len(fieldsStr) > 2 {
		fieldsStr = fieldsStr[:len(fieldsStr)-2] // Enlever dernière ", "
	}

	return fmt.Sprintf("FactDelta[%s:%s] {%s}", fd.FactType, fd.FactID, fieldsStr)
}
