# Prompt 03 : Implémentation de BindingChain

**Session** : 3/12  
**Durée estimée** : 2-3 heures  
**Pré-requis** : Avoir complété Prompt 02 et lu `BINDINGS_DESIGN.md`

---

## 🎯 Objectif de cette Session

Implémenter la structure immuable `BindingChain` avec tous ses tests unitaires en :
1. Créant le fichier `rete/binding_chain.go` avec l'implémentation complète
2. Créant le fichier `rete/binding_chain_test.go` avec une couverture > 95%
3. Validant que la structure fonctionne correctement et respecte les invariants

**Livrables finaux** :
- `tsd/rete/binding_chain.go` (~300-400 lignes)
- `tsd/rete/binding_chain_test.go` (~500-700 lignes)

---

## 📋 Tâches à Réaliser

### Tâche 1 : Implémenter la Structure BindingChain (60 min)

#### 1.1 Créer le fichier binding_chain.go

**Fichier** : `tsd/rete/binding_chain.go`

**En-tête obligatoire** :
```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete
```

**Structure de base** :
```go
// BindingChain représente une chaîne immuable de bindings variable → fact.
// Utilise le pattern "Cons list" pour le partage de structure (structural sharing).
// Une fois créée, une BindingChain ne peut jamais être modifiée.
//
// Exemple d'utilisation :
//   chain := NewBindingChain()
//   chain = chain.Add("user", userFact)
//   chain = chain.Add("order", orderFact)
//   fact := chain.Get("user")  // Retourne userFact
type BindingChain struct {
	// Variable est le nom de la variable liée (ex: "u", "order", "task")
	Variable string

	// Fact est le fait associé à cette variable
	Fact *Fact

	// Parent pointe vers la chaîne parente (nil pour une chaîne vide)
	Parent *BindingChain
}
```

**Invariants à garantir** :
```go
// Invariants de BindingChain :
// 1. Une fois créée, une BindingChain ne change jamais (immutabilité)
// 2. Add() retourne une NOUVELLE chaîne sans modifier l'existante
// 3. Parent pointe toujours vers une chaîne plus courte (pas de cycles)
// 4. Variable vide ("") est invalide sauf pour la chaîne vide (Parent == nil)
```

---

#### 1.2 Implémenter les constructeurs

**Code à écrire** :

```go
// NewBindingChain crée une chaîne de bindings vide.
func NewBindingChain() *BindingChain {
	return nil  // nil représente une chaîne vide
}

// NewBindingChainWith crée une chaîne avec un binding initial.
func NewBindingChainWith(variable string, fact *Fact) *BindingChain {
	if variable == "" {
		return nil
	}
	return &BindingChain{
		Variable: variable,
		Fact:     fact,
		Parent:   nil,
	}
}
```

**Validation** :
- Une chaîne vide est représentée par `nil`
- Variable vide retourne une chaîne vide
- Fact peut être nil (binding vers nil est valide)

---

#### 1.3 Implémenter les opérations de lecture

**Code à écrire** :

```go
// Get retourne le fait associé à une variable, ou nil si non trouvé.
// Complexité : O(n) où n = nombre de bindings.
func (bc *BindingChain) Get(variable string) *Fact {
	if bc == nil || variable == "" {
		return nil
	}
	
	current := bc
	for current != nil {
		if current.Variable == variable {
			return current.Fact
		}
		current = current.Parent
	}
	
	return nil
}

// Has vérifie si une variable existe dans la chaîne.
// Complexité : O(n).
func (bc *BindingChain) Has(variable string) bool {
	if bc == nil || variable == "" {
		return false
	}
	
	current := bc
	for current != nil {
		if current.Variable == variable {
			return true
		}
		current = current.Parent
	}
	
	return false
}

// Len retourne le nombre de bindings dans la chaîne.
// Complexité : O(n).
func (bc *BindingChain) Len() int {
	if bc == nil {
		return 0
	}
	
	count := 0
	current := bc
	for current != nil {
		count++
		current = current.Parent
	}
	
	return count
}

// Variables retourne la liste des variables dans l'ordre d'ajout (du plus ancien au plus récent).
// Complexité : O(n).
func (bc *BindingChain) Variables() []string {
	if bc == nil {
		return []string{}
	}
	
	// Parcourir pour compter
	length := bc.Len()
	vars := make([]string, length)
	
	// Remplir à l'envers (car la chaîne est dans l'ordre inverse)
	current := bc
	for i := length - 1; i >= 0; i-- {
		vars[i] = current.Variable
		current = current.Parent
	}
	
	return vars
}

// ToMap convertit la chaîne en map pour compatibilité/debug.
// En cas de variable dupliquée, le binding le plus récent (tête de chaîne) est gardé.
// Complexité : O(n).
func (bc *BindingChain) ToMap() map[string]*Fact {
	result := make(map[string]*Fact)
	
	if bc == nil {
		return result
	}
	
	// Parcourir de la fin vers le début pour que les bindings récents écrasent les anciens
	vars := bc.Variables()
	for _, v := range vars {
		result[v] = bc.Get(v)
	}
	
	return result
}
```

