# Prompt 11 : Performance et Optimisation

**Session** : 11/12  
**Durée estimée** : 2-3 heures  
**Pré-requis** : Prompt 10 complété, tous les tests E2E passent

---

## 🎯 Objectif de cette Session

Valider que le refactoring n'introduit pas de régression de performance :
1. Créer des benchmarks pour mesurer les performances
2. Comparer avec les performances théoriques attendues
3. Optimiser si nécessaire (overhead < 10%)
4. Documenter les résultats

**Livrable** : `tsd/rete/node_join_benchmark_test.go` (nouveau, ~300-400 lignes)

---

## 📋 Tâches à Réaliser

### Tâche 1 : Créer le Fichier de Benchmarks (20 min)

**Fichier** : `tsd/rete/node_join_benchmark_test.go`

**En-tête** :
```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"testing"
)
```

---

### Tâche 2 : Benchmarks de BindingChain (40 min)

#### 2.1 Benchmark Add

```go
func BenchmarkBindingChain_Add(b *testing.B) {
	fact := &Fact{ID: "f1", Type: "User"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		chain := NewBindingChain()
		chain = chain.Add("var1", fact)
	}
}

func BenchmarkBindingChain_Add_10Variables(b *testing.B) {
	facts := make([]*Fact, 10)
	for i := 0; i < 10; i++ {
		facts[i] = &Fact{ID: fmt.Sprintf("f%d", i), Type: "Type"}
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		chain := NewBindingChain()
		for j := 0; j < 10; j++ {
			chain = chain.Add(fmt.Sprintf("var%d", j), facts[j])
		}
	}
}
```

#### 2.2 Benchmark Get

```go
func BenchmarkBindingChain_Get_SmallChain(b *testing.B) {
	// Chaîne avec 3 bindings
	chain := NewBindingChain()
	for i := 0; i < 3; i++ {
		chain = chain.Add(fmt.Sprintf("var%d", i), &Fact{ID: fmt.Sprintf("f%d", i)})
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = chain.Get("var1")
	}
}

func BenchmarkBindingChain_Get_LargeChain(b *testing.B) {
	// Chaîne avec 100 bindings
	chain := NewBindingChain()
	for i := 0; i < 100; i++ {
		chain = chain.Add(fmt.Sprintf("var%d", i), &Fact{ID: fmt.Sprintf("f%d", i)})
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = chain.Get("var0") // Pire cas : chercher le plus ancien
	}
}
```

#### 2.3 Benchmark Merge

```go
func BenchmarkBindingChain_Merge(b *testing.B) {
	chain1 := NewBindingChain()
	chain2 := NewBindingChain()
	
	for i := 0; i < 5; i++ {
		chain1 = chain1.Add(fmt.Sprintf("v1_%d", i), &Fact{ID: fmt.Sprintf("f1_%d", i)})
		chain2 = chain2.Add(fmt.Sprintf("v2_%d", i), &Fact{ID: fmt.Sprintf("f2_%d", i)})
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = chain1.Merge(chain2)
	}
}
```

---

### Tâche 3 : Benchmarks de JoinNode (50 min)

#### 3.1 Benchmark jointure 2 variables (baseline)

```go
func BenchmarkJoinNode_2Variables(b *testing.B) {
	// Setup
	userFact := &Fact{
		ID:   "u1",
		Type: "User",
		Attributes: map[string]interface{}{"id": 1},
	}
	orderFact := &Fact{
		ID:   "o1",
		Type: "Order",
		Attributes: map[string]interface{}{"user_id": 1},
	}
	
	joinNode := &JoinNode{
		BaseNode: BaseNode{
			ID:       "join_bench",
			Children: []Node{},
		},
		LeftVariables:  []string{"user"},
		RightVariables: []string{"order"},
		AllVariables:   []string{"user", "order"},
		VariableTypes: map[string]string{
			"user":  "User",
			"order": "Order",
		},
		LeftMemory:     []*Token{},
		RightMemory:    []*Fact{},
		JoinConditions: nil,
	}
	
	leftToken := NewTokenWithFact(userFact, "user", "test")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		joinNode.LeftMemory = []*Token{}
		joinNode.RightMemory = []*Fact{}
		
		_ = joinNode.ActivateLeft(leftToken)
		_ = joinNode.ActivateRight(orderFact)
	}
}
```

