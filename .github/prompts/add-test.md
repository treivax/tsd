# 🧪 Ajouter des Tests (Add Test)

## Contexte

Projet TSD (Type System with Dependencies) - Moteur de règles RETE avec système de contraintes en Go.

Tu veux ajouter des tests manquants pour augmenter la couverture, tester des cas limites, ou valider de nouvelles fonctionnalités.

## Objectif

Ajouter des tests de qualité en :
- Identifiant les gaps de couverture
- Écrivant des tests complets et déterministes
- Couvrant les cas nominaux, limites et d'erreur
- Respectant les standards de test du projet
- Validant avec extraction réseau RETE réel

## 📄 RÈGLES DE LICENCE ET COPYRIGHT - OBLIGATOIRE

### 🔒 Vérification de Compatibilité de Licence

**SI les tests nécessitent du code externe ou une nouvelle bibliothèque** :

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

**SI création de nouveaux fichiers de test** :

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

- ❌ **Ne JAMAIS copier du code de test** sans vérifier la licence
- ❌ **Ne JAMAIS utiliser de code GPL/AGPL** (incompatible avec MIT)
- ❌ **Ne JAMAIS omettre les en-têtes de copyright** dans les nouveaux fichiers
- ✅ **TOUJOURS écrire des tests originaux**

## ⚠️ RÈGLES STRICTES - TESTS

### 🚫 INTERDICTIONS ABSOLUES

1. **TESTS RETE** :
   - ❌ AUCUNE simulation de résultats
   - ❌ AUCUN mock du réseau RETE
   - ❌ AUCUN calcul manuel de tokens attendus
   - ❌ AUCUNE supposition sur les résultats
   - ✅ **TOUJOURS** extraire depuis le réseau RETE réel
   - ✅ **TOUJOURS** interroger les TerminalNodes
   - ✅ **TOUJOURS** inspecter les mémoires (Left/Right/Result)

2. **TESTS GOLANG** :
   - ❌ AUCUN hardcoding de valeurs de test
   - ❌ AUCUN test non-déterministe (flaky)
   - ❌ AUCUNE dépendance entre tests
   - ✅ Tests isolés et indépendants
   - ✅ Constantes nommées pour valeurs de test
   - ✅ Setup/teardown propre

3. **RACE DETECTOR - OBLIGATOIRE** :
   - 🏁 **TOUJOURS exécuter `go test -race`** après avoir écrit/modifié des tests
   - ❌ Ne JAMAIS valider des tests sans avoir exécuté `-race`
   - ⚠️ Les race conditions ne sont détectées QUE avec le flag `-race`
   - ✅ Ajouter `make test-race` à la validation
   - ✅ Fixer toute race condition détectée avant de considérer les tests terminés
   - 📖 **Pourquoi** : Race conditions causent bugs intermittents, corruption données, crashes production
   - ⏱️ **Note** : `-race` est ~10x plus lent mais OBLIGATOIRE

4. **QUALITÉ** :
   - ❌ Pas de tests qui passent toujours
   - ❌ Pas de tests qui testent rien
   - ✅ Assertions claires et explicites
   - ✅ Messages d'erreur descriptifs
   - ✅ Tests documentés

## Instructions

### PHASE 1 : ANALYSE (Identifier les Gaps)

#### 1.1 Analyser la Couverture Actuelle

**Générer rapport de couverture** :

```bash
# Couverture globale
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | tee coverage.txt

# Identifier fichiers peu couverts
go tool cover -func=coverage.out | grep -v "100.0%" | sort -k3 -n

# Visualisation HTML
go tool cover -html=coverage.out -o coverage.html
# Ouvrir coverage.html dans navigateur
```

**Analyser par package** :

```bash
# Couverture détaillée par package
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | grep -E "^[^[:space:]]" | column -t

# Résultat typique :
# rete/node_join.go:     evaluateCondition    75.0%
# rete/propagation.go:   propagateToChildren  50.0%
# constraint/parser.go:  parseExpression      90.0%
```

**Identifier gaps** :

```
Fichiers < 80% de couverture :
1. rete/node_join.go (75%)
   - Fonctions non testées : extractRequiredVariables, evaluatePartial
   - Cas limites manquants : variables toutes absentes, condition nil

2. rete/propagation.go (50%)
   - Fonctions non testées : propagateWithRetry, handleError
   - Cas d'erreur non testés : mémoire pleine, timeout

3. constraint/parser.go (90%)
   - Cas limites : expressions très imbriquées, caractères spéciaux
```

