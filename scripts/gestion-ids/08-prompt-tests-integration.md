# Prompt 08 : Tests d'intégration et end-to-end

**Objectif** : Créer des tests d'intégration et end-to-end pour valider le fonctionnement complet de la génération automatique d'IDs avec clés primaires dans un contexte réel d'utilisation.

**Prérequis** : Prompts 01-07 complétés et validés.

---

## Contexte

Les modifications ont été apportées à plusieurs niveaux :
- Grammaire et parsing (`#` pour marquer les PK)
- Structures de données (Field.IsPrimaryKey)
- Validation (interdiction de `id` explicite, validation des PK)
- Génération d'IDs (PK ou hash)
- Intégration RETE (working memory, évaluateur, joins)

Il faut maintenant valider le système complet avec des scénarios réalistes qui combinent tous ces éléments.

---

## Tâches

### 8.1. Tests d'intégration parsing → validation → génération d'IDs

**Fichier** : `constraint/integration_test.go` (créer)

**Action** : Tester le flux complet depuis le parsing jusqu'à la génération des IDs RETE.

```go
// Copyright 2025 SekiaTech. All rights reserved.
// Use of this source code is governed by an MIT-style license.

package constraint

import (
    "strings"
    "testing"
)

func TestIntegration_ParseAndGenerateIDs(t *testing.T) {
    tests := []struct {
        name        string
        input       string
        wantErr     bool
        checkIDsFn  func(*testing.T, *Program)
    }{
        {
            name: "Programme complet avec PK simple",
            input: `
