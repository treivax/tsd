# Prompt 04 - Modification de l'Évaluation et Comparaisons

> **📋 Standards** : Ce prompt respecte les règles de [.github/prompts/common.md](../../.github/prompts/common.md) et [.github/prompts/develop.md](../../.github/prompts/develop.md)

---

## 🎯 Objectif

Modifier le moteur d'évaluation RETE pour supporter les comparaisons simplifiées de faits :
- `l.user == u` au lieu de `l.user == u.user`
- Résolution automatique des types
- Comparaison via `_id_` interne
- Validation de compatibilité des types

---

## 📋 Contexte

### État Actuel

```tsd
type User(#name: string, age: number)
type Login(userEmail: string, #email: string, password: string)

// Comparaison explicite sur champs primitifs
{u: User, l: Login} / l.userEmail == u.email ==> ...
```

### État Cible

```tsd
type User(#name: string, age: number)
type Login(user: User, #email: string, password: string)

// Comparaison directe de faits
{u: User, l: Login} / l.user == u ==> 
    Log("Login for " + u.name)
```

**Comportement interne** :
- `l.user` retourne l'ID interne du fait référencé : `"User~Alice"`
- `u` est résolu vers son ID interne : `"User~Alice"`
- La comparaison devient : `"User~Alice" == "User~Alice"` → `true`

---

## 📝 Tâches à Réaliser

### 1. Analyser l'Évaluation Actuelle

#### Identifier les Fichiers Critiques

```bash
# Rechercher les fichiers d'évaluation dans RETE
find rete/ -name "*.go" | grep -i "eval\|node\|comparison"

# Rechercher les fonctions d'évaluation de contraintes
grep -r "EvaluateConstraint\|CompareValues\|FieldAccess" rete/ --include="*.go"
```

**Fichiers attendus** :
- `rete/node_alpha.go` - Nœuds alpha (filtrage)
- `rete/node_beta.go` - Nœuds beta (jointures)
- `rete/node_join.go` - Logique de jointure
- `rete/evaluator.go` ou similaire - Évaluation des conditions
- `rete/fact_token.go` - Structures de faits et tokens

#### Questions à Répondre

1. Comment les field access sont-ils actuellement évalués ?
2. Où se fait la comparaison de valeurs ?
3. Comment les types sont-ils vérifiés lors des comparaisons ?
4. Où sont résolus les champs d'un fait ?

### 2. Comprendre la Résolution de Champs

#### Fichier : Rechercher `FieldAccess` dans RETE

**Fonction typique actuelle** :
```go
// Exemple hypothétique basé sur l'architecture RETE
func resolveFieldAccess(fact *Fact, fieldName string) (interface{}, error) {
    // Accéder au champ dans Fields
    if value, exists := fact.Fields[fieldName]; exists {
        return value, nil
    }
    return nil, fmt.Errorf("champ '%s' non trouvé dans le fait de type '%s'", fieldName, fact.Type)
}
```

**Problème** : Actuellement, on retourne la valeur brute. Pour un champ de type `User`, on a probablement l'ID stocké, mais il faut le gérer correctement.

### 3. Créer un Résolveur de Valeurs Typées

#### Nouveau Fichier : `rete/field_resolver.go`

**Objectif** : Résoudre les valeurs de champs en tenant compte de leur type

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
    "fmt"
    "github.com/resinsec/tsd/constraint"
)

// FieldResolver résout les valeurs de champs en tenant compte de leur type
type FieldResolver struct {
    // TypeMap contient les définitions de types pour la résolution
    TypeMap map[string]constraint.TypeDefinition
}

// NewFieldResolver crée un nouveau résolveur de champs
func NewFieldResolver(types []constraint.TypeDefinition) *FieldResolver {
    typeMap := make(map[string]constraint.TypeDefinition)
    for _, t := range types {
        typeMap[t.Name] = t
    }
    
    return &FieldResolver{
        TypeMap: typeMap,
    }
}

// ResolveFieldValue résout la valeur d'un champ d'un fait
// Pour les types primitifs, retourne la valeur directement
// Pour les types de faits, retourne l'ID interne du fait référencé
func (fr *FieldResolver) ResolveFieldValue(fact *Fact, fieldName string) (interface{}, string, error) {
    // Le champ _id_ est interdit
    if fieldName == constraint.FieldNameInternalID {
        return nil, "", fmt.Errorf("le champ '%s' est interne et ne peut pas être accédé", constraint.FieldNameInternalID)
    }
    
    // Vérifier que le champ existe dans le fait
    value, exists := fact.Fields[fieldName]
    if !exists {
        return nil, "", fmt.Errorf("champ '%s' non trouvé dans le fait de type '%s'", fieldName, fact.Type)
    }
    
    // Obtenir la définition du type pour connaître le type du champ
    typeDef, exists := fr.TypeMap[fact.Type]
    if !exists {
        return nil, "", fmt.Errorf("type '%s' non trouvé dans le type map", fact.Type)
    }
    
    // Trouver le champ dans la définition du type
    var fieldDef constraint.Field
    found := false
    for _, f := range typeDef.Fields {
        if f.Name == fieldName {
            fieldDef = f
            found = true
            break
        }
    }
    
    if !found {
        return nil, "", fmt.Errorf("champ '%s' non défini dans le type '%s'", fieldName, fact.Type)
    }
    
    // Déterminer le type du champ
    fieldType := fr.getFieldType(fieldDef.Type)
    
    return value, fieldType, nil
}

