# Prompt 06 : Tests du module constraint

**Objectif** : Compléter la couverture de tests du module `constraint` pour les nouvelles fonctionnalités de clés primaires et génération d'IDs. Adapter les tests existants qui pourraient être cassés par les changements.

**Prérequis** : Prompts 01-05 complétés et validés.

---

## Contexte

Les modifications apportées aux prompts 01-04 ont introduit :
- Le marquage de champs avec `#` pour les clés primaires
- La génération automatique des IDs
- L'interdiction de spécifier `id` manuellement dans les assertions
- La validation des clés primaires

Certains tests existants peuvent :
- Définir explicitement des `id` dans les faits (maintenant interdit)
- S'attendre à des IDs au format ancien (`parsed_fact_n`)
- Ne pas couvrir les nouveaux cas d'usage (PK composites, hash, etc.)

---

## Tâches

### 6.1. Inventaire des tests existants

**Fichiers à examiner** :
- `constraint/constraint_test.go`
- `constraint/constraint_facts_test.go`
- `constraint/parser_test.go`
- Tout autre fichier `*_test.go` dans `constraint/`

**Action** : Lire chaque fichier de test et identifier :
1. Les tests qui définissent explicitement des IDs dans les faits
2. Les tests qui vérifient le format des IDs
3. Les tests qui pourraient être impactés par la validation des PK

**Commande** :
```bash
cd /home/resinsec/dev/tsd
grep -n "id:" constraint/*_test.go
grep -n "parsed_fact" constraint/*_test.go
```

### 6.2. Adapter les tests existants cassés

**Pour chaque test qui définit explicitement `id` :**

**Avant** :
```go
fact := &Fact{
    ID:   "person_1",
    Type: "Person",
    Fields: []FactField{
        {Name: "nom", Value: FactValue{StringValue: ptr("Alice")}},
    },
}
```

**Après** : Retirer le champ `ID`, il sera généré automatiquement.
```go
fact := &Fact{
    Type: "Person",
    Fields: []FactField{
        {Name: "nom", Value: FactValue{StringValue: ptr("Alice")}},
    },
}
```

**Pour les tests qui vérifient les IDs** :

**Avant** :
```go
if fact.ID != "parsed_fact_1" {
    t.Errorf("Expected ID parsed_fact_1, got %s", fact.ID)
}
```

**Après** : Vérifier le nouveau format.
```go
// Si Person a une PK sur "nom"
if fact.ID != "Person~Alice" {
    t.Errorf("Expected ID Person~Alice, got %s", fact.ID)
}

// Sinon (pas de PK), vérifier le format hash
if !strings.HasPrefix(fact.ID, "Person~") {
    t.Errorf("Expected ID to start with Person~, got %s", fact.ID)
}
if len(fact.ID) != len("Person~") + 16 { // hash de 16 caractères
    t.Errorf("Expected hash ID of correct length, got %s", fact.ID)
}
```

### 6.3. Tests de parsing avec clés primaires

**Fichier** : `constraint/parser_test.go`

**Action** : Ajouter des tests pour vérifier que le parser reconnaît correctement le `#`.

```go
func TestParser_PrimaryKeyFields(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        wantErr  bool
        checkFn  func(*testing.T, *Program)
    }{
        {
            name: "PK simple sur un champ",
            input: `