---

#### 1.4 Implémenter les opérations de construction

**Code à écrire** :

```go
// Add ajoute un binding et retourne une NOUVELLE chaîne.
// L'ancienne chaîne reste inchangée (immutabilité).
// Si la variable existe déjà, le nouveau binding masque l'ancien.
// Complexité : O(1).
func (bc *BindingChain) Add(variable string, fact *Fact) *BindingChain {
	if variable == "" {
		return bc
	}
	
	return &BindingChain{
		Variable: variable,
		Fact:     fact,
		Parent:   bc,
	}
}

// Merge combine deux chaînes en retournant une NOUVELLE chaîne.
// Tous les bindings de 'other' sont ajoutés à la chaîne courante.
// En cas de conflit (même variable), le binding de 'other' est prioritaire.
// Complexité : O(m) où m = taille de 'other'.
func (bc *BindingChain) Merge(other *BindingChain) *BindingChain {
	if other == nil {
		return bc
	}
	
	// Récupérer toutes les variables de 'other' dans l'ordre
	otherVars := other.Variables()
	
	// Composer les chaînes
	result := bc
	for _, v := range otherVars {
		fact := other.Get(v)
		result = result.Add(v, fact)
	}
	
	return result
}
```

---

#### 1.5 Implémenter les opérations de debug

**Code à écrire** :

```go
import (
	"fmt"
	"strings"
)

// String retourne une représentation textuelle pour debug.
func (bc *BindingChain) String() string {
	if bc == nil {
		return "BindingChain{}"
	}
	
	vars := bc.Variables()
	if len(vars) == 0 {
		return "BindingChain{}"
	}
	
	parts := make([]string, len(vars))
	for i, v := range vars {
		fact := bc.Get(v)
		if fact != nil {
			parts[i] = fmt.Sprintf("%s→%s", v, fact.Type)
		} else {
			parts[i] = fmt.Sprintf("%s→nil", v)
		}
	}
	
	return fmt.Sprintf("BindingChain{%s}", strings.Join(parts, ", "))
}

// Chain retourne la liste des variables depuis la racine (pour traçage).
// Équivalent à Variables() mais avec un nom plus explicite pour le debugging.
func (bc *BindingChain) Chain() []string {
	return bc.Variables()
}
```

---

### Tâche 2 : Implémenter les Tests Unitaires (90 min)

#### 2.1 Créer le fichier de tests

**Fichier** : `tsd/rete/binding_chain_test.go`

**En-tête** :
```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"testing"
)
```

---

#### 2.2 Tests des constructeurs

**Code à écrire** :

