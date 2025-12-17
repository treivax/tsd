# Prompt 05 : Intégration RETE et gestion des IDs

**Objectif** : Vérifier et adapter l'intégration RETE pour que les IDs générés soient correctement utilisés dans le moteur de règles (évaluateur, tokens, working memory).

**Prérequis** : Prompts 01-04 complétés et validés.

---

## Contexte

Les IDs de faits sont maintenant générés automatiquement selon les règles de clé primaire. Le moteur RETE doit :

1. Accepter et stocker ces IDs dans les faits
2. Permettre de référencer `id` dans les expressions et règles
3. Gérer correctement les IDs dans les tokens et bindings
4. Utiliser les IDs pour l'indexation dans la working memory

Le module `rete` utilise déjà une structure `Fact` avec un champ `ID string` et des helpers `GetInternalID`, `MakeInternalID`, `ParseInternalID` dans `rete/fact_token.go`.

---

## Tâches

### 5.1. Vérifier la structure Fact dans RETE

**Fichier** : `rete/fact_token.go`

**Action** : Vérifier que la structure `Fact` contient bien :

```go
type Fact struct {
    ID     string                 // ID du fait (généré ou fourni)
    Type   string                 // Type du fait
    Fields map[string]interface{} // Champs du fait
}
```

**Validation** :
- Le champ `ID` doit être de type `string`
- Les fonctions `GetInternalID()`, `MakeInternalID()`, `ParseInternalID()` doivent exister

Si nécessaire, ajouter des commentaires clarificateurs :

```go
// ID is the unique identifier for this fact.
// It is either generated from primary keys or computed as a hash.
// Format: "TypeName~value1_value2..." or "TypeName~<hash>"
ID string
```

### 5.2. Vérifier l'accès au champ `id` dans l'évaluateur

**Fichier** : `rete/evaluator.go`

**Action** : S'assurer que l'évaluateur peut accéder au champ spécial `id` d'un fait.

Rechercher la fonction qui évalue les accès aux champs (probablement `evaluateFieldAccess` ou similaire).

**Vérification** : Quand on accède à `fact.id` dans une expression, l'évaluateur doit retourner `fact.ID` (le champ structurel), pas chercher dans `fact.Fields["id"]`.

**Code attendu** (exemple) :

```go
func (e *Evaluator) getFieldValue(fact *Fact, fieldName string) (interface{}, error) {
    // Cas spécial pour le champ id
    if fieldName == "id" {
        return fact.ID, nil
    }
    
    // Accès normal aux autres champs
    value, ok := fact.Fields[fieldName]
    if !ok {
        return nil, fmt.Errorf("field %s not found in fact %s", fieldName, fact.Type)
    }
    return value, nil
}
```

**Si cette logique n'existe pas** : l'ajouter dans la fonction appropriée de l'évaluateur.

### 5.3. Vérifier la constante FieldNameID

**Fichier** : `rete/fact_token.go` (ou `rete/constants.go` si existe)

**Action** : Vérifier qu'il existe une constante pour le nom du champ ID :

```go
const FieldNameID = "id"
```

Si elle n'existe pas, l'ajouter dans `rete/fact_token.go` en haut du fichier.

### 5.4. Tests d'intégration RETE avec IDs

**Fichier** : `rete/fact_token_test.go` (ou créer si nécessaire)

**Action** : Ajouter des tests pour vérifier la manipulation des IDs.

