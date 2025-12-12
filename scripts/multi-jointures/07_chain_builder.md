# Prompt 07 : BetaChainBuilder

**Session** : 7/12  
**Durée estimée** : 3-4 heures  
**Pré-requis** : Prompt 06 complété, ActivateLeft/Right refactorées

---

## 🎯 Objectif de cette Session

S'assurer que le BetaChainBuilder construit correctement les cascades de jointures avec :
1. AllVariables contenant TOUTES les variables cumulées à chaque niveau
2. RightVariables contenant la nouvelle variable à chaque cascade
3. LeftVariables contenant toutes les variables des niveaux précédents
4. VariableTypes correctement renseigné pour toutes les variables

**Livrables finaux** : 
- `tsd/rete/builder_beta_chain.go` (vérifié/corrigé)
- `tsd/rete/builder_join_rules_cascade.go` (vérifié/corrigé)

---

## 📋 Tâches à Réaliser

### Tâche 1 : Analyser la Construction Actuelle (40 min)

#### 1.1 Lire buildJoinPatterns

**Fichier** : `tsd/rete/builder_join_rules_cascade.go`

**Chercher la fonction** : `buildJoinPatterns` ou similaire qui crée les patterns de cascade.

**Questions à répondre** :
1. Comment les JoinPatterns sont-ils créés pour N variables ?
2. Pour une règle `{u: User, o: Order, p: Product}`, quels patterns sont créés ?
3. Les AllVars sont-ils corrects à chaque niveau ?
4. Les LeftVars/RightVars sont-ils corrects ?

**Exemple attendu pour 3 variables [u, o, p]** :

```
Pattern 1:
  LeftVars: [u]
  RightVars: [o]
  AllVars: [u, o]
  
Pattern 2:
  LeftVars: [u, o]
  RightVars: [p]
  AllVars: [u, o, p]
```

---

#### 1.2 Analyser la fonction de construction

**Chercher dans le code** :

```go
// Exemple de structure attendue
func buildJoinPatterns(variableNames []string, variableTypes map[string]string, ...) []JoinPattern {
    patterns := make([]JoinPattern, 0, len(variableNames)-1)
    
    // Pattern 1: Premières 2 variables
    // Pattern 2+: Chaque variable supplémentaire
    
    return patterns
}
```

**Vérifier** :
- Le nombre de patterns = nombre de variables - 1
- Chaque pattern a AllVars qui s'incrémente
- Les conditions de jointure sont attachées au bon pattern

---

### Tâche 2 : Corriger buildJoinPatterns si Nécessaire (60 min)

#### 2.1 Implémenter la logique correcte

**Code attendu** :

```go
func buildJoinPatterns(variableNames []string, variableTypes map[string]string, conditions []interface{}) []JoinPattern {
    if len(variableNames) < 2 {
        return []JoinPattern{}
    }
    
    patterns := make([]JoinPattern, 0, len(variableNames)-1)
    
    // Pattern 1: Joindre les 2 premières variables
    pattern1 := JoinPattern{
        LeftVars:  []string{variableNames[0]},
        RightVars: []string{variableNames[1]},
        AllVars:   []string{variableNames[0], variableNames[1]},
        VariableTypes: variableTypes,
        Conditions: extractConditionsForVariables(conditions, []string{variableNames[0], variableNames[1]}),
    }
    patterns = append(patterns, pattern1)
    
    // Patterns suivants: Ajouter chaque variable une par une
    for i := 2; i < len(variableNames); i++ {
        // LeftVars = TOUTES les variables précédentes [0..i-1]
        leftVars := make([]string, i)
        copy(leftVars, variableNames[0:i])
        
        // RightVars = La nouvelle variable seulement
        rightVars := []string{variableNames[i]}
        
        // AllVars = TOUTES les variables jusqu'à i (inclus) [0..i]
        allVars := make([]string, i+1)
        copy(allVars, variableNames[0:i+1])
        
        pattern := JoinPattern{
            LeftVars:      leftVars,
            RightVars:     rightVars,
            AllVars:       allVars,
            VariableTypes: variableTypes,
            Conditions:    extractConditionsForVariables(conditions, allVars),
        }
        
        patterns = append(patterns, pattern)
    }
    
    return patterns
}
```