// getFieldType retourne le type d'un champ (primitive ou user-defined)
func (fr *FieldResolver) getFieldType(typeName string) string {
    primitives := map[string]bool{
        "string":  true,
        "number":  true,
        "bool":    true,
        "boolean": true,
    }
    
    if primitives[typeName] {
        return "primitive"
    }
    
    // Vérifier si c'est un type utilisateur défini
    if _, exists := fr.TypeMap[typeName]; exists {
        return "fact"
    }
    
    return "unknown"
}

// ResolveFactID résout une variable de fait vers son ID interne
func (fr *FieldResolver) ResolveFactID(fact *Fact) string {
    return fact.ID
}
```

#### Tests du Résolveur

**Nouveau Fichier : `rete/field_resolver_test.go`**

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
    "strings"
    "testing"
    "github.com/resinsec/tsd/constraint"
)

func TestFieldResolver_ResolveFieldValue(t *testing.T) {
    t.Log("🧪 TEST FIELD RESOLVER - RÉSOLUTION VALEURS")
    t.Log("============================================")
    
    types := []constraint.TypeDefinition{
        {
            Name: "User",
            Fields: []constraint.Field{
                {Name: "name", Type: "string", IsPrimaryKey: true},
                {Name: "age", Type: "number"},
            },
        },
        {
            Name: "Login",
            Fields: []constraint.Field{
                {Name: "user", Type: "User"},
                {Name: "email", Type: "string", IsPrimaryKey: true},
                {Name: "password", Type: "string"},
            },
        },
    }
    
    resolver := NewFieldResolver(types)
    
    tests := []struct {
        name          string
        fact          *Fact
        fieldName     string
        expectedValue interface{}
        expectedType  string
        wantErr       bool
    }{
        {
            name: "champ primitif string",
            fact: &Fact{
                ID:   "Login~alice@ex.com",
                Type: "Login",
                Fields: map[string]interface{}{
                    "user":     "User~Alice",
                    "email":    "alice@ex.com",
                    "password": "secret",
                },
            },
            fieldName:     "email",
            expectedValue: "alice@ex.com",
            expectedType:  "primitive",
            wantErr:       false,
        },
        {
            name: "champ de type fait",
            fact: &Fact{
                ID:   "Login~alice@ex.com",
                Type: "Login",
                Fields: map[string]interface{}{
                    "user":     "User~Alice",
                    "email":    "alice@ex.com",
                    "password": "secret",
                },
            },
            fieldName:     "user",
            expectedValue: "User~Alice",
            expectedType:  "fact",
            wantErr:       false,
        },
        {
            name: "champ _id_ interdit",
            fact: &Fact{
                ID:   "User~Alice",
                Type: "User",
                Fields: map[string]interface{}{
                    "name": "Alice",
                    "age":  30.0,
                },
            },
            fieldName: constraint.FieldNameInternalID,
            wantErr:   true,
        },
        {
            name: "champ non existant",
            fact: &Fact{
                ID:   "User~Alice",
                Type: "User",
                Fields: map[string]interface{}{
                    "name": "Alice",
                    "age":  30.0,
                },
            },
            fieldName: "unknown",
            wantErr:   true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            value, fieldType, err := resolver.ResolveFieldValue(tt.fact, tt.fieldName)
            
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
            
            if value != tt.expectedValue {
                t.Errorf("❌ Valeur attendue %v, reçu %v", tt.expectedValue, value)
            }
            
            if fieldType != tt.expectedType {
                t.Errorf("❌ Type attendu '%s', reçu '%s'", tt.expectedType, fieldType)
            }
            
            t.Logf("✅ Résolution correcte: valeur=%v, type=%s", value, fieldType)
        })
    }
}

func TestFieldResolver_ResolveFactID(t *testing.T) {
    t.Log("🧪 TEST FIELD RESOLVER - RÉSOLUTION ID")
    t.Log("=======================================")
    
    resolver := NewFieldResolver(nil)
    
    fact := &Fact{
        ID:   "User~Alice",
        Type: "User",
        Fields: map[string]interface{}{
            "name": "Alice",
            "age":  30.0,
        },
    }
    
    id := resolver.ResolveFactID(fact)
    
    if id != "User~Alice" {
        t.Errorf("❌ ID attendu 'User~Alice', reçu '%s'", id)
    } else {
        t.Logf("✅ ID résolu correctement: %s", id)
    }
}
```