#### 3.2 Benchmark jointure 3 variables

```go
func BenchmarkJoinNode_3Variables(b *testing.B) {
	// Setup cascade complète
	userFact := &Fact{ID: "u1", Type: "User", Attributes: map[string]interface{}{"id": 1}}
	orderFact := &Fact{ID: "o1", Type: "Order", Attributes: map[string]interface{}{"user_id": 1}}
	productFact := &Fact{ID: "p1", Type: "Product", Attributes: map[string]interface{}{"id": 100}}
	
	joinNode1 := &JoinNode{
		BaseNode: BaseNode{ID: "join1", Children: []Node{}},
		LeftVariables:  []string{"user"},
		RightVariables: []string{"order"},
		AllVariables:   []string{"user", "order"},
		VariableTypes: map[string]string{
			"user": "User", "order": "Order", "product": "Product",
		},
		LeftMemory:  []*Token{},
		RightMemory: []*Fact{},
	}
	
	joinNode2 := &JoinNode{
		BaseNode: BaseNode{ID: "join2", Children: []Node{}},
		LeftVariables:  []string{"user", "order"},
		RightVariables: []string{"product"},
		AllVariables:   []string{"user", "order", "product"},
		VariableTypes: map[string]string{
			"user": "User", "order": "Order", "product": "Product",
		},
		LeftMemory:  []*Token{},
		RightMemory: []*Fact{},
	}
	
	joinNode1.AddChild(joinNode2)
	
	userToken := NewTokenWithFact(userFact, "user", "test")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		joinNode1.LeftMemory = []*Token{}
		joinNode1.RightMemory = []*Fact{}
		joinNode2.LeftMemory = []*Token{}
		joinNode2.RightMemory = []*Fact{}
		
		_ = joinNode1.ActivateLeft(userToken)
		_ = joinNode1.ActivateRight(orderFact)
		_ = joinNode2.ActivateRight(productFact)
	}
}
```

#### 3.3 Benchmark performJoinWithTokens

```go
func BenchmarkJoinNode_PerformJoinWithTokens(b *testing.B) {
	userFact := &Fact{ID: "u1", Type: "User"}
	orderFact := &Fact{ID: "o1", Type: "Order"}
	
	token1 := &Token{
		ID:       "t1",
		Facts:    []*Fact{userFact},
		Bindings: NewBindingChain().Add("user", userFact),
		NodeID:   "test",
	}
	
	token2 := &Token{
		ID:       "t2",
		Facts:    []*Fact{orderFact},
		Bindings: NewBindingChain().Add("order", orderFact),
		NodeID:   "test",
	}
	
	joinNode := &JoinNode{
		BaseNode:       BaseNode{ID: "join_bench"},
		AllVariables:   []string{"user", "order"},
		JoinConditions: nil,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = joinNode.performJoinWithTokens(token1, token2)
	}
}
```

---

### Tâche 4 : Benchmarks Comparatifs (40 min)

#### 4.1 Créer un benchmark avec l'ancien système (si possible)

**Si vous avez conservé une copie du code avant refactoring** :

```go
// BenchmarkOldSystem_2Variables mesure les performances de l'ancien système
func BenchmarkOldSystem_2Variables(b *testing.B) {
	// Implémenter avec l'ancienne structure Token (map[string]*Fact)
	// Pour comparaison uniquement
	
	b.Skip("Ancien code non disponible - benchmark de référence seulement")
}
```

#### 4.2 Benchmark de mémoire

```go
func BenchmarkBindingChain_Memory(b *testing.B) {
	b.ReportAllocs()
	
	fact := &Fact{ID: "f1", Type: "User"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		chain := NewBindingChain()
		for j := 0; j < 10; j++ {
			chain = chain.Add(fmt.Sprintf("var%d", j), fact)
		}
	}
}

func BenchmarkJoinNode_Memory_3Variables(b *testing.B) {
	b.ReportAllocs()
	
	// Setup identique à BenchmarkJoinNode_3Variables
	// ...
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Exécuter jointure 3 variables
	}
}
```

