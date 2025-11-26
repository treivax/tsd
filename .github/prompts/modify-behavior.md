# 🔄 Modifier un Comportement ou une Fonctionnalité

## Contexte

Projet TSD (Type System with Dependencies) - Moteur de règles RETE avec système de contraintes en Go.

Tu veux modifier le comportement d'une fonctionnalité existante, ajuster son fonctionnement, changer ses paramètres, ou adapter son implémentation sans créer de nouvelle feature.

## Objectif

Modifier proprement une fonctionnalité existante en :
- Préservant la compatibilité quand possible
- Mettant à jour tous les composants impactés
- Maintenant la qualité et les tests
- Documentant les changements

## 📄 RÈGLES DE LICENCE ET COPYRIGHT - OBLIGATOIRE

### 🔒 Vérification de Compatibilité de Licence

**SI la modification nécessite du code externe ou une nouvelle bibliothèque** :

1. **Vérifier la licence** :
   - ✅ Licences permissives acceptées : MIT, BSD, Apache-2.0, ISC
   - ⚠️ Licences à éviter : GPL, AGPL, LGPL (copyleft)
   - ❌ Code sans licence = NE PAS UTILISER
   - ❌ Code propriétaire = NE PAS UTILISER

2. **Documenter l'origine** :
   - Si code inspiré/adapté : ajouter commentaire avec source
   - Si bibliothèque tierce : mettre à jour `go.mod` et `THIRD_PARTY_LICENSES.md`
   - Si algorithme connu : citer la référence académique

### 📝 En-tête de Copyright OBLIGATOIRE

**SI création de nouveaux fichiers durant la modification** :

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package [nom_du_package]
```

**VÉRIFICATION** :
- ✅ Tous les nouveaux fichiers .go ont l'en-tête de copyright
- ✅ Les fichiers existants conservent leur en-tête
- ✅ Aucun code externe non vérifié n'est introduit

### ⚠️ INTERDICTIONS STRICTES

- ❌ **Ne JAMAIS copier du code** sans vérifier la licence
- ❌ **Ne JAMAIS utiliser de code GPL/AGPL** (incompatible avec MIT)
- ❌ **Ne JAMAIS omettre les en-têtes de copyright** dans les nouveaux fichiers
- ✅ **TOUJOURS écrire du code original** lors de modifications

## ⚠️ RÈGLES STRICTES - MODIFICATION DE CODE

### 🚫 INTERDICTIONS ABSOLUES

1. **CODE GOLANG** :
   - ❌ AUCUN HARDCODING introduit
   - ❌ AUCUNE duplication de code
   - ❌ AUCUN code mort laissé
   - ❌ AUCUNE dégradation de qualité
   - ✅ Code générique avec paramètres/interfaces
   - ✅ Constantes nommées pour toutes les valeurs
   - ✅ Respect strict Effective Go

2. **TESTS RETE** :
   - ❌ AUCUNE simulation de résultats
   - ❌ AUCUN test cassé ou obsolète
   - ✅ Extraction depuis réseau RETE réel uniquement
   - ✅ Tous les tests mis à jour
   - ✅ Nouveaux tests si comportement change

3. **COMPATIBILITÉ** :
   - ⚠️ Breaking changes identifiés et documentés
   - ✅ Migration path fourni si nécessaire
   - ✅ Deprecated markers si API change
   - ✅ Rétrocompatibilité préservée si possible

### ✅ OBLIGATIONS

1. **Mise à jour complète** :
   - ✅ Tous les tests impactés mis à jour
   - ✅ Documentation mise à jour
   - ✅ Exemples mis à jour
   - ✅ Prompts mis à jour si nécessaire

2. **Validation** :
   - ✅ Tous les tests passent (unit + integration + rete-unified)
   - ✅ go vet et golangci-lint sans erreur
   - ✅ Aucune régression introduite
   - ✅ Couverture maintenue ou améliorée

## Instructions

### PHASE 1 : ANALYSE DE L'IMPACT (Avant Modification)

#### 1.1 Identifier la Fonctionnalité

**Définir clairement** :
- **Fonctionnalité** : Quelle feature modifier ?
- **Localisation** : Module(s), fichier(s), fonction(s) concernés
- **Comportement actuel** : Comment ça marche maintenant ?
- **Comportement souhaité** : Comment ça devrait marcher ?
- **Raison** : Pourquoi modifier ?

**Exemple** :
```
Fonctionnalité : Évaluation des conditions dans JoinNode
Localisation : rete/node_join.go, evaluateJoinConditions()
Actuel : Évalue toute la condition même si variables manquantes
Souhaité : Évaluation partielle intelligente
Raison : Éviter erreurs "variable non liée"
```

#### 1.2 Analyser l'Impact

**Carte d'impact** :

```bash
# Trouver toutes les utilisations
grep -r "FonctionCible" --include="*.go" .

