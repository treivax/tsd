# 🔄 Refactoriser du Code (Refactor)

## Contexte

Projet TSD (Type System with Dependencies) - Moteur de règles RETE avec système de contraintes en Go.

Tu veux refactoriser du code existant pour améliorer sa lisibilité, sa maintenabilité, ou sa structure, **sans changer son comportement fonctionnel**. Le refactoring est une amélioration interne du code qui ne doit pas affecter les résultats externes.

## Objectif

Améliorer la qualité interne du code (structure, lisibilité, maintenabilité) tout en préservant strictement son comportement et en garantissant qu'aucune régression n'est introduite.

## ⚠️ RÈGLES STRICTES - REFACTORING

### 🚫 INTERDICTIONS ABSOLUES

1. **AUCUN CHANGEMENT DE COMPORTEMENT** :
   - ❌ Ne pas modifier la logique fonctionnelle
   - ❌ Ne pas changer les résultats produits
   - ❌ Ne pas altérer les performances (sauf si c'est le but explicite)
   - ❌ Ne pas modifier l'API publique (sauf migration explicite)
   - ✅ Uniquement améliorer la structure interne
   - ✅ Tests existants doivent passer sans modification

2. **AUCUN HARDCODING INTRODUIT** :
   - ❌ Pas de nouvelles valeurs en dur
   - ❌ Pas de nouveaux "magic numbers"
   - ✅ Extraire les magic numbers existants en constantes
   - ✅ Remplacer hardcoding par paramètres/configuration

3. **CODE TOUJOURS GÉNÉRIQUE** :
   - ✅ Améliorer la généricité si possible
   - ✅ Introduire interfaces pour découplage
   - ✅ Paramétrer ce qui était hardcodé
   - ❌ Ne pas rendre le code plus spécifique

### ✅ OBJECTIFS DU REFACTORING

1. **Lisibilité** :
   - Noms de variables/fonctions plus clairs
   - Fonctions plus courtes et focalisées
   - Commentaires améliorés
   - Structure logique plus évidente

2. **Maintenabilité** :
   - Réduction de la complexité cyclomatique
   - Élimination de la duplication (DRY)
   - Meilleure séparation des responsabilités
   - Code plus testable

3. **Architecture** :
   - Meilleur découplage
   - Interfaces appropriées
   - Composition vs héritage
   - Patterns de conception appropriés

**Exemples** :

❌ **MAUVAIS - Change le comportement** :
```go
// Avant
func calculate(x, y int) int {
    return x + y
}

// Après - INTERDIT ! Change le comportement
func calculate(x, y int) int {
    return x + y + 1  // ❌ Modifie le résultat !
}
```

✅ **BON - Améliore sans changer** :
```go
// Avant
func calculate(x, y int) int {
    return x + y
}

// Après - OK ! Juste renommé pour clarté
func sum(operand1, operand2 int) int {
    return operand1 + operand2
}
```

## Instructions

### PHASE 1 : ANALYSE (Comprendre l'Existant)

#### 1.1 Identifier le Code à Refactoriser

**Localiser précisément** :
- Quel fichier : `rete/node_join.go`
- Quelle fonction/méthode : `evaluateJoinConditions`
- Quelle portée : Fonction entière, partie, ou fichier complet

**Analyser le problème** :
- Pourquoi refactoriser ? (complexité, duplication, lisibilité, etc.)
- Quels sont les "code smells" identifiés ?
- Quelle est la priorité ? (critique, important, nice-to-have)

#### 1.2 Comprendre le Comportement Actuel

**Documentation du comportement** :
```
1. Lire le code actuel
2. Identifier les entrées et sorties
3. Comprendre la logique métier
4. Noter les cas limites gérés
5. Identifier les dépendances
```

**Vérifier les tests existants** :
```bash
# Trouver les tests associés
grep -r "TestFonctionName" ./rete/*_test.go

# Exécuter les tests existants
go test -v -run TestFonctionName ./rete
```

**Établir la baseline** :
```bash
# Tous les tests doivent passer AVANT refactoring
make test
make rete-unified

# Baseline de performance si pertinent
go test -bench=BenchmarkFonctionName -benchmem ./rete
```

#### 1.3 Analyser la Qualité du Code

**Métriques actuelles** :
- Complexité cyclomatique (< 15 idéalement)
- Longueur de la fonction (< 50 lignes idéalement)
- Nombre de paramètres (< 5 idéalement)
- Duplication de code
- Couplage / Cohésion

**Outils d'analyse** :
```bash
# Analyse statique
go vet ./rete
golangci-lint run ./rete

# Couverture de tests
go test -cover ./rete
```

### PHASE 2 : PLANIFICATION (Stratégie de Refactoring)

#### 2.1 Choisir la Technique de Refactoring

**Techniques courantes** :

1. **Extract Function** :
   - Extraire une partie de fonction en fonction séparée
   - Réduire complexité et améliorer réutilisabilité

2. **Rename** :
   - Renommer variables/fonctions pour clarté
   - Utiliser noms descriptifs et idiomatiques

3. **Extract Constant** :
   - Remplacer magic numbers par constantes nommées
   - Centraliser les valeurs de configuration

4. **Simplify Conditional** :
   - Simplifier conditions complexes
   - Extraire conditions en fonctions nommées

5. **Remove Duplication** :
   - Identifier code dupliqué
   - Extraire en fonction commune

6. **Introduce Parameter** :
   - Remplacer hardcoding par paramètres
   - Rendre le code plus générique

7. **Decompose Complex Function** :
   - Décomposer fonction longue en plusieurs petites
   - Chaque fonction = une responsabilité

8. **Replace Magic Number** :
   - Remplacer nombres magiques par constantes

**Exemple de plan** :
```
Refactoring de evaluateJoinConditions :
1. Extraire validation des paramètres → validateJoinInputs()
2. Extraire logique de matching → matchTokenVariables()
3. Extraire création résultat → createResultToken()
4. Renommer variables : t1 → leftToken, t2 → rightToken
5. Remplacer magic number 100 → const MaxConditions = 100
```

#### 2.2 Planifier les Étapes Incrémentales

**Principe** : Refactoring par petites étapes vérifiables

```
Étape 1 : Extraire validateJoinInputs()
  → Exécuter tests → ✅ Commit

Étape 2 : Extraire matchTokenVariables()
  → Exécuter tests → ✅ Commit

Étape 3 : Extraire createResultToken()
  → Exécuter tests → ✅ Commit

Étape 4 : Renommer variables
  → Exécuter tests → ✅ Commit

Étape 5 : Extraire constantes
  → Exécuter tests → ✅ Commit
```

**Avantages** :
- Facile à reverter si problème
- Tests après chaque étape
- Historique git clair
- Réduction du risque

#### 2.3 Préparer les Tests de Non-Régression

**Tests à exécuter après chaque étape** :
```bash
# Tests unitaires du module
go test -v ./rete

# Tests d'intégration
make test-integration

# Runner universel
make rete-unified

# Benchmarks (si performance critique)
go test -bench=. -benchmem ./rete
```

### PHASE 3 : EXÉCUTION (Refactoring Incrémental)

#### 3.1 Refactoring Étape par Étape

**Pour chaque étape** :

1. **Faire UN changement** (atomic refactoring)
2. **Exécuter les tests** immédiatement
3. **Vérifier** que tout passe (tests + lint)
4. **Commit** avec message clair
5. **Passer à l'étape suivante**

**Template de commit** :
```
refactor(rete): extract validateJoinInputs from evaluateJoinConditions

- Extraire la validation des entrées en fonction séparée
- Améliore la lisibilité de evaluateJoinConditions
- Aucun changement de comportement
- Tests: ✅ go test ./rete
- Lint: ✅ golangci-lint run ./rete
```

#### 3.2 Exemple Concret : Extract Function

**Avant** :
```go
func evaluateJoinConditions(left, right *Token, conditions []Condition) (*Token, error) {
    // Validation (15 lignes)
    if left == nil {
        return nil, errors.New("left token is nil")
    }
    if right == nil {
        return nil, errors.New("right token is nil")
    }
    if len(conditions) == 0 {
        return nil, errors.New("no conditions")
    }
    // ... plus de validation
    
    // Matching (30 lignes)
    for _, cond := range conditions {
        // ... logique complexe de matching
    }
    
    // Construction résultat (20 lignes)
    result := &Token{}
    // ... construction
    
    return result, nil
}
```

**Après** :
```go
// validateJoinInputs vérifie la validité des entrées pour une jointure.
func validateJoinInputs(left, right *Token, conditions []Condition) error {
    if left == nil {
        return errors.New("left token is nil")
    }
    if right == nil {
        return errors.New("right token is nil")
    }
    if len(conditions) == 0 {
        return errors.New("no conditions")
    }
    return nil
}

// matchTokenVariables effectue le matching des variables entre tokens.
func matchTokenVariables(left, right *Token, conditions []Condition) (map[string]interface{}, error) {
    matches := make(map[string]interface{})
    for _, cond := range conditions {
        // ... logique de matching
    }
    return matches, nil
}

// createResultToken construit le token résultat d'une jointure.
func createResultToken(left, right *Token, matches map[string]interface{}) *Token {
    result := &Token{}
    // ... construction
    return result
}

// evaluateJoinConditions évalue les conditions de jointure entre deux tokens.
func evaluateJoinConditions(left, right *Token, conditions []Condition) (*Token, error) {
    // Validation
    if err := validateJoinInputs(left, right, conditions); err != nil {
        return nil, err
    }
    
    // Matching
    matches, err := matchTokenVariables(left, right, conditions)
    if err != nil {
        return nil, err
    }
    
    // Construction résultat
    result := createResultToken(left, right, matches)
    
    return result, nil
}
```

**Bénéfices** :
- ✅ Fonction principale de 65 lignes → 15 lignes
- ✅ Chaque fonction a une responsabilité claire
- ✅ Fonctions réutilisables séparément
- ✅ Plus facile à tester unitairement
- ✅ Complexité réduite
- ✅ **Comportement identique**

#### 3.3 Exemple : Replace Magic Number

**Avant** :
```go
func processTokens(tokens []*Token) error {
    if len(tokens) > 1000 {  // Magic number !
        return errors.New("too many tokens")
    }
    
    timeout := 30 * time.Second  // Magic number !
    
    // ...
}
```

**Après** :
```go
const (
    // MaxTokensPerBatch est le nombre maximum de tokens traités par lot.
    MaxTokensPerBatch = 1000
    
    // DefaultProcessTimeout est le timeout par défaut pour le traitement.
    DefaultProcessTimeout = 30 * time.Second
)

func processTokens(tokens []*Token) error {
    if len(tokens) > MaxTokensPerBatch {
        return fmt.Errorf("too many tokens (max: %d)", MaxTokensPerBatch)
    }
    
    timeout := DefaultProcessTimeout
    
    // ...
}
```

#### 3.4 Exemple : Simplify Conditional

**Avant** :
```go
func shouldPropagate(token *Token, node *Node) bool {
    if token != nil && node != nil && token.IsValid() && 
       node.IsActive() && !token.IsProcessed() && 
       len(token.Bindings) > 0 && node.AcceptsToken(token) {
        return true
    }
    return false
}
```

**Après** :
```go
// isValidForPropagation vérifie si un token est valide pour propagation.
func isValidForPropagation(token *Token) bool {
    return token != nil && 
           token.IsValid() && 
           !token.IsProcessed() && 
           len(token.Bindings) > 0
}

// canNodeAcceptToken vérifie si un nœud peut accepter un token.
func canNodeAcceptToken(node *Node, token *Token) bool {
    return node != nil && 
           node.IsActive() && 
           node.AcceptsToken(token)
}

// shouldPropagate détermine si un token doit être propagé à un nœud.
func shouldPropagate(token *Token, node *Node) bool {
    return isValidForPropagation(token) && 
           canNodeAcceptToken(node, token)
}
```

### PHASE 4 : VALIDATION (Garantir Non-Régression)

#### 4.1 Tests Complets

**Exécuter TOUS les tests** :
```bash
# Tests unitaires complets
go test -v ./...

# Tests avec couverture
go test -cover ./...

# Tests d'intégration
make test-integration

# Runner universel RETE
make rete-unified

# Validation complète
make validate
```

**Tous les tests doivent passer** :
- ✅ Tests unitaires : 100%
- ✅ Tests d'intégration : 100%
- ✅ Runner universel : 58/58
- ✅ Aucune régression

#### 4.2 Analyse Statique

**Vérifier la qualité** :
```bash
# Formatage
go fmt ./...
goimports -w .

# Analyse statique
go vet ./...
golangci-lint run ./...

# Vérifier qu'il n'y a pas de nouveaux warnings
```

#### 4.3 Vérification de Performance

**Si code critique pour performance** :
```bash
# Benchmarks avant/après
go test -bench=BenchmarkFonction -benchmem ./rete > before.txt
# ... refactoring ...
go test -bench=BenchmarkFonction -benchmem ./rete > after.txt

# Comparer
benchcmp before.txt after.txt
```

**Critères acceptables** :
- Performance identique (±5%)
- Ou amélioration
- ❌ Dégradation non acceptable (sauf justification)

#### 4.4 Revue de Code

**Auto-revue** :
```
✅ Le code est-il plus lisible ?
✅ Les fonctions sont-elles plus courtes ?
✅ Les noms sont-ils plus clairs ?
✅ La complexité est-elle réduite ?
✅ Le code est-il plus testable ?
✅ Moins de duplication ?
✅ Aucun hardcoding introduit ?
✅ Comportement strictement identique ?
✅ Tous les tests passent ?
✅ Aucun warning ajouté ?
```

## Critères de Succès

### ✅ Comportement Préservé

- [ ] Tous les tests existants passent **sans modification**
- [ ] Runner universel : 58/58 ✅
- [ ] Aucune régression fonctionnelle
- [ ] Performance identique ou améliorée
- [ ] API publique inchangée (ou migration explicite)

### ✅ Qualité Améliorée

- [ ] Lisibilité améliorée (noms clairs, structure logique)
- [ ] Complexité réduite (cyclomatique < 15)
- [ ] Fonctions plus courtes (< 50 lignes)
- [ ] Duplication éliminée
- [ ] Code plus testable
- [ ] Séparation responsabilités claire

### ✅ Standards Respectés

- [ ] Aucun hardcoding introduit
- [ ] Code générique maintenu/amélioré
- [ ] Constantes nommées pour valeurs
- [ ] go fmt et goimports appliqués
- [ ] go vet et golangci-lint sans erreur
- [ ] Commentaires GoDoc à jour

### ✅ Traçabilité

- [ ] Commits atomiques avec messages clairs
- [ ] Chaque étape testée et validée
- [ ] Historique git propre
- [ ] Documentation mise à jour si nécessaire

## Format de Réponse

```markdown
# 🔄 REFACTORING : [Nom du composant]

## 📋 Résumé
- **Fichier** : `rete/node_join.go`
- **Fonction** : `evaluateJoinConditions`
- **Problème** : Fonction trop longue (150 lignes), complexité 25
- **Objectif** : Décomposer en fonctions plus petites et claires

## 🎯 Plan de Refactoring

### Étapes planifiées
1. Extraire `validateJoinInputs()` (validation des entrées)
2. Extraire `matchTokenVariables()` (logique de matching)
3. Extraire `createResultToken()` (construction résultat)
4. Renommer variables pour clarté
5. Extraire constantes (remplacer magic numbers)

## 🔨 Exécution

### Étape 1 : Extract validateJoinInputs ✅

**Changement** :
```go
// Code extrait en fonction séparée
func validateJoinInputs(left, right *Token, conditions []Condition) error {
    // ... validation
}
```

**Validation** :
- ✅ Tests unitaires : PASS
- ✅ Tests intégration : PASS
- ✅ go vet : OK
- ✅ Commit : `refactor(rete): extract validateJoinInputs`

### Étape 2 : Extract matchTokenVariables ✅

**Changement** :
```go
func matchTokenVariables(left, right *Token, conditions []Condition) (map[string]interface{}, error) {
    // ... matching
}
```

**Validation** :
- ✅ Tests unitaires : PASS
- ✅ Tests intégration : PASS
- ✅ go vet : OK
- ✅ Commit : `refactor(rete): extract matchTokenVariables`

### Étape 3 : Extract createResultToken ✅

**Changement** :
```go
func createResultToken(left, right *Token, matches map[string]interface{}) *Token {
    // ... construction
}
```

**Validation** :
- ✅ Tests unitaires : PASS
- ✅ Tests intégration : PASS
- ✅ go vet : OK
- ✅ Commit : `refactor(rete): extract createResultToken`

### Étape 4 : Rename Variables ✅

**Changements** :
- `t1` → `leftToken`
- `t2` → `rightToken`
- `conds` → `conditions`
- `res` → `result`

**Validation** :
- ✅ Tests unitaires : PASS
- ✅ Commit : `refactor(rete): improve variable names in join evaluation`

### Étape 5 : Extract Constants ✅

**Changements** :
```go
const (
    MaxJoinConditions = 100
    DefaultMatchTimeout = 30 * time.Second
)
```

**Validation** :
- ✅ Tests unitaires : PASS
- ✅ Commit : `refactor(rete): extract magic numbers to named constants`

## 📊 Résultats

### Avant Refactoring
- **Lignes** : 150
- **Complexité** : 25
- **Fonctions** : 1 (monolithique)
- **Tests** : ✅ PASS

### Après Refactoring
- **Lignes** : 15 (fonction principale) + 3 fonctions helper
- **Complexité** : 5 (fonction principale), 3-8 (helpers)
- **Fonctions** : 4 (bien séparées)
- **Tests** : ✅ PASS (identique)

### Améliorations
- ✅ Lisibilité : ++++ (structure claire, noms descriptifs)
- ✅ Maintenabilité : +++ (fonctions courtes, responsabilités claires)
- ✅ Testabilité : +++ (fonctions isolées testables)
- ✅ Complexité : --- (réduite de 80%)
- ✅ Réutilisabilité : ++ (fonctions helper réutilisables)

## ✅ Validation Finale

### Tests Complets
```bash
$ make test
✅ Tests unitaires : PASS (234/234)

$ make test-integration
✅ Tests intégration : PASS (45/45)

$ make rete-unified
✅ Runner universel : PASS (58/58)

$ make validate
✅ Format : OK
✅ Lint : OK
✅ Build : OK
✅ Tests : OK
```

### Métriques Qualité
- ✅ Complexité cyclomatique : 5 (était 25)
- ✅ Lignes par fonction : 15 (était 150)
- ✅ Duplication : 0%
- ✅ go vet : Aucun warning
- ✅ golangci-lint : Aucun warning

### Performance
```
BenchmarkEvaluateJoinConditions-8
Avant : 1000 ns/op, 500 B/op, 10 allocs/op
Après : 1020 ns/op, 500 B/op, 10 allocs/op
Impact : +2% (négligeable, dans la marge d'erreur)
```

## 📝 Documentation Mise à Jour

- ✅ Commentaires GoDoc ajoutés pour nouvelles fonctions
- ✅ Exemples d'utilisation à jour
- ✅ Code auto-documenté (noms clairs)

## 🎓 Leçons Apprises

**Ce qui a bien marché** :
- Refactoring incrémental avec tests à chaque étape
- Commits atomiques (faciles à reverter si besoin)
- Extraction de fonctions avec responsabilités claires

**Points d'attention** :
- Toujours tester après CHAQUE modification
- Ne pas hésiter à reverter si doute
- Garder les commits petits et focalisés

## 📦 Fichiers Modifiés

```
rete/node_join.go          | 180 +++++++++++++++++++++++++-----------------
rete/node_join_test.go     | 15 ++++  (tests des fonctions helper)
```

Total : 2 fichiers modifiés, 195 insertions(+), 150 suppressions(-)

## ✅ Prêt pour Merge

- [x] Tous les tests passent
- [x] Aucune régression
- [x] Performance OK
- [x] Qualité améliorée
- [x] Documentation à jour
- [x] Commits propres
- [x] Prêt pour code review
```

## Exemple d'Utilisation

```
La fonction evaluateJoinConditions dans rete/node_join.go fait 150 lignes
avec une complexité cyclomatique de 25. Elle mélange validation, matching,
et construction du résultat.

Je veux la refactoriser pour :
1. Réduire la complexité
2. Améliorer la lisibilité
3. Faciliter les tests unitaires

Utilise le prompt "refactor".
```

## Checklist de Refactoring

### Avant de Commencer

- [ ] J'ai compris le comportement actuel du code
- [ ] J'ai identifié les tests existants
- [ ] Tous les tests passent actuellement (baseline)
- [ ] J'ai un plan de refactoring clair
- [ ] J'ai prévu les étapes incrémentales
- [ ] J'ai sauvegardé/committé l'état actuel

### Pendant le Refactoring

- [ ] Je fais UN changement à la fois
- [ ] J'exécute les tests après CHAQUE changement
- [ ] Je commit après chaque étape validée
- [ ] Je ne change pas le comportement
- [ ] Je n'introduis pas de hardcoding
- [ ] Je maintiens/améliore la généricité

### Après Chaque Étape

- [ ] Tests unitaires : PASS
- [ ] Tests intégration : PASS
- [ ] go vet : OK
- [ ] golangci-lint : OK
- [ ] Commit avec message clair
- [ ] Prêt pour étape suivante

### Validation Finale

- [ ] Tous les tests passent (make test)
- [ ] Runner universel OK (make rete-unified)
- [ ] Validation complète (make validate)
- [ ] Performance vérifiée (benchmarks)
- [ ] Qualité améliorée (complexité, lisibilité)
- [ ] Documentation à jour
- [ ] Historique git propre
- [ ] Prêt pour code review

## Commandes Utiles

```bash
# Tests après chaque étape
go test -v ./rete
go test -v -run TestFonctionSpecifique ./rete

# Validation complète
make test
make test-integration
make rete-unified
make validate

# Analyse qualité
go vet ./rete
golangci-lint run ./rete
gocyclo -over 15 rete/

# Performance
go test -bench=. -benchmem ./rete

# Couverture
go test -cover ./rete
go test -coverprofile=coverage.out ./rete
go tool cover -html=coverage.out

# Formatage
go fmt ./...
goimports -w .

# Git (commits atomiques)
git add rete/node_join.go
git commit -m "refactor(rete): extract validateJoinInputs"
git log --oneline
```

## Bonnes Pratiques

### Refactoring

1. **Incrémental** : Petites étapes vérifiables
2. **Testable** : Tests après chaque changement
3. **Réversible** : Commits atomiques faciles à reverter
4. **Documenté** : Messages de commit clairs
5. **Focalisé** : Un objectif à la fois

### Code

- **OBLIGATOIRE** : Aucun hardcoding introduit
- **OBLIGATOIRE** : Comportement strictement identique
- **OBLIGATOIRE** : Code générique maintenu/amélioré
- Extract Function : Fonctions < 50 lignes
- Extract Constant : Remplacer magic numbers
- Rename : Noms descriptifs et idiomatiques
- Simplify : Réduire complexité cyclomatique

### Tests

- Exécuter après CHAQUE modification
- Ne pas modifier les tests (sauf si changement API)
- 100% de réussite obligatoire
- Vérifier performance si code critique

## Anti-Patterns à Éviter

### ❌ Big Bang Refactoring
```
❌ Tout refactoriser d'un coup
✅ Refactoriser par petites étapes incrémentales
```

### ❌ Refactoring sans Tests
```
❌ Refactoriser sans exécuter les tests
✅ Tester après CHAQUE modification
```

### ❌ Changer le Comportement
```
❌ "Améliorer" la logique en refactorant
✅ Seulement améliorer la structure, pas la logique
```

### ❌ Introduire du Hardcoding
```
❌ Remplacer code générique par hardcoding
✅ Remplacer hardcoding par code générique
```

### ❌ Refactoring "Optimisation Prématurée"
```
❌ Compliquer le code pour optimiser sans besoin
✅ Simplifier d'abord, optimiser ensuite si nécessaire
```

### ❌ Refactoring sans Plan
```
❌ Commencer sans savoir où on va
✅ Planifier les étapes avant de commencer
```

## Techniques de Refactoring

### Extract Function
**Quand** : Fonction trop longue ou complexe  
**Comment** : Extraire partie logique en fonction séparée  
**Bénéfice** : Réduction complexité, réutilisabilité

### Extract Constant
**Quand** : Magic numbers ou strings  
**Comment** : Déclarer constante nommée  
**Bénéfice** : Lisibilité, maintenabilité

### Rename
**Quand** : Noms peu clairs  
**Comment** : Renommer avec noms descriptifs  
**Bénéfice** : Code auto-documenté

### Simplify Conditional
**Quand** : Conditions complexes  
**Comment** : Extraire en fonctions nommées  
**Bénéfice** : Lisibilité, testabilité

### Remove Duplication
**Quand** : Code dupliqué  
**Comment** : Extraire en fonction commune  
**Bénéfice** : DRY, maintenabilité

### Introduce Parameter
**Quand** : Valeurs hardcodées  
**Comment** : Ajouter paramètre  
**Bénéfice** : Généricité, réutilisabilité

### Decompose Function
**Quand** : Fonction multi-responsabilité  
**Comment** : Séparer en fonctions focalisées  
**Bénéfice** : SRP, testabilité

### Introduce Interface
**Quand** : Couplage fort  
**Comment** : Définir interface  
**Bénéfice** : Découplage, testabilité

## Outils Recommandés

### Analyse Statique
- `go vet` - Détection problèmes courants
- `golangci-lint` - Linter complet
- `gocyclo` - Complexité cyclomatique
- `goconst` - Détection strings/numbers dupliqués

### Tests
- `go test -cover` - Couverture
- `go test -race` - Détection race conditions
- `go test -bench` - Benchmarks

### Refactoring IDE
- GoLand / VS Code - Refactoring automatique
- `gofmt` - Formatage
- `goimports` - Imports

## Ressources

- [Makefile](../../Makefile) - Commandes disponibles
- [Refactoring: Improving the Design of Existing Code](https://martinfowler.com/books/refactoring.html) - Martin Fowler
- [Effective Go](https://golang.org/doc/effective_go.html) - Best practices Go
- [Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) - Go standards

---

**Version** : 1.0  
**Dernière mise à jour** : Novembre 2025  
**Mainteneur** : Équipe TSD