```go
func TestBindingChain_NewBindingChain(t *testing.T) {
	t.Log("🧪 TEST BindingChain - Création chaîne vide")
	
	chain := NewBindingChain()
	
	if chain != nil {
		t.Errorf("❌ Chaîne vide devrait être nil, got %v", chain)
		return
	}
	
	t.Log("✅ Chaîne vide est nil")
}

func TestBindingChain_NewBindingChainWith(t *testing.T) {
	t.Log("🧪 TEST BindingChain - Création avec binding")
	
	fact := &Fact{ID: "f1", Type: "User"}
	chain := NewBindingChainWith("user", fact)
	
	if chain == nil {
		t.Fatal("❌ Chaîne ne devrait pas être nil")
	}
	
	if chain.Variable != "user" {
		t.Errorf("❌ Variable attendue 'user', got '%s'", chain.Variable)
	}
	
	if chain.Fact != fact {
		t.Errorf("❌ Fact incorrect")
	}
	
	if chain.Parent != nil {
		t.Errorf("❌ Parent devrait être nil")
	}
	
	t.Log("✅ Chaîne créée correctement")
}

func TestBindingChain_NewBindingChainWith_EmptyVariable(t *testing.T) {
	t.Log("🧪 TEST BindingChain - Variable vide")
	
	fact := &Fact{ID: "f1", Type: "User"}
	chain := NewBindingChainWith("", fact)
	
	if chain != nil {
		t.Errorf("❌ Variable vide devrait retourner nil")
		return
	}
	
	t.Log("✅ Variable vide retourne chaîne vide")
}
```

---

#### 2.3 Tests d'ajout (immutabilité)

**Code à écrire** :

```go
func TestBindingChain_Add_Single(t *testing.T) {
	t.Log("🧪 TEST BindingChain - Ajout simple")
	
	chain := NewBindingChain()
	fact := &Fact{ID: "f1", Type: "User"}
	
	newChain := chain.Add("user", fact)
	
	if newChain == nil {
		t.Fatal("❌ Nouvelle chaîne ne devrait pas être nil")
	}
	
	if newChain.Variable != "user" {
		t.Errorf("❌ Variable incorrecte")
	}
	
	if newChain.Fact != fact {
		t.Errorf("❌ Fact incorrect")
	}
	
	t.Log("✅ Ajout simple fonctionne")
}

func TestBindingChain_Add_Multiple(t *testing.T) {
	t.Log("🧪 TEST BindingChain - Ajouts multiples")
	
	user := &Fact{ID: "u1", Type: "User"}
	order := &Fact{ID: "o1", Type: "Order"}
	product := &Fact{ID: "p1", Type: "Product"}
	
	chain := NewBindingChain()
	chain = chain.Add("user", user)
	chain = chain.Add("order", order)
	chain = chain.Add("product", product)
	
	if chain.Len() != 3 {
		t.Errorf("❌ Longueur attendue 3, got %d", chain.Len())
	}
	
	if !chain.Has("user") || !chain.Has("order") || !chain.Has("product") {
		t.Errorf("❌ Variables manquantes")
	}
	
	t.Log("✅ Ajouts multiples fonctionnent")
}

func TestBindingChain_Add_Preserves_Parent(t *testing.T) {
	t.Log("🧪 TEST BindingChain - Immutabilité (parent préservé)")
	
	user := &Fact{ID: "u1", Type: "User"}
	order := &Fact{ID: "o1", Type: "Order"}
	
	chain1 := NewBindingChain().Add("user", user)
	chain2 := chain1.Add("order", order)
	
	// Vérifier que chain1 n'a pas été modifié
	if chain1.Len() != 1 {
		t.Errorf("❌ chain1 a été modifiée (immutabilité violée)")
	}
	
	if chain1.Has("order") {
		t.Errorf("❌ chain1 ne devrait pas avoir 'order'")
	}
	
	// Vérifier que chain2 a les deux
	if chain2.Len() != 2 {
		t.Errorf("❌ chain2 devrait avoir 2 bindings")
	}
	
	if !chain2.Has("user") || !chain2.Has("order") {
		t.Errorf("❌ chain2 devrait avoir 'user' et 'order'")
	}
	
	t.Log("✅ Immutabilité préservée")
}
```

---

#### 2.4 Tests de lecture

**Code à écrire** :