**Points critiques** :
- `copy()` pour éviter les références partagées
- AllVars s'incrémente : [u, o] puis [u, o, p] puis [u, o, p, task], etc.
- VariableTypes contient TOUS les types, pas seulement pour ce pattern

---

#### 2.2 Implémenter extractConditionsForVariables

**Fonction helper** :

```go
// extractConditionsForVariables filtre les conditions qui ne concernent que les variables données.
func extractConditionsForVariables(conditions []interface{}, variables []string) []interface{} {
    if conditions == nil || len(conditions) == 0 {
        return nil
    }
    
    // Créer un set des variables pour recherche rapide
    varSet := make(map[string]bool)
    for _, v := range variables {
        varSet[v] = true
    }
    
    filtered := make([]interface{}, 0)
    
    for _, cond := range conditions {
        // Extraire les variables mentionnées dans la condition
        varsInCondition := extractVariablesFromCondition(cond)
        
        // Vérifier que toutes les variables de la condition sont disponibles
        allAvailable := true
        for _, v := range varsInCondition {
            if !varSet[v] {
                allAvailable = false
                break
            }
        }
        
        if allAvailable {
            filtered = append(filtered, cond)
        }
    }
    
    return filtered
}
```

---

### Tâche 3 : Vérifier la Création des JoinNodes (50 min)

#### 3.1 Analyser createJoinNodesFromPatterns

**Fichier** : `tsd/rete/builder_beta_chain.go`

**Vérifier que chaque JoinNode créé a** :

```go
joinNode := &JoinNode{
    BaseNode: BaseNode{ID: generateJoinNodeID()},
    LeftVariables:  pattern.LeftVars,   // ✅ Doit être correct
    RightVariables: pattern.RightVars,  // ✅ Doit contenir la nouvelle variable
    AllVariables:   pattern.AllVars,    // ✅ CRITIQUE: toutes les variables
    VariableTypes:  pattern.VariableTypes, // ✅ Tous les types
    JoinConditions: pattern.Conditions,
    LeftMemory:     []*Token{},
    RightMemory:    []*Fact{},
}
```

**Points de vérification** :
- AllVariables doit contenir TOUTES les variables jusqu'à ce niveau
- VariableTypes doit contenir le type de chaque variable
- Les conditions sont correctement filtrées

---

#### 3.2 Vérifier la connexion entre JoinNodes

**Pour une cascade de 3 variables** :

```
JoinNode1 [u, o]  ────→  JoinNode2 [u, o, p]  ────→  TerminalNode

Configuration:
JoinNode1:
  - AllVariables: [u, o]
  - RightVariables: [o]
  
JoinNode2:
  - AllVariables: [u, o, p]
  - LeftVariables: [u, o]  ← Provient de JoinNode1
  - RightVariables: [p]
```

**S'assurer que** :
- La sortie de JoinNode1 est l'entrée (côté gauche) de JoinNode2
- Les TypeNodes sont connectés aux bons côtés (Right)

---

### Tâche 4 : Ajouter du Logging de Validation (30 min)

#### 4.1 Logger la construction des patterns

**Dans buildJoinPatterns, ajouter** :

```go
func buildJoinPatterns(...) []JoinPattern {
    // ... code existant ...
    
    // Logging de validation (TEMPORAIRE)
    fmt.Printf("\n🏗️  Building join patterns for variables: %v\n", variableNames)
    for i, pattern := range patterns {
        fmt.Printf("   Pattern %d:\n", i+1)
        fmt.Printf("     - LeftVars:  %v\n", pattern.LeftVars)
        fmt.Printf("     - RightVars: %v\n", pattern.RightVars)
        fmt.Printf("     - AllVars:   %v\n", pattern.AllVars)
        fmt.Printf("     - Conditions: %d\n", len(pattern.Conditions))
    }
    
    return patterns
}
```

---

#### 4.2 Logger la création des JoinNodes

```go
func createJoinNodesFromPatterns(...) []*JoinNode {
    // ... code existant ...
    
    // Logging (TEMPORAIRE)
    fmt.Printf("\n🔧 Created JoinNodes:\n")
    for i, jn := range joinNodes {
        fmt.Printf("   JoinNode %d (ID: %s):\n", i+1, jn.ID)
        fmt.Printf("     - LeftVariables:  %v\n", jn.LeftVariables)
        fmt.Printf("     - RightVariables: %v\n", jn.RightVariables)
        fmt.Printf("     - AllVariables:   %v\n", jn.AllVariables)
    }
    
    return joinNodes
}
```