```go
// Copyright 2025 SekiaTech. All rights reserved.
// Use of this source code is governed by an MIT-style license.

package rete

import (
    "testing"
)

func TestFact_IDHandling(t *testing.T) {
    tests := []struct {
        name     string
        fact     *Fact
        wantID   string
        wantType string
    }{
        {
            name: "fact avec PK simple",
            fact: &Fact{
                ID:   "Person~Alice",
                Type: "Person",
                Fields: map[string]interface{}{
                    "nom": "Alice",
                    "age": 30,
                },
            },
            wantID:   "Person~Alice",
            wantType: "Person",
        },
        {
            name: "fact avec PK composite",
            fact: &Fact{
                ID:   "Person~Alice_Dupont",
                Type: "Person",
                Fields: map[string]interface{}{
                    "prenom": "Alice",
                    "nom":    "Dupont",
                    "age":    30,
                },
            },
            wantID:   "Person~Alice_Dupont",
            wantType: "Person",
        },
        {
            name: "fact avec hash",
            fact: &Fact{
                ID:   "Event~a1b2c3d4e5f6g7h8",
                Type: "Event",
                Fields: map[string]interface{}{
                    "timestamp": 1234567890,
                    "message":   "test",
                },
            },
            wantID:   "Event~a1b2c3d4e5f6g7h8",
            wantType: "Event",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if tt.fact.ID != tt.wantID {
                t.Errorf("ID = %v, want %v", tt.fact.ID, tt.wantID)
            }
            if tt.fact.Type != tt.wantType {
                t.Errorf("Type = %v, want %v", tt.fact.Type, tt.wantType)
            }
        })
    }
}

func TestGetInternalID(t *testing.T) {
    fact := &Fact{
        ID:   "Person~Alice",
        Type: "Person",
        Fields: map[string]interface{}{
            "nom": "Alice",
        },
    }

    internalID := fact.GetInternalID()
    expected := "Person_Person~Alice" // Format: Type_ID

    if internalID != expected {
        t.Errorf("GetInternalID() = %v, want %v", internalID, expected)
    }
}

func TestMakeInternalID(t *testing.T) {
    tests := []struct {
        factType string
        factID   string
        want     string
    }{
        {"Person", "Person~Alice", "Person_Person~Alice"},
        {"Event", "Event~hash123", "Event_Event~hash123"},
        {"Person", "Person~Alice_Dupont", "Person_Person~Alice_Dupont"},
    }

    for _, tt := range tests {
        t.Run(tt.factType+"_"+tt.factID, func(t *testing.T) {
            got := MakeInternalID(tt.factType, tt.factID)
            if got != tt.want {
                t.Errorf("MakeInternalID(%q, %q) = %v, want %v",
                    tt.factType, tt.factID, got, tt.want)
            }
        })
    }
}

func TestParseInternalID(t *testing.T) {
    tests := []struct {
        internalID string
        wantType   string
        wantID     string
        wantErr    bool
    }{
        {
            internalID: "Person_Person~Alice",
            wantType:   "Person",
            wantID:     "Person~Alice",
            wantErr:    false,
        },
        {
            internalID: "Event_Event~hash123",
            wantType:   "Event",
            wantID:     "Event~hash123",
            wantErr:    false,
        },
        {
            internalID: "Person_Person~Alice_Dupont",
            wantType:   "Person",
            wantID:     "Person~Alice_Dupont",
            wantErr:    false,
        },
        {
            internalID: "InvalidFormat",
            wantType:   "",
            wantID:     "",
            wantErr:    true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.internalID, func(t *testing.T) {
            gotType, gotID, err := ParseInternalID(tt.internalID)
            if (err != nil) != tt.wantErr {
                t.Errorf("ParseInternalID() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if gotType != tt.wantType {
                t.Errorf("ParseInternalID() type = %v, want %v", gotType, tt.wantType)
            }
            if gotID != tt.wantID {
                t.Errorf("ParseInternalID() ID = %v, want %v", gotID, tt.wantID)
            }
        })
    }
}
```

### 5.5. Test de l'accès au champ `id` dans les expressions

**Fichier** : `rete/evaluator_test.go`

**Action** : Ajouter un test vérifiant qu'on peut accéder au champ `id` dans une expression.

```go
func TestEvaluator_AccessIDField(t *testing.T) {
    fact := &Fact{
        ID:   "Person~Alice",
        Type: "Person",
        Fields: map[string]interface{}{
            "nom": "Alice",
            "age": 30,
        },
    }

    // Créer un évaluateur (adapter selon l'API réelle)
    eval := NewEvaluator()

    // Tester l'accès au champ id
    idValue, err := eval.getFieldValue(fact, "id")
    if err != nil {
        t.Fatalf("getFieldValue(fact, 'id') error = %v", err)
    }

    if idValue != "Person~Alice" {
        t.Errorf("getFieldValue(fact, 'id') = %v, want %v", idValue, "Person~Alice")
    }

    // Tester l'accès à un champ normal
    nomValue, err := eval.getFieldValue(fact, "nom")
    if err != nil {
        t.Fatalf("getFieldValue(fact, 'nom') error = %v", err)
    }

    if nomValue != "Alice" {
        t.Errorf("getFieldValue(fact, 'nom') = %v, want %v", nomValue, "Alice")
    }
}
```

**Note** : Adapter ce test selon l'API réelle de l'évaluateur. Si la fonction `getFieldValue` n'est pas exportée, créer un test d'intégration qui évalue une expression complète comme `p.id == "Person~Alice"`.