type Person(#nom: string, age: number)

assert Person(nom: "Alice", age: 30)
assert Person(nom: "Bob", age: 25)
`,
            wantErr: false,
            checkIDsFn: func(t *testing.T, prog *Program) {
                reteFacts, err := ConvertFactsToReteFormat(prog)
                if err != nil {
                    t.Fatalf("ConvertFactsToReteFormat() error = %v", err)
                }
                
                if len(reteFacts) != 2 {
                    t.Fatalf("Expected 2 facts, got %d", len(reteFacts))
                }
                
                // Vérifier les IDs générés
                expectedIDs := map[string]bool{
                    "Person~Alice": false,
                    "Person~Bob":   false,
                }
                
                for _, fact := range reteFacts {
                    if _, ok := expectedIDs[fact.ID]; ok {
                        expectedIDs[fact.ID] = true
                    } else {
                        t.Errorf("Unexpected fact ID: %s", fact.ID)
                    }
                }
                
                for id, found := range expectedIDs {
                    if !found {
                        t.Errorf("Expected fact with ID %s not found", id)
                    }
                }
            },
        },
        {
            name: "Programme complet avec PK composite",
            input: `
type Person(#prenom: string, #nom: string, age: number)

assert Person(prenom: "Alice", nom: "Dupont", age: 30)
assert Person(prenom: "Bob", nom: "Martin", age: 25)
`,
            wantErr: false,
            checkIDsFn: func(t *testing.T, prog *Program) {
                reteFacts, err := ConvertFactsToReteFormat(prog)
                if err != nil {
                    t.Fatalf("ConvertFactsToReteFormat() error = %v", err)
                }
                
                expectedIDs := []string{
                    "Person~Alice_Dupont",
                    "Person~Bob_Martin",
                }
                
                for i, fact := range reteFacts {
                    if fact.ID != expectedIDs[i] {
                        t.Errorf("Fact %d: expected ID %s, got %s", i, expectedIDs[i], fact.ID)
                    }
                }
            },
        },
        {
            name: "Programme avec type sans PK (hash)",
            input: `
type Event(timestamp: number, message: string)

assert Event(timestamp: 1234567890, message: "test1")
assert Event(timestamp: 1234567891, message: "test2")
`,
            wantErr: false,
            checkIDsFn: func(t *testing.T, prog *Program) {
                reteFacts, err := ConvertFactsToReteFormat(prog)
                if err != nil {
                    t.Fatalf("ConvertFactsToReteFormat() error = %v", err)
                }
                
                // Vérifier le format hash
                for i, fact := range reteFacts {
                    if !strings.HasPrefix(fact.ID, "Event~") {
                        t.Errorf("Fact %d: expected ID to start with Event~, got %s", i, fact.ID)
                    }
                    
                    hashPart := strings.TrimPrefix(fact.ID, "Event~")
                    if len(hashPart) != 16 {
                        t.Errorf("Fact %d: expected hash of 16 chars, got %d: %s", i, len(hashPart), hashPart)
                    }
                }
                
                // Vérifier que les deux faits ont des IDs différents
                if reteFacts[0].ID == reteFacts[1].ID {
                    t.Errorf("Expected different IDs for different facts, got both: %s", reteFacts[0].ID)
                }
            },
        },
        {
            name: "Rejet de id explicite dans assertion",
            input: `
type Person(#nom: string)

assert Person(id: "custom_id", nom: "Alice")
`,
            wantErr: true,
        },
        {
            name: "Caractères spéciaux dans PK",
            input: `
type Resource(#path: string)

assert Resource(path: "/home/user~test_file")
`,
            wantErr: false,
            checkIDsFn: func(t *testing.T, prog *Program) {
                reteFacts, err := ConvertFactsToReteFormat(prog)
                if err != nil {
                    t.Fatalf("ConvertFactsToReteFormat() error = %v", err)
                }
                
                // Vérifier l'échappement
                expectedID := "Resource~%2Fhome%2Fuser%7Etest%5Ffile"
                if reteFacts[0].ID != expectedID {
                    t.Errorf("Expected ID %s, got %s", expectedID, reteFacts[0].ID)
                }
            },
        },
        {
            name: "Plusieurs types avec stratégies différentes",
            input: `
type Person(#nom: string, age: number)
type Event(timestamp: number, message: string)

assert Person(nom: "Alice", age: 30)
assert Event(timestamp: 1234567890, message: "User logged in")
`,
            wantErr: false,
            checkIDsFn: func(t *testing.T, prog *Program) {
                reteFacts, err := ConvertFactsToReteFormat(prog)
                if err != nil {
                    t.Fatalf("ConvertFactsToReteFormat() error = %v", err)
                }
                
                // Person doit avoir un ID basé sur PK
                personFact := reteFacts[0]
                if personFact.Type == "Person" && personFact.ID != "Person~Alice" {
                    t.Errorf("Person fact: expected ID Person~Alice, got %s", personFact.ID)
                }
                
                // Event doit avoir un ID basé sur hash
                eventFact := reteFacts[1]
                if eventFact.Type == "Event" && !strings.HasPrefix(eventFact.ID, "Event~") {
                    t.Errorf("Event fact: expected ID to start with Event~, got %s", eventFact.ID)
                }
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            prog, err := Parse(tt.input)
            if err != nil {
                if !tt.wantErr {
                    t.Fatalf("Parse() error = %v, wantErr %v", err, tt.wantErr)
                }
                return
            }
            
            // Validation
            err = ValidateFacts(prog)
            if (err != nil) != tt.wantErr {
                t.Fatalf("ValidateFacts() error = %v, wantErr %v", err, tt.wantErr)
            }
            
            if err == nil && tt.checkIDsFn != nil {
                tt.checkIDsFn(t, prog)
            }
        })
    }
}

func TestIntegration_IDDeterminism(t *testing.T) {
    input := `
type Person(#nom: string, age: number)
type Event(timestamp: number, message: string)

assert Person(nom: "Alice", age: 30)
assert Event(timestamp: 1234567890, message: "test")
`

    // Parser et convertir plusieurs fois
    var allIDs [][]string
    for i := 0; i < 5; i++ {
        prog, err := Parse(input)
        if err != nil {
            t.Fatalf("Parse() iteration %d error = %v", i, err)
        }
        
        reteFacts, err := ConvertFactsToReteFormat(prog)
        if err != nil {
            t.Fatalf("ConvertFactsToReteFormat() iteration %d error = %v", i, err)
        }
        
        ids := make([]string, len(reteFacts))
        for j, fact := range reteFacts {
            ids[j] = fact.ID
        }
        allIDs = append(allIDs, ids)
    }
    
    // Vérifier que tous les runs ont produit les mêmes IDs
    firstRun := allIDs[0]
    for i, ids := range allIDs[1:] {
        for j, id := range ids {
            if id != firstRun[j] {
                t.Errorf("Run %d, fact %d: ID mismatch: got %s, want %s", i+1, j, id, firstRun[j])
            }
        }
    }
}
```

### 8.2. Tests end-to-end avec règles

**Fichier** : `integration_test.go` (à la racine du projet, ou dans un package `integration/`)

**Action** : Tester des scénarios complets avec parsing, validation, génération d'IDs et exécution de règles.

```go
// Copyright 2025 SekiaTech. All rights reserved.
// Use of this source code is governed by an MIT-style license.

package integration

import (
    "testing"
    
    "github.com/sekiatech/tsd/constraint"
    "github.com/sekiatech/tsd/rete"
)

func TestE2E_RulesWithGeneratedIDs(t *testing.T) {
    tests := []struct {
        name           string
        program        string
        wantErr        bool
        expectedFires  int
        checkOutputFn  func(*testing.T, []interface{})
    }{
        {
            name: "Règle avec accès au champ id",
            program: `
type Person(#nom: string, age: number)

assert Person(nom: "Alice", age: 30)
assert Person(nom: "Bob", age: 25)

rule CheckAliceID {
    when {
        p: Person()
        p.id == "Person~Alice"
    }
    then {
        print("Found Alice with ID: " + p.id)
    }
}
`,
            wantErr:       false,
            expectedFires: 1,
        },
        {
            name: "Règle avec join sur IDs",
            program: `
type Person(#nom: string, age: number)
type Membership(#person_nom: string, #club: string)

assert Person(nom: "Alice", age: 30)
assert Person(nom: "Bob", age: 25)
assert Membership(person_nom: "Alice", club: "Chess")
assert Membership(person_nom: "Bob", club: "Tennis")

rule PersonWithMembership {
    when {
        p: Person()
        m: Membership()
        p.nom == m.person_nom
    }
    then {
        print(p.nom + " is member of " + m.club)
    }
}
`,
            wantErr:       false,
            expectedFires: 2,
        },
        {
            name: "Règle avec PK composite",
            program: `
type Person(#prenom: string, #nom: string, age: number)

assert Person(prenom: "Alice", nom: "Dupont", age: 30)
assert Person(prenom: "Bob", nom: "Martin", age: 25)

rule AdultWithCompositeKey {
    when {
        p: Person()
        p.age >= 18
    }
    then {
        print(p.prenom + " " + p.nom + " (ID: " + p.id + ") is an adult")
    }
}
`,
            wantErr:       false,
            expectedFires: 2,
        },
        {
            name: "Règle avec type sans PK (hash)",
            program: `
type Event(timestamp: number, message: string)

assert Event(timestamp: 1234567890, message: "Login")
assert Event(timestamp: 1234567891, message: "Logout")

rule AllEvents {
    when {
        e: Event()
    }
    then {
        print("Event ID: " + e.id)
    }
}
`,
            wantErr:       false,
            expectedFires: 2,
        },
        {
            name: "Règle avec comparaison d'IDs entre faits",
            program: `
type Person(#nom: string, age: number)
type Friend(#person1: string, #person2: string)

assert Person(nom: "Alice", age: 30)
assert Person(nom: "Bob", age: 25)
assert Friend(person1: "Alice", person2: "Bob")

rule Friendship {
    when {
        p1: Person()
        p2: Person()
        f: Friend()
        p1.nom == f.person1
        p2.nom == f.person2
        p1.id != p2.id
    }
    then {
        print(p1.nom + " is friend with " + p2.nom)
    }
}
`,
            wantErr:       false,
            expectedFires: 1,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Parser le programme
            prog, err := constraint.Parse(tt.program)
            if err != nil {
                if !tt.wantErr {
                    t.Fatalf("Parse() error = %v", err)
                }
                return
            }
            
            // Valider
            err = constraint.ValidateFacts(prog)
            if err != nil {
                if !tt.wantErr {
                    t.Fatalf("ValidateFacts() error = %v", err)
                }
                return
            }
            
            // Convertir les faits
            reteFacts, err := constraint.ConvertFactsToReteFormat(prog)
            if err != nil {
                t.Fatalf("ConvertFactsToReteFormat() error = %v", err)
            }
            
            // Créer le réseau RETE et la working memory
            network := rete.NewReteNetwork()
            wm := rete.NewWorkingMemory()
            
            // Ajouter les faits
            for _, fact := range reteFacts {
                wm.AddFact(fact)
            }
            
            // Compiler et exécuter les règles
            // (adapter selon l'API réelle du projet)
            fireCount := 0
            for _, rule := range prog.Rules {
                matches := network.EvaluateRule(rule, wm)
                fireCount += len(matches)
            }
            
            if fireCount != tt.expectedFires {
                t.Errorf("Expected %d rule fires, got %d", tt.expectedFires, fireCount)
            }
            
            if tt.checkOutputFn != nil {
                // Récupérer les outputs/actions des règles
                outputs := network.GetOutputs()
                tt.checkOutputFn(t, outputs)
            }
        })
    }
}
```

### 8.3. Tests avec fichiers .tsd

**Dossier** : `testdata/integration/` (créer)

**Action** : Créer des fichiers `.tsd` réalistes et des tests qui les chargent.

**Fichier 1** : `testdata/integration/pk_simple.tsd`

```
type Person(#nom: string, age: number, ville: string)

assert Person(nom: "Alice", age: 30, ville: "Paris")
assert Person(nom: "Bob", age: 25, ville: "Lyon")
assert Person(nom: "Charlie", age: 35, ville: "Paris")

rule AdultsInParis {
    when {
        p: Person()
        p.age >= 18
        p.ville == "Paris"
    }
    then {
        print(p.nom + " is an adult in Paris (ID: " + p.id + ")")
    }
}
```

**Fichier 2** : `testdata/integration/pk_composite.tsd`

```
type Produit(#categorie: string, #nom: string, prix: number, stock: number)

assert Produit(categorie: "Electronique", nom: "Laptop", prix: 1200, stock: 5)
assert Produit(categorie: "Electronique", nom: "Souris", prix: 25, stock: 50)
assert Produit(categorie: "Livre", nom: "TSD Guide", prix: 30, stock: 100)

rule ProduitEnRupture {
    when {
        p: Produit()
        p.stock < 10
    }
    then {
        print("Stock faible pour " + p.categorie + "/" + p.nom + " (ID: " + p.id + ")")
    }
}
```

**Fichier 3** : `testdata/integration/no_pk_hash.tsd`

```
type LogEntry(timestamp: number, level: string, message: string)

assert LogEntry(timestamp: 1704067200, level: "INFO", message: "Application started")
assert LogEntry(timestamp: 1704067201, level: "WARN", message: "High memory usage")
assert LogEntry(timestamp: 1704067202, level: "ERROR", message: "Connection failed")

rule ErrorLogs {
    when {
        log: LogEntry()
        log.level == "ERROR"
    }
    then {
        print("Error log found (ID: " + log.id + "): " + log.message)
    }
}
```

**Fichier 4** : `testdata/integration/mixed_types.tsd`

```
type User(#username: string, email: string)
type Session(session_id: string, username: string, active: bool)

assert User(username: "alice", email: "alice@example.com")
assert User(username: "bob", email: "bob@example.com")
assert Session(session_id: "sess_123", username: "alice", active: true)
assert Session(session_id: "sess_456", username: "bob", active: false)

rule ActiveUserSessions {
    when {
        u: User()
        s: Session()
        s.username == u.username
        s.active == true
    }
    then {
        print("Active session for user " + u.username + " (User ID: " + u.id + ", Session ID: " + s.id + ")")
    }
}
```

**Fichier de test** : `integration_test.go`

```go
func TestE2E_TSDFiles(t *testing.T) {
    testFiles := []struct {
        name          string
        file          string
        expectedFires int
    }{
        {"PK simple", "testdata/integration/pk_simple.tsd", 2},
        {"PK composite", "testdata/integration/pk_composite.tsd", 1},
        {"No PK (hash)", "testdata/integration/no_pk_hash.tsd", 1},
        {"Mixed types", "testdata/integration/mixed_types.tsd", 1},
    }

    for _, tt := range testFiles {
        t.Run(tt.name, func(t *testing.T) {
            // Lire le fichier
            content, err := os.ReadFile(tt.file)
            if err != nil {
                t.Fatalf("Failed to read file %s: %v", tt.file, err)
            }
            
            // Parser
            prog, err := constraint.Parse(string(content))
            if err != nil {
                t.Fatalf("Parse() error = %v", err)
            }
            
            // Valider
            err = constraint.ValidateFacts(prog)
            if err != nil {
                t.Fatalf("ValidateFacts() error = %v", err)
            }
            
            // Convertir et exécuter
            reteFacts, err := constraint.ConvertFactsToReteFormat(prog)
            if err != nil {
                t.Fatalf("ConvertFactsToReteFormat() error = %v", err)
            }
            
            // Vérifier que tous les faits ont des IDs
            for i, fact := range reteFacts {
                if fact.ID == "" {
                    t.Errorf("Fact %d has empty ID", i)
                }
                // Vérifier le format (TypeName~...)
                if !strings.Contains(fact.ID, "~") {
                    t.Errorf("Fact %d has invalid ID format: %s", i, fact.ID)
                }
            }
            
            // Exécuter les règles
            network := rete.NewReteNetwork()
            wm := rete.NewWorkingMemory()
            
            for _, fact := range reteFacts {
                wm.AddFact(fact)
            }
            
            fireCount := 0
            for _, rule := range prog.Rules {
                matches := network.EvaluateRule(rule, wm)
                fireCount += len(matches)
            }
            
            if fireCount != tt.expectedFires {
                t.Errorf("File %s: expected %d rule fires, got %d", tt.file, tt.expectedFires, fireCount)
            }
        })
    }
}
```

### 8.4. Tests de régression

**Fichier** : `integration_test.go`

**Action** : Vérifier que les programmes existants (sans `#`) continuent de fonctionner.

```go
func TestE2E_BackwardCompatibility(t *testing.T) {
    // Programme sans clés primaires (ancien format)
    program := `
type Person(nom: string, age: number)

assert Person(nom: "Alice", age: 30)
assert Person(nom: "Bob", age: 25)

rule Adults {
    when {
        p: Person()
        p.age >= 18
    }
    then {
        print(p.nom + " is an adult")
    }
}
`

    prog, err := constraint.Parse(program)
    if err != nil {
        t.Fatalf("Parse() error = %v", err)
    }
    
    // Valider
    err = constraint.ValidateFacts(prog)
    if err != nil {
        t.Fatalf("ValidateFacts() error = %v", err)
    }
    
    // Convertir
    reteFacts, err := constraint.ConvertFactsToReteFormat(prog)
    if err != nil {
        t.Fatalf("ConvertFactsToReteFormat() error = %v", err)
    }
    
    // Vérifier que les IDs sont générés avec hash (pas de PK)
    for _, fact := range reteFacts {
        if !strings.HasPrefix(fact.ID, "Person~") {
            t.Errorf("Expected ID to start with Person~, got %s", fact.ID)
        }
        hashPart := strings.TrimPrefix(fact.ID, "Person~")
        if len(hashPart) != 16 {
            t.Errorf("Expected hash of 16 chars, got %d: %s", len(hashPart), hashPart)
        }
    }
    
    // Les règles doivent continuer à fonctionner
    network := rete.NewReteNetwork()
    wm := rete.NewWorkingMemory()
    
    for _, fact := range reteFacts {
        wm.AddFact(fact)
    }
    
    fireCount := 0
    for _, rule := range prog.Rules {
        matches := network.EvaluateRule(rule, wm)
        fireCount += len(matches)
    }
    
    if fireCount != 2 {
        t.Errorf("Expected 2 rule fires, got %d", fireCount)
    }
}
```

### 8.5. Tests de cas d'erreur end-to-end

**Fichier** : `integration_test.go`

**Action** : Tester que les erreurs sont correctement détectées et reportées.

```go
func TestE2E_ErrorCases(t *testing.T) {
    tests := []struct {
        name        string
        program     string
        wantErr     bool
        errContains string
    }{
        {
            name: "ID explicite interdit",
            program: `
type Person(#nom: string)
assert Person(id: "custom", nom: "Alice")
`,
            wantErr:     true,
            errContains: "id",
        },
        {
            name: "Champ PK manquant dans fact",
            program: `
type Person(#nom: string, age: number)
assert Person(age: 30)
`,
            wantErr:     true,
            errContains: "nom",
        },
        {
            name: "Type de PK invalide (objet)",
            program: `
type Person(#data: object, age: number)
assert Person(data: {key: "value"}, age: 30)
`,
            wantErr:     true,
            errContains: "primary key",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            prog, err := constraint.Parse(tt.program)
            if err != nil {
                if !tt.wantErr {
                    t.Fatalf("Parse() error = %v", err)
                }
                if !strings.Contains(err.Error(), tt.errContains) {
                    t.Errorf("Expected error to contain %q, got: %v", tt.errContains, err)
                }
                return
            }
            
            err = constraint.ValidateFacts(prog)
            if (err != nil) != tt.wantErr {
                t.Fatalf("ValidateFacts() error = %v, wantErr %v", err, tt.wantErr)
            }
            
            if err != nil && tt.errContains != "" {
                if !strings.Contains(err.Error(), tt.errContains) {
                    t.Errorf("Expected error to contain %q, got: %v", tt.errContains, err)
                }
            }
        })
    }
}
```

---

## Validation

### Étape 1 : Créer les fichiers de test

```bash
cd /home/resinsec/dev/tsd
mkdir -p testdata/integration
# Créer les fichiers .tsd listés ci-dessus
```

### Étape 2 : Exécuter les tests d'intégration constraint

```bash
go test ./constraint/... -run TestIntegration -v
```

### Étape 3 : Exécuter les tests e2e

```bash
go test ./integration/... -v
# ou
go test -run TestE2E ./... -v
```

### Étape 4 : Exécuter tous les tests

```bash
make test-complete
# ou
go test ./... -v
```

### Étape 5 : Test manuel avec CLI (si disponible)

```bash
# Tester avec un fichier .tsd
./tsd run testdata/integration/pk_simple.tsd

# Vérifier la sortie
# Devrait afficher les IDs générés et les résultats des règles
```

### Étape 6 : Validation finale

```bash
make validate
```

---

## Checklist

- [ ] Répertoire `testdata/integration/` créé
- [ ] Fichiers `.tsd` de test créés (pk_simple, pk_composite, no_pk_hash, mixed_types)
- [ ] Tests d'intégration parsing→validation→génération ajoutés
- [ ] Tests de déterminisme ajoutés
- [ ] Tests e2e avec règles ajoutés
- [ ] Tests e2e avec fichiers .tsd ajoutés
- [ ] Tests de rétrocompatibilité ajoutés
- [ ] Tests de cas d'erreur e2e ajoutés
- [ ] `go test ./constraint/... -run TestIntegration -v` réussit
- [ ] `go test -run TestE2E ./... -v` réussit
- [ ] `make test-complete` réussit
- [ ] Test manuel avec CLI réussit (optionnel)
- [ ] `make validate` réussit

---

## Rapport

Une fois toutes les tâches complétées :

1. Lister les fichiers de test créés
2. Lister les scénarios testés
3. Indiquer le nombre de tests qui passent
4. Copier la sortie des tests réussis
5. Documenter tout problème découvert
6. Commit :

```bash
git add .
git commit -m "test(integration): tests e2e complets pour clés primaires et génération d'IDs

- Tests d'intégration parsing → validation → génération d'IDs
- Tests e2e avec règles (accès id, joins, comparaisons)
- Fichiers .tsd de test avec différents scénarios
- Tests de rétrocompatibilité (programmes sans PK)
- Tests de cas d'erreur end-to-end
- Tests de déterminisme des IDs générés

Scénarios testés:
- PK simple, composite, sans PK (hash)
- Accès au champ id dans les règles
- Joins sur IDs
- Comparaisons d'IDs
- Échappement de caractères spéciaux
- Validation des erreurs

Tous les tests passent.

Refs #<issue_number>"
```

---

## Dépendances

- **Bloque** : Prompt 09 (mise à jour exemples)
- **Bloqué par** : Prompts 01-07

---

## Notes importantes

1. **API RETE** : Les exemples de code ci-dessus supposent une certaine API pour le module RETE (`NewReteNetwork()`, `EvaluateRule()`, etc.). Adaptez selon l'API réelle du projet.

2. **Fichiers .tsd** : Les fichiers de test doivent être représentatifs de cas d'usage réels. N'hésitez pas à ajouter d'autres scénarios si nécessaire.

3. **Outputs des règles** : Si le système TSD a un mécanisme pour capturer les outputs des actions `print()` ou autres, utilisez-le dans les tests pour vérifier le comportement complet.

4. **Performance** : Les tests e2e peuvent être plus lents. Envisagez de les marquer avec un build tag `// +build integration` pour les exécuter séparément si nécessaire.

5. **Bug numérique** : Si le bug de l'évaluateur numérique n'est pas encore fixé, certains tests peuvent échouer. Documentez-les et ajoutez des `t.Skip()` temporaires si besoin, avec un commentaire expliquant le bug.

6. **CLI** : Si le CLI n'est pas encore disponible ou fonctionnel, les tests manuels peuvent être sautés. L'essentiel est que les tests Go automatisés passent.

7. **Couverture** : Visez une couverture complète des scénarios :
   - ✓ PK simple
   - ✓ PK composite
   - ✓ Sans PK (hash)
   - ✓ Types mixtes
   - ✓ Échappement de caractères
   - ✓ Accès à `id` dans règles
   - ✓ Joins sur IDs
   - ✓ Comparaisons d'IDs
   - ✓ Validation des erreurs
   - ✓ Rétrocompatibilité

---

**Prêt à passer au prompt 09 après validation complète.**