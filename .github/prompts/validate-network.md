# ✓ Valider un Réseau RETE

## Contexte

Projet TSD (Type System with Dependencies) - Moteur de règles RETE avec système de contraintes en Go.

Tu veux valider qu'un réseau RETE est correctement construit, que la propagation fonctionne comme attendu, et que les résultats sont corrects.

## Objectif

Effectuer une validation complète d'un réseau RETE : structure, propagation, conditions, et résultats.

## ⚠️ RÈGLES STRICTES - TESTS ET VALIDATION RETE

### 🚫 INTERDICTIONS ABSOLUES

1. **AUCUNE SIMULATION DE RÉSULTATS** :
   - ❌ Pas de résultats hardcodés ou simulés
   - ❌ Pas de mock des résultats du réseau RETE
   - ❌ Pas de calcul manuel des tokens attendus
   - ❌ Pas de suppositions sur le nombre de tokens
   - ✅ **TOUJOURS** extraire les résultats du réseau RETE réel
   - ✅ **TOUJOURS** interroger les TerminalNodes
   - ✅ **TOUJOURS** inspecter les mémoires (Left/Right/Result)

2. **EXTRACTION OBLIGATOIRE DEPUIS LE RÉSEAU** :
   ```go
   // ✅ BON - Extraction depuis le réseau
   terminalCount := 0
   for _, terminal := range network.TerminalNodes {
       terminalCount += len(terminal.Memory.GetTokens())
   }
   
   // ✅ BON - Inspection des tokens réels
   for _, terminal := range network.TerminalNodes {
       for _, token := range terminal.Memory.GetTokens() {
           for varName, fact := range token.Bindings {
               // Vérifier les données réelles du réseau
               t.Logf("Binding: %s -> %s", varName, fact.Type)
           }
       }
   }
   
   // ❌ MAUVAIS - Simulation
   expectedTokens := 5  // Calculé manuellement - INTERDIT !
   ```

3. **VALIDATION AVEC DONNÉES RÉSEAU RÉELLES** :
   - ✅ Compter les tokens dans les TerminalNodes
   - ✅ Vérifier les bindings dans les tokens
   - ✅ Inspecter les mémoires des JoinNodes
   - ✅ Tracer la propagation réelle avec logs
   - ✅ Extraire les actions du TupleSpace
   - ❌ Ne jamais supposer le nombre de tokens
   - ❌ Ne jamais simuler les résultats
   - ❌ Ne jamais hardcoder les résultats attendus

### ✅ BONNES PRATIQUES OBLIGATOIRES

1. **Code Golang** (si code de validation créé) :
   - ❌ Aucun hardcoding de valeurs
   - ✅ Code générique avec paramètres
   - ✅ Constantes nommées pour toutes les valeurs
   - ✅ Respect des conventions Go (Effective Go)
   - ✅ go vet et golangci-lint sans erreur
   - ✅ Gestion explicite des erreurs

2. **Tests de Validation** :
   - ✅ Extraction réelle depuis le réseau RETE
   - ✅ Validation des structures de données réelles
   - ✅ Messages d'assertion explicites et détaillés
   - ✅ Tests déterministes et isolés
   - ✅ Logs de propagation activés pour debug

**Exemples** :

❌ **MAUVAIS - Résultats simulés** :
```go
// Ne JAMAIS faire ça !
func TestNetworkValidation(t *testing.T) {
    network := buildNetwork()
    
    // ❌ Résultat hardcodé/simulé
    expectedTokens := 3  // Calculé manuellement
    
    if actualTokens != expectedTokens {
        t.Errorf("Attendu %d tokens", expectedTokens)
    }
}
```

