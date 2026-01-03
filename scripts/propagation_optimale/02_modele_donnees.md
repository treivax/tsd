# 📦 Prompt 02 - Modèle de Données Delta

> **📋 Standards** : Ce prompt respecte les règles de `.github/prompts/common.md` et `.github/prompts/develop.md`

## 🎯 Objectif

Implémenter les structures de données fondamentales pour le système de propagation delta : `FieldDelta`, `FactDelta`, `DependencyIndex`, et leurs utilitaires associés.

**⚠️ IMPORTANT** : Ce prompt génère du code. Respecter strictement les standards de `common.md`.

---

## 📋 Prérequis

Avant de commencer ce prompt :

- [x] **Prompt 01 validé** : Conception complète disponible dans `REPORTS/`
- [x] **Documents de référence** :
  - `REPORTS/conception_delta_architecture.md`
  - `REPORTS/analyse_rete_actuel.md`
- [x] **Branche feature créée** : `feature/propagation-delta-rete-ii`
- [x] **Tests passent** : `make test` (100% success)

---

## 📂 Structure des Fichiers

Créer le package `rete/delta` avec la structure suivante :

```
rete/delta/
├── field_delta.go        # Structures FieldDelta, FactDelta
├── field_delta_test.go   # Tests unitaires delta
├── types.go              # Types et constantes
├── comparison.go         # Comparaison de valeurs
├── comparison_test.go    # Tests comparaison
└── doc.go                # Documentation package
```

---

## 🔧 Tâche 1 : Types et Constantes de Base

### Fichier : `rete/delta/doc.go`

**Contenu** :

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

// Package delta implémente le système de propagation incrémentale (RETE-II/TREAT)
// pour optimiser les mises à jour de faits.
//
// Ce package fournit :
//   - Détection des changements de champs (FieldDelta, FactDelta)
//   - Indexation des dépendances (nœuds sensibles à chaque champ)
//   - Propagation sélective (uniquement vers nœuds affectés)
//
// Architecture :
//
//	Update(fact, {field: value})
//	    ↓
//	DetectDelta(oldFact, newFact) → FactDelta
//	    ↓
//	GetAffectedNodes(delta) → [nodes]
//	    ↓
//	PropagateSelective(delta, nodes)
//
// Performance :
//   - Propagation O(nœuds sensibles) au lieu de O(tous nœuds)
//   - Gain typique : 10-100x sur mises à jour partielles
//
// Compatibilité :
//   - Backward compatible (fallback Retract+Insert disponible)
//   - Feature flag pour activation progressive
package delta
```

### Fichier : `rete/delta/types.go`

**Contenu** :

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import "time"

// ChangeType représente le type de changement sur un champ
type ChangeType int

const (
	// ChangeTypeModified indique que le champ a été modifié
	ChangeTypeModified ChangeType = iota
	// ChangeTypeAdded indique que le champ a été ajouté (absent → présent)
	ChangeTypeAdded
	// ChangeTypeRemoved indique que le champ a été supprimé (présent → absent)
	ChangeTypeRemoved
)

// String retourne la représentation string du ChangeType
func (ct ChangeType) String() string {
	switch ct {
	case ChangeTypeModified:
		return "Modified"
	case ChangeTypeAdded:
		return "Added"
	case ChangeTypeRemoved:
		return "Removed"
	default:
		return "Unknown"
	}
}

// ValueType représente le type d'une valeur
type ValueType int

const (
	ValueTypeUnknown ValueType = iota
	ValueTypeString
	ValueTypeNumber
	ValueTypeBool
	ValueTypeObject
	ValueTypeArray
	ValueTypeNull
)

// String retourne la représentation string du ValueType
func (vt ValueType) String() string {
	switch vt {
	case ValueTypeString:
		return "string"
	case ValueTypeNumber:
		return "number"
	case ValueTypeBool:
		return "bool"
	case ValueTypeObject:
		return "object"
	case ValueTypeArray:
		return "array"
	case ValueTypeNull:
		return "null"
	default:
		return "unknown"
	}
}

// inferValueType déduit le ValueType depuis une valeur Go
func inferValueType(value interface{}) ValueType {
	if value == nil {
		return ValueTypeNull
	}

	switch value.(type) {
	case string:
		return ValueTypeString
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		return ValueTypeNumber
	case bool:
		return ValueTypeBool
	case map[string]interface{}:
		return ValueTypeObject
	case []interface{}:
		return ValueTypeArray
	default:
		return ValueTypeUnknown
	}
}

// Configuration par défaut
const (
	// DefaultFloatEpsilon est la tolérance par défaut pour comparaison de floats
	DefaultFloatEpsilon = 1e-9

	// DefaultMaxDeltaAge est l'âge maximum d'un delta avant expiration (pour cache)
	DefaultMaxDeltaAge = 5 * time.Minute
)
```

