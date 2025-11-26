# 🐛 Corriger un Bug (Fix Bug)

## Contexte

Projet TSD (Type System with Dependencies) - Moteur de règles RETE avec système de contraintes en Go.

Tu as identifié un bug dans le code et tu veux le corriger proprement en suivant une méthodologie rigoureuse pour éviter les régressions.

## Objectif

Corriger un bug de manière méthodique en :
- Reproduisant le bug de façon fiable
- Identifiant la cause racine
- Corrigeant sans introduire de régressions
- Ajoutant des tests de non-régression
- Documentant la correction

## ⚠️ RÈGLES STRICTES - CORRECTION DE BUG

### 🚫 INTERDICTIONS ABSOLUES

1. **CODE GOLANG** :
   - ❌ AUCUN HARDCODING introduit
   - ❌ AUCUN quick fix sans analyse
   - ❌ AUCUNE correction partielle
   - ❌ AUCUN code mort laissé
   - ✅ Code générique avec paramètres/interfaces
   - ✅ Constantes nommées pour toutes les valeurs
   - ✅ Respect strict Effective Go

2. **TESTS RETE** :
   - ❌ AUCUNE simulation de résultats
   - ❌ AUCUN test qui masque le bug
   - ✅ Extraction depuis réseau RETE réel uniquement
   - ✅ Tests de non-régression obligatoires
   - ✅ Tests reproduisant le bug avant correction

3. **MÉTHODOLOGIE** :
   - ❌ Pas de correction avant reproduction
   - ❌ Pas de commit sans tests
   - ✅ Root cause analysis obligatoire
   - ✅ Documentation de la correction
   - ✅ Validation complète

## Instructions

### PHASE 1 : REPRODUCTION (Isoler le Bug)

#### 1.1 Identifier le Bug

**Collecte d'informations** :

```
Bug ID : [Numéro si applicable]
Titre : [Description courte]
Sévérité : Critique / Majeure / Mineure / Triviale
Type : Fonctionnel / Performance / Sécurité / UI / Autre

Description :
[Description détaillée du comportement incorrect]

Comportement attendu :
[Ce qui devrait se passer]

Comportement observé :
[Ce qui se passe réellement]

Conditions de reproduction :
- Étape 1 : ...
- Étape 2 : ...
- Étape 3 : ...

Environnement :
- Go version : ...
- OS : ...
- Configuration : ...

Logs/Traces :
[Copier les messages d'erreur complets]
```

#### 1.2 Reproduire de Façon Fiable

**Créer un test de reproduction** :

```go
// rete/bug_XXXX_test.go
func TestBugXXXX_ReproduceIssue(t *testing.T) {
    t.Log("🐛 REPRODUCTION DU BUG #XXXX")
    t.Log("================================")
    
    // Arrange - Setup reproduisant le bug
    network := buildNetworkForBug()
    
    // Act - Action qui déclenche le bug
    err := network.SubmitFact(problematicFact)
    
    // Assert - Vérifier que le bug se produit
    // ⚠️ Ce test DOIT échouer avant la correction !
    if err == nil {
        t.Error("Bug ne se reproduit pas - vérifier setup")
    }
    
    t.Logf("✅ Bug reproduit : %v", err)
}
```

**Critères de reproduction** :
- [ ] Le bug se reproduit à 100%
- [ ] Les conditions minimales sont identifiées
- [ ] Le test échoue de manière prévisible
- [ ] Le test est déterministe (pas flaky)

#### 1.3 Isoler le Problème

**Réduire au cas minimal** :

```bash
# Identifier le fichier/fonction problématique
git log --oneline --all -- chemin/vers/fichier.go

# Tester avec données minimales
# Supprimer progressivement jusqu'au plus petit cas reproduisant le bug

# Vérifier avec git bisect si régression récente
git bisect start
git bisect bad HEAD
git bisect good v1.0.0
# Tester à chaque étape
```

**Questions à poser** :
- Dans quel module/fichier se trouve le bug ?
- Quelle fonction précise est impliquée ?
- Quelles données d'entrée déclenchent le bug ?
- Le bug est-il récent (régression) ou ancien ?
- Y a-t-il des bugs similaires ailleurs ?