#### 1.2 Identifier Types de Tests Manquants

**Catégories de tests** :

1. **Tests unitaires** :
   - Fonctions individuelles
   - Comportement nominal
   - Cas limites
   - Gestion d'erreurs

2. **Tests d'intégration** :
   - Interaction entre composants
   - Flux complets
   - Scénarios réels

3. **Tests RETE** :
   - Construction réseau
   - Propagation de faits
   - Évaluation de conditions
   - Validation résultats

4. **Tests de cas limites** :
   - Valeurs nulles/vides
   - Valeurs extrêmes (min/max)
   - Données invalides
   - Concurrence

5. **Tests de régression** :
   - Bugs corrigés précédemment
   - Comportements critiques
   - Optimisations

#### 1.3 Prioriser les Tests à Ajouter

**Matrice de priorisation** :

```
Criticité vs Couverture :

HAUTE PRIORITÉ (Critique + Faible couverture) :
- evaluateCondition (75%, fonction critique)
- propagateToChildren (50%, cœur du moteur)
- parseExpression cas limites (90%, parseur sensible)

MOYENNE PRIORITÉ (Moyenne criticité + Faible couverture) :
- extractRequiredVariables (0%, utilitaire important)
- handleError (0%, gestion erreurs)

BASSE PRIORITÉ (Faible criticité ou Haute couverture) :
- Fonctions d'affichage/logging
- Fonctions déjà bien testées (>95%)
```

### PHASE 2 : ÉCRITURE DES TESTS

#### 2.1 Structure de Test Standard

**Template de test** :

```go
// rete/feature_test.go
package rete

import (
    "testing"
)

// Constantes de test (pas de hardcoding)
const (
    TestUserID    = "U1"
    TestOrderID   = "O1"
    TestTimeout   = 5 * time.Second
)

func TestFeatureName_NominalCase(t *testing.T) {
    t.Log("🧪 TEST CAS NOMINAL")
    t.Log("===================")
    
    // Arrange - Setup
    input := setupTestInput()
    expected := setupExpected()
    
    // Act - Exécution
    result, err := functionToTest(input)
    
    // Assert - Vérification
    if err != nil {
        t.Fatalf("❌ Erreur inattendue : %v", err)
    }
    
    if result != expected {
        t.Errorf("❌ Attendu %v, reçu %v", expected, result)
    }
    
    t.Log("✅ Test réussi")
}

func TestFeatureName_EdgeCases(t *testing.T) {
    tests := []struct {
        name    string
        input   InputType
        want    OutputType
        wantErr bool
    }{
        {
            name:    "valeur_nil",
            input:   nil,
            want:    defaultValue,
            wantErr: false,
        },
        {
            name:    "valeur_vide",
            input:   empty,
            want:    defaultValue,
            wantErr: false,
        },
        {
            name:    "valeur_invalide",
            input:   invalid,
            want:    nil,
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := functionToTest(tt.input)
            
            if (err != nil) != tt.wantErr {
                t.Errorf("wantErr=%v, got err=%v", tt.wantErr, err)
                return
            }
            
            if got != tt.want {
                t.Errorf("want=%v, got=%v", tt.want, got)
            }
        })
    }
}
```

#### 2.2 Tests RETE Spécifiques

**Template test RETE** :