### 4. Modifier l'Évaluateur de Comparaisons

#### Fichier : Rechercher l'évaluateur actuel dans RETE

**Recherche** :
```bash
grep -r "CompareValues\|EvaluateComparison\|EvaluateBinary" rete/ --include="*.go" -A 5
```

**Fonction typique à modifier** :

```go
// Avant (exemple hypothétique)
func evaluateComparison(left, right interface{}, operator string) (bool, error) {
    switch operator {
    case "==":
        return left == right, nil
    case "!=":
        return left != right, nil
    // ...
    }
}
```

#### Nouveau : Évaluateur avec Support de Types

**Fichier : `rete/comparison_evaluator.go` (nouveau ou modification)**

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
    "fmt"
    "math"
)

// ComparisonEvaluator évalue les comparaisons entre valeurs avec support des types
type ComparisonEvaluator struct {
    resolver *FieldResolver
}

// NewComparisonEvaluator crée un nouvel évaluateur de comparaisons
func NewComparisonEvaluator(resolver *FieldResolver) *ComparisonEvaluator {
    return &ComparisonEvaluator{
        resolver: resolver,
    }
}

// EvaluateComparison évalue une comparaison entre deux valeurs
// Gère les comparaisons de primitifs ET les comparaisons de faits (via IDs)
func (ce *ComparisonEvaluator) EvaluateComparison(left, right interface{}, operator string, leftType, rightType string) (bool, error) {
    // Cas 1: Les deux valeurs sont des IDs de faits
    if leftType == "fact" && rightType == "fact" {
        return ce.compareFactIDs(left, right, operator)
    }
    
    // Cas 2: Les deux valeurs sont des primitifs
    if leftType == "primitive" && rightType == "primitive" {
        return ce.comparePrimitives(left, right, operator)
    }
    
    // Cas 3: Types incompatibles
    return false, fmt.Errorf("comparaison impossible entre types '%s' et '%s'", leftType, rightType)
}

// compareFactIDs compare deux IDs de faits
func (ce *ComparisonEvaluator) compareFactIDs(left, right interface{}, operator string) (bool, error) {
    leftID, ok1 := left.(string)
    rightID, ok2 := right.(string)
    
    if !ok1 || !ok2 {
        return false, fmt.Errorf("IDs de faits doivent être des strings")
    }
    
    switch operator {
    case "==":
        return leftID == rightID, nil
    case "!=":
        return leftID != rightID, nil
    default:
        return false, fmt.Errorf("opérateur '%s' non supporté pour les comparaisons de faits (seuls == et != sont autorisés)", operator)
    }
}

// comparePrimitives compare deux valeurs primitives
func (ce *ComparisonEvaluator) comparePrimitives(left, right interface{}, operator string) (bool, error) {
    // Essayer de comparer comme strings
    leftStr, leftIsStr := left.(string)
    rightStr, rightIsStr := right.(string)
    
    if leftIsStr && rightIsStr {
        return ce.compareStrings(leftStr, rightStr, operator)
    }
    
    // Essayer de comparer comme numbers
    leftNum, leftIsNum := toFloat64(left)
    rightNum, rightIsNum := toFloat64(right)
    
    if leftIsNum && rightIsNum {
        return ce.compareNumbers(leftNum, rightNum, operator)
    }
    
    // Essayer de comparer comme booleans
    leftBool, leftIsBool := left.(bool)
    rightBool, rightIsBool := right.(bool)
    
    if leftIsBool && rightIsBool {
        return ce.compareBooleans(leftBool, rightBool, operator)
    }
    
    // Types incompatibles
    return false, fmt.Errorf("types incompatibles pour comparaison: %T et %T", left, right)
}

// compareStrings compare deux strings
func (ce *ComparisonEvaluator) compareStrings(left, right, operator string) (bool, error) {
    switch operator {
    case "==":
        return left == right, nil
    case "!=":
        return left != right, nil
    case "<":
        return left < right, nil
    case "<=":
        return left <= right, nil
    case ">":
        return left > right, nil
    case ">=":
        return left >= right, nil
    default:
        return false, fmt.Errorf("opérateur '%s' non supporté pour strings", operator)
    }
}