### PHASE 2 : ANALYSE (Root Cause Analysis)

#### 2.1 Analyser la Cause Racine

**Techniques d'analyse** :

1. **Stack trace analysis** :
   ```bash
   # Examiner le stack trace complet
   go test -v ./... 2>&1 | tee error.log
   ```

2. **Debugging avec Delve** :
   ```bash
   dlv test ./rete -- -test.run TestBugXXXX
   # (dlv) break node_join.go:265
   # (dlv) continue
   # (dlv) print bindings
   # (dlv) print condition
   ```

3. **Logs détaillés** :
   ```go
   // Ajouter logs temporaires
   fmt.Printf("🔍 DEBUG: variable=%+v\n", variable)
   fmt.Printf("🔍 DEBUG: état=%#v\n", state)
   ```

4. **Analyse statique** :
   ```bash
   go vet ./...
   staticcheck ./...
   golangci-lint run --enable-all
   ```

**Les 5 Pourquoi** :

```
Pourquoi 1 : Le token ne se propage pas
  → Parce que evaluateCondition retourne false

Pourquoi 2 : evaluateCondition retourne false
  → Parce qu'une variable n'est pas disponible

Pourquoi 3 : La variable n'est pas disponible
  → Parce qu'on évalue la condition complète trop tôt

Pourquoi 4 : On évalue trop tôt
  → Parce qu'on ne vérifie pas les variables disponibles

Pourquoi 5 : On ne vérifie pas
  → CAUSE RACINE : Manque de vérification des variables disponibles
```

#### 2.2 Comprendre l'Impact

**Analyse d'impact** :

```bash
# Trouver code similaire (potentiellement même bug)
grep -r "pattern_problématique" --include="*.go" .

# Identifier les utilisateurs de la fonction
go list -f '{{.ImportPath}}' ./... | xargs grep "FunctionName"

# Vérifier tests impactés
go test ./... -v 2>&1 | grep FAIL
```

**Questions** :
- D'autres parties du code sont-elles affectées ?
- Y a-t-il des bugs similaires ?
- Quel est l'impact sur les utilisateurs ?
- Est-ce un problème de design ou d'implémentation ?

#### 2.3 Choisir la Stratégie de Correction

**Options** :

1. **Correction simple** :
   - Fix ponctuel dans la fonction
   - Pas d'impact API
   - Exemple : Ajouter vérification nil

2. **Correction avec refactoring** :
   - Amélioration du design
   - Impact local
   - Exemple : Extraire fonction de validation

3. **Correction architecturale** :
   - Changement de design
   - Impact large
   - Exemple : Changer ordre d'évaluation

**Choisir selon** :
- Complexité du bug
- Impact sur le code existant
- Urgence de la correction
- Risque de régression

### PHASE 3 : CORRECTION (Fix)

#### 3.1 Implémenter la Correction

**Processus** :

1. **Créer une branche dédiée** :
   ```bash
   git checkout -b fix/bug-XXXX-description
   ```

2. **Écrire le test qui doit passer** :
   ```go
   func TestBugXXXX_Fixed(t *testing.T) {
       t.Log("🔧 VALIDATION CORRECTION BUG #XXXX")
       t.Log("====================================")
       
       // Arrange
       network := buildNetworkForBug()
       
       // Act
       err := network.SubmitFact(problematicFact)
       
       // Assert - Après correction, ça doit passer
       if err != nil {
           t.Fatalf("❌ Bug non corrigé : %v", err)
       }
       
       // ✅ Vérifier résultat correct avec extraction réseau réel
       actualTokens := 0
       for _, terminal := range network.TerminalNodes {
           actualTokens += len(terminal.Memory.GetTokens())
       }
       
       if actualTokens == 0 {
           t.Error("❌ Aucun token créé après correction")
       }
       
       t.Logf("✅ Bug corrigé : %d tokens créés", actualTokens)
   }
   ```

