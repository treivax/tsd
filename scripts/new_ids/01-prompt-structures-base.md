# Prompt 01 - Modification des Structures de Base

> **📋 Standards** : Ce prompt respecte les règles de [.github/prompts/common.md](../../.github/prompts/common.md) et [.github/prompts/develop.md](../../.github/prompts/develop.md)

---

## 🎯 Objectif

Modifier les structures de base du projet pour introduire le champ interne `_id_` en remplacement du champ `id` visible.

Cette modification est **fondamentale** et impacte toutes les couches du système.

---

## 📋 Contexte

### État Actuel

```go
// constraint/constraint_constants.go
const FieldNameID = "id"

// constraint/constraint_types.go
type Fact struct {
    Type     string      `json:"type"`
    TypeName string      `json:"typeName"`
    Fields   []FactField `json:"fields"`
}

// rete/fact_token.go
type Fact struct {
    ID     string                 `json:"id"`
    Type   string                 `json:"type"`
    Fields map[string]interface{} `json:"fields"`
}

// tsdio/api.go
type Fact struct {
    ID     string                 `json:"id"`
    Type   string                 `json:"type"`
    Fields map[string]interface{} `json:"fields"`
}
```

### État Cible

```go
// constraint/constraint_constants.go
const FieldNameInternalID = "_id_"

// Les structures Fact restent similaires mais :
// - Le champ _id_ est caché (jamais dans expressions TSD)
// - Stocké en interne dans Fields map
// - Constante utilisée partout pour cohérence
```

---

## 📝 Tâches à Réaliser

### 1. Modifier les Constantes

#### Fichier : `constraint/constraint_constants.go`

**Actions** :

1. **Renommer la constante** :
   ```go
   // Avant
   const FieldNameID = "id"
   
   // Après
   const FieldNameInternalID = "_id_"
   ```

2. **Ajouter une constante pour l'ancien nom (temporaire pour migration)** :
   ```go
   // Deprecated: Use FieldNameInternalID instead
   // Kept temporarily for migration purposes
   const FieldNameIDLegacy = "id"
   ```

3. **Documenter avec GoDoc** :
   ```go
   // FieldNameInternalID is the internal identifier field name.
   // This field is automatically generated and NEVER accessible in TSD expressions.
   // It is hidden from users and used only internally by the RETE engine.
   const FieldNameInternalID = "_id_"
   ```

#### Fichier : `constraint/constraint_constants_test.go`

**Actions** :

1. **Mettre à jour les tests de constantes** :
   ```go
   func TestConstantValues(t *testing.T) {
       tests := []struct {
           name     string
           constant string
           expected string
       }{
           {"FieldNameInternalID", FieldNameInternalID, "_id_"},
           {"FieldNameIDLegacy", FieldNameIDLegacy, "id"},
           // ...
       }
       // ...
   }
   ```

2. **Ajouter tests de validation** :
   ```go
   func TestInternalIDNotAccessible(t *testing.T) {
       // Test que _id_ ne peut pas être utilisé comme nom de champ
       // dans une définition de type ou de fait
   }
   ```

### 2. Mettre à Jour les Références

#### Rechercher et Remplacer

**Commande** :
```bash
# Trouver toutes les utilisations de FieldNameID
grep -r "FieldNameID" --include="*.go" constraint/ rete/ api/ tsdio/

# Remplacer par FieldNameInternalID
# ATTENTION : Faire case par case, pas automatiquement
```

**Fichiers attendus** :
- `constraint/constraint_facts.go`
- `constraint/constraint_field_validation.go`
- `constraint/action_validator.go`
- `constraint/primary_key_validation.go`
- `constraint/id_generator.go`
- `rete/` (plusieurs fichiers)
- `api/` (plusieurs fichiers)
- `tsdio/` (plusieurs fichiers)

**Pour chaque occurrence** :

1. Vérifier le contexte
2. Remplacer `FieldNameID` → `FieldNameInternalID`
3. Vérifier que la logique reste correcte
4. Mettre à jour les commentaires si nécessaire

### 3. Modifier la Validation des Faits

#### Fichier : `constraint/primary_key_validation.go`

**Fonction actuelle** :
```go
func ValidateFactPrimaryKey(fact Fact, typeDef TypeDefinition) error {
    // Validation actuelle permet 'id' comme clé primaire
    // Interdit 'id' comme champ manuel sauf si clé primaire
}
```

