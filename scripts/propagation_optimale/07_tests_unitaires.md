# 🧪 Prompt 07 - Tests Unitaires Complets

> **📋 Standards** : Ce prompt respecte les règles de `.github/prompts/common.md` et `.github/prompts/develop.md`

## 🎯 Objectif

Développer une suite de tests unitaires exhaustive pour tous les composants du système de propagation delta. Garantir une couverture > 90% et valider tous les cas limites.

**⚠️ IMPORTANT** : Ce prompt génère du code de tests. Respecter strictement les standards de `common.md`.

---

## 📋 Prérequis

Avant de commencer ce prompt :

- [x] **Prompts 01-06 validés** : Tout le système delta implémenté et intégré
- [x] **Code fonctionnel** : `go build ./...` réussit
- [x] **Documents de référence** :
  - `REPORTS/conception_delta_architecture.md`
  - `.github/prompts/common.md` - Standards de tests
  - Tous les fichiers sources du package `rete/delta`

---

## 📂 Fichiers de Tests à Compléter/Créer

```
rete/delta/
├── field_delta_test.go              # Compléter tests edge cases
├── dependency_index_test.go         # Compléter tests concurrence
├── delta_detector_test.go           # Compléter tests complexes
├── delta_propagator_test.go         # Compléter tests propagation
├── integration_test.go              # Tests end-to-end
└── test_helpers.go                  # Utilitaires de test (nouveau)
```

---

## 🔧 Tâche 1 : Utilitaires de Test

### Fichier : `rete/delta/test_helpers.go`

**Contenu** :

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
    "testing"
)

// TestFixtures contient des données de test réutilisables.
type TestFixtures struct {
    SimpleFact      map[string]interface{}
    ComplexFact     map[string]interface{}
    LargeFact       map[string]interface{}
    NestedFact      map[string]interface{}
}

// NewTestFixtures crée un ensemble de fixtures de test.
func NewTestFixtures() *TestFixtures {
    return &TestFixtures{
        SimpleFact: map[string]interface{}{
            "id":     "123",
            "name":   "Test Product",
            "price":  100.0,
            "status": "active",
        },
        ComplexFact: map[string]interface{}{
            "id":       "456",
            "name":     "Complex Product",
            "price":    200.0,
            "quantity": 10,
            "category": "Electronics",
            "tags":     []interface{}{"new", "featured"},
            "metadata": map[string]interface{}{
                "brand": "TestBrand",
                "model": "XYZ",
            },
        },
        LargeFact: generateLargeFact(100),
        NestedFact: map[string]interface{}{
            "id": "789",
            "address": map[string]interface{}{
                "street": "123 Main St",
                "city":   "Paris",
                "country": map[string]interface{}{
                    "name": "France",
                    "code": "FR",
                },
            },
        },
    }
}

// generateLargeFact génère un fait avec N champs.
func generateLargeFact(fieldCount int) map[string]interface{} {
    fact := make(map[string]interface{})
    for i := 0; i < fieldCount; i++ {
        fieldName := "field_" + string(rune('0'+i%10)) + "_" + string(rune('0'+i/10))
        fact[fieldName] = i
    }
    return fact
}

// AssertDeltaEquals vérifie qu'un delta correspond aux attentes.
func AssertDeltaEquals(t *testing.T, delta *FactDelta, expectedFields map[string]struct{oldVal, newVal interface{}}) {
    t.Helper()
    
    if len(delta.Fields) != len(expectedFields) {
        t.Errorf("Expected %d changed fields, got %d", len(expectedFields), len(delta.Fields))
    }
    
    for fieldName, expected := range expectedFields {
        fieldDelta, exists := delta.Fields[fieldName]
        if !exists {
            t.Errorf("Expected field '%s' in delta, not found", fieldName)
            continue
        }
        
        if fieldDelta.OldValue != expected.oldVal {
            t.Errorf("Field '%s': expected old value %v, got %v", fieldName, expected.oldVal, fieldDelta.OldValue)
        }
        
        if fieldDelta.NewValue != expected.newVal {
            t.Errorf("Field '%s': expected new value %v, got %v", fieldName, expected.newVal, fieldDelta.NewValue)
        }
    }
}