### 5.6. Vérification de la working memory

**Fichier** : `rete/working_memory.go` (lire)

**Action** : Vérifier que la working memory utilise bien `GetInternalID()` ou `MakeInternalID()` pour indexer les faits.

**Rechercher** : Les appels à `wm.facts[...]` ou similaire pour s'assurer que la clé est l'internal ID.

**Exemple attendu** :

```go
func (wm *WorkingMemory) AddFact(fact *Fact) {
    internalID := fact.GetInternalID()
    wm.facts[internalID] = fact
    // ...
}
```

**Si ce n'est pas le cas** : documenter le problème dans un commentaire TODO et signaler dans le rapport de validation (pas besoin de fixer dans ce prompt, juste documenter).

---

## Validation

### Étape 1 : Compilation

```bash
cd /home/resinsec/dev/tsd
go build ./rete/...
```

Aucune erreur de compilation attendue.

### Étape 2 : Exécuter les tests RETE

```bash
go test ./rete/... -v
```

Tous les tests doivent passer, y compris les nouveaux tests sur les IDs.

### Étape 3 : Vérifier l'accès au champ `id`

Si vous avez créé un test d'expression complète, exécutez-le :

```bash
go test -run TestEvaluator_AccessIDField ./rete/... -v
```

### Étape 4 : Validation globale

```bash
make validate
```

### Étape 5 : Test manuel (optionnel)

Créer un fichier `.tsd` simple pour tester l'accès à `id` dans une règle :

**Fichier** : `test_id_access.tsd`

```
type Person(#nom: string, age: number)

assert Person(nom: "Alice", age: 30)

rule PersonWithID {
    when {
        p: Person()
        p.id == "Person~Alice"
    }
    then {
        print("Found person with ID: " + p.id)
    }
}
```

Compiler et exécuter (si le CLI est disponible) :

```bash
./tsd run test_id_access.tsd
```

Vérifier que la règle se déclenche et affiche l'ID correct.

---

## Checklist

- [ ] Structure `Fact` vérifiée dans `rete/fact_token.go`
- [ ] Commentaires ajoutés sur le champ `ID`
- [ ] Constante `FieldNameID` présente
- [ ] Évaluateur permet l'accès au champ `id` (code ajouté si nécessaire)
- [ ] Tests `TestFact_IDHandling` ajoutés et passent
- [ ] Tests `TestGetInternalID` ajoutés et passent
- [ ] Tests `TestMakeInternalID` ajoutés et passent
- [ ] Tests `TestParseInternalID` ajoutés et passent
- [ ] Test `TestEvaluator_AccessIDField` ajouté et passe
- [ ] Working memory vérifiée (utilise bien les internal IDs)
- [ ] `go build ./rete/...` réussit
- [ ] `go test ./rete/... -v` réussit
- [ ] `make validate` réussit
- [ ] Test manuel avec fichier `.tsd` (optionnel) réussit

---

## Rapport

Une fois toutes les tâches complétées :

1. Documenter tout problème découvert (par ex. si la working memory n'utilise pas les internal IDs correctement)
2. Lister les modifications apportées
3. Copier la sortie des tests réussis
4. Commit :

```bash
git add rete/
git commit -m "feat(rete): intégration des IDs générés automatiquement

- Ajout de commentaires sur Fact.ID
- Ajout de la constante FieldNameID
- Support de l'accès au champ 'id' dans l'évaluateur
- Tests pour GetInternalID, MakeInternalID, ParseInternalID
- Tests pour l'accès au champ id dans les expressions
- Vérification de l'indexation dans la working memory

Refs #<issue_number>"
```

---

## Dépendances

- **Bloque** : Prompt 06 (tests constraint), Prompt 08 (tests e2e)
- **Bloqué par** : Prompts 01-04

---

## Notes importantes

1. **Bug évaluateur numérique** : Il existe un bug connu dans l'évaluateur pour les comparaisons numériques (égalité entre int et float). Ce bug doit être fixé séparément avant d'utiliser des PK numériques dans les joins. Documenter ce risque si vous utilisez des nombres dans vos tests.

2. **Immutabilité des tokens** : Les tokens RETE utilisent des BindingChains immutables. L'accès à `id` doit fonctionner via les bindings normaux, pas de traitement spécial nécessaire au niveau des tokens.

3. **Performance** : L'accès au champ `id` doit être aussi performant qu'un accès normal (pas de réflexion, juste un if sur le nom du champ).

---

**Prêt à passer au prompt 06 après validation complète.**