3. **Implémenter la correction** :

   ```go
   // ❌ AVANT - Code buggé
   func evaluateCondition(bindings map[string]*Fact) bool {
       // Pas de vérification → Bug !
       return evaluator.Evaluate(condition)
   }
   
   // ✅ APRÈS - Code corrigé
   func evaluateCondition(bindings map[string]*Fact) bool {
       // ✅ Vérification ajoutée
       requiredVars := extractRequiredVariables(condition)
       
       for _, reqVar := range requiredVars {
           if _, exists := bindings[reqVar]; !exists {
               // Variable manquante, évaluation partielle
               return evaluatePartial(bindings, requiredVars)
           }
       }
       
       // Toutes variables disponibles, évaluation complète
       return evaluator.Evaluate(condition)
   }
   ```

**⚠️ Vérifier** :
- [ ] **AUCUN hardcoding** introduit
- [ ] **Code générique** maintenu
- [ ] Pas de quick fix, vraie correction
- [ ] Pas de code mort ajouté
- [ ] Style cohérent avec existant

#### 3.2 Ajouter Tests de Non-Régression

**Tests complets** :

```go
func TestBugXXXX_NonRegression(t *testing.T) {
    tests := []struct {
        name     string
        setup    func() *Network
        fact     *Fact
        validate func(*testing.T, *Network)
    }{
        {
            name: "cas_original_du_bug",
            setup: func() *Network {
                return buildNetworkForBug()
            },
            fact: problematicFact,
            validate: func(t *testing.T, net *Network) {
                // ✅ Extraction réseau réel
                count := 0
                for _, term := range net.TerminalNodes {
                    count += len(term.Memory.GetTokens())
                }
                if count == 0 {
                    t.Error("Régression : bug réapparu")
                }
            },
        },
        {
            name: "cas_limite_1",
            // ... autres cas
        },
        {
            name: "cas_normal_non_affecté",
            // Vérifier qu'on n'a rien cassé
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            network := tt.setup()
            network.SubmitFact(tt.fact)
            tt.validate(t, network)
        })
    }
}
```

#### 3.3 Mettre à Jour la Documentation

**GoDoc** :

```go
// EvaluateCondition évalue une condition avec gestion des variables manquantes.
//
// Fix: Correction du bug #XXXX où les variables manquantes causaient une erreur.
// Maintenant, effectue une évaluation partielle si certaines variables manquent.
//
// Paramètres:
//   - bindings: Map des variables liées aux faits
//
// Retourne:
//   - true si la condition est satisfaite (partiellement ou complètement)
//   - false sinon
//
// Version: Corrigé en v2.1.0
func EvaluateCondition(bindings map[string]*Fact) bool {
    // ...
}
```

**CHANGELOG.md** :

```markdown
## [2.1.0] - 2025-11-XX

### Fixed
- **[BUG #XXXX]** Correction erreur "variable non liée" dans evaluateCondition
  - Ajout vérification des variables disponibles
  - Implémentation évaluation partielle
  - Tests de non-régression ajoutés
  - Voir commit abc123f pour détails
```

### PHASE 4 : VALIDATION (Vérification)

#### 4.1 Tests Complets

**Checklist de validation** :

```bash
# 1. Test spécifique du bug
go test -v -run TestBugXXXX ./rete

# 2. Tests unitaires complets
go test ./...

# 3. Tests avec race detector
go test -race ./...

# 4. Tests d'intégration
make test-integration

# 5. Runner universel RETE
make rete-unified
# Doit afficher : 58/58 ✅

# 6. Formatage et lint
go fmt ./...
goimports -w .
go vet ./...
golangci-lint run

# 7. Tests de régression spécifiques
go test -run ".*Regression.*" ./...

# 8. Validation complète
make validate
```

**Tous doivent passer** ✅

#### 4.2 Tests de Régression

**Vérifier qu'on n'a rien cassé** :

```bash
# Tests existants passent toujours
go test ./... -count=1

# Benchmarks (pas de dégradation perf)
go test -bench=. ./rete -benchmem > after.txt
# Comparer avec before.txt

# Tests flaky
go test -count=10 ./...

# Tests avec différentes configurations
GOOS=linux go test ./...
GOOS=darwin go test ./...
```