# Trouver les tests
grep -r "TestFonctionCible\|Test.*FonctionCible" --include="*_test.go" .

# Trouver les imports du package
go list -f '{{.ImportPath}} {{.Imports}}' ./... | grep "package/cible"

# Vérifier l'API publique
go doc -all package/cible | grep "^func"
```

**Questions à se poser** :
- Combien de fichiers sont impactés ?
- Y a-t-il des tests existants ?
- Est-ce une API publique ou interne ?
- D'autres fonctionnalités dépendent-elles de celle-ci ?
- Y a-t-il des breaking changes ?

#### 1.3 Identifier les Dépendances

**Lister les composants impactés** :

1. **Code** :
   - Fonction(s) à modifier
   - Appelants de ces fonctions
   - Structures de données utilisées
   - Interfaces implémentées

2. **Tests** :
   - Tests unitaires directs
   - Tests d'intégration
   - Tests RETE concernés
   - Benchmarks

3. **Documentation** :
   - GoDoc
   - README
   - Fichiers .constraint d'exemple
   - Commentaires inline

4. **Prompts** :
   - Prompts décrivant cette fonctionnalité
   - Exemples dans les prompts

#### 1.4 Planifier la Migration

**Si breaking change** :

1. **Stratégie** :
   - Deprecated l'ancienne version
   - Créer nouvelle version en parallèle
   - Période de transition définie
   - Ou : modification in-place si interne

2. **Migration path** :
   ```go
   // Ancienne version (deprecated)
   // Deprecated: Use NewFunction instead
   func OldFunction(param string) error {
       return NewFunction(param, defaultValue)
   }
   
   // Nouvelle version
   func NewFunction(param string, newParam int) error {
       // Nouvelle implémentation
   }
   ```

### PHASE 2 : MODIFICATION (Action)

#### 2.1 Préparer l'Environnement

**Setup** :

```bash
# Créer une branche dédiée
git checkout -b modify-feature-name

# S'assurer que tout passe avant modification
make test
make rete-unified

# Commit de référence
git add .
git commit -m "Baseline avant modification de [feature]"
```

#### 2.2 Modifier le Code

**Processus** :

1. **Lire et comprendre le code actuel** :
   ```bash
   # Examiner la fonction
   cat rete/node_join.go | grep -A 50 "func evaluateJoinConditions"
   
   # Voir l'historique
   git log -p --follow -- rete/node_join.go
   ```

2. **Modifier le code** :
   - ⚠️ **VÉRIFIER** : Aucun hardcoding introduit
   - ⚠️ **VÉRIFIER** : Code générique maintenu
   - ⚠️ **VÉRIFIER** : Pas de duplication
   - Modifier la logique
   - Ajouter paramètres si nécessaire
   - Utiliser constantes nommées

3. **Exemple de modification** :

   ```go
   // ❌ AVANT - Problématique
   func evaluateCondition(bindings map[string]*Fact) bool {
       timeout := 30  // Hardcodé !
       result := evaluator.Evaluate(condition)
       return result
   }
   
   // ✅ APRÈS - Amélioré
   const DefaultEvaluationTimeout = 30 * time.Second
   
   func evaluateCondition(bindings map[string]*Fact, timeout time.Duration) bool {
       // Vérifier variables disponibles
       requiredVars := extractRequiredVariables(condition)
       availableVars := getAvailableVariables(bindings)
       
       if !allVariablesAvailable(requiredVars, availableVars) {
           // Évaluation partielle
           return evaluatePartial(bindings, availableVars)
       }
       
       // Évaluation complète
       result := evaluator.Evaluate(condition, timeout)
       return result
   }
   ```

#### 2.3 Mettre à Jour les Tests

**Tests unitaires** :

```go
// Trouver les tests existants
// rete/node_join_test.go