// compareNumbers compare deux numbers
func (ce *ComparisonEvaluator) compareNumbers(left, right float64, operator string) (bool, error) {
    const epsilon = 1e-9
    
    switch operator {
    case "==":
        return math.Abs(left-right) < epsilon, nil
    case "!=":
        return math.Abs(left-right) >= epsilon, nil
    case "<":
        return left < right, nil
    case "<=":
        return left <= right, nil
    case ">":
        return left > right, nil
    case ">=":
        return left >= right, nil
    default:
        return false, fmt.Errorf("opérateur '%s' non supporté pour numbers", operator)
    }
}

// compareBooleans compare deux booleans
func (ce *ComparisonEvaluator) compareBooleans(left, right bool, operator string) (bool, error) {
    switch operator {
    case "==":
        return left == right, nil
    case "!=":
        return left != right, nil
    default:
        return false, fmt.Errorf("opérateur '%s' non supporté pour booleans (seuls == et != sont autorisés)", operator)
    }
}

// toFloat64 convertit une valeur en float64 si possible
func toFloat64(v interface{}) (float64, bool) {
    switch val := v.(type) {
    case float64:
        return val, true
    case int:
        return float64(val), true
    case int64:
        return float64(val), true
    case int32:
        return float64(val), true
    case float32:
        return float64(val), true
    default:
        return 0, false
    }
}
```

#### Tests de l'Évaluateur

**Fichier : `rete/comparison_evaluator_test.go`**

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
    "testing"
)

func TestComparisonEvaluator_CompareFactIDs(t *testing.T) {
    t.Log("🧪 TEST COMPARISON EVALUATOR - IDS DE FAITS")
    t.Log("============================================")
    
    evaluator := NewComparisonEvaluator(nil)
    
    tests := []struct {
        name     string
        left     interface{}
        right    interface{}
        operator string
        expected bool
        wantErr  bool
    }{
        {
            name:     "IDs égaux avec ==",
            left:     "User~Alice",
            right:    "User~Alice",
            operator: "==",
            expected: true,
            wantErr:  false,
        },
        {
            name:     "IDs différents avec ==",
            left:     "User~Alice",
            right:    "User~Bob",
            operator: "==",
            expected: false,
            wantErr:  false,
        },
        {
            name:     "IDs différents avec !=",
            left:     "User~Alice",
            right:    "User~Bob",
            operator: "!=",
            expected: true,
            wantErr:  false,
        },
        {
            name:     "IDs égaux avec !=",
            left:     "User~Alice",
            right:    "User~Alice",
            operator: "!=",
            expected: false,
            wantErr:  false,
        },
        {
            name:     "opérateur < non supporté pour faits",
            left:     "User~Alice",
            right:    "User~Bob",
            operator: "<",
            wantErr:  true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := evaluator.compareFactIDs(tt.left, tt.right, tt.operator)
            
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
            
            if result != tt.expected {
                t.Errorf("❌ Résultat attendu %v, reçu %v", tt.expected, result)
            } else {
                t.Logf("✅ Comparaison correcte: %v %s %v = %v", tt.left, tt.operator, tt.right, result)
            }
        })
    }
}

func TestComparisonEvaluator_ComparePrimitives(t *testing.T) {
    t.Log("🧪 TEST COMPARISON EVALUATOR - PRIMITIFS")
    t.Log("=========================================")
    
    evaluator := NewComparisonEvaluator(nil)
    
    tests := []struct {
        name     string
        left     interface{}
        right    interface{}
        operator string
        expected bool
        wantErr  bool
    }{
        // Strings
        {
            name:     "strings égaux",
            left:     "alice",
            right:    "alice",
            operator: "==",
            expected: true,
        },
        {
            name:     "strings différents",
            left:     "alice",
            right:    "bob",
            operator: "!=",
            expected: true,
        },
        {
            name:     "string < string",
            left:     "alice",
            right:    "bob",
            operator: "<",
            expected: true,
        },
        
        // Numbers
        {
            name:     "numbers égaux",
            left:     42.0,
            right:    42.0,
            operator: "==",
            expected: true,
        },
        {
            name:     "numbers différents",
            left:     10.0,
            right:    20.0,
            operator: "<",
            expected: true,
        },
        {
            name:     "int et float64",
            left:     int(42),
            right:    42.0,
            operator: "==",
            expected: true,
        },
        
        // Booleans
        {
            name:     "booleans égaux",
            left:     true,
            right:    true,
            operator: "==",
            expected: true,
        },
        {
            name:     "booleans différents",
            left:     true,
            right:    false,
            operator: "!=",
            expected: true,
        },
        {
            name:     "boolean < interdit",
            left:     true,
            right:    false,
            operator: "<",
            wantErr:  true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := evaluator.comparePrimitives(tt.left, tt.right, tt.operator)
            
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
            
            if result != tt.expected {
                t.Errorf("❌ Résultat attendu %v, reçu %v", tt.expected, result)
            } else {
                t.Logf("✅ Comparaison correcte: %v %s %v = %v", tt.left, tt.operator, tt.right, result)
            }
        })
    }
}

func TestComparisonEvaluator_EvaluateComparison(t *testing.T) {
    t.Log("🧪 TEST COMPARISON EVALUATOR - GLOBAL")
    t.Log("======================================")
    
    evaluator := NewComparisonEvaluator(nil)
    
    tests := []struct {
        name      string
        left      interface{}
        right     interface{}
        operator  string
        leftType  string
        rightType string
        expected  bool
        wantErr   bool
    }{
        {
            name:      "comparaison de faits",
            left:      "User~Alice",
            right:     "User~Alice",
            operator:  "==",
            leftType:  "fact",
            rightType: "fact",
            expected:  true,
        },
        {
            name:      "comparaison de primitifs",
            left:      "alice",
            right:     "alice",
            operator:  "==",
            leftType:  "primitive",
            rightType: "primitive",
            expected:  true,
        },
        {
            name:      "types incompatibles",
            left:      "User~Alice",
            right:     "alice",
            operator:  "==",
            leftType:  "fact",
            rightType: "primitive",
            wantErr:   true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := evaluator.EvaluateComparison(tt.left, tt.right, tt.operator, tt.leftType, tt.rightType)
            
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
            
            if result != tt.expected {
                t.Errorf("❌ Résultat attendu %v, reçu %v", tt.expected, result)
            } else {
                t.Logf("✅ Évaluation correcte")
            }
        })
    }
}
```