```go
func TestBindingChain_Get_Existing(t *testing.T) {
	t.Log("🧪 TEST BindingChain - Get (variable existante)")
	
	user := &Fact{ID: "u1", Type: "User"}
	chain := NewBindingChain().Add("user", user)
	
	result := chain.Get("user")
	
	if result != user {
		t.Errorf("❌ Fact attendu %v, got %v", user, result)
	}
	
	t.Log("✅ Get retourne le bon fact")
}

func TestBindingChain_Get_NotFound(t *testing.T) {
	t.Log("🧪 TEST BindingChain - Get (variable inexistante)")
	
	user := &Fact{ID: "u1", Type: "User"}
	chain := NewBindingChain().Add("user", user)
	
	result := chain.Get("order")
	
	if result != nil {
		t.Errorf("❌ Get devrait retourner nil pour variable inexistante")
	}
	
	t.Log("✅ Get retourne nil si variable inexistante")
}

func TestBindingChain_Has(t *testing.T) {
	t.Log("🧪 TEST BindingChain - Has")
	
	tests := []struct {
		name     string
		chain    *BindingChain
		variable string
		expected bool
	}{
		{"chaîne vide", NewBindingChain(), "user", false},
		{"variable existante", NewBindingChain().Add("user", &Fact{}), "user", true},
		{"variable inexistante", NewBindingChain().Add("user", &Fact{}), "order", false},
		{"variable vide", NewBindingChain().Add("user", &Fact{}), "", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.chain.Has(tt.variable)
			if result != tt.expected {
				t.Errorf("❌ %s: attendu %v, got %v", tt.name, tt.expected, result)
			}
		})
	}
	
	t.Log("✅ Has fonctionne correctement")
}

func TestBindingChain_Len(t *testing.T) {
	t.Log("🧪 TEST BindingChain - Len")
	
	tests := []struct {
		name     string
		build    func() *BindingChain
		expected int
	}{
		{"chaîne vide", func() *BindingChain { return NewBindingChain() }, 0},
		{"1 binding", func() *BindingChain { return NewBindingChain().Add("u", &Fact{}) }, 1},
		{"3 bindings", func() *BindingChain {
			return NewBindingChain().Add("u", &Fact{}).Add("o", &Fact{}).Add("p", &Fact{})
		}, 3},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chain := tt.build()
			result := chain.Len()
			if result != tt.expected {
				t.Errorf("❌ %s: longueur attendue %d, got %d", tt.name, tt.expected, result)
			}
		})
	}
	
	t.Log("✅ Len fonctionne correctement")
}

func TestBindingChain_Variables(t *testing.T) {
	t.Log("🧪 TEST BindingChain - Variables")
	
	chain := NewBindingChain()
	chain = chain.Add("user", &Fact{})
	chain = chain.Add("order", &Fact{})
	chain = chain.Add("product", &Fact{})
	
	vars := chain.Variables()
	
	expected := []string{"user", "order", "product"}
	if len(vars) != len(expected) {
		t.Fatalf("❌ Longueur attendue %d, got %d", len(expected), len(vars))
	}
	
	for i, v := range expected {
		if vars[i] != v {
			t.Errorf("❌ Index %d: attendu '%s', got '%s'", i, v, vars[i])
		}
	}
	
	t.Log("✅ Variables dans le bon ordre")
}

func TestBindingChain_ToMap(t *testing.T) {
	t.Log("🧪 TEST BindingChain - ToMap")
	
	user := &Fact{ID: "u1", Type: "User"}
	order := &Fact{ID: "o1", Type: "Order"}
	
	chain := NewBindingChain().Add("user", user).Add("order", order)
	m := chain.ToMap()
	
	if len(m) != 2 {
		t.Errorf("❌ Map devrait avoir 2 entrées, got %d", len(m))
	}
	
	if m["user"] != user {
		t.Errorf("❌ user fact incorrect")
	}
	
	if m["order"] != order {
		t.Errorf("❌ order fact incorrect")
	}
	
	t.Log("✅ ToMap fonctionne")
}
```

---

#### 2.5 Tests de merge

**Code à écrire** :