```go
func TestRETEFeature_Propagation(t *testing.T) {
    t.Log("🎯 TEST PROPAGATION RETE")
    t.Log("========================")
    
    // Arrange - Construire réseau
    network := buildTestNetwork()
    
    // Préparer faits de test
    userFact := &Fact{
        ID:   TestUserID,
        Type: "User",
        Fields: map[string]interface{}{
            "id":  TestUserID,
            "age": 25,
        },
    }
    
    orderFact := &Fact{
        ID:   TestOrderID,
        Type: "Order",
        Fields: map[string]interface{}{
            "id":          TestOrderID,
            "customer_id": TestUserID,
            "amount":      100,
        },
    }
    
    // Act - Soumettre faits
    err := network.SubmitFact(userFact)
    if err != nil {
        t.Fatalf("❌ Erreur soumission User : %v", err)
    }
    
    err = network.SubmitFact(orderFact)
    if err != nil {
        t.Fatalf("❌ Erreur soumission Order : %v", err)
    }
    
    // Assert - ✅ OBLIGATOIRE : Extraction depuis réseau RETE réel
    actualTokens := 0
    for _, terminal := range network.TerminalNodes {
        actualTokens += len(terminal.Memory.GetTokens())
    }
    
    // ❌ INTERDIT : expectedTokens := 5 (hardcodé/simulé)
    
    t.Logf("📊 Tokens terminaux extraits : %d", actualTokens)
    
    // Vérifier qu'au moins un token a été créé
    if actualTokens == 0 {
        t.Error("❌ Aucun token terminal créé")
    }
    
    // ✅ Inspecter contenu des tokens réels
    for _, terminal := range network.TerminalNodes {
        tokens := terminal.Memory.GetTokens()
        t.Logf("TerminalNode %s : %d tokens", terminal.GetID(), len(tokens))
        
        for i, token := range tokens {
            t.Logf("  Token %d : %d faits", i, len(token.Facts))
            
            // Valider bindings
            if len(token.Bindings) == 0 {
                t.Error("❌ Token sans bindings")
            }
            
            // Valider variables
            for varName, fact := range token.Bindings {
                t.Logf("    %s -> %s (ID: %s)", varName, fact.Type, fact.ID)
                
                if fact == nil {
                    t.Errorf("❌ Binding %s est nil", varName)
                }
            }
        }
    }
    
    t.Log("✅ Test RETE réussi")
}

func TestRETEFeature_MultipleScenarios(t *testing.T) {
    scenarios := []struct {
        name     string
        facts    []*Fact
        validate func(*testing.T, *Network)
    }{
        {
            name: "scenario_simple",
            facts: []*Fact{userFact1, orderFact1},
            validate: func(t *testing.T, net *Network) {
                // ✅ Extraction réseau réel
                count := 0
                for _, term := range net.TerminalNodes {
                    count += len(term.Memory.GetTokens())
                }
                
                if count == 0 {
                    t.Error("❌ Aucun token dans scénario simple")
                }
            },
        },
        {
            name: "scenario_complexe",
            facts: []*Fact{userFact1, orderFact1, productFact1},
            validate: func(t *testing.T, net *Network) {
                // Validation spécifique au scénario
                for _, term := range net.TerminalNodes {
                    for _, token := range term.Memory.GetTokens() {
                        if len(token.Bindings) < 3 {
                            t.Error("❌ Token incomplet dans scénario complexe")
                        }
                    }
                }
            },
        },
    }
    
    for _, scenario := range scenarios {
        t.Run(scenario.name, func(t *testing.T) {
            network := buildTestNetwork()
            
            // Soumettre tous les faits
            for _, fact := range scenario.facts {
                if err := network.SubmitFact(fact); err != nil {
                    t.Fatalf("❌ Erreur soumission : %v", err)
                }
            }
            
            // Validation personnalisée
            scenario.validate(t, network)
        })
    }
}
```

#### 2.3 Tests de Cas Limites

**Cas limites à tester** :

```go
func TestFeature_EdgeCases(t *testing.T) {
    t.Log("🔍 TEST CAS LIMITES")
    t.Log("===================")
    
    t.Run("nil_input", func(t *testing.T) {
        result, err := function(nil)
        if err == nil {
            t.Error("❌ Devrait retourner erreur pour nil")
        }
    })
    
    t.Run("empty_input", func(t *testing.T) {
        result, err := function(emptyValue)
        // Vérifier comportement avec valeur vide
    })
    
    t.Run("max_value", func(t *testing.T) {
        result, err := function(math.MaxInt64)
        // Vérifier comportement avec valeur max
    })
    
    t.Run("negative_value", func(t *testing.T) {
        result, err := function(-1)
        // Vérifier comportement avec valeur négative
    })
    
    t.Run("special_characters", func(t *testing.T) {
        result, err := function("test\n\t\r\x00")
        // Vérifier caractères spéciaux
    })
    
    t.Run("concurrent_access", func(t *testing.T) {
        var wg sync.WaitGroup
        for i := 0; i < 100; i++ {
            wg.Add(1)
            go func() {
                defer wg.Done()
                function(testValue)
            }()
        }
        wg.Wait()
    })
}
```

#### 2.4 Tests d'Erreur

**Tests de gestion d'erreurs** :