#### 4.3 Revue de la Correction

**Auto-revue** :

- [ ] **Bug reproduit** avant correction
- [ ] **Cause racine** identifiée et documentée
- [ ] **Correction** implémentée sans hardcoding
- [ ] **Tests de non-régression** ajoutés
- [ ] **Tests RETE** avec extraction réseau réel
- [ ] **Documentation** mise à jour
- [ ] **CHANGELOG** mis à jour
- [ ] **Tous les tests** passent
- [ ] **Aucune régression** introduite
- [ ] **Code review** effectuée

## Critères de Succès

### ✅ Bug Corrigé

- [ ] Bug reproduit de façon fiable
- [ ] Cause racine identifiée
- [ ] Correction implémentée **sans hardcoding**
- [ ] Tests de non-régression ajoutés
- [ ] **Tests RETE avec extraction réseau réel**
- [ ] Tous les tests passent
- [ ] Aucune régression
- [ ] Documentation à jour

### ✅ Qualité Maintenue

- [ ] go vet : 0 erreur
- [ ] golangci-lint : 0 erreur
- [ ] Couverture maintenue ou améliorée
- [ ] Performance maintenue
- [ ] Pas de code mort
- [ ] Pas de duplication

### ✅ Traçabilité

- [ ] Bug ID référencé
- [ ] Commit message clair
- [ ] CHANGELOG mis à jour
- [ ] Tests documentés
- [ ] Cause racine documentée

## Format de Réponse

```
=== CORRECTION DE BUG ===

📋 IDENTIFICATION

Bug ID : #XXXX
Titre : [Titre du bug]
Sévérité : [Critique/Majeure/Mineure]
Type : [Fonctionnel/Performance/Sécurité]

Description :
[Description du bug]

Comportement attendu :
[Ce qui devrait se passer]

Comportement observé :
[Ce qui se passe]

🔬 REPRODUCTION

Test de reproduction : TestBugXXXX_ReproduceIssue
Conditions minimales :
  • Condition 1
  • Condition 2
  • Condition 3

✅ Reproduction fiable : Oui
✅ Cas minimal identifié : Oui

🔍 ANALYSE ROOT CAUSE

Fichier : rete/node_join.go
Fonction : evaluateCondition
Ligne : 265

Cause racine :
[Description de la cause racine]

Les 5 Pourquoi :
1. ...
2. ...
3. ...
4. ...
5. CAUSE RACINE : ...

Impact :
  • Modules affectés : X
  • Fonctions similaires : Y
  • Tests cassés : Z

🔧 CORRECTION

Stratégie : Correction simple avec vérification

Code modifié :
  ✅ rete/node_join.go (ligne 265-280)
  ✅ Ajout extractRequiredVariables()
  ✅ Ajout evaluatePartial()
  ⚠️ **VÉRIFIÉ** : Aucun hardcoding introduit
  ⚠️ **VÉRIFIÉ** : Code générique maintenu

Tests ajoutés :
  ✅ TestBugXXXX_Fixed
  ✅ TestBugXXXX_NonRegression
  ⚠️ **VÉRIFIÉ** : Extraction réseau RETE réel

Documentation :
  ✅ GoDoc mis à jour
  ✅ CHANGELOG mis à jour
  ✅ Commentaires explicatifs

✅ VALIDATION

Tests :
  ✅ TestBugXXXX_Fixed : PASS
  ✅ go test ./... : PASS
  ✅ go test -race ./... : PASS
  ✅ make test-integration : PASS
  ✅ make rete-unified : 58/58 ✅

Régression :
  ✅ Tests existants : PASS (X/X)
  ✅ Benchmarks : Pas de dégradation
  ✅ Tests flaky : PASS (10/10)

Qualité :
  ✅ go vet : 0 erreur
  ✅ golangci-lint : 0 erreur
  ✅ Couverture : Maintenue (X%)

📊 RÉSULTAT

Avant :
  • Bug se produit : 100% des cas
  • Tests échouent : X tests
  • Tokens créés : 0

Après :
  • Bug corrigé : 0% des cas
  • Tests passent : X/X ✅
  • Tokens créés : Y (correct)

🎯 VERDICT : BUG CORRIGÉ ✅

Commit : abc123f
Branche : fix/bug-XXXX-description
```