**Nouvelle logique** :

```go
// ValidateFactPrimaryKey validates that:
// 1. Primary key fields are present in the fact
// 2. The internal ID field (_id_) is NEVER manually defined
// 3. All primary key values are valid
func ValidateFactPrimaryKey(fact Fact, typeDef TypeDefinition) error {
    // 1. Vérifier que _id_ n'est JAMAIS défini manuellement
    for _, factField := range fact.Fields {
        if factField.Name == FieldNameInternalID {
            return fmt.Errorf(
                "fait de type '%s': le champ '%s' est réservé au système et ne peut pas être défini manuellement",
                fact.TypeName,
                FieldNameInternalID,
            )
        }
    }
    
    // 2. Valider les clés primaires (logique existante)
    // ...
    
    return nil
}
```

**Tests associés** :

```go
func TestValidateFactPrimaryKey_InternalIDForbidden(t *testing.T) {
    tests := []struct {
        name    string
        fact    Fact
        typeDef TypeDefinition
        wantErr bool
        errMsg  string
    }{
        {
            name: "interdit _id_ manuel",
            fact: Fact{
                TypeName: "User",
                Fields: []FactField{
                    {Name: FieldNameInternalID, Value: FactValue{Type: "string", Value: "manual"}},
                    {Name: "name", Value: FactValue{Type: "string", Value: "Alice"}},
                },
            },
            typeDef: TypeDefinition{
                Name: "User",
                Fields: []Field{
                    {Name: "name", Type: "string", IsPrimaryKey: true},
                },
            },
            wantErr: true,
            errMsg:  "champ '_id_' est réservé",
        },
        // Plus de tests...
    }
}
```

### 4. Modifier la Validation des Types

#### Fichier : `constraint/constraint_type_validation.go`

**Ajouter validation** :

```go
// ValidateTypeDefinition validates a type definition.
// It ensures field names are valid and don't use reserved names.
func ValidateTypeDefinition(typeDef TypeDefinition) error {
    // Validation existante...
    
    // Nouvelle validation : interdire _id_ comme nom de champ
    for _, field := range typeDef.Fields {
        if field.Name == FieldNameInternalID {
            return fmt.Errorf(
                "type '%s': le champ '%s' est réservé au système et ne peut pas être utilisé",
                typeDef.Name,
                FieldNameInternalID,
            )
        }
    }
    
    return nil
}
```

**Tests** :

```go
func TestValidateTypeDefinition_InternalIDForbidden(t *testing.T) {
    typeDef := TypeDefinition{
        Name: "User",
        Fields: []Field{
            {Name: FieldNameInternalID, Type: "string"},
            {Name: "name", Type: "string"},
        },
    }
    
    err := ValidateTypeDefinition(typeDef)
    if err == nil {
        t.Fatal("attendu une erreur pour champ _id_")
    }
    
    if !strings.Contains(err.Error(), "réservé") {
        t.Errorf("message d'erreur inattendu: %v", err)
    }
}
```

### 5. Modifier la Validation d'Accès aux Champs

#### Fichier : `constraint/constraint_field_validation.go`

**Fonction actuelle** :
```go
func GetFieldType(variable string, field string, /* ... */) (string, error) {
    // Le champ 'id' est spécial, toujours de type string
    if field == FieldNameID {
        return "string", nil
    }
    // ...
}
```

**Nouvelle logique** :

```go
func GetFieldType(variable string, field string, /* ... */) (string, error) {
    // Le champ '_id_' est INTERDIT dans les expressions
    if field == FieldNameInternalID {
        return "", fmt.Errorf(
            "le champ '%s' est interne et ne peut pas être accédé dans les expressions",
            FieldNameInternalID,
        )
    }
    
    // Note: Nous ajouterons plus tard le support pour les champs de type Fait
    // qui permettront les comparaisons comme p.user == u
    
    // Logique existante pour autres champs...
}
```

**Tests** :

```go
func TestGetFieldType_InternalIDForbidden(t *testing.T) {
    ps := NewProgramState()
    ps.AddType(TypeDefinition{
        Name: "User",
        Fields: []Field{
            {Name: "name", Type: "string", IsPrimaryKey: true},
        },
    })
    
    // Tenter d'accéder à _id_
    _, err := GetFieldType("u", FieldNameInternalID, ps.Types, map[string]string{"u": "User"})
    
    if err == nil {
        t.Fatal("attendu une erreur pour accès à _id_")
    }
    
    if !strings.Contains(err.Error(), "interne") {
        t.Errorf("message d'erreur inattendu: %v", err)
    }
}
```