```go
func TestFeature_ErrorHandling(t *testing.T) {
    t.Log("⚠️  TEST GESTION ERREURS")
    t.Log("========================")
    
    errorCases := []struct {
        name        string
        input       InputType
        expectedErr error
    }{
        {
            name:        "invalid_type",
            input:       invalidType,
            expectedErr: ErrInvalidType,
        },
        {
            name:        "out_of_range",
            input:       outOfRange,
            expectedErr: ErrOutOfRange,
        },
        {
            name:        "not_found",
            input:       notFound,
            expectedErr: ErrNotFound,
        },
    }
    
    for _, tc := range errorCases {
        t.Run(tc.name, func(t *testing.T) {
            _, err := function(tc.input)
            
            if err == nil {
                t.Errorf("❌ Attendait erreur %v, reçu nil", tc.expectedErr)
                return
            }
            
            if !errors.Is(err, tc.expectedErr) {
                t.Errorf("❌ Attendait %v, reçu %v", tc.expectedErr, err)
            }
        })
    }
}
```

### PHASE 3 : VALIDATION

#### 3.1 Exécuter et Valider Tests

**Validation complète** :

```bash
# Tests nouveaux
go test -v -run TestNewTests ./...

# Tests complets
go test ./...

# 🏁 OBLIGATOIRE : Avec race detector (détecte race conditions)
go test -race ./...
# ⚠️ CRITICAL: Ce test est OBLIGATOIRE pour détecter les race conditions
# Les race conditions ne sont détectées QUE par le flag -race
# TOUJOURS exécuter ce test, même si plus lent (~10x)

# Avec couverture
go test -cover ./...

# Tests d'intégration
make test-integration

# Runner universel RETE
make rete-unified  # Doit afficher 58/58 ✅

# Validation complète
make validate
```

#### 3.2 Vérifier Couverture Améliorée

**Mesure de l'amélioration** :

```bash
# Couverture avant
go test -coverprofile=before.out ./...
go tool cover -func=before.out | tail -1

# Couverture après
go test -coverprofile=after.out ./...
go tool cover -func=after.out | tail -1

# Comparaison
echo "Avant :"
go tool cover -func=before.out | grep "total:"
echo "Après :"
go tool cover -func=after.out | grep "total:"
```

#### 3.3 Tests de Qualité

**Validation qualité tests** :

```bash
# Tests déterministes (10 runs)
go test -count=10 ./...

# Tests ne sont pas flaky
for i in {1..20}; do go test ./... || break; done

# Tests isolés (ordre aléatoire)
go test -shuffle=on ./...

# Pas de dépendances entre tests
go test -parallel=8 ./...
```

## Critères de Succès

### ✅ Tests Ajoutés

- [ ] Gaps de couverture identifiés
- [ ] Tests écrits pour cas nominaux
- [ ] Tests écrits pour cas limites
- [ ] Tests écrits pour gestion d'erreurs
- [ ] **Tests RETE avec extraction réseau réel**
- [ ] **AUCUN hardcoding** dans les tests
- [ ] 🏁 **`go test -race` exécuté et passé (OBLIGATOIRE)**
- [ ] **Aucune race condition détectée**
- [ ] Tests déterministes (pas flaky)
- [ ] Tests isolés et indépendants

### ✅ Couverture Améliorée

- [ ] Couverture globale augmentée
- [ ] Fichiers critiques > 80%
- [ ] Fonctions importantes testées
- [ ] Cas limites couverts
- [ ] Gestion d'erreurs testée

### ✅ Qualité

- [ ] Tous les tests passent
- [ ] 🏁 **`go test -race ./...` passe sans erreur (OBLIGATOIRE)**
- [ ] Aucun test flaky (10 runs)
- [ ] go vet sans erreur
- [ ] Tests documentés
- [ ] Messages d'assertion clairs

## Format de Réponse

