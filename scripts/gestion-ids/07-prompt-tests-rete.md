# Prompt 07 : Tests du module rete

**Objectif** : Compléter la couverture de tests du module `rete` pour vérifier que les nouveaux IDs générés fonctionnent correctement dans le moteur de règles (joins, comparaisons, working memory).

**Prérequis** : Prompts 01-06 complétés et validés.

---

## Contexte

Les IDs de faits sont maintenant générés automatiquement et peuvent être :
- Basés sur des clés primaires : `TypeName~value1_value2_...`
- Basés sur un hash : `TypeName~<hash>`

Le moteur RETE doit :
1. Stocker et indexer correctement les faits avec ces nouveaux IDs
2. Permettre de comparer des IDs dans les conditions de règles
3. Gérer les joins sur des faits avec des IDs générés
4. Permettre l'accès au champ `id` dans les expressions

---

## Tâches

### 7.1. Inventaire des tests existants

**Fichiers à examiner** :
- `rete/rete_test.go`
- `rete/working_memory_test.go`
- `rete/evaluator_test.go`
- `rete/join_test.go`
- Tout autre fichier `*_test.go` dans `rete/`

**Action** : Lire chaque fichier de test et identifier :
1. Les tests qui créent des faits avec des IDs manuels
2. Les tests qui pourraient être impactés par le nouveau format d'ID
3. Les tests de joins et comparaisons qui utilisent des IDs

**Commande** :
```bash
cd /home/resinsec/dev/tsd
grep -n "ID:" rete/*_test.go
grep -n "parsed_fact" rete/*_test.go
```

### 7.2. Adapter les tests existants

**Pour chaque test qui crée des faits avec des IDs manuels** :

**Avant** :
```go
fact := &Fact{
    ID:   "person_1",
    Type: "Person",
    Fields: map[string]interface{}{
        "nom": "Alice",
    },
}
```

**Après** : Utiliser le nouveau format d'ID.
```go
fact := &Fact{
    ID:   "Person~Alice", // Format avec PK
    Type: "Person",
    Fields: map[string]interface{}{
        "nom": "Alice",
    },
}
```

Ou bien :
```go
fact := &Fact{
    ID:   "Person~a1b2c3d4e5f6g7h8", // Format avec hash
    Type: "Person",
    Fields: map[string]interface{}{
        "nom": "Alice",
    },
}
```

### 7.3. Tests de la working memory avec nouveaux IDs

**Fichier** : `rete/working_memory_test.go`

**Action** : Ajouter des tests vérifiant que la working memory gère correctement les nouveaux formats d'IDs.

