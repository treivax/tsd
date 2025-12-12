# Plan d'Action : Refactoring des Jointures Multi-Variables

**Objectif** : Résoudre le problème de perte de bindings dans les jointures à 3+ variables  
**Approche** : Refactoring du système de bindings avec architecture immuable  
**Principe clé** : **NE PAS faire cohabiter ancien et nouveau code** - Migration directe sans rétrocompatibilité

---

## 📊 État Actuel du Problème

### Symptômes
- Jointures à 3+ variables : token final ne contient que 2 variables au lieu de 3
- Erreur : `variable 'X' non trouvée (variables disponibles: [A B])`
- Exemple : `{u: User, o: Order, p: Product}` → token final = `[u, o]` au lieu de `[u, o, p]`

### Tests Affectés (3 tests E2E échouent)
1. `beta_join_complex.tsd` - Jointure User-Order-Product (3 variables)
2. `join_multi_variable_complex.tsd` - Jointure User-Team-Task (3 variables)
3. Un troisième test avec 3+ variables

### Cause Racine Suspectée
Le système actuel modifie les tokens de manière mutable, ce qui entraîne :
- Perte de bindings lors de la propagation entre niveaux de cascade
- État mutable difficile à tracer et déboguer
- Structure de données inadaptée pour les cascades multi-niveaux

---

## 🎯 Vision de la Solution

### Principes de Design

1. **Immutabilité** : Les bindings ne peuvent jamais être perdus une fois créés
2. **Chaîne de composition** : Chaque token porte la chaîne complète de ses bindings
3. **Scalabilité** : Support de N variables (pas de limite arbitraire)
4. **Simplicité** : Architecture claire et traçable
5. **Migration directe** : Remplacement complet sans code legacy

### Architecture Cible

```
TypeNode(User) ──┐
                 ├──> JoinNode1 
TypeNode(Order) ─┘     │
                       └─> Token {bindings: chain(u → o)}
                            │
TypeNode(Product) ──────────┤
                            │
                            └──> JoinNode2
                                   │
                                   └─> Token {bindings: chain(u → o → p)}
                                        │
                                        └──> TerminalNode ✅
```

---

## 📋 Plan d'Action Détaillé

### Structure des Prompts

Chaque prompt est conçu pour :
- Être exécutable en une session Zed unique (≤ 128K contexte)
- Avoir un livrable clair et testable
- Permettre une validation intermédiaire
- Respecter le principe : **pas de cohabitation ancien/nouveau**

---

## 🗂️ Phases et Prompts

### 📍 PHASE 1 : DIAGNOSTIC & DESIGN (2 sessions)

#### **Prompt 01 : Diagnostic Approfondi**
- **Durée** : 1-2 heures
- **Objectif** : Identifier le point exact de perte des bindings
- **Méthode** : Instrumentation temporaire + analyse de trace
- **Livrable** : `docs/architecture/BINDINGS_ANALYSIS.md`
- **Validation** : Cause racine identifiée avec preuves

#### **Prompt 02 : Spécification Technique du Nouveau Système**
- **Durée** : 2-3 heures
- **Objectif** : Concevoir l'architecture immuable complète
- **Livrable** : `docs/architecture/BINDINGS_DESIGN.md`
- **Contenu** :
  - Structures de données (BindingChain, ImmutableToken)
  - Interfaces et contrats
  - Stratégie de migration (remplacement, pas cohabitation)
  - Plan de test détaillé
- **Validation** : Design reviewé et approuvé

---

### 🔧 PHASE 2 : STRUCTURES IMMUABLES (2 sessions)

#### **Prompt 03 : Implémentation de BindingChain**
- **Durée** : 2-3 heures
- **Fichiers créés** :
  - `rete/binding_chain.go` (nouvelle structure)
  - `rete/binding_chain_test.go` (tests unitaires complets)
