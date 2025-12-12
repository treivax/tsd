# 🔬 Analyse et Diagnostic - Prompt Universel

> **📋 Standards** : Ce prompt respecte les règles de [common.md](./common.md)

## 🎯 Objectif

Analyser et diagnostiquer : erreurs, comportements étranges, réseaux RETE, ou problèmes de performance.

---

## ⚠️ Rappels Critiques

Avant de commencer, consulter [common.md](./common.md) :
- [🧪 Tests Fonctionnels RÉELS](./common.md#tests-fonctionnels-réels) - Résultats réels, pas de mocks
- [🔧 Outils et Commandes](./common.md#outils-et-commandes) - Validation, profiling

---

## 📋 Instructions

### 1. Définir le Problème

**Précise** :
- **Type** : [ ] Erreur  [ ] Comportement inattendu  [ ] Performance  [ ] Réseau RETE
- **Symptômes** : Que se passe-t-il exactement ?
- **Attendu** : Quel devrait être le comportement ?
- **Reproductibilité** : Toujours, parfois, conditions spécifiques ?
- **Contexte** : Changements récents, environnement

**Si erreur** :
- Message d'erreur complet
- Stack trace
- Logs pertinents

### 2. Analyser une Erreur

#### Étapes d'Analyse

1. **Lire le message d'erreur**
   - Identifier le type d'erreur
   - Localiser l'origine (fichier, ligne)
   - Comprendre le contexte

2. **Examiner la stack trace**
   - Chemin d'exécution
   - Point d'entrée
   - Propagation de l'erreur

3. **Reproduire l'erreur**
   ```bash
   # Exécuter le test qui échoue
   go test -v -run TestProblematic ./...
   
   # Avec plus de détails
   go test -v -race -run TestProblematic ./...
   ```

4. **Isoler la cause**
   - Ajouter logs de debug
   - Tester avec données simplifiées
   - Vérifier les hypothèses

5. **Vérifier l'environnement**
   - Versions dépendances
   - Configuration système
   - Variables d'environnement

#### Types d'Erreurs Courants

**Erreurs de Compilation** :
- Syntax error → Corriger syntaxe
- Type mismatch → Vérifier types
- Undefined → Imports, déclarations

**Erreurs d'Exécution** :
- Nil pointer → Vérifications nil
- Index out of range → Vérifier bounds
- Race condition → `go test -race`

**Erreurs Logiques** :
- Résultat incorrect → Revoir algorithme
- Comportement inattendu → Debug pas à pas

### 3. Investiguer un Comportement Étrange

#### Approche Systématique

1. **Documenter le comportement**
   - Entrées → Sorties observées
   - Entrées → Sorties attendues
   - Différences

2. **Éliminer les hypothèses**
   - Tester avec entrées simples
   - Vérifier un cas à la fois
   - Isoler les variables

3. **Ajouter de l'instrumentation**
   ```go
   // Logs de debug temporaires
   log.Printf("🔍 DEBUG: variable = %+v", variable)
   log.Printf("🔍 État avant: %+v", state)
   log.Printf("🔍 État après: %+v", state)
   ```

4. **Comparer avec un cas qui fonctionne**
   - Qu'est-ce qui diffère ?
   - Conditions différentes ?
   - Configuration différente ?

5. **Vérifier les dépendances**
   - État externe (fichiers, réseau, DB)
   - Ordre d'exécution
   - Concurrence

### 4. Valider un Réseau RETE

#### Vérifications Structure

```go
// Vérifier structure du réseau
func validateNetwork(network *ReteNetwork) error {
    // 1. Tous les nœuds connectés
    if err := checkConnectivity(network); err != nil {
        return fmt.Errorf("connectivité: %w", err)
    }
    
    // 2. Pas de cycles
    if hasCycles(network) {
        return errors.New("cycles détectés")
    }
    
    // 3. Terminal nodes présents
    if len(network.TerminalNodes) == 0 {
        return errors.New("aucun terminal node")
    }
    
    return nil
}
```

#### Vérifications Propagation

**Test de propagation** :
```go
func TestRetePropagation(t *testing.T) {
    // Construire réseau
    network := buildTestNetwork()
    
    // Ajouter faits
    network.AddFact(Fact{Type: "Person", Data: person1})
    network.AddFact(Fact{Type: "Person", Data: person2})
    
    // Exécuter
    network.Execute()
    
    // Vérifier résultats RÉELS (pas de mock)
    results := network.TerminalNodes[0].GetResults()
    if len(results) == 0 {
        t.Error("❌ Aucun résultat produit")
    }
    
    // Inspecter mémoires
    for _, node := range network.BetaNodes {
        leftMem := node.LeftMemory.GetTokens()
        rightMem := node.RightMemory.GetTokens()
        t.Logf("🔍 Node %s: Left=%d, Right=%d", 
               node.ID, len(leftMem), len(rightMem))
    }
}
```

#### Vérifications Résultats

```go
// Vérifier que résultats attendus sont produits
func verifyResults(network *ReteNetwork, expected []Result) error {
    actual := collectResults(network)
    
    if len(actual) != len(expected) {
        return fmt.Errorf("nombre résultats: got %d, want %d", 
                         len(actual), len(expected))
    }
    
    for i, exp := range expected {
        if !resultsEqual(actual[i], exp) {
            return fmt.Errorf("résultat %d: got %v, want %v", 
                             i, actual[i], exp)
        }
    }
    
    return nil
}

// Collecter résultats RÉELS des terminal nodes
func collectResults(network *ReteNetwork) []Result {
    var results []Result
    for _, terminal := range network.TerminalNodes {
        results = append(results, terminal.GetResults()...)
    }
    return results
}
```

### 5. Analyser la Performance

#### Profiling

```bash
# CPU profiling
go test -cpuprofile=cpu.prof -bench=.
go tool pprof cpu.prof

# Memory profiling
go test -memprofile=mem.prof -bench=.
go tool pprof mem.prof

# Trace
go test -trace=trace.out
go tool trace trace.out

# Benchmarks
go test -bench=. -benchmem ./...
```

#### Identifier Goulots

1. **Mesurer d'abord**
   - Baseline performance
   - Profiling pour identifier hotspots

2. **Analyser les résultats**
   - CPU : Fonctions les plus coûteuses
   - Mémoire : Allocations excessives
   - Goroutines : Fuites, blocages

3. **Vérifier algorithmes**
   - Complexité O(n²) → O(n log n) ?
   - Calculs redondants ?
   - Structures de données appropriées ?

---

## 🔧 Outils de Debug

### Logs Structurés

```go
// Utiliser logs avec niveaux
log.Printf("🔍 DEBUG: %s = %+v", name, value)
log.Printf("⚠️  WARN: condition inattendue: %v", condition)
log.Printf("❌ ERROR: échec traitement: %v", err)
log.Printf("✅ INFO: traitement réussi")
```

### Tests de Debug

```go
// Test minimal pour isoler problème
func TestDebug_MinimalCase(t *testing.T) {
    t.Log("🔍 Test de debug - cas minimal")
    
    // Cas le plus simple possible
    input := createMinimalInput()
    
    result, err := Function(input)
    
    t.Logf("Input: %+v", input)
    t.Logf("Result: %+v", result)
    t.Logf("Error: %v", err)
    
    // Vérifications
    if err != nil {
        t.Errorf("❌ Erreur inattendue: %v", err)
    }
}
```

### Delve Debugger

```bash
# Installer delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug test
dlv test -- -test.run TestProblematic

# Commandes utiles dans dlv
# break <file>:<line>  - Breakpoint
# continue             - Continuer
# print <var>          - Afficher variable
# step                 - Pas à pas
# next                 - Ligne suivante
```

---

## ✅ Checklist Diagnostic

- [ ] Problème clairement défini
- [ ] Reproductibilité vérifiée
- [ ] Message erreur/logs collectés
- [ ] Environnement documenté
- [ ] Hypothèses testées une par une
- [ ] Résultats RÉELS extraits (pas de mocks)
- [ ] Cause racine identifiée (pas symptôme)
- [ ] Solution testée
- [ ] Non-régression validée

---

## 🎯 Principes

1. **Méthodique** : Approche systématique, pas au hasard
2. **Isoler** : Réduire le problème au minimum
3. **Mesurer** : Données objectives, pas intuitions
4. **Documenter** : Noter observations et hypothèses
5. **Valider** : Vérifier la solution avant de conclure

---

## 🚫 Erreurs Courantes

- ❌ Chercher au hasard sans méthode
- ❌ Tester plusieurs changements simultanés
- ❌ Supposer sans vérifier
- ❌ Ignorer les logs/messages d'erreur
- ❌ Ne pas reproduire de manière fiable
- ❌ Corriger le symptôme, pas la cause
- ❌ Utiliser des mocks au lieu des résultats réels
- ❌ Ne pas tester la solution

---

## 📊 Commandes Utiles

```bash
# Exécution et debug
go test -v -run TestName ./...           # Test spécifique
go test -race ./...                       # Race conditions
go test -v -count=1 ./...                # Sans cache

# Profiling
go test -cpuprofile=cpu.prof -bench=.    # CPU
go test -memprofile=mem.prof -bench=.    # Mémoire
go tool pprof cpu.prof                    # Analyser profil

# Analyse statique
go vet ./...                              # Vérifications Go
staticcheck ./...                         # Analyse statique
errcheck ./...                            # Erreurs non gérées

# Validation
make validate                             # Validation complète
```

---

## 📝 Template Rapport d'Analyse

```markdown
## 🔬 Analyse : [Titre du Problème]

### 📋 Symptômes
[Description du problème observé]

### 🎯 Attendu
[Comportement attendu]

### 🔍 Investigation

#### Hypothèse 1: [Description]
- Test effectué: [...]
- Résultat: [...]
- Conclusion: ✅/❌

#### Hypothèse 2: [Description]
- Test effectué: [...]
- Résultat: [...]
- Conclusion: ✅/❌

### 💡 Cause Racine
[Cause identifiée]

### ✅ Solution
[Solution proposée/appliquée]

### 🧪 Validation
- [ ] Solution testée
- [ ] Tests passent
- [ ] Pas de régression
- [ ] Documentation mise à jour
```

---

## 📚 Ressources

- [common.md](./common.md) - Standards projet
- [Go Debugging](https://go.dev/doc/diagnostics) - Guide Go
- [pprof](https://github.com/google/pprof) - Profiling
- [Delve](https://github.com/go-delve/delve) - Debugger
- [Makefile](../../Makefile) - Commandes

---

**Workflow** : Observer → Reproduire → Isoler → Analyser → Résoudre → Valider