func TestEvaluateCondition(t *testing.T) {
    // ✅ Mettre à jour avec nouveaux paramètres
    result := evaluateCondition(bindings, DefaultEvaluationTimeout)
    
    // ✅ Ajouter tests pour nouveau comportement
    t.Run("partial_evaluation_missing_vars", func(t *testing.T) {
        bindings := map[string]*Fact{
            "u": userFact,
            "o": orderFact,
            // "p" manquant volontairement
        }
        
        // ✅ Extraction depuis réseau réel (pas de simulation)
        result := evaluateCondition(bindings, DefaultEvaluationTimeout)
        
        // Vérifier comportement partiel
        if !result {
            t.Error("Évaluation partielle devrait réussir")
        }
    })
}
```

**Tests RETE** :

```go
func TestRETEWithModifiedBehavior(t *testing.T) {
    network := buildNetwork()
    
    // Soumettre faits
    network.SubmitFact(fact1)
    network.SubmitFact(fact2)
    
    // ✅ OBLIGATOIRE : Extraction depuis réseau réel
    actualTokens := 0
    for _, terminal := range network.TerminalNodes {
        actualTokens += len(terminal.Memory.GetTokens())
    }
    
    // ❌ INTERDIT : expectedTokens := 5
    
    t.Logf("Tokens avec nouveau comportement : %d", actualTokens)
    
    // Vérifier que nouveau comportement fonctionne
    if actualTokens == 0 {
        t.Error("Aucun token créé avec nouveau comportement")
    }
}
```

#### 2.4 Mettre à Jour la Documentation

**GoDoc** :

```go
// EvaluateCondition évalue une condition avec gestion intelligente des variables manquantes.
//
// Si toutes les variables requises sont disponibles, effectue une évaluation complète.
// Sinon, effectue une évaluation partielle avec seulement les variables disponibles.
//
// Paramètres:
//   - bindings: Map des variables liées aux faits
//   - timeout: Délai maximum d'évaluation
//
// Retourne:
//   - true si la condition est satisfaite (partiellement ou complètement)
//   - false sinon
//
// Exemple:
//   result := evaluateCondition(bindings, 30*time.Second)
//
// Changement de comportement (v2.0):
//   Avant: Erreur si variable manquante
//   Après: Évaluation partielle intelligente
func EvaluateCondition(bindings map[string]*Fact, timeout time.Duration) bool {
    // ...
}
```

**README.md** :

```markdown
## Évaluation des Conditions (mise à jour v2.0)

Le moteur RETE évalue maintenant intelligemment les conditions même
quand certaines variables ne sont pas encore disponibles.

### Nouveau Comportement

- ✅ Évaluation partielle si variables manquantes
- ✅ Pas d'erreur "variable non liée"
- ✅ Propagation incrémentale optimisée

### Migration depuis v1.x

```go
// v1.x - Erreur si variable manquante
result := evaluateCondition(bindings)

// v2.0 - Avec timeout configurable
result := evaluateCondition(bindings, DefaultEvaluationTimeout)
```
```

**CHANGELOG.md** :

```markdown
## [2.0.0] - 2025-11-XX

### Changed
- **[BREAKING]** `evaluateCondition` prend maintenant un paramètre `timeout`
- Évaluation partielle intelligente quand variables manquantes
- Amélioration des performances pour jointures multi-variables

### Migration Guide
```go
// Avant
evaluateCondition(bindings)

// Après
evaluateCondition(bindings, DefaultEvaluationTimeout)
```

### Fixed
- Correction erreur "variable non liée" dans jointures 3-way
```

#### 2.5 Mettre à Jour les Prompts (Si Nécessaire)

**Identifier prompts impactés** :

```bash
# Chercher mentions de la fonctionnalité
grep -r "evaluateCondition\|évaluation.*condition" .github/prompts/
```

**Mettre à jour si trouvé** :

```markdown
# Dans .github/prompts/explain-code.md ou autre

## Évaluation des Conditions (MISE À JOUR v2.0)

La fonction `evaluateCondition` a été améliorée pour gérer
intelligemment les variables manquantes.

### Nouveau Comportement