✅ **BON - Extraction depuis le réseau** :
```go
// Toujours faire ça !
func TestNetworkValidation(t *testing.T) {
    network := buildNetwork()
    
    // Soumettre les faits
    network.SubmitFact(fact1)
    network.SubmitFact(fact2)
    
    // ✅ Extraire les résultats réels du réseau
    actualTokens := 0
    for _, terminal := range network.TerminalNodes {
        actualTokens += len(terminal.Memory.GetTokens())
    }
    
    t.Logf("Tokens terminaux trouvés: %d", actualTokens)
    
    // ✅ Inspecter les tokens réels
    for _, terminal := range network.TerminalNodes {
        for _, token := range terminal.Memory.GetTokens() {
            t.Logf("Token avec %d faits:", len(token.Facts))
            for varName, fact := range token.Bindings {
                t.Logf("  %s -> %s (ID: %s)", varName, fact.Type, fact.ID)
            }
        }
    }
    
    // Validation basée sur les données réelles extraites
    if actualTokens == 0 {
        t.Error("Aucun token terminal créé")
    }
}
```

## Instructions

### 1. Identifier le Réseau à Valider

**Précise** :
- **Fichier contrainte** : Chemin du fichier `.constraint`
- **Fichier faits** : Chemin du fichier `.facts` (optionnel)
- **Contexte** : Nouveau réseau, modification, ou debug ?
- **Comportement attendu** : Ce qui devrait se produire

**Exemple** :
```
Contrainte : beta_coverage_tests/join_simple.constraint
Faits : beta_coverage_tests/join_simple.facts
Contexte : Validation après ajout support multi-variables
Attendu : 5 tokens terminaux créés
```

### 2. Validation de la Structure

#### A. Parser la Contrainte

1. **Vérifier le parsing** :
   ```bash
   # Test de parsing seul
   go run cmd/tsd/main.go parse beta_coverage_tests/join_simple.constraint
   ```

2. **Analyser l'AST** :
   - Types déclarés présents
   - Variables correctement identifiées
   - Conditions bien formées
   - Actions valides

#### B. Construire le Réseau

1. **Construction du réseau** :
   - TypeNodes créés pour chaque type
   - AlphaNodes créés pour filtres
   - BetaNodes (Join/Exists/Not) créés pour jointures
   - TerminalNodes créés pour actions

2. **Vérifier la topologie** :
   - Connexions entre nœuds correctes
   - Pas de nœuds orphelins
   - Pas de cycles (sauf volontaires)
   - Mémoires initialisées

#### C. Afficher la Structure

```go
// Dans le code de test
t.Logf("✅ Réseau RETE construit")
t.Logf("   TypeNodes: %d", len(network.TypeNodes))
t.Logf("   AlphaNodes: %d", len(network.AlphaNodes))
t.Logf("   BetaNodes: %d", len(network.BetaNodes))
t.Logf("   TerminalNodes: %d", len(network.TerminalNodes))

// Afficher la hiérarchie
for _, typeNode := range network.TypeNodes {
    t.Logf("   Type: %s → %d enfants", typeNode.TypeName, len(typeNode.Children))
}
```

### 3. Validation de la Propagation

#### A. Injecter des Faits

1. **Préparer les faits de test** :
   ```json
   TestPerson(id:P1, age:25, name:Alice)
   TestOrder(id:O1, customer_id:P1, amount:100)
   ```

2. **Soumettre les faits un par un** :
   ```go
   err := network.SubmitFact(personFact)
   if err != nil {
       t.Fatalf("❌ Erreur soumission: %v", err)
   }
   ```

3. **Observer la propagation** :
   - Activer le mode verbose pour voir les logs
   - Vérifier que les faits passent par les bons nœuds
   - Contrôler les mémoires Left/Right/Result

#### B. Tracer le Flux

**En mode verbose, observer** :
```
🔥 Soumission d'un nouveau fait: Person(id:P1)
🔗 ALPHA PASSTHROUGH[rule_0_pass_p]: Propagation LEFT
    p -> Person (ID: P1)
🔗 JOINNODE[rule_0_join]: Combinaison avec Order
    p -> Person (ID: P1)
    o -> Order (ID: O1)
    Condition: p.id == o.customer_id
  ✅ Condition satisfaite
🎯 ACTION DISPONIBLE: process_order(Person(...), Order(...))
```