// AssertNodesContain vérifie qu'une liste de nœuds contient les IDs attendus.
func AssertNodesContain(t *testing.T, nodes []NodeReference, expectedIDs []string) {
    t.Helper()
    
    actualIDs := make(map[string]bool)
    for _, node := range nodes {
        actualIDs[node.NodeID] = true
    }
    
    for _, expectedID := range expectedIDs {
        if !actualIDs[expectedID] {
            t.Errorf("Expected node '%s' in result, not found", expectedID)
        }
    }
}

// BenchmarkHelper facilite l'écriture de benchmarks.
type BenchmarkHelper struct {
    detector   *DeltaDetector
    index      *DependencyIndex
    propagator *DeltaPropagator
}

// NewBenchmarkHelper crée un helper pour benchmarks.
func NewBenchmarkHelper() *BenchmarkHelper {
    index := NewDependencyIndex()
    detector := NewDeltaDetector()
    
    propagator, _ := NewDeltaPropagatorBuilder().
        WithIndex(index).
        WithDetector(detector).
        Build()
    
    return &BenchmarkHelper{
        detector:   detector,
        index:      index,
        propagator: propagator,
    }
}
```

---

## 🔧 Tâche 2 : Tests Edge Cases - FieldDelta

### Ajouts à `rete/delta/field_delta_test.go`

**Nouveaux tests** :

```go
// Test changement de type avec tracking désactivé
func TestFieldDelta_TypeChangeNoTracking(t *testing.T) {
    config := DefaultDetectorConfig()
    config.TrackTypeChanges = false
    detector := NewDeltaDetectorWithConfig(config)
    
    oldFact := map[string]interface{}{"value": 42}
    newFact := map[string]interface{}{"value": "42"}
    
    delta, _ := detector.DetectDelta(oldFact, newFact, "Test~1", "Test")
    
    // Sans tracking, pourrait être considéré comme no-change
    // Vérifier comportement selon implémentation
}

// Test avec valeurs nil
func TestFieldDelta_NilValues(t *testing.T) {
    tests := []struct {
        name     string
        oldValue interface{}
        newValue interface{}
        wantType ChangeType
    }{
        {"nil → value", nil, "test", ChangeTypeAdded},
        {"value → nil", "test", nil, ChangeTypeRemoved},
        {"nil → nil", nil, nil, ChangeTypeModified},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            delta := NewFieldDelta("field", tt.oldValue, tt.newValue)
            if delta.ChangeType != tt.wantType {
                t.Errorf("Expected %v, got %v", tt.wantType, delta.ChangeType)
            }
        })
    }
}

// Test avec unicode et caractères spéciaux
func TestFieldDelta_UnicodeFields(t *testing.T) {
    detector := NewDeltaDetector()
    
    oldFact := map[string]interface{}{
        "名前":    "古い値",
        "emoji": "😀",
    }
    
    newFact := map[string]interface{}{
        "名前":    "新しい値",
        "emoji": "😎",
    }
    
    delta, err := detector.DetectDelta(oldFact, newFact, "Test~1", "Test")
    
    if err != nil {
        t.Fatalf("Unexpected error: %v", err)
    }
    
    if len(delta.Fields) != 2 {
        t.Errorf("Expected 2 changed fields, got %d", len(delta.Fields))
    }
}