```go
func TestWorkingMemory_NewIDFormats(t *testing.T) {
    tests := []struct {
        name      string
        facts     []*Fact
        checkFn   func(*testing.T, *WorkingMemory)
    }{
        {
            name: "Ajout de fait avec ID basé sur PK simple",
            facts: []*Fact{
                {
                    ID:   "Person~Alice",
                    Type: "Person",
                    Fields: map[string]interface{}{
                        "nom": "Alice",
                        "age": 30,
                    },
                },
            },
            checkFn: func(t *testing.T, wm *WorkingMemory) {
                // Vérifier que le fait est indexé correctement
                internalID := "Person_Person~Alice"
                fact := wm.GetFact(internalID)
                if fact == nil {
                    t.Errorf("Fact not found with internal ID %s", internalID)
                }
                if fact.ID != "Person~Alice" {
                    t.Errorf("Expected fact ID Person~Alice, got %s", fact.ID)
                }
            },
        },
        {
            name: "Ajout de fait avec ID basé sur PK composite",
            facts: []*Fact{
                {
                    ID:   "Person~Alice_Dupont",
                    Type: "Person",
                    Fields: map[string]interface{}{
                        "prenom": "Alice",
                        "nom":    "Dupont",
                        "age":    30,
                    },
                },
            },
            checkFn: func(t *testing.T, wm *WorkingMemory) {
                internalID := "Person_Person~Alice_Dupont"
                fact := wm.GetFact(internalID)
                if fact == nil {
                    t.Errorf("Fact not found with internal ID %s", internalID)
                }
            },
        },
        {
            name: "Ajout de fait avec ID basé sur hash",
            facts: []*Fact{
                {
                    ID:   "Event~a1b2c3d4e5f6g7h8",
                    Type: "Event",
                    Fields: map[string]interface{}{
                        "timestamp": 1234567890,
                        "message":   "test",
                    },
                },
            },
            checkFn: func(t *testing.T, wm *WorkingMemory) {
                internalID := "Event_Event~a1b2c3d4e5f6g7h8"
                fact := wm.GetFact(internalID)
                if fact == nil {
                    t.Errorf("Fact not found with internal ID %s", internalID)
                }
            },
        },
        {
            name: "Retrait de fait avec nouveau format d'ID",
            facts: []*Fact{
                {
                    ID:   "Person~Alice",
                    Type: "Person",
                    Fields: map[string]interface{}{
                        "nom": "Alice",
                    },
                },
            },
            checkFn: func(t *testing.T, wm *WorkingMemory) {
                // Ajouter puis retirer
                fact := &Fact{
                    ID:   "Person~Alice",
                    Type: "Person",
                    Fields: map[string]interface{}{
                        "nom": "Alice",
                    },
                }
                wm.AddFact(fact)
                
                internalID := "Person_Person~Alice"
                if wm.GetFact(internalID) == nil {
                    t.Errorf("Fact should be present before retraction")
                }
                
                wm.RemoveFact(fact)
                
                if wm.GetFact(internalID) != nil {
                    t.Errorf("Fact should be removed after retraction")
                }
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            wm := NewWorkingMemory()
            
            for _, fact := range tt.facts {
                wm.AddFact(fact)
            }
            
            if tt.checkFn != nil {
                tt.checkFn(t, wm)
            }
        })
    }
}
```

### 7.4. Tests de l'évaluateur avec accès au champ `id`

**Fichier** : `rete/evaluator_test.go`

**Action** : Ajouter des tests vérifiant l'accès au champ `id` dans les expressions.

```go
func TestEvaluator_IDFieldAccess(t *testing.T) {
    tests := []struct {
        name       string
        fact       *Fact
        expression string
        wantResult bool
        wantErr    bool
    }{
        {
            name: "Comparaison d'ID avec PK simple",
            fact: &Fact{
                ID:   "Person~Alice",
                Type: "Person",
                Fields: map[string]interface{}{
                    "nom": "Alice",
                    "age": 30,
                },
            },
            expression: "p.id == \"Person~Alice\"",
            wantResult: true,
            wantErr:    false,
        },
        {
            name: "Comparaison d'ID avec PK composite",
            fact: &Fact{
                ID:   "Person~Alice_Dupont",
                Type: "Person",
                Fields: map[string]interface{}{
                    "prenom": "Alice",
                    "nom":    "Dupont",
                },
            },
            expression: "p.id == \"Person~Alice_Dupont\"",
            wantResult: true,
            wantErr:    false,
        },
        {
            name: "Comparaison d'ID négative",
            fact: &Fact{
                ID:   "Person~Alice",
                Type: "Person",
                Fields: map[string]interface{}{
                    "nom": "Alice",
                },
            },
            expression: "p.id == \"Person~Bob\"",
            wantResult: false,
            wantErr:    false,
        },
        {
            name: "Accès à ID avec hash",
            fact: &Fact{
                ID:   "Event~a1b2c3d4e5f6g7h8",
                Type: "Event",
                Fields: map[string]interface{}{
                    "message": "test",
                },
            },
            expression: "e.id == \"Event~a1b2c3d4e5f6g7h8\"",
            wantResult: true,
            wantErr:    false,
        },
        {
            name: "Utilisation de contains sur ID",
            fact: &Fact{
                ID:   "Person~Alice_Dupont",
                Type: "Person",
                Fields: map[string]interface{}{
                    "prenom": "Alice",
                    "nom":    "Dupont",
                },
            },
            expression: "contains(p.id, \"Alice\")",
            wantResult: true,
            wantErr:    false,
        },
        {
            name: "Utilisation de startsWith sur ID",
            fact: &Fact{
                ID:   "Person~Alice_Dupont",
                Type: "Person",
                Fields: map[string]interface{}{
                    "prenom": "Alice",
                    "nom":    "Dupont",
                },
            },
            expression: "startsWith(p.id, \"Person~\")",
            wantResult: true,
            wantErr:    false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Créer l'évaluateur (adapter selon l'API réelle)
            eval := NewEvaluator()
            
            // Créer un contexte de bindings avec le fait
            bindings := map[string]*Fact{
                "p": tt.fact,
                "e": tt.fact,
            }
            
            // Évaluer l'expression
            result, err := eval.Evaluate(tt.expression, bindings)
            
            if (err != nil) != tt.wantErr {
                t.Fatalf("Evaluate() error = %v, wantErr %v", err, tt.wantErr)
            }
            
            if err == nil {
                boolResult, ok := result.(bool)
                if !ok {
                    t.Fatalf("Expected bool result, got %T", result)
                }
                if boolResult != tt.wantResult {
                    t.Errorf("Evaluate() = %v, want %v", boolResult, tt.wantResult)
                }
            }
        })
    }
}
```