- **Fonctionnalités** :
  ```go
  type BindingChain struct {
      Variable string
      Fact     *Fact
      Parent   *BindingChain  // Immutable chain
  }
  
  // Méthodes :
  // - Get(variable) *Fact
  // - Has(variable) bool
  // - Add(variable, fact) *BindingChain  // Retourne NOUVELLE chaîne
  // - ToMap() map[string]*Fact
  // - Variables() []string
  // - Len() int
  ```
- **Tests** : Couverture > 95%
- **Validation** : `go test ./rete/binding_chain_test.go` passe

#### **Prompt 04 : Refactoring de Token → ImmutableToken**
- **Durée** : 2-4 heures
- **Objectif** : Remplacer complètement l'ancienne structure Token
- **Fichiers modifiés** :
  - `rete/fact_token.go` - **REMPLACEMENT COMPLET** de Token
  - Tous les fichiers utilisant Token (migration directe)
- **Nouvelle structure** :
  ```go
  type Token struct {  // Même nom, nouvelle implémentation
      ID           string
      Facts        []*Fact
      Bindings     *BindingChain  // Au lieu de map[string]*Fact
      NodeID       string
      Metadata     TokenMetadata
  }
  ```
- **Stratégie** : 
  - Renommer temporairement l'ancien Token en TokenOld
  - Créer la nouvelle implémentation
  - Fixer toutes les erreurs de compilation d'un coup
  - Supprimer TokenOld
- **Validation** : Code compile, tests unitaires de Token passent

---

### 🔗 PHASE 3 : REFACTORING DES JOINTURES (3 sessions)

#### **Prompt 05 : Refactoring JoinNode - Partie 1 (performJoinWithTokens)**
- **Durée** : 2-3 heures
- **Fichier** : `rete/node_join.go`
- **Objectif** : Réécrire la logique de jointure pour utiliser BindingChain
- **Méthode** :
  ```go
  func (jn *JoinNode) performJoinWithTokens(token1, token2 *Token) *Token {
      // Ancienne version utilisait map merge
      // Nouvelle version : composition de chaînes
      newChain := token1.Bindings
      for _, v := range token2.Bindings.Variables() {
          fact := token2.Bindings.Get(v)
          newChain = newChain.Add(v, fact)
      }
      return &Token{
          Bindings: newChain,  // Chaîne immuable complète
          // ...
      }
  }
  ```
- **Validation** : Tests unitaires de JoinNode passent

#### **Prompt 06 : Refactoring JoinNode - Partie 2 (Activation)**
- **Durée** : 2-3 heures
- **Fichier** : `rete/node_join.go`
- **Objectif** : Réécrire ActivateLeft/ActivateRight pour bindings immuables
- **Focus** :
  - Gestion des mémoires Left/Right avec BindingChain
  - Propagation correcte des tokens composés
  - getVariableForFact adapté
- **Validation** : Tests d'intégration pour jointures 2 variables passent

#### **Prompt 07 : Refactoring BetaChainBuilder**
- **Durée** : 3-4 heures
- **Fichiers** : 
  - `rete/builder_beta_chain.go`
  - `rete/builder_join_rules_cascade.go`
- **Objectif** : Assurer que les cascades sont construites correctement
- **Vérifications critiques** :
  - AllVariables contient TOUTES les variables cumulées à chaque niveau
  - RightVariables contient la nouvelle variable à chaque cascade
  - LeftVariables contient toutes les variables des niveaux précédents
- **Validation** : Construction de cascade pour 3+ variables produit la bonne structure

---

### 🎬 PHASE 4 : ACTIONS ET TERMINAL (1 session)

#### **Prompt 08 : Refactoring ExecutionContext et ActionExecutor**
- **Durée** : 2-3 heures
- **Fichiers** :
  - `rete/action_executor_context.go`
  - `rete/action_executor_evaluation.go`
  - `rete/node_terminal.go`