```
=== AJOUT DE TESTS ===

📊 ANALYSE COUVERTURE INITIALE

Couverture globale : 72%

Fichiers < 80% :
  • rete/node_join.go : 75%
  • rete/propagation.go : 50%
  • constraint/parser.go : 90%

Fonctions non testées :
  • extractRequiredVariables (0%)
  • evaluatePartial (0%)
  • handleError (0%)

🎯 TESTS AJOUTÉS

Tests unitaires :
  ✅ TestExtractRequiredVariables_NominalCase
  ✅ TestExtractRequiredVariables_EdgeCases
  ✅ TestEvaluatePartial_WithMissingVars
  ✅ TestHandleError_AllErrorTypes

Tests RETE :
  ✅ TestPropagation_MultipleVariables
  ✅ TestPropagation_IncrementalSubmission
  ⚠️ **VÉRIFIÉ** : Extraction réseau RETE réel
  ⚠️ **VÉRIFIÉ** : Aucune simulation

Tests cas limites :
  ✅ TestNilValues
  ✅ TestEmptyInputs
  ✅ TestMaxValues
  ✅ TestConcurrentAccess

Tests erreurs :
  ✅ TestErrorHandling_InvalidInput
  ✅ TestErrorHandling_OutOfRange
  ✅ TestErrorHandling_NotFound

Total : 15 nouveaux tests ajoutés

✅ VALIDATION

Tests :
  ✅ go test ./... : PASS (tous les tests)
  ✅ go test -race ./... : PASS
  ✅ go test -count=10 ./... : PASS (déterministes)
  ✅ make test-integration : PASS
  ✅ make rete-unified : 58/58 ✅

Couverture :
  Avant : 72%
  Après : 87% (+15%)
  
  Fichiers améliorés :
  • rete/node_join.go : 75% → 92% (+17%)
  • rete/propagation.go : 50% → 85% (+35%)
  • constraint/parser.go : 90% → 95% (+5%)

Qualité :
  ✅ go vet : 0 erreur
  ✅ Tests isolés : OK
  ✅ Pas de flaky tests : OK
  ✅ Messages clairs : OK

📈 RÉSULTATS

Couverture globale : 72% → 87% (+15%)
Tests ajoutés : 15
Lignes testées : +450
Branches testées : +120

🎯 VERDICT : TESTS AJOUTÉS AVEC SUCCÈS ✅
```

## Exemple d'Utilisation

```
La fonction evaluateCondition a seulement 75% de couverture.
Les cas avec variables manquantes ne sont pas testés.

Utilise le prompt "add-test" pour :
1. Analyser les gaps de couverture
2. Identifier les cas non testés
3. Ajouter tests avec extraction RETE réelle
4. Valider couverture améliorée
```

## Checklist

### Avant d'Écrire
- [ ] Couverture actuelle analysée
- [ ] Gaps identifiés
- [ ] Types de tests déterminés
- [ ] Priorités définies

### Pendant l'Écriture
- [ ] Tests isolés et indépendants
- [ ] **AUCUN hardcoding** valeurs test
- [ ] **Tests RETE extraction réseau réel**
- [ ] Cas nominaux testés
- [ ] Cas limites testés
- [ ] Gestion erreurs testée
- [ ] Messages assertion clairs

### Après l'Écriture
- [ ] **Tous les tests passent** ✅
- [ ] **Tests déterministes** (10 runs) ✅
- [ ] **Tests RETE extraction réseau réel** ✅
- [ ] Couverture améliorée ✅
- [ ] go vet sans erreur ✅
- [ ] Tests documentés ✅

## Commandes Utiles

```bash
# Couverture
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go tool cover -html=coverage.out

# Tests spécifiques
go test -v -run TestName ./...

# Tests déterministes
go test -count=10 ./...

# Race conditions
go test -race ./...

# Shuffle
go test -shuffle=on ./...

# Parallel
go test -parallel=8 ./...
```

## Bonnes Pratiques

1. **Tester comportement, pas implémentation**
2. **Tests isolés** : Aucune dépendance entre tests
3. **Tests déterministes** : Mêmes entrées → mêmes sorties
4. **Messages clairs** : Assertions explicites
5. **Table-driven** : Pour tests similaires
6. **Setup/teardown** : Propre et minimal
7. **Extraction RETE réelle** : TOUJOURS pour tests RETE

## Anti-Patterns à Éviter

❌ **Tests qui testent rien** :
```go
func TestSomething(t *testing.T) {
    result := function()
    // Aucune assertion !
}
```

❌ **Hardcoding résultats** :
```go
func TestRETETokens(t *testing.T) {
    expectedTokens := 5  // Hardcodé !
}
```

❌ **Tests dépendants** :
```go
func TestA(t *testing.T) {
    globalVar = "value"  // State partagé
}
func TestB(t *testing.T) {
    // Dépend de TestA
}
```

✅ **Bons tests** :
```go
func TestFeature_Isolated(t *testing.T) {
    // Setup propre
    input := createTestInput()
    
    // Exécution
    result := function(input)
    
    // Assertions claires
    if result == nil {
        t.Error("Result should not be nil")
    }
}
```

## Ressources

- [Go Testing](https://go.dev/doc/tutorial/add-a-test)
- [Table Driven Tests](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)
- [Testing Best Practices](https://go.dev/blog/subtests)

---

**Version** : 1.0  
**Dernière mise à jour** : Novembre 2025  
**Type** : Ajout de tests avec extraction RETE réelle obligatoire