### 5. Intégrer dans les Nœuds RETE

#### Identifier les Nœuds de Jointure

**Recherche** :
```bash
grep -r "type.*Node.*struct" rete/ --include="*.go" | grep -i "join\|beta\|alpha"
```

**Fichiers typiques** :
- `rete/node_join.go` - Nœuds de jointure
- `rete/node_beta.go` - Nœuds beta
- `rete/network.go` - Construction du réseau

#### Modification Typique d'un Nœud de Jointure

**Avant (exemple hypothétique)** :
```go
type JoinNode struct {
    // ...
    leftBinding  string
    rightBinding string
    condition    Condition
}

func (jn *JoinNode) evaluate(leftToken, rightToken *Token) bool {
    // Évaluer la condition de jointure
    leftValue := resolveValue(leftToken, jn.leftBinding)
    rightValue := resolveValue(rightToken, jn.rightBinding)
    
    return leftValue == rightValue
}
```

**Après** :
```go
type JoinNode struct {
    // ... champs existants
    resolver  *FieldResolver
    evaluator *ComparisonEvaluator
}

func (jn *JoinNode) evaluate(leftToken, rightToken *Token) (bool, error) {
    // Résoudre les valeurs avec types
    leftValue, leftType, err := jn.resolveBindingValue(leftToken, jn.leftBinding)
    if err != nil {
        return false, err
    }
    
    rightValue, rightType, err := jn.resolveBindingValue(rightToken, jn.rightBinding)
    if err != nil {
        return false, err
    }
    
    // Évaluer avec le nouveau comparateur
    return jn.evaluator.EvaluateComparison(leftValue, rightValue, jn.operator, leftType, rightType)
}

func (jn *JoinNode) resolveBindingValue(token *Token, binding BindingExpression) (interface{}, string, error) {
    // Si c'est un field access (ex: l.user)
    if binding.Type == "fieldAccess" {
        fact := token.GetFact(binding.Variable)
        if fact == nil {
            return nil, "", fmt.Errorf("fait pour variable '%s' non trouvé", binding.Variable)
        }
        
        return jn.resolver.ResolveFieldValue(fact, binding.Field)
    }
    
    // Si c'est une variable directe (ex: u)
    if binding.Type == "variable" {
        fact := token.GetFact(binding.Variable)
        if fact == nil {
            return nil, "", fmt.Errorf("fait pour variable '%s' non trouvé", binding.Variable)
        }
        
        // Retourner l'ID du fait
        return jn.resolver.ResolveFactID(fact), "fact", nil
    }
    
    // Autres cas (literals, etc.)
    return binding.Value, "primitive", nil
}
```

### 6. Tests d'Intégration RETE

#### Fichier : `rete/integration_fact_comparison_test.go` (nouveau)

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
    "testing"
    "github.com/resinsec/tsd/constraint"
)