- **Objectif** : Adapter l'exécution d'actions aux bindings immuables
- **Changements** :
  ```go
  // Avant :
  ctx.varCache = token.Bindings  // map[string]*Fact
  
  // Après :
  ctx.bindingChain = token.Bindings  // *BindingChain
  
  // Résolution de variable :
  func (ctx *ExecutionContext) resolveVariable(name string) (*Fact, error) {
      if ctx.bindingChain.Has(name) {
          return ctx.bindingChain.Get(name), nil
      }
      // Erreur avec liste complète des variables disponibles
      return nil, fmt.Errorf("variable '%s' non trouvée (disponibles: %v)",
          name, ctx.bindingChain.Variables())
  }
  ```
- **Validation** : Actions sont exécutées avec tous les bindings

---

### 🧪 PHASE 5 : TESTS COMPLETS (2 sessions)

#### **Prompt 09 : Tests Unitaires pour Cascades Multi-Variables**
- **Durée** : 2-3 heures
- **Fichier créé** : `rete/node_join_cascade_test.go`
- **Tests à implémenter** :
  1. **Test_JoinCascade_2Variables** : Régression (doit continuer à passer)
  2. **Test_JoinCascade_3Variables** : Cas principal
     - Configuration : User, Order, Product
     - Vérification : Token final contient [u, o, p]
     - Variations : Différents ordres d'arrivée des faits
  3. **Test_JoinCascade_4Variables** : Scalabilité
  4. **Test_JoinCascade_NVariables** : Test paramétrique (N=2 à 10)
- **Assertions** :
  - Nombre de bindings = nombre de variables attendu
  - Chaque variable est présente
  - Chaque binding pointe vers le bon fait
- **Validation** : Tous les tests passent

#### **Prompt 10 : Tests d'Intégration et E2E**
- **Durée** : 2-3 heures
- **Objectif** : Vérifier que TOUS les tests existants passent
- **Commandes** :
  ```bash
  make test-unit           # Tous les tests unitaires
  make test-integration    # Tests d'intégration
  make test-e2e           # Tests E2E - les 3 échouant doivent passer
  ```
- **Focus** :
  - `beta_join_complex.tsd` ✅
  - `join_multi_variable_complex.tsd` ✅
  - Tous les autres tests continuent de passer
- **Debugging** : Si échecs, ajout de logs et correction
- **Validation** : 83/83 tests E2E passent (100%)

---

### 🎯 PHASE 6 : FINALISATION (2 sessions)

#### **Prompt 11 : Optimisation et Performance**
- **Durée** : 2-3 heures
- **Objectif** : S'assurer qu'il n'y a pas de régression de performance
- **Actions** :
  1. Créer des benchmarks : `rete/node_join_benchmark_test.go`
  2. Benchmarks pour :
     - Jointure 2 variables (baseline)
     - Jointure 3 variables
     - Jointure N variables
     - Création de BindingChain
     - Recherche dans BindingChain
  3. Comparer avec les performances théoriques
  4. Optimiser si nécessaire (caching, indexation)
- **Critère** : Overhead < 10% pour jointures 2 variables
- **Validation** : `go test -bench=. ./rete/` montre des résultats acceptables

#### **Prompt 12 : Documentation et Cleanup Final**
- **Durée** : 2-3 heures
- **Objectif** : Finaliser la documentation et nettoyer le code
- **Livrables** :
  1. **Documentation technique** :
     - Mise à jour de `docs/architecture/RETE.md`
     - Documentation de BindingChain et Token immuable
     - Exemples d'utilisation
  2. **GoDoc** :
     - Commenter toutes les fonctions exportées
     - Exemples de code dans les docs
  3. **Cleanup** :
     - Supprimer tout code temporaire/debug
     - Supprimer fichiers obsolètes
     - Vérifier qu'aucun "TODO" ou "FIXME" ne reste
  4. **CHANGELOG.md** :
     - Ajouter entrée pour ce refactoring majeur
  5. **Validation finale** :
     ```bash
     make validate  # Format + Lint + Build + Tests complets
     ```
