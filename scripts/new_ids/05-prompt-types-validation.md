# Prompt 05 - Types et Validation Complète

> **📋 Standards** : Ce prompt respecte les règles de [.github/prompts/common.md](../../.github/prompts/common.md) et [.github/prompts/develop.md](../../.github/prompts/develop.md)

---

## 🎯 Objectif

Compléter et renforcer le système de types et de validation pour supporter pleinement les champs de type fait :

1. **Validation complète des types** - Vérifier cohérence et compatibilité
2. **Résolution de types** - Inférence et vérification statique
3. **Validation des affectations** - Variables et références
4. **Erreurs explicites** - Messages clairs et informatifs
5. **Validation des règles** - Cohérence des expressions

---

## 📋 Contexte

### État Actuel (après prompts précédents)

- Parser accepte les types de faits dans les champs
- Génération d'IDs avec références de faits fonctionne
- Comparaisons de faits dans RETE fonctionnent
- Validation de base existe mais incomplète

### État Cible

```tsd
type User(#name: string, age: number)
type Login(user: User, #email: string, password: string)
type Audit(login: Login, timestamp: number, action: string)

// ✅ Validation complète
alice = User("Alice", 30)          // Type correct
Login(alice, "a@ex.com", "pw")     // alice est bien un User
Audit(login_ref, 1234567890, "login") // login_ref validé

// ❌ Erreurs détectées
bob = User("Bob", "invalid_age")   // Erreur: age doit être number
Login("not_a_user", "e@ex.com", "pw") // Erreur: user doit être User
Login(unknown, "e@ex.com", "pw")   // Erreur: variable unknown non définie
```

---

## 📝 Tâches à Réaliser

### 1. Analyser la Validation Actuelle

#### Fichiers Existants

**Rechercher** :
```bash
find constraint/ -name "*validation*.go" -type f
grep -r "Validate" constraint/ --include="*.go" -l
```

**Fichiers attendus** :
- `constraint/constraint_validation.go`
- `constraint/constraint_type_validation.go`
- `constraint/constraint_field_validation.go`
- `constraint/primary_key_validation.go`
- `constraint/action_validator.go`

#### Identifier les Lacunes

**Questions à répondre** :
1. Les types de faits dans les champs sont-ils validés ?
2. Les références de variables sont-elles vérifiées ?
3. Les affectations sont-elles validées avant utilisation ?
4. Les types dans les comparaisons sont-ils vérifiés ?
5. Les messages d'erreur sont-ils explicites ?

### 2. Créer un Validateur de Types Complet

#### Nouveau Fichier : `constraint/type_system.go`

**Objectif** : Système de types centralisé pour toutes les validations

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"fmt"
	"strings"
)

// TypeSystem gère le système de types du langage TSD
// Il maintient les définitions de types et fournit des utilitaires de validation
type TypeSystem struct {
	// Types contient toutes les définitions de types
	Types map[string]TypeDefinition
	
	// Variables contient les types des variables définies par affectation
	Variables map[string]string
}

// NewTypeSystem crée un nouveau système de types
func NewTypeSystem(types []TypeDefinition) *TypeSystem {
	typeMap := make(map[string]TypeDefinition)
	for _, t := range types {
		typeMap[t.Name] = t
	}
	
	return &TypeSystem{
		Types:     typeMap,
		Variables: make(map[string]string),
	}
}

// IsPrimitiveType vérifie si un type est primitif
func (ts *TypeSystem) IsPrimitiveType(typeName string) bool {
	primitives := map[string]bool{
		ValueTypeString:  true,
		ValueTypeNumber:  true,
		ValueTypeBoolean: true,
		ValueTypeBool:    true,
	}
	return primitives[typeName]
}

// IsUserDefinedType vérifie si un type est défini par l'utilisateur
func (ts *TypeSystem) IsUserDefinedType(typeName string) bool {
	_, exists := ts.Types[typeName]
	return exists
}

// TypeExists vérifie qu'un type existe (primitif ou user-defined)
func (ts *TypeSystem) TypeExists(typeName string) bool {
	return ts.IsPrimitiveType(typeName) || ts.IsUserDefinedType(typeName)
}

// GetFieldType retourne le type d'un champ dans un type donné
func (ts *TypeSystem) GetFieldType(typeName, fieldName string) (string, error) {
	// Le champ _id_ est interdit
	if fieldName == FieldNameInternalID {
		return "", fmt.Errorf(
			"le champ '%s' est interne et ne peut pas être accédé",
			FieldNameInternalID,
		)
	}
	
	typeDef, exists := ts.Types[typeName]
	if !exists {
		return "", fmt.Errorf("type '%s' non trouvé", typeName)
	}
	
	for _, field := range typeDef.Fields {
		if field.Name == fieldName {
			return field.Type, nil
		}
	}
	
	return "", fmt.Errorf(
		"champ '%s' non trouvé dans le type '%s'",
		fieldName,
		typeName,
	)
}

// ValidateFieldType valide qu'un champ existe et retourne son type
func (ts *TypeSystem) ValidateFieldType(typeName, fieldName string) (string, error) {
	fieldType, err := ts.GetFieldType(typeName, fieldName)
	if err != nil {
		return "", err
	}
	
	if !ts.TypeExists(fieldType) {
		return "", fmt.Errorf(
			"type '%s' du champ '%s.%s' n'existe pas",
			fieldType,
			typeName,
			fieldName,
		)
	}
	
	return fieldType, nil
}

// RegisterVariable enregistre une variable avec son type
func (ts *TypeSystem) RegisterVariable(varName, typeName string) error {
	if !ts.TypeExists(typeName) {
		return fmt.Errorf(
			"impossible d'enregistrer la variable '%s': type '%s' n'existe pas",
			varName,
			typeName,
		)
	}
	
	ts.Variables[varName] = typeName
	return nil
}