---

### Tâche 5 : Tests de Construction (60 min)

#### 5.1 Test de buildJoinPatterns pour 3 variables

**Fichier** : `tsd/rete/builder_join_rules_cascade_test.go`

```go
func TestBuildJoinPatterns_3Variables(t *testing.T) {
    t.Log("🧪 TEST buildJoinPatterns - 3 variables")
    
    variableNames := []string{"u", "o", "p"}
    variableTypes := map[string]string{
        "u": "User",
        "o": "Order",
        "p": "Product",
    }
    
    patterns := buildJoinPatterns(variableNames, variableTypes, nil)
    
    // Doit créer 2 patterns (N-1)
    if len(patterns) != 2 {
        t.Fatalf("❌ Attendu 2 patterns, got %d", len(patterns))
    }
    
    // Pattern 1: [u] + [o] = [u, o]
    p1 := patterns[0]
    if !slicesEqual(p1.LeftVars, []string{"u"}) {
        t.Errorf("❌ Pattern 1 LeftVars incorrect: %v", p1.LeftVars)
    }
    if !slicesEqual(p1.RightVars, []string{"o"}) {
        t.Errorf("❌ Pattern 1 RightVars incorrect: %v", p1.RightVars)
    }
    if !slicesEqual(p1.AllVars, []string{"u", "o"}) {
        t.Errorf("❌ Pattern 1 AllVars incorrect: %v", p1.AllVars)
    }
    
    // Pattern 2: [u, o] + [p] = [u, o, p]
    p2 := patterns[1]
    if !slicesEqual(p2.LeftVars, []string{"u", "o"}) {
        t.Errorf("❌ Pattern 2 LeftVars incorrect: %v", p2.LeftVars)
    }
    if !slicesEqual(p2.RightVars, []string{"p"}) {
        t.Errorf("❌ Pattern 2 RightVars incorrect: %v", p2.RightVars)
    }
    if !slicesEqual(p2.AllVars, []string{"u", "o", "p"}) {
        t.Errorf("❌ Pattern 2 AllVars incorrect: %v", p2.AllVars)
    }
    
    t.Log("✅ Patterns corrects pour 3 variables")
}

func slicesEqual(a, b []string) bool {
    if len(a) != len(b) {
        return false
    }
    for i := range a {
        if a[i] != b[i] {
            return false
        }
    }
    return true
}
```

---

#### 5.2 Test de buildJoinPatterns pour N variables

```go
func TestBuildJoinPatterns_NVariables(t *testing.T) {
    t.Log("🧪 TEST buildJoinPatterns - N variables")
    
    for n := 2; n <= 5; n++ {
        t.Run(fmt.Sprintf("n=%d", n), func(t *testing.T) {
            // Générer N variables
            vars := make([]string, n)
            types := make(map[string]string)
            for i := 0; i < n; i++ {
                vars[i] = fmt.Sprintf("v%d", i)
                types[vars[i]] = fmt.Sprintf("Type%d", i)
            }
            
            patterns := buildJoinPatterns(vars, types, nil)
            
            // Vérifier nombre de patterns
            if len(patterns) != n-1 {
                t.Errorf("❌ Pour %d variables, attendu %d patterns, got %d", n, n-1, len(patterns))
            }
            
            // Vérifier chaque pattern
            for i, pattern := range patterns {
                expectedAllVars := i + 2 // Pattern i joint (i+2) variables
                if len(pattern.AllVars) != expectedAllVars {
                    t.Errorf("❌ Pattern %d: attendu %d AllVars, got %d", 
                        i, expectedAllVars, len(pattern.AllVars))
                }
            }
        })
    }
    
    t.Log("✅ Patterns corrects pour N variables")
}
```

---

#### 5.3 Exécuter les tests

```bash
cd tsd
go test -v ./rete/builder_join_rules_cascade_test.go
```

---

### Tâche 6 : Validation avec Test E2E (40 min)