### 7.5. Tests de joins avec IDs générés

**Fichier** : `rete/join_test.go` (ou créer si nécessaire)

**Action** : Ajouter des tests pour vérifier que les joins fonctionnent avec les nouveaux IDs.

```go
func TestJoin_WithGeneratedIDs(t *testing.T) {
    tests := []struct {
        name      string
        facts     []*Fact
        joinCond  string
        wantMatch bool
    }{
        {
            name: "Join sur ID avec PK simple",
            facts: []*Fact{
                {
                    ID:   "Person~Alice",
                    Type: "Person",
                    Fields: map[string]interface{}{
                        "nom": "Alice",
                    },
                },
                {
                    ID:   "Membership~Alice_Club1",
                    Type: "Membership",
                    Fields: map[string]interface{}{
                        "person_id": "Person~Alice",
                        "club":      "Club1",
                    },
                },
            },
            joinCond:  "p.id == m.person_id",
            wantMatch: true,
        },
        {
            name: "Join sur ID avec PK composite",
            facts: []*Fact{
                {
                    ID:   "Person~Alice_Dupont",
                    Type: "Person",
                    Fields: map[string]interface{}{
                        "prenom": "Alice",
                        "nom":    "Dupont",
                    },
                },
                {
                    ID:   "Contact~Alice_Dupont",
                    Type: "Contact",
                    Fields: map[string]interface{}{
                        "person_id": "Person~Alice_Dupont",
                        "email":     "alice@example.com",
                    },
                },
            },
            joinCond:  "p.id == c.person_id",
            wantMatch: true,
        },
        {
            name: "Join échoue avec IDs différents",
            facts: []*Fact{
                {
                    ID:   "Person~Alice",
                    Type: "Person",
                    Fields: map[string]interface{}{
                        "nom": "Alice",
                    },
                },
                {
                    ID:   "Membership~Bob_Club1",
                    Type: "Membership",
                    Fields: map[string]interface{}{
                        "person_id": "Person~Bob",
                        "club":      "Club1",
                    },
                },
            },
            joinCond:  "p.id == m.person_id",
            wantMatch: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Créer le réseau RETE et ajouter les faits
            network := NewReteNetwork()
            wm := NewWorkingMemory()
            
            for _, fact := range tt.facts {
                wm.AddFact(fact)
            }
            
            // Créer un join node avec la condition
            // (adapter selon l'API réelle)
            joinNode := network.CreateJoinNode(tt.joinCond)
            
            // Propager les tokens et vérifier les matches
            matches := joinNode.GetMatches(wm)
            
            hasMatch := len(matches) > 0
            if hasMatch != tt.wantMatch {
                t.Errorf("Join match = %v, want %v", hasMatch, tt.wantMatch)
            }
        })
    }
}
```