- **Validation** : `make validate` passe sans erreur ni warning

---

## 🔄 Ordre d'Exécution Strict

```
PHASE 1: Diagnostic & Design
  01_diagnostic.md     → BINDINGS_ANALYSIS.md
  02_design.md         → BINDINGS_DESIGN.md
  
PHASE 2: Structures Immuables
  03_binding_chain.md  → binding_chain.go + tests
  04_token_refactor.md → Migration complète de Token
  
PHASE 3: Jointures
  05_join_perform.md   → performJoinWithTokens refactoré
  06_join_activate.md  → ActivateLeft/Right refactoré
  07_chain_builder.md  → BetaChainBuilder refactoré
  
PHASE 4: Actions
  08_actions.md        → ExecutionContext + Terminal refactorés
  
PHASE 5: Tests
  09_unit_tests.md     → Tests cascade + validation
  10_e2e_tests.md      → Validation E2E complète
  
PHASE 6: Finalisation
  11_performance.md    → Benchmarks + optimisations
  12_documentation.md  → Docs + cleanup final
```

**Durée totale estimée** : 8-12 jours de travail

---

## 📦 Livrables par Phase

### Documentation
- `docs/architecture/BINDINGS_ANALYSIS.md` - Analyse du problème
- `docs/architecture/BINDINGS_DESIGN.md` - Spécification technique
- `docs/architecture/RETE.md` - Mise à jour avec nouveau système

### Code (Nouveaux fichiers)
- `rete/binding_chain.go` - Structure immuable de bindings
- `rete/binding_chain_test.go` - Tests unitaires BindingChain
- `rete/node_join_cascade_test.go` - Tests des cascades
- `rete/node_join_benchmark_test.go` - Benchmarks de performance

### Code (Fichiers modifiés - Remplacement complet)
- `rete/fact_token.go` - Token avec BindingChain
- `rete/node_join.go` - JoinNode avec bindings immuables
- `rete/builder_beta_chain.go` - Construction correcte des cascades
- `rete/builder_join_rules_cascade.go` - Patterns de cascade
- `rete/action_executor_context.go` - Résolution via BindingChain
- `rete/action_executor_evaluation.go` - Évaluation avec chaîne
- `rete/node_terminal.go` - Terminal avec bindings immuables

---

## ⚠️ Contraintes Critiques

### 1. Pas de Cohabitation Ancien/Nouveau Code

**❌ INTERDIT** :
```go
// Ne JAMAIS faire ceci :
type Token struct {
    Bindings map[string]*Fact  // Ancien
    BindingChainNew *BindingChain  // Nouveau
}

// Ou ceci :
if useNewSystem {
    // nouveau code
} else {
    // ancien code
}
```

**✅ OBLIGATOIRE** :
```go
// Remplacement direct :
type Token struct {
    Bindings *BindingChain  // Seule version
}
```

**Stratégie de migration (Prompt 04)** :
1. Créer une branche de travail
2. Renommer Token → TokenOld (commenté)
3. Créer nouvelle implémentation de Token
4. Fixer TOUTES les erreurs de compilation
5. Supprimer TokenOld
6. Valider tests
7. Commit

### 2. Tests à Chaque Étape

**Après chaque prompt, ces commandes DOIVENT passer** :
```bash
go build ./...           # Compilation sans erreur
go test ./rete/...       # Tests unitaires du module modifié
make test-unit           # Tous les tests unitaires (après Phase 2)
```

### 3. Limite de Contexte Zed (128K)

**Chaque prompt doit** :
- Lire max 15-20 fichiers
- Se concentrer sur 1-3 fichiers à modifier
- Références claires aux prompts précédents
- Pas de code dupliqué entre prompts

### 4. Backward Compatibility : NON

Ce refactoring **casse volontairement** l'ancienne API interne :
- Les tests doivent être adaptés
- Le code appelant doit être mis à jour
- Pas de support de l'ancienne structure Token
- Suppression complète de l'ancien code