Avant (v1.x) :
```go
// Erreur si variable 'p' manquante
result := evaluateCondition(bindings)
// ❌ Erreur: variable non liée: p
```

Après (v2.0) :
```go
// Évaluation partielle si variable 'p' manquante
result := evaluateCondition(bindings, timeout)
// ✅ Évalue seulement les conditions disponibles
```
```

### PHASE 3 : VALIDATION (Vérification)

#### 3.1 Tests Complets

**Checklist obligatoire** :

```bash
# 1. Formatage
go fmt ./...
goimports -w .

# 2. Analyse statique
go vet ./...
staticcheck ./...
golangci-lint run

# 3. Tests unitaires
go test ./...
go test -race ./...

# 4. Tests avec couverture
go test -cover ./... | tee coverage.txt
# Vérifier que couverture n'a pas baissé

# 5. Tests d'intégration
make test-integration

# 6. Runner universel RETE
make rete-unified
# Doit afficher : 58/58 ✅

# 7. Build
make build
make build-runners

# 8. Validation complète
make validate
```

**Tous doivent passer** ✅

#### 3.2 Tests de Régression

**Vérifier qu'on n'a rien cassé** :

```bash
# Tests de régression spécifiques
go test -run TestExisting ./...

# Benchmarks (vérifier pas de dégradation perf)
go test -bench=. ./... -benchmem

# Tests flaky
go test -count=10 ./...
```

#### 3.3 Validation Manuelle

**Tests manuels** :

1. **Cas d'usage principaux** :
   - Tester avec fichiers .constraint existants
   - Vérifier exemples du README
   - Tester cas limites

2. **Compatibilité** :
   - Anciens tests passent toujours
   - Migration path fonctionne
   - Deprecated warnings clairs

#### 3.4 Revue de Code

**Auto-revue** :

- [ ] **AUCUN hardcoding** introduit
- [ ] **Code générique** maintenu
- [ ] **Tests mis à jour** (extraction réseau réel pour RETE)
- [ ] **Documentation à jour**
- [ ] **Breaking changes documentés**
- [ ] **Migration path fourni** si nécessaire
- [ ] **go vet** sans erreur
- [ ] **golangci-lint** sans erreur
- [ ] **Tous les tests** passent
- [ ] **Aucune régression**

## Critères de Succès

### ✅ Modification Réussie

- [ ] Comportement modifié comme souhaité
- [ ] **AUCUN hardcoding** introduit
- [ ] **Code générique** préservé
- [ ] Tous les tests passent
- [ ] **Tests RETE avec extraction réseau réel**
- [ ] Aucune régression
- [ ] Documentation mise à jour
- [ ] Prompts mis à jour si nécessaire

### ✅ Qualité Maintenue

- [ ] go vet : 0 erreur
- [ ] golangci-lint : 0 erreur
- [ ] Couverture maintenue ou améliorée
- [ ] Performance maintenue ou améliorée
- [ ] Code review positive

### ✅ Compatibilité Gérée

- [ ] Breaking changes identifiés
- [ ] Migration path documenté
- [ ] Deprecated markers si nécessaire
- [ ] CHANGELOG mis à jour
- [ ] Version incrémentée correctement

## Format de Réponse

```
=== MODIFICATION DE FONCTIONNALITÉ ===

📋 IDENTIFICATION

Fonctionnalité : [Nom]
Localisation : [Module/Fichier/Fonction]
Comportement actuel : [Description]
Comportement souhaité : [Description]
Raison : [Justification]

📊 ANALYSE D'IMPACT

Fichiers impactés :
  • Code : X fichiers
  • Tests : X fichiers
  • Documentation : X fichiers
  • Prompts : X fichiers

Breaking changes : Oui/Non
Migration nécessaire : Oui/Non

🔧 MODIFICATIONS EFFECTUÉES

Code :
  ✅ Modifié evaluateCondition (rete/node_join.go)
  ✅ Ajouté extractRequiredVariables
  ✅ Ajouté evaluatePartial
  ⚠️ **VÉRIFIÉ** : Aucun hardcoding introduit
  ⚠️ **VÉRIFIÉ** : Code générique maintenu

Tests :
  ✅ Mis à jour TestEvaluateCondition
  ✅ Ajouté TestPartialEvaluation
  ✅ Corrigé TestIncrementalPropagation
  ⚠️ **VÉRIFIÉ** : Extraction réseau RETE réel