// GetVariableType retourne le type d'une variable
func (ts *TypeSystem) GetVariableType(varName string) (string, error) {
	typeName, exists := ts.Variables[varName]
	if !exists {
		return "", fmt.Errorf("variable '%s' non définie", varName)
	}
	return typeName, nil
}

// VariableExists vérifie qu'une variable existe
func (ts *TypeSystem) VariableExists(varName string) bool {
	_, exists := ts.Variables[varName]
	return exists
}

// AreTypesCompatible vérifie si deux types sont compatibles pour une opération
func (ts *TypeSystem) AreTypesCompatible(type1, type2 string, operator string) bool {
	// Même type exact
	if type1 == type2 {
		return true
	}
	
	// bool et boolean sont compatibles
	if (type1 == ValueTypeBool || type1 == ValueTypeBoolean) &&
		(type2 == ValueTypeBool || type2 == ValueTypeBoolean) {
		return true
	}
	
	// Pour les types de faits, seuls == et != sont autorisés
	if ts.IsUserDefinedType(type1) && ts.IsUserDefinedType(type2) {
		return (operator == "==" || operator == "!=") && type1 == type2
	}
	
	return false
}

// ValidateCircularReferences détecte les références circulaires dans les types
func (ts *TypeSystem) ValidateCircularReferences() error {
	// Construire le graphe de dépendances
	graph := make(map[string][]string)
	
	for typeName, typeDef := range ts.Types {
		for _, field := range typeDef.Fields {
			if ts.IsUserDefinedType(field.Type) {
				graph[typeName] = append(graph[typeName], field.Type)
			}
		}
	}
	
	// Détection de cycles avec DFS
	visited := make(map[string]bool)
	recStack := make(map[string]bool)
	
	var hasCycle func(string) bool
	hasCycle = func(node string) bool {
		visited[node] = true
		recStack[node] = true
		
		for _, neighbor := range graph[node] {
			if !visited[neighbor] {
				if hasCycle(neighbor) {
					return true
				}
			} else if recStack[neighbor] {
				return true
			}
		}
		
		recStack[node] = false
		return false
	}
	
	for typeName := range ts.Types {
		if !visited[typeName] {
			if hasCycle(typeName) {
				return fmt.Errorf(
					"référence circulaire détectée impliquant le type '%s'",
					typeName,
				)
			}
		}
	}
	
	return nil
}