**Tests** : `rete/delta/types_test.go`

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import "testing"

func TestChangeType_String(t *testing.T) {
	tests := []struct {
		ct   ChangeType
		want string
	}{
		{ChangeTypeModified, "Modified"},
		{ChangeTypeAdded, "Added"},
		{ChangeTypeRemoved, "Removed"},
		{ChangeType(999), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.ct.String(); got != tt.want {
				t.Errorf("ChangeType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValueType_String(t *testing.T) {
	tests := []struct {
		vt   ValueType
		want string
	}{
		{ValueTypeString, "string"},
		{ValueTypeNumber, "number"},
		{ValueTypeBool, "bool"},
		{ValueTypeObject, "object"},
		{ValueTypeArray, "array"},
		{ValueTypeNull, "null"},
		{ValueTypeUnknown, "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.vt.String(); got != tt.want {
				t.Errorf("ValueType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInferValueType(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		want  ValueType
	}{
		{"nil", nil, ValueTypeNull},
		{"string", "hello", ValueTypeString},
		{"int", 42, ValueTypeNumber},
		{"float64", 3.14, ValueTypeNumber},
		{"bool", true, ValueTypeBool},
		{"object", map[string]interface{}{"a": 1}, ValueTypeObject},
		{"array", []interface{}{1, 2, 3}, ValueTypeArray},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := inferValueType(tt.value); got != tt.want {
				t.Errorf("inferValueType(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}
```

---

## 🔧 Tâche 2 : Structure FieldDelta

### Fichier : `rete/delta/field_delta.go`

**Contenu** :

```go
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
```

**Tests** : `rete/delta/field_delta_test.go`

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"testing"
	"time"
)

func TestNewFieldDelta(t *testing.T) {
	tests := []struct {
		name           string
		fieldName      string
		oldValue       interface{}
		newValue       interface{}
		wantChangeType ChangeType
		wantValueType  ValueType
	}{
		{
			name:           "modification string",
			fieldName:      "status",
			oldValue:       "active",
			newValue:       "inactive",
			wantChangeType: ChangeTypeModified,
			wantValueType:  ValueTypeString,
		},
		{
			name:           "ajout champ",
			fieldName:      "email",
			oldValue:       nil,
			newValue:       "user@example.com",
			wantChangeType: ChangeTypeAdded,
			wantValueType:  ValueTypeString,
		},
		{
			name:           "suppression champ",
			fieldName:      "temp_field",
			oldValue:       42,
			newValue:       nil,
			wantChangeType: ChangeTypeRemoved,
			wantValueType:  ValueTypeNumber,
		},
		{
			name:           "modification nombre",
			fieldName:      "price",
			oldValue:       100.0,
			newValue:       150.0,
			wantChangeType: ChangeTypeModified,
			wantValueType:  ValueTypeNumber,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delta := NewFieldDelta(tt.fieldName, tt.oldValue, tt.newValue)

			if delta.FieldName != tt.fieldName {
				t.Errorf("FieldName = %v, want %v", delta.FieldName, tt.fieldName)
			}
			if delta.ChangeType != tt.wantChangeType {
				t.Errorf("ChangeType = %v, want %v", delta.ChangeType, tt.wantChangeType)
			}
			if delta.ValueType != tt.wantValueType {
				t.Errorf("ValueType = %v, want %v", delta.ValueType, tt.wantValueType)
			}
		})
	}
}

func TestFieldDelta_String(t *testing.T) {
	tests := []struct {
		name      string
		delta     FieldDelta
		wantMatch string // Substring à vérifier
	}{
		{
			name:      "modification",
			delta:     NewFieldDelta("price", 100, 150),
			wantMatch: "price: (100 → 150)",
		},
		{
			name:      "ajout",
			delta:     NewFieldDelta("email", nil, "test@example.com"),
			wantMatch: "email: (nil → test@example.com)",
		},
		{
			name:      "suppression",
			delta:     NewFieldDelta("temp", 42, nil),
			wantMatch: "temp: (42 → nil)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.delta.String()
			if got != tt.wantMatch {
				t.Errorf("String() = %v, want %v", got, tt.wantMatch)
			}
		})
	}
}

func TestNewFactDelta(t *testing.T) {
	before := time.Now()
	delta := NewFactDelta("Product~123", "Product")
	after := time.Now()

	if delta.FactID != "Product~123" {
		t.Errorf("FactID = %v, want Product~123", delta.FactID)
	}
	if delta.FactType != "Product" {
		t.Errorf("FactType = %v, want Product", delta.FactType)
	}
	if delta.Fields == nil {
		t.Error("Fields map should be initialized")
	}
	if delta.Timestamp.Before(before) || delta.Timestamp.After(after) {
		t.Errorf("Timestamp %v not in expected range", delta.Timestamp)
	}
}

func TestFactDelta_AddFieldChange(t *testing.T) {
	delta := NewFactDelta("Product~123", "Product")

	delta.AddFieldChange("price", 100.0, 150.0)
	delta.AddFieldChange("status", "active", "inactive")

	if len(delta.Fields) != 2 {
		t.Errorf("Expected 2 fields, got %d", len(delta.Fields))
	}

	priceChange, ok := delta.Fields["price"]
	if !ok {
		t.Fatal("price field not found")
	}
	if priceChange.OldValue != 100.0 || priceChange.NewValue != 150.0 {
		t.Errorf("price change incorrect: %v → %v", priceChange.OldValue, priceChange.NewValue)
	}
}

func TestFactDelta_IsEmpty(t *testing.T) {
	delta := NewFactDelta("Product~123", "Product")

	if !delta.IsEmpty() {
		t.Error("New delta should be empty")
	}

	delta.AddFieldChange("price", 100.0, 150.0)

	if delta.IsEmpty() {
		t.Error("Delta with changes should not be empty")
	}
}

func TestFactDelta_FieldsChanged(t *testing.T) {
	delta := NewFactDelta("Product~123", "Product")
	delta.AddFieldChange("price", 100.0, 150.0)
	delta.AddFieldChange("status", "active", "inactive")

	fields := delta.FieldsChanged()

	if len(fields) != 2 {
		t.Errorf("Expected 2 changed fields, got %d", len(fields))
	}

	// Vérifier que les champs sont présents (ordre non garanti)
	found := make(map[string]bool)
	for _, field := range fields {
		found[field] = true
	}
	if !found["price"] || !found["status"] {
		t.Errorf("Expected fields [price, status], got %v", fields)
	}
}

func TestFactDelta_ChangeRatio(t *testing.T) {
	tests := []struct {
		name            string
		fieldsChanged   int
		totalFields     int
		wantRatio       float64
		wantRatioApprox bool
	}{
		{"2 sur 10", 2, 10, 0.2, false},
		{"1 sur 5", 1, 5, 0.2, false},
		{"5 sur 10", 5, 10, 0.5, false},
		{"0 sur 10", 0, 10, 0.0, false},
		{"1 sur 0", 1, 0, 0.0, false}, // Edge case
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delta := NewFactDelta("Test~1", "Test")
			delta.FieldCount = tt.totalFields

			// Ajouter tt.fieldsChanged champs
			for i := 0; i < tt.fieldsChanged; i++ {
				delta.AddFieldChange(string(rune('a'+i)), i, i+1)
			}

			ratio := delta.ChangeRatio()
			if ratio != tt.wantRatio {
				t.Errorf("ChangeRatio() = %v, want %v", ratio, tt.wantRatio)
			}
		})
	}
}

func TestFactDelta_String(t *testing.T) {
	t.Run("empty delta", func(t *testing.T) {
		delta := NewFactDelta("Product~123", "Product")
		str := delta.String()
		if str != "FactDelta[Product:Product~123] (no changes)" {
			t.Errorf("String() = %v", str)
		}
	})

	t.Run("delta with changes", func(t *testing.T) {
		delta := NewFactDelta("Product~123", "Product")
		delta.AddFieldChange("price", 100, 150)

		str := delta.String()
		// Vérifier que contient les éléments clés
		if str == "" {
			t.Error("String() should not be empty")
		}
		// Vérifier format (exactitude relâchée car ordre non garanti)
		if len(str) < 10 {
			t.Errorf("String() too short: %v", str)
		}
	})
}
```

---

## 🔧 Tâche 3 : Comparaison de Valeurs

### Fichier : `rete/delta/comparison.go`

**Contenu** :

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
	"math"
	"reflect"
)

// ValuesEqual compare deux valeurs en profondeur.
//
// Cette fonction gère les cas spéciaux :
//   - Floats : comparaison avec epsilon (tolérance)
//   - Maps/slices : comparaison récursive
//   - nil : nil == nil → true
//
// Paramètres :
//   - a, b : valeurs à comparer
//   - epsilon : tolérance pour floats (utiliser DefaultFloatEpsilon)
//
// Retourne true si les valeurs sont égales.
func ValuesEqual(a, b interface{}, epsilon float64) bool {
	// Cas 1: Égalité stricte (pointeurs identiques ou valeurs simples)
	if a == b {
		return true
	}

	// Cas 2: nil vs non-nil
	if a == nil || b == nil {
		return a == b
	}

	// Cas 3: Types différents → inégaux
	typeA := reflect.TypeOf(a)
	typeB := reflect.TypeOf(b)
	if typeA != typeB {
		return false
	}

	// Cas 4: Floats avec tolérance epsilon
	switch va := a.(type) {
	case float64:
		vb, ok := b.(float64)
		if !ok {
			return false
		}
		return math.Abs(va-vb) <= epsilon

	case float32:
		vb, ok := b.(float32)
		if !ok {
			return false
		}
		return math.Abs(float64(va-vb)) <= epsilon
	}

	// Cas 5: Comparaison profonde (maps, slices, structs)
	return reflect.DeepEqual(a, b)
}

// FactsEqual compare deux faits (représentés comme maps) en profondeur.
//
// Paramètres :
//   - fact1, fact2 : faits à comparer (map[string]interface{})
//
// Retourne true si les faits sont identiques (mêmes champs, mêmes valeurs).
func FactsEqual(fact1, fact2 map[string]interface{}) bool {
	if len(fact1) != len(fact2) {
		return false
	}

	for key, val1 := range fact1 {
		val2, exists := fact2[key]
		if !exists {
			return false
		}
		if !ValuesEqual(val1, val2, DefaultFloatEpsilon) {
			return false
		}
	}

	return true
}
```

**Tests** : `rete/delta/comparison_test.go`

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import "testing"

func TestValuesEqual(t *testing.T) {
	tests := []struct {
		name    string
		a       interface{}
		b       interface{}
		epsilon float64
		want    bool
	}{
		// Cas simples
		{"int égaux", 42, 42, 0, true},
		{"int différents", 42, 43, 0, false},
		{"string égaux", "hello", "hello", 0, true},
		{"string différents", "hello", "world", 0, false},
		{"bool égaux", true, true, 0, true},
		{"bool différents", true, false, 0, false},

		// nil
		{"nil == nil", nil, nil, 0, true},
		{"nil != value", nil, 42, 0, false},
		{"value != nil", 42, nil, 0, false},

		// Floats avec epsilon
		{"float64 égaux", 1.0, 1.0, DefaultFloatEpsilon, true},
		{"float64 proches", 1.0, 1.0000000001, DefaultFloatEpsilon, true},
		{"float64 différents", 1.0, 1.1, DefaultFloatEpsilon, false},
		{"float32 égaux", float32(1.0), float32(1.0), DefaultFloatEpsilon, true},

		// Types différents
		{"int vs float", 1, 1.0, 0, false},
		{"string vs int", "1", 1, 0, false},

		// Maps
		{
			"maps égaux",
			map[string]interface{}{"a": 1, "b": 2},
			map[string]interface{}{"a": 1, "b": 2},
			0,
			true,
		},
		{
			"maps différents",
			map[string]interface{}{"a": 1, "b": 2},
			map[string]interface{}{"a": 1, "b": 3},
			0,
			false,
		},

		// Slices
		{
			"slices égaux",
			[]interface{}{1, 2, 3},
			[]interface{}{1, 2, 3},
			0,
			true,
		},
		{
			"slices différents",
			[]interface{}{1, 2, 3},
			[]interface{}{1, 2, 4},
			0,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValuesEqual(tt.a, tt.b, tt.epsilon)
			if got != tt.want {
				t.Errorf("ValuesEqual(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestFactsEqual(t *testing.T) {
	tests := []struct {
		name  string
		fact1 map[string]interface{}
		fact2 map[string]interface{}
		want  bool
	}{
		{
			name:  "faits identiques",
			fact1: map[string]interface{}{"id": "123", "price": 100.0, "status": "active"},
			fact2: map[string]interface{}{"id": "123", "price": 100.0, "status": "active"},
			want:  true,
		},
		{
			name:  "faits différents (1 champ)",
			fact1: map[string]interface{}{"id": "123", "price": 100.0, "status": "active"},
			fact2: map[string]interface{}{"id": "123", "price": 150.0, "status": "active"},
			want:  false,
		},
		{
			name:  "faits de tailles différentes",
			fact1: map[string]interface{}{"id": "123", "price": 100.0},
			fact2: map[string]interface{}{"id": "123", "price": 100.0, "status": "active"},
			want:  false,
		},
		{
			name:  "faits avec champs manquants",
			fact1: map[string]interface{}{"id": "123", "price": 100.0},
			fact2: map[string]interface{}{"id": "123", "quantity": 10},
			want:  false,
		},
		{
			name:  "faits vides",
			fact1: map[string]interface{}{},
			fact2: map[string]interface{}{},
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FactsEqual(tt.fact1, tt.fact2)
			if got != tt.want {
				t.Errorf("FactsEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkValuesEqual_Simple(b *testing.B) {
	v1 := 42
	v2 := 42
	for i := 0; i < b.N; i++ {
		_ = ValuesEqual(v1, v2, 0)
	}
}

func BenchmarkValuesEqual_Float(b *testing.B) {
	v1 := 3.14159
	v2 := 3.14159
	for i := 0; i < b.N; i++ {
		_ = ValuesEqual(v1, v2, DefaultFloatEpsilon)
	}
}

func BenchmarkFactsEqual(b *testing.B) {
	fact1 := map[string]interface{}{
		"id":     "123",
		"price":  100.0,
		"status": "active",
		"qty":    10,
	}
	fact2 := map[string]interface{}{
		"id":     "123",
		"price":  100.0,
		"status": "active",
		"qty":    10,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FactsEqual(fact1, fact2)
	}
}
```

---

## ✅ Validation

Après implémentation, exécuter :

```bash
# 1. Formattage
go fmt ./rete/delta/...
goimports -w ./rete/delta/

# 2. Validation statique
go vet ./rete/delta/...
staticcheck ./rete/delta/...

# 3. Tests
go test ./rete/delta/... -v
go test ./rete/delta/... -cover

# 4. Benchmarks
go test ./rete/delta/... -bench=. -benchmem

# 5. Race detector
go test ./rete/delta/... -race

# 6. Validation complète
make validate
```

**Critères de succès** :
- [ ] Tous les tests passent (100%)
- [ ] Couverture > 90%
- [ ] Aucune erreur `go vet`, `staticcheck`
- [ ] Benchmarks exécutés sans panic
- [ ] Aucune race condition détectée

---

## 📊 Livrables

À la fin de ce prompt :

1. **Code** :
   - ✅ `rete/delta/doc.go` - Documentation package
   - ✅ `rete/delta/types.go` - Types et constantes
   - ✅ `rete/delta/field_delta.go` - Structures FieldDelta, FactDelta
   - ✅ `rete/delta/comparison.go` - Comparaison de valeurs

2. **Tests** :
   - ✅ `rete/delta/types_test.go`
   - ✅ `rete/delta/field_delta_test.go`
   - ✅ `rete/delta/comparison_test.go`

3. **Validation** :
   - ✅ Rapport de couverture
   - ✅ Résultats benchmarks

---

## 🚀 Commit

Une fois validé :

```bash
git add rete/delta/
git commit -m "feat(rete): [Prompt 02] Implémentation modèle de données delta

- Structures FieldDelta, FactDelta pour représentation changements
- Types ChangeType, ValueType avec inférence automatique
- Comparaison valeurs avec support epsilon (floats)
- Comparaison faits en profondeur
- Tests unitaires complets (> 90% couverture)
- Benchmarks pour performance"
```

---

## 🚦 Prochaine Étape

Passer au **Prompt 03 - Indexation des Dépendances**

---

**Durée estimée** : 2-3 heures  
**Difficulté** : Moyenne  
**Prérequis** : Prompt 01 validé  
**Couverture cible** : > 90%