#### C. Vérifier les Mémoires

```go
// Vérifier la mémoire du JoinNode
leftTokens := joinNode.LeftMemory.GetTokens()
t.Logf("Tokens gauche: %d", len(leftTokens))

rightTokens := joinNode.RightMemory.GetTokens()
t.Logf("Tokens droite: %d", len(rightTokens))

resultTokens := joinNode.ResultMemory.GetTokens()
t.Logf("Tokens résultat: %d", len(resultTokens))
```

### 4. Validation des Conditions

#### A. Conditions Alpha (filtres simples)

**Vérifier** :
- Comparaisons numériques (>, <, >=, <=, ==, !=)
- Comparaisons de chaînes
- Expressions arithmétiques (p.age * 2 > 50)
- Conditions booléennes (p.active == true)

**Exemple de test** :
```go
// Fait qui doit passer
passedFact := &Fact{ID: "P1", Type: "Person", Fields: map[string]interface{}{"age": 30}}
// Fait qui doit être filtré
filteredFact := &Fact{ID: "P2", Type: "Person", Fields: map[string]interface{}{"age": 15}}
```

#### B. Conditions Beta (jointures)

**Vérifier** :
- Égalité entre champs (p.id == o.customer_id)
- Comparaisons entre faits (p.salary > o.total)
- Conditions composées (AND, OR)
- Variables disponibles au bon moment

**Exemple de test** :
```go
// Fait Person
person := &Fact{ID: "P1", Fields: map[string]interface{}{"id": "P1", "salary": 5000}}
// Fait Order qui match
matchOrder := &Fact{ID: "O1", Fields: map[string]interface{}{"customer_id": "P1", "total": 3000}}
// Fait Order qui ne match pas
noMatchOrder := &Fact{ID: "O2", Fields: map[string]interface{}{"customer_id": "P999", "total": 3000}}
```

#### C. Conditions Complexes

**Vérifier** :
- Expressions avec plusieurs AND/OR
- Parenthèses et priorités
- Fonctions d'agrégation (SUM, AVG, COUNT, MIN, MAX)
- Négations (NOT, NOT EXISTS)

### 5. Validation des Résultats

#### A. Compter les Tokens Terminaux

```go
// ⚠️ IMPORTANT : TOUJOURS extraire depuis le réseau, ne JAMAIS hardcoder expectedCount
terminalCount := 0
for _, terminal := range network.TerminalNodes {
    terminalCount += len(terminal.Memory.GetTokens())
}

t.Logf("✅ Tokens terminaux extraits du réseau: %d", terminalCount)

// Validation basée sur les données réelles
if terminalCount == 0 {
    t.Error("❌ Aucun token terminal créé")
} else {
    t.Logf("✅ %d tokens terminaux présents dans le réseau", terminalCount)
}
```

#### B. Vérifier le Contenu des Tokens

```go
// ✅ Extraction et inspection des tokens réels
for _, terminal := range network.TerminalNodes {
    tokens := terminal.Memory.GetTokens()
    t.Logf("TerminalNode %s: %d tokens", terminal.GetID(), len(tokens))
    
    for i, token := range tokens {
        t.Logf("  Token %d: %d faits", i, len(token.Facts))
        for varName, fact := range token.Bindings {
            t.Logf("    %s -> %s (ID: %s)", varName, fact.Type, fact.ID)
            // Valider les champs du fait
            for fieldName, fieldValue := range fact.Fields {
                t.Logf("      %s: %v", fieldName, fieldValue)
            }
        }
    }
}
```

#### C. Valider les Actions

```go
// ✅ Extraire les actions du TupleSpace réel
actions := terminal.TupleSpace.GetAllActions()
t.Logf("Actions extraites du TupleSpace: %d", len(actions))

for i, action := range actions {
    t.Logf("Action %d: %s avec %d arguments", i, action.Name, len(action.Args))
    // Valider les arguments en inspectant les données réelles
    for j, arg := range action.Args {
        t.Logf("  Arg %d: %v (type: %T)", j, arg, arg)
    }
}
```