// GetTypePath retourne le chemin de types pour une expression de field access
// Ex: login.user.name -> [Login, User, string]
func (ts *TypeSystem) GetTypePath(rootType, fieldPath string) ([]string, error) {
	parts := strings.Split(fieldPath, ".")
	path := []string{rootType}
	currentType := rootType
	
	for _, fieldName := range parts {
		fieldType, err := ts.GetFieldType(currentType, fieldName)
		if err != nil {
			return nil, err
		}
		
		path = append(path, fieldType)
		currentType = fieldType
	}
	
	return path, nil
}
```

#### Tests du Système de Types

**Fichier : `constraint/type_system_test.go`**

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"testing"
)

func TestTypeSystem_TypeChecks(t *testing.T) {
	t.Log("🧪 TEST TYPE SYSTEM - VÉRIFICATIONS DE TYPES")
	t.Log("=============================================")
	
	types := []TypeDefinition{
		{
			Name: "User",
			Fields: []Field{
				{Name: "name", Type: "string", IsPrimaryKey: true},
				{Name: "age", Type: "number"},
			},
		},
		{
			Name: "Login",
			Fields: []Field{
				{Name: "user", Type: "User"},
				{Name: "email", Type: "string", IsPrimaryKey: true},
			},
		},
	}
	
	ts := NewTypeSystem(types)
	
	tests := []struct {
		name     string
		testFunc func() bool
		expected bool
	}{
		{
			name:     "string est primitif",
			testFunc: func() bool { return ts.IsPrimitiveType("string") },
			expected: true,
		},
		{
			name:     "User n'est pas primitif",
			testFunc: func() bool { return ts.IsPrimitiveType("User") },
			expected: false,
		},
		{
			name:     "User est user-defined",
			testFunc: func() bool { return ts.IsUserDefinedType("User") },
			expected: true,
		},
		{
			name:     "string n'est pas user-defined",
			testFunc: func() bool { return ts.IsUserDefinedType("string") },
			expected: false,
		},
		{
			name:     "User existe",
			testFunc: func() bool { return ts.TypeExists("User") },
			expected: true,
		},
		{
			name:     "string existe",
			testFunc: func() bool { return ts.TypeExists("string") },
			expected: true,
		},
		{
			name:     "Unknown n'existe pas",
			testFunc: func() bool { return ts.TypeExists("Unknown") },
			expected: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.testFunc()
			if result != tt.expected {
				t.Errorf("❌ Attendu %v, reçu %v", tt.expected, result)
			} else {
				t.Logf("✅ Vérification correcte")
			}
		})
	}
}

func TestTypeSystem_GetFieldType(t *testing.T) {
	t.Log("🧪 TEST TYPE SYSTEM - TYPE DE CHAMP")
	t.Log("====================================")
	
	types := []TypeDefinition{
		{
			Name: "User",
			Fields: []Field{
				{Name: "name", Type: "string", IsPrimaryKey: true},
				{Name: "age", Type: "number"},
			},
		},
		{
			Name: "Login",
			Fields: []Field{
				{Name: "user", Type: "User"},
				{Name: "email", Type: "string", IsPrimaryKey: true},
			},
		},
	}
	
	ts := NewTypeSystem(types)
	
	tests := []struct {
		name      string
		typeName  string
		fieldName string
		wantType  string
		wantErr   bool
	}{
		{
			name:      "champ primitif",
			typeName:  "User",
			fieldName: "name",
			wantType:  "string",
			wantErr:   false,
		},
		{
			name:      "champ de type fait",
			typeName:  "Login",
			fieldName: "user",
			wantType:  "User",
			wantErr:   false,
		},
		{
			name:      "champ _id_ interdit",
			typeName:  "User",
			fieldName: FieldNameInternalID,
			wantErr:   true,
		},
		{
			name:      "champ inexistant",
			typeName:  "User",
			fieldName: "unknown",
			wantErr:   true,
		},
		{
			name:      "type inexistant",
			typeName:  "Unknown",
			fieldName: "field",
			wantErr:   true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fieldType, err := ts.GetFieldType(tt.typeName, tt.fieldName)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("❌ Attendu une erreur, reçu nil")
				} else {
					t.Logf("✅ Erreur attendue: %v", err)
				}
				return
			}
			
			if err != nil {
				t.Fatalf("❌ Erreur inattendue: %v", err)
			}
			
			if fieldType != tt.wantType {
				t.Errorf("❌ Type attendu '%s', reçu '%s'", tt.wantType, fieldType)
			} else {
				t.Logf("✅ Type correct: %s", fieldType)
			}
		})
	}
}

func TestTypeSystem_Variables(t *testing.T) {
	t.Log("🧪 TEST TYPE SYSTEM - VARIABLES")
	t.Log("================================")
	
	types := []TypeDefinition{
		{
			Name: "User",
			Fields: []Field{
				{Name: "name", Type: "string", IsPrimaryKey: true},
			},
		},
	}
	
	ts := NewTypeSystem(types)
	
	// Enregistrer une variable
	err := ts.RegisterVariable("alice", "User")
	if err != nil {
		t.Fatalf("❌ Erreur d'enregistrement: %v", err)
	}
	t.Log("✅ Variable 'alice' enregistrée")
	
	// Vérifier qu'elle existe
	if !ts.VariableExists("alice") {
		t.Error("❌ Variable 'alice' devrait exister")
	}
	
	// Récupérer son type
	varType, err := ts.GetVariableType("alice")
	if err != nil {
		t.Fatalf("❌ Erreur de récupération: %v", err)
	}
	
	if varType != "User" {
		t.Errorf("❌ Type attendu 'User', reçu '%s'", varType)
	} else {
		t.Logf("✅ Type de variable correct: %s", varType)
	}
	
	// Variable non définie
	_, err = ts.GetVariableType("bob")
	if err == nil {
		t.Error("❌ Attendu une erreur pour variable non définie")
	} else {
		t.Logf("✅ Erreur pour variable non définie: %v", err)
	}
	
	// Type inexistant
	err = ts.RegisterVariable("invalid", "UnknownType")
	if err == nil {
		t.Error("❌ Attendu une erreur pour type inexistant")
	} else {
		t.Logf("✅ Erreur pour type inexistant: %v", err)
	}
}

func TestTypeSystem_CircularReferences(t *testing.T) {
	t.Log("🧪 TEST TYPE SYSTEM - RÉFÉRENCES CIRCULAIRES")
	t.Log("=============================================")
	
	// Cas 1: Pas de cycle
	types1 := []TypeDefinition{
		{
			Name: "User",
			Fields: []Field{
				{Name: "name", Type: "string"},
			},
		},
		{
			Name: "Login",
			Fields: []Field{
				{Name: "user", Type: "User"},
			},
		},
	}
	
	ts1 := NewTypeSystem(types1)
	err := ts1.ValidateCircularReferences()
	if err != nil {
		t.Errorf("❌ Pas de cycle attendu, erreur reçue: %v", err)
	} else {
		t.Log("✅ Aucun cycle détecté (correct)")
	}
	
	// Cas 2: Cycle direct A -> B -> A
	types2 := []TypeDefinition{
		{
			Name: "A",
			Fields: []Field{
				{Name: "b", Type: "B"},
			},
		},
		{
			Name: "B",
			Fields: []Field{
				{Name: "a", Type: "A"},
			},
		},
	}
	
	ts2 := NewTypeSystem(types2)
	err = ts2.ValidateCircularReferences()
	if err == nil {
		t.Error("❌ Cycle attendu, aucune erreur reçue")
	} else {
		t.Logf("✅ Cycle détecté: %v", err)
	}
	
	// Cas 3: Cycle indirect A -> B -> C -> A
	types3 := []TypeDefinition{
		{
			Name: "A",
			Fields: []Field{
				{Name: "b", Type: "B"},
			},
		},
		{
			Name: "B",
			Fields: []Field{
				{Name: "c", Type: "C"},
			},
		},
		{
			Name: "C",
			Fields: []Field{
				{Name: "a", Type: "A"},
			},
		},
	}
	
	ts3 := NewTypeSystem(types3)
	err = ts3.ValidateCircularReferences()
	if err == nil {
		t.Error("❌ Cycle indirect attendu, aucune erreur reçue")
	} else {
		t.Logf("✅ Cycle indirect détecté: %v", err)
	}
}

func TestTypeSystem_TypeCompatibility(t *testing.T) {
	t.Log("🧪 TEST TYPE SYSTEM - COMPATIBILITÉ")
	t.Log("====================================")
	
	types := []TypeDefinition{
		{
			Name: "User",
			Fields: []Field{
				{Name: "name", Type: "string"},
			},
		},
	}
	
	ts := NewTypeSystem(types)
	
	tests := []struct {
		name     string
		type1    string
		type2    string
		operator string
		expected bool
	}{
		{
			name:     "même type primitif",
			type1:    "string",
			type2:    "string",
			operator: "==",
			expected: true,
		},
		{
			name:     "bool et boolean compatibles",
			type1:    "bool",
			type2:    "boolean",
			operator: "==",
			expected: true,
		},
		{
			name:     "types primitifs différents",
			type1:    "string",
			type2:    "number",
			operator: "==",
			expected: false,
		},
		{
			name:     "même type fait avec ==",
			type1:    "User",
			type2:    "User",
			operator: "==",
			expected: true,
		},
		{
			name:     "même type fait avec !=",
			type1:    "User",
			type2:    "User",
			operator: "!=",
			expected: true,
		},
		{
			name:     "type fait avec < interdit",
			type1:    "User",
			type2:    "User",
			operator: "<",
			expected: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ts.AreTypesCompatible(tt.type1, tt.type2, tt.operator)
			if result != tt.expected {
				t.Errorf("❌ Attendu %v, reçu %v", tt.expected, result)
			} else {
				t.Logf("✅ Compatibilité correcte")
			}
		})
	}
}
```