// Test avec très grandes chaînes
func TestFieldDelta_LargeStrings(t *testing.T) {
    detector := NewDeltaDetector()
    
    largeString := string(make([]byte, 1024*1024)) // 1MB
    
    oldFact := map[string]interface{}{"data": largeString}
    newFact := map[string]interface{}{"data": largeString + "x"}
    
    delta, err := detector.DetectDelta(oldFact, newFact, "Test~1", "Test")
    
    if err != nil {
        t.Fatalf("Unexpected error: %v", err)
    }
    
    if delta.IsEmpty() {
        t.Error("Expected change detected")
    }
}
```

---

## 🔧 Tâche 3 : Tests Concurrence - DependencyIndex

### Ajouts à `rete/delta/dependency_index_test.go`

**Nouveaux tests** :

```go
// Test lecture/écriture concurrente intensive
func TestDependencyIndex_ConcurrentReadWrite(t *testing.T) {
    idx := NewDependencyIndex()
    
    // Setup initial
    for i := 0; i < 10; i++ {
        nodeID := fmt.Sprintf("node%d", i)
        idx.AddAlphaNode(nodeID, "Product", []string{"field1", "field2"})
    }
    
    done := make(chan bool, 20)
    
    // 10 writers
    for i := 0; i < 10; i++ {
        go func(id int) {
            for j := 0; j < 100; j++ {
                nodeID := fmt.Sprintf("concurrent_node_%d_%d", id, j)
                idx.AddAlphaNode(nodeID, "Product", []string{"price", "status"})
            }
            done <- true
        }(i)
    }
    
    // 10 readers
    for i := 0; i < 10; i++ {
        go func() {
            for j := 0; j < 100; j++ {
                _ = idx.GetAffectedNodes("Product", "price")
            }
            done <- true
        }()
    }
    
    // Attendre fin
    for i := 0; i < 20; i++ {
        <-done
    }
    
    // Vérifier cohérence finale
    stats := idx.GetStats()
    if stats.NodeCount == 0 {
        t.Error("Expected nodes in index after concurrent operations")
    }
}

// Test clear pendant lecture
func TestDependencyIndex_ClearDuringRead(t *testing.T) {
    idx := NewDependencyIndex()
    
    for i := 0; i < 100; i++ {
        idx.AddAlphaNode(fmt.Sprintf("node%d", i), "Product", []string{"price"})
    }
    
    done := make(chan bool, 2)
    
    // Reader continu
    go func() {
        for i := 0; i < 1000; i++ {
            _ = idx.GetAffectedNodes("Product", "price")
        }
        done <- true
    }()
    
    // Clear au milieu
    go func() {
        idx.Clear()
        done <- true
    }()
    
    <-done
    <-done
    
    // Ne devrait pas paniquer
}
```

---

## 🔧 Tâche 4 : Tests Complexes - DeltaDetector

### Ajouts à `rete/delta/delta_detector_test.go`

**Nouveaux tests** :

```go
// Test avec structure profondément imbriquée
func TestDeltaDetector_DeepNesting(t *testing.T) {
    config := DefaultDetectorConfig()
    config.EnableDeepComparison = true
    config.MaxNestingLevel = 10
    detector := NewDeltaDetectorWithConfig(config)
    
    // Créer structure imbriquée sur 10 niveaux
    oldFact := createDeepNestedFact(10, "old_value")
    newFact := createDeepNestedFact(10, "new_value")
    
    delta, err := detector.DetectDelta(oldFact, newFact, "Test~1", "Test")
    
    if err != nil {
        t.Fatalf("Unexpected error: %v", err)
    }
    
    if delta.IsEmpty() {
        t.Error("Expected change in deeply nested structure")
    }
}

func createDeepNestedFact(depth int, leafValue string) map[string]interface{} {
    if depth == 0 {
        return map[string]interface{}{"value": leafValue}
    }
    return map[string]interface{}{
        "nested": createDeepNestedFact(depth-1, leafValue),
    }
}

// Test protection contre stack overflow
func TestDeltaDetector_StackOverflowProtection(t *testing.T) {
    config := DefaultDetectorConfig()
    config.MaxNestingLevel = 5
    detector := NewDeltaDetectorWithConfig(config)
    
    // Structure circulaire simulée
    oldFact := map[string]interface{}{"level": createDeepNestedFact(20, "old")}
    newFact := map[string]interface{}{"level": createDeepNestedFact(20, "new")}
    
    // Ne devrait pas paniquer
    _, err := detector.DetectDelta(oldFact, newFact, "Test~1", "Test")
    
    if err != nil {
        t.Fatalf("Unexpected error: %v", err)
    }
}

// Test avec slices de tailles différentes
func TestDeltaDetector_DifferentSliceSizes(t *testing.T) {
    detector := NewDeltaDetector()
    
    oldFact := map[string]interface{}{
        "items": []interface{}{1, 2, 3},
    }
    
    newFact := map[string]interface{}{
        "items": []interface{}{1, 2, 3, 4, 5},
    }
    
    delta, err := detector.DetectDelta(oldFact, newFact, "Test~1", "Test")
    
    if err != nil {
        t.Fatalf("Unexpected error: %v", err)
    }
    
    if delta.IsEmpty() {
        t.Error("Expected change detected for different slice sizes")
    }
}