### 6. Tests de Non-Régression

#### A. Tests Positifs

**Faits qui DOIVENT matcher** :
```go
// ⚠️ Ne JAMAIS hardcoder expectedTokens - extraire du réseau !
testCases := []struct{
    name string
    facts []*Fact
    validate func(*testing.T, *Network)
}{
    {
        "cas_nominal",
        []*Fact{person1, order1},
        func(t *testing.T, net *Network) {
            // ✅ Extraire du réseau réel
            count := 0
            for _, term := range net.TerminalNodes {
                count += len(term.Memory.GetTokens())
            }
            t.Logf("Tokens extraits: %d", count)
            if count == 0 {
                t.Error("Aucun token créé")
            }
        },
    },
}
```

#### B. Tests Négatifs

**Faits qui NE DOIVENT PAS matcher** :
```go
// ✅ Validation basée sur extraction réelle du réseau
negativeTests := []struct{
    name string
    facts []*Fact
    reason string
}{
    {"age_insuffisant", []*Fact{person_minor}, "age < 18"},
    {"customer_id_invalide", []*Fact{person1, order_wrong_id}, "IDs ne matchent pas"},
}

// Exécuter les tests
for _, tt := range negativeTests {
    t.Run(tt.name, func(t *testing.T) {
        network := buildNetwork()
        for _, fact := range tt.facts {
            network.SubmitFact(fact)
        }
        
        // ✅ Extraire les résultats réels
        tokenCount := 0
        for _, term := range network.TerminalNodes {
            tokenCount += len(term.Memory.GetTokens())
        }
        
        if tokenCount > 0 {
            t.Errorf("❌ Tokens créés alors qu'ils ne devraient pas (%s)", tt.reason)
        }
    })
}
```

#### C. Tests de Cas Limites

```go
// Cas edge à tester
edgeCases := []string{
    "Aucun fait soumis",
    "Un seul fait (jointure incomplète)",
    "Faits dans le désordre",
    "Fait soumis deux fois",
    "Fait rétracté puis re-soumis",
    "Valeurs NULL/nil",
    "Valeurs à la limite (INT_MAX, etc.)",
}
```

## Critères de Validation

### ✅ Structure Valide

- [ ] Parsing réussi sans erreur
- [ ] Tous les types déclarés présents dans TypeNodes
- [ ] AlphaNodes créés pour les filtres
- [ ] BetaNodes créés pour les jointures
- [ ] TerminalNodes créés pour les actions
- [ ] Connexions entre nœuds correctes
- [ ] Pas de nœuds orphelins

### ✅ Propagation Correcte

- [ ] Faits arrivent aux bons nœuds
- [ ] Mémoires Left/Right/Result correctement remplies
- [ ] Conditions évaluées au bon moment
- [ ] Variables liées au bon moment
- [ ] Tokens créés/propagés correctement

### ✅ Conditions Valides

- [ ] Conditions alpha (filtres) fonctionnent
- [ ] Conditions beta (jointures) fonctionnent
- [ ] Expressions arithmétiques calculées correctement
- [ ] Opérateurs logiques (AND/OR) respectés
- [ ] Cas limites gérés (nil, valeurs extrêmes)

### ✅ Résultats Corrects

- [ ] Nombre de tokens terminaux correct
- [ ] Contenu des tokens correct
- [ ] Actions déclenchées aux bons moments
- [ ] Pas de faux positifs
- [ ] Pas de faux négatifs

### ✅ Performance Acceptable

- [ ] Construction du réseau rapide (< 1s)
- [ ] Injection de faits rapide (< 10ms/fait)
- [ ] Pas de fuites mémoires
- [ ] Complexité algorithmique raisonnable

## Format de Réponse Attendu