```go
func TestBindingChain_Merge(t *testing.T) {
	t.Log("🧪 TEST BindingChain - Merge")
	
	user := &Fact{ID: "u1", Type: "User"}
	order := &Fact{ID: "o1", Type: "Order"}
	product := &Fact{ID: "p1", Type: "Product"}
	
	chain1 := NewBindingChain().Add("user", user)
	chain2 := NewBindingChain().Add("order", order).Add("product", product)
	
	merged := chain1.Merge(chain2)
	
	if merged.Len() != 3 {
		t.Errorf("❌ Merged devrait avoir 3 bindings, got %d", merged.Len())
	}
	
	if !merged.Has("user") || !merged.Has("order") || !merged.Has("product") {
		t.Errorf("❌ Variables manquantes après merge")
	}
	
	t.Log("✅ Merge fonctionne")
}

func TestBindingChain_Merge_Conflicts(t *testing.T) {
	t.Log("🧪 TEST BindingChain - Merge avec conflit")
	
	user1 := &Fact{ID: "u1", Type: "User"}
	user2 := &Fact{ID: "u2", Type: "User"}
	
	chain1 := NewBindingChain().Add("user", user1)
	chain2 := NewBindingChain().Add("user", user2)
	
	merged := chain1.Merge(chain2)
	
	// En cas de conflit, le binding de chain2 (other) est prioritaire
	result := merged.Get("user")
	if result != user2 {
		t.Errorf("❌ En cas de conflit, 'other' devrait être prioritaire")
	}
	
	t.Log("✅ Merge avec conflit - 'other' prioritaire")
}
```

---

#### 2.6 Tests edge cases

**Code à écrire** :

```go
func TestBindingChain_Nil_Operations(t *testing.T) {
	t.Log("🧪 TEST BindingChain - Opérations sur nil")
	
	var chain *BindingChain // nil
	
	// Get
	if chain.Get("user") != nil {
		t.Errorf("❌ Get sur nil devrait retourner nil")
	}
	
	// Has
	if chain.Has("user") {
		t.Errorf("❌ Has sur nil devrait retourner false")
	}
	
	// Len
	if chain.Len() != 0 {
		t.Errorf("❌ Len sur nil devrait retourner 0")
	}
	
	// Variables
	vars := chain.Variables()
	if len(vars) != 0 {
		t.Errorf("❌ Variables sur nil devrait retourner slice vide")
	}
	
	// ToMap
	m := chain.ToMap()
	if len(m) != 0 {
		t.Errorf("❌ ToMap sur nil devrait retourner map vide")
	}
	
	// Add
	newChain := chain.Add("user", &Fact{})
	if newChain == nil {
		t.Errorf("❌ Add sur nil devrait créer une chaîne")
	}
	
	t.Log("✅ Opérations sur nil gérées correctement")
}

func TestBindingChain_Long_Chain(t *testing.T) {
	t.Log("🧪 TEST BindingChain - Chaîne longue (100 bindings)")
	
	chain := NewBindingChain()
	
	// Ajouter 100 bindings
	for i := 0; i < 100; i++ {
		varName := fmt.Sprintf("var%d", i)
		fact := &Fact{ID: varName, Type: "Type"}
		chain = chain.Add(varName, fact)
	}
	
	if chain.Len() != 100 {
		t.Errorf("❌ Longueur attendue 100, got %d", chain.Len())
	}
	
	// Vérifier qu'on peut retrouver tous les bindings
	for i := 0; i < 100; i++ {
		varName := fmt.Sprintf("var%d", i)
		if !chain.Has(varName) {
			t.Errorf("❌ Variable '%s' manquante", varName)
		}
	}
	
	t.Log("✅ Chaîne longue fonctionne")
}
```

---

#### 2.7 Tests de performance (benchmarks)

**Code à écrire** :

```go
func BenchmarkBindingChain_Add(b *testing.B) {
	fact := &Fact{ID: "f1", Type: "User"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		chain := NewBindingChain()
		chain.Add("user", fact)
	}
}

func BenchmarkBindingChain_Get(b *testing.B) {
	chain := NewBindingChain()
	for i := 0; i < 10; i++ {
		chain = chain.Add(fmt.Sprintf("var%d", i), &Fact{})
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		chain.Get("var5")
	}
}

func BenchmarkBindingChain_Get_DeepChain(b *testing.B) {
	chain := NewBindingChain()
	for i := 0; i < 100; i++ {
		chain = chain.Add(fmt.Sprintf("var%d", i), &Fact{})
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		chain.Get("var0") // Chercher le plus ancien (pire cas)
	}
}

func BenchmarkBindingChain_Merge(b *testing.B) {
	chain1 := NewBindingChain()
	chain2 := NewBindingChain()
	
	for i := 0; i < 5; i++ {
		chain1 = chain1.Add(fmt.Sprintf("v1_%d", i), &Fact{})
		chain2 = chain2.Add(fmt.Sprintf("v2_%d", i), &Fact{})
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		chain1.Merge(chain2)
	}
}
```