### 6. Mettre à Jour la Génération d'IDs

#### Fichier : `constraint/constraint_facts.go`

**Fonction actuelle** :
```go
func ensureFactID(reteFact map[string]interface{}, fact Fact, typeDef TypeDefinition) (string, error) {
    // Vérifier si ID explicite
    if id, exists := reteFact[FieldNameID]; exists {
        // Backward compatibility
    }
    
    // Générer ID
    id, err := GenerateFactID(fact, typeDef)
    // ...
}
```

**Nouvelle logique** :

```go
// ensureFactID generates an internal ID for a fact.
// The ID is ALWAYS generated, never provided manually.
func ensureFactID(reteFact map[string]interface{}, fact Fact, typeDef TypeDefinition) (string, error) {
    // Vérifier que _id_ n'a PAS été fourni manuellement
    if _, exists := reteFact[FieldNameInternalID]; exists {
        return "", fmt.Errorf(
            "le champ '%s' ne peut pas être défini manuellement pour le type '%s'",
            FieldNameInternalID,
            fact.TypeName,
        )
    }
    
    // TOUJOURS générer l'ID
    id, err := GenerateFactID(fact, typeDef)
    if err != nil {
        return "", fmt.Errorf("génération d'ID pour le fait de type '%s': %v", fact.TypeName, err)
    }
    
    return id, nil
}
```

**Mise à jour de l'utilisation** :

```go
func ConvertFactsToReteFormat(program Program) ([]map[string]interface{}, error) {
    // ...
    for i, fact := range program.Facts {
        // ...
        factID, err := ensureFactID(reteFact, fact, typeDef)
        if err != nil {
            return nil, fmt.Errorf("fait %d: %v", i+1, err)
        }
        
        // Stocker avec le nouveau nom
        reteFact[FieldNameInternalID] = factID
        reteFact[FieldNameReteType] = fact.TypeName
        // ...
    }
}
```

### 7. Mettre à Jour les Structures RETE

#### Fichier : `rete/fact_token.go`

**Structure actuelle** :
```go
type Fact struct {
    ID         string                 `json:"id"`
    Type       string                 `json:"type"`
    Fields     map[string]interface{} `json:"fields"`
    Attributes map[string]interface{} `json:"attributes,omitempty"`
}
```

**Modifications** :

1. **Garder la structure** (le champ ID est nécessaire en interne)
2. **Changer la sérialisation JSON** pour cacher `_id_` :

```go
type Fact struct {
    // ID est l'identifiant interne du fait.
    // Il est généré automatiquement et JAMAIS accessible dans les expressions TSD.
    // En JSON, il est sérialisé comme "_id_" et caché de l'API publique.
    ID         string                 `json:"_id_"`
    Type       string                 `json:"type"`
    Fields     map[string]interface{} `json:"fields"`
    Attributes map[string]interface{} `json:"attributes,omitempty"` // Alias pour Fields
}
```

3. **Ajouter une méthode pour accès interne** :

```go
// InternalID returns the internal identifier of the fact.
// This should only be used internally by the RETE engine.
func (f *Fact) InternalID() string {
    return f.ID
}
```

4. **Mettre à jour NewFact** :

```go
func NewFact(id, factType string, fields map[string]interface{}) *Fact {
    return &Fact{
        ID:     id,  // ID interne
        Type:   factType,
        Fields: fields,
    }
}
```

### 8. Mettre à Jour les Structures API

#### Fichier : `tsdio/api.go`

**Structure actuelle** :
```go
type Fact struct {
    ID     string                 `json:"id"`
    Type   string                 `json:"type"`
    Fields map[string]interface{} `json:"fields"`
}
```

**Modifications** :

1. **Changer tag JSON** :
```go
type Fact struct {
    // ID est l'identifiant interne, caché de l'utilisateur
    ID     string                 `json:"_id_"`
    Type   string                 `json:"type"`
    Fields map[string]interface{} `json:"fields"`
}
```