// Test cache expiration sous charge
func TestDeltaDetector_CacheExpirationUnderLoad(t *testing.T) {
    config := DefaultDetectorConfig()
    config.CacheComparisons = true
    config.CacheTTL = 50 * time.Millisecond
    detector := NewDeltaDetectorWithConfig(config)
    
    oldFact := map[string]interface{}{"price": 100.0}
    newFact := map[string]interface{}{"price": 150.0}
    
    // Remplir cache
    for i := 0; i < 100; i++ {
        _, _ = detector.DetectDelta(oldFact, newFact, fmt.Sprintf("Product~%d", i), "Product")
    }
    
    initialMetrics := detector.GetMetrics()
    
    // Attendre expiration
    time.Sleep(100 * time.Millisecond)
    
    // Nouvelles détections
    detector.ResetMetrics()
    for i := 0; i < 100; i++ {
        _, _ = detector.DetectDelta(oldFact, newFact, fmt.Sprintf("Product~%d", i), "Product")
    }
    
    newMetrics := detector.GetMetrics()
    
    // Cache devrait être expiré, donc plus de misses que de hits
    if newMetrics.CacheHits > newMetrics.CacheMisses {
        t.Error("Expected more cache misses after expiration")
    }
}
```

---

## 🔧 Tâche 5 : Tests de Performance

### Fichier : `rete/delta/performance_test.go` (nouveau)

**Contenu** :

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

import (
    "testing"
)

// BenchmarkDeltaDetector_Various mesure la performance de détection selon le contexte
func BenchmarkDeltaDetector_Various(b *testing.B) {
    fixtures := NewTestFixtures()
    
    benchmarks := []struct {
        name     string
        oldFact  map[string]interface{}
        newFact  map[string]interface{}
    }{
        {"NoChange_Simple", fixtures.SimpleFact, fixtures.SimpleFact},
        {"SingleChange_Simple", fixtures.SimpleFact, modifyField(fixtures.SimpleFact, "price", 200.0)},
        {"MultiChange_Complex", fixtures.ComplexFact, modifyFields(fixtures.ComplexFact, map[string]interface{}{"price": 300.0, "quantity": 20})},
        {"AllChange_Large", fixtures.LargeFact, modifyAllFields(fixtures.LargeFact)},
    }
    
    for _, bm := range benchmarks {
        b.Run(bm.name, func(b *testing.B) {
            detector := NewDeltaDetector()
            b.ResetTimer()
            
            for i := 0; i < b.N; i++ {
                _, _ = detector.DetectDelta(bm.oldFact, bm.newFact, "Test~1", "Test")
            }
        })
    }
}

// BenchmarkDependencyIndex_Scaling mesure la scalabilité de l'index
func BenchmarkDependencyIndex_Scaling(b *testing.B) {
    sizes := []int{10, 100, 1000, 10000}
    
    for _, size := range sizes {
        b.Run(fmt.Sprintf("Size_%d", size), func(b *testing.B) {
            idx := NewDependencyIndex()
            
            // Setup index
            for i := 0; i < size; i++ {
                idx.AddAlphaNode(fmt.Sprintf("node%d", i), "Product", []string{"price", "status"})
            }
            
            b.ResetTimer()
            
            for i := 0; i < b.N; i++ {
                _ = idx.GetAffectedNodes("Product", "price")
            }
        })
    }
}

// BenchmarkPropagation_EndToEnd mesure la performance end-to-end
func BenchmarkPropagation_EndToEnd(b *testing.B) {
    helper := NewBenchmarkHelper()
    fixtures := NewTestFixtures()
    
    // Setup index avec nœuds
    for i := 0; i < 50; i++ {
        helper.index.AddAlphaNode(fmt.Sprintf("alpha%d", i), "Product", []string{"price"})
    }
    
    oldFact := fixtures.SimpleFact
    newFact := modifyField(fixtures.SimpleFact, "price", 200.0)
    
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        _ = helper.propagator.PropagateUpdate(oldFact, newFact, "Product~123", "Product")
    }
}

// Helpers pour benchmarks
func modifyField(fact map[string]interface{}, field string, value interface{}) map[string]interface{} {
    modified := make(map[string]interface{})
    for k, v := range fact {
        modified[k] = v
    }
    modified[field] = value
    return modified
}

func modifyFields(fact map[string]interface{}, changes map[string]interface{}) map[string]interface{} {
    modified := make(map[string]interface{})
    for k, v := range fact {
        modified[k] = v
    }
    for k, v := range changes {
        modified[k] = v
    }
    return modified
}

func modifyAllFields(fact map[string]interface{}) map[string]interface{} {
    modified := make(map[string]interface{})
    for k := range fact {
        modified[k] = "modified"
    }
    return modified
}
```