### 7.6. Tests de comparaisons d'IDs avec types différents

**Fichier** : `rete/evaluator_test.go`

**Action** : Vérifier que les comparaisons d'IDs fonctionnent correctement (rappel : les IDs sont toujours des strings).

```go
func TestEvaluator_IDComparisons(t *testing.T) {
    tests := []struct {
        name       string
        fact1      *Fact
        fact2      *Fact
        expression string
        wantResult bool
    }{
        {
            name: "Égalité d'IDs identiques",
            fact1: &Fact{
                ID:   "Person~Alice",
                Type: "Person",
                Fields: map[string]interface{}{"nom": "Alice"},
            },
            fact2: &Fact{
                ID:   "Person~Alice",
                Type: "Person",
                Fields: map[string]interface{}{"nom": "Alice"},
            },
            expression: "p1.id == p2.id",
            wantResult: true,
        },
        {
            name: "Inégalité d'IDs différents",
            fact1: &Fact{
                ID:   "Person~Alice",
                Type: "Person",
                Fields: map[string]interface{}{"nom": "Alice"},
            },
            fact2: &Fact{
                ID:   "Person~Bob",
                Type: "Person",
                Fields: map[string]interface{}{"nom": "Bob"},
            },
            expression: "p1.id != p2.id",
            wantResult: true,
        },
        {
            name: "Comparaison d'IDs avec hash",
            fact1: &Fact{
                ID:   "Event~a1b2c3d4e5f6g7h8",
                Type: "Event",
                Fields: map[string]interface{}{"message": "test1"},
            },
            fact2: &Fact{
                ID:   "Event~a1b2c3d4e5f6g7h8",
                Type: "Event",
                Fields: map[string]interface{}{"message": "test1"},
            },
            expression: "e1.id == e2.id",
            wantResult: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            eval := NewEvaluator()
            
            bindings := map[string]*Fact{
                "p1": tt.fact1,
                "p2": tt.fact2,
                "e1": tt.fact1,
                "e2": tt.fact2,
            }
            
            result, err := eval.Evaluate(tt.expression, bindings)
            if err != nil {
                t.Fatalf("Evaluate() error = %v", err)
            }
            
            boolResult, ok := result.(bool)
            if !ok {
                t.Fatalf("Expected bool result, got %T", result)
            }
            
            if boolResult != tt.wantResult {
                t.Errorf("Evaluate() = %v, want %v for expression %s",
                    boolResult, tt.wantResult, tt.expression)
            }
        })
    }
}
```

### 7.7. Tests de performance (optionnel)

**Fichier** : `rete/performance_test.go` (ou `rete/benchmark_test.go`)

**Action** : Ajouter des benchmarks pour comparer les performances avec les nouveaux IDs.