```
=== VALIDATION RÉSEAU RETE ===

📁 Fichiers
- Contrainte : beta_coverage_tests/join_simple.constraint
- Faits : beta_coverage_tests/join_simple.facts

🏗️ Structure du Réseau
✅ Parsing réussi
✅ 2 TypeNodes créés (Person, Order)
✅ 2 AlphaNodes créés (filtres age, status)
✅ 1 JoinNode créé (jointure Person-Order)
✅ 1 TerminalNode créé (action process_order)

📊 Topologie
Person (TypeNode)
  └─> AlphaNode[age >= 18]
      └─> JoinNode[p.id == o.customer_id]
          
Order (TypeNode)
  └─> PassthroughAlpha
      └─> JoinNode[p.id == o.customer_id]
          └─> TerminalNode[process_order]

🔄 Propagation
Test 1: Person(id:P1, age:25)
  ✅ Passe par TypeNode Person
  ✅ Filtre age >= 18 : PASS
  ✅ Stocké dans JoinNode.LeftMemory
  ⏸️  Attend Order correspondant

Test 2: Order(id:O1, customer_id:P1)
  ✅ Passe par TypeNode Order
  ✅ Stocké dans JoinNode.RightMemory
  ✅ Jointure avec Person P1
  ✅ Condition p.id == o.customer_id : PASS (P1 == P1)
  ✅ Token terminal créé

Test 3: Order(id:O2, customer_id:P999)
  ✅ Passe par TypeNode Order
  ✅ Stocké dans JoinNode.RightMemory
  ❌ Jointure avec Person P1
  ❌ Condition p.id == o.customer_id : FAIL (P1 != P999)
  ⏸️  Aucun token créé (filtrage correct)

✅ Conditions
✅ Alpha : age >= 18 validée
✅ Beta : p.id == o.customer_id validée
✅ Filtrage : O2 correctement rejeté

📦 Résultats
✅ 1 token terminal créé (attendu : 1)
✅ Token contient Person P1 et Order O1
✅ Action process_order disponible

⚡ Performance
✅ Construction réseau : 5ms
✅ Injection 3 faits : 12ms (4ms/fait)
✅ Mémoire utilisée : 2.3 MB

🎯 Verdict : RÉSEAU VALIDE ✅

Le réseau RETE fonctionne correctement. La propagation est conforme,
les conditions sont bien évaluées, et les résultats sont exacts.
```

## Commandes Utiles

```bash
# Valider un fichier constraint avec le runner universel
make rete-unified

# Valider un test spécifique
go test -v -run TestNomDuTest ./test/integration

# Mode verbose pour voir la propagation
go test -v -run TestNomDuTest ./test/integration 2>&1 | grep "🔥\|🔗\|✅\|❌"

# Avec profiling mémoire
go test -memprofile mem.prof -run TestNomDuTest ./test/integration
go tool pprof mem.prof

# Benchmark de propagation
go test -bench=BenchmarkPropagation -benchmem ./rete
```

## Exemple d'Utilisation

```
J'ai créé un nouveau réseau RETE dans beta_coverage_tests/join_complex.constraint
qui fait une jointure 3-way entre Person, Order et Product.

Peux-tu valider le réseau en utilisant le prompt "validate-network" ?

Je m'attends à ce que :
- 3 TypeNodes soient créés
- La jointure Person-Order-Product fonctionne
- 2 tokens terminaux soient créés avec les faits du fichier join_complex.facts
```

## Checklist de Validation

### Avant le Test
- [ ] Fichier .constraint syntaxiquement correct
- [ ] Fichier .facts préparé (si nécessaire)
- [ ] **AUCUN hardcoding de résultats attendus**
- [ ] **AUCUNE simulation de tokens**
- [ ] Cas limites identifiés

### Pendant le Test
- [ ] Mode verbose activé pour observation
- [ ] Logs de propagation analysés
- [ ] **Mémoires des nœuds EXTRAITES du réseau réel**
- [ ] **Tokens COMPTÉS depuis TerminalNodes**
- [ ] Conditions évaluées correctement