Documentation :
  ✅ Mis à jour GoDoc
  ✅ Mis à jour README.md
  ✅ Mis à jour CHANGELOG.md
  ✅ Ajouté guide de migration

Prompts :
  ✅ Mis à jour explain-code.md
  ✅ Mis à jour validate-network.md

✅ VALIDATION

Tests :
  ✅ go test ./... : PASS
  ✅ go test -race ./... : PASS
  ✅ make test-integration : PASS
  ✅ make rete-unified : 58/58 ✅

Qualité :
  ✅ go vet : 0 erreur
  ✅ golangci-lint : 0 erreur
  ✅ Couverture : 85% (maintenue)
  ✅ Benchmarks : pas de dégradation

Régression :
  ✅ Tests existants : PASS
  ✅ Exemples README : OK
  ✅ Cas limites : OK

📈 RÉSULTATS

Avant → Après :
  • Erreurs "variable non liée" : X → 0
  • Tokens terminaux créés : X → Y
  • Performance : Xms → Yms
  • Couverture : X% → Y%

🎯 VERDICT : MODIFICATION RÉUSSIE ✅

Breaking changes : [Liste si applicable]
Migration : [Lien vers guide]
```

## Exemple d'Utilisation

```
Je veux modifier le comportement de evaluateJoinConditions dans
rete/node_join.go pour gérer intelligemment les variables manquantes
au lieu de générer une erreur.

Comportement actuel : Erreur si variable non liée
Comportement souhaité : Évaluation partielle intelligente

Utilise le prompt "modify-behavior" pour :
1. Analyser l'impact
2. Modifier le code
3. Mettre à jour les tests
4. Mettre à jour la documentation
5. Valider complètement
```

## Checklist de Modification

### Avant de Commencer
- [ ] Fonctionnalité clairement identifiée
- [ ] Impact analysé (fichiers, tests, docs)
- [ ] Breaking changes identifiés
- [ ] Branche dédiée créée
- [ ] Tests passent avant modification

### Pendant la Modification
- [ ] Code modifié sans hardcoding
- [ ] Code générique maintenu
- [ ] Tests mis à jour (extraction RETE réel)
- [ ] Documentation mise à jour
- [ ] Prompts mis à jour si nécessaire
- [ ] Commits fréquents et clairs

### Après la Modification
- [ ] **Tous les tests passent** ✅
- [ ] **Aucun hardcoding** introduit ✅
- [ ] **Tests RETE avec extraction réseau réel** ✅
- [ ] go vet et golangci-lint sans erreur ✅
- [ ] Documentation complète et à jour ✅
- [ ] CHANGELOG mis à jour ✅
- [ ] Migration path documenté (si breaking) ✅
- [ ] Code review effectuée ✅

## Commandes Utiles

```bash
# Analyse d'impact
grep -r "FunctionName" --include="*.go" .
go list -f '{{.ImportPath}} {{.Imports}}' ./...

# Tests
make test
make test-coverage
make test-integration
make rete-unified
go test -race ./...
go test -count=10 ./...

# Qualité
go vet ./...
golangci-lint run
gocyclo -over 15 .

# Documentation
godoc -http=:6060
```

## Bonnes Pratiques

1. **Analyser avant d'agir** : Comprendre l'impact complet
2. **Modifier progressivement** : Petits changements testés
3. **Tester exhaustivement** : Tous les cas, y compris limites
4. **Documenter clairement** : Breaking changes, migration
5. **Respecter les règles** : Pas de hardcoding, tests RETE réels
6. **Valider complètement** : Tous les tests, toutes les validations

## Avertissements

⚠️ **ATTENTION** :
- Toujours analyser l'impact avant modification
- Ne jamais supposer qu'un changement est "petit"
- Toujours tester la compatibilité
- Documenter les breaking changes
- Préserver l'API publique si possible
- Fournir un chemin de migration clair

## Ressources

- [Semantic Versioning](https://semver.org/)
- [Go Module Version Compatibility](https://go.dev/blog/module-compatibility)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go 1 Compatibility Promise](https://go.dev/doc/go1compat)

---

**Version** : 1.0  
**Dernière mise à jour** : Novembre 2025  
**Type de changement** : Modification de comportement existant