```go
func BenchmarkWorkingMemory_AddFactWithPKID(b *testing.B) {
    wm := NewWorkingMemory()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        fact := &Fact{
            ID:   fmt.Sprintf("Person~User%d", i),
            Type: "Person",
            Fields: map[string]interface{}{
                "nom": fmt.Sprintf("User%d", i),
                "age": 30,
            },
        }
        wm.AddFact(fact)
    }
}

func BenchmarkWorkingMemory_AddFactWithHashID(b *testing.B) {
    wm := NewWorkingMemory()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        fact := &Fact{
            ID:   fmt.Sprintf("Event~%016x", i),
            Type: "Event",
            Fields: map[string]interface{}{
                "timestamp": i,
                "message":   fmt.Sprintf("Message %d", i),
            },
        }
        wm.AddFact(fact)
    }
}

func BenchmarkEvaluator_IDFieldAccess(b *testing.B) {
    eval := NewEvaluator()
    fact := &Fact{
        ID:   "Person~Alice",
        Type: "Person",
        Fields: map[string]interface{}{
            "nom": "Alice",
            "age": 30,
        },
    }
    bindings := map[string]*Fact{"p": fact}
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := eval.Evaluate("p.id == \"Person~Alice\"", bindings)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

---

## Validation

### Étape 1 : Vérifier tous les tests du module rete

```bash
cd /home/resinsec/dev/tsd
go test ./rete/... -v
```

Tous les tests doivent passer.

### Étape 2 : Vérifier la couverture de tests

```bash
go test ./rete/... -cover
```

Viser au minimum 75% de couverture pour les fichiers modifiés.

### Étape 3 : Tests spécifiques

```bash
go test -run TestWorkingMemory_NewIDFormats ./rete/... -v
go test -run TestEvaluator_IDFieldAccess ./rete/... -v
go test -run TestJoin_WithGeneratedIDs ./rete/... -v
go test -run TestEvaluator_IDComparisons ./rete/... -v
```

### Étape 4 : Benchmarks (optionnel)

```bash
go test -bench=. ./rete/... -benchmem
```

Vérifier que les performances sont acceptables.

### Étape 5 : Validation globale

```bash
make validate
```

### Étape 6 : Test de régression complet

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
- [ ] Tests de working memory avec nouveaux IDs ajoutés
- [ ] Tests d'ajout/retrait de faits avec nouveaux IDs ajoutés
- [ ] Tests d'accès au champ `id` dans l'évaluateur ajoutés
- [ ] Tests de comparaisons d'IDs ajoutés
- [ ] Tests de joins avec IDs générés ajoutés
- [ ] Tests avec contains/startsWith sur IDs ajoutés
- [ ] Tests de comparaisons d'IDs entre faits ajoutés
- [ ] Benchmarks ajoutés (optionnel)
- [ ] `go test ./rete/... -v` réussit
- [ ] Couverture de tests vérifiée (≥75%)
- [ ] `make validate` réussit
- [ ] `make test-complete` réussit

---

## Rapport

Une fois toutes les tâches complétées :

1. Lister les tests existants modifiés
2. Lister les nouveaux tests ajoutés
3. Indiquer le pourcentage de couverture atteint
4. Indiquer les résultats des benchmarks (si applicable)
5. Copier la sortie des tests réussis
6. Commit :

```bash
git add rete/
git commit -m "test(rete): tests complets pour IDs générés dans le moteur RETE

- Adaptation des tests existants aux nouveaux formats d'IDs
- Tests de working memory avec IDs basés sur PK et hash
- Tests d'accès au champ id dans l'évaluateur
- Tests de comparaisons d'IDs (égalité, inégalité)
- Tests de joins avec IDs générés
- Tests de fonctions string sur IDs (contains, startsWith)
- Benchmarks de performance (optionnel)

Couverture: XX%

Refs #<issue_number>"
```

---

## Dépendances

- **Bloque** : Prompt 08 (tests e2e)
- **Bloqué par** : Prompts 01-06

---

## Notes importantes

1. **Adaptation de l'API** : Les exemples de code ci-dessus supposent une certaine API pour le module `rete`. Adaptez les tests selon l'API réelle du projet (noms de fonctions, signatures, etc.).

2. **Bug évaluateur numérique** : Si le bug de comparaison numérique n'est pas encore fixé, évitez d'utiliser des PK numériques dans les tests de joins complexes, ou documentez les échecs attendus.

3. **Internal IDs** : Rappel que la working memory utilise des internal IDs au format `Type_ID` pour l'indexation. Les tests doivent utiliser `GetInternalID()` ou `MakeInternalID()` pour construire ces clés.

4. **Immutabilité** : Les tokens RETE sont immutables. Les tests de propagation doivent vérifier que les nouveaux tokens contiennent les bons bindings avec les bons IDs.

5. **Échappement** : Tester aussi les IDs contenant des caractères échappés (`~`, `_`, `%`) pour s'assurer que les comparaisons fonctionnent correctement.

6. **Performance** : Si les benchmarks montrent une régression significative (>20%), investiguer et optimiser si nécessaire.

---

**Prêt à passer au prompt 08 après validation complète.**