**Justification** : C'est un refactoring interne du moteur RETE, pas une API publique.

---

## 🎓 Concepts Clés

### 1. BindingChain - Chaîne Immuable

**Principe** : Structural sharing - partage de structure

```go
type BindingChain struct {
    Variable string          // Nom de la variable (ex: "u")
    Fact     *Fact          // Fait lié (ex: User{id: 1})
    Parent   *BindingChain  // Chaîne parente (nil si racine)
}

// Exemple d'utilisation :
chain1 := &BindingChain{Variable: "u", Fact: userFact, Parent: nil}
chain2 := &BindingChain{Variable: "o", Fact: orderFact, Parent: chain1}
chain3 := &BindingChain{Variable: "p", Fact: productFact, Parent: chain2}

// chain3 contient TOUTES les variables : u, o, p
// Recherche : O(n) où n = nombre de variables (acceptable pour n < 10)
```

**Avantages** :
- Impossible de perdre un binding
- Pas de copie de données (pointeurs partagés)
- Thread-safe par nature (immutable)
- Traçable (on peut remonter la chaîne)

**Inconvénient** :
- Recherche O(n) au lieu de O(1) pour map
- **Mitigation** : Cache optionnel dans Token si n > seuil

### 2. Token Immuable

```go
type Token struct {
    ID       string
    Facts    []*Fact
    Bindings *BindingChain  // Chaîne complète de bindings
    NodeID   string
    Metadata TokenMetadata   // Pour debugging
}

type TokenMetadata struct {
    CreatedAt    time.Time
    CreatedBy    string      // Node ID
    JoinLevel    int         // Profondeur de cascade
    VariableList []string    // Cache des variables
}
```

### 3. Composition de Chaînes

**Opération fondamentale** : Ajouter un binding

```go
func (bc *BindingChain) Add(variable string, fact *Fact) *BindingChain {
    // Retourne une NOUVELLE chaîne, ne modifie pas l'existante
    return &BindingChain{
        Variable: variable,
        Fact:     fact,
        Parent:   bc,  // Pointe vers la chaîne existante
    }
}

// Utilisation dans performJoinWithTokens :
func (jn *JoinNode) performJoinWithTokens(left, right *Token) *Token {
    // Partir de la chaîne gauche (bindings accumulés)
    newChain := left.Bindings
    
    // Ajouter les bindings du côté droit (nouveau fait)
    for _, v := range right.Bindings.Variables() {
        fact := right.Bindings.Get(v)
        newChain = newChain.Add(v, fact)  // Composition
    }
    
    return &Token{
        Bindings: newChain,  // Chaîne complète
        Facts:    append(left.Facts, right.Facts...),
        // ...
    }
}
```

---

## 🧪 Stratégie de Tests

### Tests Unitaires (Par Composant)

1. **BindingChain** (`binding_chain_test.go`)
   - Création de chaîne vide
   - Ajout de bindings
   - Recherche (Get, Has)
   - Conversion (ToMap, Variables)
   - Edge cases (nil, variable inexistante)

2. **Token** (`fact_token_test.go`)
   - Création avec BindingChain
   - Métadonnées
   - Sérialisation/désérialisation si applicable

3. **JoinNode** (`node_join_test.go`)
   - performJoinWithTokens avec BindingChain
   - ActivateLeft/Right avec bindings immuables
   - Mémoires Left/Right

4. **Cascades** (`node_join_cascade_test.go`)
   - 2 variables (régression)
   - 3 variables (cas principal)
   - N variables (scalabilité)

### Tests d'Intégration

1. **BetaChainBuilder** (dans `builder_beta_chain_test.go` existant)
   - Construction de cascades 2, 3, 4+ variables
   - Vérification de la structure créée

2. **End-to-End RETE** (dans `tests/integration/`)
   - Soumission de faits dans différents ordres
   - Vérification que les actions reçoivent tous les bindings

