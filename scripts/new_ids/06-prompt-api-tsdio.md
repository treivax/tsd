# Prompt 06 - Adaptation de l'API et tsdio

> **📋 Standards** : Ce prompt respecte les règles de [.github/prompts/common.md](../../.github/prompts/common.md) et [.github/prompts/develop.md](../../.github/prompts/develop.md)

---

## 🎯 Objectif

Adapter les modules `api/` et `tsdio/` pour :

1. **Cacher `_id_`** de l'API publique
2. **Supporter les nouvelles structures** - Affectations, types de faits
3. **Validation côté API** - Vérifier les entrées utilisateur
4. **Sérialisation JSON** - Format cohérent et sécurisé
5. **Documentation API** - Contrats clairs

---

## 📋 Contexte

### État Actuel

Les modules `api/` et `tsdio/` exposent :
- Structures de faits avec champ `ID` public
- API pour asserter des faits
- Récupération de résultats
- Pas de support pour affectations de variables

### État Cible

```go
// API publique ne doit jamais exposer _id_
type Fact struct {
    // ID interne caché de l'API JSON publique
    internalID string
    Type       string
    Fields     map[string]interface{}
}

// Support des affectations
type Program struct {
    Types           []TypeDef
    FactAssignments []FactAssignment  // NOUVEAU
    Facts           []Fact
    Rules           []Rule
}
```

---

## 📝 Tâches à Réaliser

### 1. Analyser l'API Actuelle

#### Fichiers à Examiner

**Recherche** :
```bash
# Structure du module api
ls -la api/

# Structure du module tsdio
ls -la tsdio/

# Trouver les structures exposées
grep -r "type.*struct" api/ tsdio/ --include="*.go"
```

**Fichiers attendus** :
- `api/result.go` - Résultats et récupération
- `api/engine.go` ou similaire - Moteur principal
- `tsdio/api.go` - API publique tsdio
- `tsdio/program.go` - Structures de programme

#### Questions à Répondre

1. Comment les faits sont-ils actuellement sérialisés en JSON ?
2. Le champ `ID` est-il exposé publiquement ?
3. Y a-t-il des méthodes pour créer des faits ?
4. Comment les programmes sont-ils construits ?
5. Y a-t-il une validation côté API ?

### 2. Modifier les Structures de Faits

#### Fichier : `tsdio/api.go`

**Structure actuelle** :
```go
type Fact struct {
    ID     string                 `json:"id"`
    Type   string                 `json:"type"`
    Fields map[string]interface{} `json:"fields"`
}
```

**Nouvelle structure** :
```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package tsdio

import (
    "fmt"
    "github.com/resinsec/tsd/constraint"
)

// Fact représente un fait dans l'API publique
// L'ID interne (_id_) n'est jamais exposé via JSON
type Fact struct {
    // internalID est l'identifiant interne du fait
    // Il n'est PAS sérialisé en JSON (pas de tag json)
    // Il est maintenu en interne pour la cohérence du système
    internalID string
    
    // Type est le type du fait (ex: "User", "Login")
    Type string `json:"type"`
    
    // Fields contient les champs du fait
    // Le champ _id_ ne doit JAMAIS être présent ici
    Fields map[string]interface{} `json:"fields"`
}

// NewFact crée un nouveau fait
// L'ID interne sera généré lors de l'insertion dans le système
func NewFact(factType string, fields map[string]interface{}) (*Fact, error) {
    // Valider que _id_ n'est pas dans les champs
    if _, exists := fields[constraint.FieldNameInternalID]; exists {
        return nil, fmt.Errorf(
            "le champ '%s' est réservé et ne peut pas être défini manuellement",
            constraint.FieldNameInternalID,
        )
    }
    
    return &Fact{
        Type:   factType,
        Fields: fields,
    }, nil
}

// GetInternalID retourne l'ID interne du fait
// Cette méthode est pour usage interne uniquement
func (f *Fact) GetInternalID() string {
    return f.internalID
}

// SetInternalID définit l'ID interne du fait
// Cette méthode est pour usage interne uniquement
func (f *Fact) SetInternalID(id string) {
    f.internalID = id
}

// Validate valide le fait
func (f *Fact) Validate() error {
    if f.Type == "" {
        return fmt.Errorf("le type du fait ne peut pas être vide")
    }
    
    if f.Fields == nil {
        return fmt.Errorf("les champs du fait ne peuvent pas être nil")
    }
    
    // Vérifier que _id_ n'est pas dans les champs
    if _, exists := f.Fields[constraint.FieldNameInternalID]; exists {
        return fmt.Errorf(
            "le champ '%s' est réservé et ne peut pas être défini",
            constraint.FieldNameInternalID,
        )
    }
    
    return nil
}

// MarshalJSON personnalise la sérialisation JSON
// S'assure que _id_ n'est jamais exposé
func (f *Fact) MarshalJSON() ([]byte, error) {
    // Créer une structure anonyme sans internalID
    type FactAlias Fact
    return json.Marshal(&struct {
        *FactAlias
    }{
        FactAlias: (*FactAlias)(f),
    })
}
```