### 3. Validateur de Faits Complet

#### Fichier : `constraint/fact_validator.go` (nouveau)

**Objectif** : Validation complète des faits avec types de faits

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"fmt"
)

// FactValidator valide les faits selon leur définition de type
type FactValidator struct {
	typeSystem *TypeSystem
}

// NewFactValidator crée un nouveau validateur de faits
func NewFactValidator(ts *TypeSystem) *FactValidator {
	return &FactValidator{
		typeSystem: ts,
	}
}

// ValidateFact valide un fait complet
func (fv *FactValidator) ValidateFact(fact Fact) error {
	// Vérifier que le type existe
	if !fv.typeSystem.TypeExists(fact.TypeName) {
		return fmt.Errorf(
			"type '%s' non défini",
			fact.TypeName,
		)
	}
	
	typeDef, _ := fv.typeSystem.Types[fact.TypeName]
	
	// Vérifier que tous les champs requis sont présents
	if err := fv.validateRequiredFields(fact, typeDef); err != nil {
		return err
	}
	
	// Vérifier que les champs fournis sont valides
	if err := fv.validateFieldDefinitions(fact, typeDef); err != nil {
		return err
	}
	
	// Vérifier les types des valeurs
	if err := fv.validateFieldValues(fact, typeDef); err != nil {
		return err
	}
	
	// Vérifier les clés primaires
	if err := ValidateFactPrimaryKey(fact, typeDef); err != nil {
		return err
	}
	
	return nil
}

// validateRequiredFields vérifie que tous les champs requis sont présents
func (fv *FactValidator) validateRequiredFields(fact Fact, typeDef TypeDefinition) error {
	// Créer un map des champs fournis
	providedFields := make(map[string]bool)
	for _, field := range fact.Fields {
		providedFields[field.Name] = true
	}
	
	// Vérifier chaque champ défini dans le type
	for _, fieldDef := range typeDef.Fields {
		if !providedFields[fieldDef.Name] {
			// Pour l'instant, tous les champs sont requis
			// TODO: Supporter les champs optionnels
			return fmt.Errorf(
				"fait de type '%s': champ requis '%s' manquant",
				fact.TypeName,
				fieldDef.Name,
			)
		}
	}
	
	return nil
}

// validateFieldDefinitions vérifie que les champs fournis sont définis
func (fv *FactValidator) validateFieldDefinitions(fact Fact, typeDef TypeDefinition) error {
	// Créer un map des champs définis
	definedFields := make(map[string]Field)
	for _, fieldDef := range typeDef.Fields {
		definedFields[fieldDef.Name] = fieldDef
	}
	
	// Vérifier chaque champ fourni
	for _, factField := range fact.Fields {
		// Interdire _id_
		if factField.Name == FieldNameInternalID {
			return fmt.Errorf(
				"fait de type '%s': le champ '%s' est réservé et ne peut pas être défini",
				fact.TypeName,
				FieldNameInternalID,
			)
		}
		
		// Vérifier que le champ est défini dans le type
		if _, exists := definedFields[factField.Name]; !exists {
			return fmt.Errorf(
				"fait de type '%s': champ '%s' non défini dans le type",
				fact.TypeName,
				factField.Name,
			)
		}
	}
	
	return nil
}

// validateFieldValues vérifie que les valeurs des champs ont le bon type
func (fv *FactValidator) validateFieldValues(fact Fact, typeDef TypeDefinition) error {
	// Créer un map pour lookup rapide
	fieldTypes := make(map[string]string)
	for _, fieldDef := range typeDef.Fields {
		fieldTypes[fieldDef.Name] = fieldDef.Type
	}
	
	for _, factField := range fact.Fields {
		expectedType := fieldTypes[factField.Name]
		
		if err := fv.validateFieldValue(factField, expectedType); err != nil {
			return fmt.Errorf(
				"fait de type '%s', champ '%s': %v",
				fact.TypeName,
				factField.Name,
				err,
			)
		}
	}
	
	return nil
}

// validateFieldValue valide une valeur de champ
func (fv *FactValidator) validateFieldValue(field FactField, expectedType string) error {
	value := field.Value
	
	// Cas 1: Référence de variable
	if value.Type == "variableReference" {
		varName, ok := value.Value.(string)
		if !ok {
			return fmt.Errorf("référence de variable invalide")
		}
		
		// Vérifier que la variable existe
		if !fv.typeSystem.VariableExists(varName) {
			return fmt.Errorf("variable '%s' non définie", varName)
		}
		
		// Vérifier que le type de la variable correspond
		varType, _ := fv.typeSystem.GetVariableType(varName)
		if varType != expectedType {
			return fmt.Errorf(
				"type incompatible: attendu '%s', la variable '%s' est de type '%s'",
				expectedType,
				varName,
				varType,
			)
		}
		
		return nil
	}
	
	// Cas 2: Valeur primitive
	return fv.validatePrimitiveValue(value, expectedType)
}

