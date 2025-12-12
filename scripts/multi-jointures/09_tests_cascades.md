# Prompt 09 : Tests Cascades Multi-Variables

**Session** : 9/12  
**Durée estimée** : 2-3 heures  
**Pré-requis** : Prompt 08 complété, ExecutionContext adapté

---

## 🎯 Objectif de cette Session

Créer des tests unitaires complets pour valider les cascades de jointures à N variables :
1. Tests de régression pour 2 variables
2. Tests exhaustifs pour 3 variables
3. Tests paramétriques pour N variables (N=2 à 10)
4. Validation que tous les bindings sont préservés

**Livrable** : `tsd/rete/node_join_cascade_test.go` (nouveau fichier, ~500-700 lignes)

---

## 📋 Tâches à Réaliser

### Tâche 1 : Créer le Fichier de Tests (20 min)

**Fichier** : `tsd/rete/node_join_cascade_test.go`

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

### Tâche 2 : Tests de Régression 2 Variables (40 min)

```go
func TestJoinCascade_2Variables_UserOrder(t *testing.T) {
	t.Log("🧪 TEST Cascade 2 variables - User-Order (régression)")
	
	// Setup : Faits
	userFact := &Fact{
		ID:   "u1",
		Type: "User",
		Attributes: map[string]interface{}{
			"id":   1,
			"name": "Alice",
		},
	}
	
	orderFact := &Fact{
		ID:   "o1",
		Type: "Order",
		Attributes: map[string]interface{}{
			"id":      100,
			"user_id": 1,
		},
	}
	
	// Setup : JoinNode
	joinNode := &JoinNode{
		BaseNode: BaseNode{
			ID:       "join_user_order",
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
		JoinConditions: nil, // Pas de conditions pour ce test
	}
	
	// Mock pour capturer le token final
	var finalToken *Token
	mockTerminal := &MockNode{
		OnActivateLeft: func(token *Token) error {
			finalToken = token
			return nil
		},
	}
	joinNode.AddChild(mockTerminal)
	
	// Act : Soumettre les faits
	// Scénario : User arrive en premier (ActivateLeft), puis Order (ActivateRight)
	userToken := NewTokenWithFact(userFact, "user", "type_node_user")
	err := joinNode.ActivateLeft(userToken)
	if err != nil {
		t.Fatalf("❌ ActivateLeft erreur: %v", err)
	}
	
	err = joinNode.ActivateRight(orderFact)
	if err != nil {
		t.Fatalf("❌ ActivateRight erreur: %v", err)
	}
	
	// Assert
	if finalToken == nil {
		t.Fatal("❌ Aucun token propagé au terminal")
	}
	
	if finalToken.Bindings.Len() != 2 {
		t.Errorf("❌ Attendu 2 bindings, got %d", finalToken.Bindings.Len())
	}
	
	if !finalToken.HasBinding("user") {
		t.Errorf("❌ Variable 'user' manquante")
	}
	
	if !finalToken.HasBinding("order") {
		t.Errorf("❌ Variable 'order' manquante")
	}
	
	if finalToken.GetBinding("user") != userFact {
		t.Errorf("❌ Binding 'user' incorrect")
	}
	
	if finalToken.GetBinding("order") != orderFact {
		t.Errorf("❌ Binding 'order' incorrect")
	}
	
	t.Log("✅ Cascade 2 variables fonctionne (régression OK)")
}
```

---

### Tâche 3 : Tests 3 Variables - Cas Principal (60 min)