---

### Tâche 5 : Exécuter les Benchmarks (30 min)

#### 5.1 Exécuter tous les benchmarks

```bash
cd tsd

# Benchmarks BindingChain
go test -bench=BenchmarkBindingChain -benchmem ./rete/

# Benchmarks JoinNode
go test -bench=BenchmarkJoinNode -benchmem ./rete/

# Tous les benchmarks
go test -bench=. -benchmem ./rete/node_join_benchmark_test.go ./rete/binding_chain.go ./rete/node_join.go ./rete/fact_token.go
```

#### 5.2 Sauvegarder les résultats

```bash
go test -bench=. -benchmem ./rete/ > benchmark_results.txt
```

**Résultats attendus** (exemple) :

```
BenchmarkBindingChain_Add-8                     50000000    25.3 ns/op    16 B/op    1 allocs/op
BenchmarkBindingChain_Get_SmallChain-8          100000000   11.2 ns/op     0 B/op    0 allocs/op
BenchmarkBindingChain_Get_LargeChain-8          10000000    150 ns/op      0 B/op    0 allocs/op
BenchmarkBindingChain_Merge-8                   20000000    85.4 ns/op    80 B/op    5 allocs/op
BenchmarkJoinNode_2Variables-8                  1000000     1200 ns/op   450 B/op   12 allocs/op
BenchmarkJoinNode_3Variables-8                  500000      2500 ns/op   920 B/op   25 allocs/op
```

---

### Tâche 6 : Analyser les Résultats (40 min)

#### 6.1 Calculer les métriques clés

**Overhead de BindingChain vs map** :

```
Temps Get() BindingChain (n=3)  : ~11 ns
Temps Get() map (théorique)     : ~5 ns
Overhead                        : ~2x (acceptable pour n < 10)
```

**Overhead jointure 3 vs 2 variables** :

```
Temps 2 variables : 1200 ns
Temps 3 variables : 2500 ns
Ratio             : 2.08x
Overhead          : ~8% (< 10% ✅)
```

#### 6.2 Identifier les goulots d'étranglement

**Si Get() est trop lent (> 50ns pour n=3)** :
- Ajouter un cache map optionnel dans BindingChain
- Lazy initialization du cache

**Si Merge() est trop lent** :
- Optimiser l'algorithme
- Utiliser un builder pattern

**Si allocations excessives** :
- Réutiliser des structures
- Pool d'objets si nécessaire

---

### Tâche 7 : Optimisations (si nécessaire) (50 min)

#### 7.1 Optimisation 1 : Cache dans BindingChain

**Si Get() est trop lent pour n > 10** :

```go
type BindingChain struct {
	Variable string
	Fact     *Fact
	Parent   *BindingChain
	
	// Cache optionnel (lazy)
	cache     map[string]*Fact
	cacheInit sync.Once
}

func (bc *BindingChain) Get(variable string) *Fact {
	if bc == nil || variable == "" {
		return nil
	}
	
	// Si chaîne courte (< 10), recherche linéaire
	if bc.Len() < 10 {
		current := bc
		for current != nil {
			if current.Variable == variable {
				return current.Fact
			}
			current = current.Parent
		}
		return nil
	}
	
	// Sinon, utiliser le cache
	bc.cacheInit.Do(func() {
		bc.cache = bc.ToMap()
	})
	
	return bc.cache[variable]
}
```

**Attention** : Ajouter le cache casse l'immutabilité stricte. Documenter cette décision.

#### 7.2 Optimisation 2 : Éviter les allocations

**Réutiliser les slices dans Variables()** :

```go
// Ajouter un champ dans Token
type Token struct {
	// ... champs existants
	cachedVars []string // Cache des variables
}

func (t *Token) GetVariables() []string {
	if t.cachedVars == nil && t.Bindings != nil {
		t.cachedVars = t.Bindings.Variables()
	}
	return t.cachedVars
}
```

#### 7.3 Re-benchmarker après optimisations