// validatePrimitiveValue valide une valeur primitive
func (fv *FactValidator) validatePrimitiveValue(value FactValue, expectedType string) error {
	// Mapping des types de valeurs vers les types de champs
	typeMapping := map[string][]string{
		ValueTypeString:  {"string"},
		ValueTypeNumber:  {"number"},
		ValueTypeBoolean: {"bool", "boolean"},
		ValueTypeBool:    {"bool", "boolean"},
	}
	
	validTypes, exists := typeMapping[value.Type]
	if !exists {
		return fmt.Errorf("type de valeur '%s' non supporté", value.Type)
	}
	
	// Vérifier que le type attendu est dans les types valides
	for _, validType := range validTypes {
		if expectedType == validType {
			return nil
		}
	}
	
	return fmt.Errorf(
		"type incompatible: attendu '%s', reçu '%s'",
		expectedType,
		value.Type,
	)
}
```

#### Tests du Validateur de Faits

**Fichier : `constraint/fact_validator_test.go`**

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"strings"
	"testing"
)

func TestFactValidator_ValidateFact(t *testing.T) {
	t.Log("🧪 TEST FACT VALIDATOR - VALIDATION COMPLÈTE")
	t.Log("=============================================")
	
	types := []TypeDefinition{
		{
			Name: "User",
			Fields: []Field{
				{Name: "name", Type: "string", IsPrimaryKey: true},
				{Name: "age", Type: "number"},
			},
		},
		{
			Name: "Login",
			Fields: []Field{
				{Name: "user", Type: "User"},
				{Name: "email", Type: "string", IsPrimaryKey: true},
				{Name: "password", Type: "string"},
			},
		},
	}
	
	ts := NewTypeSystem(types)
	ts.RegisterVariable("alice", "User")
	
	validator := NewFactValidator(ts)
	
	tests := []struct {
		name    string
		fact    Fact
		wantErr bool
		errMsg  string
	}{
		{
			name: "fait valide avec primitifs",
			fact: Fact{
				TypeName: "User",
				Fields: []FactField{
					{Name: "name", Value: FactValue{Type: "string", Value: "Alice"}},
					{Name: "age", Value: FactValue{Type: "number", Value: 30.0}},
				},
			},
			wantErr: false,
		},
		{
			name: "fait valide avec variable",
			fact: Fact{
				TypeName: "Login",
				Fields: []FactField{
					{Name: "user", Value: FactValue{Type: "variableReference", Value: "alice"}},
					{Name: "email", Value: FactValue{Type: "string", Value: "alice@ex.com"}},
					{Name: "password", Value: FactValue{Type: "string", Value: "secret"}},
				},
			},
			wantErr: false,
		},
		{
			name: "type inexistant",
			fact: Fact{
				TypeName: "UnknownType",
				Fields:   []FactField{},
			},
			wantErr: true,
			errMsg:  "non défini",
		},
		{
			name: "champ manquant",
			fact: Fact{
				TypeName: "User",
				Fields: []FactField{
					{Name: "name", Value: FactValue{Type: "string", Value: "Alice"}},
					// age manquant
				},
			},
			wantErr: true,
			errMsg:  "manquant",
		},
		{
			name: "champ non défini",
			fact: Fact{
				TypeName: "User",
				Fields: []FactField{
					{Name: "name", Value: FactValue{Type: "string", Value: "Alice"}},
					{Name: "age", Value: FactValue{Type: "number", Value: 30.0}},
					{Name: "unknown", Value: FactValue{Type: "string", Value: "test"}},
				},
			},
			wantErr: true,
			errMsg:  "non défini dans le type",
		},
		{
			name: "champ _id_ interdit",
			fact: Fact{
				TypeName: "User",
				Fields: []FactField{
					{Name: FieldNameInternalID, Value: FactValue{Type: "string", Value: "manual"}},
					{Name: "name", Value: FactValue{Type: "string", Value: "Alice"}},
					{Name: "age", Value: FactValue{Type: "number", Value: 30.0}},
				},
			},
			wantErr: true,
			errMsg:  "réservé",
		},
		{
			name: "type de valeur incompatible",
			fact: Fact{
				TypeName: "User",
				Fields: []FactField{
					{Name: "name", Value: FactValue{Type: "string", Value: "Alice"}},
					{Name: "age", Value: FactValue{Type: "string", Value: "thirty"}}, // devrait être number
				},
			},
			wantErr: true,
			errMsg:  "type incompatible",
		},
		{
			name: "variable non définie",
			fact: Fact{
				TypeName: "Login",
				Fields: []FactField{
					{Name: "user", Value: FactValue{Type: "variableReference", Value: "bob"}}, // bob non défini
					{Name: "email", Value: FactValue{Type: "string", Value: "test@ex.com"}},
					{Name: "password", Value: FactValue{Type: "string", Value: "secret"}},
				},
			},
			wantErr: true,
			errMsg:  "non définie",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateFact(tt.fact)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("❌ Attendu une erreur, reçu nil")
				} else {
					if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
						t.Errorf("❌ Message d'erreur attendu contenant '%s', reçu: %v", tt.errMsg, err)
					} else {
						t.Logf("✅ Erreur attendue: %v", err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("❌ Erreur inattendue: %v", err)
				} else {
					t.Logf("✅ Validation réussie")
				}
			}
		})
	}
}
```

### 4. Validation des Programmes Complets

#### Fichier : `constraint/program_validator.go` (nouveau)

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"fmt"
)

// ProgramValidator valide un programme TSD complet
type ProgramValidator struct {
	typeSystem    *TypeSystem
	factValidator *FactValidator
}

// NewProgramValidator crée un nouveau validateur de programme
func NewProgramValidator() *ProgramValidator {
	return &ProgramValidator{}
}