#### Tests des Structures

**Fichier : `tsdio/api_test.go` (nouveau ou modifier)**

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package tsdio

import (
    "encoding/json"
    "strings"
    "testing"
    "github.com/resinsec/tsd/constraint"
)

func TestFact_NewFact(t *testing.T) {
    t.Log("🧪 TEST FACT - CRÉATION")
    t.Log("=======================")
    
    tests := []struct {
        name    string
        fType   string
        fields  map[string]interface{}
        wantErr bool
    }{
        {
            name:  "fait valide",
            fType: "User",
            fields: map[string]interface{}{
                "name": "Alice",
                "age":  30,
            },
            wantErr: false,
        },
        {
            name:  "_id_ interdit",
            fType: "User",
            fields: map[string]interface{}{
                constraint.FieldNameInternalID: "manual-id",
                "name": "Alice",
            },
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            fact, err := NewFact(tt.fType, tt.fields)
            
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
            
            if fact.Type != tt.fType {
                t.Errorf("❌ Type attendu '%s', reçu '%s'", tt.fType, fact.Type)
            }
            
            t.Logf("✅ Fait créé: %+v", fact)
        })
    }
}

func TestFact_Validate(t *testing.T) {
    t.Log("🧪 TEST FACT - VALIDATION")
    t.Log("==========================")
    
    tests := []struct {
        name    string
        fact    *Fact
        wantErr bool
    }{
        {
            name: "fait valide",
            fact: &Fact{
                Type: "User",
                Fields: map[string]interface{}{
                    "name": "Alice",
                },
            },
            wantErr: false,
        },
        {
            name: "type vide",
            fact: &Fact{
                Type:   "",
                Fields: map[string]interface{}{},
            },
            wantErr: true,
        },
        {
            name: "fields nil",
            fact: &Fact{
                Type:   "User",
                Fields: nil,
            },
            wantErr: true,
        },
        {
            name: "_id_ dans fields",
            fact: &Fact{
                Type: "User",
                Fields: map[string]interface{}{
                    constraint.FieldNameInternalID: "manual",
                },
            },
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.fact.Validate()
            
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

func TestFact_JSONSerialization(t *testing.T) {
    t.Log("🧪 TEST FACT - SÉRIALISATION JSON")
    t.Log("==================================")
    
    fact := &Fact{
        internalID: "User~Alice",
        Type:       "User",
        Fields: map[string]interface{}{
            "name": "Alice",
            "age":  30,
        },
    }
    
    // Sérialiser en JSON
    jsonData, err := json.Marshal(fact)
    if err != nil {
        t.Fatalf("❌ Erreur de sérialisation: %v", err)
    }
    
    jsonStr := string(jsonData)
    t.Logf("JSON généré: %s", jsonStr)
    
    // Vérifier que _id_ n'est PAS dans le JSON
    if strings.Contains(jsonStr, "_id_") {
        t.Errorf("❌ Le JSON contient '_id_' alors qu'il devrait être caché")
    }
    
    // Vérifier que internalID n'est PAS dans le JSON
    if strings.Contains(jsonStr, "internalID") {
        t.Errorf("❌ Le JSON contient 'internalID' alors qu'il devrait être caché")
    }
    
    // Vérifier que les champs attendus sont présents
    if !strings.Contains(jsonStr, "User") {
        t.Error("❌ Le type 'User' devrait être dans le JSON")
    }
    
    if !strings.Contains(jsonStr, "Alice") {
        t.Error("❌ Le champ 'name' devrait être dans le JSON")
    }
    
    t.Log("✅ Sérialisation JSON correcte, _id_ caché")
}

func TestFact_InternalIDMethods(t *testing.T) {
    t.Log("🧪 TEST FACT - MÉTHODES ID INTERNE")
    t.Log("===================================")
    
    fact := &Fact{
        Type: "User",
        Fields: map[string]interface{}{
            "name": "Alice",
        },
    }
    
    // ID initial vide
    if fact.GetInternalID() != "" {
        t.Errorf("❌ ID initial devrait être vide, reçu '%s'", fact.GetInternalID())
    }
    
    // Définir un ID
    testID := "User~Alice"
    fact.SetInternalID(testID)
    
    // Vérifier que l'ID est défini
    if fact.GetInternalID() != testID {
        t.Errorf("❌ ID attendu '%s', reçu '%s'", testID, fact.GetInternalID())
    }
    
    t.Logf("✅ Méthodes d'ID interne fonctionnent correctement")
}
```

### 3. Ajouter Support des Affectations

#### Fichier : `tsdio/program.go` (nouveau ou modifier)

**Nouvelle structure** :

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package tsdio

import (
    "fmt"
)

// FactAssignment représente une affectation de fait à une variable
type FactAssignment struct {
    Variable string `json:"variable"` // Nom de la variable
    Fact     *Fact  `json:"fact"`     // Le fait assigné
}

// NewFactAssignment crée une nouvelle affectation
func NewFactAssignment(variable string, fact *Fact) (*FactAssignment, error) {
    if variable == "" {
        return nil, fmt.Errorf("le nom de la variable ne peut pas être vide")
    }
    
    if fact == nil {
        return nil, fmt.Errorf("le fait ne peut pas être nil")
    }
    
    if err := fact.Validate(); err != nil {
        return nil, fmt.Errorf("fait invalide: %v", err)
    }
    
    return &FactAssignment{
        Variable: variable,
        Fact:     fact,
    }, nil
}

// Validate valide l'affectation
func (fa *FactAssignment) Validate() error {
    if fa.Variable == "" {
        return fmt.Errorf("le nom de la variable ne peut pas être vide")
    }
    
    if fa.Fact == nil {
        return fmt.Errorf("le fait ne peut pas être nil")
    }
    
    return fa.Fact.Validate()
}

// Program représente un programme TSD complet
type Program struct {
    Types           []TypeDefinition   `json:"types,omitempty"`
    Actions         []ActionDefinition `json:"actions,omitempty"`
    FactAssignments []FactAssignment   `json:"factAssignments,omitempty"` // NOUVEAU
    Facts           []*Fact            `json:"facts,omitempty"`
    Rules           []Rule             `json:"rules,omitempty"`
}

// NewProgram crée un nouveau programme vide
func NewProgram() *Program {
    return &Program{
        Types:           make([]TypeDefinition, 0),
        Actions:         make([]ActionDefinition, 0),
        FactAssignments: make([]FactAssignment, 0),
        Facts:           make([]*Fact, 0),
        Rules:           make([]Rule, 0),
    }
}

// AddFactAssignment ajoute une affectation au programme
func (p *Program) AddFactAssignment(assignment *FactAssignment) error {
    if err := assignment.Validate(); err != nil {
        return err
    }
    
    // Vérifier que la variable n'est pas déjà définie
    for _, existing := range p.FactAssignments {
        if existing.Variable == assignment.Variable {
            return fmt.Errorf(
                "la variable '%s' est déjà définie",
                assignment.Variable,
            )
        }
    }
    
    p.FactAssignments = append(p.FactAssignments, *assignment)
    return nil
}

// AddFact ajoute un fait au programme
func (p *Program) AddFact(fact *Fact) error {
    if err := fact.Validate(); err != nil {
        return err
    }
    
    p.Facts = append(p.Facts, fact)
    return nil
}

// Validate valide le programme complet
func (p *Program) Validate() error {
    // Valider les affectations
    for i, assignment := range p.FactAssignments {
        if err := assignment.Validate(); err != nil {
            return fmt.Errorf("affectation %d: %v", i+1, err)
        }
    }
    
    // Valider les faits
    for i, fact := range p.Facts {
        if err := fact.Validate(); err != nil {
            return fmt.Errorf("fait %d: %v", i+1, err)
        }
    }
    
    return nil
}
```

#### Tests du Programme

**Fichier : `tsdio/program_test.go`**

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package tsdio

import (
    "testing"
)

func TestFactAssignment_Creation(t *testing.T) {
    t.Log("🧪 TEST FACT ASSIGNMENT - CRÉATION")
    t.Log("===================================")
    
    fact, _ := NewFact("User", map[string]interface{}{
        "name": "Alice",
    })
    
    tests := []struct {
        name     string
        variable string
        fact     *Fact
        wantErr  bool
    }{
        {
            name:     "affectation valide",
            variable: "alice",
            fact:     fact,
            wantErr:  false,
        },
        {
            name:     "variable vide",
            variable: "",
            fact:     fact,
            wantErr:  true,
        },
        {
            name:     "fait nil",
            variable: "alice",
            fact:     nil,
            wantErr:  true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            assignment, err := NewFactAssignment(tt.variable, tt.fact)
            
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
            
            if assignment.Variable != tt.variable {
                t.Errorf("❌ Variable attendue '%s', reçu '%s'", tt.variable, assignment.Variable)
            }
            
            t.Logf("✅ Affectation créée: %s = %s", assignment.Variable, assignment.Fact.Type)
        })
    }
}

func TestProgram_AddFactAssignment(t *testing.T) {
    t.Log("🧪 TEST PROGRAM - AJOUT AFFECTATION")
    t.Log("====================================")
    
    program := NewProgram()
    
    fact1, _ := NewFact("User", map[string]interface{}{"name": "Alice"})
    assignment1, _ := NewFactAssignment("alice", fact1)
    
    // Ajouter première affectation
    err := program.AddFactAssignment(assignment1)
    if err != nil {
        t.Fatalf("❌ Erreur inattendue: %v", err)
    }
    
    if len(program.FactAssignments) != 1 {
        t.Errorf("❌ Attendu 1 affectation, reçu %d", len(program.FactAssignments))
    }
    
    t.Log("✅ Première affectation ajoutée")
    
    // Essayer d'ajouter une affectation avec la même variable
    fact2, _ := NewFact("User", map[string]interface{}{"name": "Bob"})
    assignment2, _ := NewFactAssignment("alice", fact2)
    
    err = program.AddFactAssignment(assignment2)
    if err == nil {
        t.Error("❌ Attendu une erreur pour variable dupliquée")
    } else {
        t.Logf("✅ Erreur attendue pour variable dupliquée: %v", err)
    }
}

func TestProgram_Validate(t *testing.T) {
    t.Log("🧪 TEST PROGRAM - VALIDATION")
    t.Log("=============================")
    
    program := NewProgram()
    
    // Ajouter une affectation valide
    fact1, _ := NewFact("User", map[string]interface{}{"name": "Alice"})
    assignment1, _ := NewFactAssignment("alice", fact1)
    program.FactAssignments = append(program.FactAssignments, *assignment1)
    
    // Ajouter un fait valide
    fact2, _ := NewFact("Login", map[string]interface{}{
        "email": "alice@example.com",
    })
    program.Facts = append(program.Facts, fact2)
    
    // Valider
    err := program.Validate()
    if err != nil {
        t.Errorf("❌ Erreur inattendue: %v", err)
    } else {
        t.Log("✅ Programme valide")
    }
}
```

### 4. Modifier l'API de Résultats

#### Fichier : `api/result.go` (modifications)

**Cacher `_id_` dans les résultats** :

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package api

import (
    "encoding/json"
    "github.com/resinsec/tsd/constraint"
    "github.com/resinsec/tsd/xuples"
)

// Result représente un résultat de requête
// L'ID interne des faits n'est jamais exposé
type Result struct {
    Facts     []ResultFact           `json:"facts"`
    Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// ResultFact représente un fait dans les résultats
// Structure optimisée pour l'API publique
type ResultFact struct {
    Type   string                 `json:"type"`
    Fields map[string]interface{} `json:"fields"`
}

// FromXuple convertit un Xuple en ResultFact
// Cache l'ID interne
func FromXuple(xuple *xuples.Xuple) ResultFact {
    // Copier les champs sans _id_
    fields := make(map[string]interface{})
    for key, value := range xuple.Fields {
        // Exclure _id_
        if key != constraint.FieldNameInternalID {
            fields[key] = value
        }
    }
    
    return ResultFact{
        Type:   xuple.Type,
        Fields: fields,
    }
}

// MarshalJSON personnalise la sérialisation
func (rf *ResultFact) MarshalJSON() ([]byte, error) {
    // S'assurer que _id_ n'est pas dans les champs
    if _, exists := rf.Fields[constraint.FieldNameInternalID]; exists {
        // Créer une copie sans _id_
        cleanFields := make(map[string]interface{})
        for k, v := range rf.Fields {
            if k != constraint.FieldNameInternalID {
                cleanFields[k] = v
            }
        }
        
        return json.Marshal(&struct {
            Type   string                 `json:"type"`
            Fields map[string]interface{} `json:"fields"`
        }{
            Type:   rf.Type,
            Fields: cleanFields,
        })
    }
    
    // Sérialisation normale
    type Alias ResultFact
    return json.Marshal(&struct {
        *Alias
    }{
        Alias: (*Alias)(rf),
    })
}
```

#### Tests des Résultats

**Fichier : `api/result_test.go` (modifications)**

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package api

import (
    "encoding/json"
    "strings"
    "testing"
    "github.com/resinsec/tsd/constraint"
    "github.com/resinsec/tsd/xuples"
)

func TestResultFact_FromXuple(t *testing.T) {
    t.Log("🧪 TEST RESULT FACT - CONVERSION XUPLE")
    t.Log("======================================")
    
    xuple := &xuples.Xuple{
        Type: "User",
        Fields: map[string]interface{}{
            constraint.FieldNameInternalID: "User~Alice",
            "name": "Alice",
            "age":  30,
        },
    }
    
    resultFact := FromXuple(xuple)
    
    // Vérifier le type
    if resultFact.Type != "User" {
        t.Errorf("❌ Type attendu 'User', reçu '%s'", resultFact.Type)
    }
    
    // Vérifier que _id_ n'est PAS dans les champs
    if _, exists := resultFact.Fields[constraint.FieldNameInternalID]; exists {
        t.Error("❌ Le champ '_id_' ne devrait pas être présent dans ResultFact")
    }
    
    // Vérifier que les autres champs sont présents
    if name, ok := resultFact.Fields["name"].(string); !ok || name != "Alice" {
        t.Errorf("❌ Champ 'name' invalide: %v", resultFact.Fields["name"])
    }
    
    if age, ok := resultFact.Fields["age"].(int); !ok || age != 30 {
        t.Errorf("❌ Champ 'age' invalide: %v", resultFact.Fields["age"])
    }
    
    t.Log("✅ Conversion Xuple réussie, _id_ caché")
}

func TestResultFact_JSONSerialization(t *testing.T) {
    t.Log("🧪 TEST RESULT FACT - SÉRIALISATION JSON")
    t.Log("=========================================")
    
    resultFact := ResultFact{
        Type: "User",
        Fields: map[string]interface{}{
            "name": "Alice",
            "age":  30,
        },
    }
    
    jsonData, err := json.Marshal(resultFact)
    if err != nil {
        t.Fatalf("❌ Erreur de sérialisation: %v", err)
    }
    
    jsonStr := string(jsonData)
    t.Logf("JSON: %s", jsonStr)
    
    // Vérifier que _id_ n'est pas dans le JSON
    if strings.Contains(jsonStr, "_id_") {
        t.Error("❌ Le JSON contient '_id_'")
    }
    
    if strings.Contains(jsonStr, "internalID") {
        t.Error("❌ Le JSON contient 'internalID'")
    }
    
    // Vérifier les champs attendus
    if !strings.Contains(jsonStr, "User") {
        t.Error("❌ Le type devrait être dans le JSON")
    }
    
    if !strings.Contains(jsonStr, "Alice") {
        t.Error("❌ Le champ 'name' devrait être dans le JSON")
    }
    
    t.Log("✅ Sérialisation JSON correcte")
}

func TestResultFact_JSONWithInternalID(t *testing.T) {
    t.Log("🧪 TEST RESULT FACT - JSON AVEC _id_ INTERNE")
    t.Log("=============================================")
    
    // Simuler un cas où _id_ est accidentellement présent
    resultFact := ResultFact{
        Type: "User",
        Fields: map[string]interface{}{
            constraint.FieldNameInternalID: "User~Alice",
            "name": "Alice",
        },
    }
    
    jsonData, err := json.Marshal(resultFact)
    if err != nil {
        t.Fatalf("❌ Erreur de sérialisation: %v", err)
    }
    
    jsonStr := string(jsonData)
    
    // Vérifier que _id_ est filtré dans le JSON
    if strings.Contains(jsonStr, "_id_") {
        t.Error("❌ Le JSON contient '_id_' alors qu'il devrait être filtré")
    }
    
    // Vérifier que les autres champs sont présents
    if !strings.Contains(jsonStr, "Alice") {
        t.Error("❌ Le champ 'name' devrait être dans le JSON")
    }
    
    t.Log("✅ _id_ filtré correctement lors de la sérialisation")
}
```

### 5. Adapter l'API Engine

#### Fichier : `api/engine.go` (ou fichier principal du moteur)

**Ajouter support des affectations** :

```go
// AssertFactAssignment asserte une affectation de fait
func (e *Engine) AssertFactAssignment(assignment *tsdio.FactAssignment) error {
    if err := assignment.Validate(); err != nil {
        return fmt.Errorf("affectation invalide: %v", err)
    }
    
    // Convertir le fait en format interne
    internalFact, err := e.convertToInternalFact(assignment.Fact)
    if err != nil {
        return fmt.Errorf("conversion du fait: %v", err)
    }
    
    // Générer l'ID interne
    // (utilise le système de génération d'ID du prompt 03)
    factID, err := e.generateFactID(internalFact)
    if err != nil {
        return fmt.Errorf("génération d'ID: %v", err)
    }
    
    // Stocker l'ID dans le fait original
    assignment.Fact.SetInternalID(factID)
    
    // Enregistrer la variable dans le contexte
    e.variableContext.Register(assignment.Variable, factID)
    
    // Asserter le fait dans le moteur RETE
    return e.assertInternalFact(internalFact)
}

// AssertFact asserte un fait simple
func (e *Engine) AssertFact(fact *tsdio.Fact) error {
    if err := fact.Validate(); err != nil {
        return fmt.Errorf("fait invalide: %v", err)
    }
    
    // Convertir et asserter
    internalFact, err := e.convertToInternalFact(fact)
    if err != nil {
        return err
    }
    
    return e.assertInternalFact(internalFact)
}

// convertToInternalFact convertit un fait tsdio en fait interne
func (e *Engine) convertToInternalFact(fact *tsdio.Fact) (*rete.Fact, error) {
    // Vérifier que _id_ n'est pas dans les champs
    if _, exists := fact.Fields[constraint.FieldNameInternalID]; exists {
        return nil, fmt.Errorf(
            "le champ '%s' ne peut pas être défini manuellement",
            constraint.FieldNameInternalID,
        )
    }
    
    // Convertir en structure interne
    return &rete.Fact{
        Type:   fact.Type,
        Fields: fact.Fields,
        // L'ID sera généré plus tard
    }, nil
}
```

### 6. Validation Côté API

#### Nouveau Fichier : `api/validator.go`

**Validation des entrées utilisateur** :

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package api

import (
    "fmt"
    "github.com/resinsec/tsd/constraint"
)

// APIValidator valide les entrées de l'API
type APIValidator struct {
    // Configuration de validation
}

// NewAPIValidator crée un nouveau validateur API
func NewAPIValidator() *APIValidator {
    return &APIValidator{}
}

// ValidateFact valide un fait soumis via l'API
func (v *APIValidator) ValidateFact(fact interface{}) error {
    // Vérifier que c'est une map ou une structure
    factMap, ok := fact.(map[string]interface{})
    if !ok {
        return fmt.Errorf("le fait doit être un objet JSON")
    }
    
    // Vérifier que _id_ n'est pas présent
    if _, exists := factMap[constraint.FieldNameInternalID]; exists {
        return fmt.Errorf(
            "le champ '%s' est réservé et ne peut pas être défini via l'API",
            constraint.FieldNameInternalID,
        )
    }
    
    // Vérifier que le type est présent
    if _, exists := factMap["type"]; !exists {
        return fmt.Errorf("le champ 'type' est requis")
    }
    
    // Vérifier que les fields sont présents
    if _, exists := factMap["fields"]; !exists {
        return fmt.Errorf("le champ 'fields' est requis")
    }
    
    return nil
}

// ValidateFactAssignment valide une affectation via l'API
func (v *APIValidator) ValidateFactAssignment(assignment interface{}) error {
    assignmentMap, ok := assignment.(map[string]interface{})
    if !ok {
        return fmt.Errorf("l'affectation doit être un objet JSON")
    }
    
    // Vérifier la présence de 'variable'
    variable, exists := assignmentMap["variable"]
    if !exists {
        return fmt.Errorf("le champ 'variable' est requis")
    }
    
    // Vérifier que c'est une string
    if _, ok := variable.(string); !ok {
        return fmt.Errorf("le champ 'variable' doit être une string")
    }
    
    // Vérifier la présence de 'fact'
    fact, exists := assignmentMap["fact"]
    if !exists {
        return fmt.Errorf("le champ 'fact' est requis")
    }
    
    // Valider le fait
    return v.ValidateFact(fact)
}

// SanitizeFact nettoie un fait en retirant les champs interdits
func (v *APIValidator) SanitizeFact(fact map[string]interface{}) map[string]interface{} {
    cleaned := make(map[string]interface{})
    
    for key, value := range fact {
        // Exclure _id_
        if key != constraint.FieldNameInternalID {
            cleaned[key] = value
        }
    }
    
    return cleaned
}
```

#### Tests du Validateur API

**Fichier : `api/validator_test.go`**

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package api

import (
    "testing"
    "github.com/resinsec/tsd/constraint"
)

func TestAPIValidator_ValidateFact(t *testing.T) {
    t.Log("🧪 TEST API VALIDATOR - VALIDATION FAIT")
    t.Log("========================================")
    
    validator := NewAPIValidator()
    
    tests := []struct {
        name    string
        fact    interface{}
        wantErr bool
    }{
        {
            name: "fait valide",
            fact: map[string]interface{}{
                "type": "User",
                "fields": map[string]interface{}{
                    "name": "Alice",
                },
            },
            wantErr: false,
        },
        {
            name: "_id_ présent",
            fact: map[string]interface{}{
                "type": "User",
                "fields": map[string]interface{}{
                    constraint.FieldNameInternalID: "manual-id",
                    "name": "Alice",
                },
            },
            wantErr: true,
        },
        {
            name: "type manquant",
            fact: map[string]interface{}{
                "fields": map[string]interface{}{
                    "name": "Alice",
                },
            },
            wantErr: true,
        },
        {
            name: "fields manquant",
            fact: map[string]interface{}{
                "type": "User",
            },
            wantErr: true,
        },
        {
            name:    "pas un objet",
            fact:    "invalid",
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validator.ValidateFact(tt.fact)
            
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

func TestAPIValidator_SanitizeFact(t *testing.T) {
    t.Log("🧪 TEST API VALIDATOR - NETTOYAGE")
    t.Log("==================================")
    
    validator := NewAPIValidator()
    
    fact := map[string]interface{}{
        "type": "User",
        "fields": map[string]interface{}{
            constraint.FieldNameInternalID: "manual-id",
            "name": "Alice",
            "age":  30,
        },
    }
    
    cleaned := validator.SanitizeFact(fact)
    
    // Vérifier que _id_ a été retiré
    if _, exists := cleaned[constraint.FieldNameInternalID]; exists {
        t.Error("❌ Le champ '_id_' devrait être retiré")
    }
    
    // Vérifier que les autres champs sont préservés
    if _, exists := cleaned["type"]; !exists {
        t.Error("❌ Le champ 'type' devrait être préservé")
    }
    
    if _, exists := cleaned["fields"]; !exists {
        t.Error("❌ Le champ 'fields' devrait être préservé")
    }
    
    t.Log("✅ Nettoyage réussi")
}
```

### 7. Documentation de l'API

#### Nouveau Fichier : `api/README.md`

```markdown
# API TSD - Documentation

## Vue d'Ensemble

L'API TSD permet de gérer des faits et des règles dans le moteur RETE.

**Note importante** : L'identifiant interne des faits (`_id_`) est géré automatiquement par le système et n'est **jamais** exposé via l'API publique.

## Structures de Données

### Fact

Un fait représente une donnée dans le système.

```json
{
  "type": "User",
  "fields": {
    "name": "Alice",
    "age": 30
  }
}
```

**Champs** :
- `type` (string, requis) : Le type du fait
- `fields` (object, requis) : Les champs du fait

**Restrictions** :
- Le champ `_id_` ne peut **jamais** être défini manuellement
- Il est généré automatiquement lors de l'insertion

### FactAssignment

Une affectation de fait à une variable.

```json
{
  "variable": "alice",
  "fact": {
    "type": "User",
    "fields": {
      "name": "Alice",
      "age": 30
    }
  }
}
```

**Champs** :
- `variable` (string, requis) : Nom de la variable
- `fact` (Fact, requis) : Le fait à affecter

## Endpoints (Exemple)

### POST /facts

Asserte un nouveau fait.

**Request** :
```json
{
  "type": "User",
  "fields": {
    "name": "Alice",
    "age": 30
  }
}
```

**Response** :
```json
{
  "success": true,
  "message": "Fait asserté avec succès"
}
```

### POST /assignments

Asserte une affectation de variable.

**Request** :
```json
{
  "variable": "alice",
  "fact": {
    "type": "User",
    "fields": {
      "name": "Alice",
      "age": 30
    }
  }
}
```

**Response** :
```json
{
  "success": true,
  "variable": "alice",
  "message": "Affectation créée avec succès"
}
```

### GET /results

Récupère les résultats.

**Response** :
```json
{
  "facts": [
    {
      "type": "User",
      "fields": {
        "name": "Alice",
        "age": 30
      }
    }
  ],
  "metadata": {
    "count": 1
  }
}
```

**Note** : L'ID interne (`_id_`) n'est **jamais** retourné dans les résultats.

## Erreurs

### Champ Réservé

**Erreur** :
```json
{
  "error": "le champ '_id_' est réservé et ne peut pas être défini via l'API"
}
```

**Cause** : Tentative de définir manuellement le champ `_id_`.

**Solution** : Retirer le champ `_id_` de la requête. Il sera généré automatiquement.

## Exemples d'Utilisation

### Créer des Faits avec Affectations

```bash
# 1. Créer un utilisateur et l'affecter à une variable
curl -X POST http://localhost:8080/assignments \
  -H "Content-Type: application/json" \
  -d '{
    "variable": "alice",
    "fact": {
      "type": "User",
      "fields": {
        "name": "Alice",
        "age": 30
      }
    }
  }'

# 2. Créer un login qui référence la variable
curl -X POST http://localhost:8080/facts \
  -H "Content-Type: application/json" \
  -d '{
    "type": "Login",
    "fields": {
      "user": {"variable": "alice"},
      "email": "alice@example.com",
      "password": "secret"
    }
  }'
```

### Récupérer les Résultats

```bash
curl -X GET http://localhost:8080/results
```

## Sécurité

### Validation Automatique

Tous les faits soumis via l'API sont automatiquement validés :

1. **Structure** : Vérification des champs requis
2. **Champs réservés** : Interdiction de `_id_`
3. **Types** : Validation des types de données

### Sanitization

Les champs interdits sont automatiquement retirés des réponses pour garantir que les détails internes ne sont jamais exposés.
```

---

## ✅ Critères de Succès

### Compilation et Tests

```bash
# Code compile
go build ./api
go build ./tsdio

# Tests passent
go test ./api -v
go test ./tsdio -v

# Couverture > 80%
go test ./api -cover
go test ./tsdio -cover
```

### Fonctionnalités

- [ ] `_id_` caché dans toutes les API
- [ ] Support des affectations dans tsdio
- [ ] Validation côté API implémentée
- [ ] Sérialisation JSON sécurisée
- [ ] Méthodes d'accès à l'ID interne
- [ ] Documentation API complète

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

- [ ] `TestFact_NewFact`
- [ ] `TestFact_Validate`
- [ ] `TestFact_JSONSerialization`
- [ ] `TestFact_InternalIDMethods`
- [ ] `TestFactAssignment_Creation`
- [ ] `TestProgram_AddFactAssignment`
- [ ] `TestProgram_Validate`
- [ ] `TestResultFact_FromXuple`
- [ ] `TestResultFact_JSONSerialization`
- [ ] `TestAPIValidator_ValidateFact`
- [ ] `TestAPIValidator_SanitizeFact`

### Tests d'Intégration

```go
func TestAPI_CompleteFlow(t *testing.T) {
    // 1. Créer un programme
    program := tsdio.NewProgram()
    
    // 2. Ajouter une affectation
    userFact, _ := tsdio.NewFact("User", map[string]interface{}{
        "name": "Alice",
        "age":  30,
    })
    assignment, _ := tsdio.NewFactAssignment("alice", userFact)
    program.AddFactAssignment(assignment)
    
    // 3. Ajouter un fait qui référence la variable
    loginFact, _ := tsdio.NewFact("Login", map[string]interface{}{
        "user":     map[string]string{"variable": "alice"},
        "email":    "alice@example.com",
        "password": "secret",
    })
    program.AddFact(loginFact)
    
    // 4. Valider le programme
    if err := program.Validate(); err != nil {
        t.Fatalf("Erreur de validation: %v", err)
    }
    
    // 5. Sérialiser en JSON
    jsonData, err := json.Marshal(program)
    if err != nil {
        t.Fatalf("Erreur de sérialisation: %v", err)
    }
    
    // 6. Vérifier que _id_ n'est pas dans le JSON
    jsonStr := string(jsonData)
    if strings.Contains(jsonStr, "_id_") {
        t.Error("Le JSON contient '_id_'")
    }
    
    t.Log("✅ Flow complet API réussi")
}
```

---

## 🚀 Exécution

### Ordre des Modifications

1. ✅ Modifier structures Fact dans tsdio
2. ✅ Ajouter support FactAssignment
3. ✅ Modifier API Result
4. ✅ Adapter Engine
5. ✅ Créer validateur API
6. ✅ Tests unitaires
7. ✅ Tests d'intégration
8. ✅ Documentation
9. ✅ Validation finale

### Commandes

```bash
# Modifier les fichiers
vim tsdio/api.go
vim tsdio/program.go
vim api/result.go
vim api/engine.go
vim api/validator.go

# Tester
go test ./tsdio -v
go test ./api -v

# Validation
make validate
make test-complete
```

---

## 📚 Références

- `scripts/new_ids/05-prompt-types-validation.md` - Validation
- `scripts/new_ids/04-prompt-evaluation.md` - Évaluation
- `tsdio/` - API actuelle
- `api/` - API actuelle

---

## 📝 Notes

### Points d'Attention

1. **Sécurité** : `_id_` ne doit JAMAIS être exposé publiquement

2. **Sérialisation JSON** : Utiliser des méthodes personnalisées pour filtrer `_id_`

3. **Validation** : Valider côté API pour éviter les erreurs silencieuses

4. **Documentation** : API publique doit être clairement documentée

### Questions Résolues

Q: Faut-il un alias `id` en lecture seule ?
R: Non, `_id_` est complètement caché de l'API publique

Q: Comment gérer les migrations ?
R: Documenter le breaking change, fournir guide de migration

---

## 🎯 Résultat Attendu

```go
// ✅ API publique sans _id_
fact := tsdio.NewFact("User", map[string]interface{}{
    "name": "Alice",
    "age":  30,
})

// ✅ Sérialisation JSON ne contient pas _id_
json := `{"type":"User","fields":{"name":"Alice","age":30}}`

// ✅ Affectations supportées
assignment := tsdio.NewFactAssignment("alice", userFact)
program.AddFactAssignment(assignment)

// ✅ Validation automatique
validator.ValidateFact(fact) // Rejette si _id_ présent

// ❌ Erreurs détectées
NewFact("User", map[string]interface{}{
    "_id_": "manual",  // Erreur: champ réservé
})
```

---

**Prompt suivant** : `07-prompt-tests-unit.md`

**Durée estimée** : 4-5 heures

**Complexité** : 🟡 Moyenne