---

### Tâche 3 : Validation et Tests (20 min)

#### 3.1 Exécuter les tests

**Commandes** :
```bash
cd tsd

# Compiler
go build ./rete/...

# Tests unitaires
go test -v ./rete/binding_chain_test.go ./rete/binding_chain.go ./rete/fact_token.go

# Tests avec couverture
go test -coverprofile=coverage.out ./rete/binding_chain_test.go ./rete/binding_chain.go ./rete/fact_token.go
go tool cover -html=coverage.out -o coverage.html

# Benchmarks
go test -bench=. ./rete/binding_chain_test.go ./rete/binding_chain.go ./rete/fact_token.go
```

**Résultats attendus** :
- ✅ Tous les tests passent
- ✅ Couverture > 95%
- ✅ Benchmarks montrent des performances acceptables

---

#### 3.2 Vérifier la couverture

**Objectif** : > 95% de couverture

**Si couverture insuffisante** :
- Identifier les lignes non couvertes
- Ajouter des tests pour ces cas
- Vérifier les edge cases

---

#### 3.3 Vérifier le respect des standards

**Checklist** :
- [ ] En-tête de copyright présent
- [ ] GoDoc pour toutes les fonctions exportées
- [ ] Code formaté (`go fmt`)
- [ ] Pas de warnings (`go vet`)
- [ ] Conventions de nommage respectées
- [ ] Messages de test clairs avec émojis

**Commandes** :
```bash
go fmt ./rete/binding_chain.go ./rete/binding_chain_test.go
go vet ./rete/binding_chain.go ./rete/binding_chain_test.go
```

---

## ✅ Critères de Validation de cette Session

À la fin de ce prompt, vous devez avoir :

### Livrables
- [ ] ✅ Fichier `rete/binding_chain.go` complet et commenté
- [ ] ✅ Fichier `rete/binding_chain_test.go` avec tous les tests
- [ ] ✅ Tous les tests passent
- [ ] ✅ Couverture > 95%
- [ ] ✅ Benchmarks exécutés

### Qualité du Code
- [ ] Immutabilité respectée (tests le prouvent)
- [ ] Toutes les méthodes de l'API implémentées
- [ ] Edge cases gérés (nil, chaîne vide, etc.)
- [ ] Performance acceptable (benchmarks)
- [ ] Code clair et bien documenté

### Standards
- [ ] En-tête de copyright présent
- [ ] GoDoc complet
- [ ] Code formaté
- [ ] Pas de warnings
- [ ] Conventions respectées

---

## 🎯 Prochaine Étape

Une fois BindingChain **implémentée et validée**, passer au **Prompt 04 - Refactoring de Token**.

Le Prompt 04 remplacera complètement l'ancienne structure Token pour utiliser BindingChain.

---

## 💡 Conseils Pratiques

### Pour l'Implémentation
1. Commencer par les constructeurs
2. Implémenter les méthodes simples (Get, Has)
3. Tester au fur et à mesure
4. Ajouter les méthodes complexes (Merge)
5. Optimiser si nécessaire

### Pour les Tests
1. Commencer par les cas nominaux
2. Ajouter les edge cases
3. Tester l'immutabilité explicitement
4. Utiliser des noms de tests descriptifs
5. Ajouter des logs clairs

### Pour la Performance
1. Ne pas optimiser prématurément
2. Mesurer avec des benchmarks
3. Si Get() est trop lent pour n > 10, envisager un cache
4. Documenter les décisions de performance

---

**Note** : Cette session implémente une structure de données **complètement nouvelle**. Aucun code existant n'est modifié. Le but est de créer une fondation solide pour les prochains prompts.