#### 6.1 Créer un test de bout en bout

**Fichier** : `tsd/rete/builder_cascade_integration_test.go` (nouveau)

```go
func TestBetaChainBuilder_BuildCascade3Variables(t *testing.T) {
    t.Log("🧪 TEST BetaChainBuilder - Cascade 3 variables")
    
    // Créer une règle avec 3 variables
    rule := &Rule{
        Name: "test_rule",
        Variables: []Variable{
            {Name: "u", Type: "User"},
            {Name: "o", Type: "Order"},
            {Name: "p", Type: "Product"},
        },
        Conditions: []interface{}{
            // Conditions de jointure ici
        },
        Action: &Action{
            Name: "test_action",
        },
    }
    
    // Construire la cascade
    builder := NewBetaChainBuilder()
    terminalNode, err := builder.BuildCascadeForRule(rule)
    
    if err != nil {
        t.Fatalf("❌ Erreur construction: %v", err)
    }
    
    // Vérifier la structure créée
    // (Remonter depuis TerminalNode jusqu'aux JoinNodes)
    
    // TODO: Implémenter vérification de la structure
    // - Vérifier que 2 JoinNodes ont été créés
    // - Vérifier leurs configurations (AllVariables)
    
    t.Log("✅ Cascade 3 variables construite correctement")
}
```

---

#### 6.2 Tester avec une vraie fixture

**Utiliser une fixture existante** :

```bash
# Parser et construire le réseau pour une fixture 3 variables
cd tsd
go test -v -run "TestE2E.*join_multi_variable_complex" ./tests/e2e/
```

**Analyser les logs** pour vérifier :
- Les JoinNodes créés
- Leurs configurations (AllVariables)
- La propagation des tokens

---

### Tâche 7 : Nettoyage et Validation Finale (20 min)

#### 7.1 Supprimer le logging temporaire

**Supprimer ou désactiver** tous les `fmt.Printf` ajoutés pour debug.

---

#### 7.2 Vérifier la qualité du code

```bash
go fmt ./rete/builder_*.go
go vet ./rete/builder_*.go
```

---

#### 7.3 Validation complète

```bash
# Compilation
go build ./rete/...

# Tests unitaires
go test ./rete/builder_*_test.go

# Tests d'intégration
make test-integration
```

---

## ✅ Critères de Validation de cette Session

À la fin de ce prompt, vous devez avoir :

### Code
- [ ] ✅ `buildJoinPatterns` crée les bons patterns pour N variables
- [ ] ✅ AllVariables s'incrémente correctement à chaque niveau
- [ ] ✅ LeftVars contient toutes les variables précédentes
- [ ] ✅ RightVars contient la nouvelle variable seulement
- [ ] ✅ VariableTypes contient tous les types

### Tests
- [ ] ✅ `TestBuildJoinPatterns_3Variables` passe
- [ ] ✅ `TestBuildJoinPatterns_NVariables` passe (N=2 à 5)
- [ ] ✅ Tests de construction passent
- [ ] ✅ Pas de régression sur tests existants

### Validation
- [ ] ✅ Les cascades 2 variables continuent de fonctionner
- [ ] ✅ Les cascades 3 variables sont construites correctement
- [ ] ✅ Code propre et sans warnings

---

## 🎯 Prochaine Étape

Une fois le BetaChainBuilder **validé**, passer au **Prompt 08 - ExecutionContext et Actions**.

Le Prompt 08 s'assurera que les actions peuvent accéder à toutes les variables via BindingChain.

---

## 💡 Conseils Pratiques

### Pour la Construction
1. **Vérifier les copies** : Utiliser `copy()` pour éviter les références partagées
2. **Tester avec N=2 d'abord** : S'assurer de la non-régression
3. **Logger temporairement** : Voir exactement ce qui est construit

### Pour les Tests
1. **Tests paramétriques** : Tester avec N=2, 3, 4, 5, 10
2. **Vérifier chaque niveau** : AllVars doit s'incrémenter correctement
3. **Vérifier les types** : VariableTypes doit être complet

---

**Note** : Cette session est cruciale - elle garantit que les cascades sont construites avec les bonnes configurations. Si cette étape est correcte, les tokens propagés dans les Prompts 05-06 auront tous les bindings nécessaires.