### Après le Test
- [ ] **Résultats EXTRAITS du réseau (pas simulés)**
- [ ] **Validation basée sur données réelles uniquement**
- [ ] Pas d'erreurs ni de warnings
- [ ] Performance acceptable
- [ ] **Code sans hardcoding** (si code créé)
- [ ] **go vet et golangci-lint** sans erreur
- [ ] Documentation mise à jour

## Outils de Diagnostic

### Afficher la Structure du Réseau

```go
func PrintNetworkStructure(network *Network) {
    fmt.Println("=== STRUCTURE DU RÉSEAU RETE ===")
    
    for _, typeNode := range network.TypeNodes {
        fmt.Printf("📦 TypeNode: %s\n", typeNode.TypeName)
        printNodeTree(typeNode, 1)
    }
}

func printNodeTree(node Node, depth int) {
    indent := strings.Repeat("  ", depth)
    fmt.Printf("%s└─> %s [%s]\n", indent, node.GetID(), node.GetType())
    
    for _, child := range node.GetChildren() {
        printNodeTree(child, depth+1)
    }
}
```

### Afficher les Mémoires

```go
func PrintNodeMemories(node Node) {
    if joinNode, ok := node.(*JoinNode); ok {
        fmt.Printf("JoinNode %s:\n", joinNode.ID)
        fmt.Printf("  Left: %d tokens\n", len(joinNode.LeftMemory.Tokens))
        fmt.Printf("  Right: %d tokens\n", len(joinNode.RightMemory.Tokens))
        fmt.Printf("  Result: %d tokens\n", len(joinNode.ResultMemory.Tokens))
    }
}
```

### Tracer la Propagation

```go
// Activer les logs de propagation
network.EnableVerboseMode(true)

// Ou ajouter des logs personnalisés
func (n *JoinNode) ActivateRight(fact *Fact) error {
    log.Printf("🔍 JoinNode %s reçoit fait %s", n.ID, fact.ID)
    // ... reste du code
}
```

## Patterns de Validation

### Pattern 1 : Test Unitaire par Nœud

Tester chaque type de nœud isolément avant de tester le réseau complet.

### Pattern 2 : Test d'Intégration par Règle

Tester chaque règle individuellement avec des faits minimaux.

### Pattern 3 : Test End-to-End

Tester le réseau complet avec tous les faits et vérifier le résultat final.

### Pattern 4 : Test de Régression

Tester avec des cas connus qui fonctionnaient avant une modification.

## Résolution de Problèmes Courants

### Problème : Aucun token terminal créé

**Causes possibles** :
- Conditions trop restrictives
- Faits ne matchent pas
- Nœuds mal connectés
- Variables non liées

**Solution** :
1. Vérifier les logs de propagation
2. Tester les conditions isolément
3. Vérifier la topologie du réseau
4. Ajouter des logs dans les nœuds

### Problème : Trop de tokens créés

**Causes possibles** :
- Conditions trop permissives
- Cartesian product non intentionnel
- Pas assez de filtres alpha

**Solution** :
1. Vérifier les conditions de jointure
2. Ajouter des filtres alpha
3. Vérifier les types des faits

### Problème : Performance dégradée

**Causes possibles** :
- Trop de faits en mémoire
- Boucles de réévaluation
- Allocations excessives

**Solution** :
1. Profiler avec pprof
2. Vérifier la complexité algorithmique
3. Optimiser les structures de données
4. Utiliser sync.Pool si nécessaire

## Ressources

- [RETE Algorithm](https://en.wikipedia.org/wiki/Rete_algorithm)
- [Tests d'intégration](../../test/integration/)
- [Runner universel](../../cmd/universal-rete-runner/)
- [Documentation RETE](../../docs/)

---

**Rappel** : Un réseau RETE validé est un réseau fiable !