```bash
go test -bench=. -benchmem ./rete/ > benchmark_results_optimized.txt
```

**Comparer** avec les résultats précédents.

---

### Tâche 8 : Documenter les Performances (30 min)

#### 8.1 Créer un document de performance

**Fichier** : `tsd/docs/architecture/BINDINGS_PERFORMANCE.md`

**Structure** :

```markdown
# Performances du Système de Bindings Immuable

**Date** : [DATE]
**Version** : Post-refactoring

## Résumé

Le nouveau système de bindings basé sur BindingChain maintient des performances acceptables avec un overhead < 10% pour les cas d'usage typiques (N ≤ 10 variables).

## Benchmarks

### BindingChain

| Opération | n=3 | n=10 | n=100 | Complexité |
|-----------|-----|------|-------|------------|
| Add()     | 25ns | 25ns | 25ns | O(1) |
| Get()     | 11ns | 35ns | 150ns | O(n) |
| Merge()   | 85ns | 280ns | 2500ns | O(m) |

### JoinNode

| Configuration | Temps | Allocations | vs Baseline |
|---------------|-------|-------------|-------------|
| 2 variables   | 1200ns | 450 B | Baseline |
| 3 variables   | 2500ns | 920 B | +108% (+8% overhead) |
| 4 variables   | 4200ns | 1450 B | +250% |

## Analyse

### Points Forts
- Add() est O(1) : Excellent pour la construction de chaînes
- Pas de régression pour 2 variables
- Overhead acceptable pour 3-10 variables

### Points d'Attention
- Get() est O(n) : Performance dégradée pour n > 100
- Allocations proportionnelles au nombre de variables

### Optimisations Appliquées
[Liste des optimisations si appliquées]

## Conclusion

✅ Performances acceptables pour les cas d'usage réels (N ≤ 10)
✅ Overhead < 10% pour jointures 3 variables
⚠️ Surveillance recommandée pour N > 10
```

---

## ✅ Critères de Validation

### Benchmarks
- [ ] Tous les benchmarks créés et exécutés
- [ ] Résultats sauvegardés dans `benchmark_results.txt`
- [ ] Comparaison avec baseline documentée

### Performance
- [ ] Overhead < 10% pour jointures 2→3 variables
- [ ] Get() acceptable pour n ≤ 10 (< 50ns)
- [ ] Pas de fuites mémoire
- [ ] Allocations raisonnables

### Optimisations
- [ ] Optimisations appliquées si nécessaire
- [ ] Re-benchmark après optimisations
- [ ] Gains documentés

### Documentation
- [ ] `BINDINGS_PERFORMANCE.md` créé
- [ ] Résultats documentés
- [ ] Recommandations claires

---

## 🎯 Prochaine Étape

Passer au **Prompt 12 - Documentation et Cleanup Final**.

Le dernier prompt finalisera toute la documentation, nettoiera le code, et préparera le commit final.

---

## 💡 Conseils Pratiques

### Pour les Benchmarks
1. **Exécuter plusieurs fois** : Les résultats peuvent varier
2. **Mesurer avec -benchmem** : Surveiller les allocations
3. **Isoler les benchmarks** : Un aspect à la fois
4. **Comparer avec baseline** : 2 variables = référence

### Pour l'Analyse
1. **Focus sur les cas réels** : N ≤ 10 est le plus courant
2. **Accepter un overhead raisonnable** : < 10% est excellent
3. **Ne pas optimiser prématurément** : Optimiser seulement si nécessaire
4. **Documenter les décisions** : Expliquer les trade-offs

### Pour les Optimisations
1. **Mesurer avant/après** : Prouver que l'optimisation aide
2. **Garder la simplicité** : Ne pas compliquer le code inutilement
3. **Préserver l'immutabilité** : C'est la garantie de correction
4. **Tester après optimisation** : S'assurer que rien ne casse

---

**Note** : Cette session valide que le refactoring est performant. L'objectif n'est pas d'avoir le système le plus rapide possible, mais de s'assurer qu'il n'y a pas de régression significative et que les performances sont acceptables pour les cas d'usage réels.