2. **Ajouter validation** :
```go
// ValidateFact validates that a fact doesn't contain reserved fields
func ValidateFact(fact Fact) error {
    // Le champ _id_ ne doit jamais être dans Fields
    if _, exists := fact.Fields[constraint.FieldNameInternalID]; exists {
        return fmt.Errorf("le champ '%s' est réservé", constraint.FieldNameInternalID)
    }
    return nil
}
```

#### Fichier : `api/result.go`

**Actions similaires** :

1. Vérifier les structures de résultats
2. S'assurer que `_id_` n'est pas exposé publiquement
3. Mettre à jour les méthodes de sérialisation

### 9. Tests Complets

#### Fichier : `constraint/internal_id_test.go` (nouveau)

**Créer un fichier de tests dédié** :

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
    "strings"
    "testing"
)

func TestInternalID_ReservedName(t *testing.T) {
    t.Log("🧪 TEST INTERNAL ID - NOM RÉSERVÉ")
    t.Log("===================================")
    
    tests := []struct {
        name    string
        testFn  func() error
        wantErr bool
    }{
        {
            name: "interdire _id_ dans définition de type",
            testFn: func() error {
                return ValidateTypeDefinition(TypeDefinition{
                    Name: "User",
                    Fields: []Field{
                        {Name: FieldNameInternalID, Type: "string"},
                    },
                })
            },
            wantErr: true,
        },
        {
            name: "interdire _id_ dans définition de fait",
            testFn: func() error {
                return ValidateFactPrimaryKey(
                    Fact{
                        TypeName: "User",
                        Fields: []FactField{
                            {Name: FieldNameInternalID, Value: FactValue{Type: "string", Value: "test"}},
                        },
                    },
                    TypeDefinition{
                        Name: "User",
                        Fields: []Field{
                            {Name: "name", Type: "string", IsPrimaryKey: true},
                        },
                    },
                )
            },
            wantErr: true,
        },
        // Plus de tests...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.testFn()
            
            if tt.wantErr {
                if err == nil {
                    t.Errorf("❌ Attendu une erreur, reçu nil")
                } else {
                    t.Logf("✅ Erreur attendue: %v", err)
                }
            } else {
                if err != nil {
                    t.Errorf("❌ Erreur inattendue: %v", err)
                }
            }
        })
    }
}

func TestInternalID_AlwaysGenerated(t *testing.T) {
    t.Log("🧪 TEST INTERNAL ID - TOUJOURS GÉNÉRÉ")
    t.Log("======================================")
    
    // Test que l'ID est toujours généré même si non fourni
    fact := Fact{
        TypeName: "User",
        Fields: []FactField{
            {Name: "name", Value: FactValue{Type: "string", Value: "Alice"}},
        },
    }
    
    typeDef := TypeDefinition{
        Name: "User",
        Fields: []Field{
            {Name: "name", Type: "string", IsPrimaryKey: true},
        },
    }
    
    reteFact := createReteFact(fact, typeDef)
    id, err := ensureFactID(reteFact, fact, typeDef)
    
    if err != nil {
        t.Fatalf("❌ Erreur inattendue: %v", err)
    }
    
    if id == "" {
        t.Error("❌ ID généré est vide")
    }
    
    expectedID := "User~Alice"
    if id != expectedID {
        t.Errorf("❌ ID attendu '%s', reçu '%s'", expectedID, id)
    }
    
    t.Logf("✅ ID généré correctement: %s", id)
}

func TestInternalID_NeverManual(t *testing.T) {
    t.Log("🧪 TEST INTERNAL ID - JAMAIS MANUEL")
    t.Log("====================================")
    
    fact := Fact{
        TypeName: "User",
        Fields: []FactField{
            {Name: "name", Value: FactValue{Type: "string", Value: "Alice"}},
        },
    }
    
    typeDef := TypeDefinition{
        Name: "User",
        Fields: []Field{
            {Name: "name", Type: "string", IsPrimaryKey: true},
        },
    }
    
    reteFact := createReteFact(fact, typeDef)
    
    // Simuler une tentative de définition manuelle
    reteFact[FieldNameInternalID] = "manual_id"
    
    _, err := ensureFactID(reteFact, fact, typeDef)
    
    if err == nil {
        t.Error("❌ Attendu une erreur pour ID manuel")
    } else {
        t.Logf("✅ Erreur attendue pour ID manuel: %v", err)
    }
}
```

---

## ✅ Critères de Succès

### Compilation et Tests

```bash
# Doit compiler sans erreur
go build ./...