```go
func TestJoinCascade_3Variables_UserOrderProduct(t *testing.T) {
	t.Log("🧪 TEST Cascade 3 variables - User-Order-Product")
	
	// Setup : Faits
	userFact := &Fact{
		ID:   "u1",
		Type: "User",
		Attributes: map[string]interface{}{"id": 1, "name": "Alice"},
	}
	
	orderFact := &Fact{
		ID:   "o1",
		Type: "Order",
		Attributes: map[string]interface{}{"id": 100, "user_id": 1, "product_id": 200},
	}
	
	productFact := &Fact{
		ID:   "p1",
		Type: "Product",
		Attributes: map[string]interface{}{"id": 200, "name": "Laptop"},
	}
	
	// Setup : Cascade de 2 JoinNodes
	// JoinNode1 : User + Order → [user, order]
	joinNode1 := &JoinNode{
		BaseNode: BaseNode{
			ID:       "join1_user_order",
			Children: []Node{},
		},
		LeftVariables:  []string{"user"},
		RightVariables: []string{"order"},
		AllVariables:   []string{"user", "order"},
		VariableTypes: map[string]string{
			"user":    "User",
			"order":   "Order",
			"product": "Product",
		},
		LeftMemory:     []*Token{},
		RightMemory:    []*Fact{},
		JoinConditions: nil,
	}
	
	// JoinNode2 : [user, order] + Product → [user, order, product]
	joinNode2 := &JoinNode{
		BaseNode: BaseNode{
			ID:       "join2_add_product",
			Children: []Node{},
		},
		LeftVariables:  []string{"user", "order"},
		RightVariables: []string{"product"},
		AllVariables:   []string{"user", "order", "product"},
		VariableTypes: map[string]string{
			"user":    "User",
			"order":   "Order",
			"product": "Product",
		},
		LeftMemory:     []*Token{},
		RightMemory:    []*Fact{},
		JoinConditions: nil,
	}
	
	// Connecter : JoinNode1 → JoinNode2 → Terminal
	joinNode1.AddChild(joinNode2)
	
	var finalToken *Token
	mockTerminal := &MockNode{
		OnActivateLeft: func(token *Token) error {
			finalToken = token
			return nil
		},
	}
	joinNode2.AddChild(mockTerminal)
	
	// Act : Soumettre les faits dans l'ordre User → Order → Product
	userToken := NewTokenWithFact(userFact, "user", "type_user")
	err := joinNode1.ActivateLeft(userToken)
	if err != nil {
		t.Fatalf("❌ JoinNode1 ActivateLeft erreur: %v", err)
	}
	
	err = joinNode1.ActivateRight(orderFact)
	if err != nil {
		t.Fatalf("❌ JoinNode1 ActivateRight erreur: %v", err)
	}
	
	err = joinNode2.ActivateRight(productFact)
	if err != nil {
		t.Fatalf("❌ JoinNode2 ActivateRight erreur: %v", err)
	}
	
	// Assert
	if finalToken == nil {
		t.Fatal("❌ Aucun token propagé au terminal")
	}
	
	// CRITIQUE : Le token final DOIT contenir 3 bindings
	if finalToken.Bindings.Len() != 3 {
		t.Errorf("❌ CRITIQUE: Attendu 3 bindings, got %d", finalToken.Bindings.Len())
		t.Errorf("   Variables présentes: %v", finalToken.GetVariables())
	}
	
	// Vérifier chaque variable
	expectedVars := []string{"user", "order", "product"}
	for _, v := range expectedVars {
		if !finalToken.HasBinding(v) {
			t.Errorf("❌ Variable '%s' manquante dans le token final", v)
		}
	}
	
	// Vérifier les valeurs
	if finalToken.GetBinding("user") != userFact {
		t.Errorf("❌ Binding 'user' incorrect")
	}
	if finalToken.GetBinding("order") != orderFact {
		t.Errorf("❌ Binding 'order' incorrect")
	}
	if finalToken.GetBinding("product") != productFact {
		t.Errorf("❌ Binding 'product' incorrect")
	}
	
	t.Log("✅ Cascade 3 variables fonctionne - TOUS les bindings préservés")
}
```

---

### Tâche 4 : Tests avec Ordres Différents (40 min)