// Validate valide un programme complet
func (pv *ProgramValidator) Validate(program Program) error {
	// 1. Valider les définitions de types
	if err := pv.validateTypeDefinitions(program.Types); err != nil {
		return fmt.Errorf("validation des types: %v", err)
	}
	
	// 2. Créer le système de types
	pv.typeSystem = NewTypeSystem(program.Types)
	pv.factValidator = NewFactValidator(pv.typeSystem)
	
	// 3. Vérifier les références circulaires
	if err := pv.typeSystem.ValidateCircularReferences(); err != nil {
		return fmt.Errorf("validation des types: %v", err)
	}
	
	// 4. Valider les affectations de variables
	if err := pv.validateFactAssignments(program.FactAssignments); err != nil {
		return fmt.Errorf("validation des affectations: %v", err)
	}
	
	// 5. Valider les faits
	if err := pv.validateFacts(program.Facts); err != nil {
		return fmt.Errorf("validation des faits: %v", err)
	}
	
	// 6. Valider les expressions/règles
	if err := pv.validateExpressions(program.Expressions); err != nil {
		return fmt.Errorf("validation des expressions: %v", err)
	}
	
	return nil
}

// validateTypeDefinitions valide toutes les définitions de types
func (pv *ProgramValidator) validateTypeDefinitions(types []TypeDefinition) error {
	for i, typeDef := range types {
		if err := ValidateTypeDefinition(typeDef); err != nil {
			return fmt.Errorf("type %d ('%s'): %v", i+1, typeDef.Name, err)
		}
	}
	return nil
}

// validateFactAssignments valide les affectations de variables
func (pv *ProgramValidator) validateFactAssignments(assignments []FactAssignment) error {
	for i, assignment := range assignments {
		// Valider le fait
		if err := pv.factValidator.ValidateFact(assignment.Fact); err != nil {
			return fmt.Errorf("affectation %d (variable '%s'): %v", i+1, assignment.Variable, err)
		}
		
		// Enregistrer la variable
		if err := pv.typeSystem.RegisterVariable(assignment.Variable, assignment.Fact.TypeName); err != nil {
			return fmt.Errorf("affectation %d: %v", i+1, err)
		}
	}
	
	return nil
}

// validateFacts valide tous les faits
func (pv *ProgramValidator) validateFacts(facts []Fact) error {
	for i, fact := range facts {
		if err := pv.factValidator.ValidateFact(fact); err != nil {
			return fmt.Errorf("fait %d: %v", i+1, err)
		}
	}
	
	return nil
}

// validateExpressions valide toutes les expressions/règles
func (pv *ProgramValidator) validateExpressions(expressions []Expression) error {
	for i, expr := range expressions {
		if err := pv.validateExpression(expr); err != nil {
			return fmt.Errorf("expression %d (règle '%s'): %v", i+1, expr.RuleId, err)
		}
	}
	
	return nil
}

// validateExpression valide une expression
func (pv *ProgramValidator) validateExpression(expr Expression) error {
	// Créer un contexte de variables pour cette expression
	varTypes := make(map[string]string)
	
	// Collecter les variables des patterns
	patterns := expr.Patterns
	if len(patterns) == 0 && expr.Set.Type == "set" {
		patterns = []Set{expr.Set}
	}
	
	for _, pattern := range patterns {
		for _, variable := range pattern.Variables {
			if variable.Type == "typedVariable" {
				varTypes[variable.Name] = variable.DataType
			}
		}
	}
	
	// Valider les contraintes
	if expr.Constraints != nil {
		if err := pv.validateConstraints(expr.Constraints, varTypes); err != nil {
			return err
		}
	}
	
	return nil
}

// validateConstraints valide les contraintes d'une expression
func (pv *ProgramValidator) validateConstraints(constraints interface{}, varTypes map[string]string) error {
	switch c := constraints.(type) {
	case Constraint:
		return pv.validateConstraint(c, varTypes)
		
	case BinaryOperation:
		return pv.validateBinaryOperation(c, varTypes)
		
	case LogicalExpression:
		return pv.validateLogicalExpression(c, varTypes)
		
	// Autres types de contraintes...
	default:
		// Type non reconnu, ignorer pour l'instant
		return nil
	}
}

// validateConstraint valide une contrainte
func (pv *ProgramValidator) validateConstraint(constraint Constraint, varTypes map[string]string) error {
	// Valider selon le type de contrainte
	if constraint.Operator != "" {
		// C'est une comparaison
		return pv.validateComparison(constraint.Left, constraint.Right, constraint.Operator, varTypes)
	}
	
	return nil
}

// validateBinaryOperation valide une opération binaire
func (pv *ProgramValidator) validateBinaryOperation(op BinaryOperation, varTypes map[string]string) error {
	if op.Type == "comparison" {
		return pv.validateComparison(op.Left, op.Right, op.Operator, varTypes)
	}
	
	// Pour les opérations arithmétiques, vérifier que les opérandes sont des nombres
	// TODO: Implémenter validation arithmétique
	
	return nil
}

// validateLogicalExpression valide une expression logique
func (pv *ProgramValidator) validateLogicalExpression(expr LogicalExpression, varTypes map[string]string) error {
	// Valider le côté gauche
	if err := pv.validateConstraints(expr.Left, varTypes); err != nil {
		return err
	}
	
	// Valider chaque opération
	for _, op := range expr.Operations {
		if err := pv.validateConstraints(op.Right, varTypes); err != nil {
			return err
		}
	}
	
	return nil
}

// validateComparison valide une comparaison
func (pv *ProgramValidator) validateComparison(left, right interface{}, operator string, varTypes map[string]string) error {
	// Inférer le type de gauche
	leftType, err := pv.inferExpressionType(left, varTypes)
	if err != nil {
		return fmt.Errorf("expression gauche: %v", err)
	}
	
	// Inférer le type de droite
	rightType, err := pv.inferExpressionType(right, varTypes)
	if err != nil {
		return fmt.Errorf("expression droite: %v", err)
	}
	
	// Vérifier la compatibilité
	if !pv.typeSystem.AreTypesCompatible(leftType, rightType, operator) {
		return fmt.Errorf(
			"types incompatibles pour comparaison %s: '%s' et '%s'",
			operator,
			leftType,
			rightType,
		)
	}
	
	return nil
}