# Tous les tests passent
go test ./constraint/... -v
go test ./rete/... -v
go test ./tsdio/... -v
go test ./api/... -v

# Couverture > 80%
go test ./constraint/... -cover
```

### Validation

```bash
# Format
make format

# Lint
make lint

# Validation complète
make validate
```

### Vérifications Manuelles

- [ ] Constante `FieldNameInternalID = "_id_"` définie
- [ ] Toutes les références mises à jour
- [ ] Validation interdit `_id_` dans types
- [ ] Validation interdit `_id_` dans faits
- [ ] ID toujours généré automatiquement
- [ ] Tests passent (> 80% couverture)
- [ ] GoDoc à jour
- [ ] Pas de hardcoding

---

## 📊 Tests Requis

### Tests Unitaires Minimaux

- [ ] `TestConstantValues` - Vérifier valeur "_id_"
- [ ] `TestValidateTypeDefinition_InternalIDForbidden` - Interdire _id_ dans types
- [ ] `TestValidateFactPrimaryKey_InternalIDForbidden` - Interdire _id_ dans faits
- [ ] `TestGetFieldType_InternalIDForbidden` - Interdire accès à _id_
- [ ] `TestInternalID_AlwaysGenerated` - ID toujours généré
- [ ] `TestInternalID_NeverManual` - Jamais manuel
- [ ] `TestEnsureFactID_RejectsManualID` - Rejeter ID manuel

### Tests d'Intégration

- [ ] Créer un programme avec types et faits
- [ ] Vérifier que IDs sont générés
- [ ] Vérifier que `_id_` est caché
- [ ] Vérifier qu'erreur si `_id_` utilisé

---

## 🚀 Exécution

### Ordre des Modifications

1. ✅ Constantes (`constraint_constants.go`)
2. ✅ Tests constantes (`constraint_constants_test.go`)
3. ✅ Validation types (`constraint_type_validation.go`)
4. ✅ Validation faits (`primary_key_validation.go`)
5. ✅ Validation champs (`constraint_field_validation.go`)
6. ✅ Génération ID (`constraint_facts.go`)
7. ✅ Structure RETE (`rete/fact_token.go`)
8. ✅ API (`tsdio/api.go`, `api/result.go`)
9. ✅ Tests complets (`internal_id_test.go`)
10. ✅ Validation finale

### Commandes

```bash
# 1. Modifier les fichiers un par un
# 2. Tester après chaque modification
go test ./constraint -run TestConstant -v

# 3. Une fois tout modifié, tests complets
make test-unit

# 4. Validation
make validate
```

---

## 📚 Références

- `.github/prompts/common.md` - Standards
- `.github/prompts/develop.md` - Développement
- `scripts/new_ids/00-prompt-analyse.md` - Analyse préliminaire
- `scripts/new_ids/README.md` - Vue d'ensemble

---

## 📝 Notes

### Points d'Attention

1. **Parser généré** : Ne PAS modifier `constraint/parser.go`. Les changements de grammaire viendront dans le prochain prompt.

2. **Backward compatibility** : Cette modification CASSE la rétrocompatibilité. C'est assumé.

3. **Sérialisation JSON** : Le tag `json:"_id_"` cache le champ de l'API publique mais le garde en interne.

4. **Performance** : Aucun impact attendu sur les performances.

### Questions Résolues

Q: Faut-il garder un alias `id` en lecture seule ?
R: Non, `_id_` est complètement caché. Les prochains prompts ajouteront les comparaisons via types.

Q: Comment gérer la sérialisation ?
R: Tag JSON `_id_` pour le cacher dans l'API publique.

---

## 🎯 Résultat Attendu

Après ce prompt :

```go
// Constantes
const FieldNameInternalID = "_id_"

// Validation
ValidateTypeDefinition(User{_id_: "x"}) → ❌ Erreur
ValidateFact(User{_id_: "x"}) → ❌ Erreur

// Génération
fact := Fact{TypeName: "User", Fields: [...]}
id := ensureFactID(...) → "User~Alice" (toujours généré)

// Accès
GetFieldType("u", "_id_", ...) → ❌ Erreur "interne"
```

---

**Prompt suivant** : `02-prompt-parser-syntax.md`

**Durée estimée** : 4-6 heures

**Complexité** : 🔴 Élevée (modifications critiques)