func TestRETENetwork_FactComparison(t *testing.T) {
    t.Log("🧪 TEST RETE NETWORK - COMPARAISON DE FAITS")
    t.Log("============================================")
    
    // Définir le programme TSD
    input := `
        type User(#name: string, age: number)
        type Login(user: User, #email: string, password: string)
        
        alice = User("Alice", 30)
        bob = User("Bob", 25)
        
        Login(alice, "alice@example.com", "pw1")
        Login(bob, "bob@example.com", "pw2")
        
        {u: User, l: Login} / l.user == u ==> 
            Log("Login for " + u.name)
    `
    
    // Parser le programme
    program, err := constraint.ParseProgram(input)
    if err != nil {
        t.Fatalf("❌ Erreur de parsing: %v", err)
    }
    
    // Convertir en format RETE
    reteFacts, err := constraint.ConvertFactsToReteFormat(program)
    if err != nil {
        t.Fatalf("❌ Erreur de conversion: %v", err)
    }
    
    // Créer le réseau RETE
    network := NewNetwork()
    
    // Compiler les règles
    for _, expr := range program.Expressions {
        err := network.CompileExpression(expr, program.Types)
        if err != nil {
            t.Fatalf("❌ Erreur de compilation: %v", err)
        }
    }
    
    // Asserter les faits
    var activations int
    for _, fact := range reteFacts {
        activations += network.AssertFact(fact)
    }
    
    // Vérifier que les règles ont été activées
    // On attend 2 activations (une pour alice, une pour bob)
    if activations != 2 {
        t.Errorf("❌ Attendu 2 activations, reçu %d", activations)
    } else {
        t.Logf("✅ %d activations détectées", activations)
    }
    
    // Vérifier que les comparaisons ont matché correctement
    // Les détails dépendent de l'implémentation du réseau RETE
}

func TestRETENetwork_FactComparisonNoMatch(t *testing.T) {
    t.Log("🧪 TEST RETE NETWORK - PAS DE MATCH")
    t.Log("====================================")
    
    input := `
        type User(#name: string, age: number)
        type Login(user: User, #email: string)
        
        alice = User("Alice", 30)
        bob = User("Bob", 25)
        
        // Login référence bob, mais on cherche alice
        Login(bob, "someone@example.com")
        
        {u: User, l: Login} / u.name == "Alice" && l.user == u ==> 
            Log("Match")
    `
    
    program, err := constraint.ParseProgram(input)
    if err != nil {
        t.Fatalf("❌ Erreur de parsing: %v", err)
    }
    
    reteFacts, err := constraint.ConvertFactsToReteFormat(program)
    if err != nil {
        t.Fatalf("❌ Erreur de conversion: %v", err)
    }
    
    network := NewNetwork()
    
    for _, expr := range program.Expressions {
        err := network.CompileExpression(expr, program.Types)
        if err != nil {
            t.Fatalf("❌ Erreur de compilation: %v", err)
        }
    }
    
    var activations int
    for _, fact := range reteFacts {
        activations += network.AssertFact(fact)
    }
    
    // Aucune activation attendue (alice existe mais pas de login pour alice)
    if activations != 0 {
        t.Errorf("❌ Attendu 0 activations, reçu %d", activations)
    } else {
        t.Logf("✅ Aucune activation (comportement correct)")
    }
}