// inferExpressionType infère le type d'une expression
func (pv *ProgramValidator) inferExpressionType(expr interface{}, varTypes map[string]string) (string, error) {
	switch e := expr.(type) {
	case FieldAccess:
		// Type de la variable
		varType, exists := varTypes[e.Object]
		if !exists {
			return "", fmt.Errorf("variable '%s' non définie dans cette règle", e.Object)
		}
		
		// Type du champ
		fieldType, err := pv.typeSystem.GetFieldType(varType, e.Field)
		if err != nil {
			return "", err
		}
		
		return fieldType, nil
		
	case Variable:
		// Type de la variable
		varType, exists := varTypes[e.Name]
		if !exists {
			return "", fmt.Errorf("variable '%s' non définie dans cette règle", e.Name)
		}
		
		return varType, nil
		
	case StringLiteral:
		return "string", nil
		
	case NumberLiteral:
		return "number", nil
		
	case BooleanLiteral:
		return "bool", nil
		
	default:
		return "", fmt.Errorf("type d'expression non supporté: %T", expr)
	}
}
```

#### Tests du Validateur de Programme

**Fichier : `constraint/program_validator_test.go`**

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"strings"
	"testing"
)

func TestProgramValidator_ValidProgram(t *testing.T) {
	t.Log("🧪 TEST PROGRAM VALIDATOR - PROGRAMME VALIDE")
	t.Log("=============================================")
	
	program := Program{
		Types: []TypeDefinition{
			{
				Name: "User",
				Fields: []Field{
					{Name: "name", Type: "string", IsPrimaryKey: true},
					{Name: "age", Type: "number"},
				},
			},
			{
				Name: "Login",
				Fields: []Field{
					{Name: "user", Type: "User"},
					{Name: "email", Type: "string", IsPrimaryKey: true},
					{Name: "password", Type: "string"},
				},
			},
		},
		FactAssignments: []FactAssignment{
			{
				Variable: "alice",
				Fact: Fact{
					TypeName: "User",
					Fields: []FactField{
						{Name: "name", Value: FactValue{Type: "string", Value: "Alice"}},
						{Name: "age", Value: FactValue{Type: "number", Value: 30.0}},
					},
				},
			},
		},
		Facts: []Fact{
			{
				TypeName: "Login",
				Fields: []FactField{
					{Name: "user", Value: FactValue{Type: "variableReference", Value: "alice"}},
					{Name: "email", Value: FactValue{Type: "string", Value: "alice@ex.com"}},
					{Name: "password", Value: FactValue{Type: "string", Value: "secret"}},
				},
			},
		},
	}
	
	validator := NewProgramValidator()
	err := validator.Validate(program)
	
	if err != nil {
		t.Errorf("❌ Erreur inattendue: %v", err)
	} else {
		t.Log("✅ Programme valide")
	}
}

func TestProgramValidator_InvalidPrograms(t *testing.T) {
	t.Log("🧪 TEST PROGRAM VALIDATOR - PROGRAMMES INVALIDES")
	t.Log("=================================================")
	
	tests := []struct {
		name    string
		program Program
		errMsg  string
	}{
		{
			name: "référence circulaire",
			program: Program{
				Types: []TypeDefinition{
					{
						Name: "A",
						Fields: []Field{
							{Name: "b", Type: "B"},
						},
					},
					{
						Name: "B",
						Fields: []Field{
							{Name: "a", Type: "A"},
						},
					},
				},
			},
			errMsg: "circulaire",
		},
		{
			name: "type inexistant",
			program: Program{
				Types: []TypeDefinition{
					{
						Name: "Login",
						Fields: []Field{
							{Name: "user", Type: "UnknownType"},
						},
					},
				},
			},
			errMsg: "UnknownType",
		},
		{
			name: "variable non définie",
			program: Program{
				Types: []TypeDefinition{
					{
						Name: "Login",
						Fields: []Field{
							{Name: "user", Type: "User"},
							{Name: "email", Type: "string", IsPrimaryKey: true},
							{Name: "password", Type: "string"},
						},
					},
					{
						Name: "User",
						Fields: []Field{
							{Name: "name", Type: "string", IsPrimaryKey: true},
						},
					},
				},
				Facts: []Fact{
					{
						TypeName: "Login",
						Fields: []FactField{
							{Name: "user", Value: FactValue{Type: "variableReference", Value: "unknown"}},
							{Name: "email", Value: FactValue{Type: "string", Value: "test@ex.com"}},
							{Name: "password", Value: FactValue{Type: "string", Value: "pw"}},
						},
					},
				},
			},
			errMsg: "non définie",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := NewProgramValidator()
			err := validator.Validate(tt.program)
			
			if err == nil {
				t.Errorf("❌ Attendu une erreur, reçu nil")
			} else {
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("❌ Erreur attendue contenant '%s', reçu: %v", tt.errMsg, err)
				} else {
					t.Logf("✅ Erreur détectée: %v", err)
				}
			}
		})
	}
}
```

### 5. Intégration dans l'API

#### Fichier : `constraint/api.go` (modifications)

**Ajouter validation automatique** :

```go
// ParseAndValidateProgram parse et valide un programme complet
func ParseAndValidateProgram(input string) (*Program, error) {
	// Parser le programme
	program, err := ParseProgram(input)
	if err != nil {
		return nil, fmt.Errorf("erreur de parsing: %v", err)
	}
	
	// Valider le programme
	validator := NewProgramValidator()
	if err := validator.Validate(*program); err != nil {
		return nil, fmt.Errorf("erreur de validation: %v", err)
	}
	
	return program, nil
}
```

### 6. Messages d'Erreur Améliorés

#### Fichier : `constraint/errors.go` (modifications)

**Ajouter types d'erreurs structurés** :