```go
func TestJoinCascade_3Variables_DifferentOrders(t *testing.T) {
	t.Log("🧪 TEST Cascade 3 variables - Différents ordres de soumission")
	
	// Définir les 6 permutations possibles
	orders := []struct {
		name  string
		order []string // Ordre de soumission : "user", "order", "product"
	}{
		{"User→Order→Product", []string{"user", "order", "product"}},
		{"User→Product→Order", []string{"user", "product", "order"}},
		{"Order→User→Product", []string{"order", "user", "product"}},
		{"Order→Product→User", []string{"order", "product", "user"}},
		{"Product→User→Order", []string{"product", "user", "order"}},
		{"Product→Order→User", []string{"product", "order", "user"}},
	}
	
	for _, tc := range orders {
		t.Run(tc.name, func(t *testing.T) {
			// Setup identique à TestJoinCascade_3Variables_UserOrderProduct
			// mais soumettre les faits dans l'ordre spécifié par tc.order
			
			// Créer les faits
			facts := map[string]*Fact{
				"user":    {ID: "u1", Type: "User", Attributes: map[string]interface{}{"id": 1}},
				"order":   {ID: "o1", Type: "Order", Attributes: map[string]interface{}{"id": 100}},
				"product": {ID: "p1", Type: "Product", Attributes: map[string]interface{}{"id": 200}},
			}
			
			// Créer la cascade (même setup que test précédent)
			joinNode1, joinNode2, mockTerminal := setupCascade3Variables()
			
			var finalToken *Token
			mockTerminal.OnActivateLeft = func(token *Token) error {
				finalToken = token
				return nil
			}
			
			// Soumettre dans l'ordre spécifié
			for _, factName := range tc.order {
				fact := facts[factName]
				// Déterminer quel nœud activer
				// (Logique simplifiée pour le test)
				submitFactToCascade(joinNode1, joinNode2, factName, fact)
			}
			
			// Assert : Le résultat doit être le même quel que soit l'ordre
			if finalToken == nil {
				t.Errorf("❌ Aucun token final pour ordre %v", tc.order)
				return
			}
			
			if finalToken.Bindings.Len() != 3 {
				t.Errorf("❌ Ordre %v: attendu 3 bindings, got %d", tc.order, finalToken.Bindings.Len())
			}
		})
	}
	
	t.Log("✅ Résultats cohérents quel que soit l'ordre de soumission")
}

// Helpers
func setupCascade3Variables() (*JoinNode, *JoinNode, *MockNode) {
	// Retourner joinNode1, joinNode2, mockTerminal configurés
	// (Code factoriser du test précédent)
	return nil, nil, nil // TODO: Implémenter
}

func submitFactToCascade(jn1, jn2 *JoinNode, factName string, fact *Fact) {
	// Logique pour soumettre un fait au bon endroit de la cascade
	// TODO: Implémenter selon votre architecture
}
```

---

### Tâche 5 : Tests Paramétriques N Variables (50 min)