### Tests E2E (Fixtures)

- Tous les tests existants doivent continuer à passer
- Les 3 tests échouant doivent maintenant passer :
  - `beta_join_complex.tsd` ✅
  - `join_multi_variable_complex.tsd` ✅
  - Troisième test avec 3+ variables ✅

---

## 📈 Métriques de Succès

### Critères de Réussite Fonctionnels

✅ **Correction du Bug**
- [ ] Les 3 tests E2E échouant passent maintenant
- [ ] Tous les tests existants continuent de passer (non-régression)
- [ ] Nouveau test avec 4+ variables passe

✅ **Qualité du Code**
- [ ] Aucun binding perdu dans les cascades (prouvé par tests)
- [ ] Code immuable et thread-safe
- [ ] Couverture de tests > 80% sur le nouveau code
- [ ] Pas de code mort (ancien code supprimé)

✅ **Documentation**
- [ ] Architecture documentée dans `docs/architecture/`
- [ ] GoDoc pour toutes les fonctions exportées
- [ ] Exemples d'utilisation clairs

### Critères de Réussite Performance

✅ **Performance Acceptable**
- [ ] Aucune régression sur jointures 2 variables
- [ ] Overhead < 10% pour jointures 3 variables
- [ ] Scalabilité vérifiée jusqu'à N=10 variables
- [ ] Pas de memory leaks
- [ ] Benchmarks documentés

### Critères de Réussite Maintenabilité

✅ **Code Clean**
- [ ] Respect de Effective Go
- [ ] Complexité cyclomatique < 15
- [ ] Pas de duplication de code
- [ ] Pas de "magic numbers"
- [ ] Messages d'erreur clairs

✅ **Validation Automatisée**
- [ ] `make validate` passe sans erreur ni warning
- [ ] `go vet`, `staticcheck`, `errcheck` passent
- [ ] `gofmt`, `goimports` appliqués

---

## 🔍 Points d'Attention Critiques

### Risques Identifiés

1. **Complexité du Refactoring**
   - Ce refactoring touche le cœur du moteur RETE
   - **Mitigation** : Prompts séquentiels avec validation à chaque étape
   - **Rollback** : Chaque prompt est une branche Git séparée

2. **Performance**
   - BindingChain est O(n) vs map O(1)
   - **Mitigation** : n est petit (< 10 typiquement), cache optionnel
   - **Validation** : Benchmarks à chaque étape

3. **Breaking Changes**
   - L'API interne change complètement
   - **Mitigation** : C'est assumé, pas d'API publique cassée
   - **Communication** : Documentation claire des changements

4. **Tests**
   - Beaucoup de tests doivent être adaptés
   - **Mitigation** : Tests adaptés au fur et à mesure
   - **Validation** : Tests passent à chaque prompt

---

## 🚀 Démarrage

### Pré-requis

1. **Code à jour**
   ```bash
   git checkout main
   git pull
   go mod tidy
   ```

2. **Tests baseline**
   ```bash
   make test-complete  # Vérifier l'état actuel
   ```

3. **Documentation lue**
   - [ ] `RESOLUTION_TESTS_E2E.md`
   - [ ] `.github/prompts/common.md`
   - [ ] `docs/architecture/RETE.md` (si existe)

### Commencer

**Exécuter Prompt 01** : `tsd/scripts/multi-jointures/01_diagnostic.md`

---

## 📚 Références

- **Immutable Data Structures** : Okasaki, Chris. "Purely Functional Data Structures" (1998)
- **RETE Algorithm** : Forgy, Charles. "Rete: A Fast Algorithm for the Many Pattern/Many Object Pattern Match Problem" (1982)
- **Go Idioms** : Effective Go - https://go.dev/doc/effective_go

---

**Date de création** : 2025-01-XX  
**Version** : 2.0  
**Auteur** : Plan généré suite à l'analyse des échecs E2E  
**Principe directeur** : Migration directe sans cohabitation ancien/nouveau code