## Exemple d'Utilisation

```
J'ai identifié un bug où JoinNode crash avec un nil pointer quand
aucune variable n'est disponible dans les bindings.

Bug #1234
Sévérité : Majeure
Type : Fonctionnel

Utilise le prompt "fix-bug" pour :
1. Reproduire le bug de façon fiable
2. Identifier la cause racine
3. Corriger sans hardcoding
4. Ajouter tests de non-régression
5. Valider complètement
```

## Checklist de Correction

### Avant de Commencer
- [ ] Bug clairement identifié et décrit
- [ ] Environnement de reproduction disponible
- [ ] Branche dédiée créée
- [ ] Tests passent avant correction

### Pendant la Correction
- [ ] Bug reproduit de façon fiable
- [ ] Cause racine identifiée (5 Pourquoi)
- [ ] Test de reproduction écrit (doit échouer)
- [ ] Correction implémentée sans hardcoding
- [ ] Tests de non-régression ajoutés
- [ ] Tests RETE avec extraction réseau réel
- [ ] Documentation mise à jour

### Après la Correction
- [ ] **Test de reproduction passe** ✅
- [ ] **Tous les tests passent** ✅
- [ ] **Aucun hardcoding** introduit ✅
- [ ] **Tests RETE extraction réseau réel** ✅
- [ ] go vet et golangci-lint sans erreur ✅
- [ ] Aucune régression ✅
- [ ] CHANGELOG mis à jour ✅
- [ ] Commit message clair ✅
- [ ] Code review effectuée ✅

## Commandes Utiles

```bash
# Reproduction
go test -v -run TestBugXXXX ./rete

# Debugging
dlv test ./rete -- -test.run TestBugXXXX
go test -v ./rete 2>&1 | tee debug.log

# Git bisect (trouver régression)
git bisect start
git bisect bad HEAD
git bisect good v1.0.0

# Analyse statique
go vet ./...
staticcheck ./...
golangci-lint run --enable-all

# Validation
make test
make test-integration
make rete-unified
go test -race ./...
go test -count=10 ./...
```

## Bonnes Pratiques

1. **Toujours reproduire** avant de corriger
2. **Identifier la cause racine** (pas les symptômes)
3. **Corriger le problème**, pas le test
4. **Ajouter tests de non-régression** systématiquement
5. **Documenter** la correction et la cause
6. **Valider complètement** sans régression
7. **Respecter les règles** (pas de hardcoding, extraction RETE réelle)

## Anti-Patterns à Éviter

❌ **Quick fix sans analyse** :
```go
// Ne JAMAIS faire ça !
if variable == nil {
    return nil // Masque le vrai problème
}
```

❌ **Correction partielle** :
```go
// Corrige un cas mais pas les autres
if condition == "cas_spécifique" {
    // Fix temporaire
}
```

❌ **Test qui masque le bug** :
```go
// Test qui ne reproduit pas vraiment le bug
if err != nil {
    return nil // Ignore l'erreur
}
```

✅ **Correction propre** :
```go
// Identifie et corrige la cause racine
func evaluate(bindings map[string]*Fact) (bool, error) {
    if bindings == nil {
        return false, ErrNilBindings
    }
    // Vraie correction
}
```

## Ressources

- [Effective Debugging](https://www.oreilly.com/library/view/effective-debugging/9780134394909/)
- [Debugging with GDB](https://sourceware.org/gdb/documentation/)
- [Delve Debugger](https://github.com/go-delve/delve)
- [Go Testing Best Practices](https://go.dev/doc/tutorial/add-a-test)

---

**Version** : 1.0  
**Dernière mise à jour** : Novembre 2025  
**Type** : Correction de bug avec méthodologie rigoureuse