func TestRETENetwork_MultipleFactTypes(t *testing.T) {
    t.Log("🧪 TEST RETE NETWORK - TYPES MULTIPLES")
    t.Log("========================================")
    
    input := `
        type User(#name: string, age: number)
        type Order(user: User, #orderNum: number, total: number)
        type Payment(order: Order, #paymentId: string, amount: number)
        
        alice = User("Alice", 30)
        order1 = Order(alice, 1001, 150.50)
        Payment(order1, "PAY-001", 150.50)
        
        {u: User, o: Order, p: Payment} / 
            o.user == u && p.order == o ==> 
            Log("Payment " + p.paymentId + " for user " + u.name)
    `
    
    program, err := constraint.ParseProgram(input)
    if err != nil {
        t.Fatalf("❌ Erreur de parsing: %v", err)
    }
    
    reteFacts, err := constraint.ConvertFactsToReteFormat(program)
    if err != nil {
        t.Fatalf("❌ Erreur de conversion: %v", err)
    }
    
    network := NewNetwork()
    
    for _, expr := range program.Expressions {
        err := network.CompileExpression(expr, program.Types)
        if err != nil {
            t.Fatalf("❌ Erreur de compilation: %v", err)
        }
    }
    
    var activations int
    for _, fact := range reteFacts {
        activations += network.AssertFact(fact)
    }
    
    // Une activation attendue (chaîne User -> Order -> Payment)
    if activations != 1 {
        t.Errorf("❌ Attendu 1 activation, reçu %d", activations)
    } else {
        t.Logf("✅ Chaîne de faits correctement matchée")
    }
}
```

### 7. Validation de Types dans les Comparaisons

#### Fichier : `constraint/constraint_type_checking.go` (ou nouveau)

**Ajouter validation** :

```go
// ValidateFactComparison valide une comparaison impliquant des faits
// Vérifie que les types sont compatibles
func ValidateFactComparison(leftExpr, rightExpr interface{}, operator string, typeMap map[string]constraint.TypeDefinition, varTypes map[string]string) error {
    leftType, err := inferExpressionType(leftExpr, typeMap, varTypes)
    if err != nil {
        return fmt.Errorf("expression gauche: %v", err)
    }
    
    rightType, err := inferExpressionType(rightExpr, typeMap, varTypes)
    if err != nil {
        return fmt.Errorf("expression droite: %v", err)
    }
    
    // Vérifier la compatibilité des types
    if !areTypesCompatible(leftType, rightType, operator) {
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
func inferExpressionType(expr interface{}, typeMap map[string]constraint.TypeDefinition, varTypes map[string]string) (string, error) {
    switch e := expr.(type) {
    case constraint.FieldAccess:
        // Récupérer le type de la variable
        varType, exists := varTypes[e.Object]
        if !exists {
            return "", fmt.Errorf("variable '%s' non définie", e.Object)
        }
        
        // Récupérer le type du champ
        typeDef, exists := typeMap[varType]
        if !exists {
            return "", fmt.Errorf("type '%s' non trouvé", varType)
        }
        
        for _, field := range typeDef.Fields {
            if field.Name == e.Field {
                return field.Type, nil
            }
        }
        
        return "", fmt.Errorf("champ '%s' non trouvé dans type '%s'", e.Field, varType)
        
    case constraint.Variable:
        // Le type de la variable
        varType, exists := varTypes[e.Name]
        if !exists {
            return "", fmt.Errorf("variable '%s' non définie", e.Name)
        }
        return varType, nil
        
    case constraint.StringLiteral:
        return "string", nil
        
    case constraint.NumberLiteral:
        return "number", nil
        
    case constraint.BooleanLiteral:
        return "bool", nil
        
    default:
        return "", fmt.Errorf("type d'expression non supporté: %T", expr)
    }
}

// areTypesCompatible vérifie si deux types sont compatibles pour une comparaison
func areTypesCompatible(leftType, rightType, operator string) bool {
    primitives := map[string]bool{
        "string": true,
        "number": true,
        "bool":   true,
        "boolean": true,
    }
    
    // Même type primitif
    if leftType == rightType && primitives[leftType] {
        return true
    }
    
    // bool et boolean sont compatibles
    if (leftType == "bool" || leftType == "boolean") && (rightType == "bool" || rightType == "boolean") {
        return true
    }
    
    // Même type de fait
    if leftType == rightType && !primitives[leftType] {
        // Pour les faits, seuls == et != sont autorisés
        return operator == "==" || operator == "!="
    }
    
    return false
}
```

#### Tests de Validation

```go
func TestValidateFactComparison(t *testing.T) {
    t.Log("🧪 TEST VALIDATION COMPARAISON DE FAITS")
    t.Log("========================================")
    
    typeMap := map[string]constraint.TypeDefinition{
        "User": {
            Name: "User",
            Fields: []constraint.Field{
                {Name: "name", Type: "string", IsPrimaryKey: true},
                {Name: "age", Type: "number"},
            },
        },
        "Login": {
            Name: "Login",
            Fields: []constraint.Field{
                {Name: "user", Type: "User"},
                {Name: "email", Type: "string", IsPrimaryKey: true},
            },
        },
    }
    
    varTypes := map[string]string{
        "u": "User",
        "l": "Login",
    }
    
    tests := []struct {
        name      string
        leftExpr  interface{}
        rightExpr interface{}
        operator  string
        wantErr   bool
    }{
        {
            name: "comparaison fait valide",
            leftExpr: constraint.FieldAccess{
                Type:   "fieldAccess",
                Object: "l",
                Field:  "user",
            },
            rightExpr: constraint.Variable{
                Type: "variable",
                Name: "u",
            },
            operator: "==",
            wantErr:  false,
        },
        {
            name: "comparaison primitif valide",
            leftExpr: constraint.FieldAccess{
                Type:   "fieldAccess",
                Object: "u",
                Field:  "name",
            },
            rightExpr: constraint.StringLiteral{
                Type:  "string",
                Value: "Alice",
            },
            operator: "==",
            wantErr:  false,
        },
        {
            name: "comparaison types incompatibles",
            leftExpr: constraint.FieldAccess{
                Type:   "fieldAccess",
                Object: "u",
                Field:  "name",
            },
            rightExpr: constraint.NumberLiteral{
                Type:  "number",
                Value: 42.0,
            },
            operator: "==",
            wantErr:  true,
        },
        {
            name: "opérateur < interdit pour faits",
            leftExpr: constraint.FieldAccess{
                Type:   "fieldAccess",
                Object: "l",
                Field:  "user",
            },
            rightExpr: constraint.Variable{
                Type: "variable",
                Name: "u",
            },
            operator: "<",
            wantErr:  true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateFactComparison(tt.leftExpr, tt.rightExpr, tt.operator, typeMap, varTypes)
            
            if tt.wantErr {
                if err == nil {
                    t.Errorf("❌ Attendu une erreur, reçu nil")
                } else {
                    t.Logf("✅ Erreur attendue: %v", err)
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

---

## ✅ Critères de Succès

### Compilation et Tests

```bash
# Code compile
go build ./rete
go build ./constraint

# Tests unitaires passent
go test ./rete -run TestFieldResolver -v
go test ./rete -run TestComparison -v
go test ./rete -run TestRETENetwork -v

# Tests d'intégration passent
go test ./rete -run TestIntegration -v

# Couverture > 80%
go test ./rete -cover
go test ./constraint -cover
```

### Fonctionnalités

- [ ] `FieldResolver` créé et fonctionnel
- [ ] `ComparisonEvaluator` créé et fonctionnel
- [ ] Comparaisons de faits via IDs fonctionnelles
- [ ] Comparaisons de primitifs fonctionnelles
- [ ] Validation de types implémentée
- [ ] Intégration dans nœuds RETE complète
- [ ] Tests d'intégration passent

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

- [ ] `TestFieldResolver_ResolveFieldValue`
- [ ] `TestFieldResolver_ResolveFactID`
- [ ] `TestComparisonEvaluator_CompareFactIDs`
- [ ] `TestComparisonEvaluator_ComparePrimitives`
- [ ] `TestComparisonEvaluator_EvaluateComparison`
- [ ] `TestValidateFactComparison`

### Tests d'Intégration Minimaux

- [ ] `TestRETENetwork_FactComparison`
- [ ] `TestRETENetwork_FactComparisonNoMatch`
- [ ] `TestRETENetwork_MultipleFactTypes`

### Tests End-to-End

Créer un test complet avec un programme TSD réel.

---

## 🚀 Exécution

### Ordre des Modifications

1. ✅ Analyser évaluation actuelle
2. ✅ Créer `FieldResolver`
3. ✅ Créer `ComparisonEvaluator`
4. ✅ Modifier nœuds RETE
5. ✅ Ajouter validation de types
6. ✅ Tests unitaires
7. ✅ Tests d'intégration
8. ✅ Validation finale

### Commandes

```bash
# Créer les nouveaux fichiers
touch rete/field_resolver.go
touch rete/field_resolver_test.go
touch rete/comparison_evaluator.go
touch rete/comparison_evaluator_test.go
touch rete/integration_fact_comparison_test.go

# Tester au fur et à mesure
go test ./rete -run TestFieldResolver -v
go test ./rete -run TestComparison -v

# Validation complète
make validate
make test-complete
```

---

## 📚 Références

- `scripts/new_ids/03-prompt-id-generation.md` - Génération IDs
- `scripts/new_ids/02-prompt-parser-syntax.md` - Syntaxe parser
- `rete/` - Package RETE actuel
- `docs/architecture/rete.md` - Architecture RETE

---

## 📝 Notes

### Points d'Attention

1. **Performance** : La résolution de types doit être rapide (utilisée à chaque évaluation)

2. **Cache** : Considérer un cache de résolution si performances dégradées

3. **Erreurs claires** : Messages d'erreur explicites pour debugging

4. **Compatibilité** : S'assurer que les comparaisons de primitifs continuent de fonctionner

### Questions Résolues

Q: Faut-il un cache de résolution ?
R: Pas pour l'instant, optimiser plus tard si nécessaire

Q: Comment gérer les types récursifs ?
R: Validation empêche les cycles, donc pas de problème

---

## 🎯 Résultat Attendu

```tsd
type User(#name: string, age: number)
type Login(user: User, #email: string)

alice = User("Alice", 30)
Login(alice, "alice@ex.com")

// ✅ Cette syntaxe fonctionne
{u: User, l: Login} / l.user == u ==> Log("Match")

// ✅ Comparaisons mixtes
{u: User, l: Login} / l.user == u && u.age > 25 ==> Log("Senior user")

// ❌ Types incompatibles détectés
{u: User, l: Login} / l.user < u ==> ... // Erreur: < non supporté pour faits
```

---

**Prompt suivant** : `05-prompt-types-validation.md`

**Durée estimée** : 6-8 heures

**Complexité** : 🔴 Élevée (modification cœur RETE)