---

## ✅ Validation

Après implémentation, exécuter :

```bash
# 1. Tests unitaires avec verbose
go test ./rete/delta/... -v

# 2. Couverture de code
go test ./rete/delta/... -cover -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html

# Vérifier : couverture > 90%

# 3. Benchmarks
go test ./rete/delta/... -bench=. -benchmem -benchtime=5s

# 4. Race detector
go test ./rete/delta/... -race -count=10

# 5. Tests de stress
go test ./rete/delta/... -count=100 -failfast

# 6. Validation complète
make test
```

**Critères de succès** :
- [ ] Tous les tests passent (100%)
- [ ] Couverture > 90% sur tous les fichiers
- [ ] Aucune race condition détectée
- [ ] Benchmarks stables (variance < 10%)
- [ ] Tests de stress réussis (100 runs)

---

## 📊 Rapport de Couverture

Créer le rapport : `REPORTS/delta_test_coverage.md`

```markdown
# Rapport de Couverture Tests - Propagation Delta

## Résumé Général

- **Couverture globale** : XX.X%
- **Fichiers testés** : N/N
- **Assertions totales** : XXXX
- **Benchmarks** : XX

## Détail par Fichier

| Fichier | Couverture | Lignes | Branches | Commentaire |
|---------|------------|--------|----------|-------------|
| field_delta.go | XX% | XXX/XXX | XX/XX | Complet |
| dependency_index.go | XX% | XXX/XXX | XX/XX | Complet |
| delta_detector.go | XX% | XXX/XXX | XX/XX | Complet |
| delta_propagator.go | XX% | XXX/XXX | XX/XX | Complet |

## Cas Limites Testés

- [x] Valeurs nil
- [x] Types différents
- [x] Unicode et caractères spéciaux
- [x] Structures imbriquées
- [x] Stack overflow protection
- [x] Concurrence intensive
- [x] Cache expiration
- [x] Grandes données (1MB+)

## Performance

- Détection no-op : < 100ns
- Détection 1 champ : < 500ns
- Index lookup : < 200ns
- Propagation complète : < 10µs
```

---

## 📊 Livrables

À la fin de ce prompt :

1. **Tests complets** :
   - ✅ `test_helpers.go` - Utilitaires de test
   - ✅ Tests edge cases complets
   - ✅ Tests concurrence exhaustifs
   - ✅ `performance_test.go` - Benchmarks

2. **Validation** :
   - ✅ Rapport de couverture > 90%
   - ✅ Résultats benchmarks
   - ✅ Rapport race detector

3. **Documentation** :
   - ✅ `REPORTS/delta_test_coverage.md`

---

## 🚀 Commit

Une fois validé :

```bash
git add rete/delta/ REPORTS/
git commit -m "test(rete): [Prompt 07] Suite de tests unitaires complète pour propagation delta

- Utilitaires de test réutilisables
- Tests edge cases exhaustifs (nil, unicode, grandes données)
- Tests concurrence intensive
- Tests stack overflow protection
- Benchmarks performance
- Couverture > 90% sur tous les fichiers
- Aucune race condition détectée
- Rapport de couverture détaillé"
```

---

## 🚦 Prochaine Étape

Passer au **Prompt 08 - Tests d'Intégration**

---

**Durée estimée** : 3-4 heures  
**Difficulté** : Moyenne  
**Prérequis** : Prompts 01-06 validés  
**Couverture cible** : > 90%