type Person(#nom: string, age: number)
`,
            wantErr: false,
            checkFn: func(t *testing.T, prog *Program) {
                if len(prog.TypeDefinitions) != 1 {
                    t.Fatalf("Expected 1 type, got %d", len(prog.TypeDefinitions))
                }
                typeDef := prog.TypeDefinitions[0]
                if typeDef.Name != "Person" {
                    t.Errorf("Expected type Person, got %s", typeDef.Name)
                }
                if len(typeDef.Fields) != 2 {
                    t.Fatalf("Expected 2 fields, got %d", len(typeDef.Fields))
                }
                
                // Vérifier que "nom" est PK
                nomField := typeDef.Fields[0]
                if nomField.Name != "nom" {
                    t.Errorf("Expected first field to be 'nom', got %s", nomField.Name)
                }
                if !nomField.IsPrimaryKey {
                    t.Errorf("Expected 'nom' to be primary key")
                }
                
                // Vérifier que "age" n'est pas PK
                ageField := typeDef.Fields[1]
                if ageField.Name != "age" {
                    t.Errorf("Expected second field to be 'age', got %s", ageField.Name)
                }
                if ageField.IsPrimaryKey {
                    t.Errorf("Expected 'age' to not be primary key")
                }
            },
        },
        {
            name: "PK composite",
            input: `
type Person(#prenom: string, #nom: string, age: number)
`,
            wantErr: false,
            checkFn: func(t *testing.T, prog *Program) {
                typeDef := prog.TypeDefinitions[0]
                pkFields := typeDef.GetPrimaryKeyFields()
                if len(pkFields) != 2 {
                    t.Errorf("Expected 2 PK fields, got %d", len(pkFields))
                }
                if pkFields[0].Name != "prenom" {
                    t.Errorf("Expected first PK to be 'prenom', got %s", pkFields[0].Name)
                }
                if pkFields[1].Name != "nom" {
                    t.Errorf("Expected second PK to be 'nom', got %s", pkFields[1].Name)
                }
            },
        },
        {
            name: "Type sans PK",
            input: `
type Event(timestamp: number, message: string)
`,
            wantErr: false,
            checkFn: func(t *testing.T, prog *Program) {
                typeDef := prog.TypeDefinitions[0]
                if typeDef.HasPrimaryKey() {
                    t.Errorf("Expected no primary key")
                }
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            prog, err := Parse(tt.input)
            if (err != nil) != tt.wantErr {
                t.Fatalf("Parse() error = %v, wantErr %v", err, tt.wantErr)
            }
            if err == nil && tt.checkFn != nil {
                tt.checkFn(t, prog)
            }
        })
    }
}
```

### 6.4. Tests de conversion de faits avec génération d'IDs

**Fichier** : `constraint/constraint_facts_test.go`

**Action** : Ajouter des tests vérifiant que `ConvertFactsToReteFormat` génère correctement les IDs.

```go
func TestConvertFactsToReteFormat_IDGeneration(t *testing.T) {
    tests := []struct {
        name        string
        types       []TypeDefinition
        facts       []Fact
        wantErr     bool
        checkIDsFn  func(*testing.T, []*rete.Fact)
    }{
        {
            name: "Génération ID avec PK simple",
            types: []TypeDefinition{
                {
                    Name: "Person",
                    Fields: []Field{
                        {Name: "nom", Type: "string", IsPrimaryKey: true},
                        {Name: "age", Type: "number"},
                    },
                },
            },
            facts: []Fact{
                {
                    Type: "Person",
                    Fields: []FactField{
                        {Name: "nom", Value: FactValue{StringValue: ptr("Alice")}},
                        {Name: "age", Value: FactValue{NumberValue: ptr(30.0)}},
                    },
                },
            },
            wantErr: false,
            checkIDsFn: func(t *testing.T, facts []*rete.Fact) {
                if len(facts) != 1 {
                    t.Fatalf("Expected 1 fact, got %d", len(facts))
                }
                expectedID := "Person~Alice"
                if facts[0].ID != expectedID {
                    t.Errorf("Expected ID %s, got %s", expectedID, facts[0].ID)
                }
            },
        },
        {
            name: "Génération ID avec PK composite",
            types: []TypeDefinition{
                {
                    Name: "Person",
                    Fields: []Field{
                        {Name: "prenom", Type: "string", IsPrimaryKey: true},
                        {Name: "nom", Type: "string", IsPrimaryKey: true},
                        {Name: "age", Type: "number"},
                    },
                },
            },
            facts: []Fact{
                {
                    Type: "Person",
                    Fields: []FactField{
                        {Name: "prenom", Value: FactValue{StringValue: ptr("Alice")}},
                        {Name: "nom", Value: FactValue{StringValue: ptr("Dupont")}},
                        {Name: "age", Value: FactValue{NumberValue: ptr(30.0)}},
                    },
                },
            },
            wantErr: false,
            checkIDsFn: func(t *testing.T, facts []*rete.Fact) {
                expectedID := "Person~Alice_Dupont"
                if facts[0].ID != expectedID {
                    t.Errorf("Expected ID %s, got %s", expectedID, facts[0].ID)
                }
            },
        },
        {
            name: "Génération ID avec hash (pas de PK)",
            types: []TypeDefinition{
                {
                    Name: "Event",
                    Fields: []Field{
                        {Name: "timestamp", Type: "number"},
                        {Name: "message", Type: "string"},
                    },
                },
            },
            facts: []Fact{
                {
                    Type: "Event",
                    Fields: []FactField{
                        {Name: "timestamp", Value: FactValue{NumberValue: ptr(1234567890.0)}},
                        {Name: "message", Value: FactValue{StringValue: ptr("test")}},
                    },
                },
            },
            wantErr: false,
            checkIDsFn: func(t *testing.T, facts []*rete.Fact) {
                // Vérifier le format hash
                if !strings.HasPrefix(facts[0].ID, "Event~") {
                    t.Errorf("Expected ID to start with Event~, got %s", facts[0].ID)
                }
                // Hash doit être de 16 caractères hex
                hashPart := strings.TrimPrefix(facts[0].ID, "Event~")
                if len(hashPart) != 16 {
                    t.Errorf("Expected hash of 16 chars, got %d: %s", len(hashPart), hashPart)
                }
                if !isHex(hashPart) {
                    t.Errorf("Expected hex hash, got %s", hashPart)
                }
            },
        },
        {
            name: "Échappement des caractères spéciaux dans PK",
            types: []TypeDefinition{
                {
                    Name: "Resource",
                    Fields: []Field{
                        {Name: "path", Type: "string", IsPrimaryKey: true},
                    },
                },
            },
            facts: []Fact{
                {
                    Type: "Resource",
                    Fields: []FactField{
                        {Name: "path", Value: FactValue{StringValue: ptr("/home/user~test_file")}},
                    },
                },
            },
            wantErr: false,
            checkIDsFn: func(t *testing.T, facts []*rete.Fact) {
                // Le ~ et _ doivent être échappés
                expectedID := "Resource~%2Fhome%2Fuser%7Etest%5Ffile"
                if facts[0].ID != expectedID {
                    t.Errorf("Expected ID %s, got %s", expectedID, facts[0].ID)
                }
            },
        },
        {
            name: "PK avec types numériques et booléens",
            types: []TypeDefinition{
                {
                    Name: "Config",
                    Fields: []Field{
                        {Name: "id", Type: "number", IsPrimaryKey: true},
                        {Name: "enabled", Type: "bool", IsPrimaryKey: true},
                        {Name: "value", Type: "string"},
                    },
                },
            },
            facts: []Fact{
                {
                    Type: "Config",
                    Fields: []FactField{
                        {Name: "id", Value: FactValue{NumberValue: ptr(42.0)}},
                        {Name: "enabled", Value: FactValue{BoolValue: ptr(true)}},
                        {Name: "value", Value: FactValue{StringValue: ptr("test")}},
                    },
                },
            },
            wantErr: false,
            checkIDsFn: func(t *testing.T, facts []*rete.Fact) {
                expectedID := "Config~42_true"
                if facts[0].ID != expectedID {
                    t.Errorf("Expected ID %s, got %s", expectedID, facts[0].ID)
                }
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Créer un programme avec les types
            prog := &Program{
                TypeDefinitions: tt.types,
                Facts:          tt.facts,
            }

            reteFacts, err := ConvertFactsToReteFormat(prog)
            if (err != nil) != tt.wantErr {
                t.Fatalf("ConvertFactsToReteFormat() error = %v, wantErr %v", err, tt.wantErr)
            }

            if err == nil && tt.checkIDsFn != nil {
                tt.checkIDsFn(t, reteFacts)
            }
        })
    }
}

// Helper function
func isHex(s string) bool {
    for _, c := range s {
        if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
            return false
        }
    }
    return true
}

func ptr[T any](v T) *T {
    return &v
}
```

### 6.5. Tests de validation (interdiction de `id` explicite)

**Fichier** : `constraint/primary_key_validation_test.go` (créé au prompt 03)

**Action** : Vérifier que les tests de validation existent et couvrent bien :
- Interdiction de `id` dans les assertions
- Validation des types de PK (primitifs uniquement)
- Présence des champs PK dans les faits

Si ces tests manquent, les ajouter maintenant.

```go
func TestValidateFacts_ExplicitIDForbidden(t *testing.T) {
    types := []TypeDefinition{
        {
            Name: "Person",
            Fields: []Field{
                {Name: "nom", Type: "string", IsPrimaryKey: true},
            },
        },
    }

    facts := []Fact{
        {
            Type: "Person",
            Fields: []FactField{
                {Name: "id", Value: FactValue{StringValue: ptr("custom_id")}},
                {Name: "nom", Value: FactValue{StringValue: ptr("Alice")}},
            },
        },
    }

    prog := &Program{
        TypeDefinitions: types,
        Facts:          facts,
    }

    err := ValidateFacts(prog)
    if err == nil {
        t.Errorf("Expected error when fact has explicit 'id' field, got nil")
    }
    if err != nil && !strings.Contains(err.Error(), "id") {
        t.Errorf("Expected error message to mention 'id', got: %v", err)
    }
}
```

### 6.6. Tests de déterminisme des IDs

**Fichier** : `constraint/constraint_facts_test.go`

**Action** : Ajouter un test vérifiant que les mêmes données produisent toujours le même ID.

```go
func TestIDGeneration_Determinism(t *testing.T) {
    types := []TypeDefinition{
        {
            Name: "Event",
            Fields: []Field{
                {Name: "timestamp", Type: "number"},
                {Name: "message", Type: "string"},
            },
        },
    }

    facts := []Fact{
        {
            Type: "Event",
            Fields: []FactField{
                {Name: "timestamp", Value: FactValue{NumberValue: ptr(1234567890.0)}},
                {Name: "message", Value: FactValue{StringValue: ptr("test message")}},
            },
        },
    }

    prog := &Program{
        TypeDefinitions: types,
        Facts:          facts,
    }

    // Générer l'ID plusieurs fois
    var ids []string
    for i := 0; i < 10; i++ {
        reteFacts, err := ConvertFactsToReteFormat(prog)
        if err != nil {
            t.Fatalf("ConvertFactsToReteFormat() error = %v", err)
        }
        ids = append(ids, reteFacts[0].ID)
    }

    // Tous les IDs doivent être identiques
    firstID := ids[0]
    for i, id := range ids {
        if id != firstID {
            t.Errorf("ID mismatch at iteration %d: got %s, want %s", i, id, firstID)
        }
    }
}
```

### 6.7. Tests de cas limites

**Fichier** : `constraint/constraint_facts_test.go`

**Action** : Ajouter des tests pour les cas limites.

```go
func TestIDGeneration_EdgeCases(t *testing.T) {
    tests := []struct {
        name    string
        types   []TypeDefinition
        facts   []Fact
        wantErr bool
        checkFn func(*testing.T, []*rete.Fact)
    }{
        {
            name: "PK avec chaîne vide",
            types: []TypeDefinition{
                {
                    Name: "Person",
                    Fields: []Field{
                        {Name: "nom", Type: "string", IsPrimaryKey: true},
                    },
                },
            },
            facts: []Fact{
                {
                    Type: "Person",
                    Fields: []FactField{
                        {Name: "nom", Value: FactValue{StringValue: ptr("")}},
                    },
                },
            },
            wantErr: false,
            checkFn: func(t *testing.T, facts []*rete.Fact) {
                expectedID := "Person~"
                if facts[0].ID != expectedID {
                    t.Errorf("Expected ID %s, got %s", expectedID, facts[0].ID)
                }
            },
        },
        {
            name: "PK avec zéro",
            types: []TypeDefinition{
                {
                    Name: "Item",
                    Fields: []Field{
                        {Name: "id", Type: "number", IsPrimaryKey: true},
                    },
                },
            },
            facts: []Fact{
                {
                    Type: "Item",
                    Fields: []FactField{
                        {Name: "id", Value: FactValue{NumberValue: ptr(0.0)}},
                    },
                },
            },
            wantErr: false,
            checkFn: func(t *testing.T, facts []*rete.Fact) {
                expectedID := "Item~0"
                if facts[0].ID != expectedID {
                    t.Errorf("Expected ID %s, got %s", expectedID, facts[0].ID)
                }
            },
        },
        {
            name: "PK avec false",
            types: []TypeDefinition{
                {
                    Name: "Flag",
                    Fields: []Field{
                        {Name: "enabled", Type: "bool", IsPrimaryKey: true},
                    },
                },
            },
            facts: []Fact{
                {
                    Type: "Flag",
                    Fields: []FactField{
                        {Name: "enabled", Value: FactValue{BoolValue: ptr(false)}},
                    },
                },
            },
            wantErr: false,
            checkFn: func(t *testing.T, facts []*rete.Fact) {
                expectedID := "Flag~false"
                if facts[0].ID != expectedID {
                    t.Errorf("Expected ID %s, got %s", expectedID, facts[0].ID)
                }
            },
        },
        {
            name: "Float avec précision",
            types: []TypeDefinition{
                {
                    Name: "Measurement",
                    Fields: []Field{
                        {Name: "value", Type: "number", IsPrimaryKey: true},
                    },
                },
            },
            facts: []Fact{
                {
                    Type: "Measurement",
                    Fields: []FactField{
                        {Name: "value", Value: FactValue{NumberValue: ptr(3.14159)}},
                    },
                },
            },
            wantErr: false,
            checkFn: func(t *testing.T, facts []*rete.Fact) {
                // Vérifier que le float est formaté de manière déterministe
                // L'ID doit être reproductible
                id := facts[0].ID
                if !strings.HasPrefix(id, "Measurement~3.14159") {
                    t.Errorf("Expected ID to start with Measurement~3.14159, got %s", id)
                }
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            prog := &Program{
                TypeDefinitions: tt.types,
                Facts:          tt.facts,
            }

            reteFacts, err := ConvertFactsToReteFormat(prog)
            if (err != nil) != tt.wantErr {
                t.Fatalf("ConvertFactsToReteFormat() error = %v, wantErr %v", err, tt.wantErr)
            }

            if err == nil && tt.checkFn != nil {
                tt.checkFn(t, reteFacts)
            }
        })
    }
}
```

---

## Validation

### Étape 1 : Vérifier tous les tests du module constraint

```bash
cd /home/resinsec/dev/tsd
go test ./constraint/... -v
```

Tous les tests doivent passer.

### Étape 2 : Vérifier la couverture de tests

```bash
go test ./constraint/... -cover
```

Viser au minimum 80% de couverture pour les fichiers modifiés.

### Étape 3 : Tests spécifiques aux nouvelles fonctionnalités

```bash
go test -run TestParser_PrimaryKeyFields ./constraint/... -v
go test -run TestConvertFactsToReteFormat_IDGeneration ./constraint/... -v
go test -run TestIDGeneration_Determinism ./constraint/... -v
go test -run TestIDGeneration_EdgeCases ./constraint/... -v
```

### Étape 4 : Validation globale

```bash
make validate
```

### Étape 5 : Test de régression

Exécuter tous les tests du projet pour s'assurer qu'aucun autre module n'est cassé :

```bash
make test-complete
```

ou

```bash
go test ./... -v
```

---

## Checklist

- [ ] Inventaire des tests existants effectué
- [ ] Tests cassés identifiés et corrigés
- [ ] Tests de parsing avec `#` ajoutés
- [ ] Tests de génération d'ID avec PK simple ajoutés
- [ ] Tests de génération d'ID avec PK composite ajoutés
- [ ] Tests de génération d'ID avec hash ajoutés
- [ ] Tests d'échappement de caractères ajoutés
- [ ] Tests avec types numériques et booléens ajoutés
- [ ] Tests de validation (interdiction `id`) vérifiés/ajoutés
- [ ] Tests de déterminisme ajoutés
- [ ] Tests de cas limites ajoutés
- [ ] `go test ./constraint/... -v` réussit
- [ ] Couverture de tests vérifiée (≥80%)
- [ ] `make validate` réussit
- [ ] `make test-complete` réussit (ou `go test ./...`)

---

## Rapport

Une fois toutes les tâches complétées :

1. Lister les tests existants modifiés
2. Lister les nouveaux tests ajoutés
3. Indiquer le pourcentage de couverture atteint
4. Copier la sortie des tests réussis
5. Commit :

```bash
git add constraint/
git commit -m "test(constraint): tests complets pour clés primaires et génération d'IDs

- Adaptation des tests existants (suppression des IDs explicites)
- Tests de parsing des champs marqués avec #
- Tests de génération d'IDs (PK simple, composite, hash)
- Tests d'échappement de caractères spéciaux
- Tests de validation (interdiction de id explicite)
- Tests de déterminisme
- Tests de cas limites (valeurs nulles, zéro, false, floats)

Couverture: XX%

Refs #<issue_number>"
```

---

## Dépendances

- **Bloque** : Prompt 08 (tests e2e)
- **Bloqué par** : Prompts 01-05

---

## Notes importantes

1. **Déterminisme** : Tous les tests doivent être reproductibles. Les IDs générés pour les mêmes données doivent toujours être identiques.

2. **Échappement** : Bien tester l'échappement des caractères `~`, `_`, et `%` dans les valeurs de PK.

3. **Floats** : Attention à la représentation des floats. Utiliser un format déterministe (pas de notation scientifique aléatoire).

4. **Bug numérique** : Si des tests échouent à cause du bug connu de l'évaluateur numérique, les documenter mais ne pas les retirer. Le bug sera fixé séparément.

5. **Nettoyage** : Retirer tous les `TODO` et commentaires temporaires liés aux anciens IDs (`parsed_fact_n`).

---

**Prêt à passer au prompt 07 après validation complète.**