```go
// ValidationError représente une erreur de validation avec contexte
type ValidationError struct {
	Type    string // "type", "fact", "expression", etc.
	Element string // Nom de l'élément (type name, variable name, etc.)
	Field   string // Champ concerné si applicable
	Message string // Message d'erreur
}

func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("%s '%s', champ '%s': %s", e.Type, e.Element, e.Field, e.Message)
	}
	return fmt.Sprintf("%s '%s': %s", e.Type, e.Element, e.Message)
}

// NewTypeValidationError crée une erreur de validation de type
func NewTypeValidationError(typeName, fieldName, message string) *ValidationError {
	return &ValidationError{
		Type:    "type",
		Element: typeName,
		Field:   fieldName,
		Message: message,
	}
}

// NewFactValidationError crée une erreur de validation de fait
func NewFactValidationError(factType, fieldName, message string) *ValidationError {
	return &ValidationError{
		Type:    "fait",
		Element: factType,
		Field:   fieldName,
		Message: message,
	}
}

// NewExpressionValidationError crée une erreur de validation d'expression
func NewExpressionValidationError(ruleId, message string) *ValidationError {
	return &ValidationError{
		Type:    "expression",
		Element: ruleId,
		Message: message,
	}
}
```

---

## ✅ Critères de Succès

### Compilation et Tests

```bash
# Code compile
go build ./constraint

# Tests passent
go test ./constraint -run TestTypeSystem -v
go test ./constraint -run TestFactValidator -v
go test ./constraint -run TestProgramValidator -v

# Couverture > 80%
go test ./constraint -cover
```

### Fonctionnalités

- [ ] `TypeSystem` créé et fonctionnel
- [ ] `FactValidator` créé et fonctionnel
- [ ] `ProgramValidator` créé et fonctionnel
- [ ] Validation complète des types
- [ ] Validation complète des faits
- [ ] Validation complète des expressions
- [ ] Messages d'erreur clairs et informatifs
- [ ] Détection de références circulaires
- [ ] Validation de compatibilité de types

### Validation

```bash
make format
make lint
make validate
make test-complete
```

---

## 📊 Tests Requis

### Tests Unitaires Minimaux

- [ ] `TestTypeSystem_TypeChecks`
- [ ] `TestTypeSystem_GetFieldType`
- [ ] `TestTypeSystem_Variables`
- [ ] `TestTypeSystem_CircularReferences`
- [ ] `TestTypeSystem_TypeCompatibility`
- [ ] `TestFactValidator_ValidateFact`
- [ ] `TestProgramValidator_ValidProgram`
- [ ] `TestProgramValidator_InvalidPrograms`

### Tests d'Intégration

```go
func TestCompleteValidation_EndToEnd(t *testing.T) {
	input := `
		type User(#name: string, age: number)
		type Login(user: User, #email: string, password: string)
		
		alice = User("Alice", 30)
		bob = User("Bob", 25)
		
		Login(alice, "alice@ex.com", "pw1")
		Login(bob, "bob@ex.com", "pw2")
		
		{u: User, l: Login} / l.user == u && u.age > 25 ==> 
			Log("Senior user login: " + u.name)
	`
	
	program, err := ParseAndValidateProgram(input)
	if err != nil {
		t.Fatalf("Erreur: %v", err)
	}
	
	if program == nil {
		t.Fatal("Programme nil")
	}
	
	t.Log("✅ Validation complète réussie")
}
```

---

## 🚀 Exécution

### Ordre des Modifications

1. ✅ Créer `TypeSystem`
2. ✅ Créer `FactValidator`
3. ✅ Créer `ProgramValidator`
4. ✅ Améliorer messages d'erreur
5. ✅ Intégrer dans API
6. ✅ Tests unitaires complets
7. ✅ Tests d'intégration
8. ✅ Validation finale

### Commandes

```bash
# Créer les fichiers
touch constraint/type_system.go
touch constraint/type_system_test.go
touch constraint/fact_validator.go
touch constraint/fact_validator_test.go
touch constraint/program_validator.go
touch constraint/program_validator_test.go

# Tester progressivement
go test ./constraint -run TestTypeSystem -v
go test ./constraint -run TestFactValidator -v
go test ./constraint -run TestProgramValidator -v

# Validation complète
make validate
make test-complete
```

---

## 📚 Références

- `scripts/new_ids/04-prompt-evaluation.md` - Évaluation
- `scripts/new_ids/03-prompt-id-generation.md` - Génération IDs
- `constraint/constraint_validation.go` - Validation actuelle
- `docs/primary-keys.md` - Documentation

---

## 📝 Notes

### Points d'Attention

1. **Performance** : Le système de types est utilisé fréquemment, optimiser si nécessaire

2. **Messages d'erreur** : Doivent être clairs et pointer vers la source du problème

3. **Validation exhaustive** : Mieux vaut une erreur à la compilation qu'à l'exécution

4. **Références circulaires** : Algorithme DFS pour détecter les cycles

### Questions Résolues

Q: Faut-il supporter les champs optionnels ?
R: Pas pour l'instant, à ajouter plus tard si besoin

Q: Comment gérer les auto-références (ex: Tree) ?
R: Interdire pour l'instant via détection de cycles

---

## 🎯 Résultat Attendu

```tsd
// ✅ Validation complète
type User(#name: string, age: number)
type Login(user: User, #email: string)

alice = User("Alice", 30)
Login(alice, "alice@ex.com")

// ❌ Erreurs détectées clairement
type Bad(_id_: string)             // Erreur: champ '_id_' réservé
alice = User("Alice", "invalid")   // Erreur: age doit être number
Login(unknown, "test")             // Erreur: variable 'unknown' non définie
Login("string", "test")            // Erreur: type incompatible pour user

// Références circulaires détectées
type A(b: B)
type B(a: A)  // Erreur: référence circulaire A -> B -> A
```

---

**Prompt suivant** : `06-prompt-api-tsdio.md`

**Durée estimée** : 4-6 heures

**Complexité** : ⚠️ Moyenne-Élevée