```go
func TestJoinCascade_NVariables(t *testing.T) {
	t.Log("🧪 TEST Cascade N variables - Scalabilité")
	
	for n := 2; n <= 10; n++ {
		t.Run(fmt.Sprintf("n=%d_variables", n), func(t *testing.T) {
			// Générer N faits
			facts := make([]*Fact, n)
			varNames := make([]string, n)
			
			for i := 0; i < n; i++ {
				varNames[i] = fmt.Sprintf("var%d", i)
				facts[i] = &Fact{
					ID:   fmt.Sprintf("f%d", i),
					Type: fmt.Sprintf("Type%d", i),
					Attributes: map[string]interface{}{"id": i},
				}
			}
			
			// Construire une cascade de (n-1) JoinNodes
			joinNodes := buildCascade(n, varNames)
			
			// Mock terminal
			var finalToken *Token
			mockTerminal := &MockNode{
				OnActivateLeft: func(token *Token) error {
					finalToken = token
					return nil
				},
			}
			lastJoinNode := joinNodes[len(joinNodes)-1]
			lastJoinNode.AddChild(mockTerminal)
			
			// Soumettre les faits séquentiellement
			for i, fact := range facts {
				if i == 0 {
					// Premier fait : ActivateLeft du premier JoinNode
					token := NewTokenWithFact(fact, varNames[i], "type_node")
					joinNodes[0].ActivateLeft(token)
				} else if i == 1 {
					// Deuxième fait : ActivateRight du premier JoinNode
					joinNodes[0].ActivateRight(fact)
				} else {
					// Faits suivants : ActivateRight des JoinNodes suivants
					joinNodes[i-1].ActivateRight(fact)
				}
			}
			
			// Assert
			if finalToken == nil {
				t.Fatalf("❌ N=%d: Aucun token final", n)
			}
			
			if finalToken.Bindings.Len() != n {
				t.Errorf("❌ N=%d: Attendu %d bindings, got %d",
					n, n, finalToken.Bindings.Len())
				t.Errorf("   Variables présentes: %v", finalToken.GetVariables())
			}
			
			// Vérifier que chaque variable est présente
			for _, varName := range varNames {
				if !finalToken.HasBinding(varName) {
					t.Errorf("❌ N=%d: Variable '%s' manquante", n, varName)
				}
			}
			
			t.Logf("✅ N=%d variables: Tous les bindings préservés", n)
		})
	}
	
	t.Log("✅ Scalabilité validée jusqu'à N=10 variables")
}

// buildCascade construit une cascade de (n-1) JoinNodes pour n variables
func buildCascade(n int, varNames []string) []*JoinNode {
	if n < 2 {
		return []*JoinNode{}
	}
	
	joinNodes := make([]*JoinNode, n-1)
	varTypes := make(map[string]string)
	for i, v := range varNames {
		varTypes[v] = fmt.Sprintf("Type%d", i)
	}
	
	// JoinNode 1 : var0 + var1
	joinNodes[0] = &JoinNode{
		BaseNode: BaseNode{
			ID:       fmt.Sprintf("join_%d", 0),
			Children: []Node{},
		},
		LeftVariables:  []string{varNames[0]},
		RightVariables: []string{varNames[1]},
		AllVariables:   []string{varNames[0], varNames[1]},
		VariableTypes:  varTypes,
		LeftMemory:     []*Token{},
		RightMemory:    []*Fact{},
	}
	
	// JoinNodes suivants
	for i := 2; i < n; i++ {
		// LeftVars = toutes les variables précédentes
		leftVars := make([]string, i)
		copy(leftVars, varNames[0:i])
		
		// AllVars = toutes les variables jusqu'à i
		allVars := make([]string, i+1)
		copy(allVars, varNames[0:i+1])
		
		joinNodes[i-1] = &JoinNode{
			BaseNode: BaseNode{
				ID:       fmt.Sprintf("join_%d", i-1),
				Children: []Node{},
			},
			LeftVariables:  leftVars,
			RightVariables: []string{varNames[i]},
			AllVariables:   allVars,
			VariableTypes:  varTypes,
			LeftMemory:     []*Token{},
			RightMemory:    []*Fact{},
		}
		
		// Connecter au précédent
		joinNodes[i-2].AddChild(joinNodes[i-1])
	}
	
	return joinNodes
}
```

---

### Tâche 6 : Exécuter et Valider (30 min)

#### 6.1 Exécuter tous les tests

```bash
cd tsd

# Tests de cascade
go test -v ./rete/node_join_cascade_test.go ./rete/node_join.go ./rete/fact_token.go ./rete/binding_chain.go ./rete/node_base.go

# Tous les tests rete
go test -v ./rete/...
```

#### 6.2 Vérifier la couverture

```bash
go test -coverprofile=coverage.out ./rete/node_join_cascade_test.go
go tool cover -html=coverage.out
```

**Objectif** : Couverture > 90% pour le code de jointure.

---

## ✅ Critères de Validation

### Tests
- [ ] TestJoinCascade_2Variables_UserOrder passe (régression)
- [ ] TestJoinCascade_3Variables_UserOrderProduct passe (cas principal)
- [ ] TestJoinCascade_3Variables_DifferentOrders passe (robustesse)
- [ ] TestJoinCascade_NVariables passe pour N=2 à 10 (scalabilité)

### Couverture
- [ ] Couverture > 90% pour node_join.go
- [ ] Tous les cas de cascade testés

### Validation
- [ ] Aucun binding perdu dans les cascades
- [ ] Résultats cohérents quel que soit l'ordre
- [ ] Scalabilité jusqu'à N=10 validée

---

## 🎯 Prochaine Étape

Passer au **Prompt 10 - Validation E2E**.

Les tests E2E du Prompt 10 valideront que les 3 fixtures échouant passent maintenant et que tous les 83 tests E2E sont au vert.